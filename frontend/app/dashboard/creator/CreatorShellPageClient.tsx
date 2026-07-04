'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

import LearnerHeaderActions, { getLearnerHeaderActionsCopy } from '@/components/common/LearnerHeaderActions';
import LearnerAppShell from '@/components/navigation/LearnerAppShell';
import { normalizeLocale, type Locale } from '@/lib/i18n/locales';
import { getCreatorShellCopy, type CreatorShellSection } from '@/lib/i18n/pages/creator';
import {
  actionRowStyle,
  canvasStyle,
  contentGridStyle,
  eyebrowStyle,
  heroStyle,
  heroTextStyle,
  mainStyle,
  mobileStyle,
  navLinkStyle,
  navStyle,
  noteDotStyle,
  noteItemStyle,
  noteListStyle,
  pageStyle,
  primaryActionStyle,
  readinessBodyStyle,
  readinessItemStyle,
  readinessLabelStyle,
  readinessListStyle,
  readinessStyle,
  readinessSubtitleStyle,
  readinessTitleStyle,
  secondaryActionStyle,
  sectionDescriptionStyle,
  sectionEyebrowStyle,
  sectionSurfaceStyle,
  sectionTitleStyle,
  statusStyle,
  studioDotStyle,
  studioLinesStyle,
  subtitleStyle,
  titleStyle,
} from './creatorShellStyles';

interface CreatorShellPageClientProps {
  section: CreatorShellSection;
}

const navHrefs: Record<CreatorShellSection, string> = {
  overview: '/dashboard/creator',
  content: '/dashboard/creator/content',
  groupCourses: '/dashboard/creator/group-courses',
  feedback: '/dashboard/creator/feedback',
};

const primaryHrefs: Record<CreatorShellSection, string> = {
  overview: '/dashboard',
  content: '/dashboard',
  groupCourses: '/dashboard/community',
  feedback: '/dashboard',
};

export default function CreatorShellPageClient({ section }: CreatorShellPageClientProps) {
  const router = useRouter();
  const [uiLocale, setUiLocale] = useState<Locale>('ko');
  const copy = getCreatorShellCopy(uiLocale);
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
    <div style={pageStyle}>
      <style>{mobileStyle}</style>
      <LearnerAppShell locale={uiLocale} logoHref="/" rightSlot={rightSlot} contentClassName="!px-0 !pb-0">
        <main
          className="creator-main"
          style={{
            ...mainStyle,
            width: '100%',
            minHeight: 'calc(100vh - 108px)',
            margin: 0,
            padding: '20px 24px 32px',
          }}
        >
          <section className="creator-hero" style={heroStyle}>
            <div style={heroTextStyle}>
              <span style={eyebrowStyle}>{copy.shell.eyebrow}</span>
              <h1 style={titleStyle}>{copy.shell.title}</h1>
              <p style={subtitleStyle}>{copy.shell.subtitle}</p>
            </div>
            <div style={canvasStyle} aria-label={copy.shell.statusLabel}>
              <span style={statusStyle}>{copy.shell.statusLabel}</span>
              <div style={studioLinesStyle} aria-hidden="true">
                <span style={studioDotStyle} />
              </div>
            </div>
          </section>

          <nav style={navStyle} aria-label="Creator sections">
            {(Object.keys(navHrefs) as CreatorShellSection[]).map((item) => (
              <Link key={item} href={navHrefs[item]} style={navLinkStyle(item === section)}>
                {copy.nav[item]}
              </Link>
            ))}
          </nav>

          <section className="creator-content-grid" style={contentGridStyle}>
            <article className="creator-section-surface" style={sectionSurfaceStyle}>
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

            <aside style={readinessStyle}>
              <h2 style={readinessTitleStyle}>{copy.readiness.title}</h2>
              <p style={readinessSubtitleStyle}>{copy.readiness.subtitle}</p>
              <ul style={readinessListStyle}>
                {copy.readiness.items.map((item) => (
                  <li key={item.label} style={readinessItemStyle}>
                    <strong style={readinessLabelStyle}>{item.label}</strong>
                    <span style={readinessBodyStyle}>{item.body}</span>
                  </li>
                ))}
              </ul>
            </aside>
          </section>
        </main>
      </LearnerAppShell>
    </div>
  );
}
