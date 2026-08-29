package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListApplications returns all applications visible to the caller.
func ListApplications(st *store.Store) http.HandlerFunc {
	type app struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status"`
		OwnerUserID string `json:"owner_user_id,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, description, status, COALESCE(owner_user_id, '') FROM applications ORDER BY name`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var apps []app
		for rows.Next() {
			var a app
			if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.Status, &a.OwnerUserID); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			apps = append(apps, a)
		}
		if apps == nil {
			apps = []app{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(apps)
	}
}