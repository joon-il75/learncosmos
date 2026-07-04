import type { CSSProperties } from 'react'

// ── Add Item Modal (서브지역 / 탐험지점 / 연구지점 추가) ──────────────────────

export const addItemOverlayStyle: CSSProperties = {
  background: 'rgba(10, 5, 1, 0.58)',
  padding: '24px 16px',
  backdropFilter: 'blur(3px)',
}

export const addItemModalBoxStyle: CSSProperties = {
  width: 'min(420px, 100%)',
  maxWidth: '100%',
  maxHeight: 'calc(100dvh - 48px)',
  background: 'linear-gradient(180deg, #FFFAED 0%, #F8ECCE 100%)',
  borderRadius: 14,
  border: '1.5px solid rgba(196, 142, 54, 0.45)',
  boxShadow: '0 24px 60px rgba(20, 10, 2, 0.55)',
  padding: '22px 20px',
  display: 'flex',
  flexDirection: 'column',
  gap: 14,
  overflowY: 'auto',
}

export const addExplorationOverlayStyle: CSSProperties = {
  background: 'rgba(10, 5, 1, 0.58)',
  padding: '24px 16px',
  backdropFilter: 'blur(3px)',
}

export const addExplorationModalBoxStyle: CSSProperties = {
  width: 'min(860px, 100%)',
  maxWidth: '100%',
  maxHeight: 'calc(100dvh - 48px)',
  background: 'linear-gradient(180deg, #FFFAED 0%, #F8ECCE 100%)',
  borderRadius: 14,
  border: '1.5px solid rgba(196, 142, 54, 0.45)',
  boxShadow: '0 24px 60px rgba(20, 10, 2, 0.55)',
  padding: '22px 20px',
  display: 'flex',
  flexDirection: 'column',
  gap: 14,
  overflowY: 'auto',
}

export const addItemModalHeaderStyle: CSSProperties = {
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
}

export const addItemModalTitleStyle: CSSProperties = {
  fontSize: '16px',
  fontWeight: 700,
  color: '#2C1504',
  margin: 0,
}

export const addItemModalLabelStyle: CSSProperties = {
  display: 'block',
  fontSize: '12px',
  fontWeight: 700,
  color: '#6B4A1A',
  marginBottom: 5,
}

export const addItemModalInputStyle: CSSProperties = {
  width: '100%',
  height: 42,
  border: '1.5px solid rgba(122, 90, 36, 0.30)',
  borderRadius: '10px',
  background: 'rgba(255, 252, 244, 0.92)',
  color: '#1F160A',
  fontSize: '14px',
  fontWeight: 600,
  padding: '0 12px',
  outline: 'none',
  boxSizing: 'border-box',
}

export const addItemModalRowStyle: CSSProperties = {
  display: 'flex',
  gap: 8,
  alignItems: 'center',
}

export const addItemModalPrimaryBtnStyle: CSSProperties = {
  flex: 1,
  height: 42,
  border: '1.5px solid rgba(141, 190, 255, 0.82)',
  borderRadius: '10px',
  background: 'linear-gradient(180deg, rgba(58, 126, 220, 0.92), rgba(28, 82, 164, 0.96))',
  color: '#F4FAFF',
  fontSize: '14px',
  fontWeight: 700,
  cursor: 'pointer',
  padding: '0 16px',
  whiteSpace: 'nowrap',
  boxShadow: 'inset 0 1px 0 rgba(255,255,255,0.22), 0 2px 6px rgba(0,0,0,0.22)',
}

export const addItemModalCancelBtnStyle: CSSProperties = {
  height: 42,
  border: '1.5px solid rgba(122, 90, 36, 0.28)',
  borderRadius: '10px',
  background: 'rgba(122, 90, 36, 0.08)',
  color: '#6B4A1A',
  fontSize: '14px',
  fontWeight: 700,
  cursor: 'pointer',
  padding: '0 16px',
  whiteSpace: 'nowrap',
}

export const addItemModalCheckBtnStyle: CSSProperties = {
  height: 42,
  flexShrink: 0,
  border: '1.5px solid rgba(20, 140, 120, 0.30)',
  borderRadius: '10px',
  background: 'rgba(20, 140, 120, 0.07)',
  color: '#136857',
  fontSize: '13px',
  fontWeight: 700,
  cursor: 'pointer',
  padding: '0 14px',
  whiteSpace: 'nowrap',
}

export const addItemModalErrorStyle: CSSProperties = {
  fontSize: '13px',
  fontWeight: 600,
  color: '#9B2C2C',
  background: 'rgba(155, 44, 44, 0.06)',
  border: '1px solid rgba(155, 44, 44, 0.20)',
  borderRadius: 8,
  padding: '8px 12px',
}

export const addItemModalInfoStyle: CSSProperties = {
  fontSize: '12px',
  color: 'rgba(90, 61, 31, 0.65)',
  lineHeight: 1.45,
}

export const parentContextInfoStyle: CSSProperties = {
  ...addItemModalInfoStyle,
  padding: '7px 9px',
  borderRadius: '8px',
  border: '1px solid rgba(122, 90, 36, 0.22)',
  background: 'rgba(122, 90, 36, 0.08)',
  color: '#3F2A12',
  fontSize: '12px',
  fontWeight: 850,
}

export const explorationChoiceBtnStyle: CSSProperties = {
  flex: 1,
  height: 38,
  border: '1px solid transparent',
  borderBottom: '2px solid transparent',
  borderRadius: '10px 10px 0 0',
  background: 'transparent',
  color: 'rgba(90, 58, 16, 0.72)',
  fontSize: '12px',
  fontWeight: 850,
  cursor: 'pointer',
  padding: '0 10px',
  whiteSpace: 'nowrap',
  position: 'relative',
  marginBottom: -1,
  transition: 'background 160ms ease, border-color 160ms ease, color 160ms ease',
}

// ── Child Object Tabs ─────────────────────────────────────────────────────────

export const childObjectTabRowStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(3, minmax(0, 1fr))',
  gap: 0,
  padding: '4px 4px 0',
  borderBottom: '1px solid rgba(122, 90, 36, 0.24)',
  background: 'rgba(122, 90, 36, 0.07)',
  borderRadius: '12px 12px 0 0',
}

export const childObjectTabButtonStyle: CSSProperties = {
  height: 42,
  border: '1px solid transparent',
  borderBottom: '2px solid transparent',
  borderRadius: '10px 10px 0 0',
  background: 'transparent',
  color: 'rgba(78, 52, 18, 0.70)',
  fontSize: '13px',
  fontWeight: 850,
  cursor: 'pointer',
  padding: '0 10px',
  whiteSpace: 'nowrap',
  position: 'relative',
  marginBottom: -1,
  transition: 'background 160ms ease, border-color 160ms ease, color 160ms ease',
}

export const childObjectTabButtonActiveStyle: CSSProperties = {
  borderColor: 'rgba(122, 90, 36, 0.26)',
  borderBottomColor: '#FFFAED',
  background: '#FFFAED',
  color: '#1F4E8A',
  boxShadow: '0 -1px 0 rgba(255,255,255,0.65), 0 -5px 14px rgba(63, 42, 18, 0.06)',
}
