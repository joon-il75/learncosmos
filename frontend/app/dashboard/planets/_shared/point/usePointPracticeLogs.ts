import { useEffect, useState } from 'react';

import type {
  PlanetRouteKind,
  PlanetMetaState,
  ResolvedPoint,
  LearningPointMutationResponse,
} from '../pointPageTypes';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import { resolveSafetyInputMessage, type SafetyAPIErrorPayload } from '@/lib/safetyErrors';
import {
  buildPracticeLogPayload,
  createEmptyPracticeLogDraft,
  mapPracticeLogToDraft,
  type PointPracticeLogDraft,
} from '../pointPageUtils';

type UsePointPracticeLogsParams = {
  routeKind: PlanetRouteKind;
  pointID: string;
  planetID: string;
  planet: PlanetMetaState | null;
  pointDetail: ResolvedPoint | null;
  applyPointMutationPayload: (payload: LearningPointMutationResponse) => boolean;
  setPointEntryMessage: (message: string | null) => void;
  copy: PointLearningCopy['workspace']['runtime'];
};

export function usePointPracticeLogs({
  routeKind,
  pointID,
  planetID,
  planet,
  pointDetail,
  applyPointMutationPayload,
  setPointEntryMessage,
  copy,
}: UsePointPracticeLogsParams) {
  const [practiceLogDrafts, setPracticeLogDrafts] = useState<PointPracticeLogDraft[]>([]);
  const [newPracticeLog, setNewPracticeLog] = useState<PointPracticeLogDraft>(() => createEmptyPracticeLogDraft());
  const [isSavingPracticeLog, setIsSavingPracticeLog] = useState(false);

  useEffect(() => {
    setPracticeLogDrafts((pointDetail?.practice_logs ?? []).map(mapPracticeLogToDraft));
    setNewPracticeLog(createEmptyPracticeLogDraft());
  }, [pointDetail]);

  const handleChangePracticeLogDraft = (id: string, field: keyof Omit<PointPracticeLogDraft, 'id'>, value: string) => {
    const nextValue = field === 'durationMinutes' ? value.replace(/\D/g, '') : value;
    setPracticeLogDrafts((cur) => cur.map((log) => (
      log.id === id ? { ...log, [field]: nextValue } : log
    )));
  };

  const handleChangeNewPracticeLog = (field: keyof Omit<PointPracticeLogDraft, 'id'>, value: string) => {
    const nextValue = field === 'durationMinutes' ? value.replace(/\D/g, '') : value;
    setNewPracticeLog((cur) => ({ ...cur, [field]: nextValue }));
  };

  const handleCreatePracticeLog = async () => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingPracticeLog || !newPracticeLog.title.trim()) return false;
    setIsSavingPracticeLog(true);
    setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/practice-logs`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(buildPracticeLogPayload(newPracticeLog)),
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
      if (!res.ok || !applyPointMutationPayload(payload)) {
        throw new Error(resolveSafetyInputMessage(payload, payload.error ?? copy.practiceCreateFailed));
      }
      setNewPracticeLog(createEmptyPracticeLogDraft());
      setPointEntryMessage(copy.practiceCreated);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.practiceCreateFailed);
      return false;
    } finally {
      setIsSavingPracticeLog(false);
    }
  };

  const handleUpdatePracticeLog = async (log: PointPracticeLogDraft) => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingPracticeLog || !log.title.trim()) return false;
    setIsSavingPracticeLog(true);
    setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/practice-logs/${log.id}`, {
        method: 'PATCH',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(buildPracticeLogPayload(log)),
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
      if (!res.ok || !applyPointMutationPayload(payload)) {
        throw new Error(resolveSafetyInputMessage(payload, payload.error ?? copy.practiceUpdateFailed));
      }
      setPointEntryMessage(copy.practiceUpdated);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.practiceUpdateFailed);
      return false;
    } finally {
      setIsSavingPracticeLog(false);
    }
  };

  const handleDeletePracticeLog = async (log: PointPracticeLogDraft) => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingPracticeLog) return false;
    setIsSavingPracticeLog(true);
    setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/practice-logs/${log.id}`, {
        method: 'DELETE',
        credentials: 'include',
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse;
      if (!res.ok || !applyPointMutationPayload(payload)) {
        throw new Error(copy.practiceDeleteFailed);
      }
      setPointEntryMessage(copy.practiceDeleted);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.practiceDeleteFailed);
      return false;
    } finally {
      setIsSavingPracticeLog(false);
    }
  };

  return {
    practiceLogDrafts,
    newPracticeLog,
    isSavingPracticeLog,
    handleCreatePracticeLog,
    handleUpdatePracticeLog,
    handleDeletePracticeLog,
    handleChangePracticeLogDraft,
    handleChangeNewPracticeLog,
  };
}
