'use client';

import type { CSSProperties } from 'react';
import {
  pageBackgroundVideoStyle,
  pageBackgroundImageStyle,
  pageBackgroundOverlayStyle,
  starOverlayStyle,
} from './_styles';

function StarOverlay() {
  return (
    <svg aria-hidden="true" viewBox="0 0 1200 720" preserveAspectRatio="none" style={starOverlayStyle}>
      <g fill="rgba(241, 246, 255, 0.9)">
        <g>
          <circle cx="118" cy="110" r="1.7" />
          <circle cx="182" cy="302" r="2.3" />
          <circle cx="266" cy="560" r="1.2" />
          <circle cx="352" cy="214" r="3.2" />
          <circle cx="428" cy="136" r="1.5" />
          <animateTransform attributeName="transform" type="translate" values="0 0; 0 -4; 0 0" dur="8s" repeatCount="indefinite" />
        </g>
        <g>
          <circle cx="536" cy="410" r="2.5" />
          <circle cx="614" cy="174" r="1.4" />
          <circle cx="700" cy="598" r="3" />
          <circle cx="776" cy="92" r="1.1" />
          <circle cx="848" cy="210" r="2" />
          <animateTransform attributeName="transform" type="translate" values="0 0; 0 5; 0 0" dur="10s" repeatCount="indefinite" />
        </g>
        <g>
          <circle cx="922" cy="362" r="1.6" />
          <circle cx="1012" cy="144" r="3.4" />
          <circle cx="1086" cy="438" r="1.8" />
          <circle cx="1142" cy="264" r="2.7" />
          <circle cx="960" cy="560" r="1.2" />
          <animateTransform attributeName="transform" type="translate" values="0 0; 0 -3; 0 0" dur="6.5s" repeatCount="indefinite" />
        </g>
      </g>
    </svg>
  );
}

export default function DashboardBackground() {
  return (
    <>
      <style>{`
        @keyframes planetFloatA {
          0%, 100% { translate: 0 0; }
          50% { translate: 0 -10px; }
        }
        @keyframes planetFloatB {
          0%, 100% { translate: 0 0; }
          50% { translate: 0 12px; }
        }
        @keyframes planetFloatC {
          0%, 100% { translate: 0 0; }
          50% { translate: 0 -8px; }
        }
        @keyframes planetFloatD {
          0%, 100% { translate: 0 0; }
          50% { translate: 0 9px; }
        }
        @keyframes systemFloatA {
          0%, 100% { transform: translate(-50%, -50%) translateY(0px); }
          50% { transform: translate(-50%, -50%) translateY(-6px); }
        }
        @keyframes systemFloatB {
          0%, 100% { transform: translate(-50%, -50%) translateY(0px); }
          50% { transform: translate(-50%, -50%) translateY(5px); }
        }
        @keyframes systemFloatC {
          0%, 100% { transform: translate(-50%, -50%) translateY(0px); }
          50% { transform: translate(-50%, -50%) translateY(-4px); }
        }
        @keyframes activeSunPulse {
          0%, 100% { opacity: 0.72; filter: blur(0px) brightness(1); box-shadow: 0 0 18px rgba(255, 194, 95, 0.28); }
          50% { opacity: 1; filter: blur(1px) brightness(1.28); box-shadow: 0 0 30px rgba(255, 224, 155, 0.72), 0 0 44px rgba(255, 194, 95, 0.34); }
        }
        @keyframes tooltipFadeIn {
          from { opacity: 0; transform: translateY(-4px); }
          to   { opacity: 1; transform: translateY(0); }
        }
      `}</style>
      <video
        aria-hidden="true"
        style={pageBackgroundVideoStyle}
        src="/videos/LearnWeaver-Seamless-Loop.mp4"
        poster="/images/PlanetMap_background.webp"
        preload="metadata"
        autoPlay
        muted
        loop
        playsInline
      />
      <div aria-hidden="true" style={pageBackgroundImageStyle} />
      <div style={pageBackgroundOverlayStyle} />
      <StarOverlay />
    </>
  );
}
