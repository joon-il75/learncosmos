CREATE TABLE IF NOT EXISTS course_point_ai_summaries (
    course_point_id      UUID        PRIMARY KEY REFERENCES course_points(id) ON DELETE CASCADE,
    source_title         TEXT        NOT NULL DEFAULT '',
    source_description   TEXT        NOT NULL DEFAULT '',
    summary              TEXT        NOT NULL DEFAULT '',
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
