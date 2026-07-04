'use client'

import { useRouter } from 'next/navigation'
import { useCallback, useEffect, useMemo, useRef, useState, type CSSProperties } from 'react'

import { useExplorerPlan } from '@/components/explorer-plan/useExplorerPlan'
import type { DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'
import LumiModalShell, {
  lumiModalPrimaryButtonStyle,
  lumiModalSecondaryButtonStyle,
} from '@/components/common/LumiModalShell'
import type { DraftEditorState } from '../useDraftEditor'
import type { DraftLessonTree, DraftPointAggregate } from '../types'
import { PlanningMapPanel } from './planning/PlanningMapPanel'
import { PlanningRightPanel } from './planning/PlanningRightPanel'
import type { PlanningLumiGuideFocus } from './planning/PlanningLumiGuide'
import type { PlanningMapLevel } from './planning/map/usePlanningMapNavigation'
import type { ExplorerNode } from '@/components/explorer-plan/explorerPlanTypes'

// ── 반응형 분기점 ─────────────────────────────────────────────────────────────
// 이 이상이면 지도 자동 펼침 / 미만이면 디바이스만 표시
const MAP_BREAKPOINT = 1024

// ── 캔버스 기준 크기 — 접힘 이미지 기준으로 통일 ────────────────────────────
const BASE_STAGE = { width: 1395, height: 757 }
// 좁은 화면에서는 펼침 이미지의 우측 작은 디바이스 영역만 보여준다.
const SMALL_DEVICE_CROP = { x: 960, y: 48, width: 420, height: 650 }
const OPEN_BACKGROUND_SRC = '/images/explorer/Explorer_Diary_Design.webp?v=20260417-open-1395'
const CLOSED_BACKGROUND_SRC = '/images/explorer/Explorer_Diary_Design_wide.webp'

// ── 스타일 ────────────────────────────────────────────────────────────────────
const outerSectionStyle: CSSProperties = {
  position: 'relative',
  width: '100%',
}

const scaledCanvasCenterStyle: CSSProperties = {
  position: 'absolute',
  left: '50%',
  top: '50%',
  transformOrigin: 'center center',
  willChange: 'transform',
}

const frameBackdropStyle: CSSProperties = {
  position: 'relative',
  width: '100%',
  aspectRatio: `${BASE_STAGE.width} / ${BASE_STAGE.height}`,
  overflow: 'hidden',
}

const stageBaseStyle: CSSProperties = {
  position: 'relative',
  borderRadius: 28,
  overflow: 'hidden',
  border: '1px solid rgba(255,255,255,0.10)',
  boxShadow: '0 28px 70px rgba(0,0,0,0.30)',
  background: '#1b120c',
}

const backgroundWrapStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  zIndex: 0,
}

const overlayStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  zIndex: 1,
  pointerEvents: 'none',
  background: 'linear-gradient(180deg, rgba(0,0,0,0.03) 0%, rgba(0,0,0,0.10) 100%)',
}

const inactiveOverlayStyle: CSSProperties = {
  ...overlayStyle,
  background: 'linear-gradient(180deg, rgba(15, 23, 42, 0.30) 0%, rgba(15, 23, 42, 0.42) 100%)',
  backdropFilter: 'grayscale(0.45) saturate(0.72)',
}

const inactiveBannerStyle: CSSProperties = {
  position: 'absolute',
  left: '50%',
  top: '4.1%',
  transform: 'translateX(-50%)',
  zIndex: 7,
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: 36,
  padding: '0 18px',
  borderRadius: 999,
  border: '1px solid rgba(203, 213, 225, 0.48)',
  background: 'linear-gradient(180deg, rgba(51, 65, 85, 0.92), rgba(30, 41, 59, 0.94))',
  color: '#E2E8F0',
  fontSize: 13,
  fontWeight: 900,
  boxShadow: '0 12px 28px rgba(2, 6, 23, 0.30)',
  letterSpacing: '0.02em',
}

