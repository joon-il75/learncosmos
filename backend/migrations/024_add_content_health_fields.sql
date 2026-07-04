ALTER TABLE contents
ADD COLUMN IF NOT EXISTS content_status VARCHAR(20) NOT NULL DEFAULT 'active',
ADD COLUMN IF NOT EXISTS health_score NUMERIC(3,2) NOT NULL DEFAULT 1.0,
ADD COLUMN IF NOT EXISTS last_checked_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS http_status INT;

CREATE INDEX IF NOT EXISTS idx_contents_health_active
ON contents(content_status, health_score);

UPDATE contents
SET content_status = COALESCE(NULLIF(content_status, ''), 'active'),
    health_score = COALESCE(health_score, 1.0)
WHERE content_status IS NULL
   OR content_status = ''
   OR health_score IS NULL;
