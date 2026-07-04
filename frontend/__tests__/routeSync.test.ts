import { describe, it, expect } from 'vitest';
import { buildDashboardRouteParams } from '@/lib/world-ui-engine/routeSync';

describe('buildDashboardRouteParams', () => {
  it('Galaxy 모드는 명시적 경로 상태로 남기고 나머지 기본값은 제거한다', () => {
    const params = buildDashboardRouteParams({
      currentSearch: '?q=기존&status=learning&page=3&mode=star-system&system=system-2&sort=title_asc',
      nextPage: 1,
      nextQuery: '   ',
      nextStatus: 'all',
      nextMode: 'galaxy',
      nextSystemId: null,
      nextSortKey: 'updated_desc',
    });

    expect(params.toString()).toBe('mode=galaxy');
  });

  it('비기본값은 모두 직렬화한다', () => {
    const params = buildDashboardRouteParams({
      currentSearch: '',
      nextPage: 4,
      nextQuery: '  기타 독학  ',
      nextStatus: 'learning',
      nextMode: 'star-system',
      nextSystemId: 'system-4',
      nextSortKey: 'title_desc',
    });

    expect(params.get('q')).toBe('기타 독학');
    expect(params.get('status')).toBe('learning');
    expect(params.get('page')).toBe('4');
    expect(params.get('mode')).toBe('star-system');
    expect(params.get('system')).toBe('system-4');
    expect(params.get('sort')).toBe('title_desc');
  });

  it('system id가 없으면 기존 system 파라미터를 제거한다', () => {
    const params = buildDashboardRouteParams({
      currentSearch: '?system=system-9',
      nextPage: 1,
      nextQuery: '',
      nextStatus: 'all',
      nextMode: 'star-system',
      nextSystemId: null,
      nextSortKey: 'updated_desc',
    });

    expect(params.get('mode')).toBe('star-system');
    expect(params.has('system')).toBe(false);
  });
});
