import { resolveMobileDockSlot } from './lumiDockSlots';
import type {
  LumiRuntimeContext,
  LumiRuntimeCourseCounts,
} from './lumiEngineTypes';
import { resolveLumiMessageTemplate } from './lumiMessageTemplates';
import { lumiQuickActionTemplates } from './lumiQuickActions';
import { resolveLumiRuleDraft, type LumiRuleDraft } from './lumiRuleDrafts';
import type {
  LumiDockSlot,
  LumiMessageType,
  LumiPageContext,
  LumiQuickActionTemplate,
  LumiState,
  LumiTrigger,
  LumiTriggerPayload,
} from './lumiTypes';

export type DashboardPlanetStatus = 'draft' | 'ready' | 'learning' | 'completed';
export type DashboardStatusFilter = 'all' | DashboardPlanetStatus;

export interface DashboardLumiCourse {
  id: string;
  title: string;
  status: DashboardPlanetStatus;
  progress: number | null;
  lessonCount: number;
  levelCount: number;
  plannedLessonCount: number;
  completedLessonCount: number;
}

export type DashboardLumiActionId =
  | 'new-course'
  | 'view-related-course'
  | 'share-course'
  | 'view-expedition-base'
  | 'show-draft'
  | 'show-ready'
  | 'show-learning'
  | 'show-completed'
  | 'clear-filter'
  | 'open-settings';

interface BaseEngineResult {
  state: LumiState;
  context: LumiPageContext;
  dockSlot: LumiDockSlot;
  visible: boolean;
  expanded: boolean;
  message: string;
  preview: string;
  messageType: LumiMessageType;
  relatedCourseId: string | null;
  runtimeContext: LumiRuntimeContext;
}

export interface TriggerLumiEngineResult extends BaseEngineResult {
  quickActionTemplates: LumiQuickActionTemplate[];
  dashboardActionIds: [];
}

export interface DashboardLumiEngineResult extends BaseEngineResult {
  quickActionTemplates: [];
  dashboardActionIds: DashboardLumiActionId[];
}

interface TriggerEngineInput {
  kind: 'trigger';
  runtimeContext: LumiRuntimeContext;
  trigger: LumiTrigger;
  payload?: LumiTriggerPayload;
  isMobile: boolean;
  expanded?: boolean;
}

interface DashboardEngineInput {
  kind: 'dashboard';
  runtimeContext: LumiRuntimeContext;
  courses: DashboardLumiCourse[];
  selectedCourseId: string | null;
  hoveredCourseId: string | null;
  activeStatusFilter: DashboardStatusFilter;
  userName: string;
  hasByok: boolean;
  isMobile: boolean;
  expanded?: boolean;
}

export type LumiEngineInput = TriggerEngineInput | DashboardEngineInput;
export type LumiEngineResult = TriggerLumiEngineResult | DashboardLumiEngineResult;

const mobileVisibleTriggers = new Set<LumiTrigger>([
  'cta_focus',
  'cta_submit',
  'generation_loading',
  'generation_success',
  'idle',
  'manual_open',
  'course_complete',
  'base_built',
]);

function mapRuntimePageToLegacyContext(page: LumiRuntimeContext['page']): LumiPageContext {
  if (page == 'dashboard' || page == 'hero') return 'dashboard';
  if (page == 'planet_map') return 'planet-map';
  if (page == 'course_draft_editor') return 'diary';
  return 'player';
}

function defaultDockSlotForContext(context: LumiPageContext): LumiDockSlot {
  if (context == 'diary') return 'diary-header';
  if (context == 'player') return 'player-inline';
  if (context == 'planet-scene') return 'planet-scene-anchor';
  if (context == 'completion') return 'completion-center';
  if (context == 'mobile') return 'mobile-bottom-sheet-trigger';
  return 'dashboard-map-top';
}

