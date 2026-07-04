'use client';

import { createPortal } from 'react-dom';
import { useState } from 'react';
import type { CSSProperties, ComponentProps } from 'react';
import LumiAvatar from '@/components/lumi/LumiAvatar';
import AtmosphericPlanetPreview from '@/components/dashboard/AtmosphericPlanetPreview';
import type { DashboardMainCopy } from '@/lib/i18n/pages/dashboardMain';
import type { Locale } from '@/lib/i18n/locales';
import { fallbackTextureMap } from '@/components/dashboard/useActivePlanetTextureMap';
import { getPlanetLayout, planetSlotPixelSize } from './PlanetMapLayout';
import {
  getPlanetSpriteBackgroundStyle,
  getPlanetTextureFamily,
  getTerraformStage,
} from './planetSpriteSystem';
import {
  canonicalStatusColor,
  canonicalStatusRing,
  dashboardTokens,
  normalizePlanetStatus,
  type PlanetStatus,
} from './planetMapTokens';
import {
  getPlanetProgressVisualProfile,
  getProgressGridPosition,
  normalizePlanetProgressPercent,
} from '@/lib/world-ui-engine/planetProgressUtils';

interface PlanetCourseItem {
  id: string;
  title: string;
  status: PlanetStatus;
  destination: string;
  lessonCount?: number;
  completedLessonCount?: number;
  updatedAt?: string;
  progress?: number | null;
  planetTextureMapAsset?: string | null;
  planetTextureMapRotationDurationSeconds?: number | null;
  planetTextureMapRotationDirection?: 'left' | 'right' | null;
  isInactive?: boolean;
}

interface PlanetMapProps {
  courses: PlanetCourseItem[];
  hoveredCourseId: string | null;
  selectedCourseId: string | null;
  isCompactLayout: boolean;
  isPhoneLayout: boolean;
  isShortViewport?: boolean;
  objectScale?: number;
  mode: 'galaxy' | 'star-system';
  selectedSystemTitle?: string | null;
  statusLine?: string | null;
  updatedAtLabel?: string | null;
  copy: DashboardMainCopy['planetMap'];
  statusLabels: DashboardMainCopy['galaxy']['statusLabels'];
  locale: Locale;
  onHoverCourse: (courseId: string | null) => void;
  onSelectCourse: (course: PlanetCourseItem) => void;
  onReturnToGalaxy?: () => void;
}

// Radii are in SVG viewBox units (0–100), orbit container = 500px desktop / 420px mobile / 320px short
// The outermost ring (rx=36, ry=27) corresponds to the planet orbit (rx=32%, ry=27% of safe area)
const orbitEllipsePresets = [
  { rx: 14, ry: 10, opacity: 0.13 },
  { rx: 24, ry: 18, opacity: 0.10 },
  { rx: 36, ry: 27, opacity: 0.18 },
] as const;
const starSystemPlanetSlotScale = 0.7;

