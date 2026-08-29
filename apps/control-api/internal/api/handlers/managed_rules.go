package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

type managedRule struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Enabled     bool   `json:"enabled"`
	Sensitivity string `json:"sensitivity"`
	Status      string `json:"status"`
}

// ListManagedRules returns managed rule sets (e.g. OWASP CRS categories).
func ListManagedRules(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, category, enabled, sensitivity, status
			 FROM managed_rules ORDER BY name`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		rules := []managedRule{}
		for rows.Next() {
			var m managedRule
			if err := rows.Scan(&m.ID, &m.Name, &m.Category, &m.Enabled, &m.Sensitivity, &m.Status); err != nil {
				continue
			}
			rules = append(rules, m)
		}
		if rules == nil {
			rules = []managedRule{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rules)
	}
}

// ConfigureManagedRules enables/disables or sets sensitivity for a managed set.
func ConfigureManagedRules(st *store.Store) http.HandlerFunc {
	type configure struct {
		Enabled     *bool   `json:"enabled"`
		Sensitivity *string `json:"sensitivity"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body configure
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}

		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE managed_rules SET
			   enabled = COALESCE($1, enabled),
			   sensitivity = COALESCE($2, sensitivity),
			   updated_at = now()
			 WHERE id = $3`,
			body.Enabled, nullableString(body.Sensitivity), id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"managed rule not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}