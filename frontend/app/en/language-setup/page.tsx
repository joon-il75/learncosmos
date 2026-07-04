import type { Metadata } from 'next'
import LanguageSetupPageClient from '../../language-setup/LanguageSetupPageClient'

export const metadata: Metadata = {
  title: 'Language Setup | LearnCosmos',
  description: 'Choose your LearnCosmos interface and learning language preferences.',
  alternates: {
    canonical: '/en/language-setup',
    languages: {
      ko: '/language-setup',
      en: '/en/language-setup',
    },
  },
}

export default function EnglishLanguageSetupPage() {
  return <LanguageSetupPageClient locale="en" />
}
