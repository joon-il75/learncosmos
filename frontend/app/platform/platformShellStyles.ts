import type { CSSProperties } from 'react';

export const pageStyle: CSSProperties = {
  minHeight: '100vh',
  background: 'linear-gradient(135deg, #F6FAFF 0%, #F4F9EF 44%, #FFF6E8 100%)',
  color: '#162D3B',
};

export const headerStyle: CSSProperties = {
  position: 'fixed',
  top: 0,
  left: 0,
  right: 0,
  zIndex: 50,
  background: 'rgba(248, 252, 250, 0.78)',
  backdropFilter: 'blur(14px) saturate(1.08)',
  WebkitBackdropFilter: 'blur(14px) saturate(1.08)',
  borderBottom: '1px solid rgba(45, 83, 105, 0.1)',
};

export const headerInnerStyle: CSSProperties = {
  width: 'min(1120px, calc(100% - 32px))',
  margin: '0 auto',
  padding: '15px 0',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '16px',
};

export const headerLinksStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: '10px',
  flexWrap: 'wrap',
  justifyContent: 'flex-end',
};

export const headerLinkStyle: CSSProperties = {
  textDecoration: 'none',
  color: '#254B5F',
  fontSize: '13px',
  fontWeight: 900,
  padding: '8px 13px',
  borderRadius: '999px',
  background: 'rgba(255,255,255,0.58)',
};

export const mainStyle: CSSProperties = {
  width: 'min(1120px, calc(100% - 32px))',
  margin: '0 auto',
  padding: '112px 0 64px',
};

export const heroStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1.18fr) minmax(320px, 0.82fr)',
  gap: '24px',
  alignItems: 'end',
};

export const eyebrowStyle: CSSProperties = {
  width: 'fit-content',
  padding: '8px 13px',
  borderRadius: '999px',
  background: 'rgba(57, 113, 138, 0.1)',
  color: '#2B687F',
  fontSize: '12px',
  fontWeight: 900,
  letterSpacing: '0.08em',
  textTransform: 'uppercase',
};

export const titleStyle: CSSProperties = {
  margin: '16px 0 0',
  maxWidth: '820px',
  fontSize: 'clamp(34px, 7vw, 72px)',
  lineHeight: 0.96,
  letterSpacing: '-0.06em',
  color: '#162D3B',
};

export const subtitleStyle: CSSProperties = {
  margin: '18px 0 0',
  maxWidth: '720px',
  fontSize: 'clamp(16px, 2.2vw, 20px)',
  lineHeight: 1.75,
  color: '#536977',
};

export const statusPanelStyle: CSSProperties = {
  minHeight: '230px',
  borderRadius: '30px',
  padding: '28px',
  background: 'radial-gradient(circle at 24% 18%, rgba(81, 155, 181, 0.25), transparent 34%), radial-gradient(circle at 75% 72%, rgba(236, 174, 84, 0.22), transparent 34%), rgba(255,255,255,0.62)',
  boxShadow: '0 30px 86px rgba(58, 83, 99, 0.12)',
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'space-between',
};

export const statusLabelStyle: CSSProperties = {
  width: 'fit-content',
  padding: '9px 13px',
  borderRadius: '999px',
  background: 'rgba(255,255,255,0.7)',
  color: '#315D70',
  fontSize: '13px',
  fontWeight: 900,
};

export const signalStyle: CSSProperties = {
  height: '118px',
  borderRadius: '30px',
  background: 'linear-gradient(135deg, rgba(28, 66, 84, 0.12), rgba(255,255,255,0.42))',
  position: 'relative',
  overflow: 'hidden',
};

export const signalDotStyle: CSSProperties = {
  position: 'absolute',
  left: '28px',
  top: '34px',
  width: '38px',
  height: '38px',
  borderRadius: '50%',
  background: '#2E708B',
  boxShadow: '92px -12px 0 #8BC7A6, 198px 22px 0 #E9A94E, 292px 0 0 #AFC7FF',
};

export const navStyle: CSSProperties = {
  marginTop: '42px',
  display: 'flex',
  flexWrap: 'wrap',
  gap: '10px',
};

