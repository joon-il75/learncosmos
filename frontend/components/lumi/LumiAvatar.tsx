'use client';

import type { CSSProperties } from 'react';

import { getDefaultSpriteTuning } from '@/lib/lumi/lumiSpriteMap';
import type { LumiState } from '@/lib/lumi/lumiTypes';

interface LumiAvatarProps {
  state: LumiState;
  size?: number;
  reducedMotion?: boolean;
}

export default function LumiAvatar({ state, size = 48, reducedMotion = false }: LumiAvatarProps) {
  const tuning = getDefaultSpriteTuning(state);

  return (
    <div
      aria-hidden="true"
      style={{
        ...avatarFrameStyle,
        width: size,
        height: size,
        animation: reducedMotion ? 'none' : 'lumiSystemBreath 2.6s ease-in-out infinite',
      }}
    >
      <div
        style={{
          width: size,
          height: size,
          borderRadius: '50%',
          backgroundImage: "url('/images/lumi.webp')",
          backgroundRepeat: 'no-repeat',
          backgroundSize: `${tuning.backgroundSizeX}% ${tuning.backgroundSizeY}%`,
          backgroundPosition: `${tuning.backgroundPositionX}% ${tuning.backgroundPositionY}%`,
        }}
      />
      <style jsx>{`
        @keyframes lumiSystemBreath {
          0%, 100% { transform: translateY(0); }
          50% { transform: translateY(-2px); }
        }
      `}</style>
    </div>
  );
}

const avatarFrameStyle: CSSProperties = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  borderRadius: '999px',
  background: 'radial-gradient(circle at 35% 30%, rgba(224, 247, 255, 0.22), rgba(35, 58, 92, 0.14) 68%, rgba(7, 12, 20, 0))',
  boxShadow: '0 10px 28px rgba(0, 0, 0, 0.22)',
};
