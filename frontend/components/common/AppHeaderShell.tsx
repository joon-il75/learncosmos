'use client';

import type { CSSProperties, ReactNode } from 'react';
import BrandLogo from '@/components/common/BrandLogo';

interface AppHeaderShellProps {
  logoHref?: string;
  logoIconSize?: number;
  logoTextSize?: string;
  leftMeta?: ReactNode;
  centerSlot?: ReactNode;
  rightSlot?: ReactNode;
  position?: CSSProperties['position'];
  maxWidth?: string;
  headerStyle?: CSSProperties;
  innerStyle?: CSSProperties;
  leftGroupStyle?: CSSProperties;
  centerStyle?: CSSProperties;
  rightStyle?: CSSProperties;
  logoWrapperStyle?: CSSProperties;
}

export default function AppHeaderShell({
  logoHref = '/',
  logoIconSize = 32,
  logoTextSize = '18px',
  leftMeta,
  centerSlot,
  rightSlot,
  position = 'fixed',
  maxWidth = '72rem',
  headerStyle,
  innerStyle,
  leftGroupStyle,
  centerStyle,
  rightStyle,
  logoWrapperStyle,
}: AppHeaderShellProps) {
  const resolvedHeaderStyle: CSSProperties = {
    position,
    top: 0,
    left: 0,
    right: 0,
    zIndex: 50,
    borderBottom: '1px solid rgba(38, 99, 133, 0.12)',
    background: 'rgba(255,255,255,0.18)',
    backdropFilter: 'blur(12px) saturate(1.08)',
    WebkitBackdropFilter: 'blur(12px) saturate(1.08)',
    ...headerStyle,
  };

  const resolvedInnerStyle: CSSProperties = {
    width: '100%',
    maxWidth,
    margin: '0 auto',
    padding: '16px 24px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: '16px',
    minWidth: 0,
    ...innerStyle,
  };

  return (
    <header style={resolvedHeaderStyle}>
      <div style={resolvedInnerStyle}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '14px', flexShrink: 0, ...leftGroupStyle }}>
          <div style={{ flexShrink: 0, ...logoWrapperStyle }}>
            <BrandLogo href={logoHref} iconSize={logoIconSize} textSize={logoTextSize} />
          </div>
          {leftMeta}
        </div>

        {centerSlot ? (
          <div
            style={{
              flex: '1 1 auto',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              minWidth: 0,
              ...centerStyle,
            }}
          >
            {centerSlot}
          </div>
        ) : null}

        {rightSlot ? (
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '10px',
              flexWrap: 'nowrap',
              justifyContent: 'flex-end',
              marginLeft: 'auto',
              flexShrink: 0,
              ...rightStyle,
            }}
          >
            {rightSlot}
          </div>
        ) : null}
      </div>
    </header>
  );
}

export const appHeaderActionLinkStyle = (active = false): CSSProperties => ({
  fontSize: '13px',
  fontWeight: 600,
  color: active ? '#102F43' : '#254F65',
  textDecoration: 'none',
  padding: '6px 14px',
  borderRadius: '8px',
  background: active ? 'rgba(255,255,255,0.82)' : 'rgba(255,255,255,0.58)',
  border: active ? '1px solid rgba(77,183,232,0.34)' : '1px solid rgba(38,99,133,0.18)',
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: '34px',
  boxShadow: active ? '0 10px 24px rgba(42,94,139,0.14)' : '0 8px 18px rgba(42,94,139,0.08)',
  transition: 'background 180ms ease, border-color 180ms ease, color 180ms ease, box-shadow 180ms ease',
});

export const appHeaderActionButtonStyle = (active = false): CSSProperties => ({
  ...appHeaderActionLinkStyle(active),
  fontFamily: 'inherit',
  cursor: 'pointer',
});

export const appHeaderTagStyle: CSSProperties = {
  ...appHeaderActionLinkStyle(false),
  color: 'rgba(37,79,101,0.72)',
  cursor: 'default',
};
