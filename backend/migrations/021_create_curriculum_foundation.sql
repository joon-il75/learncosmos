-- 019_create_curriculum_foundation.sql
-- LearnWeaver AI 커리큘럼 1단계: 초안/확정/리소스 연결/이벤트 로그 기반 구축

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS course_drafts (
    id                    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id               UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_query          TEXT        NOT NULL,
    learning_goal         TEXT,
    current_level         TEXT,
    duration_weeks        INTEGER,
    study_hours_per_week  INTEGER,
    preferred_format      TEXT,
    title                 TEXT        NOT NULL,
    description           TEXT,
    status                TEXT        NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'confirmed', 'archived')),
    confirmed_course_id   UUID,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_course_drafts_user_id ON course_drafts(user_id);
CREATE INDEX IF NOT EXISTS idx_course_drafts_status ON course_drafts(status);
CREATE INDEX IF NOT EXISTS idx_course_drafts_created_at ON course_drafts(created_at DESC);

CREATE TABLE IF NOT EXISTS course_draft_levels (
    id                    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_draft_id       UUID        NOT NULL REFERENCES course_drafts(id) ON DELETE CASCADE,
    title                 TEXT        NOT NULL,
    description           TEXT,
    objective             TEXT,
    level_code            TEXT        NOT NULL DEFAULT 'custom',
    order_index           INTEGER     NOT NULL,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (course_draft_id, order_index)
);

CREATE INDEX IF NOT EXISTS idx_course_draft_levels_draft_id ON course_draft_levels(course_draft_id);

CREATE TABLE IF NOT EXISTS course_draft_lessons (
    id                    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_draft_level_id UUID        NOT NULL REFERENCES course_draft_levels(id) ON DELETE CASCADE,
    title                 TEXT        NOT NULL,
    objective             TEXT,
    summary               TEXT,
    difficulty_level      TEXT,
    lesson_role           TEXT        NOT NULL DEFAULT 'core' CHECK (lesson_role IN ('core', 'support', 'demo', 'reference', 'inspiration')),
    source_type           TEXT        NOT NULL DEFAULT 'manual' CHECK (source_type IN ('ai_generated', 'manual', 'recommended')),
    order_index           INTEGER     NOT NULL,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (course_draft_level_id, order_index)
);

CREATE INDEX IF NOT EXISTS idx_course_draft_lessons_level_id ON course_draft_lessons(course_draft_level_id);

CREATE TABLE IF NOT EXISTS course_draft_lesson_resources (
    id                    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_draft_lesson_id UUID       NOT NULL REFERENCES course_draft_lessons(id) ON DELETE CASCADE,
    resource_type         TEXT        NOT NULL CHECK (resource_type IN ('content', 'external', 'creator', 'ai_generated')),
    selection_state       TEXT        NOT NULL DEFAULT 'selected' CHECK (selection_state IN ('candidate', 'selected', 'rejected')),
    content_id            UUID        REFERENCES contents(id) ON DELETE SET NULL,
    external_url          TEXT,
    title                 TEXT        NOT NULL,
    description           TEXT,
    thumbnail_url         TEXT,
    price_type            TEXT        NOT NULL DEFAULT 'free' CHECK (price_type IN ('free', 'paid', 'owned')),
    rank_score            DOUBLE PRECISION,
    order_index           INTEGER     NOT NULL DEFAULT 0,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_course_draft_lesson_resources_lesson_id ON course_draft_lesson_resources(course_draft_lesson_id);
CREATE INDEX IF NOT EXISTS idx_course_draft_lesson_resources_content_id ON course_draft_lesson_resources(content_id);
CREATE INDEX IF NOT EXISTS idx_course_draft_lesson_resources_selection_state ON course_draft_lesson_resources(selection_state);

CREATE TABLE IF NOT EXISTS courses (
    id                    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id               UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_draft_id       UUID        REFERENCES course_drafts(id) ON DELETE SET NULL,
    source_query          TEXT        NOT NULL,
    learning_goal         TEXT,
    current_level         TEXT,
    duration_weeks        INTEGER,
    study_hours_per_week  INTEGER,
    preferred_format      TEXT,
    title                 TEXT        NOT NULL,
    description           TEXT,
    status                TEXT        NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'archived')),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_courses_user_id ON courses(user_id);
