-- Safety Gate Phase 2-2: audit events for super-admin moderation rule changes.

CREATE TABLE IF NOT EXISTS moderation_rule_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  rule_id UUID NOT NULL REFERENCES moderation_rules(id) ON DELETE CASCADE,
  admin_user_id TEXT NOT NULL DEFAULT '',
  action TEXT NOT NULL CHECK (action IN ('create', 'update', 'deactivate', 'reactivate')),
  before_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
  after_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_moderation_rule_events_rule_created
  ON moderation_rule_events(rule_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_moderation_rule_events_admin_created
  ON moderation_rule_events(admin_user_id, created_at DESC);
