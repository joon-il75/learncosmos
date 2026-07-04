import type { Metadata } from 'next';
import AlphaInvitePageClient from '@/app/alpha-invite/AlphaInvitePageClient';

export const metadata: Metadata = {
  title: 'Alpha Invite | LearnCosmos',
  description: 'Enter your LearnCosmos alpha invite code.',
  alternates: {
    canonical: '/en/alpha-invite',
    languages: {
      ko: '/alpha-invite',
      en: '/en/alpha-invite',
    },
  },
};

export default function EnglishAlphaInvitePage() {
  return <AlphaInvitePageClient locale="en" />;
}
