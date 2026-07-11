'use client';

import { useEffect, useMemo, useRef, useState, type CSSProperties, type FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import * as THREE from 'three';
import { useGoalCreation } from '@/app/dashboard/useGoalCreation';
import { GoalFlowLoadingOverlay } from '@/components/goal-interview/GoalFlowLoadingOverlay';
import type { Locale } from '@/lib/i18n/locales';
import type { LandingPageCopy } from '@/lib/i18n/pages/landing';

type Props = {
  copy: LandingPageCopy['hero'];
  locale: Locale;
  loginHref: string;
};

type SampleStar = {
  label: string;
  left: number;
  top: number;
  size: number;
  tone: string;
  status: 'active' | 'open' | 'next' | 'done';
};

const sampleStars: SampleStar[] = [
  { label: '통기타 입문', left: 50, top: 42, size: 38, tone: '#8B7CFC', status: 'active' },
  { label: '수채화 기초', left: 24, top: 28, size: 18, tone: '#75D7F0', status: 'done' },
  { label: '코바늘 기초', left: 73, top: 25, size: 17, tone: '#41C7B3', status: 'open' },
  { label: '라떼아트', left: 31, top: 64, size: 14, tone: '#F8C85E', status: 'next' },
  { label: '사진촬영', left: 78, top: 63, size: 18, tone: '#A38BFF', status: 'open' },
  { label: '새 학습별', left: 47, top: 74, size: 12, tone: '#7DD3FC', status: 'next' },
];

const guideVideoUrl = '';
const betaSurveyUrl = '';
const betaExperienceItems = {
  ko: [
    'AI가 주제를 하나의 학습별과 항로로 설계하는 흐름',
    '은하 지도에서 코스 별을 둘러보고 선택하는 경험',
    '학습계획, 탐험일지, 결과물 기록으로 이어지는 구조',
    '완주 중심 자기주도 학습 경험과 초기 피드백 반영',
  ],
  en: [
    'AI turns a topic into a learning star and route.',
    'Explore and select course stars on the galaxy map.',
    'Connect plans, learning logs, and artifacts into one flow.',
    'Help improve a completion-focused self-directed learning beta.',
  ],
} as const;

const routeLinks = [
  [1, 0],
  [0, 2],
  [0, 3],
  [3, 5],
  [0, 4],
] as const;

const courseStarPlacements = [
  { arm: 0, radius: 1.08, offset: 0.08 },
  { arm: 2, radius: 1.56, offset: -0.14 },
  { arm: 1, radius: 2.02, offset: 0.12 },
  { arm: 0, radius: 2.34, offset: -0.18 },
  { arm: 2, radius: 2.72, offset: 0.08 },
  { arm: 1, radius: 3.08, offset: -0.06 },
] as const;

const rootStyle: CSSProperties = {
  minHeight: '100svh',
  padding: '68px 0 12px',
  background: 'linear-gradient(180deg, #FBFFFE 0%, #F4FBFF 42%, #F8F4FF 100%)',
  color: '#12384E',
  overflowX: 'hidden',
};

const shellStyle: CSSProperties = {
  width: 'calc(100% - 12px)',
  maxWidth: '720px',
  minHeight: 'calc(100svh - 80px)',
  margin: '0 auto',
  display: 'grid',
  gridTemplateRows: 'auto minmax(430px, 1fr)',
  gap: '8px',
};

const mapFrameStyle: CSSProperties = {
  width: '100vw',
  marginLeft: 'calc(50% - 50vw)',
  marginRight: 'calc(50% - 50vw)',
};

const lumiCardStyle: CSSProperties = {
  border: '1px solid rgba(151, 190, 210, 0.34)',
  borderRadius: '16px',
  padding: '9px 10px',
  background: 'rgba(255, 255, 255, 0.82)',
  boxShadow: '0 12px 36px rgba(35, 79, 112, 0.08)',
};

const eyebrowStyle: CSSProperties = {
  margin: 0,
  color: '#6C7ABF',
  fontSize: '10px',
  fontWeight: 900,
  letterSpacing: '0.12em',
  lineHeight: 1.2,
  textTransform: 'uppercase',
};

const lumiTextStyle: CSSProperties = {
  margin: '5px 0 0',
  color: '#12384E',
  fontSize: '12px',
  lineHeight: 1.32,
  fontWeight: 850,
};

const mapStyle: CSSProperties = {
  position: 'relative',
  minHeight: '430px',
  overflow: 'hidden',
  border: '1px solid rgba(174, 198, 230, 0.40)',
  borderRadius: '22px',
  background:
    'radial-gradient(ellipse at 50% 30%, rgba(160, 108, 255, 0.18), transparent 42%), radial-gradient(ellipse at 22% 74%, rgba(82, 57, 166, 0.20), transparent 36%), radial-gradient(ellipse at 78% 16%, rgba(73, 48, 145, 0.24), transparent 34%), linear-gradient(180deg, #120B2D 0%, #1A1040 46%, #0B0822 100%)',
  boxShadow: 'inset 0 1px 0 rgba(213,196,255,0.16), 0 18px 48px rgba(48, 30, 102, 0.20)',
};

const threeLayerStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  zIndex: 1,
  cursor: 'grab',
};

const dustStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  zIndex: 2,
  opacity: 0.07,
  backgroundImage:
    'radial-gradient(circle, rgba(210, 198, 255, 0.26) 0 1px, transparent 1.6px), radial-gradient(circle, rgba(148, 103, 255, 0.22) 0 1px, transparent 1.7px)',
  backgroundPosition: '4px 10px, 21px 28px',
  backgroundSize: '38px 42px, 54px 58px',
  pointerEvents: 'none',
};

const ctaOverlayStyle: CSSProperties = {
  position: 'absolute',
  left: '10px',
  right: '10px',
  top: '10px',
  zIndex: 7,
  display: 'grid',
  gap: '8px',
  padding: '10px',
  border: '1px solid rgba(255,255,255,0.18)',
  borderRadius: '16px',
  background: 'linear-gradient(180deg, rgba(18, 11, 45, 0.74), rgba(18, 11, 45, 0.56))',
  boxShadow: '0 18px 42px rgba(19, 12, 46, 0.28)',
  backdropFilter: 'blur(12px)',
};

const compactCtaOverlayStyle: CSSProperties = {
  gap: '6px',
  padding: '8px',
};

const formStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1fr) auto',
  gap: '8px',
  alignItems: 'center',
};

const compactFormStyle: CSSProperties = {
  gap: '6px',
};

const inputStyle: CSSProperties = {
  minWidth: 0,
  height: '40px',
  border: '1px solid rgba(200, 224, 255, 0.42)',
  borderRadius: '13px',
  padding: '0 11px',
  background: 'rgba(255,255,255,0.9)',
  color: '#12384E',
  fontFamily: 'inherit',
  fontSize: '12px',
  fontWeight: 750,
  outline: 'none',
};

const compactInputStyle: CSSProperties = {
  height: '36px',
  borderRadius: '12px',
  padding: '0 9px',
  fontSize: '11px',
};

const buttonStyle: CSSProperties = {
  height: '40px',
  minWidth: '84px',
  border: 0,
  borderRadius: '13px',
  background: 'linear-gradient(135deg, #41C7B3 0%, #6E8CFB 100%)',
  color: '#FFFFFF',
  fontFamily: 'inherit',
  fontSize: '12px',
  fontWeight: 950,
  cursor: 'pointer',
  boxShadow: '0 12px 28px rgba(65, 199, 179, 0.24)',
};

const compactButtonStyle: CSSProperties = {
  height: '36px',
  minWidth: '74px',
  borderRadius: '12px',
  fontSize: '11px',
};

