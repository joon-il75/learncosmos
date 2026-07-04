ALTER TABLE course_point_material_reports
  DROP CONSTRAINT IF EXISTS course_point_material_reports_status_check;

ALTER TABLE course_point_material_reports
  ADD CONSTRAINT course_point_material_reports_status_check
  CHECK (status IN ('open', 'reviewing', 'resolved', 'dismissed', 'cancelled'));
