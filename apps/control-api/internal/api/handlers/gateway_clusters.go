package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// GetGatewayCluster returns a single cluster by id.
func GetGatewayCluster(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var cluster struct {
			ID         string   `json:"id"`
			Name       string   `json:"name"`
			Site       string   `json:"site"`
			GatewayIDs []string `json:"gateway_ids"`
			HAEnabled  bool     `json:"ha_enabled"`
			Status     string   `json:"status"`
		}
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, name, COALESCE(site,''), gateway_ids, ha_enabled, status
			 FROM clusters WHERE id = $1`, id).
			Scan(&cluster.ID, &cluster.Name, &cluster.Site, &cluster.GatewayIDs, &cluster.HAEnabled, &cluster.Status)
		if err != nil {
			http.Error(w, `{"error":"cluster not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cluster)
	}
}

// ClusterGateways returns the gateways that belong to a cluster.
func ClusterGateways(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var gatewayIDs []string
		var clusterName string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT name, gateway_ids FROM clusters WHERE id = $1`, id).Scan(&clusterName, &gatewayIDs); err != nil {
			http.Error(w, `{"error":"cluster not found"}`, http.StatusNotFound)
			return
		}

		type gw struct {
			ID         string `json:"id"`
			Hostname   string `json:"hostname"`
			Status     string `json:"status"`
			LastSeenAt *string `json:"last_seen_at,omitempty"`
		}
		gateways := []gw{}
		if len(gatewayIDs) > 0 {
			rows, err := st.Pool.Query(r.Context(),
				`SELECT id, hostname, status, last_seen_at::text FROM gateways WHERE id = ANY($1)`, gatewayIDs)
			if err == nil {
				for rows.Next() {
					var g gw
					var lastSeen *string
					if err := rows.Scan(&g.ID, &g.Hostname, &g.Status, &lastSeen); err == nil {
						g.LastSeenAt = lastSeen
						gateways = append(gateways, g)
					}
				}
				rows.Close()
			}
		}
		if gateways == nil {
			gateways = []gw{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"cluster_id": id, "cluster_name": clusterName, "gateways": gateways,
		})
	}
}

// DeployClusterConfig marks a config version for deployment to all cluster
// gateways (each gateway pulls + verifies its signed bundle via /config/current,
// honoring last-known-good semantics).
func DeployClusterConfig(st *store.Store) http.HandlerFunc {
	type deploy struct {
		ConfigVersionID string `json:"config_version_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body deploy
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ConfigVersionID == "" {
			http.Error(w, `{"error":"config_version_id required"}`, http.StatusBadRequest)
			return
		}

		var gatewayIDs []string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT gateway_ids FROM clusters WHERE id = $1`, id).Scan(&gatewayIDs); err != nil {
			http.Error(w, `{"error":"cluster not found"}`, http.StatusNotFound)
			return
		}

		var bundleHash string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT bundle_hash FROM config_versions WHERE id = $1`, body.ConfigVersionID).Scan(&bundleHash); err != nil {
			http.Error(w, `{"error":"config version not found"}`, http.StatusNotFound)
			return
		}

		created := []string{}
		for _, gwID := range gatewayIDs {
			deploymentID, err := st.NewID()
			if err != nil {
				continue
			}
			if _, err := st.Pool.Exec(r.Context(),
				`INSERT INTO config_deployments (id, config_version_id, target_gateway, status, applied_hash)
				 VALUES ($1, $2, $3, 'pending', $4)`,
				deploymentID, body.ConfigVersionID, gwID, bundleHash); err == nil {
				created = append(created, deploymentID)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]any{
			"cluster_id": id, "config_version_id": body.ConfigVersionID,
			"gateways": gatewayIDs, "deployments_created": created,
		})
	}
}

// RollbackClusterConfig rolls back all cluster gateways to the last active
// config version.
func RollbackClusterConfig(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var gatewayIDs []string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT gateway_ids FROM clusters WHERE id = $1`, id).Scan(&gatewayIDs); err != nil {
			http.Error(w, `{"error":"cluster not found"}`, http.StatusNotFound)
			return
		}

		var versionID, bundleHash string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT id, bundle_hash FROM config_versions WHERE status = 'active' ORDER BY created_at DESC LIMIT 1`).
			Scan(&versionID, &bundleHash); err != nil {
			http.Error(w, `{"error":"no active config to roll back to"}`, http.StatusNotFound)
			return
		}

		for _, gwID := range gatewayIDs {
			deploymentID, _ := st.NewID()
			_, _ = st.Pool.Exec(r.Context(),
				`INSERT INTO config_deployments (id, config_version_id, target_gateway, status, applied_hash)
				 VALUES ($1, $2, $3, 'activated', $4)`,
				deploymentID, versionID, gwID, bundleHash)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"cluster_id": id, "rollback_to": versionID, "bundle_hash": bundleHash, "gateways": gatewayIDs,
		})
	}
}

// ClusterHealth returns the aggregate health of all gateways in a cluster.
func ClusterHealth(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var gatewayIDs []string
		var clusterName string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT name, gateway_ids FROM clusters WHERE id = $1`, id).Scan(&clusterName, &gatewayIDs); err != nil {
			http.Error(w, `{"error":"cluster not found"}`, http.StatusNotFound)
			return
		}

		type gstate struct {
			ID       string `json:"id"`
			Hostname string `json:"hostname"`
			Status   string `json:"status"`
		}
		states := []gstate{}
		up := 0
		if len(gatewayIDs) > 0 {
			rows, err := st.Pool.Query(r.Context(),
				`SELECT id, hostname, status FROM gateways WHERE id = ANY($1)`, gatewayIDs)
			if err == nil {
				for rows.Next() {
					var g gstate
					if err := rows.Scan(&g.ID, &g.Hostname, &g.Status); err == nil {
						if g.Status == "active" || g.Status == "starting" {
							up++
						}
						states = append(states, g)
					}
				}
				rows.Close()
			}
		}
		total := len(states)
		overall := "degraded"
		if total > 0 && up == total {
			overall = "healthy"
		} else if up == 0 {
			overall = "down"
		} else if total == 0 {
			overall = "no_gateways"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"cluster_id": id, "cluster_name": clusterName, "overall": overall,
			"healthy": up, "total": total, "gateways": states,
		})
	}
}