const chipRowStyle: CSSProperties = {
  display: 'flex',
  gap: '6px',
  overflowX: 'auto',
  paddingBottom: '2px',
};

const chipStyle: CSSProperties = {
  flex: '0 0 auto',
  height: '32px',
  border: '1px solid rgba(151, 190, 210, 0.34)',
  borderRadius: '999px',
  padding: '0 10px',
  background: 'rgba(255,255,255,0.86)',
  color: '#31576C',
  fontFamily: 'inherit',
  fontSize: '11px',
  fontWeight: 850,
  cursor: 'pointer',
};

const compactChipStyle: CSSProperties = {
  height: '28px',
  padding: '0 8px',
  fontSize: '10px',
};

const videoSectionStyle: CSSProperties = {
  display: 'grid',
  gap: '10px',
  padding: '14px 0 4px',
};

const videoHeaderStyle: CSSProperties = {
  display: 'grid',
  gap: '4px',
  padding: '0 4px',
};

const videoEyebrowStyle: CSSProperties = {
  margin: 0,
  color: '#6C7ABF',
  fontSize: '10px',
  fontWeight: 900,
  letterSpacing: '0.12em',
  lineHeight: 1.2,
  textTransform: 'uppercase',
};

const videoTitleStyle: CSSProperties = {
  margin: 0,
  color: '#12384E',
  fontSize: '17px',
  lineHeight: 1.2,
  fontWeight: 950,
};

const videoTextStyle: CSSProperties = {
  margin: 0,
  color: 'rgba(231, 241, 255, 0.80)',
  fontSize: '12px',
  lineHeight: 1.42,
  fontWeight: 750,
};

const videoPreviewStyle: CSSProperties = {
  position: 'relative',
  display: 'grid',
  placeItems: 'center',
  minHeight: '190px',
  overflow: 'hidden',
  border: '1px solid rgba(151, 190, 210, 0.34)',
  borderRadius: '18px',
  background:
    'radial-gradient(circle at 50% 48%, rgba(255,255,255,0.78), transparent 0 4px), radial-gradient(ellipse at 52% 48%, rgba(100,223,255,0.36), transparent 0 18%), radial-gradient(ellipse at 50% 52%, rgba(160,108,255,0.28), transparent 0 44%), linear-gradient(135deg, #140D31 0%, #25145A 52%, #0F1238 100%)',
  boxShadow: '0 18px 42px rgba(48, 30, 102, 0.14)',
  textDecoration: 'none',
};

const videoPlayStyle: CSSProperties = {
  position: 'relative',
  zIndex: 2,
  display: 'grid',
  placeItems: 'center',
  width: '58px',
  height: '58px',
  borderRadius: '999px',
  border: '1px solid rgba(255,255,255,0.52)',
  background: 'linear-gradient(135deg, rgba(65,199,179,0.94), rgba(110,140,251,0.94))',
  color: '#FFFFFF',
  fontSize: '24px',
  fontWeight: 950,
  boxShadow: '0 16px 34px rgba(65, 199, 179, 0.28)',
};

const videoBadgeStyle: CSSProperties = {
  position: 'absolute',
  left: '12px',
  right: '12px',
  bottom: '12px',
  zIndex: 2,
  display: 'flex',
  justifyContent: 'space-between',
  gap: '10px',
  alignItems: 'center',
  padding: '9px 10px',
  border: '1px solid rgba(255,255,255,0.18)',
  borderRadius: '14px',
  background: 'rgba(18, 11, 45, 0.68)',
  color: '#F7FBFF',
  fontSize: '11px',
  lineHeight: 1.2,
  fontWeight: 850,
  backdropFilter: 'blur(10px)',
};

const betaSectionStyle: CSSProperties = {
  display: 'grid',
  gap: '10px',
  marginTop: '10px',
  padding: '14px',
  border: '1px solid rgba(104, 214, 197, 0.30)',
  borderRadius: '18px',
  background:
    'radial-gradient(circle at 18% 18%, rgba(65,199,179,0.20), transparent 32%), radial-gradient(circle at 82% 10%, rgba(248,200,94,0.16), transparent 30%), radial-gradient(circle at 48% 110%, rgba(110,140,251,0.26), transparent 44%), linear-gradient(180deg, #071827 0%, #111938 100%)',
  boxShadow: '0 18px 46px rgba(22, 34, 75, 0.18), inset 0 1px 0 rgba(255,255,255,0.10)',
};

const betaLabelStyle: CSSProperties = {
  margin: 0,
  color: '#F8C85E',
  fontSize: '10px',
  fontWeight: 950,
  letterSpacing: '0.12em',
  lineHeight: 1.2,
  textTransform: 'uppercase',
};

const betaTitleStyle: CSSProperties = {
  margin: 0,
  color: '#F7FBFF',
  fontSize: '18px',
  lineHeight: 1.2,
  fontWeight: 950,
};

const betaTextStyle: CSSProperties = {
  margin: 0,
  color: 'rgba(231, 241, 255, 0.80)',
  fontSize: '12px',
  lineHeight: 1.45,
  fontWeight: 750,
};

const betaListStyle: CSSProperties = {
  display: 'grid',
  gap: '8px',
  margin: '2px 0 0',
};

const betaItemStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '20px minmax(0, 1fr)',
  gap: '8px',
  alignItems: 'start',
  padding: '9px 10px',
  border: '1px solid rgba(255, 255, 255, 0.12)',
  borderRadius: '12px',
  background: 'rgba(255,255,255,0.07)',
  color: '#EAF5FF',
  fontSize: '12px',
  lineHeight: 1.38,
  fontWeight: 800,
};

const betaCheckStyle: CSSProperties = {
  display: 'inline-grid',
  placeItems: 'center',
  width: '18px',
  height: '18px',
  borderRadius: '999px',
  background: 'linear-gradient(135deg, #F8C85E, #41C7B3)',
  color: '#071827',
  fontSize: '11px',
  fontWeight: 950,
};

const betaActionsStyle: CSSProperties = {
  display: 'flex',
  gap: '8px',
  flexWrap: 'wrap',
  alignItems: 'center',
};

const betaLinkStyle: CSSProperties = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '38px',
  padding: '0 14px',
  borderRadius: '13px',
  background: 'linear-gradient(135deg, #41C7B3 0%, #6E8CFB 100%)',
  color: '#FFFFFF',
  fontSize: '12px',
  fontWeight: 950,
  textDecoration: 'none',
  boxShadow: '0 12px 28px rgba(65, 199, 179, 0.20)',
};

const betaNoteStyle: CSSProperties = {
  margin: 0,
  color: 'rgba(231, 241, 255, 0.68)',
  fontSize: '11px',
  lineHeight: 1.35,
  fontWeight: 750,
};

const miniFooterStyle: CSSProperties = {
  display: 'grid',
  gap: '10px',
  marginTop: '12px',
  padding: '16px 10px 6px',
  textAlign: 'center',
  color: '#5A7485',
};

const miniFooterBrandStyle: CSSProperties = {
  margin: 0,
  color: '#12384E',
  fontSize: '15px',
  lineHeight: 1.2,
  fontWeight: 950,
};

const miniFooterTaglineStyle: CSSProperties = {
  margin: 0,
  color: '#6C7ABF',
  fontSize: '11px',
  lineHeight: 1.35,
  fontWeight: 800,
};

const miniFooterLinksStyle: CSSProperties = {
  display: 'flex',
  justifyContent: 'center',
  gap: '8px 10px',
  flexWrap: 'wrap',
  marginTop: '2px',
};

const miniFooterLinkStyle: CSSProperties = {
  color: '#31576C',
  fontSize: '11px',
  lineHeight: 1.25,
  fontWeight: 850,
  textDecoration: 'none',
};

