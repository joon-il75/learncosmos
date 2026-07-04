import type {
  DashboardCourseRecord,
  DashboardStatusFilter,
  DashboardSortKey,
} from '@/lib/world-ui-engine/types';

function normalizeSearchQuery(query: string) {
  return query.trim().toLowerCase();
}

export function filterCourses(
  courses: DashboardCourseRecord[],
  statusFilter: DashboardStatusFilter,
  searchQuery: string,
): DashboardCourseRecord[] {
  const normalizedQuery = normalizeSearchQuery(searchQuery);

  return courses.filter((course) => {
    if (statusFilter !== 'all' && course.status !== statusFilter) {
      return false;
    }

    if (!normalizedQuery) return true;

    return (
      course.title.toLowerCase().includes(normalizedQuery) ||
      course.sourceQuery.toLowerCase().includes(normalizedQuery)
    );
  });
}

function compareStringsAsc(a: string, b: string) {
  if (a < b) return -1;
  if (a > b) return 1;
  return 0;
}

export function sortCourses(courses: DashboardCourseRecord[], sortKey: DashboardSortKey) {
  const sorted = [...courses];

  sorted.sort((a, b) => {
    if (sortKey === 'updated_desc') {
      return new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime();
    }

    if (sortKey === 'updated_asc') {
      return new Date(a.updatedAt).getTime() - new Date(b.updatedAt).getTime();
    }

    if (sortKey === 'title_asc') {
      return compareStringsAsc(a.title, b.title);
    }

    if (sortKey === 'title_desc') {
      return -compareStringsAsc(a.title, b.title);
    }

    return 0;
  });

  return sorted;
}

export interface PaginatedCourseCollection {
  pagedCourses: DashboardCourseRecord[];
  totalPages: number;
  currentPage: number;
  filteredCourses: DashboardCourseRecord[];
}

export function paginateCourses(
  courses: DashboardCourseRecord[],
  requestedPage: number,
  pageSize: number,
): PaginatedCourseCollection {
  const normalizedPage = Math.max(1, Math.floor(requestedPage));
  const totalPages = Math.max(1, Math.ceil(courses.length / pageSize));
  const currentPage = Math.min(normalizedPage, totalPages);
  const start = (currentPage - 1) * pageSize;
  const pagedCourses = courses.slice(start, start + pageSize);

  return {
    filteredCourses: courses,
    pagedCourses,
    totalPages,
    currentPage,
  };
}
