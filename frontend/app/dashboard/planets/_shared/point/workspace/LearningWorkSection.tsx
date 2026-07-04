'use client';

import { useEffect, useRef, useState } from 'react';
import LumiModalShell from '@/components/common/LumiModalShell';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../../pointPageUtils';
import {
  sectionStyle, sectionTitleStyle, sectionSubtitleStyle,
  pointTabButtonStyle,
  noticeCardStyle, formGridStyle, blockCardStyle, disclosureSummaryStyle,
} from '../../pointPageStyles';
import type { UsePointLearningResult } from '../usePointLearning';
import { makeWorkspaceStyles } from './workspaceStyles';
import { pointWorkTabs, getLearningWorkGuideDetail, scrollToWorkspaceTarget } from './workspaceTypes';
import type { NoteFieldKey, LearningWorkGuideDetail } from './workspaceTypes';
import { ResearchSectionHeader, WorkspaceDivider, WorkspaceMiniNav } from './WorkspaceSectionHeader';
import PointNotesTab from './PointNotesTab';
import PointQuestionsTab from './PointQuestionsTab';
import PointPracticeTab from './PointPracticeTab';
import PointArtifactsTab from './PointArtifactsTab';
import PointAttachmentsTab from './PointAttachmentsTab';
import PointTimelineTab from './PointTimelineTab';

interface Props {
  copy: PointLearningCopy['workspace']['learningWork'];
  editorCopy: PointLearningCopy['workspace']['editor'];
  learning: UsePointLearningResult;
  themeTokens: PointPageThemeTokens;
  canEditLearningRecords: boolean;
  isSharedRoute: boolean;
  isResearchWorkLocked: boolean;
  isSectionMiniNavSide: boolean;
  editorToolbarStickyTop: string;
  onLearningWorkGuideChange: (guide: LearningWorkGuideDetail) => void;
}

