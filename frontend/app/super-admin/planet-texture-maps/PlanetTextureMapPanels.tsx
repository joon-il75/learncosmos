import AtmosphericPlanetPreview from '@/components/dashboard/AtmosphericPlanetPreview';
import type { ChangeEvent, Dispatch, RefObject, SetStateAction } from 'react';
import {
  ATLAS_COLUMNS,
  stageLabels,
  rowLabels,
  stageProgressValues,
  VISIBLE_STAGE_ROWS,
} from './constants';
import { TextureMapCard } from './TextureMapCard';
import { TextureMapUploadFields } from './TextureMapUploadFields';
import { applyFrameSelection } from './planetTextureMapUtils';
import {
  atlasCellLabelStyle,
  atlasCellStyle,
  atlasGridStyle,
  controlGridStyle,
  emptyStateStyle,
  fieldLabelStyle,
  fieldStyle,
  frameInfoStyle,
  panelHeaderStyle,
  panelStyle,
  panelTitleStyle,
  presetButtonStyle,
  previewStageStyle,
  primaryButtonStyle,
  rangeStyle,
  secondaryActionButtonStyle,
  segmentedStyle,
  segmentButtonStyle,
  statusPillStyle,
  statusTextStyle,
  textureMapListStyle,
  toggleRowStyle,
} from './planetTextureMapStyles';
import type { PlanetTextureMapAsset, RotationDirection, TextureFrame } from './types';

type UploadPanelProps = {
  atlasIsRecommended: boolean | null | undefined;
  description: string;
  fileInputRef: RefObject<HTMLInputElement | null>;
  fileName: string | null;
  handleFileChange: (event: ChangeEvent<HTMLInputElement>) => void;
  handleUpload: () => void;
  imageSize: { width: number; height: number } | null;
  isActive: boolean;
  isUploading: boolean;
  name: string;
  selectedFile: File | null;
  setDescription: (value: string) => void;
  setIsActive: (value: boolean) => void;
  setName: (value: string) => void;
  statusMessage: string | null;
};

type StagePanelProps = {
  progress: number;
  resolvedFrame: TextureFrame;
  setProgress: (value: number) => void;
};

type RotationPreviewPanelProps = {
  isPlaying: boolean;
  previewURL: string | null;
  progress: number;
  rotationDirection: RotationDirection;
  rotationDuration: number;
  setIsPlaying: Dispatch<SetStateAction<boolean>>;
  setRotationDirection: (value: RotationDirection) => void;
  setRotationDuration: (value: number) => void;
};

type AtlasPanelProps = {
  previewURL: string | null;
  resolvedFrame: TextureFrame;
  setProgress: (value: number) => void;
};

type TextureMapListPanelProps = {
  applyRegisteredTextureMap: (textureMap: PlanetTextureMapAsset) => void;
  busyTextureMapId: string | null;
  handleDelete: (textureMap: PlanetTextureMapAsset) => void;
  handleReplaceFileChange: (textureMap: PlanetTextureMapAsset, event: ChangeEvent<HTMLInputElement>) => void;
  handleToggleActive: (textureMap: PlanetTextureMapAsset) => void;
  replaceFileInputRefs: RefObject<Record<string, HTMLInputElement | null>>;
  textureMaps: PlanetTextureMapAsset[];
};

export function TextureMapUploadPanel({
  atlasIsRecommended,
  description,
  fileInputRef,
  fileName,
  handleFileChange,
  handleUpload,
  imageSize,
  isActive,
  isUploading,
  name,
  selectedFile,
  setDescription,
  setIsActive,
  setName,
  statusMessage,
}: UploadPanelProps) {
  return (
    <article style={panelStyle}>
      <div style={panelHeaderStyle}>
        <strong style={panelTitleStyle}>텍스처맵 후보</strong>
        <span style={statusPillStyle}>{selectedFile ? '업로드 대기' : '등록 맵 프리뷰'}</span>
      </div>

      <TextureMapUploadFields
        atlasIsRecommended={atlasIsRecommended}
        description={description}
        fileInputRef={fileInputRef}
        fileName={fileName}
        handleFileChange={handleFileChange}
        imageSize={imageSize}
        name={name}
        setDescription={setDescription}
        setName={setName}
      />

      <label style={toggleRowStyle}>
        <input type="checkbox" checked={isActive} onChange={(event) => setIsActive(event.target.checked)} />
        <span>등록 직후 활성 맵으로 포함</span>
      </label>

      <button type="button" onClick={handleUpload} disabled={isUploading} style={primaryButtonStyle}>
        {isUploading ? '등록 중...' : '텍스처맵 등록'}
      </button>

      {statusMessage ? <p style={statusTextStyle}>{statusMessage}</p> : null}
    </article>
  );
}

