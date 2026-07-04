ALTER TABLE ai_usage_events
ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}'::jsonb;

ALTER TABLE course_point_ai_summaries
ADD COLUMN IF NOT EXISTS learning_language TEXT NOT NULL DEFAULT 'ko';

ALTER TABLE course_point_ai_summaries
ADD CONSTRAINT course_point_ai_summaries_learning_language_check
CHECK (learning_language IN ('ko', 'en')) NOT VALID;

ALTER TABLE course_point_ai_summaries
VALIDATE CONSTRAINT course_point_ai_summaries_learning_language_check;
