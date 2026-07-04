import type { DashboardCourseRecord } from './types';

export interface DashboardSelectionLaunchPlan {
  courseId: string;
  destination: string;
  delayMs: number;
}

const DEFAULT_SELECTION_LAUNCH_DELAY_MS = 720;

export function buildDashboardSelectionLaunchPlan(
  course: DashboardCourseRecord,
  delayMs = DEFAULT_SELECTION_LAUNCH_DELAY_MS,
): DashboardSelectionLaunchPlan {
  return {
    courseId: course.id,
    destination: course.destination,
    delayMs,
  };
}