export default function LearningWorkSection({
  copy,
  editorCopy,
  learning,
  themeTokens,
  canEditLearningRecords,
  isSharedRoute,
  isResearchWorkLocked,
  isSectionMiniNavSide,
  editorToolbarStickyTop,
  onLearningWorkGuideChange,
}: Props) {
  const [activeNoteField, setActiveNoteField] = useState<NoteFieldKey>('observation');
  const [focusedNoteGuideField, setFocusedNoteGuideField] = useState<NoteFieldKey | null>(null);
  const [isChoiceRecordOpen, setIsChoiceRecordOpen] = useState(true);
  const [questionsEditing, setQuestionsEditing] = useState(false);
  const [practiceEditing, setPracticeEditing] = useState(false);
  const [artifactsEditing, setArtifactsEditing] = useState(false);
  const [attachmentsEditing, setAttachmentsEditing] = useState(false);
  const [pendingWorkTabChange, setPendingWorkTabChange] = useState<typeof learning.activeWorkTab | null>(null);
  const previousHasSavedJournalNoteRef = useRef<boolean | null>(null);

  const { hasSavedJournalNote, activeWorkTab, setActiveWorkTab, pointQuestions, practiceLogDrafts, artifactDrafts } = learning;
  const isLearningRecordEditing = questionsEditing || practiceEditing || artifactsEditing || attachmentsEditing;
  const hasSavedLearningRecord = pointQuestions.length > 0 || practiceLogDrafts.length > 0 || artifactDrafts.length > 0;

  const { materialToolbarButtonStyle, materialToolbarPrimaryButtonStyle } = makeWorkspaceStyles(themeTokens);
  const learningWorkNarrowCardStyle = {
    width: '100%',
    maxWidth: '920px',
    margin: '0 auto',
  } as const;
  const sectionMiniNavBodyLayoutStyle = isSectionMiniNavSide
    ? {
        width: '100%',
        maxWidth: '1064px',
        margin: '0 auto',
        display: 'grid',
        gridTemplateColumns: '124px minmax(0, 920px)',
        gap: '20px',
        alignItems: 'start',
      } as const
    : { display: 'grid', gap: '14px' } as const;

  useEffect(() => {
    if (!hasSavedJournalNote && activeWorkTab !== 'notes') setActiveWorkTab('notes');
    if (!hasSavedJournalNote) setIsChoiceRecordOpen(false);
  }, [activeWorkTab, hasSavedJournalNote, setActiveWorkTab]);

  useEffect(() => {
    const previous = previousHasSavedJournalNoteRef.current;
    previousHasSavedJournalNoteRef.current = hasSavedJournalNote;
    if (previous === false && hasSavedJournalNote) setIsChoiceRecordOpen(true);
    if (previous === null && hasSavedJournalNote) setIsChoiceRecordOpen(true);
  }, [hasSavedJournalNote]);

  useEffect(() => {
    const guideNoteField = focusedNoteGuideField ?? activeNoteField;
    const guideWorkTab = focusedNoteGuideField ? 'notes' : activeWorkTab;
    onLearningWorkGuideChange(getLearningWorkGuideDetail(guideNoteField, guideWorkTab, copy));
  }, [activeNoteField, activeWorkTab, copy, focusedNoteGuideField, onLearningWorkGuideChange]);

  const scrollToLearningRecordSection = () => scrollToWorkspaceTarget('point-learning-record-section');
  const requestWorkTabChange = (tab: typeof activeWorkTab) => {
    setFocusedNoteGuideField(null);
    if (tab === activeWorkTab) return;
    if (isLearningRecordEditing) { setPendingWorkTabChange(tab); return; }
    setIsChoiceRecordOpen(true);
    setActiveWorkTab(tab);
  };
  const handleConfirmWorkTabChange = () => {
    if (!pendingWorkTabChange) return;
    const nextTab = pendingWorkTabChange;
    setPendingWorkTabChange(null);
    setIsChoiceRecordOpen(true);
    setActiveWorkTab(nextTab);
    scrollToLearningRecordSection();
  };

  return (
    <>
      {pendingWorkTabChange ? (
        <LumiModalShell
          title={copy.moveModalTitle}
          eyebrow={copy.moveModalEyebrow}
          lumiState="curious"
          tone="alert"
          width={440}
          onClose={() => setPendingWorkTabChange(null)}
          message={copy.moveModalMessage}
          actions={(
            <>
              <button type="button" onClick={() => setPendingWorkTabChange(null)} style={materialToolbarButtonStyle}>{copy.keepWriting}</button>
              <button type="button" onClick={handleConfirmWorkTabChange} style={materialToolbarPrimaryButtonStyle}>{copy.move}</button>
            </>
          )}
        />
      ) : null}
      {isResearchWorkLocked ? (
        <section id="point-work-section" style={{ ...sectionStyle, scrollMarginTop: '430px', background: themeTokens.sectionBackground, borderColor: themeTokens.sectionBorder, opacity: 0.78 }}>
          <div style={{ display: 'grid', gap: '8px' }}>
            <h2 style={{ ...sectionTitleStyle, color: themeTokens.title }}>{copy.title}</h2>
            <p style={{ ...sectionSubtitleStyle, color: themeTokens.mutedText }}>{copy.subtitle}</p>
          </div>
          <div style={{ ...noticeCardStyle, background: themeTokens.noticeBackground, borderColor: themeTokens.noticeBorder, color: themeTokens.description }}>
            {copy.lockedNotice}
          </div>
        </section>
      ) : (
        <section id="point-work-section" style={{ ...sectionStyle, scrollMarginTop: '430px', background: themeTokens.sectionBackground, borderColor: themeTokens.sectionBorder }}>
          <div style={{ display: 'grid', gap: '8px' }}>
            <h2 style={{ ...sectionTitleStyle, color: themeTokens.title }}>{copy.title}</h2>
            <p style={{ ...sectionSubtitleStyle, color: themeTokens.mutedText }}>{copy.subtitle}</p>
          </div>
          <div style={sectionMiniNavBodyLayoutStyle}>
            <WorkspaceMiniNav
              ariaLabel={copy.miniNavLabel}
              items={[
                {
                  key: 'notes',
                  label: copy.notesLabel,
                  checked: hasSavedJournalNote,
                  onClick: () => scrollToWorkspaceTarget('point-journal-note-section'),
                },
                {
                  key: 'record',
                  label: copy.recordsLabel,
                  checked: hasSavedLearningRecord,
                  disabled: !hasSavedJournalNote,
                  onClick: () => { setIsChoiceRecordOpen(true); scrollToLearningRecordSection(); },
                },
              ]}
              sticky
              layout={isSectionMiniNavSide ? 'side' : 'top'}
              themeTokens={themeTokens}
            />
            <div style={{ display: 'grid', gap: '14px', minWidth: 0 }}>
              <div style={{ ...learningWorkNarrowCardStyle, display: 'grid', gap: '14px' }}>
                <PointNotesTab
                  copy={copy.notes}
                  learning={learning}
                  themeTokens={themeTokens}
                  canEditLearningRecords={canEditLearningRecords}
                  activeNoteField={activeNoteField}
                  setActiveNoteField={setActiveNoteField}
                  focusedNoteGuideField={focusedNoteGuideField}
                  setFocusedNoteGuideField={setFocusedNoteGuideField}
                />
              </div>

              <WorkspaceDivider themeTokens={themeTokens} />

              <div style={{ ...learningWorkNarrowCardStyle, display: 'grid', gap: '14px' }}>
                {hasSavedJournalNote ? (
                  <details
                    id="point-learning-record-section"
                    open={isChoiceRecordOpen}
                    onToggle={(event) => setIsChoiceRecordOpen(event.currentTarget.open)}
                    style={{ display: 'grid', gap: '12px' }}
                  >
                    <summary style={{ ...disclosureSummaryStyle, color: themeTokens.title, listStyle: 'none' }}>
                      <ResearchSectionHeader
                        title={`${copy.recordsTitle} ${isChoiceRecordOpen ? '▲' : '▼'}`}
                        themeTokens={themeTokens}
                        action={<span style={{ color: themeTokens.mutedText, fontSize: '13px', fontWeight: 700 }}>{copy.optionalLabel}</span>}
                      />
                    </summary>
                    <div style={{ ...blockCardStyle, width: '100%', maxWidth: '860px', margin: '10px auto 0', background: themeTokens.surfaceBackground, borderColor: themeTokens.surfaceBorder }}>
                      <p style={{ margin: 0, color: themeTokens.mutedText, fontSize: '13px', lineHeight: 1.55 }}>
                        {copy.recordsDescription}
                      </p>
                      <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
                        {pointWorkTabs.filter((tab) => tab.key !== 'notes').map((tab) => {
                          const active = activeWorkTab === tab.key;
                          return (
                            <button
                              key={tab.key}
                              type="button"
                              onClick={() => requestWorkTabChange(tab.key)}
                              onFocus={() => setFocusedNoteGuideField(null)}
                              onMouseEnter={() => setFocusedNoteGuideField(null)}
                              style={{ ...pointTabButtonStyle, minHeight: '44px', padding: '8px 13px', gap: '7px', background: active ? themeTokens.primaryButtonBackground : themeTokens.secondaryButtonBackground, borderColor: active ? themeTokens.primaryButtonBorder : themeTokens.secondaryButtonBorder, color: themeTokens.buttonText, boxShadow: active ? '0 8px 18px rgba(37, 99, 235, 0.16)' : '0 4px 10px rgba(15, 23, 42, 0.06)' }}
                            >
                              <span aria-hidden="true" style={{ fontSize: '16px', lineHeight: 1 }}>{tab.icon}</span>
                              <strong>{copy.tabs[tab.key as keyof typeof copy.tabs]?.label ?? tab.label.replace(/^\d+\.\s*/, '')}</strong>
                            </button>
                          );
                        })}
                      </div>
                      {activeWorkTab === 'questions' ? (
                        <PointQuestionsTab
                          copy={copy.recordPanels.questions}
                          commonCopy={copy.recordPanels.common}
                          learning={learning}
                          themeTokens={themeTokens}
                          canEditLearningRecords={canEditLearningRecords}
                          isSharedRoute={isSharedRoute}
                          onFocusClear={() => setFocusedNoteGuideField(null)}
                          onEditingChange={setQuestionsEditing}
                        />
                      ) : null}
                      {activeWorkTab === 'practice' ? (
                        <PointPracticeTab
                          copy={copy.recordPanels.practice}
                          commonCopy={copy.recordPanels.common}
                          learning={learning}
                          themeTokens={themeTokens}
                          canEditLearningRecords={canEditLearningRecords}
                          onEditingChange={setPracticeEditing}
                        />
                      ) : null}
                      {activeWorkTab === 'artifacts' ? (
                        <PointArtifactsTab
                          copy={copy.recordPanels.artifacts}
                          commonCopy={copy.recordPanels.common}
                          editorCopy={editorCopy}
                          learning={learning}
                          themeTokens={themeTokens}
                          canEditLearningRecords={canEditLearningRecords}
                          isSharedRoute={isSharedRoute}
                          editorToolbarStickyTop={editorToolbarStickyTop}
                          onEditingChange={setArtifactsEditing}
                        />
                      ) : null}
                      {activeWorkTab === 'attachments' ? (
                        <PointAttachmentsTab
                          copy={copy.recordPanels.attachments}
                          commonCopy={copy.recordPanels.common}
                          learning={learning}
                          themeTokens={themeTokens}
                          canEditLearningRecords={canEditLearningRecords}
                          onEditingChange={setAttachmentsEditing}
                        />
                      ) : null}
                      {activeWorkTab === 'timeline' ? (
                        <PointTimelineTab
                          copy={copy.recordPanels.timeline}
                          commonCopy={copy.recordPanels.common}
                          learning={learning}
                          themeTokens={themeTokens}
                        />
                      ) : null}
                    </div>
                  </details>
                ) : (
                  <div style={formGridStyle}>
                    <ResearchSectionHeader title={copy.recordsTitle} themeTokens={themeTokens} action={<span style={{ color: themeTokens.mutedText, fontSize: '13px', fontWeight: 700 }}>{copy.optionalLabel}</span>} />
                    <div style={{ ...noticeCardStyle, background: themeTokens.noticeBackground, borderColor: themeTokens.noticeBorder, color: themeTokens.description }}>
                      {copy.recordsLockedNotice}
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        </section>
      )}
    </>
  );
}
