-- Safety Gate Phase 1 follow-up: prevent duplicate moderation rules.

CREATE UNIQUE INDEX IF NOT EXISTS idx_moderation_rules_unique_pattern_locale
  ON moderation_rules(rule_type, pattern, COALESCE(locale, ''));
