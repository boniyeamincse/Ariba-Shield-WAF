package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// CreatePolicyVersion creates a new immutable version of a security policy.
func CreatePolicyVersion(st *store.Store) http.HandlerFunc {
	type create struct {
		PolicyID  string          `json:"policy_id"`
		Document  json.RawMessage `json:"document"`
		BundleHash string         `json:"bundle_hash"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body create
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.PolicyID == "" {
			http.Error(w, `{"error":"policy_id required"}`, http.StatusBadRequest)
			return
		}

		// Get the next version number.
		var nextVersion int64
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT COALESCE(MAX(version), 0) + 1 FROM policy_versions WHERE policy_id = $1`, body.PolicyID).Scan(&nextVersion); err != nil {
			http.Error(w, `{"error":"version query failed"}`, http.StatusInternalServerError)
			return
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}
		userID := "01ARZ3NDEKTSV4RRFFQ69G5FAW"

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO policy_versions (id, policy_id, version, document, bundle_hash, created_by)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			id, body.PolicyID, nextVersion, body.Document, body.BundleHash, userID); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{"id": id, "version": nextVersion})
	}
}

// ActivatePolicyVersion sets a policy version to active and demotes prior ones.
func ActivatePolicyVersion(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		versionID := r.PathValue("id")

		// Get the policy ID from the version.
		var policyID string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT policy_id FROM policy_versions WHERE id = $1`, versionID).Scan(&policyID); err != nil {
			http.Error(w, `{"error":"version not found"}`, http.StatusNotFound)
			return
		}

		// Demote any currently active version for this policy.
		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE policy_versions SET status = 'superseded' WHERE policy_id = $1 AND status = 'active'`, policyID); err != nil {
			http.Error(w, `{"error":"demote failed"}`, http.StatusInternalServerError)
			return
		}

		// Activate the requested version.
		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE policy_versions SET status = 'active' WHERE id = $1`, versionID); err != nil {
			http.Error(w, `{"error":"activate failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "active"})
	}
}

// RollbackPolicyVersion rolls back to the last active version.
func RollbackPolicyVersion(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		versionID := r.PathValue("id")

		// Set this version to rolled_back.
		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE policy_versions SET status = 'rolled_back' WHERE id = $1`, versionID); err != nil {
			http.Error(w, `{"error":"rollback failed"}`, http.StatusInternalServerError)
			return
		}

		// Find the previous active version and reactivate it.
		// The version left as the most recent non-rolled_back superseded version.
		var prevID string
		// Reactivate the most recent superseded version (before this one).
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT id FROM policy_versions WHERE policy_id = (
				SELECT policy_id FROM policy_versions WHERE id = $1
			) AND status = 'superseded' ORDER BY version DESC LIMIT 1`, versionID).Scan(&prevID); err != nil {
			http.Error(w, `{"error":"no previous version to roll back to"}`, http.StatusNotFound)
			return
		}

		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE policy_versions SET status = 'active' WHERE id = $1`, prevID); err != nil {
			http.Error(w, `{"error":"reactivate failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "rolled_back", "active_version_id": prevID})
	}
}