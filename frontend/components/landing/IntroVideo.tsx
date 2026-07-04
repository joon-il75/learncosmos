import type { LandingPageCopy } from '@/lib/i18n/pages/landing';

export default function IntroVideo({ copy }: { copy: LandingPageCopy['introVideo'] }) {
  return (
    <section id="intro-video" className="landingIntroVideoSection relative z-10 w-full px-6 py-16 sm:py-20">
      <div className="pointer-events-none absolute inset-x-0 top-0 h-px bg-gradient-to-r from-transparent via-[#6AD2C1]/28 to-transparent" />
      <div className="mx-auto grid max-w-6xl gap-9 lg:grid-cols-[0.82fr_1.18fr] lg:items-center">
        <div className="animate-on-scroll fade-up">
          <p className="mb-3 text-sm font-bold uppercase tracking-widest text-[#2C91B8]">{copy.eyebrow}</p>
          <h2 className="text-3xl font-black leading-tight text-[#173B4E] md:text-4xl">{copy.title}</h2>
          <div className="mt-5 grid gap-3 text-base leading-8 text-[#41677B]">
            {copy.descriptionLines.map((line) => (
              <p key={line} className="m-0">{line}</p>
            ))}
          </div>
          <a
            href="/alpha"
            className="mt-7 inline-flex h-12 items-center justify-center rounded-[999px] bg-[#6AD2C1] px-5 text-sm font-black text-[#061525] shadow-[0_14px_34px_rgba(106,210,193,0.16)] transition hover:-translate-y-0.5 hover:bg-[#8CE8D8] focus:outline-none focus:ring-2 focus:ring-[#6AD2C1]/70"
          >
            {copy.alphaCta}
          </a>
        </div>

        <div className="animate-on-scroll fade-up relative overflow-hidden rounded-[28px] bg-[radial-gradient(circle_at_20%_12%,rgba(106,210,193,0.24),transparent_32%),linear-gradient(135deg,rgba(7,19,35,0.98),rgba(2,6,23,0.96))] p-4 shadow-[0_30px_90px_rgba(0,0,0,0.38),inset_0_1px_0_rgba(255,255,255,0.08)]">
          <div className="relative aspect-video w-full overflow-hidden rounded-[22px] bg-[linear-gradient(135deg,rgba(15,47,70,0.9),rgba(6,15,29,0.98))]">
            <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_42%,rgba(248,213,116,0.20),transparent_30%),linear-gradient(120deg,transparent_0%,rgba(106,210,193,0.08)_44%,transparent_60%)]" />
            <div className="absolute left-1/2 top-1/2 grid -translate-x-1/2 -translate-y-1/2 place-items-center gap-4 text-center">
              <div className="grid h-20 w-20 place-items-center rounded-full bg-[#6AD2C1] text-[#061525] shadow-[0_0_44px_rgba(106,210,193,0.34)]">
                <span className="ml-1 text-3xl font-black" aria-hidden="true">▶</span>
              </div>
              <div className="grid gap-2">
                <strong className="text-lg font-black text-white md:text-xl">{copy.videoTitle}</strong>
                <span className="text-sm font-semibold text-[rgba(220,236,245,0.72)]">{copy.placeholderLabel}</span>
              </div>
            </div>
            <div className="absolute bottom-4 left-4 right-4 h-1.5 overflow-hidden rounded-full bg-white/10">
              <div className="h-full w-1/3 rounded-full bg-[#6AD2C1]/70" />
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
