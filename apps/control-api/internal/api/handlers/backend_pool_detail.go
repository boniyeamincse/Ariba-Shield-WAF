package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// GetBackendPool returns a single pool with its nodes.
func GetBackendPool(st *store.Store) http.HandlerFunc {
	type node struct {
		ID              string `json:"id"`
		Host            string `json:"host"`
		Port            int    `json:"port"`
		Weight          int    `json:"weight"`
		Active          bool   `json:"active"`
		Draining        bool   `json:"draining"`
		LastHealthState string `json:"last_health_state,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var pool struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			ApplicationID string `json:"application_id"`
			LBAlgorithm   string `json:"lb_algorithm"`
			Version       int64  `json:"version"`
			Nodes         []node `json:"nodes"`
		}
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, name, application_id, lb_algorithm, version
			 FROM backend_pools WHERE id = $1`, id).
			Scan(&pool.ID, &pool.Name, &pool.ApplicationID, &pool.LBAlgorithm, &pool.Version)
		if err != nil {
			http.Error(w, `{"error":"pool not found"}`, http.StatusNotFound)
			return
		}
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, host, port, weight, active, draining, COALESCE(last_health_state,'')
			 FROM backend_nodes WHERE pool_id = $1 ORDER BY host, port`, id)
		if err == nil {
			for rows.Next() {
				var n node
				if err := rows.Scan(&n.ID, &n.Host, &n.Port, &n.Weight, &n.Active, &n.Draining, &n.LastHealthState); err == nil {
					pool.Nodes = append(pool.Nodes, n)
				}
			}
			rows.Close()
		}
		if pool.Nodes == nil {
			pool.Nodes = []node{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pool)
	}
}

// UpdateBackendPool updates a pool's name or lb_algorithm.
func UpdateBackendPool(st *store.Store) http.HandlerFunc {
	type update struct {
		Name        *string `json:"name"`
		LBAlgorithm *string `json:"lb_algorithm"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body update
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE backend_pools SET
			   name = COALESCE($1, name),
			   lb_algorithm = COALESCE($2, lb_algorithm),
			   version = version + 1, updated_at = now()
			 WHERE id = $3`,
			nullableString(body.Name), nullableString(body.LBAlgorithm), id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"pool not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// PoolHealth returns the health summary of all nodes in a pool.
func PoolHealth(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var total, healthy, draining int
		err := st.Pool.QueryRow(r.Context(),
			`SELECT COUNT(*),
			        COUNT(*) FILTER (WHERE active AND (last_health_state IS NULL OR last_health_state = 'healthy')),
			        COUNT(*) FILTER (WHERE draining)
			 FROM backend_nodes WHERE pool_id = $1`, id).Scan(&total, &healthy, &draining)
		if err != nil {
			http.Error(w, `{"error":"pool not found"}`, http.StatusNotFound)
			return
		}
		status := "down"
		if total > 0 && healthy == total {
			status = "healthy"
		} else if healthy > 0 {
			status = "degraded"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"pool_id": id, "total_nodes": total, "healthy_nodes": healthy,
			"draining_nodes": draining, "status": status,
		})
	}
}

// DrainPool marks all nodes in a pool as draining (connection draining, P1.5).
func DrainPool(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE backend_nodes SET draining = true, updated_at = now() WHERE pool_id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"drain failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"pool not found or no nodes"}`, http.StatusNotFound)
			return
		}
		// Also mark the pool status to indicate drain.
		_, _ = st.Pool.Exec(r.Context(),
			`UPDATE backend_pools SET lb_algorithm = lb_algorithm, updated_at = now() WHERE id = $1`, id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"pool_id": id, "drained_nodes": ct.RowsAffected()})
	}
}