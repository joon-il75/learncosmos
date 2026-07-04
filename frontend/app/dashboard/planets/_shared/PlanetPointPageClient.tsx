'use client';

import Link from 'next/link';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';

import type { PlanetRouteKind } from './pointPageTypes';
import {
  type PointPageTheme,
  getPlanetDetailHref,
  getPointPageThemeTokens,
} from './pointPageUtils';
import { usePointLearning } from './point/usePointLearning';
import PointLearningFlow from './point/PointLearningFlow';
import { getPointStageHref, pointLearningStageFromSlug, pointLearningStageToSlug, type PointLearningStage } from './point/skillFlowTypes';
import { PointLearningToolbar } from './point/PointLearningToolbar';
import { PointPageHeader } from './point/PointPageHeader';
import { getLearnerHeaderActionsCopy } from '@/components/common/LearnerHeaderActions';
import { getPointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import {
  type LearningWorkGuideDetail,
  getPointGoalContextText,
  pointToolbarOffsets,
} from './point/pointGuideCopy';
import {
  pageStyle, pageBackgroundOverlayStyle,
  loadingPageStyle, loadingCardStyle,
  mainStyle,
  sectionStyle,
  secondaryButtonStyle,
} from './pointPageStyles';

const POINT_HEADER_FIXED_TOP = 100;
const POINT_TOOLBAR_FLOW_GAP = 12;

export default function PlanetPointPageClient({ routeKind }: { routeKind: PlanetRouteKind }) {
  const router = useRouter();
  const params = useParams<{ stage?: string }>();
  const [theme, setTheme] = useState<PointPageTheme>('light');
  const [isLearningContentOpen, setIsLearningContentOpen] = useState(true);
  const [isLearningWorkOpen, setIsLearningWorkOpen] = useState(false);
  const [, setLearningWorkGuideDetail] = useState<LearningWorkGuideDetail>({
    title: 'Learning Flow - 학습하기',
    message: '먼저 내용정리에서 필요한 항목 하나를 저장해 보세요. 지점을 완료하려면 질문 답변 1개와 자기평가가 필요하고, 연습/활동기록과 결과물은 필요할 때 남기면 됩니다.',
  });
  const previousResearchMaterialConfirmedRef = useRef<boolean | null>(null);
  const [pointToolbarHeight, setPointToolbarHeight] = useState(() => Number.parseInt(pointToolbarOffsets.collapsed, 10) - POINT_HEADER_FIXED_TOP - POINT_TOOLBAR_FLOW_GAP);
  const autoOpenedConfirmedResearchPointRef = useRef<string | null>(null);
  const autoOpenedSavedJournalPointRef = useRef<string | null>(null);
  const themeTokens = useMemo(() => getPointPageThemeTokens(theme), [theme]);
  const routeStage = pointLearningStageFromSlug(typeof params.stage === 'string' ? params.stage : null);

  const learning = usePointLearning(routeKind);
  const pointCopy = useMemo(() => getPointLearningCopy(learning.uiLocale), [learning.uiLocale]);
  const learnerHeaderCopy = useMemo(() => getLearnerHeaderActionsCopy(learning.uiLocale), [learning.uiLocale]);
  const {
    planet, pointDetail, error, isLoading, planetID,
    hasSavedJournalNote,
    pointQuestions, practiceLogDrafts, artifactDrafts,
  } = learning;
  const pointToolbarOffset = `${POINT_HEADER_FIXED_TOP + pointToolbarHeight + POINT_TOOLBAR_FLOW_GAP}px`;
  const pointBottomGuideOffset = '56px';
  const hasSelfEvaluationEntryRecord =
    learning.answeredQuestionCount > 0 ||
    practiceLogDrafts.length > 0 ||
    artifactDrafts.length > 0;
  const canOpenSelfEvaluation = hasSelfEvaluationEntryRecord || Boolean(pointDetail?.self_evaluation) || pointDetail?.point.status === 'completed';
  const isExplorationPoint = pointDetail?.point.point_type === 'exploration';
  const isContentViewChecked = pointDetail ? (isExplorationPoint ? true : learning.isResearchMaterialConfirmed) : false;
  const canOpenLearningWork = isContentViewChecked;
  const getStageHref = useCallback((stage: PointLearningStage) => getPointStageHref(routeKind, planetID, learning.pointID, stage), [learning.pointID, planetID, routeKind]);
  const handleToggleTheme = () => setTheme((cur: PointPageTheme) => (cur === 'dark' ? 'light' : 'dark'));
  const handlePointToolbarHeightChange = useCallback((height: number) => {
    setPointToolbarHeight((current) => (current === height ? current : height));
  }, []);
  const scrollToPointSection = useCallback((sectionID: string) => {
    window.requestAnimationFrame(() => {
      window.requestAnimationFrame(() => {
        const section = document.getElementById(sectionID);
        if (!section) return;
        const offset = Number.parseInt(pointToolbarOffset, 10) + 16;
        const targetTop = section.getBoundingClientRect().top + window.scrollY - offset;
        window.scrollTo({ top: Math.max(0, targetTop), behavior: 'smooth' });
      });
    });
  }, [pointToolbarOffset]);
  const handleResearchMaterialConfirmed = useCallback(() => {
    setIsLearningContentOpen(false);
    setIsLearningWorkOpen(true);
    learning.setIsCompletionOpen(false);
    scrollToPointSection('point-work-section');
  }, [learning, scrollToPointSection]);
  const handleSourceContentOpen = useCallback(() => {
    learning.setActiveWorkTab('notes');
    setIsLearningContentOpen(false);
    setIsLearningWorkOpen(true);
    learning.setIsCompletionOpen(false);
    setLearningWorkGuideDetail(pointCopy.workspace.learningWork.sourceContentGuide);
    scrollToPointSection('point-work-section');
  }, [learning, pointCopy.workspace.learningWork.sourceContentGuide, scrollToPointSection]);
  const handleLogout = async () => {
    await fetch('/api/v1/auth/logout', { method: 'POST', credentials: 'include' }).catch(() => undefined);
    router.replace('/login');
  };

  const handlePointLearningStageChange = useCallback((stage: PointLearningStage) => {
    if (stage === 'point' || stage === 'goal' || stage === 'skill_discovery') {
      setIsLearningContentOpen(stage === 'point');
      setIsLearningWorkOpen(false);
      learning.setIsCompletionOpen(false);
      return;
    }
    if (stage === 'goal_iteration') {
      setIsLearningContentOpen(false);
      setIsLearningWorkOpen(true);
      learning.setIsCompletionOpen(false);
      return;
    }
    if (stage === 'complete') {
      setIsLearningContentOpen(false);
      setIsLearningWorkOpen(false);
      learning.setIsCompletionOpen(true);
      return;
    }
    setIsLearningContentOpen(false);
    setIsLearningWorkOpen(false);
    learning.setIsCompletionOpen(false);
  }, [learning]);

  useEffect(() => {
    if (!pointDetail || pointDetail.point.point_type === 'exploration') return;
    const wasConfirmed = previousResearchMaterialConfirmedRef.current;
    const isConfirmed = learning.isResearchMaterialConfirmed;
    previousResearchMaterialConfirmedRef.current = isConfirmed;
    if (isConfirmed && autoOpenedConfirmedResearchPointRef.current !== pointDetail.point.id) {
      autoOpenedConfirmedResearchPointRef.current = pointDetail.point.id;
      handleResearchMaterialConfirmed();
      return;
    }
    if (wasConfirmed === false && isConfirmed) {
      handleResearchMaterialConfirmed();
    }
  }, [handleResearchMaterialConfirmed, learning.isResearchMaterialConfirmed, pointDetail]);

  useEffect(() => {
    if (isLoading || error || !pointDetail || !hasSavedJournalNote) return;
    if (autoOpenedSavedJournalPointRef.current === pointDetail.point.id) return;
    autoOpenedSavedJournalPointRef.current = pointDetail.point.id;
    learning.setActiveWorkTab('notes');
    setIsLearningContentOpen(false);
    setIsLearningWorkOpen(true);
    learning.setIsCompletionOpen(false);
    scrollToPointSection('point-work-section');
  }, [error, hasSavedJournalNote, isLoading, learning, pointDetail, scrollToPointSection]);

  useEffect(() => {
    if (isLoading || error || !pointDetail) return;
    if (!canOpenSelfEvaluation) {
      learning.setIsCompletionOpen(false);
      if (!isLearningContentOpen && !isLearningWorkOpen) {
        setIsLearningContentOpen(!canOpenLearningWork);
        setIsLearningWorkOpen(canOpenLearningWork);
      }
    }
  }, [canOpenLearningWork, canOpenSelfEvaluation, error, isLearningContentOpen, isLearningWorkOpen, isLoading, learning.setIsCompletionOpen, pointDetail]);

  useEffect(() => {
    if (isLoading || error || !pointDetail || !planetID) return;
    const requestedStage = typeof params.stage === 'string' ? params.stage : null;
    const canonicalSlug = pointLearningStageToSlug[routeStage];
    const hasObservationNote = (pointDetail.observation_notes ?? []).some((note) => note.content.trim().length > 0);
    const isObservationComplete = pointDetail.point.status === 'completed'
      || (pointDetail.point.point_type === 'exploration' && hasObservationNote)
      || (pointDetail.point.point_type === 'research' && Boolean(pointDetail.research_material_confirmed));
    const nextStage = routeStage !== 'point' && !isObservationComplete ? 'point' : routeStage;
    const nextSlug = pointLearningStageToSlug[nextStage];
    if (requestedStage !== nextSlug || nextStage !== routeStage) {
      router.replace(getPointStageHref(routeKind, planetID, pointDetail.point.id, nextStage));
    }
  }, [error, isLoading, params.stage, planetID, pointDetail, routeKind, routeStage, router]);


  if (isLoading) {
    return (
      <div style={{ ...loadingPageStyle, background: themeTokens.loadingBackground }}>
        <div style={{ ...loadingCardStyle, background: themeTokens.loadingCardBackground, borderColor: themeTokens.loadingCardBorder, color: themeTokens.metaValue }}>
          {pointCopy.loading}
        </div>
      </div>
    );
  }

  if (error || !planet || !pointDetail) {
    return (
      <div style={{ ...loadingPageStyle, background: themeTokens.loadingBackground }}>
        <div style={{ ...loadingCardStyle, display: 'grid', gap: '14px', background: themeTokens.loadingCardBackground, borderColor: themeTokens.loadingCardBorder, color: themeTokens.metaValue }}>
          <div>{error ?? pointCopy.notFound}</div>
          <Link href={planetID ? getPlanetDetailHref(routeKind, planetID) : '/dashboard'} style={{ ...secondaryButtonStyle, background: themeTokens.secondaryButtonBackground, borderColor: themeTokens.secondaryButtonBorder, color: themeTokens.buttonText }}>
            {pointCopy.backToPlanet}
          </Link>
        </div>
      </div>
    );
  }

  const isSharedRoute = routeKind === 'shared';
  const pointGoalContextText = getPointGoalContextText(planet, pointCopy.lumiGuide.fallbackGoal);
  const planetDetailHref = getPlanetDetailHref(routeKind, planetID);
  const focusedDiaryLessonHref = `${planetDetailHref}?diaryLessonId=${encodeURIComponent(pointDetail.lessonID)}#planet-route`;
  const researchMaterialHref = `${planetDetailHref}/points/${pointDetail.point.id}/research-material`;
  const pointPathItems = [
    { key: 'system', label: '☀️', href: '/dashboard', title: 'Star System으로 이동' },
    { key: 'diary', label: `📔 ${planet.planet.title}`, href: planetDetailHref, title: '탐험일지로 이동' },
    { key: 'region', label: `🗺️ ${pointDetail.levelTitle}`, href: focusedDiaryLessonHref, title: '탐험일지에서 해당 리슨 보기' },
    ...(pointDetail.lessonTitle !== pointDetail.levelTitle
      ? [{ key: 'subregion', label: `🧭 ${pointDetail.lessonTitle}`, href: focusedDiaryLessonHref, title: '탐험일지에서 해당 리슨 보기' }]
      : []),
  ];
  const currentPointPathItem = { key: 'current', label: `${isExplorationPoint ? '🎯' : '🔬'} ${pointDetail.point.title}`, title: '현재 페이지 상단으로 이동' };
  const fullPointPathText = [...pointPathItems.map((item) => item.label), currentPointPathItem.label].join(' > ');
  return (
    <div className={`lw-point-learning-page lw-point-learning-page--${theme}`} style={{ ...pageStyle, background: themeTokens.pageBackground, color: themeTokens.metaValue }}>
      <div style={{ ...pageBackgroundOverlayStyle, background: themeTokens.pageOverlay }} />

      <PointPageHeader
        copy={pointCopy.header}
        theme={theme}
        themeTokens={themeTokens}
        planetDetailHref={planetDetailHref}
        onToggleTheme={handleToggleTheme}
        onLogout={handleLogout}
        learnerHeaderCopy={learnerHeaderCopy}
        locale={learning.uiLocale}
      />

      <PointLearningToolbar
        theme={theme}
        themeTokens={themeTokens}
        pointPathItems={pointPathItems}
        currentPointPathItem={currentPointPathItem}
        fullPointPathText={fullPointPathText}
        onCurrentPathClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
        onHeightChange={handlePointToolbarHeightChange}
      />

      <main style={{ ...mainStyle, paddingTop: pointToolbarOffset, paddingBottom: pointBottomGuideOffset }}>
        <PointLearningFlow
          copy={pointCopy}
          learning={learning}
          themeTokens={themeTokens}
          isExplorationPoint={isExplorationPoint}
          isSharedRoute={isSharedRoute}
          editorToolbarStickyTop={pointToolbarOffset}
          pointGoalContextText={pointGoalContextText}
          onResearchMaterialConfirmed={handleResearchMaterialConfirmed}
          onLearningWorkGuideChange={setLearningWorkGuideDetail}
          onSourceContentOpen={handleSourceContentOpen}
          activeStage={routeStage}
          onStageChange={handlePointLearningStageChange}
          getStageHref={getStageHref}
          diaryHref={focusedDiaryLessonHref}
          researchMaterialHref={researchMaterialHref}
        />
      </main>
      <style jsx global>{`
        .lw-point-learning-page button,
        .lw-point-learning-page main a[href],
        .lw-point-learning-page .lw-point-button-face {
          transition:
            transform 120ms ease,
            box-shadow 160ms ease,
            filter 160ms ease,
            border-color 160ms ease,
            background 160ms ease;
          -webkit-tap-highlight-color: transparent;
        }

        .lw-point-learning-page button:not(:disabled),
        .lw-point-learning-page main a[href],
        .lw-point-learning-page .lw-point-button-face {
          box-shadow:
            0 7px 16px rgba(15, 23, 42, 0.13),
            inset 0 1px 0 rgba(255, 255, 255, 0.24);
          cursor: pointer;
        }

        .lw-point-learning-page--dark button:not(:disabled),
        .lw-point-learning-page--dark main a[href],
        .lw-point-learning-page--dark .lw-point-button-face {
          box-shadow:
            0 10px 22px rgba(0, 0, 0, 0.28),
            inset 0 1px 0 rgba(255, 255, 255, 0.10);
        }

        .lw-point-learning-page button:not(:disabled):hover,
        .lw-point-learning-page main a[href]:hover,
        .lw-point-learning-page .lw-point-button-face:hover {
          transform: translateY(-1px);
          filter: brightness(1.04) saturate(1.05);
          box-shadow:
            0 11px 24px rgba(15, 23, 42, 0.18),
            inset 0 1px 0 rgba(255, 255, 255, 0.28);
        }

        .lw-point-learning-page--dark button:not(:disabled):hover,
        .lw-point-learning-page--dark main a[href]:hover,
        .lw-point-learning-page--dark .lw-point-button-face:hover {
          box-shadow:
            0 13px 28px rgba(0, 0, 0, 0.36),
            inset 0 1px 0 rgba(255, 255, 255, 0.14);
        }

        .lw-point-learning-page button:not(:disabled):active,
        .lw-point-learning-page main a[href]:active,
        .lw-point-learning-page .lw-point-button-face:active {
          transform: translateY(1px) scale(0.985);
          filter: brightness(0.98);
          box-shadow:
            0 4px 10px rgba(15, 23, 42, 0.12),
            inset 0 2px 5px rgba(15, 23, 42, 0.10);
        }

        .lw-point-learning-page button:disabled {
          box-shadow: none;
          cursor: not-allowed;
          filter: saturate(0.82);
        }

        .lw-point-learning-page button:focus-visible,
        .lw-point-learning-page main a[href]:focus-visible,
        .lw-point-learning-page .lw-point-button-face:focus-within {
          outline: 2px solid rgba(37, 99, 235, 0.46);
          outline-offset: 3px;
        }

        .lw-point-learning-toolbar-row::-webkit-scrollbar {
          height: 6px;
        }

        .lw-point-learning-toolbar-row::-webkit-scrollbar-track {
          background: transparent;
        }

        .lw-point-learning-toolbar-row::-webkit-scrollbar-thumb {
          background: rgba(100, 116, 139, 0.36);
          border-radius: 999px;
        }


        @media (max-width: 640px) {
          .lw-point-header-meta {
            display: none !important;
          }

        }
      `}</style>
    </div>
  );
}
