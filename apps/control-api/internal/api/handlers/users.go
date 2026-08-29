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
			var u userListItem
			if err := rows.Scan(&u.ID, &u.Email, &u.Status, &u.CreatedAt); err != nil {
				continue
			}
			role, err := st.LookupUserRole(r.Context(), u.ID)
			if err != nil {
				role = "Read Only"
			}
			u.Role = role
			users = append(users, u)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"users": users})
	}
}

// CreateUser creates a user and assigns a role. The role must be one of the
// seeded roles (Super Admin, Platform Admin, Security Admin, App Owner,
// SOC Analyst, Auditor, Read Only).
func CreateUser(st *store.Store) http.HandlerFunc {
	type createRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body createRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Email == "" || body.Password == "" {
			http.Error(w, `{"error":"email and password required"}`, http.StatusBadRequest)
			return
		}
		if body.Role == "" {
			body.Role = "Read Only"
		}

		// Validate the role exists in the seeded set.
		roleID, err := st.RoleID(r.Context(), body.Role)
		if err != nil {
			http.Error(w, `{"error":"invalid role"}`, http.StatusBadRequest)
			return
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		hash, err := store.HashPassword(body.Password)
		if err != nil {
			http.Error(w, `{"error":"password hashing failed"}`, http.StatusInternalServerError)
			return
		}

		orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"

		_, err = st.Pool.Exec(r.Context(),
			`INSERT INTO users (id, organization_id, email, password_hash, language, status)
			 VALUES ($1, $2, $3, $4, 'en', 'active')`,
			id, orgID, body.Email, hash)
		if err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		if err := st.AssignRole(r.Context(), id, roleID); err != nil {
			http.Error(w, `{"error":"role assignment failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id":     id,
			"email":  body.Email,
			"role":   body.Role,
			"status": "active",
		})
	}
}