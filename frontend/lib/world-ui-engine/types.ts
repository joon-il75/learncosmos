import type {
  DashboardHoveredSystemSummary,
  DashboardHoverSource,
} from '@/lib/lumi/lumiEngineTypes';
import type { DashboardLumiCourse } from '@/lib/lumi/lumiRuntimeEngine';

export type DashboardDraftStatus = 'draft' | 'learning' | 'confirmed' | 'archived';

export type DashboardPlanetStatus = 'draft' | 'ready' | 'learning' | 'completed';
export type DashboardPlanetStatusCanonical = 'draft' | 'ready' | 'exploring' | 'completed';
export type DashboardStatusFilter = 'all' | DashboardPlanetStatus;

export type DashboardCreatorPlanetState = 'idle' | 'focused' | 'launching';

export type DashboardDeviceKind = 'se' | 'mobile' | 'tablet' | 'desktop' | 'wide-desktop';
export type DashboardViewportKind = 'default' | 'short';
export type DashboardMode = 'galaxy' | 'system-selector' | 'planet-list' | 'star-system';
export type DashboardSystemLayout = 'orbit' | 'grid' | 'stack';
export type DashboardWorldVisibility = 'full' | 'partial' | 'minimal' | 'off';
export type DashboardEffectVisibility = 'full' | 'reduced' | 'off';
export type DashboardLumiMode = 'docked' | 'bottom-sheet' | 'button-only' | 'hidden';
export type DashboardSortKey = 'updated_desc' | 'updated_asc' | 'title_asc' | 'title_desc';

export interface DashboardCourseRecord {
  id: string;
  draftId?: string | null;
  title: string;
  status: DashboardPlanetStatus;
  lastAccessedAt?: string | null;
  updatedAt: string;
  levelCount: number;
  lessonCount: number;
  plannedLevelCount: number;
  plannedLessonCount: number;
  completedLessonCount: number;
  progress: number | null;
  sourceQuery: string;
  draftStatus: DashboardDraftStatus;
  isInactive?: boolean;
  planetTypeId?: string | null;
  planetTypeName?: string | null;
  planetTypeAsset?: string | null;
  planetTextureMapId?: string | null;
  planetTextureMapName?: string | null;
  planetTextureMapAsset?: string | null;
  planetTextureMapRotationDurationSeconds?: number | null;
  planetTextureMapRotationDirection?: 'left' | 'right' | null;
  destination: string;
}

export interface DashboardSelectionState {
  hoveredCourseId: string | null;
  selectedCourseId: string | null;
}

export interface DashboardFilterState {
  query: string;
  page: number;
  activeStatusFilter: DashboardStatusFilter;
  sortKey: DashboardSortKey;
}

export interface DashboardCTAState {
  creatorPlanetState: DashboardCreatorPlanetState;
  isCreatingCourse: boolean;
  creationError: string | null;
}

export interface DashboardDeviceState {
  viewportWidth: number;
  viewportHeight: number;
  isCompactLayout: boolean;
  isPhoneLayout: boolean;
  isShortViewport: boolean;
}

export interface DashboardLumiStateInput {
  hasByok: boolean;
  isPanelOpen: boolean;
  hasSeenEntryLumi: boolean;
}

export interface DashboardDataState {
  courses: DashboardCourseRecord[];
  isLoading: boolean;
}

export interface DashboardUIEngineInput {
  data: DashboardDataState;
  filter: DashboardFilterState;
  selection: DashboardSelectionState;
  cta: DashboardCTAState;
  device: DashboardDeviceState;
  lumi: DashboardLumiStateInput;
  mode: DashboardMode;
  selectedSystemId: string | null;
  hoveredSystemId: string | null;
}

export interface DashboardViewModelInput {
  courses: DashboardCourseRecord[];
  searchQuery: string;
  page: number;
  activeStatusFilter: DashboardStatusFilter;
  selectedCourseId: string | null;
  hoveredCourseId: string | null;
  pageSize?: number;
  mode: DashboardMode;
  selectedSystemId: string | null;
  hoveredSystemId: string | null;
  deviceProfile: DashboardDeviceProfile;
  sortKey: DashboardSortKey;
}

