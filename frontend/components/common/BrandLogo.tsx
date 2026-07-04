'use client';

import Link from 'next/link';
import type { CSSProperties } from 'react';

const logoSrc = '/images/learncosmos.svg';
const logoAspectRatio = 1140.8 / 298.3;

interface BrandLogoProps {
  href?: string;
  iconSize?: number;
  textSize?: string;
  style?: CSSProperties;
  textClassName?: string;
}

export default function BrandLogo({
  href = '/',
  iconSize = 32,
  textSize = '20px',
  style,
}: BrandLogoProps) {
  const parsedTextSize = Number.parseFloat(textSize);
  const minimumLogoHeight = 55;
  const logoHeight = Math.max(minimumLogoHeight, iconSize, Number.isFinite(parsedTextSize) ? parsedTextSize : iconSize);
  const logoWidth = Math.round(logoHeight * logoAspectRatio);

  return (
    <Link
      href={href}
      aria-label="LearnCosmos"
      style={{
        display: 'inline-flex',
        alignItems: 'center',
        textDecoration: 'none',
        ...style,
      }}
    >
      <img
        src={logoSrc}
        alt="LearnCosmos"
        width={logoWidth}
        height={logoHeight}
        style={{
          width: 'auto',
          height: `clamp(36px, 12vw, ${logoHeight}px)`,
          objectFit: 'contain',
          flexShrink: 0,
          filter: 'drop-shadow(1.2px 1.5px 0 rgba(22, 57, 94, 0.3))',
        }}
      />
    </Link>
  );
}
