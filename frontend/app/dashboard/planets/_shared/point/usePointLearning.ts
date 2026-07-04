'use client';

import { type Dispatch, type SetStateAction, useEffect, useMemo, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';

import type {
  PlanetRouteKind, PlanetMetaState, ResolvedPoint,
  PlanetPointDetailResponse, LearningPointMutationResponse,
  LearningPointSelfEvaluationAIDraftResponse,
  LearningPointSelfEvaluationApplicationQuestion,
  PointQuestionType, PointReplacementCandidate, PointWorkTab, PlanetPoint, ObservationNoteType,
} from '../pointPageTypes';
import {
  PointAttachmentDraft,
  PointQuestionDraft,
  mapPointDetailResponseToResolvedPoint, compactGoalMeta,
  mapAttachmentToDraft,
  mapQuestionToDraft, getPointCompletionReadiness,
  getPlanetDetailHref,
} from '../pointPageUtils';
import {
  isLearningVideoFile,
  normalizeSubtitleFile,
  validateLearningSubtitleFile,
  validateLearningThumbnailFile,
  validateLearningVideoFile,
} from '@/lib/media/videoValidation';
import {
  researchMaterialMaxAttachmentCount,
  validateResearchMaterialReference,
  inferResearchMaterialAttachmentType,
  uploadResearchMaterialAttachmentWithProgress,
} from './pointLearningUtils';
import {
  COMPLETION_CHECK_COUNT,
  type UserInfo,
} from './pointLearningDefaults';
import { getPointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import { usePointLearningSession } from './usePointLearningSession';
import { usePointAttachments } from './usePointAttachments';
import { usePointArtifacts } from './usePointArtifacts';
import { usePointJournal } from './usePointJournal';
import { usePointPracticeLogs } from './usePointPracticeLogs';
import { usePointResearchBlocks } from './usePointResearchBlocks';
import { usePointResearchMaterialAttachments } from './usePointResearchMaterialAttachments';
import { resolveSafetyInputMessage, type SafetyAPIErrorPayload } from '@/lib/safetyErrors';

type PointAIJobResponse = {
  status?: string;
  error_code?: string | null;
  result_ref?: {
    feedback?: string;
    learning_language?: string;
    draft?: LearningPointSelfEvaluationAIDraftResponse['draft'];
  } | null;
};

async function waitForPointAIJob(pollURL: string, fallbackMessage: string): Promise<PointAIJobResponse> {
  const startedAt = Date.now();
  while (Date.now() - startedAt < 45_000) {
    await new Promise((resolve) => setTimeout(resolve, 800));
    const res = await fetch(pollURL, { credentials: 'include', cache: 'no-store' });
    if (!res.ok) throw new Error(fallbackMessage);
    const payload = (await res.json().catch(() => ({}))) as PointAIJobResponse;
    if (payload.status === 'succeeded') return payload;
    if (payload.status === 'failed' || payload.status === 'canceled' || payload.status === 'expired') {
      throw new Error(fallbackMessage);
    }
  }
  throw new Error(fallbackMessage);
}

export type { UsePointLearningResult } from './pointLearningTypes';

export function usePointLearning(routeKind: PlanetRouteKind) {
  const router = useRouter();
  const params = useParams<{ id: string; pointId: string }>();

  const [planet, setPlanet] = useState<PlanetMetaState | null>(null);
  const [pointDetail, setPointDetail] = useState<ResolvedPoint | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [uiLocale, setUILocale] = useState<'ko' | 'en'>('ko');
  const [isLoading, setIsLoading] = useState(true);
  const [isSavingPointRuntime, setIsSavingPointRuntime] = useState(false);
  const [pointRuntimeMessage, setPointRuntimeMessage] = useState<string | null>(null);
  const [pointGoal, setPointGoal] = useState('');
  const [pointCategory, setPointCategory] = useState('');
  const [isSavingPointGoal, setIsSavingPointGoal] = useState(false);
  const [recordStudyMinutes, setRecordStudyMinutes] = useState('0');
  const [recordPracticeCount, setRecordPracticeCount] = useState('0');
  const [recordConfidenceLevel, setRecordConfidenceLevel] = useState('3');
  const [recordApplicationNote, setRecordApplicationNote] = useState('');
  const [isSavingRecord, setIsSavingRecord] = useState(false);
  const [isSavingObservationNote, setIsSavingObservationNote] = useState(false);
  const [isReportingMaterial, setIsReportingMaterial] = useState(false);
  const [materialReportTarget, setMaterialReportTarget] = useState<'source' | 'ai_summary' | 'attachment' | 'other'>('source');
  const [materialReportType, setMaterialReportType] = useState<'broken_link' | 'wrong_content' | 'unsafe_content' | 'copyright' | 'low_quality' | 'other'>('broken_link');
  const [materialReportMessage, setMaterialReportMessage] = useState('');
  const [replacementCandidates, setReplacementCandidates] = useState<PointReplacementCandidate[]>([]);
  const [replacementMessage, setReplacementMessage] = useState<string | null>(null);
  const [isLoadingReplacementCandidates, setIsLoadingReplacementCandidates] = useState(false);
  const [isGeneratingAISummary, setIsGeneratingAISummary] = useState(false);
  const [pointQuestions, setPointQuestions] = useState<PointQuestionDraft[]>([]);
  const [newQuestionTitle, setNewQuestionTitle] = useState('');
  const [newQuestionText, setNewQuestionText] = useState('');
  const [newQuestionType, setNewQuestionType] = useState<PointQuestionType>('reflection');
  const [selfEvalUnderstanding, setSelfEvalUnderstanding] = useState('3');
  const [selfEvalUnderstandingReason, setSelfEvalUnderstandingReason] = useState('');
  const [selfEvalApplicationScore, setSelfEvalApplicationScore] = useState('3');
  const [selfEvalApplicationReason, setSelfEvalApplicationReason] = useState('');
  const [selfEvalProficiency, setSelfEvalProficiency] = useState('3');
  const [selfEvalProficiencyReason, setSelfEvalProficiencyReason] = useState('');
  const [selfEvalProblemSolvingScore, setSelfEvalProblemSolvingScore] = useState('3');
  const [selfEvalProblemSolvingReason, setSelfEvalProblemSolvingReason] = useState('');
  const [selfEvalExpressionScore, setSelfEvalExpressionScore] = useState('3');
  const [selfEvalExpressionReason, setSelfEvalExpressionReason] = useState('');
  const [selfEvalGoalAlignmentNote, setSelfEvalGoalAlignmentNote] = useState('');
  const [selfEvalApplicationQuestions, setSelfEvalApplicationQuestions] = useState<LearningPointSelfEvaluationApplicationQuestion[]>([]);
  const [selfEvalApplicationAnswers, setSelfEvalApplicationAnswers] = useState<string[]>([]);
  const [hasSelfEvalDraft, setHasSelfEvalDraft] = useState(false);
  const [isSelfEvalLocked, setIsSelfEvalLocked] = useState(false);
  const [finalSelfEvalScore, setFinalSelfEvalScore] = useState('');
  const [isSavingQuestion, setIsSavingQuestion] = useState(false);
  const [isGeneratingFeedbackQuestionID, setIsGeneratingFeedbackQuestionID] = useState<string | null>(null);
  const [isSavingSelfEvaluation, setIsSavingSelfEvaluation] = useState(false);
  const [isGeneratingSelfEvalDraft, setIsGeneratingSelfEvalDraft] = useState(false);
  const [pointEntryMessage, setPointEntryMessage] = useState<string | null>(null);
  const [showCompletionAchievement, setShowCompletionAchievement] = useState(false);
  const [activeWorkTab, setActiveWorkTab] = useState<PointWorkTab>('notes');
  const [isCompletionOpen, setIsCompletionOpen] = useState(false);
  const [completionChecks, setCompletionChecks] = useState<boolean[]>(() => Array<boolean>(COMPLETION_CHECK_COUNT).fill(false));

  const planetID = typeof params.id === 'string' ? params.id : '';
  const pointID = typeof params.pointId === 'string' ? params.pointId : '';
  const copy = useMemo(() => getPointLearningCopy(uiLocale), [uiLocale]);
  const runtimeCopy = copy.workspace.runtime;
  const explorationContentCopy = copy.workspace.explorationContent;
  const redirectPath = useMemo(
    () => (planetID && pointID ? `/dashboard/planets/${routeKind}/${planetID}/points/${pointID}` : `/dashboard/planets/${routeKind}`),
    [planetID, pointID, routeKind],
  );

  const learningGoal = planet?.goal_context?.learning_goal?.trim() || planet?.goal_context?.confirmed_goal?.trim() || null;
  const goalMetaItems = useMemo(() => compactGoalMeta(planet?.goal_context), [planet?.goal_context]);
  const isResearchPoint = pointDetail?.point.point_type === 'research';
  const isResearchMaterialConfirmed = isResearchPoint ? Boolean(pointDetail?.research_material_confirmed) : true;
  const pointCompletionReadiness = useMemo(
    () => getPointCompletionReadiness(pointQuestions, pointDetail?.self_evaluation ?? null, pointDetail),
    [pointQuestions, pointDetail],
  );
  const previousPointID = pointDetail?.previousPointID ?? null;
  const nextPointID = pointDetail?.nextPointID ?? null;

  usePointLearningSession({
    routeKind,
    planetID,
    pointID,
    pointDetailID: pointDetail?.point.id ?? null,
  });

  const applyPointMutationPayload = (payload: LearningPointMutationResponse): boolean => {
    if (!payload.planet || !payload.point) return false;
    setPlanet({ planet: payload.planet, goal_context: payload.point.goal_context ?? null });
    setPointDetail(mapPointDetailResponseToResolvedPoint(payload.point as PlanetPointDetailResponse));
    return true;
  };

  const {
    journalObservation,
    setJournalObservation,
    journalReflection,
    setJournalReflection,
    journalExamples,
    setJournalExamples,
    journalNextStep,
    setJournalNextStep,
    hasSavedJournalNote,
    isSavingJournal,
    handleSaveJournal,
    resetJournalDraftToSaved,
  } = usePointJournal({
    routeKind,
    pointID,
    planetID,
    planet,
    pointDetail,
    applyPointMutationPayload,
    setPointEntryMessage,
    copy: runtimeCopy,
  });

  const {
    practiceLogDrafts,
    newPracticeLog,
    isSavingPracticeLog,
    handleCreatePracticeLog,
    handleUpdatePracticeLog,
    handleDeletePracticeLog,
    handleChangePracticeLogDraft,
    handleChangeNewPracticeLog,
  } = usePointPracticeLogs({
    routeKind,
    pointID,
    planetID,
    planet,
    pointDetail,
    applyPointMutationPayload,
    setPointEntryMessage,
    copy: runtimeCopy,
  });

  const {
    artifactType,
    setArtifactType,
    artifactTitle,
    setArtifactTitle,
    artifactURL,
    setArtifactURL,
    artifactDescription,
    setArtifactDescription,
    artifactPointCategory,
    setArtifactPointCategory,
    artifactProductionProcess,
    setArtifactProductionProcess,
    artifactLearnedPoints,
    setArtifactLearnedPoints,
    artifactDifficultPoints,
    setArtifactDifficultPoints,
    artifactVisibility,
    setArtifactVisibility,
    artifactDrafts,
    isSavingArtifact,
    handleSaveArtifact,
    handleChangeArtifactDraft,
    handleUpdateArtifact,
    handleDeleteArtifact,
  } = usePointArtifacts({
    routeKind,
    pointID,
    planetID,
    planet,
    pointDetail,
    applyPointMutationPayload,
    setPointEntryMessage,
    copy: runtimeCopy,
  });

  const answeredQuestionCount = useMemo(() => pointQuestions.filter((q) => q.answer.trim().length > 0).length, [pointQuestions]);

  const {
    attachmentDrafts,
    setAttachmentDrafts,
    newAttachment,
    selectedAttachmentFile,
    setSelectedAttachmentFile,
    isSavingAttachment,
    setIsSavingAttachment,
    handleCreateAttachment,
    handleCreateLinkAttachment,
    handleUploadAttachment,
    handleReplaceAttachmentFile,
    handleOpenAttachment,
    handleUpdateAttachment,
    handleDeleteAttachment,
    handleChangeAttachmentDraft,
    handleChangeNewAttachment,
  } = usePointAttachments({
    routeKind,
    pointID,
    planetID,
    pointDetailID: pointDetail?.point.id ?? null,
    pointAttachments: pointDetail?.attachments ?? null,
    planet,
    applyPointMutationPayload,
    setPointEntryMessage,
    copy: runtimeCopy,
  });

  const researchMaterialAttachmentCount = attachmentDrafts.filter((attachment) => attachment.sourceContext === 'research_material').length;
  const researchMaterialFileAttachmentCount = attachmentDrafts.filter((attachment) => attachment.sourceContext === 'research_material' && !['video', 'subtitle', 'thumbnail'].includes(attachment.attachmentType)).length;
  const hasDirtyResearchMaterialAttachment = attachmentDrafts.some((attachment) => attachment.sourceContext === 'research_material' && attachment.isDirty);
  const canAddResearchMaterialAttachment = researchMaterialFileAttachmentCount < researchMaterialMaxAttachmentCount && !hasDirtyResearchMaterialAttachment;
  const researchMaterialAttachmentLimitMessage = researchMaterialFileAttachmentCount >= researchMaterialMaxAttachmentCount
    ? copy.workspace.researchMaterial.attachmentLimitMessage(researchMaterialMaxAttachmentCount)
    : hasDirtyResearchMaterialAttachment
      ? copy.workspace.researchMaterial.attachmentDirtyMessage
      : null;

  const {
    researchBlocks,
    maxResearchBlockCount,
    canAddResearchBlock,
    researchBlockLimitMessage,
    isSavingResearchBlock,
    isConfirmingResearchMaterial,
    canConfirmResearchMaterial,
    handleAddResearchBlock,
    handleChangeResearchBlock,
    handleSaveResearchBlock,
    handleDeleteResearchBlock,
    handleConfirmResearchMaterial,
    handleUnconfirmResearchMaterial,
  } = usePointResearchBlocks({
    routeKind,
    pointID,
    planetID,
    planet,
    pointDetail,
    isResearchPoint,
    isResearchMaterialConfirmed,
    hasSavedResearchMaterialAttachment: researchMaterialAttachmentCount > 0,
    hasDirtyResearchMaterialAttachment,
    applyPointMutationPayload,
    setPointEntryMessage,
    copy: copy.workspace.researchMaterial,
  });

  const {
    newResearchMaterialAttachment,
    selectedResearchMaterialFile,
    setSelectedResearchMaterialFile,
    researchMaterialAttachmentValidationMessage,
    setResearchMaterialAttachmentValidationMessage,
    researchMaterialUploadProgress,
    setResearchMaterialUploadProgress,
    handleChangeNewResearchMaterialAttachment,
    handleCreateResearchMaterialAttachment,
    handleUploadResearchMaterialAttachment,
    handleUploadResearchMaterialInlineImage,
  } = usePointResearchMaterialAttachments({
    routeKind,
    pointID,
    planetID,
    pointDetailID: pointDetail?.point.id ?? null,
    planet,
    attachmentDrafts,
    canAddResearchMaterialAttachment,
    researchMaterialAttachmentLimitMessage,
    isSavingAttachment,
    setIsSavingAttachment,
    setAttachmentDrafts,
    setPointEntryMessage,
    applyPointMutationPayload,
    copy: copy.workspace.researchMaterial,
  });

  const refreshPointDetail = async (): Promise<boolean> => {
    if (!planetID || !pointID) return false;
    try {
      const pointRes = await fetch(`/api/v1/planets/${routeKind}/${planetID}/points/${pointID}`, { credentials: 'include', cache: 'no-store' });
      if (!pointRes.ok) return false;
      const payload = (await pointRes.json()) as { planet: PlanetMetaState['planet']; point: PlanetPointDetailResponse };
      const resolvedPoint = mapPointDetailResponseToResolvedPoint(payload.point);
      setPlanet({ planet: payload.planet, goal_context: payload.point.goal_context ?? null });
      setPointDetail(resolvedPoint);
      setPointQuestions((resolvedPoint.questions ?? []).map(mapQuestionToDraft));
      return true;
    } catch {
      return false;
    }
  };

  useEffect(() => {
    if (!planetID || !pointID) {
      setError(runtimeCopy.invalidPointPath);
      setIsLoading(false);
      return;
    }
    const load = async () => {
      const refreshRes = await fetch('/api/v1/auth/refresh', { method: 'POST', credentials: 'include' });
      if (!refreshRes.ok) { router.push(`/login?redirect_after=${redirectPath}`); return; }
      const meRes = await fetch('/api/v1/auth/me', { credentials: 'include', cache: 'no-store' });
      if (!meRes.ok) { router.push(`/login?redirect_after=${redirectPath}`); return; }
      const meData = (await meRes.json()) as UserInfo;
      const nextUILocale = meData.ui_locale === 'en' ? 'en' : 'ko';
      setUILocale(nextUILocale);
      if (meData.language_setup_required) {
        router.replace(`${nextUILocale === 'en' ? '/en/language-setup' : '/language-setup'}?redirect_after=${redirectPath}`);
        return;
      }
      if (meData.required_consent_pending) {
        router.replace(`${nextUILocale === 'en' ? '/en/agreements' : '/agreements'}?redirect_after=${redirectPath}`);
        return;
      }
      const pointRes = await fetch(`/api/v1/planets/${routeKind}/${planetID}/points/${pointID}`, { credentials: 'include', cache: 'no-store' });
      if (pointRes.status === 403) throw new Error(runtimeCopy.inactiveDiaryOnlyPlanning);
      if (!pointRes.ok) throw new Error(runtimeCopy.pointLoadFailed);
      const payload = (await pointRes.json()) as { planet: PlanetMetaState['planet']; point: PlanetPointDetailResponse };
      setPlanet({ planet: payload.planet, goal_context: payload.point.goal_context ?? null });
      setPointDetail(mapPointDetailResponseToResolvedPoint(payload.point));
      setIsLoading(false);
    };
    load().catch((e) => { setError(e instanceof Error ? e.message : runtimeCopy.pointLoadFailed); setIsLoading(false); });
  }, [planetID, pointID, redirectPath, routeKind, router, runtimeCopy]);

  useEffect(() => {
    if (!pointDetail) return;
    setPointGoal(pointDetail.point.point_goal ?? '');
    setPointCategory(pointDetail.point.point_category ?? '');
    setRecordStudyMinutes(String(pointDetail.record_entry?.study_minutes ?? 0));
    setRecordPracticeCount(String(pointDetail.record_entry?.practice_count ?? 0));
    setRecordConfidenceLevel(String(pointDetail.record_entry?.confidence_level ?? 3));
    setRecordApplicationNote(pointDetail.record_entry?.application_note ?? '');
    setPointQuestions((pointDetail.questions ?? []).map(mapQuestionToDraft));
    setNewQuestionTitle('');
    setNewQuestionText('');
    setNewQuestionType('reflection');
    setSelfEvalUnderstanding(String(pointDetail.self_evaluation?.understanding_score ?? pointDetail.self_evaluation?.understanding ?? 3));
    setSelfEvalUnderstandingReason(pointDetail.self_evaluation?.understanding_reason ?? pointDetail.self_evaluation?.application_note ?? '');
    setSelfEvalApplicationScore(String(pointDetail.self_evaluation?.application_score ?? pointDetail.self_evaluation?.understanding ?? 3));
    setSelfEvalApplicationReason(pointDetail.self_evaluation?.application_reason ?? pointDetail.self_evaluation?.application_note ?? '');
    setSelfEvalProficiency(String(pointDetail.self_evaluation?.proficiency_score ?? pointDetail.self_evaluation?.proficiency ?? 3));
    setSelfEvalProficiencyReason(pointDetail.self_evaluation?.proficiency_reason ?? pointDetail.self_evaluation?.application_note ?? '');
    setSelfEvalProblemSolvingScore(String(pointDetail.self_evaluation?.problem_solving_score ?? pointDetail.self_evaluation?.proficiency ?? 3));
    setSelfEvalProblemSolvingReason(pointDetail.self_evaluation?.problem_solving_reason ?? pointDetail.self_evaluation?.goal_alignment_note ?? '');
    setSelfEvalExpressionScore(String(pointDetail.self_evaluation?.expression_score ?? pointDetail.self_evaluation?.proficiency ?? 3));
    setSelfEvalExpressionReason(pointDetail.self_evaluation?.expression_reason ?? pointDetail.self_evaluation?.goal_alignment_note ?? '');
    setSelfEvalGoalAlignmentNote(pointDetail.self_evaluation?.goal_alignment_note ?? '');
    setSelfEvalApplicationQuestions([]);
    setSelfEvalApplicationAnswers([]);
    setHasSelfEvalDraft(Boolean(pointDetail.self_evaluation));
    setIsSelfEvalLocked(Boolean(pointDetail.self_evaluation));
    setFinalSelfEvalScore(pointDetail.self_evaluation?.final_score ? String(pointDetail.self_evaluation.final_score) : '');
    setMaterialReportTarget('source');
    setMaterialReportType('broken_link');
    setMaterialReportMessage('');
    setReplacementCandidates([]);
    setReplacementMessage(null);
    setPointEntryMessage(null);
    setShowCompletionAchievement(false);
  }, [pointDetail?.point.id]);

  // ── handlers ──────────────────────────────────────────────────────────────

  const handleSavePointRuntime = async (status: 'in_progress' | 'completed') => {
    if (!planet || !planetID || routeKind !== 'learning' || planet.planet.status !== 'learning' || isSavingPointRuntime) return;
    if (status === 'completed' && !pointCompletionReadiness.isReady) {
      setPointRuntimeMessage(runtimeCopy.completionNeedsReadiness);
      return;
    }
    const wasCompleted = pointDetail?.point.status === 'completed';
    setIsSavingPointRuntime(true);
    setPointRuntimeMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/runtime`, { method: 'PATCH', credentials: 'include', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ status }) });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse;
      if (!res.ok || !applyPointMutationPayload(payload)) {
        if (res.status === 409 && isResearchPoint && !isResearchMaterialConfirmed) throw new Error(runtimeCopy.researchContentRequired);
        if (res.status === 409) throw new Error(runtimeCopy.completionNeedsReadiness);
        if (res.status === 404) throw new Error(runtimeCopy.learningPointNotFound);
        throw new Error(runtimeCopy.pointStatusSaveFailed);
      }
      setPointRuntimeMessage(status === 'completed' ? runtimeCopy.pointCompleted : wasCompleted ? runtimeCopy.pointCompletionCancelled : runtimeCopy.pointInProgress);
      setShowCompletionAchievement(status === 'completed');
      if (status === 'completed') {
        router.push(getPlanetDetailHref('learning', planetID));
      }
    } catch (e) { setPointRuntimeMessage(e instanceof Error ? e.message : runtimeCopy.pointStatusSaveFailed); }
    finally { setIsSavingPointRuntime(false); }
  };

  const handleSavePointGoal = async () => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingPointGoal) return;
    setIsSavingPointGoal(true); setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/goal`, { method: 'PATCH', credentials: 'include', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ point_goal: pointGoal, point_category: pointCategory }) });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(resolveSafetyInputMessage(payload, payload.error ?? runtimeCopy.pointGoalSaveFailed));
      setPointEntryMessage(runtimeCopy.pointGoalSaved);
    } catch (e) { setPointEntryMessage(e instanceof Error ? e.message : runtimeCopy.pointGoalSaveFailed); }
    finally { setIsSavingPointGoal(false); }
  };

  const handleSaveRecord = async () => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingRecord) return;
    setIsSavingRecord(true); setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/record`, { method: 'PATCH', credentials: 'include', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ study_minutes: Number(recordStudyMinutes || '0'), practice_count: Number(recordPracticeCount || '0'), confidence_level: Number(recordConfidenceLevel || '3'), application_note: recordApplicationNote }) });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse;
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(runtimeCopy.pointRecordSaveFailed);
      setPointEntryMessage(runtimeCopy.pointRecordSaved);
    } catch (e) { setPointEntryMessage(e instanceof Error ? e.message : runtimeCopy.pointRecordSaveFailed); }
    finally { setIsSavingRecord(false); }
  };

  const buildReplacementRecommendationQuery = () => {
    const parts = [
      learningGoal,
      pointDetail?.levelTitle,
      pointDetail?.lessonTitle,
      pointDetail?.point.title,
      pointDetail?.point.description,
    ].map((part) => part?.trim()).filter((part): part is string => Boolean(part));
    return [...new Set(parts)].join(' ');
  };

  const handleLoadReplacementCandidates = async () => {
    if (!planet || !pointDetail || routeKind !== 'learning' || isLoadingReplacementCandidates) return;
    const query = buildReplacementRecommendationQuery();
    if (!query) {
      setReplacementMessage(explorationContentCopy.replacementQueryMissing);
      return;
    }
    setIsLoadingReplacementCandidates(true);
    setReplacementMessage(null);
    try {
      const isSubRegion = pointDetail.lessonTitle.trim() && pointDetail.lessonTitle.trim() !== pointDetail.levelTitle.trim();
      const res = await fetch('/api/v1/explorer/recommend', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          query,
          max_results: 5,
          region_title: pointDetail.levelTitle,
          subregion_title: isSubRegion ? pointDetail.lessonTitle : '',
          node_title: pointDetail.point.title,
          node_summary: pointDetail.point.description ?? '',
          course_draft_id: planet.planet.draft_id,
          point_id: pointID,
        }),
      });
      const payload = (await res.json().catch(() => ({}))) as { candidates?: PointReplacementCandidate[]; error?: string };
      if (!res.ok) {
        if (payload.error === 'insufficient_points') throw new Error(explorationContentCopy.insufficientPoints);
        throw new Error(payload.error ?? explorationContentCopy.replacementLoadFailed);
      }
      const candidates = (payload.candidates ?? []).filter((candidate) => candidate.url?.trim());
      setReplacementCandidates(candidates);
      setReplacementMessage(candidates.length ? explorationContentCopy.replacementFound : explorationContentCopy.replacementEmpty);
    } catch (e) {
      setReplacementMessage(e instanceof Error ? e.message : explorationContentCopy.replacementLoadFailed);
    } finally {
      setIsLoadingReplacementCandidates(false);
    }
  };

  const handleUseReplacementLink = async (title: string, url: string, reportID: string, candidate?: PointReplacementCandidate) => {
    const trimmedTitle = title.trim();
    const trimmedURL = url.trim();
    const trimmedReportID = reportID.trim();
    if (!planet || !planetID || routeKind !== 'learning' || isSavingAttachment || !trimmedReportID || !trimmedURL) return false;
    if (!/^https?:\/\//i.test(trimmedURL)) {
      setReplacementMessage(explorationContentCopy.replacementUrlInvalid);
      return false;
    }
    setIsSavingAttachment(true); setPointEntryMessage(null); setReplacementMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/replace-material`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          report_id: trimmedReportID,
          content_id: candidate?.content_id ?? undefined,
          title: trimmedTitle || explorationContentCopy.replacementFallbackTitle,
          url: trimmedURL,
          thumbnail_url: candidate?.thumbnail_url ?? undefined,
        }),
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse;
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(payload.error ?? explorationContentCopy.replacementFailed);
      setReplacementMessage(explorationContentCopy.replacementSuccess);
      setPointEntryMessage(explorationContentCopy.replacementEntrySuccess);
      return true;
    } catch (e) {
      const message = e instanceof Error ? e.message : explorationContentCopy.replacementFailed;
      setReplacementMessage(message);
      setPointEntryMessage(message);
      return false;
    } finally {
      setIsSavingAttachment(false);
    }
  };

  const handleUseReplacementCandidate = async (candidate: PointReplacementCandidate, reportID: string) => {
    const url = candidate.url?.trim() ?? '';
    if (!url) {
      setReplacementMessage(explorationContentCopy.candidateUrlMissing);
      return false;
    }
    const saved = await handleUseReplacementLink(candidate.title || explorationContentCopy.replacementFallbackTitle, url, reportID, candidate);
    return saved;
  };

  const handleCancelMaterialReport = async (reportID: string) => {
    const trimmedReportID = reportID.trim();
    if (!planet || !planetID || routeKind !== 'learning' || isReportingMaterial || !trimmedReportID) return false;
    setIsReportingMaterial(true); setPointEntryMessage(null); setReplacementMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/report-material/${trimmedReportID}/cancel`, {
        method: 'POST',
        credentials: 'include',
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse;
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(payload.error ?? explorationContentCopy.cancelReportFailed);
      setReplacementCandidates([]);
      setReplacementMessage(null);
      setMaterialReportMessage('');
      setPointEntryMessage(explorationContentCopy.cancelReportEntrySuccess);
      return true;
    } catch (e) {
      const message = e instanceof Error ? e.message : explorationContentCopy.cancelReportFailed;
      setPointEntryMessage(message);
      return false;
    } finally {
      setIsReportingMaterial(false);
    }
  };

  const handleUploadArtifactAttachment = async (artifactID: string, file: File, forcedAttachmentType?: 'video' | 'subtitle' | 'thumbnail', titleOverride?: string) => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingAttachment) return false;
    let uploadFile = file;
    const artifactValidationCopy = copy.workspace.researchMaterial;
    if (forcedAttachmentType === 'subtitle') {
      const subtitleValidationError = validateLearningSubtitleFile(uploadFile, artifactValidationCopy);
      if (subtitleValidationError) {
        setResearchMaterialAttachmentValidationMessage(subtitleValidationError);
        return false;
      }
      const hasSubtitle = attachmentDrafts.some((attachment) => attachment.sourceContext === 'artifact' && attachment.artifactID === artifactID && attachment.attachmentType === 'subtitle');
      if (hasSubtitle) {
        setResearchMaterialAttachmentValidationMessage(runtimeCopy.artifactSubtitleDuplicate);
        return false;
      }
      uploadFile = await normalizeSubtitleFile(uploadFile);
    }
    if (forcedAttachmentType === 'thumbnail') {
      const thumbnailValidationError = validateLearningThumbnailFile(uploadFile, artifactValidationCopy);
      if (thumbnailValidationError) {
        setResearchMaterialAttachmentValidationMessage(thumbnailValidationError);
        return false;
      }
      const hasThumbnail = attachmentDrafts.some((attachment) => attachment.sourceContext === 'artifact' && attachment.artifactID === artifactID && attachment.attachmentType === 'thumbnail');
      if (hasThumbnail) {
        setResearchMaterialAttachmentValidationMessage(runtimeCopy.artifactThumbnailDuplicate);
        return false;
      }
    }
    const validationError = forcedAttachmentType === 'thumbnail' ? null : validateResearchMaterialReference(uploadFile.name, uploadFile.size, {
      typeError: artifactValidationCopy.referenceTypeError,
      imageSizeError: artifactValidationCopy.imageSizeError,
      videoSizeError: artifactValidationCopy.videoSizeError,
      documentSizeError: artifactValidationCopy.documentSizeError,
    });
    if (validationError) {
      setResearchMaterialAttachmentValidationMessage(validationError);
      return false;
    }
    if (forcedAttachmentType === 'video' || isLearningVideoFile(uploadFile)) {
      const videoValidationError = await validateLearningVideoFile(uploadFile, artifactValidationCopy);
      if (videoValidationError) {
        setResearchMaterialAttachmentValidationMessage(videoValidationError);
        return false;
      }
      const hasVideo = attachmentDrafts.some((attachment) => attachment.sourceContext === 'artifact' && attachment.artifactID === artifactID && attachment.attachmentType === 'video');
      if (hasVideo) {
        setResearchMaterialAttachmentValidationMessage(runtimeCopy.artifactVideoDuplicate);
        return false;
      }
    }
    const artifactAttachmentCount = attachmentDrafts.filter((attachment) => (
      attachment.sourceContext === 'artifact'
      && attachment.artifactID === artifactID
      && !['video', 'subtitle', 'thumbnail'].includes(attachment.attachmentType)
    )).length;
    if (!forcedAttachmentType && !isLearningVideoFile(uploadFile) && artifactAttachmentCount >= researchMaterialMaxAttachmentCount) {
      setResearchMaterialAttachmentValidationMessage(runtimeCopy.artifactAttachmentLimit);
      return false;
    }
    setIsSavingAttachment(true); setPointEntryMessage(null);
    setResearchMaterialUploadProgress(0);
    try {
      const form = new FormData();
      form.set('file', uploadFile);
      form.set('title', titleOverride?.trim() || uploadFile.name);
      form.set('attachment_type', forcedAttachmentType ?? inferResearchMaterialAttachmentType(uploadFile.name));
      form.set('source_context', 'artifact');
      form.set('artifact_id', artifactID);
      const uploadURL = `/api/v1/planets/learning/${planetID}/points/${pointID}/attachments/upload`;
      let payload: LearningPointMutationResponse;
      try {
        payload = await uploadResearchMaterialAttachmentWithProgress(uploadURL, form, setResearchMaterialUploadProgress, {
          networkError: runtimeCopy.artifactUploadFailed,
          uploadTooLargeServer: artifactValidationCopy.uploadTooLargeServer,
          uploadFailedStatus: runtimeCopy.artifactUploadFailedStatus,
        });
      } catch (error) {
        if (error instanceof Error && error.name === 'UploadHttpError') throw error;
        setResearchMaterialUploadProgress(35);
        const res = await fetch(uploadURL, { method: 'POST', credentials: 'include', body: form });
        payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
        if (!res.ok) throw new Error(resolveSafetyInputMessage(payload, payload.error ?? runtimeCopy.artifactUploadFailedStatus(res.status)));
        setResearchMaterialUploadProgress(100);
      }
      if (!applyPointMutationPayload(payload)) throw new Error(resolveSafetyInputMessage(payload as SafetyAPIErrorPayload, (payload as SafetyAPIErrorPayload).error ?? runtimeCopy.artifactUploadFailed));
      setAttachmentDrafts(((payload.point as PlanetPointDetailResponse).point.attachments ?? []).map(mapAttachmentToDraft));
      setPointEntryMessage(runtimeCopy.artifactUploadSuccess);
      return true;
    } catch (e) {
      setResearchMaterialAttachmentValidationMessage(e instanceof Error ? e.message : runtimeCopy.artifactUploadFailed);
      return false;
    }
    finally { setIsSavingAttachment(false); setResearchMaterialUploadProgress(null); }
  };

  const handleReportMaterial = async () => {
    if (!planet || !planetID || routeKind !== 'learning' || isReportingMaterial || materialReportMessage.trim().length < 5) return;
    setIsReportingMaterial(true); setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/report-material`, { method: 'POST', credentials: 'include', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ target_type: 'source', report_type: materialReportType, message: materialReportMessage }) });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse;
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(payload.error ?? explorationContentCopy.reportFailed);
      setMaterialReportMessage('');
      setPointEntryMessage(explorationContentCopy.reportSuccess);
    } catch (e) { setPointEntryMessage(e instanceof Error ? e.message : explorationContentCopy.reportFailed); }
    finally { setIsReportingMaterial(false); }
  };

  const handleGenerateAISummary = async () => {
    if (!planet || !planetID || routeKind !== 'learning' || isGeneratingAISummary) return;
    setIsGeneratingAISummary(true); setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/ai-summary`, { method: 'POST', credentials: 'include' });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & { error_code?: string; job_id?: string; poll_url?: string };
      if (res.status === 202 || payload.error_code === 'point_ai_summary_pending') {
        await waitForPointAIJob(payload.poll_url ?? `/api/v1/llm-jobs/${payload.job_id}`, explorationContentCopy.aiSummaryFailed);
        if (!await refreshPointDetail()) throw new Error(explorationContentCopy.aiSummaryFailed);
      } else if (!res.ok || !applyPointMutationPayload(payload)) {
        if (res.status === 503) throw new Error(explorationContentCopy.aiUnavailable);
        throw new Error(payload.error ?? explorationContentCopy.aiSummaryFailed);
      }
      setPointEntryMessage(explorationContentCopy.aiSummarySuccess);
    } catch (e) { setPointEntryMessage(e instanceof Error ? e.message : explorationContentCopy.aiSummaryFailed); }
    finally { setIsGeneratingAISummary(false); }
  };

  const handleCreateQuestion = async () => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingQuestion || !newQuestionText.trim()) return false;
    setIsSavingQuestion(true); setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/questions`, { method: 'POST', credentials: 'include', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ title: '', question: newQuestionText, question_type: newQuestionType }) });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload & { question?: NonNullable<PlanetPoint['questions']>[number] };
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(resolveSafetyInputMessage(payload, payload.error ?? runtimeCopy.questionSaveFailed));
      if (payload.question) {
        const createdQuestion = mapQuestionToDraft(payload.question);
        setPointQuestions((current) => {
          const withoutDuplicate = current.filter((question) => question.id !== createdQuestion.id);
          return [...withoutDuplicate, createdQuestion];
        });
      }
      setNewQuestionTitle(''); setNewQuestionText(''); setNewQuestionType('reflection');
      setPointEntryMessage(runtimeCopy.questionCreated);
      return true;
    } catch (e) { setPointEntryMessage(e instanceof Error ? e.message : runtimeCopy.questionSaveFailed); return false; }
    finally { setIsSavingQuestion(false); }
  };

  const handleChangeQuestion = (questionID: string, value: string) => {
    setPointQuestions((cur) => cur.map((q) => q.id === questionID ? { ...q, answer: value } : q));
  };

  const handleChangeQuestionAnswerMethod = (questionID: string, value: string) => {
    setPointQuestions((cur) => cur.map((q) => q.id === questionID ? { ...q, answerMethod: value } : q));
  };

  const handleChangeQuestionDraft = (
    questionID: string,
    field: keyof Pick<PointQuestionDraft, 'title' | 'question' | 'questionType' | 'answerMethod' | 'answer' | 'aiFeedback'>,
    value: string,
  ) => {
    setPointQuestions((cur) => cur.map((q) => q.id === questionID ? { ...q, [field]: value } : q));
  };

  const handleSaveQuestion = async (question: PointQuestionDraft) => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingQuestion) return false;
    setIsSavingQuestion(true); setPointEntryMessage(null);
    try {
      const hasAnswer = question.answer.trim().length > 0;
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/questions/${question.id}`, {
        method: 'PATCH',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: question.title,
          question: question.question,
          question_type: question.questionType,
          answer: question.answer,
          answer_method: question.answerMethod,
          ai_feedback: question.aiFeedback,
          status: hasAnswer ? 'answered' : 'pending',
        }),
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload & { question?: NonNullable<PlanetPoint['questions']>[number] };
      if (!res.ok) {
        if (payload.error === 'point question saved but failed to reload' && await refreshPointDetail()) {
          setPointEntryMessage(runtimeCopy.questionAnswerSaved);
          return true;
        }
        throw new Error(resolveSafetyInputMessage(payload, payload.error ?? runtimeCopy.questionAnswerSaveFailed));
      }
      const applied = applyPointMutationPayload(payload);
      if (payload.question) {
        const savedQuestion = mapQuestionToDraft(payload.question);
        setPointQuestions((current) => current.map((item) => item.id === savedQuestion.id ? savedQuestion : item));
      }
      const refreshed = await refreshPointDetail();
      if (!applied && !refreshed) throw new Error(runtimeCopy.questionAnswerSaveFailed);
      setPointEntryMessage(runtimeCopy.questionAnswerSaved);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : runtimeCopy.questionAnswerSaveFailed);
      return false;
    }
    finally { setIsSavingQuestion(false); }
  };

  const handleDeleteQuestion = async (question: PointQuestionDraft) => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingQuestion) return;
    setIsSavingQuestion(true); setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/questions/${question.id}`, { method: 'DELETE', credentials: 'include' });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse;
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(runtimeCopy.questionDeleteFailed);
      setPointQuestions((current) => current.filter((item) => item.id !== question.id));
      setPointEntryMessage(runtimeCopy.questionDeleted);
    } catch (e) { setPointEntryMessage(e instanceof Error ? e.message : runtimeCopy.questionDeleteFailed); }
    finally { setIsSavingQuestion(false); }
  };

  const handleGenerateQuestionFeedback = async (question: PointQuestionDraft) => {
    if (!planet || !planetID || routeKind !== 'learning' || isGeneratingFeedbackQuestionID || !question.question.trim()) return;
    setIsGeneratingFeedbackQuestionID(question.id); setPointEntryMessage(null);
    try {
      const saved = await handleSaveQuestion(question);
      if (!saved) throw new Error(runtimeCopy.questionAnswerSaveFailed);
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/questions/${question.id}/ai-feedback`, {
        method: 'POST',
        credentials: 'include',
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & { error_code?: string; job_id?: string; poll_url?: string };
      let feedback = payload.feedback ?? '';
      if (res.status === 202 || payload.error_code === 'point_feedback_generate_pending') {
        const jobPayload = await waitForPointAIJob(payload.poll_url ?? `/api/v1/llm-jobs/${payload.job_id}`, runtimeCopy.aiFeedbackFailed);
        feedback = jobPayload.result_ref?.feedback ?? '';
      }
      if (!res.ok && res.status !== 202 || !feedback) {
        if (res.status === 503) throw new Error(runtimeCopy.aiFeedbackUnavailable);
        throw new Error(runtimeCopy.aiFeedbackFailed);
      }
      setPointQuestions((current) => current.map((item) => item.id === question.id ? { ...item, aiFeedback: feedback } : item));
      setPointEntryMessage(runtimeCopy.aiFeedbackCreated);
    } catch (e) { setPointEntryMessage(e instanceof Error ? e.message : runtimeCopy.aiFeedbackFailed); }
    finally { setIsGeneratingFeedbackQuestionID(null); }
  };

  const handleSaveSelfEvaluation = async () => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingSelfEvaluation) return false;
    if (!hasSelfEvalDraft) {
      setPointEntryMessage(runtimeCopy.selfEvalDraftRequired);
      return false;
    }
    setIsSavingSelfEvaluation(true); setPointEntryMessage(null);
    try {
      const finalScore = Number(finalSelfEvalScore || '0');
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/self-evaluation`, { method: 'PATCH', credentials: 'include', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ understanding: Number(selfEvalUnderstanding || '3'), application_note: selfEvalApplicationReason, proficiency: Number(selfEvalProficiency || '3'), understanding_score: Number(selfEvalUnderstanding || '3'), understanding_reason: selfEvalUnderstandingReason, application_score: Number(selfEvalApplicationScore || '3'), application_reason: selfEvalApplicationReason, proficiency_score: Number(selfEvalProficiency || '3'), proficiency_reason: selfEvalProficiencyReason, problem_solving_score: Number(selfEvalProblemSolvingScore || '3'), problem_solving_reason: selfEvalProblemSolvingReason, expression_score: Number(selfEvalExpressionScore || '3'), expression_reason: selfEvalExpressionReason, goal_alignment_note: selfEvalGoalAlignmentNote, final_score: finalScore >= 1 && finalScore <= 5 ? finalScore : null }) });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(resolveSafetyInputMessage(payload, payload.error ?? runtimeCopy.selfEvalSaveFailed));
      setPointEntryMessage(runtimeCopy.selfEvalSaved);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : runtimeCopy.selfEvalSaveFailed);
      return false;
    }
    finally { setIsSavingSelfEvaluation(false); }
  };

  const handleChangeSelfEvalApplicationAnswer = (index: number, value: string) => {
    if (isSelfEvalLocked) return;
    setSelfEvalApplicationAnswers((current) => {
      const next = [...current];
      next[index] = value;
      return next;
    });
    setHasSelfEvalDraft(false);
  };

  const handleGenerateSelfEvaluationDraft = async (forceNewQuestions = false) => {
    if (!planet || !planetID || routeKind !== 'learning' || isGeneratingSelfEvalDraft) return;
    setIsGeneratingSelfEvalDraft(true); setPointEntryMessage(null);
    try {
      const shouldGenerateQuestionsOnly = forceNewQuestions || selfEvalApplicationQuestions.length === 0;
      const applicationAnswers = shouldGenerateQuestionsOnly ? [] : selfEvalApplicationQuestions.map((question, index) => ({
        question: question.question,
        answer: selfEvalApplicationAnswers[index] ?? '',
      })).filter((item) => item.question.trim().length > 0 && item.answer.trim().length > 0);
      const requestInit: RequestInit = shouldGenerateQuestionsOnly
        ? { method: 'POST', credentials: 'include' }
        : {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ application_answers: applicationAnswers }),
          };
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/self-evaluation/ai-draft`, requestInit);
      const payload = (await res.json().catch(() => ({}))) as LearningPointSelfEvaluationAIDraftResponse & SafetyAPIErrorPayload & { error_code?: string; job_id?: string; poll_url?: string };
      let draft = payload.draft;
      if (res.status === 202 || payload.error_code === 'point_self_evaluation_draft_pending') {
        const jobPayload = await waitForPointAIJob(payload.poll_url ?? `/api/v1/llm-jobs/${payload.job_id}`, runtimeCopy.selfEvalDraftFailed);
        draft = jobPayload.result_ref?.draft;
      }
      if (!res.ok && res.status !== 202 || !draft) {
        if (res.status === 503) throw new Error(runtimeCopy.aiFeedbackUnavailable);
        throw new Error(resolveSafetyInputMessage(payload, payload.error ?? runtimeCopy.selfEvalDraftFailed));
      }
      const d = draft;
      if (d.application_questions?.length) {
        const questions = d.application_questions.slice(0, 3);
        setSelfEvalApplicationQuestions(questions);
        setSelfEvalApplicationAnswers((current) => shouldGenerateQuestionsOnly ? questions.map(() => '') : questions.map((_, index) => current[index] ?? ''));
      }
      setSelfEvalUnderstanding(String(d.understanding_score || 3));
      setSelfEvalUnderstandingReason(d.understanding_reason);
      setSelfEvalApplicationScore(String(d.application_score || 3));
      setSelfEvalApplicationReason(d.application_reason);
      setSelfEvalProficiency(String(d.proficiency_score || 3));
      setSelfEvalProficiencyReason(d.proficiency_reason);
      setSelfEvalProblemSolvingScore(String(d.problem_solving_score || 3));
      setSelfEvalProblemSolvingReason(d.problem_solving_reason);
      setSelfEvalExpressionScore(String(d.expression_score || 3));
      setSelfEvalExpressionReason(d.expression_reason);
      setSelfEvalGoalAlignmentNote(d.goal_alignment_note);
      const shouldSaveEvaluation = applicationAnswers.length > 0;
      if (shouldSaveEvaluation) {
        await saveSelfEvaluationDraft(d);
        setIsSelfEvalLocked(true);
      }
      setHasSelfEvalDraft(shouldSaveEvaluation);
      setIsCompletionOpen(true);
      setPointEntryMessage(shouldSaveEvaluation
        ? runtimeCopy.selfEvalSavedNotice
        : runtimeCopy.selfEvalQuestionsReady);
    } catch (e) { setPointEntryMessage(e instanceof Error ? e.message : runtimeCopy.selfEvalDraftFailed); }
    finally { setIsGeneratingSelfEvalDraft(false); }
  };

  const handleRestartSelfEvaluation = async () => {
    setSelfEvalApplicationQuestions([]);
    setSelfEvalApplicationAnswers([]);
    setHasSelfEvalDraft(false);
    setIsSelfEvalLocked(false);
    await handleGenerateSelfEvaluationDraft(true);
  };


  const handleCreateObservationNote = async (noteType: ObservationNoteType, content: string): Promise<boolean> => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingObservationNote) return false;
    setIsSavingObservationNote(true);
    setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/observation-notes`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ note_type: noteType, content }),
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(resolveSafetyInputMessage(payload, payload.error ?? copy.skillFlow.observe.noteSaveFailed));
      setPointEntryMessage(copy.skillFlow.observe.noteSaved);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.skillFlow.observe.noteSaveFailed);
      return false;
    } finally {
      setIsSavingObservationNote(false);
    }
  };

  const handleUpdateObservationNote = async (noteID: string, noteType: ObservationNoteType, content: string): Promise<boolean> => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingObservationNote) return false;
    setIsSavingObservationNote(true);
    setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/observation-notes/${noteID}`, {
        method: 'PATCH',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ note_type: noteType, content }),
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(resolveSafetyInputMessage(payload, payload.error ?? copy.skillFlow.observe.noteSaveFailed));
      setPointEntryMessage(copy.skillFlow.observe.noteSaved);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.skillFlow.observe.noteSaveFailed);
      return false;
    } finally {
      setIsSavingObservationNote(false);
    }
  };

  const handleDeleteObservationNote = async (noteID: string): Promise<boolean> => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingObservationNote) return false;
    setIsSavingObservationNote(true);
    setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/observation-notes/${noteID}`, {
        method: 'DELETE',
        credentials: 'include',
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse;
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(payload.error ?? copy.skillFlow.observe.noteSaveFailed);
      setPointEntryMessage(copy.skillFlow.observe.noteDeleted);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.skillFlow.observe.noteSaveFailed);
      return false;
    } finally {
      setIsSavingObservationNote(false);
    }
  };

  async function saveSelfEvaluationDraft(draft: NonNullable<LearningPointSelfEvaluationAIDraftResponse['draft']>) {
    if (!planet || !planetID || routeKind !== 'learning') return;
    const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/self-evaluation`, {
      method: 'PATCH',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        understanding: Number(draft.understanding_score || 3),
        application_note: draft.application_reason,
        proficiency: Number(draft.proficiency_score || 3),
        understanding_score: Number(draft.understanding_score || 3),
        understanding_reason: draft.understanding_reason,
        application_score: Number(draft.application_score || 3),
        application_reason: draft.application_reason,
        proficiency_score: Number(draft.proficiency_score || 3),
        proficiency_reason: draft.proficiency_reason,
        problem_solving_score: Number(draft.problem_solving_score || 3),
        problem_solving_reason: draft.problem_solving_reason,
        expression_score: Number(draft.expression_score || 3),
        expression_reason: draft.expression_reason,
        goal_alignment_note: draft.goal_alignment_note,
      }),
    });
    const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
    if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(resolveSafetyInputMessage(payload, payload.error ?? runtimeCopy.selfEvalSaveFailed));
  }

  return {
    routeKind, planetID, pointID,
    planet, pointDetail, error, isLoading, uiLocale,
    goalMetaItems, previousPointID, nextPointID, pointCompletionReadiness, learningGoal,
    isResearchMaterialConfirmed, canConfirmResearchMaterial, isConfirmingResearchMaterial, handleConfirmResearchMaterial, handleUnconfirmResearchMaterial,
    pointRuntimeMessage, pointEntryMessage,
    pointGoal, setPointGoal, pointCategory, setPointCategory, isSavingPointGoal, handleSavePointGoal,
    journalObservation, setJournalObservation, journalReflection, setJournalReflection,
    journalExamples, setJournalExamples, journalNextStep, setJournalNextStep,
    hasSavedJournalNote, isSavingJournal, handleSaveJournal, resetJournalDraftToSaved,
    isSavingObservationNote, handleCreateObservationNote, handleUpdateObservationNote, handleDeleteObservationNote,
    recordStudyMinutes, setRecordStudyMinutes, recordPracticeCount, setRecordPracticeCount,
    recordConfidenceLevel, setRecordConfidenceLevel, recordApplicationNote, setRecordApplicationNote,
    isSavingRecord, handleSaveRecord,
    artifactType, setArtifactType, artifactTitle, setArtifactTitle, artifactURL, setArtifactURL,
    artifactDescription, setArtifactDescription, artifactPointCategory, setArtifactPointCategory,
    artifactProductionProcess, setArtifactProductionProcess, artifactLearnedPoints, setArtifactLearnedPoints,
    artifactDifficultPoints, setArtifactDifficultPoints, artifactVisibility, setArtifactVisibility,
    artifactDrafts, isSavingArtifact,
    handleSaveArtifact, handleChangeArtifactDraft, handleUpdateArtifact, handleDeleteArtifact,
    attachmentDrafts, newAttachment, newResearchMaterialAttachment, selectedAttachmentFile, setSelectedAttachmentFile,
    selectedResearchMaterialFile, setSelectedResearchMaterialFile, maxResearchMaterialAttachmentCount: researchMaterialMaxAttachmentCount,
    canAddResearchMaterialAttachment, researchMaterialAttachmentLimitMessage,
    researchMaterialAttachmentValidationMessage, setResearchMaterialAttachmentValidationMessage,
    researchMaterialUploadProgress,
    isSavingAttachment,
    handleCreateAttachment, handleCreateResearchMaterialAttachment, handleUploadAttachment, handleReplaceAttachmentFile, handleUploadResearchMaterialAttachment, handleUploadArtifactAttachment, handleOpenAttachment,
    handleUploadResearchMaterialInlineImage,
    handleUpdateAttachment, handleDeleteAttachment, handleCreateLinkAttachment, handleChangeAttachmentDraft, handleChangeNewAttachment,
    handleChangeNewResearchMaterialAttachment,
    isReportingMaterial, materialReportTarget, setMaterialReportTarget,
    materialReportType, setMaterialReportType, materialReportMessage, setMaterialReportMessage, handleReportMaterial,
    replacementCandidates, replacementMessage, isLoadingReplacementCandidates,
    handleLoadReplacementCandidates, handleUseReplacementCandidate, handleUseReplacementLink, handleCancelMaterialReport,
    practiceLogDrafts, newPracticeLog, isSavingPracticeLog,
    handleCreatePracticeLog, handleUpdatePracticeLog, handleDeletePracticeLog,
    handleChangePracticeLogDraft, handleChangeNewPracticeLog,
    isGeneratingAISummary, handleGenerateAISummary,
    researchBlocks, maxResearchBlockCount, canAddResearchBlock, researchBlockLimitMessage, isSavingResearchBlock,
    handleAddResearchBlock, handleChangeResearchBlock, handleSaveResearchBlock, handleDeleteResearchBlock,
    pointQuestions, newQuestionTitle, setNewQuestionTitle, newQuestionText, setNewQuestionText,
    newQuestionType, setNewQuestionType,
    isSavingQuestion, isGeneratingFeedbackQuestionID,
    refreshPointDetail, handleCreateQuestion, handleChangeQuestionDraft, handleChangeQuestion, handleChangeQuestionAnswerMethod,
    handleSaveQuestion, handleDeleteQuestion, handleGenerateQuestionFeedback,
    selfEvalUnderstanding, setSelfEvalUnderstanding,
    selfEvalUnderstandingReason, setSelfEvalUnderstandingReason,
    selfEvalApplicationScore, setSelfEvalApplicationScore,
    selfEvalApplicationReason, setSelfEvalApplicationReason,
    selfEvalProficiency, setSelfEvalProficiency,
    selfEvalProficiencyReason, setSelfEvalProficiencyReason,
    selfEvalProblemSolvingScore, setSelfEvalProblemSolvingScore,
    selfEvalProblemSolvingReason, setSelfEvalProblemSolvingReason,
    selfEvalExpressionScore, setSelfEvalExpressionScore,
    selfEvalExpressionReason, setSelfEvalExpressionReason,
    selfEvalGoalAlignmentNote, setSelfEvalGoalAlignmentNote,
    selfEvalApplicationQuestions, selfEvalApplicationAnswers, hasSelfEvalDraft, isSelfEvalLocked,
    finalSelfEvalScore, setFinalSelfEvalScore,
    isSavingSelfEvaluation, isGeneratingSelfEvalDraft,
    handleChangeSelfEvalApplicationAnswer, handleRestartSelfEvaluation,
    handleSaveSelfEvaluation, handleGenerateSelfEvaluationDraft,
    isCompletionOpen, setIsCompletionOpen, completionChecks, setCompletionChecks,
    showCompletionAchievement, isSavingPointRuntime, handleSavePointRuntime,
    activeWorkTab, setActiveWorkTab,
    answeredQuestionCount,
  };
}
