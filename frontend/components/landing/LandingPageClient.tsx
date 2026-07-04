'use client';

import { useEffect, useRef, useState } from 'react';
import Link from 'next/link';
import type { Locale } from '@/lib/i18n/locales';
import { normalizeLocale } from '@/lib/i18n/locales';
import { getLandingPageCopy } from '@/lib/i18n/pages/landing';
import LearnerHeaderActions, { type LearnerHeaderActionsCopy } from '@/components/common/LearnerHeaderActions';
import LearnerAppShell from '@/components/navigation/LearnerAppShell';
import Hero from '@/components/landing/Hero';
import ScrollCursorGuide from '@/components/landing/ScrollCursorGuide';
import BetaNotice from '@/components/landing/BetaNotice';
import LandingJourneySections from '@/components/landing/LandingJourneySections';
import Footer from '@/components/landing/Footer';

type MeResponse = {
  ui_locale?: string | null;
};

const ALPHA_NOTICE_SUPPRESS_KEY = 'learnweaver:landing-alpha-notice-suppress-date';

const COURSE_CREATE_HASH = '#course-create-ready';

const JOURNEY_READY_PROGRESS_BY_HASH: Record<string, number> = {
  '#intro-video-ready': 0.18,
  '#public-documents-ready': 0.46,
  '#ai-points-byok-ready': 0.74,
  '#platform-news-ready': 0.97,
};

type MobileJourneyProgressConfig = {
  startRatio: number;
  distanceViewportRatio: number;
  distanceSectionRatio: number;
  distanceMaxViewportRatio: number;
};

const MOBILE_JOURNEY_PROGRESS_CONFIG_BY_SECTION_ID: Record<string, MobileJourneyProgressConfig> = {
  'intro-video': {
    startRatio: 1.08,
    distanceViewportRatio: 0.9,
    distanceSectionRatio: 0.34,
    distanceMaxViewportRatio: 1.18,
  },
  'public-documents': {
    startRatio: 1.06,
    distanceViewportRatio: 0.86,
    distanceSectionRatio: 0.32,
    distanceMaxViewportRatio: 1.14,
  },
  'ai-points-byok': {
    startRatio: 1.06,
    distanceViewportRatio: 0.88,
    distanceSectionRatio: 0.33,
    distanceMaxViewportRatio: 1.16,
  },
  'platform-news': {
    startRatio: 1.06,
    distanceViewportRatio: 0.86,
    distanceSectionRatio: 0.32,
    distanceMaxViewportRatio: 1.14,
  },
};

function isJourneyReadyHash(hash: string) {
  return JOURNEY_READY_PROGRESS_BY_HASH[hash] != null;
}

function getFixedTopOffset() {
  if (typeof document === 'undefined') return 0;

  const topHeader = document.querySelector<HTMLElement>('body > div header');
  const sectionNav = document.querySelector<HTMLElement>('nav[aria-label="Section navigation"]');
  const headerHeight = topHeader?.getBoundingClientRect().height ?? 60;
  const navHeight = sectionNav?.getBoundingClientRect().height ?? 48;

  return Math.round(headerHeight + navHeight + 8);
}

function getJourneyReadyAnchor(hash: string) {
  if (typeof document === 'undefined') return null;

  const normalizedHash = hash.startsWith('#') ? hash.slice(1) : hash;
  return document.getElementById(normalizedHash) as HTMLElement | null;
}

function getJourneyReadyPanel(hash: string) {
  const anchor = getJourneyReadyAnchor(hash);
  return anchor?.closest<HTMLElement>('.landingJourneySection') ?? null;
}

function getDesktopReadyProgress(hash: string) {
  const panelProgress = Number(getJourneyReadyPanel(hash)?.dataset.readyProgress);
  if (Number.isFinite(panelProgress)) return panelProgress;
  return JOURNEY_READY_PROGRESS_BY_HASH[hash];
}

