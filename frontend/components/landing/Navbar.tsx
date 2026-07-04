'use client';

import Link from 'next/link';
import AppHeaderShell from '@/components/common/AppHeaderShell';
import LearnerHeaderActions, { type LearnerHeaderActionsCopy } from '@/components/common/LearnerHeaderActions';
import { useState, useEffect } from 'react';
import type { LandingPageCopy } from '@/lib/i18n/pages/landing';

const learnerHeaderActionsCopy: Record<LandingPageCopy['locale'], LearnerHeaderActionsCopy> = {
  ko: {
    userMenu: {
      fallbackName: '학습자',
      profile: '학습자정보',
      points: '포인트',
      aiSettings: 'AI 설정',
      logout: 'Logout',
    },
    compactMenuAriaLabel: '학습자 메뉴 열기',
    compactPointsLoading: '포인트 확인 중',
    compactPointsOwned: (points) => `보유 ${points}pt`,
    logoutModalMessage: '지금 계정에서 나가면 다시 학습을 이어가려면 Login이 필요합니다. 계속 진행할까요?',
    logoutCancel: '취소',
    logout: 'Logout',
    loggingOut: 'Logging out...',
  },
  en: {
    userMenu: {
      fallbackName: 'Learner',
      profile: 'Profile',
      points: 'Points',
      aiSettings: 'AI Settings',
      logout: 'Logout',
    },
    compactMenuAriaLabel: 'Open learner menu',
    compactPointsLoading: 'Checking points',
    compactPointsOwned: (points) => `${points}pt available`,
    logoutModalMessage: 'If you log out now, you will need to Login again to continue learning. Continue?',
    logoutCancel: 'Cancel',
    logout: 'Logout',
    loggingOut: 'Logging out...',
  },
};

export default function Navbar({ copy }: { copy: LandingPageCopy }) {
  const [navUser, setNavUser] = useState<{
    nickname?: string; email?: string; display_id?: string;
  } | null>(null);
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  useEffect(() => {
    const cookies = document.cookie.split('; ');
    const loginCookie = cookies.find(row => row.startsWith('is_logged_in=1'));

    if (!loginCookie) return;

    const fetchUser = async () => {
      try {
        const refreshRes = await fetch('/api/v1/auth/refresh', {
          method: 'POST',
          credentials: 'include',
        });

        if (!refreshRes.ok) {
          document.cookie = 'is_logged_in=; Max-Age=0; path=/;';
          return;
        }

        const res = await fetch('/api/v1/auth/me', {
          credentials: 'include',
          cache: 'no-store',
        });

        if (res.ok) {
          const data = await res.json();
          setNavUser(data);
        }
      } catch (err) {
        // 조용히 처리
      }
    };

    fetchUser();
  }, []);

  const handleLogout = async () => {
    await fetch('/api/v1/auth/logout', { method: 'POST', credentials: 'include' });
    setNavUser(null);
    setIsMenuOpen(false);
    window.location.href = '/';
  };

  const menuAriaLabel = copy.locale === 'en' ? 'Open main menu' : '메인 메뉴 열기';
  const menuTitle = copy.locale === 'en' ? 'Main Menu' : '메인 메뉴';
  const loginLabel = 'Login';

  return (
    <>
      <AppHeaderShell
        logoHref={copy.homeHref}
        logoIconSize={36}
        logoTextSize="20px"
        innerStyle={{ position: 'relative', padding: '14px 18px' }}
        centerStyle={{
          position: 'absolute',
          left: '50%',
          transform: 'translateX(-50%)',
          flex: '0 0 auto',
          pointerEvents: 'auto',
        }}
        centerSlot={
          <nav className="hidden items-center gap-4 md:flex lg:gap-7" aria-label={menuTitle}>
            {copy.navLinks.map((link) => (
              <Link
                key={link.name}
                href={link.href}
                className="heroNavLink text-sm font-bold text-[#21485D] transition hover:text-[#0F3145]"
              >
                {link.name}
              </Link>
            ))}
          </nav>
        }
        rightSlot={
          <div className="flex items-center gap-2">
            <button
              type="button"
              className="grid h-10 w-10 place-items-center rounded-full bg-white/80 text-[#12364A] shadow-[0_10px_24px_rgba(42,94,139,0.12)] transition active:scale-95 md:hidden"
              aria-label={menuAriaLabel}
              aria-expanded={isMenuOpen}
              onClick={() => setIsMenuOpen((current) => !current)}
            >
              <span className="grid gap-1.5" aria-hidden="true">
                <span className="block h-0.5 w-5 rounded-full bg-current" />
                <span className="block h-0.5 w-5 rounded-full bg-current" />
                <span className="block h-0.5 w-5 rounded-full bg-current" />
              </span>
            </button>

            {navUser ? (
              <div className="hidden md:flex">
                <LearnerHeaderActions onLogout={handleLogout} copy={learnerHeaderActionsCopy[copy.locale]} />
              </div>
            ) : (
              <Link
                href={copy.loginHref}
                className="hidden rounded-lg border border-[#2A789A]/20 bg-white/80 px-6 py-2.5 text-sm text-[#12364A] shadow-[0_10px_24px_rgba(42,94,139,0.14)] transition hover:-translate-y-0.5 hover:bg-white active:scale-95 min-[420px]:inline-flex md:inline-flex"
              >
                {loginLabel}
              </Link>
            )}
          </div>
        }
      />

      {isMenuOpen ? (
        <div className="fixed left-3 right-3 top-[74px] z-[70] overflow-hidden rounded-[26px] bg-[#F7FCFB]/95 p-3 shadow-[0_26px_70px_rgba(21,55,74,0.22)] backdrop-blur-xl md:hidden">
          <div className="px-2 pb-2 pt-1 text-xs font-black uppercase tracking-[0.16em] text-[#467083]">
            {menuTitle}
          </div>
          <nav className="grid gap-1" aria-label={menuTitle}>
            {copy.navLinks.map((link) => (
              <Link
                key={link.name}
                href={link.href}
                className="rounded-[18px] px-4 py-3 text-base font-black text-[#12364A] transition hover:bg-white/80 active:scale-[0.99]"
                onClick={() => setIsMenuOpen(false)}
              >
                {link.name}
              </Link>
            ))}
            {navUser ? (
              <button
                type="button"
                className="rounded-[18px] px-4 py-3 text-left text-base font-black text-[#8A4B33] transition hover:bg-white/80 active:scale-[0.99]"
                onClick={() => void handleLogout()}
              >
                Logout
              </button>
            ) : (
              <Link
                href={copy.loginHref}
                className="rounded-[18px] bg-[#12364A] px-4 py-3 text-base font-black text-white transition active:scale-[0.99]"
                onClick={() => setIsMenuOpen(false)}
              >
                {loginLabel}
              </Link>
            )}
          </nav>
        </div>
      ) : null}
    </>
  );
}
