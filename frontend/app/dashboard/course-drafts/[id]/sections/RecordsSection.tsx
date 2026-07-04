'use client';

import type { DraftEditorState } from '../useDraftEditor';
import type { DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft';
import {
  diaryPageEyebrowStyle,
  diaryPageHeaderStyle,
  diaryPageTitleStyle,
  journalActionRowStyle,
  journalAnswerTextareaStyle,
  journalContextCardStyle,
  journalContextDescriptionStyle,
  journalContextGridStyle,
  journalContextTitleStyle,
  journalNoticeStyle,
  journalPrimaryButtonActiveStyle,
  journalPrimaryButtonStyle,
  journalScaffoldStyle,
  journalSecondaryButtonDisabledStyle,
  journalSecondaryButtonStyle,
  pageStatusPillStyle,
  plannerStatLabelStyle,
  recordApplicationFieldStyle,
  recordMetricCardStyle,
  recordMetricGridStyle,
  recordMetricInputStyle,
  recordMetricLabelStyle,
  recordMetricSelectStyle,
} from '../styles';

export function RecordsSection({ editor, copy }: { editor: DraftEditorState; copy: DashboardCourseDraftCopy['supportSections'] }) {
  const {
    planningSummary,
    selectedPointContext,
    hasRecordDraftChanges,
    recordStudyMinutesInput,
    setRecordStudyMinutesInput,
    recordPracticeCountInput,
    setRecordPracticeCountInput,
    recordConfidenceLevelInput,
    setRecordConfidenceLevelInput,
    recordApplicationNoteInput,
    setRecordApplicationNoteInput,
    isSavingRecordEntry,
    canSaveRecordEntry,
    recordSaveMessage,
    handleResetRecordDraft,
    handleSaveRecordEntry,
  } = editor;

  return (
    <section style={journalScaffoldStyle}>
      <div style={diaryPageHeaderStyle}>
        <div>
          <div style={diaryPageEyebrowStyle}>{copy.records.eyebrow}</div>
          <h3 style={diaryPageTitleStyle}>{copy.records.title}</h3>
        </div>
        <div style={pageStatusPillStyle}>
          {planningSummary.selectedLesson
            ? hasRecordDraftChanges
              ? copy.records.statusChanged
              : copy.records.statusConnected
            : copy.records.statusNeedRegion}
        </div>
      </div>

      <div style={journalContextGridStyle}>
        <article style={journalContextCardStyle}>
          <span style={plannerStatLabelStyle}>{copy.common.currentRegion}</span>
          <strong style={journalContextTitleStyle}>{planningSummary.selectedLesson?.lesson.title ?? copy.common.regionNotSelected}</strong>
          <p style={journalContextDescriptionStyle}>
            {planningSummary.selectedLesson?.lesson.summary ?? copy.records.regionDescriptionFallback}
          </p>
        </article>
        <article style={journalContextCardStyle}>
          <span style={plannerStatLabelStyle}>{copy.common.currentPoint}</span>
          <strong style={journalContextTitleStyle}>{selectedPointContext?.title ?? copy.common.pointNotSelected}</strong>
          <p style={journalContextDescriptionStyle}>
            {selectedPointContext?.description || copy.records.pointDescriptionFallback}
          </p>
        </article>
      </div>

      <div style={recordMetricGridStyle}>
        <label style={recordMetricCardStyle}>
          <span style={recordMetricLabelStyle}>{copy.records.focusMinutes}</span>
          <input
            type="number"
            min={0}
            value={recordStudyMinutesInput}
            onChange={(event) => setRecordStudyMinutesInput(event.target.value)}
            disabled={!planningSummary.selectedLesson || isSavingRecordEntry}
            style={recordMetricInputStyle}
          />
        </label>
        <label style={recordMetricCardStyle}>
          <span style={recordMetricLabelStyle}>{copy.records.practiceCount}</span>
          <input
            type="number"
            min={0}
            value={recordPracticeCountInput}
            onChange={(event) => setRecordPracticeCountInput(event.target.value)}
            disabled={!planningSummary.selectedLesson || isSavingRecordEntry}
            style={recordMetricInputStyle}
          />
        </label>
        <label style={recordMetricCardStyle}>
          <span style={recordMetricLabelStyle}>{copy.records.confidence}</span>
          <select
            value={recordConfidenceLevelInput}
            onChange={(event) => setRecordConfidenceLevelInput(Number(event.target.value))}
            disabled={!planningSummary.selectedLesson || isSavingRecordEntry}
            style={recordMetricSelectStyle}
          >
            {[1, 2, 3, 4, 5].map((level) => (
              <option key={level} value={level}>{`${level}/5`}</option>
            ))}
          </select>
        </label>
      </div>

      <label style={recordApplicationFieldStyle}>
        <span style={recordMetricLabelStyle}>{copy.records.applicationNote}</span>
        <textarea
          value={recordApplicationNoteInput}
          onChange={(event) => setRecordApplicationNoteInput(event.target.value)}
          rows={5}
          maxLength={600}
          placeholder={planningSummary.selectedLesson ? copy.records.applicationPlaceholder : copy.records.applicationDisabledPlaceholder}
          disabled={!planningSummary.selectedLesson || isSavingRecordEntry}
          style={journalAnswerTextareaStyle}
        />
      </label>

      <div style={journalActionRowStyle}>
        <button
          type="button"
          onClick={handleResetRecordDraft}
          disabled={!hasRecordDraftChanges || isSavingRecordEntry}
          style={{
            ...journalSecondaryButtonStyle,
            ...((!hasRecordDraftChanges || isSavingRecordEntry) ? journalSecondaryButtonDisabledStyle : null),
          }}
        >
          {copy.common.reset}
        </button>
        <button
          type="button"
          onClick={handleSaveRecordEntry}
          disabled={!canSaveRecordEntry}
          style={{
            ...journalPrimaryButtonStyle,
            ...(canSaveRecordEntry ? journalPrimaryButtonActiveStyle : null),
          }}
        >
          {isSavingRecordEntry ? copy.records.saving : copy.records.save}
        </button>
      </div>

      <div style={journalNoticeStyle}>
        {recordSaveMessage
          ?? (planningSummary.selectedLesson
            ? copy.records.noticeReady
            : copy.records.noticeNeedRegion)}
      </div>
    </section>
  );
}
