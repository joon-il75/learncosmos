'use client'

import { useState } from 'react'
import type { CSSProperties } from 'react'
import { getDashboardCourseDraftCopy, type DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'

const defaultMapCopy = getDashboardCourseDraftCopy('ko').planning.map

const navButtonBaseStyle: CSSProperties = {
  position: 'absolute',
  top: '55%',
  zIndex: 5,
  width: 42,
  height: 42,
  borderRadius: 999,
  border: '1px solid rgba(91, 70, 45, 0.28)',
  background: 'rgba(255, 248, 235, 0.82)',
  color: 'rgb(77, 56, 35)',
  fontSize: 22,
  fontWeight: 700,
  lineHeight: 1,
  boxShadow: '0 10px 20px rgba(0, 0, 0, 0.12)',
  pointerEvents: 'auto',
  cursor: 'pointer',
  transition:
    'transform 180ms ease, box-shadow 180ms ease, background-color 180ms ease, opacity 180ms ease',
}

const pageIndicatorStyle: CSSProperties = {
  position: 'absolute',
  left: '49.2%',
  bottom: '0.8%',
  transform: 'translateX(-50%)',
  zIndex: 5,
  pointerEvents: 'none',
  color: 'rgb(90, 67, 48)',
  fontSize: 13,
  fontWeight: 600,
  background: 'rgba(255, 248, 235, 0.78)',
  padding: '4px 10px',
  borderRadius: 999,
  transition: 'opacity 180ms ease, transform 180ms ease',
}

const mapLegendStyle: CSSProperties = {
  position: 'absolute',
  left: '57.4%',
  bottom: '0.8%',
  zIndex: 5,
  pointerEvents: 'none',
  display: 'flex',
  alignItems: 'center',
  gap: 8,
  color: 'rgba(82, 59, 39, 0.86)',
  fontSize: 12,
  fontWeight: 700,
  lineHeight: 1,
  background: 'rgba(255, 248, 235, 0.62)',
  padding: '5px 10px',
  borderRadius: 999,
  border: '1px solid rgba(126, 95, 62, 0.14)',
  transition: 'opacity 180ms ease, transform 180ms ease',
}

const legendDividerStyle: CSSProperties = {
  color: 'rgba(126, 95, 62, 0.45)',
  fontWeight: 700,
}

type NavDirection = 'prev' | 'next' | null

export function PlanningMapPager({
  currentPage,
  totalPages,
  isTransitioning,
  onPrev,
  onNext,
  copy = defaultMapCopy,
}: {
  currentPage: number
  totalPages: number
  isTransitioning: boolean
  onPrev: () => void
  onNext: () => void
  copy?: DashboardCourseDraftCopy['planning']['map']
}) {
  const [hoveredNav, setHoveredNav] = useState<NavDirection>(null)
  const [pressedNav, setPressedNav] = useState<NavDirection>(null)

  return (
    <>
      <button
        type="button"
        onClick={onPrev}
        disabled={currentPage === 0 || isTransitioning}
        aria-label={copy.previousPageAria}
        onMouseEnter={() => setHoveredNav('prev')}
        onMouseLeave={() => {
          setHoveredNav(null)
          setPressedNav(null)
        }}
        onMouseDown={() => setPressedNav('prev')}
        onMouseUp={() => setPressedNav(null)}
        style={{
          ...navButtonBaseStyle,
          left: '0.8%',
          opacity: currentPage === 0 ? 0.35 : 1,
          transform:
            pressedNav === 'prev'
              ? 'translateY(1px) scale(0.94)'
              : hoveredNav === 'prev'
                ? 'translateY(-2px) scale(1.06)'
                : 'translateY(0) scale(1)',
          boxShadow:
            hoveredNav === 'prev'
              ? '0 0 0 4px rgba(255,248,235,0.18), 0 14px 28px rgba(0,0,0,0.18)'
              : '0 10px 20px rgba(0,0,0,0.12)',
          background:
            hoveredNav === 'prev'
              ? 'rgba(255, 250, 240, 0.96)'
              : 'rgba(255, 248, 235, 0.82)',
        }}
      >
        <span
          style={{
            display: 'inline-block',
            transform: hoveredNav === 'prev' ? 'translateX(-2px)' : 'translateX(0)',
            transition: 'transform 180ms ease',
          }}
        >
          ←
        </span>
      </button>

      <button
        type="button"
        onClick={onNext}
        disabled={currentPage === totalPages - 1 || isTransitioning}
        aria-label={copy.nextPageAria}
        onMouseEnter={() => setHoveredNav('next')}
        onMouseLeave={() => {
          setHoveredNav(null)
          setPressedNav(null)
        }}
        onMouseDown={() => setPressedNav('next')}
        onMouseUp={() => setPressedNav(null)}
        style={{
          ...navButtonBaseStyle,
          right: '1.6%',
          opacity: currentPage === totalPages - 1 ? 0.35 : 1,
          transform:
            pressedNav === 'next'
              ? 'translateY(1px) scale(0.94)'
              : hoveredNav === 'next'
                ? 'translateY(-2px) scale(1.06)'
                : 'translateY(0) scale(1)',
          boxShadow:
            hoveredNav === 'next'
              ? '0 0 0 4px rgba(255,248,235,0.18), 0 14px 28px rgba(0,0,0,0.18)'
              : '0 10px 20px rgba(0,0,0,0.12)',
          background:
            hoveredNav === 'next'
              ? 'rgba(255, 250, 240, 0.96)'
              : 'rgba(255, 248, 235, 0.82)',
        }}
      >
        <span
          style={{
            display: 'inline-block',
            transform: hoveredNav === 'next' ? 'translateX(2px)' : 'translateX(0)',
            transition: 'transform 180ms ease',
          }}
        >
          →
        </span>
      </button>

      <div
        style={{
          ...pageIndicatorStyle,
          opacity: isTransitioning ? 0.7 : 1,
          transform: isTransitioning
            ? 'translateX(-50%) translateY(2px)'
            : 'translateX(-50%) translateY(0)',
        }}
      >
        {currentPage + 1} / {totalPages}
      </div>

      <div
        aria-hidden="true"
        style={{
          ...mapLegendStyle,
          opacity: isTransitioning ? 0.58 : 1,
          transform: isTransitioning ? 'translateY(2px)' : 'translateY(0)',
        }}
      >
        <span>{copy.subregionLegend} 📁</span>
        <span style={legendDividerStyle}>|</span>
        <span>{copy.explorationLegend} 🎯</span>
        <span style={legendDividerStyle}>|</span>
        <span>{copy.researchLegend} 🔬</span>
      </div>
    </>
  )
}
