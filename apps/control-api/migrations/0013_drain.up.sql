-- 0013_drain.up.sql
-- Backend pool draining (master plan §6.1 connection draining).

ALTER TABLE backend_nodes ADD COLUMN IF NOT EXISTS draining BOOLEAN NOT NULL DEFAULT false;