const miniFooterContactStyle: CSSProperties = {
  color: '#31576C',
  fontSize: '11px',
  lineHeight: 1.25,
  fontWeight: 850,
};

const miniFooterMetaStyle: CSSProperties = {
  margin: 0,
  color: '#7A92A0',
  fontSize: '10px',
  lineHeight: 1.35,
  fontWeight: 750,
};

function getGalaxyArmPoint(arm: number, radius: number, offset = 0, z = 0) {
  const angle = arm * ((Math.PI * 2) / 3) + radius * 1.36 + offset;
  return new THREE.Vector3(Math.cos(angle) * radius, Math.sin(angle) * radius * 0.28, z);
}

function getCourseStarPosition(index: number) {
  const placement = courseStarPlacements[index % courseStarPlacements.length];
  return getGalaxyArmPoint(placement.arm, placement.radius, placement.offset, -1.04);
}

function createStarTexture(color: string, active = false) {
  const pixelSize = 9;
  const grid = 16;
  const canvas = document.createElement('canvas');
  canvas.width = grid * pixelSize;
  canvas.height = grid * pixelSize;
  const ctx = canvas.getContext('2d');
  if (!ctx) return new THREE.CanvasTexture(canvas);

  ctx.imageSmoothingEnabled = false;
  const rgb = new THREE.Color(color);
  const r = Math.round(rgb.r * 255);
  const g = Math.round(rgb.g * 255);
  const b = Math.round(rgb.b * 255);
  const px = (x: number, y: number, w: number, h: number, fill: string) => {
    ctx.fillStyle = fill;
    ctx.fillRect(x * pixelSize, y * pixelSize, w * pixelSize, h * pixelSize);
  };

  const glow = `rgba(${r},${g},${b},0.18)`;
  const soft = `rgba(${r},${g},${b},0.42)`;
  const main = `rgba(${r},${g},${b},0.96)`;
  const hot = 'rgba(255,255,255,0.98)';
  const edge = active ? 'rgba(255,244,190,0.92)' : 'rgba(222,232,255,0.78)';

  px(7, 1, 2, 14, glow);
  px(1, 7, 14, 2, glow);
  px(5, 3, 6, 10, soft);
  px(3, 5, 10, 6, soft);
  px(6, 4, 4, 8, main);
  px(4, 6, 8, 4, main);
  px(7, 3, 2, 10, edge);
  px(3, 7, 10, 2, edge);
  px(6, 6, 4, 4, hot);
  px(7, 7, 2, 2, 'rgba(255,255,255,1)');

  if (active) {
    px(2, 2, 2, 2, 'rgba(255,220,118,0.7)');
    px(12, 3, 1, 2, 'rgba(255,220,118,0.64)');
    px(3, 12, 2, 1, 'rgba(255,220,118,0.58)');
    px(12, 11, 2, 2, 'rgba(255,220,118,0.58)');
  }

  const texture = new THREE.CanvasTexture(canvas);
  texture.colorSpace = THREE.SRGBColorSpace;
  texture.magFilter = THREE.NearestFilter;
  texture.minFilter = THREE.NearestFilter;
  texture.generateMipmaps = false;
  return texture;
}

function createCourseHaloTexture(color: string) {
  const canvas = document.createElement('canvas');
  canvas.width = 192;
  canvas.height = 192;
  const ctx = canvas.getContext('2d');
  if (!ctx) return new THREE.CanvasTexture(canvas);

  const rgb = new THREE.Color(color);
  const r = Math.round(rgb.r * 255);
  const g = Math.round(rgb.g * 255);
  const b = Math.round(rgb.b * 255);
  const gradient = ctx.createRadialGradient(96, 96, 0, 96, 96, 92);
  gradient.addColorStop(0, `rgba(255,255,255,0.44)`);
  gradient.addColorStop(0.22, `rgba(${r},${g},${b},0.34)`);
  gradient.addColorStop(0.58, `rgba(${r},${g},${b},0.12)`);
  gradient.addColorStop(1, `rgba(${r},${g},${b},0)`);
  ctx.fillStyle = gradient;
  ctx.fillRect(0, 0, 192, 192);

  const texture = new THREE.CanvasTexture(canvas);
  texture.colorSpace = THREE.SRGBColorSpace;
  return texture;
}

function createCourseRingTexture(color: string) {
  const canvas = document.createElement('canvas');
  canvas.width = 192;
  canvas.height = 192;
  const ctx = canvas.getContext('2d');
  if (!ctx) return new THREE.CanvasTexture(canvas);

  const rgb = new THREE.Color(color);
  const r = Math.round(rgb.r * 255);
  const g = Math.round(rgb.g * 255);
  const b = Math.round(rgb.b * 255);
  ctx.translate(96, 96);
  ctx.rotate(-0.28);
  ctx.scale(1, 0.52);
  ctx.strokeStyle = `rgba(${r},${g},${b},0.78)`;
  ctx.lineWidth = 4;
  ctx.beginPath();
  ctx.arc(0, 0, 54, 0, Math.PI * 2);
  ctx.stroke();
  ctx.strokeStyle = 'rgba(255,255,255,0.62)';
  ctx.lineWidth = 1.4;
  ctx.beginPath();
  ctx.arc(0, 0, 42, Math.PI * 0.18, Math.PI * 1.28);
  ctx.stroke();

  const texture = new THREE.CanvasTexture(canvas);
  texture.colorSpace = THREE.SRGBColorSpace;
  return texture;
}

function createLabelTexture(label: string) {
  const canvas = document.createElement('canvas');
  canvas.width = 256;
  canvas.height = 72;
  const ctx = canvas.getContext('2d');
  if (!ctx) return new THREE.CanvasTexture(canvas);

  ctx.fillStyle = 'rgba(255, 255, 255, 0.78)';
  ctx.strokeStyle = 'rgba(151, 190, 210, 0.38)';
  ctx.lineWidth = 2;
  ctx.beginPath();
  ctx.roundRect(8, 12, 240, 44, 18);
  ctx.fill();
  ctx.stroke();
  ctx.fillStyle = '#12384E';
  ctx.font = '700 23px sans-serif';
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';
  ctx.fillText(label, 128, 34, 210);

  const texture = new THREE.CanvasTexture(canvas);
  texture.colorSpace = THREE.SRGBColorSpace;
  return texture;
}

function seededUnit(index: number, salt: number) {
  const value = Math.sin(index * salt) * 10000;
  return value - Math.floor(value);
}

function createParticleTexture() {
  const canvas = document.createElement('canvas');
  canvas.width = 32;
  canvas.height = 32;
  const ctx = canvas.getContext('2d');
  if (!ctx) return new THREE.CanvasTexture(canvas);

  ctx.imageSmoothingEnabled = false;
  const px = (x: number, y: number, w: number, h: number, fill: string) => {
    ctx.fillStyle = fill;
    ctx.fillRect(x, y, w, h);
  };

  px(14, 14, 4, 4, 'rgba(255,255,255,1)');
  px(12, 14, 2, 4, 'rgba(210,235,255,0.58)');
  px(18, 14, 2, 4, 'rgba(210,235,255,0.5)');
  px(14, 12, 4, 2, 'rgba(210,235,255,0.52)');
  px(14, 18, 4, 2, 'rgba(210,235,255,0.42)');
  px(10, 15, 2, 2, 'rgba(160,145,255,0.28)');
  px(20, 15, 2, 2, 'rgba(160,145,255,0.24)');
  px(15, 10, 2, 2, 'rgba(160,145,255,0.22)');
  px(15, 20, 2, 2, 'rgba(160,145,255,0.2)');

  const texture = new THREE.CanvasTexture(canvas);
  texture.colorSpace = THREE.SRGBColorSpace;
  texture.magFilter = THREE.NearestFilter;
  texture.minFilter = THREE.NearestFilter;
  texture.generateMipmaps = false;
  return texture;
}

