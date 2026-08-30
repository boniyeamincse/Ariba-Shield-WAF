package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

type application struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	OwnerUserID string `json:"owner_user_id,omitempty"`
	Version     int64  `json:"version"`
}

// ListApplications returns all applications visible to the caller.
func ListApplications(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, description, status, COALESCE(owner_user_id, ''), version FROM applications ORDER BY name`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var apps []application
		for rows.Next() {
			var a application
			if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.Status, &a.OwnerUserID, &a.Version); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			apps = append(apps, a)
		}
		if apps == nil {
			apps = []application{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(apps)
	}
}

// CreateApplication creates a new application with its primary domain and
// origin, plus WAF/TLS/health-check settings from the onboarding wizard.
func CreateApplication(st *store.Store) http.HandlerFunc {
	type applicationCreate struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Environment string   `json:"environment"` // production | staging | development
		Status      string   `json:"status"`      // active | disabled
		Tags        []string `json:"tags"`

		// Domain & origin
		Domain              string `json:"domain"`
		OriginType          string `json:"origin_type"`          // ip | hostname | load_balancer
		OriginHost          string `json:"origin_host"`
		OriginPort          int    `json:"origin_port"`
		OriginProtocol      string `json:"origin_protocol"`      // http | https
		OriginPath          string `json:"origin_path"`
		OriginLoadBalancing string `json:"origin_load_balancing"`

		// WAF policy
		WAFPolicyID string `json:"waf_policy_id"`
		WAFMode     string `json:"waf_mode"` // block | detection | disabled

		// TLS
		TLSEnabled     bool   `json:"tls_enabled"`
		CertificateID  string `json:"certificate_id"`
		MinTLSVersion  string `json:"min_tls_version"`
		HTTPRedirect   bool   `json:"http_redirect"`

		// Rate limit
		RateLimitEnabled bool `json:"rate_limit_enabled"`
		RateLimit        int  `json:"rate_limit"` // requests / minute

		// Health check
		HealthCheckEnabled         bool `json:"health_check_enabled"`
		HealthCheckMethod          string `json:"health_check_method"`
		HealthCheckPath            string `json:"health_check_path"`
		HealthCheckInterval        int  `json:"health_check_interval"`
		HealthCheckTimeout         int  `json:"health_check_timeout"`
		HealthCheckRetries         int  `json:"health_check_retries"`
		HealthCheckExpectedStatus  int  `json:"health_check_expected_status"`

		// Advanced
		RequestBodyLimitMB int    `json:"request_body_limit_mb"`
		ConnectionTimeoutS int    `json:"connection_timeout_s"`
		KeepAlive          bool   `json:"keep_alive"`
		RealClientIPHeader string `json:"real_client_ip_header"`
		LogRequestHeaders  bool   `json:"log_request_headers"`
		LogResponseStatus  bool   `json:"log_response_status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body applicationCreate
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Name == "" {
			http.Error(w, `{"error":"name is required"}`, http.StatusBadRequest)
			return
		}
		if body.Environment == "" {
			body.Environment = "production"
		}
		if body.Status == "" {
			body.Status = "active"
		}
		if body.OriginProtocol == "" {
			body.OriginProtocol = "https"
		}
		if body.OriginPort == 0 {
			body.OriginPort = 443
		}
		if body.OriginPath == "" {
			body.OriginPath = "/"
		}
		if body.WAFMode == "" {
			body.WAFMode = "block"
		}
		if body.MinTLSVersion == "" {
			body.MinTLSVersion = "1.2"
		}
		if body.HealthCheckMethod == "" {
			body.HealthCheckMethod = "GET"
		}
		if body.HealthCheckPath == "" {
			body.HealthCheckPath = "/health"
		}
		if body.HealthCheckInterval == 0 {
			body.HealthCheckInterval = 30
		}
		if body.HealthCheckTimeout == 0 {
			body.HealthCheckTimeout = 5
		}
		if body.HealthCheckRetries == 0 {
			body.HealthCheckRetries = 3
		}
		if body.HealthCheckExpectedStatus == 0 {
			body.HealthCheckExpectedStatus = 200
		}
		if body.RateLimit == 0 {
			body.RateLimit = 1000
		}
		if body.RequestBodyLimitMB == 0 {
			body.RequestBodyLimitMB = 10
		}
		if body.ConnectionTimeoutS == 0 {
			body.ConnectionTimeoutS = 30
		}
		if body.RealClientIPHeader == "" {
			body.RealClientIPHeader = "X-Forwarded-For"
		}
		if body.OriginLoadBalancing == "" {
			body.OriginLoadBalancing = "single"
		}

		appID, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}
		// Single-org for 0.1.
		orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"

		// Wrap app + domain + origin in one transaction.
		tx, err := st.Pool.Begin(r.Context())
		if err != nil {
			http.Error(w, `{"error":"tx begin failed"}`, http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(r.Context())

		if _, err := tx.Exec(r.Context(),
			`INSERT INTO applications (
			   id, organization_id, name, description, environment, status, tags,
			   domain, origin_type, origin_host, origin_port, origin_protocol, origin_path, origin_load_balancing,
			   waf_policy_id, waf_mode,
			   tls_enabled, certificate_id, min_tls_version, http_redirect,
			   rate_limit_enabled, rate_limit,
			   health_check_enabled, health_check_method, health_check_path,
			   health_check_interval, health_check_timeout, health_check_retries, health_check_expected_status,
			   request_body_limit_mb, connection_timeout_s, keep_alive,
			   real_client_ip_header, log_request_headers, log_response_status
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35)`,
			appID, orgID, body.Name, body.Description, body.Environment, body.Status, body.Tags,
			body.Domain, body.OriginType, body.OriginHost, body.OriginPort, body.OriginProtocol, body.OriginPath, body.OriginLoadBalancing,
			nullIfEmpty(body.WAFPolicyID), body.WAFMode,
			body.TLSEnabled, nullIfEmpty(body.CertificateID), body.MinTLSVersion, body.HTTPRedirect,
			body.RateLimitEnabled, body.RateLimit,
			body.HealthCheckEnabled, body.HealthCheckMethod, body.HealthCheckPath,
			body.HealthCheckInterval, body.HealthCheckTimeout, body.HealthCheckRetries, body.HealthCheckExpectedStatus,
			body.RequestBodyLimitMB, body.ConnectionTimeoutS, body.KeepAlive,
			body.RealClientIPHeader, body.LogRequestHeaders, body.LogResponseStatus); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		// Create the primary domain row.
		if body.Domain != "" {
			if _, err := tx.Exec(r.Context(),
				`INSERT INTO domains (id, organization_id, application_id, hostname, enabled)
				 VALUES ($1, $2, $3, $4, true)`,
				mustNewID(st), orgID, appID, body.Domain); err != nil {
				http.Error(w, `{"error":"domain insert failed"}`, http.StatusInternalServerError)
				return
			}
		}

		// Create the primary origin row.
		if body.OriginHost != "" {
			if _, err := tx.Exec(r.Context(),
				`INSERT INTO origins (id, organization_id, application_id, name, protocol, host, port, weight)
				 VALUES ($1, $2, $3, $4, $5, $6, $7, 1)`,
				mustNewID(st), orgID, appID, "primary", body.OriginProtocol, body.OriginHost, body.OriginPort); err != nil {
				http.Error(w, `{"error":"origin insert failed"}`, http.StatusInternalServerError)
				return
			}
		}

		if err := tx.Commit(r.Context()); err != nil {
			http.Error(w, `{"error":"tx commit failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": appID})
	}
}

// UpdateApplication updates an application (PATCH semantics).
func UpdateApplication(st *store.Store) http.HandlerFunc {
	type applicationUpdate struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Status      *string `json:"status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		var body applicationUpdate
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}

		// Verify exists + read current values for partial update.
		var name, description, status string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT name, COALESCE(description,''), status FROM applications WHERE id = $1`, appID).
			Scan(&name, &description, &status); err != nil {
			http.Error(w, `{"error":"application not found"}`, http.StatusNotFound)
			return
		}

		if body.Name != nil {
			name = *body.Name
		}
		if body.Description != nil {
			description = *body.Description
		}
		if body.Status != nil {
			if *body.Status != "active" && *body.Status != "disabled" {
				http.Error(w, `{"error":"status must be active or disabled"}`, http.StatusBadRequest)
				return
			}
			status = *body.Status
		}

		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE applications SET name = $1, description = $2, status = $3, version = version + 1, updated_at = now() WHERE id = $4`,
			name, description, status, appID); err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": appID, "name": name, "status": status})
	}
}

// DeleteApplication deletes an application.
func DeleteApplication(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM applications WHERE id = $1`, appID)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"application not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}