import Link from 'next/link'
import AppHeaderShell from '@/components/common/AppHeaderShell'

const notices = [
  {
    name: 'Tiptap',
    purpose: '연구지점 학습 자료 입력용 WYSIWYG 에디터',
    license: 'MIT',
    copyright: 'Tiptap GmbH',
    packages: [
      '@tiptap/react 3.22.5',
      '@tiptap/starter-kit 3.22.5',
      '@tiptap/core 3.22.5',
      '@tiptap/extension-image 3.22.5',
      '@tiptap/extension-underline 3.22.5',
      '@tiptap/extension-highlight 3.22.5',
      '@tiptap/extension-superscript 3.22.5',
      '@tiptap/extension-subscript 3.22.5',
      '@tiptap/extension-text-align 3.22.5',
      '@tiptap/extension-table 3.22.5',
      '@tiptap/extension-table-row 3.22.5',
      '@tiptap/extension-table-cell 3.22.5',
      '@tiptap/extension-table-header 3.22.5',
    ],
    url: 'https://github.com/ueberdosis/tiptap/blob/main/LICENSE.md',
  },
  {
    name: 'ProseMirror',
    purpose: 'Tiptap이 사용하는 에디터 기반 엔진',
    license: 'MIT',
    copyright: 'ProseMirror contributors',
    packages: [
      'prosemirror-state 1.4.4',
      'prosemirror-view 1.41.8',
      'prosemirror-model 1.25.4',
      'prosemirror-transform 1.12.0',
      'prosemirror-schema-list 1.5.1',
    ],
    url: 'https://github.com/ProseMirror/prosemirror',
  },
  {
    name: 'Tiptap UI Components',
    purpose: '연구지점 Tiptap 툴바 아이콘 표시',
    license: 'MIT',
    copyright: 'Tiptap',
    packages: ['tiptap-ui-components icons'],
    url: 'https://github.com/ueberdosis/tiptap-ui-components/blob/main/LICENSE',
  },
  {
    name: 'Pretendard',
    purpose: '공개 메인 및 학습자-facing 화면 본문/설명/말풍선 타이포그래피',
    license: 'SIL Open Font License 1.1',
    copyright: 'Pretendard Project Authors',
    packages: ['PretendardVariable.woff2'],
    url: 'https://github.com/orioncactus/pretendard/blob/main/LICENSE',
  },
  {
    name: 'Gmarket Sans',
    purpose: '공개 메인 버튼/네비 포인트 타이포그래피',
    license: 'Gmarket Sans official font license',
    copyright: 'Gmarket',
    packages: ['GmarketSansTTFMedium.ttf'],
    url: 'https://corp.gmarket.com/fonts',
  },
]

const headerActionStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '40px',
  padding: '0 16px',
  borderRadius: '999px',
  border: '1px solid rgba(255,255,255,0.14)',
  background: 'rgba(255,255,255,0.08)',
  color: '#E8EAF2',
  textDecoration: 'none',
  fontSize: '14px',
  fontWeight: 700,
} as const

export default function OpenSourcePage() {
  return (
    <main className="min-h-screen bg-[linear-gradient(180deg,#123A2F_0%,#0D314E_100%)] px-4 pb-10 pt-28 text-[#E8EAF2] sm:px-6 sm:pb-16 sm:pt-32">
      <AppHeaderShell
        logoHref="/"
        logoIconSize={36}
        logoTextSize="20px"
        rightSlot={
          <Link href="/" style={headerActionStyle}>
            홈으로
          </Link>
        }
      />

      <article
        className="relative mx-auto max-w-4xl overflow-hidden rounded-[28px] p-5 shadow-[0_32px_120px_rgba(3,8,20,0.45)] sm:p-8"
        style={{
          background: 'linear-gradient(145deg, rgba(102,189,242,0.14) 0%, rgba(255,255,255,0.04) 60%, rgba(255,255,255,0.03) 100%)',
          border: '1px solid rgba(102,189,242,0.24)',
        }}
      >
        <div className="relative">
          <header className="border-b border-white/10 pb-7">
            <span className="mb-5 inline-flex rounded-full border border-white/10 bg-white/[0.06] px-3 py-1 text-xs font-semibold tracking-[0.18em] text-[#B9E5FF]">
              Open Source
            </span>
            <h1 className="mb-3 text-3xl font-bold leading-tight sm:text-4xl">오픈소스 라이선스</h1>
            <p className="max-w-2xl text-sm leading-7 text-[rgba(220,228,245,0.66)]">
              LearnCosmos 서비스에 포함된 주요 오픈소스 소프트웨어의 라이선스 고지입니다.
              각 라이브러리의 사용 조건은 해당 프로젝트의 라이선스 원문을 따릅니다.
            </p>
          </header>

          <div className="mt-8 grid gap-5">
            {notices.map((notice) => (
              <section key={notice.name} className="rounded-[20px] border border-white/10 bg-white/[0.035] p-5 sm:p-6">
                <div className="flex flex-wrap items-start justify-between gap-4">
                  <div>
                    <h2 className="text-2xl font-bold text-white">{notice.name}</h2>
                    <p className="mt-2 text-sm leading-7 text-[rgba(220,228,245,0.7)]">{notice.purpose}</p>
                  </div>
                  <span className="rounded-full border border-[#66BDF2]/30 bg-[#66BDF2]/10 px-3 py-1 text-xs font-bold text-[#B9E5FF]">
                    {notice.license}
                  </span>
                </div>

                <dl className="mt-6 grid gap-4 text-sm sm:grid-cols-2">
                  <div>
                    <dt className="font-semibold text-[rgba(220,228,245,0.52)]">저작권자</dt>
                    <dd className="mt-1 text-[rgba(220,228,245,0.82)]">{notice.copyright}</dd>
                  </div>
                  <div>
                    <dt className="font-semibold text-[rgba(220,228,245,0.52)]">라이선스 원문</dt>
                    <dd className="mt-1">
                      <a href={notice.url} target="_blank" rel="noreferrer" className="text-[#B9E5FF] underline underline-offset-4 hover:text-white">
                        원문 보기
                      </a>
                    </dd>
                  </div>
                </dl>

                <div className="mt-5">
                  <h3 className="text-sm font-semibold text-[rgba(220,228,245,0.52)]">포함 패키지</h3>
                  <ul className="mt-3 grid gap-2 text-sm text-[rgba(220,228,245,0.78)] sm:grid-cols-2">
                    {notice.packages.map((pkg) => (
                      <li key={pkg} className="rounded-xl border border-white/10 bg-black/10 px-3 py-2 font-mono text-xs">
                        {pkg}
                      </li>
                    ))}
                  </ul>
                </div>
              </section>
            ))}
          </div>

          <section className="mt-8 rounded-[20px] border border-white/10 bg-white/[0.035] p-5 text-sm leading-7 text-[rgba(220,228,245,0.72)] sm:p-6">
            LearnCosmos는 현재 Tiptap Pro, Cloud, AI, Collaboration, Comments, Versioning 계열 기능을 사용하지 않습니다.
            해당 기능을 도입할 경우 별도 라이선스 검토를 먼저 진행합니다.
          </section>
        </div>
      </article>
    </main>
  )
}
