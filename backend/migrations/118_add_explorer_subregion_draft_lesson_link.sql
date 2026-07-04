ALTER TABLE explorer_subregions
  ADD COLUMN IF NOT EXISTS course_draft_lesson_id UUID REFERENCES course_draft_lessons(id) ON DELETE SET NULL;

UPDATE explorer_subregions es
SET course_draft_lesson_id = cdl.id
FROM explorer_regions er, course_draft_lessons cdl
WHERE es.region_id = er.id
  AND cdl.id = es.id
  AND cdl.course_draft_id = er.course_draft_id
  AND es.course_draft_lesson_id IS NULL;

CREATE INDEX IF NOT EXISTS idx_explorer_subregions_draft_lesson_id
  ON explorer_subregions(course_draft_lesson_id);
