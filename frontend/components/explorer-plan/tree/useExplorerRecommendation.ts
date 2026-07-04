'use client'

import { useState, useCallback } from 'react'
import type { Dispatch, SetStateAction } from 'react'
import type { CourseAggregate, ExplorerNode, RegionAggregate, SubRegionAggregate, ParentKind } from '../explorerPlanTypes'
import type { ExplorationNodeSourceInput } from '../useExplorerPlan'
import type { ExplorerContentCandidate, ExplorerRecommendationContext } from './explorerEditFormTypes'
import type { DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'
import { candidateToExplorationSource, withTrimmedValue } from './explorerEditFormUtils'
import { buildDefaultRecommendationQuery } from '@/lib/recommendation/recommendationQueryBuilder'

function addTrimmedUnique(target: Set<string>, value?: string | null) {
  const trimmed = value?.trim()
  if (trimmed) target.add(trimmed)
}

function buildSelectedExplorerContentExclusion(courseAggregate: CourseAggregate | null): Pick<ExplorerRecommendationContext, 'excluded_content_ids' | 'excluded_urls'> {
  const contentIds = new Set<string>()
  const urls = new Set<string>()

  for (const regionAgg of courseAggregate?.regions ?? []) {
    const nodes = [
      ...regionAgg.nodes,
      ...regionAgg.subregions.flatMap((subAgg) => subAgg.nodes),
    ]
    for (const node of nodes) {
      if (node.node_type !== 'exploration') continue
      addTrimmedUnique(contentIds, node.content_id)
      addTrimmedUnique(urls, node.source_url)
    }
  }

  return {
    excluded_content_ids: [...contentIds],
    excluded_urls: [...urls],
  }
}

function findSubRegionAggregate(courseAggregate: CourseAggregate | null, subRegionId?: string | null): SubRegionAggregate | null {
  const targetId = withTrimmedValue(subRegionId)
  if (!targetId) return null
  for (const regionAgg of courseAggregate?.regions ?? []) {
    const found = regionAgg.subregions.find((subAgg) => subAgg.subregion.id === targetId)
    if (found) return found
  }
  return null
}

function resolveSubRegionLessonID(subAgg: SubRegionAggregate | null | undefined): string | undefined {
  return withTrimmedValue(subAgg?.subregion.course_draft_lesson_id) ?? withTrimmedValue(subAgg?.subregion.id)
}

function findRegionAggregate(courseAggregate: CourseAggregate | null, regionId?: string | null): RegionAggregate | null {
  const targetId = withTrimmedValue(regionId)
  if (!targetId) return null
  return courseAggregate?.regions.find((regionAgg) => regionAgg.region.id === targetId) ?? null
}

function resolveRegionLessonID(regionAgg: RegionAggregate | null | undefined): string | undefined {
  return withTrimmedValue(regionAgg?.region.course_draft_lesson_id) ?? withTrimmedValue(regionAgg?.region.id)
}

function resolveRecommendationLessonID(
  target: 'new' | 'selected',
  courseAggregate: CourseAggregate | null,
  selectedNode: ExplorerNode | null,
  selectedRegionAgg: RegionAggregate | null,
  selectedSubAgg: SubRegionAggregate | null,
): string | undefined {
  if (target === 'selected' && selectedNode?.parent_kind === 'subregion') {
    const parentSubAgg = findSubRegionAggregate(courseAggregate, selectedNode.parent_id)
    return resolveSubRegionLessonID(parentSubAgg) ?? withTrimmedValue(selectedNode.parent_id)
  }
  if (target === 'selected' && selectedNode?.parent_kind === 'region') {
    const parentRegionAgg = findRegionAggregate(courseAggregate, selectedNode.parent_id)
    return resolveRegionLessonID(parentRegionAgg) ?? withTrimmedValue(selectedNode.parent_id)
  }
  return resolveSubRegionLessonID(selectedSubAgg) ?? resolveRegionLessonID(selectedRegionAgg)
}

export function useExplorerRecommendation({
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
}: {
  courseAggregate: CourseAggregate | null
  selectedRegionAgg: RegionAggregate | null
  selectedSubAgg: SubRegionAggregate | null
  selectedNode: ExplorerNode | null
  nodeTitleInput: string
  newNodeTitleInput: string
  isRegionSelected: boolean
  isSubRegionSelected: boolean
  selectedRegionId: string | null
  selectedSubRegionId: string | null
  createExplorationNode: (parentKind: ParentKind, parentId: string, title: string, source: string | ExplorationNodeSourceInput) => Promise<void>
  updateExplorationNode: (nodeId: string, title: string, source: string | ExplorationNodeSourceInput) => Promise<void>
  setAddNodeMode: Dispatch<SetStateAction<'exploration' | 'research' | null>>
  setNewNodeTitleInput: Dispatch<SetStateAction<string>>
  setNewNodeUrlInput: Dispatch<SetStateAction<string>>
  setAddModalTypeState: Dispatch<SetStateAction<'child-object' | 'exploration-choice' | 'exploration-url' | 'research' | null>>
  setNodeTitleInput: Dispatch<SetStateAction<string>>
  setNodeUrlInput: Dispatch<SetStateAction<string>>
  setUrlCheckStatus: Dispatch<SetStateAction<'idle' | 'checking' | 'ok' | 'error'>>
  setUrlCheckMsg: Dispatch<SetStateAction<string>>
  copy: DashboardCourseDraftCopy['tree']['editPanel']['modals']
}) {
  const [isRecommendationOpen, setIsRecommendationOpen] = useState(false)
  const [recommendationQuery, setRecommendationQuery] = useState('')
  const [recommendationCandidates, setRecommendationCandidates] = useState<ExplorerContentCandidate[]>([])
  const [isLoadingRecommendations, setIsLoadingRecommendations] = useState(false)
  const [recommendationMessage, setRecommendationMessage] = useState<string | null>(null)
  const [recommendationPointError, setRecommendationPointError] = useState<string | null>(null)
  const [recommendationTarget, setRecommendationTarget] = useState<'new' | 'selected'>('selected')

  const buildRecommendationQuery = useCallback((target: 'new' | 'selected' = 'new') => buildDefaultRecommendationQuery({
    target,
    courseTitle: courseAggregate?.title,
    regionTitle: selectedRegionAgg?.region.name,
    regionDescription: selectedRegionAgg?.region.description,
    subRegionTitle: selectedSubAgg?.subregion.name,
    subRegionDescription: selectedSubAgg?.subregion.description,
    nodeTitleInput: target === 'selected' ? nodeTitleInput : newNodeTitleInput,
    selectedNodeTitle: selectedNode?.title,
    nodeSummary: selectedNode?.summary,
  }), [
    courseAggregate?.title,
    selectedRegionAgg?.region.description,
    selectedRegionAgg?.region.name,
    selectedSubAgg?.subregion.description,
    selectedSubAgg?.subregion.name,
    newNodeTitleInput,
    nodeTitleInput,
    selectedNode?.summary,
    selectedNode?.title,
  ])

  const buildRecommendationContext = useCallback((target: 'new' | 'selected'): ExplorerRecommendationContext => ({
    ...buildSelectedExplorerContentExclusion(courseAggregate),
    course_draft_id: withTrimmedValue(courseAggregate?.course_draft_id),
    region_title: withTrimmedValue(selectedRegionAgg?.region.name),
    region_description: withTrimmedValue(selectedRegionAgg?.region.description),
    subregion_title: withTrimmedValue(selectedSubAgg?.subregion.name),
    subregion_description: withTrimmedValue(selectedSubAgg?.subregion.description),
    node_title: target === 'selected'
      ? withTrimmedValue(selectedNode?.title) ?? withTrimmedValue(nodeTitleInput)
      : withTrimmedValue(newNodeTitleInput),
    node_summary: target === 'selected'
      ? withTrimmedValue(selectedNode?.summary)
      : undefined,
    lesson_id: resolveRecommendationLessonID(target, courseAggregate, selectedNode, selectedRegionAgg, selectedSubAgg),
  }), [
    courseAggregate,
    courseAggregate?.course_draft_id,
    newNodeTitleInput,
    nodeTitleInput,
    selectedNode?.summary,
    selectedNode?.title,
    selectedNode?.parent_kind,
    selectedNode?.parent_id,
    selectedRegionAgg?.region.course_draft_lesson_id,
    selectedRegionAgg?.region.description,
    selectedRegionAgg?.region.id,
    selectedRegionAgg?.region.name,
    selectedSubAgg?.subregion.course_draft_lesson_id,
    selectedSubAgg?.subregion.description,
    selectedSubAgg?.subregion.id,
    selectedSubAgg?.subregion.name,
  ])

  const loadRecommendationCandidates = useCallback(async (query: string, limit = 3) => {
    setIsLoadingRecommendations(true)
    setRecommendationMessage(null)
    setRecommendationPointError(null)

    try {
      const recommendationContext = buildRecommendationContext(recommendationTarget)
      const response = await fetch('/api/v1/explorer/recommend', {
        method: 'POST',
        credentials: 'include',
        cache: 'no-store',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          query: query.trim(),
          max_results: limit,
          ...recommendationContext,
        }),
      })
      const payload = (await response.json().catch(() => ({}))) as {
        candidates?: ExplorerContentCandidate[]
        billing_status?: string
        point_preview?: { cost: number; total_balance: number }
        error?: string
      }
      if (!response.ok) {
        if (response.status === 401) throw new Error(copy.recommendation.loginRequired)
        if (payload.error === 'insufficient_points') {
          const balance = payload.point_preview?.total_balance ?? 0
          const cost = payload.point_preview?.cost ?? 1
          setRecommendationPointError(copy.recommendation.insufficientPoints(balance, cost))
          setIsLoadingRecommendations(false)
          return
        }
        throw new Error(payload.error ?? copy.recommendation.loadFailed)
      }
      const candidates = (payload.candidates ?? []).map((candidate, index) => ({
        ...candidate,
        id: candidate.id ?? candidate.content_id ?? candidate.url ?? candidate.canonical_url ?? `${candidate.title}-${index}`,
      }))
      setRecommendationCandidates(candidates)
      setRecommendationMessage(
        candidates.length > 0
          ? copy.recommendation.loaded(candidates.length)
          : copy.recommendation.noResults
      )
    } catch (err) {
      setRecommendationCandidates([])
      setRecommendationMessage(err instanceof Error ? err.message : copy.recommendation.loadFailed)
    } finally {
      setIsLoadingRecommendations(false)
    }
  }, [buildRecommendationContext, copy.recommendation, recommendationTarget])

  const handleOpenRecommendationModal = useCallback(async (target: 'new' | 'selected' = 'selected') => {
    const nextQuery = buildRecommendationQuery(target)
    setRecommendationTarget(target)
    setRecommendationQuery(nextQuery)
    setRecommendationCandidates([])
    setRecommendationMessage(null)
    setRecommendationPointError(null)
    setIsRecommendationOpen(true)
    if (target === 'new') setAddNodeMode('exploration')
  }, [buildRecommendationQuery, setAddNodeMode])

  const handleCloseRecommendationModal = useCallback(() => {
    setIsRecommendationOpen(false)
  }, [])

  const handleSearchRecommendations = useCallback(async () => {
    await loadRecommendationCandidates(recommendationQuery, 6)
  }, [loadRecommendationCandidates, recommendationQuery])

  const handleApplyRecommendation = useCallback(async (candidate: ExplorerContentCandidate) => {
    const nextSource = candidateToExplorationSource(candidate)
    const nextUrl = nextSource?.sourceUrl ?? ''
    setRecommendationCandidates((current) => current.filter((item) => item.id !== candidate.id))
    if (recommendationTarget === 'new') {
      const title = candidate.title.trim()
      if (title && nextSource) {
        let parentKind: ParentKind
        let parentId: string
        if (isSubRegionSelected && selectedSubRegionId) {
          parentKind = 'subregion'
          parentId = selectedSubRegionId
        } else if (isRegionSelected && selectedRegionId) {
          parentKind = 'region'
          parentId = selectedRegionId
        } else {
          setIsRecommendationOpen(false)
          return
        }
        await createExplorationNode(parentKind, parentId, title, nextSource)
        setNewNodeTitleInput('')
        setNewNodeUrlInput('')
        setAddNodeMode(null)
        setAddModalTypeState(null)
      } else {
        setRecommendationMessage(copy.recommendation.missingCandidate)
        return
      }
      setIsRecommendationOpen(false)
      return
    }
    // 기존 선택된 노드 수정
    const title = candidate.title.trim()
    if (title) setNodeTitleInput(title)
    if (nextSource && selectedNode?.node_type === 'exploration') {
      await updateExplorationNode(selectedNode.id, title || selectedNode.title, nextSource)
    }
    if (nextUrl) {
      setNodeUrlInput(nextUrl)
      setUrlCheckStatus('idle')
      setUrlCheckMsg('')
    } else if (nextSource?.contentId) {
      setNodeUrlInput('')
      setUrlCheckStatus('idle')
      setUrlCheckMsg(copy.messages.contentIdBased)
    }
    setIsRecommendationOpen(false)
  }, [
    recommendationTarget, isSubRegionSelected, selectedSubRegionId,
    isRegionSelected, selectedRegionId, createExplorationNode,
    selectedNode, updateExplorationNode,
    setNewNodeTitleInput, setNewNodeUrlInput, setAddNodeMode, setAddModalTypeState,
    setNodeTitleInput, setNodeUrlInput, setUrlCheckStatus, setUrlCheckMsg,
    copy.messages.contentIdBased,
    copy.recommendation.missingCandidate,
  ])

  return {
    // State (read-only for consumers)
    isRecommendationOpen,
    recommendationQuery,
    recommendationCandidates,
    isLoadingRecommendations,
    recommendationMessage,
    recommendationPointError,
    // Setters exposed for main hook (setAddModalType, handleOpenNodeEditModal)
    setIsRecommendationOpen,
    setRecommendationTarget,
    setRecommendationQuery,
    setRecommendationCandidates,
    setRecommendationMessage,
    setRecommendationPointError,
    // Computed (used in setAddModalType, handleOpenNodeEditModal)
    buildRecommendationQuery,
    // Handlers
    handleOpenRecommendationModal,
    handleCloseRecommendationModal,
    handleSearchRecommendations,
    handleApplyRecommendation,
  }
}
