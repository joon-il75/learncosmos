import { useEffect, useState } from 'react';

import type {
  LearningPointMutationResponse,
  PlanetMetaState,
  PlanetPoint,
  PlanetPointDetailResponse,
  PlanetRouteKind,
  ResolvedPoint,
} from '../pointPageTypes';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import { resolveSafetyInputMessage, type SafetyAPIErrorPayload } from '@/lib/safetyErrors';
import {
  mapArtifactToDraft,
  type PointArtifactDraft,
} from '../pointPageUtils';

type UsePointArtifactsParams = {
  routeKind: PlanetRouteKind;
  pointID: string;
  planetID: string;
  planet: PlanetMetaState | null;
  pointDetail: ResolvedPoint | null;
  applyPointMutationPayload: (payload: LearningPointMutationResponse) => boolean;
  setPointEntryMessage: (message: string | null) => void;
  copy: PointLearningCopy['workspace']['runtime'];
};

export function usePointArtifacts({
  routeKind,
  pointID,
  planetID,
  planet,
  pointDetail,
  applyPointMutationPayload,
  setPointEntryMessage,
  copy,
}: UsePointArtifactsParams) {
  const [artifactType, setArtifactType] = useState('문서 편집자료');
  const [artifactTitle, setArtifactTitle] = useState('');
  const [artifactURL, setArtifactURL] = useState('');
  const [artifactDescription, setArtifactDescription] = useState('');
  const [artifactPointCategory, setArtifactPointCategory] = useState('');
  const [artifactProductionProcess, setArtifactProductionProcess] = useState('');
  const [artifactLearnedPoints, setArtifactLearnedPoints] = useState('');
  const [artifactDifficultPoints, setArtifactDifficultPoints] = useState('');
  const [artifactVisibility, setArtifactVisibility] = useState('private');
  const [artifactDrafts, setArtifactDrafts] = useState<PointArtifactDraft[]>([]);
  const [isSavingArtifact, setIsSavingArtifact] = useState(false);

  useEffect(() => {
    if (!pointDetail) return;
    setArtifactType('문서 편집자료');
    setArtifactTitle('');
    setArtifactURL('');
    setArtifactDescription('');
    setArtifactPointCategory(pointDetail.point.point_category ?? '');
    setArtifactProductionProcess('');
    setArtifactLearnedPoints('');
    setArtifactDifficultPoints('');
    setArtifactVisibility('private');
    setArtifactDrafts((pointDetail.artifacts ?? []).map(mapArtifactToDraft));
  }, [pointDetail]);

  const handleSaveArtifact = async (override?: { artifactType?: string; title?: string }) => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingArtifact) return null;
    setIsSavingArtifact(true);
    setPointEntryMessage(null);
    const nextArtifactType = override?.artifactType ?? artifactType;
    const fallbackTitle = nextArtifactType === '영상자료'
      ? copy.artifactFallbackTitles.video
      : nextArtifactType === '첨부자료'
        ? copy.artifactFallbackTitles.attachment
        : copy.artifactFallbackTitles.document;
    const nextArtifactTitle = override?.title?.trim() || artifactTitle.trim() || fallbackTitle;
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/artifacts`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          artifact_type: nextArtifactType,
          title: nextArtifactTitle,
          url: artifactURL,
          description: artifactDescription,
          point_category: artifactPointCategory,
          production_process: artifactProductionProcess,
          learned_points: artifactLearnedPoints,
          difficult_points: artifactDifficultPoints,
          visibility: artifactVisibility,
          order_index: artifactDrafts.length,
        }),
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload & {
        artifact?: NonNullable<PlanetPoint['artifacts']>[number];
      };
      if (!res.ok || !applyPointMutationPayload(payload)) {
        throw new Error(resolveSafetyInputMessage(payload, payload.error ?? copy.artifactSaveFailed));
      }
      const payloadArtifacts = ((payload.point as PlanetPointDetailResponse).point.artifacts ?? []).map(mapArtifactToDraft);
      const createdArtifact = payload.artifact
        ? mapArtifactToDraft(payload.artifact)
        : payloadArtifacts.find((artifact) => artifact.orderIndex === artifactDrafts.length && artifact.title === nextArtifactTitle)
          ?? payloadArtifacts.at(-1)
          ?? null;
      if (createdArtifact) {
        setArtifactDrafts((current) => {
          const next = payloadArtifacts.length ? payloadArtifacts : [...current, createdArtifact];
          return next.some((artifact) => artifact.id === createdArtifact.id) ? next : [...next, createdArtifact];
        });
      } else if (payloadArtifacts.length) {
        setArtifactDrafts(payloadArtifacts);
      }
      setArtifactType('문서 편집자료');
      setArtifactTitle('');
      setArtifactURL('');
      setArtifactDescription('');
      setArtifactProductionProcess('');
      setArtifactLearnedPoints('');
      setArtifactDifficultPoints('');
      setArtifactVisibility('private');
      setPointEntryMessage(copy.artifactCreated);
      return createdArtifact;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.artifactSaveFailed);
      return null;
    } finally {
      setIsSavingArtifact(false);
    }
  };

  const handleChangeArtifactDraft = (
    id: string,
    field: keyof Pick<PointArtifactDraft, 'artifactType' | 'title' | 'url' | 'description' | 'pointCategory' | 'productionProcess' | 'learnedPoints' | 'difficultPoints' | 'visibility'>,
    value: string,
  ) => {
    setArtifactDrafts((cur) => cur.map((artifact) => (
      artifact.id === id ? { ...artifact, [field]: value } : artifact
    )));
  };

  const handleUpdateArtifact = async (artifact: PointArtifactDraft) => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingArtifact) return false;
    setIsSavingArtifact(true);
    setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/artifacts/${artifact.id}`, {
        method: 'PATCH',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          artifact_type: artifact.artifactType,
          title: artifact.title,
          url: artifact.url,
          description: artifact.description,
          point_category: artifact.pointCategory,
          production_process: artifact.productionProcess,
          learned_points: artifact.learnedPoints,
          difficult_points: artifact.difficultPoints,
          visibility: artifact.visibility,
          order_index: artifact.orderIndex,
        }),
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
      if (!res.ok || !applyPointMutationPayload(payload)) {
        throw new Error(resolveSafetyInputMessage(payload, payload.error ?? copy.artifactUpdateFailed));
      }
      setArtifactDrafts(((payload.point as PlanetPointDetailResponse).point.artifacts ?? []).map(mapArtifactToDraft));
      setPointEntryMessage(copy.artifactUpdated);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.artifactUpdateFailed);
      return false;
    } finally {
      setIsSavingArtifact(false);
    }
  };

  const handleDeleteArtifact = async (artifact: PointArtifactDraft) => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingArtifact) return false;
    setIsSavingArtifact(true);
    setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/artifacts/${artifact.id}`, {
        method: 'DELETE',
        credentials: 'include',
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse;
      if (!res.ok || !applyPointMutationPayload(payload)) {
        throw new Error(copy.artifactDeleteFailed);
      }
      setArtifactDrafts(((payload.point as PlanetPointDetailResponse).point.artifacts ?? []).map(mapArtifactToDraft));
      setPointEntryMessage(copy.artifactDeleted);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.artifactDeleteFailed);
      return false;
    } finally {
      setIsSavingArtifact(false);
    }
  };

  return {
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
  };
}
