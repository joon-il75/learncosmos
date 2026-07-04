import type {
  DashboardAdapterResult,
  DashboardLumiActionHandlers,
  DashboardLumiControllerInput,
} from './types';
import type {
  DashboardHoveredSystemSummary,
  DashboardHoverSource,
} from '@/lib/lumi/lumiEngineTypes';

interface BuildDashboardLumiAdapterInput {
  dashboard: DashboardAdapterResult;
  userName: string;
  actionHandlers: DashboardLumiActionHandlers;
  hoverSource: DashboardHoverSource | null;
  hoveredSystemSummary: DashboardHoveredSystemSummary | null;
}

export function buildDashboardLumiAdapter(
  input: BuildDashboardLumiAdapterInput,
): DashboardLumiControllerInput {
  return {
    courses: input.dashboard.snapshot.courses.filteredCourses,
    selectedCourseId: input.dashboard.snapshot.selection.selectedCourse?.id ?? null,
    hoveredCourseId: input.dashboard.snapshot.selection.hoveredCourse?.id ?? null,
    activeStatusFilter: input.dashboard.input.filter.activeStatusFilter,
    userName: input.userName,
    hasByok: input.dashboard.input.lumi.hasByok,
    isPhoneLayout: input.dashboard.snapshot.deviceProfile.deviceKind === 'se' || input.dashboard.snapshot.deviceProfile.deviceKind === 'mobile',
    isExpanded: input.dashboard.input.lumi.isPanelOpen,
    actionHandlers: input.actionHandlers,
    hoverSource: input.hoverSource,
    hoveredSystemSummary: input.hoveredSystemSummary,
  };
}
