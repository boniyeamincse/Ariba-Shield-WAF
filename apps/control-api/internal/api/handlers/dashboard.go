package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ariba-shield/control-api/internal/store"
)

// dashboardPeriod returns the number of days to look back (default 7, max 90).
func dashboardPeriod(r *http.Request) int {
	days := 7
	if v := r.URL.Query().Get("days"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 90 {
			days = n
		}
	}
	return days
}

// dashboardOverview returns high-level counts for the Overview widget.
func DashboardOverview(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		days := dashboardPeriod(r)
		var totalEvents, blockedEvents, totalRequests int
		var totalApps, totalGateways, activeIncidents int

		_ = st.Pool.QueryRow(r.Context(),
			`SELECT COUNT(*), COUNT(*) FILTER (WHERE decision_action = 'block')
			 FROM security_events WHERE created_at > now() - make_interval(days => $1)`, days).Scan(&totalEvents, &blockedEvents)
		_ = st.Pool.QueryRow(r.Context(),
			`SELECT COUNT(*) FROM access_events WHERE created_at > now() - make_interval(days => $1)`, days).Scan(&totalRequests)
		_ = st.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM applications WHERE status = 'active'`).Scan(&totalApps)
		_ = st.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM gateways WHERE status IN ('active','starting')`).Scan(&totalGateways)
		_ = st.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM incidents WHERE status != 'resolved'`).Scan(&activeIncidents)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"period_days":      days,
			"total_events":     totalEvents,
			"blocked_events":   blockedEvents,
			"total_requests":   totalRequests,
			"applications":     totalApps,
			"gateways":         totalGateways,
			"active_incidents": activeIncidents,
		})
	}
}

// dashboardTraffic returns request volume, status distribution, and latency.
func DashboardTraffic(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		days := dashboardPeriod(r)
		var totalRequests int
		var avgLatency, p99Latency float64

		_ = st.Pool.QueryRow(r.Context(),
			`SELECT COUNT(*), COALESCE(AVG(latency_ms),0), COALESCE(PERCENTILE_CONT(0.99) WITHIN GROUP (ORDER BY latency_ms),0)
			 FROM access_events WHERE created_at > now() - make_interval(days => $1)`, days).Scan(&totalRequests, &avgLatency, &p99Latency)

		byStatus := []map[string]any{}
		rows, err := st.Pool.Query(r.Context(),
			`SELECT COALESCE(status,0)::text, COUNT(*) FROM access_events
			 WHERE created_at > now() - make_interval(days => $1)
			 GROUP BY status ORDER BY COUNT(*) DESC`, days)
		if err == nil {
			for rows.Next() {
				var status string
				var cnt int
				if err := rows.Scan(&status, &cnt); err == nil {
					byStatus = append(byStatus, map[string]any{"status": status, "count": cnt})
				}
			}
			rows.Close()
		}
		if byStatus == nil {
			byStatus = []map[string]any{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"period_days":    days,
			"total_requests": totalRequests,
			"avg_latency_ms": avgLatency,
			"p99_latency_ms": p99Latency,
			"by_status":      byStatus,
		})
	}
}

// dashboardSecurity returns security event volume and severity distribution.
func DashboardSecurity(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		days := dashboardPeriod(r)
		var totalEvents, blockedEvents, uniqueIPs int

		_ = st.Pool.QueryRow(r.Context(),
			`SELECT COUNT(*), COUNT(*) FILTER (WHERE decision_action = 'block'), COUNT(DISTINCT client_ip)
			 FROM security_events WHERE created_at > now() - make_interval(days => $1)`, days).Scan(&totalEvents, &blockedEvents, &uniqueIPs)

		bySeverity := []map[string]any{}
		rows, err := st.Pool.Query(r.Context(),
			`SELECT COALESCE(severity,'unknown'), COUNT(*) FROM security_events
			 WHERE created_at > now() - make_interval(days => $1)
			 GROUP BY severity ORDER BY COUNT(*) DESC`, days)
		if err == nil {
			for rows.Next() {
				var sev string
				var cnt int
				if err := rows.Scan(&sev, &cnt); err == nil {
					bySeverity = append(bySeverity, map[string]any{"severity": sev, "count": cnt})
				}
			}
			rows.Close()
		}
		if bySeverity == nil {
			bySeverity = []map[string]any{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"period_days":    days,
			"total_events":   totalEvents,
			"blocked_events": blockedEvents,
			"unique_ips":     uniqueIPs,
			"by_severity":    bySeverity,
		})
	}
}

// dashboardAttacks returns the top attack types (by reason) in the period.
func DashboardAttacks(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		days := dashboardPeriod(r)
		byType := []map[string]any{}
		rows, err := st.Pool.Query(r.Context(),
			`SELECT COALESCE(reason,'unknown'), COUNT(*) FROM security_events
			 WHERE created_at > now() - make_interval(days => $1)
			 GROUP BY reason ORDER BY COUNT(*) DESC LIMIT 20`, days)
		if err == nil {
			for rows.Next() {
				var attackType, cnt string
				if err := rows.Scan(&attackType, &cnt); err == nil {
					byType = append(byType, map[string]any{"type": attackType, "count": cnt})
				}
			}
			rows.Close()
		}
		if byType == nil {
			byType = []map[string]any{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"period_days": days,
			"attacks":     byType,
		})
	}
}

// dashboardTopIPs returns the top client IPs by security event volume.
func DashboardTopIPs(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		days := dashboardPeriod(r)
		top := []map[string]any{}
		rows, err := st.Pool.Query(r.Context(),
			`SELECT COALESCE(client_ip,'unknown'), COUNT(*) as hits, COUNT(*) FILTER (WHERE decision_action = 'block') as blocked
			 FROM security_events
			 WHERE created_at > now() - make_interval(days => $1)
			 GROUP BY client_ip ORDER BY hits DESC LIMIT 20`, days)
		if err == nil {
			for rows.Next() {
				var ip string
				var hits, blocked int
				if err := rows.Scan(&ip, &hits, &blocked); err == nil {
					top = append(top, map[string]any{"client_ip": ip, "hits": hits, "blocked": blocked})
				}
			}
			rows.Close()
		}
		if top == nil {
			top = []map[string]any{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"period_days": days,
			"top_ips":     top,
		})
	}
}

// dashboardTopRules returns the top rules by hit count in the period.
func DashboardTopRules(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		days := dashboardPeriod(r)
		top := []map[string]any{}
		rows, err := st.Pool.Query(r.Context(),
			`SELECT unnest(rule_ids) AS rule_id, COUNT(*) AS hits
			 FROM security_events
			 WHERE created_at > now() - make_interval(days => $1) AND rule_ids IS NOT NULL
			 GROUP BY rule_id ORDER BY hits DESC LIMIT 20`, days)
		if err == nil {
			for rows.Next() {
				var ruleID string
				var hits int
				if err := rows.Scan(&ruleID, &hits); err == nil {
					top = append(top, map[string]any{"rule_id": ruleID, "hits": hits})
				}
			}
			rows.Close()
		}
		if top == nil {
			top = []map[string]any{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"period_days": days,
			"top_rules":   top,
		})
	}
}

// dashboardApplications returns per-application traffic + security counts.
func DashboardApplications(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		days := dashboardPeriod(r)
		apps := []map[string]any{}
		rows, err := st.Pool.Query(r.Context(),
			`SELECT a.id, a.name, a.status,
			        (SELECT COUNT(*) FROM access_events ae WHERE ae.application_id = a.id AND ae.created_at > now() - make_interval(days => $1)) AS requests,
			        (SELECT COUNT(*) FROM security_events se WHERE se.application_id = a.id AND se.created_at > now() - make_interval(days => $1)) AS events,
			        (SELECT COUNT(*) FROM security_events se2 WHERE se2.application_id = a.id AND se2.decision_action = 'block' AND se2.created_at > now() - make_interval(days => $1)) AS blocked
			 FROM applications a
			 ORDER BY requests DESC`, days)
		if err == nil {
			for rows.Next() {
				var id, name, status string
				var requests, events, blocked int
				if err := rows.Scan(&id, &name, &status, &requests, &events, &blocked); err == nil {
					apps = append(apps, map[string]any{
						"id": id, "name": name, "status": status,
						"requests": requests, "events": events, "blocked": blocked,
					})
				}
			}
			rows.Close()
		}
		if apps == nil {
			apps = []map[string]any{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"period_days":  days,
			"applications": apps,
		})
	}
}

// dashboardGateways returns gateway fleet status.
func DashboardGateways(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gateways := []map[string]any{}
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, COALESCE(hostname,''), COALESCE(ip,''), COALESCE(version,''), COALESCE(status,''), COALESCE(last_seen_at::text,'')
			 FROM gateways ORDER BY created_at ASC`)
		if err == nil {
			for rows.Next() {
				var id, hostname, ip, version, status, lastSeen string
				if err := rows.Scan(&id, &hostname, &ip, &version, &status, &lastSeen); err == nil {
					gateways = append(gateways, map[string]any{
						"id": id, "hostname": hostname, "ip": ip, "version": version,
						"status": status, "last_seen_at": lastSeen,
					})
				}
			}
			rows.Close()
		}
		if gateways == nil {
			gateways = []map[string]any{}
		}

		var active, offline, total int
		_ = st.Pool.QueryRow(r.Context(),
			`SELECT COUNT(*), COUNT(*) FILTER (WHERE status IN ('active','starting')), COUNT(*) FILTER (WHERE status IN ('offline','unregistered','degraded')) FROM gateways`).Scan(&total, &active, &offline)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"gateways": gateways,
			"total":    total,
			"active":   active,
			"offline":  offline,
		})
	}
}
