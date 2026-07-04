import { describe, it, expect } from 'vitest';
import {
  filterCourses,
  paginateCourses,
  sortCourses,
} from '@/lib/dashboard/coursePipeline';
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

describe('filterCourses', () => {
  const courses = [
    makeCourse({
      id: 'guitar',
      title: '기타 코드 입문',
      status: 'ready',
      sourceQuery: '기타 독학',
    }),
    makeCourse({
      id: 'painting',
      title: '수채화 기초',
      status: 'learning',
      sourceQuery: '풍경화',
    }),
    makeCourse({
      id: 'running',
      title: '러닝 루틴 만들기',
      status: 'completed',
      sourceQuery: '마라톤 준비',
    }),
  ];

  it('statusFilter가 all이면 검색어만 적용한다', () => {
    const result = filterCourses(courses, 'all', '  기타  ');
    expect(result.map((course) => course.id)).toEqual(['guitar']);
  });

  it('title뿐 아니라 sourceQuery로도 검색한다', () => {
    const result = filterCourses(courses, 'all', '마라톤');
    expect(result.map((course) => course.id)).toEqual(['running']);
  });

  it('statusFilter와 검색어를 함께 적용한다', () => {
    const result = filterCourses(courses, 'learning', '수채화');
    expect(result.map((course) => course.id)).toEqual(['painting']);
  });

  it('검색어가 비어 있으면 상태 필터만 적용한다', () => {
    const result = filterCourses(courses, 'completed', '   ');
    expect(result.map((course) => course.id)).toEqual(['running']);
  });
});

describe('sortCourses', () => {
  const courses = [
    makeCourse({ id: 'b', title: '베이스', updatedAt: '2026-04-02T00:00:00Z' }),
    makeCourse({ id: 'a', title: '기타', updatedAt: '2026-04-03T00:00:00Z' }),
    makeCourse({ id: 'c', title: '첼로', updatedAt: '2026-04-01T00:00:00Z' }),
  ];

  it('updated_desc는 최신순으로 정렬한다', () => {
    expect(sortCourses(courses, 'updated_desc').map((course) => course.id)).toEqual(['a', 'b', 'c']);
  });

  it('updated_asc는 오래된 순으로 정렬한다', () => {
    expect(sortCourses(courses, 'updated_asc').map((course) => course.id)).toEqual(['c', 'b', 'a']);
  });

  it('title_asc는 제목 오름차순으로 정렬한다', () => {
    expect(sortCourses(courses, 'title_asc').map((course) => course.id)).toEqual(['a', 'b', 'c']);
  });

  it('title_desc는 제목 내림차순으로 정렬한다', () => {
    expect(sortCourses(courses, 'title_desc').map((course) => course.id)).toEqual(['c', 'b', 'a']);
  });

  it('원본 배열을 변경하지 않는다', () => {
    const original = courses.map((course) => course.id);
    sortCourses(courses, 'updated_desc');
    expect(courses.map((course) => course.id)).toEqual(original);
  });
});

describe('paginateCourses', () => {
  const courses = Array.from({ length: 12 }, (_, index) =>
    makeCourse({ id: `c${index + 1}` }),
  );

  it('페이지를 1보다 작게 요청하면 1로 보정한다', () => {
    const result = paginateCourses(courses, 0, 5);
    expect(result.currentPage).toBe(1);
    expect(result.pagedCourses.map((course) => course.id)).toEqual(['c1', 'c2', 'c3', 'c4', 'c5']);
  });

  it('페이지를 초과 요청하면 마지막 페이지로 보정한다', () => {
    const result = paginateCourses(courses, 99, 5);
    expect(result.currentPage).toBe(3);
    expect(result.totalPages).toBe(3);
    expect(result.pagedCourses.map((course) => course.id)).toEqual(['c11', 'c12']);
  });

  it('소수 페이지는 내림 처리한다', () => {
    const result = paginateCourses(courses, 2.8, 5);
    expect(result.currentPage).toBe(2);
    expect(result.pagedCourses.map((course) => course.id)).toEqual(['c6', 'c7', 'c8', 'c9', 'c10']);
  });

  it('빈 배열이어도 totalPages는 최소 1을 유지한다', () => {
    const result = paginateCourses([], 1, 5);
    expect(result.totalPages).toBe(1);
    expect(result.currentPage).toBe(1);
    expect(result.pagedCourses).toEqual([]);
    expect(result.filteredCourses).toEqual([]);
  });
});
