'use client';

export interface UserInfo {
  id: string;
  email: string;
  display_id: string;
  nickname: string;
  required_consent_pending: boolean;
  ui_locale?: string | null;
  learning_language?: string | null;
  language_setup_required?: boolean;
}

export type ResearchNodeTemplateType =
  | 'concept_summary'
  | 'practice_strategy'
  | 'problem_solving'
  | 'free_research';

export type ResearchNodeBlockType = 'text' | 'image' | 'link';

export interface DraftPointBlock {
  id: string;
  course_draft_point_id: string;
  block_type: ResearchNodeBlockType;
  content: Record<string, unknown>;
  order_index: number;
  created_at: string;
  updated_at: string;
}

export interface DraftPoint {
  id: string;
  course_draft_id: string;
  course_draft_lesson_id: string;
  point_type: 'exploration' | 'research';
  status: 'draft' | 'ready' | 'learning' | 'completed';
  title: string;
  description?: string | null;
  template_type?: ResearchNodeTemplateType | null;
  selection_state?: string | null;
  content_id?: string | null;
  external_url?: string | null;
  thumbnail_url?: string | null;
  price_type?: string | null;
  rank_score?: number | null;
  order_index: number;
  completed_at?: string | null;
}

export interface DraftPointAggregate {
  point: DraftPoint;
  blocks?: DraftPointBlock[];
  journal_entry?: {
    observation: string;
    reflection: string;
    next_step: string;
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
}

export interface DraftLessonMeta {
  id: string;
  title: string;
  summary?: string | null;
  operation_note?: string | null;
  difficulty_level?: string | null;
  journal_entry?: {
    observation: string;
    reflection: string;
    next_step: string;
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
  lesson_role: string;
  order_index?: number;
  parent_lesson_id?: string | null;
  course_draft_id?: string | null;
}

export interface DraftLessonTree {
  lesson: DraftLessonMeta;
  points?: DraftPointAggregate[];
  sub_lessons?: DraftLessonTree[];
}

export interface DraftAggregate {
  draft: {
    id: string;
    source_query: string;
    title: string;
    status?: 'draft' | 'confirmed' | 'learning' | 'archived';
    confirmed_course_id?: string | null;
    is_inactive?: boolean;
    description?: string | null;
    current_level?: string | null;
    preferred_format?: string | null;
    duration_weeks?: number | null;
    study_hours_per_week?: number | null;
    updated_at: string;
  };
  lessons: DraftLessonTree[];
}

export interface LessonSearchPointPreview {
  cost: number;
  free_balance: number;
  paid_balance: number;
  total_balance: number;
}

export interface LessonSearchResponse {
  draft: DraftAggregate;
  point_preview?: LessonSearchPointPreview;
  billing_status?: string;
  policy_cost?: number;
  error?: string;
}

export interface SelectDraftResourceResponse {
  draft?: DraftAggregate;
  error?: string;
}

export interface AttachDraftLessonResourceResponse {
  draft?: DraftAggregate;
  error?: string;
}

export interface UpdateDraftLessonResponse {
  draft?: DraftAggregate;
  error?: string;
}

export interface UpdateDraftLevelResponse {
  draft?: DraftAggregate;
  error?: string;
}

export interface UpdateDraftDetailMemoResponse {
  memo?: {
    target_type: string;
    target_id: string;
    main_lesson_id?: string;
    note: string;
    updated_at: string;
  };
  error?: string;
}

export interface UpdateDraftStructureResponse {
  draft?: DraftAggregate;
  error?: string;
}

export interface CreateResearchNodeResponse {
  point?: DraftPoint;
  research_node?: DraftPoint;
  error?: string;
}

export interface UpdateResearchNodeResponse {
  point?: DraftPoint;
  research_node?: DraftPoint;
  error?: string;
}

export interface DeleteResearchNodeResponse {
  error?: string;
}

export interface GetResearchPointBlocksResponse {
  blocks?: DraftPointBlock[];
  point_blocks?: DraftPointBlock[];
  error?: string;
}

export interface CreateResearchPointBlockResponse {
  block?: DraftPointBlock;
  point_block?: DraftPointBlock;
  error?: string;
}

export interface UpdateResearchPointBlockResponse {
  block?: DraftPointBlock;
  point_block?: DraftPointBlock;
  error?: string;
}

export interface UpdateDraftLessonJournalResponse {
  journal?: {
    target_id: string;
    observation: string;
    reflection: string;
    next_step: string;
    updated_at: string;
  };
  error?: string;
}

export interface UpdateDraftLessonRecordResponse {
  record?: {
    target_id: string;
    study_minutes: number;
    practice_count: number;
    confidence_level: number;
    application_note: string;
    updated_at: string;
  };
  error?: string;
}

export interface UpdateDraftLessonArtifactResponse {
  artifact?: {
    target_id: string;
    artifact_type: string;
    title: string;
    url: string;
    description: string;
    updated_at: string;
  };
  error?: string;
}

export interface ConfirmDraftResponse {
  draft?: DraftAggregate;
  error?: string;
}

export interface ManualResourceParsePreview {
  title?: string;
  description?: string;
  thumbnail_url?: string;
  author?: string;
  duration_seconds?: number;
  content_type?: string;
  source_url?: string;
  parse_error?: string;
}

export interface ManualResourceCreatedContent {
  id: string;
  title: string;
  content_type: string;
  url?: string | null;
  canonical_url?: string | null;
}

export type ManualResourceContentType = 'youtube' | 'blog' | 'article' | 'internal';

export type DiarySectionKey =
  | 'planning'
  | 'journal'
  | 'records'
  | 'artifacts'
  | 'community'
  | 'civilization';

export type DetailTargetKind = 'region' | 'point' | 'research_node';
