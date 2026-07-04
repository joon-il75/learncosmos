'use client'

import { useEffect, useMemo, useState } from 'react'
import Link from 'next/link'
import AppHeaderShell, { appHeaderActionLinkStyle } from '@/components/common/AppHeaderShell'
import type { Locale } from '@/lib/i18n/locales'

interface PolicyDocument {
  id: string
  set_id?: string
  type: string
  title: string
  version: number
  content: string
  locale?: Locale
  requested_locale?: Locale
  fallback_used?: boolean
  effective_at?: string
  updated_at?: string
}

type PolicyType = 'terms' | 'privacy'

interface PolicyDocumentPageProps {
  type: PolicyType
  locale?: Locale
  fallbackTitle: string
  loadingText: string
  errorText: string
  homeLabel?: string
  homeHref?: string
  tocLabel?: string
  historyTitle?: string
  historyText?: string
  fallbackNoticeTitle?: string
  fallbackNoticeText?: string
}

type ContentBlock =
  | { kind: 'heading2'; text: string; id: string }
  | { kind: 'heading3'; text: string }
  | { kind: 'paragraph'; text: string }
  | { kind: 'list'; items: string[] }

function formatDate(value?: string, locale: Locale = 'ko') {
  if (!value) return ''
  return new Date(value).toLocaleDateString(locale === 'en' ? 'en-US' : 'ko-KR')
}

function normalizeHeadingId(text: string, index: number) {
  return `section-${index}-${text
    .replace(/^[#\d.\s]+/, '')
    .replace(/[^\w가-힣]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .toLowerCase()}`
}

