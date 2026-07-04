import type { Dispatch, SetStateAction } from 'react';

import {
  learningVideoAllowedExtensions,
  learningVideoMaxBytes,
  learningSubtitleAllowedExtensions,
} from '@/lib/media/videoValidation';
import type { PointAttachmentDraft } from '../pointPageUtils';
import { resolveSafetyInputMessage, type SafetyAPIErrorPayload } from '@/lib/safetyErrors';
import type { LearningPointMutationResponse } from '../pointPageTypes';

export const inlineImageAllowedExtensions = ['jpg', 'jpeg', 'png', 'webp', 'gif'] as const;
export const inlineImageAllowedMimeTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/gif'];
export const inlineImageMaxBytes = 5 * 1024 * 1024;
export const maxResearchBlockCount = 3;
export const researchMaterialAllowedImageExtensions = ['jpg', 'jpeg', 'png', 'webp', 'gif'] as const;
export const researchMaterialAllowedDocumentExtensions = ['pdf', 'doc', 'docx', 'hwp', 'hwpx', 'txt', 'xls', 'xlsx', 'csv', 'ppt', 'pptx', 'rtf', 'odt', 'md', 'zip'] as const;
export const researchMaterialAllowedExtensions = [...researchMaterialAllowedImageExtensions, ...researchMaterialAllowedDocumentExtensions, ...learningVideoAllowedExtensions, ...learningSubtitleAllowedExtensions] as const;
export const researchMaterialMaxAttachmentCount = 10;
export const researchMaterialImageMaxBytes = 5 * 1024 * 1024;
export const researchMaterialDocumentMaxBytes = 20 * 1024 * 1024;
export const learningSessionHeartbeatMs = 60 * 1000;

export function formatInlineImageLimit(bytes: number): string {
  return `${Math.round((bytes / 1024 / 1024) * 10) / 10}MB`;
}

export function validateInlineImageUploadFile(file: File, messages?: { typeError: string; sizeError: string }): string | null {
  const extension = file.name.split('.').pop()?.toLowerCase() ?? '';
  if (!inlineImageAllowedExtensions.includes(extension as (typeof inlineImageAllowedExtensions)[number])) {
    return messages?.typeError ?? 'jpg, jpeg, png, webp, gif 이미지만 업로드할 수 있습니다.';
  }
  if (file.type && !inlineImageAllowedMimeTypes.includes(file.type.toLowerCase())) {
    return messages?.typeError ?? 'jpg, jpeg, png, webp, gif 이미지만 업로드할 수 있습니다.';
  }
  if (file.size > inlineImageMaxBytes) {
    return messages?.sizeError ?? `이미지는 ${formatInlineImageLimit(inlineImageMaxBytes)} 이하만 업로드할 수 있습니다.`;
  }
  return null;
}

export function formatFileLimit(bytes: number): string {
  return `${Math.round(bytes / 1024 / 1024)}MB`;
}

export function getFileExtension(value: string): string {
  const withoutQuery = value.split('?')[0]?.split('#')[0] ?? value;
  return withoutQuery.split('.').pop()?.toLowerCase() ?? '';
}

export function isResearchMaterialImageExtension(extension: string): boolean {
  return researchMaterialAllowedImageExtensions.includes(extension as (typeof researchMaterialAllowedImageExtensions)[number]);
}

export function isResearchMaterialVideoExtension(extension: string): boolean {
  return learningVideoAllowedExtensions.includes(extension as (typeof learningVideoAllowedExtensions)[number]);
}

export function isResearchMaterialSubtitleExtension(extension: string): boolean {
  return learningSubtitleAllowedExtensions.includes(extension as (typeof learningSubtitleAllowedExtensions)[number]);
}

export function inferResearchMaterialAttachmentType(filename: string): 'image' | 'file' | 'video' | 'subtitle' {
  if (isResearchMaterialVideoExtension(getFileExtension(filename))) return 'video';
  if (isResearchMaterialSubtitleExtension(getFileExtension(filename))) return 'subtitle';
  return isResearchMaterialImageExtension(getFileExtension(filename)) ? 'image' : 'file';
}

