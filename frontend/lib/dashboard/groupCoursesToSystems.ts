import type { DashboardCourseRecord } from '@/lib/world-ui-engine/types';

export interface StarSystemGroup {
  id: string;
  index: number;
  title: string;
  courses: DashboardCourseRecord[];
  statusCount: {
    draft: number;
    ready: number;
    learning: number;
    completed: number;
    inactive: number;
  };
  summaryTitle: string;
  updatedAt: string | null;
}

function buildStatusCount(courses: DashboardCourseRecord[]) {
  return courses.reduce(
    (counts, course) => {
      if (course.isInactive) {
        counts.inactive += 1;
      }
      counts[course.status] += 1;
      return counts;
    },
    { draft: 0, ready: 0, learning: 0, completed: 0, inactive: 0 },
  );
}

function getLatestUpdatedAt(courses: DashboardCourseRecord[]) {
  if (courses.length === 0) return null;
  return courses.reduce<string>(
    (latest, course) => (new Date(latest).getTime() >= new Date(course.updatedAt).getTime() ? latest : course.updatedAt),
    courses[0]!.updatedAt,
  );
}

export function chunkCoursesToSystems(
  courses: DashboardCourseRecord[],
  groupSize: number,
  startIndex = 1,
): StarSystemGroup[] {
  if (groupSize <= 0) {
    return [];
  }

  const systems: StarSystemGroup[] = [];

  for (let index = 0; index < courses.length; index += groupSize) {
    const systemIndex = startIndex + Math.floor(index / groupSize);
    const systemCourses = courses.slice(index, index + groupSize);
    const title = `제 ${systemIndex} 항성계`;
    systems.push({
      id: `system-${systemIndex}`,
      index: systemIndex,
      title,
      courses: systemCourses,
      statusCount: buildStatusCount(systemCourses),
      summaryTitle: title,
      updatedAt: getLatestUpdatedAt(systemCourses),
    });
  }

  return systems;
}
