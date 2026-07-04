import type { CSSProperties } from 'react';

export const pageStyle: CSSProperties = {
  minHeight: '100vh',
  background: 'linear-gradient(145deg, #F3FBF8 0%, #E8F3FF 48%, #FFF8E8 100%)',
  color: '#102F43',
};

export const headerStyle: CSSProperties = {
  background: 'rgba(246, 252, 250, 0.74)',
  borderBottom: '1px solid rgba(38, 99, 133, 0.1)',
};

export const headerInnerStyle: CSSProperties = {
  maxWidth: '1120px',
};

export const headerMetaStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: '2px',
  minWidth: 0,
};

export const headerEyebrowStyle: CSSProperties = {
  fontSize: '11px',
  fontWeight: 800,
  letterSpacing: '0.12em',
  textTransform: 'uppercase',
  color: '#4A7C7A',
};

export const headerTitleStyle: CSSProperties = {
  fontSize: '15px',
  fontWeight: 900,
  color: '#14384D',
};

export const mainStyle: CSSProperties = {
  width: 'min(1120px, calc(100% - 32px))',
  margin: '0 auto',
  padding: '112px 0 64px',
};

export const heroStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1.1fr) minmax(260px, 0.7fr)',
  gap: '32px',
  alignItems: 'end',
};

export const heroTextStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: '16px',
};

export const eyebrowPillStyle: CSSProperties = {
  width: 'fit-content',
  padding: '8px 13px',
  borderRadius: '999px',
  background: 'rgba(16, 74, 78, 0.09)',
  color: '#286A67',
  fontSize: '12px',
  fontWeight: 900,
  letterSpacing: '0.08em',
  textTransform: 'uppercase',
};

export const titleStyle: CSSProperties = {
  margin: 0,
  maxWidth: '720px',
  fontSize: 'clamp(34px, 7vw, 76px)',
  lineHeight: 0.95,
  letterSpacing: '-0.06em',
  color: '#102F43',
};

export const subtitleStyle: CSSProperties = {
  margin: 0,
  maxWidth: '680px',
  fontSize: 'clamp(16px, 2.2vw, 20px)',
  lineHeight: 1.75,
  color: '#3A6471',
};

export const statusPanelStyle: CSSProperties = {
  minHeight: '220px',
  borderRadius: '36px',
  padding: '28px',
  background: 'radial-gradient(circle at 20% 20%, rgba(101, 196, 178, 0.28), transparent 34%), linear-gradient(145deg, rgba(255,255,255,0.72), rgba(255,255,255,0.38))',
  boxShadow: '0 28px 80px rgba(41, 91, 111, 0.14)',
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'space-between',
  gap: '28px',
};

export const statusLabelStyle: CSSProperties = {
  fontSize: '13px',
  fontWeight: 900,
  color: '#245B62',
};

export const orbitStyle: CSSProperties = {
  height: '104px',
  borderRadius: '999px',
  background: 'linear-gradient(90deg, rgba(35, 96, 112, 0.12), rgba(255,255,255,0.42), rgba(243, 176, 88, 0.18))',
  position: 'relative',
  overflow: 'hidden',
};

export const orbitDotStyle: CSSProperties = {
  position: 'absolute',
  top: '34px',
  left: '32px',
  width: '36px',
  height: '36px',
  borderRadius: '50%',
  background: '#2A7C78',
  boxShadow: '140px 10px 0 #F2B45D, 260px -4px 0 #9BC9FF',
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
  color: active ? '#FFFFFF' : '#24546A',
  background: active ? '#163E52' : 'rgba(255,255,255,0.54)',
  fontSize: '14px',
  fontWeight: 900,
  boxShadow: active ? '0 14px 30px rgba(22, 62, 82, 0.18)' : 'none',
});

export const contentGridStyle: CSSProperties = {
  marginTop: '32px',
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1fr) minmax(260px, 360px)',
  gap: '28px',
};

export const sectionSurfaceStyle: CSSProperties = {
  borderRadius: '34px',
  padding: '32px',
  background: 'rgba(255,255,255,0.62)',
  boxShadow: '0 24px 70px rgba(41, 91, 111, 0.11)',
};

export const sectionEyebrowStyle: CSSProperties = {
  margin: 0,
  fontSize: '12px',
  fontWeight: 900,
  color: '#347A73',
  letterSpacing: '0.12em',
  textTransform: 'uppercase',
};

export const sectionTitleStyle: CSSProperties = {
  margin: '12px 0 0',
  fontSize: 'clamp(27px, 4vw, 44px)',
  lineHeight: 1.05,
  letterSpacing: '-0.04em',
  color: '#12384E',
};

export const sectionDescriptionStyle: CSSProperties = {
  margin: '18px 0 0',
  fontSize: '16px',
  lineHeight: 1.78,
  color: '#426978',
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
  background: '#12384E',
  color: '#FFFFFF',
  fontSize: '14px',
  fontWeight: 900,
};

export const secondaryActionStyle: CSSProperties = {
  padding: '13px 18px',
  borderRadius: '999px',
  background: 'rgba(18, 56, 78, 0.08)',
  color: '#24546A',
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
  alignItems: 'start',
  color: '#345C6D',
  fontSize: '14px',
  lineHeight: 1.65,
};

export const noteDotStyle: CSSProperties = {
  width: '8px',
  height: '8px',
  borderRadius: '50%',
  marginTop: '8px',
  background: '#54A99C',
};

export const roadmapStyle: CSSProperties = {
  borderRadius: '34px',
  padding: '28px',
  background: 'linear-gradient(180deg, rgba(18, 56, 78, 0.9), rgba(27, 88, 96, 0.82))',
  color: '#FFFFFF',
  boxShadow: '0 24px 70px rgba(25, 72, 89, 0.16)',
};

export const roadmapTitleStyle: CSSProperties = {
  margin: 0,
  fontSize: '20px',
  letterSpacing: '-0.02em',
};

export const roadmapListStyle: CSSProperties = {
  margin: '22px 0 0',
  padding: 0,
  listStyle: 'none',
  display: 'grid',
  gap: '20px',
};

export const roadmapItemStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '28px minmax(0, 1fr)',
  gap: '13px',
};

export const roadmapIndexStyle: CSSProperties = {
  width: '28px',
  height: '28px',
  borderRadius: '50%',
  display: 'grid',
  placeItems: 'center',
  background: 'rgba(255,255,255,0.18)',
  fontSize: '12px',
  fontWeight: 900,
};

export const roadmapLabelStyle: CSSProperties = {
  display: 'block',
  fontSize: '14px',
  fontWeight: 900,
};

export const roadmapBodyStyle: CSSProperties = {
  display: 'block',
  marginTop: '5px',
  color: 'rgba(255,255,255,0.76)',
  fontSize: '13px',
  lineHeight: 1.6,
};

export const mobileStyle = `
@media (max-width: 760px) {
  .community-hero { grid-template-columns: 1fr !important; }
  .community-content-grid { grid-template-columns: 1fr !important; }
  .community-main { width: 100% !important; padding-top: 32px !important; }
  .community-section-surface { padding: 24px !important; border-radius: 28px !important; }
}
`;
