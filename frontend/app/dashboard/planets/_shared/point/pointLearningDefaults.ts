'use client';

import type { PointAttachmentDraft } from '../pointPageUtils';

export type UserInfo = { id: string; required_consent_pending: boolean; ui_locale?: 'ko' | 'en'; language_setup_required?: boolean };

export const EMPTY_ATTACHMENT: PointAttachmentDraft = {
  id: '',
  provider: 'learner',
  attachmentType: 'link',
  sourceContext: 'work_attachment',
  artifactID: null,
  title: '',
  url: '',
  filePath: '',
  fileSize: '',
  mimeType: '',
};

export const EMPTY_RESEARCH_MATERIAL_ATTACHMENT: PointAttachmentDraft = {
  ...EMPTY_ATTACHMENT,
  sourceContext: 'research_material',
};

export const COMPLETION_CHECK_COUNT = 6;
