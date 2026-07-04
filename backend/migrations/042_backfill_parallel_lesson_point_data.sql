-- Phase 2 draft backfill for the parallel lesson/point schema.
--
-- Important:
-- - This migration is intentionally written as an idempotent backfill script.
-- - It should be applied only when the runtime is ready to tolerate main lesson
--   rows being added to course_draft_lessons / course_lessons.
-- - During the current legacy level-based read path, use dry-run validation only.

CREATE OR REPLACE FUNCTION lw_deterministic_uuid(input TEXT)
RETURNS UUID
LANGUAGE SQL
IMMUTABLE
STRICT
AS $$
    SELECT (
        SUBSTR(MD5(input), 1, 8) || '-' ||
        SUBSTR(MD5(input), 9, 4) || '-' ||
        SUBSTR(MD5(input), 13, 4) || '-' ||
        SUBSTR(MD5(input), 17, 4) || '-' ||
        SUBSTR(MD5(input), 21, 12)
    )::uuid
$$;

-- 1. Draft: legacy level -> main lesson
INSERT INTO course_draft_lessons (
    id,
    course_draft_id,
    course_draft_level_id,
    parent_lesson_id,
    title,
    objective,
    summary,
    difficulty_level,
    lesson_role,
    source_type,
    order_index,
    created_at,
    updated_at
)
SELECT
    lw_deterministic_uuid('draft-main-lesson:' || cdlv.id::text),
    cdlv.course_draft_id,
    cdlv.id,
    NULL,
    cdlv.title,
    cdlv.objective,
    cdlv.description,
    NULL,
    'core',
    'manual',
    -1,
    cdlv.created_at,
    cdlv.updated_at
FROM course_draft_levels cdlv
ON CONFLICT (id) DO NOTHING;

-- 2. Draft: legacy lesson -> sub lesson
UPDATE course_draft_lessons cdl
SET parent_lesson_id = lw_deterministic_uuid('draft-main-lesson:' || cdl.course_draft_level_id::text)
WHERE cdl.id <> lw_deterministic_uuid('draft-main-lesson:' || cdl.course_draft_level_id::text)
  AND cdl.parent_lesson_id IS DISTINCT FROM lw_deterministic_uuid('draft-main-lesson:' || cdl.course_draft_level_id::text);

-- 3. Draft: lesson resources -> exploration points
INSERT INTO course_draft_points (
    id,
    course_draft_id,
    course_draft_lesson_id,
    point_type,
    status,
    title,
    description,
    template_type,
    selection_state,
    content_id,
    external_url,
    thumbnail_url,
    price_type,
    rank_score,
    order_index,
    completed_at,
    created_at,
    updated_at
)
SELECT
    lw_deterministic_uuid('draft-exploration-point:' || cdlr.id::text),
    cdl.course_draft_id,
    cdlr.course_draft_lesson_id,
    'exploration',
    'draft',
    cdlr.title,
    cdlr.description,
    NULL,
    cdlr.selection_state,
    cdlr.content_id,
    cdlr.external_url,
    cdlr.thumbnail_url,
    cdlr.price_type,
    cdlr.rank_score,
    cdlr.order_index,
    NULL,
    cdlr.created_at,
    cdlr.updated_at
FROM course_draft_lesson_resources cdlr
JOIN course_draft_lessons cdl
  ON cdl.id = cdlr.course_draft_lesson_id
ON CONFLICT (id) DO NOTHING;

-- 4. Draft: level research nodes -> main-lesson research points
INSERT INTO course_draft_points (
    id,
    course_draft_id,
    course_draft_lesson_id,
    point_type,
    status,
    title,
    description,
    template_type,
    selection_state,
    content_id,
    external_url,
    thumbnail_url,
    price_type,
    rank_score,
    order_index,
    completed_at,
    created_at,
    updated_at
)
SELECT
    lw_deterministic_uuid('draft-research-point:' || cdrn.id::text),
    cdrn.course_draft_id,
    lw_deterministic_uuid('draft-main-lesson:' || cdrn.course_draft_level_id::text),
    'research',
    'draft',
    cdrn.title,
    NULL,
    cdrn.template_type,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL,
    cdrn.order_index,
    NULL,
    cdrn.created_at,
    cdrn.updated_at
