-- 0012_listeners.up.sql
-- Listeners / Virtual Servers: add status for enable/disable actions.

ALTER TABLE virtual_servers ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'enabled'; -- enabled | disabled
