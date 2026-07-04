ALTER TABLE course_drafts
ADD COLUMN IF NOT EXISTS completion_criteria TEXT[] NOT NULL DEFAULT '{}'::text[];

ALTER TABLE courses
ADD COLUMN IF NOT EXISTS completion_criteria TEXT[] NOT NULL DEFAULT '{}'::text[];
