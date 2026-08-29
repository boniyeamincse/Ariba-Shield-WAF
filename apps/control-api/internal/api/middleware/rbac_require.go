package middleware

import (
	"net/http"
	"strings"
)

// RequirePermission returns middleware that enforces a specific permission.
// If the user lacks the permission, returns 403 with a stable error code.
func RequirePermission(perm string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !HasPermission(r.Context(), perm) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte(`{"error":"forbidden","code":"insufficient_permission","permission":"` + perm + `"}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole returns middleware that enforces a role (OR of listed roles).
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, role := range roles {
				if HasRole(r.Context(), role) {
					next.ServeHTTP(w, r)
					return
				}
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"error":"forbidden","code":"insufficient_role"}`))
		})
	}
}

// RequireAnyPermission enforces "at least one of" the given permissions.
func RequireAnyPermission(perms ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, p := range perms {
				if HasPermission(r.Context(), p) {
					next.ServeHTTP(w, r)
					return
				}
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"error":"forbidden","code":"insufficient_permission"}`))
		})
	}
}

// NormalizeRole canonicalizes common role aliases.
func NormalizeRole(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "superadmin", "super_admin", "super-admin", "root":
		return "Super Admin"
	case "platformadmin", "platform_admin", "platform-admin":
		return "Platform Admin"
	case "securityadmin", "security_admin", "security-admin":
		return "Security Admin"
	case "appowner", "app_owner", "app-owner":
		return "App Owner"
	case "socanalyst", "soc_analyst", "soc-analyst", "analyst":
		return "SOC Analyst"
	case "auditor":
		return "Auditor"
	case "readonly", "read_only", "read-only":
		return "Read Only"
	default:
		return role
	}
}