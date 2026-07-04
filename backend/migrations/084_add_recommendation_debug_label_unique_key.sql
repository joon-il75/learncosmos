CREATE UNIQUE INDEX IF NOT EXISTS idx_recommendation_debug_labels_run_candidate
  ON recommendation_debug_labels(run_id, candidate_key)
  WHERE run_id IS NOT NULL;
