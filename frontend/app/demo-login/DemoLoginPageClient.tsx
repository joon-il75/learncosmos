'use client'

import { FormEvent, useState } from 'react'
import Link from 'next/link'
import AppHeaderShell from '@/components/common/AppHeaderShell'
import type { Locale } from '@/lib/i18n/locales'
import { getDemoLoginPageCopy } from '@/lib/i18n/pages/demoLogin'

export default function DemoLoginPageClient({ locale }: { locale: Locale }) {
  const copy = getDemoLoginPageCopy(locale)
  const [code, setCode] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const trimmedCode = code.trim()
    if (!trimmedCode) {
      setError(copy.errors.invalid_request)
      return
    }

    setSubmitting(true)
    setError('')
    try {
      const res = await fetch('/api/v1/auth/demo/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ code: trimmedCode }),
      })
      const payload = await res.json().catch(() => ({}))
      if (!res.ok) {
        const errorCode = typeof payload.error_code === 'string' ? payload.error_code : 'unknown'
        throw new Error(copy.errors[errorCode] ?? copy.errors.unknown)
      }

      const redirectURL = typeof payload.redirect_url === 'string' && payload.redirect_url ? payload.redirect_url : '/dashboard'
      window.location.href = redirectURL
    } catch (err) {
      setError(err instanceof Error ? err.message : copy.errors.unknown)
    } finally {
      setSubmitting(false)
    }
  }

  const loginHref = locale === 'en' ? '/en/login' : '/login'

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
              <rect width="18" height="14" x="3" y="5" rx="2" />
              <path d="M7 15h.01" />
              <path d="M11 15h6" />
              <path d="M7 9h10" />
            </svg>
          </div>
        </div>

        <div className="mb-8 text-center">
          <p className="mb-3 text-sm font-medium uppercase tracking-widest text-white/70">{copy.eyebrow}</p>
          <h1 className="mb-4 text-3xl font-bold leading-tight text-white">{copy.title}</h1>
          <p className="text-sm leading-7 text-white/74">{copy.description}</p>
        </div>

        <div className="mb-6 rounded-2xl border border-white/14 bg-white/8 px-5 py-4 text-center text-sm leading-7 text-white/76">
          {copy.noticeLines.map((line) => (
            <span key={line} className="block">
              {line}
            </span>
          ))}
        </div>

        <form onSubmit={handleSubmit} className="grid gap-4">
          <label className="grid gap-2">
            <span className="text-xs font-bold uppercase tracking-[0.12em] text-white/72">{copy.codeLabel}</span>
            <input
              value={code}
              onChange={(event) => setCode(event.target.value)}
              placeholder={copy.codePlaceholder}
              autoCapitalize="characters"
              spellCheck={false}
              className="min-h-14 rounded-xl border border-white/20 bg-white px-4 text-center text-base font-bold tracking-[0.08em] text-[#0D314E] outline-none transition-colors placeholder:text-gray-400 focus:border-white"
            />
          </label>
          <button
            type="submit"
            disabled={submitting}
            className="min-h-14 rounded-xl bg-white px-6 py-3.5 text-sm font-bold text-[#0D314E] shadow-[0_16px_32px_rgba(3,8,20,0.2)] transition-colors hover:bg-white/90 disabled:cursor-not-allowed disabled:opacity-65"
          >
            {submitting ? copy.submitting : copy.submit}
          </button>
        </form>

        {error ? (
          <p className="mt-5 rounded-xl border border-red-200/30 bg-red-950/35 px-4 py-3 text-center text-sm leading-6 text-red-100">
            {error}
          </p>
        ) : null}

        <p className="mt-7 text-center text-xs leading-6 text-white/68">
          <Link href={loginHref} className="underline underline-offset-4 hover:text-white">
            {copy.backToLogin}
          </Link>
        </p>
      </div>
    </main>
  )
}
