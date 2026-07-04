CREATE TABLE IF NOT EXISTS course_research_nodes (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id        UUID        NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    course_lesson_id UUID        NOT NULL REFERENCES course_lessons(id) ON DELETE CASCADE,
    title            TEXT        NOT NULL,
    template_type    TEXT        NOT NULL CHECK (template_type IN (
        'concept_summary',
        'practice_strategy',
        'problem_solving',
        'free_research'
    )),
    order_index      INTEGER     NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_course_research_nodes_course_id
    ON course_research_nodes(course_id);
CREATE INDEX IF NOT EXISTS idx_course_research_nodes_lesson_id
    ON course_research_nodes(course_lesson_id);

CREATE TABLE IF NOT EXISTS course_research_node_blocks (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    research_node_id UUID        NOT NULL REFERENCES course_research_nodes(id) ON DELETE CASCADE,
    user_id          UUID        REFERENCES users(id) ON DELETE SET NULL,
    block_type       TEXT        NOT NULL CHECK (block_type IN ('text', 'image', 'link')),
    content          JSONB       NOT NULL DEFAULT '{}'::jsonb,
    order_index      INTEGER     NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (research_node_id, order_index)
);

CREATE INDEX IF NOT EXISTS idx_course_research_node_blocks_research_node_id
    ON course_research_node_blocks(research_node_id);
