// imageSlot, prompt, detail card / edit field styles

export const imageSlotCardStyle = {
  display: 'grid',
  gap: '14px',
} as const;

export const imageSlotPreviewStyle = {
  minHeight: '180px',
  borderRadius: '18px',
  border: '1px dashed rgba(118, 92, 55, 0.28)',
  background: 'linear-gradient(135deg, rgba(255,255,255,0.46), rgba(224, 208, 180, 0.64))',
  display: 'grid',
  placeItems: 'center',
  color: '#7F6240',
  fontSize: '15px',
  fontWeight: 700,
} as const;

export const imageSlotTextBlockStyle = {
  display: 'grid',
  gap: '8px',
} as const;

export const imageSlotTitleStyle = {
  fontSize: '16px',
  color: '#2D2419',
} as const;

export const imageSlotDescriptionStyle = {
  margin: 0,
  fontSize: '14px',
  lineHeight: 1.65,
  color: '#5A4937',
} as const;

export const imageSlotListStyle = {
  margin: 0,
  paddingLeft: '18px',
  display: 'grid',
  gap: '6px',
  color: '#5D4A38',
  fontSize: '13px',
  lineHeight: 1.6,
} as const;

export const promptBoxStyle = {
  display: 'grid',
  gap: '6px',
  padding: '12px',
  borderRadius: '14px',
  border: '1px solid rgba(104, 82, 54, 0.12)',
  background: 'rgba(255, 248, 236, 0.58)',
} as const;

export const promptLabelStyle = {
  fontSize: '11px',
  fontWeight: 800,
  letterSpacing: '0.08em',
  textTransform: 'uppercase',
  color: '#8B6840',
} as const;

export const promptTextStyle = {
  margin: 0,
  fontSize: '12px',
  lineHeight: 1.6,
  color: '#5D4A38',
} as const;

export const detailCardStyle = {
  display: 'grid',
  gap: '10px',
  minHeight: '180px',
  padding: '16px',
  borderRadius: '18px',
  border: '1px solid rgba(110, 88, 58, 0.14)',
  background: 'rgba(255,255,255,0.46)',
  alignContent: 'start',
} as const;

export const detailMetaBadgeStyle = {
  display: 'inline-flex',
  width: 'fit-content',
  minHeight: '28px',
  alignItems: 'center',
  padding: '0 10px',
  borderRadius: '999px',
  background: 'rgba(86, 114, 174, 0.12)',
  color: '#4A5F92',
  fontSize: '12px',
  fontWeight: 700,
} as const;

export const detailTitleStyle = {
  margin: 0,
  fontSize: '19px',
  color: '#2E2418',
} as const;

export const detailDescriptionStyle = {
  margin: 0,
  fontSize: '14px',
  lineHeight: 1.65,
  color: '#5B4A38',
} as const;

export const detailPathStyle = {
  display: 'flex',
  gap: '8px',
  flexWrap: 'wrap',
  alignItems: 'center',
  marginTop: '4px',
  paddingTop: '12px',
  borderTop: '1px solid rgba(110, 88, 58, 0.12)',
  color: '#6A5237',
  fontSize: '12px',
  fontWeight: 700,
} as const;

export const detailEditScaffoldStyle = {
  display: 'grid',
  gap: '12px',
  marginTop: '6px',
  paddingTop: '14px',
  borderTop: '1px solid rgba(110, 88, 58, 0.12)',
} as const;

export const detailEditNoticeStyle = {
  padding: '11px 12px',
  borderRadius: '14px',
  border: '1px solid rgba(94, 120, 166, 0.14)',
  background: 'rgba(91, 117, 168, 0.1)',
  color: '#536381',
  fontSize: '12px',
  lineHeight: 1.55,
} as const;

export const detailDraftStatusRowStyle = {
  display: 'flex',
  alignItems: 'center',
  gap: '8px',
  flexWrap: 'wrap',
} as const;

export const detailDraftStatusPillStyle = {
  display: 'inline-flex',
  minHeight: '28px',
  alignItems: 'center',
  padding: '0 10px',
  borderRadius: '999px',
  border: '1px solid rgba(110, 88, 58, 0.14)',
  background: 'rgba(110, 88, 58, 0.08)',
  color: '#6B5843',
  fontSize: '12px',
  fontWeight: 900,
} as const;

export const detailDraftStatusChangedStyle = {
  borderColor: 'rgba(177, 110, 33, 0.28)',
  background: 'rgba(177, 110, 33, 0.14)',
  color: '#8B541F',
} as const;

export const detailDraftStatusTextStyle = {
  color: '#735E43',
  fontSize: '12px',
  lineHeight: 1.5,
} as const;

export const detailEditFieldStyle = {
  display: 'grid',
  gap: '7px',
} as const;

export const detailEditLabelStyle = {
  color: '#5D4A38',
  fontSize: '12px',
  fontWeight: 900,
} as const;

export const detailEditInputStyle = {
  minHeight: '42px',
  borderRadius: '13px',
  border: '1px solid rgba(110, 88, 58, 0.16)',
  background: 'rgba(255, 252, 244, 0.66)',
  color: '#2E2418',
  padding: '0 12px',
  fontSize: '14px',
  fontFamily: 'inherit',
  outline: 'none',
} as const;

export const detailEditSelectStyle = {
  ...detailEditInputStyle,
  appearance: 'none',
} as const;

export const detailEditTextareaStyle = {
  minHeight: '90px',
  borderRadius: '13px',
  border: '1px solid rgba(110, 88, 58, 0.16)',
  background: 'rgba(255, 252, 244, 0.66)',
  color: '#2E2418',
  padding: '11px 12px',
  fontSize: '13px',
  lineHeight: 1.55,
  fontFamily: 'inherit',
  outline: 'none',
  resize: 'vertical',
} as const;

export const detailEditFooterStyle = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'flex-start',
  gap: '10px',
  flexWrap: 'wrap',
} as const;

export const detailEditResetButtonStyle = {
  minHeight: '38px',
  padding: '0 14px',
  borderRadius: '999px',
  border: '1px solid rgba(112, 87, 50, 0.2)',
  background: 'rgba(255, 252, 244, 0.72)',
  color: '#5E452A',
  fontSize: '13px',
  fontWeight: 900,
  cursor: 'pointer',
} as const;

export const detailEditResetButtonDisabledStyle = {
  opacity: 0.48,
  cursor: 'not-allowed',
} as const;

export const detailEditDisabledButtonStyle = {
  minHeight: '38px',
  padding: '0 14px',
  borderRadius: '999px',
  border: '1px solid rgba(110, 88, 58, 0.14)',
  background: 'rgba(110, 88, 58, 0.08)',
  color: 'rgba(80, 63, 43, 0.58)',
  fontSize: '13px',
  fontWeight: 900,
  cursor: 'not-allowed',
} as const;

export const detailEditHintStyle = {
  color: '#735E43',
  fontSize: '12px',
  lineHeight: 1.5,
} as const;