function triggerDockSlot(trigger: LumiTrigger, context: LumiPageContext): LumiDockSlot {
  if (trigger == 'cta_focus' || trigger == 'cta_submit' || trigger == 'generation_loading' || trigger == 'generation_success') {
    return 'dashboard-cta';
  }
  if (trigger == 'planet_hover' || trigger == 'planet_select') {
    return 'dashboard-map-bottom';
  }
  if (trigger == 'search_start' || trigger == 'search_success') {
    return 'diary-sidebar';
  }
  if (trigger == 'evaluation_start' || trigger == 'evaluation_success') {
    return 'diary-header';
  }
  if (trigger == 'lesson_enter' || trigger == 'note_saved' || trigger == 'lesson_complete') {
    return 'player-inline';
  }
  if (trigger == 'course_complete' || trigger == 'base_built') {
    return 'completion-center';
  }
  return defaultDockSlotForContext(context);
}

function defaultStateForTrigger(trigger: LumiTrigger): LumiState {
  switch (trigger) {
    case 'cta_focus':
      return 'raise-hand';
    case 'cta_submit':
      return 'curious';
    case 'generation_loading':
    case 'evaluation_start':
      return 'thinking';
    case 'generation_success':
      return 'planet-hold';
    case 'planet_hover':
    case 'planet_select':
      return 'planet-point';
    case 'idle':
      return 'sleepy';
    case 'search_start':
      return 'exploring';
    case 'search_success':
      return 'surprise';
    case 'evaluation_success':
    case 'note_saved':
      return 'happy';
    case 'lesson_enter':
      return 'note-read';
    case 'lesson_complete':
    case 'level_complete':
      return 'celebrate';
    case 'course_complete':
    case 'base_built':
      return 'victory';
    case 'manual_open':
    case 'page_enter':
    default:
      return 'idle';
  }
}

function triggerPresentation(
  trigger: LumiTrigger,
  runtimeContext: LumiRuntimeContext,
  payload?: LumiTriggerPayload,
  matchedRule?: LumiRuleDraft | null,
): Pick<TriggerLumiEngineResult, 'state' | 'message' | 'messageType' | 'preview'> {
  const messageTemplate = resolveLumiMessageTemplate(trigger, payload) ?? resolveLumiMessageTemplate('page_enter');

  return {
    state: matchedRule?.state ?? defaultStateForTrigger(trigger),
    message: messageTemplate?.message ?? matchedRule?.message ?? '',
    preview: messageTemplate?.preview ?? '',
    messageType: matchedRule?.messageType ?? messageTemplate?.messageType ?? 'summary',
  };
}

// spec 기준 행성 hover 메시지 (PHASE 3)
function planetHoverMessage(course: DashboardLumiCourse): { message: string; preview: string } {
  const progressPct = course.progress != null ? Math.round(course.progress * 100) : null;
  if (course.status == 'draft') {
    return {
      message: `${course.title} 행성은 아직 탐험계획을 차근차근 정리하는 중이에요.`,
      preview: progressPct != null ? `탐험계획 ${progressPct}%를 함께 다듬고 있어요.` : `${course.lessonCount}개 지역을 준비하고 있어요.`,
    };
  }
  if (course.status == 'ready') {
    return {
      message: `${course.title} 행성은 탐험을 시작할 준비를 잘 마쳤어요.`,
      preview: `${course.lessonCount}개의 지역이 차분히 기다리고 있어요.`,
    };
  }
  if (course.status == 'learning') {
    return {
      message: `${course.title} 행성에서는 지금 탐험을 차근차근 이어가고 있어요.`,
      preview: progressPct != null ? `탐험진행 ${progressPct}%까지 함께 왔어요.` : `${course.completedLessonCount}/${course.lessonCount}개 지역을 탐험했어요.`,
    };
  }
  return {
    message: `${course.title} 행성의 탐험은 모두 잘 마무리됐어요.`,
    preview: '행성이 환하게 살아났어요.',
  };
}

function filterSummary(filter: DashboardStatusFilter): string {
  if (filter == 'draft') return '지금은 탐험계획 행성만 보고 있어요.';
  if (filter == 'ready') return '지금은 탐험대기 행성만 보고 있어요.';
  if (filter == 'learning') return '지금은 행성탐험 행성만 보고 있어요.';
  if (filter == 'completed') return '지금은 행성탐험완료 행성만 보고 있어요.';
  return '행성 전체 흐름을 보고 있어요.';
}

