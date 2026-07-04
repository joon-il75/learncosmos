'use client';

import type { CSSProperties } from 'react';

import type { LumiMessageType } from '@/lib/lumi/lumiTypes';

interface LumiBubbleProps {
  message: string;
  messageType: LumiMessageType;
  visible?: boolean;
}

const accentByType: Record<LumiMessageType, string> = {
  summary: 'rgba(126, 163, 255, 0.4)',
  hint: 'rgba(164, 198, 255, 0.38)',
  question: 'rgba(255, 215, 145, 0.42)',
  encourage: 'rgba(128, 230, 191, 0.42)',
  success: 'rgba(142, 241, 198, 0.44)',
  guide: 'rgba(151, 200, 255, 0.4)',
  discovery: 'rgba(255, 196, 122, 0.42)',
};

export default function LumiBubble({ message, messageType, visible = true }: LumiBubbleProps) {
  return (
    <div
      aria-live="polite"
      style={{
        ...bubbleStyle,
        borderColor: accentByType[messageType],
        opacity: visible ? 1 : 0,
        transform: visible ? 'translateY(0)' : 'translateY(6px)',
      }}
    >
      {message.split('\n').map((line, index) => (
        <span key={`${line}-${index}`} style={lineStyle}>
          {line}
        </span>
      ))}
    </div>
  );
}

const bubbleStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: 4,
  maxWidth: 280,
  padding: '12px 14px',
  borderRadius: 18,
  border: '1px solid rgba(126, 163, 255, 0.28)',
  background: 'linear-gradient(180deg, rgba(16, 25, 43, 0.94), rgba(11, 17, 30, 0.92))',
  boxShadow: '0 18px 42px rgba(0, 0, 0, 0.26)',
  color: 'rgba(241, 247, 255, 0.95)',
  fontSize: 13,
  lineHeight: 1.45,
  transition: 'opacity 180ms ease, transform 180ms ease',
};

const lineStyle: CSSProperties = {
  display: 'block',
};