export const navLinkStyle = (active: boolean): CSSProperties => ({
  textDecoration: 'none',
  padding: '11px 16px',
  borderRadius: '999px',
  color: active ? '#FFFFFF' : '#315D70',
  background: active ? '#162D3B' : 'rgba(255,255,255,0.6)',
  fontSize: '14px',
  fontWeight: 900,
  boxShadow: active ? '0 14px 34px rgba(22, 45, 59, 0.18)' : 'none',
});

export const contentGridStyle: CSSProperties = {
  marginTop: '32px',
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1.08fr) minmax(340px, 0.72fr)',
  gap: '24px',
};

export const sectionSurfaceStyle: CSSProperties = {
  borderRadius: '28px',
  padding: '32px',
  background: 'rgba(255,255,255,0.64)',
  boxShadow: '0 26px 74px rgba(58, 83, 99, 0.1)',
};

export const sectionEyebrowStyle: CSSProperties = {
  margin: 0,
  fontSize: '12px',
  fontWeight: 900,
  color: '#2B687F',
  letterSpacing: '0.12em',
  textTransform: 'uppercase',
};

export const sectionTitleStyle: CSSProperties = {
  margin: '12px 0 0',
  fontSize: 'clamp(27px, 4vw, 44px)',
  lineHeight: 1.07,
  letterSpacing: '-0.04em',
  color: '#162D3B',
};

export const sectionDescriptionStyle: CSSProperties = {
  margin: '18px 0 0',
  fontSize: '16px',
  lineHeight: 1.78,
  color: '#536977',
};

export const actionRowStyle: CSSProperties = {
  marginTop: '28px',
  display: 'flex',
  flexWrap: 'wrap',
  gap: '12px',
};

export const primaryActionStyle: CSSProperties = {
  textDecoration: 'none',
  padding: '13px 18px',
  borderRadius: '999px',
  background: '#162D3B',
  color: '#FFFFFF',
  fontSize: '14px',
  fontWeight: 900,
};

export const secondaryActionStyle: CSSProperties = {
  padding: '13px 18px',
  borderRadius: '999px',
  background: 'rgba(57, 113, 138, 0.1)',
  color: '#315D70',
  fontSize: '14px',
  fontWeight: 900,
};

export const noteListStyle: CSSProperties = {
  margin: '28px 0 0',
  padding: 0,
  listStyle: 'none',
  display: 'grid',
  gap: '13px',
};

export const noteItemStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '10px minmax(0, 1fr)',
  gap: '12px',
  color: '#4F6572',
  fontSize: '14px',
  lineHeight: 1.65,
};

export const noteDotStyle: CSSProperties = {
  width: '8px',
  height: '8px',
  borderRadius: '50%',
  marginTop: '8px',
  background: '#2E708B',
};


export const newsBoardStyle: CSSProperties = {
  marginTop: '28px',
  display: 'grid',
  gap: '18px',
  minHeight: '340px',
  padding: '24px',
  borderRadius: '30px',
  background: 'linear-gradient(135deg, rgba(255,255,255,0.72), rgba(235,249,255,0.54))',
  boxShadow: 'inset 0 1px 0 rgba(255,255,255,0.78), 0 24px 64px rgba(58, 83, 99, 0.08)',
};

export const newsBoardHeaderStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '14px',
  paddingBottom: '14px',
  borderBottom: '1px solid rgba(57, 113, 138, 0.14)',
};

export const newsBoardBadgeStyle: CSSProperties = {
  width: 'fit-content',
  padding: '7px 11px',
  borderRadius: '999px',
  background: 'rgba(106, 210, 193, 0.16)',
  color: '#177483',
  fontSize: '11px',
  fontWeight: 950,
  letterSpacing: '0.1em',
  textTransform: 'uppercase',
};

export const newsBoardTitleStyle: CSSProperties = {
  margin: 0,
  color: '#162D3B',
  fontSize: 'clamp(20px, 2.2vw, 28px)',
  lineHeight: 1.16,
  letterSpacing: '-0.04em',
};

