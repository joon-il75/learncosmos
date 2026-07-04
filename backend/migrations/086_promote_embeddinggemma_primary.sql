ALTER TABLE IF EXISTS contents
  DROP COLUMN IF EXISTS embedding;

DO $$
BEGIN
  IF to_regclass('public.content_embeddings') IS NULL
     AND to_regclass('public.content_embedding_shadows') IS NOT NULL THEN
    ALTER TABLE content_embedding_shadows RENAME TO content_embeddings;
  END IF;
END $$;

CREATE TABLE IF NOT EXISTS content_embeddings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  content_id UUID NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
  provider TEXT NOT NULL DEFAULT 'embedding_gemma',
  model TEXT NOT NULL DEFAULT 'embedding-gemma',
  dimension INTEGER NOT NULL,
  embedding VECTOR,
  source_text_hash TEXT NOT NULL DEFAULT '',
  source_text_preview TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'pending',
  error_message TEXT NOT NULL DEFAULT '',
  embedded_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT content_embeddings_dimension_check CHECK (dimension > 0),
  CONSTRAINT content_embeddings_status_check CHECK (
    status IN ('pending', 'ready', 'failed', 'stale')
  )
);

DROP INDEX IF EXISTS idx_content_embedding_shadows_unique;
DROP INDEX IF EXISTS idx_content_embedding_shadows_provider_status;
DROP INDEX IF EXISTS idx_content_embedding_shadows_content_id;

CREATE UNIQUE INDEX IF NOT EXISTS idx_content_embeddings_unique
  ON content_embeddings(content_id, provider, model, dimension);

CREATE INDEX IF NOT EXISTS idx_content_embeddings_provider_status
  ON content_embeddings(provider, model, dimension, status, updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_content_embeddings_content_id
  ON content_embeddings(content_id, updated_at DESC);

UPDATE recommendation_rollout_states
SET embedding_provider = 'embedding_gemma',
    embedding_model = 'embedding-gemma',
    embedding_dimension = 768,
    updated_at = NOW()
WHERE active = TRUE;
