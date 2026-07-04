'use client';

import { useRouter } from 'next/navigation';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { normalizeLocale } from '@/lib/i18n/locales';
import { getDashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft';
import { createLumiViewStateFromRuntime } from '@/lib/lumi/lumiRuntimeView';
import { useLumiRuntimeConfigBootstrap } from '@/lib/lumi/useLumiRuntimeConfigBootstrap';
import { diarySections } from './constants';
import { getExplorationPoints, getResearchPoints } from './selectors';
import type {
  DetailTargetKind, DiarySectionKey, DraftAggregate, DraftLessonTree,
  DraftPointAggregate, UpdateDraftStructureResponse, UserInfo,
} from './types';
import { useDraftStructure } from './useDraftStructure';
import { useDraftDetailEdit } from './useDraftDetailEdit';
import { useDraftResources } from './useDraftResources';
import { useDraftEntries } from './useDraftEntries';

export function useDraftEditor(draftId: string) {
  const router = useRouter();
  const lumiConfigVersion = useLumiRuntimeConfigBootstrap();

  // ── core state ────────────────────────────────────────────────────────────
  const [user, setUser] = useState<UserInfo | null>(null);
  const [draft, setDraft] = useState<DraftAggregate | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // ── navigation state ──────────────────────────────────────────────────────
  const [activeDiarySection, setActiveDiarySection] = useState<DiarySectionKey>('planning');
  const [selectedRegionId, setSelectedRegionId] = useState<string | null>(null);
  const [selectedPointId, setSelectedPointId] = useState<string | null>(null);
  const [selectedResearchNodeId, setSelectedResearchNodeId] = useState<string | null>(null);
  const [detailTargetKind, setDetailTargetKind] = useState<DetailTargetKind>('point');
  const [isNarrowViewport, setIsNarrowViewport] = useState(false);

  // ── draft meta state ──────────────────────────────────────────────────────
  const [draftTitleInput, setDraftTitleInput] = useState('');
  const [draftDescriptionInput, setDraftDescriptionInput] = useState('');
  const [isSavingDraftMeta, setIsSavingDraftMeta] = useState(false);
  const [saveMessage, setSaveMessage] = useState<string | null>(null);

  // ── computed ──────────────────────────────────────────────────────────────

  const draftLumiView = useMemo(() => {
    if (!draft) return null;
    const hasOutline = draft.lessons.length > 0;
    return createLumiViewStateFromRuntime({
      runtimeContext: { mode: hasOutline ? 'ai' : 'general', page: 'course_draft_editor', action: hasOutline ? 'create_course' : 'page_enter', scene: hasOutline ? 'success' : 'overview', userState: { currentCourseId: draft.draft.id } },
      trigger: hasOutline ? 'generation_success' : 'page_enter',
      payload: { context: 'diary', courseTitle: draft.draft.title },
    });
  }, [draft, lumiConfigVersion]);

  const uiLocale = normalizeLocale(user?.ui_locale);
  const copy = useMemo(() => getDashboardCourseDraftCopy(uiLocale), [uiLocale]);

  const planningSummary = useMemo(() => {
    if (!draft) {
      return { regionCount: 0, explorationPointCount: 0, researchNodeCount: 0, selectedLevel: null as DraftLessonTree | null, selectedLesson: null as DraftLessonTree | null, selectedResearchNode: null as DraftPointAggregate | null };
    }
    const lessons = draft.lessons;
    const selectedLevel = lessons.find((t) => t.lesson.id === selectedRegionId) ?? lessons[0] ?? null;
    const selectedLesson = detailTargetKind === 'point'
      ? (selectedLevel?.sub_lessons ?? []).find((s) => s.lesson.id === selectedPointId) ?? (selectedLevel?.sub_lessons ?? [])[0] ?? null
      : null;
    const researchPoints = selectedLevel ? getResearchPoints(selectedLevel) : [];
    const selectedResearchNode = detailTargetKind === 'research_node'
      ? researchPoints.find((p) => p.point.id === selectedResearchNodeId) ?? researchPoints[0] ?? null
      : null;
    return {
      regionCount: lessons.reduce((sum, t) => sum + (t.sub_lessons?.length ?? 0), 0),
      explorationPointCount: lessons.reduce((sum, t) => sum + (t.sub_lessons ?? []).reduce((s, sub) => s + getExplorationPoints(sub).length, 0), 0),
      researchNodeCount: lessons.reduce((sum, t) => sum + getResearchPoints(t).length, 0),
      selectedLevel, selectedLesson, selectedResearchNode,
    };
  }, [detailTargetKind, draft, selectedPointId, selectedRegionId, selectedResearchNodeId]);

  const selectedDetailKey = planningSummary.selectedLesson
    ? `point:${planningSummary.selectedLesson.lesson.id}`
    : planningSummary.selectedResearchNode
      ? `research_node:${planningSummary.selectedResearchNode.point.id}`
      : planningSummary.selectedLevel
        ? `region:${planningSummary.selectedLevel.lesson.id}`
        : 'none';

  const selectedPointResources = planningSummary.selectedLesson
    ? getExplorationPoints(planningSummary.selectedLesson).map((p) => ({
        id: p.point.id,
        title: p.point.title,
        resource_type: 'content',
        selection_state: p.point.selection_state ?? 'candidate',
        description: p.point.description ?? null,
        thumbnail_url: p.point.thumbnail_url ?? null,
        price_type: p.point.price_type ?? null,
        rank_score: p.point.rank_score ?? null,
        external_url: p.point.external_url ?? null,
      }))
    : [];

  const selectedExplorationPoint = selectedPointResources.find((r) => r.selection_state === 'selected') ?? selectedPointResources[0] ?? null;

  const selectedPointContext = useMemo(() => {
    if (planningSummary.selectedResearchNode) {
      const templateType = planningSummary.selectedResearchNode.point.template_type;
      const templateLabel = copy.supportSections.selectedPoint.researchTemplateLabel(templateType);
      return {
        kind: 'research_node' as const,
        label: copy.supportSections.selectedPoint.researchLabel,
        title: planningSummary.selectedResearchNode.point.title,
        description: templateLabel
          ? copy.supportSections.selectedPoint.researchTemplateDescription(templateLabel)
          : copy.supportSections.selectedPoint.researchDescription,
      };
    }
    if (selectedExplorationPoint) {
      const selectionLabel = selectedExplorationPoint.selection_state === 'selected'
        ? copy.supportSections.selectedPoint.selectedExploration
        : selectedExplorationPoint.selection_state === 'candidate'
          ? copy.supportSections.selectedPoint.candidateExploration
          : copy.supportSections.selectedPoint.genericExploration;
      return {
        kind: 'exploration_point' as const,
        label: copy.supportSections.selectedPoint.explorationLabel,
        title: selectedExplorationPoint.title,
        description: copy.supportSections.selectedPoint.formatExplorationDescription(selectionLabel, selectedExplorationPoint.resource_type),
      };
    }
    return null;
  }, [copy, planningSummary.selectedResearchNode, selectedExplorationPoint]);

  const resourceReadinessSummary = useMemo(() => {
    if (!draft) return { totalPoints: 0, pointsWithSelectedResources: 0, selectedResourceCount: 0, candidateResourceCount: 0, pointsWithCandidateOnlyResources: 0 };
    return draft.lessons.reduce((summary, mainTree) => {
      (mainTree.sub_lessons ?? []).forEach((subTree) => {
        const pts = getExplorationPoints(subTree);
        const selectedCount = pts.filter((p) => p.point.selection_state === 'selected').length;
        const candidateCount = pts.filter((p) => p.point.selection_state === 'candidate').length;
        summary.totalPoints += 1;
        summary.selectedResourceCount += selectedCount;
        summary.candidateResourceCount += candidateCount;
        if (selectedCount > 0) summary.pointsWithSelectedResources += 1;
        if (selectedCount === 0 && candidateCount > 0) summary.pointsWithCandidateOnlyResources += 1;
      });
      return summary;
    }, { totalPoints: 0, pointsWithSelectedResources: 0, selectedResourceCount: 0, candidateResourceCount: 0, pointsWithCandidateOnlyResources: 0 });
  }, [draft]);

  const hasDraftMetaChanges = Boolean(draft && (draftTitleInput.trim() !== draft.draft.title || draftDescriptionInput.trim() !== (draft.draft.description ?? '')));
  const isInactiveCourse = Boolean(draft?.draft.is_inactive);
  const isPlanCompleted = Boolean(draft && draft.draft.status !== 'draft');
  const isLearningStarted = Boolean(draft && !isInactiveCourse && (draft.draft.status === 'learning' || draft.draft.status === 'archived'));
  const isCourseCompleted = Boolean(draft && !isInactiveCourse && draft.draft.status === 'archived');

  const isSectionAllowed = (key: DiarySectionKey): boolean => {
    if (isInactiveCourse) return key === 'planning';
    if (key === 'planning') return true;
    if (key === 'journal') return isPlanCompleted;
    if (key === 'records' || key === 'artifacts' || key === 'community') return isLearningStarted;
    if (key === 'civilization') return isCourseCompleted;
    return false;
  };

  const setActiveDiarySectionGuarded = (key: DiarySectionKey) => {
    if (!isSectionAllowed(key)) return;
    setActiveDiarySection(key);
  };

  const activeSection = diarySections.find((s) => s.key === activeDiarySection) ?? diarySections[0];

  // ── sub-hooks ─────────────────────────────────────────────────────────────

  const structure = useDraftStructure({
    draft, setDraft, draftId,
    setSelectedRegionId, setSelectedPointId, setSelectedResearchNodeId, setDetailTargetKind,
  });

  const detailEdit = useDraftDetailEdit({
    draft, setDraft, draftId, planningSummary, selectedDetailKey,
    hasStructureDraftChanges: structure.hasStructureDraftChanges,
    setSelectedRegionId, setSelectedPointId, setSelectedResearchNodeId, setDetailTargetKind,
    copy: copy.hook.detail,
  });

  const resources = useDraftResources({
    draft, setDraft,
    selectedLesson: planningSummary.selectedLesson,
    selectedDetailKey,
    copy: copy.hook.resources,
  });

  const entries = useDraftEntries({
    draft, setDraft, draftId,
    selectedLesson: planningSummary.selectedLesson,
    selectedDetailKey,
    copy: copy.hook.entries,
  });

  // ── computed (cross-domain) ───────────────────────────────────────────────

  const hasEveryRegionPoint = Boolean(draft?.lessons.length && draft.lessons.every((t) => (t.sub_lessons?.length ?? 0) > 0));
  const hasResourceCoverage = Boolean(
    draft?.lessons.some((t) => (t.sub_lessons ?? []).some((s) => getExplorationPoints(s).length > 0))
    || resources.hasManualResourceDraft,
  );
  const hasSelectedResourceCoverage = resourceReadinessSummary.selectedResourceCount > 0;
  const pointsMissingSelectedResources = resourceReadinessSummary.totalPoints - resourceReadinessSummary.pointsWithSelectedResources;
  const hasEveryPointSelectedResource = Boolean(resourceReadinessSummary.totalPoints > 0 && pointsMissingSelectedResources === 0);
  const hasEveryResearchNodeTitle = Boolean(
    !draft?.lessons.length || draft.lessons.every((t) => getResearchPoints(t).every((p) => p.point.title.trim().length > 0)),
  );
  const hasBlockingUnsavedChanges = Boolean(
    hasDraftMetaChanges
    || structure.hasStructureDraftChanges
    || detailEdit.hasDetailDraftChanges
    || entries.hasJournalDraftChanges
    || entries.hasRecordDraftChanges
    || entries.hasArtifactDraftChanges
    || resources.hasUnattachedManualResourceDraft,
  );

  // ── useEffects ────────────────────────────────────────────────────────────

  useEffect(() => {
    if (!isSectionAllowed(activeDiarySection)) setActiveDiarySection('planning');
  }, [isPlanCompleted, isLearningStarted, isInactiveCourse, activeDiarySection]);

  useEffect(() => {
    if (!draft || draft.lessons.length === 0) { setSelectedRegionId(null); setSelectedPointId(null); setSelectedResearchNodeId(null); return; }
    const lessons = draft.lessons;
    const regionExists = selectedRegionId ? lessons.some((t) => t.lesson.id === selectedRegionId) : false;
    const nextRegion = regionExists ? selectedRegionId : lessons[0]?.lesson.id ?? null;
    const nextTree = lessons.find((t) => t.lesson.id === nextRegion) ?? null;
    const pointExists = selectedPointId ? (nextTree?.sub_lessons ?? []).some((s) => s.lesson.id === selectedPointId) : false;
    const nextPoint = detailTargetKind === 'point' ? pointExists ? selectedPointId : nextTree?.sub_lessons?.[0]?.lesson.id ?? null : null;
    const researchPoints = nextTree ? getResearchPoints(nextTree) : [];
    const researchNodeExists = selectedResearchNodeId ? researchPoints.some((p) => p.point.id === selectedResearchNodeId) : false;
    const nextResearchNode = detailTargetKind === 'research_node' ? researchNodeExists ? selectedResearchNodeId : researchPoints[0]?.point.id ?? null : null;
    if (nextRegion !== selectedRegionId) setSelectedRegionId(nextRegion);
    if (nextPoint !== selectedPointId) setSelectedPointId(nextPoint);
    if (nextResearchNode !== selectedResearchNodeId) setSelectedResearchNodeId(nextResearchNode);
  }, [detailTargetKind, draft, selectedPointId, selectedRegionId, selectedResearchNodeId]);

  useEffect(() => {
    if (!draft) return;
    setDraftTitleInput(draft.draft.title);
    setDraftDescriptionInput(draft.draft.description ?? '');
  }, [draft]);

  useEffect(() => {
    const sync = () => setIsNarrowViewport(window.innerWidth < 1024);
    sync();
    window.addEventListener('resize', sync);
    return () => window.removeEventListener('resize', sync);
  }, []);

  const loadDraft = useCallback(async () => {
    if (!draftId) { setError(copy.page.invalidDraftId); setIsLoading(false); return; }
    const refreshRes = await fetch('/api/v1/auth/refresh', { method: 'POST', credentials: 'include' });
    if (!refreshRes.ok) { router.push(copy.hook.loginPath(draftId)); return; }
    const meRes = await fetch('/api/v1/auth/me', { credentials: 'include', cache: 'no-store' });
    if (!meRes.ok) { router.push(copy.hook.loginPath(draftId)); return; }
    const meData = (await meRes.json()) as UserInfo;
    setUser(meData);
    const nextCopy = getDashboardCourseDraftCopy(meData.ui_locale);
    if (meData.language_setup_required) {
      router.replace(normalizeLocale(meData.ui_locale) === 'en'
        ? `/en/language-setup?redirect_after=/dashboard/course-drafts/${draftId}`
        : `/language-setup?redirect_after=/dashboard/course-drafts/${draftId}`);
      return;
    }
    if (meData.required_consent_pending) { router.replace(nextCopy.hook.agreementsPath(draftId)); return; }
    const draftRes = await fetch(`/api/v1/course-drafts/${draftId}`, { credentials: 'include', cache: 'no-store' });
    if (!draftRes.ok) throw new Error(nextCopy.page.loadFailed);
    const payload = (await draftRes.json()) as { draft: DraftAggregate };
    setDraft(payload.draft);
    setIsLoading(false);
  }, [copy, draftId, router]);

  useEffect(() => {
    loadDraft().catch((e) => { setError(e instanceof Error ? e.message : copy.page.loadFailed); setIsLoading(false); });
  }, [copy.page.loadFailed, loadDraft]);

  // ── handlers ──────────────────────────────────────────────────────────────

  const handleSelectRegion = (mainTree: DraftLessonTree) => {
    setSelectedRegionId(mainTree.lesson.id);
    setSelectedPointId(null);
    setSelectedResearchNodeId(null);
    setDetailTargetKind('region');
  };

  const handleSelectPoint = (mainTree: DraftLessonTree, subTree: DraftLessonTree) => {
    setSelectedRegionId(mainTree.lesson.id);
    setSelectedPointId(subTree.lesson.id);
    setSelectedResearchNodeId(null);
    setDetailTargetKind('point');
  };

  const handleSelectResearchNode = (mainTree: DraftLessonTree, point: DraftPointAggregate) => {
    setSelectedRegionId(mainTree.lesson.id);
    setSelectedPointId(null);
    setSelectedResearchNodeId(point.point.id);
    setDetailTargetKind('research_node');
  };

  const handleSaveDraftMeta = async () => {
    if (!draft || isSavingDraftMeta || !hasDraftMetaChanges) return;
    const nextTitle = draftTitleInput.trim();
    if (!nextTitle) { setSaveMessage(copy.hook.saveTitleRequired); return; }
    setIsSavingDraftMeta(true);
    setSaveMessage(null);
    try {
      const res = await fetch(`/api/v1/course-drafts/${draft.draft.id}`, {
        method: 'PATCH', credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title: nextTitle, description: draftDescriptionInput.trim() }),
      });
      if (!res.ok) throw new Error(copy.hook.saveFailed);
      const payload = (await res.json()) as { draft: DraftAggregate };
      setDraft(payload.draft);
      setSaveMessage(copy.hook.saveSuccess);
    } catch (e) {
      setSaveMessage(e instanceof Error ? e.message : copy.hook.saveFailed);
    } finally {
      setIsSavingDraftMeta(false);
    }
  };

  const handleLogout = async () => {
    await fetch('/api/v1/auth/logout', { method: 'POST', credentials: 'include' }).catch(() => undefined);
    router.replace('/login');
  };

  return {
    // Auth / loading
    user, uiLocale, copy, isLoading, error, handleLogout,
    // Core data
    draft, draftLumiView, reloadDraft: loadDraft,
    // Navigation
    activeDiarySection, setActiveDiarySection: setActiveDiarySectionGuarded, activeSection,
    isNarrowViewport, isPlanCompleted, isLearningStarted, isCourseCompleted,
    // Draft meta
    draftTitleInput, setDraftTitleInput, draftDescriptionInput, setDraftDescriptionInput,
    isSavingDraftMeta, saveMessage, hasDraftMetaChanges, handleSaveDraftMeta,
    // Structure (from useDraftStructure)
    ...structure,
    // Selection
    planningSummary, selectedDetailKey, selectedPointContext,
    handleSelectRegion, handleSelectPoint, handleSelectResearchNode,
    // Detail edit (from useDraftDetailEdit)
    ...detailEdit,
    // Resources (from useDraftResources)
    selectedPointResources, resourceReadinessSummary,
    ...resources,
    // Journal / Record / Artifact (from useDraftEntries)
    ...entries,
    // Cross-domain computed
    hasEveryRegionPoint, hasEveryPointSelectedResource, hasEveryResearchNodeTitle,
    hasBlockingUnsavedChanges, hasResourceCoverage, hasSelectedResourceCoverage,
    pointsMissingSelectedResources,
  };
}

export type DraftEditorState = ReturnType<typeof useDraftEditor>;
