import type { LoginPageCopy } from './types'

export const loginEn: LoginPageCopy = {
  locale: 'en',
  logoHref: '/en',
  eyebrow: 'Log in',
  title: 'Continue your first learning journey',
  description: 'If you do not have a LearnCosmos account, social login will create one for you.',
  aiNoticeLines: ['You can start without an API key,', 'and connect your own AI key later.'],
  providers: {
    google: 'Continue with Google',
    kakao: 'Continue with Kakao',
    naver: 'Continue with Naver',
  },
  policyNotice: {
    before: 'Required policy consent is requested before using the service.',
    terms: 'Terms',
    middle: 'and',
    privacy: 'Privacy Policy',
    after: 'are available before sign-up.',
  },
  naverNotice: {
    title: 'Naver login is not ready',
    bodyLines: ['Naver login is still in development.', 'Please use Google or Kakao for now.'],
    confirm: 'OK',
  },
}
