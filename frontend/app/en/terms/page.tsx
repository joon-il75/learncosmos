import type { Metadata } from 'next'
import PolicyDocumentPage from '@/components/policies/PolicyDocumentPage'

export const metadata: Metadata = {
  title: 'Terms of Service | LearnCosmos',
  description: 'LearnCosmos terms of service.',
  alternates: {
    canonical: '/en/terms',
    languages: {
      ko: '/terms',
      en: '/en/terms',
    },
  },
}

export default function EnglishTermsPage() {
  return (
    <PolicyDocumentPage
      type="terms"
      locale="en"
      fallbackTitle="Terms of Service"
      loadingText="Loading the terms document..."
      errorText="Could not load the terms document"
      homeLabel="Home"
      tocLabel="Contents"
      historyTitle="Revision History"
      historyText="The currently published document applies from the effective date shown above. Future changes will be announced in the service or through a separate notice."
      fallbackNoticeTitle="Korean source text"
      fallbackNoticeText="An English translation is not available yet. The Korean version is currently shown as the governing source text."
    />
  )
}
