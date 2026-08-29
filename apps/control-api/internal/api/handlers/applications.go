package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

type application struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	OwnerUserID string `json:"owner_user_id,omitempty"`
	Version     int64  `json:"version"`
}

// ListApplications returns all applications visible to the caller.
func ListApplications(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, description, status, COALESCE(owner_user_id, ''), version FROM applications ORDER BY name`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var apps []application
		for rows.Next() {
			var a application
			if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.Status, &a.OwnerUserID, &a.Version); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			apps = append(apps, a)
		}
		if apps == nil {
			apps = []application{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(apps)
	}
}

// CreateApplication creates a new application.
func CreateApplication(st *store.Store) http.HandlerFunc {
	type applicationCreate struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body applicationCreate
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Name == "" {
			http.Error(w, `{"error":"name is required"}`, http.StatusBadRequest)
			return
		}

		appID, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		// Single-org for 0.1.
		orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO applications (id, organization_id, name, description) VALUES ($1, $2, $3, $4)`,
			appID, orgID, body.Name, body.Description); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": appID})
	}
}