'use client'

import type { CSSProperties } from 'react'
import { getDashboardCourseDraftCopy, type DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'

const defaultMapCopy = getDashboardCourseDraftCopy('ko').planning.map

type PlanningMapLevel = 'course' | 'region' | 'subregion'

const headerStyle: CSSProperties = {
  position: 'absolute',
  left: '21.0%',
  top: '4.6%',
  width: '70.5%',
  display: 'flex',
  flexDirection: 'column',
  gap: 2,
  textShadow: '0 2px 10px rgba(0,0,0,0.14)',
  zIndex: 3,
}

const badgeStyle: CSSProperties = {
  alignSelf: 'flex-start',
  padding: '4px 10px',
  borderRadius: 999,
  background: 'rgba(140, 115, 80, 0.78)',
  color: '#fffaf0',
  fontSize: 13,
  fontWeight: 700,
  letterSpacing: '0.45px',
}

const titleStyle: CSSProperties = {
  margin: 0,
  fontSize: 42,
  fontWeight: 500,
  color: '#2f2419',
  fontFamily: "'Nanum Pen Script', cursive",
  letterSpacing: '0.4px',
  lineHeight: 1,
}

const descStyle: CSSProperties = {
  margin: 0,
  fontSize: 13,
  lineHeight: 1.35,
  color: '#3b2d20',
  fontWeight: 600,
  background: 'rgba(255, 248, 235, 0.58)',
  padding: '4px 8px',
  borderRadius: 6,
  display: 'inline-block',
  maxWidth: '100%',
}


function getCourseTitle(
  mapLevel: PlanningMapLevel,
  copy: DashboardCourseDraftCopy['planning']['map'],
  selectedRegionTitle?: string | null,
  selectedSubRegionTitle?: string | null,
) {
  if (mapLevel === 'course') return copy.courseMapTitle
  if (mapLevel === 'subregion') return selectedSubRegionTitle ?? copy.subregionMapFallback
  return selectedRegionTitle ?? copy.regionMapFallback
}

function getCourseDescription(
  mapLevel: PlanningMapLevel,
  courseTitle?: string | null,
  selectedRegionTitle?: string | null,
  selectedSubRegionTitle?: string | null,
) {
  if (mapLevel === 'course') return courseTitle ?? ''
  if (mapLevel === 'subregion') {
    const parts = [courseTitle, selectedRegionTitle, selectedSubRegionTitle].filter(Boolean)
    return parts.join(' / ')
  }
  const parts = [courseTitle, selectedRegionTitle].filter(Boolean)
  return parts.join(' / ')
}

export function PlanningMapHeader({
  mapLevel,
  courseTitle,
  selectedRegionTitle,
  selectedSubRegionTitle,
  copy = defaultMapCopy,
}: {
  mapLevel: PlanningMapLevel
  courseTitle?: string | null
  selectedRegionTitle?: string | null
  selectedSubRegionTitle?: string | null
  copy?: DashboardCourseDraftCopy['planning']['map']
}) {
  return (
    <>
      <header style={headerStyle}>
        <div style={badgeStyle}>{copy.headerBadge}</div>
        <h2 style={titleStyle}>{getCourseTitle(mapLevel, copy, selectedRegionTitle, selectedSubRegionTitle)}</h2>
        <p style={descStyle}>{getCourseDescription(mapLevel, courseTitle, selectedRegionTitle, selectedSubRegionTitle)}</p>
      </header>

    </>
  )
}
