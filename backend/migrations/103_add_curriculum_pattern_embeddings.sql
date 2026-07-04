CREATE TABLE IF NOT EXISTS curriculum_patterns (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  pattern_key TEXT NOT NULL,
  pattern_version TEXT NOT NULL DEFAULT 'v1',
  domain TEXT NOT NULL DEFAULT '',
  subdomain TEXT NOT NULL DEFAULT '',
  goal_type TEXT NOT NULL DEFAULT '',
  learner_level TEXT NOT NULL DEFAULT 'any',
  language TEXT NOT NULL DEFAULT 'ko',
  title TEXT NOT NULL,
  summary TEXT NOT NULL DEFAULT '',
  trigger_keywords JSONB NOT NULL DEFAULT '[]'::jsonb,
  guidance JSONB NOT NULL DEFAULT '{}'::jsonb,
  stage_rules JSONB NOT NULL DEFAULT '[]'::jsonb,
  lesson_structure JSONB NOT NULL DEFAULT '{}'::jsonb,
  recommended_sequence JSONB NOT NULL DEFAULT '[]'::jsonb,
  bad_patterns JSONB NOT NULL DEFAULT '[]'::jsonb,
  query_terms JSONB NOT NULL DEFAULT '[]'::jsonb,
  embedding_text TEXT NOT NULL,
  source JSONB NOT NULL DEFAULT '{}'::jsonb,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT curriculum_patterns_key_check CHECK (btrim(pattern_key) <> ''),
  CONSTRAINT curriculum_patterns_version_check CHECK (btrim(pattern_version) <> ''),
  CONSTRAINT curriculum_patterns_language_check CHECK (language IN ('ko', 'en')),
  CONSTRAINT curriculum_patterns_title_check CHECK (btrim(title) <> ''),
  CONSTRAINT curriculum_patterns_embedding_text_check CHECK (btrim(embedding_text) <> ''),
  CONSTRAINT curriculum_patterns_trigger_keywords_json_check CHECK (jsonb_typeof(trigger_keywords) = 'array'),
  CONSTRAINT curriculum_patterns_guidance_json_check CHECK (jsonb_typeof(guidance) = 'object'),
  CONSTRAINT curriculum_patterns_stage_rules_json_check CHECK (jsonb_typeof(stage_rules) = 'array'),
  CONSTRAINT curriculum_patterns_lesson_structure_json_check CHECK (jsonb_typeof(lesson_structure) = 'object'),
  CONSTRAINT curriculum_patterns_recommended_sequence_json_check CHECK (jsonb_typeof(recommended_sequence) = 'array'),
  CONSTRAINT curriculum_patterns_bad_patterns_json_check CHECK (jsonb_typeof(bad_patterns) = 'array'),
  CONSTRAINT curriculum_patterns_query_terms_json_check CHECK (jsonb_typeof(query_terms) = 'array'),
  CONSTRAINT curriculum_patterns_source_json_check CHECK (jsonb_typeof(source) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_curriculum_patterns_key_version
  ON curriculum_patterns(pattern_key, pattern_version);

CREATE INDEX IF NOT EXISTS idx_curriculum_patterns_active_lookup
  ON curriculum_patterns(is_active, language, pattern_version, domain, goal_type);

CREATE INDEX IF NOT EXISTS idx_curriculum_patterns_updated_at
  ON curriculum_patterns(updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_curriculum_patterns_trigger_keywords_gin
  ON curriculum_patterns USING GIN (trigger_keywords);

CREATE INDEX IF NOT EXISTS idx_curriculum_patterns_guidance_gin
  ON curriculum_patterns USING GIN (guidance);

CREATE TABLE IF NOT EXISTS curriculum_pattern_embeddings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  pattern_id UUID NOT NULL REFERENCES curriculum_patterns(id) ON DELETE CASCADE,
  provider TEXT NOT NULL DEFAULT 'embedding_gemma',
  model TEXT NOT NULL DEFAULT 'embedding-gemma',
  dimension INTEGER NOT NULL DEFAULT 768,
  embedding VECTOR,
  embedding_text_hash TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending',
  error_message TEXT NOT NULL DEFAULT '',
  embedded_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT curriculum_pattern_embeddings_provider_check CHECK (btrim(provider) <> ''),
  CONSTRAINT curriculum_pattern_embeddings_model_check CHECK (btrim(model) <> ''),
  CONSTRAINT curriculum_pattern_embeddings_dimension_check CHECK (dimension > 0),
  CONSTRAINT curriculum_pattern_embeddings_hash_check CHECK (btrim(embedding_text_hash) <> ''),
  CONSTRAINT curriculum_pattern_embeddings_status_check CHECK (
    status IN ('pending', 'ready', 'failed', 'stale')
  )
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_curriculum_pattern_embeddings_unique
  ON curriculum_pattern_embeddings(pattern_id, provider, model, dimension, embedding_text_hash);

CREATE INDEX IF NOT EXISTS idx_curriculum_pattern_embeddings_provider_status
  ON curriculum_pattern_embeddings(provider, model, dimension, status, updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_curriculum_pattern_embeddings_pattern_id
  ON curriculum_pattern_embeddings(pattern_id, updated_at DESC);