function parseContent(content: string): ContentBlock[] {
  const lines = content.split(/\r?\n/)
  const blocks: ContentBlock[] = []
  let paragraph: string[] = []
  let list: string[] = []
  let headingIndex = 0

  const flushParagraph = () => {
    if (!paragraph.length) return
    blocks.push({ kind: 'paragraph', text: paragraph.join('\n') })
    paragraph = []
  }

  const flushList = () => {
    if (!list.length) return
    blocks.push({ kind: 'list', items: list })
    list = []
  }

  for (const rawLine of lines) {
    const line = rawLine.trim()
    if (!line) {
      flushParagraph()
      flushList()
      continue
    }

    if (line.startsWith('# ')) {
      flushParagraph()
      flushList()
      continue
    }

    if (line.startsWith('## ')) {
      flushParagraph()
      flushList()
      headingIndex += 1
      const text = line.replace(/^##\s+/, '')
      blocks.push({ kind: 'heading2', text, id: normalizeHeadingId(text, headingIndex) })
      continue
    }

    if (line.startsWith('### ')) {
      flushParagraph()
      flushList()
      blocks.push({ kind: 'heading3', text: line.replace(/^###\s+/, '') })
      continue
    }

    if (line.startsWith('- ')) {
      flushParagraph()
      list.push(line.replace(/^-\s+/, ''))
      continue
    }

    flushList()
    paragraph.push(line)
  }

  flushParagraph()
  flushList()
  return blocks
}

export default function PolicyDocumentPage({
  type,
  locale = 'ko',
  fallbackTitle,
  loadingText,
  errorText,
  homeLabel = '홈으로',
  homeHref,
  tocLabel = '목차',
  historyTitle = '변경 이력',
  historyText = '현재 공개 문서는 위 시행일부터 적용됩니다. 이후 변경 사항은 서비스 화면 또는 별도 공지로 안내합니다.',
  fallbackNoticeTitle,
  fallbackNoticeText,
}: PolicyDocumentPageProps) {
  const [document, setDocument] = useState<PolicyDocument | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    fetch(`/api/v1/public/policies/${type}?locale=${locale}`, { cache: 'no-store' })
      .then(async (res) => {
        if (!res.ok) throw new Error(errorText)
        return res.json()
      })
      .then((data) => setDocument(data.document))
      .catch((err) => setError(err instanceof Error ? err.message : errorText))
  }, [errorText, locale, type])

  const blocks = useMemo(() => parseContent(document?.content ?? ''), [document?.content])
  const tableOfContents = blocks.filter((block): block is Extract<ContentBlock, { kind: 'heading2' }> => block.kind === 'heading2')
  const pageAccent = type === 'terms' ? '#66BDF2' : '#F5A524'
  const pageGlow = type === 'terms' ? 'rgba(102,189,242,0.16)' : 'rgba(245,165,36,0.14)'
  const pageBorder = type === 'terms' ? 'rgba(102,189,242,0.26)' : 'rgba(245,165,36,0.24)'

  return (
    <main className="min-h-screen bg-[linear-gradient(180deg,#123A2F_0%,#0D314E_100%)] px-4 pb-10 pt-28 text-[#E8EAF2] sm:px-6 sm:pb-16 sm:pt-32">
      <AppHeaderShell
        logoHref={homeHref ?? (locale === 'en' ? '/en' : '/')}
        logoIconSize={36}
        logoTextSize="20px"
        rightSlot={
          <Link href={homeHref ?? (locale === 'en' ? '/en' : '/')} style={appHeaderActionLinkStyle(false)}>
            {homeLabel}
          </Link>
        }
      />

      <article
        className="relative mx-auto max-w-3xl overflow-hidden rounded-[28px] p-5 shadow-[0_32px_120px_rgba(3,8,20,0.45)] sm:p-8"
        style={{
          background: `linear-gradient(145deg, ${pageGlow} 0%, rgba(255,255,255,0.04) 60%, rgba(255,255,255,0.03) 100%)`,
          border: `1px solid ${pageBorder}`,
        }}
      >
        <div className="absolute right-0 top-0 h-28 w-28 rounded-full blur-3xl" style={{ backgroundColor: pageGlow }} />
        <div className="relative">
          <header className="border-b border-white/10 pb-7">
            <span
              className="mb-5 inline-flex rounded-full px-3 py-1 text-xs font-semibold tracking-[0.18em]"
              style={{
                color: pageAccent,
                backgroundColor: 'rgba(255,255,255,0.06)',
                border: '1px solid rgba(255,255,255,0.08)',
              }}
            >
              {type === 'terms' ? 'Terms' : 'Privacy'}
            </span>
            <h1 className="mb-3 text-3xl font-bold leading-tight sm:text-4xl">{document?.title || fallbackTitle}</h1>
            {document ? (
              <p className="text-sm text-[rgba(220,228,245,0.58)]">
                버전 v{document.version}
                {document.effective_at ? ` · ${locale === 'en' ? 'Effective' : '시행'} ${formatDate(document.effective_at, locale)}` : ''}
              </p>
            ) : null}
          </header>

          {document?.fallback_used && fallbackNoticeTitle && fallbackNoticeText ? (
            <section className="mt-8 rounded-2xl border border-[#F5A524]/30 bg-[#F5A524]/10 p-5 text-sm leading-7 text-[rgba(255,244,220,0.88)]">
              <h2 className="mb-2 text-base font-bold text-white">{fallbackNoticeTitle}</h2>
              <p>{fallbackNoticeText}</p>
            </section>
          ) : null}

          {error ? <p className="mt-8 text-sm text-[#FF9A9A]">{error}</p> : null}

          {!document && !error ? (
            <p className="mt-8 text-sm text-[rgba(220,228,245,0.62)]">{loadingText}</p>
          ) : null}

          {document && tableOfContents.length > 0 ? (
            <nav className="my-8 rounded-2xl border border-white/10 bg-white/[0.035] p-5" aria-label="문서 목차">
              <p className="mb-3 text-xs font-semibold uppercase tracking-[0.16em] text-[rgba(220,228,245,0.48)]">{tocLabel}</p>
              <ol className="grid gap-2 text-sm leading-6 text-[rgba(220,228,245,0.76)] sm:grid-cols-2">
                {tableOfContents.map((item) => (
                  <li key={item.id}>
                    <a href={`#${item.id}`} className="hover:text-white">
                      {item.text}
                    </a>
                  </li>
                ))}
              </ol>
            </nav>
          ) : null}

          {document ? (
            <div className="policy-content text-[15px] leading-8 text-[rgba(220,228,245,0.8)]">
              {blocks.map((block, index) => {
                if (block.kind === 'heading2') {
                  return (
                    <h2 key={`${block.id}-${index}`} id={block.id} className="scroll-mt-8 pt-8 text-xl font-bold leading-8 text-white">
                      {block.text}
                    </h2>
                  )
                }
                if (block.kind === 'heading3') {
                  return (
                    <h3 key={`${block.text}-${index}`} className="pt-5 text-base font-semibold leading-7 text-[#B8B3FF]">
                      {block.text}
                    </h3>
                  )
                }
                if (block.kind === 'list') {
                  return (
                    <ul key={`list-${index}`} className="my-4 list-disc space-y-2 pl-5 text-[rgba(220,228,245,0.78)]">
                      {block.items.map((item) => (
                        <li key={item}>{item}</li>
                      ))}
                    </ul>
                  )
                }
                return (
                  <p key={`paragraph-${index}`} className="my-4 whitespace-pre-line">
                    {block.text}
                  </p>
                )
              })}

              <section className="mt-12 border-t border-white/10 pt-6">
                <h2 className="text-base font-semibold text-white">{historyTitle}</h2>
                <p className="mt-3 text-sm leading-7 text-[rgba(220,228,245,0.62)]">
                  {historyText}
                </p>
              </section>
            </div>
          ) : null}
        </div>
      </article>
    </main>
  )
}
