-- 0012_reports.down.sql
-- Revert Phase 4 reports.

DROP INDEX IF EXISTS idx_reports_org_created;
DROP INDEX IF EXISTS idx_reports_org_kind;

DROP TABLE IF EXISTS reports;