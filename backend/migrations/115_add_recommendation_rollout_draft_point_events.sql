ALTER TABLE recommendation_rollout_learner_events
  ADD COLUMN IF NOT EXISTS course_draft_point_id UUID REFERENCES course_draft_points(id) ON DELETE SET NULL,
  ADD COLUMN IF NOT EXISTS point_target_type TEXT NOT NULL DEFAULT '';

ALTER TABLE recommendation_rollout_learner_events
  DROP CONSTRAINT IF EXISTS recommendation_rollout_learner_events_point_target_type_check;

ALTER TABLE recommendation_rollout_learner_events
  ADD CONSTRAINT recommendation_rollout_learner_events_point_target_type_check CHECK (
    point_target_type IN ('', 'course_point', 'course_draft_point')
  );

CREATE INDEX IF NOT EXISTS idx_recommendation_rollout_learner_events_user_draft_point_time
  ON recommendation_rollout_learner_events(user_id, course_draft_point_id, created_at DESC);
