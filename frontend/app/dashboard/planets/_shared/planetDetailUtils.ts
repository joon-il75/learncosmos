import type {
  DiaryCourseAggregate,
  DiaryNode,
  DiaryRegionAggregate,
  DiarySubRegionAggregate,
} from './diary/planetDiaryTypes'

// ── 공개 타입 ──────────────────────────────────────────────────────────────────

export type PlanetRouteKind = 'learning' | 'shared'
export type DiaryStageTab = 'journal' | 'records' | 'results'
export type PlanetStatus = 'ready' | 'learning' | 'completed'

export interface PlanetListItem {
  id: string
  draft_id: string
  title: string
  status: string
}

export interface UserInfo {
  id: string
  email: string
  display_id: string
  nickname: string
  required_consent_pending: boolean
  ui_locale?: 'ko' | 'en'
}

export interface PlanetLesson {
  lesson: {
    id: string
    title: string
    summary?: string | null
    difficulty_level?: string | null
    lesson_role: string
    status?: string | null
  }
  resources: {
    id: string
    title: string
    resource_type: string
    external_url?: string | null
  }[]
  points?: PlanetPoint[]
  runtime_entry?: {
    status: 'in_progress' | 'completed'
    started_at: string
    completed_at?: string | null
    updated_at: string
  } | null
}

export interface PlanetPoint {
  point: {
    id: string
    course_id?: string | null
    point_type: 'exploration' | 'research'
    status: 'draft' | 'ready' | 'learning' | 'completed'
    title: string
    description?: string | null
    external_url?: string | null
    thumbnail_url?: string | null
    order_index: number
    item_status?: string | null
    explorer_status?: string | null
  }
}

export interface PlanetLessonTree {
  lesson: {
    id: string
    title: string
    summary?: string | null
    difficulty_level?: string | null
    lesson_role: string
    parent_lesson_id?: string | null
    status?: string | null
  }
  points?: PlanetPoint[]
  sub_lessons?: PlanetLessonTree[]
  runtime_entry?: PlanetLesson['runtime_entry']
}

export interface PlanetLevel {
  level: {
    id: string
    title: string
    objective?: string | null
  }
  lessons: PlanetLesson[]
}

export interface PlanetAggregate {
  planet: {
    id: string
    draft_id: string
    title: string
    status: PlanetStatus
    planet_type_id?: string | null
    planet_type_name?: string | null
    planet_type_asset?: string | null
    planet_texture_map_id?: string | null
    planet_texture_map_name?: string | null
    planet_texture_map_asset?: string | null
    planet_texture_map_rotation_duration_seconds?: number | null
    planet_texture_map_rotation_direction?: 'left' | 'right' | null
    lesson_count: number
    completed_lesson_count: number
    progress?: number | null
    can_complete?: boolean
    updated_at: string
    share_count?: number
  }
  goal_context?: {
    learning_goal?: string | null
    confirmed_goal?: string | null
    usage_context?: string | null
    motivation?: string | null
    goal_profile_id?: string | null
    goal_profile_version?: number | null
    goal_readiness?: number | null
  } | null
  levels: PlanetLevel[]
  lessons?: PlanetLessonTree[]
}

export interface ExplorerCourseAggregate {
  regions: {
    region: {
      id: string
      status: 'active' | 'inactive'
    }
    nodes: ExplorerNode[]
    subregions: {
      subregion: {
        id: string
        status: 'active' | 'inactive'
      }
      nodes: ExplorerNode[]
    }[]
  }[]
}

export interface ExplorerNode {
  id: string
  parent_kind: 'region' | 'subregion'
  parent_id: string
  node_type: 'exploration' | 'research'
  draft_point_id?: string | null
  status: 'active' | 'inactive'
}

export interface DiaryContextDraft {
  draft: {
    id: string
    source_query: string
    title: string
    status?: 'draft' | 'confirmed' | 'learning' | 'archived'
    description?: string | null
  }
}

// ── 유틸리티 함수 ──────────────────────────────────────────────────────────────

export function isActiveStructureItem(status?: string | null) {
  return status !== 'inactive'
}

export function getStatusLabel(status: PlanetStatus): string {
  if (status === 'ready') return '학습준비중'
  if (status === 'learning') return '학습중'
  return '공유중'
}

export function getJournalStatusLabel(status: PlanetStatus): string {
  if (status === 'ready') return '학습준비중(ready)'
  if (status === 'learning') return '학습중(learning)'
  return '탐험완료(completed)'
}

export function buildPlanetResourceFromPoint(point: PlanetPoint): PlanetLesson['resources'][number] {
  return {
    id: point.point.id,
    title: point.point.title,
    resource_type: point.point.point_type,
    external_url: point.point.external_url,
  }
}

