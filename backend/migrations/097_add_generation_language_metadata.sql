ALTER TABLE course_goal_profiles
    ADD COLUMN IF NOT EXISTS language TEXT CHECK (language IN ('ko', 'en')) DEFAULT 'ko';

ALTER TABLE course_drafts
    ADD COLUMN IF NOT EXISTS generation_language TEXT CHECK (generation_language IN ('ko', 'en')) DEFAULT 'ko';

UPDATE course_goal_profiles
SET language = 'ko'
WHERE language IS NULL;

UPDATE course_drafts
SET generation_language = 'ko'
WHERE generation_language IS NULL;
