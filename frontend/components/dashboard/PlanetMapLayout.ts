export type PlanetSlotSize = 'sm' | 'md' | 'lg';
export type PlanetLayoutBreakpoint = 'desktop' | 'mobile';

export interface PlanetLayoutPosition {
  id: string;
  position: { x: number; y: number };
  size: PlanetSlotSize;
}

// Orbit ellipse radii as % of planetSafeArea (center = 50%, 50%)
const ORBIT_DESKTOP = { rx: 32, ry: 27 };
const ORBIT_MOBILE = { rx: 30, ry: 24 };
const MOBILE_STACK_POSITIONS = [
  { x: 50, y: 8 },
  { x: 50, y: 29 },
  { x: 50, y: 50 },
  { x: 50, y: 71 },
  { x: 50, y: 92 },
] as const;

function computeOrbitPosition(
  angleDeg: number,
  orbit: { rx: number; ry: number },
): { x: number; y: number } {
  const rad = (angleDeg * Math.PI) / 180;
  return {
    x: 50 + orbit.rx * Math.cos(rad),
    y: 50 + orbit.ry * Math.sin(rad),
  };
}

// Distribute N angles evenly starting from 11 o'clock (-120°), clockwise
function distributeAngles(count: number): number[] {
  return Array.from({ length: count }, (_, i) => -120 + (360 / count) * i);
}

function orbitSlotSize(count: number): PlanetSlotSize {
  if (count <= 1) return 'lg';
  if (count <= 3) return 'md';
  return 'sm';
}

export function getPlanetLayout({
  planetIds,
  breakpoint,
}: {
  planetIds: string[];
  breakpoint: PlanetLayoutBreakpoint;
}): PlanetLayoutPosition[] {
  const capped = planetIds.slice(0, 5);
  const size = orbitSlotSize(capped.length);

  if (breakpoint === 'mobile') {
    return capped.map((id, i) => ({
      id,
      position: MOBILE_STACK_POSITIONS[i] ?? MOBILE_STACK_POSITIONS[MOBILE_STACK_POSITIONS.length - 1],
      size,
    }));
  }

  const orbit = ORBIT_DESKTOP;
  const angles = distributeAngles(capped.length);

  return capped.map((id, i) => ({
    id,
    position: computeOrbitPosition(angles[i], orbit),
    size,
  }));
}

export function planetSlotPixelSize(size: PlanetSlotSize, isCompactLayout: boolean): number {
  if (isCompactLayout) {
    if (size === 'lg') return 112;
    if (size === 'md') return 92;
    return 74;
  }
  if (size === 'lg') return 142;
  if (size === 'md') return 108;
  return 84;
}
