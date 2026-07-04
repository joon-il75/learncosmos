export type LumiMode = 'general' | 'ai';

export type LumiPage =
  | 'hero'
  | 'dashboard'
  | 'planet_map'
  | 'course_draft_editor'
  | 'course_player';

export type LumiAction =
  | 'page_enter'
  | 'open_lumi'
  | 'goal_submit'
  | 'create_course'
  | 'search_content'
  | 'ai_generate'
  | 'ai_response'
  | 'lesson_complete'
  | 'course_complete'
  | 'error';

export type LumiScene =
  | 'empty'
  | 'overview'
  | 'editing'
  | 'searching'
  | 'waiting'
  | 'success'
  | 'error';

export interface LumiRuntimeCourseCounts {
  draft: number;
  ready: number;
  learning: number;
  completed: number;
  inactive: number;
}

// 'system'       : Galaxy 모드에서 항성계 노드에 직접 hover
// 'course'        : Star System 모드에서 행성/리스트 코스에 hover
// 'course-in-galaxy': Galaxy 모드에서 코스 리스트 항목에 hover (항성계 강조 hint 포함)
export type DashboardHoverSource = 'system' | 'course' | 'course-in-galaxy' | null;

export interface DashboardHoveredSystemSummary {
  title: string;
  statusCount: LumiRuntimeCourseCounts;
  updatedAt?: string | null;
  leadCourseId?: string | null;
  leadCourseTitle?: string | null;
}

export interface LumiRuntimeUserState {
  courseCounts?: LumiRuntimeCourseCounts;
  progressPercent?: number;
  hasBYOK?: boolean;
  currentCourseId?: string;
  currentLessonId?: string;
  hoveredSystem?: DashboardHoveredSystemSummary;
  hoverSource?: DashboardHoverSource;
}

export interface LumiRuntimeContext {
  mode: LumiMode;
  page: LumiPage;
  action: LumiAction;
  scene?: LumiScene;
  userState?: LumiRuntimeUserState;
}
