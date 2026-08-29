package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

type deployment struct {
	ID              string   `json:"id"`
	PolicyVersionID string   `json:"policy_version_id,omitempty"`
	TargetGateways  []string `json:"target_gateways"`
	Status          string   `json:"status"`
	Error           string   `json:"error,omitempty"`
}

// ListDeployments returns the deployment history.
func ListDeployments(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, COALESCE(policy_version_id,''), target_gateways, status, COALESCE(error,'')
			 FROM deployments ORDER BY created_at DESC LIMIT 50`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		deps := []deployment{}
		for rows.Next() {
			var d deployment
			if err := rows.Scan(&d.ID, &d.PolicyVersionID, &d.TargetGateways, &d.Status, &d.Error); err != nil {
				continue
			}
			deps = append(deps, d)
		}
		if deps == nil {
			deps = []deployment{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deps)
	}
}

// SyncDeployment triggers a sync of a policy version to target gateways.
func SyncDeployment(st *store.Store) http.HandlerFunc {
	type sync struct {
		PolicyVersionID string   `json:"policy_version_id"`
		TargetGateways  []string `json:"target_gateways"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body sync
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.PolicyVersionID == "" || len(body.TargetGateways) == 0 {
			http.Error(w, `{"error":"policy_version_id and target_gateways required"}`, http.StatusBadRequest)
			return
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}
		orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO deployments (id, organization_id, policy_version_id, target_gateways, status)
			 VALUES ($1, $2, $3, $4, 'active')`,
			id, orgID, body.PolicyVersionID, body.TargetGateways); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "active"})
	}
}