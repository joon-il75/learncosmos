ALTER TABLE course_point_material_reports
  ADD COLUMN IF NOT EXISTS admin_note TEXT NOT NULL DEFAULT '';

ALTER TABLE course_point_material_reports
  ADD COLUMN IF NOT EXISTS reviewed_by TEXT;

ALTER TABLE course_point_material_reports
  DROP CONSTRAINT IF EXISTS course_point_material_reports_reviewed_by_fkey;

ALTER TABLE course_point_material_reports
  ALTER COLUMN reviewed_by TYPE TEXT USING reviewed_by::text;

ALTER TABLE course_point_material_reports
  ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_course_point_material_reports_reviewed_by
  ON course_point_material_reports(reviewed_by, reviewed_at DESC);
