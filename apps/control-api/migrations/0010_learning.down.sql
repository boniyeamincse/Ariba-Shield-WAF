-- 0010_learning.down.sql
-- Revert Phase 6 learning tables.

DROP INDEX IF EXISTS idx_learning_suggestions_app;
DROP INDEX IF EXISTS idx_learning_suggestions_status;
DROP INDEX IF EXISTS idx_learning_suggestions_session_id;

DROP TABLE IF EXISTS learning_suggestions;
DROP TABLE IF EXISTS learning_sessions;
