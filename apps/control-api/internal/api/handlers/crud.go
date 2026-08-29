package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// CRUDConfig describes a simple resource table for the generic CRUD factory.
type CRUDConfig struct {
	// Table is the SQL table name.
	Table string
	// Columns is the list of mutable columns (excludes id, organization_id,
	// created_at, updated_at, version).
	Columns []string
	// JSONName is the singular resource name used in responses.
	JSONName string
	// Required are column names that must be present on create.
	Required []string
}

// GetResource returns a single row by id as JSON.
func GetResource(st *store.Store, cfg CRUDConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var b []byte
		err := st.Pool.QueryRow(r.Context(),
			"SELECT row_to_json(t) FROM (SELECT * FROM "+cfg.Table+" WHERE id = $1) t", id).Scan(&b)
		if err != nil {
			http.Error(w, `{"error":"`+cfg.JSONName+` not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

// ListResource returns a handler that lists all rows of a table.
func ListResource(st *store.Store, cfg CRUDConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			"SELECT row_to_json(t) FROM (SELECT * FROM "+cfg.Table+" ORDER BY created_at DESC) t")
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		items := []json.RawMessage{}
		for rows.Next() {
			var b []byte
			if err := rows.Scan(&b); err != nil {
				continue
			}
			items = append(items, json.RawMessage(b))
		}
		if items == nil {
			items = []json.RawMessage{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

// CreateResource returns a handler that creates a row from arbitrary JSON.
// The body is stored as-is for the mutable columns (whitelisted via Columns).
func CreateResource(st *store.Store, cfg CRUDConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		for _, req := range cfg.Required {
			if body[req] == nil {
				http.Error(w, `{"error":"`+req+` is required"}`, http.StatusBadRequest)
				return
			}
		}

		// Build the column list + placeholders.
		cols := []string{"id", "organization_id"}
		vals := []any{}
		placeholders := []string{"$1", "$2"}
		vals = append(vals, "") // id set below
		vals = append(vals, "01ARZ3NDEKTSV4RRFFQ69G5FAV")

		for i, c := range cfg.Columns {
			if v, ok := body[c]; ok {
				cols = append(cols, c)
				vals = append(vals, v)
				placeholders = append(placeholders, "$"+itoa(len(vals)))
			}
			_ = i
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}
		vals[0] = id

		query := "INSERT INTO " + cfg.Table + " (" + joinComma(cols) + ") VALUES (" + joinComma(placeholders) + ")"
		if _, err := st.Pool.Exec(r.Context(), query, vals...); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// UpdateResource returns a handler that partially updates a row by id.
func UpdateResource(st *store.Store, cfg CRUDConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}

		sets := []string{}
		vals := []any{}
		for _, c := range cfg.Columns {
			if v, ok := body[c]; ok {
				vals = append(vals, v)
				sets = append(sets, c+" = $"+itoa(len(vals)))
			}
		}
		if len(sets) == 0 {
			http.Error(w, `{"error":"no updatable fields"}`, http.StatusBadRequest)
			return
		}
		vals = append(vals, id)
		query := "UPDATE " + cfg.Table + " SET " + joinComma(sets) +
			", updated_at = now() WHERE id = $" + itoa(len(vals))

		ct, err := st.Pool.Exec(r.Context(), query, vals...)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"`+cfg.JSONName+` not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// DeleteResource returns a handler that deletes a row by id.
func DeleteResource(st *store.Store, cfg CRUDConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), "DELETE FROM "+cfg.Table+" WHERE id = $1", id)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"`+cfg.JSONName+` not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}

func joinComma(s []string) string {
	out := ""
	for i, v := range s {
		if i > 0 {
			out += ", "
		}
		out += v
	}
	return out
}
