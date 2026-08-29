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
	mux.HandleFunc("POST /api/v1/policy-versions/{id}/rollback", handlers.RollbackPolicyVersion(st))

	// Apply middleware stack (outermost first).
	var h http.Handler = mux
	h = middleware.Logging(h)
	h = middleware.RequestID(h)
	h = middleware.Audit(st.Pool, h)
	h = middleware.Recovery(h)

	return h
}