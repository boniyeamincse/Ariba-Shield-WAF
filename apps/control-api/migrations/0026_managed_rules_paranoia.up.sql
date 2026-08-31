-- 0025_managed_rules_paranoia.up.sql
-- Extend managed_rules with paranoia_level, action, anomaly_threshold
-- per the managed rules implementation plan.

ALTER TABLE managed_rules ADD COLUMN IF NOT EXISTS paranoia_level INTEGER NOT NULL DEFAULT 1; -- 1=low,2=medium,3=high,4=strict
ALTER TABLE managed_rules ADD COLUMN IF NOT EXISTS action TEXT NOT NULL DEFAULT 'block'; -- block | log
ALTER TABLE managed_rules ADD COLUMN IF NOT EXISTS anomaly_threshold INTEGER NOT NULL DEFAULT 5;

-- Update existing seed rows to use the new columns.
UPDATE managed_rules SET paranoia_level = 
  CASE sensitivity
    WHEN 'low' THEN 1
    WHEN 'medium' THEN 2
    WHEN 'high' THEN 3
    WHEN 'strict' THEN 4
    ELSE 1
  END;
