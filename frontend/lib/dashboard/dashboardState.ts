import type {
  DashboardMode,
  DashboardSortKey,
  DashboardStatusFilter,
} from '@/lib/world-ui-engine/types';

export type DashboardModeContract = Extract<DashboardMode, 'galaxy' | 'star-system'>;

export interface DashboardModeState {
  mode: DashboardModeContract;
  page: number;
  searchQuery: string;
  statusFilter: DashboardStatusFilter;
  sortKey: DashboardSortKey;
  selectedSystemId: string | null;
  hoveredSystemId: string | null;
  hoveredCourseId: string | null;
  selectedCourseId: string | null;
}

export type DashboardModeAction =
  | { type: 'setMode'; payload: DashboardModeContract }
  | { type: 'setPage'; payload: number }
  | { type: 'setSearchQuery'; payload: string }
  | { type: 'setStatusFilter'; payload: DashboardStatusFilter }
  | { type: 'setSortKey'; payload: DashboardSortKey }
  | { type: 'setSelectedSystemId'; payload: string | null }
  | { type: 'setHoveredSystemId'; payload: string | null }
  | { type: 'setHoveredCourseId'; payload: string | null }
  | { type: 'setSelectedCourseId'; payload: string | null }
  | { type: 'resetSelection' };

export const initialDashboardModeState: DashboardModeState = {
  mode: 'galaxy',
  page: 1,
  searchQuery: '',
  statusFilter: 'all',
  sortKey: 'updated_desc',
  selectedSystemId: null,
  hoveredSystemId: null,
  hoveredCourseId: null,
  selectedCourseId: null,
};

export function getPageSizeForMode(mode: DashboardModeContract): number {
  return mode === 'galaxy' ? 25 : 5;
}

export function dashboardModeReducer(
  state: DashboardModeState,
  action: DashboardModeAction,
): DashboardModeState {
  switch (action.type) {
    case 'setMode':
      return {
        ...state,
        mode: action.payload,
        page: 1,
        selectedSystemId: action.payload === 'galaxy' ? null : state.selectedSystemId,
      };
    case 'setPage':
      return {
        ...state,
        page: Math.max(1, action.payload),
      };
    case 'setSearchQuery':
      return {
        ...state,
        searchQuery: action.payload,
        page: 1,
      };
    case 'setStatusFilter':
      return {
        ...state,
        statusFilter: action.payload,
        selectedSystemId: null,
        page: 1,
      };
    case 'setSortKey':
      return {
        ...state,
        sortKey: action.payload,
        page: 1,
      };
    case 'setSelectedSystemId':
      return {
        ...state,
        selectedSystemId: action.payload,
      };
    case 'setHoveredSystemId':
      return {
        ...state,
        hoveredSystemId: action.payload,
      };
    case 'setHoveredCourseId':
      return {
        ...state,
        hoveredCourseId: action.payload,
      };
    case 'setSelectedCourseId':
      return {
        ...state,
        selectedCourseId: action.payload,
      };
    case 'resetSelection':
      return {
        ...state,
        hoveredCourseId: null,
        selectedCourseId: null,
        hoveredSystemId: null,
        selectedSystemId: null,
      };
    default:
      return state;
  }
}
