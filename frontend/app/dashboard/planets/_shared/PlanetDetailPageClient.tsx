'use client';

import Link from 'next/link';
import { useParams, useRouter } from 'next/navigation';
import { useEffect, useMemo, useRef, useState } from 'react';
import AppHeaderShell from '@/components/common/AppHeaderShell';
import LearnerHeaderActions, { getLearnerHeaderActionsCopy } from '@/components/common/LearnerHeaderActions';
import ImmersiveOverlayMenu from '@/components/navigation/ImmersiveOverlayMenu';
import LumiModalShell, {
  lumiModalPrimaryButtonStyle,
  lumiModalSecondaryButtonStyle,
} from '@/components/common/LumiModalShell';
import { getPlanetDiaryCopy } from '@/lib/i18n/pages/planetDiary';
import { PlanetDiaryStage } from './diary/PlanetDiaryStage';
import {
  DiaryContextBand,
  LearningDiaryProgressAction,
  SharedDiaryNoticeCard,
} from './PlanetDetailSummary';
import {
  type PlanetRouteKind,
  type DiaryStageTab,
  type UserInfo,
  type PlanetAggregate,
  type ExplorerCourseAggregate,
  type DiaryContextDraft,
  buildPlanetLevelsFromLessons,
  collectPlanetPointStatuses,
  countProgressFromExplorerCourse,
  countPlanetPointsFromLevels,
  countPlanetPointsFromLessonTrees,
  countLessonProgressFromLessonTrees,
  countLessonProgressFromLevels,
  buildDiaryCourseAggregate,
} from './planetDetailUtils';

const explorerDiaryBackdropSrc = '/images/explorer/Explorer_Diary_Backdrop.webp';

