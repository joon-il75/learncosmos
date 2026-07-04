'use client'

import { useCallback } from 'react'
import type { Dispatch, SetStateAction } from 'react'
import type { CourseAggregate } from './explorerPlanTypes'
import {
  type PendingExplorerPlanChanges,
  readErrorMessage,
  clonePendingChanges,
  markRegionUpdated,
  markRegionMoved,
  markSubRegionDeleted,
  markNodeDeleted,
  byOrderIndex,
  nowIso,
  makeLocalRegion,
} from './explorerPlanInternalUtils'

export function useExplorerPlanRegion({
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
}: {
  courseAggregate: CourseAggregate | null
  courseDraftId: string
  isMutating: boolean
  setCourseAggregate: Dispatch<SetStateAction<CourseAggregate | null>>
  setPendingChanges: Dispatch<SetStateAction<PendingExplorerPlanChanges>>
  setIsMutating: Dispatch<SetStateAction<boolean>>
  setMutationMessage: Dispatch<SetStateAction<string | null>>
  setSelectedRegionId: Dispatch<SetStateAction<string | null>>
  setSelectedSubRegionId: Dispatch<SetStateAction<string | null>>
  setSelectedNodeId: Dispatch<SetStateAction<string | null>>
  handleSelectRegion: (regionId: string) => void
  load: () => Promise<void>
}) {
  const createRegion = useCallback(async (name: string) => {
    const trimmed = name.trim()
    if (!courseDraftId || !trimmed || isMutating) return
    const region = makeLocalRegion(courseDraftId, trimmed, courseAggregate?.regions.length ?? 0)
    setCourseAggregate((current) => {
      if (!current) return current
      return {
        ...current,
        regions: [...current.regions, { region, subregions: [], nodes: [] }],
      }
    })
    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      next.createdRegionIds.add(region.id)
      return next
    })
    handleSelectRegion(region.id)
    setMutationMessage('지역을 화면에 추가했습니다. 최종 저장 전까지 DB에는 반영되지 않습니다.')
  }, [courseAggregate?.regions.length, courseDraftId, handleSelectRegion, isMutating, setCourseAggregate, setPendingChanges, setMutationMessage])

  const updateRegion = useCallback(async (regionId: string, name: string) => {
    const trimmed = name.trim()
    if (!regionId || !trimmed || isMutating) return
    setCourseAggregate((current) => {
      if (!current) return current
      return {
        ...current,
        regions: current.regions.map((regionAgg) =>
          regionAgg.region.id === regionId
            ? { ...regionAgg, region: { ...regionAgg.region, name: trimmed, updated_at: nowIso() } }
            : regionAgg
        ),
      }
    })
    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      if (!next.createdRegionIds.has(regionId)) next.updatedRegionIds.add(regionId)
      return next
    })
    handleSelectRegion(regionId)
    setMutationMessage('지역 이름을 화면에 반영했습니다. 최종 저장 전까지 DB에는 반영되지 않습니다.')
  }, [handleSelectRegion, isMutating, setCourseAggregate, setPendingChanges, setMutationMessage])

  const deleteRegion = useCallback(async (regionId: string) => {
    if (!regionId || isMutating) return
    const deletedRegionAgg = courseAggregate?.regions.find((regionAgg) => regionAgg.region.id === regionId) ?? null
    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      next.deletedRegionIds.add(regionId)
      next.updatedRegionIds.delete(regionId)
      deletedRegionAgg?.nodes.forEach((node) => markNodeDeleted(next, node.id))
      deletedRegionAgg?.subregions.forEach((subAgg) => {
        markSubRegionDeleted(next, subAgg.subregion.id)
        subAgg.nodes.forEach((node) => markNodeDeleted(next, node.id))
      })
      return next
    })
    const childCount = (deletedRegionAgg?.nodes.length ?? 0) +
      (deletedRegionAgg?.subregions.length ?? 0) +
      (deletedRegionAgg?.subregions.reduce((sum, subAgg) => sum + subAgg.nodes.length, 0) ?? 0)
    setMutationMessage(
      childCount > 0
        ? '지역과 하위 항목을 삭제 예정으로 표시했습니다. 학습 시작 전 코스는 최종 저장 시 완전히 삭제됩니다.'
        : '지역을 삭제 예정으로 표시했습니다. 학습 시작 전 코스는 최종 저장 시 완전히 삭제됩니다.'
    )
  }, [courseAggregate?.regions, isMutating, setPendingChanges, setMutationMessage])

  const restoreDeletedRegion = useCallback((regionId: string) => {
    if (!regionId || isMutating) return
    const regionAgg = courseAggregate?.regions.find((item) => item.region.id === regionId)
    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      next.deletedRegionIds.delete(regionId)
      regionAgg?.nodes.forEach((node) => next.deletedNodeIds.delete(node.id))
      regionAgg?.subregions.forEach((subAgg) => {
        next.deletedSubRegionIds.delete(subAgg.subregion.id)
        subAgg.nodes.forEach((node) => next.deletedNodeIds.delete(node.id))
      })
      return next
    })
    setMutationMessage('지역 삭제 예정 표시를 취소했습니다.')
  }, [courseAggregate?.regions, isMutating, setPendingChanges, setMutationMessage])

  const moveRegion = useCallback((regionId: string, direction: 'up' | 'down') => {
    if (!regionId || isMutating) return
    let movedIds: string[] = []
    setCourseAggregate((current) => {
      if (!current) return current
      const regions = [...current.regions].sort((a, b) => byOrderIndex(a.region, b.region))
      const fromIndex = regions.findIndex((regionAgg) => regionAgg.region.id === regionId)
      const toIndex = direction === 'up' ? fromIndex - 1 : fromIndex + 1
      if (fromIndex < 0 || toIndex < 0 || toIndex >= regions.length) return current
      const nextRegions = [...regions]
      const [moved] = nextRegions.splice(fromIndex, 1)
      nextRegions.splice(toIndex, 0, moved)
      movedIds = [regions[fromIndex].region.id, regions[toIndex].region.id]
      return {
        ...current,
        regions: nextRegions.map((regionAgg, index) => ({
          ...regionAgg,
          region: { ...regionAgg.region, order_index: index, updated_at: nowIso() },
        })),
      }
    })
    if (movedIds.length === 0) return
    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      movedIds.forEach((id) => {
        markRegionUpdated(next, id)
        markRegionMoved(next, id)
      })
      return next
    })
    handleSelectRegion(regionId)
    setMutationMessage('지역 순서를 화면에 반영했습니다. 최종 저장 전까지 DB에는 반영되지 않습니다.')
  }, [handleSelectRegion, isMutating, setCourseAggregate, setPendingChanges, setMutationMessage])

  const activateRegion = useCallback(async (regionId: string) => {
    if (!regionId || isMutating) return
    setIsMutating(true)
    setMutationMessage(null)
    try {
      const res = await fetch(`/api/v1/explorer/region/${regionId}/status`, {
        method: 'PATCH',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status: 'active', with_children: true }),
      })
      if (!res.ok) throw new Error(await readErrorMessage(res, '지역을 활성화하지 못했습니다.'))
      await load()
      setSelectedRegionId(regionId)
      setSelectedSubRegionId(null)
      setSelectedNodeId(null)
      setMutationMessage('지역과 하위 항목을 활성화했습니다.')
    } catch (err) {
      setMutationMessage(err instanceof Error ? err.message : '지역 활성화 실패')
    } finally {
      setIsMutating(false)
    }
  }, [isMutating, load, setIsMutating, setMutationMessage, setSelectedRegionId, setSelectedSubRegionId, setSelectedNodeId])

  return { createRegion, updateRegion, deleteRegion, restoreDeletedRegion, moveRegion, activateRegion }
}
