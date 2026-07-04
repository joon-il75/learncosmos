import { type CSSProperties } from 'react';
import { SUPER_ADMIN_PAGE_WIDTH } from '@/components/super-admin/layout';

export const pageStyle: CSSProperties = {
  position: 'relative',
  minHeight: '100vh',
  overflow: 'hidden',
  background: 'radial-gradient(circle at top, rgba(32, 67, 119, 0.28), transparent 42%), linear-gradient(180deg, #07111f, #0b1321 42%, #060b12)',
};

export const backdropStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  backgroundImage: "url('/images/PlanetMap_background.webp')",
  backgroundSize: 'cover',
  backgroundPosition: 'center',
  opacity: 0.12,
  filter: 'saturate(0.9)',
  pointerEvents: 'none',
};

export const shellStyle: CSSProperties = {
  position: 'relative',
  zIndex: 10,
  width: SUPER_ADMIN_PAGE_WIDTH,
  margin: '0 auto',
  padding: '24px 0 48px',
};

export const topChromeStyle: CSSProperties = {
  position: 'relative',
  zIndex: 120,
  isolation: 'isolate',
  pointerEvents: 'auto',
  display: 'grid',
  gap: '14px',
  marginBottom: '10px',
};

export const introCardStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 18,
  flexWrap: 'wrap',
  marginTop: 24,
  padding: 24,
  borderRadius: 28,
  border: '1px solid rgba(117, 157, 234, 0.22)',
  background: 'linear-gradient(180deg, rgba(17, 27, 45, 0.9), rgba(9, 15, 26, 0.94))',
  boxShadow: '0 28px 68px rgba(0, 0, 0, 0.24)',
};

export const tabSectionStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 12,
  marginTop: 22,
};

export const tabGridStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(5, minmax(0, 1fr))',
  gap: 12,
};

export const tabButtonStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'flex-start',
  gap: 6,
  padding: '16px 18px',
  borderRadius: 22,
  border: '1px solid rgba(120, 158, 230, 0.18)',
  background: 'linear-gradient(180deg, rgba(12, 20, 33, 0.94), rgba(8, 13, 23, 0.96))',
  color: '#eff6ff',
  cursor: 'pointer',
  textAlign: 'left',
};

export const tabButtonActiveStyle: CSSProperties = {
  borderColor: 'rgba(159, 198, 255, 0.42)',
  background: 'linear-gradient(180deg, rgba(35, 56, 90, 0.96), rgba(12, 20, 34, 0.98))',
  boxShadow: '0 18px 32px rgba(0, 0, 0, 0.18)',
};

export const tabLabelStyle: CSSProperties = {
  fontSize: 15,
  fontWeight: 700,
};

export const tabDescriptionStyle: CSSProperties = {
  color: 'rgba(198, 214, 238, 0.76)',
  fontSize: 12,
  lineHeight: 1.45,
};

export const tabMetaStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: 10,
  flexWrap: 'wrap',
};

export const tabMetaCopyStyle: CSSProperties = {
  color: 'rgba(198, 214, 238, 0.76)',
  fontSize: 13,
};

export const runtimeConfigPanelStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 16,
  flexWrap: 'wrap',
  marginTop: 18,
  padding: 18,
  borderRadius: 22,
  border: '1px solid rgba(120, 158, 230, 0.18)',
  background: 'linear-gradient(180deg, rgba(14, 23, 39, 0.95), rgba(8, 13, 23, 0.96))',
};

export const runtimeConfigMetaStyle: CSSProperties = {
  display: 'grid',
  gap: 6,
};

export const draftToolbarStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1fr) 140px',
  gap: 10,
  marginBottom: 12,
};

export const diffListStyle: CSSProperties = {
  display: 'grid',
  gap: 6,
};

export const diffLineStyle: CSSProperties = {
  display: 'block',
  padding: '8px 10px',
  borderRadius: 10,
  background: 'rgba(9, 17, 28, 0.75)',
  border: '1px solid rgba(120, 158, 230, 0.12)',
  color: 'rgba(214, 230, 255, 0.88)',
  fontSize: 12,
  whiteSpace: 'pre-wrap',
};

export const runtimeConfigStatusStyle: CSSProperties = {
  color: 'rgba(180, 210, 255, 0.88)',
  fontSize: 12,
};

