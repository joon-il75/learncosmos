CREATE TABLE IF NOT EXISTS ai_usage_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source TEXT NOT NULL DEFAULT 'byok',
    provider TEXT NOT NULL,
    model TEXT,
    feature TEXT NOT NULL,
    billing_status TEXT,
    input_tokens INTEGER,
    output_tokens INTEGER,
    estimated_cost_usd NUMERIC(12, 6),
    success BOOLEAN NOT NULL DEFAULT TRUE,
    error_code TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_usage_events_user_created
    ON ai_usage_events(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_ai_usage_events_source_created
    ON ai_usage_events(source, created_at DESC);
