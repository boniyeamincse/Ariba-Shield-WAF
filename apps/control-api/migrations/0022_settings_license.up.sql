-- 0013_settings_license.up.sql
-- Phase 8: System settings + license/entitlement (commercialized).

-- system_settings: key-value settings grouped by category.
CREATE TABLE IF NOT EXISTS system_settings (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    category        TEXT NOT NULL,           -- general | security | localization | retention
    key             TEXT NOT NULL,
    value           JSONB NOT NULL DEFAULT '{}',
    updated_by      TEXT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, category, key)
);

CREATE INDEX IF NOT EXISTS idx_system_settings_org_cat
    ON system_settings (organization_id, category);

-- licenses: entitlement record for a commercial deployment.
CREATE TABLE IF NOT EXISTS licenses (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    license_key     TEXT NOT NULL UNIQUE,
    product         TEXT NOT NULL,           -- ariba-shield-waf
    edition         TEXT NOT NULL DEFAULT 'community', -- community | pro | enterprise
    status          TEXT NOT NULL DEFAULT 'inactive',  -- inactive | active | expired | revoked
    seats           INTEGER NOT NULL DEFAULT 1,
    max_gateways    INTEGER NOT NULL DEFAULT 1,
    max_applications INTEGER NOT NULL DEFAULT 10,
    issued_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at      TIMESTAMPTZ,
    activated_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_licenses_org ON licenses (organization_id);
CREATE INDEX IF NOT EXISTS idx_licenses_status ON licenses (status);
