'use client';

import { useState, type CSSProperties } from 'react';
import type { LandingPageCopy } from '@/lib/i18n/pages/landing';
import ModalPortal from '@/components/common/ModalPortal';

interface BetaNoticeProps {
  copy: LandingPageCopy['alphaNotice'];
  onVideoClick: () => void;
  onDismiss: (suppressForToday?: boolean) => void;
}

export default function BetaNotice({ copy, onVideoClick, onDismiss }: BetaNoticeProps) {
  const [suppressForToday, setSuppressForToday] = useState(false);

  const handleDismiss = () => onDismiss(suppressForToday);
  const handleVideoClick = () => {
    if (suppressForToday) onDismiss(true);
    onVideoClick();
  };

  return (
    <ModalPortal overlayStyle={overlayStyle} onMouseDown={handleDismiss}>
      <section
        role="dialog"
        aria-modal="true"
        aria-label={copy.title}
        style={cardStyle}
        onMouseDown={(event) => event.stopPropagation()}
      >
        <p style={eyebrowStyle}>{copy.eyebrow}</p>
        <h2 style={titleStyle}>{copy.title}</h2>
        <div style={bodyStyle}>
          {copy.bodyLines.map((line) => (
            <p key={line} style={bodyLineStyle}>{line}</p>
          ))}
        </div>
        <label style={suppressLabelStyle}>
          <input
            type="checkbox"
            checked={suppressForToday}
            onChange={(event) => setSuppressForToday(event.target.checked)}
            style={suppressCheckboxStyle}
          />
          <span>하루동안 다시보지 않기</span>
        </label>
        <div style={actionsStyle}>
          <a href={copy.applicationHref} style={applicationButtonStyle}>
            {copy.applicationLabel}
          </a>
          <div style={secondaryActionsStyle}>
            <button type="button" style={videoButtonStyle} onClick={handleVideoClick}>
              {copy.videoLabel}
            </button>
            <button type="button" style={confirmButtonStyle} onClick={handleDismiss}>
              {copy.confirm}
            </button>
          </div>
        </div>
      </section>
    </ModalPortal>
  );
}

const overlayStyle: CSSProperties = {
  padding: '24px 16px',
  background:
    'radial-gradient(circle at 24% 42%, rgba(67, 83, 171, 0.24), transparent 24%), radial-gradient(circle at 74% 30%, rgba(77, 183, 232, 0.14), transparent 26%), rgba(2, 8, 20, 0.78)',
  backdropFilter: 'blur(9px) saturate(0.92)',
  WebkitBackdropFilter: 'blur(9px) saturate(0.92)',
};

const cardStyle: CSSProperties = {
  width: 'min(520px, calc(100vw - 32px))',
  borderRadius: 8,
  border: '1px solid rgba(248, 213, 116, 0.38)',
  background: 'rgba(7, 24, 34, 0.94)',
  color: '#F8FBFF',
  boxShadow: '0 34px 92px rgba(0, 0, 0, 0.42), inset 0 1px 0 rgba(255, 255, 255, 0.04)',
  padding: 'clamp(28px, 5.8vw, 30px)',
};

const eyebrowStyle: CSSProperties = {
  margin: '0 0 14px',
  color: '#FFD768',
  fontSize: 13,
  fontWeight: 900,
  letterSpacing: '0.18em',
  textTransform: 'uppercase',
};

const titleStyle: CSSProperties = {
  margin: '0 0 26px',
  color: '#FFFFFF',
  fontSize: 'clamp(28px, 5.4vw, 32px)',
  lineHeight: 1.14,
  fontWeight: 950,
  letterSpacing: '-0.035em',
};

const bodyStyle: CSSProperties = {
  display: 'grid',
  gap: 12,
  marginBottom: 28,
  color: 'rgba(232, 240, 248, 0.78)',
  fontSize: 15,
  lineHeight: 1.72,
  fontWeight: 600,
};

const bodyLineStyle: CSSProperties = {
  margin: 0,
};

const suppressLabelStyle: CSSProperties = {
  display: 'inline-flex',
  alignItems: 'center',
  gap: 9,
  margin: '-8px 0 18px',
  color: 'rgba(232, 240, 248, 0.76)',
  fontSize: 13,
  fontWeight: 760,
  cursor: 'pointer',
};

const suppressCheckboxStyle: CSSProperties = {
  width: 16,
  height: 16,
  accentColor: '#6AD2C1',
};

const actionsStyle: CSSProperties = {
  display: 'grid',
  gap: 12,
};

const applicationButtonStyle: CSSProperties = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: 48,
  borderRadius: 7,
  background: '#6AD2C1',
  color: '#061525',
  textDecoration: 'none',
  fontSize: 15,
  fontWeight: 900,
  boxShadow: '0 18px 38px rgba(106, 210, 193, 0.16)',
};

const secondaryActionsStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1fr) minmax(0, 1fr)',
  gap: 12,
};

const secondaryButtonBaseStyle: CSSProperties = {
  minHeight: 48,
  borderRadius: 7,
  fontSize: 15,
  fontWeight: 900,
  cursor: 'pointer',
};

const videoButtonStyle: CSSProperties = {
  ...secondaryButtonBaseStyle,
  border: '1px solid rgba(207, 223, 236, 0.24)',
  background: 'rgba(255, 255, 255, 0.06)',
  color: 'rgba(246, 250, 255, 0.88)',
};

const confirmButtonStyle: CSSProperties = {
  ...secondaryButtonBaseStyle,
  border: '1px solid rgba(255, 214, 104, 0.46)',
  background: '#FFD768',
  color: '#2A1B04',
};
