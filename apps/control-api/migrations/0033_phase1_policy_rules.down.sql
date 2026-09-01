-- 0033_phase1_policy_rules.down.sql
-- Revert Phase 1

ALTER TABLE rules DROP COLUMN IF EXISTS lifecycle_status;
ALTER TABLE rules DROP COLUMN IF EXISTS cve_cwe;
ALTER TABLE rules DROP COLUMN IF EXISTS owasp;
ALTER TABLE rules DROP COLUMN IF EXISTS tags;
ALTER TABLE rules DROP COLUMN IF EXISTS sub_category;

DROP INDEX IF EXISTS idx_policy_rules_rule;
DROP INDEX IF EXISTS idx_policy_rules_policy;
DROP TABLE IF EXISTS policy_rules;