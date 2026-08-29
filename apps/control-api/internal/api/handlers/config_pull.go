package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// ConfigPull returns the latest active config version for the gateway's
// organization (ADR-004 D3). Returns 304 Not Modified if the gateway's
// If-None-Match matches the current bundle hash.
func ConfigPull(st *store.Store) http.HandlerFunc {
	type configResponse struct {
		ConfigID   string          `json:"config_id"`
		BundleHash string          `json:"bundle_hash"`
		Document   json.RawMessage `json:"document"`
		CreatedAt  string          `json:"created_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		gatewayID := r.PathValue("id")

		// Get the gateway's organization.
		var orgID string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT organization_id FROM gateways WHERE id = $1`, gatewayID).Scan(&orgID); err != nil {
			http.Error(w, `{"error":"gateway not found"}`, http.StatusNotFound)
			return
		}

		// Look up the latest active config version for this org.
		var configID, bundleHash, createdAt string
		var document json.RawMessage
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, bundle_hash, blob, created_at::text
			 FROM config_versions
			 WHERE organization_id = $1 AND status = 'active'
			 ORDER BY created_at DESC LIMIT 1`,
			orgID).Scan(&configID, &bundleHash, &document, &createdAt)
		if err != nil {
			http.Error(w, `{"error":"no active config"}`, http.StatusNotFound)
			return
		}

		// 304 if the gateway already has this hash.
		if r.Header.Get("If-None-Match") == `"`+bundleHash+`"` {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("ETag", `"`+bundleHash+`"`)
		json.NewEncoder(w).Encode(configResponse{
			ConfigID:   configID,
			BundleHash: bundleHash,
			Document:   document,
			CreatedAt:  createdAt,
		})
	}
}