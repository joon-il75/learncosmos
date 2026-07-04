CREATE TABLE IF NOT EXISTS course_point_practice_logs (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_point_id  UUID        NOT NULL REFERENCES course_points(id) ON DELETE CASCADE,
    user_id          UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title            TEXT        NOT NULL,
    attempt_count    INTEGER,
    success_count    INTEGER,
    failure_count    INTEGER,
    duration_minutes INTEGER,
    blocked_part     TEXT,
    changed_method   TEXT,
    achievement_note TEXT,
    next_plan        TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_course_point_practice_logs_point_created
    ON course_point_practice_logs(course_point_id, created_at);
