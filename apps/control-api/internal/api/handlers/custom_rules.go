package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

type MatchCondition struct {
	Variable string `json:"variable"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type CustomRule struct {
	ID              string           `json:"id"`
	PolicyID        string           `json:"policy_id"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	Action          string           `json:"action"` // BLOCK, LOG, CHALLENGE
	MatchConditions []MatchCondition `json:"match_conditions"`
	Status          string           `json:"status"`
}

// ListCustomRules returns custom regex/condition firewall rules
func ListCustomRules(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rules := []CustomRule{
			{
				ID:          "cr_001",
				PolicyID:    "pol_default",
				Name:        "Block SQLi in Header",
				Description: "Blocks SQL injection attempts in the User-Agent header",
				Action:      "BLOCK",
				MatchConditions: []MatchCondition{
					{Variable: "HEADER:User-Agent", Operator: "REGEX_MATCH", Value: "(?i)(union.*select|select.*from)"},
				},
				Status: "active",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rules)
	}
}

// CreateCustomRule adds a new custom firewall rule
func CreateCustomRule(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body CustomRule
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}

		if body.Name == "" || body.Action == "" || len(body.MatchConditions) == 0 {
			http.Error(w, `{"error":"name, action, and match_conditions are required"}`, http.StatusBadRequest)
			return
		}

		newID, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": newID})
	}
}
