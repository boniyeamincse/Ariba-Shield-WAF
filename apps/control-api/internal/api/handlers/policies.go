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
		ParentPolicyID  string `json:"parent_policy_id,omitempty"`
		CanaryPercent   int    `json:"canary_percent"`
		ScheduledAt     string `json:"scheduled_at,omitempty"`
		ScheduledAction string `json:"scheduled_action,omitempty"`
		Version         int64  `json:"version"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, description, enforcement_mode, COALESCE(application_id, ''),
			        COALESCE(parent_policy_id, ''), COALESCE(canary_percent, 0),
			        COALESCE(scheduled_at::text, ''), COALESCE(scheduled_action, ''), version
			 FROM security_policies ORDER BY name`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var policies []policy
		for rows.Next() {
			var p policy
			if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.EnforcementMode, &p.ApplicationID,
				&p.ParentPolicyID, &p.CanaryPercent, &p.ScheduledAt, &p.ScheduledAction, &p.Version); err != nil {
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
		ParentPolicyID  string `json:"parent_policy_id,omitempty"`
		CanaryPercent   int    `json:"canary_percent"`
		ScheduledAt     string `json:"scheduled_at,omitempty"`
		ScheduledAction string `json:"scheduled_action,omitempty"`
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
		if body.CanaryPercent < 0 || body.CanaryPercent > 100 {
			http.Error(w, `{"error":"canary_percent must be 0-100"}`, http.StatusBadRequest)
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
		parentID := &body.ParentPolicyID
		if body.ParentPolicyID == "" {
			parentID = nil
		}

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO security_policies (id, organization_id, application_id, name, description, enforcement_mode, created_by,
			                                parent_policy_id, canary_percent, scheduled_at, scheduled_action)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			policyID, orgID, appID, body.Name, body.Description, body.EnforcementMode, userID,
			parentID, body.CanaryPercent, nullIfEmpty(body.ScheduledAt), nullIfEmpty(body.ScheduledAction)); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": policyID})
	}
}
// UpdateSecurityPolicy updates a security policy (partial update).
func UpdateSecurityPolicy(st *store.Store) http.HandlerFunc {
	type update struct {
		Name            *string `json:"name"`
		Description     *string `json:"description"`
		EnforcementMode *string `json:"enforcement_mode"`
		ParentPolicyID  *string `json:"parent_policy_id"`
		CanaryPercent   *int    `json:"canary_percent"`
		ScheduledAt     *string `json:"scheduled_at"`
		ScheduledAction *string `json:"scheduled_action"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body update
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.EnforcementMode != nil {
			switch *body.EnforcementMode {
			case "transparent", "alarm", "blocking":
			default:
				http.Error(w, `{"error":"invalid enforcement_mode"}`, http.StatusBadRequest)
				return
			}
		}
		if body.CanaryPercent != nil && (*body.CanaryPercent < 0 || *body.CanaryPercent > 100) {
			http.Error(w, `{"error":"canary_percent must be 0-100"}`, http.StatusBadRequest)
			return
		}

		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE security_policies SET
			   name = COALESCE($1, name),
			   description = COALESCE($2, description),
			   enforcement_mode = COALESCE($3, enforcement_mode),
			   parent_policy_id = COALESCE($4, parent_policy_id),
			   canary_percent = COALESCE($5, canary_percent),
			   scheduled_at = COALESCE($6, scheduled_at),
			   scheduled_action = COALESCE($7, scheduled_action),
			   version = version + 1, updated_at = now()
			 WHERE id = $8`,
			nullableString(body.Name), nullableString(body.Description),
			nullableString(body.EnforcementMode), nullableString(body.ParentPolicyID),
			body.CanaryPercent, nullableString(body.ScheduledAt), nullableString(body.ScheduledAction), id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"policy not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// DeleteSecurityPolicy deletes a security policy.
func DeleteSecurityPolicy(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM security_policies WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"policy not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}