export const runtimeConfigDirtyStyle: CSSProperties = {
  color: '#f8d477',
  fontSize: 12,
  fontWeight: 600,
};

export const assetPanelStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(2, minmax(0, 1fr))',
  gap: 22,
  marginTop: 22,
};

export const assetColumnStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 14,
  padding: 22,
  borderRadius: 28,
  border: '1px solid rgba(120, 158, 230, 0.22)',
  background: 'linear-gradient(180deg, rgba(14, 23, 39, 0.95), rgba(8, 13, 23, 0.96))',
};

export const assetMetaCardStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '160px minmax(0, 1fr)',
  gap: 16,
  alignItems: 'start',
};

export const assetPreviewWrapStyle: CSSProperties = {
  borderRadius: 20,
  border: '1px solid rgba(120, 158, 230, 0.18)',
  background: 'rgba(6, 11, 19, 0.92)',
  padding: 10,
};

export const assetPreviewStyle: CSSProperties = {
  width: '100%',
  aspectRatio: '5 / 3',
  objectFit: 'cover',
  borderRadius: 12,
};

export const assetMetaTextStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 8,
  color: 'rgba(213, 225, 242, 0.84)',
  fontSize: 13,
  lineHeight: 1.5,
};

export const gridSpecStyle: CSSProperties = {
  padding: '6px 10px',
  borderRadius: 8,
  background: 'rgba(95, 131, 255, 0.1)',
  border: '1px solid rgba(95, 131, 255, 0.22)',
  fontSize: 12,
  fontFamily: 'monospace',
  color: 'rgba(180, 210, 255, 0.9)',
};

export const uploadBoxStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 10,
  padding: 16,
  borderRadius: 20,
  border: '1px solid rgba(120, 158, 230, 0.18)',
  background: 'rgba(13, 21, 35, 0.72)',
};

export const fileInputStyle: CSSProperties = {
  color: '#eff6ff',
  fontSize: 13,
};

export const uploadActionRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 12,
  flexWrap: 'wrap',
};

export const uploadHintStyle: CSSProperties = {
  color: 'rgba(190, 208, 235, 0.8)',
  fontSize: 12,
};

export const assetStatusStyle: CSSProperties = {
  color: '#a8d8ff',
  fontSize: 12,
};

export const promptCardStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 10,
  padding: 16,
  borderRadius: 20,
  border: '1px solid rgba(120, 158, 230, 0.18)',
  background: 'rgba(13, 21, 35, 0.72)',
};

export const promptHeaderStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 12,
};

export const promptTitleStyle: CSSProperties = {
  color: '#f4faff',
  fontSize: 14,
};

export const textInputStyle: CSSProperties = {
  borderRadius: 14,
  border: '1px solid rgba(124, 156, 214, 0.24)',
  background: 'rgba(13, 21, 35, 0.96)',
  color: '#eff6ff',
  padding: '12px 14px',
  fontSize: 14,
};

export const secondaryActionButtonStyle: CSSProperties = {
  border: '1px solid rgba(120, 158, 230, 0.24)',
  borderRadius: 999,
  background: 'rgba(24, 37, 60, 0.84)',
  color: '#ebf5ff',
  padding: '10px 14px',
  fontSize: 13,
  cursor: 'pointer',
};

export const copyButtonStyle: CSSProperties = {
  border: '1px solid rgba(120, 158, 230, 0.24)',
  borderRadius: 999,
  background: 'rgba(24, 37, 60, 0.84)',
  color: '#ebf5ff',
  padding: '8px 12px',
  fontSize: 12,
  cursor: 'pointer',
};

export const promptTextareaStyle: CSSProperties = {
  minHeight: 140,
  resize: 'vertical',
  borderRadius: 16,
  border: '1px solid rgba(124, 156, 214, 0.24)',
  background: 'rgba(7, 12, 20, 0.96)',
  color: '#eff6ff',
  padding: '12px 14px',
  fontSize: 13,
  lineHeight: 1.55,
};

export const promptFootnoteStyle: CSSProperties = {
  color: 'rgba(190, 208, 235, 0.78)',
  fontSize: 12,
};

export const authGateStyle: CSSProperties = {
  minHeight: '100vh',
  display: 'grid',
  placeItems: 'center',
  background: '#08111f',
};

