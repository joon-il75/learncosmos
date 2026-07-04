'use client';

import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../pointPageUtils';
import {
  formActionRowStyle,
  noticeCardStyle,
  placeholderCardStyle,
  pointShellTitleStyle,
  secondaryButtonStyle,
  sectionMiniHeaderStyle,
  textareaStyle,
} from '../pointPageStyles';
import type { UsePointLearningResult } from './usePointLearning';

function shouldHideSelfEvaluationNote(note: string): boolean {
  const normalized = note.trim();
  return normalized.includes('수익화를 위해')
    && normalized.includes('강의 주제를 선정')
    && normalized.includes('직접 연결됨');
}

export function PointSelfEvaluationPanel({
  copy,
  learning,
  themeTokens,
  canEditPoint,
}: {
  copy: PointLearningCopy['selfEvaluation'];
  learning: UsePointLearningResult;
  themeTokens: PointPageThemeTokens;
  canEditPoint: boolean;
}) {
  const {
    pointDetail,
    selfEvalUnderstandingReason,
    selfEvalApplicationReason,
    selfEvalProficiencyReason,
    selfEvalProblemSolvingReason,
    selfEvalExpressionReason,
    selfEvalGoalAlignmentNote,
    selfEvalApplicationQuestions,
    selfEvalApplicationAnswers,
    hasSelfEvalDraft,
    isSelfEvalLocked,
    finalSelfEvalScore,
    setFinalSelfEvalScore,
    isGeneratingSelfEvalDraft,
    handleChangeSelfEvalApplicationAnswer,
    handleRestartSelfEvaluation,
    handleGenerateSelfEvaluationDraft,
  } = learning;

  if (!pointDetail) return null;

  const hasApplicationQuestions = selfEvalApplicationQuestions.length > 0;
  const hasAllApplicationAnswers = hasApplicationQuestions
    && selfEvalApplicationQuestions.every((_, index) => (selfEvalApplicationAnswers[index] ?? '').trim().length > 0);
  const canRequestSelfEvaluation = !isGeneratingSelfEvalDraft && (!hasApplicationQuestions || hasAllApplicationAnswers);
  const aiEvaluationNotes = [
    selfEvalUnderstandingReason,
    selfEvalApplicationReason,
    selfEvalProficiencyReason,
    selfEvalProblemSolvingReason,
    selfEvalExpressionReason,
  ].map((item) => item.trim()).filter(Boolean);
  const visibleGoalAlignmentNote = shouldHideSelfEvaluationNote(selfEvalGoalAlignmentNote)
    ? ''
    : selfEvalGoalAlignmentNote.trim();

  return (
    <>
      <div style={sectionMiniHeaderStyle}>
        <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.title}</strong>
        <span style={{ fontSize: '14px', color: themeTokens.metaLabel }}>{copy.subtitle}</span>
      </div>

      <div style={{ ...placeholderCardStyle, background: themeTokens.placeholderBackground, borderColor: themeTokens.placeholderBorder }}>
        <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>
          {pointDetail.self_evaluation ? copy.savedTitle : copy.applicationTitle}
        </strong>
        <div style={{ marginTop: '8px', fontSize: '15px', lineHeight: 1.65, color: themeTokens.description }}>
          {pointDetail.self_evaluation
            ? copy.savedDescription
            : copy.applicationDescription}
        </div>
        {visibleGoalAlignmentNote ? (
          <div style={{ marginTop: '10px', fontSize: '14px', lineHeight: 1.6, color: themeTokens.mutedText }}>
            {visibleGoalAlignmentNote}
          </div>
        ) : null}
      </div>

      {isGeneratingSelfEvalDraft ? (
        <div style={{ ...noticeCardStyle, display: 'grid', gap: '10px', background: themeTokens.noticeBackground, borderColor: themeTokens.noticeBorder, color: themeTokens.description }}>
          <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.generatingTitle}</strong>
          <div style={{ height: '10px', overflow: 'hidden', borderRadius: '999px', background: themeTokens.placeholderBackground, border: `1px solid ${themeTokens.placeholderBorder}` }}>
            <div style={{ width: '68%', height: '100%', borderRadius: '999px', background: themeTokens.primaryButtonBackground, boxShadow: '0 0 18px rgba(59, 130, 246, 0.35)' }} />
          </div>
          <div style={{ fontSize: '14px', lineHeight: 1.6, color: themeTokens.mutedText }}>
            {copy.generatingDescription}
          </div>
        </div>
      ) : null}

      {hasApplicationQuestions ? selfEvalApplicationQuestions.map((item, index) => (
        <div key={`${item.question}-${index}`} style={{ ...placeholderCardStyle, background: themeTokens.surfaceBackground, borderColor: themeTokens.surfaceBorder }}>
          <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.questionLabel} {index + 1}</strong>
          <div style={{ marginTop: '8px', fontSize: '15px', lineHeight: 1.65, color: themeTokens.description }}>{item.question}</div>
          {item.intent ? (
            <div style={{ marginTop: '6px', fontSize: '13px', color: themeTokens.metaLabel }}>{item.intent}</div>
          ) : null}
          <textarea
            style={{ ...textareaStyle, minHeight: '86px', marginTop: '12px', background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }}
            value={selfEvalApplicationAnswers[index] ?? ''}
            onChange={(e) => handleChangeSelfEvalApplicationAnswer(index, e.target.value)}
            placeholder={copy.answerPlaceholder}
            readOnly={!canEditPoint || isSelfEvalLocked}
          />
        </div>
      )) : null}

      {hasSelfEvalDraft ? (
        <div style={{ ...noticeCardStyle, display: 'grid', gap: '10px', background: themeTokens.noticeBackground, borderColor: themeTokens.noticeBorder, color: themeTokens.description }}>
          <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.aiReferenceTitle}</strong>
          <div style={{ fontSize: '15px', lineHeight: 1.65, color: themeTokens.description }}>
            {copy.aiReferenceDescription}
          </div>
          {aiEvaluationNotes.length ? (
            <div style={{ display: 'grid', gap: '7px', fontSize: '14px', lineHeight: 1.6, color: themeTokens.description }}>
              {aiEvaluationNotes.map((note, index) => (
                <div key={`${note}-${index}`}>• {note}</div>
              ))}
            </div>
          ) : null}
        </div>
      ) : null}

      {canEditPoint ? (
        <div style={formActionRowStyle}>
          <button
            type="button"
            onClick={isSelfEvalLocked ? handleRestartSelfEvaluation : () => handleGenerateSelfEvaluationDraft()}
            disabled={isSelfEvalLocked ? isGeneratingSelfEvalDraft : !canRequestSelfEvaluation}
            style={{ ...secondaryButtonStyle, cursor: (isSelfEvalLocked ? isGeneratingSelfEvalDraft : !canRequestSelfEvaluation) ? 'not-allowed' : 'pointer', opacity: (!isSelfEvalLocked && hasApplicationQuestions && !hasAllApplicationAnswers) ? 0.55 : 1, background: themeTokens.secondaryButtonBackground, borderColor: themeTokens.secondaryButtonBorder, color: themeTokens.buttonText }}
          >
            {isGeneratingSelfEvalDraft ? copy.generatingButton : isSelfEvalLocked ? copy.restartButton : hasApplicationQuestions ? copy.aiReferenceButton : copy.createQuestionsButton}
          </button>
        </div>
      ) : null}

      <div style={{ ...placeholderCardStyle, background: themeTokens.placeholderBackground, borderColor: themeTokens.placeholderBorder }}>
        <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.finalTitle}</strong>
        <div style={{ marginTop: '8px', fontSize: '15px', lineHeight: 1.65, color: themeTokens.description }}>
          {copy.finalDescription}
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(118px, 1fr))', gap: '10px', marginTop: '14px' }}>
          {copy.options.map((option) => {
            const selected = finalSelfEvalScore === option.value;
            const disabled = !canEditPoint || !hasSelfEvalDraft;
            return (
              <label
                key={option.value}
                style={{
                  display: 'grid',
                  gap: '6px',
                  padding: '12px',
                  borderRadius: '12px',
                  border: `1px solid ${selected ? themeTokens.primaryButtonBorder : themeTokens.inputBorder}`,
                  background: selected ? themeTokens.noticeBackground : themeTokens.inputBackground,
                  color: selected ? themeTokens.title : themeTokens.description,
                  cursor: disabled ? 'not-allowed' : 'pointer',
                  opacity: disabled ? 0.55 : 1,
                }}
              >
                <span style={{ display: 'flex', alignItems: 'center', gap: '8px', fontWeight: 800 }}>
                  <input
                    type="radio"
                    name="final-self-evaluation"
                    value={option.value}
                    checked={selected}
                    onChange={(e) => setFinalSelfEvalScore(e.target.value)}
                    disabled={disabled}
                  />
                  {option.label}{copy.scoreSuffix}
                </span>
                <span style={{ fontSize: '13px', lineHeight: 1.45 }}>{option.description}</span>
              </label>
            );
          })}
        </div>
        {!hasSelfEvalDraft ? (
          <div style={{ marginTop: '10px', fontSize: '14px', lineHeight: 1.6, color: themeTokens.mutedText }}>
            {copy.lockedUntilDraft}
          </div>
        ) : null}
      </div>
    </>
  );
}