// 항성계 hover 메시지 — tooltip / 좌하단 Lumi 공통
function systemHoverMessage(
  systemTitle: string,
  counts: LumiRuntimeCourseCounts,
): { message: string; preview: string } {
  const previewParts: string[] = [];
  if (counts.draft) previewParts.push(`계획 ${counts.draft}`);
  if (counts.ready) previewParts.push(`대기 ${counts.ready}`);
  if (counts.learning) previewParts.push(`진행 ${counts.learning}`);
  if (counts.completed) previewParts.push(`완료 ${counts.completed}`);
  if (counts.inactive) previewParts.push(`비활성 ${counts.inactive}`);
  const preview = previewParts.join(' · ') || '아직 행성이 없습니다.';

  const activeKinds = [counts.draft > 0, counts.ready > 0, counts.learning > 0, counts.completed > 0, counts.inactive > 0].filter(Boolean).length;

  if (counts.learning > 0 && activeKinds === 1) {
    return { message: `${systemTitle}는 지금 탐험이 활발해서, 바로 이어가기 좋은 항성계예요.`, preview };
  }
  if (counts.completed > 0 && activeKinds === 1) {
    return { message: `${systemTitle}에는 탐험을 잘 마친 행성들이 모여 있어요.`, preview };
  }
  if (counts.ready > 0 && activeKinds === 1) {
    return { message: `${systemTitle}의 행성들이 탐험 시작을 차분히 기다리고 있어요.`, preview };
  }
  if (counts.draft > 0 && activeKinds === 1) {
    return { message: `${systemTitle}의 행성들이 탐험계획을 함께 준비하고 있어요.`, preview };
  }
  if (counts.inactive > 0 && activeKinds === 1) {
    return { message: `${systemTitle}에는 비활성화한 행성들이 모여 있어요.`, preview };
  }
  return { message: `${systemTitle}는 여러 단계의 행성이 함께 있는 항성계예요.`, preview };
}

function systemPreviewSummary(counts: LumiRuntimeCourseCounts): string {
  const pieces: string[] = [];
  if (counts.draft) pieces.push(`계획 ${counts.draft}`);
  if (counts.ready) pieces.push(`대기 ${counts.ready}`);
  if (counts.learning) pieces.push(`진행 ${counts.learning}`);
  if (counts.completed) pieces.push(`완료 ${counts.completed}`);
  if (counts.inactive) pieces.push(`비활성 ${counts.inactive}`);
  return pieces.length > 0 ? pieces.join(' · ') : '아직 행성이 없습니다.';
}

function topLevelDockSlot(isMobile: boolean, hasHighlight: boolean): LumiDockSlot {
  if (isMobile) return 'mobile-bottom-sheet-trigger';
  return hasHighlight ? 'dashboard-map-bottom' : 'dashboard-map-top';
}