export function buildPlanetLevelsFromLessons(lessons?: PlanetLessonTree[]): PlanetLevel[] {
  if (!lessons || lessons.length === 0) return []
  return lessons.map((mainLesson) => ({
    level: {
      id: mainLesson.lesson.id,
      title: mainLesson.lesson.title,
      objective: mainLesson.lesson.summary,
    },
    lessons: (mainLesson.sub_lessons ?? []).map((subLesson) => ({
      lesson: {
        id: subLesson.lesson.id,
        title: subLesson.lesson.title,
        summary: subLesson.lesson.summary,
        difficulty_level: subLesson.lesson.difficulty_level,
        lesson_role: subLesson.lesson.lesson_role,
        status: subLesson.lesson.status,
      },
      resources: (subLesson.points ?? [])
        .filter((point) => point.point.point_type === 'exploration')
        .map(buildPlanetResourceFromPoint),
      points: subLesson.points ?? [],
      runtime_entry: subLesson.runtime_entry,
    })),
  }))
}

export function collectPlanetPointStatuses(planet: PlanetAggregate): Map<string, PlanetPoint['point']['status']> {
  const statuses = new Map<string, PlanetPoint['point']['status']>()

  const addPoint = (point: PlanetPoint) => {
    statuses.set(point.point.id, point.point.status)
  }

  const visitLessonTree = (lesson: PlanetLessonTree) => {
    ;(lesson.points ?? []).forEach(addPoint)
    ;(lesson.sub_lessons ?? []).forEach(visitLessonTree)
  }

  if (planet.lessons?.length) {
    planet.lessons.forEach(visitLessonTree)
    return statuses
  }

  planet.levels.forEach((level) => {
    level.lessons.forEach((lesson) => {
      ;(lesson.points ?? []).forEach(addPoint)
    })
  })

  return statuses
}

export function countProgressFromExplorerCourse(
  explorerCourse: ExplorerCourseAggregate,
  pointStatuses: Map<string, PlanetPoint['point']['status']>,
): {
  pointProgress: { total: number; completed: number }
  lessonProgress: { total: number; completed: number; learning: number }
} {
  let totalPoints = 0
  let completedPoints = 0
  let totalLessons = 0
  let completedLessons = 0
  let learningLessons = 0

  const countRegion = (nodes: ExplorerNode[]) => {
    const activeNodes = nodes.filter((node) => node.status === 'active')
    totalLessons += 1
    totalPoints += activeNodes.length

    const containerCompletedPoints = activeNodes.filter((node) => {
      const status = pointStatuses.get(node.draft_point_id ?? node.id)
      return status === 'completed'
    }).length

    completedPoints += containerCompletedPoints
    if (activeNodes.length > 0 && containerCompletedPoints === activeNodes.length) {
      completedLessons += 1
      return
    }
    if (
      containerCompletedPoints > 0 ||
      activeNodes.some((node) => pointStatuses.get(node.draft_point_id ?? node.id) === 'learning')
    ) {
      learningLessons += 1
    }
  }

  for (const regionAgg of explorerCourse.regions) {
    if (regionAgg.region.status !== 'active') continue
    countRegion([
      ...regionAgg.nodes,
      ...regionAgg.subregions
        .filter((subAgg) => subAgg.subregion.status === 'active')
        .flatMap((subAgg) => subAgg.nodes),
    ])
  }

  return {
    pointProgress: { total: totalPoints, completed: completedPoints },
    lessonProgress: { total: totalLessons, completed: completedLessons, learning: learningLessons },
  }
}

export function countPlanetPointsFromLevels(levels: PlanetLevel[]): { total: number; completed: number } {
  let total = 0
  let completed = 0
  for (const level of levels) {
    for (const lesson of level.lessons) {
      if (!isActiveStructureItem(lesson.lesson.status)) continue
      const points = (lesson.points ?? []).filter((point) =>
        isActiveStructureItem(point.point.item_status ?? point.point.explorer_status),
      )
      total += points.length
      completed += points.filter((point) => point.point.status === 'completed').length
    }
  }
  return { total, completed }
}

export function countPlanetPointsFromLessonTrees(lessons?: PlanetLessonTree[]): { total: number; completed: number } {
  let total = 0
  let completed = 0

  const visit = (lesson: PlanetLessonTree) => {
    if (!isActiveStructureItem(lesson.lesson.status)) return
    const points = (lesson.points ?? []).filter((point) =>
      isActiveStructureItem(point.point.item_status ?? point.point.explorer_status),
    )
    total += points.length
    completed += points.filter((point) => point.point.status === 'completed').length
    ;(lesson.sub_lessons ?? []).forEach(visit)
  }

  ;(lessons ?? []).forEach(visit)
  return { total, completed }
}

export function countLessonProgressFromLessonTrees(lessons?: PlanetLessonTree[]): { total: number; completed: number; learning: number } {
  let total = 0
  let completed = 0
  let learning = 0

  const countLesson = (lesson: PlanetLessonTree) => {
    if (!isActiveStructureItem(lesson.lesson.status)) return
    total += 1
    const points = (lesson.points ?? []).filter((point) =>
      isActiveStructureItem(point.point.item_status ?? point.point.explorer_status),
    )
    const completedPoints = points.filter((point) => point.point.status === 'completed').length
    if (points.length > 0 && completedPoints === points.length) completed += 1
    else if (completedPoints > 0 || points.some((point) => point.point.status === 'learning')) learning += 1
  }

  const visit = (lesson: PlanetLessonTree) => {
    countLesson(lesson)
    ;(lesson.sub_lessons ?? []).forEach(visit)
  }

  ;(lessons ?? []).forEach(visit)
  return { total, completed, learning }
}

