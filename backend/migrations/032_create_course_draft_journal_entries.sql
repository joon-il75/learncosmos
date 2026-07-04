-- Explorer Diary journal storage for draft Point records.

CREATE TABLE IF NOT EXISTS course_draft_journal_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_draft_id UUID NOT NULL REFERENCES course_drafts(id) ON DELETE CASCADE,
    course_draft_lesson_id UUID NOT NULL REFERENCES course_draft_lessons(id) ON DELETE CASCADE,
    observation TEXT NOT NULL DEFAULT '',
    reflection TEXT NOT NULL DEFAULT '',
    next_step TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_course_draft_journal_entries_non_empty
        CHECK (
            length(btrim(observation)) > 0
            OR length(btrim(reflection)) > 0
            OR length(btrim(next_step)) > 0
        )
);

CREATE INDEX IF NOT EXISTS idx_course_draft_journal_entries_draft_id
    ON course_draft_journal_entries(course_draft_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_course_draft_journal_entries_lesson_id
    ON course_draft_journal_entries(course_draft_lesson_id);
