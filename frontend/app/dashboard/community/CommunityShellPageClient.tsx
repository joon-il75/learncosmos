'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

import LearnerHeaderActions, { getLearnerHeaderActionsCopy } from '@/components/common/LearnerHeaderActions';
import LearnerAppShell from '@/components/navigation/LearnerAppShell';
import { normalizeLocale, type Locale } from '@/lib/i18n/locales';
import { getCommunityShellCopy, type CommunityShellSection } from '@/lib/i18n/pages/community';
import {
  actionRowStyle,
  contentGridStyle,
  eyebrowPillStyle,
  heroStyle,
  heroTextStyle,
  mainStyle,
  mobileStyle,
  noteDotStyle,
  noteItemStyle,
  noteListStyle,
  orbitDotStyle,
  orbitStyle,
  pageStyle,
  primaryActionStyle,
  roadmapBodyStyle,
  roadmapIndexStyle,
  roadmapItemStyle,
  roadmapLabelStyle,
  roadmapListStyle,
  roadmapStyle,
  roadmapTitleStyle,
  secondaryActionStyle,
  sectionDescriptionStyle,
  sectionEyebrowStyle,
  sectionSurfaceStyle,
  sectionTitleStyle,
  statusLabelStyle,
  statusPanelStyle,
  subtitleStyle,
  titleStyle,
} from './communityShellStyles';

interface CommunityShellPageClientProps {
  section: CommunityShellSection;
}

const primaryHrefs: Record<CommunityShellSection, string> = {
  overview: '/dashboard',
  discover: '/dashboard',
  mine: '/dashboard',
};

export default function CommunityShellPageClient({ section }: CommunityShellPageClientProps) {
  const router = useRouter();
  const [uiLocale, setUiLocale] = useState<Locale>('ko');
  const copy = getCommunityShellCopy(uiLocale);
  const activeSection = copy.sections[section];

  useEffect(() => {
    let cancelled = false;
    const loadLocale = async () => {
      try {
        await fetch('/api/v1/auth/refresh', { method: 'POST', credentials: 'include' });
        const res = await fetch('/api/v1/auth/me', { credentials: 'include', cache: 'no-store' });
        if (!res.ok) return;
        const data = (await res.json()) as { ui_locale?: string };
        if (!cancelled) setUiLocale(normalizeLocale(data.ui_locale));
      } catch {
        // Keep Korean fallback.
      }
    };
    void loadLocale();
    return () => {
      cancelled = true;
    };
  }, []);

  const handleLogout = async () => {
    await fetch('/api/v1/auth/logout', { method: 'POST', credentials: 'include' }).catch(() => null);
    router.replace('/login');
  };

  const rightSlot = (
    <LearnerHeaderActions
      onLogout={handleLogout}
      galaxyLabel={copy.shell.galaxyLabel}
      galaxyTitle={copy.shell.galaxyTitle}
      copy={getLearnerHeaderActionsCopy(uiLocale)}
    />
  );

  return (
    <LearnerAppShell locale={uiLocale} logoHref="/" rightSlot={rightSlot} contentClassName="!px-0 !pb-0">
      <div style={pageStyle}>
        <style>{mobileStyle}</style>
        <main className="community-main" style={{ ...mainStyle, width: '100%', minHeight: 'calc(100vh - 108px)', padding: '20px 24px 32px' }}>
          <section className="community-hero" style={heroStyle}>
            <div style={heroTextStyle}>
              <span style={eyebrowPillStyle}>{copy.shell.eyebrow}</span>
              <h1 style={titleStyle}>{copy.shell.title}</h1>
              <p style={subtitleStyle}>{copy.shell.subtitle}</p>
            </div>
            <div style={statusPanelStyle} aria-label={copy.shell.statusLabel}>
              <span style={statusLabelStyle}>{copy.shell.statusLabel}</span>
              <div style={orbitStyle} aria-hidden="true">
                <span style={orbitDotStyle} />
              </div>
            </div>
          </section>

          <section className="community-content-grid" style={contentGridStyle}>
            <article className="community-section-surface" style={sectionSurfaceStyle}>
              <p style={sectionEyebrowStyle}>{activeSection.eyebrow}</p>
              <h2 style={sectionTitleStyle}>{activeSection.title}</h2>
              <p style={sectionDescriptionStyle}>{activeSection.description}</p>
              <div style={actionRowStyle}>
                <Link href={primaryHrefs[section]} style={primaryActionStyle}>{activeSection.primaryAction}</Link>
                <span style={secondaryActionStyle}>{activeSection.secondaryAction}</span>
              </div>
              <ul style={noteListStyle}>
                {activeSection.notes.map((note) => (
                  <li key={note} style={noteItemStyle}>
                    <span style={noteDotStyle} />
                    <span>{note}</span>
                  </li>
                ))}
              </ul>
            </article>

            <aside style={roadmapStyle}>
              <h2 style={roadmapTitleStyle}>{copy.roadmap.title}</h2>
              <ol style={roadmapListStyle}>
                {copy.roadmap.items.map((item, index) => (
                  <li key={item.label} style={roadmapItemStyle}>
                    <span style={roadmapIndexStyle}>{index + 1}</span>
                    <span>
                      <strong style={roadmapLabelStyle}>{item.label}</strong>
                      <span style={roadmapBodyStyle}>{item.body}</span>
                    </span>
                  </li>
                ))}
              </ol>
            </aside>
          </section>
        </main>
      </div>
    </LearnerAppShell>
  );
}
