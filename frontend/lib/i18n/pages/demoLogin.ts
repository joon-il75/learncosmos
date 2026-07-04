import type { Locale } from '@/lib/i18n/locales'

export type DemoLoginPageCopy = {
  locale: Locale
  logoHref: string
  eyebrow: string
  title: string
  description: string
  codeLabel: string
  codePlaceholder: string
  submit: string
  submitting: string
  backToLogin: string
  noticeLines: string[]
  errors: Record<string, string>
}

const demoLoginCopy: Record<Locale, DemoLoginPageCopy> = {
  ko: {
    locale: 'ko',
    logoHref: '/',
    eyebrow: 'Demo Login',
    title: 'LearnCosmos 데모 계정',
    description: '운영자가 전달한 데모 로그인 코드를 입력하면 알파 테스트용 학습자 계정으로 접속합니다.',
    codeLabel: 'Demo Code',
    codePlaceholder: 'LW-DEMO-01-X7K9Q2',
    submit: 'Continue with Demo Code',
    submitting: 'Checking...',
    backToLogin: 'Social Login으로 돌아가기',
    noticeLines: [
      '데모 계정은 테스트 목적의 학습 기록을 만들 수 있습니다.',
      '로그인 코드는 전달받은 사람만 사용해 주세요.',
    ],
    errors: {
      demo_login_disabled: '현재 데모 로그인이 열려 있지 않습니다.',
      invalid_demo_code: '데모 로그인 코드를 확인해 주세요.',
      expired_demo_code: '만료된 데모 로그인 코드입니다.',
      disabled_demo_account: '비활성화된 데모 계정입니다.',
      demo_account_unavailable: '데모 계정을 사용할 수 없습니다.',
      invalid_request: '데모 로그인 코드를 입력해 주세요.',
      unknown: '데모 로그인에 실패했습니다.',
    },
  },
  en: {
    locale: 'en',
    logoHref: '/en',
    eyebrow: 'Demo Login',
    title: 'LearnCosmos Demo Account',
    description: 'Enter the demo code provided by the LearnCosmos team to access an alpha testing learner account.',
    codeLabel: 'Demo Code',
    codePlaceholder: 'LW-DEMO-01-X7K9Q2',
    submit: 'Continue with Demo Code',
    submitting: 'Checking...',
    backToLogin: 'Back to Social Login',
    noticeLines: [
      'Demo accounts can create test learning records.',
      'Please use the demo code only if it was assigned to you.',
    ],
    errors: {
      demo_login_disabled: 'Demo login is not currently enabled.',
      invalid_demo_code: 'Please check your demo login code.',
      expired_demo_code: 'This demo login code has expired.',
      disabled_demo_account: 'This demo account is disabled.',
      demo_account_unavailable: 'This demo account is unavailable.',
      invalid_request: 'Please enter your demo login code.',
      unknown: 'Demo login failed.',
    },
  },
}

export function getDemoLoginPageCopy(locale: Locale): DemoLoginPageCopy {
  return demoLoginCopy[locale]
}
