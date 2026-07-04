'use client'

import { useCallback } from 'react'
import type { Dispatch, SetStateAction } from 'react'
import type { CourseAggregate, ExplorerNode, ExplorerRegion, ExplorerSubRegion } from './explorerPlanTypes'
import {
  type PendingExplorerPlanChanges,
  readErrorMessage,
  hasPendingChanges,
} from './explorerPlanInternalUtils'

export function useExplorerPlanSave({
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
}: {
  courseAggregate: CourseAggregate | null
  courseDraftId: string
  pendingChanges: PendingExplorerPlanChanges
  isMutating: boolean
  selectedRegionId: string | null
  selectedSubRegionId: string | null
  selectedNodeId: string | null
  setIsMutating: Dispatch<SetStateAction<boolean>>
  setMutationMessage: Dispatch<SetStateAction<string | null>>
  setSelectedRegionId: Dispatch<SetStateAction<string | null>>
  setSelectedSubRegionId: Dispatch<SetStateAction<string | null>>
  setSelectedNodeId: Dispatch<SetStateAction<string | null>>
  load: () => Promise<void>
}) {
  const savePlanChanges = useCallback(async () => {
    if (!courseAggregate || !hasPendingChanges(pendingChanges) || isMutating) return
    setIsMutating(true)
    setMutationMessage(null)

    const regionIdMap = new Map<string, string>()
    const subRegionIdMap = new Map<string, string>()
    const nodeIdMap = new Map<string, string>()

    try {
      const regionIsDeleted = (regionId: string) => pendingChanges.deletedRegionIds.has(regionId)
      const subRegionIsUnderDeletedRegion = (subRegionId: string) =>
        courseAggregate.regions.some((regionAgg) =>
          regionIsDeleted(regionAgg.region.id) &&
          regionAgg.subregions.some((subAgg) => subAgg.subregion.id === subRegionId)
        )
      const subRegionIsDeleted = (subRegionId: string) =>
        pendingChanges.deletedSubRegionIds.has(subRegionId) || subRegionIsUnderDeletedRegion(subRegionId)
      const nodeIsUnderDeletedParent = (node: ExplorerNode) =>
        node.parent_kind === 'region' ? regionIsDeleted(node.parent_id) : subRegionIsDeleted(node.parent_id)
      const nodeIsDeleted = (node: ExplorerNode) =>
        pendingChanges.deletedNodeIds.has(node.id) || nodeIsUnderDeletedParent(node)

      if (pendingChanges.updatedCourseTitle) {
        const title = courseAggregate.title.trim()
        if (!title) throw new Error('행성명을 입력해주세요.')
        const res = await fetch(`/api/v1/course-drafts/${courseDraftId}`, {
          method: 'PATCH',
          credentials: 'include',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ title }),
        })
        if (!res.ok) throw new Error(await readErrorMessage(res, '행성명을 저장하지 못했습니다.'))
      }

      for (const regionId of pendingChanges.deletedRegionIds) {
        if (pendingChanges.createdRegionIds.has(regionId)) continue
        const res = await fetch(`/api/v1/explorer/region/${regionId}`, {
          method: 'DELETE',
          credentials: 'include',
        })
        if (!res.ok) throw new Error(await readErrorMessage(res, '지역을 저장하지 못했습니다.'))
      }

      for (const regionAgg of courseAggregate.regions) {
        const region = regionAgg.region
        if (regionIsDeleted(region.id)) continue
        if (pendingChanges.createdRegionIds.has(region.id)) {
          const res = await fetch('/api/v1/explorer/region', {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              course_draft_id: courseDraftId,
              name: region.name,
              description: region.description,
              order_index: region.order_index,
            }),
          })
          if (!res.ok) throw new Error(await readErrorMessage(res, '지역을 저장하지 못했습니다.'))
          const data = await res.json() as { region?: ExplorerRegion }
          if (!data.region?.id) throw new Error('지역 저장 응답이 올바르지 않습니다.')
          regionIdMap.set(region.id, data.region.id)
        } else if (pendingChanges.updatedRegionIds.has(region.id)) {
          const res = await fetch(`/api/v1/explorer/region/${region.id}`, {
            method: 'PATCH',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name: region.name, description: region.description, order_index: region.order_index }),
          })
          if (!res.ok) throw new Error(await readErrorMessage(res, '지역을 저장하지 못했습니다.'))
        }
      }

      for (const subRegionId of pendingChanges.deletedSubRegionIds) {
        if (pendingChanges.createdSubRegionIds.has(subRegionId) || subRegionIsUnderDeletedRegion(subRegionId)) continue
        const res = await fetch(`/api/v1/explorer/subregion/${subRegionId}`, {
          method: 'DELETE',
          credentials: 'include',
        })
        if (!res.ok) throw new Error(await readErrorMessage(res, '서브지역을 저장하지 못했습니다.'))
      }

      for (const nodeId of pendingChanges.deletedNodeIds) {
        const node = courseAggregate.regions
          .flatMap((regionAgg) => [...regionAgg.nodes, ...regionAgg.subregions.flatMap((subAgg) => subAgg.nodes)])
          .find((item) => item.id === nodeId)
        if (pendingChanges.createdNodeIds.has(nodeId) || (node && nodeIsUnderDeletedParent(node))) continue
        const res = await fetch(`/api/v1/explorer/node/${nodeId}`, {
          method: 'DELETE',
          credentials: 'include',
        })
        if (!res.ok) throw new Error(await readErrorMessage(res, '지점을 저장하지 못했습니다.'))
      }

      for (const regionAgg of courseAggregate.regions) {
        if (regionIsDeleted(regionAgg.region.id)) continue
        const savedRegionId = regionIdMap.get(regionAgg.region.id) ?? regionAgg.region.id
        for (const subAgg of regionAgg.subregions) {
          const subregion = subAgg.subregion
          if (subRegionIsDeleted(subregion.id)) continue
          if (pendingChanges.createdSubRegionIds.has(subregion.id)) {
            const res = await fetch('/api/v1/explorer/subregion', {
              method: 'POST',
              credentials: 'include',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({
                region_id: savedRegionId,
                name: subregion.name,
                description: subregion.description,
                order_index: subregion.order_index,
              }),
            })
            if (!res.ok) throw new Error(await readErrorMessage(res, '서브지역을 저장하지 못했습니다.'))
            const data = await res.json() as { subregion?: ExplorerSubRegion }
            if (!data.subregion?.id) throw new Error('서브지역 저장 응답이 올바르지 않습니다.')
            subRegionIdMap.set(subregion.id, data.subregion.id)
          } else if (pendingChanges.updatedSubRegionIds.has(subregion.id)) {
            const res = await fetch(`/api/v1/explorer/subregion/${subregion.id}`, {
              method: 'PATCH',
              credentials: 'include',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({
                region_id: savedRegionId,
                name: subregion.name,
                description: subregion.description,
                order_index: subregion.order_index,
              }),
            })
            if (!res.ok) throw new Error(await readErrorMessage(res, '서브지역을 저장하지 못했습니다.'))
          }
        }
      }

      const allNodes = courseAggregate.regions.flatMap((regionAgg) => [
        ...regionAgg.nodes,
        ...regionAgg.subregions.flatMap((subAgg) => subAgg.nodes),
      ])
      for (const node of allNodes) {
        if (nodeIsDeleted(node)) continue
        const savedParentId = node.parent_kind === 'region'
          ? regionIdMap.get(node.parent_id) ?? node.parent_id
          : subRegionIdMap.get(node.parent_id) ?? node.parent_id

        if (pendingChanges.createdNodeIds.has(node.id)) {
          const isExploration = node.node_type === 'exploration'
          const res = await fetch(isExploration ? '/api/v1/explorer/exploration-node' : '/api/v1/explorer/research-node', {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(isExploration ? {
              parent_kind: node.parent_kind,
              parent_id: savedParentId,
              title: node.title,
              source_type: node.source_type ?? 'web',
              source_url: node.source_url,
              content_id: node.content_id,
              order_index: node.order_index,
            } : {
              parent_kind: node.parent_kind,
              parent_id: savedParentId,
              title: node.title,
              order_index: node.order_index,
            }),
          })
          if (!res.ok) throw new Error(await readErrorMessage(res, '지점을 저장하지 못했습니다.'))
          const data = await res.json() as { node?: ExplorerNode }
          if (!data.node?.id) throw new Error('지점 저장 응답이 올바르지 않습니다.')
          nodeIdMap.set(node.id, data.node.id)
        } else if (pendingChanges.updatedNodeIds.has(node.id)) {
          const res = await fetch(`/api/v1/explorer/node/${node.id}`, {
            method: 'PATCH',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              parent_kind: node.parent_kind,
              parent_id: savedParentId,
              title: node.title,
              source_type: node.node_type === 'exploration' ? node.source_type : undefined,
              source_url: node.node_type === 'exploration' ? node.source_url : undefined,
              content_id: node.node_type === 'exploration' ? node.content_id : undefined,
              order_index: node.order_index,
            }),
          })
          if (!res.ok) throw new Error(await readErrorMessage(res, '지점을 저장하지 못했습니다.'))
        }
      }

      const regionOrder = courseAggregate.regions
        .filter((regionAgg) => !regionIsDeleted(regionAgg.region.id))
        .map((regionAgg) => regionIdMap.get(regionAgg.region.id) ?? regionAgg.region.id)
      const res = await fetch('/api/v1/explorer/save', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          course_draft_id: courseDraftId,
          region_order: regionOrder,
        }),
      })
      if (!res.ok) throw new Error(await readErrorMessage(res, '탐험계획 저장에 실패했습니다.'))

      const nextSelectedRegionId = selectedRegionId ? regionIdMap.get(selectedRegionId) ?? selectedRegionId : null
      const nextSelectedSubRegionId = selectedSubRegionId
        ? subRegionIdMap.get(selectedSubRegionId) ?? selectedSubRegionId
        : null
      const nextSelectedNodeId = selectedNodeId ? nodeIdMap.get(selectedNodeId) ?? selectedNodeId : null
      await load()
      setSelectedRegionId(nextSelectedRegionId)
      setSelectedSubRegionId(nextSelectedSubRegionId)
      setSelectedNodeId(nextSelectedNodeId)
      setMutationMessage(null)
    } catch (err) {
      setMutationMessage(err instanceof Error ? err.message : '탐험계획 저장 실패')
    } finally {
      setIsMutating(false)
    }
  }, [courseAggregate, courseDraftId, isMutating, load, pendingChanges, selectedNodeId, selectedRegionId, selectedSubRegionId, setIsMutating, setMutationMessage, setSelectedNodeId, setSelectedRegionId, setSelectedSubRegionId])

  return { savePlanChanges }
}
