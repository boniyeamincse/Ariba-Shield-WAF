package handlers
import "encoding/json"

import (
	"encoding/csv"
	"net/http"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

// auditEvent is the response shape for a single audit event.
type auditEvent struct {
	ID        string `json:"id"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
	ResourceID string `json:"resource_id"`
	Actor     string `json:"actor_user_id"`
	IP        string `json:"ip"`
	RequestID string `json:"request_id"`
	CreatedAt string `json:"created_at"`
}

// ListAuditEvents returns immutable audit events (append-only), newest first.
func ListAuditEvents(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, action, resource_type, resource_id, COALESCE(actor_user_id,''),
			        COALESCE(ip,''), COALESCE(request_id,''), created_at
			 FROM audit_events ORDER BY created_at DESC LIMIT 100`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var events []auditEvent
		for rows.Next() {
			var e auditEvent
			var ts time.Time
			if err := rows.Scan(&e.ID, &e.Action, &e.Resource, &e.ResourceID,
				&e.Actor, &e.IP, &e.RequestID, &ts); err != nil {
				continue
			}
			e.CreatedAt = ts.Format(time.RFC3339)
			events = append(events, e)
		}
		if events == nil {
			events = []auditEvent{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	}
}

// GetAuditEvent returns a single audit event by ID.
func GetAuditEvent(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var e auditEvent
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, action, resource_type, resource_id, COALESCE(actor_user_id,''),
			        COALESCE(ip,''), COALESCE(request_id,''), created_at
			 FROM audit_events WHERE id = $1`, id).Scan(&e.ID, &e.Action, &e.Resource, &e.ResourceID,
			&e.Actor, &e.IP, &e.RequestID, &e.CreatedAt)
		if err != nil {
			http.Error(w, `{"error":"audit event not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]auditEvent{"event": e})
	}
}

// GetAuditEventExport returns audit events as a CSV download.
func GetAuditEventExport(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, action, resource_type, resource_id, COALESCE(actor_user_id,''),
			        COALESCE(ip,''), COALESCE(request_id,''), created_at
			 FROM audit_events ORDER BY created_at DESC LIMIT 100`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var events []auditEvent
		for rows.Next() {
			var e auditEvent
			var ts time.Time
			if err := rows.Scan(&e.ID, &e.Action, &e.Resource, &e.ResourceID,
				&e.Actor, &e.IP, &e.RequestID, &ts); err != nil {
				continue
			}
			e.CreatedAt = ts.Format(time.RFC3339)
			events = append(events, e)
		}
		if events == nil {
			events = []auditEvent{}
		}

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", `attachment; filename="audit-events-` + time.Now().UTC().Format("20060102") + `.csv`)
	 writer := csv.NewWriter(w)
	 for _, e := range events {
	 	writer.Write([]string{e.ID, e.Action, e.Resource, e.ResourceID, e.Actor, e.IP, e.RequestID, e.CreatedAt})
	 }
	 writer.Flush()
	}
}
