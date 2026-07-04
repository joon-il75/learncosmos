ALTER TABLE demo_accounts
  ADD COLUMN IF NOT EXISTS activated_at TIMESTAMPTZ;

UPDATE demo_accounts
SET activated_at = last_used_at
WHERE activated_at IS NULL
  AND last_used_at IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_demo_accounts_activated_at
  ON demo_accounts(activated_at DESC);
