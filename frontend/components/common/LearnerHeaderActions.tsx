'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useEffect, useMemo, useRef, useState, type CSSProperties, type ReactNode } from 'react';
import { appHeaderActionButtonStyle, appHeaderActionLinkStyle } from '@/components/common/AppHeaderShell';
import LumiModalShell, {
  lumiModalPrimaryButtonStyle,
  lumiModalSecondaryButtonStyle,
} from '@/components/common/LumiModalShell';
import LearnerUserMenu, { defaultLearnerUserMenuCopy, type LearnerUserMenuCopy } from '@/components/common/LearnerUserMenu';
import { normalizeLocale, type Locale } from '@/lib/i18n/locales';

interface CompactLearnerUser {
  email?: string;
  display_id?: string;
  nickname?: string;
  total_points?: number;
}

interface LearnerHeaderActionsProps {
  onLogout?: () => void | Promise<void>;
  galaxyHref?: string;
  galaxyLabel?: string;
  galaxyTitle?: string;
  copy?: LearnerHeaderActionsCopy;
  extraAction?: ReactNode;
  compactExtraAction?: ReactNode;
  palette?: {
    text: string;
    galaxyText: string;
    galaxyBackground: string;
    galaxyBorder: string;
    galaxyShadow?: string;
    userMenu?: {
      triggerText: string;
      pointText: string;
      pointBackground: string;
      pointBorder: string;
    };
  };
}

export interface LearnerHeaderActionsCopy {
  userMenu: LearnerUserMenuCopy;
  compactMenuAriaLabel: string;
  compactPointsLoading: string;
  compactPointsOwned: (points: number) => string;
  logoutModalMessage: string;
  logoutCancel: string;
  logout: string;
  loggingOut: string;
}

