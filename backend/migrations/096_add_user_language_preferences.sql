ALTER TABLE users
    ADD COLUMN IF NOT EXISTS ui_locale TEXT CHECK (ui_locale IN ('ko', 'en')),
    ADD COLUMN IF NOT EXISTS learning_language TEXT CHECK (learning_language IN ('ko', 'en')),
    ADD COLUMN IF NOT EXISTS language_setup_completed_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_users_language_setup_pending
ON users (id)
WHERE language_setup_completed_at IS NULL;
