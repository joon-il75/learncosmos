import type {
  DashboardCourseRecord,
  DashboardDraftStatus,
  DashboardPlanetStatus,
} from './types';
import { toCanonicalPlanetStatus } from './types';

export interface DashboardDraftListItemLike {
  id: string;
  source_query: string;
  title: string;
  status: DashboardDraftStatus;
  planet_type_id?: string | null;
  planet_type_name?: string | null;
  planet_type_asset?: string | null;
  planet_texture_map_id?: string | null;
  planet_texture_map_name?: string | null;
  planet_texture_map_asset?: string | null;
  planet_texture_map_rotation_duration_seconds?: number | null;
  planet_texture_map_rotation_direction?: 'left' | 'right' | null;
  updated_at: string;
}

export interface DashboardPlanetListItemLike {
  id: string;
  draft_id?: string;
  title: string;
  status: DashboardPlanetStatus;
  planet_type_id?: string | null;
  planet_type_name?: string | null;
  planet_type_asset?: string | null;
  planet_texture_map_id?: string | null;
  planet_texture_map_name?: string | null;
  planet_texture_map_asset?: string | null;
  planet_texture_map_rotation_duration_seconds?: number | null;
  planet_texture_map_rotation_direction?: 'left' | 'right' | null;
  completed_lesson_count?: number;
  progress?: number | null;
  is_inactive?: boolean;
  last_accessed_at?: string | null;
  updated_at: string;
  destination: string;
}

export function mapDraftStatusToPlanetStatus(
  status: DashboardDraftStatus,
): DashboardPlanetStatus {
  if (status === 'learning') return 'learning';
  if (status === 'confirmed') return 'ready';
  if (status === 'archived') return 'completed';
  return 'draft';
}

export function isDashboardPlanetStatus(
  value: string | null,
): value is DashboardPlanetStatus {
  return (
    value === 'draft' ||
    value === 'ready' ||
    value === 'learning' ||
    value === 'completed'
  );
}

export function planetStatusSortWeight(status: DashboardPlanetStatus): number {
  const canonicalStatus = toCanonicalPlanetStatus(status);

  if (canonicalStatus === 'exploring') return 0;
  if (canonicalStatus === 'ready') return 1;
  if (canonicalStatus === 'draft') return 2;
  return 3;
}

export function buildDashboardCourseRecordFromDraft(
  draft: DashboardDraftListItemLike,
  levelCount: number,
  lessonCount: number,
  plannedLevelCount = 0,
  plannedLessonCount = 0,
  completedLessonCount = 0,
  progress: number | null = null,
): DashboardCourseRecord {
  const fallbackTitle = draft.source_query.trim() || '새 행성';

  return {
    id: draft.id,
    draftId: draft.id,
    title: draft.title?.trim() || fallbackTitle,
    status: mapDraftStatusToPlanetStatus(draft.status),
    updatedAt: draft.updated_at,
    levelCount,
    lessonCount,
    plannedLevelCount,
    plannedLessonCount,
    completedLessonCount,
    progress,
    sourceQuery: draft.source_query,
    draftStatus: draft.status,
    planetTypeId: draft.planet_type_id ?? null,
    planetTypeName: draft.planet_type_name ?? null,
    planetTypeAsset: draft.planet_type_asset ?? null,
    planetTextureMapId: draft.planet_texture_map_id ?? null,
    planetTextureMapName: draft.planet_texture_map_name ?? null,
    planetTextureMapAsset: draft.planet_texture_map_asset ?? null,
    planetTextureMapRotationDurationSeconds: draft.planet_texture_map_rotation_duration_seconds ?? null,
    planetTextureMapRotationDirection: draft.planet_texture_map_rotation_direction ?? null,
    destination: `/dashboard/course-drafts/${draft.id}`,
  };
}

function mapPlanetStatusToDraftStatus(status: DashboardPlanetStatus): DashboardDraftStatus {
  if (status === 'ready') return 'confirmed';
  if (status === 'learning') return 'learning';
  if (status === 'completed') return 'archived';
  return 'draft';
}

export function buildDashboardCourseRecordFromPlanet(
  planet: DashboardPlanetListItemLike,
  levelCount: number,
  lessonCount: number,
  completedLessonCount = 0,
  progress: number | null = null,
): DashboardCourseRecord {
  const fallbackTitle = planet.title.trim() || '행성';
  const destination =
    planet.is_inactive && planet.draft_id
      ? `/dashboard/course-drafts/${planet.draft_id}?section=planning&from=diary`
      : planet.status === 'completed'
      ? `/dashboard/planets/shared/${planet.id}`
      : `/dashboard/planets/learning/${planet.id}`;

  return {
    id: planet.id,
    draftId: planet.draft_id ?? null,
    title: fallbackTitle,
    status: planet.status,
    lastAccessedAt: planet.last_accessed_at ?? null,
    updatedAt: planet.updated_at,
    levelCount,
    lessonCount,
    plannedLevelCount: 0,
    plannedLessonCount: 0,
    completedLessonCount,
    progress,
    sourceQuery: fallbackTitle,
    draftStatus: mapPlanetStatusToDraftStatus(planet.status),
    isInactive: Boolean(planet.is_inactive),
    planetTypeId: planet.planet_type_id ?? null,
    planetTypeName: planet.planet_type_name ?? null,
    planetTypeAsset: planet.planet_type_asset ?? null,
    planetTextureMapId: planet.planet_texture_map_id ?? null,
    planetTextureMapName: planet.planet_texture_map_name ?? null,
    planetTextureMapAsset: planet.planet_texture_map_asset ?? null,
    planetTextureMapRotationDurationSeconds: planet.planet_texture_map_rotation_duration_seconds ?? null,
    planetTextureMapRotationDirection: planet.planet_texture_map_rotation_direction ?? null,
    destination,
  };
}
