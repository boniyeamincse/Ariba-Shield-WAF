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
	Environment string `json:"environment"`
	Status      string `json:"status"`
	OwnerUserID string `json:"owner_user_id,omitempty"`
	Version     int64  `json:"version"`
	Tags        []string `json:"tags,omitempty"`

	Domain              string `json:"domain,omitempty"`
	OriginType          string `json:"origin_type,omitempty"`
	OriginHost          string `json:"origin_host,omitempty"`
	OriginPort          int    `json:"origin_port,omitempty"`
	OriginProtocol      string `json:"origin_protocol,omitempty"`
	OriginPath          string `json:"origin_path,omitempty"`
	OriginLoadBalancing string `json:"origin_load_balancing,omitempty"`

	WAFPolicyID string `json:"waf_policy_id,omitempty"`
	WAFMode     string `json:"waf_mode,omitempty"`

	TLSEnabled    bool   `json:"tls_enabled"`
	CertificateID string `json:"certificate_id,omitempty"`
	MinTLSVersion string `json:"min_tls_version"`
	HTTPRedirect  bool   `json:"http_redirect"`

	RateLimitEnabled bool `json:"rate_limit_enabled"`
	RateLimit        int  `json:"rate_limit"`

	HealthCheckEnabled       bool `json:"health_check_enabled"`
	HealthCheckMethod        string `json:"health_check_method"`
	HealthCheckPath          string `json:"health_check_path"`
	HealthCheckInterval      int    `json:"health_check_interval"`
	HealthCheckTimeout       int    `json:"health_check_timeout"`
	HealthCheckRetries       int    `json:"health_check_retries"`
	HealthCheckExpectedStatus int   `json:"health_check_expected_status"`

	RequestBodyLimitMB int    `json:"request_body_limit_mb"`
	ConnectionTimeoutS int    `json:"connection_timeout_s"`
	KeepAlive          bool   `json:"keep_alive"`
	RealClientIPHeader string `json:"real_client_ip_header"`
	LogRequestHeaders  bool   `json:"log_request_headers"`
	LogResponseStatus  bool   `json:"log_response_status"`
}

// appSelect is the shared column list used by list + detail queries.
const appSelect = `id, name, COALESCE(description,''), COALESCE(environment,'production'), status,
	COALESCE(owner_user_id,''), version, COALESCE(tags,'{}'),
	COALESCE(domain,''), COALESCE(origin_type,'ip'), COALESCE(origin_host,''),
	COALESCE(origin_port,0), COALESCE(origin_protocol,'https'), COALESCE(origin_path,'/'),
	COALESCE(origin_load_balancing,'single'),
	COALESCE(waf_policy_id,''), COALESCE(waf_mode,'block'),
	COALESCE(tls_enabled,false), COALESCE(certificate_id,''), COALESCE(min_tls_version,'1.2'),
	COALESCE(http_redirect,false),
	COALESCE(rate_limit_enabled,false), COALESCE(rate_limit,0),
	COALESCE(health_check_enabled,false), COALESCE(health_check_method,'GET'),
	COALESCE(health_check_path,'/health'), COALESCE(health_check_interval,30),
	COALESCE(health_check_timeout,5), COALESCE(health_check_retries,3),
	COALESCE(health_check_expected_status,200),
	COALESCE(request_body_limit_mb,10), COALESCE(connection_timeout_s,30),
	COALESCE(keep_alive,true), COALESCE(real_client_ip_header,'X-Forwarded-For'),
	COALESCE(log_request_headers,true), COALESCE(log_response_status,true)`

