import type {
  CourseAggregate,
  ExplorerNode,
  ExplorerRegion,
  ExplorerSubRegion,
  ParentKind,
  ResearchType,
  SourceType,
} from './explorerPlanTypes'
import type { ExplorationNodeSourceInput } from './useExplorerPlan'

// ── 내부 타입 ─────────────────────────────────────────────────────────────────

export type LegacyDraftPoint = {
  id: string
  point_type: 'exploration' | 'research'
  title: string
  description?: string | null
  template_type?: 'concept_summary' | 'practice_strategy' | 'problem_solving' | 'free_research' | null
  content_id?: string | null
  external_url?: string | null
  order_index: number
  created_at?: string
  updated_at?: string
}

export type LegacyDraftLessonTree = {
  lesson: {
    id: string
    course_draft_id?: string | null
    title: string
    summary?: string | null
    order_index?: number
    created_at?: string
    updated_at?: string
  }
  points?: Array<{ point: LegacyDraftPoint }>
  sub_lessons?: LegacyDraftLessonTree[]
}

export type LegacyDraftAggregate = {
  draft: {
    id: string
    title: string
    status?: 'draft' | 'confirmed' | 'learning' | 'archived'
    updated_at: string
  }
  lessons: LegacyDraftLessonTree[]
}

export type PendingExplorerPlanChanges = {
  updatedCourseTitle: boolean
  createdRegionIds: Set<string>
  updatedRegionIds: Set<string>
  deletedRegionIds: Set<string>
  movedRegionIds: Set<string>
  createdSubRegionIds: Set<string>
  updatedSubRegionIds: Set<string>
  deletedSubRegionIds: Set<string>
  movedSubRegionIds: Set<string>
  createdNodeIds: Set<string>
  updatedNodeIds: Set<string>
  deletedNodeIds: Set<string>
  movedNodeIds: Set<string>
}

// ── API 헬퍼 ─────────────────────────────────────────────────────────────────

export async function readErrorMessage(res: Response, fallback: string) {
  const data = await res.json().catch(() => ({})) as { error?: string }
  return data.error ?? fallback
}

export async function checkExplorationNodeLink(sourceUrl: string) {
  const res = await fetch('/api/v1/explorer/exploration-node/link-check', {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ url: sourceUrl }),
  })
  if (!res.ok) {
    throw new Error(await readErrorMessage(res, '탐험지점 링크를 확인하지 못했습니다.'))
  }
  const data = await res.json() as { valid?: boolean; url?: string; message?: string }
  if (!data.valid) {
    throw new Error(data.message ?? '탐험지점 링크가 유효하지 않습니다.')
  }
  return data.url ?? sourceUrl
}

// ── pending 변경 관리 ─────────────────────────────────────────────────────────

export function emptyPendingChanges(): PendingExplorerPlanChanges {
  return {
    updatedCourseTitle: false,
    createdRegionIds: new Set(),
    updatedRegionIds: new Set(),
    deletedRegionIds: new Set(),
    movedRegionIds: new Set(),
    createdSubRegionIds: new Set(),
    updatedSubRegionIds: new Set(),
    deletedSubRegionIds: new Set(),
    movedSubRegionIds: new Set(),
    createdNodeIds: new Set(),
    updatedNodeIds: new Set(),
    deletedNodeIds: new Set(),
    movedNodeIds: new Set(),
  }
}

export function clonePendingChanges(changes: PendingExplorerPlanChanges): PendingExplorerPlanChanges {
  return {
    updatedCourseTitle: changes.updatedCourseTitle,
    createdRegionIds: new Set(changes.createdRegionIds),
    updatedRegionIds: new Set(changes.updatedRegionIds),
    deletedRegionIds: new Set(changes.deletedRegionIds),
    movedRegionIds: new Set(changes.movedRegionIds),
    createdSubRegionIds: new Set(changes.createdSubRegionIds),
    updatedSubRegionIds: new Set(changes.updatedSubRegionIds),
    deletedSubRegionIds: new Set(changes.deletedSubRegionIds),
    movedSubRegionIds: new Set(changes.movedSubRegionIds),
    createdNodeIds: new Set(changes.createdNodeIds),
    updatedNodeIds: new Set(changes.updatedNodeIds),
    deletedNodeIds: new Set(changes.deletedNodeIds),
    movedNodeIds: new Set(changes.movedNodeIds),
  }
}