export function validateResearchMaterialReference(nameOrURL: string, fileSize?: number | null, messages?: {
  typeError: string;
  imageSizeError: string;
  videoSizeError: string;
  documentSizeError: string;
}): string | null {
  const extension = getFileExtension(nameOrURL);
  if (!researchMaterialAllowedExtensions.includes(extension as (typeof researchMaterialAllowedExtensions)[number])) {
    return messages?.typeError ?? `학습 콘텐츠는 아래 파일만 추가할 수 있습니다.\n이미지: ${researchMaterialAllowedImageExtensions.join(', ')}\n문서/zip: ${researchMaterialAllowedDocumentExtensions.join(', ')}`;
  }
  if (fileSize == null) return null;
  const limit = isResearchMaterialImageExtension(extension) ? researchMaterialImageMaxBytes : researchMaterialDocumentMaxBytes;
  const effectiveLimit = isResearchMaterialVideoExtension(extension) ? learningVideoMaxBytes : limit;
  if (fileSize > effectiveLimit) {
    return isResearchMaterialImageExtension(extension)
      ? messages?.imageSizeError ?? `선택한 이미지는 ${formatFileLimit(researchMaterialImageMaxBytes)}를 넘었습니다.\n이미지는 파일당 ${formatFileLimit(researchMaterialImageMaxBytes)} 이하만 업로드할 수 있습니다.`
      : isResearchMaterialVideoExtension(extension)
        ? messages?.videoSizeError ?? `선택한 동영상은 ${formatFileLimit(learningVideoMaxBytes)}를 넘었습니다.\n동영상은 파일당 ${formatFileLimit(learningVideoMaxBytes)} 이하만 업로드할 수 있습니다.`
      : messages?.documentSizeError ?? `선택한 문서/zip은 ${formatFileLimit(researchMaterialDocumentMaxBytes)}를 넘었습니다.\n문서와 zip은 파일당 ${formatFileLimit(researchMaterialDocumentMaxBytes)} 이하만 업로드할 수 있습니다.`;
  }
  return null;
}

export function createJournalNoteSignature(
  observation: string,
  reflection: string,
  examples: string,
  nextStep: string,
): string {
  return [observation, reflection, examples, nextStep].map((value) => value.trim()).join('\n');
}

export function uploadResearchMaterialAttachmentWithProgress(
  url: string,
  form: FormData,
  onProgress: Dispatch<SetStateAction<number | null>>,
  messages?: { networkError: string; uploadTooLargeServer: string; uploadFailedStatus: (status: number) => string },
): Promise<LearningPointMutationResponse> {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    xhr.open('POST', url);
    xhr.withCredentials = true;
    xhr.upload.onprogress = (event) => {
      if (!event.lengthComputable || event.total <= 0) return;
      const percent = Math.max(1, Math.min(99, Math.round((event.loaded / event.total) * 100)));
      onProgress(percent);
    };
    xhr.onerror = () => reject(new Error(messages?.networkError ?? '학습 콘텐츠 파일 업로드에 실패했습니다.'));
    xhr.onload = () => {
      let payload: LearningPointMutationResponse & SafetyAPIErrorPayload = {};
      try {
        payload = xhr.responseText ? JSON.parse(xhr.responseText) as LearningPointMutationResponse & SafetyAPIErrorPayload : {};
      } catch {
        payload = {};
      }
      if (xhr.status < 200 || xhr.status >= 300) {
        const fallbackMessage = payload.error ?? (xhr.status === 413
          ? messages?.uploadTooLargeServer ?? '파일이 앱 기준은 통과했지만 서버 업로드 허용 크기를 넘었습니다. 관리자에게 nginx 업로드 제한 설정 확인을 요청해 주세요.'
          : messages?.uploadFailedStatus?.(xhr.status) ?? `학습 콘텐츠 파일 업로드에 실패했습니다. (HTTP ${xhr.status})`);
        const message = resolveSafetyInputMessage(payload, fallbackMessage);
        const error = new Error(message);
        error.name = 'UploadHttpError';
        reject(error);
        return;
      }
      onProgress(100);
      resolve(payload);
    };
    xhr.send(form);
  });
}
