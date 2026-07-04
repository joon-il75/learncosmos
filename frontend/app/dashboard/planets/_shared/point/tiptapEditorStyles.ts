import type { CSSProperties } from 'react';

import type { PointPageThemeTokens } from '../pointPageUtils';

export function createToolbarButtonStyle(
  themeTokens: PointPageThemeTokens,
  editable: boolean,
  active = false,
): CSSProperties {
  return {
    width: '28px',
    minWidth: '28px',
    height: '28px',
    display: 'inline-flex',
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: '7px',
    border: `1px solid ${active ? themeTokens.primaryButtonBorder : themeTokens.secondaryButtonBorder}`,
    background: active ? themeTokens.primaryButtonBackground : themeTokens.secondaryButtonBackground,
    color: themeTokens.buttonText,
    padding: 0,
    cursor: editable ? 'pointer' : 'default',
    opacity: editable ? 1 : 0.62,
  };
}

export function createToolbarGroupStyle(
  themeTokens: PointPageThemeTokens,
  withDivider = true,
): CSSProperties {
  return {
    display: 'inline-flex',
    alignItems: 'center',
    gap: '4px',
    paddingRight: withDivider ? '7px' : 0,
    marginRight: withDivider ? '1px' : 0,
    borderRight: withDivider ? `1px solid ${themeTokens.surfaceBorder}` : '0',
  };
}

export function createToolbarSelectStyle(
  themeTokens: PointPageThemeTokens,
  editable: boolean,
): CSSProperties {
  return {
    height: '28px',
    minWidth: '72px',
    borderRadius: '7px',
    border: `1px solid ${themeTokens.secondaryButtonBorder}`,
    background: themeTokens.secondaryButtonBackground,
    color: themeTokens.buttonText,
    padding: '0 24px 0 8px',
    cursor: editable ? 'pointer' : 'default',
    fontSize: '12px',
    fontWeight: 800,
    outline: 'none',
    opacity: editable ? 1 : 0.62,
  };
}

export function createToolbarTextButtonStyle(
  themeTokens: PointPageThemeTokens,
  editable: boolean,
  active = false,
): CSSProperties {
  return {
    ...createToolbarButtonStyle(themeTokens, editable, active),
    width: 'auto',
    padding: '0 12px',
    fontWeight: 800,
  };
}
