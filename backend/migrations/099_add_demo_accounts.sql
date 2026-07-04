CREATE TABLE IF NOT EXISTS demo_accounts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  label TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL UNIQUE,
  login_code_hash TEXT NOT NULL UNIQUE,
  assigned_to TEXT,
  assignment_note TEXT,
  expires_at TIMESTAMPTZ,
  last_used_at TIMESTAMPTZ,
  disabled_at TIMESTAMPTZ,
  disabled_reason TEXT,
  created_by UUID REFERENCES users(id) ON DELETE SET NULL,
  created_by_actor TEXT NOT NULL DEFAULT 'super_admin',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_demo_accounts_status
  ON demo_accounts(disabled_at, expires_at, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_demo_accounts_user_id
  ON demo_accounts(user_id);
