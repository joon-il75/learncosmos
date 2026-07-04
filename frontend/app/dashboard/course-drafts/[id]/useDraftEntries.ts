'use client';

import { type Dispatch, type SetStateAction, useEffect, useState } from 'react';
import type {
  DraftAggregate, DraftLessonTree,
  UpdateDraftLessonJournalResponse, UpdateDraftLessonRecordResponse, UpdateDraftLessonArtifactResponse,
} from './types';
import type { DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft';
import { extractSavedJournalEntries, extractSavedRecordEntries, extractSavedArtifactEntries } from './selectors';

interface Input {
  draft: DraftAggregate | null;
  setDraft: Dispatch<SetStateAction<DraftAggregate | null>>;
  draftId: string;
  selectedLesson: DraftLessonTree | null;
  selectedDetailKey: string;
  copy: DashboardCourseDraftCopy['hook']['entries'];
}

export function useDraftEntries({ draft, setDraft, draftId, selectedLesson, selectedDetailKey, copy }: Input) {
  const [savedJournalEntries, setSavedJournalEntries] = useState<Record<string, { observation: string; reflection: string; next_step: string }>>({});
  const [journalObservationInput, setJournalObservationInput] = useState('');
  const [journalReflectionInput, setJournalReflectionInput] = useState('');
  const [journalNextStepInput, setJournalNextStepInput] = useState('');
  const [isSavingJournalEntry, setIsSavingJournalEntry] = useState(false);
  const [journalSaveMessage, setJournalSaveMessage] = useState<string | null>(null);

  const [savedRecordEntries, setSavedRecordEntries] = useState<Record<string, { study_minutes: number; practice_count: number; confidence_level: number; application_note: string }>>({});
  const [recordStudyMinutesInput, setRecordStudyMinutesInput] = useState('0');
  const [recordPracticeCountInput, setRecordPracticeCountInput] = useState('0');
  const [recordConfidenceLevelInput, setRecordConfidenceLevelInput] = useState(3);
  const [recordApplicationNoteInput, setRecordApplicationNoteInput] = useState('');
  const [isSavingRecordEntry, setIsSavingRecordEntry] = useState(false);
  const [recordSaveMessage, setRecordSaveMessage] = useState<string | null>(null);

  const [savedArtifactEntries, setSavedArtifactEntries] = useState<Record<string, { artifact_type: string; title: string; url: string; description: string }>>({});
  const [artifactTypeInput, setArtifactTypeInput] = useState('note');
  const [artifactTitleInput, setArtifactTitleInput] = useState('');
  const [artifactUrlInput, setArtifactUrlInput] = useState('');
  const [artifactDescriptionInput, setArtifactDescriptionInput] = useState('');
  const [isSavingArtifactEntry, setIsSavingArtifactEntry] = useState(false);
  const [artifactSaveMessage, setArtifactSaveMessage] = useState<string | null>(null);

  // 캐시 동기화: draft 변경 시 서버 저장값으로 갱신
  useEffect(() => {
    if (!draft) return;
    setSavedJournalEntries(extractSavedJournalEntries(draft));
    setSavedRecordEntries(extractSavedRecordEntries(draft));
    setSavedArtifactEntries(extractSavedArtifactEntries(draft));
  }, [draft]);

  const selectedJournalEntry = savedJournalEntries[selectedDetailKey] ?? { observation: '', reflection: '', next_step: '' };
  const selectedRecordEntry = savedRecordEntries[selectedDetailKey] ?? { study_minutes: 0, practice_count: 0, confidence_level: 3, application_note: '' };
  const selectedArtifactEntry = savedArtifactEntries[selectedDetailKey] ?? { artifact_type: 'note', title: '', url: '', description: '' };

  useEffect(() => {
    setJournalObservationInput(selectedJournalEntry.observation);
    setJournalReflectionInput(selectedJournalEntry.reflection);
    setJournalNextStepInput(selectedJournalEntry.next_step);
  }, [selectedDetailKey, selectedJournalEntry.observation, selectedJournalEntry.reflection, selectedJournalEntry.next_step]);

  useEffect(() => {
    setRecordStudyMinutesInput(String(selectedRecordEntry.study_minutes));
    setRecordPracticeCountInput(String(selectedRecordEntry.practice_count));
    setRecordConfidenceLevelInput(selectedRecordEntry.confidence_level);
    setRecordApplicationNoteInput(selectedRecordEntry.application_note);
  }, [selectedDetailKey, selectedRecordEntry.study_minutes, selectedRecordEntry.practice_count, selectedRecordEntry.confidence_level, selectedRecordEntry.application_note]);

  useEffect(() => {
    setArtifactTypeInput(selectedArtifactEntry.artifact_type);
    setArtifactTitleInput(selectedArtifactEntry.title);
    setArtifactUrlInput(selectedArtifactEntry.url);
    setArtifactDescriptionInput(selectedArtifactEntry.description);
  }, [selectedDetailKey, selectedArtifactEntry.artifact_type, selectedArtifactEntry.title, selectedArtifactEntry.url, selectedArtifactEntry.description]);

  // ── computed ──────────────────────────────────────────────────────────────

  const hasJournalDraftChanges = Boolean(
    selectedLesson && (
      journalObservationInput.trim() !== selectedJournalEntry.observation
      || journalReflectionInput.trim() !== selectedJournalEntry.reflection
      || journalNextStepInput.trim() !== selectedJournalEntry.next_step
    ),
  );
  const canSaveJournalEntry = Boolean(
    draft && selectedLesson && hasJournalDraftChanges && !isSavingJournalEntry
    && (journalObservationInput.trim() || journalReflectionInput.trim() || journalNextStepInput.trim()),
  );

  const parsedStudyMinutes = Number.parseInt(recordStudyMinutesInput, 10);
  const parsedPracticeCount = Number.parseInt(recordPracticeCountInput, 10);
  const normalizedStudyMinutes = Number.isNaN(parsedStudyMinutes) ? 0 : parsedStudyMinutes;
  const normalizedPracticeCount = Number.isNaN(parsedPracticeCount) ? 0 : parsedPracticeCount;
  const hasRecordDraftChanges = Boolean(
    selectedLesson && (
      normalizedStudyMinutes !== selectedRecordEntry.study_minutes
      || normalizedPracticeCount !== selectedRecordEntry.practice_count
      || recordConfidenceLevelInput !== selectedRecordEntry.confidence_level
      || recordApplicationNoteInput.trim() !== selectedRecordEntry.application_note
    ),
  );
  const canSaveRecordEntry = Boolean(
    draft && selectedLesson && !isSavingRecordEntry && hasRecordDraftChanges
    && normalizedStudyMinutes >= 0 && normalizedPracticeCount >= 0
    && recordConfidenceLevelInput >= 1 && recordConfidenceLevelInput <= 5
    && (normalizedStudyMinutes > 0 || normalizedPracticeCount > 0 || recordConfidenceLevelInput !== 3 || recordApplicationNoteInput.trim()),
  );

  const trimmedArtifactTitle = artifactTitleInput.trim();
  const trimmedArtifactUrl = artifactUrlInput.trim();
  const trimmedArtifactDescription = artifactDescriptionInput.trim();
  const isArtifactUrlValid = !trimmedArtifactUrl || /^https?:\/\/[^\s]+$/i.test(trimmedArtifactUrl);
  const hasArtifactDraftChanges = Boolean(
    selectedLesson && (
      artifactTypeInput !== selectedArtifactEntry.artifact_type
      || trimmedArtifactTitle !== selectedArtifactEntry.title
      || trimmedArtifactUrl !== selectedArtifactEntry.url
      || trimmedArtifactDescription !== selectedArtifactEntry.description
    ),
  );
  const canSaveArtifactEntry = Boolean(
    draft && selectedLesson && hasArtifactDraftChanges && !isSavingArtifactEntry
    && trimmedArtifactTitle && isArtifactUrlValid,
  );

  // ── internal sync helpers ─────────────────────────────────────────────────

  const syncDraftLessonJournal = (lessonID: string, entry: { observation: string; reflection: string; next_step: string; updated_at: string }) => {
    setDraft((cur) => {
      if (!cur) return cur;
      return { ...cur, lessons: cur.lessons.map((t) => ({ ...t, sub_lessons: (t.sub_lessons ?? []).map((s) => s.lesson.id === lessonID ? { ...s, lesson: { ...s.lesson, journal_entry: entry } } : s) })) };
    });
  };

  const syncDraftLessonRecord = (lessonID: string, entry: { study_minutes: number; practice_count: number; confidence_level: number; application_note: string; updated_at: string }) => {
    setDraft((cur) => {
      if (!cur) return cur;
      return { ...cur, lessons: cur.lessons.map((t) => ({ ...t, sub_lessons: (t.sub_lessons ?? []).map((s) => s.lesson.id === lessonID ? { ...s, lesson: { ...s.lesson, record_entry: entry } } : s) })) };
    });
  };

  const syncDraftLessonArtifact = (lessonID: string, entry: { artifact_type: string; title: string; url: string; description: string; updated_at: string }) => {
    setDraft((cur) => {
      if (!cur) return cur;
      return { ...cur, lessons: cur.lessons.map((t) => ({ ...t, sub_lessons: (t.sub_lessons ?? []).map((s) => s.lesson.id === lessonID ? { ...s, lesson: { ...s.lesson, artifact_entry: entry } } : s) })) };
    });
  };

  // ── handlers ──────────────────────────────────────────────────────────────

  const handleResetJournalDraft = () => {
    setJournalObservationInput(selectedJournalEntry.observation);
    setJournalReflectionInput(selectedJournalEntry.reflection);
    setJournalNextStepInput(selectedJournalEntry.next_step);
    setJournalSaveMessage(null);
  };

  const handleSaveJournalEntry = async () => {
    if (!draft || !selectedLesson || !canSaveJournalEntry) return;
    const lessonID = selectedLesson.lesson.id;
    const detailKey = `point:${lessonID}`;
    setIsSavingJournalEntry(true);
    setJournalSaveMessage(null);
    try {
      const res = await fetch(`/api/v1/course-drafts/${draft.draft.id}/lessons/${lessonID}/journal`, {
        method: 'PATCH', credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ observation: journalObservationInput.trim(), reflection: journalReflectionInput.trim(), next_step: journalNextStepInput.trim() }),
      });
      const payload = (await res.json().catch(() => ({}))) as UpdateDraftLessonJournalResponse;
      if (!res.ok) {
        if (res.status === 401) throw new Error(copy.authRequired);
        if (res.status === 404) throw new Error(copy.targetNotFound);
        if (res.status === 400 || payload.error === 'journal content is required') throw new Error(copy.journalContentRequired);
        throw new Error(copy.journalSaveFailed);
      }
      if (!payload.journal) throw new Error(copy.journalResponseMissing);
      const saved = { observation: payload.journal.observation, reflection: payload.journal.reflection, next_step: payload.journal.next_step };
      setSavedJournalEntries((cur) => ({ ...cur, [detailKey]: saved }));
      syncDraftLessonJournal(lessonID, { ...saved, updated_at: payload.journal.updated_at });
      setJournalObservationInput(saved.observation);
      setJournalReflectionInput(saved.reflection);
      setJournalNextStepInput(saved.next_step);
      setJournalSaveMessage(copy.journalSaveSuccess);
    } catch (e) {
      setJournalSaveMessage(e instanceof Error ? e.message : copy.journalSaveFailed);
    } finally {
      setIsSavingJournalEntry(false);
    }
  };

  const handleResetRecordDraft = () => {
    setRecordStudyMinutesInput(String(selectedRecordEntry.study_minutes));
    setRecordPracticeCountInput(String(selectedRecordEntry.practice_count));
    setRecordConfidenceLevelInput(selectedRecordEntry.confidence_level);
    setRecordApplicationNoteInput(selectedRecordEntry.application_note);
    setRecordSaveMessage(null);
  };

  const handleSaveRecordEntry = async () => {
    if (!draft || !selectedLesson || !canSaveRecordEntry) return;
    const lessonID = selectedLesson.lesson.id;
    const detailKey = `point:${lessonID}`;
    setIsSavingRecordEntry(true);
    setRecordSaveMessage(null);
    try {
      const res = await fetch(`/api/v1/course-drafts/${draft.draft.id}/lessons/${lessonID}/record`, {
        method: 'PATCH', credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ study_minutes: normalizedStudyMinutes, practice_count: normalizedPracticeCount, confidence_level: recordConfidenceLevelInput, application_note: recordApplicationNoteInput.trim() }),
      });
      const payload = (await res.json().catch(() => ({}))) as UpdateDraftLessonRecordResponse;
      if (!res.ok) {
        if (res.status === 401) throw new Error(copy.authRequired);
        if (res.status === 404) throw new Error(copy.targetNotFound);
        if (res.status === 400 || payload.error === 'record content is invalid') throw new Error(copy.recordInvalid);
        throw new Error(copy.recordSaveFailed);
      }
      if (!payload.record) throw new Error(copy.recordResponseMissing);
      const saved = { study_minutes: payload.record.study_minutes, practice_count: payload.record.practice_count, confidence_level: payload.record.confidence_level, application_note: payload.record.application_note };
      setSavedRecordEntries((cur) => ({ ...cur, [detailKey]: saved }));
      syncDraftLessonRecord(lessonID, { ...saved, updated_at: payload.record.updated_at });
      setRecordStudyMinutesInput(String(saved.study_minutes));
      setRecordPracticeCountInput(String(saved.practice_count));
      setRecordConfidenceLevelInput(saved.confidence_level);
      setRecordApplicationNoteInput(saved.application_note);
      setRecordSaveMessage(copy.recordSaveSuccess);
    } catch (e) {
      setRecordSaveMessage(e instanceof Error ? e.message : copy.recordSaveFailed);
    } finally {
      setIsSavingRecordEntry(false);
    }
  };

  const handleResetArtifactDraft = () => {
    setArtifactTypeInput(selectedArtifactEntry.artifact_type);
    setArtifactTitleInput(selectedArtifactEntry.title);
    setArtifactUrlInput(selectedArtifactEntry.url);
    setArtifactDescriptionInput(selectedArtifactEntry.description);
    setArtifactSaveMessage(null);
  };

  const handleSaveArtifactEntry = async () => {
    if (!draft || !selectedLesson || !canSaveArtifactEntry) return;
    const lessonID = selectedLesson.lesson.id;
    const detailKey = `point:${lessonID}`;
    setIsSavingArtifactEntry(true);
    setArtifactSaveMessage(null);
    try {
      const res = await fetch(`/api/v1/course-drafts/${draft.draft.id}/lessons/${lessonID}/artifact`, {
        method: 'PATCH', credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ artifact_type: artifactTypeInput, title: trimmedArtifactTitle, url: trimmedArtifactUrl, description: trimmedArtifactDescription }),
      });
      const payload = (await res.json().catch(() => ({}))) as UpdateDraftLessonArtifactResponse;
      if (!res.ok) {
        if (res.status === 401) throw new Error(copy.authRequired);
        if (res.status === 404) throw new Error(copy.artifactTargetNotFound);
        if (res.status === 400 || payload.error === 'artifact content is invalid') throw new Error(copy.artifactInvalid);
        throw new Error(copy.artifactSaveFailed);
      }
      if (!payload.artifact) throw new Error(copy.artifactResponseMissing);
      const saved = { artifact_type: payload.artifact.artifact_type, title: payload.artifact.title, url: payload.artifact.url, description: payload.artifact.description };
      setSavedArtifactEntries((cur) => ({ ...cur, [detailKey]: saved }));
      syncDraftLessonArtifact(lessonID, { ...saved, updated_at: payload.artifact.updated_at });
      setArtifactTypeInput(saved.artifact_type);
      setArtifactTitleInput(saved.title);
      setArtifactUrlInput(saved.url);
      setArtifactDescriptionInput(saved.description);
      setArtifactSaveMessage(copy.artifactSaveSuccess);
    } catch (e) {
      setArtifactSaveMessage(e instanceof Error ? e.message : copy.artifactSaveFailed);
    } finally {
      setIsSavingArtifactEntry(false);
    }
  };

  return {
    // journal
    journalObservationInput, setJournalObservationInput,
    journalReflectionInput, setJournalReflectionInput,
    journalNextStepInput, setJournalNextStepInput,
    hasJournalDraftChanges, canSaveJournalEntry, isSavingJournalEntry, journalSaveMessage,
    handleResetJournalDraft, handleSaveJournalEntry,
    // record
    recordStudyMinutesInput, setRecordStudyMinutesInput,
    recordPracticeCountInput, setRecordPracticeCountInput,
    recordConfidenceLevelInput, setRecordConfidenceLevelInput,
    recordApplicationNoteInput, setRecordApplicationNoteInput,
    selectedRecordSavedEntry: selectedRecordEntry,
    hasRecordDraftChanges, canSaveRecordEntry, isSavingRecordEntry, recordSaveMessage,
    handleResetRecordDraft, handleSaveRecordEntry,
    // artifact
    artifactTypeInput, setArtifactTypeInput,
    artifactTitleInput, setArtifactTitleInput,
    artifactUrlInput, setArtifactUrlInput,
    artifactDescriptionInput, setArtifactDescriptionInput,
    selectedArtifactSavedEntry: selectedArtifactEntry,
    hasArtifactDraftChanges, canSaveArtifactEntry, isSavingArtifactEntry, artifactSaveMessage,
    handleResetArtifactDraft, handleSaveArtifactEntry,
  };
}
