import type { CSSProperties } from 'react'

// ── Edit Panel ────────────────────────────────────────────────────────────────

export const editPanelStyle: CSSProperties = {
  flex: '0 1 auto',   // 트리 스크롤 영역이 남은 공간 우선, 편집 패널은 콘텐츠 높이만
  maxHeight: '52%',   // 디바이스 셸 높이의 절반 이하로 제한, 내부 스크롤
  display: 'flex',
  flexDirection: 'column',
  minHeight: 0,
  paddingRight: '10px',
  margin: '5px 6px 4px',
  borderRadius: '7px',
  background: 'linear-gradient(180deg, rgba(25, 78, 86, 0.92) 0%, rgba(14, 48, 58, 0.96) 100%)',
  border: '1px solid rgba(16, 125, 132, 0.62)',
  borderTop: '2px solid rgba(74, 214, 205, 0.78)',
  boxShadow: '0 -1px 0 rgba(255,255,255,0.22), 0 7px 16px rgba(16, 35, 44, 0.22), inset 0 1px 0 rgba(255,255,255,0.18)',
}

export const editFormAreaStyle: CSSProperties = {
  flex: '1 1 auto',
  overflowY: 'auto',
  padding: '8px 5px 7px 10px',
  display: 'flex',
  flexDirection: 'column',
  gap: '7px',
  minHeight: 0,
  background: 'linear-gradient(180deg, rgba(232, 250, 246, 0.95) 0%, rgba(204, 232, 231, 0.92) 100%)',
  borderRadius: '6px',
  border: '1px solid rgba(255, 255, 255, 0.62)',
}

export const pendingDeleteNoticeStyle: CSSProperties = {
  padding: '7px 8px',
  borderRadius: '8px',
  border: '1px solid rgba(185, 58, 42, 0.28)',
  background: 'rgba(185, 58, 42, 0.08)',
  color: '#8F2B22',
  fontSize: '11px',
  fontWeight: 800,
  lineHeight: 1.35,
}

export const pendingCreateNoticeStyle: CSSProperties = {
  padding: '7px 8px',
  borderRadius: '8px',
  border: '1px dashed rgba(37, 99, 235, 0.34)',
  background: 'rgba(37, 99, 235, 0.08)',
  color: '#1D4ED8',
  fontSize: '11px',
  fontWeight: 800,
  lineHeight: 1.35,
}

export const editSectionLabelStyle: CSSProperties = {
  fontSize: '11px',
  fontWeight: 700,
  letterSpacing: '0.05em',
  textTransform: 'uppercase',
  color: 'rgba(90, 61, 31, 0.55)',
  marginBottom: '3px',
}

export const editInputRowStyle: CSSProperties = {
  display: 'flex',
  gap: '6px',
  alignItems: 'center',
}

export const editStickyBarStyle: CSSProperties = {
  flexShrink: 0,
  display: 'flex',
  gap: '5px',
  padding: '5px 10px 7px',
  borderTop: '1px solid rgba(122, 90, 36, 0.10)',
  alignItems: 'center',
}

export const editMutationMsgStyle: CSSProperties = {
  padding: '0 10px 3px',
  fontSize: '9px',
  fontWeight: 600,
  color: 'rgba(90, 61, 31, 0.65)',
  lineHeight: 1.3,
  minHeight: 12,
}

// ── Recommendation Modal ──────────────────────────────────────────────────────

export const recommendationOverlayStyle: CSSProperties = {
  padding: '16px 12px',
  background: 'rgba(20, 12, 4, 0.50)',
  backdropFilter: 'blur(3px)',
}

export const recommendationModalStyle: CSSProperties = {
  width: 'min(520px, calc(100vw - 24px))',
  maxHeight: 'min(620px, calc(100svh - 32px))',
  display: 'grid',
  gridTemplateRows: 'auto auto 1fr auto',
  gap: 10,
  padding: 14,
  borderRadius: 8,
  border: '1px solid rgba(122, 90, 36, 0.28)',
  background: 'linear-gradient(180deg, rgba(255, 250, 236, 0.98), rgba(245, 232, 204, 0.98))',
  boxShadow: '0 18px 46px rgba(32, 20, 8, 0.34)',
}

