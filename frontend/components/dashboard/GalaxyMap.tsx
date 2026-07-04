'use client';

import Image from 'next/image';
import { useEffect, useRef, useState } from 'react';
import PlanetMap from '@/components/dashboard/PlanetMap';
import { GalaxyMobileInstruction } from '@/components/dashboard/GalaxyMobileInstruction';
import { GalaxySystemNode } from '@/components/dashboard/GalaxySystemNode';
import { GalaxyHoverTooltip } from '@/components/dashboard/GalaxyHoverTooltip';
import type { DashboardMainCopy } from '@/lib/i18n/pages/dashboardMain';
import { getDashboardMainCopy } from '@/lib/i18n/pages/dashboardMain';
import type { Locale } from '@/lib/i18n/locales';
import { defaultLocale } from '@/lib/i18n/locales';
import type { PlanetStatus } from '@/components/dashboard/planetMapTokens';
import {
  mapWindowStyle,
  mapWindowHeaderStyle,
  mapWindowHeaderGroupStyle,
  mapWindowEyebrowStyle,
  mapWindowMetaStyle,
  sortControlRowStyle,
  sortButtonStyle,
  mapWindowSinglePaneStyle,
  mapWindowFrameStyle,
  mapWindowImageStyle,
  mapWindowImageStarSystemStyle,
  starSystemMapWrapperStyle,
  starSystemMissingStateStyle,
  mapWindowGlowStyle,
} from '@/components/dashboard/GalaxyMapStyles';
import type { DashboardSortKey, DashboardCourseRecord } from '@/lib/world-ui-engine/types';
import type { StarSystemGroup } from '@/lib/dashboard/groupCoursesToSystems';

// ─── Types ──────────────────────────────────────────────────────────────────

export type PageSystemChunk = StarSystemGroup & { pageNumber: number };

export interface PageHoverLumiState {
  pageNumber: number;
  x: number;
  y: number;
}

// ─── Exported helpers ────────────────────────────────────────────────────────

export function buildSystemStatusLine(
  statusCount: PageSystemChunk['statusCount'],
  copy: DashboardMainCopy['galaxy'] = getDashboardMainCopy(defaultLocale).galaxy,
): string {
  const summaryParts = [
    statusCount.draft > 0 ? `${copy.statusLabels.draft} ${statusCount.draft}${copy.statusLine.countSuffix}` : null,
    statusCount.ready > 0 ? `${copy.statusLabels.ready} ${statusCount.ready}${copy.statusLine.countSuffix}` : null,
    statusCount.learning > 0 ? `${copy.statusLabels.learning} ${statusCount.learning}${copy.statusLine.countSuffix}` : null,
    statusCount.completed > 0 ? `${copy.statusLabels.completed} ${statusCount.completed}${copy.statusLine.countSuffix}` : null,
    statusCount.inactive > 0 ? `${copy.statusLabels.inactive} ${statusCount.inactive}${copy.statusLine.countSuffix}` : null,
  ].filter(Boolean);
  return summaryParts.length > 0 ? `${copy.statusLine.prefix} ${summaryParts.join(' · ')}` : copy.statusLine.empty;
}

export function formatSystemUpdatedAt(
  updatedAt: string | null,
  locale: Locale = defaultLocale,
  copy: DashboardMainCopy['galaxy'] = getDashboardMainCopy(defaultLocale).galaxy,
): string {
  if (!updatedAt) return copy.updatedAt.empty;
  const label = new Date(updatedAt).toLocaleDateString(locale === 'en' ? 'en-US' : 'ko-KR', { month: 'short', day: 'numeric' });
  return `${copy.updatedAt.prefix} ${label}`;
}

function getMapObjectScale(frameWidth: number, isCompactLayout: boolean): number {
  if (frameWidth <= 0) return isCompactLayout ? 0.64 : 1;
  return Math.max(0.56, Math.min(1, frameWidth / 820));
}

// ─── Props ───────────────────────────────────────────────────────────────────

interface GalaxyMapProps {
  pageSystems: PageSystemChunk[];
  mode: 'galaxy' | 'star-system';
  selectedSystemChunk: PageSystemChunk | null;
  selectedMapPage: number | null;
  hoveredSystemId: string | null;
  hoveredCourseId: string | null;
  selectedCourseId: string | null;
  pageHoverLumi: PageHoverLumiState | null;
  hoveredPageSystem: PageSystemChunk | null;
  selectedSystemTitleFallback: string | undefined;
  starSystemStatusLine: string | null;
  starSystemUpdatedLabel: string | null;
  sortKey: DashboardSortKey;
  filteredCourseCount: number;
  activeStatusFilter: string;
  isPhoneLayout: boolean;
  isCompactLayout: boolean;
  isShortViewport: boolean;
  showListOnlyFallback: boolean;
  onSystemSelect: (system: PageSystemChunk) => void;
  onPageHoverStart: (system: PageSystemChunk, x: number, y: number) => void;
  onPageHoverMove: (pageNumber: number, x: number, y: number) => void;
  onPageHoverEnd: () => void;
  onHoverCourse: (courseId: string | null) => void;
  onSelectCourse: (courseId: string) => void;
  onReturnToGalaxy: () => void;
  onSortChange: (key: DashboardSortKey) => void;
  locale: Locale;
  copy: DashboardMainCopy['galaxy'];
  planetMapCopy: DashboardMainCopy['planetMap'];
}

// ─── Main component ───────────────────────────────────────────────────────────

