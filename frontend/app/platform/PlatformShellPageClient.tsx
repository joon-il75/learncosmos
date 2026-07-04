'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';

import LearnerAppShell from '@/components/navigation/LearnerAppShell';
import { normalizeLocale, type Locale } from '@/lib/i18n/locales';
import { getPlatformShellCopy, type PlatformShellSection } from '@/lib/i18n/pages/platform';
import {
  actionRowStyle,
  contentGridStyle,
  eyebrowStyle,
  heroStyle,
  linkBodyStyle,
  linkItemStyle,
  linkLabelStyle,
  linkListStyle,
  linksPanelStyle,
  mainStyle,
  newsBoardAdminOnlyStyle,
  newsBoardBadgeStyle,
  newsBoardDescriptionStyle,
  newsBoardEmptyBodyStyle,
  newsBoardEmptyStyle,
  newsBoardEmptyTitleStyle,
  newsBoardHeaderStyle,
  newsBoardStyle,
  newsBoardTitleStyle,
  mobileStyle,
  noteDotStyle,
  noteItemStyle,
  noteListStyle,
  pageStyle,
  panelSubtitleStyle,
  panelTitleStyle,
  primaryActionStyle,
  secondaryActionStyle,
  sectionDescriptionStyle,
  sectionEyebrowStyle,
  sectionSurfaceStyle,
  sectionTitleStyle,
  sidePanelStyle,
  signalDotStyle,
  signalStyle,
  standardBodyStyle,
  standardLabelStyle,
  standardListStyle,
  standardPanelStyle,
  statusLabelStyle,
  statusPanelStyle,
  subtitleStyle,
  titleStyle,
} from './platformShellStyles';

interface PlatformShellPageClientProps {
  section: PlatformShellSection;
}

const primaryHrefs: Record<PlatformShellSection, string> = {
  overview: '/platform/policies',
  news: '/alpha',
  guide: '/dashboard/goal',
  operations: '/dashboard/community',
  policies: '/terms',
};

export default function PlatformShellPageClient({ section }: PlatformShellPageClientProps) {
  const [uiLocale, setUiLocale] = useState<Locale>('ko');
  const copy = getPlatformShellCopy(uiLocale);
  const activeSection = copy.sections[section];
  const primaryHref = uiLocale === 'en' && primaryHrefs[section] === '/terms' ? '/en/terms' : primaryHrefs[section];

  useEffect(() => {
    setUiLocale(normalizeLocale(window.navigator.language?.startsWith('en') ? 'en' : 'ko'));
  }, []);

  const utilityLinks = (
    <nav className="flex items-center gap-2" aria-label="Platform quick links">
      <Link
        href="/"
        className="rounded-full bg-white/70 px-3 py-2 text-[13px] font-black text-[#254B5F] transition hover:bg-white"
      >
        {copy.shell.homeLabel}
      </Link>
      <Link
        href="/dashboard"
        className="hidden rounded-full bg-white/70 px-3 py-2 text-[13px] font-black text-[#254B5F] transition hover:bg-white min-[420px]:inline-flex"
      >
        {copy.shell.appLabel}
      </Link>
    </nav>
  );

  return (
    <LearnerAppShell locale={uiLocale} logoHref="/" rightSlot={utilityLinks} contentClassName="!px-0 !pb-0">
      <div style={pageStyle}>
        <style>{mobileStyle}</style>
        <main className="platform-main" style={{ ...mainStyle, width: '100%', minHeight: 'calc(100vh - 108px)', padding: '20px 24px 32px' }}>
          <section className="platform-hero" style={heroStyle}>
            <div>
              <span style={eyebrowStyle}>{copy.shell.eyebrow}</span>
              <h1 style={titleStyle}>{copy.shell.title}</h1>
              <p style={subtitleStyle}>{copy.shell.subtitle}</p>
            </div>
            <div style={statusPanelStyle} aria-label={copy.shell.statusLabel}>
              <span style={statusLabelStyle}>{copy.shell.statusLabel}</span>
              <div style={signalStyle} aria-hidden="true">
                <span style={signalDotStyle} />
              </div>
            </div>
          </section>

          <section className="platform-content-grid" style={contentGridStyle}>
            <article className="platform-section-surface" style={sectionSurfaceStyle}>
              <p style={sectionEyebrowStyle}>{activeSection.eyebrow}</p>
              <h2 style={sectionTitleStyle}>{activeSection.title}</h2>
              <p style={sectionDescriptionStyle}>{activeSection.description}</p>
              <div style={actionRowStyle}>
                <Link href={primaryHref} style={primaryActionStyle}>{activeSection.primaryAction}</Link>
                <span style={secondaryActionStyle}>{activeSection.secondaryAction}</span>
              </div>

              {section === 'news' ? (
                <section style={newsBoardStyle} aria-label={copy.newsBoard.title}>
                  <div className="platform-news-board-header" style={newsBoardHeaderStyle}>
                    <div>
                      <span style={newsBoardBadgeStyle}>{copy.newsBoard.badge}</span>
                      <h3 style={newsBoardTitleStyle}>{copy.newsBoard.title}</h3>
                    </div>
                    <span style={newsBoardAdminOnlyStyle}>{copy.newsBoard.adminOnlyLabel}</span>
                  </div>
                  <p style={newsBoardDescriptionStyle}>{copy.newsBoard.description}</p>
                  <div style={newsBoardEmptyStyle}>
                    <strong style={newsBoardEmptyTitleStyle}>{copy.newsBoard.emptyTitle}</strong>
                    <p style={newsBoardEmptyBodyStyle}>{copy.newsBoard.emptyDescription}</p>
                  </div>
                </section>
              ) : null}
              <ul style={noteListStyle}>
                {activeSection.notes.map((note) => (
                  <li key={note} style={noteItemStyle}>
                    <span style={noteDotStyle} />
                    <span>{note}</span>
                  </li>
                ))}
              </ul>
            </article>

            <aside style={sidePanelStyle}>
              <section style={linksPanelStyle}>
                <h2 style={panelTitleStyle}>{copy.links.title}</h2>
                <p style={panelSubtitleStyle}>{copy.links.subtitle}</p>
                <ul style={linkListStyle}>
                  {copy.links.items.map((item) => (
                    <li key={item.href}>
                      <Link href={item.href} style={linkItemStyle}>
                        <span style={linkLabelStyle}>{item.label}</span>
                        <span style={linkBodyStyle}>{item.body}</span>
                      </Link>
                    </li>
                  ))}
                </ul>
              </section>

              <section style={standardPanelStyle}>
                <h2 style={{ ...panelTitleStyle, color: '#254B5F' }}>{copy.standards.title}</h2>
                <ul style={standardListStyle}>
                  {copy.standards.items.map((item) => (
                    <li key={item.label}>
                      <strong style={standardLabelStyle}>{item.label}</strong>
                      <span style={standardBodyStyle}>{item.body}</span>
                    </li>
                  ))}
                </ul>
              </section>
            </aside>
          </section>
        </main>
      </div>
    </LearnerAppShell>
  );
}
