package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListIncidents returns all incidents with pagination.
// Query params: limit (default 50, max 200), offset (default 0).
func ListIncidents(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 50
		if v := r.URL.Query().Get("limit"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
				limit = n
			}
		}
		offset := 0
		if v := r.URL.Query().Get("offset"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n >= 0 {
				offset = n
			}
		}

		rows, err := st.Pool.Query(r.Context(),
			`SELECT row_to_json(t) FROM (
			   SELECT id, title, severity, status, owner_user_id, notes, related_events, created_at, updated_at
			   FROM incidents ORDER BY created_at DESC LIMIT $1 OFFSET $2) t`, limit, offset)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var incidents []json.RawMessage
		for rows.Next() {
			var b []byte
			if err := rows.Scan(&b); err == nil {
				incidents = append(incidents, b)
			}
		}
		if incidents == nil {
			incidents = []json.RawMessage{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(incidents)
	}
}

// CreateIncident creates a new incident.
func CreateIncident(st *store.Store) http.HandlerFunc {
	type incident struct {
		Title         string   `json:"title"`
		Severity      string   `json:"severity"`
		Status        string   `json:"status"`
		OwnerUserID   string   `json:"owner_user_id"`
		Notes         string   `json:"notes"`
		RelatedEvents []string `json:"related_events"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body incident
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Title == "" {
			http.Error(w, `{"error":"title is required"}`, http.StatusBadRequest)
			return
		}
		if body.Severity == "" {
			body.Severity = "medium"
		}
		if body.Status == "" {
			body.Status = "open"
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		_, err = st.Pool.Exec(r.Context(),
			`INSERT INTO incidents (id, organization_id, title, severity, status, owner_user_id, notes, related_events, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now(), now())`,
			id, "01ARZ3NDEKTSV4RRFFQ69G5FAV", body.Title, body.Severity, body.Status, body.OwnerUserID, body.Notes, body.RelatedEvents)
		if err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}
func GetIncident(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var b []byte
		err := st.Pool.QueryRow(r.Context(),
			`SELECT row_to_json(t) FROM (
			   SELECT id, title, severity, status, owner_user_id, notes, related_events, created_at, updated_at
			   FROM incidents WHERE id = $1) t`, id).Scan(&b)
		if err != nil {
			http.Error(w, `{"error":"incident not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

// UpdateIncident updates an incident (partial).
func UpdateIncident(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		cols := []string{"title", "severity", "status", "owner_user_id", "notes", "related_events"}
		sets := []string{}
		vals := []any{}
		for _, c := range cols {
			if v, ok := body[c]; ok {
				vals = append(vals, v)
				sets = append(sets, c+" = $"+itoa(len(vals)))
			}
		}
		if len(sets) == 0 {
			http.Error(w, `{"error":"no updatable fields"}`, http.StatusBadRequest)
			return
		}
		vals = append(vals, id)
		ct, err := st.Pool.Exec(r.Context(),
			"UPDATE incidents SET "+joinComma(sets)+", updated_at = now() WHERE id = $"+itoa(len(vals)), vals...)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"incident not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// DeleteIncident removes an incident.
func DeleteIncident(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM incidents WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"incident not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}

// AssignIncident assigns an owner to an incident.
func AssignIncident(st *store.Store) http.HandlerFunc {
	type assign struct {
		OwnerUserID string `json:"owner_user_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body assign
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.OwnerUserID == "" {
			http.Error(w, `{"error":"owner_user_id required"}`, http.StatusBadRequest)
			return
		}
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE incidents SET owner_user_id = $1, status = 'investigating', updated_at = now() WHERE id = $2`,
			body.OwnerUserID, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"incident not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "investigating", "owner_user_id": body.OwnerUserID})
	}
}

// EscalateIncident escalates severity.
func EscalateIncident(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE incidents SET severity = 'critical', status = 'investigating', updated_at = now() WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"incident not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "severity": "critical", "status": "investigating"})
	}
}

// CloseIncident marks an incident resolved.
func CloseIncident(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE incidents SET status = 'resolved', updated_at = now() WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"incident not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "resolved"})
	}
}

// ReopenIncident reopens a resolved incident.
func ReopenIncident(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE incidents SET status = 'investigating', updated_at = now() WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"incident not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "investigating"})
	}
}

// IncidentEvents returns the security events related to an incident.
func IncidentEvents(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var relatedEvents []string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT related_events FROM incidents WHERE id = $1`, id).Scan(&relatedEvents); err != nil {
			http.Error(w, `{"error":"incident not found"}`, http.StatusNotFound)
			return
		}
		events := []map[string]any{}
		if len(relatedEvents) > 0 {
			rows, err := st.Pool.Query(r.Context(),
				`SELECT id, event_id, severity, reason, client_ip, created_at::text
				 FROM security_events WHERE id = ANY($1) ORDER BY created_at DESC`, relatedEvents)
			if err == nil {
				for rows.Next() {
					var eid, evid, sev, reason, ip, ts string
					if err := rows.Scan(&eid, &evid, &sev, &reason, &ip, &ts); err == nil {
						events = append(events, map[string]any{
							"id": eid, "event_id": evid, "severity": sev, "reason": reason, "client_ip": ip, "created_at": ts,
						})
					}
				}
				rows.Close()
			}
		}
		if events == nil {
			events = []map[string]any{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"incident_id": id, "events": events})
	}
}

// IncidentTimeline returns the status history of an incident.
func IncidentTimeline(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var title, status, severity, owner, created, updated string
		var notes []string
		err := st.Pool.QueryRow(r.Context(),
			`SELECT title, status, severity, COALESCE(owner_user_id,''), COALESCE(notes,''), created_at::text, updated_at::text
			 FROM incidents WHERE id = $1`, id).
			Scan(&title, &status, &severity, &owner, &notes, &created, &updated)
		if err != nil {
			http.Error(w, `{"error":"incident not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"incident_id": id, "title": title, "current_status": status, "severity": severity,
			"owner": owner, "notes": notes, "created_at": created, "updated_at": updated,
		})
	}
}

var _ = time.Now
