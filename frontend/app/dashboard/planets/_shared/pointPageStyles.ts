// Style objects for PlanetPointPageClient

export const pageStyle = {
  position: 'relative',
  minHeight: '100vh',
  overflowX: 'clip',
} as const;

export const pageBackgroundOverlayStyle = {
  position: 'absolute',
  inset: 0,
  background:
    'radial-gradient(circle at top, rgba(58, 95, 176, 0.14), transparent 32%), linear-gradient(rgba(7, 14, 26, 0.88), rgba(7, 14, 26, 0.96))',
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

export const mainStyle = {
  position: 'relative',
  maxWidth: '1360px',
  margin: '0 auto',
  padding: '32px 24px 56px',
  display: 'grid',
  gap: '22px',
} as const;

export const heroCardStyle = {
  display: 'grid',
  gap: '18px',
  borderRadius: '18px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  borderLeft: '2px solid rgba(50,200,100,0.35)',
  background: 'rgba(7, 18, 34, 0.58)',
  boxShadow: 'none',
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
  fontSize: '34px',
  lineHeight: 1.18,
  color: '#F7FAFF',
} as const;

export const descriptionStyle = {
  margin: 0,
  fontSize: '17px',
  lineHeight: 1.75,
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
  borderRadius: '14px',
  border: '1px solid rgba(194, 210, 245, 0.12)',
  background: 'rgba(255,255,255,0.06)',
} as const;

export const metaLabelStyle = {
  fontSize: '13px',
  color: '#A9B6D2',
} as const;

export const metaValueStyle = {
  fontSize: '16px',
  fontWeight: 700,
  color: '#F7FAFF',
} as const;

export const sectionStyle = {
  display: 'grid',
  gap: '16px',
  borderRadius: '18px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(7, 18, 34, 0.56)',
  padding: '24px',
} as const;

export const sectionTitleStyle = {
  margin: 0,
  fontSize: '26px',
  lineHeight: 1.24,
  color: '#F7FAFF',
} as const;

export const sectionSubtitleStyle = {
  margin: '6px 0 0',
  fontSize: '15px',
  color: '#C8D1E8',
  lineHeight: 1.65,
} as const;

export const secondaryButtonStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '42px',
  padding: '0 18px',
  borderRadius: '999px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(255,255,255,0.08)',
  color: '#F4F7FF',
  textDecoration: 'none',
  fontSize: '15px',
  fontWeight: 700,
} as const;

export const primaryButtonStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '42px',
  padding: '0 18px',
  borderRadius: '999px',
  border: '1px solid rgba(98, 183, 255, 0.26)',
  background: 'linear-gradient(135deg, rgba(57, 208, 255, 0.22), rgba(79, 141, 255, 0.22))',
  color: '#F4F7FF',
  textDecoration: 'none',
  fontSize: '15px',
  fontWeight: 800,
} as const;

export const pointNavRowStyle = {
  display: 'flex',
  gap: '10px',
  flexWrap: 'wrap',
  alignItems: 'center',
} as const;

export const pointNavDisabledStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: '42px',
  padding: '0 16px',
  borderRadius: '999px',
  color: '#8EA4CB',
  background: 'rgba(255,255,255,0.05)',
  border: '1px solid rgba(194, 210, 245, 0.08)',
  fontSize: '14px',
} as const;

export const pointTabListStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))',
  gap: '0',
  overflow: 'hidden',
  borderRadius: '14px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
} as const;

export const pointTabButtonStyle = {
  minHeight: '64px',
  display: 'grid',
  gap: '4px',
  alignContent: 'center',
  justifyItems: 'start',
  padding: '12px 16px',
  borderWidth: '0 1px 0 0',
  borderStyle: 'solid',
  fontFamily: 'inherit',
  fontSize: '15px',
  textAlign: 'left' as const,
  cursor: 'pointer',
} as const;

export const completionToggleStyle = {
  width: '100%',
  minHeight: '56px',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '16px',
  padding: '0 18px',
  borderRadius: '14px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  fontFamily: 'inherit',
  fontSize: '16px',
  fontWeight: 800,
  cursor: 'pointer',
} as const;

export const pointShellCardStyle = {
  display: 'grid',
  gap: '16px',
  padding: '24px',
  borderRadius: '16px',
  border: '1px solid rgba(194, 210, 245, 0.12)',
  background: 'rgba(255,255,255,0.05)',
} as const;

export const sectionMiniHeaderStyle = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '12px',
  flexWrap: 'wrap',
} as const;

export const runtimeButtonRowStyle = {
  display: 'flex',
  gap: '10px',
  flexWrap: 'wrap',
} as const;

export const runtimePrimaryButtonStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '48px',
  padding: '0 22px',
  borderRadius: '999px',
  border: '1px solid rgba(98, 183, 255, 0.26)',
  background: 'linear-gradient(135deg, rgba(57, 208, 255, 0.22), rgba(79, 141, 255, 0.22))',
  color: '#F4F7FF',
  fontSize: '16px',
  fontWeight: 800,
  cursor: 'pointer',
} as const;

export const runtimeSecondaryButtonStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '46px',
  padding: '0 20px',
  borderRadius: '999px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(255,255,255,0.08)',
  color: '#F4F7FF',
  fontSize: '15px',
  fontWeight: 700,
  cursor: 'pointer',
} as const;

export const runtimeButtonDisabledStyle = {
  opacity: 0.72,
  cursor: 'not-allowed',
} as const;

export const pointShellTitleStyle = {
  color: '#F7FAFF',
  fontSize: '19px',
  lineHeight: 1.35,
} as const;

export const pointShellMessageStyle = {
  margin: 0,
  color: '#D6E1F8',
  lineHeight: 1.75,
  fontSize: '17px',
} as const;

export const placeholderCardStyle = {
  padding: '20px',
  borderRadius: '14px',
  border: '1px dashed rgba(194, 210, 245, 0.2)',
  background: 'rgba(255,255,255,0.04)',
  color: '#C8D1E8',
  fontSize: '15px',
  lineHeight: 1.65,
} as const;

export const emptyCardStyle = {
  padding: '18px',
  borderRadius: '14px',
  border: '1px solid rgba(194, 210, 245, 0.12)',
  background: 'rgba(255,255,255,0.05)',
  color: '#C8D1E8',
  fontSize: '15px',
  lineHeight: 1.6,
} as const;

export const noticeCardStyle = {
  padding: '16px 18px',
  borderRadius: '14px',
  background: 'rgba(8, 20, 36, 0.78)',
  border: '1px solid rgba(141, 198, 255, 0.16)',
  color: '#DCEEFF',
  fontSize: '15px',
  lineHeight: 1.65,
} as const;

export const aiSummarySourceStyle = {
  fontSize: '14px',
  fontWeight: 700,
  color: '#BFD0F0',
} as const;

export const aiSummaryTextStyle = {
  color: '#E7EEFF',
  fontSize: '17px',
  lineHeight: 1.75,
  whiteSpace: 'pre-wrap' as const,
} as const;

export const formGridStyle = {
  display: 'grid',
  gap: '14px',
} as const;

export const inputRowStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))',
  gap: '12px',
} as const;

export const inputStyle = {
  width: '100%',
  minHeight: '46px',
  padding: '0 15px',
  borderRadius: '12px',
  border: '1px solid rgba(194, 210, 245, 0.16)',
  background: 'rgba(255,255,255,0.06)',
  color: '#F4F7FF',
  fontSize: '16px',
} as const;

export const textareaStyle = {
  width: '100%',
  minHeight: '110px',
  padding: '14px 16px',
  borderRadius: '12px',
  border: '1px solid rgba(194, 210, 245, 0.16)',
  background: 'rgba(255,255,255,0.06)',
  color: '#F4F7FF',
  fontSize: '16px',
  resize: 'vertical' as const,
  lineHeight: 1.75,
} as const;

export const formActionRowStyle = {
  display: 'flex',
  justifyContent: 'flex-end',
  gap: '10px',
  flexWrap: 'wrap',
} as const;

export const blockCardStyle = {
  display: 'grid',
  gap: '14px',
  padding: '20px',
  borderRadius: '14px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(255,255,255,0.05)',
} as const;

export const blockHeaderStyle = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '12px',
} as const;

export const blockMetaStyle = {
  fontSize: '13px',
  color: '#A9B6D2',
} as const;

export const blockActionRowStyle = {
  display: 'flex',
  justifyContent: 'space-between',
  gap: '10px',
  flexWrap: 'wrap',
} as const;

export const dangerButtonStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '42px',
  padding: '0 18px',
  borderRadius: '999px',
  border: '1px solid rgba(255, 126, 126, 0.26)',
  background: 'rgba(120, 20, 26, 0.36)',
  color: '#FFE4E7',
  fontSize: '15px',
  fontWeight: 700,
  cursor: 'pointer',
} as const;

export const readingPanelStyle = {
  display: 'grid',
  gap: '18px',
  padding: '28px',
  borderRadius: '16px',
  border: '1px solid rgba(194, 210, 245, 0.12)',
  background: 'rgba(255,255,255,0.06)',
} as const;

export const readingBodyStyle = {
  margin: 0,
  fontSize: '18px',
  lineHeight: 1.78,
  whiteSpace: 'pre-wrap' as const,
} as const;

export const disclosureStyle = {
  borderRadius: '14px',
  border: '1px solid rgba(194, 210, 245, 0.12)',
  padding: '16px 18px',
} as const;

export const disclosureSummaryStyle = {
  cursor: 'pointer',
  fontSize: '16px',
  fontWeight: 800,
} as const;
