export interface PlanetProgressGridPosition {
  row: 1 | 2;
  col: 1 | 2 | 3 | 4;
  stage: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7;
}

export interface PlanetProgressVisualProfile {
  stage: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7;
  brightness: number;
  saturate: number;
  glowOpacity: number;
  ringOpacity: number;
  scaleBoost: number;
}

const STAGE_RULES: ReadonlyArray<{
  min: number;
  max: number;
  row: 1 | 2;
  col: 1 | 2 | 3 | 4;
  stage: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7;
}> = [
  { min: 0, max: 0, row: 1, col: 1, stage: 0 },
  { min: 1, max: 17, row: 1, col: 2, stage: 1 },
  { min: 18, max: 33, row: 1, col: 3, stage: 2 },
  { min: 34, max: 50, row: 1, col: 4, stage: 3 },
  { min: 51, max: 66, row: 2, col: 1, stage: 4 },
  { min: 67, max: 82, row: 2, col: 2, stage: 5 },
  { min: 83, max: 99, row: 2, col: 3, stage: 6 },
  { min: 100, max: 100, row: 2, col: 4, stage: 7 },
] as const;

const VISUAL_PROFILES: Record<PlanetProgressGridPosition['stage'], PlanetProgressVisualProfile> = {
  0: { stage: 0, brightness: 0.58, saturate: 0.84, glowOpacity: 0.14, ringOpacity: 0.42, scaleBoost: 0 },
  1: { stage: 1, brightness: 0.68, saturate: 0.9, glowOpacity: 0.18, ringOpacity: 0.48, scaleBoost: 0.005 },
  2: { stage: 2, brightness: 0.76, saturate: 0.96, glowOpacity: 0.24, ringOpacity: 0.52, scaleBoost: 0.01 },
  3: { stage: 3, brightness: 0.84, saturate: 1.02, glowOpacity: 0.3, ringOpacity: 0.58, scaleBoost: 0.015 },
  4: { stage: 4, brightness: 0.9, saturate: 1.08, glowOpacity: 0.38, ringOpacity: 0.66, scaleBoost: 0.02 },
  5: { stage: 5, brightness: 0.98, saturate: 1.14, glowOpacity: 0.48, ringOpacity: 0.74, scaleBoost: 0.03 },
  6: { stage: 6, brightness: 1.06, saturate: 1.22, glowOpacity: 0.6, ringOpacity: 0.82, scaleBoost: 0.04 },
  7: { stage: 7, brightness: 1.16, saturate: 1.3, glowOpacity: 0.74, ringOpacity: 0.94, scaleBoost: 0.05 },
};

export function normalizePlanetProgressPercent(progress: number | null | undefined, status?: string): number {
  if (status === 'completed') return 100;
  if (status === 'draft' || status === 'ready') return 0;
  const numericProgress = Number.isFinite(progress) ? Number(progress) : 0;
  if (numericProgress <= 1) {
    return Math.max(0, Math.min(100, Math.round(numericProgress * 100)));
  }
  return Math.max(0, Math.min(100, Math.round(numericProgress)));
}

export function getProgressGridPosition(progressPercent: number): PlanetProgressGridPosition {
  const safePercent = Math.max(0, Math.min(100, Math.round(progressPercent)));
  const matchedRule = STAGE_RULES.find((rule) => safePercent >= rule.min && safePercent <= rule.max) ?? STAGE_RULES[0];
  return {
    row: matchedRule.row,
    col: matchedRule.col,
    stage: matchedRule.stage,
  };
}

export function getPlanetProgressVisualProfile(progressPercent: number): PlanetProgressVisualProfile {
  const stage = getProgressGridPosition(progressPercent).stage;
  return VISUAL_PROFILES[stage];
}
