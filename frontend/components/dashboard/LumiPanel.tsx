'use client';

import type { CSSProperties } from 'react';
import { getSpriteTuningByCoordinate } from '@/lib/lumi/lumiSpriteMap';
import type { LumiQuickAction } from '@/lib/lumi/lumiTypes';

interface LumiPanelProps {
  open: boolean;
  title?: string;
  message: string;
  subMessage?: string;
  sprite: string;
  isPhoneLayout: boolean;
  actions?: LumiQuickAction[];
  onClose?: () => void;
}


export default function LumiPanel({
  open,
  title = 'Lumi · 루미',
  message,
  subMessage,
  sprite,
  isPhoneLayout,
  actions = [],
  onClose,
}: LumiPanelProps) {
  if (!open) return null;

  return (
    <div style={panelWrapperStyle(isPhoneLayout)}>
      <div style={panelCardStyle(isPhoneLayout)}>
        <div style={panelHeaderStyle}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
            <div style={lumiSpriteStyle(sprite)} />
            <div style={{ display: 'grid', gap: 4 }}>
              <strong style={{ fontSize: 14, color: '#F4F7FF' }}>{title}</strong>
              <span style={{ fontSize: 12, color: '#AFC0E4' }}>탐험 파트너</span>
            </div>
          </div>
          {isPhoneLayout ? (
            <button type="button" onClick={onClose} style={closeButtonStyle} aria-label="Lumi 패널 닫기">
              닫기
            </button>
          ) : null}
        </div>

        <div style={messageBubbleStyle}>
          <div>{message}</div>
          {subMessage ? <div style={subMessageStyle}>{subMessage}</div> : null}
        </div>

        {actions.length > 0 ? (
          <div style={quickActionRowStyle}>
            {actions.slice(0, 2).map((action) => (
              <button key={action.id} type="button" style={quickActionButtonStyle} onClick={action.action}>
                {action.label}
              </button>
            ))}
          </div>
        ) : null}
      </div>
    </div>
  );
}

function panelWrapperStyle(isPhoneLayout: boolean): CSSProperties {
  return isPhoneLayout
    ? {
        position: 'fixed',
        inset: 'auto 0 0 0',
        zIndex: 60,
        padding: '16px',
        background: 'linear-gradient(180deg, rgba(7,14,26,0) 0%, rgba(7,14,26,0.48) 30%, rgba(7,14,26,0.82) 100%)',
      }
    : {
        display: 'block',
      };
}

function panelCardStyle(isPhoneLayout: boolean): CSSProperties {
  return {
    display: 'grid',
    gap: 14,
    borderRadius: isPhoneLayout ? 24 : 22,
    border: '1px solid rgba(126, 163, 255, 0.2)',
    background: 'linear-gradient(180deg, rgba(18, 31, 52, 0.84) 0%, rgba(12, 22, 38, 0.88) 100%)',
    padding: isPhoneLayout ? 18 : 20,
    boxShadow: '0 22px 44px rgba(3, 8, 18, 0.24)',
    backdropFilter: 'blur(18px)',
  };
}

const panelHeaderStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 12,
};

function lumiSpriteStyle(sprite: string): CSSProperties {
  const tuning = tuningFromSpriteToken(sprite);
  return {
    width: 40,
    height: 40,
    borderRadius: '999px',
    backgroundImage: "url('/images/lumi.webp')",
    backgroundRepeat: 'no-repeat',
    backgroundSize: `${tuning.backgroundSizeX}% ${tuning.backgroundSizeY}%`,
    backgroundPosition: `${tuning.backgroundPositionX}% ${tuning.backgroundPositionY}%`,
  };
}

function tuningFromSpriteToken(sprite: string) {
  const match = /^\[(\d+),(\d+)\]$/.exec(sprite);
  if (!match) return getSpriteTuningByCoordinate(1, 1);
  return getSpriteTuningByCoordinate(Number(match[1]), Number(match[2]));
}

const messageBubbleStyle: CSSProperties = {
  borderRadius: '8px 18px 18px 18px',
  padding: '12px 14px',
  background: 'rgba(33, 58, 92, 0.64)',
  border: '1px solid rgba(162, 191, 255, 0.1)',
  color: 'rgba(220, 232, 255, 0.92)',
  fontSize: 13,
  lineHeight: 1.7,
};

const subMessageStyle: CSSProperties = {
  marginTop: 6,
  color: 'rgba(190, 208, 235, 0.86)',
  fontSize: 12,
  lineHeight: 1.6,
};

const quickActionRowStyle: CSSProperties = {
  display: 'flex',
  flexWrap: 'wrap',
  gap: 8,
};

const quickActionButtonStyle: CSSProperties = {
  minHeight: 36,
  padding: '0 12px',
  borderRadius: 999,
  border: '1px solid rgba(194,210,245,0.18)',
  background: 'rgba(255,255,255,0.06)',
  color: 'rgba(220,232,255,0.84)',
  fontSize: 12,
  fontFamily: 'inherit',
  cursor: 'pointer',
};

const closeButtonStyle: CSSProperties = {
  minHeight: 36,
  padding: '0 12px',
  borderRadius: 999,
  border: '1px solid rgba(194,210,245,0.18)',
  background: 'rgba(255,255,255,0.06)',
  color: '#F4F7FF',
  fontSize: 12,
  fontFamily: 'inherit',
  cursor: 'pointer',
};
