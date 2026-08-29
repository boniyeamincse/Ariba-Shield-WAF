package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListSecurityPolicies returns all security policies.
func ListSecurityPolicies(st *store.Store) http.HandlerFunc {
	type policy struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		Description     string `json:"description"`
		EnforcementMode string `json:"enforcement_mode"`
		ApplicationID   string `json:"application_id,omitempty"`
		Version         int64  `json:"version"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, description, enforcement_mode, COALESCE(application_id, ''), version
			 FROM security_policies ORDER BY name`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var policies []policy
		for rows.Next() {
			var p policy
			if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.EnforcementMode, &p.ApplicationID, &p.Version); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			policies = append(policies, p)
		}
		if policies == nil {
			policies = []policy{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(policies)
	}
}

// CreateSecurityPolicy creates a policy (defaults to transparent mode in 0.1).
func CreateSecurityPolicy(st *store.Store) http.HandlerFunc {
	type policyCreate struct {
		Name            string `json:"name"`
		Description     string `json:"description"`
		EnforcementMode string `json:"enforcement_mode"`
		ApplicationID   string `json:"application_id,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body policyCreate
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Name == "" {
			http.Error(w, `{"error":"name is required"}`, http.StatusBadRequest)
			return
		}
		if body.EnforcementMode == "" {
			body.EnforcementMode = "transparent"
		}
		switch body.EnforcementMode {
		case "transparent", "alarm", "blocking":
		default:
			http.Error(w, `{"error":"invalid enforcement_mode"}`, http.StatusBadRequest)
			return
		}

		policyID, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		// Single-org for 0.1; default org is used.
		orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
		userID := "01ARZ3NDEKTSV4RRFFQ69G5FAW"

		appID := &body.ApplicationID
		if body.ApplicationID == "" {
			appID = nil
		}

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO security_policies (id, organization_id, application_id, name, description, enforcement_mode, created_by)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			policyID, orgID, appID, body.Name, body.Description, body.EnforcementMode, userID); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": policyID})
	}
}