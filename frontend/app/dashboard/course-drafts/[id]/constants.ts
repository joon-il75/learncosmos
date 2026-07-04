'use client';

import type { DiarySectionKey } from './types';
import type { ResearchNodeTemplateType } from './types';

export const researchNodeTemplateOptions: Array<{
  value: ResearchNodeTemplateType;
}> = [
  { value: 'concept_summary' },
  { value: 'practice_strategy' },
  { value: 'problem_solving' },
  { value: 'free_research' },
] as const;

export const diarySections: Array<{
  key: DiarySectionKey;
  isLocked?: boolean;
}> = [
  { key: 'planning' },
  { key: 'journal' },
  { key: 'records' },
  { key: 'artifacts' },
  { key: 'community' },
  { key: 'civilization' },
] as const;
