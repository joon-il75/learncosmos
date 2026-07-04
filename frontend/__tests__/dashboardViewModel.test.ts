import { describe, it, expect } from 'vitest';
import { buildDashboardViewModel } from '@/lib/world-ui-engine/dashboardViewModel';
import type {
  DashboardCourseRecord,
  DashboardDeviceProfile,
  DashboardViewModelInput,
} from '@/lib/world-ui-engine/types';

function makeCourse(overrides: Partial<DashboardCourseRecord> = {}): DashboardCourseRecord {
  return {
    id: 'c1',
    title: '테스트 코스',
    status: 'draft',
    draftStatus: 'draft',
    progress: null,
    lessonCount: 4,
    levelCount: 1,
    plannedLessonCount: 2,
    plannedLevelCount: 1,
    completedLessonCount: 0,
    sourceQuery: '기본 주제',
    updatedAt: '2026-04-08T00:00:00Z',
    destination: '/dashboard/course-drafts/c1',
    ...overrides,
  };
}

function makeDeviceProfile(
  overrides: Partial<DashboardDeviceProfile> = {},
): DashboardDeviceProfile {
  return {
    deviceKind: 'desktop',
    viewportKind: 'default',
    useCompactMapLabels: false,
    mapBreakpoint: 'desktop',
    ...overrides,
  };
}

function makeInput(
  overrides: Partial<DashboardViewModelInput> = {},
): DashboardViewModelInput {
  return {
    courses: [],
    searchQuery: '',
    page: 1,
    activeStatusFilter: 'all',
    selectedCourseId: null,
    hoveredCourseId: null,
    mode: 'galaxy',
    selectedSystemId: null,
    hoveredSystemId: null,
    deviceProfile: makeDeviceProfile(),
    sortKey: 'updated_desc',
    ...overrides,
  };
}

describe('buildDashboardViewModel', () => {
  it('galaxy 모드는 기본 page size 25와 현재 페이지 기반 system id를 사용한다', () => {
    const courses = Array.from({ length: 30 }, (_, index) =>
      makeCourse({
        id: `c${index + 1}`,
        updatedAt: `2026-04-${String(index + 1).padStart(2, '0')}T00:00:00Z`,
      }),
    );

    const snapshot = buildDashboardViewModel(
      makeInput({
        courses,
        page: 2,
      }),
    );

    expect(snapshot.dashboardMode).toBe('galaxy');
    expect(snapshot.courses.currentPage).toBe(2);
    expect(snapshot.courses.totalPages).toBe(2);
    expect(snapshot.courses.pagedCourses).toHaveLength(5);
    expect(snapshot.selectedSystemId).toBe('system-2');
    expect(snapshot.systems).toHaveLength(1);
  });

  it('star-system 모드는 기본 page size 5를 사용한다', () => {
    const courses = Array.from({ length: 8 }, (_, index) =>
      makeCourse({
        id: `c${index + 1}`,
        updatedAt: `2026-04-${String(index + 1).padStart(2, '0')}T00:00:00Z`,
      }),
    );

    const snapshot = buildDashboardViewModel(
      makeInput({
        courses,
        mode: 'star-system',
        selectedSystemId: 'system-1',
      }),
    );

    expect(snapshot.dashboardMode).toBe('star-system');
    expect(snapshot.courses.pagedCourses).toHaveLength(5);
    expect(snapshot.courses.totalPages).toBe(2);
    expect(snapshot.systems[0]?.isSelected).toBe(true);
  });

  it('selectedCourse와 hoveredCourse를 현재 페이지 기준으로 파생한다', () => {
    const courses = [
      makeCourse({ id: 'draft-course', title: '초안 코스', status: 'draft' }),
      makeCourse({ id: 'learning-course', title: '진행 코스', status: 'learning', progress: 0.4 }),
    ];

    const snapshot = buildDashboardViewModel(
      makeInput({
        courses,
        selectedCourseId: 'draft-course',
        hoveredCourseId: 'learning-course',
      }),
    );

    expect(snapshot.selection.selectedCourse?.id).toBe('draft-course');
    expect(snapshot.selection.hoveredCourse?.id).toBe('learning-course');
    expect(snapshot.systems[0]?.courses[0]?.isSelected).toBe(true);
    expect(snapshot.systems[0]?.courses[1]?.isHovered).toBe(true);
  });

  it('코스 카드에 progressLabel을 생성한다', () => {
    const courses = [
      makeCourse({ id: 'draft', status: 'draft', title: '초안', progress: null }),
      makeCourse({ id: 'learning', status: 'learning', title: '진행', progress: 0.42 }),
    ];

    const snapshot = buildDashboardViewModel(
      makeInput({
        courses,
        sortKey: 'title_asc',
      }),
    );

    const cards = snapshot.systems[0]!.courses;
    expect(cards.find((card) => card.id === 'learning')?.progressLabel).toBe('42% 진행');
    expect(cards.find((card) => card.id === 'draft')?.progressLabel).toBeUndefined();
  });

  it('검색 결과가 없으면 emptyState를 노출한다', () => {
    const snapshot = buildDashboardViewModel(
      makeInput({
        courses: [makeCourse({ id: 'c1', title: '기타 코스' })],
        searchQuery: '수영',
      }),
    );

    expect(snapshot.filteredCourseCount).toBe(0);
    expect(snapshot.emptyState).toEqual({
      title: '아직 배움 우주가 비어 있어요.',
      description: '새로운 행성을 만들면 이곳에서 나의 배움 우주가 확장됩니다.',
      ctaLabel: '새 행성 시작하기',
    });
  });

  it('디바이스 종류에 따라 visibleLayers, layout, lumiMode를 조정한다', () => {
    const seSnapshot = buildDashboardViewModel(
      makeInput({
        courses: [makeCourse()],
        deviceProfile: makeDeviceProfile({
          deviceKind: 'se',
          mapBreakpoint: 'mobile',
          useCompactMapLabels: true,
        }),
      }),
    );

    const tabletSnapshot = buildDashboardViewModel(
      makeInput({
        courses: [makeCourse()],
        deviceProfile: makeDeviceProfile({
          deviceKind: 'tablet',
          mapBreakpoint: 'mobile',
          useCompactMapLabels: true,
        }),
      }),
    );

    expect(seSnapshot.visibleLayers).toEqual({
      core: true,
      navigation: true,
      world: 'minimal',
      effect: 'off',
    });
    expect(seSnapshot.systemLayout).toBe('stack');
    expect(seSnapshot.lumiMode).toBe('button-only');

    expect(tabletSnapshot.visibleLayers).toEqual({
      core: true,
      navigation: true,
      world: 'full',
      effect: 'reduced',
    });
    expect(tabletSnapshot.systemLayout).toBe('grid');
    expect(tabletSnapshot.lumiMode).toBe('docked');
  });
});
