CREATE TABLE IF NOT EXISTS course_point_learning_sessions (
    id                   UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_point_id      UUID        NOT NULL REFERENCES course_points(id) ON DELETE CASCADE,
    user_id              UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    started_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at             TIMESTAMPTZ,
    active_seconds       INTEGER     NOT NULL DEFAULT 0,
    heartbeat_count      INTEGER     NOT NULL DEFAULT 0,
    end_reason           TEXT,
    last_visibility      TEXT        NOT NULL DEFAULT 'visible',
    user_agent           TEXT,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_course_point_learning_sessions_point_started
    ON course_point_learning_sessions(course_point_id, started_at DESC);

CREATE INDEX IF NOT EXISTS idx_course_point_learning_sessions_user_started
    ON course_point_learning_sessions(user_id, started_at DESC);

CREATE INDEX IF NOT EXISTS idx_course_point_learning_sessions_open
    ON course_point_learning_sessions(course_point_id, user_id, last_seen_at DESC)
    WHERE ended_at IS NULL;
