package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListAuditEvents returns immutable audit events (append-only), newest first.
func ListAuditEvents(st *store.Store) http.HandlerFunc {
	type auditEvent struct {
		ID             string `json:"id"`
		Action         string `json:"action"`
		ResourceType   string `json:"resource_type"`
		ResourceID     string `json:"resource_id"`
		ActorUserID    string `json:"actor_user_id"`
		IP             string `json:"ip"`
		RequestID      string `json:"request_id"`
		CreatedAt      string `json:"created_at"`
	}

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
			if err := rows.Scan(&e.ID, &e.Action, &e.ResourceType, &e.ResourceID,
				&e.ActorUserID, &e.IP, &e.RequestID, &ts); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
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