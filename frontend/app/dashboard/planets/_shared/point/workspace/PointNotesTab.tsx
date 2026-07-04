'use client';

import { useState } from 'react';
import LumiModalShell from '@/components/common/LumiModalShell';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../../pointPageUtils';
import type { UsePointLearningResult } from '../usePointLearning';
import {
  formGridStyle, sectionMiniHeaderStyle,
  pointShellTitleStyle, blockMetaStyle,
  blockCardStyle, blockHeaderStyle,
  emptyCardStyle, formActionRowStyle,
  inputStyle, textareaStyle,
  pointTabButtonStyle,
} from '../../pointPageStyles';
import type { NoteFieldKey } from './workspaceTypes';
import { makeWorkspaceStyles } from './workspaceStyles';
import { ResearchSectionHeader } from './WorkspaceSectionHeader';
import { scrollToWorkspaceTarget } from './workspaceTypes';

interface Props {
  copy: PointLearningCopy['workspace']['learningWork']['notes'];
  learning: UsePointLearningResult;
  themeTokens: PointPageThemeTokens;
  canEditLearningRecords: boolean;
  activeNoteField: NoteFieldKey;
  setActiveNoteField: (field: NoteFieldKey) => void;
  focusedNoteGuideField: NoteFieldKey | null;
  setFocusedNoteGuideField: (field: NoteFieldKey | null) => void;
}

const scrollToJournalNoteSection = () => scrollToWorkspaceTarget('point-journal-note-section');