export function hasPendingChanges(changes: PendingExplorerPlanChanges) {
  return (
    changes.updatedCourseTitle ||
    changes.createdRegionIds.size > 0 ||
    changes.updatedRegionIds.size > 0 ||
    changes.deletedRegionIds.size > 0 ||
    changes.movedRegionIds.size > 0 ||
    changes.createdSubRegionIds.size > 0 ||
    changes.updatedSubRegionIds.size > 0 ||
    changes.deletedSubRegionIds.size > 0 ||
    changes.movedSubRegionIds.size > 0 ||
    changes.createdNodeIds.size > 0 ||
    changes.updatedNodeIds.size > 0 ||
    changes.deletedNodeIds.size > 0 ||
    changes.movedNodeIds.size > 0
  )
}

export type PendingItemState = 'none' | 'created' | 'deleted' | 'moved'

export function pendingItemState(
  createdIds: Set<string>,
  deletedIds: Set<string>,
  movedIds: Set<string>,
  id: string,
): PendingItemState {
  if (deletedIds.has(id)) return 'deleted'
  if (createdIds.has(id)) return 'created'
  if (movedIds.has(id)) return 'moved'
  return 'none'
}

export function markRegionUpdated(changes: PendingExplorerPlanChanges, regionId: string) {
  if (!changes.createdRegionIds.has(regionId)) changes.updatedRegionIds.add(regionId)
}

export function markRegionMoved(changes: PendingExplorerPlanChanges, regionId: string) {
  if (!changes.createdRegionIds.has(regionId)) changes.movedRegionIds.add(regionId)
}

export function markSubRegionUpdated(changes: PendingExplorerPlanChanges, subRegionId: string) {
  if (!changes.createdSubRegionIds.has(subRegionId)) changes.updatedSubRegionIds.add(subRegionId)
}

export function markSubRegionMoved(changes: PendingExplorerPlanChanges, subRegionId: string) {
  if (!changes.createdSubRegionIds.has(subRegionId)) changes.movedSubRegionIds.add(subRegionId)
}

export function markNodeUpdated(changes: PendingExplorerPlanChanges, nodeId: string) {
  if (!changes.createdNodeIds.has(nodeId)) changes.updatedNodeIds.add(nodeId)
}

export function markNodeMoved(changes: PendingExplorerPlanChanges, nodeId: string) {
  if (!changes.createdNodeIds.has(nodeId)) changes.movedNodeIds.add(nodeId)
}

export function clearRegionPending(changes: PendingExplorerPlanChanges, regionId: string) {
  changes.createdRegionIds.delete(regionId)
  changes.updatedRegionIds.delete(regionId)
  changes.deletedRegionIds.delete(regionId)
  changes.movedRegionIds.delete(regionId)
}

export function clearSubRegionPending(changes: PendingExplorerPlanChanges, subRegionId: string) {
  changes.createdSubRegionIds.delete(subRegionId)
  changes.updatedSubRegionIds.delete(subRegionId)
  changes.deletedSubRegionIds.delete(subRegionId)
  changes.movedSubRegionIds.delete(subRegionId)
}

export function clearNodePending(changes: PendingExplorerPlanChanges, nodeId: string) {
  changes.createdNodeIds.delete(nodeId)
  changes.updatedNodeIds.delete(nodeId)
  changes.deletedNodeIds.delete(nodeId)
  changes.movedNodeIds.delete(nodeId)
}

export function markSubRegionDeleted(changes: PendingExplorerPlanChanges, subRegionId: string) {
  changes.deletedSubRegionIds.add(subRegionId)
  changes.updatedSubRegionIds.delete(subRegionId)
  changes.movedSubRegionIds.delete(subRegionId)
}

export function markNodeDeleted(changes: PendingExplorerPlanChanges, nodeId: string) {
  changes.deletedNodeIds.add(nodeId)
  changes.updatedNodeIds.delete(nodeId)
  changes.movedNodeIds.delete(nodeId)
}

// ── 정렬 및 정규화 ─────────────────────────────────────────────────────────────

