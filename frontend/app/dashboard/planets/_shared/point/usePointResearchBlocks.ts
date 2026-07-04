'use client';

import { type Dispatch, type SetStateAction, useEffect, useMemo, useState } from 'react';

import type {
  LearningPointMutationResponse,
  PlanetMetaState,
  PlanetPointDetailResponse,
  PlanetRouteKind,
  ResearchBlockType,
  ResolvedPoint,
} from '../pointPageTypes';
import {
  buildResearchBlockContent,
  createEmptyBlockDraft,
  hasResearchBlockBodyContent,
  mapBlockToDraft,
  type ResearchBlockDraft,
} from '../pointPageUtils';
import { maxResearchBlockCount } from './pointLearningUtils';
import { resolveSafetyInputMessage, type SafetyAPIErrorPayload } from '@/lib/safetyErrors';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';

type UsePointResearchBlocksArgs = {
  routeKind: PlanetRouteKind;
  pointID: string;
  planetID: string;
  planet: PlanetMetaState | null;
  pointDetail: ResolvedPoint | null;
  isResearchPoint: boolean;
  isResearchMaterialConfirmed: boolean;
  hasSavedResearchMaterialAttachment: boolean;
  hasDirtyResearchMaterialAttachment: boolean;
  applyPointMutationPayload: (payload: LearningPointMutationResponse) => boolean;
  setPointEntryMessage: Dispatch<SetStateAction<string | null>>;
  copy: PointLearningCopy['workspace']['researchMaterial'];
};

export function usePointResearchBlocks({
  routeKind,
  pointID,
  planetID,
  planet,
  pointDetail,
  isResearchPoint,
  isResearchMaterialConfirmed,
  hasSavedResearchMaterialAttachment,
  hasDirtyResearchMaterialAttachment,
  applyPointMutationPayload,
  setPointEntryMessage,
  copy,
}: UsePointResearchBlocksArgs) {
  const [researchBlocks, setResearchBlocks] = useState<ResearchBlockDraft[]>([]);
  const [isSavingResearchBlock, setIsSavingResearchBlock] = useState(false);
  const [isConfirmingResearchMaterial, setIsConfirmingResearchMaterial] = useState(false);

  useEffect(() => {
    setResearchBlocks((pointDetail?.blocks ?? []).map(mapBlockToDraft));
  }, [pointDetail?.point.id]);

  const canAddResearchBlock = researchBlocks.length < maxResearchBlockCount;
  const researchBlockLimitMessage = researchBlocks.length >= maxResearchBlockCount
    ? copy.blockLimitMessage
    : null;
  const canConfirmResearchMaterial = useMemo(() => (
    isResearchPoint
    && (hasSavedResearchMaterialAttachment || researchBlocks.some((block) => !block.isNew && !block.isDirty && hasResearchBlockBodyContent(block)))
    && !hasDirtyResearchMaterialAttachment
    && researchBlocks.every((block) => !block.isNew && !block.isDirty)
  ), [hasDirtyResearchMaterialAttachment, hasSavedResearchMaterialAttachment, isResearchPoint, researchBlocks]);

  const handleChangeResearchBlock = (
    id: string,
    field: keyof Pick<ResearchBlockDraft, 'text' | 'url' | 'title' | 'caption' | 'note'>,
    value: string,
  ) => {
    setResearchBlocks((cur) => cur.map((block) => (
      block.id === id ? { ...block, [field]: value, isDirty: true } : block
    )));
  };

  const handleAddResearchBlock = (blockType: ResearchBlockType) => {
    if (!canAddResearchBlock) {
      setPointEntryMessage(researchBlockLimitMessage ?? copy.blockAddUnavailable);
      return null;
    }
    const draft = createEmptyBlockDraft(blockType, researchBlocks.length);
    setResearchBlocks((cur) => [...cur, draft]);
    setPointEntryMessage(null);
    return draft.id;
  };

  const handleSaveResearchBlock = async (block: ResearchBlockDraft) => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingResearchBlock) return;
    setIsSavingResearchBlock(true);
    setPointEntryMessage(null);
    try {
      const endpoint = block.isNew
        ? `/api/v1/planets/learning/${planetID}/points/${pointID}/blocks`
        : `/api/v1/planets/learning/${planetID}/points/${pointID}/blocks/${block.id}`;
      const res = await fetch(endpoint, {
        method: block.isNew ? 'POST' : 'PATCH',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          block_type: block.blockType,
          content: buildResearchBlockContent(block),
          order_index: block.orderIndex,
        }),
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
      if (!res.ok || !applyPointMutationPayload(payload)) {
        throw new Error(resolveSafetyInputMessage(payload, payload.error ?? copy.blockSaveFailed));
      }
      setResearchBlocks(((payload.point as PlanetPointDetailResponse).point.blocks ?? []).map(mapBlockToDraft));
      setPointEntryMessage(block.isNew ? copy.blockAddSuccess : copy.blockSaveSuccess);
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.blockSaveFailed);
    } finally {
      setIsSavingResearchBlock(false);
    }
  };

  const handleDeleteResearchBlock = async (block: ResearchBlockDraft) => {
    if (block.isNew) {
      setResearchBlocks((cur) => cur.filter((item) => item.id !== block.id));
      setPointEntryMessage(null);
      return;
    }
    if (!planet || !planetID || routeKind !== 'learning' || isSavingResearchBlock) return;
    if (isResearchMaterialConfirmed) {
      setPointEntryMessage(copy.deleteLockedMessage);
      return;
    }
    setIsSavingResearchBlock(true);
    setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/blocks/${block.id}`, {
        method: 'DELETE',
        credentials: 'include',
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse;
      if (!res.ok || !applyPointMutationPayload(payload)) {
        if (res.status === 409) throw new Error(copy.deleteLockedMessage);
        throw new Error(payload.error ?? copy.blockDeleteFailed);
      }
      setResearchBlocks((cur) => cur.filter((item) => item.id !== block.id));
      setPointEntryMessage(null);
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.blockDeleteFailed);
    } finally {
      setIsSavingResearchBlock(false);
    }
  };

  const handleConfirmResearchMaterial = async () => {
    if (!planet || !planetID || routeKind !== 'learning' || isConfirmingResearchMaterial || !canConfirmResearchMaterial) {
      setPointEntryMessage(copy.confirmUnavailable);
      return false;
    }
    setIsConfirmingResearchMaterial(true);
    setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/research-material/confirm`, {
        method: 'POST',
        credentials: 'include',
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse;
      if (!res.ok || !applyPointMutationPayload(payload)) {
        if (res.status === 409) throw new Error(copy.confirmUnavailable);
        throw new Error(payload.error ?? copy.confirmFailed);
      }
      setPointEntryMessage(copy.confirmSuccess);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.confirmFailed);
      return false;
    } finally {
      setIsConfirmingResearchMaterial(false);
    }
  };

  const handleUnconfirmResearchMaterial = async () => {
    if (!planet || !planetID || routeKind !== 'learning' || isConfirmingResearchMaterial || !isResearchMaterialConfirmed) {
      return false;
    }
    setIsConfirmingResearchMaterial(true);
    setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/research-material/confirm`, {
        method: 'DELETE',
        credentials: 'include',
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse;
      if (!res.ok || !applyPointMutationPayload(payload)) {
        throw new Error(payload.error ?? copy.unconfirmFailed);
      }
      setPointEntryMessage(copy.unconfirmSuccess);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.unconfirmFailed);
      return false;
    } finally {
      setIsConfirmingResearchMaterial(false);
    }
  };

  return {
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
  };
}
