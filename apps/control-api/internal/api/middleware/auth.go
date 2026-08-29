package middleware

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Auth returns middleware that extracts the authenticated user from the
// shield_session cookie, loads their roles from PostgreSQL, and sets the RBAC
// context. Unauthenticated requests are rejected with 401 (P0.4).
func Auth(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Health and metrics are always open (no auth required).
			path := r.URL.Path
			if path == "/api/v1/health" || path == "/api/v1/metrics" {
				next.ServeHTTP(w, r)
				return
			}

			// Login/logout don't need a session.
			if path == "/api/v1/auth/login" || path == "/api/v1/auth/logout" {
				next.ServeHTTP(w, r)
				return
			}

			cookie, err := r.Cookie("shield_session")
			if err != nil || cookie.Value == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"unauthorized","code":"no_session"}`))
				return
			}

			// Parse "userID:token" from the cookie value.
			val := cookie.Value
			var userID, token string
			for i := 0; i < len(val); i++ {
				if val[i] == ':' {
					userID = val[:i]
					token = val[i+1:]
					break
				}
			}
			if userID == "" || token == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"unauthorized","code":"invalid_session"}`))
				return
			}

			// Look up the session and user.
			var email, roleName string
			err = pool.QueryRow(r.Context(),
				`SELECT u.email, COALESCE(
				  (SELECT r.name FROM roles r
				   JOIN user_group_memberships ugm ON ugm.group_id = r.id
				   WHERE ugm.user_id = u.id LIMIT 1), 'Read Only')
				 FROM sessions s
				 JOIN users u ON u.id = s.user_id
				 WHERE s.token_hash = $1 AND s.user_id = $2 AND s.expires_at > now()`,
				token, userID).Scan(&email, &roleName)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"unauthorized","code":"session_expired"}`))
				return
			}

			perms := RolePermissions[roleName]
			if perms == nil {
				perms = RolePermissions["Read Only"]
			}

			rbac := &RBACContext{
				User: RBACUser{
					ID:             userID,
					Email:          email,
					OrganizationID: "01ARZ3NDEKTSV4RRFFQ69G5FAV",
				},
				Roles:       []string{roleName},
				Permissions: perms,
			}
			ctx := ContextWithRBAC(r.Context(), rbac)
			ctx = ContextWithActor(ctx, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}