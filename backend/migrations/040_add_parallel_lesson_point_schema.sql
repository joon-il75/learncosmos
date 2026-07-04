ALTER TABLE course_draft_lessons
    ADD COLUMN IF NOT EXISTS course_draft_id UUID REFERENCES course_drafts(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS parent_lesson_id UUID REFERENCES course_draft_lessons(id) ON DELETE CASCADE;

UPDATE course_draft_lessons cdl
SET course_draft_id = cdlv.course_draft_id
FROM course_draft_levels cdlv
WHERE cdl.course_draft_level_id = cdlv.id
  AND cdl.course_draft_id IS NULL;

CREATE INDEX IF NOT EXISTS idx_course_draft_lessons_draft_id
    ON course_draft_lessons(course_draft_id);
CREATE INDEX IF NOT EXISTS idx_course_draft_lessons_parent_lesson_id
    ON course_draft_lessons(parent_lesson_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_course_draft_lessons_id_draft_id
    ON course_draft_lessons(id, course_draft_id);

ALTER TABLE course_lessons
    ADD COLUMN IF NOT EXISTS course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS parent_lesson_id UUID REFERENCES course_lessons(id) ON DELETE CASCADE;

UPDATE course_lessons cl
SET course_id = clv.course_id
FROM course_levels clv
WHERE cl.course_level_id = clv.id
  AND cl.course_id IS NULL;

CREATE INDEX IF NOT EXISTS idx_course_lessons_course_id
    ON course_lessons(course_id);
CREATE INDEX IF NOT EXISTS idx_course_lessons_parent_lesson_id
    ON course_lessons(parent_lesson_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_course_lessons_id_course_id
    ON course_lessons(id, course_id);

CREATE TABLE IF NOT EXISTS course_draft_points (
    id                     UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_draft_id        UUID        NOT NULL REFERENCES course_drafts(id) ON DELETE CASCADE,
    course_draft_lesson_id UUID        NOT NULL REFERENCES course_draft_lessons(id) ON DELETE CASCADE,
    point_type             TEXT        NOT NULL CHECK (point_type IN ('exploration', 'research')),
    status                 TEXT        NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'ready', 'learning', 'completed')),
    title                  TEXT        NOT NULL,
    description            TEXT,
    template_type          TEXT,
    selection_state        TEXT        CHECK (selection_state IN ('candidate', 'selected', 'rejected')),
    content_id             UUID        REFERENCES contents(id) ON DELETE SET NULL,
    external_url           TEXT,
    thumbnail_url          TEXT,
    price_type             TEXT        CHECK (price_type IN ('free', 'paid', 'owned')),
    rank_score             DOUBLE PRECISION,
    order_index            INTEGER     NOT NULL DEFAULT 0,
    completed_at           TIMESTAMPTZ,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (course_draft_lesson_id, order_index)
);

CREATE INDEX IF NOT EXISTS idx_course_draft_points_draft_id
    ON course_draft_points(course_draft_id);
CREATE INDEX IF NOT EXISTS idx_course_draft_points_lesson_id
    ON course_draft_points(course_draft_lesson_id);
CREATE INDEX IF NOT EXISTS idx_course_draft_points_type
    ON course_draft_points(point_type);
CREATE INDEX IF NOT EXISTS idx_course_draft_points_status
    ON course_draft_points(status);
CREATE UNIQUE INDEX IF NOT EXISTS idx_course_draft_points_id_lesson_id
    ON course_draft_points(id, course_draft_lesson_id);

CREATE TABLE IF NOT EXISTS course_points (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id        UUID        NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    course_lesson_id UUID        NOT NULL REFERENCES course_lessons(id) ON DELETE CASCADE,
    point_type       TEXT        NOT NULL CHECK (point_type IN ('exploration', 'research')),
    status           TEXT        NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'ready', 'learning', 'completed')),
    title            TEXT        NOT NULL,
    description      TEXT,
    template_type    TEXT,
    content_id       UUID        REFERENCES contents(id) ON DELETE SET NULL,
    external_url     TEXT,
    thumbnail_url    TEXT,
    price_type       TEXT        CHECK (price_type IN ('free', 'paid', 'owned')),
    rank_score       DOUBLE PRECISION,
    order_index      INTEGER     NOT NULL DEFAULT 0,
    completed_at     TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (course_lesson_id, order_index)
);

CREATE INDEX IF NOT EXISTS idx_course_points_course_id
    ON course_points(course_id);
CREATE INDEX IF NOT EXISTS idx_course_points_lesson_id
    ON course_points(course_lesson_id);
CREATE INDEX IF NOT EXISTS idx_course_points_type
    ON course_points(point_type);
CREATE INDEX IF NOT EXISTS idx_course_points_status
    ON course_points(status);
