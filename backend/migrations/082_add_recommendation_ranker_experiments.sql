CREATE TABLE IF NOT EXISTS recommendation_debug_scenarios (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_by UUID REFERENCES users(id) ON DELETE SET NULL,
  course_title TEXT NOT NULL DEFAULT '',
  initial_user_intent TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'draft',
  goal_profile_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
  generated_lessons_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
  notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT recommendation_debug_scenarios_status_check CHECK (
    status IN ('draft', 'goal_confirmed', 'lessons_generated', 'recommendation_tested', 'labelled', 'approved', 'archived')
  )
);

CREATE INDEX IF NOT EXISTS idx_recommendation_debug_scenarios_status
  ON recommendation_debug_scenarios(status, updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_recommendation_debug_scenarios_created_by
  ON recommendation_debug_scenarios(created_by, updated_at DESC);

CREATE TABLE IF NOT EXISTS recommendation_debug_runs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  scenario_id UUID REFERENCES recommendation_debug_scenarios(id) ON DELETE CASCADE,
  created_by UUID REFERENCES users(id) ON DELETE SET NULL,
  run_type TEXT NOT NULL DEFAULT 'recommendation_compare',
  request_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
  baseline_result_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
  shadow_result_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
  feature_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
  metrics_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
  provider_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
  latency_ms INTEGER,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT recommendation_debug_runs_type_check CHECK (
    run_type IN ('goal_chat', 'lesson_generation', 'recommendation_compare', 'embedding_shadow', 'ranker_shadow')
  ),
  CONSTRAINT recommendation_debug_runs_latency_check CHECK (latency_ms IS NULL OR latency_ms >= 0)
);

