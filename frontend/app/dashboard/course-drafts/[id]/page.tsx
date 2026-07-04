'use client';

import React from 'react';
import Link from 'next/link';
import { useParams, useRouter, useSearchParams } from 'next/navigation';
import AppHeaderShell from '@/components/common/AppHeaderShell';
import LearnerHeaderActions, { getLearnerHeaderActionsCopy } from '@/components/common/LearnerHeaderActions';
import ImmersiveOverlayMenu from '@/components/navigation/ImmersiveOverlayMenu';
import LumiModalShell, {
  lumiModalPrimaryButtonStyle,
  lumiModalSecondaryButtonStyle,
} from '@/components/common/LumiModalShell';
import ModalPortal from '@/components/common/ModalPortal';
import { LumiInterviewPanel, type RebuildOption } from '@/components/goal-interview/LumiInterviewPanel';
import { useGoalInterview } from '@/components/goal-interview/useGoalInterview';
import LumiAvatar from '@/components/lumi/LumiAvatar';
import { diarySections } from './constants';
import type { DiarySectionKey } from './types';
import { ArtifactsSection } from './sections/ArtifactsSection';
import { CivilizationSection } from './sections/CivilizationSection';
import { CommunitySection } from './sections/CommunitySection';
import { PlanningSection, type JournalSelectionMeta, type JournalStructureStats } from './sections/PlanningSection';
import { JournalLumiGuide } from './sections/planning/JournalLumiGuide';
import { PlanningLumiGuide, type PlanningLumiGuideFocus } from './sections/planning/PlanningLumiGuide';
import { RecordsSection } from './sections/RecordsSection';
import {
  activeSaveButtonDisabledStyle,
  activeSaveButtonStyle,
  descriptionStyle,
  draftMetaEditorHeaderStyle,
  draftMetaEditorHintStyle,
  draftMetaEditorStyle,
  draftMetaEditorTitleStyle,
  draftMetaFieldStyle,
  draftMetaFooterStyle,
  draftMetaInputStyle,
  draftMetaLabelStyle,
  draftMetaTextareaStyle,
  eyebrowStyle,
  heroCardStyle,
  loadingCardStyle,
  loadingPageStyle,
  lumiAvatarWrapStyle,
  lumiEyebrowStyle,
  lumiInsightCardStyle,
  lumiInsightHeaderStyle,
  lumiMessageStyle,
  lumiMetaPillStyle,
  lumiMetaRowStyle,
  lumiTitleStyle,
  mainStyle,
  metaCardStyle,
  metaGridStyle,
  metaLabelStyle,
  metaValueStyle,
  mobileDiaryNavItemActiveStyle,
  mobileDiaryNavItemLockedStyle,
  mobileDiaryNavItemStyle,
  mobileDiaryNavStyle,
  mobileDiaryNavVisibleStyle,
  navStyle,
  pageBackgroundOverlayStyle,
  pageStyle,
  planningGoalBandBodyStyle,
  planningGoalBandButtonStyle,
  planningGoalBandEyebrowStyle,
  planningGoalBandHeaderStyle,
  planningGoalBandHintStyle,
  planningGoalBandMetaPillStyle,
  planningGoalBandMetaRowStyle,
  planningGoalBandStyle,
  planningGoalBandTitleStyle,
  planningGoalModalCardStyle,
  planningGoalModalCloseButtonStyle,
  planningGoalModalDescriptionStyle,
  planningGoalModalEmptyStateStyle,
  planningGoalModalHeaderStyle,
  planningGoalModalOverlayStyle,
  planningGoalModalPanelWrapStyle,
  planningGoalModalTitleStyle,
  saveMessageStyle,
  secondaryButtonStyle,
  sectionHeaderStyle,
  sectionStyle,
  sectionSubtitleStyle,
  sectionTitleStyle,
  titleStyle,
} from './styles';
import { useDraftEditor } from './useDraftEditor';
import {
  countLessonPointStats,
  countDraftLessonProgress,
} from './draftPageUtils';

const showDraftOverviewPanel = false;
const explorerDiaryBackdropSrc = '/images/explorer/Explorer_Diary_Backdrop.webp';

type PlannedDiarySection = 'community' | 'civilization';

