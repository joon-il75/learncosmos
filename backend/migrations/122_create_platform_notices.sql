CREATE TABLE IF NOT EXISTS platform_notices (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  slug TEXT NOT NULL UNIQUE,
  locale TEXT NOT NULL DEFAULT 'ko',
  title TEXT NOT NULL,
  summary TEXT NOT NULL DEFAULT '',
  body TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'draft',
  pinned BOOLEAN NOT NULL DEFAULT FALSE,
  published_at TIMESTAMPTZ,
  created_by TEXT NOT NULL DEFAULT '',
  updated_by TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT platform_notices_locale_check CHECK (locale IN ('ko', 'en')),
  CONSTRAINT platform_notices_status_check CHECK (status IN ('draft', 'published', 'archived')),
  CONSTRAINT platform_notices_title_length_check CHECK (char_length(trim(title)) BETWEEN 1 AND 160),
  CONSTRAINT platform_notices_summary_length_check CHECK (char_length(summary) <= 500),
  CONSTRAINT platform_notices_body_length_check CHECK (char_length(body) <= 20000)
);

CREATE INDEX IF NOT EXISTS idx_platform_notices_public
  ON platform_notices(locale, pinned DESC, published_at DESC, created_at DESC)
  WHERE status = 'published';

CREATE INDEX IF NOT EXISTS idx_platform_notices_admin
  ON platform_notices(status, locale, updated_at DESC);
