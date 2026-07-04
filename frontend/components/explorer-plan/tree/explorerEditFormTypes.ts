'use client'

import type { SourceType } from '../explorerPlanTypes'
import type { ExplorationNodeSourceInput } from '../useExplorerPlan'

export type ExplorerContentCandidate = {
  id: string
  content_id?: string | null
  source_type?: SourceType | null
  content_type: string
  title: string
  description?: string | null
  thumbnail_url?: string | null
  url?: string | null
  canonical_url?: string | null
  author?: string | null
}

export type ExplorerRecommendationContext = {
  course_draft_id?: string
  region_title?: string
  region_description?: string
  subregion_title?: string
  subregion_description?: string
  node_title?: string
  node_summary?: string
  lesson_id?: string
  excluded_content_ids?: string[]
  excluded_urls?: string[]
}

export type DeleteConfirmState = {
  mode: 'confirm' | 'alert'
  action?: 'delete' | 'activate'
  target: 'course' | 'region' | 'subregion' | 'node'
  title: string
  message: string
  requireTitleInput: boolean
  confirmInput: string
} | null

export type { ExplorationNodeSourceInput }