export default function PlanetMap({
  courses,
  hoveredCourseId,
  selectedCourseId,
  isCompactLayout,
  isPhoneLayout,
  isShortViewport = false,
  objectScale = 1,
  mode,
  selectedSystemTitle,
  statusLine,
  updatedAtLabel,
  copy,
  statusLabels,
  locale,
  onHoverCourse,
  onSelectCourse,
  onReturnToGalaxy,
}: PlanetMapProps) {
  const breakpoint = 'desktop';
  const showStarSystemBanner = mode === 'star-system' && breakpoint === 'desktop';
  const [planetHoverPoint, setPlanetHoverPoint] = useState<{ x: number; y: number } | null>(null);
  const baseLayout = getPlanetLayout({
    planetIds: courses.map((course) => course.id),
    breakpoint,
  });
  const layout = baseLayout.map((item) => {
    return {
      ...item,
      position: {
        x: item.position.x,
        y:
          mode === 'star-system' && breakpoint === 'desktop'
            ? Math.min(item.position.y + 4, 84)
            : item.position.y,
      },
    };
  });
  const useCompactLabelMode = isPhoneLayout || isShortViewport;
  const tooltipCourse = hoveredCourseId
    ? courses.find((course) => course.id === hoveredCourseId) ?? null
    : null;
  const tooltipLayout = tooltipCourse
    ? layout.find((item) => item.id === tooltipCourse.id)
    : null;
  const shouldShowPlanetLumiTooltip = mode === 'star-system' && !isPhoneLayout && tooltipCourse && tooltipLayout && planetHoverPoint;

  return (
    <div style={planetCanvasStyle(useCompactLabelMode, isShortViewport)}>
      <div style={backgroundOverlayStyle}>
        <div style={{ ...cornerGlowStyle, left: -44, top: -44 }} />
        <div style={{ ...cornerGlowStyle, right: -44, top: -44 }} />
        <div style={{ ...cornerGlowStyle, left: -44, bottom: -44 }} />
        <div style={{ ...cornerGlowStyle, right: -44, bottom: -44 }} />
      </div>

      <div style={orbitContainerStyle(useCompactLabelMode, isShortViewport, objectScale)}>
        <svg aria-hidden="true" viewBox="0 0 100 100" style={orbitLayerStyle}>
          {orbitEllipsePresets.map((orbit, index) => {
            const isMainOrbit = index === orbitEllipsePresets.length - 1;
            return (
              <ellipse
                key={`${orbit.rx}-${orbit.ry}`}
                cx="50"
                cy="50"
                rx={orbit.rx}
                ry={orbit.ry}
                fill="none"
                stroke={`rgba(255,255,255,${orbit.opacity})`}
                strokeWidth={isMainOrbit ? 1.6 : 1.0}
              />
            );
          })}
        </svg>
      </div>

      <div style={centerHazeStyle(objectScale)} />
      <div style={sunNodeStyle(objectScale)}>
        <div style={sunAuraStyle(objectScale)} />
        <div style={sunGlowStyle(objectScale)} />
        <div style={sunCoreStyle} />
      </div>

      {showStarSystemBanner && (
        <div style={starSystemBannerStyle}>
          <div style={starSystemBannerTextStyle}>
            {copy.systemTitle(selectedSystemTitle)}
          </div>
          {(statusLine || updatedAtLabel) ? (
            <div style={starSystemBannerSummaryWrapper}>
              {statusLine ? <div style={starSystemBannerSummaryStyle}>{statusLine}</div> : null}
              {updatedAtLabel ? <div style={starSystemBannerUpdatedStyle}>{updatedAtLabel}</div> : null}
            </div>
          ) : null}
          {onReturnToGalaxy ? (
            <button type="button" style={starSystemReturnButtonStyle} onClick={onReturnToGalaxy}>
              {copy.returnToGalaxy}
            </button>
          ) : null}
        </div>
      )}

      <div style={planetSafeAreaStyle(useCompactLabelMode, isShortViewport, mode === 'star-system')}>
        {courses.length === 0 ? (
          <div style={emptyMapStyle}>{copy.empty}</div>
        ) : null}

        {courses.map((course, index) => {
          const slot = layout.find((item) => item.id === course.id);
          if (!slot) return null;

          const size =
            planetSlotPixelSize(slot.size, isCompactLayout || isShortViewport) *
            objectScale *
            (mode === 'star-system' ? starSystemPlanetSlotScale : 1);
          const isHovered = hoveredCourseId === course.id;
          const isSelected = selectedCourseId === course.id;
          const isInactive = Boolean(course.isInactive);
          const canonicalStatus = normalizePlanetStatus(course.status);
          const textureFamily = getPlanetTextureFamily(course);
          const normalizedProgressPercent = getCourseProgressPercent(course);
          const progressVisualProfile = getPlanetProgressVisualProfile(normalizedProgressPercent);
          const terraformStage = Math.max(getTerraformStage(course), progressVisualProfile.stage + (course.status === 'completed' ? 3 : 0));
          const starSystemZBoost = mode === 'star-system' ? 10 : 0;
          const planetTextureMap = getCoursePlanetTextureMap(course);
          const hasRotatingTexture = Boolean(course.planetTextureMapAsset);

          return (
            <button
              key={course.id}
              type="button"
              aria-label={copy.ariaLabel(course.title, statusLabels[course.status])}
              aria-pressed={isSelected}
              onMouseEnter={(event) => {
                setPlanetHoverPoint({ x: event.clientX, y: event.clientY });
                onHoverCourse(course.id);
              }}
              onMouseMove={(event) => {
                if (hoveredCourseId === course.id) setPlanetHoverPoint({ x: event.clientX, y: event.clientY });
              }}
              onMouseLeave={() => {
                setPlanetHoverPoint(null);
                onHoverCourse(null);
              }}
              onFocus={(event) => {
                const rect = event.currentTarget.getBoundingClientRect();
                setPlanetHoverPoint({ x: rect.left + rect.width / 2, y: rect.top });
                onHoverCourse(course.id);
              }}
              onBlur={() => {
                setPlanetHoverPoint(null);
                onHoverCourse(null);
              }}
              onClick={() => onSelectCourse(course)}
              style={{
                ...planetButtonStyle,
                position: 'absolute',
                left: `${slot.position.x}%`,
                top: `${slot.position.y}%`,
                width: `${size + (useCompactLabelMode ? 96 : 76) * objectScale}px`,
                zIndex: isHovered ? 7 + starSystemZBoost : isSelected ? 6 + starSystemZBoost : 2,
                transform: isHovered
                  ? `translate(-50%, -50%) scale(${1.095 + progressVisualProfile.scaleBoost})`
                  : isSelected
                    ? `translate(-50%, -50%) scale(${1.03 + progressVisualProfile.scaleBoost})`
                    : 'translate(-50%, -50%)',
                transition: 'transform 0.2s ease, opacity 0.2s ease, filter 0.2s ease',
                animation: `${planetFloatAnimationName(index)} ${planetFloatDuration(index)} ease-in-out infinite`,
                animationDelay: `${index * 0.26}s`,
              }}
            >
              <div style={planetShellStyle}>
                <div style={planetNodeStyle(size, isHovered, isSelected, progressVisualProfile)}>
                  <div style={planetGlowStyle(canonicalStatus, isHovered, isSelected, progressVisualProfile, isInactive)} />
                  <div style={planetOuterSoftRingStyle(canonicalStatus, isSelected, progressVisualProfile, isInactive)} />
                  <div style={planetInnerRingStyle(canonicalStatus, isSelected, progressVisualProfile, isInactive)} />
                  <div style={planetBaseSphereStyle(canonicalStatus, size, isHovered, isSelected, progressVisualProfile, isInactive)}>
                    {!isInactive && hasRotatingTexture ? (
                      <AtmosphericPlanetPreview
                        atlasURL={planetTextureMap.atlasURL}
                        progressPercent={normalizedProgressPercent}
                        size={size}
                        durationSeconds={planetTextureMap.rotationDurationSeconds}
                        direction={planetTextureMap.rotationDirection}
                        playing
                        preset="star-system"
                        style={rotatingPlanetPreviewStyle}
                        canvasStyle={rotatingPlanetCanvasStyle}
                      />
                    ) : !isInactive ? (
                      <div
                        style={planetTextureOverlayStyle(
                          textureFamily,
                          terraformStage,
                          size,
                        )}
                      />
                    ) : null}
                    <div style={planetHighlightStyle} />
                    <div style={planetSurfaceShadeStyle} />
                    <div style={planetEdgeShadeStyle} />
                  </div>
                  {canonicalStatus === 'completed' ? <div style={completedAccentStyle} /> : null}
                </div>
                <div style={planetLabelWrapStyle}>
                  <div title={course.title} style={planetLabelStyle(useCompactLabelMode, isShortViewport)}>
                    {course.title}
                  </div>
                </div>
              </div>
            </button>
          );
        })}

        {shouldShowPlanetLumiTooltip ? (
          <PlanetLumiTooltip course={tooltipCourse} point={planetHoverPoint} copy={copy} statusLabels={statusLabels} locale={locale} />
        ) : tooltipCourse && tooltipLayout && !selectedCourseId ? (
          <div
            style={{
              ...tooltipStyle,
              left: `${tooltipLayout.position.x}%`,
              top: `${tooltipLayout.position.y}%`,
              transform: useCompactLabelMode ? 'translate(-50%, calc(-100% - 70px))' : 'translate(-50%, calc(-100% - 78px))',
            }}
          >
            <div style={tooltipTitleStyle}>{tooltipCourse.title}</div>
            <div style={tooltipMetaRowStyle}>
              <span style={tooltipStatusPillStyle(tooltipCourse.status)}>{statusLabels[tooltipCourse.status]}</span>
              <span style={tooltipMetaStyle}>{copy.lessonCount(tooltipCourse.lessonCount ?? 0)}</span>
            </div>
          </div>
        ) : null}
      </div>
    </div>
  );
}

