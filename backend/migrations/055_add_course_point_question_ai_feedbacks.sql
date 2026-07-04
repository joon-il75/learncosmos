BEGIN;

CREATE TABLE IF NOT EXISTS course_point_question_ai_feedbacks (
  course_point_question_id UUID PRIMARY KEY REFERENCES course_point_questions(id) ON DELETE CASCADE,
  feedback TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_course_point_question_ai_feedbacks_updated_at
  ON course_point_question_ai_feedbacks(updated_at DESC);

COMMIT;
