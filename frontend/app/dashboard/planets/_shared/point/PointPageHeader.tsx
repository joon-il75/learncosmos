'use client';

import AppHeaderShell from '@/components/common/AppHeaderShell';
import LearnerHeaderActions, { type LearnerHeaderActionsCopy } from '@/components/common/LearnerHeaderActions';
import ImmersiveOverlayMenu from '@/components/navigation/ImmersiveOverlayMenu';
import type { Locale } from '@/lib/i18n/locales';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageTheme, getPointPageThemeTokens } from '../pointPageUtils';
import { navStyle, secondaryButtonStyle } from '../pointPageStyles';

type ThemeTokens = ReturnType<typeof getPointPageThemeTokens>;

export function PointPageHeader({
  copy,
  theme,
  themeTokens,
  planetDetailHref,
  onToggleTheme,
  onLogout,
  learnerHeaderCopy,
  locale,
}: {
  copy: PointLearningCopy['header'];
  theme: PointPageTheme;
  themeTokens: ThemeTokens;
  planetDetailHref: string;
  onToggleTheme: () => void;
  onLogout: () => void;
  learnerHeaderCopy: LearnerHeaderActionsCopy;
  locale: Locale;
}) {
  return (
    <AppHeaderShell
      logoHref="/"
      logoIconSize={28}
      logoTextSize="16px"
      maxWidth="1360px"
      headerStyle={{ ...navStyle, position: 'fixed', background: themeTokens.navBackground, borderBottomColor: themeTokens.navBorder }}
      innerStyle={{ padding: '14px 18px' }}
      leftGroupStyle={{ flex: '1 1 auto', minWidth: 0 }}
      rightSlot={
        <span style={{ display: 'inline-flex', alignItems: 'center', gap: '10px' }}>
          <ImmersiveOverlayMenu locale={locale} tone={theme === 'light' ? 'light' : 'dark'} />
          <LearnerHeaderActions
            onLogout={onLogout}
            copy={learnerHeaderCopy}
            galaxyHref={planetDetailHref}
            galaxyLabel={copy.diaryLabel}
            galaxyTitle={copy.diaryTitle}
            extraAction={
              <button
                type="button"
                onClick={onToggleTheme}
                aria-label={theme === 'dark' ? copy.lightMode : copy.darkMode}
                title={theme === 'dark' ? copy.lightMode : copy.darkMode}
                style={{
                  ...secondaryButtonStyle,
                  width: '36px',
                  minWidth: '36px',
                  minHeight: '36px',
                  padding: 0,
                  fontSize: '17px',
                  cursor: 'pointer',
                  background: themeTokens.secondaryButtonBackground,
                  borderColor: themeTokens.secondaryButtonBorder,
                  color: themeTokens.buttonText,
                }}
              >
                {theme === 'dark' ? '☀️' : '🌙'}
              </button>
            }
            compactExtraAction={
              <button
                type="button"
                onClick={onToggleTheme}
                style={{
                  ...secondaryButtonStyle,
                  width: '100%',
                  justifyContent: 'flex-start',
                  minHeight: '40px',
                  padding: '0 12px',
                  borderRadius: '12px',
                  fontSize: '14px',
                  cursor: 'pointer',
                  background: 'transparent',
                  borderColor: 'transparent',
                  color: '#E8EEFF',
                }}
              >
                {theme === 'dark' ? copy.lightMode : copy.darkMode}
              </button>
            }
            palette={theme === 'light' ? {
              text: '#0F172A',
              galaxyText: '#064E3B',
              galaxyBackground: '#DFF7EA',
              galaxyBorder: '#86D7A8',
              galaxyShadow: '0 8px 22px rgba(16, 185, 129, 0.18)',
              userMenu: { triggerText: '#0F172A', pointText: '#0F172A', pointBackground: '#F8FAFC', pointBorder: '#CBD5E1' },
            } : {
              text: '#F4F7FF',
              galaxyText: '#E9FFF3',
              galaxyBackground: 'rgba(37, 170, 118, 0.38)',
              galaxyBorder: 'rgba(134, 239, 172, 0.52)',
              galaxyShadow: '0 10px 26px rgba(16, 185, 129, 0.20)',
            }}
          />
        </span>
      }
    />
  );
}
