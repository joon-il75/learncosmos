import type { CSSProperties } from 'react';
import type { PointPageTheme, getPointPageThemeTokens } from '../pointPageUtils';
import { secondaryButtonStyle } from '../pointPageStyles';

type ThemeTokens = ReturnType<typeof getPointPageThemeTokens>;
type ToolbarTone = 'achievement' | 'content' | 'work' | 'evaluation' | 'lumi';

const toolbarToneColors: Record<ToolbarTone, {
  activeLightBackground: string;
  activeDarkBackground: string;
  activeLightBorder: string;
  activeDarkBorder: string;
}> = {
  achievement: {
    activeLightBackground: '#FFF3C4',
    activeDarkBackground: 'rgba(71, 55, 19, 0.92)',
    activeLightBorder: '#C8A63A',
    activeDarkBorder: 'rgba(234, 179, 8, 0.72)',
  },
  content: {
    activeLightBackground: '#E7F0FF',
    activeDarkBackground: 'rgba(30, 64, 175, 0.70)',
    activeLightBorder: '#7BA7E8',
    activeDarkBorder: 'rgba(147, 197, 253, 0.58)',
  },
  work: {
    activeLightBackground: '#EEF7E7',
    activeDarkBackground: 'rgba(22, 101, 52, 0.70)',
    activeLightBorder: '#8DBF72',
    activeDarkBorder: 'rgba(134, 239, 172, 0.56)',
  },
  evaluation: {
    activeLightBackground: '#DDEBFF',
    activeDarkBackground: 'rgba(51, 65, 85, 0.92)',
    activeLightBorder: '#7BA7E8',
    activeDarkBorder: 'rgba(148, 163, 184, 0.72)',
  },
  lumi: {
    activeLightBackground: '#E7F8F2',
    activeDarkBackground: 'rgba(20, 83, 45, 0.82)',
    activeLightBorder: '#72BCA1',
    activeDarkBorder: 'rgba(94, 234, 212, 0.58)',
  },
};

export function getPointToolbarButtonStyle({
  theme,
  themeTokens,
  tone,
  isActive,
  disabled = false,
}: {
  theme: PointPageTheme;
  themeTokens: ThemeTokens;
  tone: ToolbarTone;
  isActive: boolean;
  disabled?: boolean;
}): CSSProperties {
  const toneColors = toolbarToneColors[tone];
  return {
    ...secondaryButtonStyle,
    flex: '0 0 auto',
    whiteSpace: 'nowrap',
    minHeight: '38px',
    borderRadius: '3px',
    padding: '0 12px',
    gap: '7px',
    background: isActive
      ? (theme === 'light' ? toneColors.activeLightBackground : toneColors.activeDarkBackground)
      : (theme === 'light' ? '#F8FAFC' : 'rgba(15, 23, 42, 0.92)'),
    borderColor: isActive
      ? (theme === 'light' ? toneColors.activeLightBorder : toneColors.activeDarkBorder)
      : (theme === 'light' ? '#AAB6C8' : 'rgba(100, 116, 139, 0.72)'),
    color: themeTokens.buttonText,
    fontSize: '13px',
    fontWeight: 800,
    opacity: disabled ? 0.48 : 1,
    cursor: disabled ? 'not-allowed' : 'pointer',
    boxShadow: theme === 'light'
      ? 'inset 1px 1px 0 rgba(255,255,255,0.92), inset -1px -1px 0 rgba(148,163,184,0.28)'
      : 'inset 1px 1px 0 rgba(255,255,255,0.10), inset -1px -1px 0 rgba(0,0,0,0.35)',
  };
}

export function getPointToolbarIconButtonStyle(themeTokens: ThemeTokens): CSSProperties {
  return {
    ...secondaryButtonStyle,
    flex: '0 0 auto',
    width: '34px',
    minWidth: '34px',
    minHeight: '30px',
    padding: 0,
    borderRadius: '999px',
    background: themeTokens.secondaryButtonBackground,
    borderColor: themeTokens.secondaryButtonBorder,
    color: themeTokens.buttonText,
    fontSize: '14px',
    fontWeight: 900,
  };
}
