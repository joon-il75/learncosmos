'use client'

import { useCallback, useEffect, useMemo, useRef, useState, type CSSProperties } from 'react'
import { createPortal } from 'react-dom'
import { useRouter, useSearchParams } from 'next/navigation'
import LumiAvatar from '@/components/lumi/LumiAvatar'
import { useExplorerPlan } from '@/components/explorer-plan/useExplorerPlan'
import { useExplorerEditForm } from '@/components/explorer-plan/tree/useExplorerEditForm'
import {
  useResponsiveMapExpanded,
  useResponsiveCanvasScale,
  findDiaryNodeById,
  mergeDiaryCourseWithExplorer,
} from './diaryStageUtils'
import type { PlanningMapLevel } from '@/app/dashboard/course-drafts/[id]/sections/planning/map/usePlanningMapNavigation'
import {
  cancelButtonStyle,
  saveButtonStyle,
  treeButtonDisabledStyle,
} from '@/components/explorer-plan/tree/explorerTreeStyles'
import type { PlanetDiaryCopy } from '@/lib/i18n/pages/planetDiary'
import type { DiaryCourseAggregate } from './planetDiaryTypes'
import { PlanetRecordPage } from '../records/PlanetRecordPage'
import { PlanetResultPage } from '../results/PlanetResultPage'
import { DiaryTreePanel } from './DiaryTreePanel'
import { DiaryGuidePanel } from './DiaryGuidePanel'
import { DiaryMapPanel } from './DiaryMapPanel'
import { DiaryStageModals, type PlannedBookmark } from './DiaryStageModals'

type PlanetRouteKind = 'learning' | 'shared'
type DiaryStageTab = 'journal' | 'records' | 'results'
type BookmarkState = 'active' | 'inactive' | 'locked'
type FirstElementOnboardingStep = 'lesson' | 'add_element' | 'modal' | 'save' | 'open_point'

const BASE_STAGE = { width: 1395, height: 757 }
const SMALL_DEVICE_CROP = { x: 960, y: 48, width: 420, height: 650 }
const OPEN_BACKGROUND_SRC = '/images/explorer/Explorer_Diary_Design.webp?v=20260417-open-1395'
const CLOSED_BACKGROUND_SRC = '/images/explorer/Explorer_Diary_Design_wide.webp'

const OPEN_RIGHT_SHELL: CSSProperties = {
  position: 'absolute',
  left: '75.8%',
  top: '13.6%',
  width: '17.8%',
  height: '57.8%',
  zIndex: 3,
  minHeight: 0,
  overflow: 'hidden',
  padding: '8px 6px',
  boxSizing: 'border-box',
}

const CLOSED_RIGHT_SHELL: CSSProperties = {
  position: 'absolute',
  left: '20.4%',
  top: '16.7%',
  width: '70.2%',
  height: '59.4%',
  zIndex: 3,
  minHeight: 0,
  overflow: 'hidden',
  padding: '14px 16px',
  boxSizing: 'border-box',
}

const outerSectionStyle: CSSProperties = {
  position: 'relative',
  width: '100%',
  paddingBottom: 170,
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

const shellPanelStyle: CSSProperties = {
  display: 'grid',
  gridTemplateRows: 'auto 1fr',
  gap: 8,
  height: '100%',
  minHeight: 0,
  borderRadius: 16,
  background: 'rgba(255, 249, 240, 0.74)',
  border: '1px solid rgba(118, 88, 54, 0.18)',
  boxShadow: 'inset 0 1px 0 rgba(255,255,255,0.34)',
  overflow: 'hidden',
}

const panelHeaderStyle: CSSProperties = {
  display: 'grid',
  gap: 4,
  padding: '12px 12px 0',
}

const panelEyebrowStyle: CSSProperties = {
  fontSize: 10,
  fontWeight: 800,
  letterSpacing: '0.08em',
  textTransform: 'uppercase',
  color: '#7a5a24',
}

const panelTitleRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 8,
}

const panelTitleStyle: CSSProperties = {
  margin: 0,
  fontSize: 16,
  fontWeight: 900,
  color: '#1f160a',
}

const panelBadgeStyle: CSSProperties = {
  padding: '4px 9px',
  borderRadius: 999,
  background: 'rgba(122, 90, 36, 0.10)',
  color: '#6d4b18',
  fontSize: 10,
  fontWeight: 800,
}

const bookmarkRailStyle = (mapExpanded: boolean): CSSProperties => ({
  position: 'absolute',
  left: mapExpanded ? '92.9%' : '91.2%',
  top: '12.5%',
  display: 'flex',
  flexDirection: 'column',
  gap: 5,
  zIndex: 6,
})

