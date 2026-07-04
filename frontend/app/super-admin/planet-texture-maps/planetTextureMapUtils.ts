import { getProgressGridPosition } from '@/lib/world-ui-engine/planetProgressUtils';
import { ATLAS_COLUMNS, stageProgressValues } from './constants';
import type { TextureFrame } from './types';

export function resolveTextureFrame(progress: number): TextureFrame {
  const safeProgress = Math.max(0, Math.min(100, progress));
  const grid = getProgressGridPosition(safeProgress);
  return {
    row: (grid.row - 1) as 0 | 1,
    col: (grid.col - 1) as 0 | 1 | 2 | 3,
    stage: grid.stage,
    label: `학습 완료율 ${safeProgress}%`,
  };
}

export function applyFrameSelection(
  row: number,
  col: number,
  setProgress: (progress: number) => void,
) {
  setProgress(stageProgressValues[row * ATLAS_COLUMNS + col] ?? 0);
}

export function formatDate(value?: string): string {
  if (!value) return '-';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleString('ko-KR');
}

export function formatFileSize(value: number): string {
  if (!value) return '0 B';
  if (value < 1024) return `${value} B`;
  if (value < 1024 * 1024) return `${(value / 1024).toFixed(1)} KB`;
  return `${(value / (1024 * 1024)).toFixed(2)} MB`;
}
