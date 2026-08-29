package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// GetSite returns a single site by id.
func GetSite(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var site struct {
			ID          string   `json:"id"`
			Name        string   `json:"name"`
			Description string   `json:"description"`
			Location    string   `json:"location"`
			CountryCode string   `json:"country_code"`
			GatewayIDs  []string `json:"gateway_ids"`
			Status      string   `json:"status"`
		}
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, name, COALESCE(description,''), COALESCE(location,''),
			        COALESCE(country_code,''), gateway_ids, status
			 FROM sites WHERE id = $1`, id).
			Scan(&site.ID, &site.Name, &site.Description, &site.Location,
				&site.CountryCode, &site.GatewayIDs, &site.Status)
		if err != nil {
			http.Error(w, `{"error":"site not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(site)
	}
}

// SiteHealth returns the aggregated health of a site based on its gateways.
func SiteHealth(st *store.Store) http.HandlerFunc {
	type gatewayHealth struct {
		ID         string `json:"id"`
		Hostname   string `json:"hostname"`
		Status     string `json:"status"`
		LastSeenAt *string `json:"last_seen_at,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		var gatewayIDs []string
		var siteName string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT name, gateway_ids FROM sites WHERE id = $1`, id).Scan(&siteName, &gatewayIDs); err != nil {
			http.Error(w, `{"error":"site not found"}`, http.StatusNotFound)
			return
		}

		health := map[string]any{
			"site_id":   id,
			"site_name": siteName,
			"overall":   "unknown",
			"gateways":  []gatewayHealth{},
		}

		if len(gatewayIDs) == 0 {
			health["overall"] = "no_gateways"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(health)
			return
		}

		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, hostname, status, last_seen_at FROM gateways WHERE id = ANY($1)`, gatewayIDs)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var gateways []gatewayHealth
		up := 0
		for rows.Next() {
			var g gatewayHealth
			var lastSeen *string
			if err := rows.Scan(&g.ID, &g.Hostname, &g.Status, &lastSeen); err == nil {
				if lastSeen != nil {
					g.LastSeenAt = lastSeen
				}
				if g.Status == "active" || g.Status == "starting" {
					up++
				}
				gateways = append(gateways, g)
			}
		}

		total := len(gateways)
		overall := "degraded"
		if total > 0 && up == total {
			overall = "healthy"
		} else if up == 0 {
			overall = "down"
		}

		health["overall"] = overall
		health["gateways"] = gateways
		health["healthy"] = up
		health["total"] = total

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
	}
}