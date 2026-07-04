'use client'

import { useCallback, useEffect, useState } from 'react'
import type { CourseAggregate, ParentKind, SourceType } from './explorerPlanTypes'
import {
  type LegacyDraftAggregate,
  type PendingExplorerPlanChanges,
  type PendingItemState,
  emptyPendingChanges,
  hasPendingChanges,
  pendingItemState,
} from './explorerPlanInternalUtils'
import { useExplorerPlanLoader } from './useExplorerPlanLoader'
import { useExplorerPlanCourse } from './useExplorerPlanCourse'
import { useExplorerPlanRegion } from './useExplorerPlanRegion'
import { useExplorerPlanSubRegion } from './useExplorerPlanSubRegion'
import { useExplorerPlanNode } from './useExplorerPlanNode'
import { useExplorerPlanSave } from './useExplorerPlanSave'

export type ExplorationNodeSourceInput = {
  sourceType: SourceType
  sourceUrl?: string | null
  contentId?: string | null
}

export type { PendingItemState } from './explorerPlanInternalUtils'

export interface ExplorerPlanState {
  courseAggregate: CourseAggregate | null
  isLoading: boolean
  error: string | null
  selectedRegionId: string | null
  selectedSubRegionId: string | null
  selectedNodeId: string | null
  handleSelectCourse: () => void
  handleSelectRegion: (regionId: string) => void
  handleSelectSubRegion: (regionId: string, subRegionId: string) => void
  handleSelectNode: (nodeId: string) => void
  updateCourseTitle: (title: string) => Promise<void>
  deleteCourseDraft: (confirmTitle: string) => Promise<void>
  activateCourseDraft: (confirmTitle: string) => Promise<void>
  getRegionPendingState: (regionId: string) => PendingItemState
  getSubRegionPendingState: (subRegionId: string) => PendingItemState
  getNodePendingState: (nodeId: string) => PendingItemState
  cancelPlanChanges: () => Promise<void>
  restoreDeletedRegion: (regionId: string) => void
  restoreDeletedSubRegion: (subRegionId: string) => void
  restoreDeletedNode: (nodeId: string) => void
  createRegion: (name: string) => Promise<void>
  updateRegion: (regionId: string, name: string) => Promise<void>
  deleteRegion: (regionId: string) => Promise<void>
  moveRegion: (regionId: string, direction: 'up' | 'down') => void
  createSubRegion: (regionId: string, name: string) => Promise<void>
  updateSubRegion: (regionId: string, subRegionId: string, name: string) => Promise<void>
  deleteSubRegion: (subRegionId: string) => Promise<void>
  moveSubRegion: (regionId: string, subRegionId: string, direction: 'up' | 'down') => void
  createResearchNode: (parentKind: ParentKind, parentId: string, title: string) => Promise<void>
  updateResearchNode: (nodeId: string, title: string) => Promise<void>
  deleteResearchNode: (nodeId: string) => Promise<void>
  createExplorationNode: (parentKind: ParentKind, parentId: string, title: string, source: string | ExplorationNodeSourceInput) => Promise<void>
  updateExplorationNode: (nodeId: string, title: string, source: string | ExplorationNodeSourceInput) => Promise<void>
  deleteExplorationNode: (nodeId: string) => Promise<void>
  moveNode: (nodeId: string, direction: 'up' | 'down') => void
  savePlanChanges: () => Promise<void>
  showInactiveItems: boolean
  toggleInactiveItems: () => void
  activateRegion: (regionId: string) => Promise<void>
  activateSubRegion: (subRegionId: string) => Promise<void>
  activateNode: (nodeId: string) => Promise<void>
  isDirty: boolean
  isMutating: boolean
  mutationMessage: string | null
  reload: () => void
}

