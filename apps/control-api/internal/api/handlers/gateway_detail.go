package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

// GetGateway returns a single gateway with its latest applied state.
func GetGateway(st *store.Store) http.HandlerFunc {
	type gateway struct {
		ID          string   `json:"id"`
		Hostname    string   `json:"hostname"`
		IP          string   `json:"ip"`
		Version     string   `json:"version"`
		Capabilities []string `json:"capabilities"`
		Status      string   `json:"status"`
		LastSeenAt  *string  `json:"last_seen_at,omitempty"`
		AppliedHash string   `json:"applied_hash,omitempty"`
		LastError   string   `json:"last_error,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var g gateway
		var lastSeen *time.Time
		var appliedHash, lastError *string
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, hostname, COALESCE(ip,''), COALESCE(version,''), capabilities,
			        status, last_seen_at,
			        (SELECT h.applied_hash FROM gateway_heartbeats h
			         WHERE h.gateway_id = g.id ORDER BY h.created_at DESC LIMIT 1),
			        (SELECT h.last_error FROM gateway_heartbeats h
			         WHERE h.gateway_id = g.id ORDER BY h.created_at DESC LIMIT 1)
			 FROM gateways g WHERE id = $1`, id).
			Scan(&g.ID, &g.Hostname, &g.IP, &g.Version, &g.Capabilities, &g.Status,
				&lastSeen, &appliedHash, &lastError)
		if err != nil {
			http.Error(w, `{"error":"gateway not found"}`, http.StatusNotFound)
			return
		}
		if lastSeen != nil {
			s := lastSeen.Format(time.RFC3339)
			g.LastSeenAt = &s
		}
		if appliedHash != nil {
			g.AppliedHash = *appliedHash
		}
		if lastError != nil {
			g.LastError = *lastError
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(g)
	}
}

