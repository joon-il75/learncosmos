CREATE TABLE IF NOT EXISTS course_point_material_reports (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  course_point_id UUID NOT NULL REFERENCES course_points(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id),
  attachment_id UUID REFERENCES course_point_attachments(id) ON DELETE SET NULL,
  target_type TEXT NOT NULL,
  report_type TEXT NOT NULL,
  message TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'open',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT course_point_material_reports_target_type_check CHECK (target_type IN ('source', 'ai_summary', 'attachment', 'other')),
  CONSTRAINT course_point_material_reports_report_type_check CHECK (report_type IN ('broken_link', 'wrong_content', 'unsafe_content', 'copyright', 'low_quality', 'other')),
  CONSTRAINT course_point_material_reports_status_check CHECK (status IN ('open', 'reviewing', 'resolved', 'dismissed'))
);

CREATE INDEX IF NOT EXISTS idx_course_point_material_reports_point_id
  ON course_point_material_reports(course_point_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_course_point_material_reports_user_id
  ON course_point_material_reports(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_course_point_material_reports_status
  ON course_point_material_reports(status, created_at DESC);
