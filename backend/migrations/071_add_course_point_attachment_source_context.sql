ALTER TABLE course_point_attachments
  ADD COLUMN IF NOT EXISTS source_context TEXT NOT NULL DEFAULT 'work_attachment';

ALTER TABLE course_point_attachments
  DROP CONSTRAINT IF EXISTS course_point_attachments_source_context_check;

ALTER TABLE course_point_attachments
  ADD CONSTRAINT course_point_attachments_source_context_check
  CHECK (source_context IN ('research_material', 'work_attachment'));

CREATE INDEX IF NOT EXISTS idx_course_point_attachments_source_context
  ON course_point_attachments(course_point_id, source_context, created_at);
