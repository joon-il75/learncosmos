'use client'

import { useCallback, useEffect, useState } from 'react'
import type { DashboardGoalCopy } from '@/lib/i18n/pages/dashboardGoal'
import { getDashboardGoalCopy } from '@/lib/i18n/pages/dashboardGoal'
import { resolveSafetyInputMessage, type SafetyAPIErrorPayload } from '@/lib/safetyErrors'

export type InterviewState = 'listening' | 'clarifying' | 'proposing_goal' | 'confirmed' | 'revising_goal' | 'awaiting_rebuild_decision'

export type RebuildDecision = 'keep_structure' | 'rebuild_remaining' | 'rebuild_all'

export type InterviewMessage = {
  role: 'lumi' | 'user'
  content: string
}

export type GoalProfile = {
  id: string
  course_draft_id: string | null
  user_intent: string
  motivation: string | null
  usage_context: string | null
  confirmed_goal: string | null
  summarized_context: string
  rebuild_decision: RebuildDecision | null
  interview_state: InterviewState
  messages: InterviewMessage[]
  version: number
  is_active: boolean
}

type UseGoalInterviewReturn = {
  profile: GoalProfile | null
  isLoading: boolean
  isSending: boolean
  error: string | null
  startInterview: (userIntent: string) => Promise<void>
  sendMessage: (message: string) => Promise<void>
  confirmGoal: (confirmedGoal: string) => Promise<GoalProfile | null>
  reviseGoal: (message: string) => Promise<void>
  cancelRevision: () => Promise<GoalProfile | null>
  applyRebuildDecision: (decision: RebuildDecision) => Promise<void>
  attachDraft: (draftId: string) => Promise<void>
  resetActiveGoal: () => Promise<void>
}


type GoalInterviewAPIResponse = {
  goal: GoalProfile
  job_id?: string
  poll_url?: string
  status?: string
  error_code?: string
}

type LLMJobAPIResponse = {
  status?: string
  error_code?: string | null
}

function isGoalInterviewPendingResponse(status: number, data: GoalInterviewAPIResponse): boolean {
  return status === 202 || data.error_code === 'goal_interview_pending'
}

async function waitForGoalInterviewJob(pollURL: string, fallbackMessage: string): Promise<void> {
  const startedAt = Date.now()
  while (Date.now() - startedAt < 45_000) {
    await new Promise((resolve) => setTimeout(resolve, 800))
    const res = await fetch(pollURL, {
      credentials: 'include',
      cache: 'no-store',
    })
    if (!res.ok) throw new Error(fallbackMessage)
    const data = await res.json() as LLMJobAPIResponse
    if (data.status === 'succeeded') return
    if (data.status === 'failed' || data.status === 'canceled' || data.status === 'expired') {
      throw new Error(fallbackMessage)
    }
  }
  throw new Error(fallbackMessage)
}