export default function DraftDetailPage() {
  const router = useRouter();
  const params = useParams<{ id: string }>();
  const searchParams = useSearchParams();
  const draftId = typeof params.id === 'string' ? params.id : '';
  const editor = useDraftEditor(draftId);
  const copy = editor.copy;
  const [isPlanDirty, setIsPlanDirty] = React.useState(false);
  const [planningGuideFocus, setPlanningGuideFocus] = React.useState<PlanningLumiGuideFocus>(null);
  const [showGoalEditor, setShowGoalEditor] = React.useState(false);
  const [goalApplyMessage, setGoalApplyMessage] = React.useState<string | null>(null);
  const [isApplyingGoalChange, setIsApplyingGoalChange] = React.useState(false);
  const [journalSelectionMeta, setJournalSelectionMeta] = React.useState<JournalSelectionMeta | null>(null);
  const [journalStructureStats, setJournalStructureStats] = React.useState<JournalStructureStats | null>(null);
  const [showReturnToDiaryConfirm, setShowReturnToDiaryConfirm] = React.useState(false);
  const [plannedDiarySection, setPlannedDiarySection] = React.useState<PlannedDiarySection | null>(null);
  const suppressAutoRedirectRef = React.useRef(false);
  const {
    profile: goalProfile,
    isLoading: isGoalLoading,
    isSending: isGoalSending,
    error: goalError,
    startInterview,
    sendMessage,
    confirmGoal,
    reviseGoal,
    cancelRevision,
    applyRebuildDecision,
  } = useGoalInterview(draftId);

  const {
    isLoading,
    error,
    uiLocale,
    draft,
    reloadDraft,
    draftLumiView,
    activeDiarySection,
    setActiveDiarySection,
    isNarrowViewport,
    isPlanCompleted,
    isLearningStarted,
    isCourseCompleted,
    planningSummary,
    selectedPointContext,
    draftTitleInput,
    setDraftTitleInput,
    draftDescriptionInput,
    setDraftDescriptionInput,
    isSavingDraftMeta,
    saveMessage,
    hasDraftMetaChanges,
    handleSaveDraftMeta,
    handleLogout,
  } = editor;

  const requestedSection = searchParams.get('section');
  const fromSource = searchParams.get('from');
  const hasAppliedRequestedSection = React.useRef(false);
  const isInactiveCourse = Boolean(draft?.draft.is_inactive);

  React.useEffect(() => {
    if (!draft) return;
    if (suppressAutoRedirectRef.current) return;
    if (fromSource === 'diary') return;

    const confirmedCourseId = draft.draft.confirmed_course_id?.trim();
    if (!confirmedCourseId) return;

    if (draft.draft.is_inactive) {
      setActiveDiarySection('planning');
      return;
    }

    if (draft.draft.status === 'confirmed' || draft.draft.status === 'learning') {
      router.replace(`/dashboard/planets/learning/${confirmedCourseId}`);
      return;
    }

    if (draft.draft.status === 'archived') {
      router.replace(`/dashboard/planets/shared/${confirmedCourseId}`);
    }
  }, [draft, fromSource, router, setActiveDiarySection]);

  React.useEffect(() => {
    if (isInactiveCourse) {
      if (activeDiarySection !== 'planning') setActiveDiarySection('planning');
      hasAppliedRequestedSection.current = true;
      return;
    }
    if (hasAppliedRequestedSection.current) return;
    if (requestedSection !== 'journal' && requestedSection !== 'planning') return;
    if (fromSource === 'diary' && requestedSection === 'journal') return;
    if (requestedSection === 'journal' && !isPlanCompleted) return;

    setActiveDiarySection(requestedSection);
    hasAppliedRequestedSection.current = true;
  }, [activeDiarySection, fromSource, isInactiveCourse, isPlanCompleted, requestedSection, setActiveDiarySection]);

  const confirmedCourseId = draft?.draft.confirmed_course_id?.trim() ?? '';
  const returnToDiaryHref = fromSource === 'diary' && confirmedCourseId && !isInactiveCourse
    ? draft?.draft.status === 'archived'
      ? `/dashboard/planets/shared/${confirmedCourseId}`
      : `/dashboard/planets/learning/${confirmedCourseId}`
    : null;
  const handleDiarySectionChange = React.useCallback(
    (section: DiarySectionKey) => {
      if (isInactiveCourse && section !== 'planning') {
        setActiveDiarySection('planning');
        return;
      }
      if (section === 'journal' && returnToDiaryHref) {
        setShowReturnToDiaryConfirm(true);
        return;
      }
      setActiveDiarySection(section);
    },
    [isInactiveCourse, returnToDiaryHref, setActiveDiarySection],
  );

  const handleImmersiveMenuNavigate = React.useCallback(
    () => {
      if (!isPlanDirty) return true;
      return window.confirm(
        uiLocale === 'en'
          ? 'You have unsaved exploration plan changes. Leave without saving?'
          : '저장하지 않은 탐험계획 변경이 있습니다. 저장하지 않고 이동할까요?',
      );
    },
    [isPlanDirty, uiLocale],
  );

  React.useEffect(() => {
    if (!returnToDiaryHref) return;
    if (requestedSection !== 'journal') return;
    if (hasAppliedRequestedSection.current) return;

    hasAppliedRequestedSection.current = true;
    router.replace(returnToDiaryHref);
  }, [requestedSection, returnToDiaryHref, router]);

  const goalMetaPills = React.useMemo(() => {
    if (!goalProfile) return [];
    const pills: string[] = [];
    if (goalProfile.usage_context) pills.push(copy.page.goal.usage(goalProfile.usage_context));
    if (goalProfile.motivation) pills.push(copy.page.goal.motivation(goalProfile.motivation));
    pills.push(`Goal v${goalProfile.version}.0`);
    if (goalProfile.rebuild_decision === 'keep_structure') pills.push(copy.page.goal.learnerEditDecision);
    if (goalProfile.rebuild_decision === 'rebuild_remaining' || goalProfile.rebuild_decision === 'rebuild_all') {
      pills.push(copy.page.goal.rebuildAllDecision);
    }
    return pills;
  }, [copy, goalProfile]);
  const goalRebuildOptions = React.useMemo<RebuildOption[]>(
    () => [
      {
        value: 'rebuild_all',
        label: copy.page.goalRevision.rebuildOptions.rebuildAll.label,
        desc: copy.page.goalRevision.rebuildOptions.rebuildAll.desc,
      },
      {
        value: 'keep_structure',
        label: copy.page.goalRevision.rebuildOptions.keepStructure.label,
        desc: copy.page.goalRevision.rebuildOptions.keepStructure.desc,
      },
    ],
    [copy],
  );
  const goalRebuildDescription = React.useMemo(
    () =>
      isLearningStarted
        ? copy.page.goalRevision.rebuildDescriptionLearning
        : copy.page.goalRevision.rebuildDescriptionBeforeStart,
    [copy, isLearningStarted],
  );

  const handleOpenGoalEditor = React.useCallback(() => {
    if (isInactiveCourse) return;
    setGoalApplyMessage(null);
    setShowGoalEditor(true);
  }, [isInactiveCourse]);

  React.useEffect(() => {
    if (isInactiveCourse && showGoalEditor) {
      setShowGoalEditor(false);
    }
  }, [isInactiveCourse, showGoalEditor]);

  const isGoalBusy = isGoalSending || isApplyingGoalChange;

  const handleCloseGoalEditor = React.useCallback(() => {
    if (!isGoalBusy) setShowGoalEditor(false);
  }, [isGoalBusy]);

  const handleGoalMessage = React.useCallback(
    (message: string) => {
      void sendMessage(message);
    },
    [sendMessage],
  );

  const handleBeginGoalRevision = React.useCallback(() => {
    setGoalApplyMessage(null);
    void reviseGoal(copy.page.goalRevision.reviseSeed);
  }, [copy.page.goalRevision.reviseSeed, reviseGoal]);

  const applyGoalContextToDraft = React.useCallback(async () => {
    const response = await fetch(`/api/v1/course-drafts/${draftId}/goal/apply`, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
    });
    const payload = (await response.json().catch(() => ({}))) as {
      message?: string;
      applied_mode?: string;
      did_rebuild?: boolean;
      rebuild_pending?: boolean;
      error?: string;
    };
    if (!response.ok) {
      throw new Error(payload.error ?? payload.message ?? copy.page.goalRevision.applyFailed);
    }
    return payload;
  }, [copy.page.goalRevision.applyFailed, draftId]);

  const handleConfirmGoalAndApply = React.useCallback(
    async (goal: string) => {
      setGoalApplyMessage(null);
      const confirmed = await confirmGoal(goal);
      if (!confirmed) return;
      if (confirmed.interview_state === 'awaiting_rebuild_decision') {
        setGoalApplyMessage(copy.page.goalRevision.confirmedMessage);
        return;
      }
      try {
        const payload = await applyGoalContextToDraft();
        setGoalApplyMessage(payload.message ?? copy.page.goalRevision.applySuccess);
        suppressAutoRedirectRef.current = true;
        setActiveDiarySection('planning');
        setShowGoalEditor(false);
        await reloadDraft();
      } catch (applyError) {
        setGoalApplyMessage(applyError instanceof Error ? applyError.message : copy.page.goalRevision.applyFailed);
      }
    },
    [applyGoalContextToDraft, confirmGoal, copy, reloadDraft, setActiveDiarySection],
  );

  const handleApplyGoalDecision = React.useCallback(
    async (decision: RebuildOption['value']) => {
      setGoalApplyMessage(null);
      const normalizedDecision = decision === 'rebuild_remaining' ? 'rebuild_all' : decision;
      setIsApplyingGoalChange(true);
      try {
        await applyRebuildDecision(normalizedDecision);
        const payload = await applyGoalContextToDraft();
        setGoalApplyMessage(payload.did_rebuild ? null : (payload.message ?? copy.page.goalRevision.decisionSuccess));
        suppressAutoRedirectRef.current = true;
        setActiveDiarySection('planning');
        setShowGoalEditor(false);
        await reloadDraft();
      } catch (applyError) {
        setGoalApplyMessage(applyError instanceof Error ? applyError.message : copy.page.goalRevision.applyFailed);
      } finally {
        setIsApplyingGoalChange(false);
      }
    },
    [applyGoalContextToDraft, applyRebuildDecision, copy, reloadDraft, setActiveDiarySection],
  );

  const handleCancelGoalRevision = React.useCallback(async () => {
    setGoalApplyMessage(null);
    const restored = await cancelRevision();
    if (restored) {
      setGoalApplyMessage(copy.page.goalRevision.cancelSuccess);
    }
  }, [cancelRevision, copy.page.goalRevision.cancelSuccess]);

  const handleStartGoalInterview = React.useCallback(() => {
    const seedIntent =
      draft?.draft.source_query?.trim() ||
      draft?.draft.title?.trim() ||
      '';
    if (!seedIntent) return;
    void startInterview(seedIntent);
  }, [draft?.draft.source_query, draft?.draft.title, startInterview]);

  if (isLoading) {
    return (
      <div style={loadingPageStyle}>
        <div style={loadingCardStyle}>{copy.page.loading}</div>
      </div>
    );
  }

  if (error || !draft) {
    return (
      <div style={loadingPageStyle}>
        <div style={{ ...loadingCardStyle, display: 'grid', gap: '14px' }}>
          <div>{error ?? copy.page.notFound}</div>
          <Link href="/dashboard" style={secondaryButtonStyle}>
            {copy.page.backToDashboard}
          </Link>
        </div>
      </div>
    );
  }

  const journalLessonProgress = countDraftLessonProgress(draft.lessons);
  const activeJournalLessonProgress = journalStructureStats
    ? {
        total: journalStructureStats.lessonCount,
        completed: journalStructureStats.completedLessonCount,
        learning: journalStructureStats.learningLessonCount,
      }
    : journalLessonProgress;
  const completedExplorationPointCount =
    journalStructureStats?.completedExplorationCount ?? journalLessonProgress.completedExplorationPoints;
  const completedResearchPointCount =
    journalStructureStats?.completedResearchCount ?? journalLessonProgress.completedResearchPoints;
  const completedPointCount =
    journalStructureStats?.completedTotalCount ??
    (completedExplorationPointCount + completedResearchPointCount);
  const explorationPointTotal = journalStructureStats?.explorationCount ?? journalLessonProgress.explorationPoints;
  const researchPointTotal = journalStructureStats?.researchCount ?? journalLessonProgress.researchPoints;
  const totalPointCount = journalStructureStats?.totalCount ?? (explorationPointTotal + researchPointTotal);
  const journalCompletionPercent = activeJournalLessonProgress.total > 0
    ? Math.round((activeJournalLessonProgress.completed / activeJournalLessonProgress.total) * 100)
    : 0;
  const fallbackSectionCopy = copy.planning.rightPanel.bookmarks[0];
  const activeSectionCopy =
    copy.planning.rightPanel.bookmarks.find((section) => section.key === activeDiarySection) ??
    fallbackSectionCopy;
  const journalHintText = !isLearningStarted
    ? copy.page.goal.journalPreLearning
    : journalSelectionMeta?.canOpenPoint
      ? ''
      : journalSelectionMeta?.nodeTitle
        ? copy.page.goal.journalPointPending
        : '';
  const planningGoalHintText = goalProfile?.interview_state === 'awaiting_rebuild_decision'
    ? isLearningStarted
      ? copy.page.goal.planningChangedLearning
      : copy.page.goal.planningChangedBeforeStart
    : isLearningStarted
      ? ''
      : copy.page.goal.planningBeforeStart;

  const isExplorerCanvasSection = activeDiarySection === 'planning' || activeDiarySection === 'journal';
  const sharedExplorerHeaderMeta = (
    <span
      style={{
        display: 'inline-flex',
        alignItems: 'baseline',
        gap: '8px',
        flexWrap: 'wrap',
      }}
    >
      <span style={{ fontSize: '16px', fontWeight: 800, letterSpacing: '0.01em', color: '#F4F7FF' }}>
        {copy.page.headerTitle}
      </span>
      <span style={{ fontSize: '13px', fontWeight: 500, color: 'rgba(200,210,235,0.56)' }}>
        {copy.page.headerSubtitle}
      </span>
    </span>
  );

  return (
    <div
      style={{
        ...pageStyle,
        ...(isExplorerCanvasSection
          ? {
              backgroundImage: `url(${explorerDiaryBackdropSrc})`,
              backgroundSize: 'cover',
              backgroundPosition: 'center',
              backgroundRepeat: 'no-repeat',
            }
          : {}),
      }}
    >
      {isExplorerCanvasSection ? null : <div style={pageBackgroundOverlayStyle} />}

      <AppHeaderShell
        logoHref="/"
        logoIconSize={28}
        logoTextSize="18px"
        maxWidth="1360px"
        headerStyle={navStyle}
        leftMeta={
          <span style={{ display: 'inline-flex', alignItems: 'center', gap: '10px', flexWrap: 'wrap' }}>
            {sharedExplorerHeaderMeta}
          </span>
        }
        rightSlot={
          <span style={{ display: 'inline-flex', alignItems: 'center', gap: '10px' }}>
            <ImmersiveOverlayMenu locale={uiLocale} onBeforeNavigate={handleImmersiveMenuNavigate} />
            <LearnerHeaderActions onLogout={handleLogout} copy={getLearnerHeaderActionsCopy(uiLocale)} />
          </span>
        }
      />

      <main style={mainStyle}>
        {showDraftOverviewPanel ? (
          <section style={heroCardStyle}>
            <div style={{ display: 'grid', gap: '8px' }}>
              <div style={eyebrowStyle}>{copy.page.headerTitle}</div>
              <h1 style={titleStyle}>{copy.page.sectionTitle(draft.draft.title)}</h1>
              <p style={descriptionStyle}>
                {draft.draft.description ?? copy.page.overview.descriptionFallback(draft.draft.source_query)}
              </p>
            </div>

            <div style={metaGridStyle}>
              <div style={metaCardStyle}>
                <span style={metaLabelStyle}>{copy.page.overview.sourceLabel}</span>
                <span style={metaValueStyle}>{draft.draft.source_query}</span>
              </div>
              <div style={metaCardStyle}>
                <span style={metaLabelStyle}>{copy.page.overview.levelLabel}</span>
                <span style={metaValueStyle}>{draft.draft.current_level ?? copy.page.overview.unset}</span>
              </div>
              <div style={metaCardStyle}>
                <span style={metaLabelStyle}>{copy.page.overview.durationLabel}</span>
                <span style={metaValueStyle}>{draft.draft.duration_weeks ? copy.page.overview.weeks(draft.draft.duration_weeks) : copy.page.overview.unset}</span>
              </div>
              <div style={metaCardStyle}>
                <span style={metaLabelStyle}>{copy.page.overview.weeklyHoursLabel}</span>
                <span style={metaValueStyle}>{draft.draft.study_hours_per_week ? copy.page.overview.hours(draft.draft.study_hours_per_week) : copy.page.overview.unset}</span>
              </div>
            </div>

            {draftLumiView ? (
              <div style={lumiInsightCardStyle}>
                <div style={lumiInsightHeaderStyle}>
                  <div style={lumiAvatarWrapStyle}>
                    <LumiAvatar state={draftLumiView.state} size={52} reducedMotion={false} />
                  </div>
                  <div style={{ display: 'grid', gap: '6px' }}>
                    <div style={lumiEyebrowStyle}>Lumi Runtime</div>
                    <strong style={lumiTitleStyle}>{copy.page.overview.lumiTitle}</strong>
                    <p style={lumiMessageStyle}>{draftLumiView.message}</p>
                  </div>
                </div>
                <div style={lumiMetaRowStyle}>
                  <span style={lumiMetaPillStyle}>{draftLumiView.runtimeContext?.page ?? 'course_draft_editor'}</span>
                  <span style={lumiMetaPillStyle}>{draftLumiView.runtimeContext?.action ?? 'page_enter'}</span>
                  <span style={lumiMetaPillStyle}>{draftLumiView.runtimeContext?.scene ?? 'overview'}</span>
                </div>
              </div>
            ) : null}

            <div style={draftMetaEditorStyle}>
              <div style={draftMetaEditorHeaderStyle}>
                <div>
                  <div style={lumiEyebrowStyle}>{copy.page.overview.editEyebrow}</div>
                  <h2 style={draftMetaEditorTitleStyle}>{copy.page.overview.editTitle}</h2>
                </div>
                <div style={draftMetaEditorHintStyle}>{copy.page.overview.editHint}</div>
              </div>

              <label style={draftMetaFieldStyle}>
                <span style={draftMetaLabelStyle}>{copy.page.overview.titleLabel}</span>
                <input
                  value={draftTitleInput}
                  onChange={(event) => setDraftTitleInput(event.target.value)}
                  style={draftMetaInputStyle}
                  maxLength={160}
                />
              </label>

              <label style={draftMetaFieldStyle}>
                <span style={draftMetaLabelStyle}>{copy.page.overview.descriptionLabel}</span>
                <textarea
                  value={draftDescriptionInput}
                  onChange={(event) => setDraftDescriptionInput(event.target.value)}
                  style={draftMetaTextareaStyle}
                  rows={4}
                  maxLength={800}
                />
              </label>

              <div style={draftMetaFooterStyle}>
                <button
                  type="button"
                  onClick={handleSaveDraftMeta}
                  disabled={!hasDraftMetaChanges || isSavingDraftMeta}
                  style={{
                    ...activeSaveButtonStyle,
                    ...((!hasDraftMetaChanges || isSavingDraftMeta) ? activeSaveButtonDisabledStyle : null),
                  }}
                >
                  {isSavingDraftMeta ? copy.page.overview.saving : copy.page.overview.save}
                </button>
                {saveMessage ? <span style={saveMessageStyle}>{saveMessage}</span> : null}
              </div>
            </div>
          </section>
        ) : null}

        <section style={sectionStyle}>
          <div style={sectionHeaderStyle}>
            <div>
              <h2 style={sectionTitleStyle}>
                {isExplorerCanvasSection ? (
                  copy.page.sectionTitle(draft.draft.title?.trim() || copy.page.fallbackCourseName)
                ) : activeSectionCopy.label}
              </h2>
            </div>
            <div style={{ marginLeft: 'auto', minHeight: '34px', display: 'flex', alignItems: 'center', gap: '10px', flexWrap: 'wrap', justifyContent: 'flex-end' }}>
              {returnToDiaryHref ? (
                <button
                  type="button"
                  style={returnToDiaryLinkStyle}
                  onClick={() => setShowReturnToDiaryConfirm(true)}
                >
                  {copy.page.returnToDiary}
                </button>
              ) : null}
              {isExplorerCanvasSection && activeDiarySection === 'journal' ? (
                <JournalLumiGuide
                  canOpenPoint={Boolean(journalSelectionMeta?.canOpenPoint)}
                  copy={{
                    canOpen: copy.page.goal.journalGuideCanOpen,
                    explore: copy.page.goal.journalGuideExplore,
                  }}
                />
              ) : null}
            </div>
          </div>

          {activeDiarySection === 'planning' || activeDiarySection === 'journal' ? (
            <div style={planningGoalBandStyle}>
              <div style={planningGoalBandHeaderStyle}>
                <div style={{ display: 'grid', gap: '6px', minWidth: 0 }}>
                  <div style={planningGoalBandEyebrowStyle}>
                    {activeDiarySection === 'planning' ? copy.page.goal.planningEyebrow : copy.page.goal.journalEyebrow}
                  </div>
                  <div style={planningGoalBandTitleStyle}>
                    {activeDiarySection === 'planning' ? copy.page.goal.planningTitle : copy.page.goal.journalTitle}
                  </div>
                </div>
                {activeDiarySection === 'planning' && !isInactiveCourse ? (
                  <div style={{ display: 'flex', alignItems: 'center', gap: '8px', flexWrap: 'wrap', justifyContent: 'flex-end' }}>
                    <button type="button" style={planningGoalBandButtonStyle} onClick={handleOpenGoalEditor}>
                      {goalProfile?.confirmed_goal ? copy.page.goal.edit : copy.page.goal.set}
                    </button>
                  </div>
                ) : null}
              </div>
              <p style={planningGoalBandBodyStyle}>
                {isGoalLoading
                  ? copy.page.goal.loading
                  : goalProfile?.confirmed_goal?.trim() ||
                    draft.draft.source_query?.trim() ||
                    copy.page.goal.empty}
              </p>
              {goalMetaPills.length > 0 ? (
                <div style={planningGoalBandMetaRowStyle}>
                  {goalMetaPills.map((pill) => (
                    <span key={pill} style={planningGoalBandMetaPillStyle}>
                      {pill}
                    </span>
                  ))}
                  {activeDiarySection === 'journal' && journalSelectionMeta?.regionTitle ? (
                    <span style={planningGoalBandMetaPillStyle}>
                      {copy.page.goal.currentRegion(journalSelectionMeta.regionTitle)}
                    </span>
                  ) : null}
                  {activeDiarySection === 'journal' && journalSelectionMeta?.subRegionTitle ? (
                    <span style={planningGoalBandMetaPillStyle}>
                      {copy.page.goal.currentSubRegion(journalSelectionMeta.subRegionTitle)}
                    </span>
                  ) : null}
                  {activeDiarySection === 'journal' ? (
                    <span style={planningGoalBandMetaPillStyle}>
                      {!isLearningStarted
                        ? copy.page.goal.journalPreLearning
                        : journalSelectionMeta?.canOpenPoint
                          ? copy.page.goal.journalPointReady
                          : copy.page.goal.journalPointPending}
                    </span>
                  ) : null}
                </div>
              ) : null}
              {activeDiarySection === 'journal' ? (
                journalHintText ? <p style={planningGoalBandHintStyle}>{journalHintText}</p> : null
              ) : (
                planningGoalHintText ? <p style={planningGoalBandHintStyle}>{planningGoalHintText}</p> : null
              )}
              {activeDiarySection === 'journal' ? (
                <>
                  <div style={journalSummaryRowStyle}>
                    <span style={journalSummaryPillStyle}>{copy.page.goal.status(copy.page.goal.statusLabel(draft.draft.status))}</span>
                    <span style={journalSummaryPillStyle}>{copy.page.goal.completedRegions(activeJournalLessonProgress.completed, activeJournalLessonProgress.total)}</span>
                    <span style={journalSummaryPillStyle}>
                      {copy.page.goal.completedPoints(completedPointCount, totalPointCount, completedExplorationPointCount, explorationPointTotal, completedResearchPointCount, researchPointTotal)}
                    </span>
                  </div>
                  <div style={journalProgressWrapStyle}>
                  <div style={journalProgressHeaderStyle}>
                    <span>{copy.page.goal.completionRate(journalCompletionPercent)}</span>
                    <span>
                      {copy.page.goal.regionProgress(activeJournalLessonProgress.completed, activeJournalLessonProgress.total)}
                    </span>
                  </div>
                  <div style={journalProgressTrackStyle}>
                    <div
                      style={{
                        ...journalProgressFillStyle,
                        width: `${journalCompletionPercent}%`,
                      }}
                    />
                  </div>
                  <div style={journalProgressMetaStyle}>
                    <span>{copy.page.goal.learningRegions(activeJournalLessonProgress.learning)}</span>
                    <span>{copy.page.goal.remainingRegions(Math.max(0, activeJournalLessonProgress.total - activeJournalLessonProgress.completed))}</span>
                  </div>
                  </div>
                </>
              ) : null}
              {goalApplyMessage ? <p style={planningGoalBandHintStyle}>{goalApplyMessage}</p> : null}
            </div>
          ) : null}

          {activeDiarySection === 'planning' ? (
            <PlanningLumiGuide
              hasUnsavedChanges={isPlanDirty}
              focus={planningGuideFocus}
              isInactiveCourse={isInactiveCourse}
              copy={copy.planning.lumiGuide}
            />
          ) : null}

          <div
            style={{
              ...mobileDiaryNavStyle,
              ...(isNarrowViewport && !isExplorerCanvasSection ? mobileDiaryNavVisibleStyle : null),
            }}
          >
            {diarySections.map((section) => {
              const sectionCopy = copy.planning.rightPanel.bookmarks.find((item) => item.key === section.key);
              const isActive = section.key === activeDiarySection;
              const plannedKey =
                section.key === 'community'
                  ? 'community'
                  : section.key === 'civilization'
                    ? 'civilization'
                    : null;
              const isBlocked =
                section.isLocked ||
                (isInactiveCourse && section.key !== 'planning') ||
                (section.key === 'journal' && !isPlanCompleted) ||
                (['records', 'artifacts', 'community'].includes(section.key) && !isLearningStarted) ||
                (section.key === 'civilization' && !isCourseCompleted);
              const isDisabled = isBlocked && !plannedKey;
              const showInactiveSlash = isInactiveCourse && section.key !== 'planning';
              return (
                <button
                  key={section.key}
                  type="button"
                  onClick={() => {
                    if (plannedKey) {
                      setPlannedDiarySection(plannedKey);
                      return;
                    }
                    if (!isDisabled) handleDiarySectionChange(section.key);
                  }}
                  aria-disabled={isBlocked}
                  disabled={isDisabled}
                  style={{
                    ...mobileDiaryNavItemStyle,
                    ...(isActive ? mobileDiaryNavItemActiveStyle : null),
                    ...(isBlocked ? mobileDiaryNavItemLockedStyle : null),
                    ...(plannedKey ? { cursor: 'pointer' } : null),
                  }}
                >
                  {(sectionCopy ?? fallbackSectionCopy).shortLabel}
                  {showInactiveSlash ? (
                    <span aria-hidden="true" style={mobileInactiveSlashStyle}>/</span>
                  ) : null}
                </button>
              );
            })}
          </div>

          {activeDiarySection === 'planning' ? (
            <PlanningSection
              editor={editor}
              draftId={draftId}
              onDirtyChange={setIsPlanDirty}
              onGuideFocusChange={setPlanningGuideFocus}
              onAfterPlanSave={reloadDraft}
              onSectionChange={handleDiarySectionChange}
              confirmTabChange={!returnToDiaryHref}
              copy={copy.planning}
              treeCopy={copy.tree}
            />
          ) : activeDiarySection === 'journal' ? (
            <PlanningSection
              editor={editor}
              draftId={draftId}
              mode="journal"
              onJournalSelectionChange={setJournalSelectionMeta}
              onJournalStatsChange={setJournalStructureStats}
              onSectionChange={handleDiarySectionChange}
              confirmTabChange={!returnToDiaryHref}
              copy={copy.planning}
              treeCopy={copy.tree}
            />
          ) : activeDiarySection === 'records' ? (
            <RecordsSection editor={editor} copy={copy.supportSections} />
          ) : activeDiarySection === 'artifacts' ? (
            <ArtifactsSection editor={editor} copy={copy.supportSections} />
          ) : activeDiarySection === 'community' ? (
            <CommunitySection editor={editor} copy={copy.supportSections} />
          ) : (
            <CivilizationSection editor={editor} copy={copy.supportSections} />
          )}
        </section>
      </main>

      {showGoalEditor ? (
        <ModalPortal overlayStyle={planningGoalModalOverlayStyle} onMouseDown={handleCloseGoalEditor}>
          <div
            role="dialog"
            aria-modal="true"
            aria-label={copy.page.goalRevision.aria}
            style={planningGoalModalCardStyle}
            onMouseDown={(event) => event.stopPropagation()}
          >
            <div style={planningGoalModalHeaderStyle}>
              <div style={{ minWidth: 0 }}>
                <div style={planningGoalBandEyebrowStyle}>{copy.page.goalRevision.eyebrow}</div>
                <h3 style={planningGoalModalTitleStyle}>{copy.page.goalRevision.title}</h3>
                <p style={planningGoalModalDescriptionStyle}>
                  {copy.page.goalRevision.description}
                </p>
              </div>
              <button
                type="button"
                style={{
                  ...planningGoalModalCloseButtonStyle,
                  ...(isGoalBusy ? activeSaveButtonDisabledStyle : null),
                }}
                onClick={handleCloseGoalEditor}
                disabled={isGoalBusy}
              >
                ×
              </button>
            </div>

            <div style={planningGoalModalPanelWrapStyle}>
              {goalProfile ? (
                <LumiInterviewPanel
                  profile={goalProfile}
                  isSending={isGoalBusy}
                  error={goalError}
                  onSendMessage={handleGoalMessage}
                  onConfirmGoal={(goal) => void handleConfirmGoalAndApply(goal)}
                  onRebuildDecision={(decision) => void handleApplyGoalDecision(decision)}
                  onCancelRebuildDecision={() => void handleCancelGoalRevision()}
                  rebuildOptions={goalRebuildOptions}
                  rebuildTitle={copy.page.goalRevision.rebuildTitle}
                  rebuildDescription={goalRebuildDescription}
                  onConfirmedAction={handleBeginGoalRevision}
                  confirmedActionLabel={copy.page.goalRevision.confirmedAction}
                  isConfirmedActionBusy={isGoalBusy}
                />
              ) : (
                <div style={planningGoalModalEmptyStateStyle}>
                  <div style={{ fontSize: '15px', fontWeight: 800 }}>{copy.page.goalRevision.noInterviewTitle}</div>
                  <div style={{ fontSize: '13px', lineHeight: 1.7 }}>
                    {copy.page.goalRevision.noInterviewDescription}
                  </div>
                  <div>
                    <button
                      type="button"
                      onClick={handleStartGoalInterview}
                      disabled={isGoalBusy}
                      style={{
                        ...planningGoalBandButtonStyle,
                        ...(isGoalBusy ? activeSaveButtonDisabledStyle : null),
                      }}
                    >
                      {isGoalBusy ? copy.page.goalRevision.starting : copy.page.goalRevision.start}
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        </ModalPortal>
      ) : null}

      {isApplyingGoalChange ? (
        <ModalPortal overlayStyle={goalRebuildLoadingOverlayStyle}>
          <div
            role="alertdialog"
            aria-modal="true"
            aria-label={copy.page.rebuildLoading.aria}
            style={goalRebuildLoadingCardStyle}
          >
            <div style={goalRebuildLoadingAvatarWrapStyle}>
              <LumiAvatar state="thinking" size={72} />
            </div>
            <div style={goalRebuildLoadingEyebrowStyle}>{copy.page.rebuildLoading.eyebrow}</div>
            <h3 style={goalRebuildLoadingTitleStyle}>{copy.page.rebuildLoading.title}</h3>
            <p style={goalRebuildLoadingMessageStyle}>
              {copy.page.rebuildLoading.message}
            </p>
            <div style={goalRebuildLoadingStepsStyle}>
              {copy.page.rebuildLoading.steps.map((step) => (
                <span key={step} style={goalRebuildLoadingStepStyle}>{step}</span>
              ))}
            </div>
          </div>
        </ModalPortal>
      ) : null}

      {showReturnToDiaryConfirm && returnToDiaryHref ? (
        <LumiModalShell
          ariaLabel={copy.page.returnConfirm.aria}
          eyebrow={copy.page.returnConfirm.eyebrow}
          lumiState="curious"
          tone="warm"
          title={copy.page.returnConfirm.title}
          message={copy.page.returnConfirm.message}
          onClose={() => setShowReturnToDiaryConfirm(false)}
          actions={
            <>
              <button
                type="button"
                style={lumiModalSecondaryButtonStyle}
                onClick={() => setShowReturnToDiaryConfirm(false)}
              >
                {copy.page.returnConfirm.cancel}
              </button>
              <button
                type="button"
                style={lumiModalPrimaryButtonStyle}
                onClick={() => {
                  setShowReturnToDiaryConfirm(false);
                  router.push(returnToDiaryHref);
                }}
              >
                {copy.page.returnConfirm.confirm}
              </button>
            </>
          }
        />
      ) : null}

      {plannedDiarySection ? (
        <LumiModalShell
          ariaLabel={copy.page.planned[plannedDiarySection].title}
          eyebrow={copy.page.planned.eyebrow}
          lumiState="curious"
          tone="warm"
          title={copy.page.planned[plannedDiarySection].title}
          message={copy.page.planned[plannedDiarySection].message}
          onClose={() => setPlannedDiarySection(null)}
          actions={
            <button
              type="button"
              style={lumiModalPrimaryButtonStyle}
              onClick={() => setPlannedDiarySection(null)}
            >
              {copy.page.planned.confirm}
            </button>
          }
        />
      ) : null}
    </div>
  );
}

const journalProgressWrapStyle = {
  display: 'grid',
  gap: '8px',
  marginTop: '4px',
} as const;

const journalProgressHeaderStyle = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '12px',
  flexWrap: 'wrap',
  color: '#F6E9C5',
  fontSize: '12px',
  fontWeight: 700,
} as const;

