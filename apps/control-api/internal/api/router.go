package api

import (
	"net/http"
	"time"

	"github.com/ariba-shield/control-api/internal/api/handlers"
	"github.com/ariba-shield/control-api/internal/api/middleware"
	"github.com/ariba-shield/control-api/internal/config"
	"github.com/ariba-shield/control-api/internal/store"
)

// NewRouter builds the HTTP handler with all routes and middleware.
func NewRouter(st *store.Store, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Public routes.
	mux.HandleFunc("GET /api/v1/health", handlers.Health())
	mux.HandleFunc("GET /api/v1/metrics", handlers.Metrics(st))

	// Phase 1 API (placeholder auth; real auth + RBAC in Phase 3).
	mux.HandleFunc("GET /api/v1/applications", handlers.ListApplications(st))
	mux.HandleFunc("POST /api/v1/applications", handlers.CreateApplication(st))
	mux.HandleFunc("PATCH /api/v1/applications/{id}", middleware.OptimisticConcurrency(st, "applications", "id", handlers.UpdateApplication(st)))
	mux.HandleFunc("DELETE /api/v1/applications/{id}", handlers.DeleteApplication(st))
	mux.HandleFunc("GET /api/v1/applications/{id}/domains", handlers.ListDomains(st))
	mux.HandleFunc("POST /api/v1/applications/{id}/domains", handlers.CreateDomain(st))
	mux.HandleFunc("GET /api/v1/applications/{id}/origins", handlers.ListOrigins(st))
	mux.HandleFunc("POST /api/v1/applications/{id}/origins", handlers.CreateOrigin(st))
	mux.HandleFunc("GET /api/v1/security-policies", handlers.ListSecurityPolicies(st))
	mux.HandleFunc("POST /api/v1/security-policies", handlers.CreateSecurityPolicy(st))
	mux.HandleFunc("PATCH /api/v1/security-policies/{id}", handlers.UpdateSecurityPolicy(st))
	mux.HandleFunc("DELETE /api/v1/security-policies/{id}", handlers.DeleteSecurityPolicy(st))

	// Gateway fleet (Phase 2 operations).
	mux.HandleFunc("POST /api/v1/gateways/register", handlers.RegisterGateway(st))
	mux.HandleFunc("POST /api/v1/gateways/{id}/heartbeat", handlers.Heartbeat(st))
	mux.HandleFunc("GET /api/v1/gateways", handlers.ListGateways(st))
	mux.HandleFunc("GET /api/v1/gateways/{id}/config/current", handlers.ConfigPull(st))

	// Security events (Phase 2).
	mux.HandleFunc("GET /api/v1/security-events", handlers.ListSecurityEvents(st))

	// Audit log (immutable).
	mux.HandleFunc("GET /api/v1/audit-events", handlers.ListAuditEvents(st))

	// Policy binding (Phase 2).
	mux.HandleFunc("POST /api/v1/security-policies/bind", handlers.BindPolicy(st))

	// Phase 3 — safe blocking: IP lists, rate limits, policy versions.
	mux.HandleFunc("GET /api/v1/ip-lists", handlers.ListIPLists(st))
	mux.HandleFunc("POST /api/v1/ip-lists", handlers.CreateIPList(st))
	mux.HandleFunc("PATCH /api/v1/ip-lists/{id}", handlers.UpdateIPList(st))
	mux.HandleFunc("DELETE /api/v1/ip-lists/{id}", handlers.DeleteIPList(st))
	mux.HandleFunc("GET /api/v1/rate-limits", handlers.ListRateLimits(st))
	mux.HandleFunc("POST /api/v1/rate-limits", handlers.CreateRateLimit(st))
	mux.HandleFunc("PATCH /api/v1/rate-limits/{id}", handlers.UpdateRateLimit(st))
	mux.HandleFunc("DELETE /api/v1/rate-limits/{id}", handlers.DeleteRateLimit(st))
	mux.HandleFunc("POST /api/v1/security-policies/{id}/versions", handlers.CreatePolicyVersion(st))
	mux.HandleFunc("POST /api/v1/policy-versions/{id}/activate", handlers.ActivatePolicyVersion(st))
	mux.HandleFunc("POST /api/v1/policy-versions/{id}/promote", handlers.PromotePolicyVersion(st))
	mux.HandleFunc("POST /api/v1/policy-versions/{id}/rollback", handlers.RollbackPolicyVersion(st))
	mux.HandleFunc("GET /api/v1/policy-versions/diff", handlers.DiffPolicyVersions(st))

	// SRS §7 delivery endpoints.
	mux.HandleFunc("GET /api/v1/virtual-servers", handlers.ListVirtualServers(st))
	mux.HandleFunc("POST /api/v1/virtual-servers", handlers.CreateVirtualServer(st))
	mux.HandleFunc("DELETE /api/v1/virtual-servers/{id}", handlers.DeleteVirtualServer(st))
	mux.HandleFunc("GET /api/v1/backend-pools", handlers.ListBackendPools(st))
	mux.HandleFunc("POST /api/v1/backend-pools", handlers.CreateBackendPool(st))
	mux.HandleFunc("DELETE /api/v1/backend-pools/{id}", handlers.DeleteBackendPool(st))
	mux.HandleFunc("GET /api/v1/health-monitors", handlers.ListHealthMonitors(st))
	mux.HandleFunc("POST /api/v1/health-monitors", handlers.CreateHealthMonitor(st))
	mux.HandleFunc("DELETE /api/v1/health-monitors/{id}", handlers.DeleteHealthMonitor(st))
	mux.HandleFunc("GET /api/v1/config-versions", handlers.ListConfigVersions(st))
	mux.HandleFunc("GET /api/v1/config-versions/{id}", handlers.GetConfigVersion(st))
	mux.HandleFunc("GET /api/v1/traffic/requests", handlers.ListTrafficRequests(st))

	// Auth & Identity
	mux.HandleFunc("POST /api/v1/auth/login", handlers.Login(st))
	mux.HandleFunc("POST /api/v1/auth/logout", handlers.Logout(st))
	mux.HandleFunc("GET /api/v1/auth/me", handlers.Me(st))

	// Backup / restore (Phase 4).
	mux.HandleFunc("GET /api/v1/backups", handlers.ListBackups(st))
	mux.HandleFunc("POST /api/v1/backups", handlers.CreateBackup(st))
	mux.HandleFunc("POST /api/v1/backups/{id}/restore", handlers.RestoreBackup(st))

	// Users (Phase 3 — RBAC-protected, read by SUPER_ADMIN / PLATFORM_ADMIN)
	mux.HandleFunc("GET /api/v1/users", handlers.ListUsers(st))
	mux.HandleFunc("POST /api/v1/users", handlers.CreateUser(st))
	mux.HandleFunc("PATCH /api/v1/users/{id}", handlers.UpdateUser(st))
	mux.HandleFunc("DELETE /api/v1/users/{id}", handlers.DeleteUser(st))

	// Phase 3 — Webhooks, Exceptions, Rules, Certificates, Deployments
	mux.HandleFunc("GET /api/v1/webhooks", handlers.ListWebhooks(st))
	mux.HandleFunc("POST /api/v1/webhooks", handlers.CreateWebhook(st))
	mux.HandleFunc("GET /api/v1/exceptions", handlers.ListExceptions(st))
	mux.HandleFunc("POST /api/v1/exceptions", handlers.CreateException(st))
	mux.HandleFunc("GET /api/v1/managed-rules", handlers.ListManagedRules(st))
	mux.HandleFunc("POST /api/v1/managed-rules/{id}", handlers.ConfigureManagedRules(st))
	mux.HandleFunc("GET /api/v1/custom-rules", handlers.ListCustomRules(st))
	mux.HandleFunc("POST /api/v1/custom-rules", handlers.CreateCustomRule(st))
	mux.HandleFunc("GET /api/v1/deployments", handlers.ListDeployments(st))
	mux.HandleFunc("POST /api/v1/deployments", handlers.SyncDeployment(st))
	mux.HandleFunc("GET /api/v1/certificates", handlers.ListCertificates(st))
	mux.HandleFunc("POST /api/v1/certificates", handlers.UploadCertificate(st))

	// ===== API roadmap modules (endpoint.md) =====

	// Phase 2: config validation (dry-run)
	mux.HandleFunc("POST /api/v1/config-validation", handlers.ValidateConfigDryRun(st))

	// Phase 3: threat intelligence
	mux.HandleFunc("GET /api/v1/threat-intelligence", handlers.ListResource(st, handlers.CRUDConfig{Table: "threat_feeds", JSONName: "feed"}))
	mux.HandleFunc("POST /api/v1/threat-intelligence", handlers.CreateResource(st, handlers.CRUDConfig{Table: "threat_feeds", JSONName: "feed", Required: []string{"name", "source"}, Columns: []string{"name", "source", "indicator_type", "indicators", "confidence", "category", "ttl_hours", "provenance", "status"}}))
	mux.HandleFunc("PATCH /api/v1/threat-intelligence/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "threat_feeds", JSONName: "feed", Columns: []string{"name", "indicators", "confidence", "category", "ttl_hours", "status"}}))
	mux.HandleFunc("DELETE /api/v1/threat-intelligence/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "threat_feeds", JSONName: "feed"}))

	// Phase 3: API security (schema validation, discovery)
	mux.HandleFunc("GET /api/v1/api-security", handlers.ListResource(st, handlers.CRUDConfig{Table: "api_schemas", JSONName: "schema"}))
	mux.HandleFunc("POST /api/v1/api-security", handlers.CreateResource(st, handlers.CRUDConfig{Table: "api_schemas", JSONName: "schema", Required: []string{"name"}, Columns: []string{"application_id", "name", "openapi_document", "path", "method", "status", "drift_alert"}}))
	mux.HandleFunc("PATCH /api/v1/api-security/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "api_schemas", JSONName: "schema", Columns: []string{"openapi_document", "status", "drift_alert"}}))
	mux.HandleFunc("DELETE /api/v1/api-security/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "api_schemas", JSONName: "schema"}))

	// Phase 3: bot management
	mux.HandleFunc("GET /api/v1/bot-management", handlers.ListResource(st, handlers.CRUDConfig{Table: "bot_policies", JSONName: "bot policy"}))
	mux.HandleFunc("POST /api/v1/bot-management", handlers.CreateResource(st, handlers.CRUDConfig{Table: "bot_policies", JSONName: "bot policy", Required: []string{"name"}, Columns: []string{"application_id", "name", "challenge_type", "known_bots", "automation_signals", "login_protection", "scrape_protection", "status"}}))
	mux.HandleFunc("PATCH /api/v1/bot-management/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "bot_policies", JSONName: "bot policy", Columns: []string{"challenge_type", "known_bots", "automation_signals", "login_protection", "scrape_protection", "status"}}))
	mux.HandleFunc("DELETE /api/v1/bot-management/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "bot_policies", JSONName: "bot policy"}))

	// Phase 3: DLP
	mux.HandleFunc("GET /api/v1/dlp", handlers.ListResource(st, handlers.CRUDConfig{Table: "dlp_profiles", JSONName: "dlp profile"}))
	mux.HandleFunc("POST /api/v1/dlp", handlers.CreateResource(st, handlers.CRUDConfig{Table: "dlp_profiles", JSONName: "dlp profile", Required: []string{"name"}, Columns: []string{"name", "scan_targets", "rules", "status"}}))
	mux.HandleFunc("PATCH /api/v1/dlp/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "dlp_profiles", JSONName: "dlp profile", Columns: []string{"scan_targets", "rules", "status"}}))
	mux.HandleFunc("DELETE /api/v1/dlp/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "dlp_profiles", JSONName: "dlp profile"}))

	// Phase 3: integrations (SIEM/log forwarding)
	mux.HandleFunc("GET /api/v1/integrations", handlers.ListResource(st, handlers.CRUDConfig{Table: "integrations", JSONName: "integration"}))
	mux.HandleFunc("POST /api/v1/integrations", handlers.CreateResource(st, handlers.CRUDConfig{Table: "integrations", JSONName: "integration", Required: []string{"type", "name"}, Columns: []string{"type", "name", "endpoint", "log_types", "enabled", "config"}}))
	mux.HandleFunc("PATCH /api/v1/integrations/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "integrations", JSONName: "integration", Columns: []string{"endpoint", "log_types", "enabled", "config"}}))
	mux.HandleFunc("DELETE /api/v1/integrations/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "integrations", JSONName: "integration"}))

	// Phase 3: IAM / SSO
	mux.HandleFunc("GET /api/v1/iam", handlers.ListResource(st, handlers.CRUDConfig{Table: "iam_sso", JSONName: "sso config"}))
	mux.HandleFunc("POST /api/v1/iam", handlers.CreateResource(st, handlers.CRUDConfig{Table: "iam_sso", JSONName: "sso config", Required: []string{"provider_name"}, Columns: []string{"provider_name", "provider_type", "idp_entity_id", "sso_url", "config", "enabled"}}))
	mux.HandleFunc("PATCH /api/v1/iam/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "iam_sso", JSONName: "sso config", Columns: []string{"provider_type", "idp_entity_id", "sso_url", "config", "enabled"}}))
	mux.HandleFunc("DELETE /api/v1/iam/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "iam_sso", JSONName: "sso config"}))

	// Phase 3: service accounts + API keys + secrets
	mux.HandleFunc("GET /api/v1/service-accounts", handlers.ListResource(st, handlers.CRUDConfig{Table: "service_accounts", JSONName: "service account"}))
	mux.HandleFunc("POST /api/v1/service-accounts", handlers.CreateResource(st, handlers.CRUDConfig{Table: "service_accounts", JSONName: "service account", Required: []string{"name"}, Columns: []string{"name", "roles", "status"}}))
	mux.HandleFunc("DELETE /api/v1/service-accounts/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "service_accounts", JSONName: "service account"}))

	mux.HandleFunc("GET /api/v1/api-keys", handlers.ListResource(st, handlers.CRUDConfig{Table: "api_keys", JSONName: "api key"}))
	mux.HandleFunc("POST /api/v1/api-keys", handlers.CreateResource(st, handlers.CRUDConfig{Table: "api_keys", JSONName: "api key", Required: []string{"name", "key_prefix", "key_hash"}, Columns: []string{"service_account_id", "name", "key_prefix", "key_hash", "expires_at", "status"}}))
	mux.HandleFunc("DELETE /api/v1/api-keys/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "api_keys", JSONName: "api key"}))

	mux.HandleFunc("GET /api/v1/secrets", handlers.ListResource(st, handlers.CRUDConfig{Table: "secrets", JSONName: "secret"}))
	mux.HandleFunc("POST /api/v1/secrets", handlers.CreateResource(st, handlers.CRUDConfig{Table: "secrets", JSONName: "secret", Required: []string{"name", "secret_ref"}, Columns: []string{"name", "secret_ref", "provider"}}))
	mux.HandleFunc("DELETE /api/v1/secrets/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "secrets", JSONName: "secret"}))

	// Phase 4: incidents + automation + clusters + caching
	mux.HandleFunc("GET /api/v1/incident-response", handlers.ListResource(st, handlers.CRUDConfig{Table: "incidents", JSONName: "incident"}))
	mux.HandleFunc("POST /api/v1/incident-response", handlers.CreateResource(st, handlers.CRUDConfig{Table: "incidents", JSONName: "incident", Required: []string{"title"}, Columns: []string{"title", "severity", "status", "owner_user_id", "notes", "related_events"}}))
	mux.HandleFunc("PATCH /api/v1/incident-response/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "incidents", JSONName: "incident", Columns: []string{"severity", "status", "owner_user_id", "notes"}}))
	mux.HandleFunc("DELETE /api/v1/incident-response/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "incidents", JSONName: "incident"}))

	mux.HandleFunc("GET /api/v1/automation", handlers.ListResource(st, handlers.CRUDConfig{Table: "automation_rules", JSONName: "automation rule"}))
	mux.HandleFunc("POST /api/v1/automation", handlers.CreateResource(st, handlers.CRUDConfig{Table: "automation_rules", JSONName: "automation rule", Required: []string{"name", "trigger_event"}, Columns: []string{"name", "trigger_event", "action", "target", "enabled"}}))
	mux.HandleFunc("PATCH /api/v1/automation/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "automation_rules", JSONName: "automation rule", Columns: []string{"action", "target", "enabled"}}))
	mux.HandleFunc("DELETE /api/v1/automation/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "automation_rules", JSONName: "automation rule"}))

	mux.HandleFunc("GET /api/v1/clusters", handlers.ListResource(st, handlers.CRUDConfig{Table: "clusters", JSONName: "cluster"}))
	mux.HandleFunc("POST /api/v1/clusters", handlers.CreateResource(st, handlers.CRUDConfig{Table: "clusters", JSONName: "cluster", Required: []string{"name"}, Columns: []string{"name", "site", "gateway_ids", "ha_enabled", "status"}}))
	mux.HandleFunc("PATCH /api/v1/clusters/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "clusters", JSONName: "cluster", Columns: []string{"site", "gateway_ids", "ha_enabled", "status"}}))
	mux.HandleFunc("DELETE /api/v1/clusters/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "clusters", JSONName: "cluster"}))

	mux.HandleFunc("GET /api/v1/caching", handlers.ListResource(st, handlers.CRUDConfig{Table: "caching_rules", JSONName: "caching rule"}))
	mux.HandleFunc("POST /api/v1/caching", handlers.CreateResource(st, handlers.CRUDConfig{Table: "caching_rules", JSONName: "caching rule", Required: []string{"name"}, Columns: []string{"application_id", "name", "path_pattern", "ttl_seconds", "cache_methods", "status"}}))
	mux.HandleFunc("PATCH /api/v1/caching/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "caching_rules", JSONName: "caching rule", Columns: []string{"path_pattern", "ttl_seconds", "cache_methods", "status"}}))
	mux.HandleFunc("DELETE /api/v1/caching/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "caching_rules", JSONName: "caching rule"}))

	// Phase 4: analytics (aggregate queries over events)
	mux.HandleFunc("GET /api/v1/analytics", handlers.SecurityAnalytics(st))
	mux.HandleFunc("GET /api/v1/rule-analytics", handlers.RuleAnalytics(st))

	// Phase 4: organizations / workspaces (multi-tenancy foundation)
	mux.HandleFunc("GET /api/v1/organizations", handlers.ListResource(st, handlers.CRUDConfig{Table: "organizations", JSONName: "organization"}))
	mux.HandleFunc("POST /api/v1/organizations", handlers.CreateResource(st, handlers.CRUDConfig{Table: "organizations", JSONName: "organization", Required: []string{"name"}, Columns: []string{"name"}}))
	mux.HandleFunc("GET /api/v1/workspaces", handlers.ListResource(st, handlers.CRUDConfig{Table: "organizations", JSONName: "workspace"}))

	// Phase 5: graphql security, CSP, quotas, ML baselines, network protection
	mux.HandleFunc("GET /api/v1/graphql-security", handlers.ListResource(st, handlers.CRUDConfig{Table: "graphql_security", JSONName: "graphql policy"}))
	mux.HandleFunc("POST /api/v1/graphql-security", handlers.CreateResource(st, handlers.CRUDConfig{Table: "graphql_security", JSONName: "graphql policy", Required: []string{"name"}, Columns: []string{"application_id", "name", "max_depth", "max_complexity", "introspection", "status"}}))
	mux.HandleFunc("PATCH /api/v1/graphql-security/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "graphql_security", JSONName: "graphql policy", Columns: []string{"max_depth", "max_complexity", "introspection", "status"}}))
	mux.HandleFunc("DELETE /api/v1/graphql-security/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "graphql_security", JSONName: "graphql policy"}))

	mux.HandleFunc("GET /api/v1/client-side-protection", handlers.ListResource(st, handlers.CRUDConfig{Table: "csp_profiles", JSONName: "csp profile"}))
	mux.HandleFunc("POST /api/v1/client-side-protection", handlers.CreateResource(st, handlers.CRUDConfig{Table: "csp_profiles", JSONName: "csp profile", Required: []string{"name"}, Columns: []string{"application_id", "name", "csp_header", "inject_mode", "status"}}))
	mux.HandleFunc("PATCH /api/v1/client-side-protection/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "csp_profiles", JSONName: "csp profile", Columns: []string{"csp_header", "inject_mode", "status"}}))
	mux.HandleFunc("DELETE /api/v1/client-side-protection/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "csp_profiles", JSONName: "csp profile"}))

	mux.HandleFunc("GET /api/v1/api-quotas", handlers.ListResource(st, handlers.CRUDConfig{Table: "api_quotas", JSONName: "api quota"}))
	mux.HandleFunc("POST /api/v1/api-quotas", handlers.CreateResource(st, handlers.CRUDConfig{Table: "api_quotas", JSONName: "api quota", Required: []string{"name", "limit_count"}, Columns: []string{"application_id", "name", "match_by", "limit_count", "window_seconds", "status"}}))
	mux.HandleFunc("PATCH /api/v1/api-quotas/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "api_quotas", JSONName: "api quota", Columns: []string{"match_by", "limit_count", "window_seconds", "status"}}))
	mux.HandleFunc("DELETE /api/v1/api-quotas/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "api_quotas", JSONName: "api quota"}))

	mux.HandleFunc("GET /api/v1/ml-baselines", handlers.ListResource(st, handlers.CRUDConfig{Table: "ml_baselines", JSONName: "ml baseline"}))
	mux.HandleFunc("POST /api/v1/ml-baselines", handlers.CreateResource(st, handlers.CRUDConfig{Table: "ml_baselines", JSONName: "ml baseline", Required: []string{"name"}, Columns: []string{"application_id", "name", "status", "config"}}))
	mux.HandleFunc("PATCH /api/v1/ml-baselines/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "ml_baselines", JSONName: "ml baseline", Columns: []string{"status", "config"}}))
	mux.HandleFunc("DELETE /api/v1/ml-baselines/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "ml_baselines", JSONName: "ml baseline"}))

	mux.HandleFunc("GET /api/v1/network-protection", handlers.ListResource(st, handlers.CRUDConfig{Table: "network_protection", JSONName: "network protection"}))
	mux.HandleFunc("POST /api/v1/network-protection", handlers.CreateResource(st, handlers.CRUDConfig{Table: "network_protection", JSONName: "network protection", Required: []string{"name"}, Columns: []string{"name", "protection_type", "config", "status"}}))
	mux.HandleFunc("PATCH /api/v1/network-protection/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "network_protection", JSONName: "network protection", Columns: []string{"config", "status"}}))
	mux.HandleFunc("DELETE /api/v1/network-protection/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "network_protection", JSONName: "network protection"}))

	// Apply middleware stack. Execution order (outermost wraps last, runs first):
	//   Recovery -> RequestID -> Auth -> RBACEnforcer -> Idempotency -> Audit -> Logging -> mux
	// Auth MUST run before RBACEnforcer so that the RBAC context (user + permissions)
	// is populated before the enforcer reads it.
	var h http.Handler = mux
	store := middleware.NewIdempotencyStore(24 * time.Hour)
	h = middleware.Logging(h)
	h = middleware.Audit(st.Pool, h)
	h = middleware.Idempotency(store)(h)
	h = middleware.RBACEnforcer(h)
	h = middleware.Auth(st)(h)
	h = middleware.RequestID(h)
	h = middleware.Recovery(h)

	// Global CORS Middleware
	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		// Allow specific origins or just localhost for development
		if origin == "http://localhost:3002" || origin == "http://127.0.0.1:3002" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})

	return corsHandler
}