export default function GalaxyMap({
  pageSystems,
  mode,
  selectedSystemChunk,
  selectedMapPage,
  hoveredSystemId,
  hoveredCourseId,
  selectedCourseId,
  pageHoverLumi,
  hoveredPageSystem,
  selectedSystemTitleFallback,
  starSystemStatusLine,
  starSystemUpdatedLabel,
  sortKey,
  filteredCourseCount,
  activeStatusFilter,
  isPhoneLayout,
  isCompactLayout,
  isShortViewport,
  showListOnlyFallback,
  onSystemSelect,
  onPageHoverStart,
  onPageHoverMove,
  onPageHoverEnd,
  onHoverCourse,
  onSelectCourse,
  onReturnToGalaxy,
  onSortChange,
  locale,
  copy,
  planetMapCopy,
}: GalaxyMapProps) {
  const frameRef = useRef<HTMLDivElement>(null);
  const [frameWidth, setFrameWidth] = useState(0);
  const mapObjectScale = getMapObjectScale(frameWidth, isCompactLayout);

  useEffect(() => {
    const node = frameRef.current;
    if (!node || typeof ResizeObserver === 'undefined') return;
    const updateWidth = () => setFrameWidth(node.getBoundingClientRect().width);
    updateWidth();
    const observer = new ResizeObserver((entries) => {
      const width = entries[0]?.contentRect.width ?? node.getBoundingClientRect().width;
      setFrameWidth(width);
    });
    observer.observe(node);
    return () => observer.disconnect();
  }, []);

  if (showListOnlyFallback) {
    return null;
  }

  return (
    <section style={mapWindowStyle}>
      {/* Header: 타이틀 + 정렬 */}
      <div style={mapWindowHeaderStyle}>
        <div style={mapWindowHeaderGroupStyle}>
          <div style={mapWindowEyebrowStyle}>{copy.eyebrow}</div>
          <div style={mapWindowMetaStyle}>
            {copy.courseCount(filteredCourseCount)}
            {activeStatusFilter !== 'all' ? ` · ${copy.statusFilter(copy.statusLabels[activeStatusFilter as PlanetStatus] ?? activeStatusFilter)}` : ''}
          </div>
        </div>
        <div style={sortControlRowStyle}>
          {copy.sortOptions.map((option) => (
            <button
              key={option.key}
              type="button"
              onClick={() => onSortChange(option.key)}
              style={sortButtonStyle(option.key === sortKey)}
            >
              {option.label}
            </button>
          ))}
        </div>
      </div>

      {/* 모바일 안내 배너 */}
      {isPhoneLayout && mode !== 'star-system' ? (
        <GalaxyMobileInstruction mode={mode} selectedSystemChunk={selectedSystemChunk} copy={copy.mobile} />
      ) : null}

      {/* 지도 프레임 */}
      <div style={mapWindowSinglePaneStyle}>
        <div style={mapWindowFrameStyle} ref={frameRef}>
          <Image
            src={mode === 'star-system' ? '/images/PlanetMap_background.webp' : '/images/hero_galaxy.webp'}
            alt={copy.imageAlt}
            fill
            priority
            sizes="(max-width: 768px) 100vw, 1360px"
            style={mode === 'star-system' ? mapWindowImageStarSystemStyle : mapWindowImageStyle}
          />

          {/* Star System Mode */}
          {mode === 'star-system' ? (
            <div style={starSystemMapWrapperStyle}>
              {selectedSystemChunk ? (
                <PlanetMap
                  courses={selectedSystemChunk.courses}
                  hoveredCourseId={hoveredCourseId}
                  selectedCourseId={selectedCourseId}
                  isCompactLayout={isCompactLayout}
                  isPhoneLayout={isPhoneLayout}
                  isShortViewport={isShortViewport}
                  objectScale={mapObjectScale}
                  mode="star-system"
                  selectedSystemTitle={selectedSystemTitleFallback}
                  statusLine={starSystemStatusLine}
                  updatedAtLabel={starSystemUpdatedLabel}
                  copy={planetMapCopy}
                  statusLabels={copy.statusLabels}
                  locale={locale}
                  onHoverCourse={onHoverCourse}
                  onSelectCourse={(course) => onSelectCourse(course.id)}
                  onReturnToGalaxy={onReturnToGalaxy}
                />
              ) : (
                <div style={starSystemMissingStateStyle}>
                  {copy.missingSystem}
                </div>
              )}
            </div>
          ) : (
            /* Galaxy Mode */
            <>
              {pageSystems.map((pageSystem, index) => {
                const isHighlighted =
                  selectedMapPage === pageSystem.pageNumber ||
                  pageHoverLumi?.pageNumber === pageSystem.pageNumber ||
                  hoveredSystemId === pageSystem.id;
                return (
                  <GalaxySystemNode
                    key={`orbit-page-${pageSystem.pageNumber}`}
                    pageSystem={pageSystem}
                    index={index}
                    totalNodes={pageSystems.length}
                    isHighlighted={isHighlighted}
                    useColumnLayout={false}
                    objectScale={mapObjectScale}
                    onSystemSelect={onSystemSelect}
                    onPageHoverStart={onPageHoverStart}
                    onPageHoverMove={onPageHoverMove}
                    onPageHoverEnd={onPageHoverEnd}
                  />
                );
              })}

              {/* Hover Lumi tooltip */}
              {mode === 'galaxy' && hoveredPageSystem && pageHoverLumi ? (
                <GalaxyHoverTooltip
                  hoverLumi={pageHoverLumi}
                  pageSystem={hoveredPageSystem}
                  frameRect={frameRef.current ? frameRef.current.getBoundingClientRect() : null}
                  copy={copy}
                  locale={locale}
                />
              ) : null}

              <div style={mapWindowGlowStyle} />
            </>
          )}
        </div>
      </div>
    </section>
  );
}