const journalProgressTrackStyle = {
  position: 'relative',
  width: '100%',
  height: '12px',
  borderRadius: '999px',
  overflow: 'hidden',
  background: 'rgba(255,255,255,0.12)',
  border: '1px solid rgba(255,255,255,0.14)',
} as const;

const journalProgressFillStyle = {
  position: 'absolute',
  inset: 0,
  background: 'linear-gradient(90deg, rgba(103, 232, 249, 0.92), rgba(250, 204, 21, 0.96))',
  borderRadius: '999px',
} as const;

const journalProgressMetaStyle = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '12px',
  flexWrap: 'wrap',
  color: 'rgba(220,230,255,0.78)',
  fontSize: '12px',
} as const;

const journalSummaryRowStyle = {
  display: 'flex',
  alignItems: 'center',
  gap: '8px',
  flexWrap: 'wrap',
  marginTop: '4px',
} as const;

const journalSummaryPillStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: '30px',
  padding: '0 12px',
  borderRadius: '999px',
  background: 'rgba(255,255,255,0.10)',
  border: '1px solid rgba(255,255,255,0.14)',
  color: '#F6E9C5',
  fontSize: '12px',
  fontWeight: 700,
  lineHeight: 1.4,
} as const;

const returnToDiaryLinkStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '34px',
  padding: '0 14px',
  borderRadius: '999px',
  background: 'rgba(255,255,255,0.08)',
  border: '1px solid rgba(194, 210, 245, 0.18)',
  color: '#F4F7FF',
  textDecoration: 'none',
  fontSize: '12px',
  fontWeight: 700,
  lineHeight: 1,
  whiteSpace: 'nowrap',
  cursor: 'pointer',
} as const;

