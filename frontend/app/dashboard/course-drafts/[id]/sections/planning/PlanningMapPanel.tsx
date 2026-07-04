'use client'

import type { CSSProperties } from 'react'
import { useState } from 'react'
import { useRouter } from 'next/navigation'
import LumiModalShell from '@/components/common/LumiModalShell'
import PlanetReturnPlanetButton from '@/components/dashboard/PlanetReturnPlanetButton'

import type { ExplorerPlanState } from '@/components/explorer-plan/useExplorerPlan'
import type { PlanningLumiGuideFocus } from './PlanningLumiGuide'
import { PlanningMapSvg } from './map/PlanningMapSvg'
import PlanetMapFogLayer from './map/PlanetMapFogLayer'
import { PlanningRegionDetail } from './map/PlanningRegionDetail'
import { PlanningSubRegionDetail } from './map/PlanningSubRegionDetail'
import { PlanningMapHeader } from './map/PlanningMapHeader'
import { PlanningMapPager } from './map/PlanningMapPager'
import { usePlanningMapNavigation, type PlanningMapLevel } from './map/usePlanningMapNavigation'
import type { DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'

// ── 위치 상수 ─────────────────────────────────────────────────────────────────
// Explorer_Diary_Design.webp (1395×757) — 펼침 양피지 지도 영역
// x≈199~989, y≈64~641
const MAP_LEFT   = '14.3%'
const MAP_TOP    = '8.5%'
const MAP_WIDTH  = '56.6%'
const MAP_HEIGHT = '76.2%'
// 폴리곤 지도는 헤더의 코스명(예: "바리스타 배우기") 아래 약 10px부터
// 하단 페이지 표시 객체 위 약 10px까지를 실제 표시 영역으로 사용한다.
// 좌우는 페이지 이동 화살표 42px + 약 10px 여유를 피해 밝은 영역을 표시한다.
const POLYGON_AREA_SIDE = '7.4%'
const POLYGON_AREA_TOP = '23.0%'
const POLYGON_AREA_BOTTOM = '7.2%'

// ── 스타일 ────────────────────────────────────────────────────────────────────
const shellStyle: CSSProperties = {
  position: 'absolute',
  left: MAP_LEFT,
  top: MAP_TOP,
  width: MAP_WIDTH,
  height: MAP_HEIGHT,
  zIndex: 2,
  color: '#4a3520',
}

const svgAreaStyle: CSSProperties = {
  position: 'absolute',
  left: POLYGON_AREA_SIDE,
  right: POLYGON_AREA_SIDE,
  top: POLYGON_AREA_TOP,
  bottom: POLYGON_AREA_BOTTOM,
  overflow: 'hidden',
}

const emptyStateStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  padding: 24,
  color: 'rgba(90, 67, 48, 0.78)',
  textAlign: 'center',
  fontSize: 15,
  borderRadius: 12,
}

// ── 컴포넌트 ──────────────────────────────────────────────────────────────────
interface PlanningMapPanelProps {
  explorerPlan: ExplorerPlanState
  mapExpanded: boolean
  mapLevel: PlanningMapLevel
  selectedRegionId: string | null
  selectedSubRegionId: string | null
  onMapLevelChange: (level: PlanningMapLevel) => void
  onSelectedRegionIdChange: (regionId: string | null) => void
  onSelectedSubRegionIdChange: (subRegionId: string | null) => void
  onSelectionFocusChange?: (focus: PlanningLumiGuideFocus) => void
  onNodeActivate?: (nodeId: string) => void
  copy: DashboardCourseDraftCopy['planning']['map']
}

