'use client';

import LeftMainNav from './LeftMainNav';
import type { Locale } from '@/lib/i18n/locales';

interface MobileMainNavDrawerProps {
  open: boolean;
  locale?: Locale | string | null;
  title: string;
  onClose: () => void;
}

export default function MobileMainNavDrawer({ open, locale, title, onClose }: MobileMainNavDrawerProps) {
  if (!open) return null;

  return (
    <div className="fixed inset-0 z-[80] lg:hidden" role="presentation">
      <button
        type="button"
        className="absolute inset-0 bg-[#061525]/42 backdrop-blur-[2px]"
        aria-label={title}
        onClick={onClose}
      />
      <aside className="absolute bottom-3 left-3 top-3 w-[min(312px,calc(100%-24px))] overflow-hidden rounded-[30px] bg-[#F7FCFB] p-4 shadow-[0_28px_90px_rgba(5,24,38,0.30)]">
        <div className="mb-4 flex items-center justify-between gap-3">
          <strong className="text-sm font-black uppercase tracking-[0.16em] text-[#416679]">{title}</strong>
          <button
            type="button"
            className="grid h-9 w-9 place-items-center rounded-full bg-white text-xl font-black text-[#12384E]"
            onClick={onClose}
            aria-label="Close menu"
          >
            ×
          </button>
        </div>
        <LeftMainNav locale={locale} onNavigate={onClose} />
      </aside>
    </div>
  );
}
