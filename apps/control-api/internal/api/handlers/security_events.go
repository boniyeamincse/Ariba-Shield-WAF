package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListSecurityEvents returns security events with pagination.
func ListSecurityEvents(st *store.Store) http.HandlerFunc {
	type event struct {
		ID            string   `json:"id"`
		EventID       string   `json:"event_id"`
		RequestID     string   `json:"request_id"`
		Timestamp     string   `json:"timestamp"`
		Severity      string   `json:"severity"`
		DecisionAction string  `json:"decision_action"`
		Reason        string   `json:"reason"`
		RuleIDs       []string `json:"rule_ids"`
		ClientIP      string   `json:"client_ip"`
		Method        string   `json:"method"`
		Path          string   `json:"path"`
		Status        int      `json:"status"`
		CreatedAt     string   `json:"created_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		limit := 50
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, event_id, request_id, COALESCE(severity,''), COALESCE(decision_action,''),
			        COALESCE(reason,''), rule_ids, COALESCE(client_ip,''), COALESCE(method,''),
			        COALESCE(path,''), COALESCE(status,0), created_at
			 FROM security_events ORDER BY created_at DESC LIMIT $1`, limit)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var events []event
		for rows.Next() {
			var e event
			var ts time.Time
			if err := rows.Scan(&e.ID, &e.EventID, &e.RequestID, &e.Severity, &e.DecisionAction,
				&e.Reason, &e.RuleIDs, &e.ClientIP, &e.Method, &e.Path, &e.Status, &ts); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			e.CreatedAt = ts.Format(time.RFC3339)
			events = append(events, e)
		}
		if events == nil {
			events = []event{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	}
}