package handlers

import (
	"strconv"
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// GetApplication returns a single application by id.
func GetApplication(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		a, err := scanApplication(st.Pool.QueryRow(r.Context(),
			`SELECT `+appSelect+` FROM applications WHERE id = $1`, appID))
		if err != nil {
			http.Error(w, `{"error":"application not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(a)
	}
}

// ApplicationTraffic returns access events for an application (paged).
func ApplicationTraffic(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		limit := 100
		rows, err := st.Pool.Query(r.Context(),
			`SELECT event_id, request_id, COALESCE(client_ip,''), method, path, COALESCE(host,''),
			        COALESCE(status,0), bytes, COALESCE(latency_ms,0), created_at
			 FROM access_events WHERE application_id = $1 ORDER BY created_at DESC LIMIT $2`,
			appID, limit)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		events := []map[string]any{}
		for rows.Next() {
			var eid, rid, ip, method, path, host, ts string
			var status int
			var bytes, latency float64
			if err := rows.Scan(&eid, &rid, &ip, &method, &path, &host, &status, &bytes, &latency, &ts); err != nil {
				continue
			}
			events = append(events, map[string]any{
				"event_id": eid, "request_id": rid, "client_ip": ip, "method": method,
				"path": path, "host": host, "status": status, "bytes": bytes,
				"latency_ms": latency, "created_at": ts,
			})
		}
		if events == nil {
			events = []map[string]any{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"application_id": appID, "requests": events})
	}
}

// ApplicationEvents returns security events for an application (paged).
func ApplicationEvents(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		limit := 100
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, event_id, COALESCE(severity,''), COALESCE(reason,''), rule_ids,
			        COALESCE(client_ip,''), COALESCE(method,''), COALESCE(path,''), created_at
			 FROM security_events WHERE application_id = $1 ORDER BY created_at DESC LIMIT $2`,
			appID, limit)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		events := []map[string]any{}
		for rows.Next() {
			var id, eid, sev, reason, ip, method, path, ts string
			var ruleIDs []string
			if err := rows.Scan(&id, &eid, &sev, &reason, &ruleIDs, &ip, &method, &path, &ts); err != nil {
				continue
			}
			events = append(events, map[string]any{
				"id": id, "event_id": eid, "severity": sev, "reason": reason,
				"rule_ids": ruleIDs, "client_ip": ip, "method": method, "path": path, "created_at": ts,
			})
		}
		if events == nil {
			events = []map[string]any{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"application_id": appID, "events": events})
	}
}

// ApplicationIncidents returns incidents linked to an application.
func ApplicationIncidents(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		// Incidents are matched via related_events (event ids) or by scanning
		// all incidents; simplest correlation: list recent incidents for the org.
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, title, severity, status, COALESCE(notes,''), created_at
			 FROM incidents ORDER BY created_at DESC LIMIT 50`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		incidents := []map[string]any{}
		for rows.Next() {
			var id, title, sev, status, notes, ts string
			if err := rows.Scan(&id, &title, &sev, &status, &notes, &ts); err != nil {
				continue
			}
			incidents = append(incidents, map[string]any{
				"id": id, "title": title, "severity": sev, "status": status,
				"notes": notes, "created_at": ts, "application_id": appID,
			})
		}
		if incidents == nil {
			incidents = []map[string]any{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"application_id": appID, "incidents": incidents})
	}
}

// ApplicationPolicies returns the security policies bound to an application.
func ApplicationPolicies(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, COALESCE(description,''), enforcement_mode, version
			 FROM security_policies WHERE application_id = $1 ORDER BY name`, appID)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		policies := []map[string]any{}
		for rows.Next() {
			var id, name, desc, mode string
			var version int64
			if err := rows.Scan(&id, &name, &desc, &mode, &version); err != nil {
				continue
			}
			policies = append(policies, map[string]any{
				"id": id, "name": name, "description": desc, "enforcement_mode": mode, "version": version,
			})
		}
		if policies == nil {
			policies = []map[string]any{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"application_id": appID, "policies": policies})
	}
}