CREATE INDEX IF NOT EXISTS idx_recommendation_debug_runs_scenario_id
  ON recommendation_debug_runs(scenario_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_recommendation_debug_runs_type
  ON recommendation_debug_runs(run_type, created_at DESC);

CREATE TABLE IF NOT EXISTS recommendation_debug_labels (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  scenario_id UUID REFERENCES recommendation_debug_scenarios(id) ON DELETE CASCADE,
  run_id UUID REFERENCES recommendation_debug_runs(id) ON DELETE CASCADE,
  created_by UUID REFERENCES users(id) ON DELETE SET NULL,
  candidate_key TEXT NOT NULL DEFAULT '',
  content_id UUID REFERENCES contents(id) ON DELETE SET NULL,
  url TEXT NOT NULL DEFAULT '',
  label TEXT NOT NULL,
  note TEXT NOT NULL DEFAULT '',
  feature_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
  baseline_rank INTEGER,
  shadow_rank INTEGER,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT recommendation_debug_labels_label_check CHECK (
    label IN ('good_fit', 'goal_fit_stage_weak', 'stage_fit_goal_weak', 'irrelevant', 'duplicate', 'broken_link', 'low_quality', 'unsafe')
  ),
  CONSTRAINT recommendation_debug_labels_baseline_rank_check CHECK (baseline_rank IS NULL OR baseline_rank > 0),
  CONSTRAINT recommendation_debug_labels_shadow_rank_check CHECK (shadow_rank IS NULL OR shadow_rank > 0)
);

CREATE INDEX IF NOT EXISTS idx_recommendation_debug_labels_scenario_id
  ON recommendation_debug_labels(scenario_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_recommendation_debug_labels_run_id
  ON recommendation_debug_labels(run_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_recommendation_debug_labels_label
  ON recommendation_debug_labels(label, created_at DESC);

CREATE TABLE IF NOT EXISTS recommendation_rollout_states (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  mode TEXT NOT NULL DEFAULT 'shadow_only',
  traffic_percent INTEGER NOT NULL DEFAULT 0,
  embedding_provider TEXT NOT NULL DEFAULT 'openai',
  embedding_model TEXT NOT NULL DEFAULT 'text-embedding-3-small',
  embedding_dimension INTEGER NOT NULL DEFAULT 1536,
  ranker_model_version TEXT NOT NULL DEFAULT '',
  feature_schema_version TEXT NOT NULL DEFAULT 'builtin-v1',
  quality_gate_status TEXT NOT NULL DEFAULT 'pending',
  quality_gate_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
  rollback_policy_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
  approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
  approved_at TIMESTAMPTZ,
  activated_at TIMESTAMPTZ,
  rolled_back_at TIMESTAMPTZ,
  rollback_reason TEXT NOT NULL DEFAULT '',
  active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT recommendation_rollout_states_mode_check CHECK (
    mode IN ('shadow_only', 'admin_preview', 'learner_10_percent', 'learner_50_percent', 'learner_100_percent', 'rolled_back')
  ),
  CONSTRAINT recommendation_rollout_states_traffic_check CHECK (traffic_percent >= 0 AND traffic_percent <= 100),
  CONSTRAINT recommendation_rollout_states_dimension_check CHECK (embedding_dimension > 0),
  CONSTRAINT recommendation_rollout_states_quality_check CHECK (
    quality_gate_status IN ('pending', 'passed', 'failed', 'warning')
  )
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_recommendation_rollout_states_one_active
  ON recommendation_rollout_states(active)
  WHERE active = true;

CREATE INDEX IF NOT EXISTS idx_recommendation_rollout_states_mode
  ON recommendation_rollout_states(mode, updated_at DESC);

CREATE TABLE IF NOT EXISTS recommendation_rollout_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  rollout_state_id UUID REFERENCES recommendation_rollout_states(id) ON DELETE CASCADE,
  created_by UUID REFERENCES users(id) ON DELETE SET NULL,
  event_type TEXT NOT NULL,
  from_mode TEXT NOT NULL DEFAULT '',
  to_mode TEXT NOT NULL DEFAULT '',
  payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_recommendation_rollout_events_state_id
  ON recommendation_rollout_events(rollout_state_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_recommendation_rollout_events_type
  ON recommendation_rollout_events(event_type, created_at DESC);

CREATE TABLE IF NOT EXISTS recommendation_rollout_metrics (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  rollout_state_id UUID REFERENCES recommendation_rollout_states(id) ON DELETE CASCADE,
  window_started_at TIMESTAMPTZ NOT NULL,
  window_ended_at TIMESTAMPTZ NOT NULL,
  recommendation_count INTEGER NOT NULL DEFAULT 0,
  selection_rate NUMERIC(8, 4),
  broken_link_rate NUMERIC(8, 4),
  wrong_content_rate NUMERIC(8, 4),
  fallback_rate NUMERIC(8, 4),
  p95_latency_ms INTEGER,
  top5_overlap NUMERIC(8, 4),
  good_fit_rate NUMERIC(8, 4),
  irrelevant_rate NUMERIC(8, 4),
  duplicate_rate NUMERIC(8, 4),
  metrics_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT recommendation_rollout_metrics_window_check CHECK (window_ended_at >= window_started_at),
  CONSTRAINT recommendation_rollout_metrics_count_check CHECK (recommendation_count >= 0),
  CONSTRAINT recommendation_rollout_metrics_latency_check CHECK (p95_latency_ms IS NULL OR p95_latency_ms >= 0)
);

CREATE INDEX IF NOT EXISTS idx_recommendation_rollout_metrics_state_id
  ON recommendation_rollout_metrics(rollout_state_id, window_ended_at DESC);

INSERT INTO recommendation_rollout_states (
  mode,
  traffic_percent,
  embedding_provider,
  embedding_model,
  embedding_dimension,
  feature_schema_version,
  quality_gate_status,
  rollback_policy_snapshot,
  active
)
SELECT
  'shadow_only',
  0,
  'openai',
  'text-embedding-3-small',
  1536,
  'builtin-v1',
  'pending',
  '{"auto_rollback": true, "learner_expansion_requires_admin_approval": true}'::jsonb,
  true
WHERE NOT EXISTS (
  SELECT 1 FROM recommendation_rollout_states WHERE active = true
);
