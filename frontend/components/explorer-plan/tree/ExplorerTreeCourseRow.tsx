'use client'

import { useState } from 'react'
import {
  courseHeaderDisplayTitleStyle,
  courseHeaderPlanetNameStyle,
  courseHeaderPrefixStyle,
  courseHeaderRowStyle,
  courseHeaderSelectedStyle,
  treeRowHoverStyle,
} from './explorerTreeStyles'
import { getDashboardCourseDraftCopy, type DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'

interface ExplorerTreeCourseRowProps {
  displayTitle: string
  displayPlanetName: string
  isSelected: boolean
  readOnly?: boolean
  onSelect: () => void
  copy?: DashboardCourseDraftCopy['tree']
}

export function ExplorerTreeCourseRow({
  displayTitle,
  displayPlanetName,
  isSelected,
  readOnly = false,
  onSelect,
  copy = getDashboardCourseDraftCopy('ko').tree,
}: ExplorerTreeCourseRowProps) {
  const [isHovered, setIsHovered] = useState(false)
  return (
    <div
      role="button"
      tabIndex={0}
      title={displayTitle}
      onClick={onSelect}
      onKeyDown={(e) => { if (e.key === 'Enter' || e.key === ' ') onSelect() }}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      style={{
        ...courseHeaderRowStyle,
        ...(isHovered && !isSelected ? treeRowHoverStyle : undefined),
        ...(isSelected ? courseHeaderSelectedStyle : undefined),
      }}
    >
      <span style={courseHeaderDisplayTitleStyle} title={displayTitle}>
        <span style={courseHeaderPrefixStyle}>{copy.planetPrefix}</span>
        <span style={courseHeaderPlanetNameStyle}>{displayPlanetName}</span>
      </span>
      {!readOnly && <span style={{ fontSize: '12px', flexShrink: 0 }}>{copy.actionLabel}</span>}
    </div>
  )
}
