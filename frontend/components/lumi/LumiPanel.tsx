'use client';

import type { CSSProperties } from 'react';

import LumiAvatar from './LumiAvatar';
import type { LumiQuickAction, LumiState } from '@/lib/lumi/lumiTypes';

interface LumiPanelProps {
  open: boolean;
  state: LumiState;
  message: string;
  quickActions?: LumiQuickAction[];
  onClose?: () => void;
  reducedMotion?: boolean;
}

export default function LumiPanel({
  open,
  state,
  message,
  quickActions = [],
  onClose,
  reducedMotion = false,
}: LumiPanelProps) {
  if (!open) return null;

  return (
    <section style={panelStyle} aria-label="Lumi 패널">
      <div style={headerStyle}>
        <div style={identityStyle}>
          <LumiAvatar state={state} size={40} reducedMotion={reducedMotion} />
          <div style={titleGroupStyle}>
            <strong style={titleStyle}>Lumi · 루미</strong>
            <span style={subtitleStyle}>탐험 파트너 시스템</span>
          </div>
        </div>
        {onClose ? (
          <button type="button" onClick={onClose} style={closeButtonStyle} aria-label="Lumi 패널 닫기">
            닫기
          </button>
        ) : null}
      </div>
      <p style={messageStyle}>{message}</p>
      {quickActions.length > 0 ? (
        <div style={actionsStyle}>
          {quickActions.map((action) => (
            <button key={action.id} type="button" onClick={action.action} style={actionButtonStyle}>
              {action.label}
            </button>
          ))}
        </div>
      ) : null}
    </section>
  );
}

const panelStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 14,
  padding: 18,
  borderRadius: 24,
  border: '1px solid rgba(122, 164, 255, 0.22)',
  background: 'linear-gradient(180deg, rgba(18, 29, 49, 0.96), rgba(10, 16, 28, 0.95))',
  boxShadow: '0 24px 54px rgba(0, 0, 0, 0.28)',
};

const headerStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 12,
};

const identityStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: 12,
};

const titleGroupStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 2,
};

const titleStyle: CSSProperties = {
  color: '#f6fbff',
  fontSize: 15,
};

const subtitleStyle: CSSProperties = {
  color: 'rgba(198, 216, 240, 0.76)',
  fontSize: 12,
};

const closeButtonStyle: CSSProperties = {
  border: '1px solid rgba(130, 164, 239, 0.28)',
  borderRadius: 999,
  background: 'rgba(255, 255, 255, 0.03)',
  color: '#e8f2ff',
  padding: '8px 12px',
  cursor: 'pointer',
};

const messageStyle: CSSProperties = {
  margin: 0,
  color: 'rgba(236, 244, 255, 0.94)',
  fontSize: 14,
  lineHeight: 1.6,
};

const actionsStyle: CSSProperties = {
  display: 'flex',
  flexWrap: 'wrap',
  gap: 8,
};

const actionButtonStyle: CSSProperties = {
  border: '1px solid rgba(120, 162, 241, 0.24)',
  borderRadius: 999,
  background: 'rgba(31, 49, 78, 0.72)',
  color: '#eef6ff',
  padding: '8px 12px',
  fontSize: 13,
  cursor: 'pointer',
};
