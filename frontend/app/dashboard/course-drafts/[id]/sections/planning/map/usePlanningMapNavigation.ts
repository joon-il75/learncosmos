'use client'

import { useEffect, useMemo, useState } from 'react'
import type { SubRegionAggregate } from '@/components/explorer-plan/explorerPlanTypes'

const REGIONS_PER_PAGE = 3
const PAGE_TRANSITION_MS = 180

type NavDirection = 'prev' | 'next' | null
export type PlanningMapLevel = 'course' | 'region' | 'subregion'

type RegionTreeLike = {
  region: { id: string; name: string }
  subregions: SubRegionAggregate[]
}

type PlanningMapNavigationOptions = {
  mapLevel: PlanningMapLevel
  selectedRegionId: string | null
  selectedSubRegionId: string | null
  onMapLevelChange: (level: PlanningMapLevel) => void
  onSelectedRegionIdChange: (regionId: string | null) => void
  onSelectedSubRegionIdChange: (subRegionId: string | null) => void
}

export function usePlanningMapNavigation<T extends RegionTreeLike>(
  regions: T[],
  options: PlanningMapNavigationOptions,
) {
  const totalPages = Math.max(1, Math.ceil(regions.length / REGIONS_PER_PAGE))

  const [page, setPage] = useState(0)
  const [isTransitioning, setIsTransitioning] = useState(false)
  const [direction, setDirection] = useState<NavDirection>(null)
  const {
    mapLevel,
    selectedRegionId,
    selectedSubRegionId,
    onMapLevelChange,
    onSelectedRegionIdChange,
    onSelectedSubRegionIdChange,
  } = options

  const currentPage = Math.min(page, Math.max(0, totalPages - 1))

  const pagedRegions = useMemo(() => {
    return regions.slice(
      currentPage * REGIONS_PER_PAGE,
      currentPage * REGIONS_PER_PAGE + REGIONS_PER_PAGE
    )
  }, [currentPage, regions])

  const selectedRegion = useMemo(() => {
    if (selectedRegionId == null) return null
    return regions.find(r => r.region.id === selectedRegionId) ?? null
  }, [regions, selectedRegionId])

  const selectedSubRegion = useMemo(() => {
    if (!selectedRegion || selectedSubRegionId == null) return null
    return selectedRegion.subregions.find(s => s.subregion.id === selectedSubRegionId) ?? null
  }, [selectedRegion, selectedSubRegionId])

  useEffect(() => {
    setPage(prev => Math.min(prev, Math.max(0, totalPages - 1)))
  }, [totalPages])

  useEffect(() => {
    if (selectedRegionId == null) return
    const exists = regions.some(r => r.region.id === selectedRegionId)
    if (!exists) {
      onSelectedRegionIdChange(null)
      onSelectedSubRegionIdChange(null)
      onMapLevelChange('course')
    }
  }, [regions, selectedRegionId, onMapLevelChange, onSelectedRegionIdChange, onSelectedSubRegionIdChange])

  useEffect(() => {
    if (selectedRegionId == null || selectedSubRegionId == null) return
    const exists = selectedRegion?.subregions.some(s => s.subregion.id === selectedSubRegionId)
    if (!exists) {
      onSelectedSubRegionIdChange(null)
      onMapLevelChange('region')
    }
  }, [selectedRegion, selectedRegionId, selectedSubRegionId, onMapLevelChange, onSelectedSubRegionIdChange])

  useEffect(() => {
    if (mapLevel !== 'course' || selectedRegionId == null) return
    const regionIndex = regions.findIndex(r => r.region.id === selectedRegionId)
    if (regionIndex < 0) return
    setPage(Math.floor(regionIndex / REGIONS_PER_PAGE))
  }, [mapLevel, regions, selectedRegionId])

  const runTransition = (nextDirection: Exclude<NavDirection, null>, updater: () => void) => {
    if (isTransitioning || mapLevel !== 'course') return
    setDirection(nextDirection)
    setIsTransitioning(true)
    window.setTimeout(() => {
      updater()
      setIsTransitioning(false)
      setDirection(null)
    }, PAGE_TRANSITION_MS)
  }

  const goPrev = () => {
    if (currentPage === 0) return
    runTransition('prev', () => setPage(prev => Math.max(0, prev - 1)))
  }

  const goNext = () => {
    if (currentPage >= totalPages - 1) return
    runTransition('next', () => setPage(prev => Math.min(totalPages - 1, prev + 1)))
  }

  const handleRegionClick = (regionId: string) => {
    onSelectedRegionIdChange(regionId)
    onSelectedSubRegionIdChange(null)
    onMapLevelChange('region')
  }

  const handleSubRegionClick = (subRegionId: string) => {
    onSelectedSubRegionIdChange(subRegionId)
    onMapLevelChange('subregion')
  }

  const handleBackToCourseMap = () => {
    onSelectedRegionIdChange(null)
    onSelectedSubRegionIdChange(null)
    onMapLevelChange('course')
  }

  const handleBackToRegionMap = () => {
    onSelectedSubRegionIdChange(null)
    onMapLevelChange('region')
  }

  const regionIndexOffset = currentPage * REGIONS_PER_PAGE

  return {
    mapLevel,
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
    handleRegionClick,
    handleSubRegionClick,
    handleBackToCourseMap,
    handleBackToRegionMap,
  }
}
