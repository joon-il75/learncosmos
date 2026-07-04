'use client';

import AtmosphericPlanetPreview from '@/components/dashboard/AtmosphericPlanetPreview';
import { fallbackTextureMap, type ActivePlanetTextureMap } from '@/components/dashboard/useActivePlanetTextureMap';
import type { PageSystemChunk } from '@/components/dashboard/GalaxyMap';
import {
  getGalaxyCourseProgressPercent,
  mapWindowPageNodeStyle,
  mapSystemScaleStyle,
  mapSystemCoreStyle,
  mapSystemSunStyle,
  mapSystemPlanetStyle,
  mapSystemPlanetRingStyle,
  mapSystemPlanetBodyStyle,
  mapSystemPlanetCanvasStyle,
  mapSystemPlanetPreviewStyle,
  mapSystemLabelStyle,
} from '@/components/dashboard/GalaxyMapStyles';
import type { DashboardCourseRecord } from '@/lib/world-ui-engine/types';

function getCoursePlanetTextureMap(course: DashboardCourseRecord): ActivePlanetTextureMap {
  return {
    atlasURL: course.planetTextureMapAsset ?? fallbackTextureMap.atlasURL,
    rotationDurationSeconds: course.planetTextureMapRotationDurationSeconds ?? fallbackTextureMap.rotationDurationSeconds,
    rotationDirection: course.planetTextureMapRotationDirection ?? fallbackTextureMap.rotationDirection,
  };
}

function MapSystemPlanetBody({
  course,
  isHighlighted,
  objectScale,
  activeTextureMap,
}: {
  course: DashboardCourseRecord;
  isHighlighted: boolean;
  objectScale: number;
  activeTextureMap: ActivePlanetTextureMap;
}) {
  const isInactive = Boolean(course.isInactive);
  const normalizedProgressPercent = getGalaxyCourseProgressPercent(course);
  const useRotatingTexture = !isInactive && Boolean(course.planetTextureMapAsset);

  return (
    <div style={mapSystemPlanetBodyStyle(course, isHighlighted, objectScale)}>
      {useRotatingTexture ? (
        <AtmosphericPlanetPreview
          atlasURL={activeTextureMap.atlasURL}
          progressPercent={normalizedProgressPercent}
          size={14 * objectScale}
          durationSeconds={activeTextureMap.rotationDurationSeconds}
          direction={activeTextureMap.rotationDirection}
          playing
          preset="galaxy-mini"
          style={mapSystemPlanetPreviewStyle}
          canvasStyle={mapSystemPlanetCanvasStyle}
        />
      ) : null}
    </div>
  );
}

interface GalaxySystemNodeProps {
  pageSystem: PageSystemChunk;
  index: number;
  totalNodes: number;
  isHighlighted: boolean;
  useColumnLayout: boolean;
  objectScale: number;
  onSystemSelect: (system: PageSystemChunk) => void;
  onPageHoverStart: (system: PageSystemChunk, x: number, y: number) => void;
  onPageHoverMove: (pageNumber: number, x: number, y: number) => void;
  onPageHoverEnd: () => void;
}

export function GalaxySystemNode({
  pageSystem,
  index,
  totalNodes,
  isHighlighted,
  useColumnLayout,
  objectScale,
  onSystemSelect,
  onPageHoverStart,
  onPageHoverMove,
  onPageHoverEnd,
}: GalaxySystemNodeProps) {
  return (
    <div
      key={`orbit-page-${pageSystem.pageNumber}`}
      style={mapWindowPageNodeStyle(index, totalNodes, isHighlighted, useColumnLayout, objectScale)}
      onClick={() => onSystemSelect(pageSystem)}
      onMouseEnter={(e) => onPageHoverStart(pageSystem, e.clientX, e.clientY)}
      onMouseMove={(e) => onPageHoverMove(pageSystem.pageNumber, e.clientX, e.clientY)}
      onMouseLeave={onPageHoverEnd}
      onFocus={(e) => {
        const rect = e.currentTarget.getBoundingClientRect();
        onPageHoverStart(pageSystem, rect.left + rect.width / 2, rect.top);
      }}
      onBlur={onPageHoverEnd}
      role="button"
      tabIndex={0}
      aria-label={`항성계 ${pageSystem.pageNumber}: ${pageSystem.summaryTitle}`}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault();
          onSystemSelect(pageSystem);
        }
      }}
    >
      <div style={mapSystemScaleStyle(isHighlighted)}>
        <div style={mapSystemCoreStyle(isHighlighted, objectScale)}>
          <div style={mapSystemSunStyle(isHighlighted, objectScale)} />
          {pageSystem.courses.map((course: DashboardCourseRecord, courseIndex: number) => (
            <div
              key={`${pageSystem.pageNumber}-${course.id}`}
              style={mapSystemPlanetStyle(courseIndex, pageSystem.courses.length, objectScale)}
              title={course.title}
            >
              <div style={mapSystemPlanetRingStyle(course.status, isHighlighted, Boolean(course.isInactive), objectScale)} />
              <MapSystemPlanetBody
                course={course}
                isHighlighted={isHighlighted}
                objectScale={objectScale}
                activeTextureMap={getCoursePlanetTextureMap(course)}
              />
            </div>
          ))}
        </div>
        <div style={mapSystemLabelStyle(isHighlighted, objectScale)}>{pageSystem.pageNumber}</div>
      </div>
    </div>
  );
}
