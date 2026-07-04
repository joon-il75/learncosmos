'use client'

import { useState } from 'react'
import type { RegionAggregate } from '../explorerPlanTypes'
import type { PendingItemState } from '../useExplorerPlan'
import { getDashboardCourseDraftCopy, type DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'
import {
  treeInlineInputStyle,
  treeRowBaseStyle,
  treeRowSelectedStyle,
  treeRowInactiveStyle,
  treeRowHoverStyle,
  treeToggleBtnStyle,
  treeRowLabelStyle,
  treeRowPendingCreateStyle,
  treeRowPendingDeleteStyle,
  treeRowPendingMoveStyle,
  pendingCreateBadgeStyle,
  pendingDeleteBadgeStyle,
  pendingMoveIconStyle,
} from './explorerTreeStyles'

interface ExplorerTreeRegionRowProps {
  regionAgg: RegionAggregate
  regionIndex: number
  isOpen: boolean
  isSelected: boolean
  readOnly?: boolean
  compact?: boolean
  pendingState: PendingItemState
  titleInput: string
  onTitleChange: (value: string) => void
  isEditingDisabled: boolean
  onToggle: (regionId: string) => void
  onSelect: (regionId: string) => void
  copy?: DashboardCourseDraftCopy['tree']
  isOnboardingHighlighted?: boolean
  isOnboardingDisabled?: boolean
}

export function ExplorerTreeRegionRow({
  regionAgg,
  regionIndex,
  isOpen,
  isSelected,
  readOnly = false,
  compact = false,
  pendingState,
  titleInput,
  onTitleChange,
  isEditingDisabled,
  onToggle,
  onSelect,
  copy = getDashboardCourseDraftCopy('ko').tree,
  isOnboardingHighlighted = false,
  isOnboardingDisabled = false,
}: ExplorerTreeRegionRowProps) {
  const [isHovered, setIsHovered] = useState(false)
  const regionId = regionAgg.region.id
  const displayNumber = regionIndex + 1
  const name = regionAgg.region.name.trim()
  const label = copy.regionLabel(displayNumber, name)
  const isInactive = regionAgg.region.status === 'inactive'
  const isPendingDeleted = pendingState === 'deleted'
  const isPendingCreated = pendingState === 'created'
  const isPendingMoved = pendingState === 'moved'
  const labelColor = isPendingDeleted ? '#A43124' : isPendingCreated || isPendingMoved ? '#1D4ED8' : '#1A0E03'
  const isInteractionDisabled = isOnboardingDisabled

  return (
    <div
      title={label}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      style={{
        ...treeRowBaseStyle,
        paddingLeft: 6,
        ...(isInactive ? treeRowInactiveStyle : undefined),
        ...(isPendingCreated ? treeRowPendingCreateStyle : undefined),
        ...(isPendingDeleted ? treeRowPendingDeleteStyle : undefined),
        ...(isPendingMoved ? treeRowPendingMoveStyle : undefined),
        ...(isHovered && !isSelected ? treeRowHoverStyle : undefined),
        ...(isSelected ? treeRowSelectedStyle : undefined),
        ...(isOnboardingHighlighted ? onboardingRegionRowStyle : undefined),
        ...(isOnboardingDisabled ? onboardingDisabledRegionRowStyle : undefined),
      }}
    >
      <button
        type="button"
        title={label}
        onClick={(e) => { e.stopPropagation(); if (!isInteractionDisabled) onToggle(regionId) }}
        style={{
          ...treeToggleBtnStyle,
          ...(isInteractionDisabled ? onboardingDisabledControlStyle : undefined),
        }}
        aria-label={isOpen ? copy.collapseRegion : copy.expandRegion}
        disabled={isInteractionDisabled}
      >
        {isOpen ? '📂' : '📁'}
      </button>
      {isPendingMoved && !isPendingCreated && !isPendingDeleted && <span style={pendingMoveIconStyle}>🔀</span>}
      {isSelected && !readOnly ? (
        <input
          value={titleInput}
          title={label}
          onChange={(event) => onTitleChange(event.target.value)}
          onClick={(event) => event.stopPropagation()}
          onKeyDown={(event) => event.stopPropagation()}
          placeholder={copy.regionPlaceholder(displayNumber)}
          disabled={isEditingDisabled}
          style={{
            ...treeInlineInputStyle,
            color: labelColor,
            textDecoration: isPendingDeleted ? 'line-through' : 'none',
          }}
        />
      ) : (
        <span
          role="button"
          tabIndex={0}
          title={label}
          onClick={() => { if (!isInteractionDisabled) onSelect(regionId) }}
          onKeyDown={(e) => { if (!isInteractionDisabled && (e.key === 'Enter' || e.key === ' ')) onSelect(regionId) }}
          aria-disabled={isInteractionDisabled}
          style={{
            ...treeRowLabelStyle,
            fontSize: compact ? '13px' : '12px',
            fontWeight: 700,
            lineHeight: compact ? 1.35 : 1.45,
            color: labelColor,
            textDecoration: isPendingDeleted ? 'line-through' : 'none',
            ...(isInteractionDisabled ? onboardingDisabledControlStyle : undefined),
          }}
        >
          {label}
        </span>
      )}
      {isPendingDeleted && <span style={pendingDeleteBadgeStyle}>{copy.pendingDelete}</span>}
      {isPendingCreated && !isPendingDeleted && <span style={pendingCreateBadgeStyle}>{copy.pendingCreate}</span>}
    </div>
  )
}

const onboardingRegionRowStyle = {
  position: 'relative',
  zIndex: 12,
  borderRadius: 8,
  background: 'linear-gradient(135deg, rgba(255, 250, 228, 0.98), rgba(196, 236, 215, 0.96))',
  boxShadow: '0 0 0 2px rgba(16, 185, 129, 0.86), 0 0 22px rgba(16, 185, 129, 0.42)',
} as const

const onboardingDisabledRegionRowStyle = {
  opacity: 0.42,
  filter: 'grayscale(0.25)',
} as const

const onboardingDisabledControlStyle = {
  cursor: 'not-allowed',
  pointerEvents: 'none',
} as const
