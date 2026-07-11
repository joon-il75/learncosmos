'use client';

import {
  Suspense,
  useEffect,
  useMemo,
  useRef,
  useReducer,
  useState,
} from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import LearnerHeaderActions, { getLearnerHeaderActionsCopy } from '@/components/common/LearnerHeaderActions';
import LearnerAppShell from '@/components/navigation/LearnerAppShell';
import GalaxyMap, {
  buildSystemStatusLine,
  formatSystemUpdatedAt,
} from '@/components/dashboard/GalaxyMap';
import type { PageSystemChunk, PageHoverLumiState } from '@/components/dashboard/GalaxyMap';
import CourseListPanel from '@/components/dashboard/CourseListPanel';
import { useLumiController } from '@/components/dashboard/useLumiController';
import { isDashboardPlanetStatus } from '@/lib/world-ui-engine/stateNormalizer';
import { buildDashboardAdapter } from '@/lib/world-ui-engine/dashboardAdapter';
import { buildDashboardLumiAdapter } from '@/lib/world-ui-engine/lumiAdapter';
import { buildDashboardDeviceState } from '@/lib/world-ui-engine/deviceProfile';
import { resolveDashboardSelectedCourseId } from '@/lib/world-ui-engine/selectionPolicy';
import { buildDashboardRouteParams } from '@/lib/world-ui-engine/routeSync';
import { normalizeLocale } from '@/lib/i18n/locales';
import { getDashboardMainCopy } from '@/lib/i18n/pages/dashboardMain';
import {
  DEFAULT_DASHBOARD_UI_ENGINE_SETTINGS,
  normalizeDashboardUIEngineRuntimeSettings,
  type DashboardUIEngineRuntimeSettingsResponse,
} from '@/lib/world-ui-engine/uiEngineConfig';
import type { DashboardSortKey } from '@/lib/world-ui-engine/types';
import {
  dashboardModeReducer,
  getPageSizeForMode,
  initialDashboardModeState,
  type DashboardModeContract,
} from '@/lib/dashboard/dashboardState';
import { filterCourses, sortCourses } from '@/lib/dashboard/coursePipeline';
import { chunkCoursesToSystems } from '@/lib/dashboard/groupCoursesToSystems';

import { useDashboardLoader } from './useDashboardLoader';
import DashboardBackground from './DashboardBackground';
import DashboardTodayTaskSection from './DashboardTodayTaskSection';
import CosmosMobileFirstView from './CosmosMobileFirstView';
import { useGoalCreation } from './useGoalCreation';
import { useSystemNavigation } from './useSystemNavigation';
import { GoalFlowLoadingOverlay } from '@/components/goal-interview/GoalFlowLoadingOverlay';
import { PAGE_SIZE } from './_constants';
import {
  isDashboardSortKey,
  getDisplayName,
  findSystemIdForCourse,
} from './_helpers';
import type { PlanetFilter, CourseItem } from './_types';
import {
  pageStyle,
  pageChromeStyle,
  mainStyle,
  headerBlockStyle,
  heroSubtitleStyle,
  mapListGridStyle,
  mapColumnStyle,
  responsiveListColumnStyle,
  loadingPageStyle,
  loadingCardStyle,
} from './_styles';

