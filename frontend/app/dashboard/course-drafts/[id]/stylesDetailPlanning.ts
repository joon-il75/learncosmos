// planning action, selectionReadiness, transitionContract styles

export const planningActionCopyStyle = {
  display: 'grid',
  gap: '8px',
} as const;

export const planningActionTitleStyle = {
  margin: 0,
  fontSize: '22px',
  color: '#F7FAFF',
} as const;

export const planningActionDescriptionStyle = {
  margin: 0,
  maxWidth: '720px',
  color: '#CAD7EF',
  fontSize: '14px',
  lineHeight: 1.65,
} as const;

export const planningReadinessPanelStyle = {
  display: 'grid',
  gap: '12px',
  padding: '16px',
  borderRadius: '20px',
  border: '1px solid rgba(194, 210, 245, 0.12)',
  background: 'rgba(255,255,255,0.06)',
} as const;

export const planningReadinessHeaderStyle = {
  display: 'grid',
  gap: '4px',
} as const;

export const planningReadinessTitleStyle = {
  color: '#F7FAFF',
  fontSize: '16px',
} as const;

export const planningReadinessListStyle = {
  display: 'grid',
  gap: '10px',
} as const;

export const planningReadinessItemStyle = {
  display: 'flex',
  gap: '10px',
  alignItems: 'flex-start',
} as const;

export const planningReadinessDotStyle = {
  width: '10px',
  height: '10px',
  flex: '0 0 auto',
  marginTop: '5px',
  borderRadius: '999px',
} as const;

export const planningReadinessDotReadyStyle = {
  background: '#5EEAD4',
  boxShadow: '0 0 0 5px rgba(94, 234, 212, 0.12)',
} as const;

export const planningReadinessDotPendingStyle = {
  background: '#F6C76A',
  boxShadow: '0 0 0 5px rgba(246, 199, 106, 0.12)',
} as const;

export const planningReadinessItemTitleStyle = {
  color: '#F7FAFF',
  fontSize: '13px',
  lineHeight: 1.35,
} as const;

export const planningReadinessItemTextStyle = {
  color: '#C8D5EE',
  fontSize: '12px',
  lineHeight: 1.45,
} as const;

export const selectionReadinessReasonGridStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(210px, 1fr))',
  gap: '10px',
  paddingTop: '2px',
} as const;

export const selectionReadinessReasonItemStyle = {
  display: 'grid',
  gap: '6px',
  padding: '12px',
  borderRadius: '16px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(255,255,255,0.05)',
  color: '#D8E5FF',
  fontSize: '12px',
  lineHeight: 1.55,
} as const;

export const selectionReadinessReasonPillStyle = {
  justifySelf: 'start',
  minHeight: '24px',
  display: 'inline-flex',
  alignItems: 'center',
  padding: '0 9px',
  borderRadius: '999px',
  fontSize: '11px',
  fontWeight: 900,
  lineHeight: 1,
} as const;

export const selectionReadinessReasonPillReadyStyle = {
  background: 'rgba(94, 234, 212, 0.16)',
  color: '#9AF7E9',
  border: '1px solid rgba(94, 234, 212, 0.22)',
} as const;

export const selectionReadinessReasonPillPendingStyle = {
  background: 'rgba(246, 199, 106, 0.14)',
  color: '#FFE1A1',
  border: '1px solid rgba(246, 199, 106, 0.22)',
} as const;

export const planningNextActionGuideStyle = {
  display: 'grid',
  gap: '6px',
  padding: '13px 14px',
  borderRadius: '17px',
  border: '1px solid rgba(125, 211, 252, 0.18)',
  background: 'linear-gradient(135deg, rgba(14, 165, 233, 0.12), rgba(255,255,255,0.05))',
  color: '#DDF7FF',
  fontSize: '12px',
  lineHeight: 1.6,
} as const;

export const planningNextActionGuideEyebrowStyle = {
  color: '#9FE7FF',
  fontSize: '11px',
  fontWeight: 900,
  letterSpacing: '0.08em',
  textTransform: 'uppercase',
} as const;

export const transitionContractPanelStyle = {
  display: 'grid',
  gap: '12px',
  padding: '16px',
  borderRadius: '20px',
  border: '1px solid rgba(255, 218, 146, 0.16)',
  background: 'rgba(255, 218, 146, 0.07)',
} as const;

export const transitionContractGridStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(210px, 1fr))',
  gap: '10px',
} as const;

export const transitionContractItemStyle = {
  display: 'grid',
  gap: '5px',
  padding: '12px',
  borderRadius: '16px',
  border: '1px solid rgba(255, 218, 146, 0.12)',
  background: 'rgba(255,255,255,0.05)',
  color: '#F4E8C9',
  fontSize: '12px',
  lineHeight: 1.55,
} as const;

export const planningActionButtonGroupStyle = {
  display: 'flex',
  gap: '10px',
  flexWrap: 'wrap',
  justifyContent: 'flex-end',
} as const;

export const planningSecondaryButtonStyle = {
  minHeight: '44px',
  padding: '0 16px',
  borderRadius: '999px',
  border: '1px solid rgba(194, 210, 245, 0.16)',
  background: 'rgba(255,255,255,0.08)',
  color: 'rgba(244,247,255,0.62)',
  fontSize: '14px',
  fontWeight: 800,
  cursor: 'not-allowed',
} as const;

export const planningSecondaryButtonActiveStyle = {
  border: '1px solid rgba(112, 166, 255, 0.34)',
  background: 'linear-gradient(135deg, rgba(59, 130, 246, 0.36), rgba(94, 234, 212, 0.18))',
  color: '#F7FAFF',
  cursor: 'pointer',
} as const;

export const planningPrimaryButtonStyle = {
  minHeight: '44px',
  padding: '0 18px',
  borderRadius: '999px',
  border: '1px solid rgba(255, 218, 146, 0.32)',
  background: 'linear-gradient(135deg, rgba(245, 184, 83, 0.36), rgba(255, 228, 160, 0.18))',
  color: 'rgba(255,246,224,0.66)',
  fontSize: '14px',
  fontWeight: 900,
  cursor: 'not-allowed',
} as const;

export const planningPrimaryButtonReviewStyle = {
  color: '#FFF0C2',
  cursor: 'pointer',
} as const;

export const planningPrimaryButtonActiveStyle = {
  border: '1px solid rgba(255, 218, 146, 0.56)',
  background: 'linear-gradient(135deg, rgba(245, 184, 83, 0.72), rgba(255, 228, 160, 0.32))',
  color: '#FFF8DF',
  boxShadow: '0 16px 34px rgba(245, 184, 83, 0.18)',
  cursor: 'pointer',
} as const;

export const planningActionMessageStyle = {
  justifySelf: 'end',
  maxWidth: '560px',
  padding: '10px 12px',
  borderRadius: '14px',
  border: '1px solid rgba(94, 234, 212, 0.22)',
  background: 'rgba(20, 184, 166, 0.11)',
  color: '#CFFCF4',
  fontSize: '12px',
  lineHeight: 1.55,
} as const;

export const planningActionChecklistStyle = {
  gridColumn: '1 / -1',
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(210px, 1fr))',
  gap: '10px',
} as const;

export const planningActionChecklistItemStyle = {
  display: 'grid',
  gap: '5px',
  padding: '12px',
  borderRadius: '16px',
  border: '1px solid rgba(194, 210, 245, 0.1)',
  background: 'rgba(255,255,255,0.06)',
  color: '#E8EEFF',
  fontSize: '12px',
  lineHeight: 1.55,
} as const;
