package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

type Deployment struct {
	ID             string `json:"id"`
	PolicyBundleID string `json:"policy_bundle_id"`
	TargetClusters []string `json:"target_clusters"`
	Status         string `json:"status"` // PENDING, SYNCING, SUCCESS, FAILED
	DeployedBy     string `json:"deployed_by"`
	DeployedAt     string `json:"deployed_at"`
}

// ListDeployments returns the history of configuration syncs to the data plane
func ListDeployments(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deps := []Deployment{
			{
				ID:             "dep_001",
				PolicyBundleID: "bundle_v1.0.4",
				TargetClusters: []string{"cluster-us-east", "cluster-eu-west"},
				Status:         "SUCCESS",
				DeployedBy:     "admin@aribashield.local",
				DeployedAt:     "2026-08-28T14:00:00Z",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deps)
	}
}

// SyncDeployment pushes a new declarative configuration to the WAF edge gateways
func SyncDeployment(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body Deployment
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}

		if body.PolicyBundleID == "" || len(body.TargetClusters) == 0 {
			http.Error(w, `{"error":"policy_bundle_id and target_clusters are required"}`, http.StatusBadRequest)
			return
		}

		// In a real application, this triggers an asynchronous Redis/gRPC sync pipeline

		newID, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted) // Return 202 Accepted as it's async
		json.NewEncoder(w).Encode(map[string]string{
			"id":     newID,
			"status": "SYNCING",
		})
	}
}
