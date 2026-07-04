'use client'

import { useCallback } from 'react'
import type { Dispatch, SetStateAction } from 'react'
import type { CourseAggregate, ExplorerNode, ParentKind } from './explorerPlanTypes'
import type { ExplorationNodeSourceInput } from './useExplorerPlan'
import {
  type PendingExplorerPlanChanges,
  readErrorMessage,
  checkExplorationNodeLink,
  clonePendingChanges,
  markNodeUpdated,
  markNodeMoved,
  markNodeDeleted,
  byOrderIndex,
  nowIso,
  makeLocalResearchNode,
  makeLocalExplorationNode,
  normalizeExplorationNodeSource,
  shouldCheckExplorationNodeUrl,
} from './explorerPlanInternalUtils'

export function useExplorerPlanNode({
  courseAggregate,
  isMutating,
  setCourseAggregate,
  setPendingChanges,
  setIsMutating,
  setMutationMessage,
  setSelectedNodeId,
  load,
}: {
  courseAggregate: CourseAggregate | null
  isMutating: boolean
  setCourseAggregate: Dispatch<SetStateAction<CourseAggregate | null>>
  setPendingChanges: Dispatch<SetStateAction<PendingExplorerPlanChanges>>
  setIsMutating: Dispatch<SetStateAction<boolean>>
  setMutationMessage: Dispatch<SetStateAction<string | null>>
  setSelectedNodeId: Dispatch<SetStateAction<string | null>>
  load: () => Promise<void>
}) {
  const createResearchNode = useCallback(async (parentKind: ParentKind, parentId: string, title: string) => {
    const trimmed = title.trim()
    if (!parentId || !trimmed || isMutating) return

    let orderIndex: number | null = null
    if (parentKind === 'region') {
      const regionAgg = courseAggregate?.regions.find((item) => item.region.id === parentId)
      orderIndex = regionAgg?.nodes.length ?? null
    } else {
      const subAgg = courseAggregate?.regions
        .flatMap((regionAgg) => regionAgg.subregions)
        .find((item) => item.subregion.id === parentId)
      orderIndex = subAgg?.nodes.length ?? null
    }

    if (orderIndex === null) {
      setMutationMessage('상위 항목을 찾지 못했습니다.')
      return
    }

    const created = makeLocalResearchNode(parentKind, parentId, trimmed, orderIndex)
    setCourseAggregate((current) => {
      if (!current) return current
      return {
        ...current,
        regions: current.regions.map((regionAgg) => {
          if (parentKind === 'region' && regionAgg.region.id === parentId) {
            return { ...regionAgg, nodes: [...regionAgg.nodes, created] }
          }
          if (parentKind === 'subregion') {
            return {
              ...regionAgg,
              subregions: regionAgg.subregions.map((subAgg) => {
                if (subAgg.subregion.id !== parentId) return subAgg
                return { ...subAgg, nodes: [...subAgg.nodes, created] }
              }),
            }
          }
          return regionAgg
        }),
      }
    })

    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      next.createdNodeIds.add(created.id)
      return next
    })
    setSelectedNodeId(created.id)
    setMutationMessage('연구지점을 화면에 추가했습니다. 최종 저장 전까지 DB에는 반영되지 않습니다.')
  }, [courseAggregate?.regions, isMutating, setCourseAggregate, setPendingChanges, setMutationMessage, setSelectedNodeId])

  const createExplorationNode = useCallback(async (
    parentKind: ParentKind,
    parentId: string,
    title: string,
    source: string | ExplorationNodeSourceInput
  ) => {
    const trimmedTitle = title.trim()
    const sourceInput = normalizeExplorationNodeSource(source)
    const trimmedUrl = sourceInput.sourceUrl?.trim() ?? ''
    const trimmedContentId = sourceInput.contentId?.trim() ?? ''
    if (!parentId || !trimmedTitle || isMutating) return

    const nextSource: ExplorationNodeSourceInput = {
      sourceType: sourceInput.sourceType,
      sourceUrl: trimmedUrl || null,
      contentId: trimmedContentId || null,
    }

    if (shouldCheckExplorationNodeUrl(nextSource)) {
      if (!trimmedUrl) return
      if (!trimmedUrl.startsWith('http://') && !trimmedUrl.startsWith('https://')) {
        setMutationMessage('탐험지점 링크는 http:// 또는 https://로 시작해야 합니다.')
        return
      }

      try {
        nextSource.sourceUrl = await checkExplorationNodeLink(trimmedUrl)
      } catch (err) {
        setMutationMessage(err instanceof Error ? err.message : '탐험지점 링크 확인 실패')
        return
      }
    } else if (!nextSource.contentId) {
      setMutationMessage('내부 콘텐츠 기반 탐험지점은 콘텐츠 ID가 필요합니다.')
      return
    }

    let orderIndex: number | null = null
    if (parentKind === 'region') {
      const regionAgg = courseAggregate?.regions.find((item) => item.region.id === parentId)
      orderIndex = regionAgg?.nodes.length ?? null
    } else {
      const subAgg = courseAggregate?.regions
        .flatMap((regionAgg) => regionAgg.subregions)
        .find((item) => item.subregion.id === parentId)
      orderIndex = subAgg?.nodes.length ?? null
    }

    if (orderIndex === null) {
      setMutationMessage('상위 항목을 찾지 못했습니다.')
      return
    }

    const created = makeLocalExplorationNode(parentKind, parentId, trimmedTitle, nextSource, orderIndex)
    setCourseAggregate((current) => {
      if (!current) return current
      return {
        ...current,
        regions: current.regions.map((regionAgg) => {
          if (parentKind === 'region' && regionAgg.region.id === parentId) {
            return { ...regionAgg, nodes: [...regionAgg.nodes, created] }
          }
          if (parentKind === 'subregion') {
            return {
              ...regionAgg,
              subregions: regionAgg.subregions.map((subAgg) => {
                if (subAgg.subregion.id !== parentId) return subAgg
                return { ...subAgg, nodes: [...subAgg.nodes, created] }
              }),
            }
          }
          return regionAgg
        }),
      }
    })

    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      next.createdNodeIds.add(created.id)
      return next
    })
    setSelectedNodeId(created.id)
    setMutationMessage(
      nextSource.contentId
        ? '콘텐츠 기반 탐험지점을 화면에 추가했습니다. 최종 저장 전까지 DB에는 반영되지 않습니다.'
        : '탐험지점 링크를 확인하고 화면에 추가했습니다. 최종 저장 전까지 DB에는 반영되지 않습니다.'
    )
  }, [courseAggregate?.regions, isMutating, setCourseAggregate, setPendingChanges, setMutationMessage, setSelectedNodeId])

  const updateResearchNode = useCallback(async (nodeId: string, title: string) => {
    const trimmed = title.trim()
    if (!nodeId || !trimmed || isMutating) return
    setCourseAggregate((current) => {
      if (!current) return current
      return {
        ...current,
        regions: current.regions.map((regionAgg) => ({
          ...regionAgg,
          nodes: regionAgg.nodes.map((node) =>
            node.id === nodeId ? { ...node, title: trimmed, updated_at: nowIso() } : node
          ),
          subregions: regionAgg.subregions.map((subAgg) => ({
            ...subAgg,
            nodes: subAgg.nodes.map((node) =>
              node.id === nodeId ? { ...node, title: trimmed, updated_at: nowIso() } : node
            ),
          })),
        })),
      }
    })
    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      if (!next.createdNodeIds.has(nodeId)) next.updatedNodeIds.add(nodeId)
      return next
    })
    setSelectedNodeId(nodeId)
    setMutationMessage('연구지점 제목을 화면에 반영했습니다. 최종 저장 전까지 DB에는 반영되지 않습니다.')
  }, [isMutating, setCourseAggregate, setPendingChanges, setMutationMessage, setSelectedNodeId])

  const updateExplorationNode = useCallback(async (nodeId: string, title: string, source: string | ExplorationNodeSourceInput) => {
    const trimmedTitle = title.trim()
    const sourceInput = normalizeExplorationNodeSource(source)
    const trimmedUrl = sourceInput.sourceUrl?.trim() ?? ''
    const trimmedContentId = sourceInput.contentId?.trim() ?? ''
    if (!nodeId || !trimmedTitle || isMutating) return

    const nextSource: ExplorationNodeSourceInput = {
      sourceType: sourceInput.sourceType,
      sourceUrl: trimmedUrl || null,
      contentId: trimmedContentId || null,
    }

    if (shouldCheckExplorationNodeUrl(nextSource)) {
      if (!trimmedUrl) return
      if (!trimmedUrl.startsWith('http://') && !trimmedUrl.startsWith('https://')) {
        setMutationMessage('탐험지점 링크는 http:// 또는 https://로 시작해야 합니다.')
        return
      }

      try {
        nextSource.sourceUrl = await checkExplorationNodeLink(trimmedUrl)
      } catch (err) {
        setMutationMessage(err instanceof Error ? err.message : '탐험지점 링크 확인 실패')
        return
      }
    } else if (!nextSource.contentId) {
      setMutationMessage('내부 콘텐츠 기반 탐험지점은 콘텐츠 ID가 필요합니다.')
      return
    }

    setCourseAggregate((current) => {
      if (!current) return current
      return {
        ...current,
        regions: current.regions.map((regionAgg) => ({
          ...regionAgg,
          nodes: regionAgg.nodes.map((node) =>
            node.id === nodeId
              ? {
                  ...node,
                  title: trimmedTitle,
                  source_url: nextSource.sourceUrl,
                  source_type: nextSource.sourceType,
                  content_id: nextSource.contentId,
                  updated_at: nowIso(),
                }
              : node
          ),
          subregions: regionAgg.subregions.map((subAgg) => ({
            ...subAgg,
            nodes: subAgg.nodes.map((node) =>
              node.id === nodeId
                ? {
                    ...node,
                    title: trimmedTitle,
                    source_url: nextSource.sourceUrl,
                    source_type: nextSource.sourceType,
                    content_id: nextSource.contentId,
                    updated_at: nowIso(),
                  }
                : node
            ),
          })),
        })),
      }
    })
    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      if (!next.createdNodeIds.has(nodeId)) next.updatedNodeIds.add(nodeId)
      return next
    })
    setSelectedNodeId(nodeId)
    setMutationMessage(
      nextSource.contentId
        ? '콘텐츠 기반 탐험지점을 화면에 반영했습니다. 최종 저장 전까지 DB에는 반영되지 않습니다.'
        : '탐험지점 링크를 확인하고 화면에 반영했습니다. 최종 저장 전까지 DB에는 반영되지 않습니다.'
    )
  }, [isMutating, setCourseAggregate, setPendingChanges, setMutationMessage, setSelectedNodeId])

  const deleteResearchNode = useCallback(async (nodeId: string) => {
    if (!nodeId || isMutating) return
    const nodeParent = courseAggregate?.regions.find((regionAgg) =>
      regionAgg.nodes.some((node) => node.id === nodeId) ||
      regionAgg.subregions.some((subAgg) => subAgg.nodes.some((node) => node.id === nodeId))
    )
    const remainingSiblingIds = (() => {
      const regionNodes = nodeParent?.nodes ?? []
      if (regionNodes.some((node) => node.id === nodeId)) {
        return regionNodes.filter((node) => node.id !== nodeId).map((node) => node.id)
      }
      const subAgg = nodeParent?.subregions.find((item) => item.nodes.some((node) => node.id === nodeId))
      return subAgg?.nodes.filter((node) => node.id !== nodeId).map((node) => node.id) ?? []
    })()
    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      markNodeDeleted(next, nodeId)
      remainingSiblingIds.forEach((id) => markNodeUpdated(next, id))
      return next
    })
    setMutationMessage('지점을 삭제 예정으로 표시했습니다. 학습 시작 전 코스는 최종 저장 시 완전히 삭제됩니다.')
  }, [courseAggregate?.regions, isMutating, setPendingChanges, setMutationMessage])

  const restoreDeletedNode = useCallback((nodeId: string) => {
    if (!nodeId || isMutating) return
    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      next.deletedNodeIds.delete(nodeId)
      return next
    })
    setMutationMessage('지점 삭제 예정 표시를 취소했습니다.')
  }, [isMutating, setPendingChanges, setMutationMessage])

  const deleteExplorationNode = useCallback(async (nodeId: string) => {
    await deleteResearchNode(nodeId)
  }, [deleteResearchNode])

  const activateNode = useCallback(async (nodeId: string) => {
    if (!nodeId || isMutating) return
    setIsMutating(true)
    setMutationMessage(null)
    try {
      const res = await fetch(`/api/v1/explorer/node/${nodeId}/status`, {
        method: 'PATCH',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status: 'active' }),
      })
      if (!res.ok) throw new Error(await readErrorMessage(res, '지점을 활성화하지 못했습니다.'))
      await load()
      setSelectedNodeId(nodeId)
      setMutationMessage('지점을 활성화했습니다.')
    } catch (err) {
      setMutationMessage(err instanceof Error ? err.message : '지점 활성화 실패')
    } finally {
      setIsMutating(false)
    }
  }, [isMutating, load, setIsMutating, setMutationMessage, setSelectedNodeId])

  const moveNode = useCallback((nodeId: string, direction: 'up' | 'down') => {
    if (!nodeId || isMutating) return
    let movedIds: string[] = []
    setCourseAggregate((current) => {
      if (!current) return current
      const timestamp = nowIso()

      const reorderNodes = (nodes: ExplorerNode[]) =>
        nodes.map((node, index) => ({ ...node, order_index: index, updated_at: timestamp }))

      const regions = [...current.regions]
        .sort((a, b) => byOrderIndex(a.region, b.region))
        .map((regionAgg) => ({
          ...regionAgg,
          subregions: [...regionAgg.subregions]
            .sort((a, b) => byOrderIndex(a.subregion, b.subregion))
            .map((subAgg) => ({ ...subAgg, nodes: [...subAgg.nodes].sort(byOrderIndex) })),
          nodes: [...regionAgg.nodes].sort(byOrderIndex),
        }))

      type ContainerRef = {
        regionIndex: number
        parentKind: ParentKind
        parentId: string
        subIndex?: number
      }

      const getNodes = (c: ContainerRef): ExplorerNode[] =>
        c.parentKind === 'region'
          ? regions[c.regionIndex].nodes
          : regions[c.regionIndex].subregions[c.subIndex!].nodes

      const setNodes = (c: ContainerRef, nodes: ExplorerNode[]) => {
        if (c.parentKind === 'region') {
          regions[c.regionIndex] = { ...regions[c.regionIndex], nodes }
          return
        }
        const ra = regions[c.regionIndex]
        regions[c.regionIndex] = {
          ...ra,
          subregions: ra.subregions.map((s, i) =>
            i === c.subIndex ? { ...s, nodes } : s
          ),
        }
      }

      // 모든 컨테이너를 시각 순서대로 나열 (노드 없는 빈 컨테이너 포함)
      // 시각 순서: 각 지역별 [서브지역0, 서브지역1, ..., 지역 직속]
      const flatContainers: ContainerRef[] = []
      regions.forEach((ra, regionIndex) => {
        ra.subregions.forEach((s, subIndex) => {
          flatContainers.push({
            regionIndex,
            parentKind: 'subregion',
            parentId: s.subregion.id,
            subIndex,
          })
        })
        flatContainers.push({
          regionIndex,
          parentKind: 'region',
          parentId: ra.region.id,
        })
      })

      // 노드가 속한 컨테이너와 컨테이너 내 인덱스 탐색
      let sourceContainer: ContainerRef | null = null
      let sourceIndexInContainer = -1
      for (const c of flatContainers) {
        const idx = getNodes(c).findIndex((n) => n.id === nodeId)
        if (idx >= 0) {
          sourceContainer = c
          sourceIndexInContainer = idx
          break
        }
      }
      if (!sourceContainer || sourceIndexInContainer < 0) return current

      const sourceNodes = getNodes(sourceContainer)
      const sameContainerTargetIndex =
        direction === 'up' ? sourceIndexInContainer - 1 : sourceIndexInContainer + 1

      if (sameContainerTargetIndex >= 0 && sameContainerTargetIndex < sourceNodes.length) {
        // 동일 컨테이너 형제 교환
        const next = [...sourceNodes]
        const [movedNode] = next.splice(sourceIndexInContainer, 1)
        next.splice(sameContainerTargetIndex, 0, movedNode)
        movedIds = [nodeId, sourceNodes[sameContainerTargetIndex].id]
        setNodes(sourceContainer, reorderNodes(next))
      } else {
        // 인접 컨테이너로 이동
        const srcContainerIndex = flatContainers.findIndex(
          (c) => c.parentKind === sourceContainer!.parentKind && c.parentId === sourceContainer!.parentId
        )
        const dstContainerIndex = direction === 'up' ? srcContainerIndex - 1 : srcContainerIndex + 1
        if (dstContainerIndex < 0 || dstContainerIndex >= flatContainers.length) return current

        const destinationContainer = flatContainers[dstContainerIndex]
        const destinationNodes = getNodes(destinationContainer)

        const nextSource = [...sourceNodes]
        const [movedNode] = nextSource.splice(sourceIndexInContainer, 1)

        // ▲이면 목적 컨테이너 끝에, ▼이면 목적 컨테이너 앞에 삽입
        const nextDest = [...destinationNodes]
        const insertAt = direction === 'up' ? nextDest.length : 0
        nextDest.splice(insertAt, 0, {
          ...movedNode,
          parent_kind: destinationContainer.parentKind,
          parent_id: destinationContainer.parentId,
          updated_at: timestamp,
        })

        const reorderedSource = reorderNodes(nextSource)
        const reorderedDest = reorderNodes(nextDest)
        movedIds = [
          ...new Set([
            movedNode.id,
            ...reorderedSource.map((n) => n.id),
            ...reorderedDest.map((n) => n.id),
          ]),
        ]
        setNodes(sourceContainer, reorderedSource)
        setNodes(destinationContainer, reorderedDest)
      }

      return { ...current, regions }
    })
    if (movedIds.length === 0) return
    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      movedIds.forEach((id) => markNodeUpdated(next, id))
      markNodeMoved(next, nodeId)
      return next
    })
    setSelectedNodeId(nodeId)
    setMutationMessage('지점 순서를 화면에 반영했습니다. 최종 저장 전까지 DB에는 반영되지 않습니다.')
  }, [isMutating, setCourseAggregate, setPendingChanges, setMutationMessage, setSelectedNodeId])

  return {
    createResearchNode,
    createExplorationNode,
    updateResearchNode,
    updateExplorationNode,
    deleteResearchNode,
    deleteExplorationNode,
    restoreDeletedNode,
    activateNode,
    moveNode,
  }
}
