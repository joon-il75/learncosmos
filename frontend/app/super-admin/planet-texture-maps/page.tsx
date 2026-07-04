'use client';

import SuperAdminPanelHeader from '@/components/super-admin/SuperAdminPanelHeader';
import SuperAdminPanelNav from '@/components/super-admin/SuperAdminPanelNav';
import {
  TextureMapAtlasPanel,
  TextureMapListPanel,
  TextureMapRotationPreviewPanel,
  TextureMapStagePanel,
  TextureMapUploadPanel,
} from './PlanetTextureMapPanels';
import {
  gridStyle,
  pageStyle,
  previewGridStyle,
  shellStyle,
  summaryCopyBoxStyle,
  summaryLabelStyle,
  summaryMetricStyle,
  summaryNumberStyle,
  summaryStyle,
} from './planetTextureMapStyles';
import { resolveTextureFrame } from './planetTextureMapUtils';
import { usePlanetTextureMapsAdmin } from './usePlanetTextureMapsAdmin';

export default function SuperAdminPlanetTextureMapsPage() {
  const {
    activeCount,
    applyRegisteredTextureMap,
    atlasIsRecommended,
    busyTextureMapId,
    description,
    fileInputRef,
    fileName,
    handleDelete,
    handleFileChange,
    handleReplaceFileChange,
    handleToggleActive,
    handleUpload,
    imageSize,
    isActive,
    isAuthorized,
    isPlaying,
    isUploading,
    name,
    previewURL,
    progress,
    replaceFileInputRefs,
    rotationDirection,
    rotationDuration,
    selectedFile,
    setDescription,
    setIsActive,
    setIsPlaying,
    setName,
    setProgress,
    setRotationDirection,
    setRotationDuration,
    statusMessage,
    textureMaps,
  } = usePlanetTextureMapsAdmin();
  const resolvedFrame = resolveTextureFrame(progress);

  if (isAuthorized === null) {
    return <main style={pageStyle}><div style={shellStyle}>권한을 확인하는 중입니다.</div></main>;
  }

  return (
    <main style={pageStyle}>
      <style jsx global>{`
        @keyframes planetSurfaceSpinLeft {
          from { transform: translateX(0); }
          to { transform: translateX(-50%); }
        }
        @keyframes planetSurfaceSpinRight {
          from { transform: translateX(-50%); }
          to { transform: translateX(0); }
        }
        @media (prefers-reduced-motion: reduce) {
          .planetSurfaceSpin {
            animation: none !important;
          }
        }
      `}</style>

      <div style={shellStyle}>
        <SuperAdminPanelHeader
          subtitle="행성 텍스처맵"
          description="자전 가능한 4x3 행성 표면 텍스처맵을 등록하고 실제 변화 시점별 프리뷰를 확인합니다."
        />
        <SuperAdminPanelNav activeSection="experience" activeExperience="planet-texture-maps" />

        <section style={summaryStyle}>
          <div style={summaryMetricStyle}>
            <strong style={summaryNumberStyle}>{textureMaps.length}</strong>
            <span style={summaryLabelStyle}>등록된 맵</span>
          </div>
          <div style={summaryMetricStyle}>
            <strong style={summaryNumberStyle}>{activeCount}</strong>
            <span style={summaryLabelStyle}>활성 맵</span>
          </div>
          <div style={summaryCopyBoxStyle}>
            기본맵은 현재 public atlas 파일로 시드되며, 새 WebP를 등록하면 같은 기준으로 자전 프리뷰와 실제 변화 시점 검토가 가능합니다.
          </div>
        </section>

        <section style={gridStyle}>
          <TextureMapUploadPanel
            atlasIsRecommended={atlasIsRecommended}
            description={description}
            fileInputRef={fileInputRef}
            fileName={fileName}
            handleFileChange={handleFileChange}
            handleUpload={handleUpload}
            imageSize={imageSize}
            isActive={isActive}
            isUploading={isUploading}
            name={name}
            selectedFile={selectedFile}
            setDescription={setDescription}
            setIsActive={setIsActive}
            setName={setName}
            statusMessage={statusMessage}
          />
          <TextureMapStagePanel
            progress={progress}
            resolvedFrame={resolvedFrame}
            setProgress={setProgress}
          />
        </section>

        <section style={previewGridStyle}>
          <TextureMapRotationPreviewPanel
            isPlaying={isPlaying}
            previewURL={previewURL}
            progress={progress}
            rotationDirection={rotationDirection}
            rotationDuration={rotationDuration}
            setIsPlaying={setIsPlaying}
            setRotationDirection={setRotationDirection}
            setRotationDuration={setRotationDuration}
          />
          <TextureMapAtlasPanel
            previewURL={previewURL}
            resolvedFrame={resolvedFrame}
            setProgress={setProgress}
          />
        </section>

        <TextureMapListPanel
          applyRegisteredTextureMap={applyRegisteredTextureMap}
          busyTextureMapId={busyTextureMapId}
          handleDelete={handleDelete}
          handleReplaceFileChange={handleReplaceFileChange}
          handleToggleActive={handleToggleActive}
          replaceFileInputRefs={replaceFileInputRefs}
          textureMaps={textureMaps}
        />
      </div>
    </main>
  );
}
