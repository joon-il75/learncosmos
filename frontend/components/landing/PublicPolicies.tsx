import Link from 'next/link'
import type { LandingPageCopy } from '@/lib/i18n/pages/landing'

const docStyles = [
  {
    accent: '#66BDF2',
    glow: 'rgba(102,189,242,0.16)',
    border: 'rgba(102,189,242,0.26)',
  },
  {
    accent: '#F5A524',
    glow: 'rgba(245,165,36,0.14)',
    border: 'rgba(245,165,36,0.24)',
  },
  {
    accent: '#8BE0C4',
    glow: 'rgba(139,224,196,0.14)',
    border: 'rgba(139,224,196,0.24)',
  },
]

export default function PublicPolicies({ copy }: { copy: LandingPageCopy['policies'] }) {
  return (
    <section className="border-t border-white/10 bg-[linear-gradient(180deg,#123A2F_0%,#0D314E_100%)] px-6 py-24">
      <div className="mx-auto max-w-6xl">
        <div className="mb-14 max-w-3xl">
          <span
            className="mb-4 inline-block rounded-full px-4 py-1.5 text-sm font-medium"
            style={{
              backgroundColor: 'rgba(102,189,242,0.12)',
              color: '#B9E5FF',
              border: '1px solid rgba(102,189,242,0.24)',
            }}
          >
            {copy.eyebrow}
          </span>
          <h2 className="mb-4 text-3xl font-bold leading-tight text-white md:text-4xl">
            {copy.title}
          </h2>
          <p className="text-lg leading-8 text-[rgba(200,210,235,0.7)]">
            {copy.descriptionLines.map((line) => (
              <span key={line} className="block">{line}</span>
            ))}
          </p>
        </div>

        <div className="grid gap-6 md:grid-cols-3">
          {copy.docs.map((doc, index) => {
            const style = docStyles[index] ?? docStyles[0]
            return (
            <Link
              key={doc.href}
              href={doc.href}
              className="group relative overflow-hidden rounded-[28px] p-8 transition-transform duration-200 hover:-translate-y-1"
              style={{
                background: `linear-gradient(145deg, ${style.glow} 0%, rgba(255,255,255,0.04) 60%, rgba(255,255,255,0.03) 100%)`,
                border: `1px solid ${style.border}`,
              }}
            >
              <div
                className="absolute right-0 top-0 h-28 w-28 rounded-full blur-3xl"
                style={{ backgroundColor: style.glow }}
              />
              <div className="relative">
                <span
                  className="mb-5 inline-flex rounded-full px-3 py-1 text-xs font-semibold tracking-[0.18em]"
                  style={{
                    color: style.accent,
                    backgroundColor: 'rgba(255,255,255,0.06)',
                    border: '1px solid rgba(255,255,255,0.08)',
                  }}
                >
                  {doc.badge}
                </span>
                <h3 className="mb-3 text-2xl font-bold text-white">{doc.title}</h3>
                <p className="mb-8 leading-7 text-[rgba(200,210,235,0.68)]">
                  {doc.description}
                </p>
                <div
                  className="inline-flex items-center gap-2 text-sm font-semibold"
                  style={{ color: style.accent }}
                >
                  {copy.linkLabel}
                  <svg
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="1.9"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    className="transition-transform duration-200 group-hover:translate-x-1"
                  >
                    <path d="M5 12h14" />
                    <path d="m12 5 7 7-7 7" />
                  </svg>
                </div>
              </div>
            </Link>
          )})}
        </div>

        <p className="mt-7 max-w-3xl text-sm leading-7 text-[rgba(200,210,235,0.58)]">
          {copy.noteLines.map((line) => (
            <span key={line} className="block">{line}</span>
          ))}
        </p>
      </div>
    </section>
  )
}
