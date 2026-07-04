'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import type { MouseEvent } from 'react';
import { useMemo, useState, type CSSProperties } from 'react';
import type { Locale } from '@/lib/i18n/locales';

type ImmersiveOverlayMenuProps = {
  locale?: Locale | 'ko' | 'en';
  tone?: 'dark' | 'light';
  onBeforeNavigate?: (href: string) => boolean;
};

const labels = {
  ko: {
    open: '전역 메뉴 열기',
    close: '전역 메뉴 닫기',
    eyebrow: 'Global Menu',
    title: '어디로 이동할까요?',
    description: '현재 학습 흐름은 유지한 채 필요한 주요 영역으로 이동할 수 있습니다.',
    items: [
      { key: 'home', label: '홈', href: '/' },
      { key: 'personal', label: '개인학습', href: '/dashboard' },
      { key: 'community', label: '커뮤니티', href: '/dashboard/community' },
      { key: 'creator', label: '크리에이터', href: '/dashboard/creator' },
      { key: 'platform', label: '플랫폼', href: '/platform' },
    ],
  },
  en: {
    open: 'Open global menu',
    close: 'Close global menu',
    eyebrow: 'Global Menu',
    title: 'Where would you like to go?',
    description: 'Move to a main area without changing the current learning layout.',
    items: [
      { key: 'home', label: 'Home', href: '/' },
      { key: 'personal', label: 'Personal Learning', href: '/dashboard' },
      { key: 'community', label: 'Community', href: '/dashboard/community' },
      { key: 'creator', label: 'Creator', href: '/dashboard/creator' },
      { key: 'platform', label: 'Platform', href: '/platform' },
    ],
  },
} as const;

function isActive(pathname: string | null, href: string) {
  if (!pathname) return false;
  if (href === '/') return pathname === '/';
  return pathname === href || pathname.startsWith(`${href}/`);
}

export default function ImmersiveOverlayMenu({ locale = 'ko', tone = 'dark', onBeforeNavigate }: ImmersiveOverlayMenuProps) {
  const pathname = usePathname();
  const [open, setOpen] = useState(false);
  const copy = labels[locale === 'en' ? 'en' : 'ko'];
  const palette = useMemo(() => tone === 'light'
    ? {
        buttonBg: 'rgba(255,255,255,0.92)',
        buttonBorder: 'rgba(15,23,42,0.14)',
        buttonText: '#0F172A',
        panelBg: 'rgba(255,255,255,0.98)',
        panelText: '#102033',
        panelMuted: 'rgba(15,23,42,0.62)',
        activeBg: 'rgba(13,49,78,0.94)',
        activeText: '#FFFFFF',
        itemBg: 'rgba(15,23,42,0.05)',
        itemBorder: 'rgba(15,23,42,0.09)',
      }
    : {
        buttonBg: 'rgba(8,18,33,0.78)',
        buttonBorder: 'rgba(194,210,245,0.18)',
        buttonText: '#F4F7FF',
        panelBg: 'rgba(8,18,33,0.97)',
        panelText: '#F4F7FF',
        panelMuted: 'rgba(216,226,245,0.68)',
        activeBg: 'rgba(238,185,91,0.95)',
        activeText: '#2D1B06',
        itemBg: 'rgba(255,255,255,0.07)',
        itemBorder: 'rgba(194,210,245,0.14)',
      }, [tone]);

  return (
    <>
      <button
        type="button"
        aria-label={copy.open}
        aria-expanded={open}
        onClick={() => setOpen(true)}
        style={{ ...buttonStyle, background: palette.buttonBg, borderColor: palette.buttonBorder, color: palette.buttonText }}
      >
        <span style={buttonLineStyle} />
        <span style={buttonLineStyle} />
        <span style={buttonLineStyle} />
      </button>
      {open ? (
        <div style={overlayStyle} role="presentation" onMouseDown={() => setOpen(false)}>
          <nav
            aria-label={copy.eyebrow}
            style={{ ...panelStyle, background: palette.panelBg, color: palette.panelText }}
            onMouseDown={(event) => event.stopPropagation()}
          >
            <div style={panelHeaderStyle}>
              <div>
                <p style={{ ...eyebrowStyle, color: palette.panelMuted }}>{copy.eyebrow}</p>
                <h2 style={titleStyle}>{copy.title}</h2>
                <p style={{ ...descriptionStyle, color: palette.panelMuted }}>{copy.description}</p>
              </div>
              <button type="button" aria-label={copy.close} onClick={() => setOpen(false)} style={{ ...closeButtonStyle, color: palette.panelText }}>
                ×
              </button>
            </div>
            <div style={itemListStyle}>
              {copy.items.map((item) => {
                const active = isActive(pathname, item.href);
                const handleClick = (event: MouseEvent<HTMLAnchorElement>) => {
                  if (onBeforeNavigate && !onBeforeNavigate(item.href)) {
                    event.preventDefault();
                    return;
                  }
                  setOpen(false);
                };
                return (
                  <Link
                    key={item.key}
                    href={item.href}
                    onClick={handleClick}
                    style={{
                      ...itemStyle,
                      background: active ? palette.activeBg : palette.itemBg,
                      color: active ? palette.activeText : palette.panelText,
                      borderColor: active ? 'transparent' : palette.itemBorder,
                    }}
                  >
                    {item.label}
                  </Link>
                );
              })}
            </div>
          </nav>
        </div>
      ) : null}
    </>
  );
}

