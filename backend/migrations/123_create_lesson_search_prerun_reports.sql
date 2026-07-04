CREATE TABLE IF NOT EXISTS lesson_search_prerun_reports (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  run_type TEXT NOT NULL,
  status TEXT NOT NULL,
  started_at TIMESTAMPTZ NOT NULL,
  finished_at TIMESTAMPTZ,
  provider_filter TEXT NOT NULL,
  top_n INTEGER NOT NULL,
  risk_top_n INTEGER NOT NULL,
  summary JSONB NOT NULL,
  noise_taxonomy_summary JSONB NOT NULL DEFAULT '{}'::jsonb,
  report JSONB NOT NULL,
  error_code TEXT,
  error_message TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT lesson_search_prerun_reports_status_check CHECK (
    status IN (
      'running',
      'succeeded',
      'failed',
      'skipped_missing_credentials',
      'skipped_lock_not_acquired'
    )
  ),
  CONSTRAINT lesson_search_prerun_reports_summary_object_check CHECK (jsonb_typeof(summary) = 'object'),
  CONSTRAINT lesson_search_prerun_reports_noise_taxonomy_object_check CHECK (jsonb_typeof(noise_taxonomy_summary) = 'object'),
  CONSTRAINT lesson_search_prerun_reports_report_object_check CHECK (jsonb_typeof(report) = 'object'),
  CONSTRAINT lesson_search_prerun_reports_top_n_check CHECK (top_n > 0),
  CONSTRAINT lesson_search_prerun_reports_risk_top_n_check CHECK (risk_top_n >= 0)
);

CREATE INDEX IF NOT EXISTS idx_lesson_search_prerun_reports_created_at
  ON lesson_search_prerun_reports (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_lesson_search_prerun_reports_status_created_at
  ON lesson_search_prerun_reports (status, created_at DESC);
