package api

import (
	"net/http"

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
	mux.HandleFunc("GET /api/v1/applications/{id}/domains", handlers.ListDomains(st))
	mux.HandleFunc("POST /api/v1/applications/{id}/domains", handlers.CreateDomain(st))
	mux.HandleFunc("GET /api/v1/applications/{id}/origins", handlers.ListOrigins(st))
	mux.HandleFunc("POST /api/v1/applications/{id}/origins", handlers.CreateOrigin(st))
	mux.HandleFunc("GET /api/v1/security-policies", handlers.ListSecurityPolicies(st))
	mux.HandleFunc("POST /api/v1/security-policies", handlers.CreateSecurityPolicy(st))

	// Gateway fleet (Phase 2 operations).
	mux.HandleFunc("POST /api/v1/gateways/register", handlers.RegisterGateway(st))
	mux.HandleFunc("POST /api/v1/gateways/{id}/heartbeat", handlers.Heartbeat(st))
	mux.HandleFunc("GET /api/v1/gateways", handlers.ListGateways(st))

	// Security events (Phase 2).
	mux.HandleFunc("GET /api/v1/security-events", handlers.ListSecurityEvents(st))

	// Policy binding (Phase 2).
	mux.HandleFunc("POST /api/v1/security-policies/bind", handlers.BindPolicy(st))

	// Phase 3 — safe blocking: IP lists, rate limits, policy versions.
	mux.HandleFunc("GET /api/v1/ip-lists", handlers.ListIPLists(st))
	mux.HandleFunc("POST /api/v1/ip-lists", handlers.CreateIPList(st))
	mux.HandleFunc("GET /api/v1/rate-limits", handlers.ListRateLimits(st))
	mux.HandleFunc("POST /api/v1/rate-limits", handlers.CreateRateLimit(st))
	mux.HandleFunc("POST /api/v1/security-policies/{id}/versions", handlers.CreatePolicyVersion(st))
	mux.HandleFunc("POST /api/v1/policy-versions/{id}/activate", handlers.ActivatePolicyVersion(st))
	mux.HandleFunc("POST /api/v1/policy-versions/{id}/promote", handlers.PromotePolicyVersion(st))
	mux.HandleFunc("POST /api/v1/policy-versions/{id}/rollback", handlers.RollbackPolicyVersion(st))
	mux.HandleFunc("GET /api/v1/policy-versions/diff", handlers.DiffPolicyVersions(st))

	// Auth & Identity
	mux.HandleFunc("POST /api/v1/auth/login", handlers.Login(st))
	mux.HandleFunc("POST /api/v1/auth/logout", handlers.Logout(st))
	mux.HandleFunc("GET /api/v1/auth/me", handlers.Me(st))

	// Phase 3 — Webhooks, Exceptions, Rules, Certificates, Deployments
	mux.HandleFunc("GET /api/v1/webhooks", handlers.ListWebhooks(st))
	mux.HandleFunc("POST /api/v1/webhooks", handlers.CreateWebhook(st))
	mux.HandleFunc("GET /api/v1/exceptions", handlers.ListExceptions(st))
	mux.HandleFunc("POST /api/v1/exceptions", handlers.CreateException(st))
	mux.HandleFunc("GET /api/v1/managed-rules", handlers.ListManagedRules(st))
	mux.HandleFunc("POST /api/v1/managed-rules", handlers.ConfigureManagedRules(st))
	mux.HandleFunc("GET /api/v1/custom-rules", handlers.ListCustomRules(st))
	mux.HandleFunc("POST /api/v1/custom-rules", handlers.CreateCustomRule(st))
	mux.HandleFunc("GET /api/v1/deployments", handlers.ListDeployments(st))
	mux.HandleFunc("POST /api/v1/deployments", handlers.SyncDeployment(st))
	mux.HandleFunc("GET /api/v1/certificates", handlers.ListCertificates(st))
	mux.HandleFunc("POST /api/v1/certificates", handlers.UploadCertificate(st))

	// Apply middleware stack. Order matters (P0.6): Auth and RequestID must set
	// the request context BEFORE Audit reads it, otherwise the audit trail has
	// empty actor/request_id. Execution order (outermost first):
	//   Recovery -> RBACEnforcer -> RequestID -> Auth -> Audit -> Logging -> mux
	var h http.Handler = mux
	h = middleware.Logging(h)
	h = middleware.Audit(st.Pool, h)
	h = middleware.Auth(st.Pool)(h)
	h = middleware.RequestID(h)
	h = middleware.RBACEnforcer(h)
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