export default function PlanetDetailPageClient({ routeKind }: { routeKind: PlanetRouteKind }) {
  const router = useRouter();
  const params = useParams<{ id: string }>();
  const [planet, setPlanet] = useState<PlanetAggregate | null>(null);
  const [draftContext, setDraftContext] = useState<DiaryContextDraft | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isStartingPlanet, setIsStartingPlanet] = useState(false);
  const [startMessage, setStartMessage] = useState<string | null>(null);
  const [showCompleteConfirm, setShowCompleteConfirm] = useState(false);
  const [isCompletingPlanet, setIsCompletingPlanet] = useState(false);
  const [completeMessage, setCompleteMessage] = useState<string | null>(null);
  const [showPlanningConfirm, setShowPlanningConfirm] = useState(false);
  const [activeDiaryTab, setActiveDiaryTab] = useState<DiaryStageTab>('journal');
  const [explorerCourse, setExplorerCourse] = useState<ExplorerCourseAggregate | null>(null);
  const [uiLocale, setUiLocale] = useState<'ko' | 'en'>('ko');
  const promptedCompleteKeyRef = useRef<string | null>(null);
  const copy = useMemo(() => getPlanetDiaryCopy(uiLocale), [uiLocale]);

  const displayLevels = useMemo(
    () => (planet?.lessons?.length ? buildPlanetLevelsFromLessons(planet.lessons) : planet?.levels ?? []),
    [planet],
  );
  const pointProgress = useMemo(
    () => (planet?.lessons?.length ? countPlanetPointsFromLessonTrees(planet.lessons) : countPlanetPointsFromLevels(displayLevels)),
    [displayLevels, planet?.lessons],
  );
  const lessonProgress = useMemo(
    () => (planet?.lessons?.length ? countLessonProgressFromLessonTrees(planet.lessons) : countLessonProgressFromLevels(displayLevels)),
    [displayLevels, planet?.lessons],
  );
  const explorerProgress = useMemo(() => {
    if (!planet || !explorerCourse) return null;
    return countProgressFromExplorerCourse(explorerCourse, collectPlanetPointStatuses(planet));
  }, [explorerCourse, planet]);
  const backendLessonProgress = useMemo(() => {
    if (!planet) return null;
    const total = Math.max(0, planet.planet.lesson_count ?? 0);
    const completed = Math.max(0, Math.min(total, planet.planet.completed_lesson_count ?? 0));
    return {
      total,
      completed,
      learning: explorerProgress?.lessonProgress.learning ?? lessonProgress.learning,
    };
  }, [explorerProgress?.lessonProgress.learning, lessonProgress.learning, planet]);
  const activePointProgress = explorerProgress?.pointProgress ?? pointProgress;
  const activeLessonProgress = backendLessonProgress ?? explorerProgress?.lessonProgress ?? lessonProgress;
  const currentPlanetID = typeof params.id === 'string' ? params.id : planet?.planet.id ?? '';
  const diaryCourse = useMemo(() => (planet ? buildDiaryCourseAggregate(planet) : null), [planet]);
  const canCompletePlanet = useMemo(() => {
    if (planet?.planet.can_complete != null) return Boolean(planet.planet.can_complete);
    return activeLessonProgress.total > 0 && activeLessonProgress.completed === activeLessonProgress.total;
  }, [activeLessonProgress.completed, activeLessonProgress.total, planet?.planet.can_complete]);
  const completionPromptKey = useMemo(() => {
    if (!planet) return null;
    return [
      planet.planet.id,
      planet.planet.status,
      activeLessonProgress.completed,
      activeLessonProgress.total,
      activePointProgress.completed,
      activePointProgress.total,
    ].join(':');
  }, [
    activeLessonProgress.completed,
    activeLessonProgress.total,
    activePointProgress.completed,
    activePointProgress.total,
    planet,
  ]);

  const redirectPath = useMemo(() => {
    const planetID = typeof params.id === 'string' ? params.id : '';
    return planetID ? `/dashboard/planets/${routeKind}/${planetID}` : `/dashboard/planets/${routeKind}`;
  }, [params.id, routeKind]);

  const handleLogout = async () => {
    await fetch('/api/v1/auth/logout', { method: 'POST', credentials: 'include' }).catch(() => undefined);
    router.replace('/login');
  };

  const loadPlanetPageData = async (planetID: string) => {
    const refreshRes = await fetch('/api/v1/auth/refresh', { method: 'POST', credentials: 'include' });
    if (!refreshRes.ok) {
      router.push(`/login?redirect_after=${redirectPath}`);
      return;
    }

    const meRes = await fetch('/api/v1/auth/me', { credentials: 'include', cache: 'no-store' });
    if (!meRes.ok) {
      router.push(`/login?redirect_after=${redirectPath}`);
      return;
    }

    const meData = (await meRes.json()) as UserInfo;
    setUiLocale(meData.ui_locale === 'en' ? 'en' : 'ko');
    const requestCopy = getPlanetDiaryCopy(meData.ui_locale);
    if (meData.required_consent_pending) {
      router.replace(`/agreements?redirect_after=${redirectPath}`);
      return;
    }

    const planetRes = await fetch(`/api/v1/planets/${routeKind}/${planetID}`, { credentials: 'include', cache: 'no-store' });

    if (planetRes.status === 403) {
      throw new Error(requestCopy.errors.inactiveDiaryOnlyPlanning);
    }

    if (!planetRes.ok) {
      throw new Error(requestCopy.errors.loadFailed);
    }

    const payload = (await planetRes.json()) as { planet: PlanetAggregate };
    setPlanet(payload.planet);

    if (payload.planet.planet.draft_id) {
      const explorerRes = await fetch(`/api/v1/explorer/course/${payload.planet.planet.draft_id}?include_inactive=true`, {
        credentials: 'include',
        cache: 'no-store',
      });
      if (explorerRes.ok) {
        const explorerPayload = (await explorerRes.json()) as { course?: ExplorerCourseAggregate } & ExplorerCourseAggregate;
        setExplorerCourse(explorerPayload.course ?? explorerPayload);
      } else {
        setExplorerCourse(null);
      }

      const draftRes = await fetch(`/api/v1/course-drafts/${payload.planet.planet.draft_id}`, {
        credentials: 'include',
        cache: 'no-store',
      });
      if (draftRes.ok) {
        const draftPayload = (await draftRes.json()) as { draft: DiaryContextDraft };
        setDraftContext(draftPayload.draft);
      }
    }

    setIsLoading(false);
  };

  const handleStartPlanet = async () => {
    if (!planet || routeKind !== 'learning' || planet.planet.status !== 'ready' || isStartingPlanet) return;

    setIsStartingPlanet(true);
    setStartMessage(null);

    try {
      const response = await fetch(`/api/v1/planets/learning/${planet.planet.id}/start`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
      });
      const payload = (await response.json().catch(() => ({}))) as { planet?: PlanetAggregate; error?: string };

      if (!response.ok || !payload.planet) {
        if (response.status === 400) {
          throw new Error(copy.errors.startNotReady);
        }
        if (response.status === 401) {
          throw new Error(copy.errors.loginRequired);
        }
        throw new Error(payload.error ?? copy.errors.startFailed);
      }

      setPlanet(payload.planet);
      setStartMessage(copy.errors.startSuccess);
    } catch (startError) {
      setStartMessage(
        startError instanceof Error
          ? startError.message
          : copy.errors.startFallback,
      );
    } finally {
      setIsStartingPlanet(false);
    }
  };

  const handleCompletePlanet = async () => {
    if (!planet || routeKind !== 'learning' || planet.planet.status !== 'learning' || isCompletingPlanet || !canCompletePlanet) return;

    setIsCompletingPlanet(true);
    setCompleteMessage(null);

    try {
      const response = await fetch(`/api/v1/planets/learning/${planet.planet.id}/complete`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
      });
      const payload = (await response.json().catch(() => ({}))) as { planet?: PlanetAggregate; error?: string };

      if (!response.ok || !payload.planet) {
        if (response.status === 400) {
          throw new Error(copy.errors.completeNeedsAllPoints);
        }
        if (response.status === 401) {
          throw new Error(copy.errors.loginRequired);
        }
        throw new Error(payload.error ?? copy.errors.completeFailed);
      }

      router.push(`/dashboard/planets/shared/${payload.planet.planet.id}`);
    } catch (completeError) {
      setCompleteMessage(
        completeError instanceof Error
          ? completeError.message
          : copy.errors.completeFailed,
      );
      setShowCompleteConfirm(false);
    } finally {
      setIsCompletingPlanet(false);
    }
  };

  useEffect(() => {
    const planetID = typeof params.id === 'string' ? params.id : '';
    if (!planetID) {
      setError(copy.errors.invalidPlanetPath);
      setIsLoading(false);
      return;
    }

    const load = async () => {
      await loadPlanetPageData(planetID);
    };

    load().catch((loadError) => {
      setError(loadError instanceof Error ? loadError.message : copy.errors.loadFailed);
      setIsLoading(false);
    });
  }, [params.id, redirectPath, routeKind, router]);

  useEffect(() => {
    const shouldPrompt =
      routeKind === 'learning' &&
      planet?.planet.status === 'learning' &&
      canCompletePlanet &&
      completionPromptKey != null;

    if (!shouldPrompt) {
      if (!canCompletePlanet) {
        promptedCompleteKeyRef.current = null;
        setShowCompleteConfirm(false);
      }
      return;
    }

    if (promptedCompleteKeyRef.current === completionPromptKey) return;
    promptedCompleteKeyRef.current = completionPromptKey;
    setShowCompleteConfirm(true);
  }, [canCompletePlanet, completionPromptKey, planet?.planet.status, routeKind]);

  const currentDraftID = planet?.planet.draft_id;
  const planningEditHref =
    routeKind === 'learning' && currentDraftID
      ? `/dashboard/course-drafts/${currentDraftID}?from=diary&section=planning`
      : null;
  const sectionTabLabel =
    activeDiaryTab === 'records' ? copy.section.tabs.records : activeDiaryTab === 'results' ? copy.section.tabs.results : copy.section.tabs.journal;

  if (isLoading) {
    return (
      <div style={loadingPageStyle}>
        <div style={loadingCardStyle}>{copy.loading}</div>
      </div>
    );
  }

  if (error || !planet) {
    return (
      <div style={loadingPageStyle}>
        <div style={{ ...loadingCardStyle, display: 'grid', gap: '14px' }}>
          <div>{error ?? copy.notFound}</div>
          <Link href="/dashboard" style={secondaryButtonStyle}>
            {copy.backToGalaxy}
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div
      style={{
        ...pageStyle,
        backgroundImage: `url(${explorerDiaryBackdropSrc})`,
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backgroundRepeat: 'no-repeat',
      }}
    >
      <AppHeaderShell
        logoHref="/"
        logoIconSize={28}
        logoTextSize="16px"
        maxWidth="1360px"
        headerStyle={navStyle}
        innerStyle={{ padding: '14px 18px' }}
        leftGroupStyle={{ flex: '1 1 auto', minWidth: 0 }}
        leftMeta={
          <span
            className="lw-planet-detail-header-meta"
            style={{
              display: 'inline-flex',
              alignItems: 'baseline',
              gap: '8px',
              flexWrap: 'wrap',
              minWidth: 0,
            }}
          >
            <span style={{ fontSize: '16px', fontWeight: 800, letterSpacing: '0.01em', color: '#F4F7FF' }}>
              {copy.header.meta}
            </span>
            <span style={{ fontSize: '13px', fontWeight: 500, color: 'rgba(200,210,235,0.56)' }}>
              {routeKind === 'learning' ? copy.header.learning : copy.header.shared}
            </span>
          </span>
        }
        rightSlot={
          <span style={{ display: 'inline-flex', alignItems: 'center', gap: '10px' }}>
            <ImmersiveOverlayMenu locale={uiLocale} />
            <LearnerHeaderActions onLogout={handleLogout} copy={getLearnerHeaderActionsCopy(uiLocale)} />
          </span>
        }
      />

      <main style={mainStyle}>
        <section
          id="planet-route"
          style={{
            ...sectionStyle,
            ...(routeKind === 'learning'
              ? {
                  backgroundImage: `url(/images/dashboard/explorer-diary-disk.webp), linear-gradient(rgba(7, 18, 34, 0.42), rgba(7, 18, 34, 0.58)), url(${explorerDiaryBackdropSrc})`,
                  backgroundSize: 'cover, cover, cover',
                  backgroundPosition: 'center top, center, center',
                  backgroundRepeat: 'no-repeat, no-repeat, no-repeat',
                }
              : {}),
          }}
        >
          <div style={sectionHeaderStyle}>
            <div>
              <h2 style={sectionTitleStyle}>
                <span style={sectionTitlePrefixStyle}>{copy.section.planetPrefix}</span>
                <span style={sectionTitleCourseNameStyle}>{planet.planet.title}</span>
                <span style={sectionTitlePrefixStyle}>
                  {routeKind === 'learning' ? ` ${sectionTabLabel}` : `${copy.section.sharedPrefix}${sectionTabLabel}`}
                </span>
              </h2>
            </div>
            {planningEditHref ? (
              <button
                type="button"
                style={planningEditLinkStyle}
                onClick={() => setShowPlanningConfirm(true)}
              >
                {copy.section.editPlan}
              </button>
            ) : null}
          </div>
          <div style={{ display: 'grid', gap: '14px', padding: '2px 0 6px' }}>
            <div style={explorerInfoBoardStyle}>
              <DiaryContextBand
                routeKind={routeKind}
                planet={planet}
                draftContext={draftContext}
                pointProgress={activePointProgress}
                lessonProgress={activeLessonProgress}
                copy={copy.summary}
                actionSlot={
                  routeKind === 'learning' ? (
                    <LearningDiaryProgressAction
                      planetStatus={planet.planet.status}
                      canCompletePlanet={canCompletePlanet}
                      isStartingPlanet={isStartingPlanet}
                      isCompletingPlanet={isCompletingPlanet}
                      onStart={handleStartPlanet}
                      onOpenCompleteConfirm={() => {
                        if (!canCompletePlanet) return;
                        setShowCompleteConfirm(true);
                      }}
                      copy={copy.summary.progressAction}
                    />
                  ) : undefined
                }
              />
              {routeKind === 'learning' ? (
                null
              ) : (
                <SharedDiaryNoticeCard copy={copy.summary.sharedNotice} />
              )}

              {startMessage ? <div style={startNoticeStyle}>{startMessage}</div> : null}
              {completeMessage ? <div style={startNoticeStyle}>{completeMessage}</div> : null}
            </div>
          </div>
          {diaryCourse ? (
            <PlanetDiaryStage
              course={diaryCourse}
              planetId={currentPlanetID || planet.planet.id}
              routeKind={routeKind}
              activeTab={activeDiaryTab}
              canCompleteCourse={routeKind === 'learning' && planet.planet.status === 'learning' && canCompletePlanet}
              isCompletingCourse={isCompletingPlanet}
              onOpenCourseCompleteConfirm={() => {
                if (!canCompletePlanet || planet.planet.status !== 'learning') return;
                setShowCompleteConfirm(true);
              }}
              onActiveTabChange={setActiveDiaryTab}
              onRefresh={() => loadPlanetPageData(currentPlanetID || planet.planet.id)}
              copy={copy}
            />
          ) : null}
          {routeKind === 'learning' && planet.planet.status === 'learning' && showCompleteConfirm ? (
            <LumiModalShell
              ariaLabel={copy.completionModal.aria}
              eyebrow="Course Complete"
              lumiState="celebrate"
              tone="warm"
              title={copy.completionModal.title}
              message={
                <div style={completeModalContentStyle}>
                  <section style={completeModalSectionStyle}>
                    <div style={completeModalSectionTitleStyle}>{copy.completionModal.guideTitle}</div>
                    <p style={completeModalTextStyle}>{copy.completionModal.guideMessage(activeLessonProgress.completed, activeLessonProgress.total)}</p>
                  </section>
                  <section style={completeModalSectionStyle}>
                    <div style={completeModalSectionTitleStyle}>{copy.completionModal.shareTitle}</div>
                    <p style={completeModalTextStyle}>{copy.completionModal.shareMessage}</p>
                  </section>
                  <div style={completeModalActionDividerStyle}>{copy.completionModal.actionPrompt}</div>
                </div>
              }
              onClose={() => {
                if (!isCompletingPlanet) setShowCompleteConfirm(false);
              }}
              actions={
                <>
                  <button
                    type="button"
                    style={lumiModalSecondaryButtonStyle}
                    disabled={isCompletingPlanet}
                    onClick={() => setShowCompleteConfirm(false)}
                  >
                    {copy.completionModal.cancel}
                  </button>
                  <button
                    type="button"
                    style={lumiModalPrimaryButtonStyle}
                    disabled={isCompletingPlanet || !canCompletePlanet}
                    onClick={handleCompletePlanet}
                  >
                    {isCompletingPlanet ? copy.completionModal.busy : copy.completionModal.confirm}
                  </button>
                </>
              }
            />
          ) : null}
          {planningEditHref && showPlanningConfirm ? (
            <LumiModalShell
              ariaLabel={copy.planningConfirm.aria}
              eyebrow="Lumi Confirm"
              lumiState="curious"
              tone="warm"
              title={copy.planningConfirm.title}
              message={copy.planningConfirm.message}
              onClose={() => setShowPlanningConfirm(false)}
              actions={
                <>
                  <button
                    type="button"
                    style={lumiModalSecondaryButtonStyle}
                    onClick={() => setShowPlanningConfirm(false)}
                  >
                    {copy.planningConfirm.cancel}
                  </button>
                  <button
                    type="button"
                    style={lumiModalPrimaryButtonStyle}
                    onClick={() => window.location.assign(planningEditHref)}
                  >
                    {copy.planningConfirm.confirm}
                  </button>
                </>
              }
            />
          ) : null}
        </section>
      </main>
      <style jsx global>{`
        @media (max-width: 640px) {
          .lw-planet-detail-header-meta {
            display: none !important;
          }
        }
      `}</style>
    </div>
  );
}