export interface DashboardDeviceProfile {
  deviceKind: DashboardDeviceKind;
  viewportKind: DashboardViewportKind;
  useCompactMapLabels: boolean;
  mapBreakpoint: 'desktop' | 'mobile';
}

export interface DashboardPlanetStateViewModel {
  courseId: string;
  canonicalStatus: DashboardPlanetStatusCanonical;
  terraformPhase: 'barren' | 'atmosphere' | 'biosphere' | 'living' | 'civilized';
  terraformStage: 1 | 2 | 3 | 4 | 5;
  toneFamily: 'rock' | 'ocean' | 'cloud' | 'ice' | 'gas' | 'crater';
  civilizationLevel: 0 | 1 | 2 | 3;
}

export interface DashboardPlanetCardViewModel {
  id: string;
  title: string;
  originalTopic?: string;
  lessonCount: number;
  updatedAtLabel: string;
  canonicalStatus: DashboardPlanetStatusCanonical;
  progressLabel?: string;
  isSelected: boolean;
  isHovered: boolean;
  planetState: DashboardPlanetStateViewModel;
}

export interface DashboardSystemViewModel {
  id: string;
  index: number;
  title: string;
  courseCount: number;
  courses: DashboardPlanetCardViewModel[];
  isSelected: boolean;
  statusCount: {
    draft: number;
    ready: number;
    learning: number;
    completed: number;
    inactive: number;
  };
  summaryTitle: string;
  updatedAt: string | null;
}

export interface DashboardVisibleLayers {
  core: boolean;
  navigation: boolean;
  world: DashboardWorldVisibility;
  effect: DashboardEffectVisibility;
}

export interface DashboardEmptyState {
  title: string;
  description: string;
  ctaLabel: string;
}

export interface DashboardCTAViewModel {
  title: string;
  placeholder: string;
  submitLabel: string;
}

export interface DashboardCourseCollection {
  allCourses: DashboardCourseRecord[];
  filteredCourses: DashboardCourseRecord[];
  pagedCourses: DashboardCourseRecord[];
  totalPages: number;
  currentPage: number;
}

export interface DashboardStatusCounts {
  draft: number;
  ready: number;
  learning: number;
  completed: number;
}

export interface DashboardSelectionDerivedState {
  selectedCourse: DashboardCourseRecord | null;
  hoveredCourse: DashboardCourseRecord | null;
}

export interface DashboardUIEngineSnapshot {
  courses: DashboardCourseCollection;
  counts: DashboardStatusCounts;
  selection: DashboardSelectionDerivedState;
  deviceProfile: DashboardDeviceProfile;
  visibleLayers: DashboardVisibleLayers;
  dashboardMode: DashboardMode;
  systemLayout: DashboardSystemLayout;
  lumiMode: DashboardLumiMode;
  systems: DashboardSystemViewModel[];
  selectedSystemId: string | null;
  filteredCourseCount: number;
  totalCourseCount: number;
  emptyState: DashboardEmptyState | null;
  cta: DashboardCTAViewModel;
  hoveredSystemId: string | null;
}

export interface DashboardAdapterResult {
  input: DashboardUIEngineInput;
  snapshot: DashboardUIEngineSnapshot;
}

export interface DashboardLumiActionHandlers {
  newCourse: () => void;
  viewRelatedCourse: (courseId: string) => void;
  changeFilter: (status: DashboardStatusFilter) => void;
  clearFilter: () => void;
  openSettings: () => void;
}

export interface DashboardLumiControllerInput {
  courses: DashboardLumiCourse[];
  selectedCourseId: string | null;
  hoveredCourseId: string | null;
  activeStatusFilter: DashboardStatusFilter;
  userName: string;
  hasByok: boolean;
  isPhoneLayout: boolean;
  isExpanded: boolean;
  actionHandlers: DashboardLumiActionHandlers;
  hoverSource: DashboardHoverSource | null;
  hoveredSystemSummary: DashboardHoveredSystemSummary | null;
}

export function toCanonicalPlanetStatus(
  status: DashboardPlanetStatus,
): DashboardPlanetStatusCanonical {
  if (status === 'learning') return 'exploring';
  return status;
}
