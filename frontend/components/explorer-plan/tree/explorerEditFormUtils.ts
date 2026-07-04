'use client'

import type { SourceType, ExplorerNode } from '../explorerPlanTypes'
import type { ExplorerContentCandidate, ExplorationNodeSourceInput } from './explorerEditFormTypes'

export function inferCandidateUrlSourceType(url: string): SourceType {
  const normalized = url.toLowerCase()
  return normalized.includes('youtube.com') || normalized.includes('youtu.be') ? 'youtube' : 'web'
}

export function candidateToExplorationSource(candidate: ExplorerContentCandidate): ExplorationNodeSourceInput | null {
  const contentId = candidate.content_id?.trim() ?? ''
  const url = (candidate.url ?? candidate.canonical_url ?? '').trim()
  if (contentId) {
    return {
      sourceType: candidate.source_type === 'creator' ? 'creator' : 'internal',
      sourceUrl: url || null,
      contentId,
    }
  }
  if (url) {
    return {
      sourceType: inferCandidateUrlSourceType(url),
      sourceUrl: url,
      contentId: null,
    }
  }
  return null
}

export function selectedNodeSource(node: ExplorerNode, urlInput: string): string | ExplorationNodeSourceInput | null {
  const url = urlInput.trim()
  if (url) return url
  if (!node.content_id) return null
  return {
    sourceType: node.source_type === 'creator' ? 'creator' : 'internal',
    sourceUrl: node.source_url ?? null,
    contentId: node.content_id,
  }
}

export function withTrimmedValue(value?: string | null): string | undefined {
  const trimmed = value?.trim()
  return trimmed ? trimmed : undefined
}
