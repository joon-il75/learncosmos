import type {
  CourseAggregate,
  ExplorerNode,
  RegionAggregate,
  SubRegionAggregate,
} from '@/components/explorer-plan/explorerPlanTypes'
import type { DiaryPointStatusCopy } from '@/lib/i18n/pages/planetDiary'

export type DiaryPointStatus = 'draft' | 'ready' | 'learning' | 'completed'

export type DiaryNode = ExplorerNode & {
  learning_status: DiaryPointStatus
}

export type DiarySubRegionAggregate = Omit<SubRegionAggregate, 'nodes'> & {
  nodes: DiaryNode[]
}

export type DiaryRegionAggregate = Omit<RegionAggregate, 'subregions' | 'nodes'> & {
  subregions: DiarySubRegionAggregate[]
  nodes: DiaryNode[]
}

export type DiaryCourseAggregate = Omit<CourseAggregate, 'regions'> & {
  regions: DiaryRegionAggregate[]
}

export function countDiaryNodes(nodes: DiaryNode[]): {
  total: number
  completed: number
  learning: number
} {
  return nodes.reduce(
    (acc, node) => {
      acc.total += 1
      if (node.learning_status === 'completed') acc.completed += 1
      if (node.learning_status === 'learning') acc.learning += 1
      return acc
    },
    { total: 0, completed: 0, learning: 0 },
  )
}

export function countDiaryRegion(region: DiaryRegionAggregate) {
  return countDiaryNodes([
    ...region.nodes,
    ...region.subregions.flatMap((subregion) => subregion.nodes),
  ])
}

export function countDiaryCourse(course: DiaryCourseAggregate) {
  return countDiaryNodes(
    course.regions.flatMap((region) => [
      ...region.nodes,
      ...region.subregions.flatMap((subregion) => subregion.nodes),
    ]),
  )
}

export function getDiaryNodeStatusLabel(status: DiaryPointStatus, copy: DiaryPointStatusCopy): string {
  if (status === 'completed') return copy.completed
  if (status === 'learning') return copy.learning
  if (status === 'ready') return copy.ready
  return copy.draft
}
