ALTER TABLE managed_rules DROP COLUMN IF EXISTS anomaly_threshold;
ALTER TABLE managed_rules DROP COLUMN IF EXISTS action;
ALTER TABLE managed_rules DROP COLUMN IF EXISTS paranoia_level;
