'use client'

import { useState } from 'react'
import type { ExplorerNode } from '../explorerPlanTypes'
import type { PendingItemState } from '../useExplorerPlan'
import { getDashboardCourseDraftCopy, type DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'
import {
  treeRowBaseStyle,
  treeRowSelectedStyle,
  treeRowInactiveStyle,
  treeRowLearningStyle,
  treeRowHoverStyle,
  treeRowLabelStyle,
  treeRowPendingCreateStyle,
  treeRowPendingDeleteStyle,
  treeRowPendingMoveStyle,
  pendingCreateBadgeStyle,
  pendingDeleteBadgeStyle,
  pendingMoveIconStyle,
} from './explorerTreeStyles'

interface ExplorerTreeNodeRowProps {
  node: ExplorerNode
  depth: number  // 1: 지역 직속, 2: 서브지역 직속
  isSelected: boolean
  readOnly?: boolean
  journalLimitedEdit?: boolean
  compact?: boolean
  pendingState: PendingItemState
  titleInput: string
  onTitleChange: (value: string) => void
  isEditingDisabled: boolean
  onSelect: (nodeId: string) => void
  copy?: DashboardCourseDraftCopy['tree']
  isOnboardingHighlighted?: boolean
  isOnboardingDisabled?: boolean
}

export function ExplorerTreeNodeRow({
  node,
  depth,
  isSelected,
  readOnly = false,
  journalLimitedEdit = false,
  compact = false,
  pendingState,
  onSelect,
  copy = getDashboardCourseDraftCopy('ko').tree,
  isOnboardingHighlighted = false,
  isOnboardingDisabled = false,
}: ExplorerTreeNodeRowProps) {
  const [isHovered, setIsHovered] = useState(false)
  const icon = node.node_type === 'exploration' ? '🎯' : '🔬'
  const paddingLeft = depth === 1 ? 28 : 44
  const isInactive = node.status === 'inactive'
  const hasLink = node.node_type === 'exploration' && !!node.source_url
  const isLearning = (node as ExplorerNode & { learning_status?: string }).learning_status === 'learning'
  const isPendingDeleted = pendingState === 'deleted'
  const isPendingCreated = pendingState === 'created'
  const isPendingMoved = pendingState === 'moved'
  const labelColor = isPendingDeleted
    ? '#A43124'
    : isPendingCreated || isPendingMoved || isLearning
      ? '#1D4ED8'
      : '#1A0E03'

  return (
    <div
      role="button"
      tabIndex={isOnboardingDisabled ? -1 : 0}
      title={node.title}
      aria-disabled={isOnboardingDisabled}
      onClick={() => {
        if (isOnboardingDisabled) return
        onSelect(node.id)
      }}
      onKeyDown={(e) => {
        if (isOnboardingDisabled) return
        if (e.key === 'Enter' || e.key === ' ') onSelect(node.id)
      }}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      style={{
        ...treeRowBaseStyle,
        paddingLeft,
        cursor: isOnboardingDisabled ? 'default' : 'pointer',
        ...(isInactive ? treeRowInactiveStyle : undefined),
        ...(isPendingCreated ? treeRowPendingCreateStyle : undefined),
        ...(isPendingDeleted ? treeRowPendingDeleteStyle : undefined),
        ...(isPendingMoved ? treeRowPendingMoveStyle : undefined),
        ...(isLearning && !isPendingDeleted ? treeRowLearningStyle : undefined),
        ...(isHovered && !isSelected ? treeRowHoverStyle : undefined),
        ...(isSelected ? treeRowSelectedStyle : undefined),
        ...(isOnboardingHighlighted ? onboardingNodeHighlightStyle : undefined),
        ...(isOnboardingDisabled ? onboardingNodeDisabledStyle : undefined),
      }}
    >
      <span style={{ flexShrink: 0, fontSize: compact ? '14px' : '13px', lineHeight: 1.45 }}>
        {icon}
      </span>
      {isPendingMoved && !isPendingCreated && !isPendingDeleted && <span style={pendingMoveIconStyle}>🔀</span>}
      <span
        title={node.title}
        style={{
          ...treeRowLabelStyle,
          fontSize: compact ? '13px' : '12px',
          fontWeight: compact ? 700 : 600,
          lineHeight: compact ? 1.35 : 1.45,
          color: labelColor,
          flex: 1,
          minWidth: 0,
          overflow: 'hidden',
          textOverflow: 'ellipsis',
          whiteSpace: 'nowrap',
          textDecoration: isPendingDeleted ? 'line-through' : 'none',
        }}
      >
        {node.title}
      </span>
      {isPendingDeleted && <span style={pendingDeleteBadgeStyle}>{copy.pendingDelete}</span>}
      {isPendingCreated && !isPendingDeleted && <span style={pendingCreateBadgeStyle}>{copy.pendingCreate}</span>}
      {hasLink && !isPendingDeleted && !journalLimitedEdit && (
        <button
          type="button"
          title={copy.openLinkTitle}
          style={{
            flexShrink: 0,
            background: 'none',
            border: 'none',
            padding: '0 2px',
            cursor: 'pointer',
            fontSize: '11px',
            lineHeight: 1,
            opacity: 0.55,
            color: '#4A3520',
          }}
          onClick={(e) => {
            e.stopPropagation()
            window.open(node.source_url!, '_blank', 'noopener,noreferrer')
          }}
        >
          🔗
        </button>
      )}
    </div>
  )
}

const onboardingNodeHighlightStyle = {
  position: 'relative',
  zIndex: 12,
  boxShadow: '0 0 0 2px rgba(16, 185, 129, 0.82), 0 0 18px rgba(16, 185, 129, 0.48)',
  background: 'linear-gradient(90deg, rgba(236, 253, 245, 0.98), rgba(219, 234, 254, 0.86))',
} as const

const onboardingNodeDisabledStyle = {
  opacity: 0.42,
  filter: 'grayscale(0.22)',
} as const