function PlanetLumiTooltip({
  course,
  point,
  copy,
  statusLabels,
  locale,
}: {
  course: PlanetCourseItem;
  point: { x: number; y: number };
  copy: DashboardMainCopy['planetMap'];
  statusLabels: DashboardMainCopy['galaxy']['statusLabels'];
  locale: Locale;
}) {
  if (typeof document === 'undefined') return null;
  const progressPercent = getCourseProgressPercent(course);
  const updatedLabel = course.updatedAt
    ? new Date(course.updatedAt).toLocaleDateString(locale === 'en' ? 'en-US' : 'ko-KR', { month: 'short', day: 'numeric' })
    : '';
  const tooltip = (
    <div style={planetLumiTooltipStyle(point.x, point.y)}>
      <div style={planetLumiCardStyle}>
        <div style={planetLumiHeaderStyle}>
          <LumiAvatar state={lumiStateForPlanet(course)} size={32} />
          <div style={{ display: 'grid', gap: '2px', minWidth: 0 }}>
            <strong style={planetLumiTitleStyle}>Lumi · {course.title}</strong>
            <span style={planetLumiMetaStyle}>{statusLabels[course.status]} · {copy.progress(progressPercent)}</span>
          </div>
        </div>
        <div style={planetLumiBubbleStyle}>{buildPlanetLumiMessage(course, copy)}</div>
        <div style={planetLumiSummaryRowStyle}>
          <span>{copy.regionCount(course.lessonCount ?? 0)}</span>
          <span>{copy.completedCount(course.completedLessonCount ?? 0, course.lessonCount ?? 0)}</span>
          {updatedLabel ? <span>{updatedLabel}</span> : null}
        </div>
      </div>
    </div>
  );
  return createPortal(tooltip, document.body);
}

