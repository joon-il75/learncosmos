'use client'

import { useMemo, useState } from 'react'
import type { CSSProperties } from 'react'

import type { RegionAggregate } from '@/components/explorer-plan/explorerPlanTypes'
import { getBoundaryPaths, getRegionLabels, getTerritoryLayouts } from './planningMapLayout'
import { getDashboardCourseDraftCopy, type DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'

const defaultMapCopy = getDashboardCourseDraftCopy('ko').planning.map

type NavDirection = 'prev' | 'next' | null

type PlanningMapSvgProps = {
  regions: RegionAggregate[]
  regionIndexOffset?: number
  hasPreviousPage?: boolean
  hasNextPage?: boolean
  isTransitioning: boolean
  direction: NavDirection
  onRegionClick?: (regionId: string) => void
  copy?: DashboardCourseDraftCopy['planning']['map']
}

const svgStyle: CSSProperties = {
  display: 'block',
  width: '100%',
  height: '100%',
}

const TITLE_FONT_SIZE = 24
const BODY_FONT_SIZE = 17
const TITLE_LINE_HEIGHT = 30
const BODY_LINE_HEIGHT = 24
const CHARS_PER_LINE = 12
const FLOW_ARROW_COLOR = 'rgba(174, 45, 32, 0.78)'
const FLOW_ARROW_SHADOW = 'rgba(112, 30, 20, 0.18)'

function getExplorationCount(regionAgg: RegionAggregate): number {
  return [
    ...regionAgg.nodes,
    ...regionAgg.subregions.flatMap(s => s.nodes),
  ].filter(n => n.node_type === 'exploration').length
}

function getResearchCount(regionAgg: RegionAggregate): number {
  return [
    ...regionAgg.nodes,
    ...regionAgg.subregions.flatMap(s => s.nodes),
  ].filter(n => n.node_type === 'research').length
}

function splitTitleLines(title: string, maxChars: number): string[] {
  if (title.length <= maxChars) return [title]

  const words = title.split(' ')
  const lines: string[] = []
  let current = ''

  for (const word of words) {
    const next = current ? `${current} ${word}` : word
    if (next.length <= maxChars) {
      current = next
    } else {
      if (current) lines.push(current)
      current = word
    }
  }
  if (current) lines.push(current)

  if (lines.length === 1) {
    const mid = Math.ceil(title.length / 2)
    return [title.slice(0, mid), title.slice(mid)].filter(Boolean)
  }
  if (lines.length > 2) {
    const mid = Math.ceil(lines.length / 2)
    return [lines.slice(0, mid).join(' '), lines.slice(mid).join(' ')]
  }
  return lines
}

function getFlowArrowPath(from: { x: number; y: number }, to: { x: number; y: number }) {
  const dx = to.x - from.x
  const dy = to.y - from.y
  const distance = Math.sqrt(dx * dx + dy * dy) || 1
  const unitX = dx / distance
  const unitY = dy / distance
  const perpendicularX = -unitY
  const perpendicularY = unitX
  const arcLift = Math.min(92, Math.max(46, distance * 0.16))
  const startGap = 72
  const endGap = 84

  const startX = from.x + unitX * startGap
  const startY = from.y + unitY * startGap
  const endX = to.x - unitX * endGap
  const endY = to.y - unitY * endGap
  const midX = (startX + endX) / 2
  const midY = (startY + endY) / 2
  const controlX = midX + perpendicularX * arcLift
  const controlY = midY + perpendicularY * arcLift - 34
  const landingX = to.x - unitX * 46
  const landingY = to.y - unitY * 46

  return {
    d: `M ${startX.toFixed(1)} ${startY.toFixed(1)} Q ${controlX.toFixed(1)} ${controlY.toFixed(1)} ${endX.toFixed(1)} ${endY.toFixed(1)} L ${landingX.toFixed(1)} ${landingY.toFixed(1)}`,
    landingX,
    landingY,
  }
}

export function PlanningMapSvg({
  regions,
  regionIndexOffset = 0,
  hasPreviousPage = false,
  hasNextPage = false,
  isTransitioning,
  direction,
  onRegionClick,
  copy = defaultMapCopy,
}: PlanningMapSvgProps) {
  const [hoveredRegionId, setHoveredRegionId] = useState<string | null>(null)
  const [focusedRegionId, setFocusedRegionId] = useState<string | null>(null)

  const labels = useMemo(() => getRegionLabels(regions.length), [regions.length])
  const boundaries = useMemo(() => getBoundaryPaths(regions.length), [regions.length])
  const territories = useMemo(() => getTerritoryLayouts(regions.length), [regions.length])
  const showHub = regions.length >= 3

  return (
    <svg
      width="100%"
      height="100%"
      viewBox="0 0 1000 560"
      preserveAspectRatio="xMidYMid meet"
      style={{
        ...svgStyle,
        opacity: isTransitioning ? 0.42 : 1,
        transform:
          isTransitioning && direction === 'next'
            ? 'translateX(-18px)'
            : isTransitioning && direction === 'prev'
              ? 'translateX(18px)'
              : 'translateX(0)',
        transition: 'opacity 180ms ease, transform 180ms ease',
      }}
    >
      <defs>
        <filter id="region-text-shadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="1" stdDeviation="2" floodColor="rgba(255,248,230,0.55)" />
        </filter>
        <filter id="hub-shadow" x="-40%" y="-40%" width="180%" height="180%">
          <feDropShadow dx="0" dy="2" stdDeviation="3" floodColor="rgba(70,52,31,0.18)" />
        </filter>
        <marker
          id="region-flow-arrowhead"
          markerWidth="10"
          markerHeight="10"
          refX="8"
          refY="5"
          orient="auto"
          markerUnits="strokeWidth"
        >
          <path d="M 1 1 L 9 5 L 1 9" fill="none" stroke={FLOW_ARROW_COLOR} strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
        </marker>
        {/* 외부 패널에서 제한한 폴리곤 표시 영역 전체를 밝은 채움 영역으로 사용 */}
        <clipPath id="territory-clip">
          <rect x="0" y="0" width="1000" height="560" />
        </clipPath>
      </defs>

      {/* 1. 지역 영역 채움 */}
      <g clipPath="url(#territory-clip)">
        {regions.map((regionAgg, index) => {
          const territory = territories[index] ?? territories[0]
          const label = labels[index] ?? labels[0]
          const regionId = regionAgg.region.id
          const isHighlighted = hoveredRegionId === regionId || focusedRegionId === regionId
          const regionTitle = copy.formatRegionTitle(regionIndexOffset + index, regionAgg.region.name)

          return (
            <polygon
              key={`territory-${regionId}`}
              points={territory.points}
              fill={isHighlighted ? 'rgba(255, 248, 235, 0.82)' : 'rgba(255, 248, 235, 0.20)'}
              stroke="none"
              style={{
                cursor: 'pointer',
                pointerEvents: 'auto',
                transition: 'fill 160ms ease, transform 160ms ease',
                transform: isHighlighted ? 'scale(1.2)' : 'scale(1)',
                transformOrigin: `${label.x}px ${label.y}px`,
              }}
              onMouseEnter={() => setHoveredRegionId(regionId)}
              onMouseLeave={() => setHoveredRegionId(null)}
              onClick={() => onRegionClick?.(regionId)}
            >
              <title>{regionTitle}</title>
            </polygon>
          )
        })}
      </g>

      {/* 2. 경계선 — territory 위에 겹쳐 구분선으로 표시 */}
      <g style={{ opacity: hoveredRegionId ? 0 : 1, transition: 'opacity 160ms ease', pointerEvents: 'none' }}>
        {boundaries.map(boundary => (
          <path
            key={boundary.id}
            d={boundary.d}
            fill="none"
            stroke="rgba(108, 86, 57, 0.84)"
            strokeWidth="2.5"
            strokeDasharray="5 8"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        ))}
      </g>

      {/* 3. 허브 도트 */}
      {showHub ? (
        <g filter="url(#hub-shadow)" style={{ pointerEvents: 'none' }}>
          <circle cx="500" cy="314" r="7" fill="rgba(116, 88, 58, 0.18)" />
          <circle cx="500" cy="314" r="3.5" fill="rgba(116, 88, 58, 0.28)" />
        </g>
      ) : null}

      {/* 4. 지역 순서 흐름 — 얇은 붉은 점프형 화살표 */}
      {regions.length > 0 ? (
        <g style={{ pointerEvents: 'none', opacity: hoveredRegionId ? 0 : 1, transition: 'opacity 160ms ease' }}>
          {hasPreviousPage && labels[0] ? (() => {
            const arrow = getFlowArrowPath(
              { x: -60, y: Math.max(84, labels[0].y - 26) },
              labels[0],
            )

            return (
              <g key="flow-from-previous-page">
                <path
                  d={arrow.d}
                  fill="none"
                  stroke={FLOW_ARROW_SHADOW}
                  strokeWidth="5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  opacity="0.48"
                />
                <path
                  d={arrow.d}
                  fill="none"
                  stroke={FLOW_ARROW_COLOR}
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeDasharray="7 10"
                  markerEnd="url(#region-flow-arrowhead)"
                  opacity="0.9"
                />
                <circle
                  cx={arrow.landingX}
                  cy={arrow.landingY}
                  r="5.4"
                  fill="none"
                  stroke={FLOW_ARROW_COLOR}
                  strokeWidth="1.3"
                  opacity="0.76"
                />
              </g>
            )
          })() : null}

          {regions.slice(0, -1).map((regionAgg, index) => {
            const from = labels[index] ?? labels[0]
            const to = labels[index + 1] ?? labels[0]
            const arrow = getFlowArrowPath(from, to)

            return (
              <g key={`flow-${regionAgg.region.id}-${regions[index + 1]?.region.id ?? index}`}>
                <path
                  d={arrow.d}
                  fill="none"
                  stroke={FLOW_ARROW_SHADOW}
                  strokeWidth="5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  opacity="0.55"
                />
                <path
                  d={arrow.d}
                  fill="none"
                  stroke={FLOW_ARROW_COLOR}
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeDasharray="7 10"
                  markerEnd="url(#region-flow-arrowhead)"
                />
                <circle
                  cx={arrow.landingX}
                  cy={arrow.landingY}
                  r="6"
                  fill="none"
                  stroke={FLOW_ARROW_COLOR}
                  strokeWidth="1.4"
                  opacity="0.82"
                />
                <circle
                  cx={arrow.landingX}
                  cy={arrow.landingY}
                  r="2.2"
                  fill={FLOW_ARROW_COLOR}
                  opacity="0.72"
                />
              </g>
            )
          })}

          {hasNextPage && labels[regions.length - 1] ? (() => {
            const lastLabel = labels[regions.length - 1]
            const arrow = getFlowArrowPath(
              lastLabel,
              { x: 1060, y: Math.max(84, lastLabel.y - 28) },
            )

            return (
              <g key="flow-to-next-page">
                <path
                  d={arrow.d}
                  fill="none"
                  stroke={FLOW_ARROW_SHADOW}
                  strokeWidth="5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  opacity="0.48"
                />
                <path
                  d={arrow.d}
                  fill="none"
                  stroke={FLOW_ARROW_COLOR}
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeDasharray="7 10"
                  markerEnd="url(#region-flow-arrowhead)"
                  opacity="0.9"
                />
              </g>
            )
          })() : null}
        </g>
      ) : null}

      {/* 5. 지역 텍스트 레이블 — 영역 위에 표시 */}
      {regions.map((regionAgg, index) => {
        const label = labels[index] ?? labels[0]
        const regionId = regionAgg.region.id
        const isHighlighted = hoveredRegionId === regionId || focusedRegionId === regionId
        const subRegionCount = regionAgg.subregions.length
        const explorationCount = getExplorationCount(regionAgg)
        const researchCount = getResearchCount(regionAgg)
        const title = copy.formatRegionTitle(regionIndexOffset + index, regionAgg.region.name)
        const titleLines = splitTitleLines(title, CHARS_PER_LINE)
        const metaLabel = copy.formatRegionMeta(subRegionCount, explorationCount, researchCount)

        const titleBlockHeight = titleLines.length * TITLE_LINE_HEIGHT
        const firstTitleY = label.y - (titleBlockHeight / 2) - 8 - BODY_LINE_HEIGHT * 0.5
        const metaY = firstTitleY + titleBlockHeight + 10 + BODY_FONT_SIZE / 2
        const accentY = firstTitleY + titleBlockHeight - 2

        return (
          <g
            key={`label-${regionId}`}
            filter="url(#region-text-shadow)"
            role="button"
            tabIndex={0}
            aria-label={copy.formatRegionAria(title, subRegionCount, explorationCount, researchCount)}
            style={{
              pointerEvents: 'auto',
              cursor: 'pointer',
              opacity: isHighlighted ? 1 : 0.9,
              outline: 'none',
              transition: 'opacity 160ms ease, transform 160ms ease',
              transform: isHighlighted ? 'scale(1.2)' : 'scale(1)',
              transformOrigin: `${label.x}px ${label.y}px`,
            }}
            onMouseEnter={() => setHoveredRegionId(regionId)}
            onMouseLeave={() => setHoveredRegionId(null)}
            onFocus={() => setFocusedRegionId(regionId)}
            onBlur={() => setFocusedRegionId(null)}
            onClick={() => onRegionClick?.(regionId)}
            onKeyDown={event => {
              if (event.key === 'Enter' || event.key === ' ') {
                event.preventDefault()
                onRegionClick?.(regionId)
              }
            }}
          >
            <title>{title}</title>
            {isHighlighted ? (
              <g style={{ pointerEvents: 'none' }}>
                <line
                  x1={label.x - 74}
                  x2={label.x + 74}
                  y1={accentY}
                  y2={accentY}
                  stroke="rgba(136, 83, 31, 0.58)"
                  strokeWidth="2.2"
                  strokeLinecap="round"
                />
                <circle cx={label.x - 84} cy={accentY} r="2.5" fill="rgba(136, 83, 31, 0.58)" />
                <circle cx={label.x + 84} cy={accentY} r="2.5" fill="rgba(136, 83, 31, 0.58)" />
              </g>
            ) : null}

            {titleLines.map((line, lineIndex) => (
              <text
                key={`${regionId}-title-${lineIndex}`}
                x={label.x}
                y={firstTitleY + lineIndex * TITLE_LINE_HEIGHT}
                fill={isHighlighted ? '#3f2a14' : '#654520'}
                fontSize={TITLE_FONT_SIZE}
                fontWeight={isHighlighted ? '900' : '800'}
                textAnchor="middle"
                dominantBaseline="middle"
                stroke="rgba(255, 250, 238, 0.78)"
                strokeWidth="3"
                strokeLinejoin="round"
                style={{ pointerEvents: 'none', paintOrder: 'stroke fill' }}
              >
                {line}
              </text>
            ))}

            <text
              x={label.x}
              y={metaY}
              fill={isHighlighted ? 'rgba(55, 39, 24, 0.96)' : 'rgba(85, 60, 36, 0.88)'}
              fontSize={BODY_FONT_SIZE}
              fontWeight={isHighlighted ? '800' : '700'}
              textAnchor="middle"
              dominantBaseline="middle"
              stroke="rgba(255, 250, 238, 0.72)"
              strokeWidth="2.4"
              strokeLinejoin="round"
              style={{ pointerEvents: 'none', paintOrder: 'stroke fill' }}
            >
              {metaLabel}
            </text>
          </g>
        )
      })}
    </svg>
  )
}
