-- 0033_phase1_policy_rules.up.sql
-- Phase 1: Policy-Rule junction + rule lifecycle + import fields

-- 1. Policy ⇄ Rule many-to-many junction
CREATE TABLE IF NOT EXISTS policy_rules (
    id              TEXT PRIMARY KEY,
    policy_id       TEXT NOT NULL REFERENCES security_policies(id) ON DELETE CASCADE,
    rule_id         TEXT NOT NULL REFERENCES rules(id) ON DELETE CASCADE,
    action_override TEXT NOT NULL DEFAULT '', -- inherit | allow | log | block | challenge
    enabled         BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (policy_id, rule_id)
);
CREATE INDEX IF NOT EXISTS idx_policy_rules_policy ON policy_rules (policy_id);
CREATE INDEX IF NOT EXISTS idx_policy_rules_rule   ON policy_rules (rule_id);

-- 2. Extra rule fields for dataset
ALTER TABLE rules ADD COLUMN IF NOT EXISTS sub_category TEXT NOT NULL DEFAULT '';
ALTER TABLE rules ADD COLUMN IF NOT EXISTS tags        TEXT[] NOT NULL DEFAULT '{}';
ALTER TABLE rules ADD COLUMN IF NOT EXISTS owasp       TEXT NOT NULL DEFAULT '';
ALTER TABLE rules ADD COLUMN IF NOT EXISTS cve_cwe     TEXT NOT NULL DEFAULT '';

-- 3. Rule lifecycle status (extend existing status)
-- existing: active | disabled | deprecated | expired
-- new: imported | draft | validated
ALTER TABLE rules ADD COLUMN IF NOT EXISTS lifecycle_status TEXT NOT NULL DEFAULT 'active';