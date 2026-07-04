// Types and interfaces for the point page
export type PlanetRouteKind = 'learning' | 'shared';
export type PlanetStatus = 'ready' | 'learning' | 'completed';
export type PointType = 'exploration' | 'research';
export type ResearchBlockType = 'text' | 'image' | 'link';
export type PointQuestionType = 'reflection' | 'application' | 'goal_alignment';
export type PointQuestionStatus = 'pending' | 'answered';
export type PointMaterialReportTarget = 'source' | 'ai_summary' | 'attachment' | 'other';
export type PointMaterialReportType = 'broken_link' | 'wrong_content' | 'unsafe_content' | 'copyright' | 'low_quality' | 'other';
export type PointMaterialReportStatus = 'open' | 'reviewing' | 'resolved' | 'dismissed' | 'cancelled';
export type ObservationNoteType = 'core_summary' | 'revisit_part' | 'reference_material';

export interface PointReplacementCandidate {
  content_id?: string | null;
  title: string;
  description?: string | null;
  thumbnail_url?: string | null;
  url?: string | null;
  content_type: string;
  rank_score: number;
}

export interface UserInfo {
  id: string;
  email: string;
  display_id: string;
  nickname: string;
  required_consent_pending: boolean;
}

export interface PlanetPoint {
  point: {
    id: string;
    point_type: PointType;
    status: 'draft' | 'ready' | 'learning' | 'completed';
    title: string;
    description?: string | null;
    point_goal?: string | null;
    point_category?: string | null;
    template_type?: 'concept_summary' | 'practice_strategy' | 'problem_solving' | 'free_research' | null;
    external_url?: string | null;
    thumbnail_url?: string | null;
    order_index: number;
  };
  research_material_confirmed?: boolean;
  research_material_confirmed_at?: string | null;
  journal_entry?: {
    observation: string;
    reflection: string;
    next_step: string;
    core_concept?: string;
    my_explanation?: string;
    examples?: string;
    confused_parts?: string;
    reference_links?: string;
    updated_at: string;
  } | null;
  record_entry?: {
    study_minutes: number;
    practice_count: number;
    confidence_level: number;
    application_note: string;
    updated_at: string;
  } | null;
  artifact_entry?: {
    artifact_type: string;
    title: string;
    url: string;
    description: string;
    updated_at: string;
  } | null;
  artifacts?: Array<{
    id: string;
    course_point_id: string;
    artifact_type: string;
    title: string;
    url: string;
    description: string;
    point_category?: string;
    production_process?: string;
    learned_points?: string;
    difficult_points?: string;
    visibility?: string;
    order_index: number;
    created_at: string;
    updated_at: string;
  }>;
  attachments?: Array<{
    id: string;
    course_point_id: string;
    artifact_id?: string | null;
    user_id: string;
    provider: 'creator' | 'platform' | 'learner';
    attachment_type: 'image' | 'file' | 'link' | 'code' | 'other' | 'video' | 'subtitle' | 'thumbnail';
    source_context?: 'research_material' | 'work_attachment' | 'artifact';
    title: string;
    url: string;
    file_path: string;
    file_size?: number | null;
    mime_type: string;
    created_at: string;
    updated_at: string;
  }>;
  material_reports?: Array<{
    id: string;
    course_point_id: string;
    user_id: string;
    attachment_id?: string | null;
    target_type: PointMaterialReportTarget;
    report_type: PointMaterialReportType;
    message: string;
    status: PointMaterialReportStatus;
    target_content_id?: string | null;
    target_url?: string | null;
    target_title?: string | null;
    replacement_content_id?: string | null;
    replacement_url?: string | null;
    replacement_title?: string | null;
    replaced_at?: string | null;
    created_at: string;
    updated_at: string;
  }>;
  practice_logs?: Array<{
    id: string;
    course_point_id: string;
    user_id: string;
    title: string;
    activity_name?: string;
    attempt_count?: number | null;
    success_count?: number | null;
    failure_count?: number | null;
    duration_minutes?: number | null;
    blocked_part: string;
    changed_method: string;
    achievement_note: string;
    achievement?: string;
    next_plan: string;
    next_practice?: string;
    created_at: string;
    updated_at: string;
  }>;