export function byOrderIndex<T extends { order_index: number; created_at?: string }>(a: T, b: T) {
  if (a.order_index !== b.order_index) return a.order_index - b.order_index
  return (a.created_at ?? '').localeCompare(b.created_at ?? '')
}

export function normalizeCourseAggregateOrder(course: CourseAggregate): CourseAggregate {
  return {
    ...course,
    regions: [...course.regions]
      .sort((a, b) => byOrderIndex(a.region, b.region))
      .map((regionAgg, regionIndex) => ({
        ...regionAgg,
        region: { ...regionAgg.region, order_index: regionIndex },
        subregions: [...regionAgg.subregions]
          .sort((a, b) => byOrderIndex(a.subregion, b.subregion))
          .map((subAgg, subIndex) => ({
            ...subAgg,
            subregion: { ...subAgg.subregion, order_index: subIndex },
            nodes: [...subAgg.nodes]
              .sort(byOrderIndex)
              .map((node, nodeIndex) => ({ ...node, order_index: nodeIndex })),
          })),
        nodes: [...regionAgg.nodes]
          .sort(byOrderIndex)
          .map((node, nodeIndex) => ({ ...node, order_index: nodeIndex })),
      })),
  }
}

// ── 로컬 객체 팩토리 ───────────────────────────────────────────────────────────

export function nowIso() {
  return new Date().toISOString()
}

export function temporaryId(kind: 'region' | 'subregion' | 'node') {
  return `local-${kind}-${Date.now()}-${Math.random().toString(36).slice(2)}`
}

export function makeLocalRegion(courseDraftId: string, name: string, orderIndex: number): ExplorerRegion {
  const timestamp = nowIso()
  return {
    id: temporaryId('region'),
    course_draft_id: courseDraftId,
    name,
    description: null,
    order_index: orderIndex,
    status: 'active',
    created_at: timestamp,
    updated_at: timestamp,
  }
}

export function makeLocalSubRegion(regionId: string, name: string, orderIndex: number): ExplorerSubRegion {
  const timestamp = nowIso()
  return {
    id: temporaryId('subregion'),
    region_id: regionId,
    name,
    description: null,
    order_index: orderIndex,
    status: 'active',
    created_at: timestamp,
    updated_at: timestamp,
  }
}

export function makeLocalResearchNode(parentKind: ParentKind, parentId: string, title: string, orderIndex: number): ExplorerNode {
  const timestamp = nowIso()
  return {
    id: temporaryId('node'),
    parent_kind: parentKind,
    parent_id: parentId,
    node_type: 'research',
    title,
    order_index: orderIndex,
    status: 'active',
    created_at: timestamp,
    updated_at: timestamp,
    block_count: 0,
  }
}

export function makeLocalExplorationNode(
  parentKind: ParentKind,
  parentId: string,
  title: string,
  source: ExplorationNodeSourceInput,
  orderIndex: number
): ExplorerNode {
  const timestamp = nowIso()
  return {
    id: temporaryId('node'),
    parent_kind: parentKind,
    parent_id: parentId,
    node_type: 'exploration',
    title,
    order_index: orderIndex,
    status: 'active',
    created_at: timestamp,
    updated_at: timestamp,
    source_type: source.sourceType,
    source_url: source.sourceUrl ?? null,
    content_id: source.contentId ?? null,
    block_count: 0,
  }
}

// ── 소스 타입 추론 ─────────────────────────────────────────────────────────────

export function inferSourceType(sourceUrl: string): SourceType {
  const normalized = sourceUrl.toLowerCase()
  return normalized.includes('youtube.com') || normalized.includes('youtu.be') ? 'youtube' : 'web'
}

export function normalizeExplorationNodeSource(source: string | ExplorationNodeSourceInput): ExplorationNodeSourceInput {
  if (typeof source === 'string') {
    return {
      sourceType: inferSourceType(source),
      sourceUrl: source,
      contentId: null,
    }
  }
  return source
}

export function shouldCheckExplorationNodeUrl(source: ExplorationNodeSourceInput) {
  return source.sourceType === 'web' || source.sourceType === 'youtube'
}

// ── 레거시 변환 ────────────────────────────────────────────────────────────────

