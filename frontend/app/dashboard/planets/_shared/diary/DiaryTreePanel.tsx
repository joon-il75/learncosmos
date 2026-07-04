'use client'

import { useState } from 'react'
import { ExplorerTreeCourseRow } from '@/components/explorer-plan/tree/ExplorerTreeCourseRow'
import { ExplorerTreeRegionRow } from '@/components/explorer-plan/tree/ExplorerTreeRegionRow'
import { ExplorerTreeSubRegionRow } from '@/components/explorer-plan/tree/ExplorerTreeSubRegionRow'
import { ExplorerTreeNodeRow } from '@/components/explorer-plan/tree/ExplorerTreeNodeRow'
import {
  compactPanelStyle,
  panelStyle,
  scrollAreaStyle,
} from '@/components/explorer-plan/tree/explorerTreeStyles'
import type { PendingItemState } from '@/components/explorer-plan/useExplorerPlan'
import type { PlanetDiaryCopy } from '@/lib/i18n/pages/planetDiary'
import type { DiaryCourseAggregate } from './planetDiaryTypes'

export function DiaryTreePanel({
  course,
  selectedRegionId,
  selectedSubRegionId,
  canCompleteCourse = false,
  isCompletingCourse = false,
  getRegionPendingState,
  getSubRegionPendingState,
  getNodePendingState,
  onOpenCourseCompleteConfirm,
  onSelectCourse,
  onSelectRegion,
  onSelectSubRegion,
  onOpenPoint,
  copy,
  onboardingFirstRegionId = null,
  onboardingPointId = null,
  onboardingLockToFirstRegion = false,
  onboardingLockToPoint = false,
  onboardingDisableTreeInteractions = false,
}: {
  course: DiaryCourseAggregate
  selectedRegionId: string | null
  selectedSubRegionId: string | null
  canCompleteCourse?: boolean
  isCompletingCourse?: boolean
  getRegionPendingState: (regionId: string) => PendingItemState
  getSubRegionPendingState: (subRegionId: string) => PendingItemState
  getNodePendingState: (nodeId: string) => PendingItemState
  onOpenCourseCompleteConfirm?: () => void
  onSelectCourse: () => void
  onSelectRegion: (regionId: string) => void
  onSelectSubRegion: (regionId: string, subRegionId: string) => void
  onOpenPoint: (nodeId: string) => void
  copy: PlanetDiaryCopy['tree']
  onboardingFirstRegionId?: string | null
  onboardingPointId?: string | null
  onboardingLockToFirstRegion?: boolean
  onboardingLockToPoint?: boolean
  onboardingDisableTreeInteractions?: boolean
}) {
  const [collapsedRegions, setCollapsedRegions] = useState<Set<string>>(new Set())
  const [collapsedSubRegions, setCollapsedSubRegions] = useState<Set<string>>(new Set())
  const isCourseSelected = selectedRegionId == null && selectedSubRegionId == null
  const displayTitle = course.title?.trim() ? copy.courseTitle(course.title) : copy.emptyCourseTitle
  const displayPlanetName = course.title?.trim() || copy.unnamed

  const toggleRegion = (regionId: string) => {
    setCollapsedRegions((prev) => {
      const next = new Set(prev)
      if (next.has(regionId)) next.delete(regionId)
      else next.add(regionId)
      return next
    })
  }

  const toggleSubRegion = (subRegionId: string) => {
    setCollapsedSubRegions((prev) => {
      const next = new Set(prev)
      if (next.has(subRegionId)) next.delete(subRegionId)
      else next.add(subRegionId)
      return next
    })
  }

  return (
    <div style={{ ...panelStyle, ...compactPanelStyle, gap: 0, height: '100%' }}>
      <ExplorerTreeCourseRow
        displayTitle={displayTitle}
        displayPlanetName={displayPlanetName}
        isSelected={isCourseSelected}
        readOnly
        onSelect={onboardingDisableTreeInteractions ? () => {} : onSelectCourse}
      />
      {canCompleteCourse ? (
        <div style={courseCompleteActionWrapStyle}>
          <button
            type="button"
            style={{
              ...courseCompleteButtonStyle,
              ...(isCompletingCourse ? courseCompleteButtonDisabledStyle : null),
            }}
            disabled={isCompletingCourse}
            onClick={onOpenCourseCompleteConfirm}
          >
            {isCompletingCourse ? copy.completing : copy.complete}
          </button>
        </div>
      ) : null}
      <div style={{ ...scrollAreaStyle, padding: '4px 10px 8px 0' }}>
        {course.regions.map((region, index) => {
          const isSelected = selectedRegionId === region.region.id && selectedSubRegionId == null
          const isCollapsed = collapsedRegions.has(region.region.id)
          const regionPendingState = getRegionPendingState(region.region.id)
          const regionIsPendingDeleted = regionPendingState === 'deleted'
          const isOnboardingHighlighted = onboardingFirstRegionId === region.region.id
          const regionHasOnboardingPoint =
            onboardingPointId !== null &&
            (region.nodes.some((node) => node.id === onboardingPointId) ||
              region.subregions.some((subRegion) => subRegion.nodes.some((node) => node.id === onboardingPointId)))
          const isOnboardingDisabled =
            onboardingDisableTreeInteractions ||
            (onboardingLockToPoint && !regionHasOnboardingPoint) ||
            (onboardingLockToFirstRegion && !isOnboardingHighlighted)
          return (
            <div key={region.region.id}>
              <ExplorerTreeRegionRow
                regionAgg={region}
                regionIndex={index}
                isOpen={!isCollapsed}
                isSelected={isSelected}
                readOnly
                compact
                pendingState={regionPendingState}
                titleInput={region.region.name}
                onTitleChange={() => {}}
                isEditingDisabled
                onToggle={toggleRegion}
                onSelect={(regionId) => {
                  if (isOnboardingDisabled) return
                  onSelectRegion(regionId)
                }}
                isOnboardingHighlighted={isOnboardingHighlighted}
                isOnboardingDisabled={isOnboardingDisabled}
              />
              {!isCollapsed ? (
                <>
                  {region.nodes.map((node) => {
                    const isOnboardingPoint = onboardingPointId === node.id
                    return (
                      <ExplorerTreeNodeRow
                        key={node.id}
                        node={node}
                        depth={1}
                        isSelected={false}
                        readOnly
                        compact
                        pendingState={regionIsPendingDeleted ? 'deleted' : getNodePendingState(node.id)}
                        titleInput={node.title}
                        onTitleChange={() => {}}
                        isEditingDisabled
                        onSelect={(nodeId) => {
                          if (onboardingLockToPoint && !isOnboardingPoint) return
                          onOpenPoint(nodeId)
                        }}
                        isOnboardingHighlighted={isOnboardingPoint}
                        isOnboardingDisabled={onboardingLockToPoint && !isOnboardingPoint}
                      />
                    )
                  })}
                  {region.subregions.map((subRegion, subIndex) => {
                    const subSelected = selectedSubRegionId === subRegion.subregion.id
                    const subCollapsed = collapsedSubRegions.has(subRegion.subregion.id)
                    const directSubPendingState = getSubRegionPendingState(subRegion.subregion.id)
                    const subPendingState: PendingItemState = regionIsPendingDeleted ? 'deleted' : directSubPendingState
                    return (
                      <div key={subRegion.subregion.id}>
                        <ExplorerTreeSubRegionRow
                          subAgg={subRegion}
                          subIndex={subIndex}
                          regionId={region.region.id}
                          isOpen={!subCollapsed}
                          isSelected={subSelected}
                          readOnly
                          compact
                          pendingState={subPendingState}
                          titleInput={subRegion.subregion.name}
                          onTitleChange={() => {}}
                          isEditingDisabled
                          onToggle={toggleSubRegion}
                          onSelect={onSelectSubRegion}
                        />
                        {!subCollapsed
                          ? subRegion.nodes.map((node) => {
                              const isOnboardingPoint = onboardingPointId === node.id
                              return (
                                <ExplorerTreeNodeRow
                                  key={node.id}
                                  node={node}
                                  depth={2}
                                  isSelected={false}
                                  readOnly
                                  compact
                                  pendingState={subPendingState === 'deleted' ? 'deleted' : getNodePendingState(node.id)}
                                  titleInput={node.title}
                                  onTitleChange={() => {}}
                                  isEditingDisabled
                                  onSelect={(nodeId) => {
                                    if (onboardingLockToPoint && !isOnboardingPoint) return
                                    onOpenPoint(nodeId)
                                  }}
                                  isOnboardingHighlighted={isOnboardingPoint}
                                  isOnboardingDisabled={onboardingLockToPoint && !isOnboardingPoint}
                                />
                              )
                            })
                          : null}
                      </div>
                    )
                  })}
                </>
              ) : null}
            </div>
          )
        })}
      </div>
    </div>
  )
}

const courseCompleteActionWrapStyle = {
  padding: '8px 10px 7px 0',
  borderBottom: '1px solid rgba(126, 88, 36, 0.16)',
} as const

const courseCompleteButtonStyle = {
  width: '100%',
  minHeight: 36,
  border: '1px solid rgba(132, 82, 10, 0.44)',
  borderRadius: 10,
  background: 'linear-gradient(135deg, #ffe08a 0%, #d79a20 54%, #95620a 100%)',
  color: '#2c1700',
  fontSize: 13,
  fontWeight: 900,
  cursor: 'pointer',
  boxShadow: '0 8px 18px rgba(112, 74, 10, 0.22), inset 0 1px 0 rgba(255,255,255,0.46)',
} as const

const courseCompleteButtonDisabledStyle = {
  opacity: 0.72,
  cursor: 'progress',
} as const
