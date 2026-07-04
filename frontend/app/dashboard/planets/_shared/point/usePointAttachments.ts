'use client';

import { type Dispatch, type SetStateAction, useEffect, useState } from 'react';

import type {
  LearningPointMutationResponse,
  PlanetMetaState,
  PlanetPointDetailResponse,
  PlanetRouteKind,
} from '../pointPageTypes';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import { resolveSafetyInputMessage, type SafetyAPIErrorPayload } from '@/lib/safetyErrors';
import {
  buildAttachmentPayload,
  mapAttachmentToDraft,
  type PointAttachmentDraft,
} from '../pointPageUtils';
import { EMPTY_ATTACHMENT } from './pointLearningDefaults';
import {
  inferResearchMaterialAttachmentType,
  validateResearchMaterialReference,
} from './pointLearningUtils';

type UsePointAttachmentsArgs = {
  routeKind: PlanetRouteKind;
  pointID: string;
  planetID: string;
  pointDetailID?: string | null;
  pointAttachments?: PlanetPointDetailResponse['point']['attachments'] | null;
  planet: PlanetMetaState | null;
  applyPointMutationPayload: (payload: LearningPointMutationResponse) => boolean;
  setPointEntryMessage: Dispatch<SetStateAction<string | null>>;
  copy: PointLearningCopy['workspace']['runtime'];
};

