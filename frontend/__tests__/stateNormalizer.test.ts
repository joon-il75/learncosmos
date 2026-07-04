import { describe, it, expect } from 'vitest';
import {
  buildDashboardCourseRecordFromDraft,
  isDashboardPlanetStatus,
  mapDraftStatusToPlanetStatus,
  planetStatusSortWeight,
} from '@/lib/world-ui-engine/stateNormalizer';

describe('mapDraftStatusToPlanetStatus', () => {
  it('draft는 draft로 유지한다', () => {
    expect(mapDraftStatusToPlanetStatus('draft')).toBe('draft');
  });

  it('confirmed는 ready로 매핑한다', () => {
    expect(mapDraftStatusToPlanetStatus('confirmed')).toBe('ready');
  });

  it('learning은 learning으로 매핑한다', () => {
    expect(mapDraftStatusToPlanetStatus('learning')).toBe('learning');
  });

  it('archived는 completed로 매핑한다', () => {
    expect(mapDraftStatusToPlanetStatus('archived')).toBe('completed');
  });
});

describe('isDashboardPlanetStatus', () => {
  it('허용된 상태만 true를 반환한다', () => {
    expect(isDashboardPlanetStatus('draft')).toBe(true);
    expect(isDashboardPlanetStatus('ready')).toBe(true);
    expect(isDashboardPlanetStatus('learning')).toBe(true);
    expect(isDashboardPlanetStatus('completed')).toBe(true);
    expect(isDashboardPlanetStatus('confirmed')).toBe(false);
    expect(isDashboardPlanetStatus(null)).toBe(false);
  });
});

describe('planetStatusSortWeight', () => {
  it('exploring(learning)이 가장 높은 우선순위를 가진다', () => {
    expect(planetStatusSortWeight('learning')).toBe(0);
    expect(planetStatusSortWeight('ready')).toBe(1);
    expect(planetStatusSortWeight('draft')).toBe(2);
    expect(planetStatusSortWeight('completed')).toBe(3);
  });
});

describe('buildDashboardCourseRecordFromDraft', () => {
  it('draft 데이터를 대시보드 코스 레코드로 변환한다', () => {
    const result = buildDashboardCourseRecordFromDraft(
      {
        id: 'draft-1',
        source_query: '기타 독학',
        title: '기타 로드맵',
        status: 'confirmed',
        updated_at: '2026-04-08T00:00:00Z',
      },
      3,
      12,
      2,
      8,
      5,
      0.5,
    );

    expect(result).toEqual({
      id: 'draft-1',
      title: '기타 로드맵',
      status: 'ready',
      updatedAt: '2026-04-08T00:00:00Z',
      levelCount: 3,
      lessonCount: 12,
      plannedLevelCount: 2,
      plannedLessonCount: 8,
      completedLessonCount: 5,
      progress: 0.5,
      sourceQuery: '기타 독학',
      draftStatus: 'confirmed',
      draftId: 'draft-1',
      planetTypeId: null,
      planetTypeName: null,
      planetTypeAsset: null,
      planetTextureMapId: null,
      planetTextureMapName: null,
      planetTextureMapAsset: null,
      planetTextureMapRotationDurationSeconds: null,
      planetTextureMapRotationDirection: null,
      destination: '/dashboard/course-drafts/draft-1',
    });
  });

  it('title이 비어 있으면 source_query를 fallback title로 사용한다', () => {
    const result = buildDashboardCourseRecordFromDraft(
      {
        id: 'draft-2',
        source_query: '  수채화 시작하기  ',
        title: '   ',
        status: 'draft',
        updated_at: '2026-04-08T00:00:00Z',
      },
      0,
      0,
    );

    expect(result.title).toBe('수채화 시작하기');
  });

  it('source_query도 비어 있으면 새 행성을 fallback title로 사용한다', () => {
    const result = buildDashboardCourseRecordFromDraft(
      {
        id: 'draft-3',
        source_query: '   ',
        title: '',
        status: 'draft',
        updated_at: '2026-04-08T00:00:00Z',
      },
      0,
      0,
    );

    expect(result.title).toBe('새 행성');
  });
});
