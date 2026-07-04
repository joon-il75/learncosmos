'use client'

import { useState } from 'react'
import type { CSSProperties } from 'react'
import type { SubRegionAggregate, ExplorerNode } from '@/components/explorer-plan/explorerPlanTypes'
import PlanetReturnPlanetButton from '@/components/dashboard/PlanetReturnPlanetButton'
import PlanetDetailMapFogLayer from './PlanetDetailMapFogLayer'
import { getDashboardCourseDraftCopy, type DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'

const defaultMapCopy = getDashboardCourseDraftCopy('ko').planning.map

const SVG_W = 800
const SVG_H = 460

const NODE_ROW_SIZE = 5
const ZIGZAG_HALF = 22
const ROW_SPACING = 90
const NODE_H = 64
const NODE_MIN_W = 80
const NODE_GAP = 34

const ARROW_COLOR = 'rgba(174, 45, 32, 0.78)'
const ARROW_SHADOW = 'rgba(112, 30, 20, 0.18)'

function calcNodeWidth(label: string): number {
  let w = 32
  for (const ch of label) {
    w += ch.charCodeAt(0) > 127 ? 13 : 9
  }
  return Math.max(NODE_MIN_W, w)
}

function getNodePositions(nodes: ExplorerNode[]): { x: number; y: number }[] {
  const total = nodes.length
  if (total === 0) return []

  const rowCount = Math.ceil(total / NODE_ROW_SIZE)
  const zoneCenterY = SVG_H / 2

  const rowCenters =
    rowCount === 1
      ? [zoneCenterY]
      : Array.from({ length: rowCount }, (_, r) =>
          zoneCenterY + (r - (rowCount - 1) / 2) * ROW_SPACING,
        )

  return nodes.map((_, i) => {
    const row = Math.floor(i / NODE_ROW_SIZE)
    const col = i % NODE_ROW_SIZE
    const rowStart = row * NODE_ROW_SIZE
    const rowCount2 = Math.min(total - rowStart, NODE_ROW_SIZE)

    const padX = 90
    const spacing = rowCount2 > 1 ? (SVG_W - padX * 2) / (rowCount2 - 1) : 0
    const x = rowCount2 === 1 ? SVG_W / 2 : padX + col * spacing
    const y = rowCenters[row] + (col % 2 === 0 ? -ZIGZAG_HALF : ZIGZAG_HALF)

    return { x, y }
  })
}

function buildArrow(
  from: { x: number; y: number },
  to: { x: number; y: number },
  fromGap: number,
  toGap: number,
) {
  const dx = to.x - from.x
  const dy = to.y - from.y
  const dist = Math.sqrt(dx * dx + dy * dy) || 1
  const ux = dx / dist
  const uy = dy / dist
  const px = -uy
  const py = ux
  const lift = Math.min(42, Math.max(14, dist * 0.13))

  const sx = from.x + ux * fromGap
  const sy = from.y + uy * fromGap
  const ex = to.x - ux * toGap
  const ey = to.y - uy * toGap
  const mx = (sx + ex) / 2 + px * lift
  const my = (sy + ey) / 2 + py * lift - 14
  const lx = to.x - ux * (toGap - 14)
  const ly = to.y - uy * (toGap - 14)

  return {
    d: `M ${sx.toFixed(1)} ${sy.toFixed(1)} Q ${mx.toFixed(1)} ${my.toFixed(1)} ${ex.toFixed(1)} ${ey.toFixed(1)} L ${lx.toFixed(1)} ${ly.toFixed(1)}`,
    lx,
    ly,
  }
}

function truncate(text: string, max: number) {
  return text.length <= max ? text : `${text.slice(0, max - 1)}…`
}

function splitNodeLabelLines(text: string, maxChars = 9): string[] {
  if (text.length <= maxChars) return [text]

  const words = text.split(' ').filter(Boolean)
  if (words.length > 1) {
    const lines: string[] = []
    let current = ''
    for (const word of words) {
      const next = current ? `${current} ${word}` : word
      if (next.length <= maxChars) current = next
      else {
        if (current) lines.push(current)
        current = word
      }
      if (lines.length === 2) break
    }
    if (lines.length < 2 && current) lines.push(current)
    if (lines.length > 0) {
      const last = lines[1] ?? lines[0]
      if (lines.join('').length < text.replace(/\s+/g, '').length) {
        lines[Math.min(1, lines.length - 1)] = truncate(last, maxChars)
      }
      return lines.slice(0, 2)
    }
  }

  return [text.slice(0, maxChars), truncate(text.slice(maxChars), maxChars)]
}

function calcNodeLabelWidth(lines: string[]): number {
  return Math.max(...lines.map(calcNodeWidth))
}

function getNodeCounts(nodes: ExplorerNode[]) {
  return {
    exploration: nodes.filter(n => n.node_type === 'exploration').length,
    research: nodes.filter(n => n.node_type === 'research').length,
  }
}

function isLearningNode(node: ExplorerNode): boolean {
  return (node as ExplorerNode & { learning_status?: string }).learning_status === 'learning'
}

// ── 스타일 ────────────────────────────────────────────────────────────────────
const shellStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  pointerEvents: 'auto',
}

