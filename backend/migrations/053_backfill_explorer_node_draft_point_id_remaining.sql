BEGIN;

UPDATE explorer_nodes n
SET draft_point_id = p.id
FROM course_draft_points p
WHERE n.draft_point_id IS NULL
  AND n.parent_kind = 'region'
  AND n.parent_id = p.course_draft_lesson_id
  AND n.node_type = 'research'
  AND p.point_type = 'research'
  AND lower(trim(n.title)) = lower(trim(p.title));

UPDATE explorer_nodes n
SET draft_point_id = p.id
FROM course_draft_points p
WHERE n.draft_point_id IS NULL
  AND n.parent_kind = 'subregion'
  AND n.parent_id = p.course_draft_lesson_id
  AND n.node_type = 'exploration'
  AND p.point_type = 'exploration'
  AND lower(trim(n.title)) = lower(trim(p.title));

UPDATE explorer_nodes n
SET draft_point_id = p.id
FROM course_draft_points p
WHERE n.draft_point_id IS NULL
  AND n.parent_kind = 'region'
  AND n.parent_id = p.course_draft_lesson_id
  AND n.node_type = 'exploration'
  AND p.point_type = 'exploration'
  AND lower(trim(n.title)) = lower(trim(p.title));

UPDATE explorer_nodes n
SET draft_point_id = p.id
FROM course_draft_points p
WHERE n.draft_point_id IS NULL
  AND n.parent_kind = 'subregion'
  AND n.parent_id = p.course_draft_lesson_id
  AND n.node_type = 'research'
  AND p.point_type = 'research'
  AND lower(trim(n.title)) = lower(trim(p.title));

COMMIT;
