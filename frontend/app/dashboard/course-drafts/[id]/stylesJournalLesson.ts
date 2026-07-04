import type { CSSProperties } from 'react';

export const journalScaffoldStyle = {
  display: 'grid',
  gap: '18px',
  padding: '22px',
  borderRadius: '26px',
  border: '1px solid rgba(201, 179, 143, 0.18)',
  background: 'linear-gradient(180deg, rgba(248, 238, 217, 0.95), rgba(229, 212, 184, 0.94))',
  boxShadow: '0 24px 48px rgba(0, 0, 0, 0.2)',
  color: '#2E2418',
} as const;

export const journalContextGridStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
  gap: '14px',
} as const;

export const journalContextCardStyle = {
  display: 'grid',
  gap: '8px',
  padding: '16px',
  borderRadius: '18px',
  border: '1px solid rgba(107, 86, 60, 0.12)',
  background: 'rgba(255,255,255,0.46)',
} as const;

export const journalContextTitleStyle = {
  fontSize: '18px',
  color: '#2E2418',
  lineHeight: 1.4,
} as const;

export const journalContextDescriptionStyle = {
  margin: 0,
  color: '#5B4A38',
  fontSize: '14px',
  lineHeight: 1.6,
} as const;

export const journalPromptGridStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))',
  gap: '14px',
} as const;

export const journalPromptCardStyle = {
  display: 'grid',
  gap: '10px',
  padding: '18px',
  borderRadius: '20px',
  border: '1px solid rgba(107, 86, 60, 0.12)',
  background: 'rgba(255, 252, 244, 0.64)',
} as const;

export const journalPromptKickerStyle = {
  fontSize: '11px',
  fontWeight: 800,
  letterSpacing: '0.1em',
  textTransform: 'uppercase',
  color: '#8B6840',
} as const;

export const journalPromptTitleStyle = {
  margin: 0,
  fontSize: '18px',
  color: '#2E2418',
} as const;

export const journalPromptDescriptionStyle = {
  margin: 0,
  color: '#5B4A38',
  fontSize: '14px',
  lineHeight: 1.6,
} as const;

export const journalAnswerPlaceholderStyle = {
  minHeight: '96px',
  display: 'grid',
  placeItems: 'center',
  borderRadius: '16px',
  border: '1px dashed rgba(107, 86, 60, 0.2)',
  background: 'rgba(255,255,255,0.42)',
  color: '#8B6840',
  fontSize: '13px',
  fontWeight: 800,
} as const;

export const journalAnswerTextareaStyle = {
  width: '100%',
  minHeight: '132px',
  borderRadius: '16px',
  border: '1px solid rgba(107, 86, 60, 0.22)',
  background: 'rgba(255,255,255,0.72)',
  color: '#2E2418',
  padding: '14px 16px',
  fontSize: '14px',
  lineHeight: 1.7,
  resize: 'vertical',
  outline: 'none',
} as const;

export const journalAnswerTextareaDisabledStyle = {
  opacity: 0.58,
  cursor: 'not-allowed',
} as const;

export const journalActionRowStyle = {
  display: 'flex',
  justifyContent: 'flex-end',
  gap: '12px',
  flexWrap: 'wrap',
} as const;

export const journalSecondaryButtonStyle = {
  borderRadius: '999px',
  border: '1px solid rgba(107, 86, 60, 0.18)',
  background: 'rgba(255,255,255,0.62)',
  color: '#5B4A38',
  padding: '10px 16px',
  fontSize: '13px',
  fontWeight: 700,
} as const;

export const journalSecondaryButtonDisabledStyle = {
  opacity: 0.45,
  cursor: 'not-allowed',
} as const;

export const journalPrimaryButtonStyle = {
  borderRadius: '999px',
  border: '1px solid rgba(86, 114, 174, 0.18)',
  background: 'rgba(86, 114, 174, 0.14)',
  color: '#5A6FA2',
  padding: '10px 18px',
  fontSize: '13px',
  fontWeight: 800,
} as const;

export const journalPrimaryButtonActiveStyle = {
  background: 'linear-gradient(135deg, rgba(110, 136, 205, 0.92), rgba(83, 109, 178, 0.86))',
  color: '#FFFFFF',
  boxShadow: '0 14px 24px rgba(73, 95, 160, 0.2)',
} as const;

export const recordMetricGridStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(160px, 1fr))',
  gap: '12px',
} as const;

export const recordMetricCardStyle = {
  display: 'grid',
  gap: '8px',
  padding: '14px',
  borderRadius: '16px',
  border: '1px solid rgba(107, 86, 60, 0.14)',
  background: 'rgba(255,255,255,0.48)',
} as const;

export const recordMetricLabelStyle = {
  fontSize: '12px',
  fontWeight: 800,
  color: '#6D5337',
  letterSpacing: '0.02em',
} as const;

