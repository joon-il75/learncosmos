'use client'

import { Suspense, useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'

const API_BASE = ''

function CallbackContent() {
  const router = useRouter()
  const searchParams = useSearchParams()

  useEffect(() => {
    const error = searchParams.get('error')

    if (error) {
      router.replace(`/login?error=${encodeURIComponent(error)}`)
      return
    }

    const finishCookieSession = async () => {
      const refreshRes = await fetch(`${API_BASE}/api/v1/auth/refresh`, {
        method: 'POST',
        credentials: 'include',
      })
      if (!refreshRes.ok) {
        router.replace('/login?error=unsupported_callback')
        return
      }
      router.replace('/dashboard')
    }

    finishCookieSession().catch(() => router.replace('/login?error=unsupported_callback'))
  }, [searchParams, router])

  return (
    <main className="flex min-h-screen items-center justify-center bg-[#0B1629]">
      <div className="flex flex-col items-center gap-4">
        <div className="h-10 w-10 animate-spin rounded-full border-4 border-[#378ADD] border-t-transparent" />
        <p className="text-sm text-[rgba(200,210,235,0.6)]">로그인 중...</p>
      </div>
    </main>
  )
}

export default function CallbackPage() {
  return (
    <Suspense>
      <CallbackContent />
    </Suspense>
  )
}
