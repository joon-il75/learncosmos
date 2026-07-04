import type { CSSProperties } from 'react'

export const panelStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: '8px',
  height: '100%',
  maxHeight: '100%',
  minHeight: 0,
  width: '100%',
  overflow: 'hidden',
  background: 'transparent',
}

export const compactPanelStyle: CSSProperties = {
  borderRadius: '12px',
}

export const headerStyle: CSSProperties = {
  display: 'grid',
  gap: '9px',
  padding: '12px 12px 0',
}

export const headerTopRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '10px',
}

export const eyebrowStyle: CSSProperties = {
  fontSize: '10px',
  fontWeight: 700,
  letterSpacing: '0.06em',
  textTransform: 'uppercase',
  color: '#7A5A24',
}

export const titleStyle: CSSProperties = {
  margin: '2px 0 0',
  fontSize: '15px',
  fontWeight: 800,
  lineHeight: 1.2,
  color: '#1F160A',
}

export const countBadgeStyle: CSSProperties = {
  flexShrink: 0,
  padding: '4px 9px',
  borderRadius: '999px',
  background: 'rgba(122, 90, 36, 0.08)',
  color: '#7A5A24',
  fontSize: '10px',
  fontWeight: 700,
}

export const createRegionFormStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '1fr auto',
  gap: '6px',
}

export const treeInputStyle: CSSProperties = {
  minWidth: 0,
  width: '100%',
  height: 30,
  border: '1px solid rgba(122, 90, 36, 0.22)',
  borderRadius: '8px',
  background: 'rgba(255, 252, 244, 0.78)',
  color: '#1F160A',
  fontSize: '12px',
  fontWeight: 700,
  padding: '0 9px',
  outline: 'none',
  boxSizing: 'border-box',
}

export const treeInlineInputStyle: CSSProperties = {
  ...treeInputStyle,
  height: 26,
  background: 'rgba(255, 252, 244, 0.96)',
  border: '1px solid rgba(59, 130, 246, 0.26)',
  padding: '0 7px',
}

export const treePrimaryButtonStyle: CSSProperties = {
  height: 32,
  border: '1px solid rgba(122, 90, 36, 0.28)',
  borderRadius: '8px',
  background: 'rgba(122, 90, 36, 0.14)',
  color: '#5F4319',
  fontSize: '12px',
  fontWeight: 800,
  cursor: 'pointer',
  padding: '0 12px',
  whiteSpace: 'nowrap',
}

export const treeButtonDisabledStyle: CSSProperties = {
  opacity: 0.45,
  cursor: 'not-allowed',
}

export const mutationMessageStyle: CSSProperties = {
  minHeight: 14,
  fontSize: '10px',
  fontWeight: 700,
  color: 'rgba(90, 61, 31, 0.72)',
}

export const scrollAreaStyle: CSSProperties = {
  flex: '1 1 auto',
  height: 0,
  overflowY: 'auto',
  minHeight: 0,
  padding: '0 10px 12px',
  display: 'flex',
  flexDirection: 'column',
  gap: '8px',
  overscrollBehavior: 'contain',
}

export const emptyStateStyle: CSSProperties = {
  padding: '18px',
  color: 'rgba(56, 40, 18, 0.70)',
  fontSize: '13px',
}
