-- 0010_sites.up.sql
-- Sites / data centers (master plan §8: `sites`, §6.9 central management).

CREATE TABLE IF NOT EXISTS sites (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    location        TEXT NOT NULL DEFAULT '',
    country_code    TEXT NOT NULL DEFAULT '',
    gateway_ids     TEXT[] NOT NULL DEFAULT '{}',
    status          TEXT NOT NULL DEFAULT 'active', -- active | maintenance | degraded | offline
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
