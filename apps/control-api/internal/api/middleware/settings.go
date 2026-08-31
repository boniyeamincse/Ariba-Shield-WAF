package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SettingsReader fetches settings on demand (per-request, simple + correct).
type SettingsReader func(ctx context.Context) (maintenance bool, sessionTimeoutMin int)

// SettingsMiddlewareFactory creates a handler that enforces system settings:
//   - maintenance_mode: returns 503 for non-admin write requests when enabled.
//   - session_timeout_minutes: exposed for session expiry (wired by auth).
func SettingsMiddlewareFactory(pool *pgxpool.Pool, reader SettingsReader) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Always allow auth endpoints (login, MFA, etc.) even in maintenance mode.
			if strings.HasPrefix(r.URL.Path, "/api/v1/auth/") ||
				r.URL.Path == "/api/v1/health" || r.URL.Path == "/api/v1/metrics" {
				next.ServeHTTP(w, r)
				return
			}
			maintenance, _ := reader(r.Context())
			if maintenance && r.Method != "GET" && r.Method != "HEAD" && r.Method != "OPTIONS" {
				// Allow admins through even in maintenance mode.
				if !HasPermission(r.Context(), PermSystemAdmin) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusServiceUnavailable)
					_, _ = w.Write([]byte(`{"error":"maintenance mode","code":"maintenance_active"}`))
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// NewSettingsReader builds a reader backed by the system_settings table.
// Values are stored as raw JSON primitives (true, 5, "text").
func NewSettingsReader(pool *pgxpool.Pool) SettingsReader {
	return func(ctx context.Context) (bool, int) {
		maintenance := false
		timeout := 60
		_ = pool.QueryRow(ctx,
			`SELECT
			   COALESCE((SELECT value::text FROM system_settings WHERE organization_id='01ARZ3NDEKTSV4RRFFQ69G5FAV' AND category='general' AND key='maintenance_mode'), 'false')::boolean,
			   COALESCE((SELECT value::text FROM system_settings WHERE organization_id='01ARZ3NDEKTSV4RRFFQ69G5FAV' AND category='security' AND key='session_timeout_minutes'), '60')::int
			`).Scan(&maintenance, &timeout)
		return maintenance, timeout
	}
}
