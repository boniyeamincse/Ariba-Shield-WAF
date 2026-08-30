-- 0011_notification_channels.up.sql
-- Phase 4: notification channels

-- notification_channels stores configured notification channels for alerts.
CREATE TABLE IF NOT EXISTS notification_channels (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    name            TEXT NOT NULL,
    kind            TEXT NOT NULL,           -- wazuh, syslog, cef, leef, webhook, email, teams, slack, soar
    config          JSONB NOT NULL DEFAULT '{}',
    is_default      BOOLEAN NOT NULL DEFAULT false,
    enabled         BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes for common queries.
CREATE INDEX IF NOT EXISTS idx_notification_channels_org ON notification_channels (organization_id);
CREATE INDEX IF NOT EXISTS idx_notification_channels_default ON notification_channels (organization_id, is_default);
