import type { LoginPageCopy } from './types'

export const loginKo: LoginPageCopy = {
  locale: 'ko',
  logoHref: '/',
  eyebrow: 'Login',
  title: '첫 학습탐험을 이어가세요',
  description: 'LearnCosmos 계정이 없으면 같은 Social Login으로 바로 회원가입이 진행됩니다.',
  aiNoticeLines: ['API 키가 없어도 먼저 시작할 수 있고,', '나중에 내 AI 키를 연결해 사용할 수도 있습니다.'],
  providers: {
    google: 'Continue with Google',
    kakao: 'Continue with Kakao',
    naver: 'Continue with Naver',
  },
  policyNotice: {
    before: '가입 후 서비스 이용 전에 필수 약관 동의가 진행됩니다.',
    terms: '이용약관',
    middle: '과',
    privacy: '개인정보처리방침',
    after: '을 확인할 수 있습니다.',
  },
  naverNotice: {
    title: 'Naver login is not ready',
    bodyLines: ['현재 개발 중인 기능이라 Naver Login은 아직 사용할 수 없습니다.', 'Google 또는 Kakao Login을 이용해 주세요.'],
    confirm: 'OK',
  },
}
