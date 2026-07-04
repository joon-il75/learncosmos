-- Explorer Diary operation memo storage for draft Region/Point details.
-- Notes are intentionally separated from learner-facing description/objective/summary fields.

CREATE TABLE IF NOT EXISTS course_draft_detail_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_draft_id UUID NOT NULL REFERENCES course_drafts(id) ON DELETE CASCADE,
    course_draft_level_id UUID REFERENCES course_draft_levels(id) ON DELETE CASCADE,
    course_draft_lesson_id UUID REFERENCES course_draft_lessons(id) ON DELETE CASCADE,
    target_type TEXT NOT NULL CHECK (target_type IN ('level', 'lesson')),
    note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (
        (
            target_type = 'level'
            AND course_draft_level_id IS NOT NULL
            AND course_draft_lesson_id IS NULL
        )
        OR
        (
            target_type = 'lesson'
            AND course_draft_lesson_id IS NOT NULL
            AND course_draft_level_id IS NULL
        )
    )
);

CREATE INDEX IF NOT EXISTS idx_course_draft_detail_notes_draft_id
    ON course_draft_detail_notes(course_draft_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_course_draft_detail_notes_level
    ON course_draft_detail_notes(course_draft_level_id)
    WHERE target_type = 'level';

CREATE UNIQUE INDEX IF NOT EXISTS idx_course_draft_detail_notes_lesson
    ON course_draft_detail_notes(course_draft_lesson_id)
    WHERE target_type = 'lesson';
