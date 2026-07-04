'use client';

import { type Dispatch, type SetStateAction, useEffect, useState } from 'react';
import type {
  DraftAggregate, DraftLessonTree,
  AttachDraftLessonResourceResponse, LessonSearchResponse,
  ManualResourceCreatedContent, ManualResourceParsePreview, SelectDraftResourceResponse,
} from './types';
import type { DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft';
import { getExplorationPoints, normalizeManualResourceContentType } from './selectors';
import { resolveSafetyInputMessage, type SafetyAPIErrorPayload } from '@/lib/safetyErrors';

interface Input {
  draft: DraftAggregate | null;
  setDraft: Dispatch<SetStateAction<DraftAggregate | null>>;
  selectedLesson: DraftLessonTree | null;
  selectedDetailKey: string;
  copy: DashboardCourseDraftCopy['hook']['resources'];
}

export function useDraftResources({ draft, setDraft, selectedLesson, selectedDetailKey, copy }: Input) {
  const [manualResourceTitleInput, setManualResourceTitleInput] = useState('');
  const [manualResourceUrlInput, setManualResourceUrlInput] = useState('');
  const [manualResourceMemoInput, setManualResourceMemoInput] = useState('');
  const [isRefreshingCandidates, setIsRefreshingCandidates] = useState(false);
  const [candidateRefreshMessage, setCandidateRefreshMessage] = useState<string | null>(null);
  const [selectingResourceId, setSelectingResourceId] = useState<string | null>(null);
  const [isParsingManualResource, setIsParsingManualResource] = useState(false);
  const [manualResourceParsePreview, setManualResourceParsePreview] = useState<ManualResourceParsePreview | null>(null);
  const [manualResourceParseMessage, setManualResourceParseMessage] = useState<string | null>(null);
  const [isCreatingManualContent, setIsCreatingManualContent] = useState(false);
  const [manualResourceCreateMessage, setManualResourceCreateMessage] = useState<string | null>(null);
  const [manualResourceCreatedContent, setManualResourceCreatedContent] = useState<ManualResourceCreatedContent | null>(null);
  const [isAttachingManualResource, setIsAttachingManualResource] = useState(false);
  const [manualResourceAttachMessage, setManualResourceAttachMessage] = useState<string | null>(null);
  const [isManualResourceAttached, setIsManualResourceAttached] = useState(false);

  // 선택 지역 변경 시 manual resource 입력 초기화
  useEffect(() => {
    setManualResourceTitleInput('');
    setManualResourceUrlInput('');
    setManualResourceMemoInput('');
    setManualResourceParsePreview(null);
    setManualResourceParseMessage(null);
    setManualResourceCreateMessage(null);
    setManualResourceCreatedContent(null);
    setManualResourceAttachMessage(null);
    setIsManualResourceAttached(false);
    setCandidateRefreshMessage(null);
  }, [selectedDetailKey]);

  // ── computed ──────────────────────────────────────────────────────────────

  const trimmedManualResourceUrl = manualResourceUrlInput.trim();
  const trimmedManualResourceTitle = manualResourceTitleInput.trim();
  const isManualResourceUrlEmpty = trimmedManualResourceUrl.length === 0;
  const isManualResourceUrlValid = isManualResourceUrlEmpty || /^https?:\/\/[^\s]+$/i.test(trimmedManualResourceUrl);
  const hasManualResourceDraft = Boolean(manualResourceTitleInput.trim() || manualResourceUrlInput.trim() || manualResourceMemoInput.trim());
  const canCreateManualContent = Boolean(trimmedManualResourceTitle && !isManualResourceUrlEmpty && isManualResourceUrlValid && !manualResourceCreatedContent);
  const canAttachManualResource = Boolean(draft && selectedLesson && manualResourceCreatedContent && !isManualResourceAttached);
  const hasUnattachedManualResourceDraft = Boolean(!isManualResourceAttached && (hasManualResourceDraft || manualResourceCreatedContent));
  const manualResourceValidationItems = [
    { label: copy.titleLabel, message: trimmedManualResourceTitle ? copy.titleValid : copy.titleRequired, isValid: Boolean(trimmedManualResourceTitle) },
    { label: copy.urlLabel, message: isManualResourceUrlValid ? copy.urlValid : copy.urlInvalid, isValid: isManualResourceUrlValid },
    { label: copy.memoLabel, message: manualResourceMemoInput.trim() ? copy.memoValid : copy.memoOptional, isValid: true },
  ] as const;

  // ── handlers ──────────────────────────────────────────────────────────────

  const handleManualResourceTitleChange = (value: string) => {
    setManualResourceTitleInput(value);
    setManualResourceCreateMessage(null);
    setManualResourceCreatedContent(null);
    setManualResourceAttachMessage(null);
    setIsManualResourceAttached(false);
  };

  const handleManualResourceUrlChange = (value: string) => {
    setManualResourceUrlInput(value);
    setManualResourceParsePreview(null);
    setManualResourceParseMessage(null);
    setManualResourceCreateMessage(null);
    setManualResourceCreatedContent(null);
    setManualResourceAttachMessage(null);
    setIsManualResourceAttached(false);
  };

  const handleManualResourceMemoChange = (value: string) => {
    setManualResourceMemoInput(value);
    setManualResourceCreateMessage(null);
    setManualResourceCreatedContent(null);
    setManualResourceAttachMessage(null);
    setIsManualResourceAttached(false);
  };

  const handleResetManualResourceDraft = () => {
    setManualResourceTitleInput('');
    setManualResourceUrlInput('');
    setManualResourceMemoInput('');
    setManualResourceParsePreview(null);
    setManualResourceParseMessage(null);
    setManualResourceCreateMessage(null);
    setManualResourceCreatedContent(null);
    setManualResourceAttachMessage(null);
    setIsManualResourceAttached(false);
  };

  const handleParseManualResource = async () => {
    if (isParsingManualResource || manualResourceCreatedContent || !isManualResourceUrlValid || isManualResourceUrlEmpty) return;
    setIsParsingManualResource(true);
    setManualResourceParsePreview(null);
    setManualResourceParseMessage(null);
    try {
      const res = await fetch('/api/v1/contents/parse', {
        method: 'POST', credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url: trimmedManualResourceUrl }),
      });
      const payload = (await res.json().catch(() => ({}))) as ManualResourceParsePreview & { error?: string };
      if (!res.ok) {
        if (res.status === 401) throw new Error(copy.authRequired);
        throw new Error(copy.parseFailed);
      }
      setManualResourceParsePreview(payload);
      if (payload.title && !manualResourceTitleInput.trim()) setManualResourceTitleInput(payload.title);
      setManualResourceParseMessage(payload.parse_error ? copy.parsePartial : copy.parseSuccess);
    } catch (e) {
      setManualResourceParseMessage(e instanceof Error ? e.message : copy.parseFailedKeep);
    } finally {
      setIsParsingManualResource(false);
    }
  };

  const handleCreateManualContent = async () => {
    if (isCreatingManualContent || !canCreateManualContent) return;
    setIsCreatingManualContent(true);
    setManualResourceCreateMessage(null);
    const sourceUrl = manualResourceParsePreview?.source_url || trimmedManualResourceUrl;
    const description = manualResourceParsePreview?.description || manualResourceMemoInput.trim();
    const durationSeconds = typeof manualResourceParsePreview?.duration_seconds === 'number' && manualResourceParsePreview.duration_seconds > 0 ? manualResourceParsePreview.duration_seconds : undefined;
    try {
      const res = await fetch('/api/v1/contents', {
        method: 'POST', credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content_type: normalizeManualResourceContentType(manualResourceParsePreview?.content_type, sourceUrl), url: sourceUrl, title: trimmedManualResourceTitle, description: description || undefined, thumbnail_url: manualResourceParsePreview?.thumbnail_url || undefined, duration_seconds: durationSeconds, author: manualResourceParsePreview?.author || undefined, language: 'ko' }),
      });
      const payload = (await res.json().catch(() => ({}))) as ManualResourceCreatedContent & SafetyAPIErrorPayload;
      if (!res.ok) {
        if (res.status === 401) throw new Error(copy.authRequired);
        if (payload.error_code === 'safety_input_blocked' || payload.error_code === 'safety_input_soft_warn') throw new Error(resolveSafetyInputMessage(payload, payload.error ?? copy.contentSaveFailed));
        if (res.status === 400) throw new Error(copy.contentInvalid);
        throw new Error(copy.contentSaveFailed);
      }
      if (!payload.id) throw new Error(copy.contentResponseMissing);
      setManualResourceCreatedContent(payload);
      setManualResourceCreateMessage(copy.contentCreated);
      setManualResourceAttachMessage(null);
      setIsManualResourceAttached(false);
    } catch (e) {
      setManualResourceCreateMessage(e instanceof Error ? e.message : copy.contentSaveFailedKeep);
    } finally {
      setIsCreatingManualContent(false);
    }
  };

  const handleAttachManualResourceToPoint = async () => {
    if (!draft || !selectedLesson || !manualResourceCreatedContent || isAttachingManualResource || !canAttachManualResource) return;
    setIsAttachingManualResource(true);
    setManualResourceAttachMessage(null);
    try {
      const res = await fetch(`/api/v1/course-drafts/${draft.draft.id}/lessons/${selectedLesson.lesson.id}/resources`, {
        method: 'POST', credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content_id: manualResourceCreatedContent.id, selection_state: 'selected' }),
      });
      const payload = (await res.json().catch(() => ({}))) as AttachDraftLessonResourceResponse;
      if (!res.ok) {
        if (res.status === 401) throw new Error(copy.authRequired);
        if (res.status === 404) throw new Error(copy.attachNotFound);
        if (res.status === 409) throw new Error(copy.attachConflict);
        if (payload.error === 'invalid resource selection state') throw new Error(copy.attachInvalidState);
        throw new Error(copy.attachFailed);
      }
      if (!payload.draft) throw new Error(copy.attachResponseMissing);
      setDraft(payload.draft);
      setIsManualResourceAttached(true);
      setManualResourceCreateMessage(copy.contentCreatedAndAttached);
      setManualResourceAttachMessage(copy.attachSuccess);
      setCandidateRefreshMessage(copy.attachCandidateMessage);
    } catch (e) {
      setManualResourceAttachMessage(e instanceof Error ? e.message : copy.attachFailedKeep);
    } finally {
      setIsAttachingManualResource(false);
    }
  };

  const refreshLessonCandidates = async (lessonID?: string) => {
    if (!draft || isRefreshingCandidates) return;
    setIsRefreshingCandidates(true);
    setCandidateRefreshMessage(null);
    try {
      const res = await fetch(`/api/v1/course-drafts/${draft.draft.id}/search`, {
        method: 'POST', credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ lesson_id: lessonID, max_per_lesson: 5 }),
      });
      const payload = (await res.json().catch(() => ({}))) as LessonSearchResponse;
      if (!res.ok) {
        if (payload.error === 'insufficient_points') throw new Error(copy.insufficientPoints);
        if (res.status === 401) throw new Error(copy.authRequired);
        throw new Error(copy.candidateRefreshFailed);
      }
      if (!payload.draft) throw new Error(copy.candidateResponseMissing);
      setDraft(payload.draft);
      const allSubs = (payload.draft.lessons ?? []).flatMap((t) => t.sub_lessons ?? []);
      const candidateCount = lessonID
        ? (() => { const sub = allSubs.find((s) => s.lesson.id === lessonID); return sub ? getExplorationPoints(sub).filter((p) => p.point.selection_state === 'candidate').length : 0; })()
        : allSubs.reduce((sum, sub) => sum + getExplorationPoints(sub).filter((p) => p.point.selection_state === 'candidate').length, 0);
      const costText = payload.point_preview ? copy.costPreview(payload.point_preview.cost) : lessonID ? copy.selectedRegionBasis : copy.allRegionBasis;
      setCandidateRefreshMessage(lessonID ? copy.candidateRefreshSelected(costText, candidateCount) : copy.candidateRefreshAll(costText, candidateCount));
    } catch (e) {
      setCandidateRefreshMessage(e instanceof Error ? e.message : copy.candidateRefreshFailedKeep);
    } finally {
      setIsRefreshingCandidates(false);
    }
  };

  const handleRefreshLessonCandidates = async () => {
    if (!selectedLesson) return;
    await refreshLessonCandidates(selectedLesson.lesson.id);
  };
  const handleRefreshLessonCandidatesById = async (lessonId: string) => { await refreshLessonCandidates(lessonId); };
  const handleRefreshAllLessonCandidates = async () => { await refreshLessonCandidates(); };

  const handleSelectCandidateResource = async (resourceID: string) => {
    if (!draft || selectingResourceId) return;
    setSelectingResourceId(resourceID);
    setCandidateRefreshMessage(null);
    try {
      const res = await fetch(`/api/v1/course-drafts/${draft.draft.id}/resources/${resourceID}/select`, { method: 'POST', credentials: 'include' });
      const payload = (await res.json().catch(() => ({}))) as SelectDraftResourceResponse;
      if (!res.ok) {
        if (res.status === 401) throw new Error(copy.authRequired);
        if (res.status === 404) throw new Error(copy.selectNotFound);
        if (payload.error === 'resource is not a selectable candidate') throw new Error(copy.selectCandidateOnly);
        throw new Error(copy.selectFailed);
      }
      if (!payload.draft) throw new Error(copy.selectResponseMissing);
      setDraft(payload.draft);
      setCandidateRefreshMessage(copy.selectSuccess);
    } catch (e) {
      setCandidateRefreshMessage(e instanceof Error ? e.message : copy.selectFailedKeep);
    } finally {
      setSelectingResourceId(null);
    }
  };

  return {
    manualResourceTitleInput, manualResourceUrlInput, manualResourceMemoInput,
    hasManualResourceDraft, isManualResourceUrlEmpty, isManualResourceUrlValid,
    canCreateManualContent, canAttachManualResource, manualResourceValidationItems,
    hasUnattachedManualResourceDraft,
    isParsingManualResource, manualResourceParsePreview, manualResourceParseMessage,
    isCreatingManualContent, manualResourceCreateMessage,
    manualResourceCreatedContent, isAttachingManualResource,
    manualResourceAttachMessage, isManualResourceAttached,
    isRefreshingCandidates, candidateRefreshMessage, selectingResourceId,
    handleManualResourceTitleChange, handleManualResourceUrlChange, handleManualResourceMemoChange,
    handleParseManualResource, handleCreateManualContent,
    handleAttachManualResourceToPoint, handleResetManualResourceDraft,
    handleRefreshAllLessonCandidates, handleRefreshLessonCandidates,
    handleRefreshLessonCandidatesById, handleSelectCandidateResource,
  };
}
