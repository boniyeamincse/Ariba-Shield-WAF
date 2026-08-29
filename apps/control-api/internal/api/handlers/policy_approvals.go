package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListPolicyApprovals returns pending + historical approval requests.
func ListPolicyApprovals(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status") // optional filter: pending|approved|rejected
		query := `SELECT pa.id, pa.policy_id, pa.policy_version_id, pa.status,
		                 COALESCE(pa.requested_by,''), COALESCE(pa.approved_by,''),
		                 pa.reviewer_notes, pa.created_at::text
		          FROM policy_approvals pa`
		args := []any{}
		if status != "" {
			args = append(args, status)
			query += " WHERE pa.status = $1"
		}
		query += " ORDER BY pa.created_at DESC LIMIT 100"

		rows, err := st.Pool.Query(r.Context(), query, args...)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		approvals := []map[string]any{}
		for rows.Next() {
			var id, pid, pvid, st, reqBy, appBy, notes, created string
			if err := rows.Scan(&id, &pid, &pvid, &st, &reqBy, &appBy, &notes, &created); err == nil {
				approvals = append(approvals, map[string]any{
					"id": id, "policy_id": pid, "policy_version_id": pvid, "status": st,
					"requested_by": reqBy, "approved_by": appBy, "reviewer_notes": notes,
					"created_at": created,
				})
			}
		}
		if approvals == nil {
			approvals = []map[string]any{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(approvals)
	}
}

// CreatePolicyApproval requests approval for a policy change (four-eyes).
func CreatePolicyApproval(st *store.Store) http.HandlerFunc {
	type create struct {
		PolicyID        string `json:"policy_id"`
		PolicyVersionID string `json:"policy_version_id"`
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

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}
		userID := "01ARZ3NDEKTSV4RRFFQ69G5FAW"
		orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"

		// Set policy lifecycle to APPROVAL_REQUIRED.
		_, _ = st.Pool.Exec(r.Context(),
			`UPDATE security_policies SET lifecycle_status = 'approval_required', updated_at = now() WHERE id = $1`, body.PolicyID)

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO policy_approvals (id, organization_id, policy_id, policy_version_id, status, requested_by)
			 VALUES ($1, $2, $3, $4, 'pending', $5)`,
			id, orgID, body.PolicyID, nullIfEmpty(body.PolicyVersionID), userID); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "pending"})
	}
}

// GetPolicyApproval returns a single approval request.
func GetPolicyApproval(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var b []byte
		err := st.Pool.QueryRow(r.Context(),
			`SELECT row_to_json(t) FROM (
			   SELECT id, policy_id, policy_version_id, status, requested_by, approved_by,
			          reviewer_notes, created_at, reviewed_at
			   FROM policy_approvals WHERE id = $1) t`, id).Scan(&b)
		if err != nil {
			http.Error(w, `{"error":"approval not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

// ApprovePolicyApproval approves a pending request (sets policy lifecycle
// to APPROVED; the version can then be promoted to CANARY/ACTIVE).
func ApprovePolicyApproval(st *store.Store) http.HandlerFunc {
	type decision struct {
		Notes string `json:"notes"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body decision
		_ = json.NewDecoder(r.Body).Decode(&body)

		var policyID, versionID, status string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT policy_id, COALESCE(policy_version_id,''), status FROM policy_approvals WHERE id = $1`, id).
			Scan(&policyID, &versionID, &status); err != nil {
			http.Error(w, `{"error":"approval not found"}`, http.StatusNotFound)
			return
		}
		if status != "pending" {
			http.Error(w, `{"error":"approval already decided"}`, http.StatusConflict)
			return
		}

		reviewer := "01ARZ3NDEKTSV4RRFFQ69G5FAW"
		now := time.Now().UTC()
		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE policy_approvals SET status = 'approved', approved_by = $1, reviewer_notes = $2, reviewed_at = $3 WHERE id = $4`,
			reviewer, body.Notes, now, id); err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}

		// Mark the version approved.
		if versionID != "" {
			_, _ = st.Pool.Exec(r.Context(),
				`UPDATE policy_versions SET status = 'approved' WHERE id = $1`, versionID)
		}
		_, _ = st.Pool.Exec(r.Context(),
			`UPDATE security_policies SET lifecycle_status = 'approved', updated_at = now() WHERE id = $1`, policyID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "approved", "policy_id": policyID})
	}
}

// RejectPolicyApproval rejects a pending request (policy stays DRAFT).
func RejectPolicyApproval(st *store.Store) http.HandlerFunc {
	type decision struct {
		Notes string `json:"notes"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body decision
		_ = json.NewDecoder(r.Body).Decode(&body)

		var policyID, status string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT policy_id, status FROM policy_approvals WHERE id = $1`, id).Scan(&policyID, &status); err != nil {
			http.Error(w, `{"error":"approval not found"}`, http.StatusNotFound)
			return
		}
		if status != "pending" {
			http.Error(w, `{"error":"approval already decided"}`, http.StatusConflict)
			return
		}

		reviewer := "01ARZ3NDEKTSV4RRFFQ69G5FAW"
		now := time.Now().UTC()
		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE policy_approvals SET status = 'rejected', approved_by = $1, reviewer_notes = $2, reviewed_at = $3 WHERE id = $4`,
			reviewer, body.Notes, now, id); err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		// Revert lifecycle to DRAFT on rejection.
		_, _ = st.Pool.Exec(r.Context(),
			`UPDATE security_policies SET lifecycle_status = 'draft', updated_at = now() WHERE id = $1`, policyID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "rejected", "policy_id": policyID})
	}
}