BEGIN;

-- Explorer Plan: Region / SubRegion / Node
-- 기준:
-- - 상위 편집 단위는 confirmed course가 아니라 course_drafts
-- - 소프트 삭제는 row 삭제가 아니라 status='inactive'
-- - 서브지역 최대 3개 제약은 DB가 아니라 Service 레이어에서 검증

CREATE TABLE explorer_regions (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  course_draft_id  UUID NOT NULL REFERENCES course_drafts(id) ON DELETE CASCADE,
  name             TEXT NOT NULL,
  description      TEXT,
  order_index      INTEGER NOT NULL DEFAULT 0,
  status           TEXT NOT NULL DEFAULT 'active'
                   CHECK (status IN ('active', 'inactive')),
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_explorer_regions_course_draft_id
  ON explorer_regions(course_draft_id);

CREATE INDEX idx_explorer_regions_status
  ON explorer_regions(status);

CREATE TABLE explorer_subregions (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  region_id     UUID NOT NULL REFERENCES explorer_regions(id) ON DELETE CASCADE,
  name          TEXT NOT NULL,
  description   TEXT,
  order_index   INTEGER NOT NULL DEFAULT 0,
  status        TEXT NOT NULL DEFAULT 'active'
                CHECK (status IN ('active', 'inactive')),
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_explorer_subregions_region_id
  ON explorer_subregions(region_id);

CREATE INDEX idx_explorer_subregions_status
  ON explorer_subregions(status);

CREATE TABLE explorer_nodes (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  parent_kind    TEXT NOT NULL
                 CHECK (parent_kind IN ('region', 'subregion')),
  parent_id      UUID NOT NULL,
  node_type      TEXT NOT NULL
                 CHECK (node_type IN ('exploration', 'research')),

  title          TEXT NOT NULL,
  order_index    INTEGER NOT NULL DEFAULT 0,
  status         TEXT NOT NULL DEFAULT 'active'
                 CHECK (status IN ('active', 'inactive')),
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),

  -- exploration node 전용
  source_type    TEXT
                 CHECK (source_type IN ('youtube', 'web', 'creator', 'internal')),
  source_url     TEXT,
  content_id     UUID REFERENCES contents(id) ON DELETE SET NULL,
  summary        TEXT,

  -- research node 전용
  research_type  TEXT
                 CHECK (research_type IN ('concept', 'practice', 'problem', 'free')),
  layout_type    TEXT
                 CHECK (layout_type IN ('basic', 'two-column', 'note-card')),
  block_count    INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_explorer_nodes_parent
  ON explorer_nodes(parent_kind, parent_id);

CREATE INDEX idx_explorer_nodes_status
  ON explorer_nodes(status);

CREATE INDEX idx_explorer_nodes_node_type
  ON explorer_nodes(node_type);

CREATE INDEX idx_explorer_nodes_content_id
  ON explorer_nodes(content_id);

COMMIT;
