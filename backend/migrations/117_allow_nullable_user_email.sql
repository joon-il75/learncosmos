ALTER TABLE users
  ALTER COLUMN email DROP NOT NULL;

UPDATE users
SET email = NULL
WHERE btrim(COALESCE(email, '')) = '';

DROP INDEX IF EXISTS idx_users_email;
ALTER TABLE users
  DROP CONSTRAINT IF EXISTS users_email_key;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique_present
  ON users (lower(email))
  WHERE email IS NOT NULL AND btrim(email) <> '';
