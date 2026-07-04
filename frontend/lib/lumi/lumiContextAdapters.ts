import type {
  DashboardHoveredSystemSummary,
  DashboardHoverSource,
  LumiAction,
  LumiPage,
  LumiRuntimeContext,
  LumiScene,
} from './lumiEngineTypes';
import type { LumiPageContext, LumiTrigger, LumiTriggerPayload } from './lumiTypes';

export interface DashboardRuntimeContextInput {
  courses: Array<{ id: string; status: 'draft' | 'ready' | 'learning' | 'completed'; isInactive?: boolean }>;
  selectedCourseId: string | null;
  hoveredCourseId: string | null;
  activeStatusFilter: 'all' | 'draft' | 'ready' | 'learning' | 'completed';
  hasByok: boolean;
  hoverSource?: DashboardHoverSource | null;
  hoveredSystemSummary?: DashboardHoveredSystemSummary | null;
}

export function mapLegacyContextToPage(context: LumiPageContext): LumiPage {
  if (context === 'dashboard') return 'dashboard';
  if (context === 'planet-map') return 'planet_map';
  if (context === 'diary') return 'course_draft_editor';
  return 'course_player';
}

function mapTriggerToAction(trigger: LumiTrigger): LumiAction {
  if (trigger === 'manual_open') return 'open_lumi';
  if (trigger === 'cta_submit') return 'goal_submit';
  if (trigger === 'generation_success') return 'create_course';
  if (trigger === 'search_start' || trigger === 'search_success') return 'search_content';
  if (trigger === 'generation_loading' || trigger === 'evaluation_start') return 'ai_generate';
  if (trigger === 'evaluation_success' || trigger === 'note_saved') return 'ai_response';
  if (trigger === 'lesson_complete' || trigger === 'level_complete') return 'lesson_complete';
  if (trigger === 'course_complete' || trigger === 'base_built') return 'course_complete';
  if (trigger === 'cta_focus' || trigger === 'planet_hover' || trigger === 'planet_select' || trigger === 'lesson_enter') return 'open_lumi';
  return 'page_enter';
}

function mapTriggerToScene(trigger: LumiTrigger, payload?: LumiTriggerPayload): LumiScene {
  if (trigger === 'generation_loading' || trigger === 'evaluation_start') return 'waiting';
  if (trigger === 'search_start') return 'searching';
  if (trigger === 'cta_focus' || trigger === 'cta_submit') return 'editing';
  if (trigger === 'generation_success' || trigger === 'search_success' || trigger === 'evaluation_success' || trigger === 'note_saved' || trigger === 'lesson_complete' || trigger === 'level_complete' || trigger === 'course_complete' || trigger === 'base_built') return 'success';
  if (payload?.status) return 'overview';
  return 'overview';
}

function mapTriggerToMode(trigger: LumiTrigger): 'general' | 'ai' {
  if (trigger === 'search_start' || trigger === 'search_success' || trigger === 'generation_loading' || trigger === 'evaluation_start' || trigger === 'evaluation_success') {
    return 'ai';
  }
  return 'general';
}

export function fromTriggerToLumiRuntimeContext(trigger: LumiTrigger, payload?: LumiTriggerPayload): LumiRuntimeContext {
  const page = mapLegacyContextToPage(payload?.context ?? 'dashboard');
  const currentCourseId = payload?.courseTitle ? payload.courseTitle : undefined;
  const currentLessonId = payload?.lessonTitle ? payload.lessonTitle : undefined;

  return {
    mode: mapTriggerToMode(trigger),
    page,
    action: mapTriggerToAction(trigger),
    scene: mapTriggerToScene(trigger, payload),
    userState: currentCourseId || currentLessonId
      ? {
          currentCourseId,
          currentLessonId,
        }
      : undefined,
  };
}

export function fromDashboardStateToLumiRuntimeContext(input: DashboardRuntimeContextInput): LumiRuntimeContext {
  const counts = input.courses.reduce(
    (acc, course) => {
      if (course.isInactive) {
        acc.inactive += 1;
      }
      acc[course.status] += 1;
      return acc;
    },
    { draft: 0, ready: 0, learning: 0, completed: 0, inactive: 0 },
  );

  const currentCourseId = input.selectedCourseId ?? input.hoveredCourseId ?? undefined;
  const scene = input.courses.length === 0 ? 'empty' : currentCourseId ? 'overview' : 'overview';
  const action = currentCourseId ? 'open_lumi' : 'page_enter';

  return {
    mode: 'general',
    page: 'dashboard',
    action,
    scene,
    userState: {
      courseCounts: counts,
      hasBYOK: input.hasByok,
      currentCourseId,
      hoveredSystem: input.hoveredSystemSummary ?? undefined,
      hoverSource: input.hoverSource ?? null,
    },
  };
}
