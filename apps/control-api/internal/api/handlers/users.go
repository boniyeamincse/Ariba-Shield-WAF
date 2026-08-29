package handlers

import (
	"encoding/json"
	"net/http"
	"time"

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
			SELECT id, email, status, created_at
			FROM users
			ORDER BY created_at ASC
		`)
		if err != nil {
			http.Error(w, `{"error":"db query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		users := []userListItem{}
		for rows.Next() {
			var id, email, status string
			var createdAt time.Time
			if err := rows.Scan(&id, &email, &status, &createdAt); err != nil {
				continue
			}
			role, _ := st.LookupUserRole(r.Context(), id)
			if role == "" {
				role = "Read Only"
			}
			users = append(users, userListItem{
				ID:        id,
				Email:     email,
				Status:    status,
				CreatedAt: createdAt.Format(time.RFC3339),
				Role:      role,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"users": users})
	}
}
