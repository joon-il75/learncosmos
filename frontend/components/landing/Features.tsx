import Image from 'next/image';
import type { LandingPageCopy } from '@/lib/i18n/pages/landing';

export default function Features({ copy }: { copy: LandingPageCopy['features'] }) {
  return (
    <section id="features" className="relative z-10 w-full border-t border-white/8 bg-[linear-gradient(180deg,#0D314E_0%,#123A2F_100%)] px-6 py-20">
      <div className="mx-auto max-w-6xl">
        {/* 섹션 헤더 */}
        <div className="mb-14 text-center animate-on-scroll fade-up">
          <p className="mb-3 text-sm font-medium uppercase tracking-widest text-[#6AD2C1]">{copy.eyebrow}</p>
          <h2 className="text-3xl font-bold text-[#E8EAF2] md:text-4xl">{copy.title}</h2>
          <p className="mx-auto mt-4 max-w-2xl text-base leading-relaxed text-[rgba(200,210,235,0.72)]">
            {copy.description}
          </p>
        </div>

        {/* 기능 카드 그리드 */}
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
          {copy.items.map((feature, index) => (
            <div
              key={index}
              className="group flex min-h-[318px] flex-col overflow-hidden rounded-2xl border border-[#6AD2C1]/14 bg-[#06111F]/64 p-5 backdrop-blur-sm transition-all hover:border-[#6AD2C1]/32 hover:bg-[#123A2F]/76 animate-on-scroll fade-up"
              style={{ animationDelay: feature.delay }}
            >
              {/* 이미지 컨테이너 */}
              <div className="relative mb-5 h-28 w-full overflow-hidden rounded-xl border border-[#6AD2C1]/12 bg-[#06111F]/80">
                <Image
                  src={feature.img}
                  alt={feature.title}
                  fill
                  sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 25vw"
                  className="object-cover opacity-80 saturate-90 transition-transform duration-700 group-hover:scale-105"
                />
                <div className="absolute inset-0 bg-gradient-to-t from-[#081826]/80 via-[#081826]/20 to-transparent"></div>
              </div>
              {/* 텍스트 내용 */}
              <h3 className="mb-3 text-[1.02rem] font-semibold leading-snug text-[#F4F7FB]">{feature.title}</h3>
              <p className="text-sm leading-relaxed text-[rgba(200,210,235,0.68)]">{feature.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
