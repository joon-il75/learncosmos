'use client';

import Link from 'next/link';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../pointPageUtils';
import {
  getPlanetDetailHref,
} from '../pointPageUtils';
import {
  placeholderCardStyle,
  pointShellTitleStyle,
  runtimeButtonRowStyle,
  secondaryButtonStyle,
} from '../pointPageStyles';
import type { UsePointLearningResult } from './usePointLearning';

type PointDetail = NonNullable<UsePointLearningResult['pointDetail']>;
type Planet = NonNullable<UsePointLearningResult['planet']>;
type CompletionReadiness = UsePointLearningResult['pointCompletionReadiness'];

export function CompletionReadOnlyNotice({
  copy,
  themeTokens,
}: {
  copy: PointLearningCopy['completionStatus'];
  themeTokens: PointPageThemeTokens;
}) {
  return (
    <div style={{ ...placeholderCardStyle, background: themeTokens.placeholderBackground, borderColor: themeTokens.placeholderBorder }}>
      <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.readOnlyTitle}</strong>
      <div style={{ fontSize: '15px', lineHeight: 1.65, color: themeTokens.description }}>
        {copy.readOnlyDescription}
      </div>
    </div>
  );
}

export function CompletionAchievementCard({
  copy,
  themeTokens,
  routeKind,
  planetID,
  planet,
  pointDetail,
  journalObservation,
  journalReflection,
  journalNextStep,
  pointQuestionsLength,
  answeredQuestionCount,
  practiceLogDraftsLength,
  artifactDraftsLength,
}: {
  copy: PointLearningCopy['completionStatus'];
  themeTokens: PointPageThemeTokens;
  routeKind: UsePointLearningResult['routeKind'];
  planetID: UsePointLearningResult['planetID'];
  planet: Planet;
  pointDetail: PointDetail;
  journalObservation: string;
  journalReflection: string;
  journalNextStep: string;
  pointQuestionsLength: number;
  answeredQuestionCount: number;
  practiceLogDraftsLength: number;
  artifactDraftsLength: number;
}) {
  const noteCount = journalObservation.trim() || journalReflection.trim() || journalNextStep.trim() ? 1 : 0;

  return (
    <div style={{ ...placeholderCardStyle, background: themeTokens.noticeBackground, borderColor: themeTokens.noticeBorder }}>
      <div style={{ display: 'grid', gap: '12px' }}>
        <div>
          <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.achievementTitle}</strong>
          <div style={{ marginTop: '6px', fontSize: '15px', lineHeight: 1.65, color: themeTokens.description }}>{copy.achievementDescription}</div>
        </div>
        <div style={{ display: 'grid', gap: '7px', fontSize: '15px', lineHeight: 1.6, color: themeTokens.description }}>
          <div>{copy.contentChecked}</div>
          <div>{copy.noteCount(noteCount)}</div>
          <div>{copy.questionCount(pointQuestionsLength, answeredQuestionCount)}</div>
          <div>{copy.practiceCount(practiceLogDraftsLength)}</div>
          <div>{copy.artifactCount(artifactDraftsLength)}</div>
          <div>{copy.selfEvaluationState(Boolean(pointDetail.self_evaluation))}</div>
        </div>
        <div style={runtimeButtonRowStyle}>
          <Link href={getPlanetDetailHref(routeKind, planetID)} style={{ ...secondaryButtonStyle, background: themeTokens.secondaryButtonBackground, borderColor: themeTokens.secondaryButtonBorder, color: themeTokens.buttonText }}>{copy.backToDiary}</Link>
        </div>
      </div>
    </div>
  );
}

export function CompletionChecklistCard({
  copy,
  themeTokens,
  pointDetail,
  readiness,
  isSharedRoute,
  isResearchMaterialConfirmed,
}: {
  copy: PointLearningCopy['completionStatus'];
  themeTokens: PointPageThemeTokens;
  pointDetail: PointDetail;
  readiness: CompletionReadiness;
  isSharedRoute: boolean;
  isResearchMaterialConfirmed: boolean;
}) {
  return (
    <div style={{ ...placeholderCardStyle, background: themeTokens.placeholderBackground, borderColor: themeTokens.placeholderBorder }}>
      <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{isSharedRoute ? copy.sharedStatusTitle : copy.readyStatusTitle}</strong>
      <div style={{ display: 'grid', gap: '8px', color: themeTokens.description, fontSize: '15px', lineHeight: 1.6 }}>
        {pointDetail.point.point_type === 'research' ? (
          <div>{copy.researchMaterialConfirmed(isResearchMaterialConfirmed)}</div>
        ) : null}
        <div>{copy.answeredQuestionReady(readiness.answeredQuestionCount > 0, readiness.answeredQuestionCount)}</div>
        <div>{copy.selfEvaluationSaved(readiness.hasSelfEvaluation)}</div>
        <div>{copy.selfEvaluationQuality(readiness.hasSelfEvaluationQuality)}</div>
        <div>{copy.goalConnection(readiness.hasGoalConnection)}</div>
      </div>
      {!readiness.isReady ? (
        <div style={{ fontSize: '14px', lineHeight: 1.6, color: themeTokens.mutedText }}>
          {isSharedRoute
            ? copy.sharedNotReady
            : pointDetail.point.point_type === 'research' && !isResearchMaterialConfirmed
            ? copy.researchNotReady
            : copy.learningNotReady}
        </div>
      ) : (
        <div style={{ fontSize: '14px', lineHeight: 1.6, color: themeTokens.mutedText }}>
          {isSharedRoute ? copy.sharedReady : copy.learningReady}
        </div>
      )}
    </div>
  );
}