FROM course_draft_research_nodes cdrn
ON CONFLICT (id) DO NOTHING;

-- 5. Draft: research node blocks -> research point blocks
INSERT INTO course_draft_point_blocks (
    id,
    course_draft_point_id,
    block_type,
    content,
    order_index,
    created_at,
    updated_at
)
SELECT
    lw_deterministic_uuid('draft-point-block:' || cdrnb.id::text),
    lw_deterministic_uuid('draft-research-point:' || cdrnb.research_node_id::text),
    cdrnb.block_type,
    cdrnb.content,
    cdrnb.order_index,
    cdrnb.created_at,
    cdrnb.updated_at
FROM course_draft_research_node_blocks cdrnb
ON CONFLICT (id) DO NOTHING;

-- 6. Draft: lesson journal/record/artifact -> point tables
-- Only lessons with exactly one exploration resource are auto-mapped.
WITH single_resource_lessons AS (
    SELECT
        course_draft_lesson_id,
        MIN(id::text)::uuid AS resource_id
    FROM course_draft_lesson_resources
    GROUP BY course_draft_lesson_id
    HAVING COUNT(*) = 1
)
INSERT INTO course_draft_point_journal_entries (
    course_draft_point_id,
    observation,
    reflection,
    next_step,
    created_at,
    updated_at
)
SELECT
    lw_deterministic_uuid('draft-exploration-point:' || srl.resource_id::text),
    cdje.observation,
    cdje.reflection,
    cdje.next_step,
    cdje.created_at,
    cdje.updated_at
FROM course_draft_journal_entries cdje
JOIN single_resource_lessons srl
  ON srl.course_draft_lesson_id = cdje.course_draft_lesson_id
ON CONFLICT (course_draft_point_id) DO NOTHING;

WITH single_resource_lessons AS (
    SELECT
        course_draft_lesson_id,
        MIN(id::text)::uuid AS resource_id
    FROM course_draft_lesson_resources
    GROUP BY course_draft_lesson_id
    HAVING COUNT(*) = 1
)
INSERT INTO course_draft_point_record_entries (
    course_draft_point_id,
    study_minutes,
    practice_count,
    confidence_level,
    application_note,
    created_at,
    updated_at
)
SELECT
    lw_deterministic_uuid('draft-exploration-point:' || srl.resource_id::text),
    cdre.study_minutes,
    cdre.practice_count,
    cdre.confidence_level,
    cdre.application_note,
    cdre.created_at,
    cdre.updated_at
FROM course_draft_record_entries cdre
JOIN single_resource_lessons srl
  ON srl.course_draft_lesson_id = cdre.course_draft_lesson_id
ON CONFLICT (course_draft_point_id) DO NOTHING;

WITH single_resource_lessons AS (
    SELECT
        course_draft_lesson_id,
        MIN(id::text)::uuid AS resource_id
    FROM course_draft_lesson_resources
    GROUP BY course_draft_lesson_id
    HAVING COUNT(*) = 1
)
INSERT INTO course_draft_point_artifact_entries (
    course_draft_point_id,
    artifact_type,
    title,
    url,
    description,
    created_at,
    updated_at
)
SELECT
    lw_deterministic_uuid('draft-exploration-point:' || srl.resource_id::text),
    cdae.artifact_type,
    cdae.title,
    cdae.url,
    cdae.description,
    cdae.created_at,
    cdae.updated_at
FROM course_draft_artifact_entries cdae
JOIN single_resource_lessons srl
  ON srl.course_draft_lesson_id = cdae.course_draft_lesson_id
ON CONFLICT (course_draft_point_id) DO NOTHING;

