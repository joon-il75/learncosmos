'use client'

import { useEffect, useState } from 'react'
import type { CourseAggregate } from '@/components/explorer-plan/explorerPlanTypes'
import type { DiaryCourseAggregate, DiaryNode } from './planetDiaryTypes'
import type { DiaryPointStatusCopy } from '@/lib/i18n/pages/planetDiary'

const MAP_BREAKPOINT = 1024

export function useResponsiveMapExpanded() {
  const [mapExpanded, setMapExpanded] = useState(() =>
    typeof window !== 'undefined' ? window.innerWidth >= MAP_BREAKPOINT : true,
  )
  const [isNarrowViewport, setIsNarrowViewport] = useState(() =>
    typeof window !== 'undefined' ? window.innerWidth < MAP_BREAKPOINT : false,
  )

  useEffect(() => {
    if (typeof window === 'undefined') return
    const mq = window.matchMedia(`(min-width: ${MAP_BREAKPOINT}px)`)
    const handler = (e: MediaQueryListEvent) => {
      setMapExpanded(e.matches)
      setIsNarrowViewport(!e.matches)
    }
    mq.addEventListener('change', handler)
    setMapExpanded(mq.matches)
    setIsNarrowViewport(!mq.matches)
    return () => mq.removeEventListener('change', handler)
  }, [])

  return { mapExpanded, setMapExpanded, isNarrowViewport }
}

export function useResponsiveCanvasScale(baseWidth: number, baseHeight: number) {
  const [scale, setScale] = useState(1)
  const [frameNode, setFrameNode] = useState<HTMLDivElement | null>(null)

  useEffect(() => {
    if (!frameNode) return
    const updateScale = (width: number, height: number) => {
      if (width <= 0 || height <= 0) return
      setScale(Math.min(width / baseWidth, height / baseHeight))
    }
    updateScale(frameNode.clientWidth, frameNode.clientHeight)
    const observer = new ResizeObserver((entries) => {
      const entry = entries[0]
      if (!entry) return
      updateScale(entry.contentRect.width, entry.contentRect.height)
    })
    observer.observe(frameNode)
    return () => observer.disconnect()
  }, [baseHeight, baseWidth, frameNode])

  return { setFrameNode, scale }
}

export function getStatusSummary(status: string, copy: DiaryPointStatusCopy) {
  if (status === 'completed') return copy.completed
  if (status === 'learning') return copy.learning
  if (status === 'ready') return copy.ready
  return copy.draft
}

export function findDiaryNodeById(course: DiaryCourseAggregate, nodeId: string) {
  for (const region of course.regions) {
    for (const node of region.nodes) {
      if (node.id === nodeId) return { node, regionName: region.region.name, subRegionName: null as string | null }
    }
    for (const subRegion of region.subregions) {
      for (const node of subRegion.nodes) {
        if (node.id === nodeId) {
          return { node, regionName: region.region.name, subRegionName: subRegion.subregion.name }
        }
      }
    }
  }
  return null
}

function buildDiaryStatusLookup(course: DiaryCourseAggregate) {
  const statusById = new Map<string, DiaryNode['learning_status']>()
  const statusByKey = new Map<string, DiaryNode['learning_status']>()
  const courseIDById = new Map<string, string>()
  const courseIDByKey = new Map<string, string>()
  const makeKey = (regionName: string, subRegionName: string | null, nodeTitle: string, nodeType: string) =>
    [regionName.trim().toLowerCase(), subRegionName?.trim().toLowerCase() ?? '', nodeTitle.trim().toLowerCase(), nodeType].join('::')

  for (const region of course.regions) {
    for (const node of region.nodes) {
      const key = makeKey(region.region.name, null, node.title, node.node_type)
      statusById.set(node.id, node.learning_status)
      statusByKey.set(key, node.learning_status)
      if (node.course_id) {
        courseIDById.set(node.id, node.course_id)
        courseIDByKey.set(key, node.course_id)
      }
    }
    for (const subRegion of region.subregions) {
      for (const node of subRegion.nodes) {
        const key = makeKey(region.region.name, subRegion.subregion.name, node.title, node.node_type)
        statusById.set(node.id, node.learning_status)
        statusByKey.set(key, node.learning_status)
        if (node.course_id) {
          courseIDById.set(node.id, node.course_id)
          courseIDByKey.set(key, node.course_id)
        }
      }
    }
  }

  return { statusById, statusByKey, courseIDById, courseIDByKey }
}

export function mergeDiaryCourseWithExplorer(
  baseCourse: DiaryCourseAggregate,
  explorerCourse: CourseAggregate | null,
): DiaryCourseAggregate {
  if (!explorerCourse) return baseCourse

  const statusLookup = buildDiaryStatusLookup(baseCourse)
  const makeKey = (regionName: string, subRegionName: string | null, nodeTitle: string, nodeType: string) =>
    [regionName.trim().toLowerCase(), subRegionName?.trim().toLowerCase() ?? '', nodeTitle.trim().toLowerCase(), nodeType].join('::')
  const resolveStatus = (node: { id: string; draft_point_id?: string | null; title: string; node_type: string }, regionName: string, subRegionName: string | null) =>
    statusLookup.statusById.get(node.draft_point_id ?? '') ??
    statusLookup.statusById.get(node.id) ??
    statusLookup.statusByKey.get(makeKey(regionName, subRegionName, node.title, node.node_type)) ??
    'ready'
  const resolveCourseID = (node: { id: string; draft_point_id?: string | null; title: string; node_type: string }, regionName: string, subRegionName: string | null) =>
    statusLookup.courseIDById.get(node.draft_point_id ?? '') ??
    statusLookup.courseIDById.get(node.id) ??
    statusLookup.courseIDByKey.get(makeKey(regionName, subRegionName, node.title, node.node_type)) ??
    null

  return {
    ...baseCourse,
    regions: explorerCourse.regions.map((region) => ({
      region: region.region,
      nodes: region.nodes.map((node) => ({
        ...node,
        id: node.draft_point_id ?? node.id,
        course_id: resolveCourseID(node, region.region.name, null),
        learning_status: resolveStatus(node, region.region.name, null),
      })),
      subregions: region.subregions.map((subRegion) => ({
        subregion: subRegion.subregion,
        nodes: subRegion.nodes.map((node) => ({
          ...node,
          id: node.draft_point_id ?? node.id,
          course_id: resolveCourseID(node, region.region.name, subRegion.subregion.name),
          learning_status: resolveStatus(node, region.region.name, subRegion.subregion.name),
        })),
      })),
    })),
  }
}
