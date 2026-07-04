import type { DashboardSortKey } from '@/lib/world-ui-engine/types';
import type {
  UserInfo,
  DraftLessonSummary,
  DraftLevelSummary,
  CourseItem,
  DraftLessonTreeSummary,
  PlanetLessonTreeSummary,
} from './_types';

export function hasDraftLessonPlan(lesson: DraftLessonSummary): boolean {
  const objective = lesson.lesson.objective?.trim();
  const summary = lesson.lesson.summary?.trim();
  const selectedResourceCount = (lesson.resources ?? []).filter(
    (r) => r.selection_state !== 'rejected',
  ).length;
  return Boolean(objective || summary || selectedResourceCount > 0);
}

export function hasDraftLevelPlan(level: DraftLevelSummary): boolean {
  return level.lessons.length > 0 && level.lessons.every(hasDraftLessonPlan);
}

export function hasDraftLessonTreePlan(lessonTree: DraftLessonTreeSummary): boolean {
  const objective = lessonTree.lesson.objective?.trim();
  const summary = lessonTree.lesson.summary?.trim();
  const selectedPointCount = (lessonTree.points ?? []).filter(
    (point) => point.point.selection_state !== 'rejected',
  ).length;
  return Boolean(objective || summary || selectedPointCount > 0);
}

export function countDraftPlanFromLessonTrees(lessons: DraftLessonTreeSummary[]): {
  plannedLevelCount: number;
  plannedLessonCount: number;
} {
  return lessons.reduce(
    (summary, lessonTree) => {
      const subLessons = lessonTree.sub_lessons ?? [];
      const plannedSubLessons = subLessons.filter(hasDraftLessonTreePlan).length;
      return {
        plannedLevelCount:
          summary.plannedLevelCount +
          (subLessons.length > 0 && plannedSubLessons === subLessons.length ? 1 : 0),
        plannedLessonCount: summary.plannedLessonCount + plannedSubLessons,
      };
    },
    { plannedLevelCount: 0, plannedLessonCount: 0 },
  );
}

export function isLearningLessonCompleted(lesson: DraftLessonSummary): boolean {
  const summary = lesson.lesson.summary?.trim() ?? '';
  const hasLessonMarker = summary.includes('[학습완료]');
  const hasResourceMarker = (lesson.resources ?? []).some((r) =>
    (r.description?.trim() ?? '').includes('[학습완료]'),
  );
  return hasLessonMarker && hasResourceMarker;
}

export function isLearningLessonTreeCompleted(lessonTree: PlanetLessonTreeSummary): boolean {
  if (lessonTree.lesson.status === 'inactive') return false;
  const activePoints = (lessonTree.points ?? []).filter((point) => {
    const itemStatus = point.point.item_status ?? point.point.explorer_status;
    return itemStatus !== 'inactive';
  });
  return activePoints.length > 0 && activePoints.every((point) => point.point.status === 'completed');
}

export function countCompletedLessonsFromLessonTrees(
  lessons: PlanetLessonTreeSummary[],
): number {
  const countTree = (lessonTree: PlanetLessonTreeSummary): number => {
    const ownCount = isLearningLessonTreeCompleted(lessonTree) ? 1 : 0;
    return ownCount + (lessonTree.sub_lessons ?? []).reduce((count, subLesson) => count + countTree(subLesson), 0);
  };
  return lessons.reduce((count, lessonTree) => count + countTree(lessonTree), 0);
}

export function countActiveLessonsFromLessonTrees(lessons: PlanetLessonTreeSummary[]): number {
  const countTree = (lessonTree: PlanetLessonTreeSummary): number => {
    if (lessonTree.lesson.status === 'inactive') return 0;
    return 1 + (lessonTree.sub_lessons ?? []).reduce((count, subLesson) => count + countTree(subLesson), 0);
  };
  return lessons.reduce((count, lessonTree) => count + countTree(lessonTree), 0);
}

export function getDisplayName(user: UserInfo | null): string {
  if (!user) return '학습자';
  return user.nickname || user.display_id || user.email || '학습자';
}

export function formatUpdatedAt(iso: string): string {
  return new Date(iso).toLocaleDateString('ko-KR', { month: 'short', day: 'numeric' });
}

export function isDashboardSortKey(value: string | null): value is DashboardSortKey {
  return (
    value === 'updated_desc' ||
    value === 'updated_asc' ||
    value === 'title_asc' ||
    value === 'title_desc'
  );
}

export function findSystemIdForCourse(
  courseId: string | null,
  chunkedSystems: Array<{ id: string; courses: Array<{ id: string; draftId?: string | null }> }>,
): string | null {
  if (!courseId) return null;
  const containing = chunkedSystems.find((sys) =>
    sys.courses.some((c) => c.id === courseId || c.draftId === courseId),
  );
  return containing?.id ?? null;
}