const mobileInactiveSlashStyle = {
  position: 'absolute',
  inset: 0,
  display: 'grid',
  placeItems: 'center',
  color: 'rgba(226, 232, 240, 0.92)',
  fontSize: '34px',
  fontWeight: 900,
  lineHeight: 1,
  textShadow: '0 2px 6px rgba(2, 6, 23, 0.72)',
  pointerEvents: 'none',
} as const;

const goalRebuildLoadingOverlayStyle = {
  padding: '24px',
  background: 'rgba(8, 12, 20, 0.76)',
  backdropFilter: 'blur(12px)',
  cursor: 'wait',
} as const;

const goalRebuildLoadingCardStyle = {
  width: 'min(460px, 100%)',
  display: 'grid',
  justifyItems: 'center',
  gap: '14px',
  padding: '28px 28px 24px',
  borderRadius: '28px',
  border: '1px solid rgba(245, 204, 110, 0.38)',
  background: 'linear-gradient(180deg, rgba(255, 246, 226, 0.98), rgba(244, 226, 190, 0.98))',
  boxShadow: '0 30px 90px rgba(0, 0, 0, 0.42)',
  textAlign: 'center',
} as const;

const goalRebuildLoadingAvatarWrapStyle = {
  width: '88px',
  height: '88px',
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  borderRadius: '999px',
  background: 'radial-gradient(circle at 35% 30%, rgba(255,255,255,0.62), rgba(255,255,255,0.16) 72%, rgba(255,255,255,0))',
  boxShadow: 'inset 0 1px 0 rgba(255,255,255,0.72), 0 18px 38px rgba(126, 80, 18, 0.18)',
} as const;

const goalRebuildLoadingEyebrowStyle = {
  color: '#9A5B13',
  fontSize: '12px',
  fontWeight: 900,
  letterSpacing: '0.1em',
  textTransform: 'uppercase',
} as const;

const goalRebuildLoadingTitleStyle = {
  margin: 0,
  color: '#422607',
  fontSize: '23px',
  lineHeight: 1.28,
  fontWeight: 900,
} as const;

const goalRebuildLoadingMessageStyle = {
  margin: 0,
  color: 'rgba(83, 50, 12, 0.88)',
  fontSize: '14px',
  lineHeight: 1.7,
} as const;

const goalRebuildLoadingStepsStyle = {
  display: 'flex',
  justifyContent: 'center',
  gap: '8px',
  flexWrap: 'wrap',
  paddingTop: '2px',
} as const;

const goalRebuildLoadingStepStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: '28px',
  padding: '0 10px',
  borderRadius: '999px',
  border: '1px solid rgba(152, 91, 20, 0.2)',
  background: 'rgba(255, 255, 255, 0.44)',
  color: '#7A4811',
  fontSize: '12px',
  fontWeight: 800,
} as const;