const bookmarkTabStyle = (state: BookmarkState): CSSProperties => {
  const active = state === 'active'
  const locked = state === 'locked'
  return {
    width: 40,
    minHeight: 72,
    padding: '12px 0',
    background: active
      ? 'linear-gradient(135deg, #c9920e 0%, #a06c08 60%, #7a5005 100%)'
      : 'linear-gradient(135deg, rgba(92, 60, 19, 0.90) 0%, rgba(47, 30, 10, 0.94) 100%)',
    color: active ? '#fffaf0' : locked ? 'rgba(248, 223, 168, 0.86)' : 'rgba(255, 227, 170, 0.92)',
    fontSize: 12,
    fontWeight: 900,
    letterSpacing: '0.12em',
    lineHeight: 1.05,
    writingMode: 'vertical-rl',
    textOrientation: 'mixed',
    textAlign: 'center',
    textShadow: active
      ? '0 1px 3px rgba(62, 31, 0, 0.55)'
      : '0 1px 2px rgba(0, 0, 0, 0.48)',
    borderRadius: '0 11px 11px 0',
    cursor: active ? 'default' : locked ? 'not-allowed' : 'pointer',
    border: `1px solid ${active ? '#e8b020' : locked ? 'rgba(176, 126, 45, 0.50)' : 'rgba(168, 111, 37, 0.42)'}`,
    borderLeft: 'none',
    boxShadow: active
      ? 'inset 0 1px 0 rgba(255,240,180,0.35), inset 0 -2px 4px rgba(0,0,0,0.35), 3px 2px 10px rgba(0,0,0,0.60), 0 0 12px rgba(200,150,10,0.40)'
      : 'inset 0 1px 0 rgba(255,220,140,0.15), inset 0 -2px 4px rgba(0,0,0,0.40), 3px 2px 10px rgba(0,0,0,0.60)',
    userSelect: 'none',
    marginBottom: 2,
  }
}

const tabContentShellStyle = (usesSmallViewport: boolean): CSSProperties => ({
  position: 'absolute',
  left: usesSmallViewport ? '70.5%' : '14.2%',
  top: usesSmallViewport ? '8.4%' : '8.5%',
  width: usesSmallViewport ? '21.5%' : '77.2%',
  height: usesSmallViewport ? '76.8%' : '76.2%',
  zIndex: 3,
  minHeight: 0,
  overflow: 'hidden',
  boxSizing: 'border-box',
  borderRadius: usesSmallViewport ? 12 : 18,
  background: 'linear-gradient(180deg, rgba(255, 244, 224, 0.90), rgba(232, 205, 164, 0.84))',
  border: '1px solid rgba(118, 88, 54, 0.22)',
  boxShadow: '0 18px 42px rgba(37, 23, 8, 0.24), inset 0 1px 0 rgba(255,255,255,0.26)',
  backdropFilter: 'blur(2px)',
})

const diaryRightStackStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '1fr',
  gridTemplateRows: 'minmax(0, 1fr) auto',
  gap: 6,
  height: '100%',
  minHeight: 0,
}

const diaryOpenActionBarStyle: CSSProperties = {
  position: 'absolute',
  left: '75.4%',
  top: '78.4%',
  width: OPEN_RIGHT_SHELL.width,
  zIndex: 3,
  boxSizing: 'border-box',
  padding: '0 2px',
}

const diaryClosedActionBarStyle: CSSProperties = {
  position: 'absolute',
  left: CLOSED_RIGHT_SHELL.left,
  top: '77.5%',
  width: CLOSED_RIGHT_SHELL.width,
  zIndex: 3,
  boxSizing: 'border-box',
  padding: '0 4px',
}

const diaryActionBarInnerStyle: CSSProperties = {
  position: 'relative',
  display: 'grid',
  gridTemplateColumns: 'repeat(2, minmax(0, 1fr))',
  alignItems: 'center',
  gap: 6,
  padding: '10px 8px',
  background: 'linear-gradient(180deg, rgba(86, 50, 17, 0.82) 0%, rgba(50, 28, 10, 0.88) 100%)',
  borderRadius: 10,
  border: '1px solid rgba(196, 142, 54, 0.34)',
  boxShadow: 'inset 0 1px 0 rgba(255, 220, 150, 0.18), 0 6px 14px rgba(28, 14, 4, 0.34)',
  backdropFilter: 'blur(4px)',
}

const onboardingBlockerStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  zIndex: 7,
  background: 'rgba(18, 10, 3, 0.48)',
  backdropFilter: 'blur(1.5px)',
  pointerEvents: 'auto',
}

const onboardingRightShellStyle: CSSProperties = {
  zIndex: 9,
}

const onboardingActionBarTargetStyle: CSSProperties = {
  zIndex: 10,
  filter: 'drop-shadow(0 0 18px rgba(16, 185, 129, 0.56))',
}

const onboardingSaveOnlyOverlayStyle: CSSProperties = {
  position: 'absolute',
  right: 8,
  top: 10,
  bottom: 10,
  width: 'calc(50% - 11px)',
  zIndex: 13,
  borderRadius: 8,
  background: 'rgba(18, 10, 3, 0.58)',
  backdropFilter: 'blur(1px)',
  pointerEvents: 'auto',
}

const onboardingSaveButtonTargetStyle: CSSProperties = {
  position: 'relative',
  zIndex: 14,
  boxShadow: '0 0 0 2px rgba(16, 185, 129, 0.86), 0 0 18px rgba(16, 185, 129, 0.50)',
}

const onboardingGuideStyle: CSSProperties = {
  position: 'fixed',
  left: '50%',
  bottom: 18,
  width: 'min(760px, calc(100vw - 32px))',
  zIndex: 48,
  transform: 'translateX(-50%)',
  borderRadius: 8,
  border: '1px solid rgba(94, 234, 212, 0.20)',
  background: 'rgba(15, 23, 42, 0.92)',
  boxShadow: '0 18px 42px rgba(0, 0, 0, 0.38), inset 0 1px 0 rgba(255, 255, 255, 0.08)',
  backdropFilter: 'blur(14px)',
  WebkitBackdropFilter: 'blur(14px)',
  color: '#e2e8f0',
  padding: '10px 12px',
  boxSizing: 'border-box',
  pointerEvents: 'none',
}

