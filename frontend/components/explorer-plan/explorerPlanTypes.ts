// BE model.go 기준 타입 정의 — /api/v1/explorer/* 전용

export type NodeType = 'exploration' | 'research'
export type SourceType = 'youtube' | 'web' | 'creator' | 'internal'
export type ResearchType = 'concept' | 'practice' | 'problem' | 'free'
export type LayoutType = 'basic' | 'two-column' | 'note-card'
export type ParentKind = 'region' | 'subregion'
export type ItemStatus = 'active' | 'inactive'

export interface ExplorerRegion {
  id: string
  course_draft_id: string
  course_draft_lesson_id?: string | null
  name: string
  description: string | null
  order_index: number
  status: ItemStatus
  created_at: string
  updated_at: string
}

export interface ExplorerSubRegion {
  id: string
  region_id: string
  course_draft_lesson_id?: string | null
  name: string
  description: string | null
  order_index: number
  status: ItemStatus
  created_at: string
  updated_at: string
}

export interface ExplorerNode {
  id: string
  parent_kind: ParentKind
  parent_id: string
  node_type: NodeType
  draft_point_id?: string | null
  course_id?: string | null
  learning_status?: string | null
  title: string
  order_index: number
  status: ItemStatus
  created_at: string
  updated_at: string
  source_type?: SourceType
  source_url?: string | null
  content_id?: string | null
  summary?: string | null
  research_type?: ResearchType | null
  layout_type?: LayoutType | null
  block_count: number
}

export interface SubRegionAggregate {
  subregion: ExplorerSubRegion
  nodes: ExplorerNode[]
}

export interface RegionAggregate {
  region: ExplorerRegion
  subregions: SubRegionAggregate[]
  nodes: ExplorerNode[]
}

export interface CourseAggregate {
  course_draft_id: string
  title: string
  status: 'draft' | 'confirmed' | 'learning' | 'archived'
  is_inactive?: boolean
  progress?: number | null
  planet_texture_map_id?: string | null
  planet_texture_map_name?: string | null
  planet_texture_map_asset?: string | null
  planet_texture_map_rotation_duration_seconds?: number | null
  planet_texture_map_rotation_direction?: 'left' | 'right' | null
  updated_at: string
  regions: RegionAggregate[]
}
