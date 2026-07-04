CREATE EXTENSION IF NOT EXISTS pg_trgm;

ALTER TABLE contents
    ADD COLUMN IF NOT EXISTS search_text_ko TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_contents_search_text_ko_trgm
ON contents
USING gin (search_text_ko gin_trgm_ops);

UPDATE contents
SET search_text_ko = lower(
    concat_ws(
        ' ',
        COALESCE(title, ''),
        COALESCE(description, ''),
        COALESCE(author, ''),
        COALESCE(content_type, ''),
        COALESCE(language, '')
    )
)
WHERE search_text_ko = '';
