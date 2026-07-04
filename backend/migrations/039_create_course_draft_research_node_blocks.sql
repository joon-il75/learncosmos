CREATE TABLE IF NOT EXISTS course_draft_research_node_blocks (
    id               uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    research_node_id uuid        NOT NULL REFERENCES course_draft_research_nodes(id) ON DELETE CASCADE,
    block_type       varchar(20) NOT NULL
        CHECK (block_type IN ('text', 'image', 'link')),
    content          jsonb       NOT NULL DEFAULT '{}',
    order_index      int         NOT NULL DEFAULT 0,
    created_at       timestamptz NOT NULL DEFAULT NOW(),
    updated_at       timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_course_draft_research_node_blocks_node_id
    ON course_draft_research_node_blocks (research_node_id);
