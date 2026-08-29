-- 0009_auth_session.up.sql
-- Auth module completion: MFA (TOTP) secret, session refresh, break-glass.

ALTER TABLE users ADD COLUMN IF NOT EXISTS totp_secret TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS break_glass_enabled BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE users ADD COLUMN IF NOT EXISTS break_glass_expires_at TIMESTAMPTZ;

-- Track session refresh rotations.
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS refresh_token_hash TEXT;
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS last_used_at TIMESTAMPTZ;
