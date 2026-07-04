import type { PlanetRouteKind, PointWorkTab } from '../pointPageTypes';

export type PointLearningStage =
  | 'point'
  | 'goal'
  | 'skill_discovery'
  | 'goal_iteration'
  | 'final_result'
  | 'complete';

export type PointLearningStageSlug =
  | 'observe'
  | 'goal'
  | 'skills'
  | 'repeat'
  | 'final'
  | 'complete';

export type GoalObjectType = 'creative' | 'conceptual' | 'physical' | 'coding' | 'other';

export type RepeatRecordLens =
  | 'attempt'
  | 'evidence'
  | 'known'
  | 'question'
  | 'answer'
  | 'works_well'
  | 'needs_practice'
  | 'resource';

export const pointLearningStages: PointLearningStage[] = [
  'point',
  'goal',
  'skill_discovery',
  'goal_iteration',
  'final_result',
  'complete',
];

export const pointLearningStageSlugs: PointLearningStageSlug[] = [
  'observe',
  'goal',
  'skills',
  'repeat',
  'final',
  'complete',
];

export const pointLearningStageToSlug: Record<PointLearningStage, PointLearningStageSlug> = {
  point: 'observe',
  goal: 'goal',
  skill_discovery: 'skills',
  goal_iteration: 'repeat',
  final_result: 'final',
  complete: 'complete',
};

export const pointLearningSlugToStage: Record<PointLearningStageSlug, PointLearningStage> = {
  observe: 'point',
  goal: 'goal',
  skills: 'skill_discovery',
  repeat: 'goal_iteration',
  final: 'final_result',
  complete: 'complete',
};

export function pointLearningStageFromSlug(slug?: string | null): PointLearningStage {
  if (!slug) return 'point';
  return pointLearningStageSlugs.includes(slug as PointLearningStageSlug)
    ? pointLearningSlugToStage[slug as PointLearningStageSlug]
    : 'point';
}

export function getPointStageHref(
  routeKind: PlanetRouteKind,
  planetID: string,
  pointID: string,
  stage: PointLearningStage,
): string {
  return `/dashboard/planets/${routeKind}/${planetID}/points/${pointID}/${pointLearningStageToSlug[stage]}`;
}

export const goalObjectTypes: GoalObjectType[] = ['creative', 'conceptual', 'physical', 'coding', 'other'];

export const repeatRecordLensToWorkTab: Record<RepeatRecordLens, PointWorkTab> = {
  attempt: 'practice',
  evidence: 'artifacts',
  known: 'notes',
  question: 'questions',
  answer: 'questions',
  works_well: 'practice',
  needs_practice: 'practice',
  resource: 'attachments',
};

export function getNextPointLearningStage(stage: PointLearningStage): PointLearningStage {
  const index = pointLearningStages.indexOf(stage);
  return pointLearningStages[Math.min(pointLearningStages.length - 1, index + 1)] ?? 'point';
}

export function getPreviousPointLearningStage(stage: PointLearningStage): PointLearningStage {
  const index = pointLearningStages.indexOf(stage);
  return pointLearningStages[Math.max(0, index - 1)] ?? 'point';
}
