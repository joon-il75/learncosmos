BEGIN;

ALTER TABLE course_goal_profiles
  ALTER COLUMN course_draft_id DROP NOT NULL;

CREATE INDEX IF NOT EXISTS idx_goal_profiles_user_pending_active
  ON course_goal_profiles(user_id, updated_at DESC)
  WHERE is_active = true AND course_draft_id IS NULL;

COMMIT;