export function legacyResearchType(templateType?: LegacyDraftPoint['template_type']): ResearchType {
  switch (templateType) {
    case 'concept_summary': return 'concept'
    case 'practice_strategy': return 'practice'
    case 'problem_solving': return 'problem'
    default: return 'free'
  }
}

export function legacyExplorationSourceType(point: LegacyDraftPoint): SourceType {
  if (point.external_url) return inferSourceType(point.external_url)
  return point.content_id ? 'internal' : 'web'
}

export function legacyPointToNode(parentKind: ParentKind, parentId: string, point: LegacyDraftPoint): ExplorerNode {
  const timestamp = point.updated_at ?? point.created_at ?? nowIso()
  const isExploration = point.point_type === 'exploration'
  return {
    id: `legacy-node-${point.id}`,
    parent_kind: parentKind,
    parent_id: parentId,
    node_type: point.point_type,
    draft_point_id: point.id,
    title: point.title,
    order_index: point.order_index,
    status: 'active',
    created_at: point.created_at ?? timestamp,
    updated_at: timestamp,
    source_type: isExploration ? legacyExplorationSourceType(point) : undefined,
    source_url: isExploration ? point.external_url ?? null : undefined,
    content_id: isExploration ? point.content_id ?? null : undefined,
    summary: point.description ?? null,
    research_type: isExploration ? undefined : legacyResearchType(point.template_type),
    layout_type: isExploration ? undefined : 'basic',
    block_count: 0,
  }
}

export function buildLegacyCourseAggregate(courseDraftId: string, legacyDraft: LegacyDraftAggregate): CourseAggregate {
  return {
    course_draft_id: courseDraftId,
    title: legacyDraft.draft.title,
    status: legacyDraft.draft.status ?? 'draft',
    updated_at: legacyDraft.draft.updated_at,
    regions: legacyDraft.lessons.map((tree, regionIndex) => {
      const regionId = `legacy-region-${tree.lesson.id}`
      return {
        region: {
          id: regionId,
          course_draft_id: courseDraftId,
          name: tree.lesson.title,
          description: tree.lesson.summary ?? null,
          order_index: tree.lesson.order_index ?? regionIndex,
          status: 'active',
          created_at: tree.lesson.created_at ?? legacyDraft.draft.updated_at,
          updated_at: tree.lesson.updated_at ?? legacyDraft.draft.updated_at,
        },
        nodes: (tree.points ?? []).map(({ point }) => legacyPointToNode('region', regionId, point)),
        subregions: (tree.sub_lessons ?? []).map((subTree, subIndex) => {
          const subRegionId = `legacy-subregion-${subTree.lesson.id}`
          return {
            subregion: {
              id: subRegionId,
              region_id: regionId,
              name: subTree.lesson.title,
              description: subTree.lesson.summary ?? null,
              order_index: subTree.lesson.order_index ?? subIndex,
              status: 'active',
              created_at: subTree.lesson.created_at ?? legacyDraft.draft.updated_at,
              updated_at: subTree.lesson.updated_at ?? legacyDraft.draft.updated_at,
            },
            nodes: (subTree.points ?? []).map(({ point }) => legacyPointToNode('subregion', subRegionId, point)),
          }
        }),
      }
    }),
  }
}

export function shouldUseLegacyDraftAggregate(explorerAggregate: CourseAggregate, legacyDraft: LegacyDraftAggregate | null | undefined) {
  if (!legacyDraft || legacyDraft.lessons.length === 0) return false
  return explorerAggregate.regions.length === 0
}

export function pendingChangesForLegacyAggregate(
  legacyAggregate: CourseAggregate,
  previousExplorerAggregate: CourseAggregate
): PendingExplorerPlanChanges {
  const changes = emptyPendingChanges()
  for (const regionAgg of previousExplorerAggregate.regions) {
    changes.deletedRegionIds.add(regionAgg.region.id)
  }
  for (const regionAgg of legacyAggregate.regions) {
    changes.createdRegionIds.add(regionAgg.region.id)
    for (const node of regionAgg.nodes) {
      changes.createdNodeIds.add(node.id)
    }
    for (const subAgg of regionAgg.subregions) {
      changes.createdSubRegionIds.add(subAgg.subregion.id)
      for (const node of subAgg.nodes) {
        changes.createdNodeIds.add(node.id)
      }
    }
  }
  return changes
}
