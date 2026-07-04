'use client';

import type { ReactNode } from 'react';
import BrandLogo from '@/components/common/BrandLogo';

interface TopUtilityBarProps {
  logoHref?: string;
  rightSlot?: ReactNode;
  onMenuToggle?: () => void;
  menuLabel: string;
}

export default function TopUtilityBar({
  logoHref = '/',
  rightSlot,
  onMenuToggle,
  menuLabel,
}: TopUtilityBarProps) {
  return (
    <header className="fixed left-0 right-0 top-0 z-50 border-b border-[#D6E6ED]/70 bg-white/78 backdrop-blur-xl">
      <div className="flex h-[60px] w-full items-center justify-between gap-3 px-4 lg:px-6">
        <div className="flex min-w-0 items-center gap-3">
          {onMenuToggle ? (
            <button
              type="button"
              className="grid h-10 w-10 place-items-center rounded-full bg-[#EDF7F7] text-[#163E52] shadow-[0_10px_24px_rgba(33,73,93,0.10)] lg:hidden"
              aria-label={menuLabel}
              onClick={onMenuToggle}
            >
              <span className="grid gap-1.5" aria-hidden="true">
                <span className="block h-0.5 w-5 rounded-full bg-current" />
                <span className="block h-0.5 w-5 rounded-full bg-current" />
                <span className="block h-0.5 w-5 rounded-full bg-current" />
              </span>
            </button>
          ) : null}
          <BrandLogo href={logoHref} iconSize={30} textSize="18px" />
        </div>
        {rightSlot ? <div className="flex flex-shrink-0 items-center justify-end gap-2">{rightSlot}</div> : null}
      </div>
    </header>
  );
}