const onboardingGuideSmallStyle: CSSProperties = {
  bottom: 10,
  width: 'calc(100vw - 20px)',
}

const onboardingGuideHeaderStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'auto minmax(0, 1fr) auto',
  alignItems: 'center',
  gap: 8,
  padding: '3px 8px',
  borderRadius: 4,
  border: '1px solid rgba(94, 234, 212, 0.14)',
  background: 'rgba(30, 41, 59, 0.58)',
  marginBottom: 8,
}

const onboardingGuideBrandStyle: CSSProperties = {
  display: 'inline-flex',
  alignItems: 'center',
  gap: 6,
  minWidth: 0,
  color: '#f8fafc',
  fontSize: 12,
  fontWeight: 900,
}

const onboardingGuideSkipStyle: CSSProperties = {
  border: '1px solid rgba(187, 247, 208, 0.28)',
  borderRadius: 999,
  background: 'rgba(187, 247, 208, 0.10)',
  color: '#d1fae5',
  padding: '3px 8px',
  fontSize: 11,
  fontWeight: 800,
  cursor: 'pointer',
  whiteSpace: 'nowrap',
  pointerEvents: 'auto',
}

const onboardingGuideSectionStyle: CSSProperties = {
  display: 'grid',
  gap: 5,
}

const onboardingGuideLabelStyle: CSSProperties = {
  color: '#5eead4',
  fontSize: 12,
  fontWeight: 900,
  lineHeight: 1.35,
}

const onboardingGuideTextStyle: CSSProperties = {
  margin: 0,
  color: '#cbd5e1',
  fontSize: 13,
  lineHeight: 1.42,
  fontWeight: 650,
}

const onboardingGuideBodyStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'auto minmax(0, 1fr)',
  gap: 10,
  alignItems: 'center',
}

const onboardingGuideTextWrapStyle: CSSProperties = {
  display: 'grid',
  gap: 6,
  minWidth: 0,
}

const onboardingGuideDividerStyle: CSSProperties = {
  height: 1,
  width: '100%',
  margin: '1px 0',
  background: 'linear-gradient(90deg, rgba(94, 234, 212, 0.22), rgba(51, 65, 85, 0.34))',
}

const onboardingBubbleStyle: CSSProperties = {
  position: 'absolute',
  left: '61.2%',
  top: '24.2%',
  zIndex: 12,
  maxWidth: 190,
  borderRadius: 10,
  border: '1px solid rgba(16, 185, 129, 0.44)',
  background: 'rgba(236, 253, 245, 0.96)',
  color: '#064e3b',
  padding: '8px 10px',
  fontSize: 12,
  fontWeight: 900,
  boxShadow: '0 12px 24px rgba(0, 0, 0, 0.24)',
  pointerEvents: 'none',
}

const onboardingAddBubbleStyle: CSSProperties = {
  left: '57%',
  top: '67%',
}

const onboardingSaveBubbleStyle: CSSProperties = {
  left: '62%',
  top: '82.8%',
}

const onboardingPointBubbleStyle: CSSProperties = {
  left: '57%',
  top: '35%',
}