func scanApplication(row interface{ Scan(...any) error }) (*application, error) {
	var a application
	err := row.Scan(
		&a.ID, &a.Name, &a.Description, &a.Environment, &a.Status,
		&a.OwnerUserID, &a.Version, &a.Tags,
		&a.Domain, &a.OriginType, &a.OriginHost, &a.OriginPort, &a.OriginProtocol, &a.OriginPath, &a.OriginLoadBalancing,
		&a.WAFPolicyID, &a.WAFMode,
		&a.TLSEnabled, &a.CertificateID, &a.MinTLSVersion, &a.HTTPRedirect,
		&a.RateLimitEnabled, &a.RateLimit,
		&a.HealthCheckEnabled, &a.HealthCheckMethod, &a.HealthCheckPath,
		&a.HealthCheckInterval, &a.HealthCheckTimeout, &a.HealthCheckRetries, &a.HealthCheckExpectedStatus,
		&a.RequestBodyLimitMB, &a.ConnectionTimeoutS, &a.KeepAlive,
		&a.RealClientIPHeader, &a.LogRequestHeaders, &a.LogResponseStatus,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// ListApplications returns all applications visible to the caller.
func ListApplications(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(), `SELECT `+appSelect+` FROM applications ORDER BY name`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var apps []application
		for rows.Next() {
			a, err := scanApplication(rows)
			if err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			apps = append(apps, *a)
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

// UpdateApplication updates an application (PATCH semantics). All wizard
// fields are accepted as pointers so omitted fields are left unchanged.
func UpdateApplication(st *store.Store) http.HandlerFunc {
	type applicationUpdate struct {
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		Environment *string  `json:"environment"`
		Status      *string  `json:"status"`
		Tags        *[]string `json:"tags"`

		Domain              *string `json:"domain"`
		OriginType          *string `json:"origin_type"`
		OriginHost          *string `json:"origin_host"`
		OriginPort          *int    `json:"origin_port"`
		OriginProtocol      *string `json:"origin_protocol"`
		OriginPath          *string `json:"origin_path"`
		OriginLoadBalancing *string `json:"origin_load_balancing"`

		WAFPolicyID *string `json:"waf_policy_id"`
		WAFMode     *string `json:"waf_mode"`

		TLSEnabled    *bool   `json:"tls_enabled"`
		CertificateID *string `json:"certificate_id"`
		MinTLSVersion *string `json:"min_tls_version"`
		HTTPRedirect  *bool   `json:"http_redirect"`

		RateLimitEnabled *bool `json:"rate_limit_enabled"`
		RateLimit        *int  `json:"rate_limit"`

		HealthCheckEnabled        *bool   `json:"health_check_enabled"`
		HealthCheckMethod         *string `json:"health_check_method"`
		HealthCheckPath           *string `json:"health_check_path"`
		HealthCheckInterval       *int    `json:"health_check_interval"`
		HealthCheckTimeout        *int    `json:"health_check_timeout"`
		HealthCheckRetries        *int    `json:"health_check_retries"`
		HealthCheckExpectedStatus *int    `json:"health_check_expected_status"`

		RequestBodyLimitMB *int    `json:"request_body_limit_mb"`
		ConnectionTimeoutS *int    `json:"connection_timeout_s"`
		KeepAlive          *bool   `json:"keep_alive"`
		RealClientIPHeader *string `json:"real_client_ip_header"`
		LogRequestHeaders  *bool   `json:"log_request_headers"`
		LogResponseStatus  *bool   `json:"log_response_status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		var body applicationUpdate
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}

		// Verify exists.
		var exists int
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT 1 FROM applications WHERE id = $1`, appID).Scan(&exists); err != nil {
			http.Error(w, `{"error":"application not found"}`, http.StatusNotFound)
			return
		}

		// Dynamic SET using COALESCE so nil pointers leave fields unchanged.
		// Coalesce with the existing column value for text; for non-nullable
		// boolean/int columns coalesce against their default.
		_, err := st.Pool.Exec(r.Context(),
			`UPDATE applications SET
			   name = COALESCE($2, name),
			   description = COALESCE($3, description),
			   environment = COALESCE($4, environment),
			   status = COALESCE($5, status),
			   tags = COALESCE($6, tags),
			   domain = COALESCE($7, domain),
			   origin_type = COALESCE($8, origin_type),
			   origin_host = COALESCE($9, origin_host),
			   origin_port = COALESCE($10, origin_port),
			   origin_protocol = COALESCE($11, origin_protocol),
			   origin_path = COALESCE($12, origin_path),
			   origin_load_balancing = COALESCE($13, origin_load_balancing),
			   waf_policy_id = COALESCE($14, waf_policy_id),
			   waf_mode = COALESCE($15, waf_mode),
			   tls_enabled = COALESCE($16, tls_enabled),
			   certificate_id = COALESCE($17, certificate_id),
			   min_tls_version = COALESCE($18, min_tls_version),
			   http_redirect = COALESCE($19, http_redirect),
			   rate_limit_enabled = COALESCE($20, rate_limit_enabled),
			   rate_limit = COALESCE($21, rate_limit),
			   health_check_enabled = COALESCE($22, health_check_enabled),
			   health_check_method = COALESCE($23, health_check_method),
			   health_check_path = COALESCE($24, health_check_path),
			   health_check_interval = COALESCE($25, health_check_interval),
			   health_check_timeout = COALESCE($26, health_check_timeout),
			   health_check_retries = COALESCE($27, health_check_retries),
			   health_check_expected_status = COALESCE($28, health_check_expected_status),
			   request_body_limit_mb = COALESCE($29, request_body_limit_mb),
			   connection_timeout_s = COALESCE($30, connection_timeout_s),
			   keep_alive = COALESCE($31, keep_alive),
			   real_client_ip_header = COALESCE($32, real_client_ip_header),
			   log_request_headers = COALESCE($33, log_request_headers),
			   log_response_status = COALESCE($34, log_response_status),
			   version = version + 1,
			   updated_at = now()
			 WHERE id = $1`,
			appID, body.Name, body.Description, body.Environment, body.Status, body.Tags,
			body.Domain, body.OriginType, body.OriginHost, body.OriginPort, body.OriginProtocol,
			body.OriginPath, body.OriginLoadBalancing,
			body.WAFPolicyID, body.WAFMode,
			body.TLSEnabled, body.CertificateID, body.MinTLSVersion, body.HTTPRedirect,
			body.RateLimitEnabled, body.RateLimit,
			body.HealthCheckEnabled, body.HealthCheckMethod, body.HealthCheckPath,
			body.HealthCheckInterval, body.HealthCheckTimeout, body.HealthCheckRetries, body.HealthCheckExpectedStatus,
			body.RequestBodyLimitMB, body.ConnectionTimeoutS, body.KeepAlive,
			body.RealClientIPHeader, body.LogRequestHeaders, body.LogResponseStatus)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}

		// Keep the primary domain row in sync if changed.
		if body.Domain != nil && *body.Domain != "" {
			if _, err := st.Pool.Exec(r.Context(),
				`UPDATE domains SET hostname = $1 WHERE application_id = $2 AND id = (SELECT id FROM domains WHERE application_id = $2 ORDER BY created_at ASC LIMIT 1)`,
				*body.Domain, appID); err == nil {
				// ignore: domain update is best-effort
			}
		}

		// Keep the primary origin row in sync if host/port/protocol changed.
		if (body.OriginHost != nil && *body.OriginHost != "") || body.OriginPort != nil || body.OriginProtocol != nil {
			if _, err := st.Pool.Exec(r.Context(),
				`UPDATE origins SET
				   host = COALESCE($1, host),
				   port = COALESCE($2, port),
				   protocol = COALESCE($3, protocol)
				 WHERE application_id = $4 AND id = (SELECT id FROM origins WHERE application_id = $4 ORDER BY created_at ASC LIMIT 1)`,
				body.OriginHost, body.OriginPort, body.OriginProtocol, appID); err == nil {
				// ignore: origin update is best-effort
			}
		}

		name := "application"
		if body.Name != nil {
			name = *body.Name
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": appID, "name": name})
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