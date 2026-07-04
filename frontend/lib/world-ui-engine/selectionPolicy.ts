import type { DashboardCourseRecord } from './types';

export interface DashboardSelectionPolicyInput {
  pagedCourses: DashboardCourseRecord[];
  selectedCourseId: string | null;
}

export function resolveDashboardSelectedCourseId(
  input: DashboardSelectionPolicyInput,
): string | null {
  if (input.pagedCourses.length === 0) {
    return null;
  }

  if (!input.selectedCourseId) {
    return input.pagedCourses[0]?.id ?? null;
  }

  const selectedStillExists = input.pagedCourses.some(
    (course) => course.id === input.selectedCourseId,
  );

  if (selectedStillExists) {
    return input.selectedCourseId;
  }

  return input.pagedCourses[0]?.id ?? null;
}
