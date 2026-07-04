import { SUPER_ADMIN_PAGE_WIDTH } from '@/components/super-admin/layout';
import type { CSSProperties } from 'react';
import { ATLAS_COLUMNS, ATLAS_ROWS, RECOMMENDED_HEIGHT, RECOMMENDED_WIDTH } from './constants';
import type { PlanetTextureMapAsset } from './types';

export const pageStyle: CSSProperties = {
  minHeight: '100vh',
  background: 'radial-gradient(circle at top, rgba(45, 78, 132, 0.2), transparent 42%), linear-gradient(180deg, #07111f, #0b1321 42%, #060b12)',
  color: '#eff6ff',
  padding: '28px 20px 48px',
};

export const shellStyle: CSSProperties = {
  width: SUPER_ADMIN_PAGE_WIDTH,
  margin: '0 auto',
  display: 'grid',
  gap: '18px',
};

export const summaryStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))',
  gap: 14,
};

export const summaryMetricStyle: CSSProperties = {
  display: 'grid',
  gap: 8,
  padding: 18,
  borderRadius: 18,
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(255,255,255,0.055)',
};

export const summaryNumberStyle: CSSProperties = {
  color: '#f8fbff',
  fontSize: 30,
  fontWeight: 900,
};

export const summaryLabelStyle: CSSProperties = {
  color: 'rgba(210, 222, 246, 0.78)',
  fontSize: 13,
  fontWeight: 800,
};

export const summaryCopyBoxStyle: CSSProperties = {
  ...summaryMetricStyle,
  gridColumn: 'span 2',
  color: 'rgba(210, 222, 246, 0.78)',
  fontSize: 13,
  lineHeight: 1.7,
};

export const gridStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1.08fr) minmax(320px, 0.92fr)',
  gap: 16,
  minWidth: 0,
};

export const previewGridStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'minmax(360px, 0.92fr) minmax(0, 1.08fr)',
  gap: 16,
  minWidth: 0,
};

export const panelStyle: CSSProperties = {
  display: 'grid',
  gap: 16,
  minWidth: 0,
  padding: 20,
  borderRadius: 20,
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(255,255,255,0.055)',
  boxShadow: '0 18px 42px rgba(0,0,0,0.18)',
};

export const panelHeaderStyle: CSSProperties = {
  display: 'flex',
  justifyContent: 'space-between',
  gap: 12,
  alignItems: 'center',
};

export const panelTitleStyle: CSSProperties = {
  fontSize: 18,
  color: '#f8fbff',
};

export const statusPillStyle: CSSProperties = {
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: 28,
  padding: '0 10px',
  borderRadius: 999,
  background: 'rgba(54, 211, 153, 0.12)',
  border: '1px solid rgba(54, 211, 153, 0.24)',
  color: '#bff9dc',
  fontSize: 12,
  fontWeight: 800,
};

export const fieldStyle: CSSProperties = {
  display: 'grid',
  gap: 8,
};

export const fieldLabelStyle: CSSProperties = {
  color: 'rgba(210, 222, 246, 0.82)',
  fontSize: 13,
  fontWeight: 800,
};

export const inputStyle: CSSProperties = {
  minHeight: 42,
  padding: '0 12px',
  borderRadius: 10,
  border: '1px solid rgba(150, 170, 220, 0.22)',
  background: 'rgba(8, 17, 30, 0.82)',
  color: '#f8fbff',
  font: 'inherit',
};

export const textareaStyle: CSSProperties = {
  ...inputStyle,
  minHeight: 96,
  padding: 12,
  resize: 'vertical',
};

export const fileDropStyle: CSSProperties = {
  display: 'grid',
  gap: 8,
  padding: 14,
  borderRadius: 14,
  border: '1px dashed rgba(150, 190, 255, 0.28)',
  background: 'rgba(8, 17, 30, 0.48)',
};

export const fileInputStyle: CSSProperties = {
  color: '#d9e6ff',
  fontSize: 13,
};

export const hiddenFileInputStyle: CSSProperties = {
  display: 'none',
};

export const fileHelpStyle: CSSProperties = {
  color: 'rgba(190, 204, 230, 0.68)',
  fontSize: 12,
};

export const specGridStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(2, minmax(0, 1fr))',
  gap: 10,
};

