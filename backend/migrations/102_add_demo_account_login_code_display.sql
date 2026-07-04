ALTER TABLE demo_accounts
  ADD COLUMN IF NOT EXISTS login_code TEXT;

CREATE INDEX IF NOT EXISTS idx_demo_accounts_login_code
  ON demo_accounts(login_code)
  WHERE login_code IS NOT NULL;
