-- 0018_bot_abuse.up.sql
-- Bot & abuse protection (master plan §6.6): client classification,
-- challenges, progressive actions.

CREATE TABLE IF NOT EXISTS bot_clients (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    application_id  TEXT REFERENCES applications(id) ON DELETE CASCADE,
    client_ip       TEXT NOT NULL,
    classification  TEXT NOT NULL DEFAULT 'unknown', -- human | bot | known_bot | scraper | credential_stuffing | otp_abuse
    confidence      TEXT NOT NULL DEFAULT 'low',
    user_agent      TEXT NOT NULL DEFAULT '',
    score           INTEGER NOT NULL DEFAULT 0,
    action          TEXT NOT NULL DEFAULT 'observe',  -- observe | slow | challenge | temporary_block
    status          TEXT NOT NULL DEFAULT 'tracked',
    first_seen_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS bot_challenges (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    client_ip       TEXT NOT NULL,
    challenge_type  TEXT NOT NULL DEFAULT 'javascript', -- javascript | captcha | proof_of_work
    status          TEXT NOT NULL DEFAULT 'issued',  -- issued | solved | expired | revoked
    issued_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at      TIMESTAMPTZ NOT NULL DEFAULT now() + interval '10 minutes',
    solved_at       TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS bot_events (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    application_id  TEXT REFERENCES applications(id) ON DELETE CASCADE,
    client_ip       TEXT NOT NULL,
    classification  TEXT NOT NULL DEFAULT 'unknown',
    reason          TEXT NOT NULL DEFAULT '',
    action          TEXT NOT NULL DEFAULT 'observe',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
