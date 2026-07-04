import type { LandingPageCopy } from '@/lib/i18n/pages/landing';

export default function HowItWorks({ copy }: { copy: LandingPageCopy['howItWorks'] }) {
  return (
    <section id="how-it-works" className="relative w-full bg-[linear-gradient(180deg,#123A2F_0%,#0D314E_100%)] px-6 py-20">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_20%_0%,rgba(245,165,36,0.12),transparent_32%)]" />
      <div className="pointer-events-none absolute inset-x-0 top-0 h-px bg-gradient-to-r from-transparent via-[#F5A524]/30 to-transparent" />

      <div className="relative z-10 mx-auto grid max-w-6xl gap-12 lg:grid-cols-[0.82fr_1.18fr] lg:items-start">
        <div className="animate-on-scroll fade-up lg:sticky lg:top-24">
          <p className="mb-3 text-sm font-medium uppercase tracking-widest text-[#F5A524]">{copy.eyebrow}</p>
          <h2 className="text-3xl font-bold leading-tight text-[#E8EAF2] md:text-4xl">{copy.title}</h2>
          <div className="mt-5 grid gap-3 text-base leading-8 text-[rgba(220,228,245,0.72)]">
            {copy.descriptionLines.map((line) => (
              <p key={line} className="m-0">{line}</p>
            ))}
          </div>
        </div>

        <ol className="relative grid gap-6 before:absolute before:left-[19px] before:top-3 before:hidden before:h-[calc(100%-24px)] before:w-px before:bg-gradient-to-b before:from-[#F5A524]/60 before:via-[#6AD2C1]/36 before:to-transparent sm:before:block">
          {copy.steps.map((item, index) => (
            <li key={item.step} className="animate-on-scroll fade-up relative grid gap-4 sm:grid-cols-[40px_1fr]" style={{ animationDelay: `${index * 120}ms` }}>
              <div className="relative z-10 grid h-10 w-10 place-items-center rounded-full bg-[#F5A524] text-sm font-black text-[#231705] shadow-[0_12px_30px_rgba(245,165,36,0.24)]">
                {item.step}
              </div>
              <div className="min-h-[112px] rounded-[24px] bg-white/[0.055] px-5 py-5 shadow-[inset_0_1px_0_rgba(255,255,255,0.08)] backdrop-blur-sm">
                <h3 className="text-lg font-black leading-snug text-white">{item.title}</h3>
                <p className="mt-2 text-sm leading-7 text-white/68">{item.desc}</p>
              </div>
            </li>
          ))}
        </ol>
      </div>
    </section>
  );
}
