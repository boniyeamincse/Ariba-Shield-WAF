package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

type backendNode struct {
	ID              string `json:"id"`
	PoolID          string `json:"pool_id"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Weight          int    `json:"weight"`
	Active          bool   `json:"active"`
	Draining        bool   `json:"draining"`
	LastHealthState string `json:"last_health_state,omitempty"`
}

// ListBackendNodes returns the nodes of a pool.
func ListBackendNodes(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		poolID := r.PathValue("id")
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, pool_id, host, port, weight, active, draining, COALESCE(last_health_state,'')
			 FROM backend_nodes WHERE pool_id = $1 ORDER BY host, port`, poolID)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		nodes := []backendNode{}
		for rows.Next() {
			var n backendNode
			if err := rows.Scan(&n.ID, &n.PoolID, &n.Host, &n.Port, &n.Weight, &n.Active, &n.Draining, &n.LastHealthState); err != nil {
				continue
			}
			nodes = append(nodes, n)
		}
		if nodes == nil {
			nodes = []backendNode{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(nodes)
	}
}

// CreateBackendNode adds a node to a pool.
func CreateBackendNode(st *store.Store) http.HandlerFunc {
	type create struct {
		Host   string `json:"host"`
		Port   int    `json:"port"`
		Weight int    `json:"weight"`
		Active *bool  `json:"active"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		poolID := r.PathValue("id")
		var body create
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Host == "" || body.Port == 0 {
			http.Error(w, `{"error":"host and port required"}`, http.StatusBadRequest)
			return
		}
		if body.Weight == 0 {
			body.Weight = 1
		}
		active := true
		if body.Active != nil {
			active = *body.Active
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}
		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO backend_nodes (id, pool_id, host, port, weight, active)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			id, poolID, body.Host, body.Port, body.Weight, active); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// GetBackendNode returns a single node.
func GetBackendNode(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var n backendNode
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, pool_id, host, port, weight, active, draining, COALESCE(last_health_state,'')
			 FROM backend_nodes WHERE id = $1`, id).
			Scan(&n.ID, &n.PoolID, &n.Host, &n.Port, &n.Weight, &n.Active, &n.Draining, &n.LastHealthState)
		if err != nil {
			http.Error(w, `{"error":"node not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(n)
	}
}

// UpdateBackendNode updates node fields (host, port, weight, active, draining).
func UpdateBackendNode(st *store.Store) http.HandlerFunc {
	type update struct {
		Host     *string `json:"host"`
		Port     *int    `json:"port"`
		Weight   *int    `json:"weight"`
		Active   *bool   `json:"active"`
		Draining *bool   `json:"draining"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body update
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE backend_nodes SET
			   host = COALESCE($1, host),
			   port = COALESCE($2, port),
			   weight = COALESCE($3, weight),
			   active = COALESCE($4, active),
			   draining = COALESCE($5, draining),
			   updated_at = now()
			 WHERE id = $6`,
			nullableString(body.Host), nullableInt(body.Port), nullableInt(body.Weight),
			body.Active, body.Draining, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"node not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// DeleteBackendNode removes a node.
func DeleteBackendNode(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM backend_nodes WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"node not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}

// EnableBackendNode marks a node active.
func EnableBackendNode(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE backend_nodes SET active = true, draining = false, updated_at = now() WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"node not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "enabled"})
	}
}

// DisableBackendNode marks a node inactive (stops receiving traffic).
func DisableBackendNode(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE backend_nodes SET active = false, updated_at = now() WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"node not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "disabled"})
	}
}