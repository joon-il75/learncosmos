ALTER TABLE explorer_regions
  ADD COLUMN IF NOT EXISTS course_draft_lesson_id UUID REFERENCES course_draft_lessons(id) ON DELETE SET NULL;

UPDATE explorer_regions er
SET course_draft_lesson_id = cdl.id
FROM course_draft_lessons cdl
WHERE cdl.id = er.id
  AND cdl.course_draft_id = er.course_draft_id
  AND cdl.parent_lesson_id IS NULL
  AND er.course_draft_lesson_id IS NULL;

CREATE INDEX IF NOT EXISTS idx_explorer_regions_draft_lesson_id
  ON explorer_regions(course_draft_lesson_id);
