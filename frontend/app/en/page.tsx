import type { Metadata } from 'next'
import LandingPageClient from '@/components/landing/LandingPageClient'

export const metadata: Metadata = {
  title: 'LearnCosmos | Self-directed learning journeys',
  description: 'Start a self-directed learning journey with goal setup, content discovery, records, and completion.',
  alternates: {
    canonical: '/en',
    languages: {
      ko: '/',
      en: '/en',
    },
  },
  openGraph: {
    title: 'LearnCosmos | Self-directed learning journeys',
    description: 'Create AI-assisted learning journeys with goal setup, content discovery, records, and completion tracking.',
    url: '/en',
    type: 'website',
    locale: 'en_US',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'LearnCosmos | Self-directed learning journeys',
    description: 'Create AI-assisted learning journeys with goal setup, content discovery, records, and completion tracking.',
    images: ['/images/hero_galaxy.webp'],
  },
}

export default function EnglishIndexPage() {
  return <LandingPageClient locale="en" />
}
