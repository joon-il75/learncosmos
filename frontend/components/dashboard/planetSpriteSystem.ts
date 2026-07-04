import type { PlanetStatus } from './planetMapTokens';
import {
  getProgressGridPosition,
  normalizePlanetProgressPercent,
} from '@/lib/world-ui-engine/planetProgressUtils';

export const PLANET_SPRITE_ASSET_PATH = '/textures/planets/planet-spritesheet-140x140.png';
export const PLANET_TEXTURE_MAP_FALLBACK_PATH = '/textures/planets/maps/learnweaver_planet_atlas_basic_seamfixed_2048x768.webp';
export const SPRITE_COLUMNS = 11;
export const SPRITE_ROWS = 7;
export const SPRITE_SHEET_WIDTH = 1540;
export const SPRITE_SHEET_HEIGHT = 980;
export const SPRITE_CELL_WIDTH = SPRITE_SHEET_WIDTH / SPRITE_COLUMNS;
export const SPRITE_CELL_HEIGHT = SPRITE_SHEET_HEIGHT / SPRITE_ROWS;
export const SPRITE_ORIGIN_X = 0;
export const SPRITE_ORIGIN_Y = 0;
export const SPRITE_IMAGE_SCALE_X = 1;
export const SPRITE_IMAGE_SCALE_Y = 1;

export type PlanetTextureFamily = 'ocean' | 'gas' | 'cloud' | 'rock' | 'ice' | 'crater';

export const PLANET_TEXTURE_FAMILY_INDEX: Record<PlanetTextureFamily, number> = {
  ocean: 0,
  gas: 1,
  cloud: 2,
  rock: 3,
  ice: 4,
  crater: 5,
};

export const PLANET_SPRITE_INDEX_TABLE: Record<PlanetTextureFamily, number[]> = {
  ocean: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
  gas: [11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21],
  cloud: [22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32],
  rock: [33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43],
  ice: [44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54],
  crater: [55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65],
};

export function hashPlanetSeed(value: string): number {
  let hash = 0;
  for (let index = 0; index < value.length; index += 1) {
    hash = ((hash << 5) - hash + value.charCodeAt(index)) | 0;
  }
  return Math.abs(hash);
}

export function getPlanetTextureFamily(planet: { id: string }): PlanetTextureFamily {
  const families = Object.keys(PLANET_TEXTURE_FAMILY_INDEX) as PlanetTextureFamily[];
  return families[hashPlanetSeed(planet.id) % families.length] ?? 'rock';
}

export function getTerraformStage(planet: { id: string; status: PlanetStatus }): number {
  const seed = hashPlanetSeed(planet.id);

  if (planet.status === 'draft') {
    return seed % 3;
  }
  if (planet.status === 'ready') {
    return 3 + (seed % 3);
  }
  if (planet.status === 'learning') {
    return 6 + (seed % 4);
  }
  return 10;
}

export function getPlanetSpritePosition(family: PlanetTextureFamily, stage: number) {
  const maxStage = Math.max(0, (PLANET_SPRITE_INDEX_TABLE[family]?.length ?? 1) - 1);
  const safeStage = Math.max(0, Math.min(maxStage, stage));
  const spriteIndex = PLANET_SPRITE_INDEX_TABLE[family][safeStage] ?? PLANET_SPRITE_INDEX_TABLE[family][0] ?? 0;
  const spriteColumn = spriteIndex % SPRITE_COLUMNS;
  const spriteRow = Math.floor(spriteIndex / SPRITE_COLUMNS);
  return {
    x: spriteColumn * SPRITE_CELL_WIDTH,
    y: spriteRow * SPRITE_CELL_HEIGHT,
  };
}

export function getPlanetSpriteBackgroundStyle(
  family: PlanetTextureFamily,
  stage: number,
  renderedSize: number,
) {
  const position = getPlanetSpritePosition(family, stage);
  const scaleX = renderedSize / SPRITE_CELL_WIDTH;
  const scaleY = renderedSize / SPRITE_CELL_HEIGHT;
  const sheetWidth = SPRITE_SHEET_WIDTH * SPRITE_IMAGE_SCALE_X * scaleX;
  const sheetHeight = SPRITE_SHEET_HEIGHT * SPRITE_IMAGE_SCALE_Y * scaleY;

  return {
    backgroundImage: `url('${PLANET_SPRITE_ASSET_PATH}')`,
    backgroundRepeat: 'no-repeat',
    backgroundSize: `${sheetWidth}px ${sheetHeight}px`,
    backgroundPosition: `-${(position.x + SPRITE_ORIGIN_X) * scaleX}px -${(position.y + SPRITE_ORIGIN_Y) * scaleY}px`,
  };
}

export function getPlanetTypeSpriteBackgroundStyle(
  assetPath: string,
  _status: PlanetStatus,
  progressPercent: number | null,
  _renderedSize: number,
) {
  const columns = 4;
  const rows = 3;
  const effectiveProgressPercent = normalizePlanetProgressPercent(progressPercent);
  const gridPosition = getProgressGridPosition(effectiveProgressPercent);
  const row = gridPosition.row - 1;
  const column = gridPosition.col - 1;
  const backgroundPositionX = (column / (columns - 1)) * 100;
  const backgroundPositionY = (row / (rows - 1)) * 100;

  return {
    backgroundImage: `url('${assetPath}')`,
    backgroundRepeat: 'no-repeat',
    // Planet type sheets should be treated as a 4x3 grid regardless of export resolution.
    backgroundSize: `${columns * 100}% ${rows * 100}%`,
    backgroundPosition: `${backgroundPositionX}% ${backgroundPositionY}%`,
  };
}

export function getPlanetTextureMapBackgroundStyle(
  assetPath: string,
  progressPercent: number | null,
  isHovered = false,
) {
  const columns = 4;
  const rows = 3;
  const effectiveProgressPercent = normalizePlanetProgressPercent(progressPercent);
  const gridPosition = getProgressGridPosition(effectiveProgressPercent);
  const row = gridPosition.row - 1;
  const column = gridPosition.col - 1;
  const backgroundPositionX = (column / (columns - 1)) * 100;
  const backgroundPositionY = (row / (rows - 1)) * 100;
  const shade = isHovered
    ? 'radial-gradient(circle at 38% 28%, rgba(255,255,255,0.36), rgba(27, 69, 148, 0.18) 44%, rgba(3, 8, 18, 0.48) 78%)'
    : 'radial-gradient(circle at 38% 28%, rgba(255,255,255,0.26), rgba(27, 69, 148, 0.14) 44%, rgba(3, 8, 18, 0.54) 78%)';

  return {
    background: `${shade}, url('${assetPath}')`,
    backgroundImage: `${shade}, url('${assetPath}')`,
    backgroundRepeat: 'no-repeat, no-repeat',
    backgroundSize: `100% 100%, ${columns * 100}% ${rows * 100}%`,
    backgroundPosition: `center, ${backgroundPositionX}% ${backgroundPositionY}%`,
  };
}
