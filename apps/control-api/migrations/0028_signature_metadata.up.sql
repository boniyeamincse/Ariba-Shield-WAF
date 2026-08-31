-- 0028_signature_metadata.up.sql
-- F5-style signature metadata columns for the rules table.

ALTER TABLE rules ADD COLUMN IF NOT EXISTS attack_type   TEXT NOT NULL DEFAULT '';
ALTER TABLE rules ADD COLUMN IF NOT EXISTS pattern_type  TEXT NOT NULL DEFAULT 'regex'; -- regex | contains | equals | ip | cidr | starts_with | ends_with
ALTER TABLE rules ADD COLUMN IF NOT EXISTS accuracy      INTEGER NOT NULL DEFAULT 85;   -- 0-100
ALTER TABLE rules ADD COLUMN IF NOT EXISTS risk          TEXT NOT NULL DEFAULT 'medium'; -- low | medium | high | critical
ALTER TABLE rules ADD COLUMN IF NOT EXISTS confidence    INTEGER NOT NULL DEFAULT 80;   -- 0-100
ALTER TABLE rules ADD COLUMN IF NOT EXISTS source        TEXT NOT NULL DEFAULT 'ariba-core';
ALTER TABLE rules ADD COLUMN IF NOT EXISTS staging       BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE rules ADD COLUMN IF NOT EXISTS remediation   TEXT NOT NULL DEFAULT '';