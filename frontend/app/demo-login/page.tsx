import type { Metadata } from 'next'
import DemoLoginPageClient from './DemoLoginPageClient'

export const metadata: Metadata = {
  title: 'Demo Login | LearnCosmos',
  description: 'LearnCosmos 데모 로그인 코드를 입력해 알파 테스트용 학습자 계정으로 접속합니다.',
  robots: {
    index: false,
    follow: false,
  },
  alternates: {
    canonical: '/demo-login',
    languages: {
      ko: '/demo-login',
      en: '/en/demo-login',
    },
  },
}

export default function DemoLoginPage() {
  return <DemoLoginPageClient locale="ko" />
}