export const authCardStyle: CSSProperties = {
  padding: '20px 24px',
  borderRadius: 18,
  border: '1px solid rgba(120, 158, 230, 0.2)',
  background: 'rgba(12, 21, 36, 0.94)',
  color: '#eff6ff',
};

export const introTextStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 8,
  maxWidth: 760,
};

export const introTitleStyle: CSSProperties = {
  color: '#f7fbff',
  fontSize: 24,
  fontWeight: 700,
};

export const introCopyStyle: CSSProperties = {
  margin: 0,
  color: 'rgba(205, 220, 242, 0.88)',
  lineHeight: 1.6,
  fontSize: 15,
};

export const primaryButtonStyle: CSSProperties = {
  border: 0,
  borderRadius: 999,
  background: 'linear-gradient(135deg, #e9f3ff, #a8d8ff)',
  color: '#08111f',
  padding: '12px 18px',
  fontWeight: 700,
  cursor: 'pointer',
};

export const rulesLayoutStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '300px minmax(0, 1.25fr) minmax(280px, 0.95fr)',
  gap: 18,
  marginTop: 22,
  alignItems: 'start',
};

export const ruleListCardStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 14,
  padding: 22,
  borderRadius: 28,
  border: '1px solid rgba(120, 158, 230, 0.22)',
  background: 'linear-gradient(180deg, rgba(14, 23, 39, 0.95), rgba(8, 13, 23, 0.96))',
};

export const ruleEditorCardStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 16,
  padding: 22,
  borderRadius: 28,
  border: '1px solid rgba(120, 158, 230, 0.22)',
  background: 'linear-gradient(180deg, rgba(14, 23, 39, 0.95), rgba(8, 13, 23, 0.96))',
};

export const rulePreviewCardStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 16,
  padding: 22,
  borderRadius: 28,
  border: '1px solid rgba(120, 158, 230, 0.22)',
  background: 'linear-gradient(180deg, rgba(14, 23, 39, 0.95), rgba(8, 13, 23, 0.96))',
};

export const ruleCardHeaderStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 12,
  flexWrap: 'wrap',
};

export const ruleActionRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: 10,
  flexWrap: 'wrap',
};

export const ruleListStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 10,
};

export const ruleListItemStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'flex-start',
  gap: 6,
  padding: '14px 15px',
  borderRadius: 18,
  border: '1px solid rgba(121, 156, 225, 0.18)',
  background: 'rgba(16, 27, 44, 0.82)',
  color: '#eff6ff',
  cursor: 'pointer',
  textAlign: 'left',
};

export const ruleListItemActiveStyle: CSSProperties = {
  borderColor: 'rgba(159, 198, 255, 0.42)',
  background: 'rgba(32, 52, 83, 0.92)',
  boxShadow: '0 16px 28px rgba(0, 0, 0, 0.18)',
};

export const ruleListTitleStyle: CSSProperties = {
  fontSize: 14,
  fontWeight: 700,
};

export const ruleListMetaStyle: CSSProperties = {
  color: 'rgba(198, 214, 238, 0.74)',
  fontSize: 12,
  lineHeight: 1.45,
};

export const ruleEditorGridStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(2, minmax(0, 1fr))',
  gap: 14,
};

export const rulePreviewHeroStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '68px minmax(0, 1fr)',
  gap: 14,
  alignItems: 'start',
};

export const rulePreviewAvatarStyle: CSSProperties = {
  width: 68,
  height: 68,
  borderRadius: '50%',
  overflow: 'hidden',
  boxShadow: '0 16px 32px rgba(0, 0, 0, 0.25)',
  background: 'rgba(5, 12, 22, 0.65)',
};

export const rulePreviewBubbleStyle: CSSProperties = {
  display: 'grid',
  gap: 10,
};

export const rulePreviewTypeStyle: CSSProperties = {
  color: 'rgba(159, 198, 255, 0.82)',
  fontSize: 12,
  fontWeight: 700,
  textTransform: 'uppercase',
  letterSpacing: '0.08em',
};

export const rulePreviewMetaStyle: CSSProperties = {
  display: 'grid',
  gap: 10,
};

export const ruleNoteCardStyle: CSSProperties = {
  display: 'grid',
  gap: 8,
  padding: 14,
  borderRadius: 18,
  border: '1px solid rgba(120, 158, 230, 0.18)',
  background: 'rgba(13, 21, 35, 0.72)',
};

