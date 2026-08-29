package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListHealthMonitors returns all health monitors.
func ListHealthMonitors(st *store.Store) http.HandlerFunc {
	type monitor struct {
		ID              string `json:"id"`
		PoolID          string `json:"pool_id"`
		Type            string `json:"type"`
		IntervalMS      int    `json:"interval_ms"`
		TimeoutMS       int    `json:"timeout_ms"`
		FailThreshold   int    `json:"fail_threshold"`
		PassThreshold   int    `json:"pass_threshold"`
		HTTPPath        string `json:"http_path"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, pool_id, type, interval_ms, timeout_ms, fail_threshold, pass_threshold, http_path
			 FROM health_monitors ORDER BY pool_id`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var monitors []monitor
		for rows.Next() {
			var m monitor
			if err := rows.Scan(&m.ID, &m.PoolID, &m.Type, &m.IntervalMS, &m.TimeoutMS,
				&m.FailThreshold, &m.PassThreshold, &m.HTTPPath); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			monitors = append(monitors, m)
		}
		if monitors == nil {
			monitors = []monitor{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(monitors)
	}
}

// CreateHealthMonitor creates a health monitor for a pool.
func CreateHealthMonitor(st *store.Store) http.HandlerFunc {
	type create struct {
		PoolID        string `json:"pool_id"`
		Type          string `json:"type"`
		IntervalMS    int    `json:"interval_ms"`
		TimeoutMS     int    `json:"timeout_ms"`
		FailThreshold int    `json:"fail_threshold"`
		PassThreshold int    `json:"pass_threshold"`
		HTTPPath      string `json:"http_path"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body create
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.PoolID == "" {
			http.Error(w, `{"error":"pool_id required"}`, http.StatusBadRequest)
			return
		}
		if body.Type == "" {
			body.Type = "http"
		}
		if body.IntervalMS == 0 {
			body.IntervalMS = 5000
		}
		if body.TimeoutMS == 0 {
			body.TimeoutMS = 2000
		}
		if body.FailThreshold == 0 {
			body.FailThreshold = 3
		}
		if body.PassThreshold == 0 {
			body.PassThreshold = 2
		}
		if body.HTTPPath == "" {
			body.HTTPPath = "/healthz"
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO health_monitors
			   (id, pool_id, type, interval_ms, timeout_ms, fail_threshold, pass_threshold, http_path)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			id, body.PoolID, body.Type, body.IntervalMS, body.TimeoutMS,
			body.FailThreshold, body.PassThreshold, body.HTTPPath); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// DeleteHealthMonitor deletes a health monitor.
func DeleteHealthMonitor(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM health_monitors WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"monitor not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}