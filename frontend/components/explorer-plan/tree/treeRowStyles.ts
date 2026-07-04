import type { CSSProperties } from 'react'

// ── Course Header Row ─────────────────────────────────────────────────────────

export const courseHeaderRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '6px',
  padding: '8px 10px 8px 8px',
  flexShrink: 0,
  cursor: 'pointer',
  borderBottom: '1px solid rgba(100, 72, 30, 0.25)',
  userSelect: 'none',
  background: 'transparent',
  transition: 'background 120ms ease, filter 120ms ease',
}

export const courseHeaderSelectedStyle: CSSProperties = {
  background: 'rgba(59, 130, 246, 0.08)',
  borderLeft: '3px solid #2563EB',
  paddingLeft: 7,
}

export const courseHeaderTitleStyle: CSSProperties = {
  flex: 1,
  minWidth: 0,
  fontSize: '13px',
  fontWeight: 700,
  color: '#1A0E03',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
}

export const courseHeaderDisplayTitleStyle: CSSProperties = {
  flex: 1,
  minWidth: 0,
  display: 'flex',
  alignItems: 'baseline',
  gap: '4px',
  overflow: 'hidden',
  whiteSpace: 'nowrap',
}

export const courseHeaderPrefixStyle: CSSProperties = {
  flexShrink: 0,
  fontSize: '12px',
  fontWeight: 700,
  color: '#3D2A14',
}

export const courseHeaderPlanetNameStyle: CSSProperties = {
  minWidth: 0,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
  fontSize: '14px',
  fontWeight: 800,
  color: '#1E3A6D',
}

// ── Tree Row States ───────────────────────────────────────────────────────────

export const treeRowBaseStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: '4px',
  padding: '5px 10px 5px 6px',
  userSelect: 'none',
  cursor: 'pointer',
  transition: 'background 120ms ease, filter 120ms ease',
}

export const treeRowHoverStyle: CSSProperties = {
  background: 'rgba(90, 60, 20, 0.09)',
  filter: 'brightness(1.03)',
}

export const treeRowSelectedStyle: CSSProperties = {
  background: 'rgba(37, 99, 235, 0.10)',
  borderLeft: '3px solid #2563EB',
}

export const treeRowLearningStyle: CSSProperties = {
  background: 'rgba(37, 99, 235, 0.13)',
  boxShadow: 'inset 0 0 0 1px rgba(37, 99, 235, 0.22)',
}

export const treeRowInactiveStyle: CSSProperties = {
  opacity: 0.68,
  filter: 'grayscale(0.12)',
}

export const treeRowPendingCreateStyle: CSSProperties = {
  opacity: 0.72,
  background: 'rgba(37, 99, 235, 0.07)',
  outline: '1px dashed rgba(37, 99, 235, 0.34)',
  outlineOffset: '-2px',
}

export const treeRowPendingDeleteStyle: CSSProperties = {
  opacity: 0.74,
  background: 'rgba(185, 58, 42, 0.08)',
  outline: '1px solid rgba(185, 58, 42, 0.24)',
  outlineOffset: '-1px',
}

export const treeRowPendingMoveStyle: CSSProperties = {
  background: 'rgba(37, 99, 235, 0.08)',
  outline: '1px solid rgba(37, 99, 235, 0.18)',
  outlineOffset: '-1px',
}

// ── Pending Badges ────────────────────────────────────────────────────────────

export const pendingCreateBadgeStyle: CSSProperties = {
  flexShrink: 0,
  padding: '1px 4px',
  borderRadius: '999px',
  border: '1px dashed rgba(37, 99, 235, 0.42)',
  background: 'rgba(37, 99, 235, 0.08)',
  color: '#1D4ED8',
  fontSize: '9px',
  fontWeight: 800,
  lineHeight: 1.25,
}

export const pendingDeleteBadgeStyle: CSSProperties = {
  flexShrink: 0,
  padding: '1px 4px',
  borderRadius: '999px',
  border: '1px solid rgba(185, 58, 42, 0.34)',
  background: 'rgba(185, 58, 42, 0.10)',
  color: '#A43124',
  fontSize: '9px',
  fontWeight: 800,
  lineHeight: 1.25,
}

export const pendingMoveIconStyle: CSSProperties = {
  flexShrink: 0,
  fontSize: '12px',
  lineHeight: 1,
  color: '#1D4ED8',
}

// ── Row Controls ──────────────────────────────────────────────────────────────

export const treeToggleBtnStyle: CSSProperties = {
  flexShrink: 0,
  width: '20px',
  height: '20px',
  border: 'none',
  background: 'transparent',
  cursor: 'pointer',
  fontSize: '13px',
  padding: 0,
  lineHeight: 1,
}

export const treeRowLabelStyle: CSSProperties = {
  flex: 1,
  minWidth: 0,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
  lineHeight: 1.45,
  color: '#1A0E03',
}
