package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

// GetException returns a single exception.
func GetException(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var b []byte
		err := st.Pool.QueryRow(r.Context(),
			`SELECT row_to_json(t) FROM (
			   SELECT id, policy_id, application_id, rule_id, url_pattern, parameter,
			          reason, expires_at, status, created_by, created_at
			   FROM exceptions WHERE id = $1) t`, id).Scan(&b)
		if err != nil {
			http.Error(w, `{"error":"exception not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

// UpdateException updates an exception (partial).
func UpdateException(st *store.Store) http.HandlerFunc {
	type update struct {
		RuleID     *string `json:"rule_id"`
		URLPattern *string `json:"url_pattern"`
		Parameter  *string `json:"parameter"`
		Reason     *string `json:"reason"`
		ExpiresAt  *string `json:"expires_at"`
		Status     *string `json:"status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body update
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}

		var expiresAt any
		if body.ExpiresAt != nil && *body.ExpiresAt != "" {
			if t, err := time.Parse(time.RFC3339, *body.ExpiresAt); err == nil {
				expiresAt = t
			}
		}

		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE exceptions SET
			   rule_id = COALESCE($1, rule_id),
			   url_pattern = COALESCE($2, url_pattern),
			   parameter = COALESCE($3, parameter),
			   reason = COALESCE($4, reason),
			   expires_at = COALESCE($5, expires_at),
			   status = COALESCE($6, status)
			 WHERE id = $7`,
			nullableString(body.RuleID), nullableString(body.URLPattern),
			nullableString(body.Parameter), nullableString(body.Reason),
			expiresAt, nullableString(body.Status), id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"exception not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// DeleteException removes an exception.
func DeleteException(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM exceptions WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"exception not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}

// ApproveException marks an exception approved (status = approved).
func ApproveException(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE exceptions SET status = 'approved' WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"exception not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "approved"})
	}
}

// ExpireException expires an exception immediately (status = expired).
func ExpireException(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE exceptions SET status = 'expired', expires_at = now() WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"exception not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "expired"})
	}
}