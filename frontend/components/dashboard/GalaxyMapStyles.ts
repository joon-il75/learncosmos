import type { CSSProperties } from 'react';
import {
  canonicalStatusColor,
  normalizePlanetStatus,
  type PlanetStatus,
} from '@/components/dashboard/planetMapTokens';
import {
  getPlanetSpriteBackgroundStyle,
  getPlanetTextureFamily,
  getTerraformStage,
} from '@/components/dashboard/planetSpriteSystem';
import {
  getPlanetProgressVisualProfile,
  normalizePlanetProgressPercent,
} from '@/lib/world-ui-engine/planetProgressUtils';
import type { DashboardCourseRecord } from '@/lib/world-ui-engine/types';

export function getGalaxyCourseProgressPercent(course: DashboardCourseRecord): number {
  if (course.status === 'completed') return 100;
  if ((course.progress ?? null) != null) {
    return Math.min(99, normalizePlanetProgressPercent(course.progress, course.status));
  }
  const lessonCount = Math.max(0, course.lessonCount ?? 0);
  const completedLessonCount = Math.max(0, course.completedLessonCount ?? 0);
  if (lessonCount > 0) {
    return Math.min(99, normalizePlanetProgressPercent(completedLessonCount / lessonCount, course.status));
  }
  return 0;
}

export function mapSystemPlanetBodyStyle(
  course: DashboardCourseRecord,
  isHighlighted: boolean,
  objectScale: number,
): CSSProperties {
  const size = 14 * objectScale;
  const isInactive = Boolean(course.isInactive);
  const normalizedProgressPercent = getGalaxyCourseProgressPercent(course);
  const progressVisualProfile = getPlanetProgressVisualProfile(normalizedProgressPercent);
  const backgroundStyle = isInactive
    ? {
        backgroundImage: 'none',
        background: isInactive
          ? 'radial-gradient(circle at 34% 30%, rgba(226, 232, 240, 0.46), rgba(148, 163, 184, 0.55) 28%, rgba(71, 85, 105, 0.92) 70%, rgba(15, 23, 42, 0.98) 100%)'
          : 'radial-gradient(circle at 36% 32%, rgba(242, 247, 255, 0.22), rgba(88, 108, 146, 0.28) 26%, rgba(26, 35, 52, 0.9) 72%, rgba(10, 16, 28, 0.98) 100%)',
      }
    : getPlanetSpriteBackgroundStyle(
        getPlanetTextureFamily(course),
        Math.max(getTerraformStage(course), progressVisualProfile.stage + (course.status === 'completed' ? 3 : 0)),
        size,
      );

  return {
    position: 'absolute',
    inset: `${2 * objectScale}px`,
    borderRadius: '999px',
    overflow: 'hidden',
    ...backgroundStyle,
    backgroundColor: 'rgba(8, 14, 24, 0.18)',
    boxShadow: isHighlighted
      ? 'inset -2px -3px 5px rgba(8, 16, 28, 0.28), 0 0 5px rgba(255,255,255,0.22)'
      : 'inset -2px -3px 4px rgba(8, 16, 28, 0.24)',
    filter: isInactive
      ? (isHighlighted ? 'brightness(0.92) saturate(0.28) contrast(1.04)' : 'brightness(0.82) saturate(0.24)')
      : isHighlighted ? 'brightness(1.1) saturate(1.08) contrast(1.04)' : 'brightness(1.04) saturate(1.04)',
    zIndex: 2,
  };
}

export const mapSystemPlanetCanvasStyle: CSSProperties = {
  width: '100%',
  height: '100%',
  borderRadius: '999px',
  pointerEvents: 'none',
};

export const mapSystemPlanetPreviewStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  width: '100%',
  height: '100%',
};

// ─── Window / Header ─────────────────────────────────────────────────────────

export const mapWindowStyle = {
  display: 'grid',
  gap: '14px',
  borderRadius: '28px',
  border: '1px solid rgba(214, 226, 255, 0.24)',
  background: 'linear-gradient(180deg, rgba(7, 15, 30, 0.86), rgba(5, 11, 24, 0.8))',
  boxShadow: '0 24px 80px rgba(0, 0, 0, 0.3)',
  backdropFilter: 'blur(18px)',
  padding: '18px',
} satisfies CSSProperties;

export const mapWindowHeaderStyle = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '12px',
  flexWrap: 'wrap',
} satisfies CSSProperties;

export const mapWindowHeaderGroupStyle = {
  display: 'grid',
  gap: '4px',
} satisfies CSSProperties;

