'use client'

import { Suspense, useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import type { Locale } from '@/lib/i18n/locales'

type LanguagePreferences = {
  ui_locale: Locale
  learning_language: Locale
  language_setup_required: boolean
}

type MeResponse = {
  required_consent_pending?: boolean
}

type ApiErrorPayload = {
  error?: string
  error_code?: string
}

type Copy = {
  loginPath: string
  currentPath: string
  loading: string
  loadError: string
  saveError: string
  invalidRequestError: string
  eyebrow: string
  title: string
  description: string
  languageLabel: string
  languageHelp: string
  korean: string
  english: string
  saving: string
  submit: string
  home: string
  homeHref: string
}

const copyByLocale: Record<Locale, Copy> = {
  ko: {
    loginPath: '/login',
    currentPath: '/language-setup',
    loading: '언어 설정을 확인하는 중...',
    loadError: '언어 설정을 불러오지 못했습니다',
    saveError: '언어 설정을 저장하지 못했습니다',
    invalidRequestError: '잘못된 요청입니다.',
    eyebrow: 'Language Setup',
    title: 'Language Settings',
    description: '사용할 언어를 하나만 선택합니다. 화면 문구와 AI 생성 결과가 같은 언어로 저장됩니다.',
    languageLabel: 'Language',
    languageHelp: '메뉴, 안내문, 동의 화면, 목표 채팅, 리슨 생성, 추천 문구에 함께 적용됩니다.',
    korean: '한국어',
    english: 'English',
    saving: '저장 중...',
    submit: '저장하고 계속하기',
    home: '홈으로',
    homeHref: '/',
  },
  en: {
    loginPath: '/en/login',
    currentPath: '/en/language-setup',
    loading: 'Checking language preferences...',
    loadError: 'Could not load language preferences',
    saveError: 'Could not save language preferences',
    invalidRequestError: 'Invalid request.',
    eyebrow: 'Language Setup',
    title: 'Language Settings',
    description: 'Choose one language. LearnCosmos saves the same language for the interface and learning generation.',
    languageLabel: 'Language',
    languageHelp: 'Applies to menus, notices, consent screens, goal chat, lesson generation, and recommendation copy.',
    korean: 'Korean',
    english: 'English',
    saving: 'Saving...',
    submit: 'Save and continue',
    home: 'Home',
    homeHref: '/en',
  },
}

function normalizeRedirectAfter(value: string | null) {
  if (!value || !value.startsWith('/') || value.startsWith('//')) return '/dashboard'
  return value
}

function nextConsentPath(uiLocale: Locale) {
  return uiLocale === 'en' ? '/en/agreements' : '/agreements'
}

function getLanguageSetupErrorMessage(payload: ApiErrorPayload, copy: Copy) {
  if (payload.error_code === 'invalid_request') return copy.invalidRequestError
  if (payload.error_code === 'preferences_save_failed') return copy.saveError
  return payload.error || copy.saveError
}

function LanguageChoice({
  label,
  help,
  value,
  onChange,
  copy,
}: {
  label: string
  help: string
  value: Locale
  onChange: (value: Locale) => void
  copy: Copy
}) {
  return (
    <section className="rounded-lg border border-[#D8E0EA] bg-white p-5">
      <div className="mb-4">
        <h2 className="text-base font-bold text-[#122033]">{label}</h2>
        <p className="mt-1 text-sm leading-6 text-[#5C6B7A]">{help}</p>
      </div>
      <div className="grid grid-cols-2 gap-3">
        {(['ko', 'en'] as const).map((option) => {
          const active = value === option
          return (
            <button
              key={option}
              type="button"
              onClick={() => onChange(option)}
              className={[
                'min-h-12 rounded-lg border px-4 py-3 text-sm font-bold transition-colors',
                active
                  ? 'border-[#0D314E] bg-[#0D314E] text-white'
                  : 'border-[#C8D3DF] bg-white text-[#0D314E] hover:bg-[#F1F5F9]',
              ].join(' ')}
              aria-pressed={active}
            >
              {option === 'ko' ? copy.korean : copy.english}
            </button>
          )
        })}
      </div>
    </section>
  )
}

function LanguageSetupContent({ locale }: { locale: Locale }) {
  const copy = copyByLocale[locale]
  const router = useRouter()
  const searchParams = useSearchParams()
  const redirectAfter = normalizeRedirectAfter(searchParams.get('redirect_after'))
  const [uiLocale, setUILocale] = useState<Locale>(locale)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    const loadPreferences = async () => {
      const refreshRes = await fetch('/api/v1/auth/refresh', { method: 'POST', credentials: 'include' })
      if (!refreshRes.ok) {
        router.replace(`${copy.loginPath}?redirect_after=${encodeURIComponent(copy.currentPath)}`)
        return
      }

      const prefRes = await fetch(`/api/v1/users/me/preferences?locale=${locale}`, { credentials: 'include', cache: 'no-store' })
      if (!prefRes.ok) {
        router.replace(`${copy.loginPath}?redirect_after=${encodeURIComponent(copy.currentPath)}`)
        return
      }

      const data = (await prefRes.json()) as LanguagePreferences
      setUILocale(data.ui_locale || locale)
      if (!data.language_setup_required) {
        router.replace(redirectAfter)
        return
      }
      setLoading(false)
    }

    loadPreferences().catch(() => {
      setError(copy.loadError)
      setLoading(false)
    })
  }, [copy.currentPath, copy.loadError, copy.loginPath, locale, redirectAfter, router])

  const handleSubmit = async () => {
    setSaving(true)
    setError('')
    try {
      const res = await fetch('/api/v1/users/me/preferences', {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
          ui_locale: uiLocale,
          learning_language: uiLocale,
        }),
      })
      if (!res.ok) {
        const data = (await res.json().catch(() => ({}))) as ApiErrorPayload
        throw new Error(getLanguageSetupErrorMessage(data, copy))
      }

      const meRes = await fetch('/api/v1/auth/me', { credentials: 'include', cache: 'no-store' })
      if (meRes.ok) {
        const me = (await meRes.json()) as MeResponse
        if (me.required_consent_pending) {
          router.replace(`${nextConsentPath(uiLocale)}?redirect_after=${encodeURIComponent(redirectAfter)}`)
          return
        }
      }

      router.replace(redirectAfter)
    } catch (err) {
      setError(err instanceof Error ? err.message : copy.saveError)
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-[#F6F8FB] text-[#122033]">
        <p className="text-sm text-[#5C6B7A]">{copy.loading}</p>
      </main>
    )
  }

  return (
    <main className="min-h-screen bg-[#F6F8FB] px-5 py-12 text-[#122033]">
      <div className="mx-auto max-w-2xl">
        <a href={copy.homeHref} className="text-lg font-bold text-[#0D314E]">
          LearnCosmos
        </a>
        <header className="mt-10 mb-7">
          <p className="text-sm font-bold uppercase tracking-[0.18em] text-[#1C7D79]">{copy.eyebrow}</p>
          <h1 className="mt-4 text-3xl font-bold leading-tight text-[#0D314E]">{copy.title}</h1>
          <p className="mt-4 text-sm leading-7 text-[#42526B]">{copy.description}</p>
        </header>

        <div className="space-y-4">
          <LanguageChoice label={copy.languageLabel} help={copy.languageHelp} value={uiLocale} onChange={setUILocale} copy={copy} />
        </div>

        {error ? <p className="mt-5 text-sm font-semibold text-[#B42318]">{error}</p> : null}

        <div className="mt-7 flex flex-wrap gap-3">
          <button
            type="button"
            onClick={handleSubmit}
            disabled={saving}
            className="rounded-lg bg-[#0D314E] px-5 py-3 text-sm font-bold text-white transition-opacity hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {saving ? copy.saving : copy.submit}
          </button>
          <a href={copy.homeHref} className="rounded-lg border border-[#C8D3DF] px-5 py-3 text-sm font-bold text-[#0D314E]">
            {copy.home}
          </a>
        </div>
      </div>
    </main>
  )
}

export default function LanguageSetupPageClient({ locale }: { locale: Locale }) {
  const copy = copyByLocale[locale]

  return (
    <Suspense
      fallback={
        <main className="flex min-h-screen items-center justify-center bg-[#F6F8FB] text-[#122033]">
          <p className="text-sm text-[#5C6B7A]">{copy.loading}</p>
        </main>
      }
    >
      <LanguageSetupContent locale={locale} />
    </Suspense>
  )
}