const completionGuideButtonStyle: CSSProperties = {
  border: '1px solid rgba(140, 96, 22, 0.42)',
  borderRadius: 999,
  padding: '10px 18px',
  background: 'linear-gradient(180deg, rgba(246, 205, 103, 0.94), rgba(214, 153, 34, 0.96))',
  color: '#4A2C00',
  fontSize: 14,
  fontWeight: 900,
  cursor: 'pointer',
  boxShadow: '0 10px 28px rgba(138, 95, 16, 0.22)',
}

const secondaryGuideButtonStyle: CSSProperties = {
  border: '1px solid rgba(109, 73, 22, 0.28)',
  borderRadius: 999,
  padding: '10px 18px',
  background: 'rgba(255, 248, 236, 0.9)',
  color: '#6A4720',
  fontSize: 14,
  fontWeight: 800,
  cursor: 'pointer',
}

// 책갈피 공통 베이스 — 지도 세로 중간(46.6%)에서 내려오는 북마크 형태
// top = map_top(8.5%) + map_height(76.2%) / 2 = 46.6%
const bookmarkBase: CSSProperties = {
  position: 'absolute',
  zIndex: 4,
  top: '46.6%',
  background: 'rgba(160, 30, 20, 0.95)',
  color: 'rgba(255, 230, 220, 1)',
  fontSize: 11,
  fontWeight: 800,
  padding: '14px 8px 14px',
  cursor: 'pointer',
  border: '1.5px solid rgba(220, 80, 60, 0.7)',
  boxShadow: '0 4px 16px rgba(0,0,0,0.45), inset 0 1px 0 rgba(255,160,140,0.15)',
  letterSpacing: '0.05em',
  writingMode: 'vertical-rl',
  textOrientation: 'mixed',
  whiteSpace: 'nowrap',
  borderRadius: 0,
}

const collapseButtonStyle: CSSProperties = {
  ...bookmarkBase,
  left: '73.35%',
  transform: 'translateX(-50%)',
}

const expandButtonStyle: CSSProperties = {
  ...bookmarkBase,
  left: '14.4%',
}

const backgroundImageStyle: CSSProperties = {
  display: 'block',
  width: '100%',
  height: '100%',
  objectFit: 'fill',
}

// ── 반응형 지도 상태 ──────────────────────────────────────────────────────────
function useResponsiveMapExpanded() {
  const [mapExpanded, setMapExpanded] = useState(() =>
    typeof window !== 'undefined' ? window.innerWidth >= MAP_BREAKPOINT : true
  )
  const [isNarrowViewport, setIsNarrowViewport] = useState(() =>
    typeof window !== 'undefined' ? window.innerWidth < MAP_BREAKPOINT : false
  )

  useEffect(() => {
    const mq = window.matchMedia(`(min-width: ${MAP_BREAKPOINT}px)`)
    // 브레이크포인트 교차 시 자동 전환
    const handler = (e: MediaQueryListEvent) => {
      setMapExpanded(e.matches)
      setIsNarrowViewport(!e.matches)
    }
    mq.addEventListener('change', handler)
    // 마운트 시 동기화
    setMapExpanded(mq.matches)
    setIsNarrowViewport(!mq.matches)
    return () => mq.removeEventListener('change', handler)
  }, [])

  return { mapExpanded, setMapExpanded, isNarrowViewport }
}

// ── 반응형 스케일 ─────────────────────────────────────────────────────────────
function useResponsiveCanvasScale(
  baseWidth: number,
  baseHeight: number
): { frameRef: React.RefObject<HTMLDivElement | null>; scale: number } {
  const frameRef = useRef<HTMLDivElement>(null)
  const [scale, setScale] = useState(1)

  useEffect(() => {
    const element = frameRef.current
    if (!element) return

    const updateScale = (width: number, height: number) => {
      if (width <= 0 || height <= 0) return
      setScale(Math.min(width / baseWidth, height / baseHeight))
    }

    updateScale(element.clientWidth, element.clientHeight)

    const observer = new ResizeObserver(entries => {
      const entry = entries[0]
      if (!entry) return
      const { width, height } = entry.contentRect
      updateScale(width, height)
    })

    observer.observe(element)
    return () => observer.disconnect()
  }, [baseWidth, baseHeight])

  return { frameRef, scale }
}

