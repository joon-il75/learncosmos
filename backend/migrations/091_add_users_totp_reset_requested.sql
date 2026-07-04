ALTER TABLE users
    ADD COLUMN IF NOT EXISTS totp_reset_requested BOOLEAN NOT NULL DEFAULT false;

CREATE INDEX IF NOT EXISTS idx_users_totp_reset_requested
    ON users (totp_reset_requested)
    WHERE totp_reset_requested = true;
