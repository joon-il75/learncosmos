'use client'

import { useMemo } from 'react'
import type { ExplorerPlanState } from './useExplorerPlan'
import type { ExplorerEditFormState } from './tree/useExplorerEditForm'
import { ExplorerTreeCourseRow } from './tree/ExplorerTreeCourseRow'
import { ExplorerTreeRegionRow } from './tree/ExplorerTreeRegionRow'
import { ExplorerTreeSubRegionRow } from './tree/ExplorerTreeSubRegionRow'
import { ExplorerTreeNodeRow } from './tree/ExplorerTreeNodeRow'
import { ExplorerTreeEditPanel } from './tree/ExplorerTreeEditPanel'
import { DeleteConfirmModal } from './tree/ExplorerEditModals'
import { panelStyle, compactPanelStyle, scrollAreaStyle } from './tree/explorerTreeStyles'
import { useExplorerTreeCollapse } from './tree/useExplorerTreeCollapse'
import type { ExplorerNode, RegionAggregate, SubRegionAggregate } from './explorerPlanTypes'
import type { PendingItemState } from './useExplorerPlan'
import type { PlanningLumiGuideFocus } from '@/app/dashboard/course-drafts/[id]/sections/planning/PlanningLumiGuide'
import type { DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'

interface ExplorerTreePanelProps {
  explorerPlan: ExplorerPlanState
  editForm: ExplorerEditFormState
  readOnly?: boolean
  journalLimitedEdit?: boolean
  compact?: boolean
  height?: string | number
  onNodeActivate?: (nodeId: string) => void
  onSelectionFocusChange?: (focus: PlanningLumiGuideFocus) => void
  hideEditPanel?: boolean
  suppressSelectionFocus?: boolean
  copy: DashboardCourseDraftCopy['tree']
}

function byOrderIndex<T extends { order_index: number; created_at?: string }>(a: T, b: T) {
  if (a.order_index !== b.order_index) return a.order_index - b.order_index
  return (a.created_at ?? '').localeCompare(b.created_at ?? '')
}

function sortRegions(regions: RegionAggregate[]) {
  return [...regions].sort((a, b) => byOrderIndex(a.region, b.region))
}

function sortSubRegions(subregions: SubRegionAggregate[]) {
  return [...subregions].sort((a, b) => byOrderIndex(a.subregion, b.subregion))
}

function sortNodes(nodes: ExplorerNode[]) {
  return [...nodes].sort(byOrderIndex)
}

export default function ExplorerTreePanel({
  explorerPlan,
  editForm,
  readOnly = false,
  journalLimitedEdit = false,
  compact = false,
  height = '100%',
  onNodeActivate,
  onSelectionFocusChange,
  hideEditPanel = false,
  suppressSelectionFocus = false,
  copy,
}: ExplorerTreePanelProps) {
  const {
    courseAggregate,
    isLoading,
    selectedRegionId,
    selectedSubRegionId,
    selectedNodeId,
    handleSelectCourse,
    handleSelectRegion,
    handleSelectSubRegion,
    handleSelectNode,
    getRegionPendingState,
    getSubRegionPendingState,
    getNodePendingState,
  } = explorerPlan
  const rowReadOnly = readOnly || journalLimitedEdit

  const { collapsedRegionIds, collapsedSubRegionIds, toggleRegion, toggleSubRegion } =
    useExplorerTreeCollapse()

  const isCourseSelected =
    selectedRegionId === null && selectedSubRegionId === null && selectedNodeId === null

  const regions = useMemo(() => sortRegions(courseAggregate?.regions ?? []), [courseAggregate])

  if (isLoading) {
    return (
      <div style={{ ...panelStyle, ...(compact ? compactPanelStyle : undefined), height }}>
        <div style={{ padding: 12, fontSize: 12, color: 'rgba(56,40,18,0.6)' }}>
          {copy.loading}
        </div>
      </div>
    )
  }

  return (
    <div style={{ ...panelStyle, ...(compact ? compactPanelStyle : undefined), gap: 0, height }}>
      {/* 코스명 헤더 — 기본 선택 상태 */}
      <ExplorerTreeCourseRow
        displayTitle={copy.formatPlanTitle(courseAggregate?.title)}
        displayPlanetName={copy.formatPlanetName(courseAggregate?.title)}
        isSelected={isCourseSelected}
        readOnly={rowReadOnly}
        copy={copy}
        onSelect={() => {
          handleSelectCourse()
          if (!suppressSelectionFocus) onSelectionFocusChange?.('course')
        }}
      />

      {/* 트리 리스트 */}
      <div style={{ ...scrollAreaStyle, padding: '4px 10px 8px 0' }}>
        {regions.length === 0 && (
          <div style={{ padding: '10px 14px', fontSize: 12, color: 'rgba(56,40,18,0.55)' }}>
            {copy.emptyRegions}
          </div>
        )}

        {regions.map((regionAgg, regionIndex) => {
          const regionId = regionAgg.region.id
          const isRegionOpen = !collapsedRegionIds.has(regionId)
          const isRegionSelected =
            selectedRegionId === regionId &&
            selectedSubRegionId === null &&
            selectedNodeId === null
          const regionPendingState = getRegionPendingState(regionId)
          const regionIsPendingDeleted = regionPendingState === 'deleted'

          return (
            <div key={regionId}>
              <ExplorerTreeRegionRow
                regionAgg={regionAgg}
                regionIndex={regionIndex}
                isOpen={isRegionOpen}
                isSelected={isRegionSelected}
                readOnly
                compact={compact}
                pendingState={regionPendingState}
                copy={copy}
                titleInput={isRegionSelected ? editForm.regionNameInput : regionAgg.region.name}
                onTitleChange={editForm.setRegionNameInput}
                isEditingDisabled={editForm.isMutating || regionPendingState === 'deleted'}
                onToggle={toggleRegion}
                onSelect={(regionId) => {
                  handleSelectRegion(regionId)
                  if (!suppressSelectionFocus) onSelectionFocusChange?.('region')
                }}
              />

              {isRegionOpen && (
                <>
                  {/* 서브지역들 */}
                  {sortSubRegions(regionAgg.subregions).map((subAgg, subIndex) => {
                    const subId = subAgg.subregion.id
                    const isSubOpen = !collapsedSubRegionIds.has(subId)
                    const isSubSelected =
                      selectedSubRegionId === subId && selectedNodeId === null
                    const directSubPendingState = getSubRegionPendingState(subId)
                    const subPendingState: PendingItemState = regionIsPendingDeleted
                      ? 'deleted'
                      : directSubPendingState

                    return (
                      <div key={subId}>
                        <ExplorerTreeSubRegionRow
                          subAgg={subAgg}
                          subIndex={subIndex}
                          regionId={regionId}
                          isOpen={isSubOpen}
                          isSelected={isSubSelected}
                          readOnly
                          compact={compact}
                          pendingState={subPendingState}
                          copy={copy}
                          titleInput={isSubSelected ? editForm.subRegionNameInput : subAgg.subregion.name}
                          onTitleChange={editForm.setSubRegionNameInput}
                          isEditingDisabled={editForm.isMutating || subPendingState === 'deleted'}
                          onToggle={toggleSubRegion}
                          onSelect={(regionId, subId) => {
                            handleSelectSubRegion(regionId, subId)
                            if (!suppressSelectionFocus) onSelectionFocusChange?.('subregion')
                          }}
                        />

                        {isSubOpen &&
                          sortNodes(subAgg.nodes).map((node) => {
                            const nodePendingState: PendingItemState = subPendingState === 'deleted'
                              ? 'deleted'
                              : getNodePendingState(node.id)
                            return (
                              <ExplorerTreeNodeRow
                                key={node.id}
                                node={node}
                                depth={2}
                                isSelected={selectedNodeId === node.id}
                                readOnly={rowReadOnly}
                                journalLimitedEdit={journalLimitedEdit}
                                compact={compact}
                                pendingState={nodePendingState}
                                copy={copy}
                                titleInput={selectedNodeId === node.id ? editForm.nodeTitleInput : node.title}
                                onTitleChange={editForm.setNodeTitleInput}
                                isEditingDisabled={editForm.isMutating || nodePendingState === 'deleted'}
                                onSelect={(nodeId) => {
                                  handleSelectNode(nodeId)
                                  if (!suppressSelectionFocus) onSelectionFocusChange?.(node.node_type === 'research' ? 'research-node' : 'exploration-node')
                                  onNodeActivate?.(nodeId)
                                }}
                              />
                            )
                          })}
                      </div>
                    )
                  })}

                  {/* 지역 직속 지점들 */}
                  {sortNodes(regionAgg.nodes).map((node) => {
                    const nodePendingState: PendingItemState = regionIsPendingDeleted
                      ? 'deleted'
                      : getNodePendingState(node.id)
                    return (
                      <ExplorerTreeNodeRow
                        key={node.id}
                        node={node}
                        depth={1}
                        isSelected={selectedNodeId === node.id}
                        readOnly={rowReadOnly}
                        journalLimitedEdit={journalLimitedEdit}
                        compact={compact}
                        pendingState={nodePendingState}
                        copy={copy}
                        titleInput={selectedNodeId === node.id ? editForm.nodeTitleInput : node.title}
                        onTitleChange={editForm.setNodeTitleInput}
                        isEditingDisabled={editForm.isMutating || nodePendingState === 'deleted'}
                        onSelect={(nodeId) => {
                          handleSelectNode(nodeId)
                          if (!suppressSelectionFocus) onSelectionFocusChange?.(node.node_type === 'research' ? 'research-node' : 'exploration-node')
                          onNodeActivate?.(nodeId)
                        }}
                      />
                    )
                  })}
                </>
              )}
            </div>
          )
        })}
      </div>

      {/* 하단 편집 폼 (버튼 바는 device shell 밖에 별도 렌더) */}
      {!readOnly && !hideEditPanel && <ExplorerTreeEditPanel editForm={editForm} journalLimitedEdit={journalLimitedEdit} copy={copy.editPanel} />}
      {!readOnly && hideEditPanel && editForm.deleteConfirm && (
        <DeleteConfirmModal
          mode={editForm.deleteConfirm.mode}
          action={editForm.deleteConfirm.action}
          title={editForm.deleteConfirm.title}
          message={editForm.deleteConfirm.message}
          requireTitleInput={editForm.deleteConfirm.requireTitleInput}
          confirmInput={editForm.deleteConfirm.confirmInput}
          onConfirmInputChange={editForm.setDeleteConfirmInput}
          onConfirm={() => void editForm.handleConfirmDelete()}
          onClose={editForm.handleCloseDeleteConfirm}
          isMutating={editForm.isMutating}
        />
      )}
    </div>
  )
}
