package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// GetApplication returns a single application by id.
func GetApplication(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		var app application
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, name, COALESCE(description,''), status, COALESCE(owner_user_id,''), version
			 FROM applications WHERE id = $1`, appID).
			Scan(&app.ID, &app.Name, &app.Description, &app.Status, &app.OwnerUserID, &app.Version)
		if err != nil {
			http.Error(w, `{"error":"application not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(app)
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

