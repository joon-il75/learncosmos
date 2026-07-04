-- Phase 6: course_draft_points + course_points 통합
-- 목표: legacy lesson_resources / research_nodes 테이블 제거
-- dev DB는 데이터가 없으므로 데이터 마이그레이션 불필요

BEGIN;

-- 1. Draft 사이드 legacy 제거
DROP TABLE IF EXISTS course_draft_research_node_blocks;
DROP TABLE IF EXISTS course_draft_research_nodes;
DROP TABLE IF EXISTS course_draft_lesson_resources;

-- 2. Course 사이드 legacy 제거
DROP TABLE IF EXISTS course_research_node_blocks;
DROP TABLE IF EXISTS course_research_nodes;
DROP TABLE IF EXISTS course_lesson_resources;

COMMIT;
