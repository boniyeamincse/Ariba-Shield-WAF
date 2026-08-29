package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListSecurityEvents returns security events with pagination + filters.
// Query params: limit (default 50, max 200), offset (default 0),
// application_id, severity, client_ip, rule_id.
func ListSecurityEvents(st *store.Store) http.HandlerFunc {
	type event struct {
		ID             string   `json:"id"`
		EventID        string   `json:"event_id"`
		RequestID      string   `json:"request_id"`
		Timestamp      string   `json:"timestamp"`
		Severity       string   `json:"severity"`
		DecisionAction string   `json:"decision_action"`
		Reason         string   `json:"reason"`
		RuleIDs        []string `json:"rule_ids"`
		ClientIP       string   `json:"client_ip"`
		Method         string   `json:"method"`
		Path           string   `json:"path"`
		Status         int      `json:"status"`
		CreatedAt      string   `json:"created_at"`
	}

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

		// Build a parameterized WHERE clause from filters.
		where := "1=1"
		args := []any{}
		add := func(cond string, val any) {
			args = append(args, val)
			where += " AND " + cond
		}

		if v := r.URL.Query().Get("application_id"); v != "" {
			add("application_id = $"+strconv.Itoa(len(args)+1), v)
		}
		if v := r.URL.Query().Get("gateway_id"); v != "" {
			add("gateway_id = $"+strconv.Itoa(len(args)+1), v)
		}
		if v := r.URL.Query().Get("ip"); v != "" {
			add("client_ip = $"+strconv.Itoa(len(args)+1), v)
		}
		if v := r.URL.Query().Get("attack_type"); v != "" {
			add("reason ILIKE $"+strconv.Itoa(len(args)+1), "%"+v+"%")
		}
		if v := r.URL.Query().Get("severity"); v != "" {
			add("severity = $"+strconv.Itoa(len(args)+1), v)
		}
		if v := r.URL.Query().Get("action"); v != "" {
			add("decision_action = $"+strconv.Itoa(len(args)+1), v)
		}
		if v := r.URL.Query().Get("from"); v != "" {
			add("created_at >= $"+strconv.Itoa(len(args)+1)+"::timestamptz", v)
		}
		if v := r.URL.Query().Get("to"); v != "" {
			add("created_at <= $"+strconv.Itoa(len(args)+1)+"::timestamptz", v)
		}
		if v := r.URL.Query().Get("client_ip"); v != "" {
			add("client_ip = $"+strconv.Itoa(len(args)+1), v)
		}
		if v := r.URL.Query().Get("rule_id"); v != "" {
			add("$"+strconv.Itoa(len(args)+1)+" = ANY(rule_ids)", v)
		}

		args = append(args, limit, offset)
		query := `SELECT id, event_id, request_id, COALESCE(severity,''), COALESCE(decision_action,''),
			        COALESCE(reason,''), rule_ids, COALESCE(client_ip,''), COALESCE(method,''),
			        COALESCE(path,''), COALESCE(status,0), created_at
			 FROM security_events WHERE ` + where +
			` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(len(args)-1) +
			` OFFSET $` + strconv.Itoa(len(args))

		rows, err := st.Pool.Query(r.Context(), query, args...)
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
		json.NewEncoder(w).Encode(map[string]any{
			"events": events,
			"pagination": map[string]int{
				"limit":  limit,
				"offset": offset,
				"count":  len(events),
			},
		})
	}
}