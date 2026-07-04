'use client'

import { useState } from 'react'
import type { SubRegionAggregate } from '../explorerPlanTypes'
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

interface ExplorerTreeSubRegionRowProps {
  subAgg: SubRegionAggregate
  subIndex: number
  regionId: string
  isOpen: boolean
  isSelected: boolean
  readOnly?: boolean
  compact?: boolean
  pendingState: PendingItemState
  titleInput: string
  onTitleChange: (value: string) => void
  isEditingDisabled: boolean
  onToggle: (subRegionId: string) => void
  onSelect: (regionId: string, subRegionId: string) => void
  copy?: DashboardCourseDraftCopy['tree']
}

export function ExplorerTreeSubRegionRow({
  subAgg,
  subIndex,
  regionId,
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
}: ExplorerTreeSubRegionRowProps) {
  const [isHovered, setIsHovered] = useState(false)
  const subId = subAgg.subregion.id
  const displayNumber = subIndex + 1
  const name = subAgg.subregion.name.trim()
  const label = copy.subRegionLabel(displayNumber, name)
  const isInactive = subAgg.subregion.status === 'inactive'
  const isPendingDeleted = pendingState === 'deleted'
  const isPendingCreated = pendingState === 'created'
  const isPendingMoved = pendingState === 'moved'
  const labelColor = isPendingDeleted ? '#A43124' : isPendingCreated || isPendingMoved ? '#1D4ED8' : '#1A0E03'

  return (
    <div
      title={label}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      style={{
        ...treeRowBaseStyle,
        paddingLeft: 22,
        ...(isInactive ? treeRowInactiveStyle : undefined),
        ...(isPendingCreated ? treeRowPendingCreateStyle : undefined),
        ...(isPendingDeleted ? treeRowPendingDeleteStyle : undefined),
        ...(isPendingMoved ? treeRowPendingMoveStyle : undefined),
        ...(isHovered && !isSelected ? treeRowHoverStyle : undefined),
        ...(isSelected ? treeRowSelectedStyle : undefined),
      }}
    >
      <button
        type="button"
        title={label}
        onClick={(e) => { e.stopPropagation(); onToggle(subId) }}
        style={treeToggleBtnStyle}
        aria-label={isOpen ? copy.collapseSubRegion : copy.expandSubRegion}
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
          placeholder={copy.subRegionPlaceholder(displayNumber)}
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
          onClick={() => onSelect(regionId, subId)}
          onKeyDown={(e) => { if (e.key === 'Enter' || e.key === ' ') onSelect(regionId, subId) }}
          style={{
            ...treeRowLabelStyle,
            fontSize: compact ? '13px' : '12px',
            fontWeight: compact ? 700 : 600,
            lineHeight: compact ? 1.35 : 1.45,
            color: labelColor,
            textDecoration: isPendingDeleted ? 'line-through' : 'none',
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
