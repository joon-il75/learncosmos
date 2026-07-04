ALTER TABLE course_point_material_reports
  ADD COLUMN IF NOT EXISTS target_content_id UUID REFERENCES contents(id) ON DELETE SET NULL,
  ADD COLUMN IF NOT EXISTS target_url TEXT,
  ADD COLUMN IF NOT EXISTS target_title TEXT,
  ADD COLUMN IF NOT EXISTS replacement_content_id UUID REFERENCES contents(id) ON DELETE SET NULL,
  ADD COLUMN IF NOT EXISTS replacement_url TEXT,
  ADD COLUMN IF NOT EXISTS replacement_title TEXT,
  ADD COLUMN IF NOT EXISTS replaced_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_course_point_material_reports_target_content_id
  ON course_point_material_reports(target_content_id)
  WHERE target_content_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_course_point_material_reports_replaced_at
  ON course_point_material_reports(replaced_at DESC)
  WHERE replaced_at IS NOT NULL;
