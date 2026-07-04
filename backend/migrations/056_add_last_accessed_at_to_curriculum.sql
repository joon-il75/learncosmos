ALTER TABLE course_drafts
ADD COLUMN IF NOT EXISTS last_accessed_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_course_drafts_user_status_last_accessed
ON course_drafts (user_id, status, last_accessed_at DESC, updated_at DESC);
