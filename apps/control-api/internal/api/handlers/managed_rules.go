package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

type managedRule struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Category        string `json:"category"`
	Enabled         bool   `json:"enabled"`
	Sensitivity     string `json:"sensitivity"`
	ParanoiaLevel   int    `json:"paranoia_level"`
	Action          string `json:"action"`
	AnomalyThreshold int   `json:"anomaly_threshold"`
	Status          string `json:"status"`
}

// ListManagedRules returns managed rule sets (e.g. OWASP CRS categories).
func ListManagedRules(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, category, enabled, COALESCE(sensitivity,''), COALESCE(paranoia_level,1),
			        COALESCE(action,'block'), COALESCE(anomaly_threshold,5), status
			 FROM managed_rules ORDER BY name`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		rules := []managedRule{}
		for rows.Next() {
			var m managedRule
			if err := rows.Scan(&m.ID, &m.Name, &m.Category, &m.Enabled, &m.Sensitivity,
				&m.ParanoiaLevel, &m.Action, &m.AnomalyThreshold, &m.Status); err != nil {
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

// ConfigureManagedRules enables/disables or sets sensitivity/anomaly for a managed set.
func ConfigureManagedRules(st *store.Store) http.HandlerFunc {
	type configure struct {
		Enabled          *bool   `json:"enabled"`
		Sensitivity      *string `json:"sensitivity"`
		ParanoiaLevel    *int    `json:"paranoia_level"`
		Action           *string `json:"action"`
		AnomalyThreshold *int    `json:"anomaly_threshold"`
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
			   paranoia_level = COALESCE($3, paranoia_level),
			   action = COALESCE($4, action),
			   anomaly_threshold = COALESCE($5, anomaly_threshold),
			   updated_at = now()
			 WHERE id = $6`,
			body.Enabled, nullableString(body.Sensitivity), body.ParanoiaLevel,
			nullableString(body.Action), body.AnomalyThreshold, id)
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

// ManageCRSGlobally toggles the whole OWASP CRS and sets global paranoia/anomaly.
func ManageCRSGlobally(st *store.Store) http.HandlerFunc {
	type global struct {
		Enabled          *bool   `json:"enabled"`
		ParanoiaLevel    *int    `json:"paranoia_level"`
		AnomalyThreshold *int    `json:"anomaly_threshold"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body global
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Enabled != nil {
			if _, err := st.Pool.Exec(r.Context(),
				`UPDATE managed_rules SET enabled = $1, updated_at = now() WHERE category = 'owasp-crs'`,
				*body.Enabled); err != nil {
				http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
				return
			}
		}
		if body.ParanoiaLevel != nil {
			if _, err := st.Pool.Exec(r.Context(),
				`UPDATE managed_rules SET paranoia_level = $1, sensitivity = CASE $1 WHEN 1 THEN 'low' WHEN 2 THEN 'medium' WHEN 3 THEN 'high' WHEN 4 THEN 'strict' ELSE sensitivity END, updated_at = now() WHERE category = 'owasp-crs'`,
				*body.ParanoiaLevel); err != nil {
				http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
				return
			}
		}
		if body.AnomalyThreshold != nil {
			if _, err := st.Pool.Exec(r.Context(),
				`UPDATE managed_rules SET anomaly_threshold = $1, updated_at = now() WHERE category = 'owasp-crs'`,
				*body.AnomalyThreshold); err != nil {
				http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	}
}