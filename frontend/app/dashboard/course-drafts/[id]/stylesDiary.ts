import type { CSSProperties } from 'react';

export const sectionStyle = {
  display: 'grid',
  gap: '16px',
  borderRadius: '28px',
  border: '1px solid rgba(196, 163, 112, 0.22)',
  background: 'linear-gradient(rgba(62, 42, 18, 0.52), rgba(62, 42, 18, 0.52)), url("/images/dashboard/explorer-diary-disk.webp")',
  backgroundSize: 'cover',
  backgroundPosition: 'center',
  backgroundRepeat: 'no-repeat',
  backdropFilter: 'blur(14px)',
  padding: '24px',
} as const;

export const diaryShellStyle = {
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1fr)',
  gap: '0',
  alignItems: 'start',
  position: 'relative',
  paddingRight: '52px',
} as const;

export const diaryShellNarrowStyle = {
  gridTemplateColumns: '1fr',
  paddingRight: 0,
} as const;

export const mobileDiaryNavStyle = {
  display: 'none',
  gap: '8px',
  gridTemplateColumns: 'repeat(3, minmax(0, 1fr))',
} as const;

export const mobileDiaryNavVisibleStyle = {
  display: 'grid',
} as const;

export const mobileDiaryNavItemStyle = {
  position: 'relative',
  minHeight: '40px',
  borderRadius: '14px',
  border: '1px solid rgba(199, 176, 141, 0.18)',
  background: 'rgba(255,255,255,0.08)',
  color: '#F4F7FF',
  fontSize: '13px',
  fontWeight: 700,
  cursor: 'pointer',
  overflow: 'hidden',
} as const;

export const mobileDiaryNavItemActiveStyle = {
  background: 'rgba(233, 214, 183, 0.28)',
  border: '1px solid rgba(233, 214, 183, 0.48)',
  color: '#FFF7E8',
} as const;

export const mobileDiaryNavItemLockedStyle = {
  opacity: 0.54,
  cursor: 'not-allowed',
} as const;

export const diaryBookmarkRailStyle = {
  display: 'grid',
  gap: '0',
  position: 'absolute',
  top: '28px',
  right: '-4px',
  zIndex: 10,
  alignContent: 'start',
} as const;

export const diaryBookmarkRailHiddenStyle = {
  display: 'none',
} as const;

export const diaryBookmarkStyle = {
  display: 'grid',
  justifyItems: 'center',
  alignItems: 'center',
  width: '44px',
  minHeight: '112px',
  padding: '14px 8px',
  marginTop: '-8px',
  borderRadius: '0 14px 14px 0',
  border: '1px solid rgba(196, 171, 132, 0.14)',
  borderLeft: 'none',
  background: 'linear-gradient(180deg, rgba(139, 103, 57, 0.28), rgba(109, 78, 42, 0.22))',
  color: '#FCEED9',
  textAlign: 'center',
  cursor: 'pointer',
  boxShadow: '0 12px 24px rgba(0, 0, 0, 0.14)',
} as const;

export const diaryBookmarkActiveStyle = {
  transform: 'translateX(-10px)',
  background: 'linear-gradient(180deg, rgba(227, 202, 162, 0.96), rgba(196, 163, 112, 0.92))',
  color: '#3F2B16',
  boxShadow: '0 18px 28px rgba(0, 0, 0, 0.18)',
} as const;

export const diaryBookmarkLockedStyle = {
  opacity: 0.58,
  cursor: 'not-allowed',
} as const;

export const diaryBookmarkLabelStyle = {
  fontSize: '14px',
  fontWeight: 800,
  lineHeight: 1.15,
  writingMode: 'vertical-rl',
  textOrientation: 'mixed',
  letterSpacing: 0,
} as const;

export const diaryBookmarkMetaStyle = {
  display: 'none',
} as const;

export const diaryLayoutStyle = {
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 7fr) minmax(320px, 3fr)',
  gap: '18px',
  alignItems: 'start',
} as const;

export const diaryLayoutNarrowStyle = {
  gridTemplateColumns: '1fr',
} as const;