CREATE INDEX IF NOT EXISTS idx_courses_status ON courses(status);
CREATE INDEX IF NOT EXISTS idx_courses_created_at ON courses(created_at DESC);

CREATE TABLE IF NOT EXISTS course_levels (
    id                    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id             UUID        NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title                 TEXT        NOT NULL,
    description           TEXT,
    objective             TEXT,
    level_code            TEXT        NOT NULL DEFAULT 'custom',
    order_index           INTEGER     NOT NULL,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (course_id, order_index)
);

CREATE INDEX IF NOT EXISTS idx_course_levels_course_id ON course_levels(course_id);

CREATE TABLE IF NOT EXISTS course_lessons (
    id                    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_level_id       UUID        NOT NULL REFERENCES course_levels(id) ON DELETE CASCADE,
    title                 TEXT        NOT NULL,
    objective             TEXT,
    summary               TEXT,
    difficulty_level      TEXT,
    lesson_role           TEXT        NOT NULL DEFAULT 'core' CHECK (lesson_role IN ('core', 'support', 'demo', 'reference', 'inspiration')),
    source_type           TEXT        NOT NULL DEFAULT 'manual' CHECK (source_type IN ('ai_generated', 'manual', 'recommended')),
    order_index           INTEGER     NOT NULL,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (course_level_id, order_index)
);

CREATE INDEX IF NOT EXISTS idx_course_lessons_level_id ON course_lessons(course_level_id);

CREATE TABLE IF NOT EXISTS course_lesson_resources (
    id                    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_lesson_id      UUID        NOT NULL REFERENCES course_lessons(id) ON DELETE CASCADE,
    resource_type         TEXT        NOT NULL CHECK (resource_type IN ('content', 'external', 'creator', 'ai_generated')),
    content_id            UUID        REFERENCES contents(id) ON DELETE SET NULL,
    external_url          TEXT,
    title                 TEXT        NOT NULL,
    description           TEXT,
    thumbnail_url         TEXT,
    price_type            TEXT        NOT NULL DEFAULT 'free' CHECK (price_type IN ('free', 'paid', 'owned')),
    rank_score            DOUBLE PRECISION,
    order_index           INTEGER     NOT NULL DEFAULT 0,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_course_lesson_resources_lesson_id ON course_lesson_resources(course_lesson_id);
CREATE INDEX IF NOT EXISTS idx_course_lesson_resources_content_id ON course_lesson_resources(content_id);

CREATE TABLE IF NOT EXISTS recommendation_events (
    id                    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id               UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_draft_id       UUID        REFERENCES course_drafts(id) ON DELETE SET NULL,
    course_draft_level_id UUID        REFERENCES course_draft_levels(id) ON DELETE SET NULL,
    course_draft_lesson_id UUID       REFERENCES course_draft_lessons(id) ON DELETE SET NULL,
    content_id            UUID        REFERENCES contents(id) ON DELETE SET NULL,
    event_type            TEXT        NOT NULL CHECK (event_type IN (
        'curriculum_generated',
        'lesson_recommended',
        'lesson_selected',
        'lesson_rejected',
        'lesson_added_manually',
        'lesson_generated_by_ai',
        'curriculum_confirmed',
        'curriculum_edited',
        'curriculum_edited_after_start'
    )),
    source_query          TEXT,
    payload               JSONB       NOT NULL DEFAULT '{}'::jsonb,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_recommendation_events_user_id ON recommendation_events(user_id);
CREATE INDEX IF NOT EXISTS idx_recommendation_events_draft_id ON recommendation_events(course_draft_id);
CREATE INDEX IF NOT EXISTS idx_recommendation_events_type ON recommendation_events(event_type);
CREATE INDEX IF NOT EXISTS idx_recommendation_events_created_at ON recommendation_events(created_at DESC);

ALTER TABLE course_drafts
    ADD CONSTRAINT fk_course_drafts_confirmed_course_id
    FOREIGN KEY (confirmed_course_id) REFERENCES courses(id) ON DELETE SET NULL;