function lumiStateForPlanet(course: PlanetCourseItem): ComponentProps<typeof LumiAvatar>['state'] {
  if (course.isInactive) return 'sleepy';
  if (course.status === 'completed') return 'celebrate';
  if (course.status === 'learning') return 'exploring';
  if (course.status === 'ready') return 'planet-hold';
  return 'thinking';
}

function buildPlanetLumiMessage(course: PlanetCourseItem, copy: DashboardMainCopy['planetMap']): string {
  if (course.isInactive) return copy.inactiveMessage;
  if (course.status === 'completed') return copy.completedMessage;
  if (course.status === 'learning') return copy.learningMessage;
  if (course.status === 'ready') return copy.readyMessage;
  return copy.draftMessage;
}

function planetCanvasStyle(isCompactLabelMode: boolean, isShortViewport: boolean): CSSProperties {
  return {
    position: 'relative',
    width: '100%',
    minHeight: isShortViewport ? '680px' : isCompactLabelMode ? '760px' : '560px',
    overflow: 'hidden',
    borderRadius: dashboardTokens.radius.lg,
    border: `1px solid ${dashboardTokens.colors.panelBorder}`,
    background: `linear-gradient(180deg, ${dashboardTokens.colors.panelGlassBg}, rgba(6, 12, 24, 0.72))`,
    boxShadow: dashboardTokens.shadow.panel,
    backdropFilter: 'blur(8px)',
    WebkitBackdropFilter: 'blur(8px)',
  };
}

const starSystemBannerStyle: CSSProperties = {
  position: 'absolute',
  top: 0,
  left: 28,
  right: 28,
  zIndex: 12,
  padding: '12px 18px',
  borderRadius: '24px 24px 18px 18px',
  background: 'linear-gradient(180deg, rgba(10, 20, 44, 0.96), rgba(10, 20, 44, 0.82))',
  backdropFilter: 'blur(16px)',
  boxShadow: '0 12px 34px rgba(0, 0, 0, 0.24)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '12px',
  flexWrap: 'wrap',
  borderBottom: '1px solid rgba(255,255,255,0.14)',
};

const starSystemBannerTextStyle: CSSProperties = {
  fontSize: '13px',
  fontWeight: 700,
  color: '#F4F7FF',
};

const starSystemBannerSummaryWrapper: CSSProperties = {
  display: 'grid',
  gap: '3px',
  color: '#C8D7FF',
  flex: '1 1 280px',
  minWidth: 0,
};

const starSystemBannerSummaryStyle: CSSProperties = {
  fontSize: '11px',
  color: '#C8D7FF',
  lineHeight: 1.3,
};

const starSystemBannerUpdatedStyle: CSSProperties = {
  fontSize: '10px',
  color: '#A3C2FF',
  lineHeight: 1.3,
};

const starSystemReturnButtonStyle: CSSProperties = {
  padding: '6px 12px',
  fontSize: '12px',
  borderRadius: '12px',
  border: '1px solid rgba(255,255,255,0.35)',
  background: 'rgba(255,255,255,0.12)',
  color: '#FFF',
  cursor: 'pointer',
  whiteSpace: 'nowrap',
};

const backgroundOverlayStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  pointerEvents: 'none',
};

const orbitLayerStyle: CSSProperties = {
  width: '100%',
  height: '100%',
  pointerEvents: 'none',
};

function orbitContainerStyle(
  isCompactLabelMode: boolean,
  isShortViewport: boolean,
  objectScale: number,
): CSSProperties {
  const size = (isShortViewport ? 320 : isCompactLabelMode ? 420 : 500) * objectScale;
  return {
    position: 'absolute',
    left: '50%',
    top: '50%',
    width: `${size}px`,
    height: `${size}px`,
    transform: 'translate(-50%, -50%)',
    pointerEvents: 'none',
  };
}

const cornerGlowStyle: CSSProperties = {
  position: 'absolute',
  width: 140,
  height: 140,
  borderRadius: '999px',
  background: 'radial-gradient(circle, rgba(109, 160, 255, 0.18) 0%, rgba(109, 160, 255, 0.06) 42%, rgba(109, 160, 255, 0) 70%)',
  filter: 'blur(10px)',
};

function centerHazeStyle(objectScale: number): CSSProperties {
  return {
    position: 'absolute',
    left: '50%',
    top: '50%',
    width: `${260 * objectScale}px`,
    height: `${260 * objectScale}px`,
    borderRadius: '999px',
    transform: 'translate(-50%, -50%)',
    background: 'radial-gradient(circle, rgba(255, 194, 84, 0.14) 0%, rgba(255, 194, 84, 0.05) 38%, rgba(255, 194, 84, 0) 74%)',
    filter: `blur(${14 * objectScale}px)`,
    pointerEvents: 'none',
  };
}

function sunNodeStyle(objectScale: number): CSSProperties {
  return {
    position: 'absolute',
    left: '50%',
    top: '50%',
    width: `${120 * objectScale}px`,
    height: `${120 * objectScale}px`,
    transform: 'translate(-50%, -50%)',
    pointerEvents: 'none',
    zIndex: 1,
  };
}

function sunAuraStyle(objectScale: number): CSSProperties {
  return {
    position: 'absolute',
    inset: `${-28 * objectScale}px`,
    borderRadius: '999px',
    background: 'radial-gradient(circle, rgba(255, 215, 121, 0.18) 0%, rgba(255, 179, 46, 0.08) 44%, rgba(255, 179, 46, 0) 74%)',
  };
}

function sunGlowStyle(objectScale: number): CSSProperties {
  return {
    position: 'absolute',
    inset: `${-12 * objectScale}px`,
    borderRadius: '999px',
    background: 'radial-gradient(circle, rgba(255, 233, 154, 0.52) 0%, rgba(255, 181, 56, 0.3) 46%, rgba(255, 140, 30, 0) 76%)',
    filter: `blur(${10 * objectScale}px)`,
  };
}

const sunCoreStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  borderRadius: '999px',
  background: 'radial-gradient(circle at 36% 34%, #fff4b6 0%, #ffe27a 20%, #ffb400 62%, #ff7a00 100%)',
  boxShadow: '0 0 26px rgba(255,200,0,0.38), 0 0 54px rgba(255,150,0,0.26)',
};

function planetSafeAreaStyle(
  isCompactLabelMode: boolean,
  isShortViewport: boolean,
  isStarSystemMode: boolean,
): CSSProperties {
  if (isStarSystemMode) {
    return {
      position: 'absolute',
      zIndex: 2,
      inset: isShortViewport
        ? '142px 18px 64px'
        : isCompactLabelMode
          ? '150px 18px 72px'
          : '142px 28px 36px',
    };
  }

  return {
    position: 'absolute',
    zIndex: 1,
    inset: isShortViewport
      ? '20px 18px 64px'
      : isCompactLabelMode
        ? '28px 18px 72px'
        : '30px 28px 36px',
  };
}

const emptyMapStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  textAlign: 'center',
  color: '#C8D1E8',
  fontSize: '15px',
  lineHeight: 1.7,
};

const planetButtonStyle: CSSProperties = {
  display: 'grid',
  justifyItems: 'center',
  gap: '12px',
  border: 'none',
  background: 'transparent',
  cursor: 'pointer',
  fontFamily: 'inherit',
};

function planetFloatAnimationName(index: number): string {
  const names = ['planetFloatA', 'planetFloatB', 'planetFloatC', 'planetFloatD'];
  return names[index % names.length] ?? 'planetFloatA';
}