export const diaryPageStyle = {
  display: 'grid',
  gap: '18px',
  minHeight: '680px',
  padding: '22px',
  borderRadius: '26px',
  border: '1px solid rgba(201, 179, 143, 0.18)',
  backgroundImage: 'linear-gradient(180deg, rgba(248, 238, 217, 0.28), rgba(229, 212, 184, 0.32)), url("/images/dashboard/diarybook.webp")',
  backgroundSize: 'cover',
  backgroundPosition: 'center',
  backgroundRepeat: 'no-repeat',
  boxShadow: '0 24px 48px rgba(0, 0, 0, 0.2)',
  color: '#2E2418',
} as const;

export const diaryPanelPageStyle = {
  ...diaryPageStyle,
  backgroundImage:
    'linear-gradient(180deg, rgba(248, 238, 217, 0.28), rgba(232, 217, 190, 0.38)), url("/images/dashboard/panel.webp"), url("/images/dashboard/panel.webp"), url("/images/dashboard/clipboard_panel_358x64.webp")',
  backgroundSize: '100% 100%, 100% auto, 100% auto, 100% 64px',
  backgroundPosition: 'center, top center, bottom center, center 72px',
  backgroundRepeat: 'no-repeat, no-repeat, no-repeat, repeat-y',
} as const;

export const diaryPageHeaderStyle = {
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: '12px',
  flexWrap: 'wrap',
} as const;

export const diaryPageEyebrowStyle = {
  fontSize: '11px',
  fontWeight: 700,
  letterSpacing: '0.12em',
  textTransform: 'uppercase',
  color: '#8A6940',
} as const;

export const diaryPageTitleStyle = {
  margin: '4px 0 0',
  fontSize: '28px',
  color: '#2C241C',
} as const;

export const pageStatusPillStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: '30px',
  padding: '0 12px',
  borderRadius: '999px',
  border: '1px solid rgba(125, 97, 55, 0.18)',
  background: 'rgba(255,255,255,0.5)',
  color: '#6D4F2A',
  fontSize: '12px',
  fontWeight: 700,
} as const;

export const plannerMapBoardStyle = {
  position: 'relative',
  minHeight: '460px',
  borderRadius: '24px',
  overflow: 'hidden',
  border: '1px solid rgba(102, 84, 58, 0.12)',
  background: 'radial-gradient(circle at 50% 45%, rgba(57, 80, 120, 0.22), rgba(20, 28, 44, 0.06) 34%, rgba(244, 233, 209, 0.78) 66%), linear-gradient(180deg, rgba(245, 237, 219, 0.98), rgba(230, 216, 187, 0.98))',
} as const;

export const plannerMapGlowStyle = {
  position: 'absolute',
  inset: '10% 12%',
  borderRadius: '999px',
  background: 'radial-gradient(circle, rgba(108, 136, 194, 0.14), rgba(108, 136, 194, 0) 70%)',
  pointerEvents: 'none',
} as const;

export const plannerOrbitStyle = {
  position: 'absolute',
  inset: '18% 16%',
  borderRadius: '50%',
  border: '1px dashed rgba(85, 76, 63, 0.22)',
} as const;

export const plannerPlanetCoreStyle = {
  position: 'absolute',
  left: '50%',
  top: '50%',
  width: '132px',
  height: '132px',
  borderRadius: '999px',
  transform: 'translate(-50%, -50%)',
  background: 'radial-gradient(circle at 36% 32%, #fff2c4 0%, #edc56f 28%, #b7772f 72%, #65401d 100%)',
  boxShadow: '0 18px 42px rgba(99, 63, 25, 0.24)',
  display: 'grid',
  placeItems: 'center',
} as const;

export const plannerPlanetTextStyle = {
  maxWidth: '88px',
  textAlign: 'center',
  fontSize: '13px',
  fontWeight: 800,
  lineHeight: 1.35,
  color: '#2B1C10',
} as const;

export const regionNodePositions = [
  { left: '16%', top: '18%' },
  { left: '72%', top: '16%' },
  { left: '79%', top: '56%' },
  { left: '23%', top: '70%' },
  { left: '49%', top: '10%' },
  { left: '10%', top: '46%' },
] as const;

