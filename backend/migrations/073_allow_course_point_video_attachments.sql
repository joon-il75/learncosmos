ALTER TABLE course_point_attachments
  DROP CONSTRAINT IF EXISTS course_point_attachments_type_check;

ALTER TABLE course_point_attachments
  ADD CONSTRAINT course_point_attachments_type_check
  CHECK (attachment_type IN ('image', 'file', 'link', 'code', 'other', 'video'));

CREATE UNIQUE INDEX IF NOT EXISTS idx_course_point_research_material_video_one
  ON course_point_attachments(course_point_id)
  WHERE source_context = 'research_material'
    AND attachment_type = 'video';
