import type { CSSProperties } from 'react';

// ── 페이지 배경 ──────────────────────────────────────────────
export const pageStyle = {
  position: 'relative',
  minHeight: '100vh',
  overflow: 'hidden',
} satisfies CSSProperties;

export const pageBackgroundVideoStyle = {
  position: 'absolute',
  inset: 0,
  width: '100%',
  height: '100%',
  objectFit: 'cover',
  objectPosition: 'center',
  opacity: 0.22,
  zIndex: 0,
  pointerEvents: 'none',
} satisfies CSSProperties;

export const pageBackgroundImageStyle = {
  position: 'absolute',
  inset: 0,
  backgroundImage: "url('/images/PlanetMap_background.webp')",
  backgroundPosition: 'center',
  backgroundSize: 'cover',
  backgroundRepeat: 'no-repeat',
  zIndex: 1,
  opacity: 0.72,
  pointerEvents: 'none',
} satisfies CSSProperties;

export const pageBackgroundOverlayStyle = {
  position: 'absolute',
  inset: 0,
  background: 'linear-gradient(rgba(7, 14, 26, 0.3), rgba(7, 14, 26, 0.58))',
  zIndex: 2,
  pointerEvents: 'none',
} satisfies CSSProperties;

export const starOverlayStyle = {
  position: 'absolute',
  inset: 0,
  width: '100%',
  height: '100%',
  pointerEvents: 'none',
  opacity: 0.8,
} satisfies CSSProperties;

// ── 크롬 / 레이아웃 ─────────────────────────────────────────
export const pageChromeStyle = {
  position: 'relative',
  zIndex: 4,
} satisfies CSSProperties;

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
} satisfies CSSProperties;

export const mainStyle = {
  position: 'relative',
  maxWidth: '1360px',
  margin: '0 auto',
  padding: '28px 24px 40px',
  display: 'grid',
  gap: '18px',
} satisfies CSSProperties;

export const headerBlockStyle = {
  display: 'grid',
  gap: '6px',
} satisfies CSSProperties;

export const heroTitleStyle = {
  margin: 0,
  fontSize: '36px',
  color: '#F7FAFF',
} satisfies CSSProperties;

// ── CTA 섹션 ─────────────────────────────────────────────────
export const creatorSectionStyle = {
  borderRadius: '24px',
  border: '1px solid rgba(214, 226, 255, 0.28)',
  borderLeft: '2px solid rgba(92, 220, 132, 0.62)',
  background: 'rgba(5, 12, 26, 0.82)',
  boxShadow: '0 24px 80px rgba(0, 0, 0, 0.34)',
  backdropFilter: 'blur(18px)',
  padding: '14px 18px',
} satisfies CSSProperties;

export const creatorSectionInnerStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '1fr',
  gap: '10px',
  alignItems: 'start',
};

export const creatorHeroSloganStyle = {
  fontSize: '22px',
  fontWeight: 700,
  letterSpacing: '0.01em',
  color: '#FCE9D8',
} satisfies CSSProperties;

export const creatorLabelStyle = {
  fontSize: '18px',
  fontWeight: 700,
  color: '#F8FAFF',
} satisfies CSSProperties;

export const creatorDescriptionStyle = {
  fontSize: '13px',
  lineHeight: 1.45,
  color: '#C8D1E8',
  maxWidth: '620px',
} satisfies CSSProperties;

export const creatorFormStyle = {
  display: 'flex',
  gap: '10px',
  alignItems: 'center',
  flexWrap: 'wrap',
} satisfies CSSProperties;

export const creatorInputStyle = {
  flex: '1 1 420px',
  minHeight: '44px',
  padding: '0 14px',
  borderRadius: '14px',
  border: '1px solid rgba(194, 210, 245, 0.18)',
  background: 'rgba(255,255,255,0.12)',
  color: '#F4F7FF',
  fontSize: '14px',
  outline: 'none',
  fontFamily: 'inherit',
} satisfies CSSProperties;

