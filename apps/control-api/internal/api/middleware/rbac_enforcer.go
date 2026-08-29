package middleware

import (
	"context"
	"net/http"
	"strings"
)

// RoutePermission maps HTTP method + path patterns to required permissions.
// Used by RBACEnforcer to apply permission checks per route.
type RoutePermission struct {
	Method      string
	PathPrefix  string
	Permissions []string // at least one required
}

// DefaultRoutePermissions defines the permission map for Phase 3/4.
var DefaultRoutePermissions = []RoutePermission{
	{Method: "GET", PathPrefix: "/api/v1/applications", Permissions: []string{PermAppRead, PermAppWrite, PermAppAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/applications", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/applications", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/applications", Permissions: []string{PermAppAdmin}},

	{Method: "GET", PathPrefix: "/api/v1/gateways", Permissions: []string{PermGatewayRead, PermGatewayWrite}},
	{Method: "POST", PathPrefix: "/api/v1/gateways", Permissions: []string{PermGatewayWrite}},
	{Method: "PUT", PathPrefix: "/api/v1/gateways", Permissions: []string{PermGatewayWrite}},

	{Method: "GET", PathPrefix: "/api/v1/security-policies", Permissions: []string{PermPolicyRead, PermPolicyWrite, PermPolicyAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/security-policies", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/security-policies", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/security-policies", Permissions: []string{PermPolicyAdmin}},

	{Method: "GET", PathPrefix: "/api/v1/security-events", Permissions: []string{PermEventRead, PermAuditRead}},
	{Method: "GET", PathPrefix: "/api/v1/audit-events", Permissions: []string{PermAuditRead}},

	{Method: "GET", PathPrefix: "/api/v1/ip-lists", Permissions: []string{PermIPListRead, PermIPListWrite}},
	{Method: "POST", PathPrefix: "/api/v1/ip-lists", Permissions: []string{PermIPListWrite}},
	{Method: "PATCH", PathPrefix: "/api/v1/ip-lists", Permissions: []string{PermIPListWrite}},
	{Method: "DELETE", PathPrefix: "/api/v1/ip-lists", Permissions: []string{PermIPListWrite}},

	{Method: "GET", PathPrefix: "/api/v1/rate-limits", Permissions: []string{PermRateLimitRead, PermRateLimitWrite}},
	{Method: "POST", PathPrefix: "/api/v1/rate-limits", Permissions: []string{PermRateLimitWrite}},
	{Method: "PATCH", PathPrefix: "/api/v1/rate-limits", Permissions: []string{PermRateLimitWrite}},
	{Method: "DELETE", PathPrefix: "/api/v1/rate-limits", Permissions: []string{PermRateLimitWrite}},

	{Method: "GET", PathPrefix: "/api/v1/certificates", Permissions: []string{PermCertificateRead, PermCertificateWrite}},
	{Method: "POST", PathPrefix: "/api/v1/certificates", Permissions: []string{PermCertificateWrite}},

	{Method: "GET", PathPrefix: "/api/v1/users", Permissions: []string{PermUserRead, PermUserAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/users", Permissions: []string{PermUserAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/users", Permissions: []string{PermUserAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/users", Permissions: []string{PermUserAdmin}},

	// Policy version lifecycle.
	{Method: "POST", PathPrefix: "/api/v1/security-policies", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/policy-versions", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/policy-versions", Permissions: []string{PermPolicyRead, PermPolicyWrite, PermPolicyAdmin}},

	// Custom/managed rules and deployments.
	{Method: "GET", PathPrefix: "/api/v1/custom-rules", Permissions: []string{PermPolicyRead, PermPolicyWrite, PermPolicyAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/custom-rules", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/managed-rules", Permissions: []string{PermPolicyRead, PermPolicyWrite, PermPolicyAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/managed-rules", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/deployments", Permissions: []string{PermPolicyRead, PermGatewayRead, PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/deployments", Permissions: []string{PermGatewayWrite, PermSystemAdmin}},

	// Webhooks / exceptions / integrations.
	{Method: "GET", PathPrefix: "/api/v1/webhooks", Permissions: []string{PermEventRead, PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/webhooks", Permissions: []string{PermSystemAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/exceptions", Permissions: []string{PermPolicyRead, PermPolicyWrite, PermPolicyAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/exceptions", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},

	// SRS §7 delivery endpoints.
	{Method: "GET", PathPrefix: "/api/v1/virtual-servers", Permissions: []string{PermAppRead, PermAppWrite, PermAppAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/virtual-servers", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/virtual-servers", Permissions: []string{PermAppAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/backend-pools", Permissions: []string{PermAppRead, PermAppWrite, PermAppAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/backend-pools", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/backend-pools", Permissions: []string{PermAppAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/health-monitors", Permissions: []string{PermAppRead, PermAppWrite, PermAppAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/health-monitors", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/health-monitors", Permissions: []string{PermAppAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/config-versions", Permissions: []string{PermPolicyRead, PermPolicyWrite, PermPolicyAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/traffic", Permissions: []string{PermEventRead, PermAuditRead}},
}

// CheckRoutePermission returns true if the request method+path is authorized
// for the RBAC context.
func CheckRoutePermission(ctx context.Context, method, path string) bool {
	if path == "/api/v1/health" || path == "/api/v1/metrics" || strings.HasPrefix(path, "/api/v1/auth/") {
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
	// Unknown route: default to requiring admin.
	return HasPermission(ctx, PermSystemAdmin)
}

// RBACEnforcer returns middleware that checks route-level permissions using
// DefaultRoutePermissions. Must be placed AFTER the Auth middleware.
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