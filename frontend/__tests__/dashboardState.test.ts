import { describe, it, expect } from 'vitest';
import {
  dashboardModeReducer,
  initialDashboardModeState,
  getPageSizeForMode,
} from '@/lib/dashboard/dashboardState';

describe('getPageSizeForMode', () => {
  it('galaxy 모드: 25', () => {
    expect(getPageSizeForMode('galaxy')).toBe(25);
  });

  it('star-system 모드: 5', () => {
    expect(getPageSizeForMode('star-system')).toBe(5);
  });
});

describe('dashboardModeReducer', () => {
  it('초기 상태는 galaxy 모드, page 1', () => {
    expect(initialDashboardModeState.mode).toBe('galaxy');
    expect(initialDashboardModeState.page).toBe(1);
    expect(initialDashboardModeState.selectedSystemId).toBeNull();
  });

  it('setMode: star-system으로 전환', () => {
    const next = dashboardModeReducer(initialDashboardModeState, {
      type: 'setMode',
      payload: 'star-system',
    });
    expect(next.mode).toBe('star-system');
  });

  it('setPage: 페이지 변경', () => {
    const next = dashboardModeReducer(initialDashboardModeState, {
      type: 'setPage',
      payload: 3,
    });
    expect(next.page).toBe(3);
  });

  it('setSelectedSystemId: 항성계 선택', () => {
    const next = dashboardModeReducer(initialDashboardModeState, {
      type: 'setSelectedSystemId',
      payload: 'system-2',
    });
    expect(next.selectedSystemId).toBe('system-2');
  });

  it('resetSelection: selectedCourseId, hoveredCourseId 초기화', () => {
    const withSelection = dashboardModeReducer(initialDashboardModeState, {
      type: 'setSelectedCourseId',
      payload: 'course-abc',
    });
    const reset = dashboardModeReducer(withSelection, { type: 'resetSelection' });
    expect(reset.selectedCourseId).toBeNull();
    expect(reset.hoveredCourseId).toBeNull();
  });

  it('setSearchQuery: 검색어 변경', () => {
    const next = dashboardModeReducer(initialDashboardModeState, {
      type: 'setSearchQuery',
      payload: '수채화',
    });
    expect(next.searchQuery).toBe('수채화');
  });

  it('setStatusFilter: 필터 변경', () => {
    const next = dashboardModeReducer(initialDashboardModeState, {
      type: 'setStatusFilter',
      payload: 'learning',
    });
    expect(next.statusFilter).toBe('learning');
  });
});
