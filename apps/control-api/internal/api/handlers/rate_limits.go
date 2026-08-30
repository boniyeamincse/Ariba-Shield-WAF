package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListRateLimits returns all rate limit policies.
func ListRateLimits(st *store.Store) http.HandlerFunc {
	type rl struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		ApplicationID string `json:"application_id,omitempty"`
		RoutePrefix   string `json:"route_prefix"`
		LimitCount    int    `json:"limit_count"`
		WindowSeconds int    `json:"window_seconds"`
		Action        string `json:"action"`
		Version       int64  `json:"version"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, COALESCE(application_id,''), route_prefix, limit_count, window_seconds, action, version
			 FROM rate_limit_policies ORDER BY name`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var policies []rl
		for rows.Next() {
			var p rl
			if err := rows.Scan(&p.ID, &p.Name, &p.ApplicationID, &p.RoutePrefix, &p.LimitCount, &p.WindowSeconds, &p.Action, &p.Version); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			policies = append(policies, p)
		}
		if policies == nil {
			policies = []rl{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(policies)
	}
}

// CreateRateLimit creates a rate limit policy.
func CreateRateLimit(st *store.Store) http.HandlerFunc {
	type create struct {
		Name          string `json:"name"`
		ApplicationID string `json:"application_id"`
		RoutePrefix   string `json:"route_prefix"`
		LimitCount    int    `json:"limit_count"`
		WindowSeconds int    `json:"window_seconds"`
		Action        string `json:"action"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body create
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Name == "" || body.LimitCount <= 0 {
			http.Error(w, `{"error":"name and limit_count > 0 required"}`, http.StatusBadRequest)
			return
		}
		if body.WindowSeconds <= 0 {
			body.WindowSeconds = 60
		}
		if body.Action == "" {
			body.Action = "throttle"
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO rate_limit_policies (id, organization_id, name, application_id, route_prefix, limit_count, window_seconds, action)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			id, orgID, body.Name, nullIfEmpty(body.ApplicationID), body.RoutePrefix, body.LimitCount, body.WindowSeconds, body.Action); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// UpdateRateLimit updates a rate limit policy (partial update).
func UpdateRateLimit(st *store.Store) http.HandlerFunc {
	type update struct {
		Name          *string `json:"name"`
		RoutePrefix   *string `json:"route_prefix"`
		LimitCount    *int    `json:"limit_count"`
		WindowSeconds *int    `json:"window_seconds"`
		Action        *string `json:"action"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body update
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}

		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE rate_limit_policies SET
			   name = COALESCE($1, name),
			   route_prefix = COALESCE($2, route_prefix),
			   limit_count = COALESCE($3, limit_count),
			   window_seconds = COALESCE($4, window_seconds),
			   action = COALESCE($5, action),
			   version = version + 1, updated_at = now()
			 WHERE id = $6`,
			nullableString(body.Name), nullableString(body.RoutePrefix),
			nullableInt(body.LimitCount), nullableInt(body.WindowSeconds),
			nullableString(body.Action), id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"rate limit not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// DeleteRateLimit deletes a rate limit policy.
func DeleteRateLimit(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM rate_limit_policies WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"rate limit not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}

func nullableInt(i *int) any {
	if i == nil {
		return nil
	}
	return *i
}