function getMobileReadyTargetTop(panel: HTMLElement) {
  const config = MOBILE_JOURNEY_PROGRESS_CONFIG_BY_SECTION_ID[panel.id];
  if (!config) return null;

  const rect = panel.getBoundingClientRect();
  const viewportHeight = window.innerHeight || 1;
  const start = viewportHeight * config.startRatio;
  const distance = Math.min(
    Math.max(viewportHeight * config.distanceViewportRatio, rect.height * config.distanceSectionRatio),
    viewportHeight * config.distanceMaxViewportRatio,
  );
  const readyProgress = Number(panel.dataset.mobileReadyProgress);
  const progress = Number.isFinite(readyProgress) ? Math.min(1, Math.max(0, readyProgress)) : 1;
  const readyRectTop = start - distance * progress;

  return window.scrollY + rect.top - readyRectTop;
}

function scrollToCourseCreateTop(behavior: ScrollBehavior = 'smooth') {
  if (typeof window === 'undefined') return false;
  window.scrollTo({ top: 0, left: 0, behavior });
  return true;
}

function scrollToJourneyReadyHash(hash: string, behavior: ScrollBehavior = 'smooth') {
  if (typeof window === 'undefined') return false;

  const normalizedHash = hash.startsWith('#') ? hash : `#${hash}`;
  if (!isJourneyReadyHash(normalizedHash)) return false;

  const anchor = getJourneyReadyAnchor(normalizedHash);
  const panel = getJourneyReadyPanel(normalizedHash);

  if (window.innerWidth < 1024) {
    const progressTargetTop = panel ? getMobileReadyTargetTop(panel) : null;
    const anchorTargetTop = anchor
      ? window.scrollY + anchor.getBoundingClientRect().top - getFixedTopOffset()
      : null;
    const targetTop = progressTargetTop ?? anchorTargetTop;
    if (targetTop == null) return false;

    window.scrollTo({ top: Math.max(0, targetTop), left: 0, behavior });
    return true;
  }

  const progress = getDesktopReadyProgress(normalizedHash);
  if (progress == null) return false;

  const runway = document.querySelector<HTMLElement>('.landingStickyJourneyRunway');
  if (!runway) return false;

  const rect = runway.getBoundingClientRect();
  const runwayTop = window.scrollY + rect.top;
  const travel = Math.max(1, runway.offsetHeight - window.innerHeight);
  const targetTop = runwayTop + travel * progress;

  window.scrollTo({ top: Math.max(0, targetTop), left: 0, behavior });
  return true;
}

function getTodayKey() {
  return new Date().toISOString().slice(0, 10);
}

function shouldSuppressToday(key: string) {
  if (typeof window === 'undefined') return false;
  try {
    return window.localStorage.getItem(key) === getTodayKey();
  } catch {
    return false;
  }
}

function suppressToday(key: string) {
  if (typeof window === 'undefined') return;
  try {
    window.localStorage.setItem(key, getTodayKey());
  } catch {
    // Ignore storage failures.
  }
}

