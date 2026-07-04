'use client';

import Link from 'next/link';
import type { CSSProperties, ReactNode } from 'react';
import LearnerHeaderActions, { getLearnerHeaderActionsCopy } from '@/components/common/LearnerHeaderActions';
import LearnerAppShell from '@/components/navigation/LearnerAppShell';

interface DashboardSettingsShellProps {
  eyebrow: string;
  title: string;
  copy: string;
  tabLabels?: Partial<Record<'profile' | 'points' | 'ai', string>>;
  activeTab: 'profile' | 'points' | 'ai';
  uiLocale?: 'ko' | 'en';
  onLogout: () => void;
  children: ReactNode;
}

const tabs = [
  { key: 'profile', label: '학습자정보', href: '/dashboard/settings/profile' },
  { key: 'points', label: '포인트', href: '/dashboard/settings/points' },
  { key: 'ai', label: 'AI 설정', href: '/dashboard/settings/ai' },
] as const;

export default function DashboardSettingsShell({
  eyebrow,
  title,
  copy,
  tabLabels,
  activeTab,
  uiLocale = 'ko',
  onLogout,
  children,
}: DashboardSettingsShellProps) {
  return (
    <div style={pageStyle}>
      <div style={pageBackgroundOverlayStyle} />

      <LearnerAppShell
        locale={uiLocale}
        logoHref="/"
        rightSlot={
          <LearnerHeaderActions onLogout={onLogout} copy={getLearnerHeaderActionsCopy(uiLocale)} />
        }
        contentClassName="!px-0 !pb-0"
      >
      <main style={mainStyle}>
        <section style={heroCardStyle}>
          <div style={eyebrowStyle}>{eyebrow}</div>
          <h1 style={heroTitleStyle}>{title}</h1>
          <p style={heroCopyStyle}>{copy}</p>
          <div style={tabsRowStyle}>
            {tabs.map((tab, index) => (
              <Link key={tab.key} href={tab.href} style={tabStyle(tab.key === activeTab, index, tabs.length)}>
                {tabLabels?.[tab.key] ?? tab.label}
              </Link>
            ))}
          </div>
        </section>

        {children}
      </main>
      </LearnerAppShell>
    </div>
  );
}

export const pageStyle = {
  position: 'relative',
  minHeight: '100vh',
  overflow: 'hidden',
  background: 'linear-gradient(145deg, #0D314E 0%, #1C7D79 50%, #CC5216 100%)',
} satisfies CSSProperties;

const pageBackgroundOverlayStyle = {
  position: 'absolute',
  inset: 0,
  background: 'radial-gradient(circle at 50% 0%, rgba(255, 255, 255, 0.18), transparent 34%), linear-gradient(180deg, rgba(4, 15, 28, 0.08), rgba(4, 15, 28, 0.18))',
  pointerEvents: 'none',
} satisfies CSSProperties;

export const loadingPageStyle = {
  minHeight: '100vh',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  background: 'linear-gradient(145deg, #0D314E 0%, #1C7D79 50%, #CC5216 100%)',
  padding: '24px',
} satisfies CSSProperties;

export const loadingCardStyle = {
  width: '100%',
  maxWidth: '420px',
  padding: '32px',
  borderRadius: '18px',
  background: 'rgba(255, 253, 247, 0.96)',
  border: '1px solid rgba(255, 255, 255, 0.46)',
  textAlign: 'center',
  color: '#1D140B',
  fontSize: '17px',
  fontWeight: 700,
  boxShadow: '0 24px 60px rgba(3, 8, 20, 0.28)',
} satisfies CSSProperties;

const mainStyle = {
  position: 'relative',
  maxWidth: '980px',
  margin: '0 auto',
  padding: '28px 24px 40px',
  display: 'grid',
  gap: '18px',
} satisfies CSSProperties;

const heroCardStyle = {
  display: 'grid',
  gap: '12px',
  borderRadius: '20px',
  border: '1px solid rgba(255, 255, 255, 0.42)',
  borderLeft: '4px solid rgba(133, 183, 235, 0.78)',
  background: 'rgba(255, 253, 247, 0.96)',
  padding: '24px',
  boxShadow: '0 24px 60px rgba(3, 8, 20, 0.24)',
} satisfies CSSProperties;

const eyebrowStyle = {
  fontSize: '14px',
  fontWeight: 800,
  color: '#533819',
  letterSpacing: '0.08em',
  textTransform: 'uppercase',
} satisfies CSSProperties;

const heroTitleStyle = {
  margin: 0,
  fontSize: '38px',
  lineHeight: 1.18,
  color: '#160E08',
  fontWeight: 800,
} satisfies CSSProperties;

const heroCopyStyle = {
  margin: 0,
  fontSize: '15px',
  lineHeight: 1.8,
  color: '#2E241A',
  maxWidth: '760px',
  fontWeight: 500,
} satisfies CSSProperties;

const tabsRowStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))',
  alignItems: 'center',
  marginTop: '8px',
  gap: '8px',
  padding: '8px',
  borderRadius: '18px',
  border: '1px solid rgba(13, 49, 78, 0.16)',
  background: 'rgba(11, 22, 41, 0.08)',
} satisfies CSSProperties;

const tabStyle = (active: boolean, index: number, total: number): CSSProperties => ({
  minHeight: '52px',
  padding: '0 18px',
  borderRadius: '14px',
  border: active ? '2px solid rgba(13, 49, 78, 0.72)' : '1px solid rgba(13, 49, 78, 0.18)',
  marginLeft: 0,
  background: active ? 'linear-gradient(135deg, rgba(13, 49, 78, 0.96), rgba(28, 125, 121, 0.92))' : 'rgba(255, 255, 255, 0.82)',
  color: active ? '#FFFDF7' : '#1F2933',
  fontSize: '16px',
  fontWeight: active ? 900 : 800,
  textDecoration: 'none',
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  boxShadow: active ? '0 12px 26px rgba(13, 49, 78, 0.24)' : '0 6px 14px rgba(13, 49, 78, 0.08)',
  outline: active ? '2px solid rgba(255, 255, 255, 0.72)' : 'none',
  outlineOffset: '-5px',
});

