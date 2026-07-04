'use client';

import type { CSSProperties } from 'react';
import RotatingPlanetCanvas from '@/components/dashboard/RotatingPlanetCanvas';

type RotationDirection = 'left' | 'right';
type AtmospherePreset = 'admin-preview' | 'star-system' | 'galaxy-mini' | 'goal-loading';

interface AtmosphericPlanetPreviewProps {
  atlasURL: string | null;
  progressPercent: number;
  size: number;
  durationSeconds: number;
  direction: RotationDirection;
  playing?: boolean;
  preset?: AtmospherePreset;
  style?: CSSProperties;
  canvasStyle?: CSSProperties;
}

interface AtmospherePresetConfig {
  paddingRatio: number;
  minPadding: number;
  haloScale: number;
  sheenScale: number;
  haloOpacity: number;
  sheenOpacity: number;
  blurRatio: number;
  animateHalo: boolean;
  animateSheen: boolean;
  dropShadow: string;
}

const PRESET_CONFIG: Record<AtmospherePreset, AtmospherePresetConfig> = {
  'admin-preview': {
    paddingRatio: 0.12,
    minPadding: 10,
    haloScale: 1.2,
    sheenScale: 1.04,
    haloOpacity: 1,
    sheenOpacity: 0.74,
    blurRatio: 0.01,
    animateHalo: true,
    animateSheen: true,
    dropShadow: 'drop-shadow(0 18px 28px rgba(0, 0, 0, 0.34))',
  },
  'star-system': {
    paddingRatio: 0,
    minPadding: 0,
    haloScale: 1.05,
    sheenScale: 0.98,
    haloOpacity: 0.74,
    sheenOpacity: 0.42,
    blurRatio: 0.006,
    animateHalo: true,
    animateSheen: true,
    dropShadow: 'drop-shadow(0 12px 20px rgba(0, 0, 0, 0.22))',
  },
  'galaxy-mini': {
    paddingRatio: 0,
    minPadding: 0,
    haloScale: 1,
    sheenScale: 0.92,
    haloOpacity: 0.42,
    sheenOpacity: 0.18,
    blurRatio: 0,
    animateHalo: false,
    animateSheen: false,
    dropShadow: 'drop-shadow(0 0 4px rgba(151, 210, 255, 0.22))',
  },
  'goal-loading': {
    paddingRatio: 0.16,
    minPadding: 12,
    haloScale: 1.24,
    sheenScale: 1.08,
    haloOpacity: 1,
    sheenOpacity: 0.82,
    blurRatio: 0.012,
    animateHalo: true,
    animateSheen: true,
    dropShadow: 'drop-shadow(0 18px 30px rgba(120, 80, 20, 0.24))',
  },
};

export default function AtmosphericPlanetPreview({
  atlasURL,
  progressPercent,
  size,
  durationSeconds,
  direction,
  playing = true,
  preset = 'admin-preview',
  style,
  canvasStyle,
}: AtmosphericPlanetPreviewProps) {
  const config = PRESET_CONFIG[preset];
  const padding = Math.max(config.minPadding, Math.round(size * config.paddingRatio));
  const frameSize = size + padding * 2;

  return (
    <div style={{ ...planetWrapStyle(frameSize), ...style }}>
      <style jsx global>{`
        @keyframes planetAtmospherePulse {
          0%, 100% { opacity: 0.72; transform: scale(0.99); }
          50% { opacity: 0.92; transform: scale(1.03); }
        }
        @keyframes planetAtmosphereDrift {
          from { transform: rotate(0deg) scale(1.02); }
          to { transform: rotate(360deg) scale(1.02); }
        }
        @media (prefers-reduced-motion: reduce) {
          .planetAtmosphereHalo,
          .planetAtmosphereSheen {
            animation: none !important;
          }
        }
      `}</style>
      <span className="planetAtmosphereHalo" style={atmosphereHaloStyle(size, config)} aria-hidden="true" />
      <span className="planetAtmosphereSheen" style={atmosphereSheenStyle(size, config)} aria-hidden="true" />
      <RotatingPlanetCanvas
        atlasURL={atlasURL}
        size={size}
        progressPercent={progressPercent}
        durationSeconds={durationSeconds}
        direction={direction}
        playing={playing}
        style={{ ...planetCanvasLayerStyle(config), ...canvasStyle }}
      />
    </div>
  );
}

function planetWrapStyle(frameSize: number): CSSProperties {
  return {
    position: 'relative',
    width: frameSize,
    height: frameSize,
    display: 'grid',
    placeItems: 'center',
    flex: '0 0 auto',
    pointerEvents: 'none',
  };
}

function atmosphereHaloStyle(size: number, config: AtmospherePresetConfig): CSSProperties {
  const diameter = Math.round(size * config.haloScale);
  const blur = Math.max(4, Math.round(size * 0.12));
  return {
    position: 'absolute',
    width: diameter,
    height: diameter,
    borderRadius: '50%',
    background:
      'radial-gradient(circle at 50% 50%, transparent 55%, rgba(172, 222, 255, 0.18) 67%, rgba(87, 172, 255, 0.32) 78%, rgba(36, 99, 190, 0.08) 100%)',
    boxShadow: `0 0 ${blur}px rgba(96, 181, 255, 0.42), inset 0 0 ${Math.round(blur * 0.7)}px rgba(236, 248, 255, 0.22)`,
    filter: config.blurRatio > 0 ? `blur(${Math.max(1, Math.round(size * config.blurRatio))}px)` : undefined,
    animation: config.animateHalo ? 'planetAtmospherePulse 5.8s ease-in-out infinite' : 'none',
    opacity: config.haloOpacity,
    pointerEvents: 'none',
    zIndex: 1,
  };
}

function atmosphereSheenStyle(size: number, config: AtmospherePresetConfig): CSSProperties {
  const diameter = Math.round(size * config.sheenScale);
  return {
    position: 'absolute',
    width: diameter,
    height: diameter,
    borderRadius: '50%',
    background:
      'conic-gradient(from 18deg, transparent 0deg, rgba(219, 245, 255, 0.18) 42deg, transparent 88deg, transparent 210deg, rgba(96, 181, 255, 0.12) 252deg, transparent 310deg)',
    mixBlendMode: 'screen',
    opacity: config.sheenOpacity,
    animation: config.animateSheen ? 'planetAtmosphereDrift 18s linear infinite' : 'none',
    pointerEvents: 'none',
    zIndex: 3,
  };
}

function planetCanvasLayerStyle(config: AtmospherePresetConfig): CSSProperties {
  return {
    position: 'relative',
    zIndex: 2,
    borderRadius: '999px',
    filter: config.dropShadow,
  };
}
