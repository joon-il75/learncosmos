import type { Metadata } from 'next';

const formUrl = 'https://forms.gle/G2G6kZ1WfF61iiUi7';
const formEmbedUrl = 'https://docs.google.com/forms/d/e/1FAIpQLSdwFJ3lDdn5Hsu6IDo1T8QZFgu5qsfwMnnbYG9LknFRr8QdEw/viewform?embedded=true';

export const metadata: Metadata = {
  title: 'LearnCosmos 알파 테스터 2차 모집',
  description: 'LearnCosmos 알파 테스트 참여 신청 페이지입니다. AI 기반 탐험형 자기주도 학습 경험을 먼저 체험할 알파 테스터를 모집합니다.',
  alternates: {
    canonical: '/alpha',
  },
  openGraph: {
    title: 'LearnCosmos 알파 테스터 2차 모집',
    description: '30명 한정으로 LearnCosmos 알파 테스트 참여자를 모집합니다.',
    url: 'https://learncosmos.co.kr/alpha',
    siteName: 'LearnCosmos',
    type: 'website',
  },
};

const experienceItems = [
  'AI가 목표를 하나의 탐험 여정으로 설계',
  'Planet Map과 행성 텍스처의 진행 변화',
  '탐험계획, 탐험일지, Artifact 기록',
  '학습 진행도에 따른 다이어리 배경 변화',
  '완주 중심 자기주도 학습 경험',
  '피드백 반영을 통한 빠른 개선',
];

export default function AlphaPage() {
  return (
    <main className="min-h-screen bg-[#07111F] text-white">
      <section
        className="relative overflow-hidden border-b border-white/10 bg-cover bg-center"
        style={{ backgroundImage: "linear-gradient(180deg, rgba(7,17,31,0.86), rgba(7,17,31,0.96)), url('/images/hero_galaxy.webp')" }}
      >
        <div className="mx-auto max-w-5xl px-5 py-5 sm:px-6">
          <a
            href="/"
            className="inline-flex h-10 items-center rounded-[8px] border border-white/15 bg-white/[0.06] px-4 text-sm font-bold text-[rgba(232,234,242,0.86)] transition hover:-translate-y-0.5 hover:border-[#6AD2C1]/55 hover:bg-[#6AD2C1]/10 hover:text-[#8CE8D8] focus:outline-none focus:ring-2 focus:ring-[#6AD2C1]/70"
          >
            ← 메인으로 돌아가기
          </a>
        </div>
        <div className="mx-auto max-w-5xl px-5 pb-16 pt-8 sm:px-6 sm:pb-20 sm:pt-10">
          <div className="mx-auto max-w-4xl text-center">
            <p className="mb-4 inline-flex rounded-full border border-[#6AD2C1]/35 bg-[#6AD2C1]/10 px-4 py-2 text-sm font-bold text-[#8CE8D8]">
              개발자 · 교육자 · 자기주도 학습자 우선
            </p>
            <h1 className="text-4xl font-black leading-tight tracking-normal text-white sm:text-5xl lg:text-6xl">
              LearnCosmos <span className="text-[#8CE8D8]">알파 테스터 2차</span> 모집
            </h1>
            <p className="mt-5 text-xl font-semibold text-[rgba(232,234,242,0.72)] sm:text-2xl">
              2026년 7월 말
            </p>
          </div>
        </div>
      </section>

      <div className="mx-auto max-w-5xl px-5 py-12 sm:px-6 sm:py-16">
        <section className="rounded-[8px] border border-white/10 bg-[#0D1B2A] p-6 shadow-[0_18px_50px_rgba(0,0,0,0.28)] sm:p-8">
          <h2 className="text-center text-2xl font-black text-white sm:text-3xl">알파에서 경험할 수 있는 것</h2>
          <div className="mt-8 grid gap-4 md:grid-cols-2">
            {experienceItems.map((item) => (
              <div key={item} className="flex gap-3 rounded-[8px] border border-white/10 bg-white/[0.04] p-4 text-base leading-7 text-[rgba(232,234,242,0.84)]">
                <span className="mt-1 inline-flex h-5 w-5 flex-none items-center justify-center rounded-full bg-[#6AD2C1] text-xs font-black text-[#061525]">
                  ✓
                </span>
                <span>{item}</span>
              </div>
            ))}
          </div>
        </section>

        <section className="mt-10 overflow-hidden rounded-[8px] border border-white/10 bg-white shadow-[0_24px_70px_rgba(0,0,0,0.34)]">
          <iframe
            src={formEmbedUrl}
            width="100%"
            height="950"
            className="block w-full border-0"
            title="LearnCosmos Alpha Application"
          />
        </section>

        <div className="mt-8 text-center text-sm leading-7 text-[rgba(232,234,242,0.62)] sm:text-base">
          <p>
            신청 후 24시간 이내에 <strong className="text-[#F8D574]">Invite Code + 데모 계정</strong>을 이메일로 보내드립니다.
          </p>
          <p>
            설문지가 보이지 않으면{' '}
            <a href={formUrl} target="_blank" rel="noreferrer" className="font-bold text-[#8CE8D8] underline underline-offset-4">
              신청 폼을 새 창에서 열어 주세요
            </a>
            .
          </p>
          <p>문의: learnweavr@gmail.com</p>
          <p className="mt-6">
            <a
              href="/"
              className="inline-flex h-11 items-center justify-center rounded-[8px] border border-white/15 bg-white/[0.06] px-5 text-sm font-bold text-[rgba(232,234,242,0.86)] transition hover:-translate-y-0.5 hover:border-[#6AD2C1]/55 hover:bg-[#6AD2C1]/10 hover:text-[#8CE8D8] focus:outline-none focus:ring-2 focus:ring-[#6AD2C1]/70"
            >
              메인으로 돌아가기
            </a>
          </p>
        </div>
      </div>
    </main>
  );
}
