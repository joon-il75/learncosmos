-- Remove the course draft goal achievability check feature.
DELETE FROM point_settings
WHERE key = 'course_eval_cost';

ALTER TABLE course_drafts
    DROP COLUMN IF EXISTS evaluation_result,
    DROP COLUMN IF EXISTS evaluated_at;