export const specItemStyle = (active: boolean): CSSProperties => ({
  display: 'grid',
  gap: 4,
  padding: 12,
  borderRadius: 12,
  border: `1px solid ${active ? 'rgba(54, 211, 153, 0.24)' : 'rgba(239, 68, 68, 0.22)'}`,
  background: active ? 'rgba(54, 211, 153, 0.08)' : 'rgba(239, 68, 68, 0.08)',
});

export const specLabelStyle: CSSProperties = {
  color: 'rgba(200, 214, 238, 0.72)',
  fontSize: 11,
  fontWeight: 800,
};

export const specValueStyle: CSSProperties = {
  color: '#f8fbff',
  fontSize: 14,
};

export const segmentedStyle: CSSProperties = {
  display: 'flex',
  flexWrap: 'wrap',
  gap: 8,
};

export const segmentButtonStyle = (active: boolean): CSSProperties => ({
  minHeight: 34,
  padding: '0 12px',
  borderRadius: 999,
  border: active ? '1px solid rgba(54, 211, 153, 0.38)' : '1px solid rgba(130, 150, 205, 0.2)',
  background: active ? 'rgba(54, 211, 153, 0.15)' : 'rgba(255,255,255,0.055)',
  color: active ? '#d7ffe9' : 'rgba(210,222,246,0.82)',
  font: 'inherit',
  fontSize: 12,
  fontWeight: 800,
  cursor: 'pointer',
});

export const presetButtonStyle = segmentButtonStyle;

export const rangeStyle: CSSProperties = {
  width: '100%',
};

export const toggleRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: 10,
  color: 'rgba(210, 222, 246, 0.82)',
  fontSize: 13,
  fontWeight: 800,
};

export const primaryButtonStyle: CSSProperties = {
  minHeight: 44,
  borderRadius: 12,
  border: '1px solid rgba(54, 211, 153, 0.34)',
  background: 'linear-gradient(135deg, rgba(42, 118, 255, 0.68), rgba(54, 211, 153, 0.7))',
  color: '#f8fbff',
  font: 'inherit',
  fontSize: 14,
  fontWeight: 900,
  cursor: 'pointer',
};

export const statusTextStyle: CSSProperties = {
  margin: 0,
  color: '#bff9dc',
  fontSize: 13,
  lineHeight: 1.6,
};

export const frameInfoStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(3, minmax(0, 1fr))',
  gap: 8,
  color: 'rgba(220, 235, 255, 0.82)',
  fontSize: 12,
  fontWeight: 800,
};

export const previewStageStyle: CSSProperties = {
  minHeight: 280,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  gap: 34,
  borderRadius: 18,
  background: 'radial-gradient(circle at 50% 48%, rgba(52, 90, 140, 0.22), transparent 46%), rgba(4, 9, 18, 0.74)',
  overflow: 'hidden',
};

export const controlGridStyle: CSSProperties = {
  display: 'grid',
  gap: 12,
};

export const atlasGridStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: `repeat(${ATLAS_COLUMNS}, minmax(0, 1fr))`,
  gap: 8,
  minWidth: 0,
  overflow: 'hidden',
};

export const textureMapListStyle: CSSProperties = {
  display: 'grid',
  gap: 14,
};

export const emptyStateStyle: CSSProperties = {
  padding: 24,
  borderRadius: 16,
  border: '1px dashed rgba(150, 175, 225, 0.22)',
  color: 'rgba(210, 222, 246, 0.68)',
  textAlign: 'center',
};

export const textureMapCardStyle: CSSProperties = {
  display: 'grid',
  gap: 14,
  padding: 16,
  borderRadius: 16,
  border: '1px solid rgba(194, 210, 245, 0.12)',
  background: 'rgba(8, 17, 30, 0.46)',
};

export const textureMapHeaderStyle: CSSProperties = {
  display: 'flex',
  justifyContent: 'space-between',
  gap: 12,
  alignItems: 'center',
  flexWrap: 'wrap',
};

export const textureMapTitleWrapStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: 10,
  flexWrap: 'wrap',
};

export const textureMapTitleStyle: CSSProperties = {
  color: '#f8fbff',
  fontSize: 16,
};

export const textureMapBadgeStyle = (active: boolean): CSSProperties => ({
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: 26,
  padding: '0 10px',
  borderRadius: 999,
  color: active ? '#dfffe8' : '#f3d8d8',
  background: active ? 'rgba(50, 200, 100, 0.16)' : 'rgba(220, 80, 90, 0.14)',
  border: active ? '1px solid rgba(50, 200, 100, 0.28)' : '1px solid rgba(220, 80, 90, 0.24)',
  fontSize: 12,
  fontWeight: 900,
});