export default function PointNotesTab({ copy, learning, themeTokens, canEditLearningRecords, activeNoteField, setActiveNoteField, focusedNoteGuideField, setFocusedNoteGuideField }: Props) {
  const {
    journalObservation, setJournalObservation,
    journalReflection, setJournalReflection,
    journalExamples, setJournalExamples,
    journalNextStep, setJournalNextStep,
    hasSavedJournalNote, isSavingJournal,
    handleSaveJournal, resetJournalDraftToSaved,
  } = learning;

  const { materialToolbarButtonStyle, materialToolbarPrimaryButtonStyle } = makeWorkspaceStyles(themeTokens);

  const [isEditingJournalNote, setIsEditingJournalNote] = useState(!hasSavedJournalNote);
  const [isConfirmingJournalEditCancel, setIsConfirmingJournalEditCancel] = useState(false);

  const noteFieldOptions: Array<{ key: NoteFieldKey; icon: string; label: string; value: string; placeholder: string }> = [
    { key: 'observation', ...copy.fields.observation, value: journalObservation },
    { key: 'reflection', ...copy.fields.reflection, value: journalReflection },
    { key: 'examples', ...copy.fields.examples, value: journalExamples },
    { key: 'nextStep', ...copy.fields.nextStep, value: journalNextStep },
  ];
  const activeNoteOption = noteFieldOptions.find((o) => o.key === activeNoteField) ?? noteFieldOptions[0];
  const emptyNoteOption = noteFieldOptions.find((option) => !option.value.trim());
  const hasEmptyNoteOption = Boolean(emptyNoteOption);

  const handleChangeActiveNoteField = (value: string) => {
    if (activeNoteField === 'observation') setJournalObservation(value);
    if (activeNoteField === 'reflection') setJournalReflection(value);
    if (activeNoteField === 'examples') setJournalExamples(value);
    if (activeNoteField === 'nextStep') setJournalNextStep(value);
  };

  const handleSaveJournalAndClose = async () => {
    const saved = await handleSaveJournal();
    if (saved) setIsEditingJournalNote(false);
  };

  const handleConfirmJournalEditCancel = () => {
    resetJournalDraftToSaved();
    setIsConfirmingJournalEditCancel(false);
    setIsEditingJournalNote(false);
  };
  const handleAddEmptyNoteField = () => {
    if (!emptyNoteOption) return;
    setActiveNoteField(emptyNoteOption.key);
    setFocusedNoteGuideField(emptyNoteOption.key);
    setIsEditingJournalNote(true);
    scrollToJournalNoteSection();
  };

  return (
    <>
      {isConfirmingJournalEditCancel ? (
        <LumiModalShell
          title={copy.cancelModalTitle}
          eyebrow={copy.cancelModalEyebrow}
          lumiState="curious"
          tone="alert"
          width={440}
          onClose={() => setIsConfirmingJournalEditCancel(false)}
          message={copy.cancelModalMessage}
          actions={(
            <>
              <button type="button" onClick={() => setIsConfirmingJournalEditCancel(false)} style={materialToolbarButtonStyle}>{copy.back}</button>
              <button type="button" onClick={handleConfirmJournalEditCancel} style={materialToolbarButtonStyle}>{copy.cancel}</button>
            </>
          )}
        />
      ) : null}
      <div
        id="point-journal-note-section"
        style={formGridStyle}
        onMouseEnter={() => setFocusedNoteGuideField(activeNoteField)}
        onFocusCapture={() => setFocusedNoteGuideField(activeNoteField)}
      >
        <ResearchSectionHeader
          title={copy.title}
          themeTokens={themeTokens}
          action={
            canEditLearningRecords && !isEditingJournalNote ? (
              <button type="button" onClick={() => { setIsEditingJournalNote(true); scrollToJournalNoteSection(); }} style={{ ...materialToolbarButtonStyle, minHeight: '32px', padding: '0 10px', fontSize: '12px' }}>
                {copy.edit}
              </button>
            ) : null
          }
        />
        {isEditingJournalNote ? (
          <div style={{ display: 'grid', gap: '14px' }}>
            <div style={{ color: themeTokens.mutedText, fontSize: '14px', lineHeight: 1.55 }}>
              {copy.helper}
            </div>
            <div style={{ width: '100%', maxWidth: '860px', margin: '0 auto', display: 'grid', gap: '14px' }}>
              <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
                {noteFieldOptions.map((option) => {
                  const active = activeNoteField === option.key;
                  const hasValue = option.value.trim().length > 0;
                  return (
                    <button
                      key={option.key}
                      type="button"
                      onMouseEnter={() => setFocusedNoteGuideField(option.key)}
                      onFocus={() => setFocusedNoteGuideField(option.key)}
                      onClick={() => { setActiveNoteField(option.key); setFocusedNoteGuideField(option.key); }}
                      style={{
                        ...pointTabButtonStyle,
                        minHeight: '40px',
                        padding: '7px 12px',
                        gap: '7px',
                        background: active ? themeTokens.primaryButtonBackground : themeTokens.secondaryButtonBackground,
                        borderColor: active ? themeTokens.primaryButtonBorder : hasValue ? themeTokens.primaryButtonBorder : themeTokens.secondaryButtonBorder,
                        color: themeTokens.buttonText,
                        boxShadow: active ? '0 8px 18px rgba(37, 99, 235, 0.16)' : hasValue ? '0 5px 12px rgba(20, 184, 166, 0.12)' : '0 4px 10px rgba(15, 23, 42, 0.05)',
                      }}
                    >
                      <span aria-hidden="true" style={{ fontSize: '16px', lineHeight: 1 }}>{hasValue ? '✅' : option.icon}</span>
                      <strong>{option.label}</strong>
                    </button>
                  );
                })}
              </div>
              <textarea
                style={{ ...textareaStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }}
                value={activeNoteOption.value}
                onChange={(e) => handleChangeActiveNoteField(e.target.value)}
                onFocus={() => setFocusedNoteGuideField(activeNoteField)}
                onMouseEnter={() => setFocusedNoteGuideField(activeNoteField)}
                placeholder={activeNoteOption.placeholder}
                readOnly={!canEditLearningRecords}
              />
              {canEditLearningRecords ? (
                <div style={{ ...formActionRowStyle, justifyContent: 'space-between' }}>
                  <button type="button" onClick={() => setIsConfirmingJournalEditCancel(true)} disabled={isSavingJournal} style={{ ...materialToolbarButtonStyle, opacity: isSavingJournal ? 0.62 : 1 }}>{copy.cancelButton}</button>
                  <button type="button" onClick={() => void handleSaveJournalAndClose()} disabled={isSavingJournal} style={{ ...materialToolbarPrimaryButtonStyle, cursor: isSavingJournal ? 'progress' : 'pointer' }}>
                    {isSavingJournal ? copy.saving : copy.save}
                  </button>
                </div>
              ) : null}
            </div>
          </div>
        ) : (
          <div style={{ width: '100%', maxWidth: '860px', margin: '0 auto', display: 'grid', gap: '12px' }}>
            {noteFieldOptions.some((o) => o.value.trim()) ? (
              <div style={{ display: 'grid', gap: '10px' }}>
                <div style={{ ...blockCardStyle, gap: '10px', background: themeTokens.placeholderBackground, borderColor: themeTokens.placeholderBorder }}>
                  <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: '10px', flexWrap: 'wrap' }}>
                    <strong style={{ ...pointShellTitleStyle, color: themeTokens.title, fontSize: '16px' }}>{copy.progressTitle}</strong>
                    {canEditLearningRecords && hasEmptyNoteOption ? (
                      <button
                        type="button"
                        onClick={handleAddEmptyNoteField}
                        style={{ ...materialToolbarButtonStyle, minHeight: '32px', padding: '0 10px', fontSize: '12px' }}
                      >
                        {copy.addEmpty}
                      </button>
                    ) : null}
                  </div>
                  <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
                    {noteFieldOptions.map((option) => {
                      const checked = option.value.trim().length > 0;
                      return (
                        <span
                          key={option.key}
                          style={{
                            display: 'inline-flex',
                            alignItems: 'center',
                            gap: '6px',
                            minHeight: '30px',
                            padding: '0 10px',
                            borderRadius: '999px',
                            border: `1px solid ${checked ? themeTokens.primaryButtonBorder : themeTokens.secondaryButtonBorder}`,
                            background: checked ? themeTokens.primaryButtonBackground : themeTokens.secondaryButtonBackground,
                            color: themeTokens.buttonText,
                            fontSize: '13px',
                            fontWeight: 800,
                            whiteSpace: 'nowrap',
                          }}
                        >
                          <span aria-hidden="true">{checked ? '☑' : '☐'}</span>
                          {option.label}
                        </span>
                      );
                    })}
                  </div>
                  <div style={{ color: themeTokens.mutedText, fontSize: '13px', lineHeight: 1.55 }}>
                    {hasEmptyNoteOption
                      ? copy.progressIncomplete
                      : copy.progressComplete}
                  </div>
                </div>
                {noteFieldOptions.filter((o) => o.value.trim()).map((option) => (
                  <div key={option.key} style={{ ...blockCardStyle, background: themeTokens.surfaceBackground, borderColor: themeTokens.surfaceBorder }}>
                    <div style={blockHeaderStyle}>
                      <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{option.icon} {option.label}</strong>
                    </div>
                    <div style={{ color: themeTokens.description, fontSize: '16px', lineHeight: 1.75, whiteSpace: 'pre-wrap' }}>{option.value}</div>
                  </div>
                ))}
              </div>
            ) : (
              <div style={{ ...emptyCardStyle, background: themeTokens.emptyBackground, borderColor: themeTokens.emptyBorder, color: themeTokens.mutedText }}>{copy.empty}</div>
            )}
            {canEditLearningRecords ? (
              <div style={{ ...formActionRowStyle, justifyContent: 'flex-end' }}>
                <button type="button" onClick={() => { setIsEditingJournalNote(true); scrollToJournalNoteSection(); }} style={materialToolbarButtonStyle}>{copy.edit}</button>
              </div>
            ) : null}
          </div>
        )}
      </div>
    </>
  );
}
