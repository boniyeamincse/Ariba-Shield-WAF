-- 0023_app_wizard.down.sql
-- Revert application onboarding fields.

ALTER TABLE applications DROP COLUMN IF EXISTS log_response_status;
ALTER TABLE applications DROP COLUMN IF EXISTS log_request_headers;
ALTER TABLE applications DROP COLUMN IF EXISTS real_client_ip_header;
ALTER TABLE applications DROP COLUMN IF EXISTS keep_alive;
ALTER TABLE applications DROP COLUMN IF EXISTS connection_timeout_s;
ALTER TABLE applications DROP COLUMN IF EXISTS request_body_limit_mb;
ALTER TABLE applications DROP COLUMN IF EXISTS health_check_expected_status;
ALTER TABLE applications DROP COLUMN IF EXISTS health_check_retries;
ALTER TABLE applications DROP COLUMN IF EXISTS health_check_timeout;
ALTER TABLE applications DROP COLUMN IF EXISTS health_check_interval;
ALTER TABLE applications DROP COLUMN IF EXISTS health_check_path;
ALTER TABLE applications DROP COLUMN IF EXISTS health_check_method;
ALTER TABLE applications DROP COLUMN IF EXISTS health_check_enabled;
ALTER TABLE applications DROP COLUMN IF EXISTS rate_limit;
ALTER TABLE applications DROP COLUMN IF EXISTS rate_limit_enabled;
ALTER TABLE applications DROP COLUMN IF EXISTS http_redirect;
ALTER TABLE applications DROP COLUMN IF EXISTS min_tls_version;
ALTER TABLE applications DROP COLUMN IF EXISTS certificate_id;
ALTER TABLE applications DROP COLUMN IF EXISTS tls_enabled;
ALTER TABLE applications DROP COLUMN IF EXISTS waf_mode;
ALTER TABLE applications DROP COLUMN IF EXISTS waf_policy_id;
ALTER TABLE applications DROP COLUMN IF EXISTS origin_load_balancing;
ALTER TABLE applications DROP COLUMN IF EXISTS origin_path;
ALTER TABLE applications DROP COLUMN IF EXISTS origin_protocol;
ALTER TABLE applications DROP COLUMN IF EXISTS origin_port;
ALTER TABLE applications DROP COLUMN IF EXISTS origin_host;
ALTER TABLE applications DROP COLUMN IF EXISTS origin_type;
ALTER TABLE applications DROP COLUMN IF EXISTS domain;
ALTER TABLE applications DROP COLUMN IF EXISTS tags;
ALTER TABLE applications DROP COLUMN IF EXISTS environment;
