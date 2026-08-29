package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// groupListItem is the response shape for listing groups.
type groupListItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	OrganizationID string `json:"organization_id"`
	CreatedAt string `json:"created_at"`
}

// Group represents a group entity.
type Group struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	OrganizationID string `json:"organization_id"`
	CreatedAt   string `json:"created_at"`
}

// ListGroups returns all groups for the organization.
func ListGroups(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, organization_id, created_at FROM groups ORDER BY created_at DESC`)
		if err != nil {
			http.Error(w, `{"error":"db query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var groups []groupListItem
		for rows.Next() {
			var g groupListItem
			if err := rows.Scan(&g.ID, &g.Name, &g.OrganizationID, &g.CreatedAt); err != nil {
				continue
			}
			groups = append(groups, g)
		}
		if groups == nil {
			groups = []groupListItem{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"groups": groups})
	}
}

// GetGroup returns a single group by ID.
func GetGroup(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var idVar, name, orgID, created string
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, name, organization_id, created_at FROM groups WHERE id = $1`, id).Scan(&idVar, &name, &orgID, &created)
		if err != nil {
			http.Error(w, `{"error":"group not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":         idVar,
			"name":       name,
			"organization_id": orgID,
			"created_at": created,
		})
	}
}

// CreateGroup creates a new group.
func CreateGroup(st *store.Store) http.HandlerFunc {
	type createGroupReq struct {
		Name string `json:"name"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req createGroupReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if req.Name == "" {
			http.Error(w, `{"error":"name is required"}`, http.StatusBadRequest)
			return
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		_, err = st.Pool.Exec(r.Context(),
			`INSERT INTO groups (id, organization_id, name, created_at)
			 VALUES ($1, $2, $3, now())`,
			id, "01ARZ3NDEKTSV4RRFFQ69G5FAV", req.Name)
		if err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// UpdateGroup updates a group's name.
func UpdateGroup(st *store.Store) http.HandlerFunc {
	type updateGroupReq struct {
		Name string `json:"name"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var req updateGroupReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if req.Name == "" {
			http.Error(w, `{"error":"name is required"}`, http.StatusBadRequest)
			return
		}
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE groups SET name = $1 WHERE id = $2`, req.Name, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"group not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// DeleteGroup deletes a group and its memberships.
func DeleteGroup(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		// Delete group memberships first, then the group.
		st.Pool.Exec(r.Context(), `DELETE FROM user_group_memberships WHERE group_id = $1`, id)
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM groups WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"group not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}

// roleListItem is the response shape for listing roles.
type roleListItem struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Permissions []string `json:"permissions"`
}

// ListRoles returns all roles.
func ListRoles(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, permissions FROM roles ORDER BY name`)
		if err != nil {
			http.Error(w, `{"error":"db query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var roles []roleListItem
		for rows.Next() {
			var id, name string
			// pgx v5 supports []string for TEXT[] columns natively.
			var perms []string
			if err := rows.Scan(&id, &name, &perms); err != nil {
				continue
			}
			roles = append(roles, roleListItem{
				ID:        id,
				Name:      name,
				Permissions: perms,
			})
		}
		if roles == nil {
			roles = []roleListItem{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"roles": roles})
	}
}

// GetRole returns a single role by ID.
func GetRole(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var idVar, name string
		var perms []string
		// pgx v5 supports []string for TEXT[] columns natively.
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, name, permissions FROM roles WHERE id = $1`, id).Scan(&idVar, &name, &perms)
		if err != nil {
			http.Error(w, `{"error":"role not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":        idVar,
			"name":      name,
			"permissions": perms,
		})
	}
}

// assignRoleToUserReq is the request body for assigning a role to a user.
type assignRoleToUserReq struct {
	RoleID string `json:"role_id"`
}

// AssignRoleToUser assigns a role to a user.
func AssignRoleToUser(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("id")
		var req assignRoleToUserReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if req.RoleID == "" {
			http.Error(w, `{"error":"role_id is required"}`, http.StatusBadRequest)
			return
		}
		// Check role exists.
		_, err := st.RoleID(r.Context(), req.RoleID)
		if err != nil {
			http.Error(w, `{"error":"role not found"}`, http.StatusNotFound)
			return
		}
		err = st.AssignRole(r.Context(), userID, req.RoleID)
		if err != nil {
			http.Error(w, `{"error":"role assignment failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": userID, "role_id": req.RoleID})
	}
}

// removeRoleFromUserReq is the request body for removing a role from a user.
type removeRoleFromUserReq struct {
	RoleID string `json:"role_id"`
}

// RemoveRoleFromUser removes a role from a user.
func RemoveRoleFromUser(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("id")
		var req removeRoleFromUserReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if req.RoleID == "" {
			http.Error(w, `{"error":"role_id is required"}`, http.StatusBadRequest)
			return
		}
		_, err := st.Pool.Exec(r.Context(),
			`DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`, userID, req.RoleID)
		if err != nil {
			http.Error(w, `{"error":"role removal failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": userID, "role_id": req.RoleID})
	}
}
