import type { DashboardCourseRecord } from './types';

export interface DashboardSystemGroup<TCourse = DashboardCourseRecord> {
  id: string;
  index: number;
  title: string;
  courses: TCourse[];
}

export interface GroupDashboardSystemsInput<TCourse = DashboardCourseRecord> {
  courses: TCourse[];
  groupSize?: number;
}

const DEFAULT_SYSTEM_GROUP_SIZE = 5;

export function groupDashboardSystems<TCourse = DashboardCourseRecord>({
  courses,
  groupSize = DEFAULT_SYSTEM_GROUP_SIZE,
}: GroupDashboardSystemsInput<TCourse>): Array<DashboardSystemGroup<TCourse>> {
  if (groupSize <= 0) {
    return [];
  }

  const systems: Array<DashboardSystemGroup<TCourse>> = [];

  for (let index = 0; index < courses.length; index += groupSize) {
    const systemIndex = Math.floor(index / groupSize) + 1;
    systems.push({
      id: `system-${systemIndex}`,
      index: systemIndex,
      title: `System ${systemIndex}`,
      courses: courses.slice(index, index + groupSize),
    });
  }

  return systems;
}
