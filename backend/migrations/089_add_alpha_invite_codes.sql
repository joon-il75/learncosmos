ALTER TABLE users
  ADD COLUMN IF NOT EXISTS alpha_access_granted_at TIMESTAMPTZ NULL;

CREATE TABLE IF NOT EXISTS alpha_invite_codes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT NOT NULL UNIQUE,
  status TEXT NOT NULL DEFAULT 'active',
  max_uses INTEGER NOT NULL DEFAULT 1,
  used_count INTEGER NOT NULL DEFAULT 0,
  expires_at TIMESTAMPTZ NOT NULL,
  sent_to_note TEXT NOT NULL DEFAULT '',
  admin_note TEXT NOT NULL DEFAULT '',
  created_by UUID NULL REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT alpha_invite_codes_status_check CHECK (status IN ('active', 'revoked')),
  CONSTRAINT alpha_invite_codes_max_uses_check CHECK (max_uses > 0),
  CONSTRAINT alpha_invite_codes_used_count_check CHECK (used_count >= 0 AND used_count <= max_uses)
);

CREATE INDEX IF NOT EXISTS idx_alpha_invite_codes_expires_at
  ON alpha_invite_codes(expires_at DESC);

CREATE TABLE IF NOT EXISTS alpha_invite_code_uses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code_id UUID NOT NULL REFERENCES alpha_invite_codes(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  provider TEXT NULL,
  provider_id TEXT NULL,
  used_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT alpha_invite_code_uses_code_user_unique UNIQUE (code_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_alpha_invite_code_uses_user_id
  ON alpha_invite_code_uses(user_id);

CREATE INDEX IF NOT EXISTS idx_alpha_invite_code_uses_code_id
  ON alpha_invite_code_uses(code_id);