// ── 컴포넌트 ──────────────────────────────────────────────────────────────────
interface PlanningSectionProps {
  editor: DraftEditorState
  draftId: string
  onDirtyChange?: (isDirty: boolean) => void
  onGuideFocusChange?: (focus: PlanningLumiGuideFocus) => void
  onAfterPlanSave?: () => void
  mode?: 'planning' | 'journal'
  onJournalSelectionChange?: (selection: JournalSelectionMeta | null) => void
  onJournalStatsChange?: (stats: JournalStructureStats | null) => void
  onSectionChange?: (section: 'planning' | 'journal') => void
  confirmTabChange?: boolean
  copy: DashboardCourseDraftCopy['planning']
  treeCopy: DashboardCourseDraftCopy['tree']
}

export interface JournalSelectionMeta {
  regionTitle: string | null
  subRegionTitle: string | null
  nodeTitle: string | null
  nodeType: 'exploration' | 'research' | null
  pointRoute: string | null
  canOpenPoint: boolean
}

export interface JournalStructureStats {
  lessonCount: number
  completedLessonCount: number
  learningLessonCount: number
  explorationCount: number
  researchCount: number
  totalCount: number
  completedExplorationCount: number
  completedResearchCount: number
  completedTotalCount: number
}

type SelectedNodeContext = {
  regionTitle: string
  subRegionTitle: string | null
  node: ExplorerNode
}

type JournalPointModalState = {
  regionTitle: string
  subRegionTitle: string | null
  node: ExplorerNode
  pointRoute: string
}

type TransitionModalState =
  {
    kind: 'tab_change'
    target: 'planning' | 'journal'
  }

function normalizeTitle(value?: string | null) {
  return value?.trim().toLowerCase() ?? ''
}

function findSelectedNodeContext(
  courseAggregate: ReturnType<typeof useExplorerPlan>['courseAggregate'],
  selectedNodeId: string | null,
): SelectedNodeContext | null {
  if (!courseAggregate || !selectedNodeId) return null

  for (const regionAgg of courseAggregate.regions) {
    for (const node of regionAgg.nodes) {
      if (node.id === selectedNodeId) {
        return {
          regionTitle: regionAgg.region.name,
          subRegionTitle: null,
          node,
        }
      }
    }
    for (const subAgg of regionAgg.subregions) {
      for (const node of subAgg.nodes) {
        if (node.id === selectedNodeId) {
          return {
            regionTitle: regionAgg.region.name,
            subRegionTitle: subAgg.subregion.name,
            node,
          }
        }
      }
    }
  }

  return null
}

function findPointByNodeContext(
  lessons: DraftLessonTree[],
  context: SelectedNodeContext,
): DraftPointAggregate | null {
  const region = lessons.find((lesson) => normalizeTitle(lesson.lesson.title) === normalizeTitle(context.regionTitle))
  if (!region) return null

  const pointMatchesNode = (pointAgg: DraftPointAggregate) =>
    normalizeTitle(pointAgg.point.title) === normalizeTitle(context.node.title) &&
    pointAgg.point.point_type === context.node.node_type

  if (context.subRegionTitle) {
    const subLesson = (region.sub_lessons ?? []).find(
      (sub) => normalizeTitle(sub.lesson.title) === normalizeTitle(context.subRegionTitle),
    )
    if (!subLesson) return null
    return (subLesson.points ?? []).find(pointMatchesNode) ?? null
  }

  return (region.points ?? []).find(pointMatchesNode) ?? null
}

function flattenDraftPoints(lessons: DraftLessonTree[]): DraftPointAggregate[] {
  const points: DraftPointAggregate[] = []

  const visit = (lessonTree: DraftLessonTree) => {
    points.push(...(lessonTree.points ?? []))
    ;(lessonTree.sub_lessons ?? []).forEach(visit)
  }

  lessons.forEach(visit)
  return points
}

function findDraftPointForExplorerNode(
  points: DraftPointAggregate[],
  node: ExplorerNode,
): DraftPointAggregate | null {
  const linkedPointId = node.draft_point_id?.trim()
  if (linkedPointId) {
    const linkedPoint = points.find((pointAgg) => pointAgg.point.id === linkedPointId)
    if (linkedPoint) return linkedPoint
  }

  return (
    points.find(
      (pointAgg) =>
        normalizeTitle(pointAgg.point.title) === normalizeTitle(node.title) &&
        pointAgg.point.point_type === node.node_type,
    ) ?? null
  )
}