export const ruleNoteCopyStyle: CSSProperties = {
  margin: 0,
  color: 'rgba(205, 220, 242, 0.84)',
  fontSize: 13,
  lineHeight: 1.55,
  whiteSpace: 'pre-wrap',
};

export const infoGridStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(3, minmax(0, 1fr))',
  gap: 18,
  marginTop: 22,
};

export const infoCardStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 16,
  minHeight: 240,
  padding: 22,
  borderRadius: 28,
  border: '1px solid rgba(120, 158, 230, 0.22)',
  background: 'linear-gradient(180deg, rgba(14, 23, 39, 0.95), rgba(8, 13, 23, 0.96))',
};

export const infoCardHeaderStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 6,
};

export const infoCardSubtitleStyle: CSSProperties = {
  margin: 0,
  color: 'rgba(197, 213, 238, 0.78)',
  fontSize: 13,
  lineHeight: 1.5,
};

export const detailRowStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '96px minmax(0, 1fr)',
  gap: 12,
  alignItems: 'start',
};

export const detailLabelStyle: CSSProperties = {
  color: 'rgba(197, 213, 238, 0.72)',
  fontSize: 12,
  textTransform: 'uppercase',
  letterSpacing: '0.04em',
};

export const detailValueStyle: CSSProperties = {
  display: 'block',
  padding: '8px 10px',
  borderRadius: 12,
  background: 'rgba(95, 131, 255, 0.1)',
  border: '1px solid rgba(95, 131, 255, 0.18)',
  color: 'rgba(189, 214, 255, 0.94)',
  fontSize: 12,
  whiteSpace: 'pre-wrap',
};

export const tagGroupStyle: CSSProperties = {
  display: 'flex',
  flexWrap: 'wrap',
  gap: 8,
};

export const tagStyle: CSSProperties = {
  borderRadius: 999,
  padding: '8px 10px',
  background: 'rgba(28, 43, 68, 0.9)',
  border: '1px solid rgba(123, 158, 230, 0.2)',
  color: '#eef6ff',
  fontSize: 12,
};

export const bulletListStyle: CSSProperties = {
  margin: 0,
  paddingLeft: 18,
  color: 'rgba(216, 226, 242, 0.86)',
  fontSize: 14,
  lineHeight: 1.6,
};

export const messagePreviewStyle: CSSProperties = {
  padding: 16,
  borderRadius: 18,
  border: '1px solid rgba(122, 162, 235, 0.2)',
  background: 'rgba(11, 17, 30, 0.92)',
  color: '#eef6ff',
  fontSize: 14,
  lineHeight: 1.6,
  whiteSpace: 'pre-wrap',
};

export const layoutStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'minmax(320px, 420px) minmax(0, 1fr)',
  gap: 22,
  marginTop: 22,
};

export const controlsPanelStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 16,
  padding: 22,
  borderRadius: 28,
  border: '1px solid rgba(120, 158, 230, 0.22)',
  background: 'linear-gradient(180deg, rgba(14, 23, 39, 0.95), rgba(8, 13, 23, 0.96))',
};

export const previewPanelStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 18,
  padding: 22,
  borderRadius: 28,
  border: '1px solid rgba(120, 158, 230, 0.22)',
  background: 'linear-gradient(180deg, rgba(14, 23, 39, 0.95), rgba(8, 13, 23, 0.96))',
};

export const sectionTitleStyle: CSSProperties = {
  margin: 0,
  color: '#f5faff',
  fontSize: 18,
};

export const fieldStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 8,
};

export const labelStyle: CSSProperties = {
  color: 'rgba(197, 213, 238, 0.82)',
  fontSize: 13,
  fontWeight: 600,
};

export const selectStyle: CSSProperties = {
  borderRadius: 14,
  border: '1px solid rgba(124, 156, 214, 0.24)',
  background: 'rgba(13, 21, 35, 0.96)',
  color: '#eff6ff',
  padding: '12px 14px',
  fontSize: 14,
};

export const textareaStyle: CSSProperties = {
  minHeight: 96,
  resize: 'vertical',
  borderRadius: 16,
  border: '1px solid rgba(124, 156, 214, 0.24)',
  background: 'rgba(13, 21, 35, 0.96)',
  color: '#eff6ff',
  padding: '12px 14px',
  fontSize: 14,
  lineHeight: 1.5,
};

