'use client';

import Link from 'next/link';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import LumiAvatar from '@/components/lumi/LumiAvatar';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../pointPageUtils';
import type { ObservationNoteType, PointWorkTab } from '../pointPageTypes';
import { NodeSection } from '../pointPageUtils';
import { noticeCardStyle, primaryButtonStyle, secondaryButtonStyle } from '../pointPageStyles';
import type { UsePointLearningResult } from './usePointLearning';
import type { LearningWorkGuideDetail } from './workspace/workspaceTypes';
import PointWorkspace from './PointWorkspace';
import PointCompletion from './PointCompletion';
import { PointFlowProgress } from './PointFlowProgress';
import {
  getNextPointLearningStage,
  getPreviousPointLearningStage,
  goalObjectTypes,
  pointLearningStages,
  repeatRecordLensToWorkTab,
  type GoalObjectType,
  type PointLearningStage,
  type RepeatRecordLens,
} from './skillFlowTypes';

const observationNoteTypes: ObservationNoteType[] = ['core_summary', 'revisit_part', 'reference_material'];

interface Props {
  copy: PointLearningCopy;
  learning: UsePointLearningResult;
  themeTokens: PointPageThemeTokens;
  isExplorationPoint: boolean;
  isSharedRoute: boolean;
  editorToolbarStickyTop: string;
  pointGoalContextText: string;
  onResearchMaterialConfirmed: () => void;
  onLearningWorkGuideChange: (guide: LearningWorkGuideDetail) => void;
  onSourceContentOpen: () => void;
  activeStage: PointLearningStage;
  onStageChange?: (stage: PointLearningStage) => void;
  getStageHref: (stage: PointLearningStage) => string;
  diaryHref: string;
  researchMaterialHref: string;
}

