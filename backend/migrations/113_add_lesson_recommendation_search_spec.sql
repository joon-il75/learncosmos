ALTER TABLE course_draft_lessons
  ADD COLUMN IF NOT EXISTS recommendation_search_spec JSONB NOT NULL DEFAULT '{}'::jsonb;

ALTER TABLE course_lessons
  ADD COLUMN IF NOT EXISTS recommendation_search_spec JSONB NOT NULL DEFAULT '{}'::jsonb;
