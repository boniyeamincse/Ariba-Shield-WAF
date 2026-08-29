package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

type ManagedRuleSet struct {
	ID                     string `json:"id"`
	PolicyID               string `json:"policy_id"`
	RuleSet                string `json:"rule_set"`
	ParanoiaLevel          int    `json:"paranoia_level"`
	EnforcementMode        string `json:"enforcement_mode"`
	AnomalyScoreThreshold  int    `json:"anomaly_score_threshold"`
}

// ListManagedRules returns the configured managed rule sets (like OWASP CRS) for a policy
func ListManagedRules(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Mock response
		rules := []ManagedRuleSet{
			{
				ID:                    "mrs_001",
				PolicyID:              "pol_default",
				RuleSet:               "OWASP_CRS_V3",
				ParanoiaLevel:         1,
				EnforcementMode:       "DETECTION",
				AnomalyScoreThreshold: 5,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rules)
	}
}

// ConfigureManagedRules enables or updates a managed rule set for a security policy
func ConfigureManagedRules(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body ManagedRuleSet
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}

		if body.RuleSet == "" || body.EnforcementMode == "" {
			http.Error(w, `{"error":"rule_set and enforcement_mode are required"}`, http.StatusBadRequest)
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