export function PlanningMapPanel({
  explorerPlan,
  mapExpanded,
  mapLevel,
  selectedRegionId,
  selectedSubRegionId,
  onMapLevelChange,
  onSelectedRegionIdChange,
  onSelectedSubRegionIdChange,
  onSelectionFocusChange,
  onNodeActivate,
  copy,
}: PlanningMapPanelProps) {
  const router = useRouter()
  const [showExitConfirm, setShowExitConfirm] = useState(false)
  const regions = explorerPlan.courseAggregate?.regions ?? []

  const {
    currentPage,
    totalPages,
    pagedRegions,
    regionIndexOffset,
    selectedRegion,
    selectedSubRegion,
    isTransitioning,
    direction,
    goPrev,
    goNext,
    handleRegionClick: navigateToRegion,
    handleSubRegionClick: navigateToSubRegion,
    handleBackToCourseMap,
    handleBackToRegionMap,
  } = usePlanningMapNavigation(regions, {
    mapLevel,
    selectedRegionId,
    selectedSubRegionId,
    onMapLevelChange,
    onSelectedRegionIdChange,
    onSelectedSubRegionIdChange,
  })

  const handleRegionClick = (regionId: string) => {
    navigateToRegion(regionId)
    explorerPlan.handleSelectRegion(regionId)
    onSelectionFocusChange?.('region')
  }

  const handleSubRegionClick = (subRegionId: string) => {
    const regionId = selectedRegion?.region.id
    if (!regionId) return
    navigateToSubRegion(subRegionId)
    explorerPlan.handleSelectSubRegion(regionId, subRegionId)
    onSelectionFocusChange?.('subregion')
  }

  const handleRegionNodeClick = (nodeId: string) => {
    const node = selectedRegion?.nodes.find((item) => item.id === nodeId)
    explorerPlan.handleSelectNode(nodeId)
    if (node) onSelectionFocusChange?.(node.node_type === 'research' ? 'research-node' : 'exploration-node')
    onNodeActivate?.(nodeId)
    if (!selectedRegion) return
    onSelectedRegionIdChange(selectedRegion.region.id)
    onSelectedSubRegionIdChange(null)
    onMapLevelChange('region')
  }

  const handleSubRegionNodeClick = (nodeId: string) => {
    const node = selectedSubRegion?.nodes.find((item) => item.id === nodeId)
    explorerPlan.handleSelectNode(nodeId)
    if (node) onSelectionFocusChange?.(node.node_type === 'research' ? 'research-node' : 'exploration-node')
    onNodeActivate?.(nodeId)
    if (!selectedRegion || !selectedSubRegion) return
    onSelectedRegionIdChange(selectedRegion.region.id)
    onSelectedSubRegionIdChange(selectedSubRegion.subregion.id)
    onMapLevelChange('subregion')
  }

  const handleCourseBack = () => {
    handleBackToCourseMap()
    explorerPlan.handleSelectCourse()
    onSelectionFocusChange?.('course')
  }

  const handleRegionBack = () => {
    handleBackToRegionMap()
    if (selectedRegion) {
      explorerPlan.handleSelectRegion(selectedRegion.region.id)
      onSelectionFocusChange?.('region')
    }
  }

  if (!mapExpanded) return null

  return (
    <div style={shellStyle}>
      <PlanningMapHeader
        mapLevel={mapLevel}
        courseTitle={explorerPlan.courseAggregate?.title ?? null}
        selectedRegionTitle={selectedRegion?.region.name ?? null}
        selectedSubRegionTitle={selectedSubRegion?.subregion.name ?? null}
        copy={copy}
      />

      {mapLevel === 'course' && (
        <PlanetReturnPlanetButton
          ariaLabel={copy.returnToStarSystemAria}
          textureMapAsset={explorerPlan.courseAggregate?.planet_texture_map_asset}
          progressPercent={explorerPlan.courseAggregate?.progress ?? 0}
          style={{
            position: 'absolute',
            left: '2.2%',
            top: '4.6%',
          }}
          onClick={() => setShowExitConfirm(true)}
        >
          {copy.returnToStarSystemLine1}
          <br />
          {copy.returnToStarSystemLine2}
        </PlanetReturnPlanetButton>
      )}

      {mapLevel === 'course' && (
        <PlanningMapPager
          currentPage={currentPage}
          totalPages={totalPages}
          isTransitioning={isTransitioning}
          onPrev={goPrev}
          onNext={goNext}
          copy={copy}
        />
      )}

      {mapLevel === 'course' ? (
        <div style={svgAreaStyle}>
          <PlanetMapFogLayer
            regions={pagedRegions}
            isTransitioning={isTransitioning}
            progressPercent={explorerPlan.courseAggregate?.progress ?? 0}
          />
          {pagedRegions.length === 0 ? (
            <div style={emptyStateStyle}>
              {copy.emptyCourseLine1}
              <br />
              {copy.emptyCourseLine2}
            </div>
          ) : (
            <PlanningMapSvg
              regions={pagedRegions}
              regionIndexOffset={regionIndexOffset}
              hasPreviousPage={currentPage > 0}
              hasNextPage={currentPage < totalPages - 1}
              isTransitioning={isTransitioning}
              direction={direction}
              onRegionClick={handleRegionClick}
              copy={copy}
            />
          )}
        </div>
      ) : mapLevel === 'subregion' ? (
        selectedSubRegion ? (
          <PlanningSubRegionDetail
            subRegion={selectedSubRegion}
            onBack={handleRegionBack}
            onNodeClick={handleSubRegionNodeClick}
            textureMapAsset={explorerPlan.courseAggregate?.planet_texture_map_asset}
            progressPercent={explorerPlan.courseAggregate?.progress ?? 0}
            copy={copy}
          />
        ) : (
          <div style={svgAreaStyle}>
            <div style={emptyStateStyle}>{copy.missingSubregion}</div>
          </div>
        )
      ) : selectedRegion ? (
        <PlanningRegionDetail
          region={selectedRegion}
          onBack={handleCourseBack}
          onSubRegionClick={handleSubRegionClick}
          onNodeClick={handleRegionNodeClick}
          textureMapAsset={explorerPlan.courseAggregate?.planet_texture_map_asset}
          progressPercent={explorerPlan.courseAggregate?.progress ?? 0}
          copy={copy}
        />
      ) : (
        <div style={svgAreaStyle}>
          <div style={emptyStateStyle}>
            {copy.missingRegion}
          </div>
        </div>
      )}
      {showExitConfirm ? (
        <LumiModalShell
          title={copy.exitTitle}
          eyebrow={copy.exitEyebrow}
          lumiState="planet-hold"
          message={copy.exitMessage}
          onClose={() => setShowExitConfirm(false)}
          actions={
            <>
              <button type="button" onClick={() => setShowExitConfirm(false)} style={modalSecondaryButtonStyle}>
                {copy.exitCancel}
              </button>
              <button
                type="button"
                onClick={() => {
                  const courseId = explorerPlan.courseAggregate?.course_draft_id
                  router.push(courseId ? `/dashboard?returnCourse=${courseId}` : '/dashboard')
                }}
                style={modalPrimaryButtonStyle}
              >
                {copy.exitConfirm}
              </button>
            </>
          }
        />
      ) : null}
    </div>
  )
}

const modalSecondaryButtonStyle: CSSProperties = {
  minHeight: 40,
  padding: '0 14px',
  borderRadius: 999,
  border: '1px solid rgba(122, 90, 36, 0.22)',
  background: 'rgba(255,255,255,0.68)',
  color: '#5F4319',
  fontSize: 13,
  fontWeight: 800,
  cursor: 'pointer',
}

const modalPrimaryButtonStyle: CSSProperties = {
  minHeight: 40,
  padding: '0 16px',
  borderRadius: 999,
  border: '1px solid rgba(173, 120, 36, 0.34)',
  background: 'linear-gradient(180deg, rgba(244, 214, 153, 0.96), rgba(205, 154, 67, 0.96))',
  color: '#4A2C00',
  fontSize: 13,
  fontWeight: 900,
  cursor: 'pointer',
}
