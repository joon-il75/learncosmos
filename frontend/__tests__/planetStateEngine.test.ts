import { describe, it, expect } from 'vitest';
import { buildDashboardPlanetState } from '@/lib/world-ui-engine/planetStateEngine';
import type { DashboardCourseRecord } from '@/lib/world-ui-engine/types';

function makeCourse(overrides: Partial<DashboardCourseRecord> = {}): DashboardCourseRecord {
  return {
    id: 'c1',
    title: '테스트 코스',
    status: 'draft',
    draftStatus: 'draft',
    progress: null,
    lessonCount: 0,
    levelCount: 0,
    plannedLessonCount: 0,
    plannedLevelCount: 0,
    completedLessonCount: 0,
    sourceQuery: '',
    updatedAt: '2026-04-07T00:00:00Z',
    destination: '/dashboard/course-drafts/c1',
    ...overrides,
  };
}

describe('buildDashboardPlanetState', () => {
  it('draft는 barren stage 1로 파생된다', () => {
    expect(buildDashboardPlanetState(makeCourse({ status: 'draft' }))).toMatchObject({
      canonicalStatus: 'draft',
      terraformStage: 1,
      terraformPhase: 'barren',
      toneFamily: 'rock',
      civilizationLevel: 0,
    });
  });

  it('ready는 atmosphere stage 2로 파생된다', () => {
    expect(buildDashboardPlanetState(makeCourse({ status: 'ready' }))).toMatchObject({
      canonicalStatus: 'ready',
      terraformStage: 2,
      terraformPhase: 'atmosphere',
      toneFamily: 'cloud',
      civilizationLevel: 0,
    });
  });

  it('learning progress 65% 미만은 biosphere stage 3으로 파생된다', () => {
    expect(buildDashboardPlanetState(makeCourse({ status: 'learning', progress: 0.64 }))).toMatchObject({
      canonicalStatus: 'exploring',
      terraformStage: 3,
      terraformPhase: 'biosphere',
      toneFamily: 'ocean',
      civilizationLevel: 1,
    });
  });

  it('learning progress 65% 이상은 living stage 4로 파생된다', () => {
    expect(buildDashboardPlanetState(makeCourse({ status: 'learning', progress: 0.65 }))).toMatchObject({
      canonicalStatus: 'exploring',
      terraformStage: 4,
      terraformPhase: 'living',
      toneFamily: 'ice',
      civilizationLevel: 2,
    });
  });

  it('completed는 civilized stage 5로 파생된다', () => {
    expect(buildDashboardPlanetState(makeCourse({ status: 'completed', progress: 1 }))).toMatchObject({
      canonicalStatus: 'completed',
      terraformStage: 5,
      terraformPhase: 'civilized',
      toneFamily: 'gas',
      civilizationLevel: 3,
    });
  });
});