export const mapWindowEyebrowStyle = {
  fontSize: '12px',
  letterSpacing: '0.16em',
  textTransform: 'uppercase',
  color: '#8CC7FF',
} satisfies CSSProperties;

export const mapWindowMetaStyle = {
  fontSize: '13px',
  color: '#8EA4C8',
} satisfies CSSProperties;

export const sortControlRowStyle = {
  display: 'flex',
  flexWrap: 'wrap',
  gap: '8px',
} satisfies CSSProperties;

export function sortButtonStyle(isActive: boolean): CSSProperties {
  return {
    minHeight: '32px',
    padding: '0 14px',
    borderRadius: '999px',
    border: isActive ? '1px solid #FCD34D' : '1px solid rgba(255, 255, 255, 0.2)',
    background: isActive ? 'rgba(252, 211, 77, 0.16)' : 'rgba(255, 255, 255, 0.04)',
    color: isActive ? '#FDE68A' : '#E8EEFF',
    fontSize: '12px',
    fontWeight: 600,
    cursor: 'pointer',
    fontFamily: 'inherit',
  };
}

export const mapWindowSinglePaneStyle = {
  display: 'grid',
  gap: '14px',
} satisfies CSSProperties;

export const mapWindowFrameStyle = {
  position: 'relative',
  minHeight: '600px',
  overflow: 'hidden',
  borderRadius: '22px',
  border: '1px solid rgba(214, 226, 255, 0.18)',
  background: 'rgba(0, 0, 0, 0.2)',
} satisfies CSSProperties;

export const mapWindowImageStyle = {
  objectFit: 'cover',
  objectPosition: '54% -30%',
  opacity: 0.4,
} satisfies CSSProperties;

export const mapWindowImageStarSystemStyle = {
  objectFit: 'cover',
  objectPosition: 'center',
  opacity: 0.55,
} satisfies CSSProperties;

export const starSystemMapWrapperStyle = {
  position: 'relative',
  minHeight: '560px',
  width: '100%',
  display: 'grid',
  placeItems: 'center',
} satisfies CSSProperties;

export const starSystemMissingStateStyle = {
  position: 'relative',
  zIndex: 3,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '560px',
  padding: '24px',
  textAlign: 'center',
  color: '#E5ECFF',
  fontSize: '15px',
  lineHeight: 1.7,
} satisfies CSSProperties;

export const mapWindowGlowStyle = {
  position: 'absolute',
  inset: 0,
  background: 'linear-gradient(180deg, rgba(4, 9, 20, 0.06), rgba(4, 9, 20, 0.24))',
  pointerEvents: 'none',
} satisfies CSSProperties;

// ─── Mobile instruction ───────────────────────────────────────────────────────

export const mobileInstructionStyle = {
  display: 'grid',
  gap: '6px',
  padding: '12px',
  borderRadius: '18px',
  border: '1px solid rgba(255, 255, 255, 0.12)',
  background: 'rgba(10, 18, 35, 0.85)',
} satisfies CSSProperties;

export const mobileInstructionTitleStyle = {
  fontSize: '13px',
  fontWeight: 700,
  color: '#FDE68A',
} satisfies CSSProperties;

export const mobileInstructionTextStyle = {
  fontSize: '12px',
  color: '#CFE1FF',
  lineHeight: 1.5,
} satisfies CSSProperties;

export const mobileInstructionSubtitleStyle = {
  fontSize: '12px',
  color: '#95B2FF',
  lineHeight: 1.4,
} satisfies CSSProperties;

// ─── Galaxy node ──────────────────────────────────────────────────────────────

export function mapWindowPageNodeStyle(
  index: number,
  nodeCount: number,
  isHighlighted: boolean,
  useColumnLayout: boolean,
  objectScale: number,
): CSSProperties {
  let left: string;
  let top: string;

  if (useColumnLayout) {
    const spacingPct = 14;
    const totalSpan = (nodeCount - 1) * spacingPct;
    const yPct = nodeCount > 1 ? -(totalSpan / 2) + index * spacingPct : 0;
    left = '50%';
    top = `calc(50% + ${yPct}%)`;
  } else {
    const orbitSequence = ['inner', 'outer', 'middle'] as const;
    const orbit = orbitSequence[index % orbitSequence.length];
    const startAngle = -Math.PI / 3;
    const angle = ((Math.PI * 2) / Math.max(nodeCount, 1)) * index + startAngle;
    const radiusX = orbit === 'inner' ? 148 : orbit === 'middle' ? 238 : 295;
    const radiusY = orbit === 'inner' ? 66 : orbit === 'middle' ? 108 : 120;
    const x = -Math.cos(angle) * radiusX;
    const y = Math.sin(angle) * radiusY;
    left = `calc(50% + ${x * objectScale}px)`;
    top = `calc(50% + ${y * objectScale}px)`;
  }

  return {
    position: 'absolute',
    left,
    top,
    width: `${76 * objectScale}px`,
    height: `${76 * objectScale}px`,
    transform: 'translate(-50%, -50%)',
    animation: `${index % 3 === 0 ? 'systemFloatA' : index % 3 === 1 ? 'systemFloatB' : 'systemFloatC'} ${7 + (index % 3)}s ease-in-out infinite`,
    display: 'grid',
    placeItems: 'center',
    opacity: 1,
    filter: 'brightness(1)',
    zIndex: isHighlighted ? 3 : 1,
    cursor: 'pointer',
  };
}