function planetFloatDuration(index: number): string {
  const durations = ['8.2s', '9.1s', '7.6s', '8.8s'];
  return durations[index % durations.length] ?? '8.4s';
}

const planetShellStyle: CSSProperties = {
  display: 'grid',
  justifyItems: 'center',
  alignContent: 'center',
  gap: 12,
  textAlign: 'center',
};

const planetLabelWrapStyle: CSSProperties = {
  display: 'grid',
  gap: 4,
  minHeight: '44px',
  alignContent: 'start',
  justifyItems: 'center',
  width: '100%',
};

function planetLabelStyle(isCompactLabelMode: boolean, isShortViewport: boolean): CSSProperties {
  return {
    maxWidth: isShortViewport ? '118px' : isCompactLabelMode ? '132px' : '176px',
    padding: '3px 10px',
    borderRadius: '999px',
    background: 'rgba(4, 10, 20, 0.46)',
    fontSize: isShortViewport ? 12 : isCompactLabelMode ? 13 : 14,
    fontWeight: 700,
    lineHeight: 1.35,
    color: '#F8FAFF',
    textAlign: 'center',
    textShadow: '0 1px 6px rgba(4, 10, 20, 0.64)',
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  };
}

function planetNodeStyle(
  size: number,
  isHovered: boolean,
  isSelected: boolean,
  visualProfile: ReturnType<typeof getPlanetProgressVisualProfile>,
): CSSProperties {
  return {
    position: 'relative',
    width: `${size}px`,
    height: `${size}px`,
    display: 'grid',
    placeItems: 'center',
    transition: 'transform 0.2s ease',
    transform: isHovered
      ? `scale(${1.1 + visualProfile.scaleBoost})`
      : isSelected
        ? `scale(${1.03 + visualProfile.scaleBoost})`
        : `scale(${1 + visualProfile.scaleBoost})`,
  };
}

function planetGlowStyle(
  status: ReturnType<typeof normalizePlanetStatus>,
  isHovered: boolean,
  isSelected: boolean,
  visualProfile: ReturnType<typeof getPlanetProgressVisualProfile>,
  isInactive: boolean,
): CSSProperties {
  const accent = isInactive ? '#94A3B8' : canonicalStatusColor(status);
  return {
    position: 'absolute',
    inset: '-10%',
    borderRadius: '999px',
    background: `radial-gradient(circle, ${accent}${isSelected ? '52' : isHovered ? '42' : '24'} 0%, transparent 72%)`,
    filter: `blur(${8 + visualProfile.stage * 0.8}px)`,
    opacity: isInactive
      ? (isSelected ? 0.34 : isHovered ? 0.28 : 0.16)
      : Math.min(0.98, (isSelected ? 0.42 : isHovered ? 0.34 : 0.22) + visualProfile.glowOpacity),
  };
}

function planetOuterSoftRingStyle(
  status: ReturnType<typeof normalizePlanetStatus>,
  isSelected: boolean,
  visualProfile: ReturnType<typeof getPlanetProgressVisualProfile>,
  isInactive: boolean,
): CSSProperties {
  const accent = isInactive ? '#94A3B8' : canonicalStatusColor(status);
  return {
    position: 'absolute',
    inset: '-13%',
    borderRadius: '999px',
    border: `8px solid ${accent}${isInactive ? '28' : '1f'}`,
    boxShadow: isSelected
      ? `0 0 ${34 + visualProfile.stage * 2}px ${accent}${isInactive ? '44' : '36'}`
      : `0 0 ${22 + visualProfile.stage * 1.5}px ${accent}${isInactive ? '2e' : '24'}`,
    opacity: isInactive ? (isSelected ? 0.7 : 0.52) : Math.min(1, (isSelected ? 0.56 : 0.42) + visualProfile.ringOpacity * 0.45),
  };
}

function planetInnerRingStyle(
  status: ReturnType<typeof normalizePlanetStatus>,
  isSelected: boolean,
  visualProfile: ReturnType<typeof getPlanetProgressVisualProfile>,
  isInactive: boolean,
): CSSProperties {
  const accent = isInactive ? '#94A3B8' : canonicalStatusColor(status);
  return {
    position: 'absolute',
    inset: '-8%',
    borderRadius: '999px',
    border: `3px solid ${isSelected ? accent : `${accent}${isInactive ? '86' : 'a8'}`}`,
    boxShadow: isSelected
      ? `0 0 ${24 + visualProfile.stage * 1.8}px ${accent}${isInactive ? '46' : '42'}`
      : `0 0 ${16 + visualProfile.stage * 1.2}px ${accent}${isInactive ? '30' : '28'}`,
  };
}

