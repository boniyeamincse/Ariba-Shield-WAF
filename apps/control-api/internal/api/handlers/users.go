package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

type userListItem struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// ListUsers returns all users in the organization.
func ListUsers(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(), `
			SELECT u.id, u.email, u.status, u.created_at,
			       COALESCE(r.name, 'Read Only') AS role
			FROM users u
			LEFT JOIN user_role_assignments ura ON ura.user_id = u.id
			LEFT JOIN roles r ON r.id = ura.role_id
			ORDER BY u.created_at ASC
		`)
		if err != nil {
			http.Error(w, `{"error":"db query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		users := []userListItem{}
		for rows.Next() {
			var u userListItem
			if err := rows.Scan(&u.ID, &u.Email, &u.Status, &u.CreatedAt, &u.Role); err != nil {
				continue
			}
			users = append(users, u)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"users": users})
	}
}
