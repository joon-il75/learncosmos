DROP INDEX IF EXISTS idx_curriculum_patterns_key_version;

CREATE UNIQUE INDEX IF NOT EXISTS idx_curriculum_patterns_key_version_language
  ON curriculum_patterns(pattern_key, pattern_version, language);
