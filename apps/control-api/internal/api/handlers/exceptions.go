package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

type ExceptionCondition struct {
	Variable string `json:"variable"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type ExceptionRule struct {
	ID             string               `json:"id"`
	TargetRuleID   string               `json:"target_rule_id"`
	ApplicationID  string               `json:"application_id"`
	MatchConditions []ExceptionCondition `json:"match_conditions"`
	Reason         string               `json:"reason"`
	ExpiryDays     int                  `json:"expiry_days"`
	Status         string               `json:"status"`
}

// ListExceptions returns all configured WAF exceptions
func ListExceptions(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Mock response for now, to be replaced with PostgreSQL query
		exceptions := []ExceptionRule{
			{
				ID:            "exc_001",
				TargetRuleID:  "ARIBA-SQLI-001",
				ApplicationID: "app_internal_hr",
				MatchConditions: []ExceptionCondition{
					{Variable: "URI_PATH", Operator: "EQUALS", Value: "/api/reports/query"},
				},
				Reason:     "Internal HR reporting payload contains SQL-like syntax",
				ExpiryDays: 90,
				Status:     "active",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exceptions)
	}
}

// CreateException registers a new exclusion rule to mitigate false positives
func CreateException(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body ExceptionRule
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}

		if body.TargetRuleID == "" || body.ApplicationID == "" {
			http.Error(w, `{"error":"target_rule_id and application_id are required"}`, http.StatusBadRequest)
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
