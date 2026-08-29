-- 0001_init_schema.up.sql
-- Release 0.1 schema v0 (SRS §6 Data model v0).
-- Immutable once merged. New changes = new migration, never edit this file.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ULID generation: uses gen_random_uuid() as a stand-in until a ULID extension.
-- Public IDs are ULID (26-char Crockford base32); generated in the application.
-- This migration defines tables; IDs are app-generated ULID text values.

CREATE TABLE IF NOT EXISTS organizations (
    id          TEXT PRIMARY KEY,            -- ULID
    name        TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users (
    id             TEXT PRIMARY KEY,         -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    email          TEXT NOT NULL UNIQUE,
    password_hash  TEXT NOT NULL,
    language       TEXT NOT NULL DEFAULT 'en',
    totp_enabled   BOOLEAN NOT NULL DEFAULT false,
    status         TEXT NOT NULL DEFAULT 'active',
    version        BIGINT NOT NULL DEFAULT 1,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS groups (
    id              TEXT PRIMARY KEY,        -- ULID
    name            TEXT NOT NULL,
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, name)
);

CREATE TABLE IF NOT EXISTS roles (
    id   TEXT PRIMARY KEY,                   -- ULID
    name TEXT NOT NULL UNIQUE,
    permissions TEXT[] NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS user_group_memberships (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id TEXT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, group_id)
);

CREATE TABLE IF NOT EXISTS sessions (
    id          TEXT PRIMARY KEY,            -- ULID
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    request_id  TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions (expires_at);

CREATE TABLE IF NOT EXISTS audit_events (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    actor_user_id   TEXT REFERENCES users(id),
    action          TEXT NOT NULL,
    resource_type   TEXT NOT NULL,
    resource_id     TEXT NOT NULL,
    before_ref      TEXT,                    -- JSON reference/checksum of prior state
    after_ref       TEXT,                    -- JSON reference/checksum of new state
    ip              TEXT,
    request_id      TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_audit_events_org_created
    ON audit_events (organization_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_events_actor
    ON audit_events (actor_user_id);

CREATE TABLE IF NOT EXISTS applications (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    owner_user_id   TEXT REFERENCES users(id),
    status          TEXT NOT NULL DEFAULT 'active',
    version         BIGINT NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, name)
);

CREATE TABLE IF NOT EXISTS virtual_servers (
    id                      TEXT PRIMARY KEY, -- ULID
    organization_id         TEXT NOT NULL REFERENCES organizations(id),
    name                    TEXT NOT NULL,
    listen_addr             TEXT NOT NULL DEFAULT '0.0.0.0',
    listen_port             INTEGER NOT NULL,
    tls_enabled             BOOLEAN NOT NULL DEFAULT true,
    certificate_ref         TEXT,             -- reference to certificate metadata, never key material
    default_backend_pool_id TEXT REFERENCES backend_pools(id),
    version                 BIGINT NOT NULL DEFAULT 1,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, listen_addr, listen_port)
);

CREATE TABLE IF NOT EXISTS routes (
    id              TEXT PRIMARY KEY,        -- ULID
    virtual_server_id TEXT NOT NULL REFERENCES virtual_servers(id) ON DELETE CASCADE,
    path            TEXT NOT NULL,
    match_type      TEXT NOT NULL DEFAULT 'prefix',   -- prefix | exact
    backend_pool_id TEXT NOT NULL REFERENCES backend_pools(id),
    priority        INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS backend_pools (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    application_id  TEXT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    lb_algorithm    TEXT NOT NULL DEFAULT 'round_robin',
    version         BIGINT NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, name)
);

CREATE TABLE IF NOT EXISTS backend_nodes (
    id              TEXT PRIMARY KEY,        -- ULID
    pool_id         TEXT NOT NULL REFERENCES backend_pools(id) ON DELETE CASCADE,
    host            TEXT NOT NULL,
    port            INTEGER NOT NULL,
    weight          INTEGER NOT NULL DEFAULT 1,
    active          BOOLEAN NOT NULL DEFAULT true,
    last_health_state TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS health_monitors (
    id              TEXT PRIMARY KEY,        -- ULID
    pool_id         TEXT NOT NULL REFERENCES backend_pools(id) ON DELETE CASCADE,
    type            TEXT NOT NULL DEFAULT 'http',   -- tcp | http
    interval_ms     INTEGER NOT NULL DEFAULT 5000,
    timeout_ms      INTEGER NOT NULL DEFAULT 2000,
    fail_threshold  INTEGER NOT NULL DEFAULT 3,
    pass_threshold  INTEGER NOT NULL DEFAULT 2,
    http_path       TEXT NOT NULL DEFAULT '/healthz',
    http_expected_status INTEGER[] NOT NULL DEFAULT '{200}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS config_versions (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    version         BIGINT NOT NULL,
    bundle_hash     TEXT NOT NULL UNIQUE,
    blob            JSONB NOT NULL,          -- immutable policy document (ADR-002)
    status          TEXT NOT NULL DEFAULT 'draft', -- draft | validating | approved | active | superseded | rolled_back
    created_by      TEXT NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, version)
);

CREATE TABLE IF NOT EXISTS config_deployments (
    id               TEXT PRIMARY KEY,       -- ULID
    config_version_id TEXT NOT NULL REFERENCES config_versions(id),
    target_gateway   TEXT NOT NULL,
    status           TEXT NOT NULL DEFAULT 'pending', -- pending | validating | activated | failed | rolled_back
    error            TEXT,
    applied_hash     TEXT,
    deployed_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Access events (best-effort, high-volume; moves to ClickHouse in later phases).
CREATE TABLE IF NOT EXISTS access_events (
    id                TEXT PRIMARY KEY,      -- ULID
    request_id        TEXT NOT NULL,
    event_id          TEXT NOT NULL,
    gateway_id        TEXT NOT NULL,
    virtual_server_id TEXT NOT NULL,
    application_id    TEXT NOT NULL,
    client_ip         TEXT,
    method            TEXT NOT NULL,
    path              TEXT NOT NULL,
    http_version      TEXT,
    host              TEXT,
    body_size         BIGINT NOT NULL DEFAULT 0,
    status            INTEGER,
    bytes             BIGINT NOT NULL DEFAULT 0,
    latency_ms        NUMERIC(10,3),
    backend_node      TEXT,
    decision_action   TEXT NOT NULL DEFAULT 'pass',
    raw               JSONB,                 -- full event document
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_access_events_created
    ON access_events (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_access_events_application
    ON access_events (application_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_access_events_request_id
    ON access_events (request_id);