export function creatorButtonStyle(isCreatingCourse: boolean): CSSProperties {
  return {
    minHeight: '44px',
    padding: '0 16px',
    borderRadius: '14px',
    border: '1px solid #D48314',
    background: isCreatingCourse
      ? 'linear-gradient(135deg, #D08A21 0%, #B96F12 100%)'
      : 'linear-gradient(135deg, #F3A533 0%, #D98619 100%)',
    color: '#FFFFFF',
    fontSize: '13px',
    fontWeight: 700,
    cursor: isCreatingCourse ? 'progress' : 'pointer',
    fontFamily: 'inherit',
    minWidth: '160px',
  };
}

export const creatorErrorStyle = {
  fontSize: '13px',
  color: '#FFD3CC',
} satisfies CSSProperties;

export const queryInputTooltipStyle = {
  position: 'absolute',
  top: 'calc(100% + 6px)',
  left: 0,
  display: 'flex',
  alignItems: 'center',
  gap: '10px',
  padding: '8px 14px 8px 10px',
  borderRadius: '12px',
  border: '1px solid rgba(126, 163, 255, 0.24)',
  background: 'rgba(8, 18, 40, 0.92)',
  backdropFilter: 'blur(12px)',
  boxShadow: '0 4px 20px rgba(0,0,0,0.3)',
  zIndex: 20,
  pointerEvents: 'none',
  animation: 'tooltipFadeIn 0.15s ease',
} satisfies CSSProperties;

export const queryInputTooltipTextStyle = {
  fontSize: '13px',
  fontWeight: 500,
  color: 'rgba(200, 220, 255, 0.9)',
  whiteSpace: 'nowrap',
} satisfies CSSProperties;

export const chipRowStyle = {
  display: 'flex',
  gap: '8px',
  flexWrap: 'wrap',
} satisfies CSSProperties;

export const chipStyle = {
  minHeight: '34px',
  padding: '0 12px',
  borderRadius: '999px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(255,255,255,0.09)',
  color: '#E8EEFF',
  fontSize: '13px',
  cursor: 'pointer',
  fontFamily: 'inherit',
} satisfies CSSProperties;

// ── 맵 / 리스트 레이아웃 ────────────────────────────────────
export function mapListGridStyle(
  isPhoneLayout: boolean,
  _isShortViewport: boolean,
  hideSpaceMap = false,
): CSSProperties {
  return {
    display: 'flex',
    flexDirection: isPhoneLayout ? 'column' : 'row',
    alignItems: 'stretch',
    width: '100%',
    gap: '24px',
    ...(hideSpaceMap ? { gap: '0' } : null),
  };
}

export const mapColumnStyle = {
  flex: 1,
  minWidth: 0,
} satisfies CSSProperties;

export const listColumnStyle = {
  width: '260px',
  height: '690px',
  flexShrink: 0,
  display: 'flex',
  flexDirection: 'column',
} satisfies CSSProperties;

export function responsiveListColumnStyle(hideSpaceMap: boolean): CSSProperties {
  if (!hideSpaceMap) return listColumnStyle;
  return {
    ...listColumnStyle,
    width: '100%',
    maxWidth: '100%',
    flex: '1 1 auto',
  };
}

// Lumi 패널 (shouldRenderDockedLumi 활성화 시 사용)
export const desktopLumiPanelStyle = {
  position: 'absolute',
  bottom: '72px',
  left: '18px',
  width: '260px',
  zIndex: 9,
} satisfies CSSProperties;

// ── 로딩 ────────────────────────────────────────────────────
export const loadingPageStyle = {
  minHeight: '100vh',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  background: 'linear-gradient(rgba(7, 14, 26, 0.76), rgba(7, 14, 26, 0.88))',
  padding: '24px',
} satisfies CSSProperties;

export const loadingCardStyle = {
  width: '100%',
  maxWidth: '420px',
  padding: '32px',
  borderRadius: '20px',
  background: 'rgba(8, 18, 33, 0.7)',
  border: '1px solid rgba(194, 210, 245, 0.16)',
  textAlign: 'center',
  color: '#F4F7FF',
  fontSize: '16px',
  backdropFilter: 'blur(10px)',
} satisfies CSSProperties;


export const heroSubtitleStyle = {
  margin: '4px 0 0',
  maxWidth: '760px',
  color: 'rgba(221, 232, 255, 0.78)',
  fontSize: '15px',
  lineHeight: 1.55,
  fontWeight: 500,
} satisfies CSSProperties;
