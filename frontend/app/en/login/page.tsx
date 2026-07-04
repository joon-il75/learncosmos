import type { Metadata } from 'next'
import LoginPageClient from '@/app/(auth)/login/LoginPageClient'

export const metadata: Metadata = {
  title: 'Log in | LearnCosmos',
  description: 'Log in to LearnCosmos with a social account and continue your learning journey.',
  alternates: {
    canonical: '/en/login',
    languages: {
      ko: '/login',
      en: '/en/login',
    },
  },
}

export default function EnglishLoginPage() {
  return <LoginPageClient locale="en" />
}
