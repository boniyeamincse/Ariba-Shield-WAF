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
	mux.HandleFunc("GET /api/v1/applications/{id}", handlers.GetApplication(st))
	mux.HandleFunc("GET /api/v1/applications/{id}/traffic", handlers.ApplicationTraffic(st))
	mux.HandleFunc("GET /api/v1/applications/{id}/events", handlers.ApplicationEvents(st))
	mux.HandleFunc("GET /api/v1/applications/{id}/incidents", handlers.ApplicationIncidents(st))
	mux.HandleFunc("GET /api/v1/applications/{id}/policies", handlers.ApplicationPolicies(st))
	mux.HandleFunc("GET /api/v1/applications/{id}/health", handlers.ApplicationHealth(st))
	mux.HandleFunc("DELETE /api/v1/applications/{id}", handlers.DeleteApplication(st))
	mux.HandleFunc("GET /api/v1/applications/{id}/domains", handlers.ListDomains(st))
	mux.HandleFunc("POST /api/v1/applications/{id}/domains", handlers.CreateDomain(st))
	mux.HandleFunc("GET /api/v1/applications/{id}/origins", handlers.ListOrigins(st))
	mux.HandleFunc("POST /api/v1/applications/{id}/origins", handlers.CreateOrigin(st))
	mux.HandleFunc("GET /api/v1/security-policies", handlers.ListSecurityPolicies(st))
	mux.HandleFunc("POST /api/v1/security-policies", handlers.CreateSecurityPolicy(st))
	mux.HandleFunc("GET /api/v1/security-policies/{id}", handlers.GetSecurityPolicy(st))
	mux.HandleFunc("PATCH /api/v1/security-policies/{id}", handlers.UpdateSecurityPolicy(st))
	mux.HandleFunc("DELETE /api/v1/security-policies/{id}", handlers.DeleteSecurityPolicy(st))
	mux.HandleFunc("POST /api/v1/security-policies/{id}/validate", handlers.ValidatePolicy(st))
	mux.HandleFunc("POST /api/v1/security-policies/{id}/activate", handlers.ActivatePolicy(st))
	mux.HandleFunc("POST /api/v1/security-policies/{id}/disable", handlers.DisablePolicy(st))
	mux.HandleFunc("POST /api/v1/security-policies/{id}/rollback", handlers.RollbackPolicy(st))
	mux.HandleFunc("GET /api/v1/security-policies/{id}/versions", handlers.ListPolicyVersions(st))
	mux.HandleFunc("GET /api/v1/security-policies/{id}/diff", handlers.DiffPolicy(st))
	mux.HandleFunc("POST /api/v1/security-policies/{id}/clone", handlers.ClonePolicy(st))

	// Policy approvals (four-eyes workflow, §7.1)
	mux.HandleFunc("GET /api/v1/policy-approvals", handlers.ListPolicyApprovals(st))
	mux.HandleFunc("POST /api/v1/policy-approvals", handlers.CreatePolicyApproval(st))
	mux.HandleFunc("GET /api/v1/policy-approvals/{id}", handlers.GetPolicyApproval(st))
	mux.HandleFunc("POST /api/v1/policy-approvals/{id}/approve", handlers.ApprovePolicyApproval(st))
	mux.HandleFunc("POST /api/v1/policy-approvals/{id}/reject", handlers.RejectPolicyApproval(st))

	// WAF Rules / Signatures (master plan §6.3)
	mux.HandleFunc("GET /api/v1/rules", handlers.ListResource(st, handlers.CRUDConfig{Table: "rules", JSONName: "rule"}))
	mux.HandleFunc("POST /api/v1/rules", handlers.CreateResource(st, handlers.CRUDConfig{Table: "rules", JSONName: "rule", Required: []string{"rule_id", "name"}, Columns: []string{"rule_id", "name", "description", "action", "severity", "phase", "source", "status"}}))
	mux.HandleFunc("GET /api/v1/rules/{id}", handlers.GetRule(st))
	mux.HandleFunc("PATCH /api/v1/rules/{id}", handlers.UpdateRule(st))
	mux.HandleFunc("DELETE /api/v1/rules/{id}", handlers.DeleteRule(st))
	mux.HandleFunc("GET /api/v1/rules/{id}/versions", handlers.ListRuleVersions(st))
	mux.HandleFunc("POST /api/v1/rules/{id}/test", handlers.TestRule(st))
	mux.HandleFunc("POST /api/v1/rules/{id}/enable", handlers.EnableRule(st))
	mux.HandleFunc("POST /api/v1/rules/{id}/disable", handlers.DisableRule(st))

	// Rule Bundles
	mux.HandleFunc("GET /api/v1/rule-bundles", handlers.ListResource(st, handlers.CRUDConfig{Table: "rule_bundles", JSONName: "bundle"}))
	mux.HandleFunc("POST /api/v1/rule-bundles", handlers.CreateResource(st, handlers.CRUDConfig{Table: "rule_bundles", JSONName: "bundle", Required: []string{"name"}, Columns: []string{"name", "description", "rule_ids", "status"}}))
	mux.HandleFunc("GET /api/v1/rule-bundles/{id}", handlers.GetRuleBundle(st))
	mux.HandleFunc("POST /api/v1/rule-bundles/{id}/sign", handlers.SignRuleBundle(st))
	mux.HandleFunc("POST /api/v1/rule-bundles/{id}/deploy", handlers.DeployRuleBundle(st))
	mux.HandleFunc("POST /api/v1/rule-bundles/{id}/rollback", handlers.RollbackRuleBundle(st))

	// Gateway fleet (Phase 2 operations).
	mux.HandleFunc("POST /api/v1/gateways/register", handlers.RegisterGateway(st))
	mux.HandleFunc("POST /api/v1/gateways/{id}/heartbeat", handlers.Heartbeat(st))
	mux.HandleFunc("GET /api/v1/gateways", handlers.ListGateways(st))
	mux.HandleFunc("GET /api/v1/gateways/{id}/config/current", handlers.ConfigPull(st))
	mux.HandleFunc("GET /api/v1/gateways/{id}", handlers.GetGateway(st))
	mux.HandleFunc("PATCH /api/v1/gateways/{id}", handlers.UpdateGateway(st))
	mux.HandleFunc("DELETE /api/v1/gateways/{id}", handlers.DeleteGateway(st))
	mux.HandleFunc("GET /api/v1/gateways/{id}/status", handlers.GatewayStatus(st))
	mux.HandleFunc("GET /api/v1/gateways/{id}/config", handlers.GatewayConfig(st))
	mux.HandleFunc("POST /api/v1/gateways/{id}/config/apply", handlers.ApplyGatewayConfig(st))
	mux.HandleFunc("POST /api/v1/gateways/{id}/config/rollback", handlers.RollbackGatewayConfig(st))
	mux.HandleFunc("GET /api/v1/gateways/{id}/metrics", handlers.GatewayMetrics(st))

	// Security events (Phase 2).
	mux.HandleFunc("GET /api/v1/security-events", handlers.ListSecurityEvents(st))
	mux.HandleFunc("GET /api/v1/security-events/{id}", handlers.GetSecurityEvent(st))
	mux.HandleFunc("GET /api/v1/security-events/{id}/matches", handlers.EventMatches(st))
	mux.HandleFunc("GET /api/v1/security-events/{id}/timeline", handlers.EventTimeline(st))
	mux.HandleFunc("POST /api/v1/security-events/{id}/mask", handlers.MaskSecurityEvent(st))
	mux.HandleFunc("POST /api/v1/security-events/{id}/export", handlers.ExportSecurityEvent(st))

	// Audit log (immutable).
	mux.HandleFunc("GET /api/v1/audit-events", handlers.ListAuditEvents(st))
	mux.HandleFunc("GET /api/v1/audit-events/{id}", handlers.GetAuditEvent(st))
	mux.HandleFunc("GET /api/v1/audit-events/export", handlers.GetAuditEventExport(st))

	// Policy binding (Phase 2).
	mux.HandleFunc("POST /api/v1/security-policies/bind", handlers.BindPolicy(st))

	// Phase 3 — safe blocking: IP lists, rate limits, policy versions.
	mux.HandleFunc("GET /api/v1/ip-lists", handlers.ListIPLists(st))
	mux.HandleFunc("POST /api/v1/ip-lists", handlers.CreateIPList(st))
	mux.HandleFunc("GET /api/v1/ip-lists/{id}", handlers.GetIPList(st))
	mux.HandleFunc("PATCH /api/v1/ip-lists/{id}", handlers.UpdateIPList(st))
	mux.HandleFunc("DELETE /api/v1/ip-lists/{id}", handlers.DeleteIPList(st))
	mux.HandleFunc("GET /api/v1/ip-lists/{id}/entries", handlers.ListIPEntries(st))
	mux.HandleFunc("POST /api/v1/ip-lists/{id}/entries", handlers.AddIPEntry(st))
	mux.HandleFunc("DELETE /api/v1/ip-lists/{id}/entries/{entryId}", handlers.DeleteIPEntry(st))
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
	mux.HandleFunc("GET /api/v1/backend-pools/{id}", handlers.GetBackendPool(st))
	mux.HandleFunc("PATCH /api/v1/backend-pools/{id}", handlers.UpdateBackendPool(st))
	mux.HandleFunc("DELETE /api/v1/backend-pools/{id}", handlers.DeleteBackendPool(st))
	mux.HandleFunc("GET /api/v1/backend-pools/{id}/health", handlers.PoolHealth(st))
	mux.HandleFunc("POST /api/v1/backend-pools/{id}/drain", handlers.DrainPool(st))
	mux.HandleFunc("GET /api/v1/backend-pools/{id}/nodes", handlers.ListBackendNodes(st))
	mux.HandleFunc("POST /api/v1/backend-pools/{id}/nodes", handlers.CreateBackendNode(st))
	mux.HandleFunc("GET /api/v1/backend-nodes/{id}", handlers.GetBackendNode(st))
	mux.HandleFunc("PATCH /api/v1/backend-nodes/{id}", handlers.UpdateBackendNode(st))
	mux.HandleFunc("DELETE /api/v1/backend-nodes/{id}", handlers.DeleteBackendNode(st))
	mux.HandleFunc("POST /api/v1/backend-nodes/{id}/enable", handlers.EnableBackendNode(st))
	mux.HandleFunc("POST /api/v1/backend-nodes/{id}/disable", handlers.DisableBackendNode(st))
	mux.HandleFunc("GET /api/v1/health-monitors", handlers.ListHealthMonitors(st))
	mux.HandleFunc("POST /api/v1/health-monitors", handlers.CreateHealthMonitor(st))
	mux.HandleFunc("DELETE /api/v1/health-monitors/{id}", handlers.DeleteHealthMonitor(st))

	// Routes (URL routing rules)
	mux.HandleFunc("GET /api/v1/routes", handlers.ListResource(st, handlers.CRUDConfig{Table: "routes", JSONName: "route"}))
	mux.HandleFunc("POST /api/v1/routes", handlers.CreateResource(st, handlers.CRUDConfig{Table: "routes", JSONName: "route", Required: []string{"virtual_server_id", "path", "backend_pool_id"}, Columns: []string{"virtual_server_id", "path", "match_type", "backend_pool_id", "priority"}}))
	mux.HandleFunc("GET /api/v1/routes/{id}", handlers.GetResource(st, handlers.CRUDConfig{Table: "routes", JSONName: "route"}))
	mux.HandleFunc("PATCH /api/v1/routes/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "routes", JSONName: "route", Columns: []string{"path", "match_type", "backend_pool_id", "priority"}}))
	mux.HandleFunc("DELETE /api/v1/routes/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "routes", JSONName: "route"}))
	mux.HandleFunc("GET /api/v1/config-versions", handlers.ListConfigVersions(st))
	mux.HandleFunc("GET /api/v1/config-versions/{id}", handlers.GetConfigVersion(st))
	mux.HandleFunc("GET /api/v1/traffic/requests", handlers.ListTrafficRequests(st))

	// Auth & Identity
	mux.HandleFunc("POST /api/v1/auth/login", handlers.Login(st))
	mux.HandleFunc("POST /api/v1/auth/logout", handlers.Logout(st))
	mux.HandleFunc("POST /api/v1/auth/refresh", handlers.Refresh(st))
	mux.HandleFunc("GET /api/v1/auth/me", handlers.Me(st))
	mux.HandleFunc("POST /api/v1/auth/mfa/enable", handlers.EnableMFA(st))
	mux.HandleFunc("POST /api/v1/auth/mfa/verify", handlers.VerifyMFA(st))
	mux.HandleFunc("POST /api/v1/auth/mfa/disable", handlers.DisableMFA(st))
	mux.HandleFunc("POST /api/v1/auth/password/change", handlers.ChangePassword(st))
	mux.HandleFunc("POST /api/v1/auth/break-glass", handlers.BreakGlass(st))

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
	mux.HandleFunc("GET /api/v1/exceptions/{id}", handlers.GetException(st))
	mux.HandleFunc("PATCH /api/v1/exceptions/{id}", handlers.UpdateException(st))
	mux.HandleFunc("DELETE /api/v1/exceptions/{id}", handlers.DeleteException(st))
	mux.HandleFunc("POST /api/v1/exceptions/{id}/approve", handlers.ApproveException(st))
	mux.HandleFunc("POST /api/v1/exceptions/{id}/expire", handlers.ExpireException(st))
	mux.HandleFunc("GET /api/v1/managed-rules", handlers.ListManagedRules(st))
	mux.HandleFunc("POST /api/v1/managed-rules/{id}", handlers.ConfigureManagedRules(st))
	mux.HandleFunc("GET /api/v1/custom-rules", handlers.ListCustomRules(st))
	mux.HandleFunc("POST /api/v1/custom-rules", handlers.CreateCustomRule(st))
	mux.HandleFunc("GET /api/v1/deployments", handlers.ListDeployments(st))
	mux.HandleFunc("POST /api/v1/deployments", handlers.SyncDeployment(st))
	mux.HandleFunc("GET /api/v1/certificates", handlers.ListCertificates(st))
	mux.HandleFunc("POST /api/v1/certificates", handlers.UploadCertificate(st))
	mux.HandleFunc("GET /api/v1/certificates/{id}", handlers.GetCertificate(st))
	mux.HandleFunc("DELETE /api/v1/certificates/{id}", handlers.DeleteCertificate(st))
	mux.HandleFunc("POST /api/v1/certificates/import", handlers.ImportCertificate(st))
	mux.HandleFunc("POST /api/v1/certificates/{id}/renew", handlers.RenewCertificate(st))
	mux.HandleFunc("GET /api/v1/certificates/{id}/expiry", handlers.CertExpiry(st))

	// TLS profiles
	mux.HandleFunc("GET /api/v1/tls-profiles", handlers.ListResource(st, handlers.CRUDConfig{Table: "tls_profiles", JSONName: "tls profile"}))
	mux.HandleFunc("POST /api/v1/tls-profiles", handlers.CreateResource(st, handlers.CRUDConfig{Table: "tls_profiles", JSONName: "tls profile", Required: []string{"name"}, Columns: []string{"name", "min_version", "max_version", "cipher_suites", "certificate_ref", "renegotiation", "status"}}))
	mux.HandleFunc("GET /api/v1/tls-profiles/{id}", handlers.GetResource(st, handlers.CRUDConfig{Table: "tls_profiles", JSONName: "tls profile"}))
	mux.HandleFunc("PATCH /api/v1/tls-profiles/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "tls_profiles", JSONName: "tls profile", Columns: []string{"name", "min_version", "max_version", "cipher_suites", "certificate_ref", "renegotiation", "status"}}))
	mux.HandleFunc("DELETE /api/v1/tls-profiles/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "tls_profiles", JSONName: "tls profile"}))

	// ===== API roadmap modules (endpoint.md) =====

	// Phase 2: config validation (dry-run)
	mux.HandleFunc("POST /api/v1/config-validation", handlers.ValidateConfigDryRun(st))

	// Phase 3: threat intelligence
	mux.HandleFunc("GET /api/v1/threat-intelligence", handlers.ListResource(st, handlers.CRUDConfig{Table: "threat_feeds", JSONName: "feed"}))
	mux.HandleFunc("POST /api/v1/threat-intelligence", handlers.CreateResource(st, handlers.CRUDConfig{Table: "threat_feeds", JSONName: "feed", Required: []string{"name", "source"}, Columns: []string{"name", "source", "indicator_type", "indicators", "confidence", "category", "ttl_hours", "provenance", "status"}}))
	mux.HandleFunc("PATCH /api/v1/threat-intelligence/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "threat_feeds", JSONName: "feed", Columns: []string{"name", "indicators", "confidence", "category", "ttl_hours", "status"}}))
	mux.HandleFunc("DELETE /api/v1/threat-intelligence/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "threat_feeds", JSONName: "feed"}))
	// Threat feeds (canonical path per API list)
	mux.HandleFunc("GET /api/v1/threat-feeds", handlers.ListResource(st, handlers.CRUDConfig{Table: "threat_feeds", JSONName: "feed"}))
	mux.HandleFunc("POST /api/v1/threat-feeds", handlers.CreateResource(st, handlers.CRUDConfig{Table: "threat_feeds", JSONName: "feed", Required: []string{"name", "source"}, Columns: []string{"name", "source", "indicator_type", "indicators", "confidence", "category", "ttl_hours", "provenance", "status"}}))
	mux.HandleFunc("GET /api/v1/threat-feeds/{id}", handlers.GetThreatFeed(st))
	mux.HandleFunc("PATCH /api/v1/threat-feeds/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "threat_feeds", JSONName: "feed", Columns: []string{"name", "indicators", "confidence", "category", "ttl_hours", "status"}}))
	mux.HandleFunc("DELETE /api/v1/threat-feeds/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "threat_feeds", JSONName: "feed"}))
	mux.HandleFunc("POST /api/v1/threat-feeds/{id}/sync", handlers.SyncThreatFeed(st))
	mux.HandleFunc("GET /api/v1/threat-feeds/{id}/indicators", handlers.ListFeedIndicators(st))
	mux.HandleFunc("POST /api/v1/threat-feeds/{id}/test", handlers.TestThreatFeed(st))

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
	// Bot & abuse protection (canonical paths per API list)
	mux.HandleFunc("GET /api/v1/bot-policies", handlers.ListResource(st, handlers.CRUDConfig{Table: "bot_policies", JSONName: "bot policy"}))
	mux.HandleFunc("POST /api/v1/bot-policies", handlers.CreateResource(st, handlers.CRUDConfig{Table: "bot_policies", JSONName: "bot policy", Required: []string{"name"}, Columns: []string{"application_id", "name", "challenge_type", "known_bots", "automation_signals", "login_protection", "scrape_protection", "status"}}))
	mux.HandleFunc("GET /api/v1/bot-policies/{id}", handlers.GetBotPolicy(st))
	mux.HandleFunc("PATCH /api/v1/bot-policies/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "bot_policies", JSONName: "bot policy", Columns: []string{"challenge_type", "known_bots", "automation_signals", "login_protection", "scrape_protection", "status"}}))
	mux.HandleFunc("DELETE /api/v1/bot-policies/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "bot_policies", JSONName: "bot policy"}))
	mux.HandleFunc("GET /api/v1/bot/events", handlers.ListBotEvents(st))
	mux.HandleFunc("GET /api/v1/bot/clients", handlers.ListBotClients(st))
	mux.HandleFunc("POST /api/v1/bot/challenges/{id}/revoke", handlers.RevokeBotChallenge(st))

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

	// Gateway Clusters (aliased to clusters table, HA fleet management)
	mux.HandleFunc("GET /api/v1/gateway-clusters", handlers.ListResource(st, handlers.CRUDConfig{Table: "clusters", JSONName: "cluster"}))
	mux.HandleFunc("POST /api/v1/gateway-clusters", handlers.CreateResource(st, handlers.CRUDConfig{Table: "clusters", JSONName: "cluster", Required: []string{"name"}, Columns: []string{"name", "site", "gateway_ids", "ha_enabled", "status"}}))
	mux.HandleFunc("GET /api/v1/gateway-clusters/{id}", handlers.GetGatewayCluster(st))
	mux.HandleFunc("PATCH /api/v1/gateway-clusters/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "clusters", JSONName: "cluster", Columns: []string{"name", "site", "gateway_ids", "ha_enabled", "status"}}))
	mux.HandleFunc("DELETE /api/v1/gateway-clusters/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "clusters", JSONName: "cluster"}))
	mux.HandleFunc("GET /api/v1/gateway-clusters/{id}/gateways", handlers.ClusterGateways(st))
	mux.HandleFunc("POST /api/v1/gateway-clusters/{id}/deploy", handlers.DeployClusterConfig(st))
	mux.HandleFunc("POST /api/v1/gateway-clusters/{id}/rollback", handlers.RollbackClusterConfig(st))
	mux.HandleFunc("GET /api/v1/gateway-clusters/{id}/health", handlers.ClusterHealth(st))

	// Phase 4: incidents
	mux.HandleFunc("GET /api/v1/incidents", handlers.ListIncidents(st))
	mux.HandleFunc("POST /api/v1/incidents", handlers.CreateIncident(st))
	mux.HandleFunc("GET /api/v1/incidents/{id}", handlers.GetIncident(st))
	mux.HandleFunc("PATCH /api/v1/incidents/{id}", handlers.UpdateIncident(st))
	mux.HandleFunc("DELETE /api/v1/incidents/{id}", handlers.DeleteIncident(st))
	mux.HandleFunc("POST /api/v1/incidents/{id}/assign", handlers.AssignIncident(st))
	mux.HandleFunc("POST /api/v1/incidents/{id}/escalate", handlers.EscalateIncident(st))
	mux.HandleFunc("POST /api/v1/incidents/{id}/close", handlers.CloseIncident(st))
	mux.HandleFunc("POST /api/v1/incidents/{id}/reopen", handlers.ReopenIncident(st))
	mux.HandleFunc("GET /api/v1/incidents/{id}/events", handlers.IncidentEvents(st))
	mux.HandleFunc("GET /api/v1/incidents/{id}/timeline", handlers.IncidentTimeline(st))

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

	// Phase 4: organizations / workspaces / tenants
	mux.HandleFunc("GET /api/v1/organizations", handlers.ListResource(st, handlers.CRUDConfig{Table: "organizations", JSONName: "organization"}))
	mux.HandleFunc("POST /api/v1/organizations", handlers.CreateResource(st, handlers.CRUDConfig{Table: "organizations", JSONName: "organization", Required: []string{"name"}, Columns: []string{"name", "description", "contact_email", "plan", "quotas", "status"}}))
	mux.HandleFunc("GET /api/v1/organizations/{id}", handlers.GetResource(st, handlers.CRUDConfig{Table: "organizations", JSONName: "organization"}))
	mux.HandleFunc("PATCH /api/v1/organizations/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "organizations", JSONName: "organization", Columns: []string{"name", "description", "contact_email", "plan", "quotas", "status"}}))
	mux.HandleFunc("DELETE /api/v1/organizations/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "organizations", JSONName: "organization"}))

	mux.HandleFunc("GET /api/v1/tenants", handlers.ListResource(st, handlers.CRUDConfig{Table: "tenants", JSONName: "tenant"}))
	mux.HandleFunc("POST /api/v1/tenants", handlers.CreateResource(st, handlers.CRUDConfig{Table: "tenants", JSONName: "tenant", Required: []string{"name"}, Columns: []string{"organization_id", "name", "description", "contact_email", "plan", "quotas", "status"}}))
	mux.HandleFunc("GET /api/v1/tenants/{id}", handlers.GetResource(st, handlers.CRUDConfig{Table: "tenants", JSONName: "tenant"}))
	mux.HandleFunc("PATCH /api/v1/tenants/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "tenants", JSONName: "tenant", Columns: []string{"name", "description", "contact_email", "plan", "quotas", "status"}}))
	mux.HandleFunc("DELETE /api/v1/tenants/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "tenants", JSONName: "tenant"}))

	// Listeners / Virtual Servers (aliased to virtual_servers table)
	mux.HandleFunc("GET /api/v1/listeners", handlers.ListResource(st, handlers.CRUDConfig{Table: "virtual_servers", JSONName: "listener"}))
	mux.HandleFunc("POST /api/v1/listeners", handlers.CreateResource(st, handlers.CRUDConfig{Table: "virtual_servers", JSONName: "listener", Required: []string{"name", "listen_port"}, Columns: []string{"name", "listen_addr", "listen_port", "tls_enabled", "certificate_ref", "default_backend_pool_id", "status"}}))
	mux.HandleFunc("GET /api/v1/listeners/{id}", handlers.GetResource(st, handlers.CRUDConfig{Table: "virtual_servers", JSONName: "listener"}))
	mux.HandleFunc("PATCH /api/v1/listeners/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "virtual_servers", JSONName: "listener", Columns: []string{"name", "listen_addr", "listen_port", "tls_enabled", "certificate_ref", "default_backend_pool_id", "status"}}))
	mux.HandleFunc("DELETE /api/v1/listeners/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "virtual_servers", JSONName: "listener"}))
	mux.HandleFunc("POST /api/v1/listeners/{id}/enable", handlers.EnableListener(st))
	mux.HandleFunc("POST /api/v1/listeners/{id}/disable", handlers.DisableListener(st))

	// Sites / data centers
	mux.HandleFunc("GET /api/v1/sites", handlers.ListResource(st, handlers.CRUDConfig{Table: "sites", JSONName: "site"}))
	mux.HandleFunc("POST /api/v1/sites", handlers.CreateResource(st, handlers.CRUDConfig{Table: "sites", JSONName: "site", Required: []string{"name"}, Columns: []string{"name", "description", "location", "country_code", "gateway_ids", "status"}}))
	mux.HandleFunc("GET /api/v1/sites/{id}", handlers.GetSite(st))
	mux.HandleFunc("PATCH /api/v1/sites/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "sites", JSONName: "site", Columns: []string{"name", "description", "location", "country_code", "gateway_ids", "status"}}))
	mux.HandleFunc("DELETE /api/v1/sites/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "sites", JSONName: "site"}))
	mux.HandleFunc("GET /api/v1/sites/{id}/health", handlers.SiteHealth(st))
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
	mux.HandleFunc("GET /api/v1/notification-channels", handlers.ListNotificationChannels(st))
	mux.HandleFunc("POST /api/v1/notification-channels", handlers.CreateNotificationChannel(st))
	mux.HandleFunc("GET /api/v1/notification-channels/{id}", handlers.GetNotificationChannel(st))
	mux.HandleFunc("PATCH /api/v1/notification-channels/{id}", handlers.UpdateNotificationChannel(st))
	mux.HandleFunc("DELETE /api/v1/notification-channels/{id}", handlers.TestNotificationChannel(st))
	mux.HandleFunc("POST /api/v1/notification-channels/{id}/test", handlers.TestNotificationChannel(st))
	mux.HandleFunc("GET /api/v1/integrations", handlers.ListIntegrations(st))
	mux.HandleFunc("POST /api/v1/integrations", handlers.CreateResource(st, handlers.CRUDConfig{Table: "integrations", JSONName: "integration", Required: []string{"type", "name"}, Columns: []string{"type", "name", "endpoint", "log_types", "enabled", "config"}}))
	mux.HandleFunc("GET /api/v1/integrations/{id}", handlers.GetIntegration(st))
	mux.HandleFunc("PATCH /api/v1/integrations/{id}", handlers.UpdateResource(st, handlers.CRUDConfig{Table: "integrations", JSONName: "integration", Columns: []string{"endpoint", "log_types", "enabled", "config"}}))
	mux.HandleFunc("DELETE /api/v1/integrations/{id}", handlers.DeleteResource(st, handlers.CRUDConfig{Table: "integrations", JSONName: "integration"}))
	mux.HandleFunc("POST /api/v1/integrations/{id}/test", handlers.TestIntegration(st))
	mux.HandleFunc("POST /api/v1/integrations/{id}/enable", handlers.EnableIntegration(st))
	mux.HandleFunc("POST /api/v1/integrations/{id}/disable", handlers.DisableIntegration(st))
	mux.HandleFunc("GET /api/v1/groups", handlers.ListGroups(st))
	mux.HandleFunc("POST /api/v1/groups", handlers.CreateGroup(st))
	mux.HandleFunc("GET /api/v1/groups/{id}", handlers.GetGroup(st))
	mux.HandleFunc("PATCH /api/v1/groups/{id}", handlers.UpdateGroup(st))
	mux.HandleFunc("DELETE /api/v1/groups/{id}", handlers.DeleteGroup(st))
	mux.HandleFunc("GET /api/v1/roles", handlers.ListRoles(st))
	mux.HandleFunc("GET /api/v1/roles/{id}", handlers.GetRole(st))
	mux.HandleFunc("GET /api/v1/permissions", handlers.ListRoles(st))
	mux.HandleFunc("POST /api/v1/users/{id}/roles", handlers.AssignRoleToUser(st))
	mux.HandleFunc("DELETE /api/v1/users/{id}/roles/{roleId}", handlers.RemoveRoleFromUser(st))
	mux.HandleFunc("GET /api/v1/learning/sessions", handlers.ListLearningSessions(st))
	mux.HandleFunc("POST /api/v1/learning/sessions", handlers.CreateLearningSession(st))
	mux.HandleFunc("GET /api/v1/learning/sessions/{id}", handlers.GetLearningSession(st))
	mux.HandleFunc("POST /api/v1/learning/sessions/{id}/start", handlers.LearnSessionStartStop(st))
	mux.HandleFunc("POST /api/v1/learning/sessions/{id}/stop", handlers.LearnSessionStartStop(st))
	mux.HandleFunc("GET /api/v1/learning/suggestions", handlers.ListLearningSuggestions(st))
	mux.HandleFunc("GET /api/v1/learning/suggestions/{id}", handlers.GetLearningSuggestion(st))
	mux.HandleFunc("POST /api/v1/learning/suggestions/{id}/accept", handlers.AcceptLearningSuggestion(st))
	mux.HandleFunc("POST /api/v1/learning/suggestions/{id}/reject", handlers.RejectLearningSuggestion(st))

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
