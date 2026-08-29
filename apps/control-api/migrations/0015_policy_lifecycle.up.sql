-- 0015_policy_lifecycle.up.sql
-- Policy lifecycle state (master plan §11 deployment state machine):
-- DRAFT -> VALIDATING -> APPROVAL_REQUIRED -> APPROVED -> CANARY -> ACTIVE
-- -> SUPERSEDED -> ROLLED_BACK

ALTER TABLE security_policies ADD COLUMN IF NOT EXISTS lifecycle_status TEXT NOT NULL DEFAULT 'draft';
ALTER TABLE security_policies ADD COLUMN IF NOT EXISTS active_version_id TEXT;