type ReferenceGalaxyMode = 'disk' | 'innerDust' | 'blueRim' | 'outerArm' | 'core' | 'upperPulse' | 'lowerSparkles';

function createReferenceGalaxyLayer(count: number, texture: THREE.Texture, mode: ReferenceGalaxyMode, isMobileViewport: boolean) {
  const positions = new Float32Array(count * 3);
  const colors = new Float32Array(count * 3);
  const white = new THREE.Color('#FFFFFF');
  const violet = new THREE.Color('#A06CFF');
  const magenta = new THREE.Color('#E57BFF');
  const cyan = new THREE.Color('#64DFFF');
  const blue = new THREE.Color('#3787FF');
  const warm = new THREE.Color('#FFB777');
  const peach = new THREE.Color('#FFE0C2');

  for (let index = 0; index < count; index += 1) {
    let x = 0;
    let y = 0;
    let z = -1.4;
    let color = white.clone();

    if (mode === 'disk') {
      const rSeed = seededUnit(index + 7, 13.17);
      const radius = Math.pow(rSeed, 0.42) * 3.72;
      const turn = radius * 0.72;
      const angle = seededUnit(index + 11, 19.91) * Math.PI * 2 + turn;
      const band = (seededUnit(index + 23, 29.71) - 0.5) * (0.1 + radius * 0.12);
      x = Math.cos(angle) * radius + Math.cos(angle + Math.PI / 2) * band;
      y = Math.sin(angle) * radius * 0.24 + Math.sin(angle + Math.PI / 2) * band * 0.12;
      z = -1.48 + (seededUnit(index + 31, 37.7) - 0.5) * 0.46;
      const center = Math.max(0, 1 - radius / 3.72);
      color = cyan.clone().lerp(blue, radius / 3.72 * 0.38);
      color.lerp(violet, Math.pow(center, 1.2) * 0.62);
      color.lerp(peach, Math.pow(center, 3.2) * 0.46);
    } else if (mode === 'innerDust') {
      const t = seededUnit(index + 157, 113.7);
      const radius = 0.42 + Math.pow(t, 0.7) * 1.76;
      const angle = -0.28 + radius * 0.56 + (seededUnit(index + 163, 127.1) - 0.5) * 1.35;
      const band = (seededUnit(index + 167, 131.3) - 0.5) * (0.08 + radius * 0.16);
      x = Math.cos(angle) * radius * 1.38 + Math.cos(angle + Math.PI / 2) * band;
      y = Math.sin(angle) * radius * 0.2 + Math.sin(angle + Math.PI / 2) * band * 0.1;
      z = -1.02 + (seededUnit(index + 173, 137.9) - 0.5) * 0.16;
      color = warm.clone().lerp(peach, 0.28).lerp(magenta, t * 0.18).lerp(white, seededUnit(index + 179, 139.1) * 0.12);
    } else if (mode === 'blueRim') {
      const t = seededUnit(index + 181, 149.3);
      const side = index % 2 === 0 ? 1 : -1;
      const angle = side > 0 ? -0.58 + t * 2.85 : Math.PI + 0.45 + t * 2.5;
      const radius = 2.25 + Math.pow(t, 0.82) * 1.92;
      const width = (seededUnit(index + 191, 151.7) - 0.5) * (0.08 + t * 0.34);
      x = Math.cos(angle) * radius + Math.cos(angle + Math.PI / 2) * width;
      y = Math.sin(angle) * radius * 0.26 + Math.sin(angle + Math.PI / 2) * width * 0.12;
      z = -1.24 + (seededUnit(index + 193, 157.9) - 0.5) * 0.22;
      color = cyan.clone().lerp(blue, t * 0.28).lerp(white, seededUnit(index + 197, 163.3) * 0.2);
    } else if (mode === 'outerArm') {
      const side = index % 2 === 0 ? 1 : -1;
      const t = seededUnit(index + 41, 11.33);
      const angle = side > 0 ? -0.2 + t * 2.55 : Math.PI + 0.18 + t * 2.18;
      const radius = 1.28 + Math.pow(t, 0.78) * 3.38;
      const width = (seededUnit(index + 53, 17.19) - 0.5) * (0.2 + t * 0.48);
      x = Math.cos(angle) * radius + Math.cos(angle + Math.PI / 2) * width;
      y = Math.sin(angle) * radius * 0.28 + Math.sin(angle + Math.PI / 2) * width * 0.18;
      z = -1.34 + (seededUnit(index + 61, 23.57) - 0.5) * 0.38;
      color = cyan.clone().lerp(blue, t * 0.5).lerp(white, seededUnit(index + 67, 31.7) * 0.16);
    } else if (mode === 'core') {
      const radius = Math.pow(seededUnit(index + 71, 41.13), 1.8) * 0.92;
      const angle = seededUnit(index + 79, 43.61) * Math.PI * 2;
      x = Math.cos(angle) * radius * 1.55 + (seededUnit(index + 83, 47.89) - 0.5) * 0.18;
      y = Math.sin(angle) * radius * 0.33 + (seededUnit(index + 89, 53.21) - 0.5) * 0.08;
      z = -1.06 + (seededUnit(index + 97, 59.31) - 0.5) * 0.2;
      color = warm.clone().lerp(peach, 0.42).lerp(white, 0.28).lerp(violet, seededUnit(index + 101, 61.11) * 0.18);
    } else if (mode === 'upperPulse') {
      const upward = index % 3 !== 0 ? 1 : -1;
      const t = seededUnit(index + 103, 67.13);
      const height = Math.pow(t, 0.72) * (upward > 0 ? 2.46 : 1.5);
      const taper = Math.max(0.06, 1 - height / 2.46);
      const angle = seededUnit(index + 107, 71.41) * Math.PI * 2;
      const radius = Math.pow(seededUnit(index + 109, 73.17), 1.9) * 0.46 * taper;
      x = Math.cos(angle) * radius + (seededUnit(index + 113, 79.37) - 0.5) * 0.18 * (1 - taper);
      y = Math.sin(angle) * radius * 0.24 + (seededUnit(index + 127, 83.43) - 0.5) * 0.12;
      z = -1.18 + upward * (0.16 + height);
      color = magenta.clone().lerp(violet, t * 0.42).lerp(cyan, (1 - taper) * 0.22).lerp(white, taper * 0.2);
    } else {
      const t = seededUnit(index + 131, 89.11);
      const angle = 3.75 + t * 1.35;
      const radius = 2.1 + t * 1.55;
      const width = (seededUnit(index + 137, 97.7) - 0.5) * 0.42;
      x = Math.cos(angle) * radius + width;
      y = Math.sin(angle) * radius * 0.34 - 0.35 + (seededUnit(index + 139, 101.9) - 0.5) * 0.28;
      z = -1.12 + (seededUnit(index + 149, 107.3) - 0.5) * 0.28;
      color = cyan.clone().lerp(white, seededUnit(index + 151, 109.9) * 0.42).lerp(blue, 0.22);
    }

    positions[index * 3] = x;
    positions[index * 3 + 1] = y;
    positions[index * 3 + 2] = z;
    colors[index * 3] = color.r;
    colors[index * 3 + 1] = color.g;
    colors[index * 3 + 2] = color.b;
  }

  const geometry = new THREE.BufferGeometry();
  geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
  geometry.setAttribute('color', new THREE.BufferAttribute(colors, 3));

  const sizeByMode: Record<ReferenceGalaxyMode, number> = {
    disk: isMobileViewport ? 0.044 : 0.034,
    innerDust: isMobileViewport ? 0.064 : 0.05,
    blueRim: isMobileViewport ? 0.074 : 0.058,
    outerArm: isMobileViewport ? 0.066 : 0.052,
    core: isMobileViewport ? 0.078 : 0.064,
    upperPulse: isMobileViewport ? 0.07 : 0.056,
    lowerSparkles: isMobileViewport ? 0.092 : 0.074,
  };
  const opacityByMode: Record<ReferenceGalaxyMode, number> = {
    disk: 0.66,
    innerDust: 0.7,
    blueRim: 0.74,
    outerArm: 0.78,
    core: 0.96,
    upperPulse: 0.44,
    lowerSparkles: 0.72,
  };

  return new THREE.Points(
    geometry,
    new THREE.PointsMaterial({
      map: texture,
      size: sizeByMode[mode],
      transparent: true,
      opacity: opacityByMode[mode],
      vertexColors: true,
      blending: THREE.AdditiveBlending,
      depthWrite: false,
      alphaTest: 0.02,
    }),
  );
}