export function countLessonProgressFromLevels(levels: PlanetLevel[]): { total: number; completed: number; learning: number } {
  let total = 0
  let completed = 0
  let learning = 0

  for (const level of levels) {
    for (const lesson of level.lessons) {
      if (!isActiveStructureItem(lesson.lesson.status)) continue
      total += 1
      const points = (lesson.points ?? []).filter((point) =>
        isActiveStructureItem(point.point.item_status ?? point.point.explorer_status),
      )
      const completedPoints = points.filter((point) => point.point.status === 'completed').length
      if (points.length > 0 && completedPoints === points.length) completed += 1
      else if (completedPoints > 0 || points.some((point) => point.point.status === 'learning')) learning += 1
    }
  }

  return { total, completed, learning }
}

export function buildDiaryNode(
  point: PlanetPoint,
  parentKind: 'region' | 'subregion',
  parentId: string,
): DiaryNode {
  return {
    id: point.point.id,
    course_id: point.point.course_id ?? null,
    parent_kind: parentKind,
    parent_id: parentId,
    node_type: point.point.point_type,
    title: point.point.title,
    order_index: point.point.order_index,
    status: 'active',
    created_at: '',
    updated_at: '',
    source_type: point.point.external_url ? 'web' : undefined,
    source_url: point.point.external_url ?? null,
    summary: point.point.description ?? null,
    block_count: 0,
    learning_status: point.point.status,
  }
}

export function buildDiaryCourseAggregate(planet: PlanetAggregate): DiaryCourseAggregate {
  const regions: DiaryRegionAggregate[] = []

  if (planet.lessons?.length) {
    planet.lessons.forEach((mainLesson, regionIndex) => {
      const regionId = mainLesson.lesson.id
      const regionNodes = (mainLesson.points ?? []).map((point) => buildDiaryNode(point, 'region', regionId))
      const subregions: DiarySubRegionAggregate[] = (mainLesson.sub_lessons ?? []).map((subLesson, subIndex) => {
        const subRegionId = subLesson.lesson.id
        return {
          subregion: {
            id: subRegionId,
            region_id: regionId,
            name: subLesson.lesson.title,
            description: subLesson.lesson.summary ?? null,
            order_index: subIndex,
            status: 'active',
            created_at: '',
            updated_at: '',
          },
          nodes: (subLesson.points ?? []).map((point) => buildDiaryNode(point, 'subregion', subRegionId)),
        }
      })

      regions.push({
        region: {
          id: regionId,
          course_draft_id: planet.planet.draft_id,
          name: mainLesson.lesson.title,
          description: mainLesson.lesson.summary ?? null,
          order_index: regionIndex,
          status: 'active',
          created_at: '',
          updated_at: '',
        },
        subregions,
        nodes: regionNodes,
      })
    })
  } else {
    planet.levels.forEach((level, regionIndex) => {
      const regionId = level.level.id
      const subregions: DiarySubRegionAggregate[] = level.lessons.map((lesson, subIndex) => {
        const subRegionId = lesson.lesson.id
        return {
          subregion: {
            id: subRegionId,
            region_id: regionId,
            name: lesson.lesson.title,
            description: lesson.lesson.summary ?? null,
            order_index: subIndex,
            status: 'active',
            created_at: '',
            updated_at: '',
          },
          nodes: (lesson.points ?? []).map((point) => buildDiaryNode(point, 'subregion', subRegionId)),
        }
      })

      regions.push({
        region: {
          id: regionId,
          course_draft_id: planet.planet.draft_id,
          name: level.level.title,
          description: level.level.objective ?? null,
          order_index: regionIndex,
          status: 'active',
          created_at: '',
          updated_at: '',
        },
        subregions,
        nodes: [],
      })
    })
  }

  return {
    course_draft_id: planet.planet.draft_id,
    title: planet.planet.title,
    status: planet.planet.status === 'completed' ? 'learning' : 'confirmed',
    progress: planet.planet.progress ?? null,
    planet_texture_map_id: planet.planet.planet_texture_map_id ?? null,
    planet_texture_map_name: planet.planet.planet_texture_map_name ?? null,
    planet_texture_map_asset: planet.planet.planet_texture_map_asset ?? null,
    planet_texture_map_rotation_duration_seconds: planet.planet.planet_texture_map_rotation_duration_seconds ?? null,
    planet_texture_map_rotation_direction: planet.planet.planet_texture_map_rotation_direction ?? null,
    updated_at: planet.planet.updated_at,
    regions,
  }
}
