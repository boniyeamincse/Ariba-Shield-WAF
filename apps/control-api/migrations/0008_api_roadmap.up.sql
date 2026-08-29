-- 0008_api_roadmap.up.sql
-- Completes the API master plan (endpoint.md) modules Phase 2-5.

-- Phase 2: config-validation (dry-run)
CREATE TABLE IF NOT EXISTS config_validations (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    policy_version_id TEXT REFERENCES policy_versions(id),
    document        JSONB,
    result          TEXT NOT NULL DEFAULT 'pending', -- pending | valid | invalid
    errors          JSONB,
    created_by      TEXT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Phase 3: threat intelligence
CREATE TABLE IF NOT EXISTS threat_feeds (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    source          TEXT NOT NULL,
    indicator_type  TEXT NOT NULL DEFAULT 'ip', -- ip | domain | url | asn
    indicators      JSONB NOT NULL DEFAULT '[]',
    confidence      TEXT NOT NULL DEFAULT 'low', -- low | medium | high
    category        TEXT NOT NULL DEFAULT '',
    ttl_hours       INTEGER NOT NULL DEFAULT 24,
    provenance      TEXT NOT NULL DEFAULT '',
    status          TEXT NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Phase 3: API security
CREATE TABLE IF NOT EXISTS api_schemas (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    application_id  TEXT REFERENCES applications(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    openapi_document JSONB,
    path            TEXT NOT NULL DEFAULT '',
    method          TEXT NOT NULL DEFAULT '',
    status          TEXT NOT NULL DEFAULT 'active',
    drift_alert     BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Phase 3: bot management
CREATE TABLE IF NOT EXISTS bot_policies (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    application_id  TEXT REFERENCES applications(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    challenge_type  TEXT NOT NULL DEFAULT 'javascript', -- javascript | captcha | proof_of_work
    known_bots      TEXT NOT NULL DEFAULT 'allow',
    automation_signals BOOLEAN NOT NULL DEFAULT true,
    login_protection   BOOLEAN NOT NULL DEFAULT false,
    scrape_protection  BOOLEAN NOT NULL DEFAULT true,
    status          TEXT NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Phase 3: DLP
CREATE TABLE IF NOT EXISTS dlp_profiles (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    scan_targets    TEXT[] NOT NULL DEFAULT '{}',
    rules           JSONB NOT NULL DEFAULT '[]',
    status          TEXT NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Phase 3: integrations (SIEM/log forwarding)
CREATE TABLE IF NOT EXISTS integrations (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    type            TEXT NOT NULL,           -- splunk_hec | wazuh | syslog | webhook | teams | slack
    name            TEXT NOT NULL,
    endpoint        TEXT NOT NULL DEFAULT '',
    log_types       TEXT[] NOT NULL DEFAULT '{}',
    enabled         BOOLEAN NOT NULL DEFAULT true,
    config          JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Phase 3: IAM / SSO
CREATE TABLE IF NOT EXISTS iam_sso (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    provider_name   TEXT NOT NULL,
    provider_type   TEXT NOT NULL DEFAULT 'saml', -- saml | oidc
    idp_entity_id   TEXT NOT NULL DEFAULT '',
    sso_url         TEXT NOT NULL DEFAULT '',
    config          JSONB NOT NULL DEFAULT '{}',
    enabled         BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Phase 3: service accounts + API keys + secrets
CREATE TABLE IF NOT EXISTS service_accounts (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    roles           TEXT[] NOT NULL DEFAULT '{}',
    status          TEXT NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS api_keys (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    service_account_id TEXT REFERENCES service_accounts(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    key_prefix      TEXT NOT NULL,
    key_hash        TEXT NOT NULL,           -- never store the full key
    expires_at      TIMESTAMPTZ,
    last_used_at    TIMESTAMPTZ,
    status          TEXT NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS secrets (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    secret_ref      TEXT NOT NULL,           -- reference to KMS/Vault, never the value
    provider        TEXT NOT NULL DEFAULT 'envelope',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Phase 4: analytics + incidents + automation + clusters + caching
CREATE TABLE IF NOT EXISTS incidents (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    title           TEXT NOT NULL,
    severity        TEXT NOT NULL DEFAULT 'medium',
    status          TEXT NOT NULL DEFAULT 'open', -- open | investigating | resolved
    owner_user_id   TEXT REFERENCES users(id),
    notes           TEXT NOT NULL DEFAULT '',
    related_events  TEXT[] NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS automation_rules (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    trigger_event   TEXT NOT NULL,
    action          TEXT NOT NULL DEFAULT 'notify', -- notify | block_ip | escalate | webhook
    target          TEXT NOT NULL DEFAULT '',
    enabled         BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS clusters (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    site            TEXT NOT NULL DEFAULT '',
    gateway_ids     TEXT[] NOT NULL DEFAULT '{}',
    ha_enabled      BOOLEAN NOT NULL DEFAULT false,
    status          TEXT NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS caching_rules (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    application_id  TEXT REFERENCES applications(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    path_pattern    TEXT NOT NULL DEFAULT '',
    ttl_seconds     INTEGER NOT NULL DEFAULT 60,
    cache_methods   TEXT[] NOT NULL DEFAULT '{}',
    status          TEXT NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Phase 5: graphql security, CSP, quotas, ML baselines, network protection
CREATE TABLE IF NOT EXISTS graphql_security (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    application_id  TEXT REFERENCES applications(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    max_depth       INTEGER NOT NULL DEFAULT 10,
    max_complexity  INTEGER NOT NULL DEFAULT 1000,
    introspection   TEXT NOT NULL DEFAULT 'deny',
    status          TEXT NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS csp_profiles (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    application_id  TEXT REFERENCES applications(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    csp_header      TEXT NOT NULL DEFAULT '',
    inject_mode     TEXT NOT NULL DEFAULT 'header', -- header | meta
    status          TEXT NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS api_quotas (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    application_id  TEXT REFERENCES applications(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    match_by        TEXT NOT NULL DEFAULT 'client_ip', -- client_ip | jwt_claim | api_key
    limit_count     INTEGER NOT NULL,
    window_seconds  INTEGER NOT NULL DEFAULT 60,
    status          TEXT NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ml_baselines (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    application_id  TEXT REFERENCES applications(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'training', -- training | active | inactive
    config          JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS network_protection (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    protection_type TEXT NOT NULL DEFAULT 'bgp_route', -- bgp_route | anycast | scrubbing
    config          JSONB NOT NULL DEFAULT '{}',
    status          TEXT NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
