ALTER TABLE IF EXISTS course_drafts
    DROP COLUMN IF EXISTS evaluation_result,
    DROP COLUMN IF EXISTS evaluated_at;

CREATE INDEX IF NOT EXISTS idx_users_totp_reset_requested
    ON users (totp_reset_requested)
    WHERE totp_reset_requested = true;

DO $$
BEGIN
    IF to_regclass('public.content_embeddings') IS NOT NULL THEN
        IF EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'content_embedding_shadows_pkey'
              AND conrelid = 'public.content_embeddings'::regclass
        ) THEN
            ALTER TABLE content_embeddings
                RENAME CONSTRAINT content_embedding_shadows_pkey TO content_embeddings_pkey;
        END IF;

        IF EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'content_embedding_shadows_dimension_check'
              AND conrelid = 'public.content_embeddings'::regclass
        ) THEN
            ALTER TABLE content_embeddings
                RENAME CONSTRAINT content_embedding_shadows_dimension_check TO content_embeddings_dimension_check;
        END IF;

        IF EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'content_embedding_shadows_status_check'
              AND conrelid = 'public.content_embeddings'::regclass
        ) THEN
            ALTER TABLE content_embeddings
                RENAME CONSTRAINT content_embedding_shadows_status_check TO content_embeddings_status_check;
        END IF;

        IF EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'content_embedding_shadows_content_id_fkey'
              AND conrelid = 'public.content_embeddings'::regclass
        ) THEN
            ALTER TABLE content_embeddings
                RENAME CONSTRAINT content_embedding_shadows_content_id_fkey TO content_embeddings_content_id_fkey;
        END IF;
    END IF;
END $$;