export function mapSystemScaleStyle(isHighlighted: boolean): CSSProperties {
  return {
    display: 'grid',
    justifyItems: 'center',
    transform: `scale(${isHighlighted ? 1.16 : 1})`,
    filter: isHighlighted ? 'brightness(1.1)' : 'brightness(1)',
    transition: 'transform 180ms ease, filter 180ms ease',
  };
}

export function mapSystemCoreStyle(isHighlighted: boolean, objectScale: number): CSSProperties {
  return {
    position: 'relative',
    width: `${76 * objectScale}px`,
    height: `${76 * objectScale}px`,
    borderRadius: '999px',
    border: isHighlighted
      ? '1px solid rgba(196, 222, 255, 0.42)'
      : '1px solid rgba(168, 202, 255, 0.26)',
    background: 'transparent',
    boxShadow: isHighlighted
      ? '0 0 26px rgba(122, 170, 255, 0.28), inset 0 0 16px rgba(196, 222, 255, 0.12)'
      : '0 0 22px rgba(122, 170, 255, 0.2), inset 0 0 14px rgba(168, 202, 255, 0.08)',
    backdropFilter: 'none',
    transition: 'border-color 180ms ease, box-shadow 180ms ease',
  };
}

export function mapSystemSunStyle(isHighlighted: boolean, objectScale: number): CSSProperties {
  return {
    position: 'absolute',
    left: '50%',
    top: '50%',
    width: `${18 * objectScale}px`,
    height: `${18 * objectScale}px`,
    transform: 'translate(-50%, -50%)',
    borderRadius: '999px',
    background:
      'radial-gradient(circle, rgba(255,247,212,1) 0%, rgba(255,194,95,0.95) 45%, rgba(237,142,26,0.78) 100%)',
    boxShadow: isHighlighted
      ? '0 0 24px rgba(255, 194, 95, 0.44)'
      : '0 0 14px rgba(255, 194, 95, 0.18)',
    opacity: isHighlighted ? 1 : 0.8,
    animation: isHighlighted ? 'activeSunPulse 2.2s ease-in-out infinite' : 'none',
  };
}

export function mapSystemPlanetStyle(index: number, totalPlanets: number, objectScale: number): CSSProperties {
  const safeCount = Math.max(totalPlanets, 1);
  const startAngle = -(5 * Math.PI) / 6;
  const angle = ((Math.PI * 2) / safeCount) * index + startAngle;
  const x = Math.cos(angle) * 24 * objectScale;
  const y = Math.sin(angle) * 16 * objectScale;
  return {
    position: 'absolute',
    left: `calc(50% + ${x}px)`,
    top: `calc(50% + ${y}px)`,
    width: `${18 * objectScale}px`,
    height: `${18 * objectScale}px`,
    transform: 'translate(-50%, -50%)',
    pointerEvents: 'none',
  };
}

export function mapSystemPlanetRingStyle(
  status: PlanetStatus,
  isHighlighted: boolean,
  isInactive: boolean,
  objectScale: number,
): CSSProperties {
  const canonicalStatus = normalizePlanetStatus(status);
  const accent = isInactive ? '#94A3B8' : canonicalStatusColor(canonicalStatus);
  return {
    position: 'absolute',
    inset: `${-4 * objectScale}px`,
    borderRadius: '999px',
    border: `${Math.max(1.5, 2.2 * objectScale)}px solid ${accent}${isInactive ? 'a8' : 'cc'}`,
    boxShadow: isHighlighted
      ? `0 0 16px ${accent}8a, inset 0 0 6px ${accent}24`
      : `0 0 10px ${accent}48`,
    opacity: isHighlighted ? 1 : 0.98,
    zIndex: 1,
  };
}

