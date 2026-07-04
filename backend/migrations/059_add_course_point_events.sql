CREATE TABLE IF NOT EXISTS course_point_events (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_point_id UUID        NOT NULL REFERENCES course_points(id) ON DELETE CASCADE,
    user_id         UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_type      TEXT        NOT NULL,
    event_payload   JSONB       NOT NULL DEFAULT '{}'::jsonb,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_course_point_events_point_created
    ON course_point_events(course_point_id, created_at);