export const recommendationHeaderStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 10,
}

export const recommendationTitleStyle: CSSProperties = {
  margin: 0,
  fontSize: '15px',
  fontWeight: 800,
  color: '#1F160A',
}

export const recommendationSearchRowStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '1fr auto',
  gap: 6,
}

export const recommendationListStyle: CSSProperties = {
  display: 'grid',
  gap: 8,
  overflowY: 'auto',
  minHeight: 0,
  paddingRight: 2,
}

export const recommendationItemStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '1fr auto',
  gap: 10,
  padding: '12px 12px',
  borderRadius: 8,
  border: '1px solid rgba(122, 90, 36, 0.18)',
  background: 'rgba(255, 252, 244, 0.82)',
}

export const recommendationItemTitleStyle: CSSProperties = {
  margin: 0,
  fontSize: '14px',
  fontWeight: 700,
  color: '#1A0E03',
  lineHeight: 1.4,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
}

export const recommendationItemMetaStyle: CSSProperties = {
  marginTop: 4,
  fontSize: '11px',
  fontWeight: 600,
  color: 'rgba(90, 61, 31, 0.75)',
  lineHeight: 1.4,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
}

export const recommendationMessageStyle: CSSProperties = {
  minHeight: 14,
  fontSize: '10px',
  fontWeight: 700,
  color: 'rgba(90, 61, 31, 0.72)',
}

// ── Action Buttons ────────────────────────────────────────────────────────────

export const saveButtonStyle: CSSProperties = {
  width: '100%',
  height: 32,
  border: '1px solid rgba(141, 190, 255, 0.82)',
  borderRadius: '8px',
  background: 'linear-gradient(180deg, rgba(58, 126, 220, 0.88), rgba(28, 82, 164, 0.92))',
  color: '#F4FAFF',
  fontSize: '12px',
  fontWeight: 700,
  cursor: 'pointer',
  padding: '0 8px',
  whiteSpace: 'nowrap',
  boxShadow: 'inset 0 1px 0 rgba(255, 255, 255, 0.28), 0 1px 3px rgba(0, 0, 0, 0.30)',
  textShadow: '0 1px 1px rgba(0, 0, 0, 0.28)',
}

export const cancelButtonStyle: CSSProperties = {
  width: '100%',
  height: 32,
  border: '1px solid rgba(255, 222, 155, 0.62)',
  borderRadius: '8px',
  background: 'rgba(255, 236, 190, 0.18)',
  color: '#FFF2CE',
  fontSize: '12px',
  fontWeight: 700,
  cursor: 'pointer',
  padding: '0 8px',
  whiteSpace: 'nowrap',
  boxShadow: 'inset 0 1px 0 rgba(255, 255, 255, 0.16)',
  textShadow: '0 1px 1px rgba(0, 0, 0, 0.30)',
}

export const modalCloseButtonStyle: CSSProperties = {
  height: 30,
  border: '1px solid rgba(122, 90, 36, 0.32)',
  borderRadius: '8px',
  background: 'rgba(200, 160, 80, 0.14)',
  color: '#5A3A10',
  fontSize: '12px',
  fontWeight: 800,
  cursor: 'pointer',
  padding: '0 12px',
  whiteSpace: 'nowrap',
}

export const panelCancelButtonStyle: CSSProperties = {
  height: 30,
  border: '1px solid rgba(122, 90, 36, 0.28)',
  borderRadius: '8px',
  background: 'rgba(180, 140, 80, 0.10)',
  color: '#5A3A10',
  fontSize: '12px',
  fontWeight: 700,
  cursor: 'pointer',
  padding: '0 12px',
  whiteSpace: 'nowrap',
}

