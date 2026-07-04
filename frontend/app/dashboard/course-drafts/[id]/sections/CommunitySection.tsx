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

export function CommunitySection({ editor, copy }: { editor: DraftEditorState; copy: DashboardCourseDraftCopy['supportSections'] }) {
  const { selectedPointContext } = editor;

  return (
    <section style={journalScaffoldStyle}>
      <div style={diaryPageHeaderStyle}>
        <div>
          <div style={diaryPageEyebrowStyle}>{copy.community.eyebrow}</div>
          <h3 style={diaryPageTitleStyle}>{copy.community.title}</h3>
        </div>
        <div style={pageStatusPillStyle}>{copy.community.status}</div>
      </div>

      <div style={journalContextGridStyle}>
        <article style={journalContextCardStyle}>
          <span style={plannerStatLabelStyle}>{copy.common.currentPoint}</span>
          <strong style={journalContextTitleStyle}>{selectedPointContext?.title ?? copy.common.pointNotSelected}</strong>
          <p style={journalContextDescriptionStyle}>
            {copy.community.pointDescription}
          </p>
        </article>
        <article style={journalContextCardStyle}>
          <span style={plannerStatLabelStyle}>{copy.common.currentStage}</span>
          <strong style={journalContextTitleStyle}>{copy.common.draftPreview}</strong>
          <p style={journalContextDescriptionStyle}>
            {copy.community.stageDescription}
          </p>
        </article>
      </div>

      <div style={previewChecklistStyle}>
        {copy.community.items.map((item) => (
          <div key={item.title} style={previewChecklistItemStyle}>
            <strong>{item.title}</strong>
            <span>{item.body}</span>
          </div>
        ))}
      </div>
    </section>
  );
}
