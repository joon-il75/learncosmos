import type { CSSProperties } from 'react';

export const planningGoalBandStyle = {
  display: 'grid',
  gap: '14px',
  marginTop: '8px',
  padding: '18px 20px',
  borderRadius: '24px',
  border: '1px solid rgba(191, 154, 81, 0.28)',
  background: 'linear-gradient(180deg, rgba(38, 24, 10, 0.76), rgba(24, 15, 7, 0.82))',
  boxShadow: '0 18px 48px rgba(10, 6, 3, 0.28)',
} as const;

export const planningGoalBandHeaderStyle = {
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: '16px',
  flexWrap: 'wrap',
} as const;

export const planningGoalBandEyebrowStyle = {
  fontSize: '11px',
  fontWeight: 800,
  letterSpacing: '0.12em',
  textTransform: 'uppercase',
  color: 'rgba(231, 199, 137, 0.78)',
} as const;

export const planningGoalBandTitleStyle = {
  margin: '4px 0 0',
  fontSize: '15px',
  fontWeight: 800,
  color: '#F4E9D2',
} as const;

export const planningGoalBandBodyStyle = {
  margin: 0,
  fontSize: '18px',
  lineHeight: 1.7,
  color: 'rgba(255, 244, 223, 0.94)',
} as const;

export const planningGoalBandMetaRowStyle = {
  display: 'flex',
  flexWrap: 'wrap',
  gap: '8px',
} as const;

export const planningGoalBandMetaPillStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: '30px',
  padding: '0 12px',
  borderRadius: '999px',
  border: '1px solid rgba(212, 177, 105, 0.24)',
  background: 'rgba(255, 247, 231, 0.08)',
  color: 'rgba(243, 226, 192, 0.92)',
  fontSize: '12px',
  fontWeight: 600,
} as const;

export const planningGoalBandHintStyle = {
  margin: 0,
  fontSize: '13px',
  lineHeight: 1.6,
  color: 'rgba(218, 198, 161, 0.82)',
} as const;

export const planningGoalBandButtonStyle = {
  minHeight: '40px',
  padding: '0 16px',
  borderRadius: '999px',
  border: '1px solid rgba(215, 180, 106, 0.34)',
  background: 'linear-gradient(180deg, rgba(244, 214, 153, 0.96), rgba(205, 154, 67, 0.96))',
  color: '#4A2C00',
  fontSize: '13px',
  fontWeight: 900,
  cursor: 'pointer',
  boxShadow: '0 10px 28px rgba(70, 44, 7, 0.22)',
} as const;

export const planningGoalModalOverlayStyle = {
  padding: '24px',
  background: 'rgba(18, 11, 8, 0.62)',
  backdropFilter: 'blur(8px)',
} as const;

export const planningGoalModalCardStyle = {
  width: 'min(880px, 100%)',
  maxHeight: 'min(860px, calc(100vh - 48px))',
  display: 'grid',
  gridTemplateRows: 'auto minmax(0, 1fr)',
  gap: '14px',
  padding: '18px',
  borderRadius: '28px',
  border: '1px solid rgba(215, 180, 106, 0.28)',
  background: 'linear-gradient(180deg, rgba(43, 28, 14, 0.98), rgba(27, 18, 9, 0.98))',
  boxShadow: '0 28px 80px rgba(0, 0, 0, 0.38)',
} as const;

export const planningGoalModalHeaderStyle = {
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: '16px',
} as const;

export const planningGoalModalTitleStyle = {
  margin: '4px 0 0',
  fontSize: '22px',
  lineHeight: 1.35,
  color: '#FFF3DA',
} as const;

export const planningGoalModalDescriptionStyle = {
  margin: '8px 0 0',
  fontSize: '13px',
  lineHeight: 1.65,
  color: 'rgba(230, 214, 184, 0.76)',
} as const;

export const planningGoalModalCloseButtonStyle = {
  minWidth: '42px',
  minHeight: '42px',
  borderRadius: '999px',
  border: '1px solid rgba(215, 180, 106, 0.24)',
  background: 'rgba(255, 247, 231, 0.08)',
  color: '#F5E7CB',
  fontSize: '18px',
  fontWeight: 700,
  cursor: 'pointer',
} as const;

export const planningGoalModalPanelWrapStyle = {
  minHeight: 0,
  overflow: 'hidden',
  borderRadius: '22px',
  background: '#FFFDF7',
} as const;

export const planningGoalModalEmptyStateStyle = {
  display: 'grid',
  gap: '14px',
  padding: '24px',
  borderRadius: '22px',
  background: 'linear-gradient(180deg, rgba(255, 248, 235, 0.98), rgba(247, 234, 204, 0.98))',
  color: '#5B3912',
} as const;

export const saveMessageStyle = {
  color: '#CFE0FF',
  fontSize: '13px',
  lineHeight: 1.5,
} as const;