function resolveDashboard(input: DashboardEngineInput): DashboardLumiEngineResult {
  const { courses, selectedCourseId, hoveredCourseId, activeStatusFilter, userName, hasByok, isMobile, expanded, runtimeContext } = input;
  const selectedCourse = selectedCourseId ? courses.find((course) => course.id == selectedCourseId) ?? null : null;
  const hoveredCourse = hoveredCourseId ? courses.find((course) => course.id == hoveredCourseId) ?? null : null;
  const hoveredSystem = runtimeContext.userState?.hoveredSystem ?? null;
  const hoverSource = runtimeContext.userState?.hoverSource ?? null;

  // PHASE 2: Star System 행성/코스 hover — 최우선 (hoveredCourseCard > hoveredPlanet > hoveredSystem > selectedCourse)
  if (hoverSource === 'course' && hoveredCourse) {
    const { message, preview } = planetHoverMessage(hoveredCourse);
    return {
      state: 'planet-point',
      context: isMobile ? 'mobile' : 'dashboard',
      dockSlot: topLevelDockSlot(isMobile, true),
      visible: !isMobile,
      expanded: expanded ?? false,
      message,
      preview,
      messageType: 'guide',
      dashboardActionIds: ['view-related-course'],
      quickActionTemplates: [],
      relatedCourseId: hoveredCourse.id,
      runtimeContext,
    };
  }

  // Galaxy 모드 리스트 코스 hover → 코스 요약 + 항성계 맥락 hint (PHASE 3)
  if (hoverSource === 'course-in-galaxy' && hoveredCourse) {
    const { message, preview } = planetHoverMessage(hoveredCourse);
    const systemHint = hoveredSystem ? ` (${hoveredSystem.title})` : '';
    return {
      state: 'planet-point',
      context: isMobile ? 'mobile' : 'dashboard',
      dockSlot: topLevelDockSlot(isMobile, true),
      visible: !isMobile,
      expanded: expanded ?? false,
      message: message + systemHint,
      preview,
      messageType: 'guide',
      dashboardActionIds: ['view-related-course'],
      quickActionTemplates: [],
      relatedCourseId: hoveredCourse.id,
      runtimeContext,
    };
  }

  // Galaxy 모드 항성계 노드 hover → 항성계 요약 (PHASE 4)
  if (hoverSource === 'system' && hoveredSystem) {
    const { message, preview } = systemHoverMessage(hoveredSystem.title, hoveredSystem.statusCount);
    return {
      state: 'planet-point',
      context: isMobile ? 'mobile' : 'dashboard',
      dockSlot: topLevelDockSlot(isMobile, true),
      visible: !isMobile,
      expanded: expanded ?? false,
      message,
      preview,
      messageType: 'guide',
      dashboardActionIds: hoveredSystem.leadCourseId ? ['view-related-course'] : [],
      quickActionTemplates: [],
      relatedCourseId: hoveredSystem.leadCourseId ?? null,
      runtimeContext,
    };
  }

  // hoveredCourse fallback (hoverSource 미지정 시)
  if (hoveredCourse) {
    const { message, preview } = planetHoverMessage(hoveredCourse);
    return {
      state: 'planet-point',
      context: isMobile ? 'mobile' : 'dashboard',
      dockSlot: topLevelDockSlot(isMobile, true),
      visible: !isMobile,
      expanded: expanded ?? false,
      message,
      preview,
      messageType: 'guide',
      dashboardActionIds: ['view-related-course'],
      quickActionTemplates: [],
      relatedCourseId: hoveredCourse.id,
      runtimeContext,
    };
  }

  if (selectedCourse) {
    const { message, preview } = planetHoverMessage(selectedCourse);
    // PHASE 5: completed 코스 선택 시 share-course / view-expedition-base 노출
    const actionIds: DashboardLumiActionId[] = selectedCourse.status === 'completed'
      ? ['view-related-course', 'share-course', 'view-expedition-base']
      : ['view-related-course', activeStatusFilter == 'all' ? (`show-${selectedCourse.status}` as DashboardLumiActionId) : 'clear-filter'];
    return {
      state: 'planet-hold',
      context: isMobile ? 'mobile' : 'dashboard',
      dockSlot: topLevelDockSlot(isMobile, true),
      visible: !isMobile,
      expanded: expanded ?? false,
      message,
      preview,
      messageType: 'guide',
      dashboardActionIds: actionIds,
      quickActionTemplates: [],
      relatedCourseId: selectedCourse.id,
      runtimeContext,
    };
  }

  const draftCourses = courses.filter((course) => course.status == 'draft');
  const readyCourses = courses.filter((course) => course.status == 'ready');
  const learningCourses = courses.filter((course) => course.status == 'learning');
  const completedCourses = courses.filter((course) => course.status == 'completed');

  if (courses.length == 0) {
    return {
      state: hasByok ? 'raise-hand' : 'curious',
      context: isMobile ? 'mobile' : 'dashboard',
      dockSlot: topLevelDockSlot(isMobile, false),
      visible: !isMobile,
      expanded: expanded ?? false,
      message: `${userName}님, 첫 행성을 함께 만들어 볼까요?`,
      preview: hasByok
        ? '아직 생성된 행성은 없어요. 원하는 주제를 적어 주시면 바로 탐험 계획을 안내해 드릴게요.'
        : 'BYOK 없이도 바로 시작할 수 있어요. 원하시면 설정에서 본인 키를 연결하셔도 좋아요.',
      messageType: 'question',
      dashboardActionIds: hasByok ? ['new-course'] : ['new-course', 'open-settings'],
      quickActionTemplates: [],
      relatedCourseId: null,
      runtimeContext,
    };
  }

  if (learningCourses.length > 0) {
    return {
      state: 'exploring',
      context: isMobile ? 'mobile' : 'dashboard',
      dockSlot: topLevelDockSlot(isMobile, false),
      visible: !isMobile,
      expanded: expanded ?? false,
      message: `${userName}님이 이어가고 있는 탐험 행성이 ${learningCourses.length}개 있어요.`,
      preview: filterSummary(activeStatusFilter),
      messageType: 'guide',
      dashboardActionIds: ['view-related-course', activeStatusFilter == 'learning' ? 'new-course' : 'show-learning'],
      quickActionTemplates: [],
      relatedCourseId: learningCourses[0]?.id ?? null,
      runtimeContext,
    };
  }

  if (readyCourses.length > 0) {
    return {
      state: 'focus',
      context: isMobile ? 'mobile' : 'dashboard',
      dockSlot: topLevelDockSlot(isMobile, false),
      visible: !isMobile,
      expanded: expanded ?? false,
      message: `출발 준비를 마친 행성이 ${readyCourses.length}개 있어요.`,
      preview: filterSummary(activeStatusFilter),
      messageType: 'guide',
      dashboardActionIds: ['view-related-course', activeStatusFilter == 'ready' ? 'new-course' : 'show-ready'],
      quickActionTemplates: [],
      relatedCourseId: readyCourses[0]?.id ?? null,
      runtimeContext,
    };
  }

  if (draftCourses.length > 0) {
    return {
      state: 'thinking',
      context: isMobile ? 'mobile' : 'dashboard',
      dockSlot: topLevelDockSlot(isMobile, false),
      visible: !isMobile,
      expanded: expanded ?? false,
      message: `차근차근 다듬고 있는 탐험 계획이 ${draftCourses.length}개 있어요.`,
      preview: filterSummary(activeStatusFilter),
      messageType: 'guide',
      dashboardActionIds: ['view-related-course', activeStatusFilter == 'draft' ? 'new-course' : 'show-draft'],
      quickActionTemplates: [],
      relatedCourseId: draftCourses[0]?.id ?? null,
      runtimeContext,
    };
  }

  return {
    state: 'victory',
    context: isMobile ? 'mobile' : 'dashboard',
    dockSlot: topLevelDockSlot(isMobile, false),
    visible: !isMobile,
    expanded: expanded ?? false,
    message: `탐험을 잘 마친 행성이 ${completedCourses.length}개 있어요.`,
    preview: filterSummary(activeStatusFilter),
    messageType: 'success',
    dashboardActionIds: ['view-related-course', 'share-course', activeStatusFilter == 'completed' ? 'new-course' : 'show-completed'],
    quickActionTemplates: [],
    relatedCourseId: completedCourses[0]?.id ?? null,
    runtimeContext,
  };
}