export const recordMetricInputStyle = {
  width: '100%',
  minHeight: '42px',
  borderRadius: '12px',
  border: '1px solid rgba(107, 86, 60, 0.18)',
  background: 'rgba(255,255,255,0.9)',
  color: '#2E2418',
  padding: '0 12px',
  fontSize: '14px',
} as const;

export const recordMetricSelectStyle = {
  ...recordMetricInputStyle,
} as const;

export const recordApplicationFieldStyle = {
  display: 'grid',
  gap: '8px',
} as const;

export const previewChecklistStyle = {
  display: 'grid',
  gap: '10px',
} as const;

export const previewChecklistItemStyle = {
  display: 'grid',
  gap: '6px',
  padding: '14px 16px',
  borderRadius: '16px',
  border: '1px solid rgba(86, 114, 174, 0.14)',
  background: 'rgba(255,255,255,0.44)',
  color: '#4A5F92',
  fontSize: '13px',
  lineHeight: 1.6,
} as const;

export const journalNoticeStyle = {
  padding: '14px 16px',
  borderRadius: '16px',
  border: '1px solid rgba(86, 114, 174, 0.16)',
  background: 'rgba(86, 114, 174, 0.1)',
  color: '#4A5F92',
  fontSize: '13px',
  lineHeight: 1.6,
  fontWeight: 700,
} as const;

export const sectionHeaderStyle = {
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: '12px',
  flexWrap: 'wrap',
} as const;

export const sectionTitleStyle = {
  margin: 0,
  fontSize: '28px',
  color: '#F7FAFF',
} as const;

export const sectionSubtitleStyle = {
  margin: '6px 0 0',
  fontSize: '14px',
  color: '#C8D1E8',
} as const;

export const secondaryButtonStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '42px',
  padding: '0 16px',
  borderRadius: '999px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(255,255,255,0.08)',
  color: '#F4F7FF',
  textDecoration: 'none',
  fontSize: '14px',
  fontWeight: 700,
} as const;

export const levelCardStyle = {
  display: 'grid',
  gap: '14px',
  padding: '18px',
  borderRadius: '20px',
  border: '1px solid rgba(194, 210, 245, 0.12)',
  background: 'rgba(255,255,255,0.05)',
} as const;

export const levelCardSelectedStyle = {
  border: '1px solid rgba(77, 104, 166, 0.36)',
  background: 'rgba(255,255,255,0.34)',
  boxShadow: '0 16px 28px rgba(68, 87, 133, 0.12)',
} as const;

export const levelBadgeStyle = {
  display: 'inline-flex',
  width: 'fit-content',
  alignItems: 'center',
  minHeight: '28px',
  padding: '0 10px',
  borderRadius: '999px',
  background: 'rgba(64, 160, 255, 0.14)',
  border: '1px solid rgba(64, 160, 255, 0.24)',
  color: '#D6E8FF',
  fontSize: '12px',
  fontWeight: 700,
} as const;

export const levelTitleStyle = {
  margin: 0,
  fontSize: '22px',
  color: '#F7FAFF',
} as const;

export const levelObjectiveStyle = {
  margin: 0,
  fontSize: '14px',
  color: '#D6E1F8',
  lineHeight: 1.6,
} as const;

export const lessonCardStyle = {
  display: 'grid',
  gap: '10px',
  padding: '16px',
  borderRadius: '16px',
  border: '1px solid rgba(194, 210, 245, 0.08)',
  background: 'rgba(5, 13, 25, 0.44)',
  width: '100%',
  textAlign: 'left',
  cursor: 'pointer',
  fontFamily: 'inherit',
} as const;

export const lessonCardSelectedStyle = {
  border: '1px solid rgba(77, 104, 166, 0.46)',
  background: 'rgba(72, 94, 145, 0.18)',
  boxShadow: 'inset 4px 0 0 rgba(77, 104, 166, 0.44)',
} as const;

export const lessonMetaStyle = {
  fontSize: '12px',
  color: '#9FB4DA',
} as const;

export const lessonTitleStyle = {
  fontSize: '17px',
  fontWeight: 700,
  color: '#F7FAFF',
} as const;

export const lessonSummaryStyle = {
  fontSize: '14px',
  lineHeight: 1.6,
  color: '#D0DCF6',
} as const;

export const resourceListStyle = {
  display: 'flex',
  gap: '8px',
  flexWrap: 'wrap',
} as const;

export const resourceChipStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: '28px',
  padding: '0 10px',
  borderRadius: '999px',
  background: 'rgba(243, 165, 51, 0.14)',
  border: '1px solid rgba(243, 165, 51, 0.28)',
  color: '#FFE5B4',
  fontSize: '12px',
} as const;
