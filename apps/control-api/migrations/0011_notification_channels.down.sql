-- 0011_notification_channels.down.sql
-- Revert Phase 4 notification channels.

DROP INDEX IF EXISTS idx_notification_channels_default;
DROP INDEX IF EXISTS idx_notification_channels_org;

DROP TABLE IF EXISTS notification_channels;
