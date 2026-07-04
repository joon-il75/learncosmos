// resourceNotebook + contentSaveContract styles

export const resourceNotebookStyle = {
  display: 'grid',
  gap: '12px',
  marginTop: '6px',
  paddingTop: '14px',
  borderTop: '1px solid rgba(110, 88, 58, 0.12)',
} as const;

export const resourceNotebookHeaderStyle = {
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: '10px',
  flexWrap: 'wrap',
} as const;

export const resourceNotebookEyebrowStyle = {
  color: '#8B6840',
  fontSize: '11px',
  fontWeight: 900,
  letterSpacing: '0.08em',
  textTransform: 'uppercase',
} as const;

export const resourceNotebookTitleStyle = {
  display: 'block',
  marginTop: '3px',
  color: '#2E2418',
  fontSize: '16px',
} as const;

export const resourceNotebookCountStyle = {
  display: 'inline-flex',
  minHeight: '28px',
  alignItems: 'center',
  padding: '0 10px',
  borderRadius: '999px',
  background: 'rgba(86, 114, 174, 0.12)',
  color: '#4A5F92',
  fontSize: '12px',
  fontWeight: 900,
} as const;

export const resourceNotebookListStyle = {
  display: 'grid',
  gap: '9px',
} as const;

export const resourceNotebookItemStyle = {
  display: 'grid',
  gap: '10px',
  padding: '12px',
  borderRadius: '15px',
  border: '1px solid rgba(110, 88, 58, 0.12)',
  background: 'rgba(255, 252, 244, 0.58)',
} as const;

export const resourceNotebookItemSelectedStyle = {
  border: '1px solid rgba(34, 197, 94, 0.44)',
  background: 'linear-gradient(135deg, rgba(236, 253, 245, 0.82), rgba(255, 252, 244, 0.62))',
  boxShadow: '0 12px 30px rgba(34, 197, 94, 0.12)',
} as const;

export const resourceNotebookItemCandidateStyle = {
  border: '1px solid rgba(59, 130, 246, 0.38)',
  background: 'linear-gradient(135deg, rgba(239, 246, 255, 0.84), rgba(255, 252, 244, 0.62))',
  boxShadow: '0 12px 30px rgba(59, 130, 246, 0.1)',
} as const;

export const resourceNotebookItemRejectedStyle = {
  border: '1px solid rgba(148, 163, 184, 0.34)',
  background: 'rgba(248, 250, 252, 0.52)',
  opacity: 0.72,
} as const;

export const resourceNotebookItemUnknownStyle = {
  border: '1px solid rgba(180, 83, 9, 0.28)',
  background: 'linear-gradient(135deg, rgba(254, 243, 199, 0.76), rgba(255, 252, 244, 0.62))',
} as const;

export const resourceNotebookItemHeaderStyle = {
  display: 'flex',
  gap: '10px',
  alignItems: 'flex-start',
} as const;

export const resourceNotebookNumberStyle = {
  display: 'inline-grid',
  width: '26px',
  height: '26px',
  flex: '0 0 auto',
  placeItems: 'center',
  borderRadius: '999px',
  background: 'rgba(112, 87, 50, 0.12)',
  color: '#6A5237',
  fontSize: '12px',
  fontWeight: 900,
} as const;

export const resourceNotebookItemContentStyle = {
  display: 'grid',
  gap: '5px',
  minWidth: 0,
  flex: '1 1 auto',
} as const;

export const resourceNotebookItemTitleRowStyle = {
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: '8px',
  flexWrap: 'wrap',
} as const;

export const resourceNotebookItemTitleStyle = {
  color: '#342819',
  fontSize: '13px',
  lineHeight: 1.35,
  overflowWrap: 'anywhere',
} as const;

export const resourceNotebookStatusBadgeStyle = {
  display: 'inline-flex',
  minHeight: '24px',
  alignItems: 'center',
  padding: '0 9px',
  borderRadius: '999px',
  fontSize: '11px',
  fontWeight: 900,
  lineHeight: 1,
  whiteSpace: 'nowrap',
} as const;

export const resourceNotebookStatusSelectedStyle = {
  border: '1px solid rgba(22, 163, 74, 0.22)',
  background: 'rgba(220, 252, 231, 0.92)',
  color: '#166534',
} as const;

export const resourceNotebookStatusCandidateStyle = {
  border: '1px solid rgba(37, 99, 235, 0.22)',
  background: 'rgba(219, 234, 254, 0.94)',
  color: '#1D4ED8',
} as const;

export const resourceNotebookStatusRejectedStyle = {
  border: '1px solid rgba(100, 116, 139, 0.2)',
  background: 'rgba(241, 245, 249, 0.9)',
  color: '#64748B',
} as const;

export const resourceNotebookStatusUnknownStyle = {
  border: '1px solid rgba(180, 83, 9, 0.2)',
  background: 'rgba(254, 243, 199, 0.9)',
  color: '#92400E',
} as const;

export const resourceNotebookItemMetaStyle = {
  color: '#7B664C',
  fontSize: '12px',
  lineHeight: 1.4,
} as const;

export const resourceNotebookLinkStyle = {
  justifySelf: 'start',
  color: '#355FA7',
  fontSize: '12px',
  fontWeight: 900,
  textDecoration: 'none',
} as const;

export const resourceNotebookSelectButtonStyle = {
  justifySelf: 'start',
  minHeight: '32px',
  padding: '0 12px',
  borderRadius: '999px',
  border: '1px solid rgba(37, 99, 235, 0.22)',
  background: 'rgba(219, 234, 254, 0.94)',
  color: '#1D4ED8',
  fontSize: '12px',
  fontWeight: 900,
  cursor: 'pointer',
} as const;

export const resourceNotebookMutedStyle = {
  color: '#8A765D',
  fontSize: '12px',
  lineHeight: 1.4,
} as const;

export const resourceNotebookEmptyStyle = {
  padding: '14px',
  borderRadius: '15px',
  border: '1px dashed rgba(110, 88, 58, 0.18)',
  background: 'rgba(255, 252, 244, 0.42)',
  color: '#6D5A41',
  fontSize: '13px',
  lineHeight: 1.6,
} as const;

export const resourceNotebookFooterStyle = {
  display: 'flex',
  alignItems: 'center',
  gap: '10px',
  flexWrap: 'wrap',
} as const;

export const resourceNotebookMessageStyle = {
  padding: '10px 12px',
  borderRadius: '14px',
  border: '1px solid rgba(56, 139, 253, 0.16)',
  background: 'rgba(219, 234, 254, 0.34)',
  color: '#4C5F7A',
  fontSize: '12px',
  lineHeight: 1.5,
} as const;

export const contentSaveContractStyle = {
  display: 'grid',
  gap: '12px',
  padding: '14px',
  borderRadius: '18px',
  border: '1px solid rgba(125, 211, 252, 0.18)',
  background: 'rgba(14, 165, 233, 0.08)',
} as const;

export const contentSaveContractGridStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(190px, 1fr))',
  gap: '10px',
} as const;

export const contentSaveContractItemStyle = {
  display: 'grid',
  gap: '5px',
  padding: '11px',
  borderRadius: '14px',
  border: '1px solid rgba(125, 211, 252, 0.14)',
  background: 'rgba(255,255,255,0.05)',
  color: '#D8F3FF',
  fontSize: '12px',
  lineHeight: 1.55,
} as const;
