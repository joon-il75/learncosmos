export type Candidate = {
  content_id?: string | null
  title: string
  description?: string | null
  thumbnail_url?: string | null
  external_url?: string | null
  rank_score: number
  selection_reason?: string
  selection_signals?: string[]
  lexical_matched?: boolean
  vector_matched?: boolean
  content_type?: string
  resource_type?: string
  quality_score?: number
  language?: string
  feature_snapshot?: Record<string, unknown>
  ranker_score?: number
  ranker_rank?: number
  ranker_route_rank?: number
  ranker_rerank_rank?: number
  ranker_rank_delta?: number
  ranker_reason?: string
  ranker_provider?: string
  ranker_model_version?: string
  source_type?: string
  external_source?: string
  external_content_id?: string
  author?: string
  already_saved?: boolean
  eligible_for_save?: boolean
}

export type EmbeddingStage = {
  provider: string
  used: boolean
  fallback_text_only: boolean
  dimension: number
  error?: string
  preview?: number[]
}

export type ExternalPreviewItem = {
  source?: string
  title: string
  description?: string
  url?: string
  external_url?: string
  thumbnail_url?: string
  author?: string
  language?: string
  external_content_id?: string
  content_id?: string
  already_saved: boolean
  duplicate_with_internal: boolean
  duplicate_with_external: boolean
  eligible_for_supplement: boolean
  reason?: string
}

export type PreviewResponse = {
  stages: {
    query: {
      detected_intent: string
      raw_query: string
      dense_query: string
      lexical_query: string
      lexical_query_ko: string
      lexical_query_expanded: string
      tokens: string[]
      expanded_terms: string[]
      goal_pattern_key: string
      goal_subpattern_key: string
      goal_query_hint: string
      goal_query_terms: string[]
      dictionary_version: string
      recommendation_query?: string
    }
    embedding: EmbeddingStage
    search: {
      max_per_lesson: number
      retrieval_limit: number
      lexical_hit_count: number
      vector_hit_count: number
      merged_hit_count: number
      failure_reason?: string
    }
    external: {
      provider: string
      attempted: boolean
      reason?: string
      external_query: string
      raw_hit_count: number
      existing_content_match_count: number
      dedupe_dropped_count: number
      supplement_hit_count: number
      final_merged_with_external_count: number
      preview_results: ExternalPreviewItem[]
    }
  }
  candidates: Candidate[]
}

export type RolloutState = {
  id: string
  mode: string
  traffic_percent: number
  embedding_provider: string
  embedding_model: string
  embedding_dimension: number
  ranker_model_version: string
  feature_schema_version: string
  quality_gate_status: string
  rollback_reason: string
  active: boolean
  updated_at: string
}

export type RankerArtifactStatus = {
  ready: boolean
  reason: string
  provider: string
  ranker_model_version: string
  feature_schema_version: string
  expected_model_version: string
  expected_schema_version: string
  config_source: string
  model_path: string
  manifest_path: string
  model_file_exists: boolean
  manifest_exists: boolean
  model_size_bytes: number
  model_updated_at: string
  manifest_version_match: boolean
  manifest_schema_match: boolean
  checked_paths: string[]
}

export type DebugScenario = {
  id: string
  created_by: string
  course_title: string
  initial_user_intent: string
  status: string
  notes: string
  run_count?: number
  label_count?: number
  created_at: string
  updated_at: string
}

export type DebugRun = {
  id: string
  scenario_id: string
  run_type: string
  request_snapshot: string
  baseline_result_snapshot: string
  shadow_result_snapshot: string
  feature_snapshot: string
  metrics_snapshot: string
  provider_snapshot: string
  latency_ms?: number | null
  created_at: string
}

export type DebugScenarioDetail = DebugScenario & {
  goal_profile_snapshot: string
  generated_lessons_snapshot: string
  runs: DebugRun[]
  labels: DebugLabel[]
}

export type GoalMessage = {
  role: 'user' | 'lumi'
  content: string
}

export type GoalProfile = {
  user_intent: string
  motivation?: string | null
  usage_context?: string | null
  confirmed_goal?: string | null
  goal_type?: string | null
  output_type?: string | null
  difficulty_level?: string | null
  time_horizon?: string | null
  interview_state: string
  messages: GoalMessage[]
  version: number
}

export type GeneratedLessonItem = {
  lesson_id: string
  title: string
  objective: string
  order_index: number
  source_type: string
  lesson_role: string
}

export type GeneratedLessonsSnapshot = {
  draft_title: string
  draft_description: string
  confirmed_goal: string
  completion_criteria: string[]
  lessons: GeneratedLessonItem[]
  lesson_count: number
  provider: string
  model: string
  latency_ms: number
  generated_at: string
}