const pageStyle = {
  position: 'relative',
  minHeight: '100vh',
  overflow: 'hidden',
} as const;

const completeModalContentStyle = {
  display: 'grid',
  gap: '12px',
} as const;

const completeModalSectionStyle = {
  display: 'grid',
  gap: '6px',
  padding: '12px 14px',
  borderRadius: '14px',
  background: 'rgba(255, 250, 239, 0.56)',
  border: '1px solid rgba(148, 94, 22, 0.16)',
} as const;

const completeModalSectionTitleStyle = {
  fontSize: '12px',
  fontWeight: 900,
  letterSpacing: '0.04em',
  color: '#78480F',
} as const;

const completeModalTextStyle = {
  margin: 0,
  color: 'rgba(74, 44, 10, 0.92)',
  fontSize: '14px',
  lineHeight: 1.68,
} as const;

const completeModalActionDividerStyle = {
  paddingTop: '10px',
  borderTop: '1px solid rgba(148, 94, 22, 0.20)',
  color: 'rgba(91, 55, 12, 0.78)',
  fontSize: '12px',
  fontWeight: 800,
} as const;

const loadingPageStyle = {
  minHeight: '100vh',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  background: 'linear-gradient(rgba(7, 14, 26, 0.76), rgba(7, 14, 26, 0.88))',
  padding: '24px',
} as const;

