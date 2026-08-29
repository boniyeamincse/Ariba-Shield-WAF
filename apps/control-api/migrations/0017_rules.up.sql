-- 0017_rules.up.sql
-- Signature lifecycle (master plan §6.3): stable IDs, tags, versions,
-- tests, bundles with signing/staging/canary/rollback.

CREATE TABLE IF NOT EXISTS rules (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    rule_id         TEXT NOT NULL,           -- stable public rule ID (e.g. 950001)
    name            TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    action          TEXT NOT NULL DEFAULT 'LOG',     -- LOG | BLOCK | ALLOW
    severity        TEXT NOT NULL DEFAULT 'medium',  -- low | medium | high | critical
    phase           INTEGER NOT NULL DEFAULT 2,
    source          TEXT NOT NULL DEFAULT '',        -- rule source (human-readable)
    status          TEXT NOT NULL DEFAULT 'active',  -- active | disabled | deprecated | expired
    version         BIGINT NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, rule_id)
);

CREATE TABLE IF NOT EXISTS rule_versions (
    id              TEXT PRIMARY KEY,        -- ULID
    rule_id         TEXT NOT NULL REFERENCES rules(id) ON DELETE CASCADE,
    version         BIGINT NOT NULL,
    source          TEXT NOT NULL,
    bundle_hash     TEXT NOT NULL DEFAULT '',
    status          TEXT NOT NULL DEFAULT 'draft', -- draft | staging | canary | active | superseded | rolled_back
    created_by      TEXT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (rule_id, version)
);

CREATE TABLE IF NOT EXISTS rule_tags (
    rule_id         TEXT NOT NULL REFERENCES rules(id) ON DELETE CASCADE,
    tag             TEXT NOT NULL,           -- attack-class | technology | cve | confidence | sensitivity
    PRIMARY KEY (rule_id, tag)
);

CREATE TABLE IF NOT EXISTS rule_tests (
    id              TEXT PRIMARY KEY,        -- ULID
    rule_id         TEXT NOT NULL REFERENCES rules(id) ON DELETE CASCADE,
    test_type       TEXT NOT NULL DEFAULT 'positive', -- positive | negative
    payload         TEXT NOT NULL,
    expected        TEXT NOT NULL DEFAULT 'block',
    passed          BOOLEAN,
    last_run_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS rule_bundles (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    version         BIGINT NOT NULL DEFAULT 1,
    rule_ids        TEXT[] NOT NULL DEFAULT '{}',
    status          TEXT NOT NULL DEFAULT 'draft', -- draft | signed | staged | canary | active | rolled_back
    signature       TEXT NOT NULL DEFAULT '',
    sign_key_id     TEXT NOT NULL DEFAULT '',
    deployed_gateways TEXT[] NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
