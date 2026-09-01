CREATE TABLE IF NOT EXISTS geo_blocking (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name TEXT NOT NULL,
    blocked_countries TEXT[] NOT NULL DEFAULT '{}',
    allowed_countries TEXT[] NOT NULL DEFAULT '{}',
    action TEXT NOT NULL DEFAULT 'block',
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_geo_org ON geo_blocking (organization_id);
