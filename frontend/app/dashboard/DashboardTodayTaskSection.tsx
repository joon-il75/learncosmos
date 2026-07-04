'use client';

import type { CSSProperties } from 'react';
import type { DashboardMainCopy } from '@/lib/i18n/pages/dashboardMain';
import type { TodayTask, TodayTaskKind } from './_types';

interface Props {
  task: TodayTask | null;
  copy: DashboardMainCopy['todayTask'];
  onOpenTask: (href: string) => void;
}

const taskAccentByKind: Record<TodayTaskKind, string> = {
  create_course: '#F8B84E',
  start_planet: '#7DD3FC',
  continue_point: '#86EFAC',
  complete_planet: '#FDE68A',
  review_records: '#C4B5FD',
};

const sectionStyle: CSSProperties = {
  position: 'relative',
  overflow: 'hidden',
  borderRadius: '28px',
  padding: '18px 20px',
  background:
    'radial-gradient(circle at 16% 8%, rgba(134, 239, 172, 0.22), transparent 34%), linear-gradient(135deg, rgba(8, 21, 41, 0.9), rgba(5, 12, 26, 0.74))',
  boxShadow: '0 26px 80px rgba(0, 0, 0, 0.28), inset 0 1px 0 rgba(255,255,255,0.08)',
  backdropFilter: 'blur(18px)',
};

const innerStyle: CSSProperties = {
  position: 'relative',
  zIndex: 1,
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1fr) auto',
  gap: '18px',
  alignItems: 'center',
};

const eyebrowStyle: CSSProperties = {
  margin: 0,
  color: 'rgba(211, 226, 255, 0.78)',
  fontSize: '12px',
  fontWeight: 800,
  letterSpacing: '0.16em',
  textTransform: 'uppercase',
};

const titleStyle: CSSProperties = {
  margin: '6px 0 0',
  color: '#F8FAFF',
  fontSize: 'clamp(20px, 2.4vw, 30px)',
  lineHeight: 1.16,
  fontWeight: 850,
};

const descriptionStyle: CSSProperties = {
  margin: '8px 0 0',
  maxWidth: '720px',
  color: 'rgba(219, 229, 250, 0.82)',
  fontSize: '14px',
  lineHeight: 1.55,
};

const actionWrapStyle: CSSProperties = {
  display: 'grid',
  gap: '10px',
  justifyItems: 'end',
};

const buttonStyle: CSSProperties = {
  minHeight: '44px',
  padding: '0 18px',
  border: 0,
  borderRadius: '999px',
  background: 'linear-gradient(135deg, #F8B84E 0%, #E77922 100%)',
  color: '#201103',
  fontSize: '13px',
  fontWeight: 850,
  cursor: 'pointer',
  fontFamily: 'inherit',
  boxShadow: '0 14px 30px rgba(231, 121, 34, 0.28)',
};

const metaStyle: CSSProperties = {
  color: 'rgba(226, 236, 255, 0.7)',
  fontSize: '12px',
  fontWeight: 700,
};

export default function DashboardTodayTaskSection({ task, copy, onOpenTask }: Props) {
  if (!task) return null;

  const accent = taskAccentByKind[task.kind] ?? '#86EFAC';

  return (
    <section aria-label={copy.ariaLabel} style={sectionStyle}>
      <div
        aria-hidden="true"
        style={{
          position: 'absolute',
          right: '-48px',
          top: '-64px',
          width: '210px',
          height: '210px',
          borderRadius: '999px',
          background: `radial-gradient(circle, ${accent}55, transparent 68%)`,
          filter: 'blur(2px)',
        }}
      />
      <div style={innerStyle}>
        <div>
          <p style={eyebrowStyle}>{copy.eyebrow}</p>
          <h2 style={titleStyle}>{task.title}</h2>
          <p style={descriptionStyle}>{task.description}</p>
        </div>
        <div style={actionWrapStyle}>
          <span style={{ ...metaStyle, color: accent }}>{copy.kindLabels[task.kind]}</span>
          <button type="button" style={buttonStyle} onClick={() => onOpenTask(task.href)}>
            {task.cta_label}
          </button>
        </div>
      </div>
    </section>
  );
}