export const regionNodeStyle = {
  position: 'absolute',
  width: '148px',
  padding: '12px 14px',
  borderRadius: '18px',
  border: '1px solid rgba(110, 88, 58, 0.18)',
  background: 'rgba(255, 250, 240, 0.86)',
  boxShadow: '0 16px 28px rgba(79, 61, 34, 0.12)',
  transform: 'translate(-50%, -50%)',
  textAlign: 'left',
  cursor: 'pointer',
  fontFamily: 'inherit',
  transition: 'transform 0.18s ease, box-shadow 0.18s ease, border-color 0.18s ease',
} as const;

export const regionNodeSelectedStyle = {
  border: '2px solid rgba(69, 93, 153, 0.42)',
  background: 'rgba(255, 252, 244, 0.96)',
  boxShadow: '0 20px 34px rgba(58, 76, 132, 0.18)',
  transform: 'translate(-50%, -50%) scale(1.04)',
} as const;

export const regionNodeBadgeStyle = {
  fontSize: '11px',
  fontWeight: 700,
  color: '#8B6840',
  textTransform: 'uppercase',
  letterSpacing: '0.08em',
} as const;

export const regionNodeTitleStyle = {
  marginTop: '6px',
  fontSize: '15px',
  fontWeight: 800,
  color: '#2F2418',
  lineHeight: 1.4,
} as const;

export const regionNodeMetaStyle = {
  marginTop: '4px',
  fontSize: '12px',
  color: '#6D5942',
} as const;

export const plannerStatsGridStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(120px, 1fr))',
  gap: '10px',
} as const;

export const plannerStatCardStyle = {
  display: 'grid',
  gap: '4px',
  padding: '14px',
  borderRadius: '16px',
  border: '1px solid rgba(107, 86, 60, 0.12)',
  background: 'rgba(255,255,255,0.46)',
} as const;

export const plannerStatLabelStyle = {
  fontSize: '11px',
  textTransform: 'uppercase',
  letterSpacing: '0.08em',
  color: '#8A6D48',
} as const;

export const plannerStatValueStyle = {
  fontSize: '20px',
  color: '#322619',
} as const;

export const diaryRightStackStyle = {
  display: 'grid',
  gap: '16px',
} as const;

export const plannerPanelStyle = {
  display: 'grid',
  gap: '14px',
  padding: '18px',
  borderRadius: '20px',
  border: '1px solid rgba(107, 86, 60, 0.18)',
  background: 'transparent',
} as const;

export const treeRegionButtonStyle = {
  display: 'grid',
  gap: '6px',
  padding: 0,
  border: 'none',
  background: 'transparent',
  textAlign: 'left',
  cursor: 'pointer',
  fontFamily: 'inherit',
} as const;

export const structureEditableHeaderStyle = {
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1fr) auto',
  gap: '10px',
  alignItems: 'start',
} as const;

export const structureLessonHeaderStyle = {
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1fr) auto',
  gap: '10px',
  alignItems: 'start',
} as const;

export const structureMoveButtonGroupStyle = {
  display: 'flex',
  gap: '6px',
  flexWrap: 'wrap',
  justifyContent: 'flex-end',
} as const;

export const structureMoveButtonStyle = {
  minHeight: '28px',
  padding: '0 9px',
  borderRadius: '999px',
  border: '1px solid rgba(94, 234, 212, 0.28)',
  background: 'rgba(20, 184, 166, 0.14)',
  color: '#DDFCF8',
  fontSize: '11px',
  fontWeight: 800,
  cursor: 'pointer',
  fontFamily: 'inherit',
} as const;

export const structureMoveButtonDisabledStyle = {
  opacity: 0.4,
  cursor: 'not-allowed',
} as const;

export const panelHeaderStyle = {
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: '10px',
  flexWrap: 'wrap',
} as const;

export const panelEyebrowStyle = {
  fontSize: '11px',
  fontWeight: 700,
  letterSpacing: '0.12em',
  textTransform: 'uppercase',
  color: '#8A6940',
} as const;

export const panelTitleStyle = {
  margin: '4px 0 0',
  fontSize: '20px',
  color: '#2E2418',
} as const;

export const panelHintStyle = {
  fontSize: '12px',
  color: '#6B5843',
} as const;

export const plannerTreeStyle = {
  display: 'grid',
  gap: '12px',
  maxHeight: '420px',
  overflowY: 'auto',
  padding: '12px',
} as const;
