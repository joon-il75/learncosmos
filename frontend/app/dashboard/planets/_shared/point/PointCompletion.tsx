'use client';

import { useState } from 'react';
import LumiModalShell, { lumiModalPrimaryButtonStyle, lumiModalSecondaryButtonStyle } from '@/components/common/LumiModalShell';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../pointPageUtils';
import type { UsePointLearningResult } from './usePointLearning';
import {
  CompletionAchievementCard,
  CompletionChecklistCard,
  CompletionReadOnlyNotice,
} from './PointCompletionStatusCards';
import { PointSelfEvaluationPanel } from './PointSelfEvaluationPanel';
import {
  sectionStyle,
  pointShellTitleStyle,
  noticeCardStyle,
  formGridStyle,
  runtimeButtonRowStyle, runtimePrimaryButtonStyle, runtimeButtonDisabledStyle,
} from '../pointPageStyles';


interface Props {
  copy: PointLearningCopy;
  learning: UsePointLearningResult;
  themeTokens: PointPageThemeTokens;
  isSharedRoute: boolean;
}

export default function PointCompletion({ copy, learning, themeTokens, isSharedRoute }: Props) {
  const {
    routeKind, planetID, planet, pointDetail,
    pointCompletionReadiness,
    isCompletionOpen,
    showCompletionAchievement,
    isSavingPointRuntime, handleSavePointRuntime,
    isResearchMaterialConfirmed,
    finalSelfEvalScore,
    isSavingSelfEvaluation,
    handleSaveSelfEvaluation,
    // summary
    answeredQuestionCount, practiceLogDrafts, artifactDrafts,
    journalObservation, journalReflection, journalNextStep,
  } = learning;
  const [isCompletionConfirmOpen, setIsCompletionConfirmOpen] = useState(false);
  const [isCompletionCancelConfirmOpen, setIsCompletionCancelConfirmOpen] = useState(false);
  const [isCompletionCancelExecuting, setIsCompletionCancelExecuting] = useState(false);

  if (!planet || !pointDetail) return null;
  if (!isCompletionOpen) return null;

  const isPointCompleted = pointDetail.point.status === 'completed';
  const canEditPoint = routeKind === 'learning' && !isPointCompleted;
  const hasFinalSelfEvalScore = finalSelfEvalScore.trim().length > 0;
  const isCompletionButtonDisabled = planet.planet.status !== 'learning'
    || isSavingPointRuntime
    || isSavingSelfEvaluation
    || isPointCompleted
    || !pointCompletionReadiness.isReady
    || !hasFinalSelfEvalScore;
  return (
    <section id="point-completion-section" style={{ ...sectionStyle, scrollMarginTop: '430px', background: themeTokens.sectionBackground, borderColor: themeTokens.sectionBorder }}>
      <div style={formGridStyle}>
        {isSharedRoute ? <CompletionReadOnlyNotice copy={copy.completionStatus} themeTokens={themeTokens} /> : null}
        {showCompletionAchievement || pointDetail.point.status === 'completed' ? (
          <CompletionAchievementCard
            copy={copy.completionStatus}
            themeTokens={themeTokens}
            routeKind={routeKind}
            planetID={planetID}
            planet={planet}
            pointDetail={pointDetail}
            journalObservation={journalObservation}
            journalReflection={journalReflection}
            journalNextStep={journalNextStep}
            pointQuestionsLength={learning.pointQuestions.length}
            answeredQuestionCount={answeredQuestionCount}
            practiceLogDraftsLength={practiceLogDrafts.length}
            artifactDraftsLength={artifactDrafts.length}
          />
        ) : null}

        <PointSelfEvaluationPanel
          copy={copy.selfEvaluation}
          learning={learning}
          themeTokens={themeTokens}
          canEditPoint={canEditPoint}
        />

        <div style={{ ...noticeCardStyle, background: themeTokens.noticeBackground, borderColor: themeTokens.noticeBorder, color: themeTokens.description }}>
          {copy.completion.encouragement}
        </div>

        <CompletionChecklistCard
          copy={copy.completionStatus}
          themeTokens={themeTokens}
          pointDetail={pointDetail}
          readiness={pointCompletionReadiness}
          isSharedRoute={isSharedRoute}
          isResearchMaterialConfirmed={isResearchMaterialConfirmed}
        />

        {routeKind === 'learning' ? (
          <div style={{ ...noticeCardStyle, display: 'grid', gap: '14px', background: themeTokens.noticeBackground, borderColor: themeTokens.noticeBorder, color: themeTokens.description }}>
            <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.completion.finalAction}</strong>
            <div style={{ fontSize: '15px', lineHeight: 1.65, color: themeTokens.description }}>
              {isPointCompleted ? copy.completion.completedLocked : copy.completion.readyInstruction}
            </div>
            <div style={runtimeButtonRowStyle}>
              {isPointCompleted ? (
                <button
                  type="button"
                  onClick={() => setIsCompletionCancelConfirmOpen(true)}
                  disabled={isSavingPointRuntime}
                  style={{
                    ...runtimePrimaryButtonStyle,
                    minHeight: '54px',
                    padding: '0 28px',
                    background: isSavingPointRuntime
                      ? themeTokens.secondaryButtonBackground
                      : 'linear-gradient(180deg, #FBBF24, #F97316)',
                    borderColor: isSavingPointRuntime ? themeTokens.secondaryButtonBorder : 'rgba(194, 65, 12, 0.45)',
                    color: isSavingPointRuntime ? themeTokens.buttonText : '#431407',
                    boxShadow: isSavingPointRuntime ? 'none' : '0 12px 28px rgba(234, 88, 12, 0.22)',
                    cursor: isSavingPointRuntime ? 'progress' : 'pointer',
                    opacity: isSavingPointRuntime ? 0.62 : 1,
                  }}
                >
                  {isSavingPointRuntime ? copy.completion.cancelProcessing : copy.completion.cancelCompletion}
                </button>
              ) : (
                <button
                  type="button"
                  onClick={() => setIsCompletionConfirmOpen(true)}
                  disabled={isCompletionButtonDisabled}
                  style={{
                    ...runtimePrimaryButtonStyle,
                    minHeight: '54px',
                    padding: '0 28px',
                    background: isCompletionButtonDisabled
                      ? themeTokens.primaryButtonBackground
                      : 'linear-gradient(180deg, rgba(246, 205, 103, 0.96), rgba(214, 153, 34, 0.98))',
                    borderColor: isCompletionButtonDisabled ? themeTokens.primaryButtonBorder : 'rgba(140, 96, 22, 0.42)',
                    color: isCompletionButtonDisabled ? themeTokens.buttonText : '#4A2C00',
                    boxShadow: isCompletionButtonDisabled
                      ? 'none'
                      : '0 14px 34px rgba(138, 95, 16, 0.24), 0 0 0 3px rgba(246, 205, 103, 0.22)',
                    transform: isCompletionButtonDisabled ? 'none' : 'translateY(-1px)',
                    ...((isCompletionButtonDisabled) ? runtimeButtonDisabledStyle : null),
                  }}
                >
                  {isSavingPointRuntime ? copy.completion.saving : copy.completion.completePoint}
                </button>
              )}
            </div>
          </div>
        ) : null}
        {isCompletionConfirmOpen ? (
          <LumiModalShell
            title={copy.completion.completeModalTitle}
            eyebrow={copy.completion.completeModalEyebrow}
            lumiState="celebrate"
            tone="warm"
            width={440}
            onClose={() => setIsCompletionConfirmOpen(false)}
            message={copy.completion.completeModalMessage}
            actions={(
              <>
                <button
                  type="button"
                  onClick={() => setIsCompletionConfirmOpen(false)}
                  style={lumiModalSecondaryButtonStyle}
                >
                  {copy.completion.cancel}
                </button>
                <button
                  type="button"
                  onClick={async () => {
                    setIsCompletionConfirmOpen(false);
                    const saved = await handleSaveSelfEvaluation();
                    if (saved) await handleSavePointRuntime('completed');
                  }}
                  style={lumiModalPrimaryButtonStyle}
                >
                  {copy.completion.complete}
                </button>
              </>
            )}
          />
        ) : null}
        {isCompletionCancelConfirmOpen ? (
          <LumiModalShell
            title={copy.completion.cancelModalTitle}
            eyebrow={copy.completion.completeModalEyebrow}
            lumiState="curious"
            tone="warm"
            width={440}
            onClose={() => {
              if (!isCompletionCancelExecuting) setIsCompletionCancelConfirmOpen(false);
            }}
            message={copy.completion.cancelModalMessage}
            actions={(
              <>
                <button
                  type="button"
                  onClick={() => setIsCompletionCancelConfirmOpen(false)}
                  disabled={isCompletionCancelExecuting}
                  style={lumiModalSecondaryButtonStyle}
                >
                  {copy.completion.close}
                </button>
                <button
                  type="button"
                  onClick={async () => {
                    setIsCompletionCancelExecuting(true);
                    await handleSavePointRuntime('in_progress');
                    setIsCompletionCancelExecuting(false);
                    setIsCompletionCancelConfirmOpen(false);
                  }}
                  disabled={isCompletionCancelExecuting}
                  style={{ ...lumiModalPrimaryButtonStyle, opacity: isCompletionCancelExecuting ? 0.64 : 1, cursor: isCompletionCancelExecuting ? 'progress' : 'pointer' }}
                >
                  {isCompletionCancelExecuting ? copy.completion.cancelProcessing : copy.completion.cancelDone}
                </button>
              </>
            )}
          />
        ) : null}
      </div>
    </section>
  );
}
