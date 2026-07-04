-- Safety Gate AI output review workflow decisions.

CREATE TABLE IF NOT EXISTS moderation_ai_output_review_decisions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  candidate_key TEXT NOT NULL UNIQUE,
  status TEXT NOT NULL DEFAULT 'unreviewed' CHECK (status IN ('unreviewed', 'false_positive', 'needs_prompt_guard', 'needs_rule_tuning', 'needs_masking', 'needs_regeneration', 'resolved')),
  note TEXT NOT NULL DEFAULT '',
  reviewed_by TEXT NOT NULL DEFAULT '',
  reviewed_at TIMESTAMPTZ,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_moderation_ai_output_review_decisions_status
  ON moderation_ai_output_review_decisions(status, updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_moderation_ai_output_review_decisions_reviewed_by
  ON moderation_ai_output_review_decisions(reviewed_by, reviewed_at DESC);
