import type { Metadata } from 'next'
import DemoLoginPageClient from '@/app/demo-login/DemoLoginPageClient'

export const metadata: Metadata = {
  title: 'Demo Login | LearnCosmos',
  description: 'Enter your LearnCosmos demo login code to access an alpha testing learner account.',
  robots: {
    index: false,
    follow: false,
  },
  alternates: {
    canonical: '/en/demo-login',
    languages: {
      ko: '/demo-login',
      en: '/en/demo-login',
    },
  },
}

export default function EnglishDemoLoginPage() {
  return <DemoLoginPageClient locale="en" />
}
