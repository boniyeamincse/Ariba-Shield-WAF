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

	// Authenticated routes (placeholder auth; real auth in Sprint 3).
	mux.HandleFunc("GET /api/v1/applications", handlers.ListApplications(st))

	// Apply middleware stack (outermost first).
	var h http.Handler = mux
	h = middleware.Logging(h)
	h = middleware.RequestID(h)
	h = middleware.Recovery(h)

	return h
}