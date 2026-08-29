package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListIPLists returns all IP lists.
func ListIPLists(st *store.Store) http.HandlerFunc {
	type iplist struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		ListType    string   `json:"list_type"`
		Entries     []string `json:"entries"`
		Description string   `json:"description"`
		Version     int64    `json:"version"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, list_type, entries, description, version FROM ip_lists ORDER BY name`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var lists []iplist
		for rows.Next() {
			var l iplist
			if err := rows.Scan(&l.ID, &l.Name, &l.ListType, &l.Entries, &l.Description, &l.Version); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			lists = append(lists, l)
		}
		if lists == nil {
			lists = []iplist{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lists)
	}
}

// CreateIPList creates an IP list.
func CreateIPList(st *store.Store) http.HandlerFunc {
	type create struct {
		Name        string   `json:"name"`
		ListType    string   `json:"list_type"`
		Entries     []string `json:"entries"`
		Description string   `json:"description"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body create
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Name == "" || body.ListType == "" {
			http.Error(w, `{"error":"name and list_type required"}`, http.StatusBadRequest)
			return
		}
		if body.ListType != "allowed" && body.ListType != "blocked" {
			http.Error(w, `{"error":"list_type must be allowed or blocked"}`, http.StatusBadRequest)
			return
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}
		orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO ip_lists (id, organization_id, name, list_type, entries, description)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			id, orgID, body.Name, body.ListType, body.Entries, body.Description); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// UpdateIPList updates an IP list (entries, description, list type).
func UpdateIPList(st *store.Store) http.HandlerFunc {
	type update struct {
		Name        *string   `json:"name"`
		ListType    *string   `json:"list_type"`
		Entries     *[]string `json:"entries"`
		Description *string   `json:"description"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body update
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.ListType != nil && *body.ListType != "allowed" && *body.ListType != "blocked" {
			http.Error(w, `{"error":"list_type must be allowed or blocked"}`, http.StatusBadRequest)
			return
		}

		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE ip_lists SET
			   name = COALESCE($1, name),
			   list_type = COALESCE($2, list_type),
			   entries = COALESCE($3, entries),
			   description = COALESCE($4, description),
			   version = version + 1, updated_at = now()
			 WHERE id = $5`,
			nullableString(body.Name), nullableString(body.ListType),
			body.Entries, nullableString(body.Description), id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"ip list not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// DeleteIPList deletes an IP list.
func DeleteIPList(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM ip_lists WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"ip list not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}

func nullableString(s *string) any {
	if s == nil {
		return nil
	}
	return *s
}