function resolveJournalPointRoute(
  draft: DraftEditorState['draft'],
  selectedNodeId: string | null,
  selectedContext: SelectedNodeContext | null,
): string | null {
  const confirmedCourseId = draft?.draft.confirmed_course_id?.trim()
  if (!confirmedCourseId || !selectedNodeId || !selectedContext) return null
  const routeKind = draft?.draft.status === 'archived' ? 'shared' : 'learning'

  const linkedPointId = selectedContext.node.draft_point_id?.trim()
  if (linkedPointId) {
    return `/dashboard/planets/${routeKind}/${confirmedCourseId}/points/${linkedPointId}`
  }

  if (selectedNodeId.startsWith('legacy-node-')) {
    const pointId = selectedNodeId.slice('legacy-node-'.length)
    return pointId ? `/dashboard/planets/${routeKind}/${confirmedCourseId}/points/${pointId}` : null
  }

  const matchedPoint = findPointByNodeContext(draft?.lessons ?? [], selectedContext)
  if (!matchedPoint) return null

  return `/dashboard/planets/${routeKind}/${confirmedCourseId}/points/${matchedPoint.point.id}`
}

function PlanningBackground({
  src,
  alt,
}: {
  src: string
  alt: string
}) {
  return (
    <img
      src={src}
      alt={alt}
      style={backgroundImageStyle}
    />
  )
}

