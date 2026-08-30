package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

// RegisterGateway registers a new gateway. The ID is always provisioned by the
// server (ignoring any client-supplied gateway_id) so a gateway cannot claim
// another gateway's ID.
func RegisterGateway(st *store.Store) http.HandlerFunc {
	type register struct {
		Hostname     string   `json:"hostname"`
		IP           string   `json:"ip"`
		Version      string   `json:"version"`
		Capabilities []string `json:"capabilities"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body register
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Hostname == "" {
			http.Error(w, `{"error":"hostname is required"}`, http.StatusBadRequest)
			return
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO gateways (id, organization_id, hostname, ip, version, capabilities, status)
			 VALUES ($1, $2, $3, $4, $5, $6, 'active')`,
			id, orgID, body.Hostname, body.IP, body.Version, body.Capabilities); err != nil {
			http.Error(w, `{"error":"upsert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "active"})
	}
}

// Heartbeat records a gateway heartbeat and updates its last_seen/status.
func Heartbeat(st *store.Store) http.HandlerFunc {
	type heartbeat struct {
		Status      string          `json:"status"`
		AppliedHash string          `json:"applied_hash"`
		Version     string          `json:"version"`
		Health      json.RawMessage `json:"health"`
		LastError   string          `json:"last_error"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		gatewayID := r.PathValue("id")

		var body heartbeat
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Status == "" {
			body.Status = "active"
		}

		hID, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO gateway_heartbeats (id, gateway_id, status, applied_hash, version, health, last_error)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			hID, gatewayID, body.Status, body.AppliedHash, body.Version, body.Health, body.LastError); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		now := time.Now().UTC()
		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE gateways SET status = $2, last_seen_at = $3, version = $4, updated_at = $3 WHERE id = $1`,
			gatewayID, body.Status, now, body.Version); err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

// ListGateways returns gateway fleet status.
func ListGateways(st *store.Store) http.HandlerFunc {
	type gateway struct {
		ID          string     `json:"id"`
		Hostname    string     `json:"hostname"`
		IP          string     `json:"ip"`
		Version     string     `json:"version"`
		Status      string     `json:"status"`
		LastSeenAt  *time.Time `json:"last_seen_at"`
		AppliedHash string     `json:"applied_hash,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT g.id, g.hostname, COALESCE(g.ip, ''), COALESCE(g.version, ''), g.status, g.last_seen_at,
			        (SELECT h.applied_hash FROM gateway_heartbeats h WHERE h.gateway_id = g.id ORDER BY h.created_at DESC LIMIT 1)
			 FROM gateways g ORDER BY g.hostname`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var gateways []gateway
		for rows.Next() {
			var g gateway
			var appliedHash *string
			if err := rows.Scan(&g.ID, &g.Hostname, &g.IP, &g.Version, &g.Status, &g.LastSeenAt, &appliedHash); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			if appliedHash != nil {
				g.AppliedHash = *appliedHash
			}
			gateways = append(gateways, g)
		}
		if gateways == nil {
			gateways = []gateway{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(gateways)
	}
}