function planetBaseSphereStyle(
  status: ReturnType<typeof normalizePlanetStatus>,
  size: number,
  isHovered: boolean,
  isSelected: boolean,
  visualProfile: ReturnType<typeof getPlanetProgressVisualProfile>,
  isInactive: boolean,
): CSSProperties {
  return {
    width: `${size}px`,
    height: `${size}px`,
    borderRadius: '999px',
    position: 'relative',
    overflow: 'hidden',
    opacity: 0.98,
    background: isInactive
      ? 'radial-gradient(circle at 34% 30%, rgba(226, 232, 240, 0.45) 0%, rgba(148, 163, 184, 0.52) 24%, rgba(71, 85, 105, 0.9) 66%, rgba(15, 23, 42, 0.98) 100%)'
      : 'radial-gradient(circle at 36% 32%, rgba(242, 247, 255, 0.22) 0%, rgba(88, 108, 146, 0.28) 26%, rgba(26, 35, 52, 0.9) 72%, rgba(10, 16, 28, 0.98) 100%)',
    boxShadow: planetShadow(status, isHovered, isSelected, visualProfile, isInactive),
    filter: isInactive
      ? `brightness(${0.82 + (isSelected ? 0.08 : isHovered ? 0.04 : 0)}) saturate(0.28)`
      : `brightness(${visualProfile.brightness + (isSelected ? 0.08 : isHovered ? 0.04 : 0)}) saturate(${visualProfile.saturate})`,
  };
}

function planetTextureOverlayStyle(
  family: ReturnType<typeof getPlanetTextureFamily>,
  stage: number,
  size: number,
): CSSProperties {
  return {
    position: 'absolute',
    inset: 0,
    borderRadius: '999px',
    opacity: 0.92,
    mixBlendMode: 'normal',
    filter: 'saturate(1.08) contrast(1.08) brightness(1.08)',
    ...getPlanetSpriteBackgroundStyle(family, stage, size),
    pointerEvents: 'none',
  };
}

const rotatingPlanetCanvasStyle: CSSProperties = {
  width: '100%',
  height: '100%',
  borderRadius: '999px',
  opacity: 0.96,
  filter: 'saturate(1.08) contrast(1.08) brightness(1.04)',
  pointerEvents: 'none',
};

const rotatingPlanetPreviewStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  width: '100%',
  height: '100%',
};

function planetShadow(
  status: ReturnType<typeof normalizePlanetStatus>,
  isHovered: boolean,
  isSelected: boolean,
  visualProfile: ReturnType<typeof getPlanetProgressVisualProfile>,
  isInactive: boolean,
): string {
  const accent = isInactive ? '#94A3B8' : canonicalStatusColor(status);
  const glow = visualProfile.glowOpacity + (isSelected ? 0.18 : isHovered ? 0.12 : 0.06);
  return `0 18px 36px rgba(0,0,0,0.24), 0 0 24px color-mix(in srgb, ${accent} ${Math.round(glow * 100)}%, transparent)`;
}

const planetHighlightStyle: CSSProperties = {
  position: 'absolute',
  inset: '8% 14% 34% 14%',
  borderRadius: '50%',
  background: 'radial-gradient(circle at 42% 38%, rgba(255,255,255,0.58) 0%, rgba(255,255,255,0.18) 34%, rgba(255,255,255,0) 72%)',
  pointerEvents: 'none',
};

const planetSurfaceShadeStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  borderRadius: '999px',
  background: 'linear-gradient(180deg, rgba(255,255,255,0.04) 0%, rgba(255,255,255,0) 42%, rgba(0,0,0,0.16) 100%)',
  pointerEvents: 'none',
};

const planetEdgeShadeStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  borderRadius: '999px',
  boxShadow: 'inset -12px -14px 22px rgba(8, 16, 28, 0.24)',
  pointerEvents: 'none',
};

const completedAccentStyle: CSSProperties = {
  position: 'absolute',
  inset: '10%',
  borderRadius: '999px',
  border: '1px dashed rgba(255, 235, 173, 0.34)',
  pointerEvents: 'none',
};

