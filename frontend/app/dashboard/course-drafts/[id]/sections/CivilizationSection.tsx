'use client';

import type { DraftEditorState } from '../useDraftEditor';
import type { DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft';
import {
  diaryPageEyebrowStyle,
  diaryPageHeaderStyle,
  diaryPageTitleStyle,
  journalContextCardStyle,
  journalContextDescriptionStyle,
  journalContextGridStyle,
  journalContextTitleStyle,
  journalScaffoldStyle,
  pageStatusPillStyle,
  plannerStatLabelStyle,
  previewChecklistItemStyle,
  previewChecklistStyle,
} from '../styles';

export function CivilizationSection({ editor: _editor, copy }: { editor: DraftEditorState; copy: DashboardCourseDraftCopy['supportSections'] }) {
  return (
    <section style={journalScaffoldStyle}>
      <div style={diaryPageHeaderStyle}>
        <div>
          <div style={diaryPageEyebrowStyle}>{copy.civilization.eyebrow}</div>
          <h3 style={diaryPageTitleStyle}>{copy.civilization.title}</h3>
        </div>
        <div style={pageStatusPillStyle}>{copy.civilization.status}</div>
      </div>

      <div style={journalContextGridStyle}>
        <article style={journalContextCardStyle}>
          <span style={plannerStatLabelStyle}>{copy.common.currentStage}</span>
          <strong style={journalContextTitleStyle}>{copy.common.draftPreview}</strong>
          <p style={journalContextDescriptionStyle}>
            {copy.civilization.stageDescription}
          </p>
        </article>
        <article style={journalContextCardStyle}>
          <span style={plannerStatLabelStyle}>{copy.civilization.metricLabel}</span>
          <strong style={journalContextTitleStyle}>{copy.civilization.metricValue}</strong>
          <p style={journalContextDescriptionStyle}>
            {copy.civilization.metricDescription}
          </p>
        </article>
      </div>

      <div style={previewChecklistStyle}>
        {copy.civilization.items.map((item) => (
          <div key={item.title} style={previewChecklistItemStyle}>
            <strong>{item.title}</strong>
            <span>{item.body}</span>
          </div>
        ))}
      </div>
    </section>
  );
}
