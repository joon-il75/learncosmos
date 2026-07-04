ALTER TABLE course_goal_profiles
    ADD COLUMN IF NOT EXISTS revision_snapshot JSONB;
