import type { Metadata } from 'next';
import LandingPageClient from '@/components/landing/LandingPageClient';

const pageTitle = 'LearnCosmos - AI 기반 탐험형 자기주도 학습 플랫폼';
const pageDescription = 'LearnCosmos는 AI와 함께 학습 목표를 탐험 여정으로 설계하고, 코스 생성·탐험 계획·Explorer Diary를 통해 자기주도 학습을 지원하는 AI 학습 플랫폼입니다.';

export const metadata: Metadata = {
  title: pageTitle,
  description: pageDescription,
  alternates: {
    canonical: '/',
    languages: {
      ko: '/',
      en: '/en',
    },
  },
  openGraph: {
    title: pageTitle,
    description: pageDescription,
    url: '/',
    type: 'website',
    locale: 'ko_KR',
  },
}

export default function Home() {
  return <LandingPageClient locale="ko" preferUserLocale />;
}
