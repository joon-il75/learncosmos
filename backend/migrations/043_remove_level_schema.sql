-- Phase 5: level 스키마 완전 제거
-- 목표: course_draft_levels, course_levels 테이블 및 관련 *_level_id 컬럼 삭제
-- dev DB는 데이터가 없으므로 데이터 마이그레이션 불필요

BEGIN;

-- 1. course_draft_research_nodes: course_draft_level_id → course_draft_lesson_id (main lesson FK)
ALTER TABLE course_draft_research_nodes
    DROP CONSTRAINT IF EXISTS course_draft_research_nodes_course_draft_level_id_fkey;

DROP INDEX IF EXISTS idx_course_draft_research_nodes_level_id;

ALTER TABLE course_draft_research_nodes
    RENAME COLUMN course_draft_level_id TO course_draft_lesson_id;

ALTER TABLE course_draft_research_nodes
    ADD CONSTRAINT course_draft_research_nodes_course_draft_lesson_id_fkey
    FOREIGN KEY (course_draft_lesson_id) REFERENCES course_draft_lessons(id) ON DELETE CASCADE;

CREATE INDEX idx_course_draft_research_nodes_lesson_id
    ON course_draft_research_nodes (course_draft_lesson_id);

-- 2. course_draft_lessons: course_draft_level_id 컬럼 및 관련 제약 제거
ALTER TABLE course_draft_lessons
    DROP CONSTRAINT IF EXISTS course_draft_lessons_course_draft_level_id_fkey;

ALTER TABLE course_draft_lessons
    DROP CONSTRAINT IF EXISTS course_draft_lessons_course_draft_level_id_order_index_key;

DROP INDEX IF EXISTS idx_course_draft_lessons_level_id;

ALTER TABLE course_draft_lessons
    DROP COLUMN IF EXISTS course_draft_level_id;

-- course_draft_id: NULL → NOT NULL
ALTER TABLE course_draft_lessons
    ALTER COLUMN course_draft_id SET NOT NULL;

-- order_index uniqueness를 (course_draft_id, parent_lesson_id, order_index) 기준으로 재설정
-- (동일 부모 안에서 고유)
CREATE UNIQUE INDEX idx_course_draft_lessons_parent_order
    ON course_draft_lessons (course_draft_id, COALESCE(parent_lesson_id, '00000000-0000-0000-0000-000000000000'::uuid), order_index);

-- 3. course_draft_detail_notes: course_draft_level_id 컬럼 및 관련 제약 제거
ALTER TABLE course_draft_detail_notes
    DROP CONSTRAINT IF EXISTS course_draft_detail_notes_course_draft_level_id_fkey;

ALTER TABLE course_draft_detail_notes
    DROP CONSTRAINT IF EXISTS course_draft_detail_notes_check;

ALTER TABLE course_draft_detail_notes
    DROP CONSTRAINT IF EXISTS course_draft_detail_notes_target_type_check;

DROP INDEX IF EXISTS idx_course_draft_detail_notes_level;

ALTER TABLE course_draft_detail_notes
    DROP COLUMN IF EXISTS course_draft_level_id;

-- 단순화: course_draft_lesson_id NOT NULL, target_type = 'lesson' 고정
ALTER TABLE course_draft_detail_notes
    ALTER COLUMN course_draft_lesson_id SET NOT NULL;

ALTER TABLE course_draft_detail_notes
    ADD CONSTRAINT course_draft_detail_notes_target_type_check
    CHECK (target_type = 'lesson');

ALTER TABLE course_draft_detail_notes
    ALTER COLUMN target_type SET DEFAULT 'lesson';

-- 4. recommendation_events: course_draft_level_id 컬럼 제거
ALTER TABLE recommendation_events
    DROP CONSTRAINT IF EXISTS recommendation_events_course_draft_level_id_fkey;

ALTER TABLE recommendation_events
    DROP COLUMN IF EXISTS course_draft_level_id;

-- 5. course_lessons: course_level_id 컬럼 및 관련 제약 제거
ALTER TABLE course_lessons
    DROP CONSTRAINT IF EXISTS course_lessons_course_level_id_fkey;

ALTER TABLE course_lessons
    DROP CONSTRAINT IF EXISTS course_lessons_course_level_id_order_index_key;

DROP INDEX IF EXISTS idx_course_lessons_level_id;

ALTER TABLE course_lessons
    DROP COLUMN IF EXISTS course_level_id;

-- course_id: NULL → NOT NULL
ALTER TABLE course_lessons
    ALTER COLUMN course_id SET NOT NULL;

-- order_index uniqueness를 (course_id, parent_lesson_id, order_index) 기준으로 재설정
CREATE UNIQUE INDEX idx_course_lessons_parent_order
    ON course_lessons (course_id, COALESCE(parent_lesson_id, '00000000-0000-0000-0000-000000000000'::uuid), order_index);

-- 6. course_draft_levels 테이블 삭제
DROP TABLE IF EXISTS course_draft_levels;

-- 7. course_levels 테이블 삭제
DROP TABLE IF EXISTS course_levels;

COMMIT;
