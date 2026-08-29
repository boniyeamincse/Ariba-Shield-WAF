-- 0006_backups.up.sql
-- Phase 4: backup/restore (master plan §12.2, §6.9).

CREATE TABLE IF NOT EXISTS backups (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    status          TEXT NOT NULL DEFAULT 'pending', -- pending | running | completed | failed
    artifact_ref    TEXT,                    -- object-storage reference
    size_bytes      BIGINT NOT NULL DEFAULT 0,
    triggered_by    TEXT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_backups_org ON backups (organization_id, created_at DESC);
