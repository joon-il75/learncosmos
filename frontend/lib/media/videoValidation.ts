export const learningVideoAllowedExtensions = ['mp4', 'webm', 'mov'] as const;
export const learningVideoAllowedMimeTypes = ['video/mp4', 'video/webm', 'video/quicktime'];
export const learningVideoMaxBytes = 500 * 1024 * 1024;
export const learningVideoMaxDurationSeconds = 25 * 60;
export const learningSubtitleAllowedExtensions = ['vtt', 'srt'] as const;
export const learningSubtitleMaxBytes = 2 * 1024 * 1024;
export const learningThumbnailAllowedExtensions = ['jpg', 'jpeg', 'png', 'webp', 'gif'] as const;
export const learningThumbnailAllowedMimeTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/gif'];
export const learningThumbnailMaxBytes = 5 * 1024 * 1024;

export function formatVideoLimit(bytes: number): string {
  return `${Math.round(bytes / 1024 / 1024)}MB`;
}

export interface LearningVideoValidationMessages {
  videoTypeError: string;
  videoSizeError: string;
  videoDurationError: string;
  videoMetadataError: string;
  subtitleTypeError: string;
  subtitleSizeError: string;
  thumbnailTypeError: string;
  thumbnailSizeError: string;
}

export function isLearningVideoFile(file: File): boolean {
  const extension = file.name.split('.').pop()?.toLowerCase() ?? '';
  return learningVideoAllowedExtensions.includes(extension as (typeof learningVideoAllowedExtensions)[number]);
}

export function validateLearningVideoFile(file: File, messages?: Pick<LearningVideoValidationMessages, 'videoTypeError' | 'videoSizeError' | 'videoDurationError' | 'videoMetadataError'>): Promise<string | null> {
  const extension = file.name.split('.').pop()?.toLowerCase() ?? '';
  if (!learningVideoAllowedExtensions.includes(extension as (typeof learningVideoAllowedExtensions)[number])) {
    return Promise.resolve(messages?.videoTypeError ?? '동영상은 mp4, webm, mov 파일만 업로드할 수 있습니다.');
  }
  if (file.type && !learningVideoAllowedMimeTypes.includes(file.type.toLowerCase())) {
    return Promise.resolve(messages?.videoTypeError ?? '동영상은 mp4, webm, mov 파일만 업로드할 수 있습니다.');
  }
  if (file.size > learningVideoMaxBytes) {
    return Promise.resolve(messages?.videoSizeError ?? `동영상은 ${formatVideoLimit(learningVideoMaxBytes)} 이하만 업로드할 수 있습니다.`);
  }

  return new Promise((resolve) => {
    const video = document.createElement('video');
    const objectURL = URL.createObjectURL(file);
    const cleanup = () => {
      URL.revokeObjectURL(objectURL);
      video.removeAttribute('src');
      video.load();
    };
    video.preload = 'metadata';
    video.onloadedmetadata = () => {
      const duration = Number.isFinite(video.duration) ? video.duration : 0;
      cleanup();
      if (duration > learningVideoMaxDurationSeconds) {
        resolve(messages?.videoDurationError ?? '동영상은 25분 이하만 업로드할 수 있습니다.');
        return;
      }
      resolve(null);
    };
    video.onerror = () => {
      cleanup();
      resolve(messages?.videoMetadataError ?? '동영상 정보를 읽지 못했습니다. mp4, webm, mov 파일인지 확인해 주세요.');
    };
    video.src = objectURL;
  });
}

export function validateLearningSubtitleFile(file: File, messages?: Pick<LearningVideoValidationMessages, 'subtitleTypeError' | 'subtitleSizeError'>): string | null {
  const extension = file.name.split('.').pop()?.toLowerCase() ?? '';
  if (!learningSubtitleAllowedExtensions.includes(extension as (typeof learningSubtitleAllowedExtensions)[number])) {
    return messages?.subtitleTypeError ?? '자막은 vtt 또는 srt 파일만 업로드할 수 있습니다.';
  }
  if (file.size > learningSubtitleMaxBytes) {
    return messages?.subtitleSizeError ?? `자막은 ${formatVideoLimit(learningSubtitleMaxBytes)} 이하만 업로드할 수 있습니다.`;
  }
  return null;
}

export function validateLearningThumbnailFile(file: File, messages?: Pick<LearningVideoValidationMessages, 'thumbnailTypeError' | 'thumbnailSizeError'>): string | null {
  const extension = file.name.split('.').pop()?.toLowerCase() ?? '';
  if (!learningThumbnailAllowedExtensions.includes(extension as (typeof learningThumbnailAllowedExtensions)[number])) {
    return messages?.thumbnailTypeError ?? '썸네일은 jpg, jpeg, png, webp, gif 이미지만 업로드할 수 있습니다.';
  }
  if (file.type && !learningThumbnailAllowedMimeTypes.includes(file.type.toLowerCase())) {
    return messages?.thumbnailTypeError ?? '썸네일은 jpg, jpeg, png, webp, gif 이미지만 업로드할 수 있습니다.';
  }
  if (file.size > learningThumbnailMaxBytes) {
    return messages?.thumbnailSizeError ?? `썸네일은 ${formatVideoLimit(learningThumbnailMaxBytes)} 이하만 업로드할 수 있습니다.`;
  }
  return null;
}

export async function normalizeSubtitleFile(file: File): Promise<File> {
  const extension = file.name.split('.').pop()?.toLowerCase() ?? '';
  if (extension !== 'srt') return file;
  const text = await file.text();
  const body = text
    .replace(/^\uFEFF/, '')
    .replace(/\r\n/g, '\n')
    .replace(/\r/g, '\n')
    .replace(/(\d{2}:\d{2}:\d{2}),(\d{3})/g, '$1.$2')
    .replace(/^\d+\n(?=\d{2}:\d{2}:\d{2}\.\d{3}\s+-->\s+\d{2}:\d{2}:\d{2}\.\d{3})/gm, '');
  const vtt = `WEBVTT\n\n${body.trim()}\n`;
  const filename = file.name.replace(/\.srt$/i, '.vtt');
  return new File([vtt], filename, { type: 'text/vtt' });
}
