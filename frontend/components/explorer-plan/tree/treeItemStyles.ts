import type { CSSProperties } from 'react'

// ── Region ────────────────────────────────────────────────────────────────────

export const regionCardStyle: CSSProperties = {
  flexShrink: 0,
  borderRadius: '8px',
  border: '1px solid rgba(150, 120, 72, 0.05)',
  background: 'transparent',
  overflow: 'hidden',
}

export const regionHeaderRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'flex-start',
  gap: '6px',
  padding: '7px 7px 5px 7px',
  background: 'rgba(196, 157, 91, 0.05)',
}

export const regionActionColumnStyle: CSSProperties = {
  flexShrink: 0,
  display: 'grid',
  gap: '4px',
  marginTop: '1px',
}

export const regionActionButtonStyle: CSSProperties = {
  minWidth: 34,
  height: 20,
  border: '1px solid rgba(122, 90, 36, 0.18)',
  borderRadius: '7px',
  background: 'rgba(255, 252, 244, 0.62)',
  color: '#6D4B18',
  fontSize: '10px',
  fontWeight: 800,
  cursor: 'pointer',
  padding: '0 6px',
}

export const regionDangerButtonStyle: CSSProperties = {
  color: '#9B2C2C',
  borderColor: 'rgba(155, 44, 44, 0.18)',
}

export const regionEditFormStyle: CSSProperties = {
  flex: 1,
  display: 'grid',
  gap: '5px',
  minWidth: 0,
}

export const regionEditActionRowStyle: CSSProperties = {
  display: 'flex',
  gap: '5px',
}

export const collapseButtonStyle: CSSProperties = {
  flexShrink: 0,
  width: '20px',
  height: '22px',
  border: 'none',
  background: 'transparent',
  color: 'rgba(79, 56, 18, 0.68)',
  cursor: 'pointer',
  fontSize: '10px',
  padding: 0,
  marginTop: '2px',
}

export const regionButtonStyle: CSSProperties = {
  flex: 1,
  display: 'grid',
  gap: '2px',
  minWidth: 0,
  textAlign: 'left',
  background: 'transparent',
  border: 'none',
  borderRadius: '8px',
  padding: '1px 3px',
  cursor: 'pointer',
}

export const selectedButtonStyle: CSSProperties = {
  background: 'rgba(59, 130, 246, 0.08)',
  boxShadow: 'inset 0 0 0 1px rgba(59, 130, 246, 0.16)',
}

export const regionTopRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: '6px',
}

export const regionBadgeStyle: CSSProperties = {
  display: 'inline-flex',
  alignItems: 'center',
  width: 'fit-content',
  padding: '2px 7px',
  borderRadius: '999px',
  background: 'rgba(122, 90, 36, 0.10)',
  color: '#6D4B18',
  fontSize: '9px',
  fontWeight: 700,
}

export const regionTitleStyle: CSSProperties = {
  fontSize: '13px',
  fontWeight: 800,
  color: '#1F160A',
  lineHeight: 1.25,
  wordBreak: 'keep-all',
}

export const regionMetaStyle: CSSProperties = {
  fontSize: '10px',
  fontWeight: 700,
  color: 'rgba(90, 61, 31, 0.68)',
  lineHeight: 1.3,
}

export const childrenWrapStyle: CSSProperties = {
  display: 'grid',
  gap: '7px',
  padding: '6px 8px 8px 18px',
  marginLeft: '12px',
  borderLeft: '2px solid rgba(120, 98, 64, 0.18)',
}

export const emptyChildrenStyle: CSSProperties = {
  padding: '4px 0 2px 30px',
  fontSize: '12px',
  color: 'rgba(56, 40, 18, 0.55)',
}

export const createSubRegionFormStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '1fr',
  gap: '6px',
}

export const subRegionLimitStyle: CSSProperties = {
  padding: '1px 0 0 1px',
  fontSize: '10px',
  fontWeight: 700,
  color: 'rgba(90, 61, 31, 0.58)',
}

// ── SubRegion ─────────────────────────────────────────────────────────────────

export const treeItemBlockStyle: CSSProperties = {
  display: 'grid',
  gap: '5px',
}

export const treeRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'flex-start',
  gap: '6px',
}

export const smallCollapseButtonStyle: CSSProperties = {
  flexShrink: 0,
  width: '18px',
  height: '20px',
  border: 'none',
  borderRadius: '6px',
  background: 'transparent',
  color: '#5E4521',
  cursor: 'pointer',
  fontSize: '10px',
  padding: 0,
  marginTop: '3px',
}

export const subRegionButtonStyle: CSSProperties = {
  flex: 1,
  minWidth: 0,
  display: 'grid',
  gap: '3px',
  textAlign: 'left',
  background: 'transparent',
  border: 'none',
  borderRadius: '8px',
  padding: '4px 6px',
  cursor: 'pointer',
}

export const subRegionTopRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: '6px',
}

export const subRegionBadgeStyle: CSSProperties = {
  display: 'inline-flex',
  width: 'fit-content',
  padding: '2px 6px',
  borderRadius: '999px',
  background: 'rgba(59, 130, 246, 0.10)',
  color: '#215C99',
  fontSize: '9px',
  fontWeight: 700,
}

export const subRegionTitleStyle: CSSProperties = {
  fontSize: '12px',
  fontWeight: 700,
  color: '#23170A',
  lineHeight: 1.3,
}

export const subRegionMetaStyle: CSSProperties = {
  fontSize: '10px',
  fontWeight: 700,
  color: 'rgba(69, 90, 120, 0.68)',
  lineHeight: 1.3,
}

export const subRegionInnerWrapStyle: CSSProperties = {
  marginLeft: '28px',
  padding: '4px 0 0 12px',
  borderLeft: '1.5px dashed rgba(120, 98, 64, 0.28)',
  position: 'relative',
  display: 'grid',
  gap: '7px',
}

// ── Research Node ─────────────────────────────────────────────────────────────

export const subResearchRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'flex-start',
  gap: '6px',
}

export const researchButtonStyle: CSSProperties = {
  flex: 1,
  minWidth: 0,
  display: 'grid',
  gap: '2px',
  textAlign: 'left',
  background: 'transparent',
  border: 'none',
  borderRadius: '8px',
  padding: '4px 6px',
  cursor: 'pointer',
}

export const researchBadgeStyle: CSSProperties = {
  display: 'inline-flex',
  width: 'fit-content',
  padding: '2px 6px',
  borderRadius: '999px',
  background: 'rgba(20, 140, 120, 0.10)',
  color: '#136857',
  fontSize: '9px',
  fontWeight: 700,
}

export const researchTitleStyle: CSSProperties = {
  fontSize: '12px',
  fontWeight: 700,
  color: '#23170A',
  lineHeight: 1.28,
}

export const researchMetaStyle: CSSProperties = {
  fontSize: '10px',
  color: 'rgba(56, 40, 18, 0.54)',
  lineHeight: 1.25,
}
