ALTER TABLE curriculum_patterns
  ADD COLUMN IF NOT EXISTS recommendation_search_spec_template JSONB NOT NULL DEFAULT '{}'::jsonb;

ALTER TABLE curriculum_patterns
  DROP CONSTRAINT IF EXISTS curriculum_patterns_recommendation_search_spec_template_json_check;

ALTER TABLE curriculum_patterns
  ADD CONSTRAINT curriculum_patterns_recommendation_search_spec_template_json_check
  CHECK (jsonb_typeof(recommendation_search_spec_template) = 'object');
