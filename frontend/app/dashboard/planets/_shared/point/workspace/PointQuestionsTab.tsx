'use client';

import { useEffect, useState } from 'react';
import LumiModalShell from '@/components/common/LumiModalShell';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../../pointPageUtils';
import type { UsePointLearningResult } from '../usePointLearning';
import { getQuestionAuthorLabel } from '../../pointPageUtils';
import {
  formGridStyle, sectionMiniHeaderStyle, blockMetaStyle,
  blockCardStyle, blockHeaderStyle,
  placeholderCardStyle, emptyCardStyle, formActionRowStyle,
  textareaStyle,
  pointShellTitleStyle,
} from '../../pointPageStyles';
import { makeWorkspaceStyles } from './workspaceStyles';
import { scrollToLearningRecordSection } from './workspaceTypes';

interface Props {
  copy: PointLearningCopy['workspace']['learningWork']['recordPanels']['questions'];
  commonCopy: PointLearningCopy['workspace']['learningWork']['recordPanels']['common'];
  learning: UsePointLearningResult;
  themeTokens: PointPageThemeTokens;
  canEditLearningRecords: boolean;
  isSharedRoute: boolean;
  onFocusClear: () => void;
  onEditingChange: (editing: boolean) => void;
}

const QUESTIONS_PER_PAGE = 5;

type DisplayAiCoachFeedback = {
  summary: string;
  nextAction: string;
  fallbackText: string;
};

function asString(value: unknown): string {
  return typeof value === 'string' ? value.trim() : '';
}

function pickString(record: Record<string, unknown>, keys: string[]): string {
  for (const key of keys) {
    const value = asString(record[key]);
    if (value) return value;
  }
  return '';
}

function pickQuotedField(source: string, keys: string[]): string {
  for (const key of keys) {
    const escapedKey = key.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    const pattern = new RegExp(`"?${escapedKey}"?\\s*:\\s*"([^"]*)"`, 'u');
    const match = source.match(pattern);
    const value = match?.[1]?.trim();
    if (value) return value;
  }
  return '';
}

function formatAiCoachFeedbackFromLooseText(source: string): DisplayAiCoachFeedback {
  const feedbackBlockMatch = source.match(/"?피드백"?\s*:\s*\{([\s\S]*?)\}\s*$/u)
    ?? source.match(/"?feedback"?\s*:\s*\{([\s\S]*?)\}\s*$/iu);
  const feedbackSource = feedbackBlockMatch?.[1] ?? source;
  const summary = pickQuotedField(feedbackSource, ['요약', '핵심관점', '핵심 관점', 'summary', 'Summary', 'feedback_summary']);
  const nextAction = pickQuotedField(feedbackSource, ['다음 행동 제안', '다음 행동', '다음행동', 'next_action', 'nextAction', 'Next Action', 'recommendation']);
  return summary || nextAction
    ? { summary, nextAction, fallbackText: '' }
    : { summary: '', nextAction: '', fallbackText: source.trim() };
}

function formatAiCoachFeedback(rawFeedback: string): DisplayAiCoachFeedback {
  const trimmed = rawFeedback.trim();
  if (!trimmed) return { summary: '', nextAction: '', fallbackText: '' };

  const looseResult = formatAiCoachFeedbackFromLooseText(trimmed);
  if (looseResult.summary || looseResult.nextAction) return looseResult;

  const parseCandidates = [
    trimmed,
    trimmed.startsWith('{') ? '' : `{${trimmed}`,
  ].filter(Boolean);

  for (const candidate of parseCandidates) {
    try {
      const parsed = JSON.parse(candidate) as unknown;
    if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) {
        continue;
    }

    const root = parsed as Record<string, unknown>;
    const feedback = root['피드백'] && typeof root['피드백'] === 'object' && !Array.isArray(root['피드백'])
      ? root['피드백'] as Record<string, unknown>
      : root.feedback && typeof root.feedback === 'object' && !Array.isArray(root.feedback)
        ? root.feedback as Record<string, unknown>
        : root;
      const summary = pickString(feedback, ['요약', '핵심관점', '핵심 관점', 'summary', 'Summary', 'feedback_summary']);
      const nextAction = pickString(feedback, ['다음 행동 제안', '다음 행동', '다음행동', 'next_action', 'nextAction', 'Next Action', 'recommendation']);
    if (summary || nextAction) return { summary, nextAction, fallbackText: '' };
    } catch {
      // Try the next parse candidate; non-JSON feedback falls through below.
    }
  }

  return { summary: '', nextAction: '', fallbackText: trimmed };
}

