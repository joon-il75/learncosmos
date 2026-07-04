CREATE TABLE IF NOT EXISTS course_draft_research_nodes (
    id                   uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_draft_id      uuid        NOT NULL REFERENCES course_drafts(id) ON DELETE CASCADE,
    course_draft_level_id uuid       NOT NULL REFERENCES course_draft_levels(id) ON DELETE CASCADE,
    title                varchar(200) NOT NULL,
    template_type        varchar(50) NOT NULL DEFAULT 'free_research'
        CHECK (template_type IN ('concept_summary', 'practice_strategy', 'problem_solving', 'free_research')),
    order_index          int         NOT NULL DEFAULT 0,
    created_at           timestamptz NOT NULL DEFAULT NOW(),
    updated_at           timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_course_draft_research_nodes_draft_id
    ON course_draft_research_nodes (course_draft_id);

CREATE INDEX IF NOT EXISTS idx_course_draft_research_nodes_level_id
    ON course_draft_research_nodes (course_draft_level_id);
