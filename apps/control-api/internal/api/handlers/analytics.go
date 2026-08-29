package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// SecurityAnalytics returns aggregate security event stats (Phase 4).
func SecurityAnalytics(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		days := 7
		var totalEvents, blockedEvents, uniqueIPs int
		_ = st.Pool.QueryRow(r.Context(),
			`SELECT COUNT(*), COUNT(*) FILTER (WHERE decision_action = 'block'),
			        COUNT(DISTINCT client_ip)
			 FROM security_events WHERE created_at > now() - make_interval(days => $1)`, days).Scan(&totalEvents, &blockedEvents, &uniqueIPs)

		bySeverity := []map[string]any{}
		rows, _ := st.Pool.Query(r.Context(),
			`SELECT COALESCE(severity,'unknown'), COUNT(*) FROM security_events
			 WHERE created_at > now() - make_interval(days => $1)
			 GROUP BY severity ORDER BY COUNT(*) DESC`, days)
		if rows != nil {
			for rows.Next() {
				var sev string
				var cnt int
				if err := rows.Scan(&sev, &cnt); err == nil {
					bySeverity = append(bySeverity, map[string]any{"severity": sev, "count": cnt})
				}
			}
			rows.Close()
		}

		byType := []map[string]any{}
		rows2, _ := st.Pool.Query(r.Context(),
			`SELECT COALESCE(reason,'unknown'), COUNT(*) FROM security_events
			 WHERE created_at > now() - make_interval(days => $1)
			 GROUP BY reason ORDER BY COUNT(*) DESC LIMIT 10`, days)
		if rows2 != nil {
			for rows2.Next() {
				var r, cnt string
				if err := rows2.Scan(&r, &cnt); err == nil {
					byType = append(byType, map[string]any{"type": r, "count": cnt})
				}
			}
			rows2.Close()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"period_days":    days,
			"total_events":   totalEvents,
			"blocked_events": blockedEvents,
			"unique_ips":     uniqueIPs,
			"by_severity":    bySeverity,
			"by_reason":      byType,
		})
	}
}

// RuleAnalytics returns rule hit statistics (Phase 4).
func RuleAnalytics(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		days := 7
		// De-duplicate: each event has rule_ids[]; count each match.
		rows, err := st.Pool.Query(r.Context(),
			`SELECT unnest(rule_ids) AS rule_id, COUNT(*) AS hits
			 FROM security_events
			 WHERE created_at > now() - make_interval(days => $1) AND rule_ids IS NOT NULL
			 GROUP BY rule_id ORDER BY hits DESC LIMIT 20`, days)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		rules := []map[string]any{}
		for rows.Next() {
			var ruleID string
			var hits int
			if err := rows.Scan(&ruleID, &hits); err == nil {
				rules = append(rules, map[string]any{"rule_id": ruleID, "hits": hits})
			}
		}
		if rules == nil {
			rules = []map[string]any{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"period_days": days,
			"rules":       rules,
		})
	}
}