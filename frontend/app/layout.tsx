import type { Metadata } from 'next'
import Script from 'next/script'
import './globals.css'

const siteTitle = 'LearnCosmos - AI 기반 탐험형 자기주도 학습 플랫폼'
const siteDescription = 'LearnCosmos는 AI와 함께 학습 목표를 탐험 여정으로 설계하고, 코스 생성·탐험 계획·Explorer Diary를 통해 자기주도 학습을 지원하는 AI 학습 플랫폼입니다.'
const googleAnalyticsID = 'G-XQQLP382KZ'
const siteUrl = process.env.NEXT_PUBLIC_SITE_URL || 'https://learncosmos.co.kr'
const isProd = process.env.NODE_ENV === 'production'

export const metadata: Metadata = {
  title: {
    default: siteTitle,
    template: '%s | LearnCosmos',
  },
  description: siteDescription,
  metadataBase: new URL(siteUrl),
  applicationName: 'LearnCosmos',
  authors: [{ name: 'LearnCosmos', url: siteUrl }],
  creator: 'LearnCosmos',
  publisher: 'LearnCosmos',
  category: 'education',
  keywords: [
    'LearnCosmos',
    'self-directed learning',
    'AI learning journey',
    'learning goals',
    'learning records',
    'AI 학습',
    '자기주도 학습',
    '학습탐험',
    '학습 목표',
  ],
  alternates: {
    canonical: '/',
    languages: {
      ko: '/',
      en: '/en',
    },
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      'max-image-preview': 'large',
      'max-snippet': -1,
      'max-video-preview': -1,
    },
  },
  icons: {
    icon: [
      { url: '/favicon.ico' },
      { url: '/favicon.svg', type: 'image/svg+xml' },
      { url: '/favicon-16x16.png', sizes: '16x16', type: 'image/png' },
      { url: '/favicon-32x32.png', sizes: '32x32', type: 'image/png' },
      { url: '/favicon-48x48.png', sizes: '48x48', type: 'image/png' },
    ],
  },
  openGraph: {
    title: siteTitle,
    description: siteDescription,
    url: siteUrl,
    siteName: 'LearnCosmos',
    images: [{ url: '/images/hero_galaxy.webp', width: 1200, height: 630, alt: 'LearnCosmos learning journey galaxy' }],
    locale: 'ko_KR',
    alternateLocale: ['en_US'],
    type: 'website',
  },
  twitter: {
    card: 'summary_large_image',
    title: siteTitle,
    description: siteDescription,
    images: ['/images/hero_galaxy.webp'],
  },
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ko">
      <head>
        <meta name="ai-content-declaration" content="public landing pages and policy documents are available for indexing on learncosmos.co.kr; private learner data is not public" />
        <meta name="llm-summary" content="LearnCosmos is a Korean-first self-directed learning platform that helps learners set goals, generate AI-assisted learning journeys, collect content, keep records, and complete learning paths." />
        <link
          rel="preconnect"
          href="https://fonts.googleapis.com"
        />
        <link
          rel="preconnect"
          href="https://fonts.gstatic.com"
          crossOrigin="anonymous"
        />
        <link
          href="https://fonts.googleapis.com/css2?family=Nanum+Pen+Script&display=swap"
          rel="stylesheet"
        />
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{
            __html: JSON.stringify({
              '@context': 'https://schema.org',
              '@type': 'WebApplication',
              name: 'LearnCosmos',
              url: siteUrl,
              applicationCategory: 'EducationalApplication',
              operatingSystem: 'Web',
              inLanguage: ['ko', 'en'],
              description: siteDescription,
              offers: {
                '@type': 'Offer',
                price: '0',
                priceCurrency: 'KRW',
                availability: 'https://schema.org/InStock',
              },
              creator: {
                '@type': 'Organization',
                name: 'LearnCosmos',
                url: siteUrl,
              },
            }),
          }}
        />
      </head>
      <body className="bg-[#081826] text-[#E8EAF2] antialiased">
        {isProd ? (
          <>
            <Script
              src={`https://www.googletagmanager.com/gtag/js?id=${googleAnalyticsID}`}
              strategy="afterInteractive"
            />
            <Script id="google-analytics" strategy="afterInteractive">
              {`
                window.dataLayer = window.dataLayer || [];
                function gtag(){dataLayer.push(arguments);}
                gtag('js', new Date());
                gtag('config', '${googleAnalyticsID}');
              `}
            </Script>
          </>
        ) : null}
        {children}
      </body>
    </html>
  )
}
