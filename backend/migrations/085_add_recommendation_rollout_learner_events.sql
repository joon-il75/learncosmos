CREATE TABLE IF NOT EXISTS recommendation_rollout_learner_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  rollout_state_id UUID REFERENCES recommendation_rollout_states(id) ON DELETE SET NULL,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  course_draft_id UUID REFERENCES course_drafts(id) ON DELETE SET NULL,
  course_point_id UUID REFERENCES course_points(id) ON DELETE SET NULL,
  event_type TEXT NOT NULL,
  rollout_mode TEXT NOT NULL DEFAULT '',
  traffic_percent INTEGER NOT NULL DEFAULT 0,
  bucket INTEGER,
  would_apply BOOLEAN NOT NULL DEFAULT false,
  active_applied BOOLEAN NOT NULL DEFAULT false,
  candidate_count INTEGER NOT NULL DEFAULT 0,
  selected_content_id UUID REFERENCES contents(id) ON DELETE SET NULL,
  selected_url TEXT NOT NULL DEFAULT '',
  report_id UUID REFERENCES course_point_material_reports(id) ON DELETE SET NULL,
  report_type TEXT NOT NULL DEFAULT '',
  request_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
  candidates_snapshot JSONB NOT NULL DEFAULT '[]'::jsonb,
  event_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT recommendation_rollout_learner_events_type_check CHECK (
    event_type IN ('recommendation_exposed', 'material_reported', 'material_replaced', 'material_report_cancelled')
  ),
  CONSTRAINT recommendation_rollout_learner_events_traffic_check CHECK (traffic_percent >= 0 AND traffic_percent <= 100),
  CONSTRAINT recommendation_rollout_learner_events_bucket_check CHECK (bucket IS NULL OR (bucket >= 0 AND bucket < 100)),
  CONSTRAINT recommendation_rollout_learner_events_candidate_count_check CHECK (candidate_count >= 0)
);

CREATE INDEX IF NOT EXISTS idx_recommendation_rollout_learner_events_state_time
  ON recommendation_rollout_learner_events(rollout_state_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_recommendation_rollout_learner_events_type_time
  ON recommendation_rollout_learner_events(event_type, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_recommendation_rollout_learner_events_user_point_time
  ON recommendation_rollout_learner_events(user_id, course_point_id, created_at DESC);
