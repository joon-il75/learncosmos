CREATE TABLE IF NOT EXISTS content_recommendation_blocks (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  content_id UUID REFERENCES contents(id) ON DELETE SET NULL,
  url TEXT,
  url_key TEXT,
  source_report_id UUID REFERENCES course_point_material_reports(id) ON DELETE SET NULL,
  reason TEXT NOT NULL DEFAULT '',
  created_by TEXT,
  active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT content_recommendation_blocks_target_check CHECK (
    content_id IS NOT NULL OR COALESCE(url_key, '') <> ''
  )
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_content_recommendation_blocks_content_active
  ON content_recommendation_blocks(content_id)
  WHERE active = true AND content_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_content_recommendation_blocks_url_active
  ON content_recommendation_blocks(url_key)
  WHERE active = true AND COALESCE(url_key, '') <> '';

CREATE INDEX IF NOT EXISTS idx_content_recommendation_blocks_source_report
  ON content_recommendation_blocks(source_report_id);
