import type { CSSProperties } from 'react';

export const pageStyle = {
  position: 'relative',
  minHeight: '100vh',
  overflow: 'hidden',
} as const;

export const pageBackgroundOverlayStyle = {
  position: 'absolute',
  inset: 0,
  background: 'radial-gradient(ellipse at 50% 40%, rgba(160, 110, 45, 0.22), transparent 55%), linear-gradient(rgba(36, 24, 12, 0.72), rgba(28, 18, 8, 0.82))',
} as const;

export const loadingPageStyle = {
  minHeight: '100vh',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  background: 'linear-gradient(rgba(7, 14, 26, 0.76), rgba(7, 14, 26, 0.88))',
  padding: '24px',
} as const;

export const loadingCardStyle = {
  width: '100%',
  maxWidth: '540px',
  padding: '32px',
  borderRadius: '24px',
  background: 'rgba(8, 18, 33, 0.76)',
  border: '1px solid rgba(194, 210, 245, 0.16)',
  textAlign: 'center',
  color: '#F4F7FF',
  fontSize: '16px',
  backdropFilter: 'blur(10px)',
} as const;

export const navStyle = {
  position: 'relative',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '16px',
  padding: '16px 24px',
  borderBottom: '1px solid rgba(180, 205, 255, 0.12)',
  background: 'rgba(6, 15, 29, 0.58)',
  backdropFilter: 'blur(14px)',
} as const;

export const userBadgeStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: '38px',
  padding: '0 14px',
  borderRadius: '999px',
  border: '1px solid rgba(180, 205, 255, 0.14)',
  color: '#E8EEFF',
  fontSize: '14px',
  background: 'rgba(255, 255, 255, 0.08)',
} as const;

export const mainStyle = {
  position: 'relative',
  maxWidth: '1360px',
  margin: '0 auto',
  padding: '32px 24px 48px',
  display: 'grid',
  gap: '18px',
} as const;

export const heroCardStyle = {
  display: 'grid',
  gap: '18px',
  borderRadius: '28px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  borderLeft: '2px solid rgba(50,200,100,0.35)',
  background: 'rgba(7, 18, 34, 0.58)',
  backdropFilter: 'blur(16px)',
  padding: '24px',
} as const;

export const eyebrowStyle = {
  fontSize: '13px',
  fontWeight: 700,
  color: '#C8D1E8',
  letterSpacing: '0.08em',
  textTransform: 'uppercase',
} as const;

export const titleStyle = {
  margin: 0,
  fontSize: '36px',
  color: '#F7FAFF',
} as const;

export const descriptionStyle = {
  margin: 0,
  fontSize: '15px',
  lineHeight: 1.7,
  color: '#D9E3FB',
  maxWidth: '760px',
} as const;

export const metaGridStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))',
  gap: '12px',
} as const;

export const metaCardStyle = {
  display: 'grid',
  gap: '6px',
  padding: '16px',
  borderRadius: '18px',
  border: '1px solid rgba(194, 210, 245, 0.12)',
  background: 'rgba(255,255,255,0.06)',
} as const;

export const metaLabelStyle = {
  fontSize: '12px',
  color: '#A9B6D2',
} as const;

export const metaValueStyle = {
  fontSize: '15px',
  fontWeight: 700,
  color: '#F7FAFF',
} as const;

export const lumiInsightCardStyle = {
  display: 'grid',
  gap: '14px',
  marginTop: '18px',
  padding: '18px 20px',
  borderRadius: '22px',
  border: '1px solid rgba(118, 170, 255, 0.18)',
  background: 'linear-gradient(180deg, rgba(16, 24, 40, 0.86), rgba(10, 16, 28, 0.94))',
  boxShadow: '0 20px 40px rgba(0,0,0,0.18)',
} as const;

export const lumiInsightHeaderStyle = {
  display: 'grid',
  gridTemplateColumns: '52px minmax(0, 1fr)',
  gap: '14px',
  alignItems: 'center',
} as const;

export const lumiAvatarWrapStyle = {
  width: '52px',
  height: '52px',
  borderRadius: '999px',
  overflow: 'hidden',
  boxShadow: '0 14px 28px rgba(0, 0, 0, 0.22)',
} as const;

export const lumiEyebrowStyle = {
  fontSize: '11px',
  fontWeight: 700,
  letterSpacing: '0.12em',
  textTransform: 'uppercase',
  color: 'rgba(159, 200, 255, 0.78)',
} as const;

export const lumiTitleStyle = {
  color: '#F3F7FF',
  fontSize: '16px',
} as const;

export const lumiMessageStyle = {
  margin: 0,
  color: 'rgba(214, 225, 242, 0.82)',
  fontSize: '14px',
  lineHeight: 1.6,
  whiteSpace: 'pre-wrap',
} as const;

export const lumiMetaRowStyle = {
  display: 'flex',
  flexWrap: 'wrap',
  gap: '8px',
} as const;

export const lumiMetaPillStyle = {
  padding: '6px 10px',
  borderRadius: '999px',
  background: 'rgba(84, 120, 255, 0.12)',
  border: '1px solid rgba(118, 170, 255, 0.18)',
  color: 'rgba(201, 220, 255, 0.92)',
  fontSize: '12px',
} as const;

export const draftMetaEditorStyle = {
  display: 'grid',
  gap: '14px',
  padding: '18px',
  borderRadius: '22px',
  border: '1px solid rgba(194, 210, 245, 0.12)',
  background: 'rgba(255,255,255,0.06)',
} as const;

export const draftMetaEditorHeaderStyle = {
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: '12px',
  flexWrap: 'wrap',
} as const;

export const draftMetaEditorTitleStyle = {
  margin: '4px 0 0',
  fontSize: '20px',
  color: '#F7FAFF',
} as const;

export const draftMetaEditorHintStyle = {
  maxWidth: '340px',
  color: '#AFC0DD',
  fontSize: '12px',
  lineHeight: 1.55,
} as const;

export const draftMetaFieldStyle = {
  display: 'grid',
  gap: '8px',
} as const;

export const draftMetaLabelStyle = {
  color: '#C8D1E8',
  fontSize: '12px',
  fontWeight: 800,
} as const;

export const draftMetaInputStyle = {
  minHeight: '44px',
  borderRadius: '14px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(4, 10, 20, 0.42)',
  color: '#F7FAFF',
  padding: '0 14px',
  fontSize: '15px',
  fontFamily: 'inherit',
  outline: 'none',
} as const;

export const draftMetaTextareaStyle = {
  minHeight: '104px',
  borderRadius: '14px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(4, 10, 20, 0.42)',
  color: '#F7FAFF',
  padding: '12px 14px',
  fontSize: '14px',
  lineHeight: 1.6,
  fontFamily: 'inherit',
  outline: 'none',
  resize: 'vertical',
} as const;

export const draftMetaFooterStyle = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '12px',
  flexWrap: 'wrap',
} as const;

export const activeSaveButtonStyle = {
  minHeight: '42px',
  padding: '0 16px',
  borderRadius: '999px',
  border: '1px solid rgba(112, 166, 255, 0.34)',
  background: 'linear-gradient(135deg, rgba(59, 130, 246, 0.44), rgba(94, 234, 212, 0.22))',
  color: '#F7FAFF',
  fontSize: '14px',
  fontWeight: 900,
  cursor: 'pointer',
} as const;

export const activeSaveButtonDisabledStyle = {
  opacity: 0.52,
  cursor: 'not-allowed',
} as const;
