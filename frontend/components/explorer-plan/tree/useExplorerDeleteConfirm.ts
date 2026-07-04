'use client'

import { useState, useCallback, useEffect } from 'react'
import type { ExplorerNode, RegionAggregate, SubRegionAggregate } from '../explorerPlanTypes'
import type { DeleteConfirmState } from './explorerEditFormTypes'
import type { DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'

export function useExplorerDeleteConfirm({
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
}: {
  isMutating: boolean
  isModalBlocked: boolean
  courseAggregate: { title: string; regions: { region: { id: string }; subregions: { subregion: { id: string } }[] }[] } | null
  selectedRegionId: string | null
  selectedSubRegionId: string | null
  selectedNodeId: string | null
  selectedNode: ExplorerNode | null
  selectedRegionAgg: RegionAggregate | null
  selectedSubAgg: SubRegionAggregate | null
  isRegionSelected: boolean
  isSubRegionSelected: boolean
  isNodeSelected: boolean
  remainingRegionCount: number
  hasChildren: boolean
  canDeleteCourse: boolean
  canActivateCourse: boolean
  isSelectedPendingDeleted: boolean
  selectedRegionPendingState: string
  getRegionPendingState: (regionId: string) => string
  getSubRegionPendingState: (subRegionId: string) => string
  restoreDeletedRegion: (regionId: string) => void
  restoreDeletedSubRegion: (subRegionId: string) => void
  restoreDeletedNode: (nodeId: string) => void
  deleteRegion: (regionId: string) => Promise<void>
  deleteSubRegion: (subRegionId: string) => Promise<void>
  deleteExplorationNode: (nodeId: string) => Promise<void>
  deleteResearchNode: (nodeId: string) => Promise<void>
  deleteCourseDraft: (confirmTitle: string) => Promise<void>
  activateCourseDraft: (confirmTitle: string) => Promise<void>
  copy: DashboardCourseDraftCopy['tree']['editPanel']['modals']
}) {
  const [deleteConfirm, setDeleteConfirm] = useState<DeleteConfirmState>(null)

  // selectedNodeId 변경 시 삭제 확인 모달 초기화
  useEffect(() => {
    setDeleteConfirm(null)
  }, [selectedNodeId]) // eslint-disable-line react-hooks/exhaustive-deps

  const setDeleteConfirmInput = useCallback((value: string) => {
    setDeleteConfirm((current) => current ? { ...current, confirmInput: value } : current)
  }, [])

  const handleDelete = useCallback(async () => {
    if (isMutating || isModalBlocked) return
    if (isSelectedPendingDeleted) {
      if (isRegionSelected && selectedRegionId) {
        restoreDeletedRegion(selectedRegionId)
      } else if (isSubRegionSelected && selectedSubRegionId) {
        if (selectedRegionPendingState === 'deleted' && selectedRegionId) {
          restoreDeletedRegion(selectedRegionId)
        } else {
          restoreDeletedSubRegion(selectedSubRegionId)
        }
      } else if (isNodeSelected && selectedNodeId) {
        if (selectedNode?.parent_kind === 'region' && getRegionPendingState(selectedNode.parent_id) === 'deleted') {
          restoreDeletedRegion(selectedNode.parent_id)
        } else if (selectedNode?.parent_kind === 'subregion') {
          const parentRegionAgg = courseAggregate?.regions.find((regionAgg) =>
            regionAgg.subregions.some((subAgg) => subAgg.subregion.id === selectedNode.parent_id)
          )
          if (parentRegionAgg && getRegionPendingState(parentRegionAgg.region.id) === 'deleted') {
            restoreDeletedRegion(parentRegionAgg.region.id)
          } else if (getSubRegionPendingState(selectedNode.parent_id) === 'deleted') {
            restoreDeletedSubRegion(selectedNode.parent_id)
          } else {
            restoreDeletedNode(selectedNodeId)
          }
        } else {
          restoreDeletedNode(selectedNodeId)
        }
      }
      return
    }
    if (canDeleteCourse && courseAggregate) {
      setDeleteConfirm({
        mode: 'confirm',
        action: 'delete',
        target: 'course',
        title: courseAggregate.title,
        message: copy.messages.deleteCourse,
        requireTitleInput: true,
        confirmInput: '',
      })
      return
    }
    if (isRegionSelected && selectedRegionId) {
      if (remainingRegionCount <= 1) {
        setDeleteConfirm({
          mode: 'alert',
          action: 'delete',
          target: 'region',
          title: selectedRegionAgg?.region.name ?? copy.messages.lastRegionTitle,
          message: copy.messages.lastRegionBlock,
          requireTitleInput: false,
          confirmInput: '',
        })
        return
      }
      setDeleteConfirm({
        mode: 'confirm',
        action: 'delete',
        target: 'region',
        title: selectedRegionAgg?.region.name ?? copy.messages.fallbackRegion,
        message: hasChildren
          ? copy.messages.deleteRegionWithChildren
          : copy.messages.deleteRegion,
        requireTitleInput: false,
        confirmInput: '',
      })
    } else if (isSubRegionSelected && selectedSubRegionId) {
      setDeleteConfirm({
        mode: 'confirm',
        action: 'delete',
        target: 'subregion',
        title: selectedSubAgg?.subregion.name ?? copy.messages.fallbackSubregion,
        message: hasChildren
          ? copy.messages.deleteSubregionWithChildren
          : copy.messages.deleteSubregion,
        requireTitleInput: false,
        confirmInput: '',
      })
    } else if (isNodeSelected && selectedNodeId && selectedNode) {
      setDeleteConfirm({
        mode: 'confirm',
        action: 'delete',
        target: 'node',
        title: selectedNode.title,
        message: copy.messages.deleteNode,
        requireTitleInput: false,
        confirmInput: '',
      })
    }
  }, [
    isMutating, isModalBlocked, canDeleteCourse, courseAggregate,
    hasChildren, isRegionSelected, selectedRegionId, selectedRegionAgg?.region.name, remainingRegionCount,
    isSubRegionSelected, selectedSubRegionId, selectedSubAgg?.subregion.name,
    isNodeSelected, selectedNodeId, selectedNode,
    isSelectedPendingDeleted, restoreDeletedRegion, restoreDeletedSubRegion, restoreDeletedNode,
    selectedRegionPendingState, getRegionPendingState, getSubRegionPendingState,
    copy.messages,
  ])

  const handleActivateCourse = useCallback(async () => {
    if (isMutating || isModalBlocked || !canActivateCourse || !courseAggregate) return
    setDeleteConfirm({
      mode: 'confirm',
      action: 'activate',
      target: 'course',
      title: courseAggregate.title,
      message: copy.messages.activateCourse,
      requireTitleInput: true,
      confirmInput: '',
    })
  }, [canActivateCourse, courseAggregate, copy.messages.activateCourse, isModalBlocked, isMutating])

  const handleCloseDeleteConfirm = useCallback(() => {
    if (isMutating) return
    setDeleteConfirm(null)
  }, [isMutating])

  const handleConfirmDelete = useCallback(async () => {
    if (!deleteConfirm || isMutating) return

    if (deleteConfirm.target === 'course' && deleteConfirm.action === 'activate') {
      await activateCourseDraft(deleteConfirm.confirmInput)
      return
    }

    if (deleteConfirm.target === 'course') {
      await deleteCourseDraft(deleteConfirm.confirmInput)
      return
    }

    setDeleteConfirm(null)
    if (deleteConfirm.target === 'region' && selectedRegionId) {
      await deleteRegion(selectedRegionId)
    } else if (deleteConfirm.target === 'subregion' && selectedSubRegionId) {
      await deleteSubRegion(selectedSubRegionId)
    } else if (deleteConfirm.target === 'node' && selectedNodeId && selectedNode) {
      if (selectedNode.node_type === 'exploration') {
        await deleteExplorationNode(selectedNodeId)
      } else {
        await deleteResearchNode(selectedNodeId)
      }
    }
  }, [
    deleteConfirm, isMutating, activateCourseDraft, deleteCourseDraft,
    selectedRegionId, deleteRegion,
    selectedSubRegionId, deleteSubRegion,
    selectedNodeId, selectedNode, deleteExplorationNode, deleteResearchNode,
  ])

  const resetDeleteConfirm = useCallback(() => {
    setDeleteConfirm(null)
  }, [])

  return {
    deleteConfirm,
    setDeleteConfirmInput,
    resetDeleteConfirm,
    handleDelete,
    handleActivateCourse,
    handleCloseDeleteConfirm,
    handleConfirmDelete,
  }
}
