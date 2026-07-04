'use client';

import { useEffect, useState, type CSSProperties, type ReactNode } from 'react';
import LumiAvatar from '@/components/lumi/LumiAvatar';
import type { LumiState } from '@/lib/lumi/lumiTypes';
import ModalPortal from '@/components/common/ModalPortal';

interface LumiModalShellProps {
  title: string;
  message: ReactNode;
  lumiState: LumiState;
  eyebrow?: string;
  onClose: () => void;
  actions?: ReactNode;
  width?: number;
  ariaLabel?: string;
  tone?: 'warm' | 'alert';
}

export default function LumiModalShell({
  title,
  message,
  lumiState,
  eyebrow = 'Lumi Guide',
  onClose,
  actions,
  width = 420,
  ariaLabel,
  tone = 'warm',
}: LumiModalShellProps) {
  const [isMounted, setIsMounted] = useState(false);
  const palette = tone === 'alert' ? alertPalette : warmPalette;

  useEffect(() => {
    setIsMounted(true);
  }, []);

  if (!isMounted) return null;

  return (
    <ModalPortal overlayStyle={overlayStyle} onMouseDown={onClose}>
      <div
        role="dialog"
        aria-modal="true"
        aria-label={ariaLabel ?? title}
        style={{
          ...modalStyle,
          width,
          maxWidth: 'calc(100% - 32px)',
          background: palette.background,
          border: palette.border,
        }}
        onMouseDown={(event) => event.stopPropagation()}
      >
        <div style={headerStyle}>
          <div style={avatarWrapStyle}>
            <LumiAvatar state={lumiState} size={64} />
          </div>
          <div style={headerTextWrapStyle}>
            <div style={{ ...eyebrowStyle, color: palette.eyebrow }}>{eyebrow}</div>
            <h3 style={{ ...titleStyle, color: palette.title }}>{title}</h3>
          </div>
        </div>
        <div style={{ ...messageStyle, color: palette.body }}>{message}</div>
        <div style={actionsWrapStyle}>{actions}</div>
      </div>
    </ModalPortal>
  );
}

export const lumiModalSecondaryButtonStyle: CSSProperties = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '40px',
  padding: '0 16px',
  borderRadius: '999px',
  border: '1px solid rgba(128, 88, 28, 0.28)',
  background: 'rgba(255, 249, 238, 0.92)',
  color: '#6B461B',
  fontSize: '14px',
  fontWeight: 800,
  lineHeight: 1,
  whiteSpace: 'nowrap',
  cursor: 'pointer',
  boxShadow: '0 10px 24px rgba(92, 58, 14, 0.08)',
};

export const lumiModalPrimaryButtonStyle: CSSProperties = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '40px',
  padding: '0 16px',
  borderRadius: '999px',
  border: '1px solid rgba(140, 96, 22, 0.42)',
  background: 'linear-gradient(180deg, rgba(246, 205, 103, 0.96), rgba(214, 153, 34, 0.98))',
  color: '#4A2C00',
  fontSize: '14px',
  fontWeight: 900,
  lineHeight: 1,
  whiteSpace: 'nowrap',
  cursor: 'pointer',
  boxShadow: '0 12px 28px rgba(138, 95, 16, 0.18)',
};

const overlayStyle: CSSProperties = {
  padding: '24px 16px',
  background: 'rgba(22, 12, 8, 0.48)',
  backdropFilter: 'blur(4px)',
};

const modalStyle: CSSProperties = {
  borderRadius: 22,
  padding: '22px 22px 18px',
  boxShadow: '0 24px 70px rgba(0, 0, 0, 0.36)',
  display: 'grid',
  gap: 14,
  maxHeight: 'calc(100dvh - 48px)',
  overflowY: 'auto',
};

const headerStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: 14,
};

const avatarWrapStyle: CSSProperties = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  width: 74,
  height: 74,
  borderRadius: 999,
  background: 'radial-gradient(circle at 35% 30%, rgba(255,255,255,0.44), rgba(255,255,255,0.08) 72%, rgba(255,255,255,0))',
  boxShadow: 'inset 0 1px 0 rgba(255,255,255,0.55)',
  flexShrink: 0,
};

const headerTextWrapStyle: CSSProperties = {
  display: 'grid',
  gap: 4,
  minWidth: 0,
};

const eyebrowStyle: CSSProperties = {
  fontSize: 12,
  fontWeight: 800,
  letterSpacing: '0.08em',
  textTransform: 'uppercase',
};

const titleStyle: CSSProperties = {
  margin: 0,
  fontSize: 24,
  lineHeight: 1.25,
  fontWeight: 900,
};

const messageStyle: CSSProperties = {
  fontSize: 14,
  lineHeight: 1.7,
  whiteSpace: 'pre-line',
};

const actionsWrapStyle: CSSProperties = {
  display: 'flex',
  justifyContent: 'flex-end',
  gap: 10,
  flexWrap: 'wrap',
  rowGap: 10,
  paddingTop: 4,
};

const warmPalette = {
  background: 'linear-gradient(180deg, rgba(255, 245, 224, 0.98), rgba(245, 227, 192, 0.98))',
  border: '1px solid rgba(173, 120, 36, 0.42)',
  eyebrow: 'rgba(120, 70, 18, 0.86)',
  title: '#4B2E09',
  body: 'rgba(74, 44, 10, 0.92)',
};

const alertPalette = {
  background: 'linear-gradient(180deg, rgba(255, 240, 232, 0.98), rgba(248, 223, 213, 0.98))',
  border: '1px solid rgba(172, 88, 62, 0.34)',
  eyebrow: 'rgba(145, 58, 36, 0.84)',
  title: '#5C2418',
  body: 'rgba(92, 36, 24, 0.92)',
};