const buttonStyle: CSSProperties = {
  width: '38px',
  minWidth: '38px',
  height: '38px',
  borderRadius: '14px',
  border: '1px solid',
  display: 'inline-grid',
  placeItems: 'center',
  gap: '3px',
  padding: '9px',
  cursor: 'pointer',
  backdropFilter: 'blur(12px)',
  WebkitBackdropFilter: 'blur(12px)',
};

const buttonLineStyle: CSSProperties = {
  width: '16px',
  height: '2px',
  borderRadius: '999px',
  background: 'currentColor',
};

const overlayStyle: CSSProperties = {
  position: 'fixed',
  inset: 0,
  zIndex: 100,
  display: 'flex',
  justifyContent: 'flex-start',
  background: 'rgba(3, 7, 18, 0.48)',
  backdropFilter: 'blur(5px)',
  WebkitBackdropFilter: 'blur(5px)',
};

const panelStyle: CSSProperties = {
  width: 'min(340px, calc(100vw - 28px))',
  minHeight: '100%',
  padding: '24px 18px',
  boxShadow: '24px 0 70px rgba(0,0,0,0.32)',
  overflowY: 'auto',
};

const panelHeaderStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1fr) 40px',
  gap: '12px',
  alignItems: 'start',
};

const eyebrowStyle: CSSProperties = {
  margin: 0,
  fontSize: '11px',
  fontWeight: 900,
  letterSpacing: '0.12em',
  textTransform: 'uppercase',
};

const titleStyle: CSSProperties = {
  margin: '8px 0 0',
  fontSize: '25px',
  lineHeight: 1.12,
  letterSpacing: '-0.04em',
};

const descriptionStyle: CSSProperties = {
  margin: '10px 0 0',
  fontSize: '13px',
  lineHeight: 1.55,
};

const closeButtonStyle: CSSProperties = {
  width: '38px',
  height: '38px',
  border: 0,
  borderRadius: '999px',
  background: 'rgba(255,255,255,0.08)',
  fontSize: '25px',
  cursor: 'pointer',
};

const itemListStyle: CSSProperties = {
  display: 'grid',
  gap: '10px',
  marginTop: '24px',
};

const itemStyle: CSSProperties = {
  minHeight: '48px',
  display: 'flex',
  alignItems: 'center',
  padding: '0 15px',
  borderRadius: '16px',
  border: '1px solid',
  textDecoration: 'none',
  fontSize: '15px',
  fontWeight: 900,
};
