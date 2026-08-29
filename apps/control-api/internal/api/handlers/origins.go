package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListOrigins returns the origins (upstreams) for an application.
func ListOrigins(st *store.Store) http.HandlerFunc {
	type origin struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Protocol string `json:"protocol"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Weight   int    `json:"weight"`
		Enabled  bool   `json:"enabled"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, protocol, host, port, weight, enabled
			 FROM origins WHERE application_id = $1 ORDER BY name`, appID)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var origins []origin
		for rows.Next() {
			var o origin
			if err := rows.Scan(&o.ID, &o.Name, &o.Protocol, &o.Host, &o.Port, &o.Weight, &o.Enabled); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			origins = append(origins, o)
		}
		if origins == nil {
			origins = []origin{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(origins)
	}
}

// CreateOrigin adds an origin (upstream) to an application.
func CreateOrigin(st *store.Store) http.HandlerFunc {
	type originCreate struct {
		Name     string `json:"name"`
		Protocol string `json:"protocol"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Weight   int    `json:"weight"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")

		var body originCreate
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Name == "" || body.Host == "" {
			http.Error(w, `{"error":"name and host are required"}`, http.StatusBadRequest)
			return
		}
		if body.Protocol == "" {
			body.Protocol = "http"
		}
		if body.Port == 0 {
			if body.Protocol == "https" {
				body.Port = 443
			} else {
				body.Port = 80
			}
		}
		if body.Weight == 0 {
			body.Weight = 1
		}

		var orgID string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT organization_id FROM applications WHERE id = $1`, appID).Scan(&orgID); err != nil {
			http.Error(w, `{"error":"application not found"}`, http.StatusNotFound)
			return
		}

		originID, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO origins (id, organization_id, application_id, name, protocol, host, port, weight)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			originID, orgID, appID, body.Name, body.Protocol, body.Host, body.Port, body.Weight); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": originID})
	}
}