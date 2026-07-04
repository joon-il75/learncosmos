import type { LumiState } from './lumiTypes';

// 스프라이트 시트 규격 (public/images/lumi.webp 기준)
export const LUMI_SPRITE_SHEET_WIDTH = 1619;
export const LUMI_SPRITE_SHEET_HEIGHT = 971;
export const LUMI_SPRITE_CELL_WIDTH = LUMI_SPRITE_SHEET_WIDTH / 5;
export const LUMI_SPRITE_CELL_HEIGHT = LUMI_SPRITE_SHEET_HEIGHT / 3;
export const LUMI_SPRITE_COLUMNS = LUMI_SPRITE_SHEET_WIDTH / LUMI_SPRITE_CELL_WIDTH;   // 5
export const LUMI_SPRITE_ROWS = LUMI_SPRITE_SHEET_HEIGHT / LUMI_SPRITE_CELL_HEIGHT;    // 3

export interface LumiSpriteTuning {
  backgroundSizeX: number;
  backgroundSizeY: number;
  backgroundPositionX: number;
  backgroundPositionY: number;
}

// 상태별 스프라이트 시트 셀 좌표 (1-based, row × col)
// 교체용 시트는 LumiState 코드 순서대로 5열 x 3행에 배치한다.
export const lumiSpriteMap: Record<LumiState, { row: number; col: number }> = {
  idle:           { row: 1, col: 1 },
  thinking:       { row: 1, col: 2 },
  happy:          { row: 1, col: 3 },
  celebrate:      { row: 1, col: 4 },
  victory:        { row: 1, col: 5 },
  encourage:      { row: 2, col: 1 },
  surprise:       { row: 2, col: 2 },
  curious:        { row: 2, col: 3 },
  'raise-hand':   { row: 2, col: 4 },
  sleepy:         { row: 2, col: 5 },
  exploring:      { row: 3, col: 1 },
  focus:          { row: 3, col: 2 },
  'planet-point': { row: 3, col: 3 },
  'planet-hold':  { row: 3, col: 4 },
  'note-read':    { row: 3, col: 5 },
};

export const LUMI_DEFAULT_SPRITE_TUNING: Pick<LumiSpriteTuning, "backgroundSizeX" | "backgroundSizeY"> = {
  backgroundSizeX: 500,
  backgroundSizeY: 300,
};

// 새 lumi.webp는 5열 x 3행의 균일 셀 시트이므로 셀 좌표 기반 기본 position을 사용한다.
export const lumiSpriteTuningMap: Partial<Record<LumiState, LumiSpriteTuning>> = {};

export function getSpriteTuningByCoordinate(row: number, col: number): LumiSpriteTuning {
  const entry = Object.entries(lumiSpriteMap).find(([, coord]) => coord.row === row && coord.col === col);
  if (!entry) return getDefaultSpriteTuning('idle');
  return getDefaultSpriteTuning(entry[0] as LumiState);
}

export function getDefaultSpriteTuning(state: LumiState): LumiSpriteTuning {
  const custom = lumiSpriteTuningMap[state];
  if (custom) return custom;

  const sprite = lumiSpriteMap[state];
  const baseX = sprite.col === 1 ? 0 : ((sprite.col - 1) / (LUMI_SPRITE_COLUMNS - 1)) * 100;
  const baseY = sprite.row === 1 ? 0 : ((sprite.row - 1) / (LUMI_SPRITE_ROWS - 1)) * 100;

  return {
    backgroundSizeX: LUMI_DEFAULT_SPRITE_TUNING.backgroundSizeX,
    backgroundSizeY: LUMI_DEFAULT_SPRITE_TUNING.backgroundSizeY,
    backgroundPositionX: baseX,
    backgroundPositionY: baseY,
  };
}
