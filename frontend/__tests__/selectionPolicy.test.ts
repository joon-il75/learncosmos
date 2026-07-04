import { describe, it, expect } from 'vitest';
import { resolveDashboardSelectedCourseId } from '@/lib/world-ui-engine/selectionPolicy';
import type { DashboardCourseRecord } from '@/lib/world-ui-engine/types';

function makeCourse(overrides: Partial<DashboardCourseRecord> = {}): DashboardCourseRecord {
  return {
    id: 'c1',
    title: '테스트 코스',
    status: 'draft',
    draftStatus: 'draft',
    progress: null,
    lessonCount: 0,
    levelCount: 0,
    plannedLessonCount: 0,
    plannedLevelCount: 0,
    completedLessonCount: 0,
    sourceQuery: '',
    updatedAt: '2026-04-07T00:00:00Z',
    destination: '/dashboard/course-drafts/c1',
    ...overrides,
  };
}

describe('resolveDashboardSelectedCourseId', () => {
  it('페이지가 비어 있으면 null을 반환한다', () => {
    expect(
      resolveDashboardSelectedCourseId({
        pagedCourses: [],
        selectedCourseId: 'c1',
      }),
    ).toBeNull();
  });

  it('선택된 코스가 없으면 첫 코스를 선택한다', () => {
    const pagedCourses = [makeCourse({ id: 'c1' }), makeCourse({ id: 'c2' })];
    expect(
      resolveDashboardSelectedCourseId({
        pagedCourses,
        selectedCourseId: null,
      }),
    ).toBe('c1');
  });

  it('선택된 코스가 현재 페이지에 남아 있으면 유지한다', () => {
    const pagedCourses = [makeCourse({ id: 'c1' }), makeCourse({ id: 'c2' })];
    expect(
      resolveDashboardSelectedCourseId({
        pagedCourses,
        selectedCourseId: 'c2',
      }),
    ).toBe('c2');
  });

  it('선택된 코스가 현재 페이지에서 사라지면 첫 코스로 fallback한다', () => {
    const pagedCourses = [makeCourse({ id: 'c3' }), makeCourse({ id: 'c4' })];
    expect(
      resolveDashboardSelectedCourseId({
        pagedCourses,
        selectedCourseId: 'c1',
      }),
    ).toBe('c3');
  });
});
