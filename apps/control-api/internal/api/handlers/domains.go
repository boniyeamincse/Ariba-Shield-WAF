package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListDomains returns the domains for an application.
func ListDomains(st *store.Store) http.HandlerFunc {
	type domain struct {
		ID        string `json:"id"`
		Hostname  string `json:"hostname"`
		Enabled   bool   `json:"enabled"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, hostname, enabled FROM domains WHERE application_id = $1 ORDER BY hostname`, appID)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var domains []domain
		for rows.Next() {
			var d domain
			if err := rows.Scan(&d.ID, &d.Hostname, &d.Enabled); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			domains = append(domains, d)
		}
		if domains == nil {
			domains = []domain{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(domains)
	}
}

// CreateDomain adds a domain to an application.
func CreateDomain(st *store.Store) http.HandlerFunc {
	type domainCreate struct {
		Hostname string `json:"hostname"`
		Enabled  *bool  `json:"enabled"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")

		var body domainCreate
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Hostname == "" {
			http.Error(w, `{"error":"hostname is required"}`, http.StatusBadRequest)
			return
		}
		enabled := true
		if body.Enabled != nil {
			enabled = *body.Enabled
		}

		var orgID string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT organization_id FROM applications WHERE id = $1`, appID).Scan(&orgID); err != nil {
			http.Error(w, `{"error":"application not found"}`, http.StatusNotFound)
			return
		}

		domainID, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO domains (id, organization_id, application_id, hostname, enabled)
			 VALUES ($1, $2, $3, $4, $5)`,
			domainID, orgID, appID, body.Hostname, enabled); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": domainID})
	}
}