export default function PointQuestionsTab({ copy, commonCopy, learning, themeTokens, canEditLearningRecords, isSharedRoute, onFocusClear, onEditingChange }: Props) {
  const {
    pointQuestions, newQuestionText, setNewQuestionText,
    isSavingQuestion, isGeneratingFeedbackQuestionID,
    refreshPointDetail, handleCreateQuestion, handleChangeQuestionDraft,
    handleSaveQuestion, handleDeleteQuestion, handleGenerateQuestionFeedback,
  } = learning;

  const {
    materialToolbarButtonStyle, materialToolbarPrimaryButtonStyle, materialToolbarDangerButtonStyle,
    listBoardStyle, getListHeaderRowStyle, getListDataRowStyle, listColumnHeaderStyle,
    getListIndexCellStyle, getListTextCellStyle, readonlyRecordFieldStyle, readonlyRecordValueStyle,
  } = makeWorkspaceStyles(themeTokens);

  const [selectedQuestionID, setSelectedQuestionID] = useState<string | null>(null);
  const [questionMode, setQuestionMode] = useState<'view' | 'edit'>('view');
  const [questionPage, setQuestionPage] = useState(1);
  const [visibleQuestionAnswerIDs, setVisibleQuestionAnswerIDs] = useState<Set<string>>(() => new Set());
  const [isAddingQuestion, setIsAddingQuestion] = useState(false);
  const [isConfirmingQuestionListReturn, setIsConfirmingQuestionListReturn] = useState(false);
  const [isConfirmingQuestionAddCancel, setIsConfirmingQuestionAddCancel] = useState(false);
  const [pendingDeleteQuestion, setPendingDeleteQuestion] = useState<typeof pointQuestions[number] | null>(null);
  const [questionSaveModal, setQuestionSaveModal] = useState<{ title: string; message: string; tone: 'warm' | 'alert' } | null>(null);
  const [hoveredListRowID, setHoveredListRowID] = useState<string | null>(null);

  useEffect(() => {
    const maxPage = Math.max(1, Math.ceil(pointQuestions.length / QUESTIONS_PER_PAGE));
    setQuestionPage((current) => Math.min(current, maxPage));
    if (selectedQuestionID && !pointQuestions.some((q) => q.id === selectedQuestionID)) {
      setSelectedQuestionID(null);
      setQuestionMode('view');
    }
  }, [pointQuestions, selectedQuestionID]);

  const selectedQuestion = selectedQuestionID ? pointQuestions.find((q) => q.id === selectedQuestionID) ?? null : null;
  const isQuestionEditing = Boolean(selectedQuestion && questionMode === 'edit' && canEditLearningRecords && !isSharedRoute);
  const isAnswerEditing = Boolean(selectedQuestion && visibleQuestionAnswerIDs.has(selectedQuestion.id) && canEditLearningRecords && !isSharedRoute);
  const isEditing = Boolean(isQuestionEditing || isAnswerEditing || isAddingQuestion);
  useEffect(() => { onEditingChange(isEditing); }, [isEditing, onEditingChange]);
  const selectedAiCoachFeedback = selectedQuestion ? formatAiCoachFeedback(selectedQuestion.aiFeedback) : null;
  const renderAiCoachFeedback = (feedback: DisplayAiCoachFeedback) => (
    <div style={{ display: 'grid', gap: '10px', fontSize: '15px', lineHeight: 1.65 }}>
      {feedback.summary ? (
        <div>
          <strong style={{ display: 'block', color: themeTokens.title, fontSize: '13px', marginBottom: '3px' }}>{copy.detail.aiCoachSummary}</strong>
          <span style={{ whiteSpace: 'pre-wrap' }}>{feedback.summary}</span>
        </div>
      ) : null}
      {feedback.nextAction ? (
        <div>
          <strong style={{ display: 'block', color: themeTokens.title, fontSize: '13px', marginBottom: '3px' }}>{copy.detail.aiCoachNextAction}</strong>
          <span style={{ whiteSpace: 'pre-wrap' }}>{feedback.nextAction}</span>
        </div>
      ) : null}
      {feedback.fallbackText ? <span style={{ whiteSpace: 'pre-wrap' }}>{feedback.fallbackText}</span> : null}
    </div>
  );

  const getRecentTime = (r: { updatedAt?: string; createdAt?: string }) => {
    const parsed = Date.parse(r.updatedAt || r.createdAt || '');
    return Number.isFinite(parsed) ? parsed : 0;
  };
  const sortedPointQuestions = [...pointQuestions].sort((a, b) => getRecentTime(b) - getRecentTime(a));
  const questionPageCount = Math.max(1, Math.ceil(sortedPointQuestions.length / QUESTIONS_PER_PAGE));
  const pagedQuestions = sortedPointQuestions.slice((questionPage - 1) * QUESTIONS_PER_PAGE, questionPage * QUESTIONS_PER_PAGE);

  const handleConfirmQuestionListReturn = () => {
    setIsConfirmingQuestionListReturn(false);
    setSelectedQuestionID(null);
    setQuestionMode('view');
    scrollToLearningRecordSection();
  };

  const handleSaveQuestionDetail = async (question: typeof pointQuestions[number]) => {
    const saved = await handleSaveQuestion(question);
    setQuestionSaveModal(saved
      ? { title: copy.detail.saveSuccessTitle, message: copy.detail.saveSuccessMessage, tone: 'warm' }
      : { title: copy.detail.saveFailTitle, message: copy.detail.saveFailMessage, tone: 'alert' });
    if (saved) {
      setSelectedQuestionID(null);
      setQuestionMode('view');
      scrollToLearningRecordSection();
    }
  };

  const handleSaveQuestionAnswer = async (question: typeof pointQuestions[number]) => {
    const saved = await handleSaveQuestion(question);
    setQuestionSaveModal(saved
      ? { title: copy.detail.answerSaveSuccessTitle, message: copy.detail.answerSaveSuccessMessage, tone: 'warm' }
      : { title: copy.detail.answerSaveFailTitle, message: copy.detail.answerSaveFailMessage, tone: 'alert' });
    if (saved) {
      setVisibleQuestionAnswerIDs((current) => {
        const next = new Set(current);
        next.delete(question.id);
        return next;
      });
      scrollToLearningRecordSection();
    }
  };

  const handleCreateQuestionAndClose = async () => {
    const created = await handleCreateQuestion();
    if (created) {
      setIsAddingQuestion(false);
      scrollToLearningRecordSection();
    }
  };

  const handleConfirmDeleteQuestion = async () => {
    if (!pendingDeleteQuestion || isSavingQuestion) return;
    const question = pendingDeleteQuestion;
    setPendingDeleteQuestion(null);
    if (selectedQuestionID === question.id) { setSelectedQuestionID(null); setQuestionMode('view'); }
    await handleDeleteQuestion(question);
    scrollToLearningRecordSection();
  };

  const handleConfirmQuestionAddCancel = () => {
    setIsConfirmingQuestionAddCancel(false);
    setIsAddingQuestion(false);
    scrollToLearningRecordSection();
  };

  const handleOpenQuestionDetail = async (questionID: string) => {
    await refreshPointDetail();
    setIsAddingQuestion(false);
    setSelectedQuestionID(questionID);
    setQuestionMode('view');
    scrollToLearningRecordSection();
  };

  return (
    <>
      {pendingDeleteQuestion ? (
        <LumiModalShell
          title={copy.detail.deleteModalTitle} eyebrow="Lumi Confirm" lumiState="curious" tone="alert" width={440}
          onClose={() => { if (!isSavingQuestion) setPendingDeleteQuestion(null); }}
          message={copy.detail.deleteModalMessage(pendingDeleteQuestion.title || pendingDeleteQuestion.question || copy.fallbackTitle)}
          actions={(
            <>
              <button type="button" onClick={() => setPendingDeleteQuestion(null)} disabled={isSavingQuestion} style={{ ...materialToolbarButtonStyle, opacity: isSavingQuestion ? 0.62 : 1 }}>{copy.detail.cancel}</button>
              <button type="button" onClick={() => void handleConfirmDeleteQuestion()} disabled={isSavingQuestion} style={{ ...materialToolbarDangerButtonStyle, opacity: isSavingQuestion ? 0.72 : 1 }}>{isSavingQuestion ? copy.detail.deleting : copy.detail.delete}</button>
            </>
          )}
        />
      ) : null}
      {isConfirmingQuestionListReturn ? (
        <LumiModalShell
          title={copy.detail.listReturnTitle} eyebrow="Lumi Confirm" lumiState="curious" tone="alert" width={440}
          onClose={() => setIsConfirmingQuestionListReturn(false)}
          message={copy.detail.listReturnMessage}
          actions={(
            <>
              <button type="button" onClick={() => setIsConfirmingQuestionListReturn(false)} style={materialToolbarButtonStyle}>{copy.detail.back}</button>
              <button type="button" onClick={handleConfirmQuestionListReturn} style={materialToolbarButtonStyle}>{copy.detail.list}</button>
            </>
          )}
        />
      ) : null}
      {isConfirmingQuestionAddCancel ? (
        <LumiModalShell
          title={copy.detail.addCancelTitle} eyebrow="Lumi Confirm" lumiState="curious" tone="alert" width={440}
          onClose={() => setIsConfirmingQuestionAddCancel(false)}
          message={copy.detail.addCancelMessage}
          actions={(
            <>
              <button type="button" onClick={() => setIsConfirmingQuestionAddCancel(false)} style={materialToolbarButtonStyle}>{copy.detail.back}</button>
              <button type="button" onClick={handleConfirmQuestionAddCancel} style={materialToolbarButtonStyle}>{copy.detail.cancel}</button>
            </>
          )}
        />
      ) : null}
      {questionSaveModal ? (
        <LumiModalShell
          title={questionSaveModal.title} eyebrow="Lumi" lumiState={questionSaveModal.tone === 'warm' ? 'happy' : 'curious'} tone={questionSaveModal.tone} width={440}
          onClose={() => setQuestionSaveModal(null)}
          message={questionSaveModal.message}
          actions={<button type="button" onClick={() => setQuestionSaveModal(null)} style={materialToolbarButtonStyle}>{copy.detail.ok}</button>}
        />
      ) : null}
      <div style={formGridStyle} onMouseEnter={onFocusClear} onFocusCapture={onFocusClear}>
        <div style={sectionMiniHeaderStyle}>
          <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.title}</strong>
          {!selectedQuestion && canEditLearningRecords && !isAddingQuestion ? (
            <button
              type="button"
              onClick={() => { setIsAddingQuestion(true); scrollToLearningRecordSection(); }}
              style={materialToolbarPrimaryButtonStyle}
            >
              {copy.add}
            </button>
          ) : null}
        </div>
        {selectedQuestion ? (
          <div style={{ ...blockCardStyle, background: themeTokens.surfaceBackground, borderColor: themeTokens.surfaceBorder }}>
            <div style={blockHeaderStyle}>
              <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{isQuestionEditing ? copy.detail.editTitle : copy.detail.viewTitle}</strong>
              <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{getQuestionAuthorLabel(selectedQuestion.createdBy)} / {selectedQuestion.status}</span>
            </div>
            {isQuestionEditing ? (
              <>
                <label style={{ display: 'grid', gap: '7px' }}>
                  <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.content}</span>
                  <textarea value={selectedQuestion.question} onChange={(e) => handleChangeQuestionDraft(selectedQuestion.id, 'question', e.target.value)} placeholder={copy.detail.contentPlaceholder} readOnly={isSharedRoute || !canEditLearningRecords} style={{ ...textareaStyle, minHeight: '92px', background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} />
                </label>
              </>
            ) : (
              <div style={{ display: 'grid', gap: '12px' }}>
                <div style={readonlyRecordFieldStyle}>
                  <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.content}</span>
                  <p style={readonlyRecordValueStyle}>{selectedQuestion.question.trim() || copy.detail.noContent}</p>
                </div>
                <div style={{ ...placeholderCardStyle, background: themeTokens.placeholderBackground, borderColor: themeTokens.placeholderBorder, color: themeTokens.description }}>
                  <div style={sectionMiniHeaderStyle}>
                    <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.detail.answer}</strong>
                    <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{selectedQuestion.answer.trim() ? copy.detail.answerMeta : copy.detail.answerEmptyMeta}</span>
                  </div>
                  {isAnswerEditing ? (
                    <div style={{ display: 'grid', gap: '10px' }}>
                      <textarea style={{ ...textareaStyle, minHeight: '118px', background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={selectedQuestion.answer} onChange={(e) => handleChangeQuestionDraft(selectedQuestion.id, 'answer', e.target.value)} placeholder={copy.detail.answerPlaceholder} readOnly={isSharedRoute || !canEditLearningRecords} />
                      <div style={{ ...formActionRowStyle, justifyContent: 'space-between' }}>
                        <button
                          type="button"
                          onClick={() => {
                            setVisibleQuestionAnswerIDs((current) => {
                              const next = new Set(current);
                              next.delete(selectedQuestion.id);
                              return next;
                            });
                          }}
                          disabled={isSavingQuestion || Boolean(isGeneratingFeedbackQuestionID)}
                          style={{ ...materialToolbarButtonStyle, opacity: isSavingQuestion || Boolean(isGeneratingFeedbackQuestionID) ? 0.62 : 1 }}
                        >
                          {copy.detail.cancel}
                        </button>
                        <div style={{ display: 'flex', gap: '8px', justifyContent: 'flex-end', flexWrap: 'wrap' }}>
                          <button type="button" onClick={() => { handleGenerateQuestionFeedback(selectedQuestion); scrollToLearningRecordSection(); }} disabled={Boolean(isGeneratingFeedbackQuestionID) || !selectedQuestion.question.trim()} style={{ ...materialToolbarButtonStyle, cursor: isGeneratingFeedbackQuestionID === selectedQuestion.id ? 'progress' : 'pointer', opacity: Boolean(isGeneratingFeedbackQuestionID) || !selectedQuestion.question.trim() ? 0.62 : 1 }}>
                            {isGeneratingFeedbackQuestionID === selectedQuestion.id ? copy.detail.generatingAiCoach : selectedQuestion.answer.trim() ? copy.detail.createAiCoachRevision : copy.detail.createAiCoachStarter}
                          </button>
                          <button type="button" onClick={() => void handleSaveQuestionAnswer(selectedQuestion)} disabled={isSavingQuestion || !selectedQuestion.answer.trim()} style={{ ...materialToolbarPrimaryButtonStyle, cursor: isSavingQuestion ? 'progress' : 'pointer', opacity: isSavingQuestion || !selectedQuestion.answer.trim() ? 0.62 : 1 }}>
                            {isSavingQuestion ? copy.detail.savingAnswer : copy.detail.saveAnswer}
                          </button>
                        </div>
                      </div>
                      <div style={{ color: themeTokens.mutedText, fontSize: '12px', fontWeight: 750, lineHeight: 1.45, textAlign: 'right' }}>
                        {selectedQuestion.answer.trim() ? copy.detail.aiCoachRevisionHint : copy.detail.aiCoachStarterHint}
                      </div>
                    </div>
                  ) : selectedQuestion.answer.trim() ? (
                    <p style={{ ...readonlyRecordValueStyle, whiteSpace: 'pre-wrap', margin: 0 }}>{selectedQuestion.answer}</p>
                  ) : (
                    <div style={{ ...emptyCardStyle, background: themeTokens.emptyBackground, borderColor: themeTokens.emptyBorder, color: themeTokens.mutedText }}>{copy.detail.noAnswer}</div>
                  )}
                </div>
                {selectedAiCoachFeedback && selectedQuestion.aiFeedback.trim() ? (
                  <div style={{ ...placeholderCardStyle, background: themeTokens.placeholderBackground, borderColor: themeTokens.placeholderBorder, color: themeTokens.description }}>
                    <div style={sectionMiniHeaderStyle}>
                      <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.detail.aiCoach}</strong>
                      <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.aiCoachMeta}</span>
                    </div>
                    {renderAiCoachFeedback(selectedAiCoachFeedback)}
                  </div>
                ) : null}
              </div>
            )}
            <div style={{ display: 'grid', gap: '8px' }}>
              <div style={{ ...formActionRowStyle, justifyContent: 'space-between', alignItems: 'center' }}>
                {canEditLearningRecords ? (
                  <>
                    <div style={{ display: 'flex', justifyContent: 'flex-start', gap: '8px' }}>
                      <button type="button" onClick={() => { if (isQuestionEditing) setIsConfirmingQuestionListReturn(true); else handleConfirmQuestionListReturn(); }} style={materialToolbarButtonStyle}>{copy.detail.list}</button>
                      <button type="button" onClick={() => setPendingDeleteQuestion(selectedQuestion)} disabled={isSavingQuestion} style={{ ...materialToolbarDangerButtonStyle, opacity: isSavingQuestion ? 0.62 : 1 }}>{copy.detail.delete}</button>
                    </div>
                    <div style={{ flex: 1 }} />
                    <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '8px' }}>
                      {isQuestionEditing ? (
                        <button type="button" onClick={() => void handleSaveQuestionDetail(selectedQuestion)} disabled={isSavingQuestion || !selectedQuestion.question.trim()} style={{ ...materialToolbarPrimaryButtonStyle, cursor: isSavingQuestion ? 'progress' : 'pointer', opacity: isSavingQuestion || !selectedQuestion.question.trim() ? 0.62 : 1 }}>
                          {isSavingQuestion ? copy.detail.savingDetail : copy.detail.saveDetail}
                        </button>
                      ) : !isAnswerEditing ? (
                        <>
                          <button type="button" onClick={() => { setVisibleQuestionAnswerIDs((s) => new Set(s).add(selectedQuestion.id)); scrollToLearningRecordSection(); }} style={materialToolbarButtonStyle}>
                            {selectedQuestion.answer.trim() ? copy.detail.editAnswer : copy.detail.addAnswer}
                          </button>
                          <button type="button" onClick={() => { setQuestionMode('edit'); scrollToLearningRecordSection(); }} style={materialToolbarButtonStyle}>{copy.detail.edit}</button>
                        </>
                      ) : null}
                    </div>
                  </>
                ) : (
                  <button type="button" onClick={handleConfirmQuestionListReturn} style={materialToolbarButtonStyle}>{copy.detail.list}</button>
                )}
              </div>
            </div>
          </div>
        ) : !isAddingQuestion && pointQuestions.length ? (
          <div id="point-question-list-section" style={{ ...formGridStyle, scrollMarginTop: '430px' }}>
            <div style={sectionMiniHeaderStyle}>
              <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{commonCopy.countRange(pointQuestions.length, (questionPage - 1) * QUESTIONS_PER_PAGE + 1, Math.min(questionPage * QUESTIONS_PER_PAGE, pointQuestions.length))}</span>
            </div>
            <div role="table" aria-label={copy.listAria} style={listBoardStyle}>
              <div role="row" style={getListHeaderRowStyle('58px minmax(0, 1fr) 132px')}>
                <span role="columnheader" style={listColumnHeaderStyle}>{commonCopy.number}</span>
                <span role="columnheader" style={listColumnHeaderStyle}>{copy.detail.content}</span>
                <span role="columnheader" style={listColumnHeaderStyle}>{commonCopy.status}</span>
              </div>
              {pagedQuestions.map((question, index) => {
                const isAnswered = question.status === 'answered';
                const statusColor = isAnswered ? (themeTokens.pageBackground === '#F8FAFC' ? '#047857' : '#A7F3D0') : (themeTokens.pageBackground === '#F8FAFC' ? '#B45309' : '#FDE68A');
                const rowID = `question-${question.id}`;
                const hovered = hoveredListRowID === rowID;
                const questionIndex = sortedPointQuestions.length - ((questionPage - 1) * QUESTIONS_PER_PAGE + index);
                return (
                  <button key={question.id} type="button" role="row" onClick={() => void handleOpenQuestionDetail(question.id)} onMouseEnter={() => setHoveredListRowID(rowID)} onMouseLeave={() => setHoveredListRowID(null)} style={getListDataRowStyle('58px minmax(0, 1fr) 132px', hovered)}>
                    <span role="cell" style={getListIndexCellStyle(hovered)}>{questionIndex}</span>
                    <span role="cell" style={getListTextCellStyle(hovered)}>{question.question.trim() || copy.detail.noContent}</span>
                    <span role="cell" style={{ ...blockMetaStyle, color: hovered ? themeTokens.title : statusColor, fontWeight: 900 }}>{isAnswered ? copy.answered : copy.pending}</span>
                  </button>
                );
              })}
            </div>
            {questionPageCount > 1 ? (
              <div style={{ ...formActionRowStyle, justifyContent: 'center' }}>
                <button type="button" onClick={() => setQuestionPage((p) => Math.max(1, p - 1))} disabled={questionPage <= 1} style={{ ...materialToolbarButtonStyle, opacity: questionPage <= 1 ? 0.55 : 1 }}>{commonCopy.previous}</button>
                <span style={{ color: themeTokens.mutedText, fontSize: '14px', fontWeight: 800 }}>{questionPage} / {questionPageCount}</span>
                <button type="button" onClick={() => setQuestionPage((p) => Math.min(questionPageCount, p + 1))} disabled={questionPage >= questionPageCount} style={{ ...materialToolbarButtonStyle, opacity: questionPage >= questionPageCount ? 0.55 : 1 }}>{commonCopy.next}</button>
              </div>
            ) : null}
          </div>
        ) : !isAddingQuestion ? (
          <div style={{ ...emptyCardStyle, background: themeTokens.emptyBackground, borderColor: themeTokens.emptyBorder, color: themeTokens.mutedText }}>{copy.empty}</div>
        ) : null}
        {!selectedQuestion && canEditLearningRecords && isAddingQuestion ? (
          <div style={{ ...placeholderCardStyle, background: themeTokens.placeholderBackground, borderColor: themeTokens.placeholderBorder }}>
            <div style={formGridStyle}>
              <div style={sectionMiniHeaderStyle}>
                <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.detail.addTitle}</strong>
                <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.addSubtitle}</span>
              </div>
              <label style={{ display: 'grid', gap: '7px' }}>
                <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.content}</span>
                <textarea value={newQuestionText} onChange={(e) => setNewQuestionText(e.target.value)} placeholder={copy.detail.contentPlaceholder} style={{ ...textareaStyle, minHeight: '92px', background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} />
              </label>
              <div style={{ color: themeTokens.mutedText, fontSize: '14px', lineHeight: 1.55 }}>{copy.detail.answerInstruction}</div>
              <div style={formActionRowStyle}>
                <button type="button" disabled style={{ ...materialToolbarButtonStyle, background: themeTokens.disabledButtonBackground, borderColor: themeTokens.disabledButtonBorder, color: themeTokens.disabledButtonText, opacity: 0.72 }}>{copy.detail.communityDisabled}</button>
                <button type="button" onClick={() => setIsConfirmingQuestionAddCancel(true)} disabled={isSavingQuestion} style={{ ...materialToolbarButtonStyle, opacity: isSavingQuestion ? 0.62 : 1 }}>{copy.detail.cancel}</button>
                <button type="button" onClick={() => void handleCreateQuestionAndClose()} disabled={isSavingQuestion || !newQuestionText.trim()} style={{ ...materialToolbarPrimaryButtonStyle, cursor: isSavingQuestion ? 'progress' : 'pointer', opacity: isSavingQuestion || !newQuestionText.trim() ? 0.62 : 1 }}>
                  {isSavingQuestion ? copy.detail.saving : copy.detail.register}
                </button>
              </div>
            </div>
          </div>
        ) : null}
      </div>
    </>
  );
}