const mapBodyStyle: CSSProperties = {
  position: 'absolute',
  left: '7.4%',
  right: '7.4%',
  top: '23%',
  bottom: '7.2%',
  overflow: 'hidden',
}

const summaryBarStyle: CSSProperties = {
  position: 'absolute',
  left: '50%',
  bottom: '2%',
  transform: 'translateX(-50%)',
  display: 'flex',
  gap: 10,
  alignItems: 'center',
  justifyContent: 'center',
  flexWrap: 'nowrap',
  whiteSpace: 'nowrap',
  padding: '8px 12px',
  borderRadius: 999,
  border: '1px solid rgba(126, 95, 62, 0.18)',
  background: 'rgba(255, 248, 235, 0.56)',
  color: '#5c3f22',
  fontSize: 13,
  fontWeight: 800,
}

const emptyStateStyle: CSSProperties = {
  position: 'absolute',
  left: '50%',
  top: '48%',
  transform: 'translate(-50%, -50%)',
  width: '68%',
  minHeight: 160,
  borderRadius: 16,
  border: '1px dashed rgba(126, 95, 62, 0.38)',
  background: 'rgba(255,248,235,0.24)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: 'rgba(86, 61, 36, 0.82)',
  fontSize: 16,
  fontWeight: 800,
  textAlign: 'center',
  lineHeight: 1.6,
  padding: 20,
}

// ── 컴포넌트 ──────────────────────────────────────────────────────────────────
type PlanningSubRegionDetailProps = {
  subRegion: SubRegionAggregate
  onBack: () => void
  onNodeClick?: (nodeId: string) => void
  textureMapAsset?: string | null
  progressPercent?: number | null
  copy?: DashboardCourseDraftCopy['planning']['map']
}