export function PlanetDiaryStage({
  course,
  planetId,
  routeKind,
  activeTab,
  canCompleteCourse = false,
  isCompletingCourse = false,
  onOpenCourseCompleteConfirm,
  onActiveTabChange,
  onRefresh,
  copy,
}: {
  course: DiaryCourseAggregate
  planetId: string
  routeKind: PlanetRouteKind
  activeTab: DiaryStageTab
  canCompleteCourse?: boolean
  isCompletingCourse?: boolean
  onOpenCourseCompleteConfirm?: () => void
  onActiveTabChange: (tab: DiaryStageTab) => void
  onRefresh: () => Promise<void>
  copy: PlanetDiaryCopy
}) {
  const router = useRouter()
  const searchParams = useSearchParams()
  const explorerPlan = useExplorerPlan(course.course_draft_id)
  const editForm = useExplorerEditForm(explorerPlan)
  const { mapExpanded, setMapExpanded, isNarrowViewport } = useResponsiveMapExpanded()
  const [mapLevel, setMapLevel] = useState<PlanningMapLevel>('course')
  const [selectedRegionId, setSelectedRegionId] = useState<string | null>(null)
  const [selectedSubRegionId, setSelectedSubRegionId] = useState<string | null>(null)
  const [showExitConfirm, setShowExitConfirm] = useState(false)
  const [pendingPointId, setPendingPointId] = useState<string | null>(null)
  const [blockedPointId, setBlockedPointId] = useState<string | null>(null)
  const [showPlanningConfirm, setShowPlanningConfirm] = useState(false)
  const [plannedBookmark, setPlannedBookmark] = useState<PlannedBookmark | null>(null)
  const [firstElementOnboardingStep, setFirstElementOnboardingStep] = useState<FirstElementOnboardingStep | null>(null)
  const [recentlyCompletedOnboarding, setRecentlyCompletedOnboarding] = useState(false)
  const appliedDiaryFocusRef = useRef<string | null>(null)

  const displayCourse = useMemo(
    () => mergeDiaryCourseWithExplorer(course, explorerPlan.courseAggregate),
    [course, explorerPlan.courseAggregate],
  )
  const selectedRegion = useMemo(
    () => displayCourse.regions.find((region) => region.region.id === selectedRegionId) ?? null,
    [displayCourse.regions, selectedRegionId],
  )
  const selectedSubRegion = useMemo(
    () => selectedRegion?.subregions.find((subRegion) => subRegion.subregion.id === selectedSubRegionId) ?? null,
    [selectedRegion, selectedSubRegionId],
  )
  const pendingPointContext = useMemo(
    () => (pendingPointId ? findDiaryNodeById(displayCourse, pendingPointId) : null),
    [displayCourse, pendingPointId],
  )
  const blockedPointContext = useMemo(
    () => (blockedPointId ? findDiaryNodeById(displayCourse, blockedPointId) : null),
    [blockedPointId, displayCourse],
  )
  const planningEditHref =
    routeKind === 'learning'
      ? `/dashboard/course-drafts/${course.course_draft_id}?from=diary&section=planning`
      : null
  const firstRegion = useMemo(
    () => displayCourse.regions
      .filter((region) =>
        region.region.status !== 'inactive' &&
        explorerPlan.getRegionPendingState(region.region.id) !== 'deleted'
      )
      .sort((a, b) => a.region.order_index - b.region.order_index)[0] ?? null,
    [displayCourse.regions, explorerPlan],
  )
  const savedExplorerElementCount = useMemo(() => {
    let count = 0
    displayCourse.regions.forEach((region) => {
      region.subregions.forEach((subRegion) => {
        if (
          subRegion.subregion.status !== 'inactive' &&
          explorerPlan.getSubRegionPendingState(subRegion.subregion.id) !== 'created' &&
          explorerPlan.getSubRegionPendingState(subRegion.subregion.id) !== 'deleted'
        ) {
          count += 1
        }
        subRegion.nodes.forEach((node) => {
          if (node.status !== 'inactive' && explorerPlan.getNodePendingState(node.id) !== 'created' && explorerPlan.getNodePendingState(node.id) !== 'deleted') {
            count += 1
          }
        })
      })
      region.nodes.forEach((node) => {
        if (node.status !== 'inactive' && explorerPlan.getNodePendingState(node.id) !== 'created' && explorerPlan.getNodePendingState(node.id) !== 'deleted') {
          count += 1
        }
      })
    })
    return count
  }, [displayCourse.regions, explorerPlan])
  const pendingExplorerElementCount = useMemo(() => {
    let count = 0
    displayCourse.regions.forEach((region) => {
      region.subregions.forEach((subRegion) => {
        if (explorerPlan.getSubRegionPendingState(subRegion.subregion.id) === 'created') count += 1
        subRegion.nodes.forEach((node) => {
          if (explorerPlan.getNodePendingState(node.id) === 'created') count += 1
        })
      })
      region.nodes.forEach((node) => {
        if (explorerPlan.getNodePendingState(node.id) === 'created') count += 1
      })
    })
    return count
  }, [displayCourse.regions, explorerPlan])
  const firstSavedPoint = useMemo(() => {
    const sortedRegions = displayCourse.regions
      .filter((region) =>
        region.region.status !== 'inactive' &&
        explorerPlan.getRegionPendingState(region.region.id) !== 'deleted'
      )
      .sort((a, b) => a.region.order_index - b.region.order_index)
    for (const region of sortedRegions) {
      const directNode = [...region.nodes]
        .filter((node) =>
          node.status !== 'inactive' &&
          explorerPlan.getNodePendingState(node.id) !== 'created' &&
          explorerPlan.getNodePendingState(node.id) !== 'deleted'
        )
        .sort((a, b) => a.order_index - b.order_index)[0]
      if (directNode) return directNode
      const sortedSubRegions = [...region.subregions]
        .filter((subRegion) =>
          subRegion.subregion.status !== 'inactive' &&
          explorerPlan.getSubRegionPendingState(subRegion.subregion.id) !== 'created' &&
          explorerPlan.getSubRegionPendingState(subRegion.subregion.id) !== 'deleted'
        )
        .sort((a, b) => a.subregion.order_index - b.subregion.order_index)
      for (const subRegion of sortedSubRegions) {
        const subNode = [...subRegion.nodes]
          .filter((node) =>
            node.status !== 'inactive' &&
            explorerPlan.getNodePendingState(node.id) !== 'created' &&
            explorerPlan.getNodePendingState(node.id) !== 'deleted'
          )
          .sort((a, b) => a.order_index - b.order_index)[0]
        if (subNode) return subNode
      }
    }
    return null
  }, [displayCourse.regions, explorerPlan])
  const firstElementOnboardingKey = `learnweaver:onboarding:diary-first-element:${course.course_draft_id}`
  const canShowFirstElementOnboarding =
    routeKind === 'learning' &&
    activeTab === 'journal' &&
    firstRegion !== null &&
    (savedExplorerElementCount === 0 || firstElementOnboardingStep === 'open_point')
  const starSystemReturnTarget = course.course_draft_id || planetId
  const starSystemReturnHref = `/dashboard?returnCourse=${encodeURIComponent(starSystemReturnTarget)}`
  const handleSelectCourse = () => {
    explorerPlan.handleSelectCourse()
    setSelectedRegionId(null)
    setSelectedSubRegionId(null)
    setMapLevel('course')
  }

  const handleSelectRegion = (regionId: string) => {
    if (isFirstElementLessonStep && firstRegion && regionId !== firstRegion.region.id) return
    explorerPlan.handleSelectRegion(regionId)
    setSelectedRegionId(regionId)
    setSelectedSubRegionId(null)
    setMapLevel('region')
  }

  const handleSelectSubRegion = (regionId: string, subRegionId: string) => {
    explorerPlan.handleSelectSubRegion(regionId, subRegionId)
    setSelectedRegionId(regionId)
    setSelectedSubRegionId(subRegionId)
    setMapLevel('subregion')
  }

  const handleMapSelectedRegionChange = (regionId: string | null) => {
    if (regionId) explorerPlan.handleSelectRegion(regionId)
    else explorerPlan.handleSelectCourse()
    setSelectedRegionId(regionId)
  }

  const handleMapSelectedSubRegionChange = (subRegionId: string | null) => {
    if (subRegionId) {
      const parentRegion = displayCourse.regions.find((region) =>
        region.subregions.some((subRegion) => subRegion.subregion.id === subRegionId),
      )
      if (parentRegion) explorerPlan.handleSelectSubRegion(parentRegion.region.id, subRegionId)
    } else if (selectedRegionId) {
      explorerPlan.handleSelectRegion(selectedRegionId)
    }
    setSelectedSubRegionId(subRegionId)
  }

  useEffect(() => {
    const focusedLessonId = searchParams.get('diaryLessonId')
    if (!focusedLessonId || appliedDiaryFocusRef.current === focusedLessonId) return

    const focusedRegion = displayCourse.regions.find((region) => region.region.id === focusedLessonId)
    if (focusedRegion) {
      appliedDiaryFocusRef.current = focusedLessonId
      if (activeTab !== 'journal') onActiveTabChange('journal')
      explorerPlan.handleSelectRegion(focusedRegion.region.id)
      setSelectedRegionId(focusedRegion.region.id)
      setSelectedSubRegionId(null)
      setMapLevel('region')
      return
    }

    const parentRegion = displayCourse.regions.find((region) =>
      region.subregions.some((subRegion) => subRegion.subregion.id === focusedLessonId),
    )
    if (!parentRegion) return

    appliedDiaryFocusRef.current = focusedLessonId
    if (activeTab !== 'journal') onActiveTabChange('journal')
    explorerPlan.handleSelectSubRegion(parentRegion.region.id, focusedLessonId)
    setSelectedRegionId(parentRegion.region.id)
    setSelectedSubRegionId(focusedLessonId)
    setMapLevel('subregion')
  }, [activeTab, displayCourse.regions, explorerPlan, onActiveTabChange, searchParams])

  useEffect(() => {
    if (!canShowFirstElementOnboarding) {
      setFirstElementOnboardingStep(null)
      return
    }

    const stored = window.localStorage.getItem(firstElementOnboardingKey)
    if (stored === 'skipped') {
      setFirstElementOnboardingStep(null)
      return
    }
    if (stored === 'done') {
      setFirstElementOnboardingStep(null)
      return
    }
    if (stored === 'saved') {
      setFirstElementOnboardingStep('open_point')
      return
    }

    if (savedExplorerElementCount > 0) {
      if (firstElementOnboardingStep === 'open_point') {
        setFirstElementOnboardingStep('open_point')
      } else {
        setFirstElementOnboardingStep(null)
      }
      return
    }

    if (editForm.addModalType !== null) {
      setFirstElementOnboardingStep('modal')
      return
    }

    if (pendingExplorerElementCount > 0 || explorerPlan.isDirty) {
      setFirstElementOnboardingStep('save')
      return
    }

    if (selectedRegionId === firstRegion?.region.id && selectedSubRegionId == null) {
      setFirstElementOnboardingStep('add_element')
      return
    }

    setFirstElementOnboardingStep('lesson')
  }, [
    canShowFirstElementOnboarding,
    editForm.addModalType,
    explorerPlan.isDirty,
    firstElementOnboardingStep,
    firstElementOnboardingKey,
    firstRegion?.region.id,
    pendingExplorerElementCount,
    savedExplorerElementCount,
    selectedRegionId,
    selectedSubRegionId,
  ])

  useEffect(() => {
    if (!recentlyCompletedOnboarding) return
    const timeoutId = window.setTimeout(() => setRecentlyCompletedOnboarding(false), 5200)
    return () => window.clearTimeout(timeoutId)
  }, [recentlyCompletedOnboarding])

  useEffect(() => {
    if (!firstElementOnboardingStep) return
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key !== 'Escape') return
      window.localStorage.setItem(firstElementOnboardingKey, 'skipped')
      setFirstElementOnboardingStep(null)
    }
    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [firstElementOnboardingKey, firstElementOnboardingStep])

  const handleSkipFirstElementOnboarding = useCallback(() => {
    window.localStorage.setItem(firstElementOnboardingKey, 'skipped')
    setFirstElementOnboardingStep(null)
  }, [firstElementOnboardingKey])

  const handleOnboardingElementAdded = useCallback(() => {
    setFirstElementOnboardingStep('save')
  }, [])

  const markFirstElementOnboardingDone = useCallback(() => {
    window.localStorage.setItem(firstElementOnboardingKey, 'done')
    setFirstElementOnboardingStep(null)
    setRecentlyCompletedOnboarding(true)
  }, [firstElementOnboardingKey])

  const handleOpenPoint = (pointId: string) => {
    if (firstElementOnboardingStep === 'open_point' && firstSavedPoint && pointId !== firstSavedPoint.id) return
    if (firstElementOnboardingStep === 'open_point') {
      markFirstElementOnboardingDone()
    }
    if (explorerPlan.isDirty) {
      setBlockedPointId(pointId)
      return
    }
    setPendingPointId(pointId)
  }

  const bookmarkTabs: {
    label: string
    shortLabel: string
    key: DiaryStageTab | 'planning' | 'community' | 'civilization'
    state: BookmarkState
  }[] = [
    {
      label: copy.stage.tabs.planning.label,
      shortLabel: copy.stage.tabs.planning.shortLabel,
      key: 'planning',
      state: routeKind === 'learning' ? 'inactive' : 'locked',
    },
    {
      label: copy.stage.tabs.journal.label,
      shortLabel: copy.stage.tabs.journal.shortLabel,
      key: 'journal',
      state: activeTab === 'journal' ? 'active' : 'inactive',
    },
    {
      label: copy.stage.tabs.records.label,
      shortLabel: copy.stage.tabs.records.shortLabel,
      key: 'records',
      state: activeTab === 'records' ? 'active' : 'inactive',
    },
    {
      label: copy.stage.tabs.results.label,
      shortLabel: copy.stage.tabs.results.shortLabel,
      key: 'results',
      state: activeTab === 'results' ? 'active' : 'inactive',
    },
    { label: copy.stage.tabs.community.label, shortLabel: copy.stage.tabs.community.shortLabel, key: 'community', state: 'locked' },
    { label: copy.stage.tabs.civilization.label, shortLabel: copy.stage.tabs.civilization.shortLabel, key: 'civilization', state: 'locked' },
  ]

  const viewportBase = isNarrowViewport ? SMALL_DEVICE_CROP : BASE_STAGE
  const { setFrameNode, scale } = useResponsiveCanvasScale(viewportBase.width, viewportBase.height)
  const panelUsesSmallDevice = mapExpanded || isNarrowViewport
  const backgroundSrc = panelUsesSmallDevice ? OPEN_BACKGROUND_SRC : CLOSED_BACKGROUND_SRC
  const isJournalActionDisabled = explorerPlan.isMutating || editForm.isMutating
  const showJournalActionBar = activeTab === 'journal' && routeKind === 'learning'
  const isFirstElementOnboardingActive = firstElementOnboardingStep !== null
  const isFirstElementLessonStep = firstElementOnboardingStep === 'lesson'
  const isFirstElementAddStep = firstElementOnboardingStep === 'add_element'
  const isFirstElementSaveStep = firstElementOnboardingStep === 'save'
  const isFirstElementModalStep = firstElementOnboardingStep === 'modal'
  const isFirstElementOpenPointStep = firstElementOnboardingStep === 'open_point'

  const handleSaveJournalChanges = async () => {
    await editForm.handleSave()
    await onRefresh()
    if (firstElementOnboardingStep === 'save') {
      window.localStorage.setItem(firstElementOnboardingKey, 'saved')
      setFirstElementOnboardingStep('open_point')
    }
  }

  const handleCancelJournalChanges = () => {
    editForm.handleCancel()
    handleSelectCourse()
  }
  const firstElementGuideStep = firstElementOnboardingStep ?? (recentlyCompletedOnboarding ? 'done' : null)
  const diaryGuideStep = firstElementGuideStep ?? (savedExplorerElementCount === 0 ? 'empty' : 'default')

  return (
    <section style={outerSectionStyle}>
      <div
        ref={setFrameNode}
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
            <div style={backgroundWrapStyle}>
              <img
                src={backgroundSrc}
                alt={copy.stage.alt}
                style={{ display: 'block', width: '100%', height: '100%', objectFit: 'fill' }}
              />
            </div>
            <div style={overlayStyle} />
            {isFirstElementOnboardingActive && !isFirstElementModalStep ? (
              <div
                style={onboardingBlockerStyle}
                aria-hidden="true"
                onMouseDown={(event) => event.preventDefault()}
              />
            ) : null}

            {activeTab === 'journal' ? (
              <>
                <DiaryMapPanel
                  course={displayCourse}
                  mapExpanded={mapExpanded}
                  mapLevel={mapLevel}
                  selectedRegionId={selectedRegionId}
                  selectedSubRegionId={selectedSubRegionId}
                  onMapLevelChange={setMapLevel}
                  onSelectedRegionIdChange={handleMapSelectedRegionChange}
                  onSelectedSubRegionIdChange={handleMapSelectedSubRegionChange}
                  onOpenPoint={handleOpenPoint}
                  onRequestExit={() => setShowExitConfirm(true)}
                  copy={copy.mapPanel}
                />

                <div style={{
                  ...(panelUsesSmallDevice ? OPEN_RIGHT_SHELL : CLOSED_RIGHT_SHELL),
                  ...(isFirstElementOnboardingActive && !isFirstElementSaveStep ? onboardingRightShellStyle : undefined),
                }}>
                  <div style={diaryRightStackStyle}>
                    <DiaryTreePanel
                      course={displayCourse}
                      selectedRegionId={selectedRegionId}
                      selectedSubRegionId={selectedSubRegionId}
                      canCompleteCourse={canCompleteCourse}
                      isCompletingCourse={isCompletingCourse}
                      getRegionPendingState={explorerPlan.getRegionPendingState}
                      getSubRegionPendingState={explorerPlan.getSubRegionPendingState}
                      getNodePendingState={explorerPlan.getNodePendingState}
                      onOpenCourseCompleteConfirm={onOpenCourseCompleteConfirm}
                      onSelectCourse={handleSelectCourse}
                      onSelectRegion={handleSelectRegion}
                      onSelectSubRegion={handleSelectSubRegion}
                      onOpenPoint={handleOpenPoint}
                      copy={copy.tree}
                      onboardingFirstRegionId={isFirstElementLessonStep ? firstRegion?.region.id ?? null : null}
                      onboardingPointId={isFirstElementOpenPointStep ? firstSavedPoint?.id ?? null : null}
                      onboardingLockToFirstRegion={isFirstElementLessonStep}
                      onboardingLockToPoint={isFirstElementOpenPointStep}
                      onboardingDisableTreeInteractions={isFirstElementAddStep}
                    />
                    <DiaryGuidePanel
                      explorerPlan={explorerPlan}
                      editForm={editForm}
                      selectedRegion={selectedRegion}
                      selectedSubRegion={selectedSubRegion}
                      routeKind={routeKind}
                      copy={copy.guidePanel}
                      isOnboardingAddTarget={isFirstElementAddStep}
                      onOnboardingAddClick={() => setFirstElementOnboardingStep('modal')}
                      onOnboardingElementAdded={handleOnboardingElementAdded}
                    />
                  </div>
                </div>

                {showJournalActionBar ? (
                  <div style={{
                    ...(panelUsesSmallDevice ? diaryOpenActionBarStyle : diaryClosedActionBarStyle),
                    ...(isFirstElementSaveStep ? onboardingActionBarTargetStyle : undefined),
                  }}>
                    <div style={diaryActionBarInnerStyle}>
                      <button
                        type="button"
                        style={{
                          ...saveButtonStyle,
                          ...(isFirstElementSaveStep ? onboardingSaveButtonTargetStyle : undefined),
                          ...(isJournalActionDisabled ? treeButtonDisabledStyle : undefined),
                        }}
                        disabled={isJournalActionDisabled}
                        onClick={() => void handleSaveJournalChanges()}
                      >
                        {copy.stage.save}
                      </button>
                      <button
                        type="button"
                        style={{
                          ...cancelButtonStyle,
                          ...(isJournalActionDisabled ? treeButtonDisabledStyle : undefined),
                        }}
                        disabled={isJournalActionDisabled}
                        onClick={handleCancelJournalChanges}
                      >
                        {copy.stage.cancel}
                      </button>
                      {isFirstElementSaveStep ? <div style={onboardingSaveOnlyOverlayStyle} aria-hidden="true" /> : null}
                    </div>
                  </div>
                ) : null}

                {isFirstElementLessonStep ? (
                  <div style={onboardingBubbleStyle}>{copy.firstElementOnboarding.lessonBubble}</div>
                ) : null}
                {isFirstElementAddStep ? (
                  <div style={{ ...onboardingBubbleStyle, ...onboardingAddBubbleStyle }}>
                    {copy.firstElementOnboarding.addBubble}
                  </div>
                ) : null}
                {isFirstElementSaveStep ? (
                  <div style={{ ...onboardingBubbleStyle, ...onboardingSaveBubbleStyle }}>
                    {copy.firstElementOnboarding.saveBubble}
                  </div>
                ) : null}
                {isFirstElementOpenPointStep ? (
                  <div style={{ ...onboardingBubbleStyle, ...onboardingPointBubbleStyle }}>
                    {copy.firstElementOnboarding.openPointBubble}
                  </div>
                ) : null}

                <button
                  type="button"
                  style={mapExpanded ? collapseButtonStyle : expandButtonStyle}
                  onClick={() => setMapExpanded((prev) => !prev)}
                >
                  {mapExpanded ? copy.stage.collapseMap : copy.stage.expandMap}
                </button>
              </>
            ) : (
              <div style={tabContentShellStyle(isNarrowViewport)}>
                {activeTab === 'records' ? (
                  <PlanetRecordPage planetId={planetId} routeKind={routeKind} />
                ) : (
                  <PlanetResultPage planetId={planetId} routeKind={routeKind} />
                )}
              </div>
            )}

            <div style={bookmarkRailStyle(panelUsesSmallDevice)}>
              {bookmarkTabs.map((tab) => {
                const plannedKey =
                  tab.key === 'community'
                    ? 'community'
                    : tab.key === 'civilization'
                      ? 'civilization'
                      : null
                const isDisabled = tab.state === 'locked' && plannedKey == null
                return (
                  <button
                    type="button"
                    key={tab.label}
                    style={{
                      ...bookmarkTabStyle(tab.state),
                      ...(plannedKey ? { cursor: 'pointer' } : null),
                    }}
                    role="tab"
                    aria-label={tab.label}
                    aria-selected={tab.state === 'active'}
                    aria-disabled={tab.state === 'locked'}
                    disabled={isDisabled}
                    title={tab.state === 'locked' ? copy.stage.readyTitle(tab.label) : tab.label}
                    onClick={() => {
                      if (plannedKey) {
                        setPlannedBookmark(plannedKey)
                        return
                      }
                      if (tab.key === 'planning' && planningEditHref) {
                        setShowPlanningConfirm(true)
                        return
                      }
                      if (tab.key === 'journal' || tab.key === 'records' || tab.key === 'results') onActiveTabChange(tab.key)
                    }}
                  >
                    {tab.shortLabel}
                  </button>
                )
              })}
            </div>
          </div>
        </div>
      </div>
      {routeKind === 'learning' && activeTab === 'journal' ? (
        <DiaryFirstElementOnboardingGuide
          step={diaryGuideStep}
          usesSmallViewport={isNarrowViewport}
          copy={copy.firstElementOnboarding}
          onSkip={handleSkipFirstElementOnboarding}
        />
      ) : null}
      <DiaryStageModals
        showExitConfirm={showExitConfirm}
        showPlanningConfirm={showPlanningConfirm}
        planningEditHref={planningEditHref}
        pendingPointContext={pendingPointContext}
        blockedPointContext={blockedPointContext}
        plannedBookmark={plannedBookmark}
        onCloseExitConfirm={() => setShowExitConfirm(false)}
        onConfirmExit={() => router.push(starSystemReturnHref)}
        onClosePlanningConfirm={() => setShowPlanningConfirm(false)}
        onConfirmPlanning={() => planningEditHref ? window.location.assign(planningEditHref) : undefined}
        onClosePendingPoint={() => setPendingPointId(null)}
        onConfirmPendingPoint={() => {
          if (!pendingPointContext) return
          router.push(`/dashboard/planets/${routeKind}/${pendingPointContext.node.course_id ?? planetId}/points/${pendingPointContext.node.id}`)
        }}
        onClosePlannedBookmark={() => setPlannedBookmark(null)}
        onCloseBlockedPoint={() => setBlockedPointId(null)}
        copy={copy}
      />
    </section>
  )
}

