-- 0029_acme_certificates.up.sql
-- ACME auto-provisioning columns for the certificates table.

ALTER TABLE certificates ADD COLUMN IF NOT EXISTS acme_enabled      BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE certificates ADD COLUMN IF NOT EXISTS acme_account_email TEXT NOT NULL DEFAULT '';
ALTER TABLE certificates ADD COLUMN IF NOT EXISTS acme_staging      BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE certificates ADD COLUMN IF NOT EXISTS acme_last_error   TEXT NOT NULL DEFAULT '';