export function useExplorerPlan(courseDraftId: string, legacyDraft?: LegacyDraftAggregate | null): ExplorerPlanState {
  const [courseAggregate, setCourseAggregate] = useState<CourseAggregate | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedRegionId, setSelectedRegionId] = useState<string | null>(null)
  const [selectedSubRegionId, setSelectedSubRegionId] = useState<string | null>(null)
  const [selectedNodeId, setSelectedNodeId] = useState<string | null>(null)
  const [isMutating, setIsMutating] = useState(false)
  const [mutationMessage, setMutationMessage] = useState<string | null>(null)
  const [pendingChanges, setPendingChanges] = useState<PendingExplorerPlanChanges>(() => emptyPendingChanges())
  const [showInactiveItems, setShowInactiveItems] = useState(false)

  const { load } = useExplorerPlanLoader({
    courseDraftId,
    legacyDraft,
    showInactiveItems,
    setCourseAggregate,
    setPendingChanges,
    setMutationMessage,
    setIsLoading,
    setError,
  })

  useEffect(() => {
    void load()
  }, [load])

  const { handleSelectCourse, updateCourseTitle, deleteCourseDraft, activateCourseDraft } = useExplorerPlanCourse({
    courseAggregate,
    courseDraftId,
    isMutating,
    setCourseAggregate,
    setPendingChanges,
    setIsMutating,
    setMutationMessage,
    setSelectedRegionId,
    setSelectedSubRegionId,
    setSelectedNodeId,
  })

  const handleSelectRegion = useCallback((regionId: string) => {
    setSelectedRegionId(regionId)
    setSelectedSubRegionId(null)
    setSelectedNodeId(null)
  }, [])

  const handleSelectSubRegion = useCallback((regionId: string, subRegionId: string) => {
    setSelectedRegionId(regionId)
    setSelectedSubRegionId(subRegionId)
    setSelectedNodeId(null)
  }, [])

  const handleSelectNode = useCallback((nodeId: string) => {
    setSelectedNodeId(nodeId)
  }, [])

  const getRegionPendingState = useCallback((regionId: string): PendingItemState => {
    return pendingItemState(
      pendingChanges.createdRegionIds,
      pendingChanges.deletedRegionIds,
      pendingChanges.movedRegionIds,
      regionId,
    )
  }, [pendingChanges])

  const getSubRegionPendingState = useCallback((subRegionId: string): PendingItemState => {
    return pendingItemState(
      pendingChanges.createdSubRegionIds,
      pendingChanges.deletedSubRegionIds,
      pendingChanges.movedSubRegionIds,
      subRegionId,
    )
  }, [pendingChanges])

  const getNodePendingState = useCallback((nodeId: string): PendingItemState => {
    return pendingItemState(
      pendingChanges.createdNodeIds,
      pendingChanges.deletedNodeIds,
      pendingChanges.movedNodeIds,
      nodeId,
    )
  }, [pendingChanges])

  const cancelPlanChanges = useCallback(async () => {
    if (isMutating) return
    setSelectedRegionId(null)
    setSelectedSubRegionId(null)
    setSelectedNodeId(null)
    await load()
  }, [isMutating, load])

  const { createRegion, updateRegion, deleteRegion, restoreDeletedRegion, moveRegion, activateRegion } = useExplorerPlanRegion({
    courseAggregate,
    courseDraftId,
    isMutating,
    setCourseAggregate,
    setPendingChanges,
    setIsMutating,
    setMutationMessage,
    setSelectedRegionId,
    setSelectedSubRegionId,
    setSelectedNodeId,
    handleSelectRegion,
    load,
  })

  const { createSubRegion, updateSubRegion, deleteSubRegion, restoreDeletedSubRegion, moveSubRegion, activateSubRegion } = useExplorerPlanSubRegion({
    courseAggregate,
    isMutating,
    setCourseAggregate,
    setPendingChanges,
    setIsMutating,
    setMutationMessage,
    setSelectedSubRegionId,
    setSelectedNodeId,
    handleSelectSubRegion,
    load,
  })

  const {
    createResearchNode,
    createExplorationNode,
    updateResearchNode,
    updateExplorationNode,
    deleteResearchNode,
    deleteExplorationNode,
    restoreDeletedNode,
    activateNode,
    moveNode,
  } = useExplorerPlanNode({
    courseAggregate,
    isMutating,
    setCourseAggregate,
    setPendingChanges,
    setIsMutating,
    setMutationMessage,
    setSelectedNodeId,
    load,
  })

  const toggleInactiveItems = useCallback(() => {
    setShowInactiveItems((current) => !current)
  }, [])


  const { savePlanChanges } = useExplorerPlanSave({
    courseAggregate,
    courseDraftId,
    pendingChanges,
    isMutating,
    selectedRegionId,
    selectedSubRegionId,
    selectedNodeId,
    setIsMutating,
    setMutationMessage,
    setSelectedRegionId,
    setSelectedSubRegionId,
    setSelectedNodeId,
    load,
  })

  return {
    courseAggregate,
    isLoading,
    error,
    selectedRegionId,
    selectedSubRegionId,
    selectedNodeId,
    handleSelectCourse,
    handleSelectRegion,
    handleSelectSubRegion,
    handleSelectNode,
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
    isDirty: hasPendingChanges(pendingChanges),
    isMutating,
    mutationMessage,
    reload: load,
  }
}
