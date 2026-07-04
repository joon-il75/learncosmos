import type { PlanetMetaState, ResolvedPoint } from '../pointPageTypes';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';

export type LumiGuideMode = 'point' | 'learning-content' | 'learning-work' | 'self-evaluation';
export type LearningWorkGuideDetail = { title: string; message: string };

export const LUMI_REST_NOTICE_DELAY_MS = 25 * 60 * 1000;

export const pointToolbarOffsets = {
  collapsed: '166px',
  expanded: '304px',
} as const;

export function buildPointLumiGuideCopy({
  isExplorationPoint,
  learningWorkGuideDetail,
  copy,
}: {
  isExplorationPoint: boolean;
  learningWorkGuideDetail: LearningWorkGuideDetail;
  copy: PointLearningCopy['lumiGuide'];
}): Record<LumiGuideMode, { title: string; message: string }> {
  return {
    point: { title: copy.pointTitle(isExplorationPoint), message: copy.pointMessage(isExplorationPoint) },
    'learning-content': {
      title: copy.learningContentTitle,
      message: copy.learningContentMessage(isExplorationPoint),
    },
    'learning-work': learningWorkGuideDetail,
    'self-evaluation': {
      title: copy.selfEvaluationTitle,
      message: copy.selfEvaluationMessage,
    },
  };
}

export function getPointGoalContextText(planet: PlanetMetaState, fallback: string): string {
  return (
    planet.goal_context?.confirmed_goal?.trim() ||
    planet.goal_context?.learning_goal?.trim() ||
    planet.planet.title.trim() ||
    fallback
  );
}

export function getCurrentFlowGuideMessage({
  pointDetail,
  isExplorationPoint,
  isResearchMaterialConfirmed,
  hasSavedJournalNote,
  hasSelfEvaluationEntryRecord,
  hasSelfEvaluationQuality,
  pointQuestionsLength,
  answeredQuestionCount,
  isPointCompletionReady,
  copy,
}: {
  pointDetail: ResolvedPoint;
  isExplorationPoint: boolean;
  isResearchMaterialConfirmed: boolean;
  hasSavedJournalNote: boolean;
  hasSelfEvaluationEntryRecord: boolean;
  hasSelfEvaluationQuality: boolean;
  pointQuestionsLength: number;
  answeredQuestionCount: number;
  isPointCompletionReady: boolean;
  copy: PointLearningCopy['lumiGuide'];
}): string {
  if (pointDetail.point.status === 'completed') {
    return copy.completedFlow;
  }
  if (!isExplorationPoint && !isResearchMaterialConfirmed) {
    return copy.researchMaterialFlow;
  }
  if (!hasSavedJournalNote) {
    if (isExplorationPoint && pointDetail.point.external_url) {
      return copy.openExternalFlow;
    }
    return copy.firstNoteFlow;
  }
  if (hasSelfEvaluationEntryRecord && !hasSelfEvaluationQuality) {
    return copy.selfEvaluationReadyFlow;
  }
  if (pointQuestionsLength === 0) {
    return copy.noQuestionFlow;
  }
  if (answeredQuestionCount === 0) {
    return copy.noAnswerFlow;
  }
  if (isPointCompletionReady) {
    return copy.completionReadyFlow;
  }
  return copy.fillMissingFlow;
}