export function useGoalInterview(
  courseDraftId?: string,
  options?: { autoLoad?: boolean; errorCopy?: DashboardGoalCopy['hookErrors'] },
): UseGoalInterviewReturn {
  const [profile, setProfile] = useState<GoalProfile | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isSending, setIsSending] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const hasDraftScope = Boolean(courseDraftId)
  const autoLoad = options?.autoLoad ?? true
  const errorCopy = options?.errorCopy ?? getDashboardGoalCopy('ko').hookErrors

  const buildPath = useCallback((kind: 'load' | 'start' | 'interview' | 'confirm' | 'revise' | 'cancel-revision' | 'rebuild-decision') => {
    if (courseDraftId) {
      if (kind === 'load') return `/api/v1/course-drafts/${courseDraftId}/goal`
      return `/api/v1/course-drafts/${courseDraftId}/goal/${kind}`
    }
    if (kind === 'load') return '/api/v1/goals/active'
    return `/api/v1/goals/${kind}`
  }, [courseDraftId])

  const load = useCallback(async () => {
    setIsLoading(true)
    setError(null)
    try {
      const res = await fetch(buildPath('load'), {
        credentials: 'include',
        cache: 'no-store',
      })
      if (!res.ok) throw new Error(errorCopy.loadFailed)
      const data = await res.json() as { goal: GoalProfile | null }
      setProfile(data.goal)
    } catch (err) {
      setError(err instanceof Error ? err.message : errorCopy.loadFallback)
    } finally {
      setIsLoading(false)
    }
  }, [buildPath, errorCopy])

  useEffect(() => {
    if (!autoLoad) {
      setIsLoading(false)
      return
    }
    void load()
  }, [autoLoad, load])

  const startInterview = useCallback(async (userIntent: string) => {
    if (!userIntent.trim() || isSending) return
    setIsSending(true)
    setError(null)
    try {
      const res = await fetch(buildPath('start'), {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_intent: userIntent }),
      })
      if (!res.ok) {
        const d = await res.json().catch(() => ({})) as SafetyAPIErrorPayload
        throw new Error(resolveSafetyInputMessage(d, d.error ?? errorCopy.startFailed))
      }
      const data = await res.json() as GoalInterviewAPIResponse
      setProfile(data.goal)
      if (isGoalInterviewPendingResponse(res.status, data)) {
        await waitForGoalInterviewJob(data.poll_url ?? `/api/v1/llm-jobs/${data.job_id}`, errorCopy.startFailed)
        await load()
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : errorCopy.startFallback)
    } finally {
      setIsSending(false)
    }
  }, [buildPath, errorCopy, isSending, load])

  const sendMessage = useCallback(async (message: string) => {
    if (!message.trim() || isSending) return
    setIsSending(true)
    setError(null)
    // 낙관적 업데이트: 사용자 메시지 즉시 표시
    setProfile((prev) => prev ? {
      ...prev,
      messages: [...prev.messages, { role: 'user', content: message }],
    } : prev)
    try {
      const res = await fetch(buildPath('interview'), {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ message }),
      })
      if (!res.ok) {
        const d = await res.json().catch(() => ({})) as SafetyAPIErrorPayload
        throw new Error(resolveSafetyInputMessage(d, d.error ?? errorCopy.sendFailed))
      }
      const data = await res.json() as GoalInterviewAPIResponse
      setProfile(data.goal)
      if (isGoalInterviewPendingResponse(res.status, data)) {
        await waitForGoalInterviewJob(data.poll_url ?? `/api/v1/llm-jobs/${data.job_id}`, errorCopy.sendFailed)
        await load()
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : errorCopy.sendFallback
      // 낙관적 업데이트 롤백 후에도 사용자 안내는 유지한다.
      await load()
      setError(message)
    } finally {
      setIsSending(false)
    }
  }, [buildPath, errorCopy, isSending, load])

  const confirmGoal = useCallback(async (confirmedGoal: string) => {
    if (!confirmedGoal.trim() || isSending) return null
    setIsSending(true)
    setError(null)
    try {
      const res = await fetch(buildPath('confirm'), {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ confirmed_goal: confirmedGoal }),
      })
      if (!res.ok) {
        const d = await res.json().catch(() => ({})) as { error?: string }
        throw new Error(d.error ?? errorCopy.confirmFailed)
      }
      const data = await res.json() as { goal: GoalProfile }
      setProfile(data.goal)
      return data.goal
    } catch (err) {
      setError(err instanceof Error ? err.message : errorCopy.confirmFallback)
      return null
    } finally {
      setIsSending(false)
    }
  }, [buildPath, errorCopy, isSending])

  const reviseGoal = useCallback(async (message: string) => {
    if (!message.trim() || isSending) return
    setIsSending(true)
    setError(null)
    setProfile((prev) => prev ? {
      ...prev,
      messages: [...prev.messages, { role: 'user', content: message }],
    } : prev)
    try {
      const res = await fetch(buildPath('revise'), {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ message }),
      })
      if (!res.ok) {
        const d = await res.json().catch(() => ({})) as { error?: string }
        throw new Error(d.error ?? errorCopy.reviseFailed)
      }
      const data = await res.json() as GoalInterviewAPIResponse
      setProfile(data.goal)
      if (isGoalInterviewPendingResponse(res.status, data)) {
        await waitForGoalInterviewJob(data.poll_url ?? `/api/v1/llm-jobs/${data.job_id}`, errorCopy.reviseFailed)
        await load()
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : errorCopy.reviseFallback
      await load()
      setError(message)
    } finally {
      setIsSending(false)
    }
  }, [buildPath, errorCopy, isSending, load])

  const applyRebuildDecision = useCallback(async (decision: RebuildDecision) => {
    if (isSending) return
    setIsSending(true)
    setError(null)
    try {
      const res = await fetch(buildPath('rebuild-decision'), {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ decision }),
      })
      if (!res.ok) {
        const d = await res.json().catch(() => ({})) as { error?: string }
        throw new Error(d.error ?? errorCopy.rebuildFailed)
      }
      setProfile((prev) => prev ? { ...prev, rebuild_decision: decision, interview_state: 'confirmed' } : prev)
    } catch (err) {
      setError(err instanceof Error ? err.message : errorCopy.rebuildFallback)
    } finally {
      setIsSending(false)
    }
  }, [buildPath, errorCopy, isSending])

  const cancelRevision = useCallback(async () => {
    if (isSending) return null
    setIsSending(true)
    setError(null)
    try {
      const res = await fetch(buildPath('cancel-revision'), {
        method: 'POST',
        credentials: 'include',
      })
      if (!res.ok) {
        const d = await res.json().catch(() => ({})) as { error?: string }
        throw new Error(d.error ?? errorCopy.cancelFailed)
      }
      const data = await res.json() as { goal: GoalProfile }
      setProfile(data.goal)
      return data.goal
    } catch (err) {
      setError(err instanceof Error ? err.message : errorCopy.cancelFallback)
      await load()
      return null
    } finally {
      setIsSending(false)
    }
  }, [buildPath, errorCopy, isSending, load])

  const attachDraft = useCallback(async (draftId: string) => {
    if (hasDraftScope) return
    const res = await fetch('/api/v1/goals/attach-draft', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ draft_id: draftId }),
    })
    if (!res.ok) {
      const d = await res.json().catch(() => ({})) as { error?: string }
      throw new Error(d.error ?? errorCopy.attachFailed)
    }
    const data = await res.json() as { goal: GoalProfile }
    setProfile(data.goal)
  }, [errorCopy, hasDraftScope])

  const resetActiveGoal = useCallback(async () => {
    if (hasDraftScope) return
    setError(null)
    setProfile(null)
    const res = await fetch('/api/v1/goals/active', {
      method: 'DELETE',
      credentials: 'include',
    })
    if (!res.ok) {
      const d = await res.json().catch(() => ({})) as { error?: string }
      throw new Error(d.error ?? errorCopy.resetFailed)
    }
  }, [errorCopy, hasDraftScope])

  return { profile, isLoading, isSending, error, startInterview, sendMessage, confirmGoal, reviseGoal, cancelRevision, applyRebuildDecision, attachDraft, resetActiveGoal }
}
