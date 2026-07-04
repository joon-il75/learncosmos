-- Extend learning point runtime data for the redesigned exploration point screen.
-- This migration is intentionally additive so existing learning records remain intact.

ALTER TABLE course_points
  ADD COLUMN IF NOT EXISTS point_goal TEXT,
  ADD COLUMN IF NOT EXISTS point_category TEXT;

ALTER TABLE course_point_journal_entries
  ADD COLUMN IF NOT EXISTS core_concept TEXT,
  ADD COLUMN IF NOT EXISTS my_explanation TEXT,
  ADD COLUMN IF NOT EXISTS examples TEXT,
  ADD COLUMN IF NOT EXISTS confused_parts TEXT,
  ADD COLUMN IF NOT EXISTS reference_links TEXT;

ALTER TABLE course_point_questions
  ADD COLUMN IF NOT EXISTS title TEXT,
  ADD COLUMN IF NOT EXISTS answer_method TEXT;

ALTER TABLE course_point_practice_logs
  ADD COLUMN IF NOT EXISTS activity_name TEXT,
  ADD COLUMN IF NOT EXISTS achievement TEXT,
  ADD COLUMN IF NOT EXISTS next_practice TEXT;

ALTER TABLE course_point_artifacts
  ADD COLUMN IF NOT EXISTS point_category TEXT,
  ADD COLUMN IF NOT EXISTS production_process TEXT,
  ADD COLUMN IF NOT EXISTS learned_points TEXT,
  ADD COLUMN IF NOT EXISTS difficult_points TEXT,
  ADD COLUMN IF NOT EXISTS visibility TEXT NOT NULL DEFAULT 'private';

CREATE TABLE IF NOT EXISTS course_point_attachments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  course_point_id UUID NOT NULL REFERENCES course_points(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id),
  provider TEXT NOT NULL,
  attachment_type TEXT NOT NULL,
  title TEXT,
  url TEXT,
  file_path TEXT,
  file_size BIGINT,
  mime_type TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT course_point_attachments_provider_check CHECK (provider IN ('creator', 'platform', 'learner')),
  CONSTRAINT course_point_attachments_type_check CHECK (attachment_type IN ('image', 'file', 'link', 'code', 'other'))
);

CREATE INDEX IF NOT EXISTS idx_course_point_attachments_point_id
  ON course_point_attachments(course_point_id, created_at);

CREATE INDEX IF NOT EXISTS idx_course_point_attachments_user_id
  ON course_point_attachments(user_id, created_at);
