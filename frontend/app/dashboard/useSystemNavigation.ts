'use client';

import type { MutableRefObject, Dispatch } from 'react';
import type { PageSystemChunk, PageHoverLumiState } from '@/components/dashboard/GalaxyMap';
import type { DashboardModeAction, DashboardModeContract } from '@/lib/dashboard/dashboardState';
import type { DashboardSortKey } from '@/lib/world-ui-engine/types';
import type { PlanetFilter } from './_types';

// ── 입력 ────────────────────────────────────────────────────
export interface UseSystemNavigationOptions {
  mode: DashboardModeContract;
  query: string;
  activeStatusFilter: PlanetFilter;
  sortKey: DashboardSortKey;
  selectedSystemChunk: PageSystemChunk | null;
  selectedMapPage: number | null;
  pageSystems: PageSystemChunk[];
  totalPages: number;
  systemsPerGalaxyPage: number;
  pageHoverTimerRef: MutableRefObject<number | null>;
  dispatchDashboardMode: Dispatch<DashboardModeAction>;
  setSelectedMapPage: (p: number | null) => void;
  setPageHoverLumi: (v: PageHoverLumiState | null) => void;
  setHoverSource: (v: 'system' | 'course' | 'course-in-galaxy' | null) => void;
  syncRouteState: (
    nextPage: number,
    nextQuery: string,
    nextStatus: PlanetFilter,
    nextMode?: DashboardModeContract,
    nextSystemId?: string | null,
    nextSortKey?: DashboardSortKey,
  ) => void;
}

// ── 출력 ────────────────────────────────────────────────────
export interface UseSystemNavigationReturn {
  prevSystemChunk: PageSystemChunk | null;
  nextSystemChunk: PageSystemChunk | null;
  handleSystemSelect: (pageSystem: PageSystemChunk) => void;
  handleReturnToGalaxy: () => void;
  handlePrevSystem: () => void;
  handleNextSystem: () => void;
  handleModeToggle: () => void;
}

const EMPTY_STATUS_COUNT = { draft: 0, ready: 0, learning: 0, completed: 0, inactive: 0 };

export function useSystemNavigation({
  mode,
  query,
  activeStatusFilter,
  sortKey,
  selectedSystemChunk,
  selectedMapPage,
  pageSystems,
  totalPages,
  systemsPerGalaxyPage,
  pageHoverTimerRef,
  dispatchDashboardMode,
  setSelectedMapPage,
  setPageHoverLumi,
  setHoverSource,
  syncRouteState,
}: UseSystemNavigationOptions): UseSystemNavigationReturn {
  // ── 공통 hover 상태 정리 ─────────────────────────────────
  const clearHoverState = () => {
    if (pageHoverTimerRef.current !== null) {
      window.clearTimeout(pageHoverTimerRef.current);
      pageHoverTimerRef.current = null;
    }
    setPageHoverLumi(null);
    setHoverSource(null);
  };

  // ── 항성계 선택 (Galaxy → Star-System 진입) ──────────────
  const handleSystemSelect = (pageSystem: PageSystemChunk) => {
    clearHoverState();
    dispatchDashboardMode({ type: 'setMode', payload: 'star-system' });
    dispatchDashboardMode({ type: 'setSelectedSystemId', payload: pageSystem.id });
    dispatchDashboardMode({ type: 'setHoveredSystemId', payload: null });
    dispatchDashboardMode({ type: 'resetSelection' });
    dispatchDashboardMode({ type: 'setPage', payload: pageSystem.index });
    setSelectedMapPage(pageSystem.pageNumber);
    syncRouteState(pageSystem.index, query, activeStatusFilter, 'star-system', pageSystem.id, sortKey);
  };

  // ── Galaxy 복귀 ──────────────────────────────────────────
  const handleReturnToGalaxy = () => {
    const fallbackSystemIndex = selectedSystemChunk?.index ?? selectedMapPage ?? 1;
    const targetGalaxyPage = Math.max(
      1,
      Math.floor((fallbackSystemIndex - 1) / systemsPerGalaxyPage) + 1,
    );
    clearHoverState();
    dispatchDashboardMode({ type: 'setMode', payload: 'galaxy' });
    dispatchDashboardMode({ type: 'setSelectedSystemId', payload: null });
    dispatchDashboardMode({ type: 'setHoveredSystemId', payload: null });
    dispatchDashboardMode({ type: 'resetSelection' });
    dispatchDashboardMode({ type: 'setPage', payload: targetGalaxyPage });
    // In galaxy mode selectedMapPage tracks the visible galaxy page, not the
    // absolute system index. Keeping these in the same coordinate system avoids
    // a brief highlight/state mismatch while the route sync catches up.
    setSelectedMapPage(targetGalaxyPage);
    syncRouteState(targetGalaxyPage, query, activeStatusFilter, 'galaxy', null, sortKey);
  };

  // ── 이전·다음 항성계 ─────────────────────────────────────
  const currentSystemIndex = pageSystems.findIndex((s) => s.id === selectedSystemChunk?.id);

  let prevSystemChunk: PageSystemChunk | null;
  let nextSystemChunk: PageSystemChunk | null;

  if (mode === 'star-system' && selectedSystemChunk) {
    const idx = selectedSystemChunk.index;
    prevSystemChunk =
      idx > 1
        ? {
            id: `system-${idx - 1}`,
            index: idx - 1,
            pageNumber: idx - 1,
            title: '',
            courses: [],
            statusCount: EMPTY_STATUS_COUNT,
            summaryTitle: '',
            updatedAt: null,
          }
        : null;
    nextSystemChunk =
      idx < totalPages
        ? {
            id: `system-${idx + 1}`,
            index: idx + 1,
            pageNumber: idx + 1,
            title: '',
            courses: [],
            statusCount: EMPTY_STATUS_COUNT,
            summaryTitle: '',
            updatedAt: null,
          }
        : null;
  } else {
    prevSystemChunk =
      currentSystemIndex > 0 ? (pageSystems[currentSystemIndex - 1] ?? null) : null;
    nextSystemChunk =
      currentSystemIndex < pageSystems.length - 1
        ? (pageSystems[currentSystemIndex + 1] ?? null)
        : null;
  }

  const handlePrevSystem = () => {
    if (prevSystemChunk) handleSystemSelect(prevSystemChunk);
  };

  const handleNextSystem = () => {
    if (nextSystemChunk) handleSystemSelect(nextSystemChunk);
  };

  // ── 모드 토글 (Galaxy ↔ Star-System) ────────────────────
  const handleModeToggle = () => {
    if (mode === 'star-system') {
      handleReturnToGalaxy();
      return;
    }
    const selectedPageSystem = selectedMapPage
      ? (pageSystems.find((ps) => ps.pageNumber === selectedMapPage) ?? null)
      : null;
    const target = selectedPageSystem ?? pageSystems[0];
    if (target) handleSystemSelect(target);
  };

  return {
    prevSystemChunk,
    nextSystemChunk,
    handleSystemSelect,
    handleReturnToGalaxy,
    handlePrevSystem,
    handleNextSystem,
    handleModeToggle,
  };
}
