'use client'

import { useCallback } from 'react'
import type { Dispatch, SetStateAction } from 'react'
import type { CourseAggregate } from './explorerPlanTypes'
import {
  type PendingExplorerPlanChanges,
  readErrorMessage,
  clonePendingChanges,
  nowIso,
} from './explorerPlanInternalUtils'

export function useExplorerPlanCourse({
  courseAggregate,
  courseDraftId,
  isMutating,
  setCourseAggregate,
  setPendingChanges,
  setIsMutating,
  setMutationMessage,
  setSelectedRegionId,
  setSelectedSubRegionId,
  setSelectedNodeId,
}: {
  courseAggregate: CourseAggregate | null
  courseDraftId: string
  isMutating: boolean
  setCourseAggregate: Dispatch<SetStateAction<CourseAggregate | null>>
  setPendingChanges: Dispatch<SetStateAction<PendingExplorerPlanChanges>>
  setIsMutating: Dispatch<SetStateAction<boolean>>
  setMutationMessage: Dispatch<SetStateAction<string | null>>
  setSelectedRegionId: Dispatch<SetStateAction<string | null>>
  setSelectedSubRegionId: Dispatch<SetStateAction<string | null>>
  setSelectedNodeId: Dispatch<SetStateAction<string | null>>
}) {
  const handleSelectCourse = useCallback(() => {
    setSelectedRegionId(null)
    setSelectedSubRegionId(null)
    setSelectedNodeId(null)
  }, [setSelectedRegionId, setSelectedSubRegionId, setSelectedNodeId])

  const updateCourseTitle = useCallback(async (title: string) => {
    const trimmed = title.trim()
    if (!trimmed || isMutating) return
    setMutationMessage(null)
    setCourseAggregate((current) => current ? { ...current, title: trimmed, updated_at: nowIso() } : current)
    setPendingChanges((current) => {
      const next = clonePendingChanges(current)
      next.updatedCourseTitle = true
      return next
    })
    handleSelectCourse()
    setMutationMessage('행성명을 화면에 반영했습니다. 탐험계획 저장을 누르면 최종 반영됩니다.')
  }, [handleSelectCourse, isMutating, setCourseAggregate, setPendingChanges, setMutationMessage])

  const deleteCourseDraft = useCallback(async (confirmTitle: string) => {
    if (!courseAggregate || !courseDraftId || isMutating) return
    if (courseAggregate.status === 'archived') {
      setMutationMessage('이미 비활성화된 코스입니다.')
      return
    }

    const trimmed = confirmTitle.trim()
    if (!trimmed) return

    setIsMutating(true)
    setMutationMessage(null)
    try {
      const res = await fetch(`/api/v1/course-drafts/${courseDraftId}`, {
        method: 'DELETE',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ confirm_title: trimmed }),
      })
      if (!res.ok) throw new Error(await readErrorMessage(res, '코스를 삭제하지 못했습니다.'))
      window.location.assign('/dashboard')
    } catch (err) {
      setMutationMessage(err instanceof Error ? err.message : '코스 비활성화 실패')
      setIsMutating(false)
    }
  }, [courseAggregate, courseDraftId, isMutating, setIsMutating, setMutationMessage])

  const activateCourseDraft = useCallback(async (confirmTitle: string) => {
    if (!courseAggregate || !courseDraftId || isMutating) return
    if (courseAggregate.status !== 'archived') {
      setMutationMessage('이미 활성화된 코스입니다.')
      return
    }

    const trimmed = confirmTitle.trim()
    if (!trimmed) return

    setIsMutating(true)
    setMutationMessage(null)
    try {
      const res = await fetch(`/api/v1/course-drafts/${courseDraftId}/activate`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ confirm_title: trimmed }),
      })
      const payload = (await res.json().catch(() => ({}))) as { destination?: string; error?: string }
      if (!res.ok) {
        throw new Error(payload.error ?? '코스를 활성화하지 못했습니다.')
      }
      window.location.assign(payload.destination ?? `/dashboard/course-drafts/${courseDraftId}?section=planning`)
    } catch (err) {
      setMutationMessage(err instanceof Error ? err.message : '코스 활성화 실패')
      setIsMutating(false)
    }
  }, [courseAggregate, courseDraftId, isMutating, setIsMutating, setMutationMessage])

  return { handleSelectCourse, updateCourseTitle, deleteCourseDraft, activateCourseDraft }
}
