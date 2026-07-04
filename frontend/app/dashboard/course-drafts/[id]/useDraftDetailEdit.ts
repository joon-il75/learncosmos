'use client';

import { type Dispatch, type SetStateAction, useEffect, useMemo, useState } from 'react';
import type { DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft';
import type {
  DraftAggregate, DraftLessonTree, DraftPointAggregate, DraftPoint,
  ResearchNodeTemplateType, DetailTargetKind,
  UpdateDraftDetailMemoResponse, UpdateDraftLessonResponse,
  UpdateDraftLevelResponse, UpdateResearchNodeResponse,
  CreateResearchNodeResponse, DeleteResearchNodeResponse,
} from './types';
import { extractSavedDetailNotes, getResearchPoints } from './selectors';
import { researchNodeTemplateOptions } from './constants';

type PlanningSummary = {
  selectedLevel: DraftLessonTree | null;
  selectedLesson: DraftLessonTree | null;
  selectedResearchNode: DraftPointAggregate | null;
};

interface Input {
  draft: DraftAggregate | null;
  setDraft: Dispatch<SetStateAction<DraftAggregate | null>>;
  draftId: string;
  planningSummary: PlanningSummary;
  selectedDetailKey: string;
  hasStructureDraftChanges: boolean;
  setSelectedRegionId: Dispatch<SetStateAction<string | null>>;
  setSelectedPointId: Dispatch<SetStateAction<string | null>>;
  setSelectedResearchNodeId: Dispatch<SetStateAction<string | null>>;
  setDetailTargetKind: Dispatch<SetStateAction<DetailTargetKind>>;
  copy: DashboardCourseDraftCopy['hook']['detail'];
}

export function useDraftDetailEdit({
  draft, setDraft, draftId, planningSummary, selectedDetailKey,
  hasStructureDraftChanges,
  setSelectedRegionId, setSelectedPointId, setSelectedResearchNodeId, setDetailTargetKind,
  copy,
}: Input) {
  const [savedDetailNotes, setSavedDetailNotes] = useState<Record<string, string>>({});
  const [detailTitleInput, setDetailTitleInput] = useState('');
  const [detailDescriptionInput, setDetailDescriptionInput] = useState('');
  const [detailDifficultyInput, setDetailDifficultyInput] = useState('medium');
  const [researchNodeTemplateInput, setResearchNodeTemplateInput] = useState<ResearchNodeTemplateType>('free_research');
  const [detailNoteInput, setDetailNoteInput] = useState('');
  const [isSavingPointDetail, setIsSavingPointDetail] = useState(false);
  const [pointDetailSaveMessage, setPointDetailSaveMessage] = useState<string | null>(null);
  const [isSavingRegionDetail, setIsSavingRegionDetail] = useState(false);
  const [regionDetailSaveMessage, setRegionDetailSaveMessage] = useState<string | null>(null);
  const [isSavingResearchNodeDetail, setIsSavingResearchNodeDetail] = useState(false);
  const [researchNodeDetailSaveMessage, setResearchNodeDetailSaveMessage] = useState<string | null>(null);
  const [isCreatingResearchNode, setIsCreatingResearchNode] = useState(false);
  const [isCreatingSubLesson, setIsCreatingSubLesson] = useState(false);
  const [subLessonActionMessage, setSubLessonActionMessage] = useState<string | null>(null);
  const [isDeletingResearchNodeId, setIsDeletingResearchNodeId] = useState<string | null>(null);
  const [researchNodeActionMessage, setResearchNodeActionMessage] = useState<string | null>(null);

  // 캐시 동기화: draft 변경 시 메모 캐시 갱신
  useEffect(() => {
    if (!draft) return;
    setSavedDetailNotes(extractSavedDetailNotes(draft));
  }, [draft]);

  // selectedDetailSource: 선택 항목의 서버 저장 값
  const selectedDetailSource = useMemo(() => ({
    title: planningSummary.selectedLesson?.lesson.title
      ?? planningSummary.selectedResearchNode?.point.title
      ?? planningSummary.selectedLevel?.lesson.title ?? '',
    description: planningSummary.selectedLesson?.lesson.summary
      ?? (planningSummary.selectedResearchNode
        ? copy.templateDescription(planningSummary.selectedResearchNode?.point.template_type)
        : planningSummary.selectedLevel?.lesson.summary) ?? '',
    difficulty: planningSummary.selectedLesson?.lesson.difficulty_level ?? 'medium',
    templateType: planningSummary.selectedResearchNode?.point.template_type ?? 'free_research',
  }), [copy, planningSummary.selectedLesson, planningSummary.selectedLevel, planningSummary.selectedResearchNode]);

  const selectedDetailSavedNote = savedDetailNotes[selectedDetailKey] ?? '';

  // 선택 항목 변경 시 입력 값 초기화
  useEffect(() => {
    setDetailTitleInput(selectedDetailSource.title);
    setDetailDescriptionInput(selectedDetailSource.description);
    setDetailDifficultyInput(selectedDetailSource.difficulty);
    setResearchNodeTemplateInput(selectedDetailSource.templateType);
    setDetailNoteInput(selectedDetailSavedNote);
  }, [selectedDetailKey, selectedDetailSavedNote, selectedDetailSource.title, selectedDetailSource.description, selectedDetailSource.difficulty, selectedDetailSource.templateType]);

  // 선택 항목 변경 시 저장 메시지 초기화
  useEffect(() => {
    setPointDetailSaveMessage(null);
    setRegionDetailSaveMessage(null);
    setResearchNodeDetailSaveMessage(null);
  }, [selectedDetailKey]);

  // ── computed ──────────────────────────────────────────────────────────────

  const hasDetailDraftChanges = Boolean(
    detailTitleInput !== selectedDetailSource.title
    || (!planningSummary.selectedResearchNode && detailDescriptionInput !== selectedDetailSource.description)
    || (planningSummary.selectedLesson && detailDifficultyInput !== selectedDetailSource.difficulty)
    || (planningSummary.selectedResearchNode && researchNodeTemplateInput !== selectedDetailSource.templateType)
    || detailNoteInput.trim() !== selectedDetailSavedNote,
  );

  const hasPointDetailSaveChanges = Boolean(
    planningSummary.selectedLesson
    && (detailTitleInput.trim() !== selectedDetailSource.title || detailDescriptionInput.trim() !== selectedDetailSource.description || detailDifficultyInput !== selectedDetailSource.difficulty),
  );
  const hasPointMemoSaveChanges = Boolean(planningSummary.selectedLesson && detailNoteInput.trim() !== selectedDetailSavedNote);
  const canSavePointDetail = Boolean(draft && planningSummary.selectedLesson && (hasPointDetailSaveChanges || hasPointMemoSaveChanges) && !isSavingPointDetail);

  const hasRegionDetailSaveChanges = Boolean(
    planningSummary.selectedLevel && !planningSummary.selectedLesson
    && (detailTitleInput.trim() !== selectedDetailSource.title || detailDescriptionInput.trim() !== selectedDetailSource.description),
  );
  const hasRegionMemoSaveChanges = Boolean(planningSummary.selectedLevel && !planningSummary.selectedLesson && detailNoteInput.trim() !== selectedDetailSavedNote);
  const canSaveRegionDetail = Boolean(draft && planningSummary.selectedLevel && !planningSummary.selectedLesson && (hasRegionDetailSaveChanges || hasRegionMemoSaveChanges) && !isSavingRegionDetail);

  const hasResearchNodeDetailSaveChanges = Boolean(
    planningSummary.selectedResearchNode
    && (detailTitleInput.trim() !== selectedDetailSource.title || researchNodeTemplateInput !== selectedDetailSource.templateType),
  );
  const canSaveResearchNodeDetail = Boolean(draft && planningSummary.selectedResearchNode && hasResearchNodeDetailSaveChanges && !isSavingResearchNodeDetail);

  const canSaveSelectedDetail = planningSummary.selectedResearchNode ? canSaveResearchNodeDetail : planningSummary.selectedLesson ? canSavePointDetail : canSaveRegionDetail;
  const selectedDetailSaveLabel = planningSummary.selectedLesson
    ? isSavingPointDetail ? copy.savingPoint : copy.savePoint
    : planningSummary.selectedResearchNode
      ? isSavingResearchNodeDetail ? copy.savingResearch : copy.saveResearch
      : isSavingRegionDetail ? copy.savingRegion : copy.saveRegion;

  // ── internal helpers ──────────────────────────────────────────────────────

  const saveDetailMemo = async (target: 'point' | 'region', targetID: string, note: string): Promise<string> => {
    if (!draft) throw new Error(copy.draftMissing);
    const path = target === 'point'
      ? `/api/v1/course-drafts/${draft.draft.id}/lessons/${targetID}/memo`
      : `/api/v1/course-drafts/${draft.draft.id}/lessons/main/${targetID}/memo`;
    const res = await fetch(path, { method: 'PATCH', credentials: 'include', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ note }) });
    const payload = (await res.json().catch(() => ({}))) as UpdateDraftDetailMemoResponse;
    if (!res.ok) {
      if (res.status === 401) throw new Error(copy.authRequired);
      if (res.status === 404) throw new Error(target === 'point' ? copy.memoPointNotFound : copy.memoRegionNotFound);
      if (res.status === 400 || payload.error === 'memo note is required') throw new Error(copy.memoInvalid);
      throw new Error(copy.memoSaveFailed);
    }
    if (!payload.memo) throw new Error(copy.memoResponseMissing);
    return payload.memo.note;
  };

  const syncDraftResearchPoint = (updatedPoint: DraftPoint) => {
    setDraft((cur) => {
      if (!cur) return cur;
      return { ...cur, lessons: cur.lessons.map((t) => ({ ...t, points: (t.points ?? []).map((p) => p.point.id === updatedPoint.id ? { ...p, point: updatedPoint } : p) })) };
    });
  };

  // ── handlers ──────────────────────────────────────────────────────────────

  const handleResetDetailDraft = () => {
    setDetailTitleInput(selectedDetailSource.title);
    setDetailDescriptionInput(selectedDetailSource.description);
    setDetailDifficultyInput(selectedDetailSource.difficulty);
    setResearchNodeTemplateInput(selectedDetailSource.templateType);
    setDetailNoteInput(selectedDetailSavedNote);
    setPointDetailSaveMessage(null);
    setRegionDetailSaveMessage(null);
    setResearchNodeDetailSaveMessage(null);
  };

  const handleSavePointDetail = async () => {
    if (!draft || !planningSummary.selectedLesson || !canSavePointDetail) return;
    const lessonID = planningSummary.selectedLesson.lesson.id;
    const nextTitle = detailTitleInput.trim();
    const nextNote = detailNoteInput.trim();
    if (!nextTitle) { setPointDetailSaveMessage(copy.pointTitleRequired); return; }
    setIsSavingPointDetail(true);
    setPointDetailSaveMessage(null);
    try {
      const savedParts: string[] = [];
      if (hasPointDetailSaveChanges) {
        const res = await fetch(`/api/v1/course-drafts/${draft.draft.id}/lessons/${lessonID}`, {
          method: 'PATCH', credentials: 'include',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ title: nextTitle, summary: detailDescriptionInput.trim(), difficulty_level: detailDifficultyInput }),
        });
        const payload = (await res.json().catch(() => ({}))) as UpdateDraftLessonResponse;
        if (!res.ok) {
          if (res.status === 401) throw new Error(copy.authRequired);
          if (res.status === 404) throw new Error(copy.pointNotFound);
          if (res.status === 400 || payload.error === 'lesson title is required') throw new Error(copy.pointTitleRequired);
          throw new Error(copy.pointSaveFailed);
        }
        if (!payload.draft) throw new Error(copy.pointResponseMissing);
        setDraft(payload.draft);
        savedParts.push(copy.savedTitleDescriptionDifficulty);
      }
      if (hasPointMemoSaveChanges) {
        const savedNote = await saveDetailMemo('point', lessonID, nextNote);
        setSavedDetailNotes((cur) => ({ ...cur, [selectedDetailKey]: savedNote }));
        setDetailNoteInput(savedNote);
        savedParts.push(copy.savedMemo);
      }
      setPointDetailSaveMessage(copy.pointSaved(savedParts.join(', ')));
    } catch (e) {
      setPointDetailSaveMessage(e instanceof Error ? e.message : copy.pointSaveFailed);
    } finally {
      setIsSavingPointDetail(false);
    }
  };

  const handleSaveRegionDetail = async () => {
    if (!draft || !planningSummary.selectedLevel || planningSummary.selectedLesson || !canSaveRegionDetail) return;
    const levelID = planningSummary.selectedLevel.lesson.id;
    const nextTitle = detailTitleInput.trim();
    const nextNote = detailNoteInput.trim();
    if (!nextTitle) { setRegionDetailSaveMessage(copy.regionTitleRequired); return; }
    setIsSavingRegionDetail(true);
    setRegionDetailSaveMessage(null);
    try {
      const savedParts: string[] = [];
      if (hasRegionDetailSaveChanges) {
        const res = await fetch(`/api/v1/course-drafts/${draft.draft.id}/lessons/main/${levelID}`, {
          method: 'PATCH', credentials: 'include',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ title: nextTitle, objective: detailDescriptionInput.trim() }),
        });
        const payload = (await res.json().catch(() => ({}))) as UpdateDraftLevelResponse;
        if (!res.ok) {
          if (res.status === 401) throw new Error(copy.authRequired);
          if (res.status === 404) throw new Error(copy.regionNotFound);
          if (res.status === 400 || payload.error === 'main lesson title is required') throw new Error(copy.regionTitleRequired);
          throw new Error(copy.regionSaveFailed);
        }
        if (!payload.draft) throw new Error(copy.regionResponseMissing);
        setDraft(payload.draft);
        savedParts.push(copy.savedTitleDescription);
      }
      if (hasRegionMemoSaveChanges) {
        const savedNote = await saveDetailMemo('region', levelID, nextNote);
        setSavedDetailNotes((cur) => ({ ...cur, [selectedDetailKey]: savedNote }));
        setDetailNoteInput(savedNote);
        savedParts.push(copy.savedMemo);
      }
      setRegionDetailSaveMessage(copy.regionSaved(savedParts.join(', ')));
    } catch (e) {
      setRegionDetailSaveMessage(e instanceof Error ? e.message : copy.regionSaveFailed);
    } finally {
      setIsSavingRegionDetail(false);
    }
  };

  const handleSaveResearchNodeDetail = async () => {
    if (!draft || !planningSummary.selectedResearchNode || !canSaveResearchNodeDetail) return;
    const nodeID = planningSummary.selectedResearchNode.point.id;
    const nextTitle = detailTitleInput.trim();
    if (!nextTitle) { setResearchNodeDetailSaveMessage(copy.researchTitleRequired); return; }
    setIsSavingResearchNodeDetail(true);
    setResearchNodeDetailSaveMessage(null);
    try {
      const res = await fetch(`/api/v1/course-drafts/${draft.draft.id}/points/${nodeID}`, {
        method: 'PATCH', credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title: nextTitle, template_type: researchNodeTemplateInput }),
      });
      const payload = (await res.json().catch(() => ({}))) as UpdateResearchNodeResponse;
      if (!res.ok) {
        if (res.status === 401) throw new Error(copy.authRequired);
        if (res.status === 404) throw new Error(copy.researchNotFound);
        if (res.status === 400) throw new Error(copy.researchInvalid);
        throw new Error(copy.researchSaveFailed);
      }
      const savedPoint = payload.point ?? payload.research_node;
      if (!savedPoint) throw new Error(copy.researchResponseMissing);
      syncDraftResearchPoint(savedPoint);
      setDetailTitleInput(savedPoint.title);
      setResearchNodeTemplateInput(savedPoint.template_type ?? 'free_research');
      setResearchNodeDetailSaveMessage(copy.researchSaved);
    } catch (e) {
      setResearchNodeDetailSaveMessage(e instanceof Error ? e.message : copy.researchSaveFailed);
    } finally {
      setIsSavingResearchNodeDetail(false);
    }
  };

  const handleCreateResearchNode = async (mainTree: DraftLessonTree, templateType: ResearchNodeTemplateType = 'free_research', customTitle?: string) => {
    if (!draft || isCreatingResearchNode) return;
    if (hasStructureDraftChanges) { setResearchNodeActionMessage(copy.structureSaveFirstAdd); return; }
    const safeTemplateType = researchNodeTemplateOptions.some((o) => o.value === templateType) ? templateType : 'free_research';
    const existingCount = getResearchPoints(mainTree).length;
    const title = customTitle?.trim() || copy.defaultResearchTitle(safeTemplateType, existingCount + 1);
    setIsCreatingResearchNode(true);
    setResearchNodeActionMessage(null);
    try {
      const res = await fetch(`/api/v1/course-drafts/${draft.draft.id}/lessons/main/${mainTree.lesson.id}/points`, {
        method: 'POST', credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, template_type: safeTemplateType, order_index: existingCount }),
      });
      const payload = (await res.json().catch(() => ({}))) as CreateResearchNodeResponse;
      if (!res.ok) {
        if (res.status === 401) throw new Error(copy.authRequired);
        if (res.status === 404) throw new Error(copy.researchAddNotFound);
        if (res.status === 400) throw new Error(copy.researchInvalid);
        throw new Error(copy.researchAddFailed);
      }
      const created = payload.point ?? payload.research_node;
      if (!created) throw new Error(copy.researchAddResponseMissing);
      setDraft((cur) => {
        if (!cur) return cur;
        return { ...cur, lessons: cur.lessons.map((t) => t.lesson.id === mainTree.lesson.id ? { ...t, points: [...(t.points ?? []), { point: created }] } : t) };
      });
      setSelectedRegionId(mainTree.lesson.id);
      setSelectedPointId(null);
      setSelectedResearchNodeId(created.id);
      setDetailTargetKind('research_node');
      setResearchNodeActionMessage(copy.researchAddSuccess);
    } catch (e) {
      setResearchNodeActionMessage(e instanceof Error ? e.message : copy.researchAddFailed);
    } finally {
      setIsCreatingResearchNode(false);
    }
  };

  const handleDeleteResearchNode = async (mainTree: DraftLessonTree, point: DraftPointAggregate) => {
    if (!draft || isDeletingResearchNodeId) return;
    if (hasStructureDraftChanges) { setResearchNodeActionMessage(copy.structureSaveFirstDelete); return; }
    setIsDeletingResearchNodeId(point.point.id);
    setResearchNodeActionMessage(null);
    try {
      const res = await fetch(`/api/v1/course-drafts/${draft.draft.id}/points/${point.point.id}`, { method: 'DELETE', credentials: 'include' });
      const payload = (await res.json().catch(() => ({}))) as DeleteResearchNodeResponse;
      if (!res.ok) {
        if (res.status === 401) throw new Error(copy.authRequired);
        if (res.status === 404) throw new Error(copy.researchDeleteNotFound);
        throw new Error(payload.error || copy.researchDeleteFailed);
      }
      setDraft((cur) => {
        if (!cur) return cur;
        return { ...cur, lessons: cur.lessons.map((t) => t.lesson.id === mainTree.lesson.id ? { ...t, points: (t.points ?? []).filter((p) => p.point.id !== point.point.id) } : t) };
      });
      setSelectedRegionId(mainTree.lesson.id);
      setSelectedPointId(null);
      setSelectedResearchNodeId(null);
      setDetailTargetKind('region');
      setResearchNodeActionMessage(copy.researchDeleteSuccess);
    } catch (e) {
      setResearchNodeActionMessage(e instanceof Error ? e.message : copy.researchDeleteFailed);
    } finally {
      setIsDeletingResearchNodeId(null);
    }
  };

  const handleCreateSubLesson = async (mainTree: DraftLessonTree, title: string) => {
    if (!draft || isCreatingSubLesson) return;
    const trimmedTitle = title.trim();
    if (!trimmedTitle) { setSubLessonActionMessage(copy.subregionTitleRequired); return; }
    setIsCreatingSubLesson(true);
    setSubLessonActionMessage(null);
    try {
      const res = await fetch(`/api/v1/course-drafts/${draft.draft.id}/lessons/main/${mainTree.lesson.id}/sub`, {
        method: 'POST', credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title: trimmedTitle, order_index: mainTree.sub_lessons?.length ?? 0 }),
      });
      const payload = (await res.json().catch(() => ({}))) as { lesson?: DraftLessonTree['lesson'] };
      if (!res.ok) {
        if (res.status === 401) throw new Error(copy.authRequired);
        if (res.status === 404) throw new Error(copy.subregionAddNotFound);
        if (res.status === 400) throw new Error(copy.subregionInvalid);
        throw new Error(copy.subregionAddFailed);
      }
      const newLesson = payload.lesson;
      if (!newLesson) throw new Error(copy.subregionResponseMissing);
      setDraft((cur) => {
        if (!cur) return cur;
        return { ...cur, lessons: cur.lessons.map((t) => t.lesson.id === mainTree.lesson.id ? { ...t, sub_lessons: [...(t.sub_lessons ?? []), { lesson: newLesson, points: [], sub_lessons: [] }] } : t) };
      });
      setSubLessonActionMessage(copy.subregionAddSuccess);
    } catch (e) {
      setSubLessonActionMessage(e instanceof Error ? e.message : copy.subregionAddFailed);
    } finally {
      setIsCreatingSubLesson(false);
    }
  };

  return {
    detailTitleInput, setDetailTitleInput,
    detailDescriptionInput, setDetailDescriptionInput,
    detailDifficultyInput, setDetailDifficultyInput,
    researchNodeTemplateInput, setResearchNodeTemplateInput,
    detailNoteInput, setDetailNoteInput,
    hasDetailDraftChanges, canSaveSelectedDetail, selectedDetailSaveLabel,
    isSavingPointDetail, pointDetailSaveMessage,
    isSavingRegionDetail, regionDetailSaveMessage,
    isSavingResearchNodeDetail, researchNodeDetailSaveMessage,
    isCreatingResearchNode, isDeletingResearchNodeId, researchNodeActionMessage,
    isCreatingSubLesson, subLessonActionMessage,
    handleResetDetailDraft, handleSavePointDetail, handleSaveRegionDetail, handleSaveResearchNodeDetail,
    handleCreateResearchNode, handleDeleteResearchNode, handleCreateSubLesson,
  };
}
