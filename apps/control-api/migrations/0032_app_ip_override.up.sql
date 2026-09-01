-- App-level IP blacklist/whitelist override.
CREATE TABLE IF NOT EXISTS app_ip_overrides (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    application_id TEXT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    list_type TEXT NOT NULL DEFAULT 'block', -- block | allow
    ip TEXT NOT NULL,
    reason TEXT NOT NULL DEFAULT '',
    created_by TEXT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (application_id, ip)
);
CREATE INDEX IF NOT EXISTS idx_app_ip_overrides_app ON app_ip_overrides (application_id);
