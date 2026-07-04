import type { Metadata } from 'next'
import PolicyDocumentPage from '@/components/policies/PolicyDocumentPage'

export const metadata: Metadata = {
  title: 'Privacy Policy | LearnCosmos',
  description: 'LearnCosmos privacy policy.',
  alternates: {
    canonical: '/en/privacy',
    languages: {
      ko: '/privacy',
      en: '/en/privacy',
    },
  },
}

export default function EnglishPrivacyPage() {
  return (
    <PolicyDocumentPage
      type="privacy"
      locale="en"
      fallbackTitle="Privacy Policy"
      loadingText="Loading the privacy policy..."
      errorText="Could not load the privacy policy"
      homeLabel="Home"
      tocLabel="Contents"
      historyTitle="Revision History"
      historyText="The currently published document applies from the effective date shown above. Future changes will be announced in the service or through a separate notice."
      fallbackNoticeTitle="Korean source text"
      fallbackNoticeText="An English translation is not available yet. The Korean version is currently shown as the governing source text."
    />
  )
}
