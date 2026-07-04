'use client';

import { type Dispatch, type SetStateAction, useState } from 'react';
import type { DetailTargetKind, DraftAggregate, DraftLessonTree, UpdateDraftStructureResponse } from './types';
import { getResearchPoints } from './selectors';

interface Input {
  draft: DraftAggregate | null;
  setDraft: Dispatch<SetStateAction<DraftAggregate | null>>;
  draftId: string;
  setSelectedRegionId: Dispatch<SetStateAction<string | null>>;
  setSelectedPointId: Dispatch<SetStateAction<string | null>>;
  setSelectedResearchNodeId: Dispatch<SetStateAction<string | null>>;
  setDetailTargetKind: Dispatch<SetStateAction<DetailTargetKind>>;
}

export function useDraftStructure({
  draft, setDraft,
  setSelectedRegionId, setSelectedPointId, setSelectedResearchNodeId, setDetailTargetKind,
}: Input) {
  const [hasStructureDraftChanges, setHasStructureDraftChanges] = useState(false);
  const [isSavingDraftStructure, setIsSavingDraftStructure] = useState(false);
  const [structureSaveMessage, setStructureSaveMessage] = useState<string | null>(null);

  const canSaveDraftStructure = Boolean(
    draft && draft.lessons.length > 0 && hasStructureDraftChanges && !isSavingDraftStructure,
  );

  const handleMoveRegion = (levelID: string, direction: -1 | 1) => {
    if (isSavingDraftStructure) return;
    setDraft((cur) => {
      if (!cur) return cur;
      const idx = cur.lessons.findIndex((t) => t.lesson.id === levelID);
      const next = idx + direction;
      if (idx < 0 || next < 0 || next >= cur.lessons.length) return cur;
      const arr = [...cur.lessons];
      [arr[idx], arr[next]] = [arr[next], arr[idx]];
      return { ...cur, lessons: arr };
    });
    setHasStructureDraftChanges(true);
    setStructureSaveMessage('지역 순서를 조정했습니다. 탐험 구조 저장을 눌러 반영해 주세요.');
  };

  const handleMovePoint = (levelID: string, lessonID: string, direction: -1 | 1) => {
    if (isSavingDraftStructure) return;
    setDraft((cur) => {
      if (!cur) return cur;
      const nextLessons = cur.lessons.map((tree) => {
        if (tree.lesson.id !== levelID) return tree;
        const subs = tree.sub_lessons ?? [];
        const idx = subs.findIndex((s) => s.lesson.id === lessonID);
        const next = idx + direction;
        if (idx < 0 || next < 0 || next >= subs.length) return tree;
        const arr = [...subs];
        [arr[idx], arr[next]] = [arr[next], arr[idx]];
        return { ...tree, sub_lessons: arr };
      });
      return { ...cur, lessons: nextLessons };
    });
    setHasStructureDraftChanges(true);
    setStructureSaveMessage('지역 순서를 조정했습니다. 탐험 구조 저장을 눌러 반영해 주세요.');
  };

  const handleMovePointToAdjacentRegion = (levelID: string, lessonID: string, direction: -1 | 1) => {
    if (isSavingDraftStructure) return;
    let targetLevelID: string | null = null;
    setDraft((cur) => {
      if (!cur) return cur;
      const srcIdx = cur.lessons.findIndex((t) => t.lesson.id === levelID);
      const tgtIdx = srcIdx + direction;
      if (srcIdx < 0 || tgtIdx < 0 || tgtIdx >= cur.lessons.length) return cur;
      const srcTree = cur.lessons[srcIdx];
      const moving = (srcTree.sub_lessons ?? []).find((s) => s.lesson.id === lessonID);
      if (!moving) return cur;
      targetLevelID = cur.lessons[tgtIdx].lesson.id;
      const nextLessons = cur.lessons.map((tree, i) => {
        if (i === srcIdx) return { ...tree, sub_lessons: (tree.sub_lessons ?? []).filter((s) => s.lesson.id !== lessonID) };
        if (i === tgtIdx) return { ...tree, sub_lessons: [...(tree.sub_lessons ?? []), moving] };
        return tree;
      });
      return { ...cur, lessons: nextLessons };
    });
    if (targetLevelID) {
      setSelectedRegionId(targetLevelID);
      setSelectedPointId(lessonID);
      setDetailTargetKind('point');
      setHasStructureDraftChanges(true);
      setStructureSaveMessage('서브지역을 다른 지역으로 이동했습니다. 탐험 구조 저장을 눌러 반영해 주세요.');
    }
  };

  const handleMoveResearchNode = (levelID: string, nodeID: string, direction: -1 | 1) => {
    if (isSavingDraftStructure) return;
    setDraft((cur) => {
      if (!cur) return cur;
      const nextLessons = cur.lessons.map((tree) => {
        if (tree.lesson.id !== levelID) return tree;
        const pts = getResearchPoints(tree);
        const idx = pts.findIndex((p) => p.point.id === nodeID);
        const next = idx + direction;
        if (idx < 0 || next < 0 || next >= pts.length) return tree;
        const arr = [...pts];
        [arr[idx], arr[next]] = [arr[next], arr[idx]];
        const nonResearch = (tree.points ?? []).filter((p) => p.point.point_type !== 'research');
        return { ...tree, points: [...nonResearch, ...arr] };
      });
      return { ...cur, lessons: nextLessons };
    });
    setHasStructureDraftChanges(true);
    setStructureSaveMessage('연구지점 순서를 조정했습니다. 탐험 구조 저장을 눌러 반영해 주세요.');
  };

  const handleMoveResearchNodeToAdjacentRegion = (levelID: string, nodeID: string, direction: -1 | 1) => {
    if (isSavingDraftStructure) return;
    let targetLevelID: string | null = null;
    setDraft((cur) => {
      if (!cur) return cur;
      const srcIdx = cur.lessons.findIndex((t) => t.lesson.id === levelID);
      const tgtIdx = srcIdx + direction;
      if (srcIdx < 0 || tgtIdx < 0 || tgtIdx >= cur.lessons.length) return cur;
      const srcTree = cur.lessons[srcIdx];
      const moving = getResearchPoints(srcTree).find((p) => p.point.id === nodeID);
      if (!moving) return cur;
      targetLevelID = cur.lessons[tgtIdx].lesson.id;
      const tgtMainLessonID = targetLevelID;
      const nextLessons = cur.lessons.map((tree, i) => {
        if (i === srcIdx) return { ...tree, points: (tree.points ?? []).filter((p) => p.point.id !== nodeID) };
        if (i === tgtIdx) return { ...tree, points: [...(tree.points ?? []), { ...moving, point: { ...moving.point, course_draft_lesson_id: tgtMainLessonID } }] };
        return tree;
      });
      return { ...cur, lessons: nextLessons };
    });
    if (targetLevelID) {
      setSelectedRegionId(targetLevelID);
      setSelectedPointId(null);
      setSelectedResearchNodeId(nodeID);
      setDetailTargetKind('research_node');
      setHasStructureDraftChanges(true);
      setStructureSaveMessage('연구지점을 다른 지역으로 이동했습니다. 탐험 구조 저장을 눌러 반영해 주세요.');
    }
  };

  const handleSaveDraftStructure = async () => {
    if (!draft || !canSaveDraftStructure) return;
    const mainLessons = draft.lessons.map((tree, levelIndex) => ({
      main_lesson_id: tree.lesson.id,
      id: tree.lesson.id,
      order_index: levelIndex,
      lessons: (tree.sub_lessons ?? []).map((sub, i) => ({ id: sub.lesson.id, order_index: i })),
      points: getResearchPoints(tree).map((p, i) => ({ id: p.point.id, order_index: i })),
    }));
    setIsSavingDraftStructure(true);
    setStructureSaveMessage(null);
    try {
      const res = await fetch(`/api/v1/course-drafts/${draft.draft.id}/structure`, {
        method: 'PATCH', credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ main_lessons: mainLessons }),
      });
      const payload = (await res.json().catch(() => ({}))) as UpdateDraftStructureResponse;
      if (!res.ok) {
        if (res.status === 401) throw new Error('로그인 상태를 확인한 뒤 다시 시도해 주세요.');
        if (res.status === 404) throw new Error('탐험 구조 대상을 찾지 못했습니다. 새로고침 후 다시 확인해 주세요.');
        if (res.status === 400 || payload.error === 'invalid draft structure') throw new Error('탐험 구조를 저장할 수 없습니다. 그룹/지역 구성을 다시 확인해 주세요.');
        throw new Error('탐험 구조 저장에 실패했습니다. 기존 화면은 유지됩니다.');
      }
      if (!payload.draft) throw new Error('탐험 구조 저장 응답을 확인하지 못했습니다.');
      setDraft(payload.draft);
      setHasStructureDraftChanges(false);
      setStructureSaveMessage('탐험 구조를 저장했습니다. 그룹/지역 ID와 연결 자료는 유지됩니다.');
    } catch (e) {
      setStructureSaveMessage(e instanceof Error ? e.message : '탐험 구조 저장에 실패했습니다. 기존 화면은 유지됩니다.');
    } finally {
      setIsSavingDraftStructure(false);
    }
  };

  const handleResetDraftStructure = async () => {
    if (!draft) return;
    try {
      const res = await fetch(`/api/v1/course-drafts/${draft.draft.id}`, { credentials: 'include' });
      const payload = (await res.json().catch(() => ({}))) as { draft?: DraftAggregate };
      if (res.ok && payload.draft) {
        setDraft(payload.draft);
        setHasStructureDraftChanges(false);
        setStructureSaveMessage(null);
      }
    } catch { /* 실패 시 화면 유지 */ }
  };

  return {
    hasStructureDraftChanges, canSaveDraftStructure, isSavingDraftStructure, structureSaveMessage,
    handleMoveRegion, handleMovePoint, handleMovePointToAdjacentRegion,
    handleMoveResearchNode, handleMoveResearchNodeToAdjacentRegion,
    handleSaveDraftStructure, handleResetDraftStructure,
  };
}