function createBackgroundStarfield(count: number, texture: THREE.Texture) {
  const positions = new Float32Array(count * 3);
  const colors = new Float32Array(count * 3);
  const palette = [new THREE.Color('#FFFFFF'), new THREE.Color('#BDEFFF'), new THREE.Color('#CEC2FF')];

  for (let index = 0; index < count; index += 1) {
    positions[index * 3] = (seededUnit(index + 5, 7.31) - 0.5) * 9.8;
    positions[index * 3 + 1] = (seededUnit(index + 9, 11.77) - 0.5) * 5.6;
    positions[index * 3 + 2] = -3.4 - seededUnit(index + 12, 19.17) * 1.8;
    const color = palette[index % palette.length];
    colors[index * 3] = color.r;
    colors[index * 3 + 1] = color.g;
    colors[index * 3 + 2] = color.b;
  }

  const geometry = new THREE.BufferGeometry();
  geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
  geometry.setAttribute('color', new THREE.BufferAttribute(colors, 3));

  return new THREE.Points(
    geometry,
    new THREE.PointsMaterial({
      map: texture,
      size: 0.026,
      transparent: true,
      opacity: 0.42,
      vertexColors: true,
      blending: THREE.AdditiveBlending,
      depthWrite: false,
    }),
  );
}

function createGalaxyCoreTexture() {
  const pixelSize = 8;
  const grid = 32;
  const canvas = document.createElement('canvas');
  canvas.width = grid * pixelSize;
  canvas.height = grid * pixelSize;
  const ctx = canvas.getContext('2d');
  if (!ctx) return new THREE.CanvasTexture(canvas);

  ctx.imageSmoothingEnabled = false;
  const px = (x: number, y: number, w: number, h: number, fill: string) => {
    ctx.fillStyle = fill;
    ctx.fillRect(x * pixelSize, y * pixelSize, w * pixelSize, h * pixelSize);
  };

  px(9, 13, 14, 6, 'rgba(255,178,112,0.16)');
  px(7, 14, 18, 4, 'rgba(160,108,255,0.16)');
  px(11, 12, 10, 8, 'rgba(255,210,158,0.32)');
  px(12, 13, 8, 6, 'rgba(255,238,213,0.72)');
  px(14, 14, 4, 4, 'rgba(255,255,255,0.98)');
  px(10, 15, 2, 2, 'rgba(229,123,255,0.5)');
  px(21, 15, 2, 2, 'rgba(100,223,255,0.45)');
  px(15, 10, 2, 2, 'rgba(229,123,255,0.34)');
  px(15, 20, 2, 2, 'rgba(100,223,255,0.3)');

  const texture = new THREE.CanvasTexture(canvas);
  texture.colorSpace = THREE.SRGBColorSpace;
  texture.magFilter = THREE.NearestFilter;
  texture.minFilter = THREE.NearestFilter;
  texture.generateMipmaps = false;
  return texture;
}

type ThreeCosmosMapProps = {
  selectedLabel: string | null;
  onSelect: (label: string | null) => void;
};

