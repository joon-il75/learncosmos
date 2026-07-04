import { buildDashboardPlanetState } from './planetStateEngine';
import type {
  DashboardCourseCollection,
  DashboardCourseRecord,
  DashboardDeviceProfile,
  DashboardPlanetCardViewModel,
  DashboardSelectionDerivedState,
  DashboardStatusCounts,
  DashboardStatusFilter,
  DashboardSystemViewModel,
  DashboardUIEngineSnapshot,
  DashboardMode,
  DashboardViewModelInput,
} from './types';
import { toCanonicalPlanetStatus } from './types';
import { getPageSizeForMode, type DashboardModeContract } from '@/lib/dashboard/dashboardState';
import { chunkCoursesToSystems } from '@/lib/dashboard/groupCoursesToSystems';
import { filterCourses, paginateCourses, sortCourses } from '@/lib/dashboard/coursePipeline';

function buildStatusCounts(courses: DashboardCourseRecord[]): DashboardStatusCounts {
  return courses.reduce<DashboardStatusCounts>(
    (counts, course) => {
      counts[course.status] += 1;
      return counts;
    },
    { draft: 0, ready: 0, learning: 0, completed: 0 },
  );
}

function buildCourseCollection(input: DashboardViewModelInput): DashboardCourseCollection {
  const normalizedMode: DashboardModeContract = input.mode === 'star-system' ? 'star-system' : 'galaxy';
  const pageSize = input.pageSize ?? getPageSizeForMode(normalizedMode);
  const filteredCourses = filterCourses(input.courses, input.activeStatusFilter, input.searchQuery);
  const sortedCourses = sortCourses(filteredCourses, input.sortKey);
  const paginated = paginateCourses(sortedCourses, input.page, pageSize);

  return {
    allCourses: input.courses,
    filteredCourses: paginated.filteredCourses,
    pagedCourses: paginated.pagedCourses,
    totalPages: paginated.totalPages,
    currentPage: paginated.currentPage,
  };
}

function buildSelectionState(
  courses: DashboardCourseRecord[],
  selectedCourseId: string | null,
  hoveredCourseId: string | null,
): DashboardSelectionDerivedState {
  return {
    selectedCourse: courses.find((course) => course.id === selectedCourseId) ?? null,
    hoveredCourse: courses.find((course) => course.id === hoveredCourseId) ?? null,
  };
}

function formatUpdatedAtLabel(iso: string): string {
  return new Date(iso).toLocaleDateString('ko-KR', {
    month: 'short',
    day: 'numeric',
  });
}


function buildPlanetCardViewModel(
  course: DashboardCourseRecord,
  selectedCourseId: string | null,
  hoveredCourseId: string | null,
): DashboardPlanetCardViewModel {
  const planetState = buildDashboardPlanetState(course);

  return {
    id: course.id,
    title: course.title,
    originalTopic: course.sourceQuery,
    lessonCount: course.lessonCount,
    updatedAtLabel: formatUpdatedAtLabel(course.updatedAt),
    canonicalStatus: planetState.canonicalStatus,
    progressLabel: course.progress != null ? `${Math.round(course.progress * 100)}% 진행` : undefined,
    isSelected: selectedCourseId === course.id,
    isHovered: hoveredCourseId === course.id,
    planetState,
  };
}

function buildSystems(
  pagedCourses: DashboardCourseRecord[],
  selectedSystemId: string | null,
  selectedCourseId: string | null,
  hoveredCourseId: string | null,
): DashboardSystemViewModel[] {
  const systems = chunkCoursesToSystems(pagedCourses, 5);

  return systems.map((system) => ({
    id: system.id,
    index: system.index,
    title: system.title,
    courseCount: system.courses.length,
    isSelected: system.id === selectedSystemId,
    statusCount: system.statusCount,
    summaryTitle: system.summaryTitle,
    updatedAt: system.updatedAt,
    courses: system.courses.map((course) =>
      buildPlanetCardViewModel(course, selectedCourseId, hoveredCourseId),
    ),
  }));
}

function buildVisibleLayers(deviceKind: DashboardDeviceProfile['deviceKind']) {
  if (deviceKind === 'se') {
    return {
      core: true,
      navigation: true,
      world: 'minimal',
      effect: 'off',
    } as const;
  }
  if (deviceKind === 'mobile') {
    return {
      core: true,
      navigation: true,
      world: 'partial',
      effect: 'off',
    } as const;
  }
  if (deviceKind === 'tablet') {
    return {
      core: true,
      navigation: true,
      world: 'full',
      effect: 'reduced',
    } as const;
  }
  return {
    core: true,
    navigation: true,
    world: 'full',
    effect: deviceKind === 'wide-desktop' ? 'full' : 'reduced',
  } as const;
}

function buildDashboardMode(deviceKind: DashboardDeviceProfile['deviceKind']) {
  if (deviceKind === 'se') return 'planet-list' as const;
  if (deviceKind === 'mobile' || deviceKind === 'tablet') return 'system-selector' as const;
  return 'galaxy' as const;
}

function buildSystemLayout(deviceKind: DashboardDeviceProfile['deviceKind']) {
  if (deviceKind === 'se') return 'stack' as const;
  if (deviceKind === 'mobile') return 'grid' as const;
  if (deviceKind === 'tablet') return 'grid' as const;
  return 'orbit' as const;
}

function buildLumiMode(deviceKind: DashboardDeviceProfile['deviceKind']) {
  if (deviceKind === 'se') return 'button-only' as const;
  if (deviceKind === 'mobile') return 'bottom-sheet' as const;
  return 'docked' as const;
}

function buildEmptyState(filteredCourseCount: number) {
  if (filteredCourseCount > 0) return null;

  return {
    title: '아직 배움 우주가 비어 있어요.',
    description: '새로운 행성을 만들면 이곳에서 나의 배움 우주가 확장됩니다.',
    ctaLabel: '새 행성 시작하기',
  };
}

export function buildDashboardViewModel(
  input: DashboardViewModelInput,
): DashboardUIEngineSnapshot {
  const normalizedMode = input.mode === 'star-system' ? 'star-system' : 'galaxy';
  const sanitizedInput: DashboardViewModelInput = {
    ...input,
    mode: normalizedMode,
  };
  const courseCollection = buildCourseCollection(sanitizedInput);
  const selectedSystemId = input.selectedSystemId ?? `system-${courseCollection.currentPage}`;

  return {
    courses: courseCollection,
    counts: buildStatusCounts(input.courses),
    selection: buildSelectionState(
      courseCollection.pagedCourses,
      input.selectedCourseId,
      input.hoveredCourseId,
    ),
    deviceProfile: input.deviceProfile,
    visibleLayers: buildVisibleLayers(input.deviceProfile.deviceKind),
    dashboardMode: normalizedMode,
    systemLayout: buildSystemLayout(input.deviceProfile.deviceKind),
    lumiMode: buildLumiMode(input.deviceProfile.deviceKind),
    systems: buildSystems(
      courseCollection.pagedCourses,
      selectedSystemId,
      input.selectedCourseId,
      input.hoveredCourseId,
    ),
    selectedSystemId,
    hoveredSystemId: input.hoveredSystemId,
    filteredCourseCount: courseCollection.filteredCourses.length,
    totalCourseCount: input.courses.length,
    emptyState: buildEmptyState(courseCollection.filteredCourses.length),
    cta: {
      title: '새 학습탐험 경로 만들기',
      placeholder: '예: 기타 코드 3개월 안에 마스터하기',
      submitLabel: '학습탐험 시작하기',
    },
  };
}
