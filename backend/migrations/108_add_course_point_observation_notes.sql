CREATE TABLE IF NOT EXISTS course_point_observation_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_point_id UUID NOT NULL REFERENCES course_points(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    note_type TEXT NOT NULL,
    content TEXT NOT NULL,
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_course_point_observation_notes_type
        CHECK (note_type IN ('core_summary', 'revisit_part', 'reference_material')),
    CONSTRAINT chk_course_point_observation_notes_content
        CHECK (length(trim(content)) > 0)
);

CREATE INDEX IF NOT EXISTS idx_course_point_observation_notes_point
    ON course_point_observation_notes(course_point_id, order_index, created_at);

CREATE INDEX IF NOT EXISTS idx_course_point_observation_notes_user
    ON course_point_observation_notes(user_id, created_at DESC);
