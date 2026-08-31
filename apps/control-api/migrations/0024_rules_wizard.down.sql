-- 0024_rules_wizard.down.sql

DROP INDEX IF EXISTS idx_waf_rule_scopes_rule;
DROP TABLE IF EXISTS waf_rule_scopes;

DROP INDEX IF EXISTS idx_waf_rule_conditions_rule;
DROP TABLE IF EXISTS waf_rule_conditions;

ALTER TABLE rules DROP COLUMN IF EXISTS logic;
ALTER TABLE rules DROP COLUMN IF EXISTS priority;
ALTER TABLE rules DROP COLUMN IF EXISTS category;