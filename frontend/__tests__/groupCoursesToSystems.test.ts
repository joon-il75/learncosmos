import { describe, it, expect } from 'vitest';
import { chunkCoursesToSystems } from '@/lib/dashboard/groupCoursesToSystems';
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

describe('chunkCoursesToSystems', () => {
  it('빈 배열은 빈 결과를 반환한다', () => {
    expect(chunkCoursesToSystems([], 5)).toEqual([]);
  });

  it('groupSize가 0 이하이면 빈 결과를 반환한다', () => {
    const courses = [makeCourse()];
    expect(chunkCoursesToSystems(courses, 0)).toEqual([]);
    expect(chunkCoursesToSystems(courses, -1)).toEqual([]);
  });

  it('코스 5개를 groupSize 5로 나누면 항성계 1개가 된다', () => {
    const courses = Array.from({ length: 5 }, (_, i) =>
      makeCourse({ id: `c${i + 1}` }),
    );
    const result = chunkCoursesToSystems(courses, 5);
    expect(result).toHaveLength(1);
    expect(result[0]!.index).toBe(1);
    expect(result[0]!.courses).toHaveLength(5);
    expect(result[0]!.id).toBe('system-1');
  });

  it('코스 6개를 groupSize 5로 나누면 항성계 2개가 된다', () => {
    const courses = Array.from({ length: 6 }, (_, i) =>
      makeCourse({ id: `c${i + 1}` }),
    );
    const result = chunkCoursesToSystems(courses, 5);
    expect(result).toHaveLength(2);
    expect(result[0]!.courses).toHaveLength(5);
    expect(result[1]!.courses).toHaveLength(1);
    expect(result[1]!.index).toBe(2);
  });

  it('statusCount가 올바르게 집계된다', () => {
    const courses = [
      makeCourse({ id: 'c1', status: 'draft' }),
      makeCourse({ id: 'c2', status: 'draft' }),
      makeCourse({ id: 'c3', status: 'learning' }),
      makeCourse({ id: 'c4', status: 'completed' }),
      makeCourse({ id: 'c5', status: 'ready' }),
    ];
    const result = chunkCoursesToSystems(courses, 5);
    expect(result[0]!.statusCount).toEqual({
      draft: 2,
      ready: 1,
      learning: 1,
      completed: 1,
      inactive: 0,
    });
  });

  it('updatedAt은 가장 최신 코스의 날짜를 반환한다', () => {
    const courses = [
      makeCourse({ id: 'c1', updatedAt: '2026-04-01T00:00:00Z' }),
      makeCourse({ id: 'c2', updatedAt: '2026-04-07T00:00:00Z' }),
      makeCourse({ id: 'c3', updatedAt: '2026-04-03T00:00:00Z' }),
    ];
    const result = chunkCoursesToSystems(courses, 5);
    expect(result[0]!.updatedAt).toBe('2026-04-07T00:00:00Z');
  });

  it('항성계 id와 title이 index 기반으로 생성된다', () => {
    const courses = Array.from({ length: 10 }, (_, i) =>
      makeCourse({ id: `c${i + 1}` }),
    );
    const result = chunkCoursesToSystems(courses, 5);
    expect(result[0]!.id).toBe('system-1');
    expect(result[0]!.title).toBe('제 1 항성계');
    expect(result[1]!.id).toBe('system-2');
    expect(result[1]!.title).toBe('제 2 항성계');
  });
});