CREATE UNIQUE INDEX IF NOT EXISTS idx_course_points_id_lesson_id
    ON course_points(id, course_lesson_id);

CREATE TABLE IF NOT EXISTS course_draft_point_blocks (
    id                     UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_draft_point_id  UUID        NOT NULL REFERENCES course_draft_points(id) ON DELETE CASCADE,
    block_type             TEXT        NOT NULL CHECK (block_type IN ('text', 'image', 'link')),
    content                JSONB       NOT NULL DEFAULT '{}'::jsonb,
    order_index            INTEGER     NOT NULL DEFAULT 0,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (course_draft_point_id, order_index)
);

CREATE INDEX IF NOT EXISTS idx_course_draft_point_blocks_point_id
    ON course_draft_point_blocks(course_draft_point_id);

CREATE TABLE IF NOT EXISTS course_point_blocks (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_point_id  UUID        NOT NULL REFERENCES course_points(id) ON DELETE CASCADE,
    user_id          UUID        REFERENCES users(id) ON DELETE SET NULL,
    block_type       TEXT        NOT NULL CHECK (block_type IN ('text', 'image', 'link')),
    content          JSONB       NOT NULL DEFAULT '{}'::jsonb,
    order_index      INTEGER     NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (course_point_id, order_index)
);

CREATE INDEX IF NOT EXISTS idx_course_point_blocks_point_id
    ON course_point_blocks(course_point_id);

CREATE TABLE IF NOT EXISTS course_draft_point_journal_entries (
    course_draft_point_id UUID        PRIMARY KEY REFERENCES course_draft_points(id) ON DELETE CASCADE,
    observation           TEXT        NOT NULL DEFAULT '',
    reflection            TEXT        NOT NULL DEFAULT '',
    next_step             TEXT        NOT NULL DEFAULT '',
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS course_draft_point_record_entries (
    course_draft_point_id UUID        PRIMARY KEY REFERENCES course_draft_points(id) ON DELETE CASCADE,
    study_minutes         INTEGER     NOT NULL DEFAULT 0,
    practice_count        INTEGER     NOT NULL DEFAULT 0,
    confidence_level      INTEGER     NOT NULL DEFAULT 3,
    application_note      TEXT        NOT NULL DEFAULT '',
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS course_draft_point_artifact_entries (
    course_draft_point_id UUID        PRIMARY KEY REFERENCES course_draft_points(id) ON DELETE CASCADE,
    artifact_type         TEXT        NOT NULL DEFAULT '',
    title                 TEXT        NOT NULL DEFAULT '',
    url                   TEXT        NOT NULL DEFAULT '',
    description           TEXT        NOT NULL DEFAULT '',
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS course_draft_community_posts (
    id                     UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_draft_id        UUID        NOT NULL REFERENCES course_drafts(id) ON DELETE CASCADE,
    course_draft_lesson_id UUID        NOT NULL REFERENCES course_draft_lessons(id) ON DELETE CASCADE,
    course_draft_point_id  UUID,
    user_id                UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_type              TEXT        NOT NULL DEFAULT 'discussion',
    title                  TEXT        NOT NULL DEFAULT '',
    body                   TEXT        NOT NULL DEFAULT '',
    visibility             TEXT        NOT NULL DEFAULT 'private',
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_course_draft_community_posts_point_lesson
        FOREIGN KEY (course_draft_point_id, course_draft_lesson_id)
        REFERENCES course_draft_points(id, course_draft_lesson_id)
        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_course_draft_community_posts_draft_id
    ON course_draft_community_posts(course_draft_id);
CREATE INDEX IF NOT EXISTS idx_course_draft_community_posts_lesson_id
    ON course_draft_community_posts(course_draft_lesson_id);
CREATE INDEX IF NOT EXISTS idx_course_draft_community_posts_point_id
    ON course_draft_community_posts(course_draft_point_id);

CREATE TABLE IF NOT EXISTS course_community_posts (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id        UUID        NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    course_lesson_id UUID        NOT NULL REFERENCES course_lessons(id) ON DELETE CASCADE,
    course_point_id  UUID,
    user_id          UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_type        TEXT        NOT NULL DEFAULT 'discussion',
    title            TEXT        NOT NULL DEFAULT '',
    body             TEXT        NOT NULL DEFAULT '',
    visibility       TEXT        NOT NULL DEFAULT 'private',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_course_community_posts_point_lesson
        FOREIGN KEY (course_point_id, course_lesson_id)
        REFERENCES course_points(id, course_lesson_id)
        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_course_community_posts_course_id
    ON course_community_posts(course_id);
CREATE INDEX IF NOT EXISTS idx_course_community_posts_lesson_id
    ON course_community_posts(course_lesson_id);
CREATE INDEX IF NOT EXISTS idx_course_community_posts_point_id
    ON course_community_posts(course_point_id);