export function mapSystemLabelStyle(isHighlighted: boolean, objectScale: number): CSSProperties {
  return {
    marginTop: `${6 * objectScale}px`,
    fontSize: isHighlighted ? '13px' : '12px',
    fontWeight: 700,
    color: isHighlighted ? '#FFFFFF' : '#F4F8FF',
    textAlign: 'center',
    textShadow: isHighlighted
      ? '0 0 14px rgba(156, 205, 255, 0.55)'
      : '0 1px 6px rgba(0, 0, 0, 0.45)',
    transition: 'font-size 180ms ease, color 180ms ease, text-shadow 180ms ease',
  };
}

// ─── Hover tooltip ────────────────────────────────────────────────────────────

export function mapHoverLumiTooltipStyle(x: number, y: number, frameRect: DOMRect | null): CSSProperties {
  const tooltipWidth = 280;
  const gap = 24;
  const padding = 8;

  const bounds =
    frameRect ??
    (typeof window !== 'undefined'
      ? { left: padding, right: window.innerWidth - padding, top: padding, bottom: window.innerHeight - padding }
      : { left: 0, right: 9999, top: 0, bottom: 9999 });

  let left: number;
  if (x + gap + tooltipWidth + padding <= bounds.right) {
    left = x + gap;
  } else {
    left = x - tooltipWidth - gap;
  }
  left = Math.max(bounds.left + padding, Math.min(left, bounds.right - tooltipWidth - padding));

  const top = bounds.top + padding;

  return {
    position: 'fixed',
    left: `${left}px`,
    top: `${top}px`,
    zIndex: 9999,
    pointerEvents: 'none',
    animation: 'planetTooltipFade 120ms ease both',
  };
}

export const mapHoverLumiStackStyle = {
  display: 'grid',
  gap: '8px',
  width: '280px',
} satisfies CSSProperties;

export const mapHoverLumiCardStyle = {
  display: 'grid',
  gap: '6px',
  padding: '10px',
  borderRadius: '14px',
  border: '1px solid rgba(148, 195, 255, 0.26)',
  background: 'rgba(6, 14, 28, 0.94)',
  backdropFilter: 'blur(16px)',
  boxShadow: '0 8px 32px rgba(0, 0, 0, 0.4)',
} satisfies CSSProperties;

export const mapHoverLumiHeaderStyle = {
  display: 'flex',
  alignItems: 'center',
  gap: '10px',
} satisfies CSSProperties;

export const mapHoverLumiTitleStyle = {
  fontSize: '13px',
  fontWeight: 700,
  color: '#E8F2FF',
} satisfies CSSProperties;

export const mapHoverLumiMetaStyle = {
  fontSize: '12px',
  color: '#8EA4C8',
} satisfies CSSProperties;

export const mapHoverLumiBubbleStyle = {
  fontSize: '11px',
  color: '#C8DCFF',
  lineHeight: 1.5,
  padding: '6px 8px',
  borderRadius: '10px',
  background: 'rgba(30, 60, 100, 0.38)',
} satisfies CSSProperties;

export const mapHoverLumiSummaryRowStyle = {
  display: 'flex',
  flexWrap: 'wrap',
  gap: '6px',
  fontSize: '11px',
  color: '#7A96C0',
  lineHeight: 1.3,
} satisfies CSSProperties;

export const mapHoverLumiCoursePanelStyle = {
  display: 'grid',
  gap: '6px',
  padding: '10px 12px',
  borderRadius: '14px',
  border: '1px solid rgba(120, 160, 255, 0.18)',
  background: 'rgba(8, 18, 38, 0.9)',
  backdropFilter: 'blur(10px)',
} satisfies CSSProperties;

export const mapHoverLumiCoursePanelHeaderStyle = {
  fontSize: '11px',
  fontWeight: 700,
  color: '#8CC7FF',
  letterSpacing: '0.12em',
  textTransform: 'uppercase',
} satisfies CSSProperties;

export const mapHoverLumiCourseListStyle = {
  display: 'grid',
  gap: '4px',
  maxHeight: '260px',
  overflowY: 'auto',
} satisfies CSSProperties;

export function mapHoverLumiCourseItemStyle(): CSSProperties {
  return {
    display: 'grid',
    gap: '2px',
    padding: '6px 8px',
    borderRadius: '10px',
    background: 'rgba(255, 255, 255, 0.04)',
  };
}

export function mapHoverLumiCourseTitleStyle(): CSSProperties {
  return {
    fontSize: '12px',
    fontWeight: 600,
    color: '#E4EEFF',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  };
}

export function mapHoverLumiCourseMetaStyle(): CSSProperties {
  return {
    fontSize: '11px',
    color: '#7A96C0',
  };
}
