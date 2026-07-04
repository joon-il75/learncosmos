BEGIN;

ALTER TABLE course_goal_profiles
  ADD COLUMN IF NOT EXISTS learning_intent_profile JSONB NOT NULL DEFAULT '{}'::jsonb;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'course_goal_profiles_learning_intent_profile_object'
  ) THEN
    ALTER TABLE course_goal_profiles
      ADD CONSTRAINT course_goal_profiles_learning_intent_profile_object
      CHECK (jsonb_typeof(learning_intent_profile) = 'object');
  END IF;
END $$;

COMMIT;
