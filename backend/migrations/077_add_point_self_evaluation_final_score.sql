-- Store the learner's final self-rating selected after AI self-evaluation.

ALTER TABLE course_point_self_evaluations
  ADD COLUMN IF NOT EXISTS final_score INT;

ALTER TABLE course_point_self_evaluations
  DROP CONSTRAINT IF EXISTS course_point_self_evaluations_final_score_check;

ALTER TABLE course_point_self_evaluations
  ADD CONSTRAINT course_point_self_evaluations_final_score_check
  CHECK (final_score IS NULL OR final_score BETWEEN 1 AND 5);
