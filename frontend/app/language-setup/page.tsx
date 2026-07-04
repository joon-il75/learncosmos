import type { Metadata } from 'next'
import LanguageSetupPageClient from './LanguageSetupPageClient'

export const metadata: Metadata = {
  title: 'Language Settings | LearnCosmos',
  description: 'LearnCosmos 언어 설정을 저장합니다.',
  alternates: {
    canonical: '/language-setup',
    languages: {
      ko: '/language-setup',
      en: '/en/language-setup',
    },
  },
}

export default function LanguageSetupPage() {
  return <LanguageSetupPageClient locale="ko" />
}
