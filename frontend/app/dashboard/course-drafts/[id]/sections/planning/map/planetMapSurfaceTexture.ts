import type { CSSProperties } from 'react'
import {
  getProgressGridPosition,
  normalizePlanetProgressPercent,
} from '@/lib/world-ui-engine/planetProgressUtils'

export const PLANET_MAP_SURFACE_FALLBACK_PATH =
  '/textures/planets/map-surfaces/learnweaver_planet_map_surface_basic_8stage_4x2.webp'

export function getPlanetMapSurfaceTextureStyle(
  progressPercent: number | null | undefined,
  assetPath = PLANET_MAP_SURFACE_FALLBACK_PATH,
): CSSProperties {
  const effectiveProgress = normalizePlanetProgressPercent(progressPercent)
  const gridPosition = getProgressGridPosition(effectiveProgress)
  const column = gridPosition.col - 1
  const row = gridPosition.row - 1
  const backgroundPositionX = (column / 3) * 100
  const backgroundPositionY = row * 100

  return {
    backgroundImage: `url('${assetPath}')`,
    backgroundRepeat: 'no-repeat',
    backgroundSize: '400% 200%',
    backgroundPosition: `${backgroundPositionX}% ${backgroundPositionY}%`,
  }
}
