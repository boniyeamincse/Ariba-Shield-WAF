-- 0012_reports.up.sql
-- Phase 4: Reports

-- reports stores generated security, traffic, incident, and compliance reports.
CREATE TABLE IF NOT EXISTS reports (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    kind            TEXT NOT NULL,           -- security | traffic | incidents | compliance
    status          TEXT NOT NULL DEFAULT 'pending', -- pending | ready | failed
    params          JSONB NOT NULL DEFAULT '{}',     -- query params used to generate
    summary         JSONB NOT NULL DEFAULT '{}',     -- aggregate data of the report
    error           TEXT,
    created_by      TEXT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_reports_org_kind ON reports (organization_id, kind);
CREATE INDEX IF NOT EXISTS idx_reports_org_created ON reports (organization_id, created_at DESC);