export default function PointLearningFlow({
  copy,
  learning,
  themeTokens,
  isExplorationPoint,
  isSharedRoute,
  editorToolbarStickyTop,
  pointGoalContextText,
  onResearchMaterialConfirmed,
  onLearningWorkGuideChange,
  onSourceContentOpen,
  activeStage,
  onStageChange,
  getStageHref,
  diaryHref,
  researchMaterialHref,
}: Props) {
  const router = useRouter();
  const [goalType, setGoalType] = useState<GoalObjectType>('conceptual');
  const [goalObjectText, setGoalObjectText] = useState('');
  const [selectedSkills, setSelectedSkills] = useState<string[]>([]);
  const [activeLens, setActiveLens] = useState<RepeatRecordLens>('attempt');
  const [isObservationContentOpen, setIsObservationContentOpen] = useState(false);
  const [isObservationSourceOpened, setIsObservationSourceOpened] = useState(false);
  const [observationNoteType, setObservationNoteType] = useState<ObservationNoteType>('core_summary');
  const [observationNoteContent, setObservationNoteContent] = useState('');
  const [editingObservationNoteID, setEditingObservationNoteID] = useState<string | null>(null);
  const [maxVisitedStageIndex, setMaxVisitedStageIndex] = useState(0);
  const [researchEditModeSignal, setResearchEditModeSignal] = useState(0);
  const [isCompactObservationLayout, setIsCompactObservationLayout] = useState(false);
  const flowRootRef = useRef<HTMLDivElement>(null);

  const stageCopy = copy.skillFlow.stages[activeStage];
  const shouldHideObservationChrome = activeStage === 'point' && !isExplorationPoint && isObservationContentOpen;
  const suggestedSkills = copy.skillFlow.skills.presets[goalType];
  const activeLensCopy = copy.skillFlow.repeat.lenses[activeLens];
  const activeLensWorkTab = repeatRecordLensToWorkTab[activeLens];
  const activeLensRecordCount = (() => {
    if (activeLensWorkTab === 'questions') return learning.pointQuestions.length;
    if (activeLensWorkTab === 'practice') return learning.practiceLogDrafts.length;
    if (activeLensWorkTab === 'artifacts') return learning.artifactDrafts.length;
    if (activeLensWorkTab === 'attachments') return learning.attachmentDrafts.filter((attachment) => attachment.sourceContext === 'work_attachment').length;
    return learning.hasSavedJournalNote ? 1 : 0;
  })();
  const observationStorageKey = learning.pointID ? 'lw:point-observation-opened:' + learning.pointID : '';
  const isPointCompleted = learning.pointDetail?.point.status === 'completed';
  const observationNotes = learning.pointDetail?.observation_notes ?? [];
  const hasSavedObservationNote = observationNotes.some((note) => note.content.trim().length > 0);
  const canEditObservationNotes = isExplorationPoint && learning.routeKind === 'learning' && !isPointCompleted;
  const isObservationGateComplete = isPointCompleted || (!isExplorationPoint
    ? Boolean(learning.pointDetail?.research_material_confirmed)
    : hasSavedObservationNote);
  const activeStageIndex = Math.max(0, pointLearningStages.indexOf(activeStage));
  const effectiveVisitedStageIndex = Math.max(maxVisitedStageIndex, activeStageIndex);
  const checkedStages = useMemo(() => {
    const entries = pointLearningStages.map((stage, index) => [
      stage,
      index === 0 ? isObservationGateComplete : effectiveVisitedStageIndex >= index,
    ] as const);
    return Object.fromEntries(entries) as Record<PointLearningStage, boolean>;
  }, [effectiveVisitedStageIndex, isObservationGateComplete]);
  const unlockedStages = useMemo(() => {
    const entries = pointLearningStages.map((stage, index) => [
      stage,
      index === 0 || checkedStages[pointLearningStages[index - 1]],
    ] as const);
    return Object.fromEntries(entries) as Record<PointLearningStage, boolean>;
  }, [checkedStages]);
  useEffect(() => {
    if (unlockedStages[activeStage]) return;
    const activeIndex = pointLearningStages.indexOf(activeStage);
    const fallbackStage = [...pointLearningStages]
      .slice(0, Math.max(1, activeIndex))
      .reverse()
      .find((stage) => unlockedStages[stage]) ?? 'point';
    router.replace(getStageHref(fallbackStage));
  }, [activeStage, getStageHref, router, unlockedStages]);

  useEffect(() => {
    onStageChange?.(activeStage);
    if (activeStage === 'complete') learning.setIsCompletionOpen(true);
    else learning.setIsCompletionOpen(false);
  }, [activeStage, learning.setIsCompletionOpen, onStageChange]);

  const getWorkTabLabel = (tab: PointWorkTab) => {
    if (tab === 'notes') return copy.workspace.learningWork.notesLabel;
    return copy.workspace.learningWork.tabs[tab]?.label ?? tab;
  };
  const observationSpeech = !isExplorationPoint
    ? copy.skillFlow.observe.researchSpeech
    : learning.pointDetail?.point.external_url?.trim()
      ? copy.skillFlow.observe.externalSpeech
      : copy.skillFlow.observe.internalSpeech;
  const observationLumiSize = isCompactObservationLayout ? 104 : 148;
  const observationBubbleStyle = {
    position: 'relative' as const,
    flex: isCompactObservationLayout ? '1 1 0' : '0 1 340px',
    maxWidth: isCompactObservationLayout ? 'min(220px, calc(100vw - 174px))' : '340px',
    minWidth: isCompactObservationLayout ? '0' : '220px',
    padding: isCompactObservationLayout ? '11px 12px' : '14px 16px',
    borderRadius: '14px',
    border: `1px solid ${themeTokens.noticeBorder}`,
    background: themeTokens.noticeBackground,
    color: themeTokens.title,
    boxShadow: '0 14px 28px rgba(15, 23, 42, 0.12)',
    textAlign: 'left' as const,
  };
  const observationBubbleArrowStyle = {
    position: 'absolute' as const,
    left: '-9px',
    top: isCompactObservationLayout ? '28px' : '34px',
    width: '16px',
    height: '16px',
    transform: 'rotate(45deg)',
    borderLeft: `1px solid ${themeTokens.noticeBorder}`,
    borderBottom: `1px solid ${themeTokens.noticeBorder}`,
    background: themeTokens.noticeBackground,
  };
  useEffect(() => {
    if (!observationStorageKey || !isExplorationPoint) return;
    setIsObservationSourceOpened(window.localStorage.getItem(observationStorageKey) === '1');
  }, [isExplorationPoint, observationStorageKey]);

  useEffect(() => {
    const stageIndex = pointLearningStages.indexOf(activeStage);
    setMaxVisitedStageIndex((current) => Math.max(current, stageIndex));
  }, [activeStage]);

  useEffect(() => {
    const query = window.matchMedia('(max-width: 520px)');
    const updateLayout = () => setIsCompactObservationLayout(query.matches);
    updateLayout();
    query.addEventListener('change', updateLayout);
    return () => query.removeEventListener('change', updateLayout);
  }, []);

  const setStage = (stage: PointLearningStage) => {
    if (!unlockedStages[stage]) return;
    if (stage !== activeStage) setIsObservationContentOpen(false);
    onStageChange?.(stage);
    if (stage === 'complete') learning.setIsCompletionOpen(true);
    else learning.setIsCompletionOpen(false);
    router.push(getStageHref(stage));
  };

  const goNext = () => {
    const nextStage = getNextPointLearningStage(activeStage);
    if (unlockedStages[nextStage]) setStage(nextStage);
  };
  const goBack = () => setStage(getPreviousPointLearningStage(activeStage));
  const toggleSkill = (skill: string) => {
    setSelectedSkills((current) => (current.includes(skill) ? current.filter((item) => item !== skill) : [...current, skill]));
  };
  const selectLens = (lens: RepeatRecordLens) => {
    setActiveLens(lens);
    learning.setActiveWorkTab(repeatRecordLensToWorkTab[lens]);
  };
  const markObservationSourceOpened = () => {
    setIsObservationSourceOpened(true);
    if (observationStorageKey) window.localStorage.setItem(observationStorageKey, '1');
  };
  const handleObservationAction = () => {
    if (isExplorationPoint) {
      markObservationSourceOpened();
      const url = learning.pointDetail?.point.external_url?.trim();
      if (url) {
        window.open(url, '_blank', 'noopener,noreferrer');
        onSourceContentOpen();
        return;
      }
      setIsObservationContentOpen(true);
      return;
    }
    router.push(researchMaterialHref);
  };


  const resetObservationNoteForm = () => {
    setObservationNoteType('core_summary');
    setObservationNoteContent('');
    setEditingObservationNoteID(null);
  };

  const handleSubmitObservationNote = async () => {
    const content = observationNoteContent.trim();
    if (!content || !canEditObservationNotes || learning.isSavingObservationNote) return;
    const saved = editingObservationNoteID
      ? await learning.handleUpdateObservationNote(editingObservationNoteID, observationNoteType, content)
      : await learning.handleCreateObservationNote(observationNoteType, content);
    if (saved) resetObservationNoteForm();
  };

  const handleEditObservationNote = (note: NonNullable<NonNullable<typeof learning.pointDetail>['observation_notes']>[number]) => {
    setObservationNoteType(note.note_type);
    setObservationNoteContent(note.content);
    setEditingObservationNoteID(note.id);
  };

  const renderObservationNotes = () => {
    if (!isExplorationPoint) return null;
    const observeCopy = copy.skillFlow.observe;
    return (
      <section style={{ width: 'min(860px, 100%)', display: 'grid', gap: '14px', padding: isCompactObservationLayout ? '14px' : '18px', borderRadius: '10px', border: `1px solid ${themeTokens.sectionBorder}`, background: themeTokens.sectionBackground, textAlign: 'left' }}>
        <div style={{ display: 'grid', gap: '4px' }}>
          <strong style={{ color: themeTokens.title, fontSize: '17px' }}>{observeCopy.notePanelTitle}</strong>
          <p style={{ margin: 0, color: themeTokens.mutedText, fontSize: '14px', lineHeight: 1.55 }}>{observeCopy.notePanelDescription}</p>
        </div>
        {learning.pointEntryMessage ? (
          <div style={{ padding: '12px', borderRadius: '8px', border: `1px solid ${themeTokens.noticeBorder}`, background: themeTokens.noticeBackground, color: themeTokens.description, fontSize: '14px', lineHeight: 1.55 }}>
            {learning.pointEntryMessage}
          </div>
        ) : null}
        {canEditObservationNotes ? (
          <div style={{ display: 'grid', gap: '10px' }}>
            <label style={{ display: 'grid', gap: '6px', color: themeTokens.title, fontWeight: 800, fontSize: '13px' }}>
              {observeCopy.noteTypeLabel}
              <select
                value={observationNoteType}
                onChange={(event) => setObservationNoteType(event.target.value as ObservationNoteType)}
                disabled={learning.isSavingObservationNote}
                style={{ minHeight: '42px', borderRadius: '8px', border: `1px solid ${themeTokens.inputBorder}`, background: themeTokens.inputBackground, color: themeTokens.inputText, padding: '0 10px', font: 'inherit' }}
              >
                {observationNoteTypes.map((type) => (
                  <option key={type} value={type}>{observeCopy.noteTypes[type].label}</option>
                ))}
              </select>
            </label>
            <label style={{ display: 'grid', gap: '6px', color: themeTokens.title, fontWeight: 800, fontSize: '13px' }}>
              {observeCopy.noteContentLabel}
              <textarea
                value={observationNoteContent}
                onChange={(event) => setObservationNoteContent(event.target.value)}
                placeholder={observeCopy.noteTypes[observationNoteType].placeholder}
                disabled={learning.isSavingObservationNote}
                rows={3}
                style={{ width: '100%', resize: 'vertical', borderRadius: '8px', border: `1px solid ${themeTokens.inputBorder}`, background: themeTokens.inputBackground, color: themeTokens.inputText, padding: '12px', font: 'inherit', lineHeight: 1.5, boxSizing: 'border-box' }}
              />
            </label>
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '8px', flexWrap: 'wrap' }}>
              {editingObservationNoteID ? (
                <button type="button" onClick={resetObservationNoteForm} disabled={learning.isSavingObservationNote} style={{ ...secondaryButtonStyle, minHeight: '38px', background: themeTokens.secondaryButtonBackground, borderColor: themeTokens.secondaryButtonBorder, color: themeTokens.buttonText }}>
                  {observeCopy.noteCancel}
                </button>
              ) : null}
              <button type="button" onClick={() => void handleSubmitObservationNote()} disabled={!observationNoteContent.trim() || learning.isSavingObservationNote} style={{ ...primaryButtonStyle, minHeight: '38px', background: observationNoteContent.trim() && !learning.isSavingObservationNote ? themeTokens.primaryButtonBackground : themeTokens.disabledButtonBackground, borderColor: observationNoteContent.trim() && !learning.isSavingObservationNote ? themeTokens.primaryButtonBorder : themeTokens.disabledButtonBorder, color: observationNoteContent.trim() && !learning.isSavingObservationNote ? themeTokens.buttonText : themeTokens.disabledButtonText, cursor: learning.isSavingObservationNote ? 'progress' : observationNoteContent.trim() ? 'pointer' : 'not-allowed' }}>
                {editingObservationNoteID ? observeCopy.noteSave : observeCopy.noteAdd}
              </button>
            </div>
          </div>
        ) : null}
        <div style={{ display: 'grid', gap: '12px' }}>
          {observationNotes.length === 0 ? (
            <div style={{ padding: '14px', borderRadius: '8px', border: `1px dashed ${themeTokens.placeholderBorder}`, background: themeTokens.placeholderBackground, color: themeTokens.mutedText, fontSize: '14px', lineHeight: 1.55 }}>
              {observeCopy.noteEmpty}
            </div>
          ) : observationNoteTypes.map((type) => {
            const notes = observationNotes.filter((note) => note.note_type === type);
            if (notes.length === 0) return null;
            return (
              <div key={type} style={{ display: 'grid', gap: '8px' }}>
                <strong style={{ color: themeTokens.title, fontSize: '14px' }}>{observeCopy.noteTypes[type].label}</strong>
                <div style={{ display: 'grid', gap: '8px' }}>
                  {notes.map((note) => (
                    <article key={note.id} style={{ display: 'grid', gap: '8px', padding: '12px', borderRadius: '8px', border: `1px solid ${themeTokens.placeholderBorder}`, background: themeTokens.placeholderBackground }}>
                      <p style={{ margin: 0, whiteSpace: 'pre-wrap', color: themeTokens.metaValue, lineHeight: 1.6 }}>{note.content}</p>
                      {canEditObservationNotes ? (
                        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '6px', flexWrap: 'wrap' }}>
                          <button type="button" onClick={() => handleEditObservationNote(note)} disabled={learning.isSavingObservationNote} style={{ ...secondaryButtonStyle, minHeight: '32px', padding: '0 10px', fontSize: '12px', background: themeTokens.secondaryButtonBackground, borderColor: themeTokens.secondaryButtonBorder, color: themeTokens.buttonText }}>{observeCopy.noteEdit}</button>
                          <button type="button" onClick={() => void learning.handleDeleteObservationNote(note.id)} disabled={learning.isSavingObservationNote} style={{ ...secondaryButtonStyle, minHeight: '32px', padding: '0 10px', fontSize: '12px', background: themeTokens.secondaryButtonBackground, borderColor: themeTokens.secondaryButtonBorder, color: themeTokens.buttonText }}>{observeCopy.noteDelete}</button>
                        </div>
                      ) : null}
                    </article>
                  ))}
                </div>
              </div>
            );
          })}
        </div>
      </section>
    );
  };

  const bottomActionBaseStyle = {
    ...secondaryButtonStyle,
    minHeight: '46px',
    minWidth: 'min(190px, 48%)',
    display: 'inline-flex',
    alignItems: 'center',
    justifyContent: 'center',
    gap: '8px',
    textDecoration: 'none',
    background: themeTokens.secondaryButtonBackground,
    borderColor: themeTokens.secondaryButtonBorder,
    color: themeTokens.buttonText,
  };
  const renderObservationNavigation = () => (
    <div className="lw-point-flow-fixed-bottom-actions" style={{ position: 'fixed', left: '50%', bottom: '14px', zIndex: 70, width: 'min(1120px, calc(100vw - 24px))', transform: 'translateX(-50%)', display: 'flex', justifyContent: 'space-between', gap: '10px', alignItems: 'center', padding: '10px', borderRadius: '12px', border: `1px solid ${themeTokens.sectionBorder}`, background: themeTokens.sectionBackground, boxShadow: '0 16px 36px rgba(15, 23, 42, 0.22)', backdropFilter: 'blur(14px)' }}>
      <Link href={diaryHref} style={bottomActionBaseStyle}>
        ⬅️ 탐험일지
      </Link>
      <button type="button" disabled={!unlockedStages.goal} onClick={() => setStage('goal')} style={{ ...bottomActionBaseStyle, background: unlockedStages.goal ? '#EF9F27' : themeTokens.disabledButtonBackground, borderColor: unlockedStages.goal ? 'rgba(239, 159, 39, 0.76)' : themeTokens.disabledButtonBorder, color: unlockedStages.goal ? '#241505' : themeTokens.disabledButtonText, cursor: unlockedStages.goal ? 'pointer' : 'not-allowed', opacity: unlockedStages.goal ? 1 : 0.72 }}>
        목표물 발견 ➡️
      </button>
    </div>
  );

  const renderStageEntryMessage = () => {
    if (!learning.pointEntryMessage || activeStage === 'point') return null;
    return (
      <div style={{ ...noticeCardStyle, background: themeTokens.noticeBackground, borderColor: themeTokens.noticeBorder, color: themeTokens.description }}>
        {learning.pointEntryMessage}
      </div>
    );
  };

  const renderNavigation = () => (
    <div className="lw-point-flow-bottom-actions" style={{ display: 'flex', justifyContent: 'space-between', gap: '10px', alignItems: 'center', flexWrap: 'wrap' }}>
      <button type="button" disabled={activeStage === 'point'} onClick={goBack} style={{ ...secondaryButtonStyle, background: themeTokens.secondaryButtonBackground, borderColor: themeTokens.secondaryButtonBorder, color: activeStage === 'point' ? themeTokens.disabledButtonText : themeTokens.buttonText }}>
        {copy.skillFlow.back}
      </button>
      <button type="button" onClick={goNext} disabled={activeStage === 'complete' || !unlockedStages[getNextPointLearningStage(activeStage)]} style={{ ...primaryButtonStyle, background: activeStage === 'complete' || !unlockedStages[getNextPointLearningStage(activeStage)] ? themeTokens.disabledButtonBackground : themeTokens.primaryButtonBackground, borderColor: activeStage === 'complete' || !unlockedStages[getNextPointLearningStage(activeStage)] ? themeTokens.disabledButtonBorder : themeTokens.primaryButtonBorder, color: activeStage === 'complete' || !unlockedStages[getNextPointLearningStage(activeStage)] ? themeTokens.disabledButtonText : themeTokens.buttonText }}>
        {activeStage === 'final_result' ? copy.skillFlow.completePoint : copy.skillFlow.next}
      </button>
    </div>
  );

  return (
    <div ref={flowRootRef} className="lw-point-flow-root" style={{ display: 'grid', gap: '18px', maxWidth: '1120px', margin: '0 auto', width: '100%', paddingBottom: activeStage === 'point' && !shouldHideObservationChrome ? '86px' : 0, scrollMarginTop: '340px' }}>
      <section style={{ display: 'grid', gap: '14px', padding: activeStage === 'point' ? 0 : '18px', borderRadius: activeStage === 'point' ? 0 : '10px', border: activeStage === 'point' ? 'none' : `1px solid ${themeTokens.sectionBorder}`, background: activeStage === 'point' ? 'transparent' : themeTokens.sectionBackground }}>
        <PointFlowProgress copy={copy.skillFlow} activeStage={activeStage} themeTokens={themeTokens} checkedStages={checkedStages} unlockedStages={unlockedStages} onStageSelect={setStage} />
        {activeStage !== 'point' ? (
          <div style={{ display: 'grid', gap: '4px' }}>
            <h1 style={{ margin: 0, fontSize: 'clamp(22px, 4vw, 30px)', lineHeight: 1.2, color: themeTokens.title }}>{stageCopy.title}</h1>
            <p style={{ margin: 0, color: themeTokens.mutedText, lineHeight: 1.55 }}>{stageCopy.subtitle}</p>
          </div>
        ) : null}
      </section>
      {renderStageEntryMessage()}

      {activeStage === 'point' ? (
        <>
          {!shouldHideObservationChrome ? (
            <div style={{ display: 'grid', justifyItems: 'center', alignContent: 'center', gap: isCompactObservationLayout ? '18px' : '24px', minHeight: 'clamp(300px, calc(100dvh - 430px), 520px)', padding: isCompactObservationLayout ? '24px 0' : '32px 0', boxSizing: 'border-box', textAlign: 'center' }}>
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', flexWrap: isCompactObservationLayout ? 'nowrap' : 'wrap', gap: isCompactObservationLayout ? '10px' : '18px', maxWidth: 'min(620px, 100%)', width: '100%' }}>
                <LumiAvatar state={isExplorationPoint ? 'curious' : 'focus'} size={observationLumiSize} />
                <div style={observationBubbleStyle}>
                  <span aria-hidden="true" style={observationBubbleArrowStyle} />
                  <strong style={{ display: 'block', marginBottom: '5px', color: themeTokens.title, fontSize: isCompactObservationLayout ? '13px' : '14px', lineHeight: 1.35 }}>Lumi</strong>
                  <span style={{ display: 'block', color: themeTokens.mutedText, fontSize: isCompactObservationLayout ? '12.5px' : '13px', lineHeight: 1.5 }}>{observationSpeech}</span>
                </div>
              </div>
              <button
                type="button"
                onClick={handleObservationAction}
                style={{
                  ...primaryButtonStyle,
                  width: 'min(280px, 100%)',
                  minHeight: '56px',
                  justifySelf: 'center',
                  background: isExplorationPoint ? '#378ADD' : '#7F77DD',
                  borderColor: isExplorationPoint ? 'rgba(55, 138, 221, 0.76)' : 'rgba(127, 119, 221, 0.76)',
                  color: '#FFFDF7',
                  boxShadow: isExplorationPoint ? '0 10px 20px rgba(55, 138, 221, 0.18)' : '0 10px 20px rgba(127, 119, 221, 0.18)',
                }}
              >
                {isExplorationPoint ? '자료 열기' : '연구자료 만들기'}
              </button>
              {renderObservationNotes()}
            </div>
          ) : null}
          {isObservationContentOpen ? (
            <PointWorkspace
              copy={copy.workspace}
              learning={learning}
              themeTokens={themeTokens}
              isExplorationPoint={isExplorationPoint}
              isSharedRoute={isSharedRoute}
              isLearningContentOpen
              isLearningWorkOpen={false}
              editorToolbarStickyTop={editorToolbarStickyTop}
              onResearchMaterialConfirmed={onResearchMaterialConfirmed}
              onLearningWorkGuideChange={onLearningWorkGuideChange}
              onSourceContentOpen={onSourceContentOpen}
              researchEditModeSignal={researchEditModeSignal}
            />
          ) : null}
          {!shouldHideObservationChrome ? renderObservationNavigation() : null}
        </>
      ) : null}


      {activeStage === 'goal' ? (
        <>
          <NodeSection id="point-flow-goal" title={copy.skillFlow.stages.goal.title} subtitle={copy.skillFlow.stages.goal.subtitle} tokens={themeTokens}>
            <div style={{ display: 'grid', gap: '14px' }}>
              <div style={{ padding: '14px', borderRadius: '8px', border: `1px solid ${themeTokens.noticeBorder}`, background: themeTokens.noticeBackground }}>
                <strong style={{ display: 'block', color: themeTokens.title, marginBottom: '6px' }}>{copy.skillFlow.goal.contextLabel}</strong>
                <p style={{ margin: 0, color: themeTokens.mutedText, lineHeight: 1.55 }}>{pointGoalContextText}</p>
              </div>
              <label style={{ display: 'grid', gap: '8px', color: themeTokens.title, fontWeight: 800 }}>
                {copy.skillFlow.goal.titleLabel}
                <textarea
                  value={goalObjectText}
                  onChange={(event) => setGoalObjectText(event.target.value)}
                  placeholder={copy.skillFlow.goal.titlePlaceholder}
                  rows={3}
                  style={{ width: '100%', resize: 'vertical', borderRadius: '8px', border: `1px solid ${themeTokens.inputBorder}`, background: themeTokens.inputBackground, color: themeTokens.inputText, padding: '12px', font: 'inherit', lineHeight: 1.5 }}
                />
              </label>
              <div style={{ display: 'grid', gap: '8px' }}>
                <strong style={{ color: themeTokens.title }}>{copy.skillFlow.goal.typeLabel}</strong>
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '8px' }}>
                  {goalObjectTypes.map((type) => {
                    const active = goalType === type;
                    return (
                      <button key={type} type="button" onClick={() => setGoalType(type)} style={{ display: 'grid', gap: '5px', minHeight: '88px', padding: '12px', borderRadius: '8px', border: `1px solid ${active ? themeTokens.primaryButtonBorder : themeTokens.secondaryButtonBorder}`, background: active ? themeTokens.primaryButtonBackground : themeTokens.secondaryButtonBackground, color: themeTokens.buttonText, textAlign: 'left', font: 'inherit' }}>
                        <strong>{copy.skillFlow.goal.types[type].label}</strong>
                        <span style={{ color: themeTokens.mutedText, fontSize: '13px', lineHeight: 1.45 }}>{copy.skillFlow.goal.types[type].description}</span>
                      </button>
                    );
                  })}
                </div>
              </div>
            </div>
          </NodeSection>
          {renderNavigation()}
        </>
      ) : null}

      {activeStage === 'skill_discovery' ? (
        <>
          <NodeSection id="point-flow-skills" title={copy.skillFlow.stages.skill_discovery.title} subtitle={copy.skillFlow.skills.prompt} tokens={themeTokens}>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
              {suggestedSkills.map((skill) => {
                const active = selectedSkills.includes(skill);
                return (
                  <button key={skill} type="button" onClick={() => toggleSkill(skill)} style={{ ...secondaryButtonStyle, minHeight: '42px', borderRadius: '8px', background: active ? themeTokens.primaryButtonBackground : themeTokens.secondaryButtonBackground, borderColor: active ? themeTokens.primaryButtonBorder : themeTokens.secondaryButtonBorder, color: themeTokens.buttonText }}>
                    {skill}
                  </button>
                );
              })}
            </div>
            <p style={{ margin: 0, color: themeTokens.mutedText, lineHeight: 1.55 }}>{goalObjectText.trim() || copy.skillFlow.skills.emptyGoalHint}</p>
          </NodeSection>
          {renderNavigation()}
        </>
      ) : null}

      {activeStage === 'goal_iteration' ? (
        <>
          <NodeSection id="point-flow-repeat" title={copy.skillFlow.stages.goal_iteration.title} subtitle={copy.skillFlow.repeat.requiredHint} tokens={themeTokens}>
            <div style={{ display: 'grid', gap: '10px' }}>
              <strong style={{ color: themeTokens.title }}>{copy.skillFlow.repeat.lensLabel}</strong>
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(160px, 1fr))', gap: '8px' }}>
                {(Object.keys(copy.skillFlow.repeat.lenses) as RepeatRecordLens[]).map((lens) => {
                  const active = activeLens === lens;
                  const workTab = repeatRecordLensToWorkTab[lens];
                  const recordCount = workTab === 'questions'
                    ? learning.pointQuestions.length
                    : workTab === 'practice'
                      ? learning.practiceLogDrafts.length
                      : workTab === 'artifacts'
                        ? learning.artifactDrafts.length
                        : workTab === 'attachments'
                          ? learning.attachmentDrafts.filter((attachment) => attachment.sourceContext === 'work_attachment').length
                          : learning.hasSavedJournalNote ? 1 : 0;
                  return (
                    <button key={lens} type="button" onClick={() => selectLens(lens)} style={{ minHeight: '84px', display: 'grid', gap: '5px', alignContent: 'start', padding: '12px', borderRadius: '8px', border: `1px solid ${active ? themeTokens.primaryButtonBorder : themeTokens.secondaryButtonBorder}`, background: active ? themeTokens.primaryButtonBackground : themeTokens.secondaryButtonBackground, color: themeTokens.buttonText, textAlign: 'left', font: 'inherit' }}>
                      <span style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: '8px' }}>
                        <strong>{copy.skillFlow.repeat.lenses[lens].label}</strong>
                        <span style={{ flex: '0 0 auto', borderRadius: '999px', border: `1px solid ${active ? themeTokens.primaryButtonBorder : themeTokens.surfaceBorder}`, padding: '2px 7px', color: themeTokens.mutedText, fontSize: '11px', fontWeight: 850 }}>
                          {copy.skillFlow.repeat.savedCountLabel(recordCount)}
                        </span>
                      </span>
                      <span style={{ color: themeTokens.mutedText, fontSize: '13px', lineHeight: 1.4 }}>{copy.skillFlow.repeat.lenses[lens].description}</span>
                    </button>
                  );
                })}
              </div>
              <div style={{ display: 'grid', gap: '12px', padding: '14px', borderRadius: '8px', border: `1px solid ${themeTokens.noticeBorder}`, background: themeTokens.noticeBackground }}>
                <div style={{ display: 'grid', gap: '5px' }}>
                  <span style={{ color: themeTokens.metaLabel, fontSize: '12px', fontWeight: 850 }}>{copy.skillFlow.repeat.selectedLensLabel}</span>
                  <strong style={{ color: themeTokens.title, fontSize: '16px', lineHeight: 1.35 }}>{activeLensCopy.label}</strong>
                  <span style={{ color: themeTokens.mutedText, fontSize: '13px', lineHeight: 1.5 }}>{activeLensCopy.description}</span>
                </div>
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '10px' }}>
                  <div style={{ display: 'grid', gap: '4px' }}>
                    <span style={{ color: themeTokens.metaLabel, fontSize: '12px', fontWeight: 850 }}>{copy.skillFlow.repeat.connectedRecordLabel}</span>
                    <strong style={{ color: themeTokens.title, fontSize: '14px' }}>{getWorkTabLabel(activeLensWorkTab)}</strong>
                    <span style={{ color: themeTokens.mutedText, fontSize: '12px', fontWeight: 800 }}>{copy.skillFlow.repeat.savedCountLabel(activeLensRecordCount)}</span>
                  </div>
                  <div style={{ display: 'grid', gap: '4px' }}>
                    <span style={{ color: themeTokens.metaLabel, fontSize: '12px', fontWeight: 850 }}>{copy.skillFlow.repeat.goalObjectLabel}</span>
                    <strong style={{ color: themeTokens.title, fontSize: '14px', lineHeight: 1.45 }}>{goalObjectText.trim() || copy.skillFlow.repeat.emptyGoalObject}</strong>
                  </div>
                </div>
                <div style={{ display: 'grid', gap: '6px' }}>
                  <span style={{ color: themeTokens.metaLabel, fontSize: '12px', fontWeight: 850 }}>{copy.skillFlow.repeat.selectedSkillsLabel}</span>
                  {selectedSkills.length ? (
                    <div style={{ display: 'flex', flexWrap: 'wrap', gap: '6px' }}>
                      {selectedSkills.map((skill) => (
                        <span key={skill} style={{ borderRadius: '999px', border: `1px solid ${themeTokens.surfaceBorder}`, background: themeTokens.surfaceBackground, color: themeTokens.title, padding: '5px 9px', fontSize: '12px', fontWeight: 850 }}>
                          {skill}
                        </span>
                      ))}
                    </div>
                  ) : (
                    <span style={{ color: themeTokens.mutedText, fontSize: '13px' }}>{copy.skillFlow.repeat.noSelectedSkills}</span>
                  )}
                </div>
              </div>
            </div>
          </NodeSection>
          <PointWorkspace
            copy={copy.workspace}
            learning={learning}
            themeTokens={themeTokens}
            isExplorationPoint={isExplorationPoint}
            isSharedRoute={isSharedRoute}
            isLearningContentOpen={false}
            isLearningWorkOpen
            editorToolbarStickyTop={editorToolbarStickyTop}
            onResearchMaterialConfirmed={onResearchMaterialConfirmed}
            onLearningWorkGuideChange={onLearningWorkGuideChange}
            onSourceContentOpen={onSourceContentOpen}
          />
          {renderNavigation()}
        </>
      ) : null}

      {activeStage === 'final_result' ? (
        <>
          <NodeSection id="point-flow-final-result" title={copy.skillFlow.stages.final_result.title} subtitle={copy.skillFlow.finalResult.artifactHint} tokens={themeTokens}>
            <button type="button" onClick={() => { selectLens('evidence'); setStage('goal_iteration'); }} style={{ ...primaryButtonStyle, width: 'fit-content', background: themeTokens.primaryButtonBackground, borderColor: themeTokens.primaryButtonBorder, color: themeTokens.buttonText }}>
              {copy.skillFlow.finalResult.openArtifacts}
            </button>
          </NodeSection>
          {renderNavigation()}
        </>
      ) : null}

      {activeStage === 'complete' ? (
        <>
          <NodeSection id="point-flow-complete-hint" title={copy.skillFlow.stages.complete.title} subtitle={copy.skillFlow.complete.readinessHint} tokens={themeTokens}>
            <span style={{ color: themeTokens.mutedText }}>{copy.skillFlow.stages.complete.subtitle}</span>
          </NodeSection>
          <PointCompletion copy={copy} learning={learning} themeTokens={themeTokens} isSharedRoute={isSharedRoute} />
          {renderNavigation()}
        </>
      ) : null}
    </div>
  );
}
