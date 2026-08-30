package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

type customRule struct {
	ID              string          `json:"id"`
	PolicyID        string          `json:"policy_id,omitempty"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Action          string          `json:"action"`
	MatchConditions json.RawMessage `json:"match_conditions,omitempty"`
	Status          string          `json:"status"`
}

// ListCustomRules returns user-defined custom rules.
func ListCustomRules(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, COALESCE(policy_id,''), name, description, action, COALESCE(match_conditions,'{}'::jsonb), status
			 FROM custom_rules ORDER BY name`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		rules := []customRule{}
		for rows.Next() {
			var c customRule
			if err := rows.Scan(&c.ID, &c.PolicyID, &c.Name, &c.Description, &c.Action, &c.MatchConditions, &c.Status); err != nil {
				continue
			}
			rules = append(rules, c)
		}
		if rules == nil {
			rules = []customRule{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rules)
	}
}

// CreateCustomRule creates a custom regex/condition rule.
func CreateCustomRule(st *store.Store) http.HandlerFunc {
	type create struct {
		PolicyID        string          `json:"policy_id"`
		Name            string          `json:"name"`
		Description     string          `json:"description"`
		Action          string          `json:"action"`
		MatchConditions json.RawMessage `json:"match_conditions"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body create
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Name == "" {
			http.Error(w, `{"error":"name required"}`, http.StatusBadRequest)
			return
		}
		if body.Action == "" {
			body.Action = "BLOCK"
		}
		switch body.Action {
		case "BLOCK", "LOG", "ALLOW":
		default:
			http.Error(w, `{"error":"action must be BLOCK, LOG or ALLOW"}`, http.StatusBadRequest)
			return
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO custom_rules (id, organization_id, policy_id, name, description, action, match_conditions)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			id, orgID, nullIfEmpty(body.PolicyID), body.Name, body.Description, body.Action, body.MatchConditions); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}
