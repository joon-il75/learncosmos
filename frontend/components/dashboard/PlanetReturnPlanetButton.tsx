'use client';

import { useState, type CSSProperties, type ReactNode } from 'react';
import AtmosphericPlanetPreview from '@/components/dashboard/AtmosphericPlanetPreview';
import { PLANET_TEXTURE_MAP_FALLBACK_PATH } from '@/components/dashboard/planetSpriteSystem';
import { normalizePlanetProgressPercent } from '@/lib/world-ui-engine/planetProgressUtils';

type PlanetReturnPlanetButtonVariant = 'star' | 'map';

interface PlanetReturnPlanetButtonProps {
  ariaLabel: string;
  children: ReactNode;
  textureMapAsset?: string | null;
  progressPercent?: number | null;
  size?: number;
  variant?: PlanetReturnPlanetButtonVariant;
  style?: CSSProperties;
  onClick: () => void;
}

export default function PlanetReturnPlanetButton({
  ariaLabel,
  children,
  textureMapAsset,
  progressPercent,
  size = 107,
  variant = 'star',
  style,
  onClick,
}: PlanetReturnPlanetButtonProps) {
  const [active, setActive] = useState(false);
  const atlasPath = textureMapAsset ?? PLANET_TEXTURE_MAP_FALLBACK_PATH;
  const progress = normalizePlanetProgressPercent(progressPercent);

  return (
    <button
      type="button"
      aria-label={ariaLabel}
      onMouseEnter={() => setActive(true)}
      onMouseLeave={() => setActive(false)}
      onFocus={() => setActive(true)}
      onBlur={() => setActive(false)}
      onClick={onClick}
      style={{
        ...buttonBaseStyle(size, active, variant),
        ...style,
      }}
    >
      <span style={rotatingLayerStyle(active)} aria-hidden="true">
        <AtmosphericPlanetPreview
          atlasURL={atlasPath}
          progressPercent={progress}
          size={Math.max(70, Math.round(size * 0.84))}
          durationSeconds={18}
          direction="left"
          playing={active}
          preset="star-system"
          style={{ width: size, height: size, placeItems: 'center' }}
        />
      </span>
      <span style={contrastLayerStyle(active)} aria-hidden="true" />
      <span style={labelLayerStyle}>{children}</span>
    </button>
  );
}

function buttonBaseStyle(
  size: number,
  active: boolean,
  variant: PlanetReturnPlanetButtonVariant,
): CSSProperties {
  const glowColor = variant === 'star' ? '100, 170, 255' : '180, 138, 72';
  return {
    position: 'absolute',
    zIndex: 5,
    width: size,
    height: size,
    display: 'grid',
    placeItems: 'center',
    overflow: 'hidden',
    borderRadius: 999,
    border: `1.5px solid rgba(198, 226, 255, ${active ? '0.88' : '0.62'})`,
    color: '#ffffff',
    fontSize: 15,
    fontWeight: 900,
    lineHeight: 1.18,
    cursor: 'pointer',
    pointerEvents: 'auto',
    textAlign: 'center',
    padding: 0,
    background: active
      ? 'radial-gradient(circle at 44% 30%, rgba(60, 120, 220, 0.20), rgba(5, 12, 28, 0.82) 72%)'
      : 'radial-gradient(circle at 44% 30%, rgba(58, 102, 176, 0.16), rgba(5, 12, 28, 0.88) 72%)',
    boxShadow: active
      ? `0 0 0 6px rgba(${glowColor}, 0.28), 0 16px 36px rgba(30, 80, 180, 0.42), inset 0 2px 0 rgba(255,255,255,0.62), inset 0 -12px 18px rgba(0, 0, 0, 0.34)`
      : `0 0 0 4px rgba(${glowColor}, 0.14), 0 12px 28px rgba(30, 80, 180, 0.28), inset 0 2px 0 rgba(255,255,255,0.46), inset 0 -12px 18px rgba(0, 0, 0, 0.38)`,
    textShadow:
      '0 0 6px rgba(20, 60, 180, 0.9), 0 1px 3px rgba(0, 0, 0, 0.6), 0 2px 8px rgba(0, 30, 120, 0.7)',
  };
}

function rotatingLayerStyle(active: boolean): CSSProperties {
  return {
    position: 'absolute',
    inset: 0,
    display: 'grid',
    placeItems: 'center',
    pointerEvents: 'none',
    zIndex: 1,
    opacity: active ? 1 : 0.94,
    transition: 'opacity 180ms ease',
  };
}

function contrastLayerStyle(active: boolean): CSSProperties {
  return {
    position: 'absolute',
    inset: 0,
    pointerEvents: 'none',
    zIndex: 2,
    background: active
      ? 'radial-gradient(circle at 42% 30%, rgba(255,255,255,0.06), rgba(2, 8, 24, 0.22) 58%, rgba(1, 4, 14, 0.46) 100%)'
      : 'radial-gradient(circle at 42% 30%, rgba(255,255,255,0.04), rgba(2, 8, 24, 0.18) 58%, rgba(1, 4, 14, 0.52) 100%)',
  };
}

const labelLayerStyle: CSSProperties = {
  position: 'relative',
  zIndex: 3,
  display: 'block',
  whiteSpace: 'normal',
};
