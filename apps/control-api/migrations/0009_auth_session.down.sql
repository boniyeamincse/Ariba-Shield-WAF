ALTER TABLE sessions DROP COLUMN IF EXISTS last_used_at;
ALTER TABLE sessions DROP COLUMN IF EXISTS refresh_token_hash;
ALTER TABLE users DROP COLUMN IF EXISTS break_glass_expires_at;
ALTER TABLE users DROP COLUMN IF EXISTS break_glass_enabled;
ALTER TABLE users DROP COLUMN IF EXISTS totp_secret;