function ThreeCosmosMap({ selectedLabel, onSelect }: ThreeCosmosMapProps) {
  const mountRef = useRef<HTMLDivElement | null>(null);
  const onSelectRef = useRef(onSelect);
  const selectedLabelRef = useRef(selectedLabel);

  useEffect(() => {
    onSelectRef.current = onSelect;
  }, [onSelect]);

  useEffect(() => {
    selectedLabelRef.current = selectedLabel;
  }, [selectedLabel]);

  useEffect(() => {
    const mount = mountRef.current;
    if (!mount) return undefined;

    const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(42, 1, 0.1, 100);
    camera.position.set(0, 0.08, 6.85);

    const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true, powerPreference: 'high-performance' });
    renderer.setClearColor(0xffffff, 0);
    renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 1.7));
    renderer.domElement.style.width = '100%';
    renderer.domElement.style.height = '100%';
    renderer.domElement.style.display = 'block';
    mount.appendChild(renderer.domElement);

    const cameraTarget = new THREE.Vector3(0, 0.08, 6.85);
    const pointer = new THREE.Vector2();
    const raycaster = new THREE.Raycaster();
    raycaster.params.Sprite = { threshold: 0.16 };
    const clock = new THREE.Clock();
    const galaxyRoot = new THREE.Group();
    galaxyRoot.position.y = 0.32;
    scene.add(galaxyRoot);
    const dragState = {
      active: false,
      moved: false,
      pointerId: -1,
      startX: 0,
      startY: 0,
      lastX: 0,
      lastY: 0,
      rotationX: -0.82,
      rotationY: 0.08,
      targetX: -0.82,
      targetY: 0.08,
    };
    galaxyRoot.rotation.x = dragState.rotationX;
    galaxyRoot.rotation.y = dragState.rotationY;
    const starSprites: THREE.Sprite[] = [];
    const courseHalos: THREE.Sprite[] = [];
    const courseRings: THREE.Sprite[] = [];
    const labelSprites: THREE.Sprite[] = [];
    const starTextures: THREE.Texture[] = [];
    const labelTextures: THREE.Texture[] = [];
    const effectTextures: THREE.Texture[] = [];

    const ambient = new THREE.AmbientLight(0xffffff, 1.1);
    scene.add(ambient);

    const isMobileViewport = window.innerWidth < 480;
    const particleTexture = createParticleTexture();
    const backgroundStars = createBackgroundStarfield(isMobileViewport ? 1200 : 2400, particleTexture);
    galaxyRoot.add(backgroundStars);



    const galaxyDisk = createReferenceGalaxyLayer(isMobileViewport ? 7200 : 15600, particleTexture, 'disk', isMobileViewport);
    galaxyDisk.rotation.z = -0.18;
    galaxyRoot.add(galaxyDisk);

    const innerDust = createReferenceGalaxyLayer(isMobileViewport ? 1100 : 2600, particleTexture, 'innerDust', isMobileViewport);
    innerDust.rotation.z = -0.18;
    galaxyRoot.add(innerDust);

    const blueRim = createReferenceGalaxyLayer(isMobileViewport ? 1300 : 3000, particleTexture, 'blueRim', isMobileViewport);
    blueRim.rotation.z = -0.18;
    galaxyRoot.add(blueRim);

    const outerArm = createReferenceGalaxyLayer(isMobileViewport ? 1800 : 4400, particleTexture, 'outerArm', isMobileViewport);
    outerArm.rotation.z = -0.18;
    galaxyRoot.add(outerArm);

    const coreCluster = createReferenceGalaxyLayer(isMobileViewport ? 1400 : 3200, particleTexture, 'core', isMobileViewport);
    coreCluster.rotation.z = -0.18;
    galaxyRoot.add(coreCluster);

    const upperPulse = createReferenceGalaxyLayer(isMobileViewport ? 820 : 1800, particleTexture, 'upperPulse', isMobileViewport);
    upperPulse.rotation.z = -0.18;
    galaxyRoot.add(upperPulse);

    const lowerSparkles = createReferenceGalaxyLayer(isMobileViewport ? 420 : 900, particleTexture, 'lowerSparkles', isMobileViewport);
    lowerSparkles.rotation.z = -0.18;
    galaxyRoot.add(lowerSparkles);

    const coreTexture = createGalaxyCoreTexture();
    const core = new THREE.Sprite(
      new THREE.SpriteMaterial({
        map: coreTexture,
        transparent: true,
        opacity: 0.98,
        blending: THREE.AdditiveBlending,
        depthWrite: false,
      }),
    );
    core.position.set(0, 0.02, -0.9);
    core.scale.set(2.62, 1.18, 1);
    galaxyRoot.add(core);

    const routePositions: number[] = [];
    routeLinks.forEach(([from, to]) => {
      const start = getCourseStarPosition(from);
      const end = getCourseStarPosition(to);
      routePositions.push(start.x, start.y, -1.08, end.x, end.y, -1.08);
    });
    const routeGeometry = new THREE.BufferGeometry();
    routeGeometry.setAttribute('position', new THREE.Float32BufferAttribute(routePositions, 3));
    const routeMaterial = new THREE.LineBasicMaterial({
      color: '#8B7CFC',
      transparent: true,
      opacity: 0.18,
      blending: THREE.AdditiveBlending,
    });
    const routes = new THREE.LineSegments(routeGeometry, routeMaterial);
    galaxyRoot.add(routes);

    sampleStars.forEach((star, index) => {
      const position = getCourseStarPosition(index);
      const isActive = star.status === 'active';
      const haloTexture = createCourseHaloTexture(star.tone);
      const ringTexture = createCourseRingTexture(star.tone);
      effectTextures.push(haloTexture, ringTexture);

      const halo = new THREE.Sprite(
        new THREE.SpriteMaterial({
          map: haloTexture,
          transparent: true,
          opacity: 0.42,
          blending: THREE.AdditiveBlending,
          depthWrite: false,
        }),
      );
      const haloScale = 0.28 + star.size / 88;
      halo.scale.set(haloScale, haloScale, 1);
      halo.position.set(position.x, position.y, position.z - 0.03);
      halo.userData = { label: star.label, baseScale: haloScale };
      galaxyRoot.add(halo);
      courseHalos.push(halo);

      const ring = new THREE.Sprite(
        new THREE.SpriteMaterial({
          map: ringTexture,
          transparent: true,
          opacity: 0.52,
          blending: THREE.AdditiveBlending,
          depthWrite: false,
        }),
      );
      const ringScale = 0.24 + star.size / 120;
      ring.scale.set(ringScale, ringScale, 1);
      ring.position.set(position.x, position.y, position.z + 0.02);
      ring.userData = { label: star.label, baseScale: ringScale };
      galaxyRoot.add(ring);
      courseRings.push(ring);

      const texture = createStarTexture(star.tone, isActive);
      starTextures.push(texture);
      const material = new THREE.SpriteMaterial({
        map: texture,
        transparent: true,
        opacity: star.status === 'done' ? 0.92 : 1,
        blending: THREE.AdditiveBlending,
        depthWrite: false,
      });
      const sprite = new THREE.Sprite(material);
      const scale = 0.14 + star.size / 128;
      sprite.scale.set(scale, scale, 1);
      sprite.position.copy(position);
      sprite.userData = { label: star.label, baseScale: scale, tone: star.tone };
      galaxyRoot.add(sprite);
      starSprites.push(sprite);

      const labelTexture = createLabelTexture(star.label);
      labelTextures.push(labelTexture);
      const labelMaterial = new THREE.SpriteMaterial({ map: labelTexture, transparent: true, opacity: isActive ? 0.38 : 0.18 });
      const labelSprite = new THREE.Sprite(labelMaterial);
      labelSprite.scale.set(0.76, 0.21, 1);
      labelSprite.position.set(position.x, position.y - 0.28, -0.82);
      galaxyRoot.add(labelSprite);
      labelSprites.push(labelSprite);
    });

    const resize = () => {
      const rect = mount.getBoundingClientRect();
      renderer.setSize(Math.max(1, rect.width), Math.max(1, rect.height), false);
      camera.aspect = Math.max(1, rect.width) / Math.max(1, rect.height);
      camera.updateProjectionMatrix();
    };

    const showOverview = () => {
      cameraTarget.set(0, 0.08, 6.85);
    };

    const focusStar = (label: string | null) => {
      if (!label) {
        showOverview();
        return;
      }
      const star = sampleStars.find((item) => item.label === label);
      if (!star) {
        showOverview();
        return;
      }
      const position = getCourseStarPosition(sampleStars.indexOf(star)).clone();
      position.applyEuler(galaxyRoot.rotation);
      cameraTarget.set(position.x * 0.5, (position.y + galaxyRoot.position.y) * 0.5 + 0.05, 3.45);
    };

    const handlePointerDown = (event: PointerEvent) => {
      dragState.active = true;
      dragState.moved = false;
      dragState.pointerId = event.pointerId;
      dragState.startX = event.clientX;
      dragState.startY = event.clientY;
      dragState.lastX = event.clientX;
      dragState.lastY = event.clientY;
      renderer.domElement.style.cursor = 'grabbing';
      renderer.domElement.setPointerCapture?.(event.pointerId);
    };

    const handlePointerMove = (event: PointerEvent) => {
      if (!dragState.active || event.pointerId !== dragState.pointerId) return;
      const dx = event.clientX - dragState.lastX;
      const dy = event.clientY - dragState.lastY;
      const totalDx = event.clientX - dragState.startX;
      const totalDy = event.clientY - dragState.startY;
      if (Math.hypot(totalDx, totalDy) > 5) dragState.moved = true;
      dragState.targetY += dx * 0.006;
      dragState.targetX += dy * 0.004;
      dragState.targetX = Math.max(-1.18, Math.min(-0.42, dragState.targetX));
      dragState.lastX = event.clientX;
      dragState.lastY = event.clientY;
    };

    const handlePointerUp = (event: PointerEvent) => {
      if (!dragState.active || event.pointerId !== dragState.pointerId) return;
      dragState.active = false;
      dragState.pointerId = -1;
      renderer.domElement.style.cursor = 'grab';
      renderer.domElement.releasePointerCapture?.(event.pointerId);

      if (dragState.moved) return;

      const rect = renderer.domElement.getBoundingClientRect();
      pointer.x = ((event.clientX - rect.left) / rect.width) * 2 - 1;
      pointer.y = -((event.clientY - rect.top) / rect.height) * 2 + 1;
      raycaster.setFromCamera(pointer, camera);
      const intersections = raycaster.intersectObjects(starSprites, false);
      if (!intersections.length) {
        onSelectRef.current(null);
        showOverview();
        return;
      }
      const selected = intersections[0].object.userData.label as string;
      onSelectRef.current(selected);
      focusStar(selected);
    };

    const animate = () => {
      const elapsed = clock.getElapsedTime();
      dragState.rotationX += (dragState.targetX - dragState.rotationX) * 0.08;
      dragState.rotationY += (dragState.targetY - dragState.rotationY) * 0.08;
      galaxyRoot.rotation.x = dragState.rotationX + Math.sin(elapsed * 0.16) * 0.012;
      galaxyRoot.rotation.y = dragState.rotationY + Math.sin(elapsed * 0.12) * 0.018;
      galaxyRoot.rotation.z = elapsed * 0.006;

      backgroundStars.rotation.z = 0;
      galaxyDisk.rotation.z = -0.18;
      innerDust.rotation.z = -0.18;
      blueRim.rotation.z = -0.18;
      outerArm.rotation.z = -0.18;
      coreCluster.rotation.z = -0.18;
      upperPulse.rotation.z = -0.18;
      (upperPulse.material as THREE.PointsMaterial).opacity = 0.44;
      lowerSparkles.rotation.z = -0.18;
      core.material.rotation = Math.sin(elapsed * 0.18) * 0.035;
      routes.rotation.z = 0;

      if (selectedLabelRef.current != null) {
        const selectedStar = sampleStars.find((item) => item.label === selectedLabelRef.current);
        if (selectedStar) {
          const selectedPosition = getCourseStarPosition(sampleStars.indexOf(selectedStar)).clone();
          selectedPosition.applyEuler(galaxyRoot.rotation);
          cameraTarget.set(selectedPosition.x * 0.5, (selectedPosition.y + galaxyRoot.position.y) * 0.5 + 0.05, 3.45);
        }
      }

      starSprites.forEach((sprite, index) => {
        const base = sprite.userData.baseScale as number;
        const hasSelection = selectedLabelRef.current != null;
        const isSelected = sprite.userData.label === selectedLabelRef.current;
        const pulse = 1 + Math.sin(elapsed * 2.6 + index) * (isSelected ? 0.09 : 0.045);
        const target = base * pulse * (isSelected ? 1.2 : 1);
        sprite.scale.lerp(new THREE.Vector3(target, target, 1), 0.16);
        const material = sprite.material as THREE.SpriteMaterial;
        material.opacity = hasSelection ? (isSelected ? 1 : 0.72) : 0.92;
      });

      courseHalos.forEach((sprite, index) => {
        const base = sprite.userData.baseScale as number;
        const isSelected = sprite.userData.label === selectedLabelRef.current;
        const pulse = 1 + Math.sin(elapsed * 1.8 + index * 0.7) * 0.08;
        const target = base * pulse * (isSelected ? 1.18 : 1);
        sprite.scale.lerp(new THREE.Vector3(target, target, 1), 0.12);
        const material = sprite.material as THREE.SpriteMaterial;
        material.opacity += ((isSelected ? 0.72 : 0.36) - material.opacity) * 0.12;
      });

      courseRings.forEach((sprite, index) => {
        const base = sprite.userData.baseScale as number;
        const isSelected = sprite.userData.label === selectedLabelRef.current;
        const pulse = 1 + Math.sin(elapsed * 2.2 + index) * 0.045;
        const target = base * pulse * (isSelected ? 1.24 : 1);
        sprite.scale.lerp(new THREE.Vector3(target, target, 1), 0.14);
        const material = sprite.material as THREE.SpriteMaterial;
        material.rotation = elapsed * (isSelected ? 0.42 : 0.18) + index * 0.4;
        material.opacity += ((isSelected ? 0.86 : 0.44) - material.opacity) * 0.12;
      });

      labelSprites.forEach((sprite, index) => {
        const label = sampleStars[index].label;
        const material = sprite.material as THREE.SpriteMaterial;
        const targetOpacity = selectedLabelRef.current == null ? 0.035 : label === selectedLabelRef.current ? 0.62 : 0.11;
        material.opacity += (targetOpacity - material.opacity) * 0.12;
      });

      camera.position.lerp(cameraTarget, reduceMotion ? 1 : 0.075);
      camera.lookAt(camera.position.x * 0.36, camera.position.y * 0.36, 0);
      renderer.render(scene, camera);
      if (!reduceMotion) frame = window.requestAnimationFrame(animate);
    };

    let frame = 0;
    resize();
    showOverview();
    renderer.domElement.addEventListener('pointerdown', handlePointerDown);
    renderer.domElement.addEventListener('pointermove', handlePointerMove);
    renderer.domElement.addEventListener('pointerup', handlePointerUp);
    renderer.domElement.addEventListener('pointercancel', handlePointerUp);
    window.addEventListener('resize', resize);
    animate();

    return () => {
      if (frame) window.cancelAnimationFrame(frame);
      window.removeEventListener('resize', resize);
      renderer.domElement.removeEventListener('pointerdown', handlePointerDown);
      renderer.domElement.removeEventListener('pointermove', handlePointerMove);
      renderer.domElement.removeEventListener('pointerup', handlePointerUp);
      renderer.domElement.removeEventListener('pointercancel', handlePointerUp);
      starTextures.forEach((texture) => texture.dispose());
      labelTextures.forEach((texture) => texture.dispose());
      effectTextures.forEach((texture) => texture.dispose());
      starSprites.forEach((sprite) => (sprite.material as THREE.SpriteMaterial).dispose());
      courseHalos.forEach((sprite) => (sprite.material as THREE.SpriteMaterial).dispose());
      courseRings.forEach((sprite) => (sprite.material as THREE.SpriteMaterial).dispose());
      labelSprites.forEach((sprite) => (sprite.material as THREE.SpriteMaterial).dispose());
      routeGeometry.dispose();
      routeMaterial.dispose();
      backgroundStars.geometry.dispose();
      (backgroundStars.material as THREE.Material).dispose();
      galaxyDisk.geometry.dispose();
      (galaxyDisk.material as THREE.Material).dispose();
      innerDust.geometry.dispose();
      (innerDust.material as THREE.Material).dispose();
      blueRim.geometry.dispose();
      (blueRim.material as THREE.Material).dispose();
      outerArm.geometry.dispose();
      (outerArm.material as THREE.Material).dispose();
      coreCluster.geometry.dispose();
      (coreCluster.material as THREE.Material).dispose();
      upperPulse.geometry.dispose();
      (upperPulse.material as THREE.Material).dispose();
      lowerSparkles.geometry.dispose();
      (lowerSparkles.material as THREE.Material).dispose();
      particleTexture.dispose();
      coreTexture.dispose();
      (core.material as THREE.SpriteMaterial).dispose();
      renderer.dispose();
      renderer.domElement.remove();
    };
  }, []);

  return <div ref={mountRef} style={threeLayerStyle} aria-hidden="true" />;
}

