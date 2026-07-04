'use client'

import { Suspense, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import AppHeaderShell from '@/components/common/AppHeaderShell'
import type { Locale } from '@/lib/i18n/locales'
import { getLoginPageCopy } from '@/lib/i18n/pages/login'

function GoogleIcon() {
  return (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" aria-hidden="true">
      <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4" />
      <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
      <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z" fill="#FBBC05" />
      <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
    </svg>
  )
}

function KakaoIcon() {
  return (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" aria-hidden="true">
      <path fillRule="evenodd" clipRule="evenodd" d="M12 2C6.48 2 2 5.82 2 10.5c0 3.01 1.87 5.65 4.7 7.18L5.63 22l5.05-2.73c.43.06.87.09 1.32.09 5.52 0 10-3.82 10-8.5S17.52 2 12 2z" fill="#3C1E1E" />
    </svg>
  )
}

function NaverIcon() {
  return (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" aria-hidden="true">
      <path d="M13.37 12.28L10.43 7H7.5v10h3.13v-5.28L13.57 17H16.5V7h-3.13v5.28z" fill="white" />
    </svg>
  )
}

function normalizeRedirectAfter(value: string | null): string {
  const trimmed = (value ?? '').trim()
  if (!trimmed) return '/dashboard'
  if (!trimmed.startsWith('/') || trimmed.startsWith('//')) return '/dashboard'
  if (trimmed.includes('\\') || trimmed.includes('\r') || trimmed.includes('\n') || trimmed.includes('\t')) return '/dashboard'
  try {
    const parsed = new URL(trimmed, 'https://learncosmos.co.kr')
    return `${parsed.pathname}${parsed.search}${parsed.hash}` || '/dashboard'
  } catch {
    return '/dashboard'
  }
}

function LoginContent({ locale }: { locale: Locale }) {
  const copy = getLoginPageCopy(locale)
  const searchParams = useSearchParams()
  const redirectAfter = normalizeRedirectAfter(searchParams.get('redirect_after'))
  const [showNaverNotice, setShowNaverNotice] = useState(false)
  const termsHref = locale === 'en' ? '/en/terms' : '/terms'
  const privacyHref = locale === 'en' ? '/en/privacy' : '/privacy'

  const handleLogin = (provider: string) => {
    if (provider === 'naver') {
      setShowNaverNotice(true)
      return
    }

    const params = new URLSearchParams({
      redirect_after: redirectAfter,
      locale,
    })
    window.location.href = `/api/v1/auth/${provider}/login?${params.toString()}`
  }

  return (
    <main className="relative flex min-h-screen items-center justify-center overflow-hidden bg-[linear-gradient(145deg,#0D314E_0%,#1C7D79_52%,#CC5216_100%)] px-6 py-28 text-[#E8EAF2]">
      <AppHeaderShell
        logoHref={copy.logoHref}
        logoIconSize={36}
        logoTextSize="20px"
        centerSlot={null}
        rightSlot={null}
      />
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_50%_0%,rgba(255,255,255,0.16),transparent_34%)]" />
      <div className="relative z-10 w-full max-w-sm">
        <div className="mb-6 flex justify-center">
          <div className="flex h-16 w-16 items-center justify-center rounded-full border border-white/35 bg-white/18">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="32"
              height="32"
              viewBox="0 0 24 24"
              fill="none"
              stroke="#ffffff"
              strokeWidth="1.8"
              strokeLinecap="round"
              strokeLinejoin="round"
              aria-hidden="true"
            >
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
              <circle cx="12" cy="7" r="4" />
              <path d="M16 11l4 4-4 4" />
              <path d="M20 15H9" />
            </svg>
          </div>
        </div>
        <div className="mb-10 text-center">
          <p className="mb-3 text-sm font-medium uppercase tracking-widest text-white/70">{copy.eyebrow}</p>
          <h1 className="mb-4 text-3xl font-bold leading-tight text-white">{copy.title}</h1>
          <p className="text-sm leading-7 text-white/74">{copy.description}</p>
        </div>

        <div className="mb-8 rounded-2xl border border-white/14 bg-white/8 px-5 py-4 text-center text-sm leading-7 text-white/76">
          {copy.aiNoticeLines.map((line) => (
            <span key={line} className="block">
              {line}
            </span>
          ))}
        </div>

        <div className="flex flex-col gap-4">
          <button
            onClick={() => handleLogin('google')}
            className="flex min-h-14 w-full items-center justify-center gap-3 rounded-xl border border-gray-300 bg-white px-6 py-3.5 text-sm font-medium text-gray-700 shadow-[0_16px_32px_rgba(3,8,20,0.2)] transition-colors hover:bg-gray-50"
          >
            <GoogleIcon />
            {copy.providers.google}
          </button>

          <button
            onClick={() => handleLogin('kakao')}
            className="flex min-h-14 w-full items-center justify-center gap-3 rounded-xl bg-[#FEE500] px-6 py-3.5 text-sm font-medium text-[#3C1E1E] transition-opacity hover:opacity-90"
          >
            <KakaoIcon />
            {copy.providers.kakao}
          </button>

          <button
            onClick={() => handleLogin('naver')}
            className="flex min-h-14 w-full items-center justify-center gap-3 rounded-xl bg-[#03C75A] px-6 py-3.5 text-sm font-medium text-white transition-opacity hover:opacity-90"
          >
            <NaverIcon />
            {copy.providers.naver}
          </button>
        </div>

        <p className="mt-7 text-center text-xs leading-6 text-white/68">
          {copy.policyNotice.before}
          <br />
          <a href={termsHref} className="underline underline-offset-4 hover:text-white">{copy.policyNotice.terms}</a>
          {' '}{copy.policyNotice.middle}{' '}
          <a href={privacyHref} className="underline underline-offset-4 hover:text-white">{copy.policyNotice.privacy}</a>
          {' '}{copy.policyNotice.after}
        </p>
      </div>
      {showNaverNotice && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-[#081826]/70 px-5 backdrop-blur-sm"
          role="dialog"
          aria-modal="true"
          aria-labelledby="naver-login-notice-title"
        >
          <div className="w-full max-w-sm rounded-2xl border border-white/16 bg-[#10283B] p-6 text-center shadow-[0_28px_72px_rgba(0,0,0,0.42)]">
            <h2 id="naver-login-notice-title" className="text-lg font-bold text-white">
              {copy.naverNotice.title}
            </h2>
            <p className="mt-3 text-sm leading-7 text-white/76">
              {copy.naverNotice.bodyLines.map((line) => (
                <span key={line} className="block">
                  {line}
                </span>
              ))}
            </p>
            <button
              type="button"
              onClick={() => setShowNaverNotice(false)}
              className="mt-6 min-h-11 w-full rounded-xl bg-white px-5 py-3 text-sm font-bold text-[#0D314E] transition-colors hover:bg-white/90"
            >
              {copy.naverNotice.confirm}
            </button>
          </div>
        </div>
      )}
    </main>
  )
}

export default function LoginPageClient({ locale }: { locale: Locale }) {
  return (
    <Suspense>
      <LoginContent locale={locale} />
    </Suspense>
  )
}
