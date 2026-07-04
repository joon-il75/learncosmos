ALTER TABLE users
ADD COLUMN IF NOT EXISTS display_id TEXT UNIQUE;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_display_id
ON users (display_id)
WHERE display_id IS NOT NULL;
