import type { Metadata } from 'next'
import AgreementsPageClient from '@/app/agreements/AgreementsPageClient'

export const metadata: Metadata = {
  title: 'Required Consent | LearnCosmos',
  description: 'Review and agree to required LearnCosmos policy documents.',
}

export default function EnglishAgreementsPage() {
  return <AgreementsPageClient locale="en" />
}
