import type { DashboardCourseRecord, DashboardDraftStatus } from '@/lib/world-ui-engine/types';
import type { PlanetStatus } from '@/components/dashboard/planetMapTokens';

export interface UserInfo {
  id: string;
  email: string;
  display_id: string;
  nickname: string;
  avatar_url?: string;
  role: string;
  premium_access: boolean;
  provider: string;
  created_at: string;
  free_points: number;
  paid_points: number;
  total_points: number;
  terms_agreed: boolean;
  privacy_agreed: boolean;
  required_consent_pending: boolean;
  ui_locale?: 'ko' | 'en';
  learning_language?: 'ko' | 'en';
  language_setup_required?: boolean;
}

export interface UserAISettings {
  mode: string;
  provider: string;
  has_api_key: boolean;
  endpoint_url?: string | null;
  updated_at?: string;
  last_validation_status?: string;
  last_validated_at?: string;
  last_validation_error?: string;
  last_failed_at?: string;
  next_retry_at?: string;
}

export type DraftStatus = DashboardDraftStatus;
export type PlanetFilter = 'all' | PlanetStatus;
export type CourseItem = DashboardCourseRecord;

export type TodayTaskKind =
  | 'create_course'
  | 'start_planet'
  | 'continue_point'
  | 'complete_planet'
  | 'review_records';

export interface TodayTask {
  kind: TodayTaskKind;
  title: string;
  description: string;
  cta_label: string;
  href: string;
  planet_id?: string | null;
  point_id?: string | null;
}

export interface TodayTaskResponse {
  task?: TodayTask | null;
}

export interface DraftListItem {
  id: string;
  source_query: string;
  title: string;
  status: DraftStatus;
  confirmed_course_id?: string | null;
  planet_type_id?: string | null;
  planet_type_name?: string | null;
  planet_type_asset?: string | null;
  planet_texture_map_id?: string | null;
  planet_texture_map_name?: string | null;
  planet_texture_map_asset?: string | null;
  planet_texture_map_rotation_duration_seconds?: number | null;
  planet_texture_map_rotation_direction?: 'left' | 'right' | null;
  updated_at: string;
}

export interface DraftResourceSummary {
  selection_state?: string;
  description?: string | null;
}

export interface DraftLessonSummary {
  lesson: {
    objective?: string | null;
    summary?: string | null;
  };
  resources?: DraftResourceSummary[];
}

export interface DraftLevelSummary {
  lessons: DraftLessonSummary[];
}

export interface DraftAggregateSummary {
  draft: {
    id: string;
    source_query: string;
    title: string;
    status?: DraftStatus;
    confirmed_course_id?: string | null;
    planet_type_id?: string | null;
    planet_type_name?: string | null;
    planet_type_asset?: string | null;
    planet_texture_map_id?: string | null;
    planet_texture_map_name?: string | null;
    planet_texture_map_asset?: string | null;
    planet_texture_map_rotation_duration_seconds?: number | null;
    planet_texture_map_rotation_direction?: 'left' | 'right' | null;
    updated_at: string;
  };
  levels: DraftLevelSummary[];
  lessons?: DraftLessonTreeSummary[];
}

export interface DraftPointTreeSummary {
  point: {
    selection_state?: string | null;
    description?: string | null;
  };
}

export interface DraftLessonTreeSummary {
  lesson: {
    objective?: string | null;
    summary?: string | null;
  };
  points?: DraftPointTreeSummary[];
  sub_lessons?: DraftLessonTreeSummary[];
}

export interface PlanetListItem {
  id: string;
  draft_id: string;
  title: string;
  status: PlanetStatus;
  planet_type_id?: string | null;
  planet_type_name?: string | null;
  planet_type_asset?: string | null;
  planet_texture_map_id?: string | null;
  planet_texture_map_name?: string | null;
  planet_texture_map_asset?: string | null;
  planet_texture_map_rotation_duration_seconds?: number | null;
  planet_texture_map_rotation_direction?: 'left' | 'right' | null;
  lesson_count: number;
  completed_lesson_count: number;
  progress?: number | null;
  can_complete?: boolean;
  is_inactive?: boolean;
  last_accessed_at?: string | null;
  updated_at: string;
  destination: string;
  share_count?: number;
}

export interface PlanetLessonSummary {
  lesson: {
    summary?: string | null;
  };
  resources?: Array<{
    description?: string | null;
  }>;
}

export interface PlanetLevelSummary {
  lessons: PlanetLessonSummary[];
}

export interface PlanetAggregateSummary {
  planet: PlanetListItem;
  levels: PlanetLevelSummary[];
  lessons?: PlanetLessonTreeSummary[];
}

export interface PlanetPointTreeSummary {
  point: {
    status?: 'draft' | 'ready' | 'learning' | 'completed';
    description?: string | null;
    item_status?: string | null;
    explorer_status?: string | null;
  };
}

export interface PlanetLessonTreeSummary {
  lesson: {
    summary?: string | null;
    status?: string | null;
  };
  points?: PlanetPointTreeSummary[];
  sub_lessons?: PlanetLessonTreeSummary[];
}