export const sectionCardStyle = {
  display: 'grid',
  gap: '16px',
  borderRadius: '20px',
  border: '1px solid rgba(255, 255, 255, 0.42)',
  borderLeft: '4px solid rgba(133, 183, 235, 0.78)',
  background: 'rgba(255, 253, 247, 0.96)',
  padding: '22px',
  boxShadow: '0 24px 60px rgba(3, 8, 20, 0.22)',
} satisfies CSSProperties;

export const sectionHeaderStyle = {
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: '12px',
  flexWrap: 'wrap',
} satisfies CSSProperties;

export const sectionTitleStyle = {
  margin: 0,
  fontSize: '28px',
  lineHeight: 1.24,
  color: '#160E08',
  fontWeight: 800,
} satisfies CSSProperties;

export const sectionSubtitleStyle = {
  margin: '6px 0 0',
  fontSize: '15px',
  lineHeight: 1.7,
  color: '#35291E',
  fontWeight: 500,
} satisfies CSSProperties;

export const pillStyle = (accent: string): CSSProperties => ({
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: '34px',
  padding: '0 12px',
  borderRadius: '999px',
  background: `${accent}22`,
  color: accent,
  border: `1px solid ${accent}33`,
  fontSize: '13px',
  fontWeight: 700,
});

export const summaryRowStyle = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '12px',
  flexWrap: 'wrap',
} satisfies CSSProperties;

export const summaryTextStyle = {
  fontSize: '14px',
  color: '#4C4033',
} satisfies CSSProperties;

export const profileMetaGridStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))',
  gap: '12px',
} satisfies CSSProperties;

export const profileMetaCardStyle = {
  display: 'grid',
  gap: '8px',
  padding: '16px',
  borderRadius: '14px',
  border: '1px solid rgba(146, 111, 62, 0.16)',
  background: 'rgba(255, 248, 235, 0.62)',
} satisfies CSSProperties;

export const profileMetaLabelStyle = {
  fontSize: '12px',
  color: '#7A6A55',
  textTransform: 'uppercase',
  letterSpacing: '0.06em',
} satisfies CSSProperties;

export const profileMetaValueStyle = {
  fontSize: '15px',
  fontWeight: 600,
  color: '#2E2112',
} satisfies CSSProperties;

export const fieldGridStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
  gap: '14px',
} satisfies CSSProperties;

export const fieldLabelStyle = {
  display: 'grid',
  gap: '8px',
  fontSize: '14px',
  fontWeight: 800,
  color: '#21160E',
} satisfies CSSProperties;

export const settingsInputStyle = {
  minHeight: '48px',
  borderRadius: '14px',
  border: '1px solid rgba(13, 49, 78, 0.22)',
  background: 'rgba(255, 255, 255, 0.92)',
  color: '#160E08',
  fontSize: '15px',
  fontWeight: 700,
  padding: '0 14px',
  outline: 'none',
} satisfies CSSProperties;

export const profilePointsRowStyle = {
  display: 'flex',
  alignItems: 'center',
  gap: '10px',
  flexWrap: 'wrap',
  alignSelf: 'end',
} satisfies CSSProperties;

export const miniPointChipStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: '32px',
  padding: '0 12px',
  borderRadius: '999px',
  border: '1px solid rgba(146, 111, 62, 0.18)',
  background: 'rgba(255, 248, 235, 0.62)',
  color: '#4C4033',
  fontSize: '13px',
  fontWeight: 600,
} satisfies CSSProperties;

export const inlineInfoStyle = (accent: string): CSSProperties => ({
  padding: '12px 14px',
  borderRadius: '14px',
  border: `1px solid ${accent}33`,
  background: `${accent}15`,
  color: accent === '#FFB4A2' ? '#9B2C2C' : accent === '#AFC0E4' ? '#4C4033' : accent,
  fontSize: '13px',
  lineHeight: 1.6,
});

export const actionRowStyle = {
  display: 'flex',
  alignItems: 'center',
  gap: '10px',
  flexWrap: 'wrap',
} satisfies CSSProperties;

export const primaryActionButtonStyle = (disabled: boolean): CSSProperties => ({
  minHeight: '42px',
  padding: '0 18px',
  borderRadius: '999px',
  border: 'none',
  background: disabled ? 'rgba(146, 111, 62, 0.18)' : 'linear-gradient(135deg, #45D483 0%, #2CB67D 100%)',
  color: '#06131E',
  fontWeight: 800,
  fontSize: '14px',
  cursor: disabled ? 'default' : 'pointer',
  fontFamily: 'inherit',
});

export const secondaryActionButtonStyle = (disabled: boolean): CSSProperties => ({
  minHeight: '42px',
  padding: '0 18px',
  borderRadius: '999px',
  border: '1px solid rgba(146, 111, 62, 0.22)',
  background: disabled ? 'rgba(146, 111, 62, 0.12)' : 'rgba(255, 253, 247, 0.76)',
  color: '#3F3428',
  fontWeight: 700,
  fontSize: '14px',
  cursor: disabled ? 'default' : 'pointer',
  fontFamily: 'inherit',
});

export const helperTextStyle = {
  fontSize: '14px',
  lineHeight: 1.7,
  color: '#35291E',
  fontWeight: 500,
} satisfies CSSProperties;
