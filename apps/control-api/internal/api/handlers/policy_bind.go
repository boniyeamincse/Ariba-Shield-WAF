package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// BindPolicy binds a security policy to an application.
func BindPolicy(st *store.Store) http.HandlerFunc {
	type bindRequest struct {
		PolicyID      string `json:"policy_id"`
		ApplicationID string `json:"application_id"`
		EnforcementMode string `json:"enforcement_mode"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body bindRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.PolicyID == "" || body.ApplicationID == "" {
			http.Error(w, `{"error":"policy_id and application_id are required"}`, http.StatusBadRequest)
			return
		}

		// Verify the policy and application exist.
		var exists bool
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT EXISTS(SELECT 1 FROM security_policies WHERE id = $1)`, body.PolicyID).Scan(&exists); err != nil || !exists {
			http.Error(w, `{"error":"policy not found"}`, http.StatusNotFound)
			return
		}
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT EXISTS(SELECT 1 FROM applications WHERE id = $1)`, body.ApplicationID).Scan(&exists); err != nil || !exists {
			http.Error(w, `{"error":"application not found"}`, http.StatusNotFound)
			return
		}

		mode := body.EnforcementMode
		if mode == "" {
			mode = "transparent"
		}
		switch mode {
		case "transparent", "alarm", "blocking":
		default:
			http.Error(w, `{"error":"invalid enforcement_mode"}`, http.StatusBadRequest)
			return
		}

		// Update the policy's application binding and enforcement mode.
		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE security_policies SET application_id = $1, enforcement_mode = $2, version = version + 1, updated_at = now() WHERE id = $3`,
			body.ApplicationID, mode, body.PolicyID); err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "bound"})
	}
}