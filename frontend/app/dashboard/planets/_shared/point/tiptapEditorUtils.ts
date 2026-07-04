export const emptyContent = '<p></p>';
export const allowedInlineImageExtensions = ['jpg', 'jpeg', 'png', 'webp', 'gif'] as const;
export const allowedInlineImageMimeTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/gif'];
export const inlineImageMaxBytes = 5 * 1024 * 1024;
export const inlineImageAccept = allowedInlineImageExtensions.map((ext) => `.${ext}`).join(',');

export function formatFileSize(bytes: number): string {
  return `${Math.round((bytes / 1024 / 1024) * 10) / 10}MB`;
}

export function getInlineImageUploadHint(): string {
  return `업로드 가능: jpg, jpeg, png, webp, gif / ${formatFileSize(inlineImageMaxBytes)} 이하`;
}

export function validateInlineImageFile(file: File, messages?: { typeError: string; sizeError: string }): string | null {
  const extension = file.name.split('.').pop()?.toLowerCase() ?? '';
  if (!allowedInlineImageExtensions.includes(extension as (typeof allowedInlineImageExtensions)[number])) {
    return messages?.typeError ?? 'jpg, jpeg, png, webp, gif 이미지만 업로드할 수 있습니다.';
  }
  if (file.type && !allowedInlineImageMimeTypes.includes(file.type.toLowerCase())) {
    return messages?.typeError ?? 'jpg, jpeg, png, webp, gif 이미지만 업로드할 수 있습니다.';
  }
  if (file.size > inlineImageMaxBytes) {
    return messages?.sizeError ?? `이미지는 ${formatFileSize(inlineImageMaxBytes)} 이하만 업로드할 수 있습니다.`;
  }
  return null;
}

export function normalizeContent(value: string): string {
  return value.trim().length ? value : emptyContent;
}

export function isHttpsUrl(value: string): boolean {
  const trimmed = value.trim();
  if (!trimmed.startsWith('https://')) return false;
  try {
    return new URL(trimmed).protocol === 'https:';
  } catch {
    return false;
  }
}
