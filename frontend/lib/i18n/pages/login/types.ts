export type LoginPageCopy = {
  locale: 'ko' | 'en'
  logoHref: string
  eyebrow: string
  title: string
  description: string
  aiNoticeLines: string[]
  providers: {
    google: string
    kakao: string
    naver: string
  }
  policyNotice: {
    before: string
    terms: string
    middle: string
    privacy: string
    after: string
  }
  naverNotice: {
    title: string
    bodyLines: string[]
    confirm: string
  }
}
