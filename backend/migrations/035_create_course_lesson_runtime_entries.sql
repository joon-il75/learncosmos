-- Learning planet lesson runtime state for actual progress/completion tracking.

CREATE TABLE IF NOT EXISTS course_lesson_runtime_entries (
    id                UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id         UUID        NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    course_lesson_id  UUID        NOT NULL REFERENCES course_lessons(id) ON DELETE CASCADE,
    status            TEXT        NOT NULL CHECK (status IN ('in_progress', 'completed')),
    started_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at      TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, course_lesson_id)
);

CREATE INDEX IF NOT EXISTS idx_course_lesson_runtime_entries_course_id
    ON course_lesson_runtime_entries(course_id);

CREATE INDEX IF NOT EXISTS idx_course_lesson_runtime_entries_lesson_id
    ON course_lesson_runtime_entries(course_lesson_id);
