'use client';

import Link from 'next/link';
import { useEffect, useMemo, useRef, useState, type CSSProperties } from 'react';
import { appHeaderActionLinkStyle } from '@/components/common/AppHeaderShell';

interface LearnerUser {
  email?: string;
  display_id?: string;
  nickname?: string;
  avatar_url?: string;
  total_points?: number;
}

interface LearnerUserMenuPalette {
  triggerText: string;
  pointText: string;
  pointBackground: string;
  pointBorder: string;
}

export interface LearnerUserMenuCopy {
  fallbackName: string;
  profile: string;
  points: string;
  aiSettings: string;
  logout: string;
}

export const defaultLearnerUserMenuCopy: LearnerUserMenuCopy = {
  fallbackName: '학습자',
  profile: '학습자정보',
  points: '포인트',
  aiSettings: 'AI 설정',
  logout: 'Logout',
};

function getDisplayName(user: LearnerUser | null, fallbackName: string) {
  if (!user) return fallbackName;
  return user.nickname || user.display_id || user.email || fallbackName;
}

function formatPointLabel(points?: number) {
  const safePoints = Number.isFinite(points) ? Math.max(0, Math.floor(points ?? 0)) : 0;
  return `${String(safePoints).padStart(2, '0')}pt`;
}

export default function LearnerUserMenu({
  palette,
  copy = defaultLearnerUserMenuCopy,
  onLogoutRequest,
}: {
  palette?: LearnerUserMenuPalette;
  copy?: LearnerUserMenuCopy;
  onLogoutRequest?: () => void;
}) {
  const [user, setUser] = useState<LearnerUser | null>(null);
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    let cancelled = false;

    const load = async () => {
      try {
        const refreshRes = await fetch('/api/v1/auth/refresh', { method: 'POST', credentials: 'include' });
        if (!refreshRes.ok) return;
        const meRes = await fetch('/api/v1/auth/me', { credentials: 'include', cache: 'no-store' });
        if (!meRes.ok) return;
        const data = (await meRes.json()) as LearnerUser;
        if (!cancelled) {
          setUser(data);
        }
      } catch {
        // noop
      }
    };

    load();
    const intervalId = window.setInterval(load, 30000);
    const handleFocus = () => void load();
    const handleUserRefresh = () => void load();
    const handlePointerDown = (event: MouseEvent) => {
      if (!containerRef.current?.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    window.addEventListener('focus', handleFocus);
    window.addEventListener('learnweaver:user-refresh', handleUserRefresh as EventListener);
    document.addEventListener('mousedown', handlePointerDown);

    return () => {
      cancelled = true;
      window.clearInterval(intervalId);
      window.removeEventListener('focus', handleFocus);
      window.removeEventListener('learnweaver:user-refresh', handleUserRefresh as EventListener);
      document.removeEventListener('mousedown', handlePointerDown);
    };
  }, []);

  const pointLabel = useMemo(() => formatPointLabel(user?.total_points), [user?.total_points]);

  return (
    <div ref={containerRef} style={wrapperStyle}>
      <button type="button" onClick={() => setIsOpen((current) => !current)} style={triggerStyle(isOpen)}>
        {user?.avatar_url ? (
          <img src={user.avatar_url} alt="" style={avatarStyle} />
        ) : null}
        <span style={nameStyle(palette)}>{getDisplayName(user, copy.fallbackName)}</span>
        <span style={pointStyle(palette)}>{pointLabel}</span>
      </button>

      {isOpen ? (
        <div style={menuStyle}>
          <Link href="/dashboard/settings/profile" style={menuItemStyle} onClick={() => setIsOpen(false)}>
            {copy.profile}
          </Link>
          <Link href="/dashboard/settings/points" style={menuItemStyle} onClick={() => setIsOpen(false)}>
            {copy.points}
          </Link>
          <Link href="/dashboard/settings/ai" style={menuItemStyle} onClick={() => setIsOpen(false)}>
            {copy.aiSettings}
          </Link>
          {onLogoutRequest ? (
            <>
              <div style={menuDividerStyle} />
              <button
                type="button"
                style={logoutMenuItemStyle}
                onClick={() => {
                  setIsOpen(false);
                  onLogoutRequest();
                }}
              >
                {copy.logout}
              </button>
            </>
          ) : null}
        </div>
      ) : null}
    </div>
  );
}

const wrapperStyle: CSSProperties = {
  position: 'relative',
};

const triggerStyle = (isOpen: boolean): CSSProperties => ({
  ...appHeaderActionLinkStyle(isOpen),
  minHeight: '38px',
  padding: '6px 8px 6px 10px',
  cursor: 'pointer',
  fontFamily: 'inherit',
  display: 'inline-flex',
  alignItems: 'center',
  gap: '8px',
  background: isOpen ? 'rgba(255,255,255,0.84)' : 'rgba(255,255,255,0.62)',
  border: '1px solid rgba(38,99,133,0.2)',
  boxShadow: isOpen ? '0 10px 28px rgba(42, 94, 139, 0.16)' : '0 8px 18px rgba(42, 94, 139, 0.1)',
});

const nameStyle = (palette?: LearnerUserMenuPalette): CSSProperties => ({
  maxWidth: '120px',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
  fontWeight: 900,
  color: palette?.triggerText ?? '#12364A',
});

const avatarStyle: CSSProperties = {
  width: '24px',
  height: '24px',
  borderRadius: '999px',
  objectFit: 'cover',
  border: '1px solid rgba(38,99,133,0.24)',
};

const pointStyle = (palette?: LearnerUserMenuPalette): CSSProperties => ({
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '22px',
  minWidth: '44px',
  padding: '0 8px',
  borderRadius: '8px',
  background: palette?.pointBackground ?? 'rgba(255, 247, 214, 0.96)',
  border: `1px solid ${palette?.pointBorder ?? 'rgba(245,158,11,0.58)'}`,
  color: palette?.pointText ?? '#78350F',
  fontSize: '12px',
  fontWeight: 900,
});

const menuStyle: CSSProperties = {
  position: 'absolute',
  top: 'calc(100% + 10px)',
  right: 0,
  minWidth: '180px',
  padding: '8px',
  borderRadius: '16px',
  border: '1px solid rgba(180, 205, 255, 0.16)',
  background: 'rgba(9, 19, 35, 0.94)',
  backdropFilter: 'blur(16px)',
  WebkitBackdropFilter: 'blur(16px)',
  boxShadow: '0 20px 45px rgba(0, 0, 0, 0.28)',
  display: 'grid',
  gap: '4px',
  zIndex: 80,
};

const menuItemStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  minHeight: '40px',
  padding: '0 12px',
  borderRadius: '12px',
  color: '#E8EEFF',
  textDecoration: 'none',
  fontSize: '14px',
  background: 'transparent',
};

const menuDividerStyle: CSSProperties = {
  height: '1px',
  margin: '4px 2px',
  background: 'rgba(180, 205, 255, 0.14)',
};

const logoutMenuItemStyle: CSSProperties = {
  ...menuItemStyle,
  width: '100%',
  border: 0,
  fontFamily: 'inherit',
  cursor: 'pointer',
  color: '#FCA5A5',
  textAlign: 'left',
};
