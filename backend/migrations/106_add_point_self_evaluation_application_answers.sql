CREATE TABLE IF NOT EXISTS course_point_self_evaluation_application_answers (
  id uuid PRIMARY KEY,
  course_point_id uuid NOT NULL REFERENCES course_points(id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  goal_profile_version int,
  answers jsonb NOT NULL DEFAULT '[]',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT course_point_self_eval_app_answers_array CHECK (jsonb_typeof(answers) = 'array')
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_course_point_self_eval_app_answers_point
  ON course_point_self_evaluation_application_answers(course_point_id);

CREATE INDEX IF NOT EXISTS idx_course_point_self_eval_app_answers_user_updated
  ON course_point_self_evaluation_application_answers(user_id, updated_at DESC);
