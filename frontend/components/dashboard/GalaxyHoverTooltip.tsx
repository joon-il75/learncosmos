'use client';

import { createPortal } from 'react-dom';
import LumiAvatar from '@/components/lumi/LumiAvatar';
import type { PageSystemChunk, PageHoverLumiState } from '@/components/dashboard/GalaxyMap';
import type { DashboardMainCopy } from '@/lib/i18n/pages/dashboardMain';
import {
  mapHoverLumiTooltipStyle,
  mapHoverLumiStackStyle,
  mapHoverLumiCardStyle,
  mapHoverLumiHeaderStyle,
  mapHoverLumiTitleStyle,
  mapHoverLumiMetaStyle,
  mapHoverLumiBubbleStyle,
  mapHoverLumiSummaryRowStyle,
  mapHoverLumiCoursePanelStyle,
  mapHoverLumiCoursePanelHeaderStyle,
  mapHoverLumiCourseListStyle,
  mapHoverLumiCourseItemStyle,
  mapHoverLumiCourseTitleStyle,
  mapHoverLumiCourseMetaStyle,
} from '@/components/dashboard/GalaxyMapStyles';
import type { DashboardCourseRecord } from '@/lib/world-ui-engine/types';
import type { LumiState } from '@/lib/lumi/lumiTypes';
import type { Locale } from '@/lib/i18n/locales';
import { buildSystemStatusLine, formatSystemUpdatedAt } from '@/components/dashboard/GalaxyMap';

function lumiHoverStateForPage(statusCount: PageSystemChunk['statusCount']): LumiState {
  const { draft = 0, ready = 0, learning = 0, completed = 0, inactive = 0 } = statusCount;
  const total = draft + ready + learning + completed + inactive;
  if (total === 0) return 'curious';

  if (draft > 0 && ready === 0 && learning === 0 && completed === 0 && inactive === 0) return 'thinking';
  if (ready > 0 && draft === 0 && learning === 0 && completed === 0 && inactive === 0) return 'planet-hold';
  if (learning > 0 && draft === 0 && ready === 0 && completed === 0 && inactive === 0) return 'exploring';
  if (completed > 0 && draft === 0 && ready === 0 && learning === 0 && inactive === 0) return 'celebrate';

  return 'focus';
}

function buildPageHoverLumiMessage(
  pageSystem: PageSystemChunk,
  copy: DashboardMainCopy['galaxy']['hover'],
): string {
  if (pageSystem.courses.length === 0) {
    return copy.emptySystem(pageSystem.summaryTitle);
  }
  const { draft, ready, learning, completed, inactive } = pageSystem.statusCount;
  const title = pageSystem.summaryTitle;
  const total = pageSystem.courses.length;
  const activeKinds = [draft > 0, ready > 0, learning > 0, completed > 0, inactive > 0].filter(Boolean).length;

  if (learning > 0 && activeKinds === 1) {
    return copy.learningOnly(title);
  }
  if (inactive === total) {
    return copy.inactiveOnly(title);
  }
  if (completed === total) {
    return copy.completedOnly(title);
  }
  if (ready > 0 && activeKinds === 1) {
    return copy.readyOnly(title);
  }
  if (draft === total) {
    return copy.draftOnly(title);
  }
  return copy.mixed(title);
}

interface GalaxyHoverTooltipProps {
  hoverLumi: PageHoverLumiState;
  pageSystem: PageSystemChunk;
  frameRect: DOMRect | null;
  copy: DashboardMainCopy['galaxy'];
  locale: Locale;
}

export function GalaxyHoverTooltip({ hoverLumi, pageSystem, frameRect, copy, locale }: GalaxyHoverTooltipProps) {
  // backdrop-filter가 있는 조상 요소는 position:fixed의 containing block이 되어
  // clientX/clientY 좌표가 어긋나는 브라우저 버그가 있다.
  // createPortal로 document.body에 직접 렌더해 뷰포트 기준 좌표를 보장한다.
  if (typeof document === 'undefined') return null;

  const tooltip = (
    <div style={mapHoverLumiTooltipStyle(hoverLumi.x, hoverLumi.y, frameRect)}>
      <div style={mapHoverLumiStackStyle}>
        <div style={mapHoverLumiCardStyle}>
          <div style={mapHoverLumiHeaderStyle}>
            <LumiAvatar state={lumiHoverStateForPage(pageSystem.statusCount)} size={32} />
            <div style={{ display: 'grid', gap: '2px' }}>
              <strong style={mapHoverLumiTitleStyle}>{copy.hover.title(pageSystem.pageNumber)}</strong>
              <span style={mapHoverLumiMetaStyle}>{copy.hover.courseCount(pageSystem.courses.length)}</span>
            </div>
          </div>
          <div style={mapHoverLumiBubbleStyle}>{buildPageHoverLumiMessage(pageSystem, copy.hover)}</div>
          <div style={mapHoverLumiSummaryRowStyle}>
            <span>{buildSystemStatusLine(pageSystem.statusCount, copy)}</span>
            <span>{formatSystemUpdatedAt(pageSystem.updatedAt, locale, copy)}</span>
          </div>
        </div>
        {pageSystem.courses.length > 0 ? (
          <div style={mapHoverLumiCoursePanelStyle}>
            <div style={mapHoverLumiCoursePanelHeaderStyle}>{copy.hover.courseList}</div>
            <div style={mapHoverLumiCourseListStyle}>
              {pageSystem.courses.map((course: DashboardCourseRecord) => (
                <div key={course.id} style={mapHoverLumiCourseItemStyle()}>
                  <div style={mapHoverLumiCourseTitleStyle()}>{course.title}</div>
                  <div style={mapHoverLumiCourseMetaStyle()}>{copy.statusLabels[course.status]}</div>
                </div>
              ))}
            </div>
          </div>
        ) : null}
      </div>
    </div>
  );

  return createPortal(tooltip, document.body);
}
