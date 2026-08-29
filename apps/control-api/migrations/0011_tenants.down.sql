DROP TABLE IF EXISTS tenants;

ALTER TABLE organizations DROP COLUMN IF EXISTS updated_at;
ALTER TABLE organizations DROP COLUMN IF EXISTS version;
ALTER TABLE organizations DROP COLUMN IF EXISTS status;
ALTER TABLE organizations DROP COLUMN IF EXISTS quotas;
ALTER TABLE organizations DROP COLUMN IF EXISTS plan;
ALTER TABLE organizations DROP COLUMN IF EXISTS contact_email;
ALTER TABLE organizations DROP COLUMN IF EXISTS description;
