ALTER TABLE course_drafts
    ADD COLUMN IF NOT EXISTS goal_profile_id UUID REFERENCES course_goal_profiles(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS goal_profile_version INTEGER;

CREATE INDEX IF NOT EXISTS idx_course_drafts_goal_profile_id
    ON course_drafts(goal_profile_id);
