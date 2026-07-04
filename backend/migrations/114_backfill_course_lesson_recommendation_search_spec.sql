UPDATE course_lessons cl
SET recommendation_search_spec = cdl.recommendation_search_spec,
    updated_at = NOW()
FROM courses c, course_draft_lessons cdl
WHERE cl.course_id = c.id
  AND cdl.course_draft_id = c.source_draft_id
  AND cdl.order_index = cl.order_index
  AND cdl.title = cl.title
  AND cl.recommendation_search_spec = '{}'::jsonb
  AND cdl.recommendation_search_spec <> '{}'::jsonb;