export default function PublicCosmosMobileFirstView({ copy, locale, loginHref }: Props) {
  const router = useRouter();
  const [query, setQuery] = useState('');
  const [selectedLabel, setSelectedLabel] = useState<string | null>(null);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [isCompactCta, setIsCompactCta] = useState(false);
  const { isGoalSubmitting, goalFlowStage, goalCreationError, handleGoalSubmit } = useGoalCreation({});
  const chips = useMemo(() => copy.chips.slice(0, 4), [copy.chips]);

  useEffect(() => {
    setIsLoggedIn(document.cookie.includes('is_logged_in=1'));
  }, []);

  useEffect(() => {
    const updateCompactCta = () => setIsCompactCta(window.innerWidth <= 344);
    updateCompactCta();
    window.addEventListener('resize', updateCompactCta);
    return () => window.removeEventListener('resize', updateCompactCta);
  }, []);

  const selectStar = (label: string | null) => {
    setSelectedLabel(label);
    if (label) setQuery(label);
  };

  const useTopicChip = (label: string) => {
    setSelectedLabel(null);
    setQuery(label);
  };

  const getGoalPath = () => {
    const trimmed = query.trim();
    return trimmed ? `/dashboard/goal?intent=${encodeURIComponent(trimmed)}` : '/dashboard/goal';
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!query.trim()) {
      await handleGoalSubmit(query);
      return;
    }
    if (!isLoggedIn) {
      router.push(`${loginHref}?redirect_after=${encodeURIComponent(getGoalPath())}`);
      return;
    }
    await handleGoalSubmit(query);
  };

  return (
    <section style={rootStyle} aria-label="LearnCosmos mobile first view">
      {isGoalSubmitting && goalFlowStage ? <GoalFlowLoadingOverlay stage={goalFlowStage} /> : null}
      <div style={shellStyle}>
        <div style={lumiCardStyle}>
          <p style={eyebrowStyle}>Lumi Route</p>
          <p style={lumiTextStyle}>
            {locale === 'en'
              ? 'Tell Lumi what you want to learn. We will shape it into your first learning route.'
              : '배우고 싶은 주제를 편하게 적어 주세요. Lumi가 첫 번째 학습탐험 경로를 함께 열어드릴게요.'}
          </p>
        </div>

        <div style={mapFrameStyle}>
          <div style={mapStyle} aria-label="밝은 코스모스 학습 지도">
            <ThreeCosmosMap selectedLabel={selectedLabel} onSelect={selectStar} />
            <div aria-hidden="true" style={dustStyle} />
            <div style={{ ...ctaOverlayStyle, ...(isCompactCta ? compactCtaOverlayStyle : null) }}>
              <form
                style={{ ...formStyle, ...(isCompactCta ? compactFormStyle : null) }}
                onSubmit={handleSubmit}
                aria-label={copy.inputLabel}
              >
                <input
                  value={query}
                  onChange={(event) => setQuery(event.target.value)}
                  style={{ ...inputStyle, ...(isCompactCta ? compactInputStyle : null) }}
                  placeholder={copy.placeholder}
                  aria-label={copy.inputLabel}
                />
                <button
                  type="submit"
                  style={{ ...buttonStyle, ...(isCompactCta ? compactButtonStyle : null) }}
                  disabled={isGoalSubmitting}
                >
                  {isGoalSubmitting ? (locale === 'en' ? 'Creating' : '생성중') : locale === 'en' ? 'New star' : '새 별'}
                </button>
                {goalCreationError ? (
                  <p style={{ gridColumn: '1 / -1', margin: 0, color: '#FFD4D4', fontSize: '11px', fontWeight: 850 }}>
                    {goalCreationError}
                  </p>
                ) : null}
              </form>

              <div style={chipRowStyle} aria-label="예시 주제">
                {chips.map((chip) => (
                  <button
                    key={chip}
                    type="button"
                    style={{ ...chipStyle, ...(isCompactCta ? compactChipStyle : null) }}
                    onClick={() => useTopicChip(chip)}
                  >
                    {chip}
                  </button>
                ))}
              </div>
            </div>
          </div>
        </div>

        <section style={videoSectionStyle} aria-label="LearnCosmos 안내영상">
          <div style={videoHeaderStyle}>
            <p style={videoEyebrowStyle}>Guide Video</p>
            <h2 style={videoTitleStyle}>{locale === 'en' ? 'Start with a short tour.' : '짧은 안내영상으로 먼저 둘러보세요.'}</h2>
            <p style={videoTextStyle}>
              {locale === 'en'
                ? 'See the flow from topic input to goal chat, learning route, Explorer Diary, and completion.'
                : '주제 입력에서 목표 채팅, 학습탐험 경로, Explorer Diary, 완료까지 이어지는 흐름을 보여줍니다.'}
            </p>
          </div>
          <a
            href={guideVideoUrl || undefined}
            style={{ ...videoPreviewStyle, cursor: guideVideoUrl ? 'pointer' : 'default' }}
            aria-disabled={!guideVideoUrl}
            target={guideVideoUrl ? '_blank' : undefined}
            rel={guideVideoUrl ? 'noreferrer' : undefined}
          >
            <span style={videoPlayStyle} aria-hidden="true">▶</span>
            <span style={videoBadgeStyle}>
              <span>{locale === 'en' ? 'LearnCosmos YouTube guide' : 'LearnCosmos YouTube 안내영상'}</span>
              <span>{locale === 'en' ? 'Watch' : '보기'}</span>
            </span>
          </a>
        </section>

        <section style={betaSectionStyle} aria-label="LearnCosmos 베타 테스트 참가 안내">
          <p style={betaLabelStyle}>{locale === 'en' ? 'Beta Test' : '초기 베타 테스트'}</p>
          <h2 style={betaTitleStyle}>{locale === 'en' ? 'Shape the first learning galaxy.' : '첫 번째 학습 은하를 함께 다듬어 주세요.'}</h2>
          <p style={betaTextStyle}>
            {locale === 'en'
              ? 'LearnCosmos is still stabilizing its early learning experience. Some screens, recommendations, and AI behavior may change as feedback comes in.'
              : 'LearnCosmos는 아직 초기 베타 단계입니다. 일부 화면, 추천 결과, AI 동작은 피드백을 반영하며 계속 조정될 수 있습니다.'}
          </p>
          <div style={betaListStyle}>
            {betaExperienceItems[locale === 'en' ? 'en' : 'ko'].map((item) => (
              <div key={item} style={betaItemStyle}>
                <span style={betaCheckStyle} aria-hidden="true">✓</span>
                <span>{item}</span>
              </div>
            ))}
          </div>
          <p style={betaTextStyle}>
            {locale === 'en'
              ? 'After trying it, tell us what felt useful, confusing, or missing through the survey.'
              : '직접 둘러본 뒤 좋았던 점, 헷갈린 점, 꼭 필요하다고 느낀 기능을 설문으로 알려주세요.'}
          </p>
          <div style={betaActionsStyle}>
            <a
              href={betaSurveyUrl || undefined}
              style={{ ...betaLinkStyle, opacity: betaSurveyUrl ? 1 : 0.62, cursor: betaSurveyUrl ? 'pointer' : 'default' }}
              aria-disabled={!betaSurveyUrl}
              target={betaSurveyUrl ? '_blank' : undefined}
              rel={betaSurveyUrl ? 'noreferrer' : undefined}
            >
              {locale === 'en' ? 'Open survey' : '설문조사 참여하기'}
            </a>
            <p style={betaNoteStyle}>
              {locale === 'en'
                ? betaSurveyUrl
                  ? 'The survey opens in a new tab.'
                  : 'Survey link will be connected soon.'
                : betaSurveyUrl
                  ? '설문은 새 탭에서 열립니다.'
                  : '설문 링크는 곧 연결됩니다.'}
            </p>
          </div>
        </section>

        <footer style={miniFooterStyle} aria-label="LearnCosmos footer">
          <p style={miniFooterBrandStyle}>LearnCosmos</p>
          <p style={miniFooterTaglineStyle}>
            {locale === 'en' ? 'From topic input to records and completion.' : '목표 정리부터 기록과 완료까지 함께합니다.'}
          </p>
          <nav style={miniFooterLinksStyle} aria-label={locale === 'en' ? 'Footer links' : '푸터 링크'}>
            <a href={locale === 'en' ? '/en/terms' : '/terms'} style={miniFooterLinkStyle}>
              {locale === 'en' ? 'Terms' : '이용약관'}
            </a>
            <a href={locale === 'en' ? '/en/privacy' : '/privacy'} style={miniFooterLinkStyle}>
              {locale === 'en' ? 'Privacy' : '개인정보처리방침'}
            </a>
            <a href={locale === 'en' ? '/en/open-source' : '/open-source'} style={miniFooterLinkStyle}>
              {locale === 'en' ? 'Open Source' : '오픈소스'}
            </a>
            <a href={locale === 'en' ? '/en/youtube-api-use' : '/youtube-api-use'} style={miniFooterLinkStyle}>
              {locale === 'en' ? 'YouTube API' : 'YouTube API 사용'}
            </a>
            <span style={miniFooterContactStyle}>{locale === 'en' ? 'Contact: learnweavr@gmail.com' : '문의: learnweavr@gmail.com'}</span>
          </nav>
          <p style={miniFooterMetaStyle}>© 2026 LearnCosmos. AI based learning tool platform.</p>
        </footer>
      </div>
    </section>
  );
}
