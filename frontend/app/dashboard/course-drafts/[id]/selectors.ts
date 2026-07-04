'use client';

import type {
  DraftAggregate,
  DraftLessonTree,
  ManualResourceContentType,
} from './types';

// ── Selectors ─────────────────────────────────────────────────────────────

export const extractSavedDetailNotes = (draft: DraftAggregate): Record<string, string> => {
  const notes: Record<string, string> = {};
  for (const main of draft.lessons) {
    notes[`region:${main.lesson.id}`] = main.lesson.operation_note ?? '';
    for (const sub of (main.sub_lessons ?? [])) {
      notes[`point:${sub.lesson.id}`] = sub.lesson.operation_note ?? '';
    }
  }
  return notes;
};

export const extractSavedJournalEntries = (
  draft: DraftAggregate,
): Record<string, { observation: string; reflection: string; next_step: string }> => {
  const entries: Record<string, { observation: string; reflection: string; next_step: string }> = {};
  for (const main of draft.lessons) {
    for (const sub of (main.sub_lessons ?? [])) {
      entries[`point:${sub.lesson.id}`] = {
        observation: sub.lesson.journal_entry?.observation ?? '',
        reflection: sub.lesson.journal_entry?.reflection ?? '',
        next_step: sub.lesson.journal_entry?.next_step ?? '',
      };
    }
  }
  return entries;
};

export const extractSavedRecordEntries = (
  draft: DraftAggregate,
): Record<string, { study_minutes: number; practice_count: number; confidence_level: number; application_note: string }> => {
  const entries: Record<string, { study_minutes: number; practice_count: number; confidence_level: number; application_note: string }> = {};
  for (const main of draft.lessons) {
    for (const sub of (main.sub_lessons ?? [])) {
      entries[`point:${sub.lesson.id}`] = {
        study_minutes: sub.lesson.record_entry?.study_minutes ?? 0,
        practice_count: sub.lesson.record_entry?.practice_count ?? 0,
        confidence_level: sub.lesson.record_entry?.confidence_level ?? 3,
        application_note: sub.lesson.record_entry?.application_note ?? '',
      };
    }
  }
  return entries;
};

export const extractSavedArtifactEntries = (
  draft: DraftAggregate,
): Record<string, { artifact_type: string; title: string; url: string; description: string }> => {
  const entries: Record<string, { artifact_type: string; title: string; url: string; description: string }> = {};
  for (const main of draft.lessons) {
    for (const sub of (main.sub_lessons ?? [])) {
      entries[`point:${sub.lesson.id}`] = {
        artifact_type: sub.lesson.artifact_entry?.artifact_type ?? 'note',
        title: sub.lesson.artifact_entry?.title ?? '',
        url: sub.lesson.artifact_entry?.url ?? '',
        description: sub.lesson.artifact_entry?.description ?? '',
      };
    }
  }
  return entries;
};

export const getResourceSelectionStateLabel = (selectionState: string) => {
  switch (selectionState) {
    case 'selected':
      return '선택됨';
    case 'candidate':
      return '추천 후보';
    case 'rejected':
      return '제외됨';
    default:
      return '상태 확인 필요';
  }
};

export const normalizeManualResourceContentType = (
  contentType: string | undefined,
  sourceUrl: string,
): ManualResourceContentType => {
  if (contentType === 'youtube' || contentType === 'blog' || contentType === 'article' || contentType === 'internal') {
    return contentType;
  }

  if (/youtube\.com|youtu\.be/i.test(sourceUrl)) return 'youtube';
  return 'article';
};

/** 탐험지점(exploration) DraftPointAggregate 목록을 sub_lessons에서 추출 */
export function getExplorationPoints(subTree: DraftLessonTree) {
  return (subTree.points ?? []).filter((p) => p.point.point_type === 'exploration');
}

/** 연구지점(research) DraftPointAggregate 목록을 mainTree에서 추출 */
export function getResearchPoints(mainTree: DraftLessonTree) {
  return (mainTree.points ?? []).filter((p) => p.point.point_type === 'research');
}
