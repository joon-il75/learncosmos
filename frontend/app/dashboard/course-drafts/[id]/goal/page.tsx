'use client'

import { useEffect, type CSSProperties } from 'react'
import { useParams, useRouter, useSearchParams } from 'next/navigation'

export default function LegacyGoalRedirectPage() {
  const router = useRouter()
  const params = useParams<{ id: string }>()
  const draftId = typeof params.id === 'string' ? params.id : ''
  const searchParams = useSearchParams()
  const intent = searchParams.get('intent') ?? ''

  useEffect(() => {
    let cancelled = false

    const cleanupAndRedirect = async () => {
      if (draftId) {
        try {
          const draftRes = await fetch(`/api/v1/course-drafts/${draftId}`, {
            credentials: 'include',
            cache: 'no-store',
          })
          if (draftRes.ok) {
            const draftData = await draftRes.json() as { draft?: { draft?: { title?: string } } }
            const title = draftData.draft?.draft?.title?.trim()
            if (title) {
              await fetch(`/api/v1/course-drafts/${draftId}`, {
                method: 'DELETE',
                credentials: 'include',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ confirm_title: title }),
              })
            }
          }
        } catch {
          // cleanup 실패는 redirect를 막지 않는다.
        }
      }

      if (cancelled) return
      const target = intent ? `/dashboard/goal?intent=${encodeURIComponent(intent)}` : '/dashboard/goal'
      router.replace(target)
    }

    void cleanupAndRedirect()

    return () => {
      cancelled = true
    }
  }, [draftId, intent, router])

  return (
    <div style={loadingStyle}>
      <div style={loadingTextStyle}>목표 인터뷰 화면으로 이동하는 중...</div>
    </div>
  )
}

const loadingStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  height: '100dvh',
  background: '#FFF8E7',
}

const loadingTextStyle: CSSProperties = {
  fontSize: 16,
  color: '#92400E',
}
