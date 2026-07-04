'use client';

import { type Dispatch, type SetStateAction, useEffect, useState } from 'react';

import {
  isLearningVideoFile,
  normalizeSubtitleFile,
  validateLearningSubtitleFile,
  validateLearningThumbnailFile,
  validateLearningVideoFile,
} from '@/lib/media/videoValidation';
import type {
  LearningPointMutationResponse,
  PlanetMetaState,
  PlanetPointDetailResponse,
  PlanetRouteKind,
  PlanetPoint,
} from '../pointPageTypes';
import {
  buildAttachmentPayload,
  mapAttachmentToDraft,
  type PointAttachmentDraft,
} from '../pointPageUtils';
import {
  EMPTY_RESEARCH_MATERIAL_ATTACHMENT,
} from './pointLearningDefaults';
import {
  inferResearchMaterialAttachmentType,
  uploadResearchMaterialAttachmentWithProgress,
  validateInlineImageUploadFile,
  validateResearchMaterialReference,
} from './pointLearningUtils';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import { resolveSafetyInputMessage, type SafetyAPIErrorPayload } from '@/lib/safetyErrors';

type UsePointResearchMaterialAttachmentsArgs = {
  routeKind: PlanetRouteKind;
  pointID: string;
  planetID: string;
  pointDetailID?: string | null;
  planet: PlanetMetaState | null;
  attachmentDrafts: PointAttachmentDraft[];
  canAddResearchMaterialAttachment: boolean;
  researchMaterialAttachmentLimitMessage: string | null;
  isSavingAttachment: boolean;
  setIsSavingAttachment: Dispatch<SetStateAction<boolean>>;
  setAttachmentDrafts: Dispatch<SetStateAction<PointAttachmentDraft[]>>;
  setPointEntryMessage: Dispatch<SetStateAction<string | null>>;
  applyPointMutationPayload: (payload: LearningPointMutationResponse) => boolean;
  copy: PointLearningCopy['workspace']['researchMaterial'];
};

