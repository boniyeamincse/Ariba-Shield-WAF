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

// UpdateApplication updates an application (PATCH semantics).
func UpdateApplication(st *store.Store) http.HandlerFunc {
	type applicationUpdate struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Status      *string `json:"status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		var body applicationUpdate
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}

		// Verify exists + read current values for partial update.
		var name, description, status string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT name, COALESCE(description,''), status FROM applications WHERE id = $1`, appID).
			Scan(&name, &description, &status); err != nil {
			http.Error(w, `{"error":"application not found"}`, http.StatusNotFound)
			return
		}

		if body.Name != nil {
			name = *body.Name
		}
		if body.Description != nil {
			description = *body.Description
		}
		if body.Status != nil {
			if *body.Status != "active" && *body.Status != "disabled" {
				http.Error(w, `{"error":"status must be active or disabled"}`, http.StatusBadRequest)
				return
			}
			status = *body.Status
		}

		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE applications SET name = $1, description = $2, status = $3, version = version + 1, updated_at = now() WHERE id = $4`,
			name, description, status, appID); err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": appID, "name": name, "status": status})
	}
}

// DeleteApplication deletes an application.
func DeleteApplication(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM applications WHERE id = $1`, appID)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"application not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}