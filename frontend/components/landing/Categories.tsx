'use client'

import type { LandingPageCopy } from '@/lib/i18n/pages/landing'

export default function Categories({ copy }: { copy: LandingPageCopy['categories'] }) {
  const handleCategoryClick = (query: string) => {
    window.dispatchEvent(new CustomEvent('learnweaver:category-select', { detail: { query } }))
  }

  return (
    <section id="categories" className="relative w-full overflow-hidden border-t border-white/5 bg-[linear-gradient(180deg,#222C4D_0%,#123A2F_100%)] px-6 py-20">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top,rgba(106,210,193,0.14),transparent_36%)]" />
      <div className="relative mx-auto max-w-6xl text-center z-10 animate-on-scroll fade-up">
        <p className="mb-3 text-sm font-medium uppercase tracking-widest text-[#6AD2C1]">{copy.eyebrow}</p>
        <h2 className="text-3xl font-bold md:text-4xl text-[#E8EAF2]">
          {copy.title}
        </h2>
        <p className="mx-auto mt-4 max-w-2xl text-base leading-relaxed text-[rgba(200,210,235,0.72)]">
          {copy.descriptionLines[0]}
          <br className="hidden md:block" />
          {' '}
          {copy.descriptionLines[1]}
        </p>

        <div className="mt-12 flex flex-wrap justify-center gap-3.5">
          {copy.items.map((category, index) => (
            <button
              key={index}
              type="button"
              className={`group min-h-[44px] rounded-full border px-5 py-2.5 text-sm font-medium backdrop-blur-sm transition-all hover:-translate-y-0.5 active:scale-[0.98] animate-on-scroll fade-up ${
                category.direct
                  ? 'border-[#FF6E61]/55 bg-[#FF6E61]/14 text-[#FFD2CE] hover:border-[#FF6E61]/75 hover:bg-[#FF6E61]/20'
                  : 'border-white/10 bg-[#081826]/58 text-[rgba(232,234,242,0.9)] hover:border-[#6AD2C1]/38 hover:bg-[#123A2F]/78'
              }`}
              style={{ animationDelay: `${index * 100}ms` }}
              onClick={() => handleCategoryClick(category.query)}
            >
              <span className="mr-2 group-hover:scale-110 transition-transform inline-block">
                {category.icon}
              </span>
              {category.name}
            </button>
          ))}
        </div>
      </div>
    </section>
  );
}
