import type { Candidate, DebugRun, RecommendationCompareLesson, RecommendationComparison } from './types'

export const parseScenarioJSON = <T,>(raw?: string | null): T | null => {
  const text = String(raw || '').trim()
  if (!text) return null
  try {
    return JSON.parse(text) as T
  } catch {
    return null
  }
}

export const buildStoredRecommendationComparison = (runs: DebugRun[]): RecommendationComparison | null => {
  const compareRuns = runs.filter(run => run.run_type === 'recommendation_compare')
  const parsedRuns = compareRuns
    .map(run => {
      const comparison = parseScenarioJSON<RecommendationComparison>(run.baseline_result_snapshot)
      if (!comparison) return null
      return { run, comparison }
    })
    .filter((item): item is { run: DebugRun; comparison: RecommendationComparison } => Boolean(item))
  if (parsedRuns.length === 0) return null

  const base = parsedRuns[0]
  const lessons: RecommendationCompareLesson[] = []
  const seenLessonIDs = new Set<string>()
  for (const item of parsedRuns) {
    for (const lesson of item.comparison.lessons ?? []) {
      const key = lesson.lesson.lesson_id || `${lesson.lesson.order_index}:${lesson.lesson.title}`
      if (seenLessonIDs.has(key)) continue
      seenLessonIDs.add(key)
      lessons.push({ ...lesson, run_id: item.run.id })
    }
  }

  return {
    ...base.comparison,
    run_id: base.run.id,
    lessons,
    lesson_count: lessons.length,
    total_candidates: lessons.reduce((sum, item) => sum + item.baseline.candidates.length, 0),
    total_lexical_hits: lessons.reduce((sum, item) => sum + item.baseline.stages.search.lexical_hit_count, 0),
    total_vector_hits: lessons.reduce((sum, item) => sum + item.baseline.stages.search.vector_hit_count, 0),
  }
}

export const candidateLabelKey = (lessonId: string, candidate: Candidate, candidateIndex: number) => (
  `${lessonId}:${candidate.content_id || candidate.external_url || candidate.title}:${candidateIndex + 1}`
)

export const candidateFeatureSnapshotForLabel = (
  lesson: RecommendationCompareLesson['lesson'],
  candidate: Candidate,
  candidateIndex: number,
) => ({
  ...(candidate.feature_snapshot ?? {}),
  lesson_id: lesson.lesson_id,
  lesson_title: lesson.title,
  label_candidate_index: candidateIndex + 1,
  ranker_score: candidate.ranker_score ?? null,
  ranker_route_rank: candidate.ranker_route_rank ?? candidateIndex + 1,
  ranker_rerank_rank: candidate.ranker_rerank_rank ?? candidate.ranker_rank ?? candidateIndex + 1,
  ranker_rank_delta: candidate.ranker_rank_delta ?? 0,
  ranker_reason: candidate.ranker_reason ?? '',
  ranker_provider: candidate.ranker_provider ?? '',
  ranker_model_version: candidate.ranker_model_version ?? '',
})

export const recommendationQueriesByLesson = (comparison: RecommendationComparison | null): Record<string, string> => {
  const entries = (comparison?.lessons ?? [])
    .map(item => {
      const lessonId = item.lesson.lesson_id
      const recommendationQuery = item.lesson.recommendation_query ?? item.baseline.stages.query.recommendation_query ?? ''
      return [lessonId, recommendationQuery.trim()] as const
    })
    .filter((entry): entry is readonly [string, string] => Boolean(entry[0] && entry[1]))
  return Object.fromEntries(entries)
}
