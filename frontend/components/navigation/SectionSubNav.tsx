'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import {
  getLearnerMainSectionFromPath,
  getLearnerNavigationModel,
  isLearnerSubNavActive,
  type LearnerMainSection,
} from './learnerNavigationModel';
import type { Locale } from '@/lib/i18n/locales';

interface SectionSubNavProps {
  locale?: Locale | string | null;
  mainSection?: LearnerMainSection;
}

export default function SectionSubNav({ locale, mainSection }: SectionSubNavProps) {
  const pathname = usePathname();
  const model = getLearnerNavigationModel(locale);
  const activeMain = mainSection ?? getLearnerMainSectionFromPath(pathname);
  const items = model.subItemsByMain[activeMain];

  return (
    <nav aria-label="Section navigation" className="sticky top-[60px] z-30 border-b border-[#D6E6ED]/70 bg-[#F7FBFB]/86 backdrop-blur-xl">
      <div className="flex min-h-12 items-center gap-2 overflow-x-auto px-4 lg:px-6">
        {items.map((item) => {
          const active = isLearnerSubNavActive(pathname, item);
          return (
            <Link
              key={item.key}
              href={item.href}
              className={`whitespace-nowrap rounded-full px-3.5 py-2 text-[13px] font-black transition ${active ? 'bg-[#12384E] text-white shadow-[0_10px_24px_rgba(18,56,78,0.16)]' : 'bg-white/64 text-[#416679] hover:bg-white hover:text-[#12384E]'}`}
            >
              {item.label}
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