export function usePointResearchMaterialAttachments({
  routeKind,
  pointID,
  planetID,
  pointDetailID,
  planet,
  attachmentDrafts,
  canAddResearchMaterialAttachment,
  researchMaterialAttachmentLimitMessage,
  isSavingAttachment,
  setIsSavingAttachment,
  setAttachmentDrafts,
  setPointEntryMessage,
  applyPointMutationPayload,
  copy,
}: UsePointResearchMaterialAttachmentsArgs) {
  const [newResearchMaterialAttachment, setNewResearchMaterialAttachment] = useState<PointAttachmentDraft>(EMPTY_RESEARCH_MATERIAL_ATTACHMENT);
  const [selectedResearchMaterialFile, setSelectedResearchMaterialFile] = useState<File | null>(null);
  const [researchMaterialAttachmentValidationMessage, setResearchMaterialAttachmentValidationMessage] = useState<string | null>(null);
  const [researchMaterialUploadProgress, setResearchMaterialUploadProgress] = useState<number | null>(null);

  useEffect(() => {
    setNewResearchMaterialAttachment(EMPTY_RESEARCH_MATERIAL_ATTACHMENT);
    setSelectedResearchMaterialFile(null);
    setResearchMaterialAttachmentValidationMessage(null);
    setResearchMaterialUploadProgress(null);
  }, [pointDetailID]);

  const handleChangeNewResearchMaterialAttachment = (field: keyof Omit<PointAttachmentDraft, 'id'>, value: string) => {
    setNewResearchMaterialAttachment((cur) => ({
      ...cur,
      [field]: value,
      provider: 'learner',
      sourceContext: 'research_material',
    }));
  };

  const handleCreateResearchMaterialAttachment = async () => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingAttachment || !newResearchMaterialAttachment.title.trim() || !newResearchMaterialAttachment.url.trim()) return;
    if (!canAddResearchMaterialAttachment) {
      setResearchMaterialAttachmentValidationMessage(researchMaterialAttachmentLimitMessage ?? copy.attachmentAddUnavailable);
      return;
    }
    const validationError = validateResearchMaterialReference(
      newResearchMaterialAttachment.url,
      newResearchMaterialAttachment.fileSize ? Number(newResearchMaterialAttachment.fileSize) : null,
      {
        typeError: copy.referenceTypeError,
        imageSizeError: copy.imageSizeError,
        videoSizeError: copy.videoSizeError,
        documentSizeError: copy.documentSizeError,
      },
    );
    if (validationError) {
      setResearchMaterialAttachmentValidationMessage(validationError);
      return;
    }
    setIsSavingAttachment(true);
    setPointEntryMessage(null);
    try {
      const attachmentType = inferResearchMaterialAttachmentType(newResearchMaterialAttachment.url);
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/attachments`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(buildAttachmentPayload({
          ...newResearchMaterialAttachment,
          provider: 'learner',
          attachmentType,
          sourceContext: 'research_material',
        })),
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(resolveSafetyInputMessage(payload, payload.error ?? copy.attachmentAddFailed));
      setAttachmentDrafts(((payload.point as PlanetPointDetailResponse).point.attachments ?? []).map(mapAttachmentToDraft));
      setNewResearchMaterialAttachment(EMPTY_RESEARCH_MATERIAL_ATTACHMENT);
      setPointEntryMessage(copy.attachmentAddSuccess);
    } catch (e) {
      setResearchMaterialAttachmentValidationMessage(e instanceof Error ? e.message : copy.attachmentAddFailed);
    } finally {
      setIsSavingAttachment(false);
    }
  };

  const handleUploadResearchMaterialAttachment = async (
    file?: File,
    forcedAttachmentType?: 'video' | 'subtitle' | 'thumbnail',
    titleOverride?: string,
  ) => {
    let uploadFile = file ?? selectedResearchMaterialFile;
    if (!planet || !planetID || routeKind !== 'learning' || isSavingAttachment || !uploadFile) return false;
    if (!canAddResearchMaterialAttachment) {
      setSelectedResearchMaterialFile(null);
      setResearchMaterialAttachmentValidationMessage(researchMaterialAttachmentLimitMessage ?? copy.attachmentAddUnavailable);
      return false;
    }
    if (forcedAttachmentType === 'subtitle') {
      const subtitleValidationError = validateLearningSubtitleFile(uploadFile, {
        subtitleTypeError: copy.subtitleTypeError,
        subtitleSizeError: copy.subtitleSizeError,
      });
      if (subtitleValidationError) {
        setSelectedResearchMaterialFile(null);
        setResearchMaterialAttachmentValidationMessage(subtitleValidationError);
        return false;
      }
      const hasSubtitle = attachmentDrafts.some((attachment) => attachment.sourceContext === 'research_material' && attachment.attachmentType === 'subtitle');
      if (hasSubtitle) {
        setSelectedResearchMaterialFile(null);
        setResearchMaterialAttachmentValidationMessage(copy.subtitleDuplicate);
        return false;
      }
      uploadFile = await normalizeSubtitleFile(uploadFile);
    }
    if (forcedAttachmentType === 'thumbnail') {
      const thumbnailValidationError = validateLearningThumbnailFile(uploadFile, {
        thumbnailTypeError: copy.thumbnailTypeError,
        thumbnailSizeError: copy.thumbnailSizeError,
      });
      if (thumbnailValidationError) {
        setSelectedResearchMaterialFile(null);
        setResearchMaterialAttachmentValidationMessage(thumbnailValidationError);
        return false;
      }
      const hasThumbnail = attachmentDrafts.some((attachment) => attachment.sourceContext === 'research_material' && attachment.attachmentType === 'thumbnail');
      if (hasThumbnail) {
        setSelectedResearchMaterialFile(null);
        setResearchMaterialAttachmentValidationMessage(copy.thumbnailDuplicate);
        return false;
      }
    }
    const validationError = forcedAttachmentType === 'thumbnail'
      ? null
      : validateResearchMaterialReference(uploadFile.name, uploadFile.size, {
        typeError: copy.referenceTypeError,
        imageSizeError: copy.imageSizeError,
        videoSizeError: copy.videoSizeError,
        documentSizeError: copy.documentSizeError,
      });
    if (validationError) {
      setSelectedResearchMaterialFile(null);
      setResearchMaterialAttachmentValidationMessage(validationError);
      return false;
    }
    if (forcedAttachmentType === 'video' || isLearningVideoFile(uploadFile)) {
      const videoValidationError = await validateLearningVideoFile(uploadFile, {
        videoTypeError: copy.videoTypeError,
        videoSizeError: copy.videoSizeError,
        videoDurationError: copy.videoDurationError,
        videoMetadataError: copy.videoMetadataError,
      });
      if (videoValidationError) {
        setSelectedResearchMaterialFile(null);
        setResearchMaterialAttachmentValidationMessage(videoValidationError);
        return false;
      }
      const hasVideo = attachmentDrafts.some((attachment) => attachment.sourceContext === 'research_material' && attachment.attachmentType === 'video');
      if (hasVideo) {
        setSelectedResearchMaterialFile(null);
        setResearchMaterialAttachmentValidationMessage(copy.videoDuplicate);
        return false;
      }
    }
    setIsSavingAttachment(true);
    setPointEntryMessage(null);
    setResearchMaterialUploadProgress(0);
    try {
      const form = new FormData();
      form.set('file', uploadFile);
      form.set('title', titleOverride?.trim() || newResearchMaterialAttachment.title.trim() || uploadFile.name);
      form.set('attachment_type', forcedAttachmentType ?? inferResearchMaterialAttachmentType(uploadFile.name));
      form.set('source_context', 'research_material');
      const uploadURL = `/api/v1/planets/learning/${planetID}/points/${pointID}/attachments/upload`;
      let payload: LearningPointMutationResponse;
      try {
        payload = await uploadResearchMaterialAttachmentWithProgress(uploadURL, form, setResearchMaterialUploadProgress, {
          networkError: copy.fileUploadFailed,
          uploadTooLargeServer: copy.uploadTooLargeServer,
          uploadFailedStatus: copy.uploadFailedStatus,
        });
      } catch (error) {
        if (error instanceof Error && error.name === 'UploadHttpError') throw error;
        setResearchMaterialUploadProgress(35);
        const res = await fetch(uploadURL, { method: 'POST', credentials: 'include', body: form });
        payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
        if (!res.ok) {
          throw new Error(resolveSafetyInputMessage(payload, payload.error ?? (res.status === 413
            ? copy.uploadTooLargeServer
            : copy.uploadFailedStatus(res.status))));
        }
        setResearchMaterialUploadProgress(100);
      }
      if (!applyPointMutationPayload(payload)) throw new Error(resolveSafetyInputMessage(payload as SafetyAPIErrorPayload, (payload as SafetyAPIErrorPayload).error ?? copy.fileUploadFailed));
      setAttachmentDrafts(((payload.point as PlanetPointDetailResponse).point.attachments ?? []).map(mapAttachmentToDraft));
      setSelectedResearchMaterialFile(null);
      setNewResearchMaterialAttachment(EMPTY_RESEARCH_MATERIAL_ATTACHMENT);
      setPointEntryMessage(copy.fileUploadSuccess);
      return true;
    } catch (e) {
      setSelectedResearchMaterialFile(null);
      setResearchMaterialAttachmentValidationMessage(e instanceof Error ? e.message : copy.fileUploadFailed);
      return false;
    } finally {
      setIsSavingAttachment(false);
      setResearchMaterialUploadProgress(null);
    }
  };

  const handleUploadResearchMaterialInlineImage = async (file: File): Promise<{ src: string; alt?: string }> => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingAttachment) throw new Error(copy.inlineImageUnavailable);
    const validationError = validateInlineImageUploadFile(file, {
      typeError: copy.inlineImageTypeError,
      sizeError: copy.inlineImageSizeError,
    });
    if (validationError) throw new Error(validationError);
    setIsSavingAttachment(true);
    setPointEntryMessage(null);
    try {
      const form = new FormData();
      form.set('file', file);
      form.set('title', file.name);
      form.set('attachment_type', 'image');
      form.set('source_context', 'research_material');
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/attachments/upload`, {
        method: 'POST',
        credentials: 'include',
        body: form,
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload & { attachment?: NonNullable<PlanetPoint['attachments']>[number] };
      if (!res.ok || !payload.attachment || !applyPointMutationPayload(payload)) throw new Error(resolveSafetyInputMessage(payload, payload.error ?? copy.inlineImageUploadFailed));
      setPointEntryMessage(copy.inlineImageInserted);
      return {
        src: `/api/v1/planets/learning/${planetID}/points/${pointID}/attachments/${payload.attachment.id}/inline`,
        alt: payload.attachment.title || file.name,
      };
    } finally {
      setIsSavingAttachment(false);
    }
  };

  return {
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
  };
}
