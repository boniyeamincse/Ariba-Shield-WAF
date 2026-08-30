-- 0013_settings_license.down.sql
-- Revert Phase 8 system settings + license.

DROP INDEX IF EXISTS idx_licenses_status;
DROP INDEX IF EXISTS idx_licenses_org;
DROP TABLE IF EXISTS licenses;

DROP INDEX IF EXISTS idx_system_settings_org_cat;
DROP TABLE IF EXISTS system_settings;
