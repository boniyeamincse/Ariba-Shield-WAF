-- 0010_learning.up.sql
-- Phase 6: Learning and policy builder

-- learning_sessions tracks trusted-source learning sessions.
-- Learning system cannot directly see live traffic; only trusted sessions.
CREATE TABLE IF NOT EXISTS learning_sessions (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    source          TEXT NOT NULL,           -- trusted | csv | file | api
    description     TEXT NOT NULL DEFAULT '',
    confidence_threshold  TEXT NOT NULL DEFAULT '0.7', -- 0.0 to 1.0
    status          TEXT NOT NULL DEFAULT 'active', -- active | paused | completed
    created_by      TEXT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- learning_suggestions engine-generated recommendations.
-- Only after /accept does policy get updated (exit criteria: learning cannot directly weaken policy).
CREATE TABLE IF NOT EXISTS learning_suggestions (
    id              TEXT PRIMARY KEY,        -- ULID
    session_id      TEXT NOT NULL REFERENCES learning_sessions(id) ON DELETE CASCADE,
    application_id  TEXT REFERENCES applications(id) ON DELETE SET NULL,
    rule_id         TEXT NOT NULL,           -- proposed rule ID (e.g., "ARIBA-SQLI-001")
    severity        TEXT NOT NULL DEFAULT 'medium', -- low | medium | high | critical
    confidence      TEXT NOT NULL DEFAULT '0.7', -- 0.0 to 1.0, from engine
    rationale       TEXT NOT NULL DEFAULT '',    -- why the engine flagged this
    status          TEXT NOT NULL DEFAULT 'pending', -- pending | accepted | rejected
    applied_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes for suggestion queries.
CREATE INDEX IF NOT EXISTS idx_learning_suggestions_session_id ON learning_suggestions (session_id);
CREATE INDEX IF NOT EXISTS idx_learning_suggestions_status ON learning_suggestions (status);
CREATE INDEX IF NOT EXISTS idx_learning_suggestions_app ON learning_suggestions (application_id, status);