// ─────────────────────────────────────────────────────────────
// DashboardContent
// ─────────────────────────────────────────────────────────────
function DashboardContent() {
  const router = useRouter();
  const searchParams = useSearchParams();

  // ── 데이터 로딩 ─────────────────────────────────────────
  const { user, aiSettings, courses, todayTask, isLoading } = useDashboardLoader();
  const locale = normalizeLocale(user?.ui_locale);
  const dashboardCopy = getDashboardMainCopy(locale);
  const learnerHeaderCopy = getLearnerHeaderActionsCopy(locale);

  // ── 대시보드 모드 상태 ───────────────────────────────────
  const [dashboardModeState, dispatchDashboardMode] = useReducer(
    dashboardModeReducer,
    initialDashboardModeState,
  );
  const {
    mode,
    page,
    searchQuery: query,
    statusFilter: activeStatusFilter,
    sortKey,
    selectedSystemId,
    hoveredSystemId,
    hoveredCourseId,
    selectedCourseId,
  } = dashboardModeState;

  // ── UI 상태 ──────────────────────────────────────────────
  const [selectedMapPage, setSelectedMapPage] = useState<number | null>(null);
  const [ctaQuery, setCtaQuery] = useState('');
  const [isQueryInputFocused, setIsQueryInputFocused] = useState(false);
  const [isLumiPanelOpen, setIsLumiPanelOpen] = useState(false);
  const [hasSeenEntryLumi, setHasSeenEntryLumi] = useState(false);
  const [showInactiveCourses, setShowInactiveCourses] = useState(false);
  const [hoverSource, setHoverSource] = useState<'system' | 'course' | 'course-in-galaxy' | null>(null);
  const [pageHoverLumi, setPageHoverLumi] = useState<PageHoverLumiState | null>(null);
  const [uiEngineSettings, setUIEngineSettings] = useState(DEFAULT_DASHBOARD_UI_ENGINE_SETTINGS);

  // ── 디바이스 상태 ────────────────────────────────────────
  const [isCompactLayout, setIsCompactLayout] = useState(false);
  const [isPhoneLayout, setIsPhoneLayout] = useState(false);
  const [isShortViewport, setIsShortViewport] = useState(false);
  const [viewportWidth, setViewportWidth] = useState(1280);
  const [viewportHeight, setViewportHeight] = useState(720);

  const queryInputRef = useRef<HTMLInputElement>(null);
  const pageHoverTimerRef = useRef<number | null>(null);

  const {
    isGoalSubmitting,
    goalFlowStage,
    goalCreationError,
    handleGoalSubmit,
    clearGoalError,
  } = useGoalCreation({
    queryInputFocusFn: () => queryInputRef.current?.focus(),
    errorMessages: dashboardCopy.cta.errors,
  });

  // ── URL → 상태 동기화 ────────────────────────────────────
  useEffect(() => {
    const nextQuery = searchParams.get('q') ?? '';
    const rawPage = searchParams.get('page');
    const nextPage = Number(rawPage ?? '1');
    const nextStatus = searchParams.get('status');
    const nextModeValue = searchParams.get('mode');
    const normalizedNextMode: DashboardModeContract =
      nextModeValue === 'star-system' ? 'star-system' : 'galaxy';
    const nextSystemId =
      normalizedNextMode === 'star-system' ? searchParams.get('system') ?? null : null;
    const normalizedPage =
      Number.isFinite(nextPage) && nextPage > 0 ? Math.floor(nextPage) : 1;
    const nextSortValue = searchParams.get('sort');

    dispatchDashboardMode({ type: 'setMode', payload: normalizedNextMode });
    setSelectedMapPage(
      rawPage && Number.isFinite(nextPage) && nextPage > 0 ? Math.floor(nextPage) : null,
    );
    dispatchDashboardMode({
      type: 'setStatusFilter',
      payload: isDashboardPlanetStatus(nextStatus) ? nextStatus : 'all',
    });
    dispatchDashboardMode({ type: 'setSearchQuery', payload: nextQuery });
    dispatchDashboardMode({
      type: 'setSortKey',
      payload: isDashboardSortKey(nextSortValue) ? nextSortValue : 'updated_desc',
    });
    dispatchDashboardMode({ type: 'setPage', payload: normalizedPage });
    dispatchDashboardMode({ type: 'setSelectedSystemId', payload: nextSystemId });
  }, [searchParams]);

  // ── UI 엔진 설정 로드 ────────────────────────────────────
  useEffect(() => {
    fetch('/api/v1/public/ui-engine-settings', { cache: 'no-store' })
      .then((r) => r.json())
      .then((payload: DashboardUIEngineRuntimeSettingsResponse) => {
        setUIEngineSettings(normalizeDashboardUIEngineRuntimeSettings(payload.settings));
      })
      .catch(() => undefined);
  }, []);

  // ── 디바이스 크기 감지 ───────────────────────────────────
  useEffect(() => {
    const update = () => {
      const s = buildDashboardDeviceState(
        { width: window.innerWidth, height: window.innerHeight },
        uiEngineSettings,
      );
      setViewportWidth(s.viewportWidth);
      setViewportHeight(s.viewportHeight);
      setIsCompactLayout(s.isCompactLayout);
      setIsPhoneLayout(s.isPhoneLayout);
      setIsShortViewport(s.isShortViewport);
    };
    update();
    window.addEventListener('resize', update);
    return () => window.removeEventListener('resize', update);
  }, [uiEngineSettings]);

  // ── 엔트리 Lumi 딜레이 ───────────────────────────────────
  useEffect(() => {
    if (!hasSeenEntryLumi && !isLoading) {
      const timer = window.setTimeout(() => setHasSeenEntryLumi(true), 1200);
      return () => window.clearTimeout(timer);
    }
  }, [hasSeenEntryLumi, isLoading]);

  // ── hover 타이머 cleanup ─────────────────────────────────
  useEffect(
    () => () => {
      if (pageHoverTimerRef.current !== null) {
        window.clearTimeout(pageHoverTimerRef.current);
      }
    },
    [],
  );

  const visibleCourses = useMemo(
    () => showInactiveCourses ? courses : courses.filter((course) => !course.isInactive),
    [courses, showInactiveCourses],
  );
  const inactiveCourseCount = useMemo(
    () => courses.filter((course) => course.isInactive).length,
    [courses],
  );

  // ── 뷰모델 ──────────────────────────────────────────────
  const dashboardAdapter = useMemo(
    () =>
      buildDashboardAdapter({
        data: { courses: visibleCourses, isLoading },
        filter: { query, page, activeStatusFilter, sortKey },
        selection: { hoveredCourseId, selectedCourseId },
        cta: { creatorPlanetState: 'idle', isCreatingCourse: false, creationError: null },
        device: { viewportWidth, viewportHeight, isCompactLayout, isPhoneLayout, isShortViewport },
        lumi: {
          hasByok: Boolean(aiSettings?.has_api_key),
          isPanelOpen: isLumiPanelOpen,
          hasSeenEntryLumi,
        },
        mode,
        selectedSystemId,
        hoveredSystemId,
      }),
    [
      activeStatusFilter, aiSettings?.has_api_key, visibleCourses, hasSeenEntryLumi,
      hoveredCourseId, hoveredSystemId, mode, selectedSystemId, isCompactLayout,
      isLoading, isLumiPanelOpen, isPhoneLayout, isShortViewport,
      page, query, sortKey, selectedCourseId, viewportHeight, viewportWidth,
    ],
  );

  const dashboardViewModel = dashboardAdapter.snapshot;
  const filteredCourses = dashboardViewModel.courses.filteredCourses;
  const pagedCourses = dashboardViewModel.courses.pagedCourses;
  const selectedSystem =
    dashboardViewModel.systems.find((s) => s.id === dashboardViewModel.selectedSystemId) ?? null;
  const totalPages = dashboardViewModel.courses.totalPages;
  const systemsPerGalaxyPage = Math.max(1, Math.ceil(getPageSizeForMode('galaxy') / PAGE_SIZE));
  const systemStartIndex = mode === 'galaxy'
    ? (page - 1) * systemsPerGalaxyPage + 1
    : page;

  const chunkedPagedCourses = useMemo(
    () => chunkCoursesToSystems(pagedCourses, PAGE_SIZE, systemStartIndex),
    [pagedCourses, systemStartIndex],
  );
  const pageSystems = useMemo(
    () => chunkedPagedCourses.map((sys) => ({ ...sys, pageNumber: sys.index })),
    [chunkedPagedCourses],
  );

  const hoveredPageSystem = pageHoverLumi
    ? pageSystems.find((ps) => ps.pageNumber === pageHoverLumi.pageNumber) ?? null
    : null;

  const initialEntryTarget = useMemo(() => {
    const learningCourses = visibleCourses
      .filter((course) => course.status === 'learning')
      .sort((a, b) => {
        const aLastAccessed = a.lastAccessedAt ? new Date(a.lastAccessedAt).getTime() : 0;
        const bLastAccessed = b.lastAccessedAt ? new Date(b.lastAccessedAt).getTime() : 0;
        if (aLastAccessed !== bLastAccessed) {
          return bLastAccessed - aLastAccessed;
        }
        return new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime();
      });
    const completedCourses = visibleCourses
      .filter((course) => course.status === 'completed')
      .sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime());
    const targetCourse = learningCourses[0] ?? completedCourses[0] ?? null;

    if (!targetCourse) return null;

    const allSystems = chunkCoursesToSystems(visibleCourses, PAGE_SIZE, 1);
    const systemId = findSystemIdForCourse(targetCourse.id, allSystems);
    if (!systemId) return null;

    return {
      courseId: targetCourse.id,
      systemId,
    };
  }, [visibleCourses]);

  // ── 로그인 직후 기본 진입: 코스 없으면 Galaxy, 있으면 최근 코스 Star Map ──
  useEffect(() => {
    if (isLoading) return;

    const hasExplicitRouteState =
      Boolean(searchParams.get('mode')) ||
      Boolean(searchParams.get('system')) ||
      Boolean(searchParams.get('returnCourse')) ||
      Boolean(searchParams.get('q')) ||
      Boolean(searchParams.get('status')) ||
      Boolean(searchParams.get('page')) ||
      Boolean(searchParams.get('sort'));

    if (hasExplicitRouteState) {
      return;
    }

    if (!initialEntryTarget) {
      return;
    }

    router.replace(
      `/dashboard?mode=star-system&system=${initialEntryTarget.systemId}`,
      { scroll: false },
    );
  }, [initialEntryTarget, isLoading, router, searchParams]);

  // ── returnCourse 파라미터 처리 (코스 드래프트 → 소속 스타 시스템으로 복귀) ──
  useEffect(() => {
    const returnCourseId = searchParams.get('returnCourse');
    if (!returnCourseId || visibleCourses.length === 0) return;
    const allSystems = chunkCoursesToSystems(visibleCourses, PAGE_SIZE, 1);
    const systemId = findSystemIdForCourse(returnCourseId, allSystems);
    if (!systemId) {
      router.replace('/dashboard', { scroll: false });
      return;
    }
    router.replace(`/dashboard?mode=star-system&system=${systemId}`, { scroll: false });
  }, [router, searchParams, visibleCourses]);

  const pageSystemMap = useMemo(() => {
    const map = new Map<string, PageSystemChunk>();
    pageSystems.forEach((ps) => map.set(ps.id, ps));
    return map;
  }, [pageSystems]);

  const selectedSystemChunk =
    selectedSystemId
      ? pageSystemMap.get(selectedSystemId) ?? null
      : mode === 'star-system'
        ? (pageSystems[0] ?? null)
        : null;

  const starSystemStatusLine = selectedSystemChunk
    ? buildSystemStatusLine(selectedSystemChunk.statusCount, dashboardCopy.galaxy)
    : null;
  const starSystemUpdatedLabel = selectedSystemChunk
    ? formatSystemUpdatedAt(selectedSystemChunk.updatedAt, locale, dashboardCopy.galaxy)
    : null;
  const selectedSystemTitleFallback =
    selectedSystem?.title ??
    selectedSystemChunk?.summaryTitle ??
    (selectedSystemChunk ? dashboardCopy.galaxy.selectedSystemFallback(selectedSystemChunk.index) : undefined);
  const hoveredSystemChunk = hoveredSystemId
    ? pageSystemMap.get(hoveredSystemId) ?? null
    : null;
  const hoveredSystemSummary = hoveredSystemChunk
    ? {
        title: hoveredSystemChunk.summaryTitle,
        statusCount: hoveredSystemChunk.statusCount,
        updatedAt: hoveredSystemChunk.updatedAt,
        leadCourseId: hoveredSystemChunk.courses[0]?.id ?? null,
        leadCourseTitle: hoveredSystemChunk.courses[0]?.title ?? null,
      }
    : null;

  // ── selectedCourseId 동기화 ──────────────────────────────
  useEffect(() => {
    const next = resolveDashboardSelectedCourseId({ pagedCourses, selectedCourseId });
    if (next !== selectedCourseId) {
      dispatchDashboardMode({ type: 'setSelectedCourseId', payload: next });
    }
  }, [dispatchDashboardMode, pagedCourses, selectedCourseId]);

  // ── 라우트 동기화 헬퍼 ───────────────────────────────────
  const syncRouteState = (
    nextPage: number,
    nextQuery: string,
    nextStatus: PlanetFilter,
    nextMode: DashboardModeContract = mode,
    nextSystemId: string | null = selectedSystemId,
    nextSortKey: DashboardSortKey = sortKey,
  ) => {
    const params = buildDashboardRouteParams({
      currentSearch: searchParams.toString(),
      nextPage,
      nextQuery,
      nextStatus,
      nextMode,
      nextSystemId,
      nextSortKey,
    });
    const qs = params.toString();
    const nextUrl = qs ? `/dashboard?${qs}` : '/dashboard';

    if (typeof window !== 'undefined') {
      const s = window.history.state;
      const next =
        s && typeof s === 'object'
          ? { ...s, as: nextUrl, url: nextUrl }
          : { as: nextUrl, url: nextUrl };
      window.history.replaceState(next, '', nextUrl);
    }
    router.replace(nextUrl, { scroll: false });
  };

  // ── 이벤트 핸들러 ────────────────────────────────────────
  const navigateToCourse = (course: CourseItem) => {
    router.push(course.destination || `/dashboard/course-drafts/${course.id}`);
  };

  const handleLumiCourseFocus = (courseId: string) => {
    const target = visibleCourses.find((c) => c.id === courseId);
    if (!target) return;

    const nextStatus =
      activeStatusFilter === 'all' || activeStatusFilter === target.status
        ? activeStatusFilter
        : target.status;
    const filtered = filterCourses(visibleCourses, nextStatus, query.trim());
    const sorted = sortCourses(filtered, sortKey);
    const idx = sorted.findIndex((c) => c.id === courseId);
    const pageSize = getPageSizeForMode(mode);
    const nextPage = idx >= 0 ? Math.floor(idx / pageSize) + 1 : 1;
    const systemId = findSystemIdForCourse(courseId, chunkedPagedCourses);

    dispatchDashboardMode({ type: 'setHoveredCourseId', payload: courseId });
    dispatchDashboardMode({ type: 'setSelectedCourseId', payload: courseId });
    dispatchDashboardMode({ type: 'setHoveredSystemId', payload: systemId });
    dispatchDashboardMode({ type: 'setSelectedSystemId', payload: systemId });
    dispatchDashboardMode({ type: 'setPage', payload: nextPage });
    setIsLumiPanelOpen(false);
    syncRouteState(nextPage, query, nextStatus, mode, systemId, sortKey);
    navigateToCourse(target);
  };

  const handleCourseHover = (courseId: string | null) => {
    const systemId = findSystemIdForCourse(courseId, chunkedPagedCourses);
    dispatchDashboardMode({ type: 'setHoveredCourseId', payload: courseId });
    dispatchDashboardMode({ type: 'setHoveredSystemId', payload: systemId });
    setHoverSource(courseId ? (mode === 'galaxy' ? 'course-in-galaxy' : 'course') : null);
  };

  const handleStatusFilterChange = (nextStatus: PlanetFilter) => {
    dispatchDashboardMode({ type: 'setHoveredCourseId', payload: null });
    dispatchDashboardMode({ type: 'setSelectedCourseId', payload: null });
    setSelectedMapPage(null);
    dispatchDashboardMode({ type: 'setStatusFilter', payload: nextStatus });
    dispatchDashboardMode({ type: 'setPage', payload: 1 });
    syncRouteState(1, query, nextStatus, mode, null, sortKey);
  };

  const handleSearchChange = (nextSearch: string) => {
    dispatchDashboardMode({ type: 'setSearchQuery', payload: nextSearch });
    dispatchDashboardMode({ type: 'setPage', payload: 1 });
    syncRouteState(1, nextSearch, activeStatusFilter, mode, selectedSystemId, sortKey);
  };

  const handlePageSelect = (nextPage: number) => {
    const p = Math.min(Math.max(nextPage, 1), Math.max(totalPages, 1));
    dispatchDashboardMode({ type: 'setPage', payload: p });
    setSelectedMapPage(p);
    syncRouteState(p, query, activeStatusFilter, mode, selectedSystemId, sortKey);
  };

  const handlePageHoverStart = (pageSystem: PageSystemChunk, clientX: number, clientY: number) => {
    if (isPhoneLayout) return;
    if (pageHoverTimerRef.current !== null) window.clearTimeout(pageHoverTimerRef.current);
    pageHoverTimerRef.current = window.setTimeout(() => {
      setPageHoverLumi({ pageNumber: pageSystem.pageNumber, x: clientX, y: clientY });
      pageHoverTimerRef.current = null;
      dispatchDashboardMode({ type: 'setHoveredSystemId', payload: pageSystem.id });
      setHoverSource('system');
    }, 260);
  };

  const handlePageHoverMove = (pageNumber: number, clientX: number, clientY: number) => {
    if (isPhoneLayout) return;
    setPageHoverLumi((cur) => {
      if (!cur || cur.pageNumber !== pageNumber) return cur;
      return { ...cur, x: clientX, y: clientY };
    });
  };

  const handlePageHoverEnd = () => {
    if (pageHoverTimerRef.current !== null) {
      window.clearTimeout(pageHoverTimerRef.current);
      pageHoverTimerRef.current = null;
    }
    setPageHoverLumi(null);
    dispatchDashboardMode({ type: 'setHoveredSystemId', payload: null });
    setHoverSource(null);
  };

  const handleSortChange = (nextSortKey: DashboardSortKey) => {
    if (nextSortKey === sortKey) return;
    dispatchDashboardMode({ type: 'setSortKey', payload: nextSortKey });
    dispatchDashboardMode({ type: 'setPage', payload: 1 });
    setSelectedMapPage(null);
    syncRouteState(1, query, activeStatusFilter, mode, selectedSystemId, nextSortKey);
  };

  // ── 시스템 네비게이션 ────────────────────────────────────
  const {
    prevSystemChunk,
    nextSystemChunk,
    handleSystemSelect,
    handleReturnToGalaxy,
    handlePrevSystem,
    handleNextSystem,
    handleModeToggle,
  } = useSystemNavigation({
    mode,
    query,
    activeStatusFilter,
    sortKey,
    selectedSystemChunk,
    selectedMapPage,
    pageSystems,
    totalPages,
    systemsPerGalaxyPage,
    pageHoverTimerRef,
    dispatchDashboardMode,
    setSelectedMapPage,
    setPageHoverLumi,
    setHoverSource,
    syncRouteState,
  });

  const handleLogout = async () => {
    await fetch('/api/v1/auth/logout', { method: 'POST', credentials: 'include' }).catch(
      () => undefined,
    );
    router.push('/');
  };

  // ── Lumi ─────────────────────────────────────────────────
  const dashboardLumiInput = useMemo(
    () =>
      buildDashboardLumiAdapter({
        dashboard: dashboardAdapter,
        userName: getDisplayName(user),
        actionHandlers: {
          newCourse: () => {
            setIsLumiPanelOpen(false);
            queryInputRef.current?.scrollIntoView({ behavior: 'smooth', block: 'center' });
            queryInputRef.current?.focus();
          },
          viewRelatedCourse: handleLumiCourseFocus,
          changeFilter: (status) => handleStatusFilterChange(status),
          clearFilter: () => handleStatusFilterChange('all'),
          openSettings: () => router.push('/dashboard/settings'),
        },
        hoverSource,
        hoveredSystemSummary,
      }),
    [dashboardAdapter, handleLumiCourseFocus, router, user, hoverSource, hoveredSystemSummary],
  );

  const lumiState = useLumiController({
    courses: dashboardLumiInput.courses,
    selectedCourseId: dashboardLumiInput.selectedCourseId,
    hoveredCourseId: dashboardLumiInput.hoveredCourseId,
    activeStatusFilter: dashboardLumiInput.activeStatusFilter,
    userName: dashboardLumiInput.userName,
    hasByok: dashboardLumiInput.hasByok,
    isPhoneLayout: dashboardLumiInput.isPhoneLayout,
    isExpanded: dashboardLumiInput.isExpanded,
    actionHandlers: { ...dashboardLumiInput.actionHandlers },
    hoverSource,
    hoveredSystemSummary,
  });

  // ── 파생 표시값 ──────────────────────────────────────────
  const shouldHideSpaceMap = isPhoneLayout;
  const showListOnlyFallback = shouldHideSpaceMap;
  const starSystemIndexLabel = selectedSystem?.index ?? selectedSystemChunk?.index ?? '-';
  const displayCourses =
    mode === 'star-system' && selectedSystemChunk ? selectedSystemChunk.courses : pagedCourses;
  const listModeLabel =
    mode === 'star-system' ? dashboardCopy.list.modeLabel.starSystem(starSystemIndexLabel) : dashboardCopy.list.modeLabel.galaxy;
  const listPageLabel =
    mode === 'star-system'
      ? dashboardCopy.list.pageLabel.starSystem(starSystemIndexLabel)
      : dashboardCopy.list.pageLabel.galaxy(page, Math.max(totalPages, 1));
  const navEnabled = mode === 'galaxy' && totalPages > 1;
  const isCompletedOnly = visibleCourses.length > 0 && visibleCourses.every((course) => course.status === 'completed');

  // ── 로딩 ────────────────────────────────────────────────
  if (isLoading) {
    return (
      <div style={loadingPageStyle}>
        <div style={loadingCardStyle}>{dashboardCopy.loading}</div>
      </div>
    );
  }

  // ── 렌더 ─────────────────────────────────────────────────
  return (
    <div style={pageStyle}>
      {isGoalSubmitting && goalFlowStage ? (
        <GoalFlowLoadingOverlay stage={goalFlowStage} />
      ) : null}
      <DashboardBackground />
      <div style={pageChromeStyle}>
        <LearnerAppShell
          locale={locale}
          logoHref="/"
          rightSlot={<LearnerHeaderActions onLogout={handleLogout} copy={learnerHeaderCopy} />}
          contentClassName="!px-0 !pb-0"
          showSectionSubNav={!isPhoneLayout}
        >
        {isPhoneLayout ? (
          <CosmosMobileFirstView
            courses={visibleCourses}
            todayTask={todayTask}
            userName={getDisplayName(user)}
            isSubmitting={isGoalSubmitting}
            creationError={goalCreationError}
            onCreateStar={handleGoalSubmit}
            onOpenCourse={navigateToCourse}
            onOpenTask={(href) => router.push(href)}
          />
        ) : (
        <main style={{ ...mainStyle, padding: '20px 24px 40px' }}>
          <section style={headerBlockStyle}>
            <p style={heroSubtitleStyle}>{dashboardCopy.heroSubtitle}</p>
          </section>


          <DashboardTodayTaskSection
            task={todayTask}
            copy={dashboardCopy.todayTask}
            onOpenTask={(href) => router.push(href)}
          />

          <section style={mapListGridStyle(isPhoneLayout, isShortViewport, shouldHideSpaceMap)}>
            {!shouldHideSpaceMap ? (
              <div style={{ ...mapColumnStyle, position: 'relative' }}>
                <GalaxyMap
                  pageSystems={pageSystems}
                  mode={mode}
                  selectedSystemChunk={selectedSystemChunk}
                  selectedMapPage={selectedMapPage}
                  hoveredSystemId={hoveredSystemId}
                  hoveredCourseId={hoveredCourseId}
                  selectedCourseId={selectedCourseId}

                  pageHoverLumi={pageHoverLumi}
                  hoveredPageSystem={hoveredPageSystem}
                  selectedSystemTitleFallback={selectedSystemTitleFallback}
                  starSystemStatusLine={starSystemStatusLine}
                  starSystemUpdatedLabel={starSystemUpdatedLabel}
                  sortKey={sortKey}
                  filteredCourseCount={filteredCourses.length}
                  activeStatusFilter={activeStatusFilter}
                  isPhoneLayout={isPhoneLayout}
                  isCompactLayout={isCompactLayout}
                  isShortViewport={isShortViewport}
                  showListOnlyFallback={showListOnlyFallback}
                  onSystemSelect={handleSystemSelect}
                  onPageHoverStart={handlePageHoverStart}
                  onPageHoverMove={handlePageHoverMove}
                  onPageHoverEnd={handlePageHoverEnd}
                  onHoverCourse={handleCourseHover}
                  onSelectCourse={handleLumiCourseFocus}
                  onReturnToGalaxy={handleReturnToGalaxy}
                  onSortChange={handleSortChange}
                  locale={locale}
                  copy={dashboardCopy.galaxy}
                  planetMapCopy={dashboardCopy.planetMap}
                />
              </div>
            ) : null}

            <div style={responsiveListColumnStyle(shouldHideSpaceMap)}>
              <CourseListPanel
                courses={displayCourses}
                mode={mode}
                isCompactLayout={isCompactLayout}
                isListOnly={shouldHideSpaceMap}
                page={page}
                totalPages={totalPages}
                hoveredCourseId={hoveredCourseId}
                selectedCourseId={selectedCourseId}
                selectedSystemChunk={selectedSystemChunk}
                searchQuery={query}
                listModeLabel={listModeLabel}
                listPageLabel={listPageLabel}
                showInactiveCourses={showInactiveCourses}
                inactiveCourseCount={inactiveCourseCount}
                navEnabled={navEnabled}
                prevPageDisabled={!navEnabled || page <= 1}
                nextPageDisabled={!navEnabled || page >= totalPages}
                activeStatusFilter={activeStatusFilter}
                copy={dashboardCopy.list}
                galaxyCopy={dashboardCopy.galaxy}
                locale={locale}
                onStatusFilterChange={handleStatusFilterChange}
                onSearchChange={handleSearchChange}
                onCourseHover={handleCourseHover}
                onCourseClick={handleLumiCourseFocus}
                onPagePrev={() => handlePageSelect(page - 1)}
                onPageNext={() => handlePageSelect(page + 1)}
                onToggleInactiveCourses={() => {
                  setShowInactiveCourses((current) => !current);
                  setSelectedMapPage(1);
                  setPageHoverLumi(null);
                  setHoverSource(null);
                  dispatchDashboardMode({ type: 'setMode', payload: 'galaxy' });
                  dispatchDashboardMode({ type: 'setHoveredCourseId', payload: null });
                  dispatchDashboardMode({ type: 'setSelectedCourseId', payload: null });
                  dispatchDashboardMode({ type: 'setHoveredSystemId', payload: null });
                  dispatchDashboardMode({ type: 'setSelectedSystemId', payload: null });
                  dispatchDashboardMode({ type: 'setPage', payload: 1 });
                  syncRouteState(1, query, activeStatusFilter, 'galaxy', null, sortKey);
                }}
                onModeToggle={handleModeToggle}
                onPrevSystem={handlePrevSystem}
                onNextSystem={handleNextSystem}
                prevSystemDisabled={!prevSystemChunk}
                nextSystemDisabled={!nextSystemChunk}
              />
            </div>
          </section>
        </main>
        )}
        </LearnerAppShell>
      </div>
    </div>
  );
}

// ─────────────────────────────────────────────────────────────
// Page export
// ─────────────────────────────────────────────────────────────
export default function DashboardPage() {
  return (
    <Suspense>
      <DashboardContent />
    </Suspense>
  );
}