const loadingCardStyle = {
  width: '100%',
  maxWidth: '540px',
  padding: '32px',
  borderRadius: '24px',
  background: 'rgba(8, 18, 33, 0.76)',
  border: '1px solid rgba(194, 210, 245, 0.16)',
  textAlign: 'center',
  color: '#F4F7FF',
  fontSize: '16px',
  backdropFilter: 'blur(10px)',
} as const;

const planningEditLinkStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  borderRadius: '999px',
  padding: '10px 16px',
  background: 'rgba(255,255,255,0.08)',
  border: '1px solid rgba(194, 210, 245, 0.18)',
  color: '#F4F7FF',
  textDecoration: 'none',
  fontSize: '13px',
  fontWeight: 700,
  whiteSpace: 'nowrap' as const,
  cursor: 'pointer',
} as const;

const navStyle = {
  position: 'relative',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '16px',
  padding: '16px 24px',
  borderBottom: '1px solid rgba(180, 205, 255, 0.12)',
  background: 'rgba(6, 15, 29, 0.58)',
  backdropFilter: 'blur(14px)',
} as const;

const mainStyle = {
  position: 'relative',
  maxWidth: '1360px',
  margin: '0 auto',
  padding: '32px 24px 48px',
  display: 'grid',
  gap: '18px',
} as const;

