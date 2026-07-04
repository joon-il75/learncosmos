'use client'

import type { CSSProperties } from 'react'
import type { ExplorerNode } from '@/components/explorer-plan/explorerPlanTypes'
import { getPlanetMapSurfaceTextureStyle } from './planetMapSurfaceTexture'

type DetailFogNode = Pick<ExplorerNode, 'id'> & {
  learning_status?: string | null
}

type DetailFogPoint = {
  id: string
  x: number
  y: number
  state: 'unstarted' | 'in_progress' | 'completed'
}

type PlanetDetailMapFogLayerProps = {
  nodes: DetailFogNode[]
  positions: { x: number; y: number }[]
  viewBoxWidth: number
  viewBoxHeight: number
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

const mistBaseStyle: CSSProperties = {
  position: 'absolute',
  inset: '-8%',
  background:
    'radial-gradient(circle at 22% 22%, rgba(255,255,255,0.26), transparent 30%), radial-gradient(circle at 78% 28%, rgba(246,243,232,0.22), transparent 32%), radial-gradient(circle at 48% 74%, rgba(235,230,214,0.24), transparent 38%), linear-gradient(180deg, rgba(246,240,224,0.20), rgba(92,76,60,0.10))',
  filter: 'blur(8px)',
  opacity: 0,
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
    'linear-gradient(135deg, rgba(255,252,240,0.12), rgba(255,252,240,0.02) 45%, rgba(65,55,48,0.10))',
  mixBlendMode: 'screen',
  opacity: 0,
}

export default function PlanetDetailMapFogLayer({
  nodes,
  positions,
  viewBoxWidth,
  viewBoxHeight,
  progressPercent,
  mapSurfaceAsset,
}: PlanetDetailMapFogLayerProps) {
  if (positions.length === 0) return null

  const fogPoints: DetailFogPoint[] = positions.map((position, index) => {
    const node = nodes[index]
    return {
      id: node?.id ?? `point-${index}`,
      x: position.x,
      y: position.y,
      state: getPointState(node),
    }
  })

  return (
    <div aria-hidden="true" style={layerStyle}>
      <div
        data-planet-detail-fog-surface="true"
        style={{
          ...detailSurfaceStyle(viewBoxWidth, viewBoxHeight),
          ...getPlanetMapSurfaceTextureStyle(progressPercent, mapSurfaceAsset ?? undefined),
        }}
      >
        <div style={surfaceTextureOverlayStyle} />
        <div style={mistBaseStyle} />
        <div style={veilStyle} />
        {fogPoints.map((point) => (
          <span
            key={point.id}
            style={detailFogPatchStyle(point, viewBoxWidth, viewBoxHeight)}
          />
        ))}
      </div>
    </div>
  )
}

function getPointState(node: DetailFogNode | undefined): DetailFogPoint['state'] {
  if (!node) return 'unstarted'
  if (node.learning_status === 'completed') return 'completed'
  if (node.learning_status === 'learning') return 'in_progress'
  return 'unstarted'
}

function detailSurfaceStyle(viewBoxWidth: number, viewBoxHeight: number): CSSProperties {
  return {
    position: 'absolute',
    left: '50%',
    top: '50%',
    width: '100%',
    aspectRatio: `${viewBoxWidth} / ${viewBoxHeight}`,
    transform: 'translate(-50%, -50%)',
    overflow: 'hidden',
  }
}

function detailFogPatchStyle(
  point: DetailFogPoint,
  viewBoxWidth: number,
  viewBoxHeight: number,
): CSSProperties {
  const left = `${(point.x / viewBoxWidth) * 100}%`
  const top = `${(point.y / viewBoxHeight) * 100}%`
  const opacity = point.state === 'completed' ? 0 : point.state === 'in_progress' ? 0.52 : 0.98
  const blur = point.state === 'completed' ? 0 : point.state === 'in_progress' ? 6 : 7
  const width = point.state === 'completed' ? 34 : point.state === 'in_progress' ? 40 : 48

  return {
    position: 'absolute',
    left,
    top,
    width: `${width}%`,
    aspectRatio: '1 / 0.72',
    transform: 'translate(-50%, -50%)',
    borderRadius: '50%',
    opacity,
    background:
      'radial-gradient(ellipse at center, rgba(255,255,255,0.90) 0%, rgba(246,242,230,0.76) 42%, rgba(210,202,184,0.44) 68%, transparent 84%)',
    filter: `blur(${blur}px)`,
    transition: 'opacity 220ms ease, filter 220ms ease',
  }
}