export const newsBoardDescriptionStyle: CSSProperties = {
  margin: 0,
  color: '#536977',
  fontSize: '14px',
  lineHeight: 1.7,
};

export const newsBoardEmptyStyle: CSSProperties = {
  minHeight: '190px',
  display: 'grid',
  placeItems: 'center',
  alignContent: 'center',
  gap: '10px',
  borderRadius: '24px',
  background: 'repeating-linear-gradient(180deg, rgba(22,45,59,0.055) 0 1px, transparent 1px 56px), rgba(255,255,255,0.38)',
  textAlign: 'center',
};

export const newsBoardEmptyTitleStyle: CSSProperties = {
  margin: 0,
  color: '#162D3B',
  fontSize: 'clamp(20px, 2.4vw, 30px)',
  lineHeight: 1.18,
  letterSpacing: '-0.04em',
};

export const newsBoardEmptyBodyStyle: CSSProperties = {
  margin: 0,
  maxWidth: '460px',
  color: '#5B6E78',
  fontSize: '14px',
  lineHeight: 1.7,
};

export const newsBoardAdminOnlyStyle: CSSProperties = {
  width: 'fit-content',
  padding: '10px 13px',
  borderRadius: '999px',
  background: 'rgba(22,45,59,0.08)',
  color: '#315D70',
  fontSize: '13px',
  fontWeight: 900,
};

export const sidePanelStyle: CSSProperties = {
  display: 'grid',
  gap: '18px',
};

export const linksPanelStyle: CSSProperties = {
  borderRadius: '28px',
  padding: '28px',
  background: 'linear-gradient(180deg, rgba(22, 45, 59, 0.94), rgba(47, 86, 102, 0.84))',
  color: '#FFFFFF',
  boxShadow: '0 24px 70px rgba(22, 45, 59, 0.16)',
};

export const panelTitleStyle: CSSProperties = {
  margin: 0,
  fontSize: '20px',
  letterSpacing: '-0.02em',
};

export const panelSubtitleStyle: CSSProperties = {
  margin: '8px 0 0',
  color: 'rgba(255,255,255,0.74)',
  fontSize: '13px',
  lineHeight: 1.6,
};

export const linkListStyle: CSSProperties = {
  margin: '20px 0 0',
  padding: 0,
  listStyle: 'none',
  display: 'grid',
  gap: '12px',
};

export const linkItemStyle: CSSProperties = {
  display: 'block',
  textDecoration: 'none',
  color: '#FFFFFF',
  padding: '14px 0',
  borderTop: '1px solid rgba(255,255,255,0.14)',
};

export const linkLabelStyle: CSSProperties = {
  display: 'block',
  fontSize: '14px',
  fontWeight: 900,
};

export const linkBodyStyle: CSSProperties = {
  display: 'block',
  marginTop: '5px',
  color: 'rgba(255,255,255,0.72)',
  fontSize: '13px',
  lineHeight: 1.55,
};

export const standardPanelStyle: CSSProperties = {
  borderRadius: '30px',
  padding: '24px',
  background: 'rgba(255,255,255,0.58)',
};

export const standardListStyle: CSSProperties = {
  margin: '18px 0 0',
  padding: 0,
  listStyle: 'none',
  display: 'grid',
  gap: '14px',
};

export const standardLabelStyle: CSSProperties = {
  display: 'block',
  fontSize: '13px',
  fontWeight: 900,
  color: '#254B5F',
};

export const standardBodyStyle: CSSProperties = {
  display: 'block',
  marginTop: '4px',
  color: '#5B6E78',
  fontSize: '13px',
  lineHeight: 1.55,
};

export const mobileStyle = `
@media (max-width: 760px) {
  .platform-hero { grid-template-columns: 1fr !important; }
  .platform-content-grid { grid-template-columns: 1fr !important; }
  .platform-main { width: 100% !important; padding-top: 32px !important; }
  .platform-section-surface { padding: 24px !important; border-radius: 28px !important; }
  .platform-news-board-header { align-items: flex-start !important; flex-direction: column !important; }
  .platform-header-inner { width: min(100% - 24px, 1120px) !important; }
}
`;
