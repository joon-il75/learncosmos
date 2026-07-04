'use client'

import { useCallback } from 'react'
import type { Dispatch, SetStateAction } from 'react'
import type { CourseAggregate } from './explorerPlanTypes'
import {
  type LegacyDraftAggregate,
  type PendingExplorerPlanChanges,
  emptyPendingChanges,
  normalizeCourseAggregateOrder,
  buildLegacyCourseAggregate,
  shouldUseLegacyDraftAggregate,
  pendingChangesForLegacyAggregate,
} from './explorerPlanInternalUtils'

type CourseAggregateResponse = CourseAggregate | { course: CourseAggregate }

function unwrapCourseAggregate(data: CourseAggregateResponse): CourseAggregate {
  return 'course' in data ? data.course : data
}

export function useExplorerPlanLoader({
  courseDraftId,
  legacyDraft,
  showInactiveItems,
  setCourseAggregate,
  setPendingChanges,
  setMutationMessage,
  setIsLoading,
  setError,
}: {
  courseDraftId: string
  legacyDraft?: LegacyDraftAggregate | null
  showInactiveItems: boolean
  setCourseAggregate: Dispatch<SetStateAction<CourseAggregate | null>>
  setPendingChanges: Dispatch<SetStateAction<PendingExplorerPlanChanges>>
  setMutationMessage: Dispatch<SetStateAction<string | null>>
  setIsLoading: Dispatch<SetStateAction<boolean>>
  setError: Dispatch<SetStateAction<string | null>>
}) {
  const load = useCallback(async () => {
    if (!courseDraftId) return
    setIsLoading(true)
    setError(null)
    try {
      const query = showInactiveItems ? '?include_inactive=true' : ''
      const res = await fetch(`/api/v1/explorer/course/${courseDraftId}${query}`, {
        credentials: 'include',
        cache: 'no-store',
      })
      if (!res.ok) {
        const data = await res.json().catch(() => ({})) as { error?: string }
        throw new Error(data.error ?? '탐험계획을 불러오지 못했습니다.')
      }
      const data = await res.json() as CourseAggregateResponse
      const explorerAggregate = unwrapCourseAggregate(data)
      if (legacyDraft && shouldUseLegacyDraftAggregate(explorerAggregate, legacyDraft)) {
        const legacyAggregate = normalizeCourseAggregateOrder(buildLegacyCourseAggregate(courseDraftId, legacyDraft))
        setCourseAggregate(legacyAggregate)
        setPendingChanges(pendingChangesForLegacyAggregate(legacyAggregate, explorerAggregate))
        setMutationMessage(null)
      } else {
        setCourseAggregate(normalizeCourseAggregateOrder(explorerAggregate))
        setPendingChanges(emptyPendingChanges())
        setMutationMessage(null)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '탐험계획 로드 실패')
    } finally {
      setIsLoading(false)
    }
  }, [courseDraftId, legacyDraft, showInactiveItems]) // eslint-disable-line react-hooks/exhaustive-deps

  return { load }
}
