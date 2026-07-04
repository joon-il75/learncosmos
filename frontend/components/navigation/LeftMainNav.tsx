'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { getLearnerNavigationModel, isLearnerMainNavActive, type LearnerMainSection } from './learnerNavigationModel';
import type { SVGProps } from 'react';
import type { Locale } from '@/lib/i18n/locales';

interface LeftMainNavProps {
  locale?: Locale | string | null;
  collapsed?: boolean;
  onNavigate?: () => void;
}

const activeBarStyle = 'before:absolute before:left-0 before:top-1/2 before:h-7 before:w-1 before:-translate-y-1/2 before:rounded-r-full before:bg-[#2E708B]';

function NavIcon({ section, ...props }: SVGProps<SVGSVGElement> & { section: LearnerMainSection }) {
  const common = {
    fill: 'none',
    stroke: 'currentColor',
    strokeWidth: 2,
    strokeLinecap: 'round' as const,
    strokeLinejoin: 'round' as const,
  };

  if (section === 'home') {
    return (
      <svg viewBox="0 0 24 24" aria-hidden="true" {...props}>
        <path {...common} d="M4 11.5 12 5l8 6.5" />
        <path {...common} d="M6.5 10.5V19h11v-8.5" />
        <path {...common} d="M10 19v-5h4v5" />
      </svg>
    );
  }
  if (section === 'personal') {
    return (
      <svg viewBox="0 0 24 24" aria-hidden="true" {...props}>
        <path {...common} d="M5 6.5A2.5 2.5 0 0 1 7.5 4H19v14.5H7.5A2.5 2.5 0 0 0 5 21z" />
        <path {...common} d="M5 6.5v14" />
        <path {...common} d="M9 8h6" />
        <path {...common} d="M9 11h4" />
      </svg>
    );
  }
  if (section === 'community') {
    return (
      <svg viewBox="0 0 24 24" aria-hidden="true" {...props}>
        <circle {...common} cx="9" cy="8" r="3" />
        <circle {...common} cx="17" cy="10" r="2.5" />
        <path {...common} d="M3.8 19a5.2 5.2 0 0 1 10.4 0" />
        <path {...common} d="M14.5 18.5a4.1 4.1 0 0 1 5.7 0" />
      </svg>
    );
  }
  if (section === 'creator') {
    return (
      <svg viewBox="0 0 24 24" aria-hidden="true" {...props}>
        <path {...common} d="M4 20l4.2-1 10-10a2.1 2.1 0 0 0-3-3l-10 10z" />
        <path {...common} d="m13.8 7.2 3 3" />
        <path {...common} d="M5 13.5 10.5 19" />
      </svg>
    );
  }
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true" {...props}>
      <path {...common} d="M12 4v16" />
      <path {...common} d="M4 8h16" />
      <path {...common} d="M6 16h12" />
      <path {...common} d="M7.5 4.5a16 16 0 0 0 0 15" />
      <path {...common} d="M16.5 4.5a16 16 0 0 1 0 15" />
    </svg>
  );
}

export default function LeftMainNav({ locale, collapsed = false, onNavigate }: LeftMainNavProps) {
  const pathname = usePathname();
  const model = getLearnerNavigationModel(locale);

  return (
    <nav aria-label="Main navigation" className="grid gap-1">
      {model.mainItems.map((item) => {
        const active = isLearnerMainNavActive(pathname, item.key as LearnerMainSection);
        return (
          <Link
            key={item.key}
            href={item.href}
            title={item.label}
            onClick={onNavigate}
            className={`relative flex min-h-11 items-center gap-3 rounded-[18px] px-3 py-2 text-sm font-black transition ${active ? `${activeBarStyle} bg-[#E8F5F7] text-[#12384E]` : 'text-[#416679] hover:bg-white/70 hover:text-[#12384E]'}`}
          >
            <span className={`grid h-8 w-8 flex-shrink-0 place-items-center rounded-full ${active ? 'bg-[#12384E] text-white' : 'bg-white/70 text-[#416679]'}`}>
              <NavIcon section={item.key as LearnerMainSection} className="h-[18px] w-[18px]" />
            </span>
            {!collapsed ? <span className="truncate">{item.label}</span> : null}
          </Link>
        );
      })}
    </nav>
  );
}