export const addNodeModalStyle: CSSProperties = {
  width: 'min(380px, calc(100vw - 24px))',
  display: 'flex',
  flexDirection: 'column',
  gap: 10,
  padding: 16,
  borderRadius: 10,
  border: '1px solid rgba(122, 90, 36, 0.28)',
  background: 'linear-gradient(180deg, rgba(255, 250, 236, 0.98), rgba(245, 232, 204, 0.98))',
  boxShadow: '0 18px 46px rgba(32, 20, 8, 0.34)',
}

export const deleteButtonStyle: CSSProperties = {
  width: '100%',
  height: 32,
  border: '1px solid rgba(255, 164, 139, 0.72)',
  borderRadius: '8px',
  background: 'rgba(181, 55, 36, 0.34)',
  color: '#FFE3DA',
  fontSize: '12px',
  fontWeight: 700,
  cursor: 'pointer',
  padding: '0 8px',
  whiteSpace: 'nowrap',
  boxShadow: 'inset 0 1px 0 rgba(255, 255, 255, 0.14)',
  textShadow: '0 1px 1px rgba(0, 0, 0, 0.30)',
}

export const inactiveToggleButtonStyle: CSSProperties = {
  width: '100%',
  height: 32,
  border: '1px solid rgba(224, 214, 190, 0.52)',
  borderRadius: '8px',
  background: 'rgba(250, 244, 226, 0.14)',
  color: '#F2E8D0',
  fontSize: '12px',
  fontWeight: 700,
  cursor: 'pointer',
  padding: '0 8px',
  whiteSpace: 'nowrap',
  boxShadow: 'inset 0 1px 0 rgba(255, 255, 255, 0.12)',
  textShadow: '0 1px 1px rgba(0, 0, 0, 0.30)',
}

export const addNodeButtonStyle: CSSProperties = {
  height: 30,
  border: '1px solid rgba(122, 90, 36, 0.22)',
  borderRadius: '8px',
  background: 'rgba(122, 90, 36, 0.07)',
  color: '#5F4319',
  fontSize: '12px',
  fontWeight: 600,
  cursor: 'pointer',
  padding: '0 12px',
  whiteSpace: 'nowrap',
}

export const addSubRegionButtonStyle: CSSProperties = {
  border: '1px solid rgba(69, 124, 194, 0.58)',
  background: 'linear-gradient(180deg, rgba(72, 137, 214, 0.24), rgba(45, 95, 168, 0.18))',
  color: '#1F5C99',
  boxShadow: 'inset 0 1px 0 rgba(255, 255, 255, 0.20), 0 1px 3px rgba(35, 70, 120, 0.14)',
}

export const addExplorationButtonStyle: CSSProperties = {
  borderColor: 'rgba(174, 109, 22, 0.34)',
  borderBottomColor: '#FFFAED',
  background: '#FFFAED',
  color: '#7A3F06',
  boxShadow: '0 -1px 0 rgba(255, 255, 255, 0.68), 0 -4px 12px rgba(124, 70, 10, 0.08)',
}

export const addResearchButtonStyle: CSSProperties = {
  border: '1px solid rgba(21, 128, 105, 0.58)',
  background: 'linear-gradient(180deg, rgba(28, 155, 126, 0.24), rgba(13, 105, 85, 0.18))',
  color: '#0C6755',
  boxShadow: 'inset 0 1px 0 rgba(230, 255, 248, 0.22), 0 1px 3px rgba(12, 90, 72, 0.14)',
}

export const checkButtonStyle: CSSProperties = {
  flexShrink: 0,
  height: 30,
  border: '1px solid rgba(20, 140, 120, 0.25)',
  borderRadius: '8px',
  background: 'rgba(20, 140, 120, 0.07)',
  color: '#136857',
  fontSize: '12px',
  fontWeight: 700,
  cursor: 'pointer',
  padding: '0 10px',
  whiteSpace: 'nowrap',
}
