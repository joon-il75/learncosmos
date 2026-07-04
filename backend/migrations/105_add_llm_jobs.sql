CREATE TABLE IF NOT EXISTS llm_jobs (
  id uuid PRIMARY KEY,
  user_id uuid REFERENCES users(id),
  feature text NOT NULL,
  idempotency_key text,
  status text NOT NULL,
  priority int NOT NULL DEFAULT 0,
  phase text,
  request_ref jsonb NOT NULL DEFAULT '{}',
  prompt_template_key text,
  prompt_input_ref jsonb NOT NULL DEFAULT '{}',
  provider text,
  model text,
  provider_mode text,
  api_key_ref text,
  result_ref jsonb NOT NULL DEFAULT '{}',
  error_code text,
  error_message text,
  retryable boolean NOT NULL DEFAULT false,
  point_cost int,
  billing_status text,
  queue_wait_ms int,
  provider_wait_ms int,
  llm_generation_elapsed_ms int,
  parse_validate_elapsed_ms int,
  persist_elapsed_ms int,
  total_job_elapsed_ms int,
  provider_inflight_count int,
  provider_queue_depth int,
  worker_id text,
  attempts int NOT NULL DEFAULT 0,
  created_at timestamptz NOT NULL DEFAULT now(),
  queued_at timestamptz,
  started_at timestamptz,
  finished_at timestamptz,
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT llm_jobs_status_check CHECK (status IN ('queued', 'running', 'succeeded', 'failed', 'canceled', 'expired')),
  CONSTRAINT llm_jobs_feature_not_blank CHECK (length(btrim(feature)) > 0)
);

CREATE INDEX IF NOT EXISTS idx_llm_jobs_status_created_at
  ON llm_jobs(status, created_at);

CREATE INDEX IF NOT EXISTS idx_llm_jobs_user_created_at
  ON llm_jobs(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_llm_jobs_feature_status
  ON llm_jobs(feature, status, created_at);

CREATE UNIQUE INDEX IF NOT EXISTS idx_llm_jobs_idempotency_key
  ON llm_jobs(feature, idempotency_key)
  WHERE idempotency_key IS NOT NULL AND status IN ('queued', 'running', 'succeeded');
