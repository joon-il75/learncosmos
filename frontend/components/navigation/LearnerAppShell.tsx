'use client';

import type { ReactNode } from 'react';
import { useState } from 'react';
import type { Locale } from '@/lib/i18n/locales';
import LeftMainNav from './LeftMainNav';
import MobileMainNavDrawer from './MobileMainNavDrawer';
import SectionSubNav from './SectionSubNav';
import TopUtilityBar from './TopUtilityBar';

interface LearnerAppShellProps {
  locale?: Locale | string | null;
  children: ReactNode;
  contentClassName?: string;
  pageHeader?: ReactNode;
  rightSlot?: ReactNode;
  showSectionSubNav?: boolean;
  logoHref?: string;
}

export default function LearnerAppShell({
  locale,
  children,
  contentClassName,
  pageHeader,
  rightSlot,
  showSectionSubNav = true,
  logoHref = '/',
}: LearnerAppShellProps) {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const menuLabel = locale === 'en' ? 'Open main menu' : '메인 메뉴 열기';
  const drawerTitle = locale === 'en' ? 'Main Menu' : '메인 메뉴';

  return (
    <div className="min-h-screen bg-[#F5FAFA] text-[#12384E]">
      <TopUtilityBar
        logoHref={logoHref}
        rightSlot={rightSlot}
        menuLabel={menuLabel}
        onMenuToggle={() => setDrawerOpen(true)}
      />
      <MobileMainNavDrawer open={drawerOpen} locale={locale} title={drawerTitle} onClose={() => setDrawerOpen(false)} />

      <div className="flex w-full pt-[60px]">
        <aside className="sticky top-[60px] hidden h-[calc(100vh-60px)] w-[240px] flex-shrink-0 border-r border-[#D6E6ED]/70 bg-[#EEF8F7]/76 p-4 lg:block">
          <LeftMainNav locale={locale} />
        </aside>
        <main className="min-w-0 flex-1">
          {showSectionSubNav ? <SectionSubNav locale={locale} /> : null}
          {pageHeader ? <div className="px-4 py-5 lg:px-6">{pageHeader}</div> : null}
          <div className={`px-4 pb-10 lg:px-6 ${contentClassName ?? ''}`}>{children}</div>
        </main>
      </div>
    </div>
  );
}
