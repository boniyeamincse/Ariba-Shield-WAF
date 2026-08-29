-- 0007_mock_handlers.up.sql
-- Phase 4: real DB backing for previously-mock handler resources.
-- webhooks, exceptions, certificates, custom_rules, managed_rules, deployments.

CREATE TABLE IF NOT EXISTS webhooks (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    url             TEXT NOT NULL,
    severity        TEXT[] NOT NULL DEFAULT '{}',
    event_types     TEXT[] NOT NULL DEFAULT '{}',
    custom_headers  JSONB,
    max_retry_attempts INTEGER NOT NULL DEFAULT 3,
    status          TEXT NOT NULL DEFAULT 'active',
    version         BIGINT NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS exceptions (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    policy_id       TEXT REFERENCES security_policies(id) ON DELETE CASCADE,
    application_id  TEXT REFERENCES applications(id) ON DELETE CASCADE,
    rule_id         TEXT,
    url_pattern     TEXT,
    parameter       TEXT,
    reason          TEXT NOT NULL DEFAULT '',
    expires_at      TIMESTAMPTZ,
    status          TEXT NOT NULL DEFAULT 'active',
    created_by      TEXT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS certificates (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    domain          TEXT NOT NULL,
    issuer          TEXT NOT NULL DEFAULT '',
    serial          TEXT NOT NULL DEFAULT '',
    not_before      TIMESTAMPTZ,
    not_after       TIMESTAMPTZ,
    status          TEXT NOT NULL DEFAULT 'active',  -- active | expiring | expired
    version         BIGINT NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS custom_rules (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    policy_id       TEXT REFERENCES security_policies(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    action          TEXT NOT NULL DEFAULT 'BLOCK',   -- BLOCK | LOG | ALLOW
    match_conditions JSONB,
    status          TEXT NOT NULL DEFAULT 'active',
    version         BIGINT NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS managed_rules (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    category        TEXT NOT NULL DEFAULT 'owasp-crs',
    enabled         BOOLEAN NOT NULL DEFAULT true,
    sensitivity     TEXT NOT NULL DEFAULT 'paranoia-1',
    status          TEXT NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deployments (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    policy_version_id TEXT REFERENCES policy_versions(id),
    target_gateways TEXT[] NOT NULL DEFAULT '{}',
    status          TEXT NOT NULL DEFAULT 'pending',  -- pending | syncing | active | failed
    error           TEXT,
    created_by      TEXT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
