'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import type { Locale } from '@/lib/i18n/locales';
import { normalizeLocale } from '@/lib/i18n/locales';
import { getLandingPageCopy } from '@/lib/i18n/pages/landing';
import BrandLogo from '@/components/common/BrandLogo';
import MobileMainNavDrawer from '@/components/navigation/MobileMainNavDrawer';
import PublicCosmosMobileFirstView from '@/components/landing/PublicCosmosMobileFirstView';

type MeResponse = {
  ui_locale?: string | null;
};

export default function LandingPageClient({
  locale,
  preferUserLocale = false,
}: {
  locale: Locale;
  preferUserLocale?: boolean;
}) {
  const [effectiveLocale, setEffectiveLocale] = useState<Locale>(locale);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [navUser, setNavUser] = useState<{ nickname?: string; email?: string; display_id?: string } | null>(null);

  useEffect(() => {
    setEffectiveLocale(locale);
  }, [locale]);

  useEffect(() => {
    if (typeof window === 'undefined') return;
    window.history.scrollRestoration = 'manual';
    window.scrollTo({ top: 0, left: 0, behavior: 'auto' });
  }, []);

  useEffect(() => {
    if (!document.cookie.includes('is_logged_in=1')) return;

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
      } catch {
        // Keep the public home header usable even if auth probing fails.
      }
    };

    fetchUser();
  }, []);

  useEffect(() => {
    if (!preferUserLocale) return;
    if (!document.cookie.includes('is_logged_in=1')) return;

    let cancelled = false;

    const loadUserLocale = async () => {
      const res = await fetch('/api/v1/auth/me', {
        credentials: 'include',
        cache: 'no-store',
      });
      if (!res.ok) return;

      const data = (await res.json()) as MeResponse;
      if (!cancelled) {
        setEffectiveLocale(normalizeLocale(data.ui_locale));
      }
    };

    loadUserLocale().catch(() => {});

    return () => {
      cancelled = true;
    };
  }, [preferUserLocale]);

  const copy = getLandingPageCopy(effectiveLocale);

  const handleLogout = async () => {
    await fetch('/api/v1/auth/logout', { method: 'POST', credentials: 'include' });
    setNavUser(null);
    window.location.href = '/';
  };

  const drawerTitle = copy.locale === 'en' ? 'Main Menu' : '메인 메뉴';
  const menuLabel = copy.locale === 'en' ? 'Open main menu' : '메인 메뉴 열기';
  const accountLabel = navUser?.nickname || navUser?.display_id || navUser?.email || (copy.locale === 'en' ? 'Learner' : '학습자');

  const accountSlot = navUser ? (
    <div className="rounded-[22px] border border-[#D6E6ED]/80 bg-white/78 p-3 shadow-[0_14px_34px_rgba(33,73,93,0.10)]">
      <p className="m-0 text-[11px] font-black uppercase tracking-[0.14em] text-[#6C7ABF]">Lumi Account</p>
      <p className="mt-1 truncate text-sm font-black text-[#12384E]">{accountLabel}</p>
      <div className="mt-3 grid grid-cols-2 gap-2">
        <Link
          href="/dashboard/settings/profile"
          onClick={() => setDrawerOpen(false)}
          className="grid h-9 place-items-center rounded-full bg-[#EDF7F7] text-xs font-black text-[#12384E]"
        >
          {copy.locale === 'en' ? 'Profile' : '사용자 정보'}
        </Link>
        <button
          type="button"
          onClick={handleLogout}
          className="h-9 rounded-full bg-white text-xs font-black text-[#416679] shadow-[inset_0_0_0_1px_rgba(151,190,210,0.44)]"
        >
          Logout
        </button>
      </div>
    </div>
  ) : (
    <div className="rounded-[22px] border border-[#D6E6ED]/80 bg-white/78 p-3 shadow-[0_14px_34px_rgba(33,73,93,0.10)]">
      <p className="m-0 text-[11px] font-black uppercase tracking-[0.14em] text-[#6C7ABF]">Lumi Account</p>
      <p className="mt-1 text-sm font-black text-[#12384E]">
        {copy.locale === 'en' ? 'Login to continue your learning map.' : '로그인하면 학습 지도를 이어갈 수 있어요.'}
      </p>
      <Link
        href={copy.loginHref}
        onClick={() => setDrawerOpen(false)}
        className="mt-3 grid h-10 place-items-center rounded-full bg-[#12384E] text-sm font-black text-white"
      >
        Login
      </Link>
    </div>
  );

  return (
    <div className="relative min-h-screen overflow-hidden bg-[#F5FAFA] text-[#12384E]">
      <header className="fixed left-0 right-0 top-0 z-[60] border-b border-white/60 bg-white/70 backdrop-blur-xl">
        <div className="mx-auto flex h-[60px] w-full max-w-[720px] items-center justify-between gap-3 px-3">
          <button
            type="button"
            className="grid h-10 w-10 flex-shrink-0 place-items-center rounded-full border border-white/70 bg-white/72 text-[#12384E] shadow-[0_12px_28px_rgba(33,73,93,0.14)] transition hover:bg-white"
            aria-label={menuLabel}
            onClick={() => setDrawerOpen(true)}
          >
            <span className="grid gap-1.5" aria-hidden="true">
              <span className="block h-0.5 w-5 rounded-full bg-current" />
              <span className="block h-0.5 w-5 rounded-full bg-current" />
              <span className="block h-0.5 w-5 rounded-full bg-current" />
            </span>
          </button>
          <div className="min-w-0 flex-1">
            <BrandLogo href={copy.homeHref} iconSize={30} textSize="18px" />
          </div>
          <div className="h-10 w-10 flex-shrink-0" aria-hidden="true" />
        </div>
      </header>

      <MobileMainNavDrawer
        open={drawerOpen}
        locale={copy.locale}
        title={drawerTitle}
        onClose={() => setDrawerOpen(false)}
        accountSlot={accountSlot}
        immersive
      />

      <PublicCosmosMobileFirstView copy={copy.hero} locale={copy.locale} loginHref={copy.loginHref} />
    </div>
  );
}
