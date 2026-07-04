import type { ChangeEvent, RefObject } from 'react';
import { formatDate, formatFileSize } from './planetTextureMapUtils';
import {
  dangerButtonStyle,
  hiddenFileInputStyle,
  secondaryActionButtonStyle,
  textureMapActionRowStyle,
  textureMapBadgeStyle,
  textureMapCardStyle,
  textureMapDescriptionStyle,
  textureMapHeaderStyle,
  textureMapInfoStyle,
  textureMapMetaGridStyle,
  textureMapMetaStyle,
  textureMapPathStyle,
  textureMapPreviewRowStyle,
  textureMapThumbStyle,
  textureMapTitleStyle,
  textureMapTitleWrapStyle,
} from './planetTextureMapStyles';
import type { PlanetTextureMapAsset } from './types';

type TextureMapCardProps = {
  applyRegisteredTextureMap: (textureMap: PlanetTextureMapAsset) => void;
  busyTextureMapId: string | null;
  handleDelete: (textureMap: PlanetTextureMapAsset) => void;
  handleReplaceFileChange: (textureMap: PlanetTextureMapAsset, event: ChangeEvent<HTMLInputElement>) => void;
  handleToggleActive: (textureMap: PlanetTextureMapAsset) => void;
  replaceFileInputRefs: RefObject<Record<string, HTMLInputElement | null>>;
  textureMap: PlanetTextureMapAsset;
};

export function TextureMapCard({
  applyRegisteredTextureMap,
  busyTextureMapId,
  handleDelete,
  handleReplaceFileChange,
  handleToggleActive,
  replaceFileInputRefs,
  textureMap,
}: TextureMapCardProps) {
  return (
    <article style={textureMapCardStyle}>
      <div style={textureMapHeaderStyle}>
        <div style={textureMapTitleWrapStyle}>
          <strong style={textureMapTitleStyle}>{textureMap.name}</strong>
          <span style={textureMapBadgeStyle(textureMap.is_active)}>
            {textureMap.is_active ? '활성' : '비활성'}
          </span>
        </div>
        <span style={textureMapMetaStyle}>{formatFileSize(textureMap.size_bytes)}</span>
      </div>

      <div style={textureMapPreviewRowStyle}>
        <button
          type="button"
          onClick={() => applyRegisteredTextureMap(textureMap)}
          style={textureMapThumbStyle(textureMap)}
          title={`${textureMap.name} 프리뷰`}
        />
        <div style={textureMapInfoStyle}>
          <p style={textureMapDescriptionStyle}>{textureMap.description || '설명 없음'}</p>
          <div style={textureMapMetaGridStyle}>
            <span>{textureMap.width} × {textureMap.height}</span>
            <span>{textureMap.columns}열 × {textureMap.rows}행</span>
            <span>Cell {textureMap.cell_width} × {textureMap.cell_height}</span>
            <span>자전 {textureMap.rotation_duration_seconds}초 · {textureMap.rotation_direction === 'left' ? '좌회전' : '우회전'}</span>
          </div>
        </div>
      </div>

      <div style={textureMapActionRowStyle}>
        <input
          ref={(node) => {
            replaceFileInputRefs.current[textureMap.id] = node;
          }}
          type="file"
          accept=".png,.jpg,.jpeg,.webp,image/png,image/jpeg,image/webp"
          onChange={(event) => handleReplaceFileChange(textureMap, event)}
          style={hiddenFileInputStyle}
        />
        <button
          type="button"
          onClick={() => applyRegisteredTextureMap(textureMap)}
          style={secondaryActionButtonStyle(false)}
        >
          프리뷰
        </button>
        <button
          type="button"
          onClick={() => replaceFileInputRefs.current[textureMap.id]?.click()}
          disabled={busyTextureMapId === textureMap.id || textureMap.is_builtin}
          style={secondaryActionButtonStyle(false)}
        >
          {textureMap.is_builtin ? '기본맵 보호' : '이미지 교체'}
        </button>
        <button
          type="button"
          onClick={() => handleToggleActive(textureMap)}
          disabled={busyTextureMapId === textureMap.id}
          style={secondaryActionButtonStyle(textureMap.is_active)}
        >
          {textureMap.is_active ? '비활성화' : '활성화'}
        </button>
        <button
          type="button"
          onClick={() => handleDelete(textureMap)}
          disabled={busyTextureMapId === textureMap.id || textureMap.is_builtin}
          style={dangerButtonStyle}
        >
          {textureMap.is_builtin ? '기본맵' : '삭제'}
        </button>
      </div>

      <div style={textureMapPathStyle}>경로: {textureMap.public_url} · 업데이트: {formatDate(textureMap.updated_at)}</div>
    </article>
  );
}
