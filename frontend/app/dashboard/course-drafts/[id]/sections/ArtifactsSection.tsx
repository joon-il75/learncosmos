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

export function ArtifactsSection({ editor, copy }: { editor: DraftEditorState; copy: DashboardCourseDraftCopy['supportSections'] }) {
  const {
    planningSummary,
    selectedPointContext,
    hasArtifactDraftChanges,
    artifactTypeInput,
    setArtifactTypeInput,
    artifactTitleInput,
    setArtifactTitleInput,
    artifactUrlInput,
    setArtifactUrlInput,
    artifactDescriptionInput,
    setArtifactDescriptionInput,
    isSavingArtifactEntry,
    canSaveArtifactEntry,
    artifactSaveMessage,
    handleResetArtifactDraft,
    handleSaveArtifactEntry,
  } = editor;

  return (
    <section style={journalScaffoldStyle}>
      <div style={diaryPageHeaderStyle}>
        <div>
          <div style={diaryPageEyebrowStyle}>{copy.artifacts.eyebrow}</div>
          <h3 style={diaryPageTitleStyle}>{copy.artifacts.title}</h3>
        </div>
        <div style={pageStatusPillStyle}>
          {planningSummary.selectedLesson
            ? hasArtifactDraftChanges
              ? copy.artifacts.statusChanged
              : copy.artifacts.statusConnected
            : copy.artifacts.statusNeedRegion}
        </div>
      </div>

      <div style={journalContextGridStyle}>
        <article style={journalContextCardStyle}>
          <span style={plannerStatLabelStyle}>{copy.common.currentRegion}</span>
          <strong style={journalContextTitleStyle}>{planningSummary.selectedLesson?.lesson.title ?? copy.common.regionNotSelected}</strong>
          <p style={journalContextDescriptionStyle}>
            {planningSummary.selectedLesson?.lesson.summary ?? copy.artifacts.regionDescriptionFallback}
          </p>
        </article>
        <article style={journalContextCardStyle}>
          <span style={plannerStatLabelStyle}>{copy.common.currentPoint}</span>
          <strong style={journalContextTitleStyle}>{selectedPointContext?.title ?? copy.common.pointNotSelected}</strong>
          <p style={journalContextDescriptionStyle}>
            {selectedPointContext?.description || copy.artifacts.pointDescriptionFallback}
          </p>
        </article>
      </div>

      <div style={recordMetricGridStyle}>
        <label style={recordMetricCardStyle}>
          <span style={recordMetricLabelStyle}>{copy.artifacts.type}</span>
          <select
            value={artifactTypeInput}
            onChange={(event) => setArtifactTypeInput(event.target.value)}
            disabled={!planningSummary.selectedLesson || isSavingArtifactEntry}
            style={recordMetricSelectStyle}
          >
            {copy.artifacts.typeOptions.map((option) => (
              <option key={option.value} value={option.value}>{option.label}</option>
            ))}
          </select>
        </label>
        <label style={recordMetricCardStyle}>
          <span style={recordMetricLabelStyle}>{copy.artifacts.artifactTitle}</span>
          <input
            value={artifactTitleInput}
            onChange={(event) => setArtifactTitleInput(event.target.value)}
            disabled={!planningSummary.selectedLesson || isSavingArtifactEntry}
            style={recordMetricInputStyle}
            maxLength={160}
          />
        </label>
        <label style={recordMetricCardStyle}>
          <span style={recordMetricLabelStyle}>{copy.artifacts.link}</span>
          <input
            value={artifactUrlInput}
            onChange={(event) => setArtifactUrlInput(event.target.value)}
            disabled={!planningSummary.selectedLesson || isSavingArtifactEntry}
            style={recordMetricInputStyle}
            placeholder="https://..."
          />
        </label>
      </div>

      <label style={recordApplicationFieldStyle}>
        <span style={recordMetricLabelStyle}>{copy.artifacts.description}</span>
        <textarea
          value={artifactDescriptionInput}
          onChange={(event) => setArtifactDescriptionInput(event.target.value)}
          rows={5}
          maxLength={800}
          placeholder={planningSummary.selectedLesson ? copy.artifacts.descriptionPlaceholder : copy.artifacts.descriptionDisabledPlaceholder}
          disabled={!planningSummary.selectedLesson || isSavingArtifactEntry}
          style={journalAnswerTextareaStyle}
        />
      </label>

      <div style={journalActionRowStyle}>
        <button
          type="button"
          onClick={handleResetArtifactDraft}
          disabled={!hasArtifactDraftChanges || isSavingArtifactEntry}
          style={{
            ...journalSecondaryButtonStyle,
            ...((!hasArtifactDraftChanges || isSavingArtifactEntry) ? journalSecondaryButtonDisabledStyle : null),
          }}
        >
          {copy.common.reset}
        </button>
        <button
          type="button"
          onClick={handleSaveArtifactEntry}
          disabled={!canSaveArtifactEntry}
          style={{
            ...journalPrimaryButtonStyle,
            ...(canSaveArtifactEntry ? journalPrimaryButtonActiveStyle : null),
          }}
        >
          {isSavingArtifactEntry ? copy.artifacts.saving : copy.artifacts.save}
        </button>
      </div>

      <div style={journalNoticeStyle}>
        {artifactSaveMessage
          ?? (planningSummary.selectedLesson
            ? copy.artifacts.noticeReady
            : copy.artifacts.noticeNeedRegion)}
      </div>
    </section>
  );
}