// UpdateGateway updates a gateway's metadata (hostname, version, status).
func UpdateGateway(st *store.Store) http.HandlerFunc {
	type update struct {
		Hostname     *string  `json:"hostname"`
		IP           *string  `json:"ip"`
		Version      *string  `json:"version"`
		Capabilities *[]string `json:"capabilities"`
		Status       *string  `json:"status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body update
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}

		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE gateways SET
			   hostname = COALESCE($1, hostname),
			   ip = COALESCE($2, ip),
			   version = COALESCE($3, version),
			   capabilities = COALESCE($4, capabilities),
			   status = COALESCE($5, status),
			   updated_at = now()
			 WHERE id = $6`,
			nullableString(body.Hostname), nullableString(body.IP), nullableString(body.Version),
			body.Capabilities, nullableString(body.Status), id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"gateway not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// DeleteGateway removes a gateway (and its heartbeats via cascade).
func DeleteGateway(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM gateways WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"gateway not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}

// GatewayStatus returns the live status + latest heartbeat of a gateway.
func GatewayStatus(st *store.Store) http.HandlerFunc {
	type status struct {
		GatewayID   string    `json:"gateway_id"`
		Status      string    `json:"status"`
		AppliedHash string    `json:"applied_hash"`
		Version     string    `json:"version"`
		Health      json.RawMessage `json:"health"`
		LastError   string    `json:"last_error"`
		LastSeenAt  string    `json:"last_seen_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var s status
		var health json.RawMessage
		var lastError *string
		var lastSeen time.Time
		err := st.Pool.QueryRow(r.Context(),
			`SELECT gateway_id, status, COALESCE(applied_hash,''), COALESCE(version,''),
			        COALESCE(health,'{}'), COALESCE(last_error,''), created_at
			 FROM gateway_heartbeats
			 WHERE gateway_id = $1 ORDER BY created_at DESC LIMIT 1`, id).
			Scan(&s.GatewayID, &s.Status, &s.AppliedHash, &s.Version, &health, &lastError, &lastSeen)
		if err != nil {
			http.Error(w, `{"error":"no heartbeat yet"}`, http.StatusNotFound)
			return
		}
		s.GatewayID = id
		s.Health = health
		if lastError != nil {
			s.LastError = *lastError
		}
		s.LastSeenAt = lastSeen.Format(time.RFC3339)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s)
	}
}

// GatewayConfig returns the config currently applied to a gateway (the
// immutable signed bundle hash + metadata; the bundle itself is served by
// /config/current per ADR-004 D3).
func GatewayConfig(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var appliedHash string
		var lastError *string
		var status string
		err := st.Pool.QueryRow(r.Context(),
			`SELECT COALESCE((SELECT h.applied_hash FROM gateway_heartbeats h
			                  WHERE h.gateway_id = $1 ORDER BY h.created_at DESC LIMIT 1), ''),
			        COALESCE((SELECT h.last_error FROM gateway_heartbeats h
			                  WHERE h.gateway_id = $1 ORDER BY h.created_at DESC LIMIT 1), ''),
			        status FROM gateways WHERE id = $1`, id).
			Scan(&appliedHash, &lastError, &status)
		if err != nil {
			http.Error(w, `{"error":"gateway not found"}`, http.StatusNotFound)
			return
		}

		// Map applied hash to its config version.
		var versionID, created string
		_ = st.Pool.QueryRow(r.Context(),
			`SELECT id, created_at::text FROM config_versions WHERE bundle_hash = $1`, appliedHash).
			Scan(&versionID, &created)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"gateway_id":     id,
			"applied_hash":   appliedHash,
			"config_version": versionID,
			"applied_at":     created,
			"last_error":     lastError,
			"gateway_status": status,
		})
	}
}

// ApplyGatewayConfig marks a gateway as applying a config version (records a
// deployment; the gateway pulls + verifies the signed bundle via /config/current,
// honoring last-known-good semantics in apply-config.sh).
func ApplyGatewayConfig(st *store.Store) http.HandlerFunc {
	type apply struct {
		ConfigVersionID string `json:"config_version_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		gatewayID := r.PathValue("id")
		var body apply
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ConfigVersionID == "" {
			http.Error(w, `{"error":"config_version_id required"}`, http.StatusBadRequest)
			return
		}

		// Verify the config version exists.
		var bundleHash string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT bundle_hash FROM config_versions WHERE id = $1`, body.ConfigVersionID).Scan(&bundleHash); err != nil {
			http.Error(w, `{"error":"config version not found"}`, http.StatusNotFound)
			return
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}
		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO config_deployments (id, config_version_id, target_gateway, status, applied_hash)
			 VALUES ($1, $2, $3, 'pending', $4)`,
			id, body.ConfigVersionID, gatewayID, bundleHash); err != nil {
			http.Error(w, `{"error":"deployment failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"deployment_id": id, "gateway_id": gatewayID, "status": "pending",
		})
	}
}

// RollbackGatewayConfig rolls a gateway back to the previous active config
// (re-points the applied config and records a deployment).
func RollbackGatewayConfig(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gatewayID := r.PathValue("id")

		// Find the most recent ACTIVE config version (the known-good one).
		var versionID, bundleHash string
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, bundle_hash FROM config_versions
			 WHERE status = 'active' ORDER BY created_at DESC LIMIT 1`).Scan(&versionID, &bundleHash)
		if err != nil {
			http.Error(w, `{"error":"no active config to roll back to"}`, http.StatusNotFound)
			return
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}
		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO config_deployments (id, config_version_id, target_gateway, status, applied_hash)
			 VALUES ($1, $2, $3, 'activated', $4)`,
			id, versionID, gatewayID, bundleHash); err != nil {
			http.Error(w, `{"error":"rollback failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"deployment_id": id, "gateway_id": gatewayID,
			"rollback_to": versionID, "bundle_hash": bundleHash,
		})
	}
}

// GatewayMetrics returns gateway-side counters (from heartbeats + health).
func GatewayMetrics(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var total, blocked, lastSeen string
		err := st.Pool.QueryRow(r.Context(),
			`SELECT COUNT(*), COUNT(*) FILTER (WHERE decision_action = 'block'),
			        COALESCE(MAX(created_at)::text, '')
			 FROM security_events WHERE gateway_id = $1`, id).Scan(&total, &blocked, &lastSeen)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"gateway_id":   id,
			"events_total": total,
			"events_blocked": blocked,
			"last_event_at": lastSeen,
		})
	}
}