export function TextureMapStagePanel({ progress, resolvedFrame, setProgress }: StagePanelProps) {
  return (
    <article style={panelStyle}>
      <div style={panelHeaderStyle}>
        <strong style={panelTitleStyle}>완료율 8단계</strong>
        <span style={statusPillStyle}>Row {resolvedFrame.row + 1} · Col {resolvedFrame.col + 1}</span>
      </div>

      <label style={fieldStyle}>
        <span style={fieldLabelStyle}>학습 완료율 {progress}%</span>
        <input
          type="range"
          min={0}
          max={100}
          step={1}
          value={progress}
          onChange={(event) => setProgress(Number(event.target.value))}
          style={rangeStyle}
        />
      </label>

      <div style={segmentedStyle}>
        {stageProgressValues.map((value, index) => (
          <button key={value} type="button" onClick={() => setProgress(value)} style={presetButtonStyle(progress === value)}>
            {stageLabels[index]}
          </button>
        ))}
      </div>

      <div style={frameInfoStyle}>
        <span>{rowLabels[resolvedFrame.row]}</span>
        <span>{stageLabels[resolvedFrame.stage]}</span>
        <span>{resolvedFrame.label}</span>
      </div>
    </article>
  );
}

export function TextureMapRotationPreviewPanel({
  isPlaying,
  previewURL,
  progress,
  rotationDirection,
  rotationDuration,
  setIsPlaying,
  setRotationDirection,
  setRotationDuration,
}: RotationPreviewPanelProps) {
  return (
    <article style={panelStyle}>
      <div style={panelHeaderStyle}>
        <strong style={panelTitleStyle}>자전 프리뷰</strong>
        <span style={statusPillStyle}>{isPlaying ? '재생 중' : '정지'}</span>
      </div>

      <div style={previewStageStyle}>
        <AtmosphericPlanetPreview
          atlasURL={previewURL}
          size={220}
          progressPercent={progress}
          durationSeconds={rotationDuration}
          direction={rotationDirection}
          playing={isPlaying}
        />
        <AtmosphericPlanetPreview
          atlasURL={previewURL}
          size={68}
          progressPercent={progress}
          durationSeconds={Math.round(rotationDuration * 1.4)}
          direction={rotationDirection}
          playing={isPlaying}
        />
      </div>

      <div style={controlGridStyle}>
        <label style={fieldStyle}>
          <span style={fieldLabelStyle}>대표 자전 속도 {rotationDuration}초</span>
          <input
            type="range"
            min={18}
            max={90}
            step={1}
            value={rotationDuration}
            onChange={(event) => setRotationDuration(Number(event.target.value))}
            style={rangeStyle}
          />
        </label>
        <div style={segmentedStyle}>
          <button type="button" onClick={() => setRotationDirection('left')} style={segmentButtonStyle(rotationDirection === 'left')}>좌회전</button>
          <button type="button" onClick={() => setRotationDirection('right')} style={segmentButtonStyle(rotationDirection === 'right')}>우회전</button>
          <button type="button" onClick={() => setIsPlaying((current) => !current)} style={segmentButtonStyle(isPlaying)}>
            {isPlaying ? '정지' : '재생'}
          </button>
        </div>
      </div>
    </article>
  );
}

export function TextureMapAtlasPanel({ previewURL, resolvedFrame, setProgress }: AtlasPanelProps) {
  return (
    <article style={panelStyle}>
      <div style={panelHeaderStyle}>
        <strong style={panelTitleStyle}>4x3 atlas</strong>
        <span style={statusPillStyle}>{previewURL ? '8단계 사용' : '파일 선택 전'}</span>
      </div>
      <div style={atlasGridStyle}>
        {Array.from({ length: ATLAS_COLUMNS * VISIBLE_STAGE_ROWS }, (_, index) => {
          const row = Math.floor(index / ATLAS_COLUMNS);
          const col = index % ATLAS_COLUMNS;
          const selected = row === resolvedFrame.row && col === resolvedFrame.col;
          return (
            <button
              key={`${row}-${col}`}
              type="button"
              onClick={() => applyFrameSelection(row, col, setProgress)}
              style={atlasCellStyle(previewURL, row, col, selected)}
              title={`${rowLabels[row]} · ${stageLabels[index]}`}
            >
              <span style={atlasCellLabelStyle}>{row + 1}-{col + 1}</span>
            </button>
          );
        })}
      </div>
    </article>
  );
}

export function TextureMapListPanel({
  applyRegisteredTextureMap,
  busyTextureMapId,
  handleDelete,
  handleReplaceFileChange,
  handleToggleActive,
  replaceFileInputRefs,
  textureMaps,
}: TextureMapListPanelProps) {
  return (
    <section style={panelStyle}>
      <div style={panelHeaderStyle}>
        <strong style={panelTitleStyle}>등록된 텍스처맵</strong>
        <span style={statusPillStyle}>운영 목록</span>
      </div>

      <div style={textureMapListStyle}>
        {textureMaps.length === 0 ? (
          <div style={emptyStateStyle}>아직 등록된 텍스처맵이 없습니다.</div>
        ) : (
          textureMaps.map((textureMap) => (
            <TextureMapCard
              key={textureMap.id}
              applyRegisteredTextureMap={applyRegisteredTextureMap}
              busyTextureMapId={busyTextureMapId}
              handleDelete={handleDelete}
              handleReplaceFileChange={handleReplaceFileChange}
              handleToggleActive={handleToggleActive}
              replaceFileInputRefs={replaceFileInputRefs}
              textureMap={textureMap}
            />
          ))
        )}
      </div>
    </section>
  );
}