const learnerHeaderActionsCopy: Record<Locale, LearnerHeaderActionsCopy> = {
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

export default function LandingPageClient({
  locale,
  preferUserLocale = false,
}: {
  locale: Locale;
  preferUserLocale?: boolean;
}) {
  const observer = useRef<IntersectionObserver | null>(null);
  const [effectiveLocale, setEffectiveLocale] = useState<Locale>(locale);
  const [isAlphaNoticeOpen, setIsAlphaNoticeOpen] = useState(true);
  const [navUser, setNavUser] = useState<{ nickname?: string; email?: string; display_id?: string } | null>(null);

  useEffect(() => {
    if (typeof window === 'undefined') return;

    window.history.scrollRestoration = 'manual';
    if (window.location.hash) return;

    const scrollToTop = () => window.scrollTo({ top: 0, left: 0, behavior: 'auto' });
    scrollToTop();
    const frameOne = window.requestAnimationFrame(scrollToTop);
    const frameTwo = window.requestAnimationFrame(() => window.requestAnimationFrame(scrollToTop));
    const timeout = window.setTimeout(scrollToTop, 180);

    return () => {
      window.cancelAnimationFrame(frameOne);
      window.cancelAnimationFrame(frameTwo);
      window.clearTimeout(timeout);
    };
  }, []);

  useEffect(() => {
    if (typeof window === 'undefined') return;

    const hash = window.location.hash;
    if (hash === COURSE_CREATE_HASH) {
      const frame = window.requestAnimationFrame(() => scrollToCourseCreateTop('auto'));
      const timeout = window.setTimeout(() => scrollToCourseCreateTop('auto'), 220);

      return () => {
        window.cancelAnimationFrame(frame);
        window.clearTimeout(timeout);
      };
    }
    if (!isJourneyReadyHash(hash)) return;

    setIsAlphaNoticeOpen(false);
    const frame = window.requestAnimationFrame(() => scrollToJourneyReadyHash(hash, 'auto'));
    const timeout = window.setTimeout(() => scrollToJourneyReadyHash(hash, 'auto'), 220);

    return () => {
      window.cancelAnimationFrame(frame);
      window.clearTimeout(timeout);
    };
  }, []);

  useEffect(() => {
    if (typeof window === 'undefined') return;

    const handleHashClick = (event: MouseEvent) => {
      const target = event.target;
      if (!(target instanceof Element)) return;

      const anchor = target.closest<HTMLAnchorElement>('a[href]');
      if (!anchor) return;

      const url = new URL(anchor.href, window.location.href);
      if (url.origin !== window.location.origin || url.pathname !== window.location.pathname) return;
      if (url.hash === COURSE_CREATE_HASH) {
        event.preventDefault();
        window.history.pushState(null, '', url.hash);
        window.requestAnimationFrame(() => scrollToCourseCreateTop());
        return;
      }
      if (!isJourneyReadyHash(url.hash)) return;

      event.preventDefault();
      setIsAlphaNoticeOpen(false);
      window.history.pushState(null, '', url.hash);
      window.requestAnimationFrame(() => scrollToJourneyReadyHash(url.hash));
    };

    document.addEventListener('click', handleHashClick);
    return () => document.removeEventListener('click', handleHashClick);
  }, []);

  const copy = getLandingPageCopy(effectiveLocale);

  const dismissAlphaNotice = (suppressForToday = false) => {
    if (suppressForToday) suppressToday(ALPHA_NOTICE_SUPPRESS_KEY);
    setIsAlphaNoticeOpen(false);
  };

  const handleIntroVideoClick = () => {
    setIsAlphaNoticeOpen(false);
    window.history.pushState(null, '', '#intro-video-ready');
    window.requestAnimationFrame(() => {
      if (!scrollToJourneyReadyHash('#intro-video-ready')) {
        document.getElementById('intro-video-ready')?.scrollIntoView({ behavior: 'smooth', block: 'start' });
      }
    });
  };

  useEffect(() => {
    if (shouldSuppressToday(ALPHA_NOTICE_SUPPRESS_KEY)) setIsAlphaNoticeOpen(false);
  }, []);

  useEffect(() => {
    observer.current = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add('show');
        }
      });
    }, { threshold: 0.15 });

    const hiddenElements = document.querySelectorAll('.animate-on-scroll');
    hiddenElements.forEach((el) => observer.current?.observe(el));

    return () => observer.current?.disconnect();
  }, []);

  useEffect(() => {
    setEffectiveLocale(locale);
  }, [locale]);

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
        // Keep the public landing header usable even if auth probing fails.
      }
    };

    fetchUser();
  }, []);

  const handleLogout = async () => {
    await fetch('/api/v1/auth/logout', { method: 'POST', credentials: 'include' });
    setNavUser(null);
    window.location.href = '/';
  };

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

  const headerActions = (
    <div className="flex items-center gap-2">
      {navUser ? (
        <LearnerHeaderActions onLogout={handleLogout} copy={learnerHeaderActionsCopy[copy.locale]} />
      ) : (
        <Link
          href={copy.loginHref}
          onClick={(event) => {
            if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey || event.button !== 0) return;
            event.preventDefault();
            window.location.href = copy.loginHref;
          }}
          className="inline-flex h-10 items-center rounded-full bg-white/80 px-4 text-sm font-black text-[#12364A] shadow-[0_10px_24px_rgba(42,94,139,0.12)] transition hover:bg-white"
        >
          Login
        </Link>
      )}
    </div>
  );

  return (
    <div className="landingAppShell">
    <LearnerAppShell
      locale={copy.locale}
      logoHref={copy.homeHref}
      rightSlot={headerActions}
      contentClassName="landingAppContent !px-0 !pb-0"
    >
      <div className="landingPageRoot">
        <ScrollCursorGuide disabled={false} />
      {isAlphaNoticeOpen ? (
        <BetaNotice
          copy={copy.alphaNotice}
          onVideoClick={handleIntroVideoClick}
          onDismiss={dismissAlphaNotice}
        />
      ) : null}
      <section className="landingJourneyLayer" aria-label="LearnCosmos journey intro">
        <svg className="landingJourneyVectorSvg" viewBox="0 0 1120 6600" aria-hidden="true">
          <defs>
            <linearGradient id="landingJourneyBluePath" x1="0" y1="0" x2="1" y2="1">
              <stop offset="0" stopColor="#6AD2C1" stopOpacity="0" />
              <stop offset="0.1" stopColor="#88DAFF" stopOpacity="0.48" />
              <stop offset="0.26" stopColor="#6AD2C1" stopOpacity="0.62" />
              <stop offset="0.52" stopColor="#7BC9FF" stopOpacity="0.54" />
              <stop offset="0.78" stopColor="#6AD2C1" stopOpacity="0.38" />
              <stop offset="1" stopColor="#88DAFF" stopOpacity="0" />
            </linearGradient>
            <filter id="landingJourneyVectorGlow" x="-20%" y="-10%" width="140%" height="120%">
              <feGaussianBlur stdDeviation="10" result="blur" />
              <feMerge>
                <feMergeNode in="blur" />
                <feMergeNode in="SourceGraphic" />
              </feMerge>
            </filter>
          </defs>
          <path className="landingJourneyVectorGlow" d="M506 1068 C592 1262 616 1374 430 1490 C214 1742 356 2096 648 2378 C928 2648 930 2970 666 3280 C420 3570 264 3930 472 4250 C682 4574 874 4864 604 5450 C430 5712 364 6014 552 6400" />
          <path className="landingJourneyVectorBase" d="M514 820 C476 850 466 928 506 1018 C592 1212 616 1324 430 1440 C214 1692 356 2046 648 2328 C928 2598 930 2920 666 3230 C420 3520 264 3880 472 4200 C682 4524 874 4814 604 5400 C430 5662 364 5964 552 6350" />
          <path className="landingJourneyVectorCore" d="M514 820 C476 850 466 928 506 1018 C592 1212 616 1324 430 1440 C214 1692 356 2046 648 2328 C928 2598 930 2920 666 3230 C420 3520 264 3880 472 4200 C682 4524 874 4814 604 5400 C430 5662 364 5964 552 6350" />
        </svg>
        <span id="course-create-ready" className="landingCourseCreateReadyAnchor" aria-hidden="true" />
        <Hero copy={copy.hero} locale={copy.locale} />
        <LandingJourneySections
          copy={{ introVideo: copy.introVideo, policies: copy.policies, aiUsage: copy.aiUsage }}
          locale={copy.locale}
          videoGateDisabled={isAlphaNoticeOpen}
        />
      </section>
      <Footer copy={copy.footer} homeHref={copy.homeHref} />
      </div>
    </LearnerAppShell>
    </div>
  );
}
