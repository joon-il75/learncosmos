'use client'

import { useState, useEffect, useCallback } from 'react'
import type { ExplorerPlanState } from '../useExplorerPlan'
import type { ParentKind, ExplorerNode, RegionAggregate, SubRegionAggregate } from '../explorerPlanTypes'
import { getDashboardCourseDraftCopy, type DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'
import type {
  ExplorerContentCandidate,
  DeleteConfirmState,
} from './explorerEditFormTypes'
import {
  selectedNodeSource,
} from './explorerEditFormUtils'
import { useExplorerRecommendation } from './useExplorerRecommendation'
import { useExplorerDeleteConfirm } from './useExplorerDeleteConfirm'
export type { ExplorerContentCandidate } from './explorerEditFormTypes'

const defaultModalCopy = getDashboardCourseDraftCopy('ko').tree.editPanel.modals

export type ExplorerEditFormState = {
  // Derived selection
  isCourseSelected: boolean
  isRegionSelected: boolean
  isSubRegionSelected: boolean
  isNodeSelected: boolean
  selectedRegionAgg: RegionAggregate | null
  selectedSubAgg: SubRegionAggregate | null
  selectedNode: ExplorerNode | null
  selectedNodeRegionAgg: RegionAggregate | null
  selectedNodeSubAgg: SubRegionAggregate | null
  canDelete: boolean
  canDeleteCourse: boolean
  canActivateCourse: boolean
  hasChildren: boolean
  isSelectedInactive: boolean
  isSelectedPendingCreated: boolean
  isSelectedPendingDeleted: boolean
  isActionBarBlocked: boolean
  canAddSubRegion: boolean
  canMoveUp: boolean
  canMoveDown: boolean
  showInactiveItems: boolean

  // Course form
  courseTitleInput: string
  setCourseTitleInput: (v: string) => void
  newRegionNameInput: string
  setNewRegionNameInput: (v: string) => void

  // Region form
  regionNameInput: string
  setRegionNameInput: (v: string) => void
  showAddSubRegion: boolean
  setShowAddSubRegion: (v: boolean | ((prev: boolean) => boolean)) => void
  newSubRegionNameInput: string
  setNewSubRegionNameInput: (v: string) => void

  // SubRegion form
  subRegionNameInput: string
  setSubRegionNameInput: (v: string) => void

  // Add-node form
  addNodeMode: 'exploration' | 'research' | null
  setAddNodeMode: (v: 'exploration' | 'research' | null) => void
  newNodeTitleInput: string
  setNewNodeTitleInput: (v: string) => void
  newNodeUrlInput: string
  setNewNodeUrlInput: (v: string) => void

  // Node form
  nodeTitleInput: string
  setNodeTitleInput: (v: string) => void
  nodeUrlInput: string
  setNodeUrlInput: (v: string) => void
  urlCheckStatus: 'idle' | 'checking' | 'ok' | 'error'
  setUrlCheckStatus: (v: 'idle' | 'checking' | 'ok' | 'error') => void
  urlCheckMsg: string
  setUrlCheckMsg: (v: string) => void
  newUrlCheckStatus: 'idle' | 'checking' | 'ok' | 'error'
  newUrlCheckMsg: string
  isRecommendationOpen: boolean
  recommendationQuery: string
  setRecommendationQuery: (v: string) => void
  recommendationCandidates: ExplorerContentCandidate[]
  isLoadingRecommendations: boolean
  recommendationMessage: string | null
  recommendationPointError: string | null
  deleteConfirm: DeleteConfirmState
  setDeleteConfirmInput: (value: string) => void

  // Add item modal
  addModalType: 'child-object' | 'exploration-choice' | 'exploration-url' | 'research' | null
  setAddModalType: (v: 'child-object' | 'exploration-choice' | 'exploration-url' | 'research' | null) => void
  editNodeModalType: 'exploration' | 'research' | null

  // Handlers
  handleSave: () => Promise<void>
  handleCancel: () => void
  handleDelete: () => Promise<void>
  handleActivateCourse: () => Promise<void>
  handleConfirmDelete: () => Promise<void>
  handleCloseDeleteConfirm: () => void
  handleInactiveToggle: () => Promise<void> | void
  handleMoveUp: () => void
  handleMoveDown: () => void
  handleCheckUrl: (target?: 'new' | 'selected') => Promise<void>
  handleOpenRecommendationModal: (target?: 'new' | 'selected') => Promise<void>
  handleCloseRecommendationModal: () => void
  handleSearchRecommendations: () => Promise<void>
  handleApplyRecommendation: (candidate: ExplorerContentCandidate) => Promise<void>
  handleOpenNodeEditModal: () => void
  handleCloseNodeEditModal: () => void
  handleSubmitNodeEdit: () => Promise<void>
  handleRenameCourse: () => Promise<void>
  handleRenameRegion: () => Promise<void>
  handleAddRegion: () => Promise<void>
  handleAddSubRegion: () => Promise<void>
  handleAddNode: (modeOverride?: 'exploration' | 'research') => Promise<void>

  // Forwarded from explorerPlan
  isMutating: boolean
  mutationMessage: string | null
}

export function useExplorerEditForm(
  explorerPlan: ExplorerPlanState,
  copy: DashboardCourseDraftCopy['tree']['editPanel']['modals'] = defaultModalCopy,
): ExplorerEditFormState {
  const {
    courseAggregate,
    selectedRegionId,
    selectedSubRegionId,
    selectedNodeId,
    updateCourseTitle,
    deleteCourseDraft,
    activateCourseDraft,
    getRegionPendingState,
    getSubRegionPendingState,
    getNodePendingState,
    cancelPlanChanges,
    restoreDeletedRegion,
    restoreDeletedSubRegion,
    restoreDeletedNode,
    createRegion,
    updateRegion,
    deleteRegion,
    moveRegion,
    createSubRegion,
    updateSubRegion,
    deleteSubRegion,
    moveSubRegion,
    createResearchNode,
    updateResearchNode,
    deleteResearchNode,
    createExplorationNode,
    updateExplorationNode,
    deleteExplorationNode,
    moveNode,
    savePlanChanges,
    showInactiveItems,
    toggleInactiveItems,
    activateRegion,
    activateSubRegion,
    activateNode,
    isMutating,
    mutationMessage,
  } = explorerPlan

  const isCourseSelected =
    selectedRegionId === null && selectedSubRegionId === null && selectedNodeId === null

  const selectedRegionAgg: RegionAggregate | null = selectedRegionId
    ? courseAggregate?.regions.find((r) => r.region.id === selectedRegionId) ?? null
    : null

  const selectedSubAgg: SubRegionAggregate | null =
    selectedSubRegionId && selectedRegionAgg
      ? selectedRegionAgg.subregions.find((s) => s.subregion.id === selectedSubRegionId) ?? null
      : null

  const selectedNode: ExplorerNode | null = selectedNodeId
    ? courseAggregate?.regions
        .flatMap((r) => [...r.nodes, ...r.subregions.flatMap((s) => s.nodes)])
        .find((n) => n.id === selectedNodeId) ?? null
    : null

  const isRegionSelected =
    selectedRegionAgg !== null && selectedSubRegionId === null && selectedNodeId === null
  const isSubRegionSelected = selectedSubAgg !== null && selectedNodeId === null
  const isNodeSelected = selectedNode !== null

  const selectedRegionPendingState = selectedRegionId ? getRegionPendingState(selectedRegionId) : 'none'
  const selectedSubRegionDirectPendingState = selectedSubRegionId ? getSubRegionPendingState(selectedSubRegionId) : 'none'
  const selectedSubRegionPendingState =
    selectedRegionPendingState === 'deleted' ? 'deleted' : selectedSubRegionDirectPendingState
  const selectedNodeDirectPendingState = selectedNodeId ? getNodePendingState(selectedNodeId) : 'none'
  const selectedNodeParentPendingState = (() => {
    if (!selectedNode) return 'none'
    if (selectedNode.parent_kind === 'region') return getRegionPendingState(selectedNode.parent_id)
    const parentRegionAgg = courseAggregate?.regions.find((regionAgg) =>
      regionAgg.subregions.some((subAgg) => subAgg.subregion.id === selectedNode.parent_id)
    )
    if (parentRegionAgg && getRegionPendingState(parentRegionAgg.region.id) === 'deleted') return 'deleted'
    return getSubRegionPendingState(selectedNode.parent_id)
  })()
  const selectedNodePendingState =
    selectedNodeParentPendingState === 'deleted' ? 'deleted' : selectedNodeDirectPendingState
  const isSelectedPendingDeleted = isRegionSelected
    ? selectedRegionPendingState === 'deleted'
    : isSubRegionSelected
    ? selectedSubRegionPendingState === 'deleted'
    : isNodeSelected
    ? selectedNodePendingState === 'deleted'
    : false
  const isSelectedPendingCreated = isRegionSelected
    ? selectedRegionPendingState === 'created'
    : isSubRegionSelected
    ? selectedSubRegionPendingState === 'created'
    : isNodeSelected
    ? selectedNodePendingState === 'created'
    : false

  const canDeleteCourse =
    isCourseSelected &&
    courseAggregate?.status !== 'archived'
  const canActivateCourse =
    Boolean(courseAggregate?.is_inactive)
  const canDelete = canDeleteCourse || isRegionSelected || isSubRegionSelected || isNodeSelected
  const isSelectedInactive = isRegionSelected
    ? selectedRegionAgg?.region.status === 'inactive'
    : isSubRegionSelected
    ? selectedSubAgg?.subregion.status === 'inactive'
    : isNodeSelected
    ? selectedNode?.status === 'inactive'
    : false
  const hasChildren = isRegionSelected
    ? (selectedRegionAgg?.subregions.length ?? 0) > 0 ||
      (selectedRegionAgg?.nodes.length ?? 0) > 0
    : isSubRegionSelected
    ? (selectedSubAgg?.nodes.length ?? 0) > 0
    : false
  const canAddSubRegion =
    isRegionSelected &&
    !isSelectedPendingDeleted &&
    (selectedRegionAgg?.subregions.filter((subAgg) => getSubRegionPendingState(subAgg.subregion.id) !== 'deleted').length ?? 0) < 3
  const remainingRegionCount =
    courseAggregate?.regions.filter((regionAgg) => getRegionPendingState(regionAgg.region.id) !== 'deleted').length ?? 0

  const selectedRegionIndex = selectedRegionId && courseAggregate
    ? courseAggregate.regions.findIndex((regionAgg) => regionAgg.region.id === selectedRegionId)
    : -1
  const selectedSubRegionIndex = selectedSubRegionId && selectedRegionAgg
    ? selectedRegionAgg.subregions.findIndex((subAgg) => subAgg.subregion.id === selectedSubRegionId)
    : -1
  const selectedNodeContext = (() => {
    if (!selectedNode || !courseAggregate) {
      return {
        regionAgg: null as RegionAggregate | null,
        selectedSubAgg: null as SubRegionAggregate | null,
        previousSubAgg: null as SubRegionAggregate | null,
        nextSubAgg: null as SubRegionAggregate | null,
        previousRegionAgg: null as RegionAggregate | null,
        nextRegionAgg: null as RegionAggregate | null,
        siblings: [] as ExplorerNode[],
      }
    }

    const sortedRegions = [...courseAggregate.regions].sort(
      (a, b) => a.region.order_index - b.region.order_index
    )

    for (let regionIndex = 0; regionIndex < sortedRegions.length; regionIndex += 1) {
      const regionAgg = sortedRegions[regionIndex]
      if (selectedNode.parent_kind === 'region' && regionAgg.region.id === selectedNode.parent_id) {
        const sortedSubRegions = [...regionAgg.subregions].sort(
          (a, b) => a.subregion.order_index - b.subregion.order_index
        )
        const sortedRegionNodes = [...regionAgg.nodes].sort(
          (a, b) => a.order_index - b.order_index
        )
        return {
          regionAgg,
          selectedSubAgg: null,
          previousSubAgg: sortedSubRegions.at(-1) ?? null,
          nextSubAgg: null,
          previousRegionAgg: regionIndex > 0 ? sortedRegions[regionIndex - 1] : null,
          nextRegionAgg: regionIndex < sortedRegions.length - 1 ? sortedRegions[regionIndex + 1] : null,
          siblings: sortedRegionNodes,
        }
      }

      const sortedSubRegions = [...regionAgg.subregions].sort(
        (a, b) => a.subregion.order_index - b.subregion.order_index
      )
      const subIndex = sortedSubRegions.findIndex((item) => item.subregion.id === selectedNode.parent_id)
      if (subIndex >= 0) {
        const sortedSubNodes = [...sortedSubRegions[subIndex].nodes].sort(
          (a, b) => a.order_index - b.order_index
        )
        return {
          regionAgg,
          selectedSubAgg: sortedSubRegions[subIndex],
          previousSubAgg: subIndex > 0 ? sortedSubRegions[subIndex - 1] : null,
          nextSubAgg: subIndex < sortedSubRegions.length - 1 ? sortedSubRegions[subIndex + 1] : null,
          previousRegionAgg: regionIndex > 0 ? sortedRegions[regionIndex - 1] : null,
          nextRegionAgg: regionIndex < sortedRegions.length - 1 ? sortedRegions[regionIndex + 1] : null,
          siblings: sortedSubNodes,
        }
      }
    }

    return {
      regionAgg: null as RegionAggregate | null,
      selectedSubAgg: null as SubRegionAggregate | null,
      previousSubAgg: null as SubRegionAggregate | null,
      nextSubAgg: null as SubRegionAggregate | null,
      previousRegionAgg: null as RegionAggregate | null,
      nextRegionAgg: null as RegionAggregate | null,
      siblings: [] as ExplorerNode[],
    }
  })()
  const selectedNodeSiblings = selectedNodeContext.siblings
  const selectedNodeRegionAgg = selectedNodeContext.regionAgg
  const selectedNodeSubAgg = selectedNodeContext.selectedSubAgg
  const selectedNodeIndex = selectedNode
    ? selectedNodeSiblings.findIndex((node) => node.id === selectedNode.id)
    : -1
  const canMoveNodeUpAcrossBoundary = selectedNode !== null && (
    (selectedNode.parent_kind === 'region' &&
      selectedNodeIndex === 0 &&
      (selectedNodeContext.previousSubAgg !== null || selectedNodeContext.previousRegionAgg !== null)) ||
    (selectedNode.parent_kind === 'subregion' &&
      selectedNodeIndex === 0 &&
      (selectedNodeContext.previousSubAgg !== null || selectedNodeContext.previousRegionAgg !== null))
  )
  const canMoveNodeDownAcrossBoundary = selectedNode !== null && selectedNodeIndex === selectedNodeSiblings.length - 1 && (
    (selectedNode.parent_kind === 'region' && selectedNodeContext.nextRegionAgg !== null) ||
    (selectedNode.parent_kind === 'subregion' && (
      selectedNodeContext.nextSubAgg !== null ||
      (selectedNodeContext.regionAgg?.nodes.length ?? 0) > 0 ||
      selectedNodeContext.nextRegionAgg !== null
    ))
  )
  const canMoveUp = isSelectedPendingDeleted
    ? false
    : isRegionSelected
    ? selectedRegionIndex > 0
    : isSubRegionSelected
    ? selectedSubRegionIndex > 0
    : isNodeSelected
    ? selectedNodeIndex > 0 || canMoveNodeUpAcrossBoundary
    : false
  const canMoveDown = isSelectedPendingDeleted
    ? false
    : isRegionSelected
    ? courseAggregate !== null && selectedRegionIndex >= 0 && selectedRegionIndex < courseAggregate.regions.length - 1
    : isSubRegionSelected
    ? selectedRegionAgg !== null && selectedSubRegionIndex >= 0 && selectedSubRegionIndex < selectedRegionAgg.subregions.length - 1
    : isNodeSelected
    ? selectedNodeIndex >= 0 && (selectedNodeIndex < selectedNodeSiblings.length - 1 || canMoveNodeDownAcrossBoundary)
    : false

  // ── Form state ─────────────────────────────────────────────────────────────
  const [courseTitleInput, setCourseTitleInput] = useState(courseAggregate?.title ?? '')
  const [newRegionNameInput, setNewRegionNameInput] = useState('')
  const [regionNameInput, setRegionNameInput] = useState('')
  const [showAddSubRegion, setShowAddSubRegion] = useState(false)
  const [newSubRegionNameInput, setNewSubRegionNameInput] = useState('')
  const [subRegionNameInput, setSubRegionNameInput] = useState('')
  const [addNodeMode, setAddNodeMode] = useState<'exploration' | 'research' | null>(null)
  const [newNodeTitleInput, setNewNodeTitleInput] = useState('')
  const [newNodeUrlInput, setNewNodeUrlInput] = useState('')
  const [nodeTitleInput, setNodeTitleInput] = useState('')
  const [nodeUrlInput, setNodeUrlInput] = useState('')
  const [urlCheckStatus, setUrlCheckStatus] = useState<'idle' | 'checking' | 'ok' | 'error'>('idle')
  const [urlCheckMsg, setUrlCheckMsg] = useState('')
  const [newUrlCheckStatus, setNewUrlCheckStatus] = useState<'idle' | 'checking' | 'ok' | 'error'>('idle')
  const [newUrlCheckMsg, setNewUrlCheckMsg] = useState('')
  const [addModalType, setAddModalTypeState] = useState<'child-object' | 'exploration-choice' | 'exploration-url' | 'research' | null>(null)
  const [editNodeModalType, setEditNodeModalType] = useState<'exploration' | 'research' | null>(null)
  const {
    isRecommendationOpen,
    recommendationQuery, setRecommendationQuery,
    recommendationCandidates,
    isLoadingRecommendations,
    recommendationMessage,
    recommendationPointError,
    setIsRecommendationOpen,
    setRecommendationTarget,
    setRecommendationCandidates,
    setRecommendationMessage,
    setRecommendationPointError,
    buildRecommendationQuery,
    handleOpenRecommendationModal,
    handleCloseRecommendationModal,
    handleSearchRecommendations,
    handleApplyRecommendation,
  } = useExplorerRecommendation({
    courseAggregate,
    selectedRegionAgg,
    selectedSubAgg,
    selectedNode,
    nodeTitleInput,
    newNodeTitleInput,
    isRegionSelected,
    isSubRegionSelected,
    selectedRegionId,
    selectedSubRegionId,
    createExplorationNode,
    updateExplorationNode,
    setAddNodeMode,
    setNewNodeTitleInput,
    setNewNodeUrlInput,
    setAddModalTypeState,
    setNodeTitleInput,
    setNodeUrlInput,
    setUrlCheckStatus,
    setUrlCheckMsg,
    copy,
  })

  // Reset on selection change
  useEffect(() => {
    setCourseTitleInput(courseAggregate?.title ?? '')
    setNewRegionNameInput('')
  }, [isCourseSelected, courseAggregate?.title])

  useEffect(() => {
    setRegionNameInput(selectedRegionAgg?.region.name ?? '')
    setNewSubRegionNameInput('')
    setShowAddSubRegion(false)
    setAddNodeMode(null)
    setNewNodeTitleInput('')
    setNewNodeUrlInput('')
    setAddModalTypeState(null)
  }, [selectedRegionId]) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    setSubRegionNameInput(selectedSubAgg?.subregion.name ?? '')
    setAddNodeMode(null)
    setNewNodeTitleInput('')
    setNewNodeUrlInput('')
    setAddModalTypeState(null)
  }, [selectedSubRegionId]) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    setNodeTitleInput(selectedNode?.title ?? '')
    setNodeUrlInput(selectedNode?.source_url ?? '')
    setUrlCheckStatus('idle')
    setUrlCheckMsg('')
    setIsRecommendationOpen(false)
    setEditNodeModalType(null)
    setRecommendationCandidates([])
    setRecommendationMessage(null)
  }, [selectedNodeId]) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    setNewUrlCheckStatus('idle')
    setNewUrlCheckMsg('')
  }, [newNodeUrlInput])

  const setAddModalType = useCallback((type: 'child-object' | 'exploration-choice' | 'exploration-url' | 'research' | null) => {
    setAddModalTypeState(type)
    if (type === 'exploration-choice' || type === 'exploration-url') {
      setAddNodeMode('exploration')
      setRecommendationTarget('new')
    } else if (type === 'research') {
      setAddNodeMode('research')
      setRecommendationTarget('new')
    } else if (type === 'child-object') {
      setAddNodeMode(null)
      setRecommendationTarget('new')
    } else {
      setAddNodeMode(null)
    }
    setNewNodeTitleInput('')
    setNewNodeUrlInput('')
    setNewUrlCheckStatus('idle')
    setNewUrlCheckMsg('')
    if (type === 'exploration-choice' || type === 'child-object') {
      setRecommendationTarget('new')
      setRecommendationQuery(buildRecommendationQuery('new'))
      setRecommendationCandidates([])
      setRecommendationMessage(null)
      setRecommendationPointError(null)
      setIsRecommendationOpen(false)
    }
  }, [buildRecommendationQuery])

  // ── Handlers ───────────────────────────────────────────────────────────────

  const handleSave = useCallback(async () => {
    if (isMutating || isModalBlocked) return
    if (isSelectedPendingDeleted) {
      await savePlanChanges()
      return
    }

    if (isCourseSelected) {
      const trimmed = courseTitleInput.trim()
      if (trimmed && trimmed !== courseAggregate?.title) {
        await updateCourseTitle(trimmed)
      } else {
        await savePlanChanges()
      }
      return
    }

    if (isRegionSelected && selectedRegionId) {
      const trimmed = regionNameInput.trim()
      if (trimmed && trimmed !== selectedRegionAgg?.region.name) {
        await updateRegion(selectedRegionId, trimmed)
      }
      await savePlanChanges()
      return
    }

    if (isSubRegionSelected && selectedRegionId && selectedSubRegionId) {
      const trimmed = subRegionNameInput.trim()
      if (trimmed && trimmed !== selectedSubAgg?.subregion.name) {
        await updateSubRegion(selectedRegionId, selectedSubRegionId, trimmed)
      }
      await savePlanChanges()
      return
    }

    if (isNodeSelected && selectedNodeId && selectedNode) {
      if (selectedNode.node_type === 'exploration') {
        const source = selectedNodeSource(selectedNode, nodeUrlInput)
        if (nodeTitleInput.trim() && source) {
          await updateExplorationNode(selectedNodeId, nodeTitleInput.trim(), source)
        }
      } else {
        if (nodeTitleInput.trim()) {
          await updateResearchNode(selectedNodeId, nodeTitleInput.trim())
        }
      }
      await savePlanChanges()
      return
    }

    await savePlanChanges()
  }, [
    isMutating, isCourseSelected, courseTitleInput, courseAggregate?.title, updateCourseTitle,
    isRegionSelected, selectedRegionId, regionNameInput, selectedRegionAgg?.region.name, updateRegion,
    isSubRegionSelected, selectedSubRegionId, subRegionNameInput, selectedSubAgg?.subregion.name, updateSubRegion,
    isNodeSelected, selectedNodeId, selectedNode, nodeTitleInput, nodeUrlInput,
    updateExplorationNode, updateResearchNode, savePlanChanges,
    isSelectedPendingDeleted,
  ])

  const handleCancel = useCallback(() => {
    void cancelPlanChanges()
    if (isCourseSelected) {
      setCourseTitleInput(courseAggregate?.title ?? '')
      setNewRegionNameInput('')
    } else if (isRegionSelected) {
      setRegionNameInput(selectedRegionAgg?.region.name ?? '')
      setNewSubRegionNameInput('')
      setShowAddSubRegion(false)
      setAddNodeMode(null)
    } else if (isSubRegionSelected) {
      setSubRegionNameInput(selectedSubAgg?.subregion.name ?? '')
      setAddNodeMode(null)
    } else if (isNodeSelected) {
      setNodeTitleInput(selectedNode?.title ?? '')
      setNodeUrlInput(selectedNode?.source_url ?? '')
      setUrlCheckStatus('idle')
      setUrlCheckMsg('')
    }
  }, [
    isCourseSelected, isRegionSelected, isSubRegionSelected, isNodeSelected,
    courseAggregate?.title, selectedRegionAgg?.region.name, selectedSubAgg?.subregion.name, selectedNode,
    cancelPlanChanges,
  ])

  const isModalBlocked = addModalType !== null || editNodeModalType !== null || isRecommendationOpen

  const {
    deleteConfirm,
    setDeleteConfirmInput,
    handleDelete,
    handleActivateCourse,
    handleCloseDeleteConfirm,
    handleConfirmDelete,
  } = useExplorerDeleteConfirm({
    isMutating,
    isModalBlocked,
    courseAggregate,
    selectedRegionId,
    selectedSubRegionId,
    selectedNodeId,
    selectedNode,
    selectedRegionAgg,
    selectedSubAgg,
    isRegionSelected,
    isSubRegionSelected,
    isNodeSelected,
    remainingRegionCount,
    hasChildren,
    canDeleteCourse,
    canActivateCourse,
    isSelectedPendingDeleted,
    selectedRegionPendingState,
    getRegionPendingState,
    getSubRegionPendingState,
    restoreDeletedRegion,
    restoreDeletedSubRegion,
    restoreDeletedNode,
    deleteRegion,
    deleteSubRegion,
    deleteExplorationNode,
    deleteResearchNode,
    deleteCourseDraft,
    activateCourseDraft,
    copy,
  })

  const isActionBarBlocked = isModalBlocked || deleteConfirm !== null

  const handleInactiveToggle = useCallback(async () => {
    if (isMutating || isActionBarBlocked || isSelectedPendingDeleted) return
    if (isSelectedInactive) {
      if (isRegionSelected && selectedRegionId) {
        await activateRegion(selectedRegionId)
      } else if (isSubRegionSelected && selectedSubRegionId) {
        await activateSubRegion(selectedSubRegionId)
      } else if (isNodeSelected && selectedNodeId) {
        await activateNode(selectedNodeId)
      }
      return
    }
    toggleInactiveItems()
  }, [
    isMutating, isSelectedInactive,
    isRegionSelected, selectedRegionId, activateRegion,
    isSubRegionSelected, selectedSubRegionId, activateSubRegion,
    isNodeSelected, selectedNodeId, activateNode,
    toggleInactiveItems, isActionBarBlocked, isSelectedPendingDeleted,
  ])

  const handleMoveUp = useCallback(() => {
    if (isMutating || isSelectedPendingDeleted || !canMoveUp) return
    if (isRegionSelected && selectedRegionId) {
      moveRegion(selectedRegionId, 'up')
    } else if (isSubRegionSelected && selectedRegionId && selectedSubRegionId) {
      moveSubRegion(selectedRegionId, selectedSubRegionId, 'up')
    } else if (isNodeSelected && selectedNodeId) {
      moveNode(selectedNodeId, 'up')
    }
  }, [
    canMoveUp, isMutating, isRegionSelected, selectedRegionId, moveRegion,
    isSubRegionSelected, selectedSubRegionId, moveSubRegion,
    isNodeSelected, selectedNodeId, moveNode, isSelectedPendingDeleted,
  ])

  const handleMoveDown = useCallback(() => {
    if (isMutating || isSelectedPendingDeleted || !canMoveDown) return
    if (isRegionSelected && selectedRegionId) {
      moveRegion(selectedRegionId, 'down')
    } else if (isSubRegionSelected && selectedRegionId && selectedSubRegionId) {
      moveSubRegion(selectedRegionId, selectedSubRegionId, 'down')
    } else if (isNodeSelected && selectedNodeId) {
      moveNode(selectedNodeId, 'down')
    }
  }, [
    canMoveDown, isMutating, isRegionSelected, selectedRegionId, moveRegion,
    isSubRegionSelected, selectedSubRegionId, moveSubRegion,
    isNodeSelected, selectedNodeId, moveNode, isSelectedPendingDeleted,
  ])

  const handleCheckUrl = useCallback(async (target: 'new' | 'selected' = 'selected') => {
    const url = (target === 'new' ? newNodeUrlInput : nodeUrlInput).trim()
    if (!url) return
    const setStatus = target === 'new' ? setNewUrlCheckStatus : setUrlCheckStatus
    const setMsg = target === 'new' ? setNewUrlCheckMsg : setUrlCheckMsg
    setStatus('checking')
    setMsg('')
    try {
      const res = await fetch('/api/v1/explorer/exploration-node/link-check', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url }),
      })
      const data = await res.json() as { valid?: boolean; message?: string }
      if (res.ok && data.valid) {
        setStatus('ok')
        setMsg(copy.messages.urlValid)
      } else {
        setStatus('error')
        setMsg(data.message ?? copy.messages.urlInvalid)
      }
    } catch {
      setStatus('error')
      setMsg(copy.messages.urlCheckFailed)
    }
  }, [copy.messages.urlCheckFailed, copy.messages.urlInvalid, copy.messages.urlValid, nodeUrlInput, newNodeUrlInput])


  const handleOpenNodeEditModal = useCallback(() => {
    if (!selectedNode || isMutating || isSelectedPendingDeleted) return
    setNodeTitleInput(selectedNode.title ?? '')
    setNodeUrlInput(selectedNode.source_url ?? '')
    setUrlCheckStatus('idle')
    setUrlCheckMsg('')
    setRecommendationTarget('selected')
    setRecommendationQuery(buildRecommendationQuery('selected'))
    setEditNodeModalType(selectedNode.node_type)
  }, [
    buildRecommendationQuery,
    isMutating,
    isSelectedPendingDeleted,
    selectedNode,
    setRecommendationTarget,
    setRecommendationQuery,
  ])

  const handleCloseNodeEditModal = useCallback(() => {
    if (isMutating) return
    setEditNodeModalType(null)
  }, [isMutating])

  const handleSubmitNodeEdit = useCallback(async () => {
    if (!selectedNode || !selectedNodeId || isMutating) return

    if (selectedNode.node_type === 'exploration') {
      const source = selectedNodeSource(selectedNode, nodeUrlInput)
      if (!nodeTitleInput.trim() || !source) return
      await updateExplorationNode(selectedNodeId, nodeTitleInput.trim(), source)
    } else {
      if (!nodeTitleInput.trim()) return
      await updateResearchNode(selectedNodeId, nodeTitleInput.trim())
    }

    setEditNodeModalType(null)
  }, [
    isMutating,
    nodeTitleInput,
    nodeUrlInput,
    selectedNode,
    selectedNodeId,
    updateExplorationNode,
    updateResearchNode,
  ])

  const handleRenameCourse = useCallback(async () => {
    const trimmed = courseTitleInput.trim()
    if (!trimmed || isMutating) return
    if (trimmed === courseAggregate?.title) return
    await updateCourseTitle(trimmed)
  }, [courseAggregate?.title, courseTitleInput, isMutating, updateCourseTitle])

  const handleRenameRegion = useCallback(async () => {
    if (isSubRegionSelected) {
      const trimmed = subRegionNameInput.trim()
      if (!selectedRegionId || !selectedSubRegionId || !trimmed || isMutating) return
      if (trimmed === selectedSubAgg?.subregion.name) return
      await updateSubRegion(selectedRegionId, selectedSubRegionId, trimmed)
      return
    }

    const trimmed = regionNameInput.trim()
    if (!selectedRegionId || !trimmed || isMutating) return
    if (trimmed === selectedRegionAgg?.region.name) return
    await updateRegion(selectedRegionId, trimmed)
  }, [
    isMutating,
    isSubRegionSelected,
    regionNameInput,
    selectedRegionAgg?.region.name,
    selectedRegionId,
    selectedSubAgg?.subregion.name,
    selectedSubRegionId,
    subRegionNameInput,
    updateRegion,
    updateSubRegion,
  ])

  const handleAddRegion = useCallback(async () => {
    if (!newRegionNameInput.trim() || isMutating) return
    await createRegion(newRegionNameInput.trim())
    setNewRegionNameInput('')
  }, [newRegionNameInput, createRegion, isMutating])

  const handleAddSubRegion = useCallback(async () => {
    if (!selectedRegionId || !newSubRegionNameInput.trim() || isMutating) return
    await createSubRegion(selectedRegionId, newSubRegionNameInput.trim())
    setNewSubRegionNameInput('')
    setShowAddSubRegion(false)
    setAddModalTypeState(null)
  }, [selectedRegionId, newSubRegionNameInput, createSubRegion, isMutating])

  const handleAddNode = useCallback(async (modeOverride?: 'exploration' | 'research') => {
    const nextMode = modeOverride ?? addNodeMode
    if (!nextMode || !newNodeTitleInput.trim() || isMutating) return

    let parentKind: ParentKind
    let parentId: string
    if (isSubRegionSelected && selectedSubRegionId) {
      parentKind = 'subregion'
      parentId = selectedSubRegionId
    } else if (isRegionSelected && selectedRegionId) {
      parentKind = 'region'
      parentId = selectedRegionId
    } else return

    if (nextMode === 'research') {
      await createResearchNode(parentKind, parentId, newNodeTitleInput.trim())
    } else {
      if (!newNodeUrlInput.trim()) return
      await createExplorationNode(parentKind, parentId, newNodeTitleInput.trim(), newNodeUrlInput.trim())
    }
    setNewNodeTitleInput('')
    setNewNodeUrlInput('')
    setAddNodeMode(null)
    setAddModalTypeState(null)
  }, [
    addNodeMode, newNodeTitleInput, newNodeUrlInput, isMutating,
    isSubRegionSelected, selectedSubRegionId, isRegionSelected, selectedRegionId,
    createResearchNode, createExplorationNode,
  ])

  return {
    isCourseSelected,
    isRegionSelected,
    isSubRegionSelected,
    isNodeSelected,
    selectedRegionAgg,
    selectedSubAgg,
    selectedNode,
    selectedNodeRegionAgg,
    selectedNodeSubAgg,
    canDelete,
    canDeleteCourse,
    canActivateCourse,
    hasChildren,
    isSelectedInactive,
    isSelectedPendingCreated,
    isSelectedPendingDeleted,
    isActionBarBlocked,
    canAddSubRegion,
    canMoveUp,
    canMoveDown,
    showInactiveItems,
    courseTitleInput,
    setCourseTitleInput,
    newRegionNameInput,
    setNewRegionNameInput,
    regionNameInput,
    setRegionNameInput,
    showAddSubRegion,
    setShowAddSubRegion,
    newSubRegionNameInput,
    setNewSubRegionNameInput,
    subRegionNameInput,
    setSubRegionNameInput,
    addNodeMode,
    setAddNodeMode,
    newNodeTitleInput,
    setNewNodeTitleInput,
    newNodeUrlInput,
    setNewNodeUrlInput,
    nodeTitleInput,
    setNodeTitleInput,
    nodeUrlInput,
    setNodeUrlInput,
    urlCheckStatus,
    setUrlCheckStatus,
    urlCheckMsg,
    setUrlCheckMsg,
    newUrlCheckStatus,
    newUrlCheckMsg,
    isRecommendationOpen,
    recommendationQuery,
    setRecommendationQuery,
    recommendationCandidates,
    isLoadingRecommendations,
    recommendationMessage,
    recommendationPointError,
    deleteConfirm,
    setDeleteConfirmInput,
    addModalType,
    setAddModalType,
    editNodeModalType,
    handleSave,
    handleCancel,
    handleDelete,
    handleActivateCourse,
    handleConfirmDelete,
    handleCloseDeleteConfirm,
    handleInactiveToggle,
    handleMoveUp,
    handleMoveDown,
    handleCheckUrl,
    handleOpenRecommendationModal,
    handleCloseRecommendationModal,
    handleSearchRecommendations,
    handleApplyRecommendation,
    handleOpenNodeEditModal,
    handleCloseNodeEditModal,
    handleSubmitNodeEdit,
    handleRenameCourse,
    handleRenameRegion,
    handleAddRegion,
    handleAddSubRegion,
    handleAddNode,
    isMutating,
    mutationMessage,
  }
}
