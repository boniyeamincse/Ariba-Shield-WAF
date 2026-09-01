package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// GeoEnforcement returns Geo-IP enforcement configuration.
func GetGeoConfig(st *store.Store) http.HandlerFunc {
	type geo struct {
		ID              string   `json:"id"`
		Name            string   `json:"name"`
		BlockedCountries []string `json:"blocked_countries"`
		AllowedCountries []string `json:"allowed_countries"`
		Action          string   `json:"action"`
		Status          string   `json:"status"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, COALESCE(blocked_countries,'{}'), COALESCE(allowed_countries,'{}'), action, status
			 FROM geo_blocking ORDER BY name`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		list := []geo{}
		for rows.Next() {
			var g geo
			var blocked, allowed []string
			if err := rows.Scan(&g.ID, &g.Name, &blocked, &allowed, &g.Action, &g.Status); err == nil {
				g.BlockedCountries = blocked
				g.AllowedCountries = allowed
				list = append(list, g)
			}
		}
		if list == nil {
			list = []geo{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	}
}