const learnerHeaderActionsCopyByLocale: Record<Locale, LearnerHeaderActionsCopy> = {
  ko: {
    userMenu: defaultLearnerUserMenuCopy,
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

export const defaultLearnerHeaderActionsCopy: LearnerHeaderActionsCopy = learnerHeaderActionsCopyByLocale.ko;

export function getLearnerHeaderActionsCopy(locale?: string | null): LearnerHeaderActionsCopy {
  return learnerHeaderActionsCopyByLocale[normalizeLocale(locale)];
}

const COMPACT_BREAKPOINT = 720;

export default function LearnerHeaderActions({
  onLogout,
  galaxyHref = '/dashboard?mode=galaxy',
  galaxyLabel = 'Galaxy',
  galaxyTitle = '갤럭시 대시보드로 이동',
  copy = defaultLearnerHeaderActionsCopy,
  extraAction,
  compactExtraAction,
  palette,
}: LearnerHeaderActionsProps) {
  const pathname = usePathname();
  const [isCompact, setIsCompact] = useState(false);
  const [isOpen, setIsOpen] = useState(false);
  const [showLogoutConfirm, setShowLogoutConfirm] = useState(false);
  const [isLoggingOut, setIsLoggingOut] = useState(false);
  const [compactPoints, setCompactPoints] = useState<number | null>(null);
  const [compactName, setCompactName] = useState(copy.userMenu.fallbackName);
  const containerRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    setCompactName((current) => (current === defaultLearnerHeaderActionsCopy.userMenu.fallbackName ? copy.userMenu.fallbackName : current));
  }, [copy.userMenu.fallbackName]);

  useEffect(() => {
    const syncCompactMode = () => {
      const nextCompact = window.innerWidth <= COMPACT_BREAKPOINT;
      setIsCompact(nextCompact);
      if (!nextCompact) {
        setIsOpen(false);
      }
    };

    const handlePointerDown = (event: MouseEvent) => {
      if (!containerRef.current?.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setIsOpen(false);
      }
    };

    syncCompactMode();
    window.addEventListener('resize', syncCompactMode);
    document.addEventListener('mousedown', handlePointerDown);
    window.addEventListener('keydown', handleEscape);

    return () => {
      window.removeEventListener('resize', syncCompactMode);
      document.removeEventListener('mousedown', handlePointerDown);
      window.removeEventListener('keydown', handleEscape);
    };
  }, []);

  useEffect(() => {
    if (!isCompact) return;
    let cancelled = false;
    const loadCompactPoints = async () => {
      try {
        const refreshRes = await fetch('/api/v1/auth/refresh', { method: 'POST', credentials: 'include' });
        if (!refreshRes.ok) return;
        const meRes = await fetch('/api/v1/auth/me', { credentials: 'include', cache: 'no-store' });
        if (!meRes.ok) return;
        const data = (await meRes.json()) as CompactLearnerUser;
        if (!cancelled) {
          const value = Number.isFinite(data.total_points) ? Math.max(0, Math.floor(data.total_points ?? 0)) : 0;
          setCompactPoints(value);
          setCompactName(data.nickname || data.display_id || data.email || copy.userMenu.fallbackName);
        }
      } catch {
        if (!cancelled) setCompactPoints(null);
      }
    };
    void loadCompactPoints();
    const handleUserRefresh = () => void loadCompactPoints();
    window.addEventListener('learnweaver:user-refresh', handleUserRefresh as EventListener);
    return () => {
      cancelled = true;
      window.removeEventListener('learnweaver:user-refresh', handleUserRefresh as EventListener);
    };
  }, [copy.userMenu.fallbackName, isCompact]);

  const isGalaxyActive = useMemo(
    () => pathname === '/dashboard' || pathname?.startsWith('/dashboard/course-drafts/') === true,
    [pathname],
  );

  const requestLogout = () => {
    if (!onLogout) return;
    setShowLogoutConfirm(true);
  };

  const confirmLogout = async () => {
    if (!onLogout || isLoggingOut) return;
    setIsLoggingOut(true);
    try {
      await onLogout();
    } finally {
      setIsLoggingOut(false);
      setShowLogoutConfirm(false);
    }
  };

  const logoutConfirmModal = showLogoutConfirm && onLogout ? (
    <LumiModalShell
      ariaLabel="Logout confirmation"
      eyebrow="Lumi Confirm"
      lumiState="curious"
      tone="warm"
      title="Log out?"
      message={copy.logoutModalMessage}
      onClose={() => {
        if (!isLoggingOut) {
          setShowLogoutConfirm(false);
        }
      }}
      actions={
        <>
          <button
            type="button"
            style={lumiModalSecondaryButtonStyle}
            disabled={isLoggingOut}
            onClick={() => setShowLogoutConfirm(false)}
          >
            {copy.logoutCancel}
          </button>
          <button
            type="button"
            style={{
              ...lumiModalPrimaryButtonStyle,
              opacity: isLoggingOut ? 0.72 : 1,
              cursor: isLoggingOut ? 'wait' : 'pointer',
            }}
            disabled={isLoggingOut}
            onClick={() => void confirmLogout()}
          >
            {isLoggingOut ? copy.loggingOut : copy.logout}
          </button>
        </>
      }
    />
  ) : null;

  if (!isCompact) {
    return (
      <>
        <LearnerUserMenu palette={palette?.userMenu} copy={copy.userMenu} onLogoutRequest={onLogout ? requestLogout : undefined} />
        <Link href={galaxyHref} style={galaxyLinkStyle(isGalaxyActive, palette)} title={galaxyTitle}>
          {galaxyLabel}
        </Link>
        {extraAction}
        {logoutConfirmModal}
      </>
    );
  }

  return (
    <div ref={containerRef} style={compactWrapperStyle}>
      <Link href="/dashboard/settings/points" style={compactPointChipStyle}>
        {compactPoints === null ? 'pt' : `${compactPoints}pt`}
      </Link>
      <button
        type="button"
        aria-label={copy.compactMenuAriaLabel}
        aria-expanded={isOpen}
        onClick={() => setIsOpen((current) => !current)}
        style={compactTriggerStyle(isOpen, palette)}
      >
        <span style={hamburgerLineStyle(palette)} />
        <span style={hamburgerLineStyle(palette)} />
        <span style={hamburgerLineStyle(palette)} />
      </button>

      {isOpen ? (
        <div style={compactMenuStyle}>
          <div style={compactMenuHeaderStyle}>
            <strong style={compactMenuTitleStyle}>{compactName}</strong>
            <span style={compactMenuCaptionStyle}>
              {compactPoints === null ? copy.compactPointsLoading : copy.compactPointsOwned(compactPoints)}
            </span>
          </div>

          <Link href={galaxyHref} style={compactMenuItemStyle(isGalaxyActive)} onClick={() => setIsOpen(false)}>
            {galaxyLabel}
          </Link>
          {compactExtraAction ? (
            <div style={compactMenuActionWrapperStyle} onClick={() => setIsOpen(false)}>
              {compactExtraAction}
            </div>
          ) : null}
          <Link
            href="/dashboard/settings/profile"
            style={compactMenuItemStyle(pathname === '/dashboard/settings/profile')}
            onClick={() => setIsOpen(false)}
          >
            {copy.userMenu.profile}
          </Link>
          <Link
            href="/dashboard/settings/points"
            style={compactMenuItemStyle(pathname === '/dashboard/settings/points')}
            onClick={() => setIsOpen(false)}
          >
            {copy.userMenu.points}
          </Link>
          <Link
            href="/dashboard/settings/ai"
            style={compactMenuItemStyle(pathname === '/dashboard/settings/ai')}
            onClick={() => setIsOpen(false)}
          >
            {copy.userMenu.aiSettings}
          </Link>
          {onLogout ? (
            <button
              type="button"
              style={compactLogoutButtonStyle}
              onClick={() => {
                setIsOpen(false);
                requestLogout();
              }}
            >
              {copy.logout}
            </button>
          ) : null}
        </div>
      ) : null}
      {logoutConfirmModal}
    </div>
  );
}

