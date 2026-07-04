-- Safety Gate Phase 1: rule-based moderation logs and rules.

CREATE TABLE IF NOT EXISTS moderation_rules (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  rule_type TEXT NOT NULL CHECK (rule_type IN ('keyword', 'regex')),
  pattern TEXT NOT NULL,
  risk_type TEXT NOT NULL,
  action TEXT NOT NULL CHECK (action IN ('soft_warn', 'block')),
  locale TEXT,
  description TEXT NOT NULL DEFAULT '',
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_moderation_rules_active
  ON moderation_rules(is_active, locale, rule_type);

CREATE INDEX IF NOT EXISTS idx_moderation_rules_risk_action
  ON moderation_rules(risk_type, action);

CREATE TABLE IF NOT EXISTS moderation_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  target_type TEXT NOT NULL,
  target_id UUID,
  input_text_hash TEXT NOT NULL,
  risk_type TEXT NOT NULL,
  action TEXT NOT NULL CHECK (action IN ('allow', 'soft_warn', 'block')),
  matched_rule_id UUID REFERENCES moderation_rules(id) ON DELETE SET NULL,
  route TEXT NOT NULL DEFAULT '',
  locale TEXT,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_moderation_logs_created
  ON moderation_logs(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_moderation_logs_user_created
  ON moderation_logs(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_moderation_logs_target
  ON moderation_logs(target_type, target_id);

CREATE INDEX IF NOT EXISTS idx_moderation_logs_action_risk
  ON moderation_logs(action, risk_type, created_at DESC);

INSERT INTO moderation_rules (rule_type, pattern, risk_type, action, locale, description)
VALUES
  ('keyword', '죽이는 법', 'unsafe_instruction', 'block', 'ko', 'Harmful how-to phrasing'),
  ('keyword', '폭탄 만들기', 'unsafe_instruction', 'block', 'ko', 'Weapon construction phrasing'),
  ('keyword', '자살 방법', 'self_harm', 'block', 'ko', 'Self-harm instruction phrasing'),
  ('keyword', '혐오', 'abuse_or_harassment', 'soft_warn', 'ko', 'Potential harassment phrasing'),
  ('keyword', 'how to kill', 'unsafe_instruction', 'block', 'en', 'Harmful how-to phrasing'),
  ('keyword', 'make a bomb', 'unsafe_instruction', 'block', 'en', 'Weapon construction phrasing'),
  ('keyword', 'suicide method', 'self_harm', 'block', 'en', 'Self-harm instruction phrasing'),
  ('keyword', 'hate speech', 'abuse_or_harassment', 'soft_warn', 'en', 'Potential harassment phrasing')
ON CONFLICT DO NOTHING;
