export type RegionLabelLayout = {
  x: number
  y: number
  textAnchor?: 'start' | 'middle' | 'end'
}

export type BoundaryLayout = {
  id: string
  d: string
}

export type TerritoryLayout = {
  // SVG polygon points string (viewBox 0 0 1000 560 기준)
  points: string
}

const SINGLE_LABELS: RegionLabelLayout[] = [
  { x: 500, y: 280, textAnchor: 'middle' },
]

const DOUBLE_LABELS: RegionLabelLayout[] = [
  { x: 245, y: 292, textAnchor: 'middle' },
  { x: 745, y: 292, textAnchor: 'middle' },
]

const TRIPLE_LABELS: RegionLabelLayout[] = [
  { x: 500, y: 225, textAnchor: 'middle' },
  { x: 258, y: 390, textAnchor: 'middle' },
  { x: 742, y: 390, textAnchor: 'middle' },
]

const DOUBLE_BOUNDARIES: BoundaryLayout[] = [
  {
    id: 'diagonal-divider',
    d: 'M 480 120 C 460 220, 470 300, 500 360 C 530 420, 520 480, 510 540',
  },
]

const TRIPLE_BOUNDARIES: BoundaryLayout[] = [
  {
    id: 'left-branch',
    d: 'M 500 314 C 438 302, 388 280, 334 242 C 286 208, 236 178, 170 148',
  },
  {
    id: 'right-branch',
    d: 'M 500 314 C 562 302, 612 280, 666 242 C 714 208, 764 178, 830 148',
  },
  {
    id: 'bottom-branch',
    d: 'M 500 314 C 494 356, 488 396, 492 438 C 496 476, 506 512, 518 548',
  },
]

// ── 지역 영역 폴리곤 (경계선에 밀착) ─────────────────────────────────────────
// viewBox 0 0 1000 560 기준 — boundary path를 선분으로 근사해 영역 분할

// Single: 전체 맵 (경계선 없음)
const SINGLE_TERRITORIES: TerritoryLayout[] = [
  { points: '4,4 996,4 996,556 4,556' },
]

// Double: 경계선 M 480 120 C ... 510 540 근사
const DOUBLE_TERRITORIES: TerritoryLayout[] = [
  // 왼쪽 영역
  { points: '4,4 490,4 480,120 460,220 470,300 500,360 530,420 510,540 510,556 4,556' },
  // 오른쪽 영역
  { points: '490,4 996,4 996,556 510,556 510,540 530,420 500,360 470,300 460,220 480,120' },
]

// Triple: 허브(500,314)에서 세 방향으로 분기
// 왼쪽 분기: (500,314)→(170,148)→(0,62) 연장
// 오른쪽 분기: (500,314)→(830,148)→(1000,62) 연장
// 아래 분기: (500,314)→(519,556) 연장
const TRIPLE_TERRITORIES: TerritoryLayout[] = [
  // 위쪽 영역
  { points: '4,4 996,4 996,62 830,148 500,314 170,148 0,62' },
  // 왼쪽 아래 영역
  { points: '4,62 170,148 500,314 519,556 4,556' },
  // 오른쪽 아래 영역
  { points: '996,62 830,148 500,314 519,556 996,556' },
]

export function getRegionLabels(total: number): RegionLabelLayout[] {
  if (total <= 1) return SINGLE_LABELS
  if (total === 2) return DOUBLE_LABELS
  return TRIPLE_LABELS
}

export function getBoundaryPaths(total: number): BoundaryLayout[] {
  if (total <= 1) return []
  if (total === 2) return DOUBLE_BOUNDARIES
  return TRIPLE_BOUNDARIES
}

export function getTerritoryLayouts(total: number): TerritoryLayout[] {
  if (total <= 1) return SINGLE_TERRITORIES
  if (total === 2) return DOUBLE_TERRITORIES
  return TRIPLE_TERRITORIES
}
