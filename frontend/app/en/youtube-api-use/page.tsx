import type { Metadata } from 'next'
import Link from 'next/link'
import { getYoutubeApiUseCopy } from '@/lib/i18n/pages/youtube-api-use'

export const metadata: Metadata = {
  title: 'YouTube Data API Use | LearnCosmos',
  description: 'How LearnCosmos uses YouTube search to recommend educational learning materials.',
  alternates: {
    canonical: '/en/youtube-api-use',
  },
  openGraph: {
    title: 'YouTube Data API Use | LearnCosmos',
    description: 'How LearnCosmos uses YouTube search to recommend educational learning materials.',
    url: '/en/youtube-api-use',
    siteName: 'LearnCosmos',
    locale: 'en_US',
    type: 'website',
  },
}

export default function YoutubeApiUsePage() {
  const copy = getYoutubeApiUseCopy()

  return (
    <main className="min-h-screen bg-[#F6F8FB] text-[#122033]">
      <header className="border-b border-[#D8E0EA] bg-white/95">
        <div className="mx-auto flex w-full max-w-6xl items-center justify-between px-5 py-4">
          <Link href="/en" className="text-lg font-bold text-[#0D314E]">
            LearnCosmos
          </Link>
          <nav className="flex items-center gap-3 text-sm font-semibold">
            <Link href="/en" className="rounded-lg px-3 py-2 text-[#42526B] hover:bg-[#EEF3F8]">
              {copy.nav.home}
            </Link>
            <Link href="/en/login" className="rounded-lg bg-[#0D314E] px-4 py-2 text-white hover:bg-[#174B73]">
              {copy.nav.login}
            </Link>
          </nav>
        </div>
      </header>

      <section className="border-b border-[#D8E0EA] bg-white">
        <div className="mx-auto grid w-full max-w-6xl gap-8 px-5 py-12 md:grid-cols-[minmax(0,1fr)_320px] md:py-16">
          <div>
            <p className="mb-4 text-sm font-bold uppercase tracking-[0.18em] text-[#1C7D79]">{copy.hero.eyebrow}</p>
            <h1 className="max-w-3xl text-4xl font-bold leading-tight tracking-normal text-[#0D314E] md:text-5xl">
              {copy.hero.title}
            </h1>
            <p className="mt-6 max-w-3xl text-base leading-8 text-[#42526B] md:text-lg">{copy.hero.description}</p>
          </div>
          <dl className="grid gap-3 rounded-lg border border-[#D8E0EA] bg-[#F6F8FB] p-5">
            {copy.facts.map((fact) => (
              <div key={fact.label} className="border-b border-[#D8E0EA] pb-3 last:border-b-0 last:pb-0">
                <dt className="text-xs font-bold uppercase tracking-[0.14em] text-[#6B778C]">{fact.label}</dt>
                <dd className="mt-1 text-sm font-semibold leading-6 text-[#172B4D]">{fact.value}</dd>
              </div>
            ))}
          </dl>
        </div>
      </section>

      <div className="mx-auto grid w-full max-w-6xl gap-6 px-5 py-10 md:grid-cols-[minmax(0,1fr)_320px]">
        <div className="space-y-5">
          {copy.sections.map((section) => (
            <section key={section.title} className="rounded-lg border border-[#D8E0EA] bg-white p-6">
              <h2 className="text-xl font-bold text-[#0D314E]">{section.title}</h2>
              <div className="mt-4 space-y-3 text-sm leading-7 text-[#42526B]">
                {section.body.map((paragraph) => (
                  <p key={paragraph}>{paragraph}</p>
                ))}
              </div>
            </section>
          ))}
        </div>

        <aside className="space-y-5">
          <section className="rounded-lg border border-[#D8E0EA] bg-white p-5">
            <h2 className="text-lg font-bold text-[#0D314E]">{copy.policyFallback.title}</h2>
            <p className="mt-3 text-sm leading-7 text-[#42526B]">{copy.policyFallback.body}</p>
            <div className="mt-4 flex flex-col gap-2 text-sm font-semibold">
              <Link href="/terms" className="rounded-lg border border-[#C8D3DF] px-3 py-2 text-[#0D314E] hover:bg-[#EEF3F8]">
                {copy.policyFallback.terms}
              </Link>
              <Link href="/privacy" className="rounded-lg border border-[#C8D3DF] px-3 py-2 text-[#0D314E] hover:bg-[#EEF3F8]">
                {copy.policyFallback.privacy}
              </Link>
            </div>
          </section>

          <section className="rounded-lg border border-[#D8E0EA] bg-white p-5">
            <h2 className="text-lg font-bold text-[#0D314E]">{copy.contact.title}</h2>
            <p className="mt-3 text-sm leading-7 text-[#42526B]">{copy.contact.body}</p>
            <a className="mt-4 block text-sm font-bold text-[#1C7D79] underline underline-offset-4" href={`mailto:${copy.contact.email}`}>
              {copy.contact.email}
            </a>
          </section>
        </aside>
      </div>
    </main>
  )
}
