import type { PointPageThemeTokens } from '../../pointPageUtils';
import {
  secondaryButtonStyle,
  blockMetaStyle,
} from '../../pointPageStyles';

export function makeWorkspaceStyles(themeTokens: PointPageThemeTokens) {
  const materialToolbarButtonStyle = {
    ...secondaryButtonStyle,
    minHeight: '38px',
    borderRadius: '3px',
    padding: '0 12px',
    gap: '7px',
    background: themeTokens.secondaryButtonBackground,
    borderColor: themeTokens.secondaryButtonBorder,
    color: themeTokens.buttonText,
    fontSize: '13px',
    fontWeight: 800,
    whiteSpace: 'nowrap',
  } as const;

  const materialToolbarPrimaryButtonStyle = {
    ...materialToolbarButtonStyle,
    background: themeTokens.pageBackground === '#F8FAFC' ? '#E7F0FF' : 'rgba(30, 64, 175, 0.70)',
    borderColor: themeTokens.pageBackground === '#F8FAFC' ? '#7BA7E8' : 'rgba(147, 197, 253, 0.58)',
    color: themeTokens.buttonText,
  } as const;

  const materialToolbarDangerButtonStyle = {
    ...materialToolbarButtonStyle,
    background: themeTokens.dangerButtonBackground,
    borderColor: themeTokens.dangerButtonBorder,
    color: themeTokens.dangerButtonText,
  } as const;

  const listBoardStyle = {
    display: 'grid',
    overflow: 'hidden',
    border: `1px solid ${themeTokens.surfaceBorder}`,
    borderRadius: '4px',
    background: themeTokens.surfaceBackground,
  } as const;

  const listHoverBackground = themeTokens.pageBackground === '#F8FAFC' ? '#EFF6FF' : 'rgba(96, 165, 250, 0.13)';

  const getListHeaderRowStyle = (columns: string) => ({
    display: 'grid',
    gridTemplateColumns: columns,
    alignItems: 'center',
    minHeight: '38px',
    padding: '0 16px',
    borderBottom: `1px solid ${themeTokens.surfaceBorder}`,
    background: themeTokens.pageBackground === '#F8FAFC' ? '#F1F5F9' : 'rgba(148, 163, 184, 0.10)',
  } as const);

  const getListDataRowStyle = (columns: string, hovered: boolean) => ({
    display: 'grid',
    gridTemplateColumns: columns,
    alignItems: 'center',
    width: '100%',
    minHeight: '46px',
    padding: '9px 16px',
    border: 'none',
    borderLeft: `3px solid ${hovered ? themeTokens.primaryButtonBorder : 'transparent'}`,
    borderBottom: `1px solid ${themeTokens.surfaceBorder}`,
    background: hovered ? listHoverBackground : themeTokens.surfaceBackground,
    color: themeTokens.title,
    cursor: 'pointer',
    textAlign: 'left',
    font: 'inherit',
    transition: 'background 140ms ease, border-color 140ms ease',
  } as const);

  const listColumnHeaderStyle = {
    ...blockMetaStyle,
    color: themeTokens.metaLabel,
    fontWeight: 850,
  } as const;

  const getListIndexCellStyle = (hovered: boolean) => ({
    ...blockMetaStyle,
    color: hovered ? themeTokens.title : themeTokens.metaLabel,
    fontWeight: 850,
  } as const);

  const getListTextCellStyle = (hovered: boolean) => ({
    color: hovered ? themeTokens.title : themeTokens.description,
    fontSize: '14px',
    lineHeight: 1.6,
    fontWeight: hovered ? 750 : 600,
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  } as const);

  const readonlyRecordFieldStyle = {
    display: 'grid',
    gap: '6px',
    padding: '12px 0',
    borderBottom: `1px solid ${themeTokens.surfaceBorder}`,
  } as const;

  const readonlyRecordValueStyle = {
    margin: 0,
    color: themeTokens.description,
    fontSize: '15px',
    lineHeight: 1.65,
    whiteSpace: 'pre-wrap',
  } as const;

  return {
    materialToolbarButtonStyle,
    materialToolbarPrimaryButtonStyle,
    materialToolbarDangerButtonStyle,
    listBoardStyle,
    listHoverBackground,
    getListHeaderRowStyle,
    getListDataRowStyle,
    listColumnHeaderStyle,
    getListIndexCellStyle,
    getListTextCellStyle,
    readonlyRecordFieldStyle,
    readonlyRecordValueStyle,
  };
}

export type WorkspaceStyles = ReturnType<typeof makeWorkspaceStyles>;
