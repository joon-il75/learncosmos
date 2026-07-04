'use client'

import { useEffect, useMemo, useState, type CSSProperties } from 'react'
import type { DashboardGoalCopy } from '@/lib/i18n/pages/dashboardGoal'
import { getDashboardGoalCopy } from '@/lib/i18n/pages/dashboardGoal'

const COURSE_GEN_COST = 5

interface GoalPointUser {
  total_points?: number
}

interface GoalAISettings {
  provider?: string
  has_api_key?: boolean
  is_enabled?: boolean
}

function normalizePoints(value?: number) {
  return Number.isFinite(value) ? Math.max(0, Math.floor(value ?? 0)) : 0
}

interface GoalPointStatusBarProps {
  copy?: DashboardGoalCopy['points']
}

export default function GoalPointStatusBar({
  copy = getDashboardGoalCopy('ko').points,
}: GoalPointStatusBarProps) {
  const [totalPoints, setTotalPoints] = useState<number | null>(null)
  const [aiSettings, setAISettings] = useState<GoalAISettings | null>(null)

  useEffect(() => {
    let cancelled = false

    const load = async () => {
      try {
        const refreshRes = await fetch('/api/v1/auth/refresh', { method: 'POST', credentials: 'include' })
        if (!refreshRes.ok) return
        const meRes = await fetch('/api/v1/auth/me', { credentials: 'include', cache: 'no-store' })
        if (!meRes.ok) return
        const data = (await meRes.json()) as GoalPointUser
        const aiRes = await fetch('/api/v1/users/me/ai-settings', { credentials: 'include', cache: 'no-store' })
        const aiData = aiRes.ok ? ((await aiRes.json()) as GoalAISettings) : null
        if (!cancelled) {
          setTotalPoints(normalizePoints(data.total_points))
          setAISettings(aiData)
        }
      } catch {
        if (!cancelled) setTotalPoints(null)
      }
    }

    void load()
    const handleUserRefresh = () => void load()
    window.addEventListener('learnweaver:user-refresh', handleUserRefresh as EventListener)

    return () => {
      cancelled = true
      window.removeEventListener('learnweaver:user-refresh', handleUserRefresh as EventListener)
    }
  }, [])

  const status = useMemo(() => {
    const usesOpenAIByok = Boolean(
      aiSettings?.has_api_key &&
      aiSettings?.is_enabled &&
      (aiSettings.provider ?? '').trim().toLowerCase() === 'openai',
    )
    if (usesOpenAIByok) {
      return {
        title: totalPoints === null ? copy.byokTitle : copy.points(totalPoints),
        detail: copy.byokDetail,
        tone: 'byok' as const,
      }
    }
    if (totalPoints === null) {
      return {
        title: copy.checkingTitle,
        detail: copy.checkingDetail(COURSE_GEN_COST),
        tone: 'neutral' as const,
      }
    }
    if (totalPoints < COURSE_GEN_COST) {
      return {
        title: copy.points(totalPoints),
        detail: copy.insufficientDetail(COURSE_GEN_COST),
        tone: 'warning' as const,
      }
    }
    return {
      title: copy.points(totalPoints),
      detail: copy.readyDetail(COURSE_GEN_COST, totalPoints - COURSE_GEN_COST),
      tone: 'ready' as const,
    }
  }, [aiSettings, copy, totalPoints])

  return (
    <div style={barStyle(status.tone)} aria-live="polite">
      <div style={titleStyle(status.tone)}>{status.title}</div>
      <div style={detailStyle(status.tone)}>{status.detail}</div>
    </div>
  )
}

const barStyle = (tone: 'ready' | 'warning' | 'neutral' | 'byok'): CSSProperties => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  flexWrap: 'wrap',
  gap: 10,
  padding: '9px 12px',
  borderTop: '1px solid rgba(245,158,11,0.16)',
  background: tone === 'warning' ? '#FFF7ED' : tone === 'byok' ? '#ECFDF5' : '#FFFBF0',
  borderBottom: `1px solid ${tone === 'warning' ? 'rgba(234,88,12,0.26)' : tone === 'byok' ? 'rgba(5,150,105,0.22)' : 'rgba(245,158,11,0.14)'}`,
})

const titleStyle = (tone: 'ready' | 'warning' | 'neutral' | 'byok'): CSSProperties => ({
  flexShrink: 0,
  fontSize: 13,
  fontWeight: 900,
  color: tone === 'warning' ? '#C2410C' : tone === 'byok' ? '#047857' : '#78350F',
})

const detailStyle = (tone: 'ready' | 'warning' | 'neutral' | 'byok'): CSSProperties => ({
  minWidth: 0,
  flex: '1 1 220px',
  fontSize: 12,
  fontWeight: 700,
  color: tone === 'warning' ? '#9A3412' : tone === 'byok' ? '#065F46' : '#92400E',
  textAlign: 'right',
  lineHeight: 1.35,
})
