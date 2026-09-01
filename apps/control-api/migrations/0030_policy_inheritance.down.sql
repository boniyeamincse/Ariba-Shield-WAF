DROP INDEX IF EXISTS idx_sec_policies_parent;
ALTER TABLE security_policies DROP COLUMN IF EXISTS scheduled_action;
ALTER TABLE security_policies DROP COLUMN IF EXISTS scheduled_at;
ALTER TABLE security_policies DROP COLUMN IF EXISTS canary_percent;
ALTER TABLE security_policies DROP COLUMN IF EXISTS parent_policy_id;
