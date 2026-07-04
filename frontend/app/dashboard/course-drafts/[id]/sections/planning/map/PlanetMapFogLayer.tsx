'use client'

import { useMemo, type CSSProperties } from 'react'
import type { RegionAggregate } from '@/components/explorer-plan/explorerPlanTypes'
import { getRegionLabels, getTerritoryLayouts } from './planningMapLayout'
import { getPlanetMapSurfaceTextureStyle } from './planetMapSurfaceTexture'

type RegionFogState = 'unstarted' | 'in_progress' | 'completed'

type PlanetMapFogLayerProps = {
  regions: RegionAggregate[]
  isTransitioning?: boolean
  progressPercent?: number | null
  mapSurfaceAsset?: string | null
}

const layerStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  overflow: 'hidden',
  pointerEvents: 'none',
  zIndex: 0,
  borderRadius: 12,
}

const surfaceStyle: CSSProperties = {
  position: 'absolute',
  left: '50%',
  top: '50%',
  width: '100%',
  aspectRatio: '1000 / 560',
  transform: 'translate(-50%, -50%)',
  overflow: 'hidden',
}

const mistBaseStyle: CSSProperties = {
  position: 'absolute',
  inset: '-8%',
  background:
    'radial-gradient(circle at 18% 24%, rgba(255,255,255,0.32), transparent 28%), radial-gradient(circle at 76% 18%, rgba(246,243,232,0.26), transparent 30%), radial-gradient(circle at 54% 72%, rgba(235,230,214,0.28), transparent 36%), linear-gradient(180deg, rgba(246,240,224,0.22), rgba(92,76,60,0.12))',
  filter: 'blur(9px)',
  opacity: 0.14,
}

const surfaceTextureOverlayStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  background:
    'linear-gradient(180deg, rgba(255,248,230,0.10), rgba(35,27,20,0.12))',
  mixBlendMode: 'multiply',
  opacity: 0.36,
}

const veilStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  background:
    'linear-gradient(135deg, rgba(255,252,240,0.18), rgba(255,252,240,0.02) 44%, rgba(65,55,48,0.12))',
  mixBlendMode: 'screen',
  opacity: 0.16,
}

export default function PlanetMapFogLayer({
  regions,
  isTransitioning = false,
  progressPercent,
  mapSurfaceAsset,
}: PlanetMapFogLayerProps) {
  const labels = useMemo(() => getRegionLabels(regions.length), [regions.length])
  const territories = useMemo(() => getTerritoryLayouts(regions.length), [regions.length])

  if (regions.length === 0) return null

  return (
    <div
      aria-hidden="true"
      style={{
        ...layerStyle,
        opacity: isTransitioning ? 0.48 : 1,
        transition: 'opacity 180ms ease',
      }}
    >
      <div
        data-planet-map-fog-surface="true"
        style={{
          ...surfaceStyle,
          ...getPlanetMapSurfaceTextureStyle(progressPercent, mapSurfaceAsset ?? undefined),
        }}
      >
        <div style={surfaceTextureOverlayStyle} />
        <div style={mistBaseStyle} />
        <div style={veilStyle} />
        {regions.map((regionAgg, index) => {
          const label = labels[index] ?? labels[0]
          const state = getRegionFogState(regionAgg)
          const territory = territories[index] ?? territories[0]
          return (
            <span
              key={regionAgg.region.id}
              style={regionFogPatchStyle(
                label?.x ?? 500,
                label?.y ?? 280,
                state,
                regions.length,
                territory?.points ?? '',
              )}
            />
          )
        })}
        {regions.map((regionAgg, index) => {
          const label = labels[index] ?? labels[0]
          const state = getRegionFogState(regionAgg)
          if (state !== 'completed') return null
          return (
            <span
              key={`${regionAgg.region.id}-reveal`}
              style={regionRevealPatchStyle(
                label?.x ?? 500,
                label?.y ?? 280,
                regions.length,
                territories[index]?.points ?? '',
              )}
            />
          )
        })}
      </div>
    </div>
  )
}

function getRegionFogState(region: RegionAggregate): RegionFogState {
  const nodes = [
    ...region.nodes,
    ...region.subregions.flatMap((subregion) => subregion.nodes),
  ]

  if (nodes.length === 0) return 'unstarted'

  const completed = nodes.filter((node) => getNodeLearningStatus(node) === 'completed').length
  const active = nodes.filter((node) => {
    const status = getNodeLearningStatus(node)
    return status === 'learning' || status === 'completed'
  }).length

  if (completed === nodes.length) return 'completed'
  if (active > 0) return 'in_progress'
  return 'unstarted'
}

function getNodeLearningStatus(node: RegionAggregate['nodes'][number]): string | null {
  return (node as RegionAggregate['nodes'][number] & { learning_status?: string | null }).learning_status ?? null
}

function regionRevealPatchStyle(
  labelX: number,
  labelY: number,
  regionCount: number,
  territoryPoints: string,
): CSSProperties {
  const left = `${(labelX / 1000) * 100}%`
  const top = `${(labelY / 560) * 100}%`
  const baseSize = regionCount <= 1 ? 58 : regionCount === 2 ? 50 : 44

  return {
    position: 'absolute',
    inset: 0,
    opacity: 0.42,
    background: `radial-gradient(ellipse ${baseSize}% ${Math.round(baseSize * 0.74)}% at ${left} ${top}, transparent 0%, rgba(255, 245, 210, 0.08) 38%, rgba(255, 229, 160, 0.16) 62%, transparent 82%)`,
    boxShadow: '0 0 34px rgba(255, 235, 170, 0.20), inset 0 0 26px rgba(255, 249, 225, 0.12)',
    filter: 'blur(2px) saturate(1.12)',
    mixBlendMode: 'screen',
    transition: 'opacity 220ms ease',
    clipPath: toCssPolygon(territoryPoints),
  }
}

function regionFogPatchStyle(
  labelX: number,
  labelY: number,
  state: RegionFogState,
  regionCount: number,
  territoryPoints: string,
): CSSProperties {
  const left = `${(labelX / 1000) * 100}%`
  const top = `${(labelY / 560) * 100}%`
  const baseSize = regionCount <= 1 ? 58 : regionCount === 2 ? 48 : 42
  const size = state === 'completed' ? baseSize : state === 'in_progress' ? baseSize + 6 : baseSize + 14
  const opacity = state === 'completed' ? 0 : state === 'in_progress' ? 0.54 : 0.98
  const blur = 0

  return {
    position: 'absolute',
    inset: 0,
    opacity,
    background: `radial-gradient(ellipse ${size}% ${Math.round(size * 0.72)}% at ${left} ${top}, rgba(255,255,255,0.90) 0%, rgba(246,242,230,0.78) 42%, rgba(210,202,184,0.46) 68%, transparent 84%)`,
    filter: `blur(${blur}px)`,
    transition: 'opacity 220ms ease, filter 220ms ease',
    clipPath: toCssPolygon(territoryPoints),
  }
}

function toCssPolygon(points: string): string | undefined {
  const converted = points
    .trim()
    .split(/\s+/)
    .map((point) => {
      const [rawX, rawY] = point.split(',')
      const x = Number(rawX)
      const y = Number(rawY)
      if (!Number.isFinite(x) || !Number.isFinite(y)) return null
      return `${(x / 1000) * 100}% ${(y / 560) * 100}%`
    })
    .filter((point): point is string => Boolean(point))

  return converted.length >= 3 ? `polygon(${converted.join(', ')})` : undefined
}
