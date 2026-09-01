-- 0030_policy_inheritance.up.sql
-- Policy inheritance hierarchy + canary deployment support.

ALTER TABLE security_policies ADD COLUMN IF NOT EXISTS parent_policy_id TEXT REFERENCES security_policies(id) ON DELETE SET NULL;
ALTER TABLE security_policies ADD COLUMN IF NOT EXISTS canary_percent INTEGER NOT NULL DEFAULT 0; -- 0=all, 1-100=% traffic
ALTER TABLE security_policies ADD COLUMN IF NOT EXISTS scheduled_at TIMESTAMPTZ;
ALTER TABLE security_policies ADD COLUMN IF NOT EXISTS scheduled_action TEXT NOT NULL DEFAULT ''; -- activate | rollback

CREATE INDEX IF NOT EXISTS idx_sec_policies_parent ON security_policies (parent_policy_id);
