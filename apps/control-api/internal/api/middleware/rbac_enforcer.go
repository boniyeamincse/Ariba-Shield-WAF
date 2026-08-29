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
	{Method: "PATCH", PathPrefix: "/api/v1/gateways", Permissions: []string{PermGatewayWrite}},
	{Method: "DELETE", PathPrefix: "/api/v1/gateways", Permissions: []string{PermGatewayWrite}},

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
	{Method: "GET", PathPrefix: "/api/v1/listeners", Permissions: []string{PermAppRead, PermAppWrite, PermAppAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/listeners", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/listeners", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/listeners", Permissions: []string{PermAppAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/backend-pools", Permissions: []string{PermAppRead, PermAppWrite, PermAppAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/backend-pools", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/backend-pools", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/backend-pools", Permissions: []string{PermAppAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/backend-nodes", Permissions: []string{PermAppRead, PermAppWrite, PermAppAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/backend-nodes", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/backend-nodes", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/backend-nodes", Permissions: []string{PermAppAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/health-monitors", Permissions: []string{PermAppRead, PermAppWrite, PermAppAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/health-monitors", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/health-monitors", Permissions: []string{PermAppAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/routes", Permissions: []string{PermAppRead, PermAppWrite, PermAppAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/routes", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/routes", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/routes", Permissions: []string{PermAppAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/config-versions", Permissions: []string{PermPolicyRead, PermPolicyWrite, PermPolicyAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/traffic", Permissions: []string{PermEventRead, PermAuditRead}},

	// API roadmap modules (endpoint.md): GET read, write POST/PATCH/DELETE.
	// Route-prefix wildcards: declare each resource prefix once.
	{Method: "POST", PathPrefix: "/api/v1/config-validation", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/threat-intelligence", Permissions: []string{PermPolicyRead, PermIPListRead, PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/threat-intelligence", Permissions: []string{PermPolicyAdmin, PermSystemAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/threat-intelligence", Permissions: []string{PermPolicyAdmin, PermSystemAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/threat-intelligence", Permissions: []string{PermPolicyAdmin, PermSystemAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/api-security", Permissions: []string{PermPolicyRead, PermPolicyAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/api-security", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/api-security", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/api-security", Permissions: []string{PermPolicyAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/bot-management", Permissions: []string{PermPolicyRead, PermPolicyAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/bot-management", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/bot-management", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/bot-management", Permissions: []string{PermPolicyAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/dlp", Permissions: []string{PermPolicyRead, PermPolicyAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/dlp", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/dlp", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/dlp", Permissions: []string{PermPolicyAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/integrations", Permissions: []string{PermEventRead, PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/integrations", Permissions: []string{PermSystemAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/integrations", Permissions: []string{PermSystemAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/integrations", Permissions: []string{PermSystemAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/iam", Permissions: []string{PermUserRead, PermUserAdmin, PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/iam", Permissions: []string{PermUserAdmin, PermSystemAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/iam", Permissions: []string{PermUserAdmin, PermSystemAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/iam", Permissions: []string{PermUserAdmin, PermSystemAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/service-accounts", Permissions: []string{PermUserRead, PermUserAdmin, PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/service-accounts", Permissions: []string{PermUserAdmin, PermSystemAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/service-accounts", Permissions: []string{PermUserAdmin, PermSystemAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/api-keys", Permissions: []string{PermUserRead, PermUserAdmin, PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/api-keys", Permissions: []string{PermUserAdmin, PermSystemAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/api-keys", Permissions: []string{PermUserAdmin, PermSystemAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/secrets", Permissions: []string{PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/secrets", Permissions: []string{PermSystemAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/secrets", Permissions: []string{PermSystemAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/incident-response", Permissions: []string{PermEventRead, PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/incident-response", Permissions: []string{PermSystemAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/incident-response", Permissions: []string{PermSystemAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/incident-response", Permissions: []string{PermSystemAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/automation", Permissions: []string{PermPolicyRead, PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/automation", Permissions: []string{PermSystemAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/automation", Permissions: []string{PermSystemAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/automation", Permissions: []string{PermSystemAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/clusters", Permissions: []string{PermGatewayRead, PermGatewayWrite, PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/clusters", Permissions: []string{PermGatewayWrite, PermSystemAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/clusters", Permissions: []string{PermGatewayWrite, PermSystemAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/clusters", Permissions: []string{PermGatewayWrite, PermSystemAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/gateway-clusters", Permissions: []string{PermGatewayRead, PermGatewayWrite, PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/gateway-clusters", Permissions: []string{PermGatewayWrite, PermSystemAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/gateway-clusters", Permissions: []string{PermGatewayWrite, PermSystemAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/gateway-clusters", Permissions: []string{PermGatewayWrite, PermSystemAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/caching", Permissions: []string{PermAppRead, PermAppWrite, PermAppAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/caching", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/caching", Permissions: []string{PermAppWrite, PermAppAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/caching", Permissions: []string{PermAppAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/analytics", Permissions: []string{PermEventRead, PermAuditRead}},
	{Method: "GET", PathPrefix: "/api/v1/rule-analytics", Permissions: []string{PermEventRead, PermAuditRead}},
	{Method: "GET", PathPrefix: "/api/v1/organizations", Permissions: []string{PermUserAdmin, PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/organizations", Permissions: []string{PermSystemAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/organizations", Permissions: []string{PermSystemAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/organizations", Permissions: []string{PermSystemAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/tenants", Permissions: []string{PermUserRead, PermUserAdmin, PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/tenants", Permissions: []string{PermUserAdmin, PermSystemAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/tenants", Permissions: []string{PermUserAdmin, PermSystemAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/tenants", Permissions: []string{PermUserAdmin, PermSystemAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/sites", Permissions: []string{PermGatewayRead, PermGatewayWrite, PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/sites", Permissions: []string{PermGatewayWrite, PermSystemAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/sites", Permissions: []string{PermGatewayWrite, PermSystemAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/sites", Permissions: []string{PermGatewayWrite, PermSystemAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/graphql-security", Permissions: []string{PermPolicyRead, PermPolicyAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/graphql-security", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/graphql-security", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/graphql-security", Permissions: []string{PermPolicyAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/client-side-protection", Permissions: []string{PermPolicyRead, PermPolicyAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/client-side-protection", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/client-side-protection", Permissions: []string{PermPolicyWrite, PermPolicyAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/client-side-protection", Permissions: []string{PermPolicyAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/api-quotas", Permissions: []string{PermRateLimitRead, PermRateLimitWrite, PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/api-quotas", Permissions: []string{PermRateLimitWrite, PermSystemAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/api-quotas", Permissions: []string{PermRateLimitWrite, PermSystemAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/api-quotas", Permissions: []string{PermRateLimitWrite, PermSystemAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/ml-baselines", Permissions: []string{PermPolicyRead, PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/ml-baselines", Permissions: []string{PermPolicyWrite, PermSystemAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/ml-baselines", Permissions: []string{PermPolicyWrite, PermSystemAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/ml-baselines", Permissions: []string{PermPolicyAdmin, PermSystemAdmin}},
	{Method: "GET", PathPrefix: "/api/v1/network-protection", Permissions: []string{PermSystemAdmin}},
	{Method: "POST", PathPrefix: "/api/v1/network-protection", Permissions: []string{PermSystemAdmin}},
	{Method: "PATCH", PathPrefix: "/api/v1/network-protection", Permissions: []string{PermSystemAdmin}},
	{Method: "DELETE", PathPrefix: "/api/v1/network-protection", Permissions: []string{PermSystemAdmin}},
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