function DiaryFirstElementOnboardingGuide({
  step,
  usesSmallViewport,
  copy,
  onSkip,
}: {
  step: FirstElementOnboardingStep | 'done' | 'empty' | 'default'
  usesSmallViewport: boolean
  copy: PlanetDiaryCopy['firstElementOnboarding']
  onSkip: () => void
}) {
  const [mounted, setMounted] = useState(false)
  useEffect(() => setMounted(true), [])

  const message = (() => {
    if (step === 'lesson') return { task: copy.lessonTask, feature: copy.lessonFeature }
    if (step === 'add_element') return { task: copy.addTask, feature: copy.addFeature }
    if (step === 'modal') return { task: copy.modalTask, feature: copy.modalFeature }
    if (step === 'save') return { task: copy.saveTask, feature: copy.saveFeature }
    if (step === 'open_point') return { task: copy.openPointTask, feature: copy.openPointFeature }
    if (step === 'empty') return { task: copy.emptyTask, feature: copy.emptyFeature }
    if (step === 'default') return { task: copy.defaultTask, feature: copy.defaultFeature }
    return { task: copy.doneTask, feature: copy.doneFeature }
  })()
  const isOverlayStep = step === 'lesson' || step === 'add_element' || step === 'modal' || step === 'save' || step === 'open_point'

  if (!mounted) return null

  return createPortal(
    <aside
      style={{
        ...onboardingGuideStyle,
        ...(usesSmallViewport ? onboardingGuideSmallStyle : undefined),
      }}
      aria-live="polite"
    >
      <div style={onboardingGuideHeaderStyle}>
        <div style={onboardingGuideBrandStyle}>
          <span aria-hidden="true">🎯</span>
          <span>Lumi Guide</span>
        </div>
        <span
          title={message.task}
          style={{
            minWidth: 0,
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
            fontSize: 13,
            lineHeight: 1.35,
            fontWeight: 650,
            color: '#cbd5e1',
          }}
        >
          {message.task}
        </span>
        {isOverlayStep ? (
          <button type="button" style={onboardingGuideSkipStyle} onClick={onSkip}>
            {copy.skip}
          </button>
        ) : null}
      </div>
      <div style={onboardingGuideBodyStyle}>
        <LumiAvatar state={step === 'done' ? 'happy' : 'exploring'} size={34} reducedMotion={false} />
        <div style={onboardingGuideTextWrapStyle}>
          <div style={onboardingGuideSectionStyle}>
            <span style={onboardingGuideLabelStyle}>{copy.nowTitle}</span>
            <p style={onboardingGuideTextStyle}>{message.task}</p>
          </div>
          <div style={onboardingGuideDividerStyle} aria-hidden="true" />
          <div style={onboardingGuideSectionStyle}>
            <span style={onboardingGuideLabelStyle}>{copy.featureTitle}</span>
            <p style={onboardingGuideTextStyle}>{message.feature}</p>
          </div>
        </div>
      </div>
    </aside>,
    document.body,
  )
}
