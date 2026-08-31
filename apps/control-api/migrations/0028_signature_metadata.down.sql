-- 0028_signature_metadata.down.sql
-- Revert F5-style signature metadata columns.

ALTER TABLE rules DROP COLUMN IF EXISTS remediation;
ALTER TABLE rules DROP COLUMN IF EXISTS staging;
ALTER TABLE rules DROP COLUMN IF EXISTS source;
ALTER TABLE rules DROP COLUMN IF EXISTS confidence;
ALTER TABLE rules DROP COLUMN IF EXISTS risk;
ALTER TABLE rules DROP COLUMN IF EXISTS accuracy;
ALTER TABLE rules DROP COLUMN IF EXISTS pattern_type;
ALTER TABLE rules DROP COLUMN IF EXISTS attack_type;