export function PlanningSection({
  editor,
  draftId,
  onDirtyChange,
  onGuideFocusChange,
  onAfterPlanSave,
  mode = 'planning',
  onJournalSelectionChange,
  onJournalStatsChange,
  onSectionChange,
  confirmTabChange = true,
  copy,
  treeCopy,
}: PlanningSectionProps) {
  const router = useRouter()
  const explorerPlan = useExplorerPlan(draftId, editor.draft)
  const isInactiveCourse = Boolean(editor.draft?.draft.is_inactive)
  const isJournalMode = mode === 'journal'
  const isJournalLimitedEdit = isJournalMode
  const isStructureReadOnly = false
  const { mapExpanded, setMapExpanded, isNarrowViewport } = useResponsiveMapExpanded()
  const [mapLevel, setMapLevel] = useState<PlanningMapLevel>('course')
  const [mapSelectedRegionId, setMapSelectedRegionId] = useState<string | null>(null)
  const [mapSelectedSubRegionId, setMapSelectedSubRegionId] = useState<string | null>(null)
  const [journalPointModal, setJournalPointModal] = useState<JournalPointModalState | null>(null)
  const [transitionModal, setTransitionModal] = useState<TransitionModalState | null>(null)

  const regions = explorerPlan.courseAggregate?.regions ?? []

  useEffect(() => {
    onDirtyChange?.(isStructureReadOnly ? false : explorerPlan.isDirty)
  }, [explorerPlan.isDirty, isStructureReadOnly, onDirtyChange])

  const selectedNodeContext = useMemo(
    () => findSelectedNodeContext(explorerPlan.courseAggregate, explorerPlan.selectedNodeId),
    [explorerPlan.courseAggregate, explorerPlan.selectedNodeId],
  )

  const journalSelectionMeta = useMemo<JournalSelectionMeta | null>(() => {
    if (!isJournalMode) return null

    const pointRoute = editor.isLearningStarted
      ? resolveJournalPointRoute(editor.draft, explorerPlan.selectedNodeId, selectedNodeContext)
      : null
    return {
      regionTitle: selectedNodeContext?.regionTitle ?? null,
      subRegionTitle: selectedNodeContext?.subRegionTitle ?? null,
      nodeTitle: selectedNodeContext?.node.title ?? null,
      nodeType: selectedNodeContext?.node.node_type ?? null,
      pointRoute,
      canOpenPoint: Boolean(pointRoute),
    }
  }, [editor.draft, editor.isLearningStarted, explorerPlan.selectedNodeId, isJournalMode, selectedNodeContext])

  const journalStructureStats = useMemo<JournalStructureStats | null>(() => {
    if (!isJournalMode || !explorerPlan.courseAggregate) return null

    let explorationCount = 0
    let researchCount = 0
    let completedExplorationCount = 0
    let completedResearchCount = 0
    let lessonCount = 0
    let completedLessonCount = 0
    let learningLessonCount = 0
    const draftPoints = flattenDraftPoints(editor.draft?.lessons ?? [])

    const getNodeCompletion = (node: ExplorerNode) => {
      const matchedPoint = findDraftPointForExplorerNode(draftPoints, node)
      if (node.node_type === 'exploration') {
        explorationCount += 1
        if (matchedPoint?.point.status === 'completed') {
          completedExplorationCount += 1
          return 'completed'
        }
      }
      if (node.node_type === 'research') {
        researchCount += 1
        if (matchedPoint?.point.status === 'completed') {
          completedResearchCount += 1
          return 'completed'
        }
      }
      return matchedPoint?.point.status === 'learning' ? 'learning' : 'ready'
    }

    const countRegionLesson = (nodes: ExplorerNode[]) => {
      const activeNodes = nodes.filter((node) => node.status === 'active')
      lessonCount += 1
      if (activeNodes.length === 0) return

      const nodeStatuses = activeNodes.map(getNodeCompletion)
      if (nodeStatuses.every((status) => status === 'completed')) {
        completedLessonCount += 1
      } else if (nodeStatuses.some((status) => status === 'completed' || status === 'learning')) {
        learningLessonCount += 1
      }
    }

    explorerPlan.courseAggregate.regions.forEach((region) => {
      if (region.region.status !== 'active') return
      const activeRegionNodes = [
        ...region.nodes,
        ...region.subregions
          .filter((subregion) => subregion.subregion.status === 'active')
          .flatMap((subregion) => subregion.nodes),
      ]
      countRegionLesson(activeRegionNodes)
    })

    return {
      lessonCount,
      completedLessonCount,
      learningLessonCount,
      explorationCount,
      researchCount,
      totalCount: explorationCount + researchCount,
      completedExplorationCount,
      completedResearchCount,
      completedTotalCount: completedExplorationCount + completedResearchCount,
    }
  }, [editor.draft?.lessons, explorerPlan.courseAggregate, isJournalMode])

  useEffect(() => {
    if (!isJournalMode) return
    onJournalSelectionChange?.(journalSelectionMeta)
  }, [isJournalMode, journalSelectionMeta, onJournalSelectionChange])

  useEffect(() => {
    if (!isJournalMode) return
    onJournalStatsChange?.(journalStructureStats)
  }, [isJournalMode, journalStructureStats, onJournalStatsChange])

  const handleJournalNodeActivate = useCallback((nodeId: string) => {
    if (!isJournalMode || !editor.isLearningStarted) return
    const selectedContext = findSelectedNodeContext(explorerPlan.courseAggregate, nodeId)
    if (!selectedContext) return
    const pointRoute = resolveJournalPointRoute(editor.draft, nodeId, selectedContext)
    if (!pointRoute) return
    setJournalPointModal({
      regionTitle: selectedContext.regionTitle,
      subRegionTitle: selectedContext.subRegionTitle,
      node: selectedContext.node,
      pointRoute,
    })
  }, [editor.draft, editor.isLearningStarted, explorerPlan.courseAggregate, isJournalMode])

  const handleRequestTabSelect = useCallback((tab: 'planning' | 'journal') => {
    if (!onSectionChange) return
    if (tab === mode) return
    if (!confirmTabChange) {
      onSectionChange(tab)
      return
    }
    setTransitionModal({ kind: 'tab_change', target: tab })
  }, [confirmTabChange, mode, onSectionChange])

  const handleConfirmTransition = useCallback(() => {
    if (!transitionModal) return

    onSectionChange?.(transitionModal.target)
    setTransitionModal(null)
  }, [onSectionChange, transitionModal])

  useEffect(() => {
    if (explorerPlan.selectedNodeId) {
      for (const regionAgg of regions) {
        if (regionAgg.nodes.some(node => node.id === explorerPlan.selectedNodeId)) {
          setMapLevel('region')
          setMapSelectedRegionId(regionAgg.region.id)
          setMapSelectedSubRegionId(null)
          return
        }
        for (const subAgg of regionAgg.subregions) {
          if (subAgg.nodes.some(node => node.id === explorerPlan.selectedNodeId)) {
            setMapLevel('subregion')
            setMapSelectedRegionId(regionAgg.region.id)
            setMapSelectedSubRegionId(subAgg.subregion.id)
            return
          }
        }
      }
    }

    if (explorerPlan.selectedSubRegionId) {
      const parentRegion = regions.find(regionAgg =>
        regionAgg.subregions.some(subAgg => subAgg.subregion.id === explorerPlan.selectedSubRegionId)
      )
      if (parentRegion) {
        setMapLevel('subregion')
        setMapSelectedRegionId(parentRegion.region.id)
        setMapSelectedSubRegionId(explorerPlan.selectedSubRegionId)
      }
      return
    }

    if (explorerPlan.selectedRegionId) {
      setMapLevel('region')
      setMapSelectedRegionId(explorerPlan.selectedRegionId)
      setMapSelectedSubRegionId(null)
      return
    }

    setMapLevel('course')
    setMapSelectedRegionId(null)
    setMapSelectedSubRegionId(null)
  }, [
    explorerPlan.selectedNodeId,
    explorerPlan.selectedRegionId,
    explorerPlan.selectedSubRegionId,
    regions,
  ])

  const toggleMap = useCallback(() => setMapExpanded(v => !v), [setMapExpanded])

  const viewportBase = isNarrowViewport ? SMALL_DEVICE_CROP : BASE_STAGE
  const panelUsesSmallDevice = mapExpanded || isNarrowViewport
  const backgroundSrc = panelUsesSmallDevice ? OPEN_BACKGROUND_SRC : CLOSED_BACKGROUND_SRC
  const backgroundAlt = isNarrowViewport
    ? isJournalMode ? copy.backgroundAlt.smallJournal : copy.backgroundAlt.smallPlanning
    : mapExpanded
      ? isJournalMode ? copy.backgroundAlt.expandedJournal : copy.backgroundAlt.expandedPlanning
      : isJournalMode ? copy.backgroundAlt.collapsedJournal : copy.backgroundAlt.collapsedPlanning

  const { frameRef, scale } = useResponsiveCanvasScale(viewportBase.width, viewportBase.height)

  return (
    <section
      style={outerSectionStyle}
    >
      <div
        ref={frameRef}
        style={{
          ...frameBackdropStyle,
          aspectRatio: `${viewportBase.width} / ${viewportBase.height}`,
        }}
      >
        <div
          style={isNarrowViewport
            ? {
                position: 'absolute',
                left: -SMALL_DEVICE_CROP.x * scale,
                top: -SMALL_DEVICE_CROP.y * scale,
                width: BASE_STAGE.width,
                height: BASE_STAGE.height,
                transform: `scale(${scale})`,
                transformOrigin: 'top left',
                willChange: 'transform',
              }
            : {
                ...scaledCanvasCenterStyle,
                width: BASE_STAGE.width,
                height: BASE_STAGE.height,
                transform: `translate(-50%, -50%) scale(${scale})`,
              }}
        >
          <div style={{ ...stageBaseStyle, width: BASE_STAGE.width, height: BASE_STAGE.height }}>
            {/* 배경 이미지 — 모드에 따라 교차 */}
            {/* TODO: 스팀펑크 스타일 이미지로 교체 필요 */}
            <div style={backgroundWrapStyle}>
              <PlanningBackground
                src={backgroundSrc}
                alt={backgroundAlt}
              />
            </div>

            <div style={isInactiveCourse ? inactiveOverlayStyle : overlayStyle} />

            {isInactiveCourse ? (
              <div style={inactiveBannerStyle}>
                {copy.inactiveBanner}
              </div>
            ) : null}

            {/* 지도 토글 버튼 */}
            {!isNarrowViewport ? (
              mapExpanded ? (
                <button
                  type="button"
                  style={collapseButtonStyle}
                  onClick={toggleMap}
                  aria-label={copy.collapseMap}
                >
                  {copy.collapseMap}
                </button>
              ) : (
                <button
                  type="button"
                  style={expandButtonStyle}
                  onClick={toggleMap}
                  aria-label={copy.expandMap}
                >
                  {copy.expandMap}
                </button>
              )
            ) : null}

            {/* 패널 — stage 기준 절대 위치 */}
            <PlanningMapPanel
              explorerPlan={explorerPlan}
              mapExpanded={mapExpanded}
              mapLevel={mapLevel}
              selectedRegionId={mapSelectedRegionId}
              selectedSubRegionId={mapSelectedSubRegionId}
              onMapLevelChange={setMapLevel}
              onSelectedRegionIdChange={setMapSelectedRegionId}
              onSelectedSubRegionIdChange={setMapSelectedSubRegionId}
              onSelectionFocusChange={isInactiveCourse ? undefined : onGuideFocusChange}
              onNodeActivate={handleJournalNodeActivate}
              copy={copy.map}
            />
            <PlanningRightPanel
              explorerPlan={explorerPlan}
              mapExpanded={panelUsesSmallDevice}
              isNarrowViewport={isNarrowViewport}
              readOnly={isStructureReadOnly}
              journalLimitedEdit={isJournalLimitedEdit}
              activeTab={mode}
              isJournalAvailable={editor.isPlanCompleted}
              onTabSelect={handleRequestTabSelect}
              onNodeActivate={handleJournalNodeActivate}
              onSelectionFocusChange={onGuideFocusChange}
              onAfterSave={onAfterPlanSave}
              copy={copy.rightPanel}
              treeCopy={treeCopy}
            />
            {isJournalMode && journalPointModal && (
              <LumiModalShell
                ariaLabel={copy.pointModal.aria}
                eyebrow={journalPointModal.node.node_type === 'research' ? 'Research Point' : 'Exploration Point'}
                lumiState="curious"
                title={journalPointModal.node.title}
                message={
                  <>
                    <div>{copy.pointModal.region(journalPointModal.regionTitle)}</div>
                    {journalPointModal.subRegionTitle ? <div>{copy.pointModal.subRegion(journalPointModal.subRegionTitle)}</div> : null}
                    <div>{copy.pointModal.type(journalPointModal.node.node_type)}</div>
                    <div>{copy.pointModal.question}</div>
                  </>
                }
                onClose={() => setJournalPointModal(null)}
                actions={
                  <>
                    <button
                      type="button"
                      style={lumiModalSecondaryButtonStyle}
                      onClick={() => setJournalPointModal(null)}
                    >
                      {copy.pointModal.cancel}
                    </button>
                    <button
                      type="button"
                      style={lumiModalPrimaryButtonStyle}
                      onClick={() => {
                        const pointRoute = journalPointModal.pointRoute
                        setJournalPointModal(null)
                        router.push(pointRoute)
                      }}
                    >
                      {copy.pointModal.open}
                    </button>
                  </>
                }
              />
            )}
            {transitionModal && (
              <LumiModalShell
                ariaLabel={copy.transition.aria}
                eyebrow="Lumi Confirm"
                lumiState="curious"
                tone="warm"
                title={
                  transitionModal.target === 'planning'
                    ? copy.transition.toPlanningTitle
                    : copy.transition.toJournalTitle
                }
                message={
                  transitionModal.target === 'planning'
                    ? copy.transition.toPlanningMessage
                    : copy.transition.toJournalMessage
                }
                onClose={() => setTransitionModal(null)}
                actions={
                  <>
                    <button
                      type="button"
                      style={lumiModalSecondaryButtonStyle}
                      onClick={() => setTransitionModal(null)}
                    >
                      {copy.transition.cancel}
                    </button>
                    <button
                      type="button"
                      style={lumiModalPrimaryButtonStyle}
                      onClick={handleConfirmTransition}
                    >
                      {copy.transition.confirm}
                    </button>
                  </>
                }
              />
            )}
          </div>
        </div>
      </div>
    </section>
  )
}
