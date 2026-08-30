-- 0023_app_wizard.up.sql
-- Extend applications with full WAF onboarding fields (wizard).

ALTER TABLE applications ADD COLUMN IF NOT EXISTS environment        TEXT NOT NULL DEFAULT 'production'; -- production | staging | development
ALTER TABLE applications ADD COLUMN IF NOT EXISTS tags              TEXT[] NOT NULL DEFAULT '{}';

-- Domain & origin (single primary origin inline for the wizard; full
-- multi-origin pools are managed via /applications/{id}/origins).
ALTER TABLE applications ADD COLUMN IF NOT EXISTS domain            TEXT NOT NULL DEFAULT '';
ALTER TABLE applications ADD COLUMN IF NOT EXISTS origin_type       TEXT NOT NULL DEFAULT 'ip'; -- ip | hostname | load_balancer
ALTER TABLE applications ADD COLUMN IF NOT EXISTS origin_host       TEXT NOT NULL DEFAULT '';
ALTER TABLE applications ADD COLUMN IF NOT EXISTS origin_port       INTEGER NOT NULL DEFAULT 443;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS origin_protocol   TEXT NOT NULL DEFAULT 'https'; -- http | https
ALTER TABLE applications ADD COLUMN IF NOT EXISTS origin_path       TEXT NOT NULL DEFAULT '/';
ALTER TABLE applications ADD COLUMN IF NOT EXISTS origin_load_balancing TEXT NOT NULL DEFAULT 'single'; -- single | round_robin | ip_hash

-- WAF policy binding
ALTER TABLE applications ADD COLUMN IF NOT EXISTS waf_policy_id     TEXT REFERENCES security_policies(id) ON DELETE SET NULL;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS waf_mode          TEXT NOT NULL DEFAULT 'block'; -- block | detection | disabled

-- TLS / SSL
ALTER TABLE applications ADD COLUMN IF NOT EXISTS tls_enabled       BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS certificate_id    TEXT REFERENCES certificates(id) ON DELETE SET NULL;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS min_tls_version   TEXT NOT NULL DEFAULT '1.2'; -- 1.0 | 1.1 | 1.2 | 1.3
ALTER TABLE applications ADD COLUMN IF NOT EXISTS http_redirect     BOOLEAN NOT NULL DEFAULT false;

-- Rate limiting
ALTER TABLE applications ADD COLUMN IF NOT EXISTS rate_limit_enabled BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS rate_limit         INTEGER NOT NULL DEFAULT 1000; -- requests / minute

-- Health check
ALTER TABLE applications ADD COLUMN IF NOT EXISTS health_check_enabled BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS health_check_method  TEXT NOT NULL DEFAULT 'GET';
ALTER TABLE applications ADD COLUMN IF NOT EXISTS health_check_path    TEXT NOT NULL DEFAULT '/health';
ALTER TABLE applications ADD COLUMN IF NOT EXISTS health_check_interval INTEGER NOT NULL DEFAULT 30; -- seconds
ALTER TABLE applications ADD COLUMN IF NOT EXISTS health_check_timeout INTEGER NOT NULL DEFAULT 5;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS health_check_retries INTEGER NOT NULL DEFAULT 3;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS health_check_expected_status INTEGER NOT NULL DEFAULT 200;

-- Advanced
ALTER TABLE applications ADD COLUMN IF NOT EXISTS request_body_limit_mb INTEGER NOT NULL DEFAULT 10;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS connection_timeout_s  INTEGER NOT NULL DEFAULT 30;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS keep_alive            BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS real_client_ip_header TEXT NOT NULL DEFAULT 'X-Forwarded-For';
ALTER TABLE applications ADD COLUMN IF NOT EXISTS log_request_headers   BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS log_response_status  BOOLEAN NOT NULL DEFAULT true;
