package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListBackendPools returns all backend pools with their nodes.
func ListBackendPools(st *store.Store) http.HandlerFunc {
	type node struct {
		ID        string `json:"id"`
		Host      string `json:"host"`
		Port      int    `json:"port"`
		Weight    int    `json:"weight"`
		Active    bool   `json:"active"`
		LastState string `json:"last_health_state,omitempty"`
	}
	type pool struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		ApplicationID string `json:"application_id"`
		LBAlgorithm   string `json:"lb_algorithm"`
		Version       int64  `json:"version"`
		Nodes         []node `json:"nodes"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, application_id, lb_algorithm, version FROM backend_pools ORDER BY name`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var pools []pool
		for rows.Next() {
			var p pool
			if err := rows.Scan(&p.ID, &p.Name, &p.ApplicationID, &p.LBAlgorithm, &p.Version); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			pools = append(pools, p)
		}

		// Attach nodes per pool.
		for i := range pools {
			nodeRows, err := st.Pool.Query(r.Context(),
				`SELECT id, host, port, weight, active, COALESCE(last_health_state,'')
				 FROM backend_nodes WHERE pool_id = $1 ORDER BY host, port`, pools[i].ID)
			if err != nil {
				continue
			}
			for nodeRows.Next() {
				var n node
				if err := nodeRows.Scan(&n.ID, &n.Host, &n.Port, &n.Weight, &n.Active, &n.LastState); err == nil {
					pools[i].Nodes = append(pools[i].Nodes, n)
				}
			}
			nodeRows.Close()
		}

		if pools == nil {
			pools = []pool{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pools)
	}
}

// CreateBackendPool creates a pool (and optional initial nodes).
func CreateBackendPool(st *store.Store) http.HandlerFunc {
	type create struct {
		Name          string `json:"name"`
		ApplicationID string `json:"application_id"`
		LBAlgorithm   string `json:"lb_algorithm"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body create
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Name == "" || body.ApplicationID == "" {
			http.Error(w, `{"error":"name and application_id required"}`, http.StatusBadRequest)
			return
		}
		if body.LBAlgorithm == "" {
			body.LBAlgorithm = "round_robin"
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}
		orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO backend_pools (id, organization_id, application_id, name, lb_algorithm)
			 VALUES ($1, $2, $3, $4, $5)`,
			id, orgID, body.ApplicationID, body.Name, body.LBAlgorithm); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// DeleteBackendPool deletes a pool (cascades to nodes).
func DeleteBackendPool(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM backend_pools WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"pool not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}