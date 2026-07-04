ALTER TABLE course_point_attachments
  ADD COLUMN IF NOT EXISTS artifact_id UUID REFERENCES course_point_artifacts(id) ON DELETE CASCADE;

ALTER TABLE course_point_attachments
  DROP CONSTRAINT IF EXISTS course_point_attachments_source_context_check;

ALTER TABLE course_point_attachments
  ADD CONSTRAINT course_point_attachments_source_context_check
  CHECK (source_context IN ('research_material', 'work_attachment', 'artifact'));

ALTER TABLE course_point_attachments
  DROP CONSTRAINT IF EXISTS course_point_attachments_artifact_context_check;

ALTER TABLE course_point_attachments
  ADD CONSTRAINT course_point_attachments_artifact_context_check
  CHECK (
    (source_context = 'artifact' AND artifact_id IS NOT NULL)
    OR (source_context <> 'artifact' AND artifact_id IS NULL)
  );

CREATE INDEX IF NOT EXISTS idx_course_point_attachments_artifact
  ON course_point_attachments(artifact_id, created_at)
  WHERE source_context = 'artifact';

CREATE UNIQUE INDEX IF NOT EXISTS idx_course_point_artifact_video_one
  ON course_point_attachments(artifact_id)
  WHERE source_context = 'artifact'
    AND attachment_type = 'video';

CREATE UNIQUE INDEX IF NOT EXISTS idx_course_point_artifact_subtitle_one
  ON course_point_attachments(artifact_id)
  WHERE source_context = 'artifact'
    AND attachment_type = 'subtitle';

CREATE UNIQUE INDEX IF NOT EXISTS idx_course_point_artifact_thumbnail_one
  ON course_point_attachments(artifact_id)
  WHERE source_context = 'artifact'
    AND attachment_type = 'thumbnail';
