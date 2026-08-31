-- 0024_rules_wizard.up.sql
-- Extend rules with full wizard fields (category, priority) plus conditions
-- and scopes tables per the Rules module spec.

ALTER TABLE rules ADD COLUMN IF NOT EXISTS category TEXT NOT NULL DEFAULT '';
ALTER TABLE rules ADD COLUMN IF NOT EXISTS priority INTEGER NOT NULL DEFAULT 0;
ALTER TABLE rules ADD COLUMN IF NOT EXISTS logic TEXT NOT NULL DEFAULT 'AND'; -- AND | OR
ALTER TABLE rules ADD COLUMN IF NOT EXISTS type TEXT NOT NULL DEFAULT 'custom'; -- managed | custom

-- waf_rule_conditions: one rule can have multiple conditions.
CREATE TABLE IF NOT EXISTS waf_rule_conditions (
    id              TEXT PRIMARY KEY,        -- ULID
    rule_id         TEXT NOT NULL REFERENCES rules(id) ON DELETE CASCADE,
    group_id        INTEGER NOT NULL DEFAULT 0, -- for grouping (AND/OR groups)
    field           TEXT NOT NULL,           -- request_body | method | url | header | query_param | cookie | source_ip | user_agent | host
    operator        TEXT NOT NULL,           -- equals | not_equals | contains | not_contains | starts_with | ends_with | regex | ip_match | cidr_match | gt | lt
    value           TEXT NOT NULL,
    transformation  TEXT NOT NULL DEFAULT '', -- url_decode | base64_decode | lowercase | none
    case_sensitive  BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_waf_rule_conditions_rule ON waf_rule_conditions (rule_id);

-- waf_rule_scopes: which applications/paths/methods a rule applies to.
CREATE TABLE IF NOT EXISTS waf_rule_scopes (
    id              TEXT PRIMARY KEY,        -- ULID
    rule_id         TEXT NOT NULL REFERENCES rules(id) ON DELETE CASCADE,
    application_id  TEXT REFERENCES applications(id) ON DELETE CASCADE,
    path_pattern    TEXT NOT NULL DEFAULT '*',  -- glob pattern
    methods         TEXT[] NOT NULL DEFAULT '{}', -- e.g. {GET,POST,PUT,PATCH,DELETE}
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_waf_rule_scopes_rule ON waf_rule_scopes (rule_id);