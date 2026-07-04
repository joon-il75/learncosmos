'use client'

import { useCallback } from 'react'
import type { Dispatch, SetStateAction } from 'react'
import type { CourseAggregate } from './explorerPlanTypes'
import {
  type PendingExplorerPlanChanges,
  readErrorMessage,
  clonePendingChanges,
  markSubRegionUpdated,
  markSubRegionMoved,
  markSubRegionDeleted,
  markNodeDeleted,
  byOrderIndex,
  nowIso,
  makeLocalSubRegion,
} from './explorerPlanInternalUtils'

export function useExplorerPlanSubRegion({
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
}: {
  courseAggregate: CourseAggregate | null
  isMutating: boolean
  setCourseAggregate: Dispatch<SetStateAction<CourseAggregate | null>>
  setPendingChanges: Dispatch<SetStateAction<PendingExplorerPlanChanges>>
  setIsMutating: Dispatch<SetStateAction<boolean>>
  setMutationMessage: Dispatch<SetStateAction<string | null>>
  setSelectedSubRegionId: Dispatch<SetStateAction<string | null>>
  setSelectedNodeId: Dispatch<SetStateAction<string | null>>
  handleSelectSubRegion: (regionId: string, subRegionId: string) => void
  load: () => Promise<void>
}) {
  const createSubRegion = useCallback(async (regionId: string, name: string) => {
    const trimmed = name.trim()
    if (!regionId || !trimmed || isMutating) return

    const regionAgg = courseAggregate?.regions.find((item) => item.region.id === regionId)
    if (!regionAgg) {
      setMutationMessage('상위 지역을 찾지 못했습니다.')
      return
    }
    if (regionAgg.subregions.length >= 3) {
      setMutationMessage('서브지역은 지역마다 최대 3개까지 만들 수 있습니다.')
      return
    }

    const created = makeLocalSubRegion(regionId, trimmed, regionAgg.subregions.length)
    setCourseAggregate((current) => {
      if (!current) return current
      return {
        ...current,
        regions: current.regions.map((item) => {
          if (item.region.id !== regionId) return item
          return {
            ...item,
            subregions: [...item.subregions, { subregion: created, nodes: [] }],
          }
        }),
      }
    })

    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      next.createdSubRegionIds.add(created.id)
      return next
    })
    handleSelectSubRegion(regionId, created.id)
    setMutationMessage('서브지역을 화면에 추가했습니다. 최종 저장 전까지 DB에는 반영되지 않습니다.')
  }, [courseAggregate?.regions, handleSelectSubRegion, isMutating, setCourseAggregate, setPendingChanges, setMutationMessage])

  const updateSubRegion = useCallback(async (regionId: string, subRegionId: string, name: string) => {
    const trimmed = name.trim()
    if (!regionId || !subRegionId || !trimmed || isMutating) return
    setCourseAggregate((current) => {
      if (!current) return current
      return {
        ...current,
        regions: current.regions.map((regionAgg) => {
          if (regionAgg.region.id !== regionId) return regionAgg
          return {
            ...regionAgg,
            subregions: regionAgg.subregions.map((subAgg) =>
              subAgg.subregion.id === subRegionId
                ? { ...subAgg, subregion: { ...subAgg.subregion, name: trimmed, updated_at: nowIso() } }
                : subAgg
            ),
          }
        }),
      }
    })
    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      if (!next.createdSubRegionIds.has(subRegionId)) next.updatedSubRegionIds.add(subRegionId)
      return next
    })
    handleSelectSubRegion(regionId, subRegionId)
    setMutationMessage('서브지역 이름을 화면에 반영했습니다. 최종 저장 전까지 DB에는 반영되지 않습니다.')
  }, [handleSelectSubRegion, isMutating, setCourseAggregate, setPendingChanges, setMutationMessage])

  const deleteSubRegion = useCallback(async (subRegionId: string) => {
    if (!subRegionId || isMutating) return
    const parentRegionAgg = courseAggregate?.regions.find((regionAgg) =>
      regionAgg.subregions.some((subAgg) => subAgg.subregion.id === subRegionId)
    )
    const deletedSubAgg = parentRegionAgg?.subregions.find((subAgg) => subAgg.subregion.id === subRegionId) ?? null
    const remainingSiblingIds = parentRegionAgg?.subregions
      .filter((subAgg) => subAgg.subregion.id !== subRegionId)
      .map((subAgg) => subAgg.subregion.id) ?? []
    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      markSubRegionDeleted(next, subRegionId)
      deletedSubAgg?.nodes.forEach((node) => markNodeDeleted(next, node.id))
      remainingSiblingIds.forEach((id) => markSubRegionUpdated(next, id))
      return next
    })
    setMutationMessage(
      (deletedSubAgg?.nodes.length ?? 0) > 0
        ? '서브지역과 하위 지점을 삭제 예정으로 표시했습니다. 학습 시작 전 코스는 최종 저장 시 완전히 삭제됩니다.'
        : '서브지역을 삭제 예정으로 표시했습니다. 학습 시작 전 코스는 최종 저장 시 완전히 삭제됩니다.'
    )
  }, [courseAggregate?.regions, isMutating, setPendingChanges, setMutationMessage])

  const restoreDeletedSubRegion = useCallback((subRegionId: string) => {
    if (!subRegionId || isMutating) return
    const subAgg = courseAggregate?.regions
      .flatMap((regionAgg) => regionAgg.subregions)
      .find((item) => item.subregion.id === subRegionId)
    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      next.deletedSubRegionIds.delete(subRegionId)
      subAgg?.nodes.forEach((node) => next.deletedNodeIds.delete(node.id))
      return next
    })
    setMutationMessage('서브지역 삭제 예정 표시를 취소했습니다.')
  }, [courseAggregate?.regions, isMutating, setPendingChanges, setMutationMessage])

  const moveSubRegion = useCallback((regionId: string, subRegionId: string, direction: 'up' | 'down') => {
    if (!regionId || !subRegionId || isMutating) return
    let movedSubIds: string[] = []
    let nextRegionId = regionId
    setCourseAggregate((current) => {
      if (!current) return current
      const regions = [...current.regions].sort((a, b) => byOrderIndex(a.region, b.region))
      const regionIndex = regions.findIndex((r) => r.region.id === regionId)
      if (regionIndex < 0) return current
      const regionAgg = regions[regionIndex]
      const subregions = [...regionAgg.subregions].sort((a, b) => byOrderIndex(a.subregion, b.subregion))
      const fromIndex = subregions.findIndex((s) => s.subregion.id === subRegionId)
      if (fromIndex < 0) return current

      const toIndex = direction === 'up' ? fromIndex - 1 : fromIndex + 1

      if (toIndex >= 0 && toIndex < subregions.length) {
        // 동일 지역 내 교환
        const next = [...subregions]
        const [moved] = next.splice(fromIndex, 1)
        next.splice(toIndex, 0, moved)
        movedSubIds = [subregions[fromIndex].subregion.id, subregions[toIndex].subregion.id]
        regions[regionIndex] = {
          ...regionAgg,
          subregions: next.map((s, i) => ({
            ...s,
            subregion: { ...s.subregion, order_index: i, updated_at: nowIso() },
          })),
        }
      } else {
        // 크로스 지역 이동
        const targetRegionIndex = direction === 'up' ? regionIndex - 1 : regionIndex + 1
        if (targetRegionIndex < 0 || targetRegionIndex >= regions.length) return current
        const targetRegionAgg = regions[targetRegionIndex]
        const targetSubs = [...targetRegionAgg.subregions].sort((a, b) => byOrderIndex(a.subregion, b.subregion))

        // 원본 지역에서 제거
        const remainingSubs = subregions.filter((s) => s.subregion.id !== subRegionId)
        const movedSubAgg = subregions[fromIndex]
        regions[regionIndex] = {
          ...regionAgg,
          subregions: remainingSubs.map((s, i) => ({
            ...s,
            subregion: { ...s.subregion, order_index: i, updated_at: nowIso() },
          })),
        }

        // 목적 지역에 삽입 (▲이면 끝, ▼이면 앞)
        const insertIndex = direction === 'up' ? targetSubs.length : 0
        const updatedSubAgg = {
          ...movedSubAgg,
          subregion: {
            ...movedSubAgg.subregion,
            region_id: targetRegionAgg.region.id,
            updated_at: nowIso(),
          },
        }
        const nextTargetSubs = [...targetSubs]
        nextTargetSubs.splice(insertIndex, 0, updatedSubAgg)
        regions[targetRegionIndex] = {
          ...targetRegionAgg,
          subregions: nextTargetSubs.map((s, i) => ({
            ...s,
            subregion: { ...s.subregion, order_index: i, updated_at: nowIso() },
          })),
        }

        movedSubIds = [
          subRegionId,
          ...remainingSubs.map((s) => s.subregion.id),
          ...nextTargetSubs.map((s) => s.subregion.id),
        ]
        nextRegionId = targetRegionAgg.region.id
      }

      return { ...current, regions }
    })
    if (movedSubIds.length === 0) return
    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      movedSubIds.forEach((id) => {
        markSubRegionUpdated(next, id)
        markSubRegionMoved(next, id)
      })
      return next
    })
    handleSelectSubRegion(nextRegionId, subRegionId)
    setMutationMessage('서브지역 순서를 화면에 반영했습니다. 최종 저장 전까지 DB에는 반영되지 않습니다.')
  }, [handleSelectSubRegion, isMutating, setCourseAggregate, setPendingChanges, setMutationMessage])

  const activateSubRegion = useCallback(async (subRegionId: string) => {
    if (!subRegionId || isMutating) return
    setIsMutating(true)
    setMutationMessage(null)
    try {
      const res = await fetch(`/api/v1/explorer/subregion/${subRegionId}/status`, {
        method: 'PATCH',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status: 'active' }),
      })
      if (!res.ok) throw new Error(await readErrorMessage(res, '서브지역을 활성화하지 못했습니다.'))
      await load()
      setSelectedSubRegionId(subRegionId)
      setSelectedNodeId(null)
      setMutationMessage('서브지역을 활성화했습니다.')
    } catch (err) {
      setMutationMessage(err instanceof Error ? err.message : '서브지역 활성화 실패')
    } finally {
      setIsMutating(false)
    }
  }, [isMutating, load, setIsMutating, setMutationMessage, setSelectedSubRegionId, setSelectedNodeId])

  return { createSubRegion, updateSubRegion, deleteSubRegion, restoreDeletedSubRegion, moveSubRegion, activateSubRegion }
}
