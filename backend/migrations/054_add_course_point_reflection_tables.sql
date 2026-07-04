BEGIN;

CREATE TABLE IF NOT EXISTS course_point_questions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  course_point_id UUID NOT NULL REFERENCES course_points(id) ON DELETE CASCADE,
  goal_profile_version INTEGER,
  question TEXT NOT NULL,
  question_type TEXT NOT NULL DEFAULT 'reflection',
  answer TEXT,
  status TEXT NOT NULL DEFAULT 'pending',
  created_by TEXT NOT NULL DEFAULT 'learner',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT course_point_questions_question_type_check CHECK (question_type IN ('reflection', 'application', 'goal_alignment')),
  CONSTRAINT course_point_questions_status_check CHECK (status IN ('pending', 'answered'))
);

CREATE INDEX IF NOT EXISTS idx_course_point_questions_course_point_id
  ON course_point_questions(course_point_id, created_at);

CREATE TABLE IF NOT EXISTS course_point_self_evaluations (
  course_point_id UUID PRIMARY KEY REFERENCES course_points(id) ON DELETE CASCADE,
  goal_profile_version INTEGER,
  understanding INTEGER NOT NULL,
  application_note TEXT NOT NULL DEFAULT '',
  proficiency INTEGER NOT NULL,
  goal_alignment_note TEXT NOT NULL DEFAULT '',
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT course_point_self_evaluations_understanding_check CHECK (understanding BETWEEN 1 AND 5),
  CONSTRAINT course_point_self_evaluations_proficiency_check CHECK (proficiency BETWEEN 1 AND 5)
);

COMMIT;
