-- Explorer Diary record storage for draft Point progress metrics.

CREATE TABLE IF NOT EXISTS course_draft_record_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_draft_id UUID NOT NULL REFERENCES course_drafts(id) ON DELETE CASCADE,
    course_draft_lesson_id UUID NOT NULL REFERENCES course_draft_lessons(id) ON DELETE CASCADE,
    study_minutes INTEGER NOT NULL DEFAULT 0 CHECK (study_minutes >= 0),
    practice_count INTEGER NOT NULL DEFAULT 0 CHECK (practice_count >= 0),
    confidence_level INTEGER NOT NULL DEFAULT 3 CHECK (confidence_level BETWEEN 1 AND 5),
    application_note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_course_draft_record_entries_non_empty
        CHECK (
            study_minutes > 0
            OR practice_count > 0
            OR length(btrim(application_note)) > 0
            OR confidence_level <> 3
        )
);

CREATE INDEX IF NOT EXISTS idx_course_draft_record_entries_draft_id
    ON course_draft_record_entries(course_draft_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_course_draft_record_entries_lesson_id
    ON course_draft_record_entries(course_draft_lesson_id);
