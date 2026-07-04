import type {
  DashboardCourseRecord,
  DashboardPlanetStateViewModel,
} from './types';
import { toCanonicalPlanetStatus } from './types';

function resolveTerraformStage(course: DashboardCourseRecord): 1 | 2 | 3 | 4 | 5 {
  const canonicalStatus = toCanonicalPlanetStatus(course.status);

  if (canonicalStatus === 'completed') return 5;
  if (canonicalStatus === 'exploring') {
    if ((course.progress ?? 0) >= 0.65) return 4;
    return 3;
  }
  if (canonicalStatus === 'ready') return 2;
  return 1;
}

function resolveTerraformPhase(stage: 1 | 2 | 3 | 4 | 5): DashboardPlanetStateViewModel['terraformPhase'] {
  if (stage === 1) return 'barren';
  if (stage === 2) return 'atmosphere';
  if (stage === 3) return 'biosphere';
  if (stage === 4) return 'living';
  return 'civilized';
}

function resolveToneFamily(stage: 1 | 2 | 3 | 4 | 5): DashboardPlanetStateViewModel['toneFamily'] {
  if (stage === 1) return 'rock';
  if (stage === 2) return 'cloud';
  if (stage === 3) return 'ocean';
  if (stage === 4) return 'ice';
  return 'gas';
}

function resolveCivilizationLevel(stage: 1 | 2 | 3 | 4 | 5): 0 | 1 | 2 | 3 {
  if (stage <= 2) return 0;
  if (stage === 3) return 1;
  if (stage === 4) return 2;
  return 3;
}

export function buildDashboardPlanetState(
  course: DashboardCourseRecord,
): DashboardPlanetStateViewModel {
  const canonicalStatus = toCanonicalPlanetStatus(course.status);
  const terraformStage = resolveTerraformStage(course);

  return {
    courseId: course.id,
    canonicalStatus,
    terraformPhase: resolveTerraformPhase(terraformStage),
    terraformStage,
    toneFamily: resolveToneFamily(terraformStage),
    civilizationLevel: resolveCivilizationLevel(terraformStage),
  };
}