export const textureMapMetaStyle: CSSProperties = {
  color: 'rgba(190, 204, 230, 0.68)',
  fontSize: 12,
};

export const textureMapPreviewRowStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '180px minmax(0, 1fr)',
  gap: 14,
  alignItems: 'stretch',
};

export const textureMapInfoStyle: CSSProperties = {
  display: 'grid',
  alignContent: 'start',
  gap: 10,
};

export const textureMapDescriptionStyle: CSSProperties = {
  margin: 0,
  color: 'rgba(220, 235, 255, 0.82)',
  fontSize: 13,
  lineHeight: 1.6,
};

export const textureMapMetaGridStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))',
  gap: 8,
  color: 'rgba(190, 204, 230, 0.72)',
  fontSize: 12,
  fontWeight: 800,
};

export const textureMapActionRowStyle: CSSProperties = {
  display: 'flex',
  gap: 10,
  flexWrap: 'wrap',
};

export const secondaryActionButtonStyle = (warning: boolean): CSSProperties => ({
  minHeight: 36,
  padding: '0 14px',
  borderRadius: 12,
  border: warning ? '1px solid rgba(220, 170, 80, 0.28)' : '1px solid rgba(130, 185, 255, 0.28)',
  background: warning ? 'rgba(220, 170, 80, 0.12)' : 'rgba(130, 185, 255, 0.12)',
  color: '#f7faff',
  font: 'inherit',
  fontSize: 13,
  fontWeight: 800,
  cursor: 'pointer',
});

export const dangerButtonStyle: CSSProperties = {
  minHeight: 36,
  padding: '0 14px',
  borderRadius: 12,
  border: '1px solid rgba(225, 92, 114, 0.26)',
  background: 'rgba(225, 92, 114, 0.12)',
  color: '#ffd9df',
  font: 'inherit',
  fontSize: 13,
  fontWeight: 800,
  cursor: 'pointer',
};

export const textureMapPathStyle: CSSProperties = {
  color: 'rgba(150, 170, 205, 0.75)',
  fontSize: 12,
  overflowWrap: 'anywhere',
};

export const atlasCellLabelStyle: CSSProperties = {
  position: 'absolute',
  left: 8,
  bottom: 6,
  padding: '3px 6px',
  borderRadius: 999,
  background: 'rgba(0,0,0,0.48)',
  color: '#eaf2ff',
  fontSize: 11,
  fontWeight: 900,
};

export function textureMapThumbStyle(textureMap: PlanetTextureMapAsset): CSSProperties {
  const versionSuffix = textureMap.version ? `?v=${textureMap.version}` : '';
  return {
    width: '100%',
    minWidth: 0,
    aspectRatio: `${RECOMMENDED_WIDTH} / ${RECOMMENDED_HEIGHT}`,
    borderRadius: 12,
    border: '1px solid rgba(150, 170, 220, 0.16)',
    backgroundColor: 'rgba(8, 17, 30, 0.88)',
    backgroundImage: `url(${textureMap.public_url}${versionSuffix})`,
    backgroundSize: 'contain',
    backgroundPosition: 'center',
    backgroundRepeat: 'no-repeat',
    overflow: 'hidden',
    cursor: 'pointer',
  };
}

export function atlasCellStyle(atlasURL: string | null, row: number, col: number, selected: boolean): CSSProperties {
  return {
    position: 'relative',
    width: '100%',
    minWidth: 0,
    aspectRatio: '2 / 1',
    borderRadius: 12,
    border: selected ? '2px solid rgba(54, 211, 153, 0.92)' : '1px solid rgba(150, 170, 220, 0.16)',
    backgroundColor: 'rgba(8, 17, 30, 0.88)',
    backgroundImage: atlasURL ? `url(${atlasURL})` : 'linear-gradient(135deg, rgba(34, 54, 84, 0.8), rgba(10, 18, 30, 0.9))',
    backgroundSize: `${ATLAS_COLUMNS * 100}% ${ATLAS_ROWS * 100}%`,
    backgroundPosition: `${(col / (ATLAS_COLUMNS - 1)) * 100}% ${(row / (ATLAS_ROWS - 1)) * 100}%`,
    backgroundRepeat: 'no-repeat',
    overflow: 'hidden',
    cursor: 'pointer',
    boxShadow: selected ? '0 0 0 2px rgba(54, 211, 153, 0.16), 0 0 24px rgba(54, 211, 153, 0.18)' : 'none',
  };
}
