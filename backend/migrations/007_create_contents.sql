-- 006_create_contents.sql
-- LearnWeaver 콘텐츠 원본 저장소 테이블

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS contents (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content_type     TEXT        NOT NULL CHECK (content_type IN ('youtube', 'blog', 'article', 'internal')),
    url              TEXT,
    title            TEXT        NOT NULL,
    description      TEXT,
    thumbnail_url    TEXT,
    duration_seconds INTEGER,
    author           TEXT,
    language         TEXT        NOT NULL DEFAULT 'ko',
    is_public        BOOLEAN     NOT NULL DEFAULT false,
    quality_score    FLOAT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contents_user_id    ON contents(user_id);
CREATE INDEX IF NOT EXISTS idx_contents_type       ON contents(content_type);
CREATE INDEX IF NOT EXISTS idx_contents_created_at ON contents(created_at DESC);
