package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListConfigVersions returns config version history with status + hashes.
func ListConfigVersions(st *store.Store) http.HandlerFunc {
	type version struct {
		ID         string `json:"id"`
		Version    int64  `json:"version"`
		BundleHash string `json:"bundle_hash"`
		Status     string `json:"status"`
		CreatedBy  string `json:"created_by"`
		CreatedAt  string `json:"created_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, version, bundle_hash, status, created_by, created_at
			 FROM config_versions ORDER BY version DESC LIMIT 50`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var versions []version
		for rows.Next() {
			var v version
			var ts time.Time
			if err := rows.Scan(&v.ID, &v.Version, &v.BundleHash, &v.Status, &v.CreatedBy, &ts); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			v.CreatedAt = ts.Format(time.RFC3339)
			versions = append(versions, v)
		}
		if versions == nil {
			versions = []version{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(versions)
	}
}

// GetConfigVersion returns the full config document for a version.
func GetConfigVersion(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var blob json.RawMessage
		var bundleHash string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT blob, bundle_hash FROM config_versions WHERE id = $1`, id).Scan(&blob, &bundleHash); err != nil {
			http.Error(w, `{"error":"config version not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":          id,
			"bundle_hash": bundleHash,
			"document":    json.RawMessage(blob),
		})
	}
}

// ListTrafficRequests queries access events (structured request log).
func ListTrafficRequests(st *store.Store) http.HandlerFunc {
	type event struct {
		EventID       string  `json:"event_id"`
		RequestID     string  `json:"request_id"`
		GatewayID     string  `json:"gateway_id"`
		ApplicationID string  `json:"application_id"`
		ClientIP      string  `json:"client_ip"`
		Method        string  `json:"method"`
		Path          string  `json:"path"`
		Host          string  `json:"host"`
		Status        int     `json:"status"`
		Bytes         int64   `json:"bytes"`
		LatencyMS     float64 `json:"latency_ms"`
		BackendNode   string  `json:"backend_node"`
		Decision      string  `json:"decision_action"`
		CreatedAt     string  `json:"created_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		limit := 100
		where := "1=1"
		args := []any{}

		// Filters: application_id, client_ip, from/to (RFC3339).
		if v := r.URL.Query().Get("application_id"); v != "" {
			args = append(args, v)
			where += " AND application_id = $" + itoa(len(args))
		}
		if v := r.URL.Query().Get("client_ip"); v != "" {
			args = append(args, v)
			where += " AND client_ip = $" + itoa(len(args))
		}
		if v := r.URL.Query().Get("from"); v != "" {
			args = append(args, v)
			where += " AND created_at >= $" + itoa(len(args)) + "::timestamptz"
		}
		if v := r.URL.Query().Get("to"); v != "" {
			args = append(args, v)
			where += " AND created_at <= $" + itoa(len(args)) + "::timestamptz"
		}
		args = append(args, limit)

		query := `SELECT event_id, request_id, gateway_id, COALESCE(application_id,''),
		          COALESCE(client_ip,''), method, path, COALESCE(host,''),
		          COALESCE(status,0), bytes, COALESCE(latency_ms,0),
		          COALESCE(backend_node,''), COALESCE(decision_action,''), created_at
		          FROM access_events WHERE ` + where +
			` ORDER BY created_at DESC LIMIT $` + itoa(len(args))

		rows, err := st.Pool.Query(r.Context(), query, args...)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var events []event
		for rows.Next() {
			var e event
			var ts time.Time
			if err := rows.Scan(&e.EventID, &e.RequestID, &e.GatewayID, &e.ApplicationID,
				&e.ClientIP, &e.Method, &e.Path, &e.Host, &e.Status, &e.Bytes,
				&e.LatencyMS, &e.BackendNode, &e.Decision, &ts); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			e.CreatedAt = ts.Format(time.RFC3339)
			events = append(events, e)
		}
		if events == nil {
			events = []event{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"requests": events})
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [12]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}