export type RecommendationCompareLesson = {
  run_id?: string
  lesson: {
    lesson_id: string
    title: string
    objective: string
    order_index: number
    recommendation_query?: string
  }
  baseline: PreviewResponse
  shadow: {
    enabled: boolean
    reason: string
    embedding?: {
      provider: string
      model: string
      attempted: boolean
      used: boolean
      dimension: number
      latency_ms: number
      reason: string
      search_used: boolean
    }
    search?: {
      route?: string
      used: boolean
      candidate_count: number
      reason: string
      fallback_reason?: string
      ranker?: {
        enabled: boolean
        provider?: string
        model_version?: string
        feature_schema_version?: string
        candidate_count?: number
        rerank_applied?: boolean
        candidate_order_mutated?: boolean
        reason: string
      }
    }
    candidates?: Candidate[]
    ranker?: {
      enabled: boolean
      provider?: string
      model_version?: string
      feature_schema_version?: string
      candidate_count?: number
      rerank_applied?: boolean
      candidate_order_mutated?: boolean
      reason: string
    }
  }
  metrics: {
    candidate_count: number
    shadow_overlap?: number | null
  }
}

export type RecommendationComparison = {
  run_id?: string
  scenario_id: string
  course_title: string
  confirmed_goal: string
  max_per_lesson: number
  lesson_count: number
  total_candidates: number
  total_lexical_hits: number
  total_vector_hits: number
  shadow_enabled: boolean
  shadow_attempted?: number
  shadow_succeeded?: number
  lessons: RecommendationCompareLesson[]
  latency_ms: number
}

export type DebugLabel = {
  id: string
  scenario_id: string
  run_id: string
  candidate_key: string
  content_id?: string
  url?: string
  label: string
  note: string
  baseline_rank?: number | null
  shadow_rank?: number | null
  updated_at: string
}

export type LabelSummary = {
  scenario_id: string
  latest_run_id: string
  latest_run_at?: string
  total_labels: number
  label_counts: Record<string, number>
  good_fit_rate: number
  irrelevant_rate: number
  duplicate_rate: number
  problem_rate: number
  quality_status: 'pending' | 'warning' | 'passed' | 'failed'
  quality_reason: string
  minimum_required: number
}

export type ScenarioInput = {
  course_title: string
  notes: string
}

export type SavedExternalCandidate = {
  content?: {
    id: string
    title: string
    url?: string | null
    external_source?: string | null
    external_content_id?: string | null
  }
  embedding_status?: string
}


export type RecommendationSpecMetricQuery = {
  query: string
  count: number
}

export type RecommendationSpecRolloutMetrics = {
  recommendation_exposed: number
  material_replaced: number
  broken_link_reported: number
  wrong_content_reported: number
  average_candidate_count?: number | null
}

export type RecommendationSpecMetrics = {
  window_days: number
  window_started_at: string
  window_ended_at: string
  recommendation_searches: number
  search_spec_source_counts: Record<string, number>
  resolution_source_counts: Record<string, number>
  empty_effective_query_rate?: number | null
  short_effective_query_rate?: number | null
  top_effective_queries: RecommendationSpecMetricQuery[]
  rollout_events: RecommendationSpecRolloutMetrics
}

export type LLMJobObservationStatus = 'queued' | 'running' | 'succeeded' | 'failed' | 'canceled' | 'expired'

export type LLMJobObservationCount = {
  feature: string
  status: LLMJobObservationStatus
  total_count: number
  stale_count: number
  oldest_created_at?: string | null
  oldest_activity_at?: string | null
}

export type LLMJobObservations = {
  generated_at: string
  policy: {
    queued_timeout_ms: number
    running_timeout_ms: number
  }
  counts: LLMJobObservationCount[]
}

export type LessonSearchPrerunTokenCount = {
  token: string
  count: number
}

export type LessonSearchPrerunSummary = {
  fixtures_total?: number
  fixtures_processed?: number
  fixtures_skipped?: number
  provider_errors?: number
  noise_suspects?: number
  relevant_like_candidates?: number
  risk_reports?: number
  risk_provider_errors?: number
  risk_noise_suspects?: number
  estimated_provider_requests?: number
  estimated_candidates_max?: number
  main_provider_requests?: number
  risk_provider_requests?: number
  automatic_correction_applied?: boolean
  risk_fixture_ids_configured?: string[]
  risk_fixture_ids_processed?: string[]
  quota?: {
    provider_filter?: string
    top_n?: number
    risk_top_n?: number
    risk_fixture_ids?: string[]
    main_provider_requests?: number
    risk_provider_requests?: number
    estimated_provider_requests?: number
    estimated_candidates_max?: number
  }
  [key: string]: unknown
}

export type LessonSearchPrerunNoiseTaxonomy = {
  registered_avoid_hits?: LessonSearchPrerunTokenCount[]
  unregistered_noise_hints?: LessonSearchPrerunTokenCount[]
  [key: string]: unknown
}

export type LessonSearchPrerunReport = {
  id: string
  run_type: string
  status: string
  started_at: string
  finished_at?: string | null
  provider_filter: string
  top_n: number
  risk_top_n: number
  summary: LessonSearchPrerunSummary
  noise_taxonomy_summary: LessonSearchPrerunNoiseTaxonomy
  error_code?: string | null
  error_message?: string | null
  created_at: string
}

export type LessonSearchPrerunReportsResponse = {
  limit: number
  include_report: boolean
  reports: LessonSearchPrerunReport[]
}
