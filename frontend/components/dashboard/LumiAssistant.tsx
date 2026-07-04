'use client';

import type { CSSProperties } from 'react';
import { getSpriteTuningByCoordinate } from '@/lib/lumi/lumiSpriteMap';
import type { LumiDockSlot } from './useLumiController';


interface LumiAssistantProps {
  dockSlot: LumiDockSlot;
  message: string;
  preview: string;
  sprite: string;
  visible: boolean;
  onTogglePanel: () => void;
  isPanelOpen: boolean;
  isPhoneLayout: boolean;
}

export default function LumiAssistant({
  dockSlot,
  message,
  preview,
  sprite,
  visible,
  onTogglePanel,
  isPanelOpen,
  isPhoneLayout,
}: LumiAssistantProps) {
  if (!visible && dockSlot !== 'mobile-bottom-sheet-trigger') return null;

  return (
    <div style={assistantWrapperStyle(dockSlot, visible, isPhoneLayout)}>
      <button
        type="button"
        aria-label="Lumi 패널 열기"
        onClick={onTogglePanel}
        style={assistantButtonStyle}
      >
        <div style={avatarPulseStyle}>
          <div style={lumiSpriteStyle(sprite)} />
        </div>
        <div style={assistantTextStyle}>
          <span style={assistantNameStyle}>Lumi</span>
          {(isPhoneLayout ? preview : message).split('\n').map((line, i) => (
            <span key={i} style={assistantPreviewStyle}>{line}</span>
          ))}
        </div>
        <span style={assistantChevronStyle(isPanelOpen)}>⌄</span>
      </button>
    </div>
  );
}

function assistantWrapperStyle(dockSlot: LumiDockSlot, visible: boolean, isPhoneLayout: boolean): CSSProperties {
  if (dockSlot === 'mobile-bottom-sheet-trigger') {
    return {
      position: isPhoneLayout ? 'sticky' : 'relative',
      bottom: isPhoneLayout ? '20px' : undefined,
      display: 'flex',
      justifyContent: isPhoneLayout ? 'flex-end' : 'flex-start',
      zIndex: 7,
      opacity: visible || isPhoneLayout ? 1 : 0,
      pointerEvents: visible || isPhoneLayout ? 'auto' : 'none',
    };
  }

  return {
    position: 'absolute',
    top: dockSlot === 'dashboard-cta' ? '16px' : dockSlot === 'dashboard-map-bottom' ? 'auto' : '18px',
    bottom: dockSlot === 'dashboard-map-bottom' ? '18px' : 'auto',
    right: dockSlot === 'dashboard-floating' ? '18px' : '16px',
    left: dockSlot === 'dashboard-map-bottom' ? '18px' : 'auto',
    zIndex: 7,
    opacity: visible ? 1 : 0,
    pointerEvents: visible ? 'auto' : 'none',
    transition: 'opacity 0.22s ease, transform 0.22s ease',
  };
}

const assistantButtonStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: '10px',
  minHeight: '48px',
  maxWidth: '320px',
  padding: '8px 12px 8px 8px',
  borderRadius: '999px',
  border: '1px solid rgba(126, 163, 255, 0.26)',
  background: 'rgba(14, 26, 45, 0.82)',
  boxShadow: '0 18px 40px rgba(3, 8, 18, 0.26)',
  backdropFilter: 'blur(16px)',
  color: '#F4F7FF',
  cursor: 'pointer',
  fontFamily: 'inherit',
  textAlign: 'left',
};

const avatarPulseStyle: CSSProperties = {
  width: '40px',
  height: '40px',
  borderRadius: '999px',
  flexShrink: 0,
  boxShadow: '0 0 0 3px rgba(50,200,100,0.0), 0 0 10px 2px rgba(50,200,100,0.22)',
};

function lumiSpriteStyle(sprite: string): CSSProperties {
  const tuning = tuningFromSpriteToken(sprite);
  return {
    width: '40px',
    height: '40px',
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

const assistantTextStyle: CSSProperties = {
  display: 'grid',
  gap: '2px',
  minWidth: 0,
  flex: 1,
};

const assistantNameStyle: CSSProperties = {
  fontSize: '11px',
  fontWeight: 700,
  color: '#E7EDFF',
  letterSpacing: '0.04em',
  textTransform: 'uppercase',
};

const assistantPreviewStyle: CSSProperties = {
  fontSize: '12px',
  color: 'rgba(200, 218, 245, 0.9)',
  lineHeight: 1.45,
};

function assistantChevronStyle(isOpen: boolean): CSSProperties {
  return {
    fontSize: '18px',
    color: 'rgba(200, 210, 235, 0.45)',
    transform: isOpen ? 'rotate(180deg)' : 'rotate(0deg)',
    transition: 'transform 0.2s ease',
    lineHeight: 1,
  };
}