export const rangeGroupStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 12,
};

export const rangeStyle: CSSProperties = {
  width: '100%',
};

export const rangeFieldRowStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1fr) 88px',
  gap: 10,
  alignItems: 'center',
};

export const rangeInputStyle: CSSProperties = {
  width: '100%',
};

export const numberInputStyle: CSSProperties = {
  width: '100%',
  minHeight: 40,
  padding: '0 10px',
  borderRadius: 10,
  border: '1px solid rgba(124, 156, 214, 0.24)',
  background: 'rgba(13, 21, 35, 0.96)',
  color: '#eff6ff',
  fontSize: 14,
};

export const toggleGridStyle: CSSProperties = {
  display: 'flex',
  flexWrap: 'wrap',
  gap: 10,
};

export const toggleChipStyle: CSSProperties = {
  border: '1px solid rgba(124, 156, 214, 0.24)',
  borderRadius: 999,
  background: 'rgba(18, 29, 47, 0.86)',
  color: 'rgba(222, 234, 252, 0.88)',
  padding: '10px 14px',
  cursor: 'pointer',
};

export const toggleChipActiveStyle: CSSProperties = {
  background: 'rgba(115, 172, 255, 0.18)',
  borderColor: 'rgba(148, 194, 255, 0.44)',
  color: '#f4f9ff',
};

export const triggerGroupStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 10,
};

export const triggerButtonsStyle: CSSProperties = {
  display: 'flex',
  flexWrap: 'wrap',
  gap: 8,
};

export const triggerButtonStyle: CSSProperties = {
  border: '1px solid rgba(121, 156, 225, 0.2)',
  borderRadius: 999,
  background: 'rgba(24, 37, 60, 0.84)',
  color: '#ebf5ff',
  padding: '8px 12px',
  fontSize: 12,
  cursor: 'pointer',
};

export const previewGridStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(2, minmax(0, 1fr))',
  gap: 16,
};

export const previewCardStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  gap: 12,
  minHeight: 220,
  padding: 18,
  borderRadius: 24,
  border: '1px solid rgba(122, 162, 235, 0.2)',
  background: 'radial-gradient(circle at top, rgba(39, 64, 101, 0.3), rgba(11, 17, 30, 0.92) 70%)',
};

export const miniLabelStyle: CSSProperties = {
  color: 'rgba(196, 212, 238, 0.84)',
  fontSize: 12,
  fontWeight: 600,
};

export const codeStyle: CSSProperties = {
  whiteSpace: 'pre-wrap',
  color: 'rgba(172, 206, 255, 0.86)',
  fontSize: 12,
  textAlign: 'center',
};

export const dockStageStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 12,
  padding: 18,
  borderRadius: 24,
  border: '1px solid rgba(122, 162, 235, 0.2)',
  background: 'linear-gradient(180deg, rgba(18, 29, 48, 0.94), rgba(9, 14, 24, 0.96))',
};

export const dockHeaderStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 10,
};

export const codePillStyle: CSSProperties = {
  borderRadius: 999,
  background: 'rgba(255, 255, 255, 0.06)',
  color: 'rgba(206, 221, 242, 0.9)',
  padding: '6px 10px',
  fontSize: 12,
};

export const panelPreviewRowStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1fr) 320px',
  gap: 16,
};

export const panelColumnStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 12,
};

export const sheetColumnStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 12,
};

export const spriteSheetFrameStyle: CSSProperties = {
  borderRadius: 22,
  border: '1px solid rgba(120, 158, 226, 0.18)',
  background: 'rgba(5, 10, 18, 0.9)',
  padding: 14,
};

export const spriteSheetStyle: CSSProperties = {
  width: '100%',
  aspectRatio: '5 / 3',
  borderRadius: 14,
  backgroundImage: "url('/images/lumi.webp')",
  backgroundSize: 'cover',
  backgroundPosition: 'center',
  backgroundRepeat: 'no-repeat',
};

export const sheetHelpStyle: CSSProperties = {
  margin: 0,
  color: 'rgba(193, 212, 238, 0.78)',
  fontSize: 13,
  lineHeight: 1.55,
};
