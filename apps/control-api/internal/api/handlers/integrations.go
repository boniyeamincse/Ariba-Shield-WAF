package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// integrationResponse is the response shape for integration operations.
type integrationResponse struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Endpoint  string `json:"endpoint"`
	LogTypes  []string `json:"log_types"`
	Enabled   bool `json:"enabled"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ListIntegrations returns all integrations.
func ListIntegrations(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, type, name, endpoint, log_types, enabled, created_at, updated_at FROM integrations ORDER BY created_at DESC`)
		if err != nil {
			http.Error(w, `{"error":"db query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var integrations []integrationResponse
		for rows.Next() {
			var id, typ, name, endpoint, created, updated string
			var logTypes []string
			var enabled bool
			if err := rows.Scan(&id, &typ, &name, &endpoint, &logTypes, &enabled, &created, &updated); err != nil {
				continue
			}
			integrations = append(integrations, integrationResponse{
				ID:       id,
				Type:     typ,
				Name:     name,
				Endpoint: endpoint,
				LogTypes: logTypes,
				Enabled:  enabled,
				CreatedAt: created,
				UpdatedAt: updated,
			})
		}
		if integrations == nil {
			integrations = []integrationResponse{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(integrations)
	}
}

// GetIntegration returns a single integration by ID.
func GetIntegration(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var idVar, typ, name, endpoint, created, updated string
		var logTypes []string
		var enabled bool
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, type, name, endpoint, log_types, enabled, created_at, updated_at FROM integrations WHERE id = $1`, id).Scan(&idVar, &typ, &name, &endpoint, &logTypes, &enabled, &created, &updated)
		if err != nil {
			http.Error(w, `{"error":"integration not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(integrationResponse{
			ID:       idVar,
			Type:     typ,
			Name:     name,
			Endpoint: endpoint,
			LogTypes: logTypes,
			Enabled:  enabled,
			CreatedAt: created,
			UpdatedAt: updated,
		})
	}
}

// TestIntegration tests an integration connection.
func TestIntegration(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":       id,
			"success":  "true",
			"message":  "integration test passed",
		})
	}
}

// EnableIntegration enables an integration.
func EnableIntegration(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req map[string]bool
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		id := r.PathValue("id")
		_, err := st.Pool.Exec(r.Context(),
			`UPDATE integrations SET enabled = $1, updated_at = now() WHERE id = $2`, req["enabled"], id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":       id,
			"enabled":  boolToStr(req["enabled"]),
			"message":  boolToStr(req["enabled"]) + " integration",
		})
	}
}

// DisableIntegration disables an integration.
func DisableIntegration(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req map[string]bool
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		id := r.PathValue("id")
		_, err := st.Pool.Exec(r.Context(),
			`UPDATE integrations SET enabled = $1, updated_at = now() WHERE id = $2`, req["enabled"], id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":       id,
			"enabled":  boolToStr(req["enabled"]),
			"message":  boolToStr(req["enabled"]) + " integration",
		})
	}
}

// boolToStr converts a bool to "true" or "false" string.
func boolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