const compactWrapperStyle: CSSProperties = {
  position: 'relative',
  display: 'inline-flex',
  alignItems: 'center',
  gap: '8px',
};

const compactPointChipStyle: CSSProperties = {
  minHeight: '34px',
  minWidth: '50px',
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  padding: '0 9px',
  borderRadius: '11px',
  background: 'rgba(255, 247, 214, 0.96)',
  border: '1px solid rgba(245,158,11,0.58)',
  color: '#78350F',
  fontSize: '12px',
  fontWeight: 900,
  textDecoration: 'none',
  boxShadow: '0 8px 18px rgba(42,94,139,0.1)',
};

const compactTriggerStyle = (
  isOpen: boolean,
  palette?: LearnerHeaderActionsProps['palette'],
): CSSProperties => ({
  ...appHeaderActionButtonStyle(isOpen),
  width: '42px',
  minWidth: '42px',
  padding: '0',
  justifyContent: 'center',
  gap: '3px',
  flexDirection: 'column',
  color: palette?.text ?? '#12364A',
  background: isOpen
    ? palette?.galaxyBackground ?? 'rgba(255,255,255,0.86)'
    : 'rgba(255,255,255,0.62)',
  borderColor: palette?.galaxyBorder ?? 'rgba(38,99,133,0.2)',
});

const hamburgerLineStyle = (palette?: LearnerHeaderActionsProps['palette']): CSSProperties => ({
  width: '16px',
  height: '2px',
  borderRadius: '999px',
  background: palette?.text ?? '#12364A',
});

const compactMenuStyle: CSSProperties = {
  position: 'absolute',
  top: 'calc(100% + 10px)',
  right: 0,
  minWidth: '200px',
  padding: '10px',
  borderRadius: '16px',
  border: '1px solid rgba(180, 205, 255, 0.16)',
  background: 'rgba(9, 19, 35, 0.94)',
  backdropFilter: 'blur(16px)',
  WebkitBackdropFilter: 'blur(16px)',
  boxShadow: '0 20px 45px rgba(0, 0, 0, 0.28)',
  display: 'grid',
  gap: '6px',
  zIndex: 90,
};

const compactMenuHeaderStyle: CSSProperties = {
  display: 'grid',
  gap: '2px',
  padding: '4px 6px 10px',
  borderBottom: '1px solid rgba(180, 205, 255, 0.12)',
  marginBottom: '2px',
};

const compactMenuTitleStyle: CSSProperties = {
  color: '#F4F7FF',
  fontSize: '13px',
};

const compactMenuCaptionStyle: CSSProperties = {
  color: 'rgba(232,234,242,0.64)',
  fontSize: '12px',
};

const compactMenuItemStyle = (active: boolean): CSSProperties => ({
  display: 'flex',
  alignItems: 'center',
  minHeight: '40px',
  padding: '0 12px',
  borderRadius: '12px',
  color: '#E8EEFF',
  textDecoration: 'none',
  fontSize: '14px',
  fontWeight: active ? 700 : 500,
  background: active ? 'rgba(255,255,255,0.10)' : 'transparent',
  border: active ? '1px solid rgba(255,255,255,0.14)' : '1px solid transparent',
});

const compactMenuActionWrapperStyle: CSSProperties = {
  display: 'grid',
};

const galaxyLinkStyle = (
  active: boolean,
  palette?: LearnerHeaderActionsProps['palette'],
): CSSProperties => ({
  ...appHeaderActionLinkStyle(active),
  minHeight: '36px',
  padding: '7px 18px',
  fontWeight: 900,
  background: palette?.galaxyBackground ?? (active ? '#37266F' : '#4F36A6'),
  border: `1px solid ${palette?.galaxyBorder ?? (active ? 'rgba(226, 216, 255, 0.86)' : 'rgba(226, 216, 255, 0.62)')}`,
  color: palette?.galaxyText ?? '#FFFFFF',
  boxShadow: palette?.galaxyShadow ?? (active ? '0 12px 30px rgba(64, 42, 128, 0.34)' : '0 9px 24px rgba(64, 42, 128, 0.26)'),
});

const compactLogoutButtonStyle: CSSProperties = {
  ...compactMenuItemStyle(false),
  fontFamily: 'inherit',
  cursor: 'pointer',
  width: '100%',
  justifyContent: 'flex-start',
};