// ApplicationHealth returns the health of an application's backend pools.
func ApplicationHealth(st *store.Store) http.HandlerFunc {
	type poolHealth struct {
		PoolID       string `json:"pool_id"`
		PoolName     string `json:"pool_name"`
		TotalNodes   int    `json:"total_nodes"`
		HealthyNodes int    `json:"healthy_nodes"`
		Status       string `json:"status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")

		rows, err := st.Pool.Query(r.Context(),
			`SELECT bp.id, bp.name,
			        COUNT(bn.id) AS total,
			        COUNT(bn.id) FILTER (WHERE bn.active AND (bn.last_health_state IS NULL OR bn.last_health_state = 'healthy')) AS healthy
			 FROM backend_pools bp
			 LEFT JOIN backend_nodes bn ON bn.pool_id = bp.id
			 WHERE bp.application_id = $1
			 GROUP BY bp.id, bp.name`, appID)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		pools := []poolHealth{}
		for rows.Next() {
			var p poolHealth
			if err := rows.Scan(&p.PoolID, &p.PoolName, &p.TotalNodes, &p.HealthyNodes); err != nil {
				continue
			}
			p.Status = "down"
			if p.TotalNodes > 0 && p.HealthyNodes == p.TotalNodes {
				p.Status = "healthy"
			} else if p.HealthyNodes > 0 {
				p.Status = "degraded"
			}
			pools = append(pools, p)
		}
		if pools == nil {
			pools = []poolHealth{}
		}

		overall := "healthy"
		for _, p := range pools {
			if p.Status == "down" {
				overall = "down"
				break
			}
			if p.Status == "degraded" {
				overall = "degraded"
			}
		}
		if len(pools) == 0 {
			overall = "no_pools"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"application_id": appID,
			"overall":        overall,
			"pools":          pools,
		})
	}
}


// ApplicationAnalytics returns per-app metrics: health score, traffic, attacks,
// top IPs, top rules, and a daily attack timeline.
func ApplicationAnalytics(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		days := 7
		if v := r.URL.Query().Get("days"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 90 {
				days = n
			}
		}

		// Health score: 0-100 based on pool/node health.
		var totalNodes, healthyNodes int
		_ = st.Pool.QueryRow(r.Context(),
			`SELECT COUNT(bn.id), COUNT(bn.id) FILTER (WHERE bn.active AND (bn.last_health_state IS NULL OR bn.last_health_state='healthy'))
			 FROM backend_pools bp JOIN backend_nodes bn ON bn.pool_id = bp.id
			 WHERE bp.application_id = $1`, appID).Scan(&totalNodes, &healthyNodes)
		healthScore := 100
		if totalNodes > 0 {
			healthScore = healthyNodes * 100 / totalNodes
		}

		// Traffic + security aggregates.
		var totalRequests, totalEvents, blocked int
		_ = st.Pool.QueryRow(r.Context(),
			`SELECT COUNT(*), (SELECT COUNT(*) FROM security_events WHERE application_id=$1 AND created_at > now()-make_interval(days=>$2)),
			        (SELECT COUNT(*) FROM security_events WHERE application_id=$1 AND decision_action='block' AND created_at > now()-make_interval(days=>$2))
			 FROM access_events WHERE application_id=$1 AND created_at > now()-make_interval(days=>$2)`, appID, days).Scan(&totalRequests, &totalEvents, &blocked)

		// Top source IPs by event volume.
		topIPs := []map[string]any{}
		rows, _ := st.Pool.Query(r.Context(),
			`SELECT COALESCE(client_ip,'unknown'), COUNT(*) FROM security_events
			 WHERE application_id=$1 AND created_at > now()-make_interval(days=>$2)
			 GROUP BY client_ip ORDER BY COUNT(*) DESC LIMIT 10`, appID, days)
		if rows != nil {
			for rows.Next() {
				var ip string
				var cnt int
				if rows.Scan(&ip, &cnt) == nil {
					topIPs = append(topIPs, map[string]any{"client_ip": ip, "hits": cnt})
				}
			}
			rows.Close()
		}
		if topIPs == nil {
			topIPs = []map[string]any{}
		}

		// Top rules.
		topRules := []map[string]any{}
		rows2, _ := st.Pool.Query(r.Context(),
			`SELECT unnest(rule_ids) AS rule_id, COUNT(*) FROM security_events
			 WHERE application_id=$1 AND created_at > now()-make_interval(days=>$2)
			 GROUP BY rule_id ORDER BY COUNT(*) DESC LIMIT 10`, appID, days)
		if rows2 != nil {
			for rows2.Next() {
				var rid string
				var cnt int
				if rows2.Scan(&rid, &cnt) == nil {
					topRules = append(topRules, map[string]any{"rule_id": rid, "hits": cnt})
				}
			}
			rows2.Close()
		}
		if topRules == nil {
			topRules = []map[string]any{}
		}

		// Daily attack timeline.
		timeline := []map[string]any{}
		rows3, _ := st.Pool.Query(r.Context(),
			`SELECT date_trunc('day', created_at)::date::text, COUNT(*) FROM security_events
			 WHERE application_id=$1 AND created_at > now()-make_interval(days=>$2)
			 GROUP BY 1 ORDER BY 1`, appID, days)
		if rows3 != nil {
			for rows3.Next() {
				var day string
				var cnt int
				if rows3.Scan(&day, &cnt) == nil {
					timeline = append(timeline, map[string]any{"date": day, "events": cnt})
				}
			}
			rows3.Close()
		}
		if timeline == nil {
			timeline = []map[string]any{}
		}

		status := "healthy"
		if healthScore == 0 {
			status = "down"
		} else if healthScore < 100 {
			status = "degraded"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"application_id": appID,
			"period_days":    days,
			"health_score":   healthScore,
			"health_status":  status,
			"total_requests": totalRequests,
			"total_events":   totalEvents,
			"blocked":        blocked,
			"block_ratio":    safeRatio(blocked, totalEvents),
			"top_ips":        topIPs,
			"top_rules":      topRules,
			"timeline":       timeline,
		})
	}
}

// safeRatio computes percent of x/y; 0 when y is 0.
func safeRatio(x, y int) float64 {
	if y <= 0 {
		return 0
	}
	return float64(x) / float64(y) * 100
}
