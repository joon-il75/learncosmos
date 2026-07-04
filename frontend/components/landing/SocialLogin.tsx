'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import type { Locale } from '@/lib/i18n/locales';
import type { LandingPageCopy } from '@/lib/i18n/pages/landing';

const OAUTH_URLS: Record<string, string> = {
  google: '/api/v1/auth/google/login',
  kakao: '/api/v1/auth/kakao/login',
  naver: '/api/v1/auth/naver/login',
};

export default function SocialLogin({ copy, locale }: { copy: LandingPageCopy['socialLogin']; locale: Locale }) {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [showNaverNotice, setShowNaverNotice] = useState(false);
  const termsHref = locale === 'en' ? '/en/terms' : '/terms';
  const privacyHref = locale === 'en' ? '/en/privacy' : '/privacy';

  useEffect(() => {
    setIsLoggedIn(document.cookie.includes('is_logged_in=1'));
  }, []);

  const handleLogin = (provider: string) => {
    if (provider === 'naver') {
      setShowNaverNotice(true);
      return;
    }

    const params = new URLSearchParams({
      redirect_after: '/dashboard',
      locale,
    });
    window.location.href = `${OAUTH_URLS[provider]}?${params.toString()}`;
  };

  return (
    <section
      id="signup"
      className="relative w-full overflow-hidden border-t border-white/10 bg-[linear-gradient(145deg,#0D314E_0%,#1C7D79_50%,#CC5216_100%)] px-6 py-24 text-center"
    >
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_0%,rgba(255,255,255,0.16),transparent_34%)]" />
      <div className="relative z-10 mx-auto max-w-2xl">
        {/* Login icon */}
        <div className="flex justify-center mb-6">
          <div className="w-16 h-16 rounded-full bg-white/20 border border-white/35 flex items-center justify-center">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="32"
              height="32"
              viewBox="0 0 24 24"
              fill="none"
              stroke="#ffffff"
              strokeWidth="1.8"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
              <circle cx="12" cy="7" r="4" />
              <path d="M16 11l4 4-4 4" />
              <path d="M20 15H9" />
            </svg>
          </div>
        </div>

        <p className="mb-3 text-sm font-medium uppercase tracking-widest text-white/70">{copy.eyebrow}</p>
        <h2 className="mb-5 text-3xl font-bold leading-tight text-white md:whitespace-nowrap md:text-4xl">{copy.title}</h2>

        {isLoggedIn ? (
          /* 1. 로그인 후: 대시보드 이동 버튼 */
          <div className="mx-auto flex max-w-md flex-col items-center">
            <p className="mb-10 text-lg text-white/80">
              {copy.loggedInMessage}
            </p>

            <Link
              href="/dashboard"
              className="group relative flex w-full items-center justify-center overflow-hidden rounded-xl bg-white px-8 py-5 text-lg font-bold text-[#0B3D91] shadow-lg transition-all hover:scale-[1.02] hover:bg-opacity-90"
            >
              <span className="flex items-center gap-2">
                {copy.dashboardCta}
                <svg
                  width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"
                  className="transition-transform group-hover:translate-x-1"
                >
                  <path d="M5 12h14m-7-7 7 7-7 7"/>
                </svg>
              </span>
            </Link>

            <p className="mt-6 text-xs text-white/60">
              {copy.loggedInHint}
            </p>
          </div>
        ) : (
          /* 2. 로그인 전: 소셜 로그인 버튼 세트 (Fragment로 감싸서 에러 해결) */
          <div className="mx-auto max-w-md">
            <p className="mb-5 text-lg leading-8 text-white/82">
              {copy.signupDescriptionLines.map((line, index) => (
                <span key={line} className="block">
                  {line}
                </span>
              ))}
            </p>

            <p className="mb-5 text-sm leading-7 text-white/74">
              {copy.socialDescriptionLines.map((line) => (
                <span key={line} className="block">{line}</span>
              ))}
            </p>

            <p className="mb-8 rounded-2xl border border-white/14 bg-white/8 px-5 py-4 text-sm leading-7 text-white/76">
              {copy.aiNoticeLines.map((line) => (
                <span key={line} className="block">{line}</span>
              ))}
            </p>

            <div className="flex flex-col gap-4">
              {/* 구글 버튼 */}
              <button
                onClick={() => handleLogin('google')}
                className="flex min-h-14 w-full items-center gap-0 overflow-hidden rounded-xl border border-[#081826]/20 bg-white shadow-[0_16px_32px_rgba(3,8,20,0.2)] transition-colors hover:bg-gray-50"
              >
                <div className="w-14 h-14 flex items-center justify-center bg-white border-r border-[#081826]/10 flex-shrink-0">
                  <svg width="22" height="22" viewBox="0 0 24 24">
                    <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
                    <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
                    <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81.38z" fill="#FBBC05"/>
                    <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
                  </svg>
                </div>
                <span className="flex-1 text-center pr-14 text-[#111111] font-medium">{copy.providers.google}</span>
              </button>

              {/* 카카오 버튼 */}
              <button
                onClick={() => handleLogin('kakao')}
                className="flex min-h-14 w-full items-center gap-0 overflow-hidden rounded-xl border border-[#c9a800]/20 bg-[#FEE500] transition-colors hover:bg-[#FADA00]"
              >
                <div className="w-14 h-14 flex items-center justify-center bg-[#FEE500]/80 border-r border-[#191919]/10 flex-shrink-0">
                  <svg width="22" height="22" viewBox="0 0 24 24">
                    <path fillRule="evenodd" clipRule="evenodd" d="M12 2C6.48 2 2 5.82 2 10.5c0 3.01 1.87 5.65 4.7 7.18L5.5 22l4.82-3.21c.55.07 1.1.11 1.68.11 5.52 0 10-3.82 10-8.5S17.52 2 12 2z" fill="#191919"/>
                  </svg>
                </div>
                <span className="flex-1 text-center pr-14 text-[#191919] font-medium">{copy.providers.kakao}</span>
              </button>

              {/* 네이버 버튼 */}
              <button
                onClick={() => handleLogin('naver')}
                className="flex min-h-14 w-full items-center gap-0 overflow-hidden rounded-xl border border-[#029a46]/20 bg-[#03C75A] transition-colors hover:bg-[#02b852]"
              >
                <div className="w-14 h-14 flex items-center justify-center bg-[#02b050] border-r border-white/10 flex-shrink-0">
                  <svg width="22" height="22" viewBox="0 0 24 24">
                    <rect width="24" height="24" rx="4" fill="#03C75A"/>
                    <path d="M13.37 12.28L10.43 7H7.5v10h3.13v-5.28L13.57 17H16.5V7h-3.13v5.28z" fill="white"/>
                  </svg>
                </div>
                <span className="flex-1 text-center pr-14 text-white font-medium">{copy.providers.naver}</span>
              </button>
            </div>

            <p className="mt-7 text-xs leading-6 text-white/68">
              {copy.policyNotice.before}
              <br />
              <Link href={termsHref} className="underline underline-offset-4 hover:text-white">
                {copy.policyNotice.terms}
              </Link>{' '}
              {copy.policyNotice.middle}{' '}
              <Link href={privacyHref} className="underline underline-offset-4 hover:text-white">
                {copy.policyNotice.privacy}
              </Link>
              {' '}{copy.policyNotice.after}
            </p>
          </div>
        )}
      </div>
      {showNaverNotice && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-[#081826]/70 px-5 backdrop-blur-sm"
          role="dialog"
          aria-modal="true"
          aria-labelledby="landing-naver-login-notice-title"
        >
          <div className="w-full max-w-sm rounded-2xl border border-white/16 bg-[#10283B] p-6 text-center shadow-[0_28px_72px_rgba(0,0,0,0.42)]">
            <h2 id="landing-naver-login-notice-title" className="text-lg font-bold text-white">
              {copy.naverNotice.title}
            </h2>
            <p className="mt-3 text-sm leading-7 text-white/76">
              {copy.naverNotice.bodyLines.map((line) => (
                <span key={line} className="block">{line}</span>
              ))}
            </p>
            <button
              type="button"
              onClick={() => setShowNaverNotice(false)}
              className="mt-6 min-h-11 w-full rounded-xl bg-white px-5 py-3 text-sm font-bold text-[#0D314E] transition-colors hover:bg-white/90"
            >
              {copy.naverNotice.confirm}
            </button>
          </div>
        </div>
      )}
    </section>
  );
}
