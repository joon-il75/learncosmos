ALTER TABLE contents
ADD COLUMN IF NOT EXISTS external_source TEXT,
ADD COLUMN IF NOT EXISTS external_content_id TEXT,
ADD COLUMN IF NOT EXISTS canonical_url TEXT;

UPDATE contents
SET canonical_url = LOWER(REGEXP_REPLACE(TRIM(TRAILING '/' FROM SPLIT_PART(COALESCE(url, ''), '#', 1)), '^https?://', 'https://'))
WHERE COALESCE(url, '') <> ''
  AND COALESCE(canonical_url, '') = '';

CREATE INDEX IF NOT EXISTS idx_contents_external_source_id
ON contents(external_source, external_content_id);

CREATE INDEX IF NOT EXISTS idx_contents_canonical_url
ON contents(canonical_url);
