'use client'

import { Suspense, useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import type { Locale } from '@/lib/i18n/locales'

interface ConsentStatus {
  terms_agreed: boolean
  privacy_agreed: boolean
  required_consent_pending: boolean
  required_documents?: Array<{
    id: string
    set_id?: string
    type: string
    title: string
    version: number
    locale?: Locale
    requested_locale?: Locale
    fallback_used?: boolean
    effective_at?: string
  }>
}

type ApiErrorPayload = {
  error?: string
  error_code?: string
}

type AgreementsCopy = {
  loginPath: string
  currentPath: string
  loading: string
  loadError: string
  submitError: string
  requiredError: string
  invalidRequestError: string
  eyebrow: string
  title: string
  description: string
  notice: string[]
  requiredDocumentsLabel: string
  requiredBadge: string
  versionLabel: string
  effectiveLabel: string
  fallbackNotice: string
  termsAgree: string
  termsLink: string
  termsHref: string
  privacyAgree: string
  privacyLink: string
  privacyHref: string
  submitting: string
  submit: string
  home: string
  homeHref: string
}

const agreementsCopy: Record<Locale, AgreementsCopy> = {
  ko: {
    loginPath: '/login',
    currentPath: '/agreements',
    loading: '동의 상태를 확인하는 중...',
    loadError: '동의 상태를 불러오지 못했습니다',
    submitError: '동의 저장에 실패했습니다',
    requiredError: '서비스 이용약관과 개인정보처리방침에 모두 동의해야 합니다',
    invalidRequestError: '잘못된 요청입니다.',
    eyebrow: 'Required Consent',
    title: '서비스 이용 전 약관 동의가 필요합니다',
    description: '신규 가입이거나 운영 중인 약관 문서가 변경된 경우, 최신 서비스 이용약관과 개인정보처리방침에 다시 동의해야 대시보드와 학습 기능을 사용할 수 있습니다.',
    notice: [
      'LearnCosmos는 AI 학습 기능 제공을 위해 사용자의 학습 입력, 코스 생성 정보, AI 사용량 정보를 처리할 수 있습니다.',
      'BYOK 기능을 사용하는 경우 사용자가 등록한 API Key는 암호화 저장되며, AI 호출에 필요한 순간에만 사용됩니다. API Key 원문, 프롬프트 본문, AI 응답 본문은 사용량 기록에 저장하지 않습니다.',
      '콘텐츠 검색과 추천은 서비스 서버의 설치형 임베딩 모델을 기반으로 처리하며, 사용자가 등록한 BYOK LLM 키와 사용 범위를 분리합니다.',
      'OpenAI를 시작으로 Claude, Gemini, Grok, HyperCLOVA X, Solar, Llama, EXAONE 등은 단계적 연결 예정 제공자이며, 실제 지원 제공자와 모델 범위는 서비스 화면과 운영 정책에 따라 제공됩니다.',
      'BYOK는 선택 기능이며, API Key를 등록하지 않아도 서비스 기본 AI 포인트 정책에 따라 AI 기능을 사용할 수 있습니다.',
      '연구지점 등에서 직접 작성하거나 첨부하는 학습 콘텐츠의 저작권 및 이용 권한 확인 책임은 작성자에게 있습니다. 타인의 자료는 관련 법령과 허용된 범위 안에서만 사용해야 합니다.',
      'LearnCosmos는 자체 플랫폼 콘텐츠를 제외한 외부 콘텐츠의 무료/유료 구분과 무관하게 동일한 학습 흐름으로 제공합니다.',
      '외부 콘텐츠는 영상 자체 과금·판매 항목이 아니며, 결제 혜택은 AI 피드백·학습 기록 분석·복습 루틴·결과물 코칭·개인 학습 리포트에만 연결됩니다.',
      '학습 추천 결과에는 YouTube 출처, 원본 링크, 공식 플레이어 링크가 항상 함께 노출되며, 영상 위나 재생 버튼 앞에 유료 CTA를 배치하지 않습니다.',
    ],
    requiredDocumentsLabel: '최신 필수 문서',
    requiredBadge: '필수',
    versionLabel: '버전',
    effectiveLabel: '시행',
    fallbackNotice: '요청한 언어의 번역본이 아직 없어 한국어 원문을 표시합니다.',
    termsAgree: '서비스 이용약관에 동의합니다.',
    termsLink: '약관 보기',
    termsHref: '/terms',
    privacyAgree: '개인정보처리방침에 동의합니다.',
    privacyLink: '방침 보기',
    privacyHref: '/privacy',
    submitting: '저장 중...',
    submit: '동의하고 계속하기',
    home: '홈으로',
    homeHref: '/',
  },
  en: {
    loginPath: '/en/login',
    currentPath: '/en/agreements',
    loading: 'Checking consent status...',
    loadError: 'Could not load consent status',
    submitError: 'Could not save consent',
    requiredError: 'You must agree to both the Terms of Service and Privacy Policy.',
    invalidRequestError: 'Invalid request.',
    eyebrow: 'Required Consent',
    title: 'Policy consent is required before using LearnCosmos',
    description: 'If you are a new user or required policy documents have changed, you must agree to the current Terms of Service and Privacy Policy before using the dashboard and learning features.',
    notice: [
      'LearnCosmos may process your learning inputs, course generation information, and AI usage data to provide AI learning features.',
      'If you use BYOK, your API key is stored encrypted and used only when needed for AI calls. Raw API keys, prompt text, and AI response text are not stored in usage records.',
      'Content search and recommendations are handled through an installed embedding model on the LearnCosmos server, separate from the scope of user-registered BYOK LLM keys.',
      'OpenAI starts first, while Claude, Gemini, Grok, HyperCLOVA X, Solar, Llama, and EXAONE are planned providers. Supported providers and models are introduced according to service screens and operating policy.',
      'BYOK is optional. You can use AI features under the default AI point policy without registering an API key.',
      'You are responsible for confirming copyright and usage rights for learning content you write or attach in research points and related areas.',
      'LearnCosmos provides access to external content without dividing learning flow by paid/free external videos.',
      'Payment benefits are attached only to AI feedback, learning-record analytics, review routines, artifact coaching, and personal learning reports.',
      'Learning screens must show YouTube source, original URL, and official player link together, and must not place paid CTAs above the play button or over the player area.',
    ],
    requiredDocumentsLabel: 'Required Documents',
    requiredBadge: 'Required',
    versionLabel: 'Version',
    effectiveLabel: 'Effective',
    fallbackNotice: 'A translation for the requested language is not available yet, so the Korean source text is shown.',
    termsAgree: 'I agree to the Terms of Service.',
    termsLink: 'View terms',
    termsHref: '/en/terms',
    privacyAgree: 'I agree to the Privacy Policy.',
    privacyLink: 'View policy',
    privacyHref: '/en/privacy',
    submitting: 'Saving...',
    submit: 'Agree and continue',
    home: 'Home',
    homeHref: '/en',
  },
}

function normalizeRedirectAfter(value: string | null) {
  if (!value || !value.startsWith('/') || value.startsWith('//')) return '/dashboard'
  return value
}

function formatDate(value: string | undefined, locale: Locale) {
  if (!value) return ''
  return new Date(value).toLocaleDateString(locale === 'en' ? 'en-US' : 'ko-KR')
}

function getConsentErrorMessage(payload: ApiErrorPayload, copy: AgreementsCopy) {
  if (payload.error_code === 'invalid_request') return copy.invalidRequestError
  if (payload.error_code === 'consent_required_terms_privacy') return copy.requiredError
  if (payload.error_code === 'consent_save_failed' || payload.error_code === 'consent_status_load_failed') return copy.submitError
  return payload.error || copy.submitError
}

function AgreementsContent({ locale }: { locale: Locale }) {
  const copy = agreementsCopy[locale]
  const router = useRouter()
  const searchParams = useSearchParams()
  const redirectAfter = normalizeRedirectAfter(searchParams.get('redirect_after'))
  const [agreeTerms, setAgreeTerms] = useState(false)
  const [agreePrivacy, setAgreePrivacy] = useState(false)
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [requiredDocuments, setRequiredDocuments] = useState<ConsentStatus['required_documents']>([])

  useEffect(() => {
    const loadStatus = async () => {
      const refreshRes = await fetch('/api/v1/auth/refresh', { method: 'POST', credentials: 'include' })
      if (!refreshRes.ok) {
        router.replace(`${copy.loginPath}?redirect_after=${encodeURIComponent(copy.currentPath)}`)
        return
      }

      const consentRes = await fetch(`/api/v1/users/me/consents?locale=${locale}`, { credentials: 'include', cache: 'no-store' })
      if (!consentRes.ok) {
        router.replace(`${copy.loginPath}?redirect_after=${encodeURIComponent(copy.currentPath)}`)
        return
      }

      const consentData: ConsentStatus = await consentRes.json()
      if (!consentData.required_consent_pending) {
        router.replace(redirectAfter)
        return
      }

      setAgreeTerms(consentData.terms_agreed)
      setAgreePrivacy(consentData.privacy_agreed)
      setRequiredDocuments(consentData.required_documents || [])
      setLoading(false)
    }

    loadStatus().catch(() => {
      setError(copy.loadError)
      setLoading(false)
    })
  }, [copy.currentPath, copy.loadError, copy.loginPath, locale, redirectAfter, router])

  const handleSubmit = async () => {
    if (!agreeTerms || !agreePrivacy) {
      setError(copy.requiredError)
      return
    }

    setSubmitting(true)
    setError('')
    try {
      const res = await fetch('/api/v1/users/me/consents', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
          agree_terms: true,
          agree_privacy: true,
          locale,
        }),
      })

      if (!res.ok) {
        const data = (await res.json().catch(() => ({}))) as ApiErrorPayload
        throw new Error(getConsentErrorMessage(data, copy))
      }

      router.replace(redirectAfter)
    } catch (err) {
      setError(err instanceof Error ? err.message : copy.submitError)
      setSubmitting(false)
    }
  }

  if (loading) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-[#0B1629] text-[#E8EAF2]">
        <p className="text-sm text-[rgba(200,210,235,0.6)]">{copy.loading}</p>
      </main>
    )
  }

  return (
    <main className="min-h-screen bg-[#0B1629] px-6 py-16 text-[#E8EAF2]">
      <div className="mx-auto max-w-2xl rounded-[28px] border border-[rgba(160,186,224,0.18)] bg-[rgba(11,22,41,0.92)] p-8 shadow-[0_32px_120px_rgba(3,8,20,0.45)]">
        <p className="mb-3 text-xs font-semibold uppercase tracking-[0.28em] text-[#7AA2D6]">{copy.eyebrow}</p>
        <h1 className="mb-4 text-3xl font-bold">{copy.title}</h1>
        <p className="mb-8 text-sm leading-7 text-[rgba(220,228,245,0.72)]">{copy.description}</p>

        <div className="mb-8 rounded-2xl border border-[rgba(122,162,214,0.22)] bg-[rgba(122,162,214,0.08)] p-4 text-sm leading-7 text-[rgba(232,234,242,0.82)]">
          {copy.notice.map((line, index) => (
            <p key={line} className={index === 0 ? '' : 'mt-3'}>
              {line}
            </p>
          ))}
        </div>

        {requiredDocuments && requiredDocuments.length > 0 ? (
          <div className="mb-8 rounded-2xl border border-[rgba(160,186,224,0.16)] bg-[rgba(255,255,255,0.03)] p-4">
            <p className="mb-3 text-xs font-semibold uppercase tracking-[0.18em] text-[rgba(122,162,214,0.9)]">
              {copy.requiredDocumentsLabel}
            </p>
            <div className="space-y-2">
              {requiredDocuments.map((doc) => (
                <div key={doc.set_id || doc.id} className="flex flex-wrap items-center justify-between gap-4 rounded-xl bg-[rgba(255,255,255,0.03)] px-4 py-3">
                  <div>
                    <p className="text-sm font-semibold text-[#E8EAF2]">{doc.title}</p>
                    <p className="mt-1 text-xs text-[rgba(220,228,245,0.58)]">
                      {copy.versionLabel} v{doc.version}
                      {doc.effective_at ? ` · ${copy.effectiveLabel} ${formatDate(doc.effective_at, locale)}` : ''}
                    </p>
                    {doc.fallback_used ? <p className="mt-1 text-xs text-[#F5C76B]">{copy.fallbackNotice}</p> : null}
                  </div>
                  <span className="rounded-full border border-[rgba(122,162,214,0.24)] px-3 py-1 text-[11px] font-semibold text-[#7AA2D6]">
                    {copy.requiredBadge}
                  </span>
                </div>
              ))}
            </div>
          </div>
        ) : null}

        <div className="space-y-4">
          <label className="flex items-start gap-3 rounded-2xl border border-[rgba(160,186,224,0.16)] bg-[rgba(255,255,255,0.03)] p-4">
            <input type="checkbox" className="mt-1 h-4 w-4" checked={agreeTerms} onChange={(e) => setAgreeTerms(e.target.checked)} />
            <span className="text-sm leading-6 text-[rgba(232,234,242,0.9)]">
              {copy.termsAgree} <a href={copy.termsHref} className="underline underline-offset-4">{copy.termsLink}</a>
            </span>
          </label>

          <label className="flex items-start gap-3 rounded-2xl border border-[rgba(160,186,224,0.16)] bg-[rgba(255,255,255,0.03)] p-4">
            <input type="checkbox" className="mt-1 h-4 w-4" checked={agreePrivacy} onChange={(e) => setAgreePrivacy(e.target.checked)} />
            <span className="text-sm leading-6 text-[rgba(232,234,242,0.9)]">
              {copy.privacyAgree} <a href={copy.privacyHref} className="underline underline-offset-4">{copy.privacyLink}</a>
            </span>
          </label>
        </div>

        {error ? <p className="mt-5 text-sm text-[#FF9A9A]">{error}</p> : null}

        <div className="mt-8 flex flex-wrap gap-3">
          <button
            type="button"
            onClick={handleSubmit}
            disabled={submitting}
            className="rounded-xl bg-[#378ADD] px-5 py-3 text-sm font-semibold text-white transition-opacity hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {submitting ? copy.submitting : copy.submit}
          </button>
          <a href={copy.homeHref} className="rounded-xl border border-[rgba(160,186,224,0.2)] px-5 py-3 text-sm font-semibold text-[rgba(232,234,242,0.88)]">
            {copy.home}
          </a>
        </div>
      </div>
    </main>
  )
}

export default function AgreementsPageClient({ locale }: { locale: Locale }) {
  const copy = agreementsCopy[locale]

  return (
    <Suspense
      fallback={
        <main className="flex min-h-screen items-center justify-center bg-[#0B1629] text-[#E8EAF2]">
          <p className="text-sm text-[rgba(200,210,235,0.6)]">{copy.loading}</p>
        </main>
      }
    >
      <AgreementsContent locale={locale} />
    </Suspense>
  )
}