export function PlanningSubRegionDetail({
  subRegion,
  onBack,
  onNodeClick,
  textureMapAsset,
  progressPercent,
  copy = defaultMapCopy,
}: PlanningSubRegionDetailProps) {
  const [hoveredNodeId, setHoveredNodeId] = useState<string | null>(null)
  const hasNodes = subRegion.nodes.length > 0
  const counts = getNodeCounts(subRegion.nodes)
  const nodePosArr = getNodePositions(subRegion.nodes)

  type FlowItem = { x: number; y: number; gap: number }
  const flowItems: FlowItem[] = subRegion.nodes.map((node, i) => ({
    ...nodePosArr[i],
    gap: calcNodeLabelWidth(splitNodeLabelLines(node.title)) / 2 + 4,
  }))

  return (
    <div style={shellStyle}>
      <PlanetReturnPlanetButton
        onClick={onBack}
        textureMapAsset={textureMapAsset}
        progressPercent={progressPercent ?? 0}
        style={{
          position: 'absolute',
          left: '2.2%',
          top: '4.6%',
        }}
        ariaLabel={copy.returnToRegionMapAria}
      >
        {copy.returnToRegionMapLine1}
        <br />
        {copy.returnToRegionMapLine2}
      </PlanetReturnPlanetButton>

      <div style={mapBodyStyle}>
        {!hasNodes ? (
          <div style={emptyStateStyle}>
            {copy.subregionEmptyLine1}
            <br />
            {copy.subregionEmptyLine2}
            <br />
            {copy.subregionEmptyLine3}
          </div>
        ) : (
          <>
            <PlanetDetailMapFogLayer
              nodes={subRegion.nodes}
              positions={nodePosArr}
              viewBoxWidth={SVG_W}
              viewBoxHeight={SVG_H}
              progressPercent={progressPercent ?? 0}
            />
            <svg
              width="100%"
              height="100%"
              viewBox={`0 0 ${SVG_W} ${SVG_H}`}
              preserveAspectRatio="xMidYMid meet"
              style={{ position: 'relative', zIndex: 1 }}
            >
              <defs>
                <filter id="srd-shadow" x="-20%" y="-20%" width="140%" height="140%">
                  <feDropShadow dx="0" dy="1" stdDeviation="1.8" floodColor="rgba(255,248,230,0.52)" />
                </filter>
                <marker
                  id="srd-arrowhead"
                  markerWidth="10"
                  markerHeight="10"
                  refX="8"
                  refY="5"
                  orient="auto"
                  markerUnits="strokeWidth"
                >
                  <path
                    d="M 1 1 L 9 5 L 1 9"
                    fill="none"
                    stroke={ARROW_COLOR}
                    strokeWidth="1.5"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </marker>
              </defs>

            {/* 1. 흐름 화살표 */}
            <g style={{ pointerEvents: 'none' }}>
              {flowItems.slice(0, -1).map((item, i) => {
                const next = flowItems[i + 1]
                const arrow = buildArrow(item, next, item.gap, next.gap)
                return (
                  <g key={`arrow-${i}`}>
                    <path
                      d={arrow.d}
                      fill="none"
                      stroke={ARROW_SHADOW}
                      strokeWidth="4"
                      strokeLinecap="round"
                      opacity="0.5"
                    />
                    <path
                      d={arrow.d}
                      fill="none"
                      stroke={ARROW_COLOR}
                      strokeWidth="1.8"
                      strokeLinecap="round"
                      strokeDasharray="6 9"
                      markerEnd="url(#srd-arrowhead)"
                    />
                    <circle
                      cx={arrow.lx}
                      cy={arrow.ly}
                      r="4.5"
                      fill="none"
                      stroke={ARROW_COLOR}
                      strokeWidth="1.2"
                      opacity="0.8"
                    />
                  </g>
                )
              })}
            </g>

            {/* 2. 노드 배지 */}
            {subRegion.nodes.map((node, i) => {
              const { x, y } = nodePosArr[i]
              const icon = node.node_type === 'exploration' ? '🎯' : '🔬'
              const labelLines = splitNodeLabelLines(node.title)
              const nw = calcNodeLabelWidth(labelLines)
              const isNodeHovered = hoveredNodeId === node.id
              const isNodeLearning = isLearningNode(node)
              return (
                <g
                  key={node.id}
                  filter="url(#srd-shadow)"
                  style={{
                    transform: isNodeHovered ? 'scale(1.12)' : 'scale(1)',
                    transformOrigin: `${x}px ${y}px`,
                    transition: 'transform 150ms ease',
                    cursor: 'pointer',
                  }}
                  onMouseEnter={() => setHoveredNodeId(node.id)}
                  onMouseLeave={() => setHoveredNodeId(null)}
                  onClick={() => onNodeClick?.(node.id)}
                >
                  <title>{node.title}</title>
                  <rect
                    x={x - nw / 2}
                    y={y - NODE_H / 2}
                    width={nw}
                    height={NODE_H}
                    rx="10"
                    fill={
                      isNodeLearning
                        ? isNodeHovered
                          ? 'rgba(219, 234, 254, 0.98)'
                          : 'rgba(191, 219, 254, 0.92)'
                        : isNodeHovered
                          ? 'rgba(255, 248, 235, 0.98)'
                          : 'rgba(255, 248, 235, 0.90)'
                    }
                    stroke={
                      isNodeLearning
                        ? isNodeHovered
                          ? 'rgba(37, 99, 235, 0.82)'
                          : 'rgba(37, 99, 235, 0.56)'
                        : isNodeHovered
                          ? 'rgba(112, 84, 48, 0.72)'
                          : 'rgba(112, 84, 48, 0.40)'
                    }
                    strokeWidth={isNodeHovered ? '2' : '1.5'}
                  />
                  <text
                    x={x}
                    y={y - 8}
                    textAnchor="middle"
                    dominantBaseline="middle"
                    fontSize="17"
                    style={{ pointerEvents: 'none' }}
                  >
                    {icon}
                  </text>
                  {labelLines.map((line, lineIndex) => (
                    <text
                      key={`${node.id}-label-${lineIndex}`}
                      x={x}
                      y={labelLines.length === 1 ? y + 14 : y + 7 + lineIndex * 16}
                      textAnchor="middle"
                      dominantBaseline="middle"
                      fontSize="14.5"
                      fontWeight="900"
                      fill={isNodeLearning ? '#1D4ED8' : isNodeHovered ? '#2b1909' : '#3d2712'}
                      stroke="rgba(255, 252, 244, 0.82)"
                      strokeWidth="3"
                      strokeLinejoin="round"
                      paintOrder="stroke"
                      style={{ pointerEvents: 'none' }}
                    >
                      {line}
                    </text>
                  ))}
                </g>
              )
            })}
            </svg>
          </>
        )}
      </div>

      {hasNodes && (
        <div style={summaryBarStyle}>
          {copy.formatSubregionSummary(counts.exploration, counts.research).map((label) => (
            <span key={label}>{label}</span>
          ))}
        </div>
      )}
    </div>
  )
}
