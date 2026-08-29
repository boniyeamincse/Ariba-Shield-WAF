-- 0004_phase3.up.sql
-- Phase 3 (safe blocking): IP lists, rate limit policies, policy versioning.

CREATE TABLE IF NOT EXISTS ip_lists (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    list_type       TEXT NOT NULL,           -- allowed | blocked
    entries         TEXT[] NOT NULL DEFAULT '{}',  -- IP/CIDR
    description     TEXT NOT NULL DEFAULT '',
    version         BIGINT NOT NULL DEFAULT 1,
    created_by      TEXT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, name)
);

CREATE TABLE IF NOT EXISTS rate_limit_policies (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    application_id  TEXT REFERENCES applications(id) ON DELETE CASCADE,
    route_prefix    TEXT NOT NULL DEFAULT '',
    limit_count     INTEGER NOT NULL,
    window_seconds  INTEGER NOT NULL DEFAULT 60,
    action          TEXT NOT NULL DEFAULT 'throttle',  -- throttle | block
    version         BIGINT NOT NULL DEFAULT 1,
    created_by      TEXT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, name)
);

-- Policy versions: immutable snapshots for draft/approve/activate/rollback.
CREATE TABLE IF NOT EXISTS policy_versions (
    id              TEXT PRIMARY KEY,        -- ULID
    policy_id       TEXT NOT NULL REFERENCES security_policies(id) ON DELETE CASCADE,
    version         BIGINT NOT NULL,
    document        JSONB NOT NULL,          -- immutable policy document snapshot
    bundle_hash     TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'draft',  -- draft | approved | active | superseded | rolled_back
    created_by      TEXT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (policy_id, version)
);
