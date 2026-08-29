-- 0002_api_v1.up.sql
-- Phase 1 API expansion: domains, origins, security policies, gateways.
-- Additive-only. New changes = new migration.

CREATE TABLE IF NOT EXISTS domains (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    application_id  TEXT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    hostname        TEXT NOT NULL,
    enabled         BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (application_id, hostname)
);

CREATE INDEX IF NOT EXISTS idx_domains_app ON domains (application_id);

CREATE TABLE IF NOT EXISTS origins (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    application_id  TEXT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    protocol        TEXT NOT NULL DEFAULT 'http',   -- http | https
    host            TEXT NOT NULL,
    port            INTEGER NOT NULL DEFAULT 80,
    weight          INTEGER NOT NULL DEFAULT 1,
    enabled         BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (application_id, name)
);

CREATE INDEX IF NOT EXISTS idx_origins_app ON origins (application_id);

CREATE TABLE IF NOT EXISTS security_policies (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    application_id  TEXT REFERENCES applications(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    enforcement_mode TEXT NOT NULL DEFAULT 'transparent', -- transparent | alarm | blocking
    version         BIGINT NOT NULL DEFAULT 1,
    created_by      TEXT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, name)
);

CREATE INDEX IF NOT EXISTS idx_policies_app ON security_policies (application_id);

CREATE TABLE IF NOT EXISTS gateways (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    hostname        TEXT NOT NULL,
    ip              TEXT,
    version         TEXT,
    capabilities    TEXT[] NOT NULL DEFAULT '{}',
    status          TEXT NOT NULL DEFAULT 'unregistered', -- unregistered | starting | active | degraded | stopping | offline
    last_seen_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_gateways_org ON gateways (organization_id, status);

CREATE TABLE IF NOT EXISTS gateway_heartbeats (
    id              TEXT PRIMARY KEY,        -- ULID
    gateway_id      TEXT NOT NULL REFERENCES gateways(id) ON DELETE CASCADE,
    status          TEXT NOT NULL,
    applied_hash    TEXT,
    version         TEXT,
    health          JSONB,
    last_error      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_heartbeats_gateway ON gateway_heartbeats (gateway_id, created_at DESC);
