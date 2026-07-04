'use client'

import { PlanningMapHeader } from '@/app/dashboard/course-drafts/[id]/sections/planning/map/PlanningMapHeader'
import { PlanningMapPager } from '@/app/dashboard/course-drafts/[id]/sections/planning/map/PlanningMapPager'
import { PlanningMapSvg } from '@/app/dashboard/course-drafts/[id]/sections/planning/map/PlanningMapSvg'
import PlanetMapFogLayer from '@/app/dashboard/course-drafts/[id]/sections/planning/map/PlanetMapFogLayer'
import { PlanningRegionDetail } from '@/app/dashboard/course-drafts/[id]/sections/planning/map/PlanningRegionDetail'
import { PlanningSubRegionDetail } from '@/app/dashboard/course-drafts/[id]/sections/planning/map/PlanningSubRegionDetail'
import { usePlanningMapNavigation, type PlanningMapLevel } from '@/app/dashboard/course-drafts/[id]/sections/planning/map/usePlanningMapNavigation'
import PlanetReturnPlanetButton from '@/components/dashboard/PlanetReturnPlanetButton'
import type { PlanetDiaryCopy } from '@/lib/i18n/pages/planetDiary'
import type { DiaryCourseAggregate } from './planetDiaryTypes'

export function DiaryMapPanel({
  course,
  mapExpanded,
  mapLevel,
  selectedRegionId,
  selectedSubRegionId,
  onMapLevelChange,
  onSelectedRegionIdChange,
  onSelectedSubRegionIdChange,
  onOpenPoint,
  onRequestExit,
  copy,
}: {
  course: DiaryCourseAggregate
  mapExpanded: boolean
  mapLevel: PlanningMapLevel
  selectedRegionId: string | null
  selectedSubRegionId: string | null
  onMapLevelChange: (level: PlanningMapLevel) => void
  onSelectedRegionIdChange: (regionId: string | null) => void
  onSelectedSubRegionIdChange: (subRegionId: string | null) => void
  onOpenPoint: (nodeId: string) => void
  onRequestExit: () => void
  copy: PlanetDiaryCopy['mapPanel']
}) {
  const regions = course.regions

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

  if (!mapExpanded) return null

  return (
    <div
      style={{
        position: 'absolute',
        left: '14.3%',
        top: '8.5%',
        width: '56.6%',
        height: '76.2%',
        zIndex: 2,
        color: '#4a3520',
      }}
    >
      <PlanningMapHeader
        mapLevel={mapLevel}
        courseTitle={course.title}
        selectedRegionTitle={selectedRegion?.region.name ?? null}
        selectedSubRegionTitle={selectedSubRegion?.subregion.name ?? null}
      />

      {mapLevel === 'course' ? (
        <PlanetReturnPlanetButton
          ariaLabel={copy.backToStarSystemAria}
          onClick={onRequestExit}
          textureMapAsset={course.planet_texture_map_asset}
          progressPercent={course.progress ?? 0}
          style={{
            position: 'absolute',
            left: '2.2%',
            top: '4.6%',
          }}
        >
          {copy.backToStarSystemLines[0]}
          <br />
          {copy.backToStarSystemLines[1]}
        </PlanetReturnPlanetButton>
      ) : null}

      {mapLevel === 'course' ? (
        <PlanningMapPager
          currentPage={currentPage}
          totalPages={totalPages}
          isTransitioning={isTransitioning}
          onPrev={goPrev}
          onNext={goNext}
        />
      ) : null}

      {mapLevel === 'course' ? (
        <div style={{ position: 'absolute', left: '7.4%', right: '7.4%', top: '23%', bottom: '7.2%', overflow: 'hidden' }}>
          <PlanetMapFogLayer
            regions={pagedRegions}
            isTransitioning={isTransitioning}
            progressPercent={course.progress ?? 0}
          />
          <PlanningMapSvg
            regions={pagedRegions}
            regionIndexOffset={regionIndexOffset}
            hasPreviousPage={currentPage > 0}
            hasNextPage={currentPage < totalPages - 1}
            isTransitioning={isTransitioning}
            direction={direction}
            onRegionClick={navigateToRegion}
          />
        </div>
      ) : mapLevel === 'subregion' ? (
        selectedSubRegion ? (
          <PlanningSubRegionDetail
            subRegion={selectedSubRegion}
            onBack={handleBackToRegionMap}
            onNodeClick={onOpenPoint}
            textureMapAsset={course.planet_texture_map_asset}
            progressPercent={course.progress ?? 0}
          />
        ) : null
      ) : selectedRegion ? (
        <PlanningRegionDetail
          region={selectedRegion}
          onBack={handleBackToCourseMap}
          onSubRegionClick={navigateToSubRegion}
          onNodeClick={onOpenPoint}
          textureMapAsset={course.planet_texture_map_asset}
          progressPercent={course.progress ?? 0}
        />
      ) : null}
    </div>
  )
}
