'use client';

import { useState } from 'react';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../../pointPageUtils';
import type { UsePointLearningResult } from '../usePointLearning';
import {
  formGridStyle, sectionMiniHeaderStyle, blockMetaStyle,
  emptyCardStyle, formActionRowStyle,
  pointShellTitleStyle,
} from '../../pointPageStyles';
import { getPointEventLabel, formatPointEventTime } from '../../pointPageUtils';
import { makeWorkspaceStyles } from './workspaceStyles';

interface Props {
  copy: PointLearningCopy['workspace']['learningWork']['recordPanels']['timeline'];
  commonCopy: PointLearningCopy['workspace']['learningWork']['recordPanels']['common'];
  learning: UsePointLearningResult;
  themeTokens: PointPageThemeTokens;
}

const EVENTS_PER_PAGE = 10;

export default function PointTimelineTab({ copy, commonCopy, learning, themeTokens }: Props) {
  const { pointDetail } = learning;
  const { materialToolbarButtonStyle } = makeWorkspaceStyles(themeTokens);

  const [eventPage, setEventPage] = useState(1);
  const [hoveredEventID, setHoveredEventID] = useState<string | null>(null);

  const pointEvents = [...(pointDetail?.events ?? [])].sort((left, right) => {
    const rightTime = Date.parse(right.created_at);
    const leftTime = Date.parse(left.created_at);
    return (Number.isFinite(rightTime) ? rightTime : 0) - (Number.isFinite(leftTime) ? leftTime : 0);
  });
  const eventPageCount = Math.max(1, Math.ceil(pointEvents.length / EVENTS_PER_PAGE));
  const pagedPointEvents = pointEvents.slice((eventPage - 1) * EVENTS_PER_PAGE, eventPage * EVENTS_PER_PAGE);

  const listHoverBackground = themeTokens.pageBackground === '#F8FAFC' ? '#EFF6FF' : 'rgba(96, 165, 250, 0.13)';

  return (
    <div style={formGridStyle}>
      <div style={sectionMiniHeaderStyle}>
        <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.title}</strong>
        <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.subtitle}</span>
      </div>
      {pointEvents.length ? (
        <div style={formGridStyle}>
          <div style={sectionMiniHeaderStyle}>
            <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>
              {commonCopy.countRange(pointEvents.length, (eventPage - 1) * EVENTS_PER_PAGE + 1, Math.min(eventPage * EVENTS_PER_PAGE, pointEvents.length))}
            </span>
          </div>
          <div
            role="table"
            aria-label={copy.listAria}
            style={{
              display: 'grid',
              overflow: 'hidden',
              border: `1px solid ${themeTokens.surfaceBorder}`,
              borderRadius: '4px',
              background: themeTokens.surfaceBackground,
            }}
          >
            <div
              role="row"
              style={{
                display: 'grid',
                gridTemplateColumns: '58px 132px minmax(0, 1fr)',
                alignItems: 'center',
                minHeight: '38px',
                padding: '0 16px',
                borderBottom: `1px solid ${themeTokens.surfaceBorder}`,
                background: themeTokens.pageBackground === '#F8FAFC' ? '#F1F5F9' : 'rgba(148, 163, 184, 0.10)',
              }}
            >
              <span role="columnheader" style={{ ...blockMetaStyle, color: themeTokens.metaLabel, fontWeight: 850 }}>{commonCopy.number}</span>
              <span role="columnheader" style={{ ...blockMetaStyle, color: themeTokens.metaLabel, fontWeight: 850 }}>{copy.time}</span>
              <span role="columnheader" style={{ ...blockMetaStyle, color: themeTokens.metaLabel, fontWeight: 850 }}>{copy.activity}</span>
            </div>
            {pagedPointEvents.map((event, index) => {
              const hovered = hoveredEventID === event.id;
              const eventIndex = pointEvents.length - ((eventPage - 1) * EVENTS_PER_PAGE + index);
              return (
                <div
                  key={event.id}
                  role="row"
                  onMouseEnter={() => setHoveredEventID(event.id)}
                  onMouseLeave={() => setHoveredEventID(null)}
                  style={{
                    display: 'grid',
                    gridTemplateColumns: '58px 132px minmax(0, 1fr)',
                    alignItems: 'center',
                    minHeight: '46px',
                    padding: '9px 16px',
                    borderLeft: `3px solid ${hovered ? themeTokens.primaryButtonBorder : 'transparent'}`,
                    borderBottom: `1px solid ${themeTokens.surfaceBorder}`,
                    background: hovered ? listHoverBackground : themeTokens.surfaceBackground,
                    transition: 'background 140ms ease, border-color 140ms ease',
                  }}
                >
                  <span role="cell" style={{ ...blockMetaStyle, color: hovered ? themeTokens.title : themeTokens.metaLabel, fontWeight: 850 }}>{eventIndex}</span>
                  <span role="cell" style={{ ...blockMetaStyle, color: hovered ? themeTokens.title : themeTokens.metaLabel }}>{formatPointEventTime(event.created_at)}</span>
                  <span role="cell" style={{ color: hovered ? themeTokens.title : themeTokens.description, fontSize: '14px', lineHeight: 1.6, fontWeight: hovered ? 750 : 600 }}>{getPointEventLabel(event.event_type, event.event_payload)}</span>
                </div>
              );
            })}
          </div>
          {eventPageCount > 1 ? (
            <div style={{ ...formActionRowStyle, justifyContent: 'center' }}>
              <button type="button" onClick={() => setEventPage((p) => Math.max(1, p - 1))} disabled={eventPage <= 1} style={{ ...materialToolbarButtonStyle, opacity: eventPage <= 1 ? 0.55 : 1 }}>{commonCopy.previous}</button>
              <span style={{ color: themeTokens.mutedText, fontSize: '14px', fontWeight: 800 }}>{eventPage} / {eventPageCount}</span>
              <button type="button" onClick={() => setEventPage((p) => Math.min(eventPageCount, p + 1))} disabled={eventPage >= eventPageCount} style={{ ...materialToolbarButtonStyle, opacity: eventPage >= eventPageCount ? 0.55 : 1 }}>{commonCopy.next}</button>
            </div>
          ) : null}
        </div>
      ) : (
        <div style={{ ...emptyCardStyle, background: themeTokens.emptyBackground, borderColor: themeTokens.emptyBorder, color: themeTokens.mutedText }}>{copy.empty}</div>
      )}
    </div>
  );
}
