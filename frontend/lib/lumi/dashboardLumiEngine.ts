import type { LumiDockSlot, LumiState } from './lumiTypes';

export type DashboardPlanetStatus = 'draft' | 'ready' | 'learning' | 'completed';
export type DashboardStatusFilter = 'all' | DashboardPlanetStatus;

export interface DashboardLumiCourse {
  id: string;
  title: string;
  status: DashboardPlanetStatus;
}

export type DashboardLumiActionId =
  | 'new-course'
  | 'view-related-course'
  | 'show-draft'
  | 'show-ready'
  | 'show-learning'
  | 'show-completed'
  | 'clear-filter'
  | 'open-settings';

export interface DashboardLumiViewModelInput {
  courses: DashboardLumiCourse[];
  selectedCourseId: string | null;
  hoveredCourseId: string | null;
  activeStatusFilter: DashboardStatusFilter;
  userName: string;
  hasByok: boolean;
  isMobile: boolean;
}

export interface DashboardLumiViewModel {
  mood: LumiState;
  dockSlot: LumiDockSlot;
  visible: boolean;
  message: string;
  subMessage: string;
  actions: DashboardLumiActionId[];
  relatedCourseId: string | null;
}

function statusSummary(status: DashboardPlanetStatus): string {
  if (status === 'draft') return '탐험 계획을 조금만 더 다듬으면 바로 출발할 수 있어요.';
  if (status === 'ready') return '준비가 끝난 행성이에요. 지금 바로 탐험을 시작할 수 있어요.';
  if (status === 'learning') return '지금 진행 중인 탐험이에요. 이어서 계속 가면 됩니다.';
  return '탐험을 완료한 행성이에요. 다시 돌아보거나 다음 행성으로 넘어갈 수 있어요.';
}

function filterSummary(filter: DashboardStatusFilter): string {
  if (filter === 'draft') return '지금은 탐험계획 행성만 보고 있어요.';
  if (filter === 'ready') return '지금은 탐험대기 행성만 보고 있어요.';
  if (filter === 'learning') return '지금은 행성탐험 행성만 보고 있어요.';
  if (filter === 'completed') return '지금은 행성탐험완료 행성만 보고 있어요.';
  return '행성 전체 흐름을 보고 있어요.';
}

function topLevelDockSlot(isMobile: boolean, hasHighlight: boolean): LumiDockSlot {
  if (isMobile) return 'mobile-bottom-sheet-trigger';
  return hasHighlight ? 'dashboard-map-bottom' : 'dashboard-map-top';
}

export function resolveLumiViewModel(input: DashboardLumiViewModelInput): DashboardLumiViewModel {
  const { courses, selectedCourseId, hoveredCourseId, activeStatusFilter, userName, hasByok, isMobile } = input;
  const selectedCourse = selectedCourseId ? courses.find((course) => course.id === selectedCourseId) ?? null : null;
  const hoveredCourse = hoveredCourseId ? courses.find((course) => course.id === hoveredCourseId) ?? null : null;

  if (selectedCourse) {
    return {
      mood: 'planet-hold',
      dockSlot: topLevelDockSlot(isMobile, true),
      visible: !isMobile,
      message: `${selectedCourse.title}`,
      subMessage: statusSummary(selectedCourse.status),
      actions: ['view-related-course', activeStatusFilter === 'all' ? `show-${selectedCourse.status}` as DashboardLumiActionId : 'clear-filter'],
      relatedCourseId: selectedCourse.id,
    };
  }

  if (hoveredCourse) {
    return {
      mood: 'planet-point',
      dockSlot: topLevelDockSlot(isMobile, true),
      visible: !isMobile,
      message: `${hoveredCourse.title}`,
      subMessage: statusSummary(hoveredCourse.status),
      actions: ['view-related-course'],
      relatedCourseId: hoveredCourse.id,
    };
  }

  const draftCourses = courses.filter((course) => course.status === 'draft');
  const readyCourses = courses.filter((course) => course.status === 'ready');
  const learningCourses = courses.filter((course) => course.status === 'learning');
  const completedCourses = courses.filter((course) => course.status === 'completed');

  if (courses.length === 0) {
    return {
      mood: hasByok ? 'raise-hand' : 'curious',
      dockSlot: topLevelDockSlot(isMobile, false),
      visible: !isMobile,
      message: `${userName}님, 첫 행성을 만들어볼까요?`,
      subMessage: hasByok
        ? '아직 생성된 행성은 없어요. 원하는 주제를 입력하면 바로 탐험 계획을 만들 수 있어요.'
        : 'BYOK 없이도 시작할 수 있어요. 원하면 설정에서 본인 키를 연결할 수도 있어요.',
      actions: hasByok ? ['new-course'] : ['new-course', 'open-settings'],
      relatedCourseId: null,
    };
  }

  if (learningCourses.length > 0) {
    return {
      mood: 'exploring',
      dockSlot: topLevelDockSlot(isMobile, false),
      visible: !isMobile,
      message: `${userName}님이 탐험 중인 행성이 ${learningCourses.length}개 있어요.`,
      subMessage: filterSummary(activeStatusFilter),
      actions: ['view-related-course', activeStatusFilter === 'learning' ? 'new-course' : 'show-learning'],
      relatedCourseId: learningCourses[0]?.id ?? null,
    };
  }

  if (readyCourses.length > 0) {
    return {
      mood: 'focus',
      dockSlot: topLevelDockSlot(isMobile, false),
      visible: !isMobile,
      message: `출발 준비가 끝난 행성이 ${readyCourses.length}개 있어요.`,
      subMessage: filterSummary(activeStatusFilter),
      actions: ['view-related-course', activeStatusFilter === 'ready' ? 'new-course' : 'show-ready'],
      relatedCourseId: readyCourses[0]?.id ?? null,
    };
  }

  if (draftCourses.length > 0) {
    return {
      mood: 'thinking',
      dockSlot: topLevelDockSlot(isMobile, false),
      visible: !isMobile,
      message: `정리 중인 탐험 계획이 ${draftCourses.length}개 있어요.`,
      subMessage: filterSummary(activeStatusFilter),
      actions: ['view-related-course', activeStatusFilter === 'draft' ? 'new-course' : 'show-draft'],
      relatedCourseId: draftCourses[0]?.id ?? null,
    };
  }

  return {
    mood: 'victory',
    dockSlot: topLevelDockSlot(isMobile, false),
    visible: !isMobile,
    message: `완료한 행성이 ${completedCourses.length}개 있어요.`,
    subMessage: filterSummary(activeStatusFilter),
    actions: ['view-related-course', activeStatusFilter === 'completed' ? 'new-course' : 'show-completed'],
    relatedCourseId: completedCourses[0]?.id ?? null,
  };
}
