package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Auth returns middleware that extracts the authenticated user from the
// session cookie or Authorization header, loads their roles/permissions from
// PostgreSQL, and sets the RBAC context.
//
// Phase 3: uses X-User-Email + X-User-Role headers for development; when
// real auth is added (Phase 4+), this reads the session cookie instead.
func Auth(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Phase 3 dev-mode: read identity from headers.
			userEmail := r.Header.Get("X-User-Email")
			userRole := r.Header.Get("X-User-Role")

			if userEmail == "" {
				// No identity: serve as anonymous (read-only access).
				rbac := &RBACContext{
					User:        RBACUser{ID: "", Email: "anonymous", OrganizationID: "01ARZ3NDEKTSV4RRFFQ69G5FAV"},
					Roles:       []string{"Read Only"},
					Permissions: RolePermissions["Read Only"],
				}
				next.ServeHTTP(w, r.WithContext(ContextWithRBAC(r.Context(), rbac)))
				return
			}

			// Normalize the role.
			role := NormalizeRole(userRole)
			perms := RolePermissions[role]
			if perms == nil {
				perms = RolePermissions["Read Only"]
			}

			// Look up the user in DB for the real ID.
			var userID string
			userID = "" // default empty
			if pool != nil {
				if err := pool.QueryRow(r.Context(),
					`SELECT id FROM users WHERE email = $1`, userEmail).Scan(&userID); err != nil {
					// User not found — still allow with a best-effort ID.
					slog.Warn("auth: user not found", "email", userEmail)
				}
			}

			rbac := &RBACContext{
				User: RBACUser{
					ID:             userID,
					Email:          userEmail,
					OrganizationID: "01ARZ3NDEKTSV4RRFFQ69G5FAV",
				},
				Roles:       []string{role},
				Permissions: perms,
			}
			ctx := ContextWithRBAC(r.Context(), rbac)
			ctx = ContextWithActor(ctx, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RoutePermission maps HTTP method + path patterns to required permissions.
// Used by the router to automatically apply permission checks.
type RoutePermission struct {
	Method      string
	PathPrefix  string
	Permissions []string // at least one required
}

// DefaultRoutePermissions defines the permission map for Phase 3.
var DefaultRoutePermissions = []RoutePermission{
	// Applications
	{Method: "GET", PathPrefix: "/api/v1/applications", Permissions: []string{PermAppRead, PermAppWrite, PermAppAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/applications", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "PUT", PathPrefix: "/api/v1/applications", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/applications", Permissions: []string{PermAppAdmin}},

	// Gateways
	{Method: "GET", PathPrefix: "/api/v1/gateways", Permissions: []string{PermGatewayRead, PermGatewayWrite}},
	{Method: "POST", PathPrefix: "/api/v1/gateways", Permissions: []string{PermGatewayWrite}},
	{Method: "PUT", PathPrefix: "/api/v1/gateways", Permissions: []string{PermGatewayWrite}},

	// Security policies
	{Method: "GET", PathPrefix: "/api/v1/security-policies", Permissions: []string{PermPolicyRead, PermPolicyWrite, PermPolicyAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/security-policies", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "PUT", PathPrefix: "/api/v1/security-policies", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},

	// Security events
	{Method: "GET", PathPrefix: "/api/v1/security-events", Permissions: []string{PermEventRead, PermAuditRead}},

	// Audit log
	{Method: "GET", PathPrefix: "/api/v1/audit-events", Permissions: []string{PermAuditRead}},

	// IP lists
	{Method: "GET", PathPrefix: "/api/v1/ip-lists", Permissions: []string{PermIPListRead, PermIPListWrite}},
	{Method: "POST", PathPrefix: "/api/v1/ip-lists", Permissions: []string{PermIPListWrite}},

	// Rate limits
	{Method: "GET", PathPrefix: "/api/v1/rate-limits", Permissions: []string{PermRateLimitRead, PermRateLimitWrite}},
	{Method: "POST", PathPrefix: "/api/v1/rate-limits", Permissions: []string{PermRateLimitWrite}},

	// Certificates
	{Method: "GET", PathPrefix: "/api/v1/certificates", Permissions: []string{PermCertificateRead, PermCertificateWrite}},
	{Method: "POST", PathPrefix: "/api/v1/certificates", Permissions: []string{PermCertificateWrite}},

	// Users (admin only)
	{Method: "GET", PathPrefix: "/api/v1/users", Permissions: []string{PermUserRead, PermUserAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/users", Permissions: []string{PermUserAdmin}},
}

// CheckRoutePermission returns true if the request method+path matches any
// route permission rule and the context has at least one of the required perms.
func CheckRoutePermission(ctx context.Context, method, path string) bool {
	// Health and metrics are always open.
	if path == "/api/v1/health" || path == "/api/v1/metrics" {
		return true
	}
	for _, rp := range DefaultRoutePermissions {
		if rp.Method == method && strings.HasPrefix(path, rp.PathPrefix) {
			for _, p := range rp.Permissions {
				if HasPermission(ctx, p) {
					return true
				}
			}
			return false
		}
	}
	// Unknown route: default to requiring any permission.
	return HasPermission(ctx, "system:admin")
}

// RBACEnforcer returns middleware that checks route-level permissions using
// the DefaultRoutePermissions map. Must be placed AFTER the Auth middleware.
func RBACEnforcer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !CheckRoutePermission(r.Context(), r.Method, r.URL.Path) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"error":"forbidden","code":"route_not_authorized"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}