const sectionStyle = {
  display: 'grid',
  gap: '16px',
  borderRadius: '28px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(7, 18, 34, 0.56)',
  backdropFilter: 'blur(14px)',
  padding: '24px',
} as const;

const sectionHeaderStyle = {
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'flex-start',
  gap: '12px',
  flexWrap: 'wrap',
} as const;

const sectionTitleStyle = {
  margin: 0,
  fontSize: '28px',
  color: '#F7FAFF',
} as const;

const sectionTitlePrefixStyle = {
  color: '#F7FAFF',
  fontWeight: 700,
} as const;

const sectionTitleCourseNameStyle = {
  color: '#FFF1B8',
  fontWeight: 900,
  textShadow: '0 2px 12px rgba(145, 92, 17, 0.28)',
} as const;

const explorerInfoBoardStyle = {
  display: 'grid',
  gap: '14px',
  borderRadius: '24px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(7, 18, 34, 0.48)',
  backdropFilter: 'blur(14px)',
  padding: '18px 18px 20px',
} as const;

const secondaryButtonStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '42px',
  padding: '0 16px',
  borderRadius: '999px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(255,255,255,0.08)',
  color: '#F4F7FF',
  textDecoration: 'none',
  fontSize: '14px',
  fontWeight: 700,
} as const;

const startNoticeStyle = {
  marginTop: '18px',
  padding: '14px 16px',
  borderRadius: '16px',
  background: 'rgba(8, 20, 36, 0.78)',
  border: '1px solid rgba(141, 198, 255, 0.16)',
  color: '#DCEEFF',
  fontSize: '14px',
  lineHeight: 1.6,
} as const;
