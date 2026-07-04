import { buildDashboardViewModel } from './dashboardViewModel';
import { buildDashboardDeviceProfile } from './deviceProfile';
import type {
  DashboardAdapterResult,
  DashboardUIEngineInput,
} from './types';

export function buildDashboardAdapter(
  input: DashboardUIEngineInput,
): DashboardAdapterResult {
  const deviceProfile = buildDashboardDeviceProfile({
    width: input.device.viewportWidth,
    height: input.device.viewportHeight,
  });

  const snapshot = buildDashboardViewModel({
    courses: input.data.courses,
    searchQuery: input.filter.query,
    page: input.filter.page,
    activeStatusFilter: input.filter.activeStatusFilter,
    sortKey: input.filter.sortKey,
    selectedCourseId: input.selection.selectedCourseId,
    hoveredCourseId: input.selection.hoveredCourseId,
    mode: input.mode,
    selectedSystemId: input.selectedSystemId,
    hoveredSystemId: input.hoveredSystemId,
    deviceProfile,
  });

  return {
    input,
    snapshot: {
      ...snapshot,
      deviceProfile,
    },
  };
}
