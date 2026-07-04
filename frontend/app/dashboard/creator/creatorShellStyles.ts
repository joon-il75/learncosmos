import type { CSSProperties } from 'react';

export const pageStyle: CSSProperties = {
  minHeight: '100vh',
  background: 'linear-gradient(140deg, #FFF7E8 0%, #EEF8F3 46%, #E8F1FF 100%)',
  color: '#182C3A',
};

export const headerStyle: CSSProperties = {
  background: 'rgba(255, 250, 240, 0.74)',
  borderBottom: '1px solid rgba(122, 83, 28, 0.1)',
};

export const headerInnerStyle: CSSProperties = { maxWidth: '1120px' };

export const headerMetaStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: '2px',
  minWidth: 0,
};

export const headerEyebrowStyle: CSSProperties = {
  fontSize: '11px',
  fontWeight: 900,
  letterSpacing: '0.12em',
  textTransform: 'uppercase',
  color: '#8A672D',
};

export const headerTitleStyle: CSSProperties = {
  fontSize: '15px',
  fontWeight: 900,
  color: '#263B4A',
};

export const mainStyle: CSSProperties = {
  width: 'min(1120px, calc(100% - 32px))',
  margin: '0 auto',
  padding: '112px 0 64px',
};

export const heroStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 0.95fr) minmax(280px, 0.75fr)',
  gap: '34px',
  alignItems: 'stretch',
};

export const heroTextStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'center',
  gap: '16px',
};

export const eyebrowStyle: CSSProperties = {
  width: 'fit-content',
  padding: '8px 13px',
  borderRadius: '999px',
  background: 'rgba(176, 120, 42, 0.11)',
  color: '#926B27',
  fontSize: '12px',
  fontWeight: 900,
  letterSpacing: '0.08em',
  textTransform: 'uppercase',
};

export const titleStyle: CSSProperties = {
  margin: 0,
  fontSize: 'clamp(34px, 7vw, 72px)',
  lineHeight: 0.97,
  letterSpacing: '-0.06em',
  color: '#1B3244',
};

export const subtitleStyle: CSSProperties = {
  margin: 0,
  maxWidth: '700px',
  fontSize: 'clamp(16px, 2.2vw, 20px)',
  lineHeight: 1.75,
  color: '#5B6870',
};

export const canvasStyle: CSSProperties = {
  minHeight: '260px',
  borderRadius: '40px',
  padding: '28px',
  background: 'radial-gradient(circle at 25% 25%, rgba(255, 194, 93, 0.3), transparent 32%), radial-gradient(circle at 78% 68%, rgba(98, 171, 164, 0.28), transparent 34%), rgba(255,255,255,0.58)',
  boxShadow: '0 30px 86px rgba(76, 64, 47, 0.13)',
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'space-between',
};

export const statusStyle: CSSProperties = {
  width: 'fit-content',
  padding: '9px 13px',
  borderRadius: '999px',
  background: 'rgba(255,255,255,0.72)',
  color: '#6E552B',
  fontSize: '13px',
  fontWeight: 900,
};

export const studioLinesStyle: CSSProperties = {
  height: '132px',
  borderRadius: '28px',
  background: 'linear-gradient(135deg, rgba(27, 50, 68, 0.13), rgba(255,255,255,0.34))',
  position: 'relative',
  overflow: 'hidden',
};

export const studioDotStyle: CSSProperties = {
  position: 'absolute',
  width: '48px',
  height: '48px',
  borderRadius: '18px',
  left: '26px',
  top: '28px',
  background: '#1B3244',
  boxShadow: '104px 28px 0 #E7A84E, 220px -2px 0 #73B7AF',
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
  color: active ? '#FFFFFF' : '#5B4B31',
  background: active ? '#1B3244' : 'rgba(255,255,255,0.58)',
  fontSize: '14px',
  fontWeight: 900,
  boxShadow: active ? '0 14px 34px rgba(27, 50, 68, 0.18)' : 'none',
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
  background: 'rgba(255,255,255,0.64)',
  boxShadow: '0 26px 74px rgba(76, 64, 47, 0.11)',
};

export const sectionEyebrowStyle: CSSProperties = {
  margin: 0,
  fontSize: '12px',
  fontWeight: 900,
  color: '#A5752F',
  letterSpacing: '0.12em',
  textTransform: 'uppercase',
};

export const sectionTitleStyle: CSSProperties = {
  margin: '12px 0 0',
  fontSize: 'clamp(27px, 4vw, 44px)',
  lineHeight: 1.08,
  letterSpacing: '-0.04em',
  color: '#1B3244',
};

export const sectionDescriptionStyle: CSSProperties = {
  margin: '18px 0 0',
  fontSize: '16px',
  lineHeight: 1.78,
  color: '#5B6870',
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
  background: '#1B3244',
  color: '#FFFFFF',
  fontSize: '14px',
  fontWeight: 900,
};

export const secondaryActionStyle: CSSProperties = {
  padding: '13px 18px',
  borderRadius: '999px',
  background: 'rgba(176, 120, 42, 0.11)',
  color: '#6E552B',
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
  color: '#4F606A',
  fontSize: '14px',
  lineHeight: 1.65,
};

export const noteDotStyle: CSSProperties = {
  width: '8px',
  height: '8px',
  borderRadius: '50%',
  marginTop: '8px',
  background: '#E7A84E',
};

export const readinessStyle: CSSProperties = {
  borderRadius: '34px',
  padding: '28px',
  background: 'linear-gradient(180deg, rgba(27, 50, 68, 0.94), rgba(69, 87, 96, 0.84))',
  color: '#FFFFFF',
  boxShadow: '0 24px 70px rgba(27, 50, 68, 0.16)',
};

export const readinessTitleStyle: CSSProperties = {
  margin: 0,
  fontSize: '20px',
  letterSpacing: '-0.02em',
};

export const readinessSubtitleStyle: CSSProperties = {
  margin: '10px 0 0',
  color: 'rgba(255,255,255,0.74)',
  fontSize: '13px',
  lineHeight: 1.62,
};

export const readinessListStyle: CSSProperties = {
  margin: '22px 0 0',
  padding: 0,
  listStyle: 'none',
  display: 'grid',
  gap: '18px',
};

export const readinessItemStyle: CSSProperties = {
  paddingTop: '18px',
  borderTop: '1px solid rgba(255,255,255,0.14)',
};

export const readinessLabelStyle: CSSProperties = {
  display: 'block',
  fontSize: '14px',
  fontWeight: 900,
};

export const readinessBodyStyle: CSSProperties = {
  display: 'block',
  marginTop: '5px',
  color: 'rgba(255,255,255,0.76)',
  fontSize: '13px',
  lineHeight: 1.6,
};

export const mobileStyle = `
@media (max-width: 760px) {
  .creator-hero { grid-template-columns: 1fr !important; }
  .creator-content-grid { grid-template-columns: 1fr !important; }
  .creator-main { width: 100% !important; padding-top: 32px !important; }
  .creator-section-surface { padding: 24px !important; border-radius: 28px !important; }
}
`;