function getCoursePlanetTextureMap(course: {
  planetTextureMapAsset?: string | null;
  planetTextureMapRotationDurationSeconds?: number | null;
  planetTextureMapRotationDirection?: 'left' | 'right' | null;
}) {
  return {
    atlasURL: course.planetTextureMapAsset ?? fallbackTextureMap.atlasURL,
    rotationDurationSeconds: course.planetTextureMapRotationDurationSeconds ?? fallbackTextureMap.rotationDurationSeconds,
    rotationDirection: course.planetTextureMapRotationDirection ?? fallbackTextureMap.rotationDirection,
  };
}

function getCourseProgressPercent(course: PlanetCourseItem): number {
  if (course.status === 'completed') return 100;
  if ((course.progress ?? null) != null) {
    return Math.min(99, normalizePlanetProgressPercent(course.progress, course.status));
  }
  const lessonCount = Math.max(0, course.lessonCount ?? 0);
  const completedLessonCount = Math.max(0, course.completedLessonCount ?? 0);
  if (lessonCount > 0) {
    return Math.min(99, normalizePlanetProgressPercent(completedLessonCount / lessonCount, course.status));
  }
  if (course.status === 'learning') return 0;
  if (course.status === 'ready') return 0;
  return 0;
}

const tooltipStyle: CSSProperties = {
  position: 'absolute',
  minWidth: 156,
  maxWidth: 214,
  padding: '10px 12px',
  borderRadius: 14,
  background: 'rgba(7, 15, 29, 0.9)',
  border: '1px solid rgba(140, 176, 255, 0.18)',
  boxShadow: '0 18px 42px rgba(0, 0, 0, 0.28)',
  pointerEvents: 'none',
  zIndex: 8,
  animation: 'planetTooltipFade 0.18s ease',
};

const tooltipTitleStyle: CSSProperties = {
  fontSize: 13,
  fontWeight: 700,
  color: '#F7FAFF',
  lineHeight: 1.4,
};

const tooltipMetaRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: 8,
  flexWrap: 'wrap',
  marginTop: 6,
};

function tooltipStatusPillStyle(status: PlanetStatus): CSSProperties {
  return {
    display: 'inline-flex',
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: '22px',
    padding: '0 8px',
    borderRadius: '999px',
    background: `${canonicalStatusColor(normalizePlanetStatus(status))}22`,
    color: canonicalStatusColor(normalizePlanetStatus(status)),
    fontSize: 11,
    fontWeight: 700,
  };
}

const tooltipMetaStyle: CSSProperties = {
  fontSize: 12,
  color: 'rgba(200, 214, 240, 0.76)',
  lineHeight: 1.5,
};

function planetLumiTooltipStyle(clientX: number, clientY: number): CSSProperties {
  const tooltipWidth = 280;
  const padding = 14;
  const left = Math.max(
    padding,
    Math.min(window.innerWidth - tooltipWidth - padding, clientX + 16),
  );
  const top = Math.max(112, Math.min(window.innerHeight - 210, clientY - 84));
  return {
    position: 'fixed',
    left: `${left}px`,
    top: `${top}px`,
    zIndex: 9999,
    width: `${tooltipWidth}px`,
    pointerEvents: 'none',
    animation: 'planetTooltipFade 120ms ease both',
  };
}

const planetLumiCardStyle: CSSProperties = {
  display: 'grid',
  gap: '6px',
  padding: '10px',
  borderRadius: '14px',
  border: '1px solid rgba(148, 195, 255, 0.26)',
  background: 'rgba(6, 14, 28, 0.94)',
  backdropFilter: 'blur(16px)',
  WebkitBackdropFilter: 'blur(16px)',
  boxShadow: '0 8px 32px rgba(0, 0, 0, 0.4)',
};

const planetLumiHeaderStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: '10px',
};

const planetLumiTitleStyle: CSSProperties = {
  fontSize: '13px',
  fontWeight: 700,
  color: '#E8F2FF',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
};

const planetLumiMetaStyle: CSSProperties = {
  fontSize: '12px',
  color: '#8EA4C8',
};

const planetLumiBubbleStyle: CSSProperties = {
  fontSize: '11px',
  color: '#C8DCFF',
  lineHeight: 1.5,
  padding: '6px 8px',
  borderRadius: '10px',
  background: 'rgba(30, 60, 100, 0.38)',
};

const planetLumiSummaryRowStyle: CSSProperties = {
  display: 'flex',
  flexWrap: 'wrap',
  gap: '6px',
  fontSize: '11px',
  color: '#7A96C0',
  lineHeight: 1.3,
};