-- 7. Confirmed: legacy level -> main lesson
INSERT INTO course_lessons (
    id,
    course_id,
    course_level_id,
    parent_lesson_id,
    title,
    objective,
    summary,
    difficulty_level,
    lesson_role,
    source_type,
    order_index,
    created_at,
    updated_at
)
SELECT
    lw_deterministic_uuid('course-main-lesson:' || clv.id::text),
    clv.course_id,
    clv.id,
    NULL,
    clv.title,
    clv.objective,
    clv.description,
    NULL,
    'core',
    'manual',
    -1,
    clv.created_at,
    clv.updated_at
FROM course_levels clv
ON CONFLICT (id) DO NOTHING;

-- 8. Confirmed: legacy lesson -> sub lesson
UPDATE course_lessons cl
SET parent_lesson_id = lw_deterministic_uuid('course-main-lesson:' || cl.course_level_id::text)
WHERE cl.id <> lw_deterministic_uuid('course-main-lesson:' || cl.course_level_id::text)
  AND cl.parent_lesson_id IS DISTINCT FROM lw_deterministic_uuid('course-main-lesson:' || cl.course_level_id::text);

-- 9. Confirmed: lesson resources -> exploration points
INSERT INTO course_points (
    id,
    course_id,
    course_lesson_id,
    point_type,
    status,
    title,
    description,
    template_type,
    content_id,
    external_url,
    thumbnail_url,
    price_type,
    rank_score,
    order_index,
    completed_at,
    created_at,
    updated_at
)
SELECT
    lw_deterministic_uuid('course-exploration-point:' || clr.id::text),
    cl.course_id,
    clr.course_lesson_id,
    'exploration',
    CASE
        WHEN cd.status = 'learning' THEN 'learning'
        WHEN cd.status = 'archived' THEN 'completed'
        ELSE 'ready'
    END,
    clr.title,
    clr.description,
    NULL,
    clr.content_id,
    clr.external_url,
    clr.thumbnail_url,
    clr.price_type,
    clr.rank_score,
    clr.order_index,
    NULL,
    clr.created_at,
    clr.updated_at
FROM course_lesson_resources clr
JOIN course_lessons cl
  ON cl.id = clr.course_lesson_id
JOIN courses c
  ON c.id = cl.course_id
LEFT JOIN course_drafts cd
  ON cd.id = c.source_draft_id
ON CONFLICT (id) DO NOTHING;

-- 10. Confirmed: research nodes -> research points
INSERT INTO course_points (
    id,
    course_id,
    course_lesson_id,
    point_type,
    status,
    title,
    description,
    template_type,
    content_id,
    external_url,
    thumbnail_url,
    price_type,
    rank_score,
    order_index,
    completed_at,
    created_at,
    updated_at
)
SELECT
    lw_deterministic_uuid('course-research-point:' || crn.id::text),
    crn.course_id,
    crn.course_lesson_id,
    'research',
    CASE
        WHEN cd.status = 'learning' THEN 'learning'
        WHEN cd.status = 'archived' THEN 'completed'
        ELSE 'ready'
    END,
    crn.title,
    NULL,
    crn.template_type,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL,
    crn.order_index,
    NULL,
    crn.created_at,
    crn.updated_at
FROM course_research_nodes crn
JOIN courses c
  ON c.id = crn.course_id
LEFT JOIN course_drafts cd
  ON cd.id = c.source_draft_id
ON CONFLICT (id) DO NOTHING;

-- 11. Confirmed: research node blocks -> point blocks
INSERT INTO course_point_blocks (
    id,
    course_point_id,
    user_id,
    block_type,
    content,
    order_index,
    created_at,
    updated_at
)
SELECT
    lw_deterministic_uuid('course-point-block:' || crnb.id::text),
    lw_deterministic_uuid('course-research-point:' || crnb.research_node_id::text),
    crnb.user_id,
    crnb.block_type,
    crnb.content,
    crnb.order_index,
    crnb.created_at,
    crnb.updated_at
FROM course_research_node_blocks crnb
ON CONFLICT (id) DO NOTHING;

DROP FUNCTION IF EXISTS lw_deterministic_uuid(TEXT);
