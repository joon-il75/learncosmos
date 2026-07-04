ALTER TABLE users
    ADD COLUMN IF NOT EXISTS premium_access BOOLEAN NOT NULL DEFAULT false;

CREATE INDEX IF NOT EXISTS idx_users_premium_access ON users(premium_access);
