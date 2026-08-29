-- 0003_security_events.up.sql
-- Phase 2: security events table (ADR-002 event schema v0 extension).
-- High-volume, best-effort like access_events; masked payloads only.

CREATE TABLE IF NOT EXISTS security_events (
    id                TEXT PRIMARY KEY,      -- ULID
    event_id          TEXT NOT NULL,
    request_id        TEXT NOT NULL,
    gateway_id        TEXT,
    application_id    TEXT,
    virtual_server_id TEXT,
    client_ip         TEXT,
    method            TEXT,
    path              TEXT,
    host              TEXT,
    status            INTEGER,
    severity          TEXT,
    decision_action   TEXT,
    reason            TEXT,
    rule_ids          TEXT[] NOT NULL DEFAULT '{}',
    match_details     JSONB,
    masked            BOOLEAN NOT NULL DEFAULT true,   -- payloads always masked/omitted
    raw               JSONB,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_sec_events_created
    ON security_events (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sec_events_application
    ON security_events (application_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sec_events_rule
    ON security_events (rule_ids);
CREATE INDEX IF NOT EXISTS idx_sec_events_ip
    ON security_events (client_ip, created_at DESC);
