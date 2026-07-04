// manualResource styles

export const manualResourceDraftStyle = {
  display: 'grid',
  gap: '12px',
  padding: '13px',
  borderRadius: '17px',
  border: '1px solid rgba(91, 117, 168, 0.14)',
  background: 'linear-gradient(180deg, rgba(255, 252, 244, 0.64), rgba(236, 229, 213, 0.54))',
} as const;

export const manualResourceHeaderStyle = {
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: '10px',
  flexWrap: 'wrap',
} as const;

export const manualResourceNoticeStyle = {
  padding: '11px 12px',
  borderRadius: '14px',
  border: '1px solid rgba(94, 120, 166, 0.12)',
  background: 'rgba(91, 117, 168, 0.08)',
  color: '#536381',
  fontSize: '12px',
  lineHeight: 1.55,
} as const;

export const manualResourceValidationStyle = {
  display: 'grid',
  gap: '9px',
  padding: '12px',
  borderRadius: '15px',
  border: '1px solid rgba(110, 88, 58, 0.12)',
  background: 'rgba(255, 252, 244, 0.46)',
} as const;

export const manualResourceValidationTitleStyle = {
  color: '#5D4A38',
  fontSize: '12px',
  fontWeight: 900,
} as const;

export const manualResourceValidationItemStyle = {
  display: 'flex',
  gap: '9px',
  alignItems: 'flex-start',
} as const;

export const manualResourceValidationDotStyle = {
  width: '9px',
  height: '9px',
  flex: '0 0 auto',
  marginTop: '5px',
  borderRadius: '999px',
} as const;

export const manualResourceValidationDotValidStyle = {
  background: '#4A7C59',
  boxShadow: '0 0 0 4px rgba(74, 124, 89, 0.12)',
} as const;

export const manualResourceValidationDotInvalidStyle = {
  background: '#B45332',
  boxShadow: '0 0 0 4px rgba(180, 83, 50, 0.12)',
} as const;

export const manualResourceValidationLabelStyle = {
  color: '#342819',
  fontSize: '12px',
  lineHeight: 1.3,
} as const;

export const manualResourceValidationMessageStyle = {
  color: '#735E43',
  fontSize: '12px',
  lineHeight: 1.45,
} as const;

export const manualResourceParseMessageStyle = {
  padding: '10px 12px',
  borderRadius: '14px',
  border: '1px solid rgba(56, 139, 253, 0.16)',
  background: 'rgba(219, 234, 254, 0.34)',
  color: '#4C5F7A',
  fontSize: '12px',
  lineHeight: 1.5,
} as const;

export const manualResourceCreateMessageStyle = {
  padding: '10px 12px',
  borderRadius: '14px',
  border: '1px solid rgba(74, 124, 89, 0.18)',
  background: 'rgba(236, 253, 245, 0.58)',
  color: '#34523F',
  fontSize: '12px',
  lineHeight: 1.5,
} as const;

export const manualResourceAttachMessageStyle = {
  padding: '10px 12px',
  borderRadius: '14px',
  border: '1px solid rgba(37, 99, 235, 0.18)',
  background: 'rgba(239, 246, 255, 0.62)',
  color: '#2D4368',
  fontSize: '12px',
  lineHeight: 1.5,
} as const;

export const manualResourcePreviewStyle = {
  display: 'grid',
  gap: '7px',
  padding: '12px',
  borderRadius: '15px',
  border: '1px solid rgba(74, 124, 89, 0.16)',
  background: 'rgba(236, 253, 245, 0.54)',
  color: '#2E4638',
  fontSize: '12px',
  lineHeight: 1.55,
} as const;

export const manualResourceCreatedContentStyle = {
  display: 'grid',
  gap: '7px',
  padding: '12px',
  borderRadius: '15px',
  border: '1px solid rgba(37, 99, 235, 0.18)',
  background: 'rgba(239, 246, 255, 0.64)',
  color: '#2D4368',
  fontSize: '12px',
  lineHeight: 1.55,
  overflowWrap: 'anywhere',
} as const;

export const manualResourceAttachActionStyle = {
  display: 'flex',
  alignItems: 'center',
  gap: '9px',
  flexWrap: 'wrap',
  color: '#3F5F8D',
} as const;

export const manualResourcePreviewMetaStyle = {
  display: 'flex',
  gap: '8px',
  flexWrap: 'wrap',
  color: '#4F6F5C',
  overflowWrap: 'anywhere',
} as const;

export const manualResourcePreviewWarningStyle = {
  color: '#9A5A2C',
  overflowWrap: 'anywhere',
} as const;
