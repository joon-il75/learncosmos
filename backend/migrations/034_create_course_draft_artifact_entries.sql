-- Explorer Diary artifact storage for draft Point outputs.

CREATE TABLE IF NOT EXISTS course_draft_artifact_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_draft_id UUID NOT NULL REFERENCES course_drafts(id) ON DELETE CASCADE,
    course_draft_lesson_id UUID NOT NULL REFERENCES course_draft_lessons(id) ON DELETE CASCADE,
    artifact_type TEXT NOT NULL DEFAULT 'note',
    title TEXT NOT NULL,
    url TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_course_draft_artifact_entries_title
        CHECK (length(btrim(title)) > 0)
);

CREATE INDEX IF NOT EXISTS idx_course_draft_artifact_entries_draft_id
    ON course_draft_artifact_entries(course_draft_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_course_draft_artifact_entries_lesson_id
    ON course_draft_artifact_entries(course_draft_lesson_id);