export function usePointAttachments({
  routeKind,
  pointID,
  planetID,
  pointDetailID,
  pointAttachments,
  planet,
  applyPointMutationPayload,
  setPointEntryMessage,
  copy,
}: UsePointAttachmentsArgs) {
  const [attachmentDrafts, setAttachmentDrafts] = useState<PointAttachmentDraft[]>([]);
  const [newAttachment, setNewAttachment] = useState<PointAttachmentDraft>(EMPTY_ATTACHMENT);
  const [selectedAttachmentFile, setSelectedAttachmentFile] = useState<File | null>(null);
  const [isSavingAttachment, setIsSavingAttachment] = useState(false);

  useEffect(() => {
    setAttachmentDrafts((pointAttachments ?? []).map(mapAttachmentToDraft));
    setNewAttachment(EMPTY_ATTACHMENT);
    setSelectedAttachmentFile(null);
  }, [pointDetailID, pointAttachments]);

  useEffect(() => {
    if (!planet || !planetID || routeKind !== 'learning' || !pointID || !pointDetailID) return;
    let cancelled = false;
    const loadAttachments = async () => {
      try {
        const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/attachments`, {
          credentials: 'include',
          cache: 'no-store',
        });
        const payload = (await res.json().catch(() => ({}))) as { attachments?: PlanetPointDetailResponse['point']['attachments'] };
        if (!cancelled && res.ok && Array.isArray(payload.attachments)) {
          setAttachmentDrafts(payload.attachments.map(mapAttachmentToDraft));
        }
      } catch {
        // The point detail payload remains the fallback source.
      }
    };
    void loadAttachments();
    return () => {
      cancelled = true;
    };
  }, [planetID, pointID, pointDetailID, routeKind]);

  const handleChangeAttachmentDraft = (id: string, field: keyof Omit<PointAttachmentDraft, 'id'>, value: string) => {
    setAttachmentDrafts((cur) => cur.map((attachment) => (
      attachment.id === id ? { ...attachment, [field]: value, isDirty: true } : attachment
    )));
  };

  const handleChangeNewAttachment = (field: keyof Omit<PointAttachmentDraft, 'id'>, value: string) => {
    setNewAttachment((cur) => ({ ...cur, [field]: value, sourceContext: 'work_attachment' }));
  };

  const handleCreateAttachment = async () => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingAttachment || !newAttachment.title.trim() || !newAttachment.url.trim()) return false;
    setIsSavingAttachment(true);
    setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/attachments`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(buildAttachmentPayload(newAttachment)),
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(resolveSafetyInputMessage(payload, payload.error ?? copy.attachmentCreateFailed));
      setNewAttachment(EMPTY_ATTACHMENT);
      setPointEntryMessage(copy.attachmentCreated);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.attachmentCreateFailed);
      return false;
    } finally {
      setIsSavingAttachment(false);
    }
  };

  const handleCreateLinkAttachment = async (title: string, url: string) => {
    const trimmedTitle = title.trim();
    const trimmedURL = url.trim();
    if (!planet || !planetID || routeKind !== 'learning' || isSavingAttachment || !trimmedTitle || !trimmedURL) return false;
    if (!/^https?:\/\//i.test(trimmedURL)) {
      setPointEntryMessage(copy.attachmentLinkInvalid);
      return false;
    }
    setIsSavingAttachment(true);
    setPointEntryMessage(null);
    try {
      const attachment: PointAttachmentDraft = {
        ...EMPTY_ATTACHMENT,
        title: trimmedTitle,
        url: trimmedURL,
        mimeType: 'text/uri-list',
      };
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/attachments`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(buildAttachmentPayload(attachment)),
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(resolveSafetyInputMessage(payload, payload.error ?? copy.replacementAttachmentCreateFailed));
      setAttachmentDrafts(((payload.point as PlanetPointDetailResponse).point.attachments ?? []).map(mapAttachmentToDraft));
      setPointEntryMessage(copy.attachmentCreated);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.replacementAttachmentCreateFailed);
      return false;
    } finally {
      setIsSavingAttachment(false);
    }
  };

  const handleUploadAttachment = async (fileOverride?: File | null, titleOverride?: string) => {
    const uploadFile = fileOverride ?? selectedAttachmentFile;
    if (!planet || !planetID || routeKind !== 'learning' || isSavingAttachment || !uploadFile) return false;
    setIsSavingAttachment(true);
    setPointEntryMessage(null);
    try {
      const form = new FormData();
      form.set('file', uploadFile);
      form.set('title', titleOverride?.trim() || newAttachment.title.trim() || uploadFile.name);
      form.set('attachment_type', newAttachment.attachmentType === 'link' ? 'file' : newAttachment.attachmentType);
      form.set('source_context', 'work_attachment');
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/attachments/upload`, {
        method: 'POST',
        credentials: 'include',
        body: form,
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(resolveSafetyInputMessage(payload, payload.error ?? copy.attachmentUploadFailed));
      setAttachmentDrafts(((payload.point as PlanetPointDetailResponse).point.attachments ?? []).map(mapAttachmentToDraft));
      setSelectedAttachmentFile(null);
      setNewAttachment(EMPTY_ATTACHMENT);
      setPointEntryMessage(copy.attachmentUploadSuccess);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.attachmentUploadFailed);
      return false;
    } finally {
      setIsSavingAttachment(false);
    }
  };

  const handleReplaceAttachmentFile = async (attachment: PointAttachmentDraft, file: File, titleOverride?: string) => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingAttachment || attachment.sourceContext !== 'work_attachment') return false;
    setIsSavingAttachment(true);
    setPointEntryMessage(null);
    try {
      const form = new FormData();
      form.set('file', file);
      form.set('title', titleOverride?.trim() || attachment.title.trim() || file.name);
      form.set('attachment_type', 'file');
      form.set('source_context', 'work_attachment');
      form.set('replace_attachment_id', attachment.id);
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/attachments/upload`, {
        method: 'POST',
        credentials: 'include',
        body: form,
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(resolveSafetyInputMessage(payload, payload.error ?? copy.attachmentReplaceFailed));
      setAttachmentDrafts(((payload.point as PlanetPointDetailResponse).point.attachments ?? []).map(mapAttachmentToDraft));
      setSelectedAttachmentFile(null);
      setPointEntryMessage(copy.attachmentReplaced);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.attachmentReplaceFailed);
      return false;
    } finally {
      setIsSavingAttachment(false);
    }
  };

  const handleOpenAttachment = async (attachment: PointAttachmentDraft) => {
    if (!planet || !planetID) return;
    if (!attachment.filePath.trim()) {
      if (attachment.url.trim()) window.open(attachment.url, '_blank', 'noopener,noreferrer');
      return;
    }
    setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/attachments/${attachment.id}/open`, {
        credentials: 'include',
      });
      const payload = (await res.json().catch(() => ({}))) as { url?: string; error?: string };
      if (!res.ok || !payload.url) throw new Error(copy.attachmentOpenFailed);
      window.open(payload.url, '_blank', 'noopener,noreferrer');
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.attachmentOpenFailed);
    }
  };

  const handleUpdateAttachment = async (attachment: PointAttachmentDraft) => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingAttachment || !attachment.title.trim() || (!attachment.url.trim() && !attachment.filePath.trim())) return false;
    if (attachment.sourceContext === 'research_material') {
      const validationTarget = attachment.filePath.trim() || attachment.url.trim();
      const parsedFileSize = attachment.fileSize.trim() ? Number(attachment.fileSize.trim()) : null;
      const validationError = validateResearchMaterialReference(validationTarget, Number.isFinite(parsedFileSize) ? parsedFileSize : null);
      if (validationError) {
        setPointEntryMessage(validationError);
        return false;
      }
    }
    setIsSavingAttachment(true);
    setPointEntryMessage(null);
    try {
      const payloadAttachment = attachment.sourceContext === 'research_material'
        ? { ...attachment, provider: 'learner' as const, attachmentType: inferResearchMaterialAttachmentType(attachment.filePath.trim() || attachment.url.trim()) }
        : attachment;
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/attachments/${attachment.id}`, {
        method: 'PATCH',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(buildAttachmentPayload(payloadAttachment)),
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
      if (!res.ok || !applyPointMutationPayload(payload)) throw new Error(resolveSafetyInputMessage(payload, payload.error ?? copy.attachmentUpdateFailed));
      setAttachmentDrafts(((payload.point as PlanetPointDetailResponse).point.attachments ?? []).map(mapAttachmentToDraft));
      setPointEntryMessage(copy.attachmentUpdated);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.attachmentUpdateFailed);
      return false;
    } finally {
      setIsSavingAttachment(false);
    }
  };

  const handleDeleteAttachment = async (attachment: PointAttachmentDraft) => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingAttachment) return false;
    setIsSavingAttachment(true);
    setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/attachments/${attachment.id}`, {
        method: 'DELETE',
        credentials: 'include',
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse;
      if (!res.ok || !applyPointMutationPayload(payload)) {
        if (res.status === 409) throw new Error(copy.attachmentDeleteLocked);
        throw new Error(copy.attachmentDeleteFailed);
      }
      setAttachmentDrafts(((payload.point as PlanetPointDetailResponse).point.attachments ?? []).map(mapAttachmentToDraft));
      setPointEntryMessage(null);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.attachmentDeleteFailed);
      return false;
    } finally {
      setIsSavingAttachment(false);
    }
  };

  return {
    attachmentDrafts,
    setAttachmentDrafts,
    newAttachment,
    selectedAttachmentFile,
    setSelectedAttachmentFile,
    isSavingAttachment,
    setIsSavingAttachment,
    handleCreateAttachment,
    handleCreateLinkAttachment,
    handleUploadAttachment,
    handleReplaceAttachmentFile,
    handleOpenAttachment,
    handleUpdateAttachment,
    handleDeleteAttachment,
    handleChangeAttachmentDraft,
    handleChangeNewAttachment,
  };
}
