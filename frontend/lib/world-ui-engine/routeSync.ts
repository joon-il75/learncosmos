import type {
  DashboardMode,
  DashboardSortKey,
  DashboardStatusFilter,
} from './types';

export interface DashboardRouteSyncInput {
  currentSearch: string;
  nextPage: number;
  nextQuery: string;
  nextStatus: DashboardStatusFilter;
  nextMode: DashboardMode;
  nextSystemId?: string | null;
  nextSortKey: DashboardSortKey;
}

export function buildDashboardRouteParams(
  input: DashboardRouteSyncInput,
): URLSearchParams {
  const params = new URLSearchParams(input.currentSearch);
  const normalizedQuery = input.nextQuery.trim();

  if (normalizedQuery) params.set('q', normalizedQuery);
  else params.delete('q');

  if (input.nextStatus === 'all') params.delete('status');
  else params.set('status', input.nextStatus);

  if (input.nextPage <= 1) params.delete('page');
  else params.set('page', String(input.nextPage));

  if (input.nextMode === 'galaxy') {
    params.set('mode', 'galaxy');
  } else {
    params.set('mode', input.nextMode);
  }

  if (input.nextSystemId) {
    params.set('system', input.nextSystemId);
  } else {
    params.delete('system');
  }

  if (input.nextSortKey === 'updated_desc') {
    params.delete('sort');
  } else {
    params.set('sort', input.nextSortKey);
  }

  return params;
}