function resolveTrigger(input: TriggerEngineInput): TriggerLumiEngineResult {
  const matchedRule = resolveLumiRuleDraft({
    runtimeContext: input.runtimeContext,
    trigger: input.trigger,
    payload: input.payload,
  });
  const baseContext = mapRuntimePageToLegacyContext(input.runtimeContext.page);
  const context = input.payload?.context ?? matchedRule?.context ?? baseContext;
  const rawDockSlot = input.payload?.dockSlot ?? matchedRule?.dockSlot ?? triggerDockSlot(input.trigger, context);
  const dockSlot = input.isMobile ? resolveMobileDockSlot(rawDockSlot) : rawDockSlot;
  const outputContext = input.isMobile ? 'mobile' : context;
  const base = triggerPresentation(input.trigger, input.runtimeContext, input.payload, matchedRule);

  return {
    ...base,
    context: outputContext,
    dockSlot,
    visible: input.isMobile ? mobileVisibleTriggers.has(input.trigger) : true,
    expanded: input.expanded ?? false,
    quickActionTemplates: lumiQuickActionTemplates[outputContext] ?? [],
    dashboardActionIds: [],
    relatedCourseId: null,
    runtimeContext: input.runtimeContext,
  };
}

export function resolveLumiRuntime(input: LumiEngineInput): LumiEngineResult {
  if (input.kind == 'dashboard') return resolveDashboard(input);
  return resolveTrigger(input);
}
