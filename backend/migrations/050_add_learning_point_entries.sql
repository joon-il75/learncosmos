CREATE TABLE IF NOT EXISTS course_point_journal_entries (
    course_point_id UUID        PRIMARY KEY REFERENCES course_points(id) ON DELETE CASCADE,
    observation     TEXT        NOT NULL DEFAULT '',
    reflection      TEXT        NOT NULL DEFAULT '',
    next_step       TEXT        NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS course_point_record_entries (
    course_point_id   UUID        PRIMARY KEY REFERENCES course_points(id) ON DELETE CASCADE,
    study_minutes     INTEGER     NOT NULL DEFAULT 0,
    practice_count    INTEGER     NOT NULL DEFAULT 0,
    confidence_level  INTEGER     NOT NULL DEFAULT 3,
    application_note  TEXT        NOT NULL DEFAULT '',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS course_point_artifact_entries (
    course_point_id UUID        PRIMARY KEY REFERENCES course_points(id) ON DELETE CASCADE,
    artifact_type   TEXT        NOT NULL DEFAULT '',
    title           TEXT        NOT NULL DEFAULT '',
    url             TEXT        NOT NULL DEFAULT '',
    description     TEXT        NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