  observation_notes?: Array<{
    id: string;
    course_point_id: string;
    user_id: string;
    note_type: ObservationNoteType;
    content: string;
    order_index: number;
    created_at: string;
    updated_at: string;
  }>;
  events?: Array<{
    id: string;
    course_point_id: string;
    user_id: string;
    event_type: string;
    event_payload: Record<string, unknown> | null;
    created_at: string;
  }>;
  ai_summary_entry?: {
    source_title: string;
    source_description: string;
    summary: string;
    updated_at: string;
  } | null;
  blocks?: Array<{
    id: string;
    course_point_id: string;
    block_type: ResearchBlockType;
    content: Record<string, unknown> | null;
    order_index: number;
    created_at: string;
    updated_at: string;
  }>;
  questions?: Array<{
    id: string;
    course_point_id: string;
    goal_profile_version?: number | null;
    title?: string | null;
    question: string;
    question_type: PointQuestionType;
    answer_method?: string | null;
    answer?: string | null;
    ai_feedback?: string | null;
    status: PointQuestionStatus;
    created_by: string;
    created_at: string;
    updated_at: string;
  }>;
  self_evaluation?: {
    course_point_id: string;
    goal_profile_version?: number | null;
    understanding: number;
    application_note: string;
    proficiency: number;
    understanding_score?: number | null;
    understanding_reason?: string;
    application_score?: number | null;
    application_reason?: string;
    proficiency_score?: number | null;
    proficiency_reason?: string;
    problem_solving_score?: number | null;
    problem_solving_reason?: string;
    expression_score?: number | null;
    expression_reason?: string;
    goal_alignment_note: string;
    final_score?: number | null;
    updated_at: string;
  } | null;
}

export interface PlanetMetaState {
  planet: {
    id: string;
    draft_id: string;
    title: string;
    status: PlanetStatus;
    progress?: number | null;
  };
  goal_context?: {
    learning_goal?: string | null;
    goal_profile_id?: string | null;
    goal_profile_version?: number | null;
    confirmed_goal?: string | null;
    usage_context?: string | null;
    motivation?: string | null;
  } | null;
}

export interface ResolvedPoint {
  levelTitle: string;
  lessonTitle: string;
  lessonID: string;
  previousPointID?: string | null;
  nextPointID?: string | null;
  point: PlanetPoint['point'];
  research_material_confirmed?: boolean;
  research_material_confirmed_at?: string | null;
  journal_entry?: PlanetPoint['journal_entry'];
  record_entry?: PlanetPoint['record_entry'];
  artifact_entry?: PlanetPoint['artifact_entry'];
  artifacts?: NonNullable<PlanetPoint['artifacts']>;
  attachments?: NonNullable<PlanetPoint['attachments']>;
  material_reports?: NonNullable<PlanetPoint['material_reports']>;
  practice_logs?: NonNullable<PlanetPoint['practice_logs']>;
  events?: NonNullable<PlanetPoint['events']>;
  observation_notes?: NonNullable<PlanetPoint['observation_notes']>;
  ai_summary_entry?: PlanetPoint['ai_summary_entry'];
  blocks?: NonNullable<PlanetPoint['blocks']>;
  questions?: NonNullable<PlanetPoint['questions']>;
  self_evaluation?: PlanetPoint['self_evaluation'];
  template_type?: PlanetPoint['point']['template_type'];
}

export interface PlanetPointDetailResponse {
  planet: PlanetMetaState['planet'];
  goal_context?: PlanetMetaState['goal_context'];
  level_title: string;
  lesson_title: string;
  lesson_id: string;
  point: PlanetPoint;
  previous_point_id?: string | null;
  next_point_id?: string | null;
}

export interface LearningPointMutationResponse {
  planet?: PlanetMetaState['planet'];
  point?: PlanetPointDetailResponse;
  error?: string;
  feedback?: string;
  report?: {
    id: string;
    target_type: string;
    report_type: string;
    message: string;
    status: string;
    created_at: string;
  };
}

export interface LearningPointSelfEvaluationApplicationQuestion {
  question: string;
  intent: string;
}

export interface LearningPointSelfEvaluationAIDraft {
  application_questions?: LearningPointSelfEvaluationApplicationQuestion[];
  understanding_score: number;
  understanding_reason: string;
  application_score: number;
  application_reason: string;
  proficiency_score: number;
  proficiency_reason: string;
  problem_solving_score: number;
  problem_solving_reason: string;
  expression_score: number;
  expression_reason: string;
  goal_alignment_note: string;
}

export interface LearningPointSelfEvaluationAIDraftResponse {
  draft?: LearningPointSelfEvaluationAIDraft;
  error?: string;
}

export type PointWorkTab = 'notes' | 'questions' | 'practice' | 'artifacts' | 'attachments' | 'timeline';
