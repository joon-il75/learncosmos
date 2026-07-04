'use client';

import { useMemo } from 'react';
import LumiAvatar from '@/components/lumi/LumiAvatar';
import LumiBubble from '@/components/lumi/LumiBubble';
import LumiDock from '@/components/lumi/LumiDock';
import LumiPanel from '@/components/lumi/LumiPanel';
import {
  lumiSpriteMap,
  LUMI_SPRITE_COLUMNS, LUMI_SPRITE_ROWS,
  LUMI_SPRITE_CELL_WIDTH, LUMI_SPRITE_CELL_HEIGHT,
  LUMI_SPRITE_SHEET_WIDTH, LUMI_SPRITE_SHEET_HEIGHT,
  getDefaultSpriteTuning,
} from '@/lib/lumi/lumiSpriteMap';
import type { LumiState } from '@/lib/lumi/lumiTypes';
import { stateOptions, contextOptions, dockSlotOptions, messageTypeOptions, triggerOptions } from '../lumiLabTypes';
import type { LumiLabAssetsResult } from '../useLumiLabAssets';
import type { LumiLabUIPreviewResult } from '../useLumiLabUIPreview';
import {
  assetPanelStyle, assetColumnStyle, assetMetaCardStyle, assetPreviewWrapStyle, assetPreviewStyle,
  assetMetaTextStyle, gridSpecStyle, uploadBoxStyle, fileInputStyle, uploadActionRowStyle,
  uploadHintStyle, assetStatusStyle, promptCardStyle, promptHeaderStyle, promptTitleStyle,
  copyButtonStyle, promptTextareaStyle, promptFootnoteStyle,
  layoutStyle, controlsPanelStyle, previewPanelStyle, sectionTitleStyle, fieldStyle, labelStyle,
  selectStyle, textareaStyle, rangeGroupStyle, toggleGridStyle, triggerGroupStyle,
  triggerButtonsStyle, triggerButtonStyle, previewGridStyle, previewCardStyle, miniLabelStyle,
  codeStyle, dockStageStyle, dockHeaderStyle, codePillStyle, panelPreviewRowStyle,
  panelColumnStyle, sheetColumnStyle, spriteSheetFrameStyle, spriteSheetStyle, sheetHelpStyle,
  primaryButtonStyle,
} from '../lumiLabStyles';
import { RangeField, ToggleChip } from '../LumiLabComponents';

interface Props {
  assets: LumiLabAssetsResult;
  uip: LumiLabUIPreviewResult;
}

export default function LumiSpritesTab({ assets, uip }: Props) {
  const {
    assetMeta, assetVersion, assetStatus, selectedUpload, setSelectedUpload,
    isUploading, copiedPromptId, handleUpload, handleCopyPrompt,
  } = assets;
  const {
    lumi, selectedState, setSelectedState, selectedContext, setSelectedContext,
    selectedDockSlot, setSelectedDockSlot, messageType, setMessageType, message, setMessage,
    size, setSize, manualVisible, setManualVisible, manualExpanded, setManualExpanded,
    manualMobile, setManualMobile, backgroundSizeX, setBackgroundSizeX, backgroundSizeY,
    setBackgroundSizeY, backgroundPositionX, setBackgroundPositionX, backgroundPositionY,
    setBackgroundPositionY, manualSpriteStyle, spriteCoordinate, selectedStateTuning,
    handleTrigger,
  } = uip;

  const spriteSheetPrompt = useMemo(() => {
    const rows = stateOptions.map((state) => {
      const coord = lumiSpriteMap[state];
      const tuning = getDefaultSpriteTuning(state);
      return `- ${state}: row ${coord.row}, col ${coord.col}, background-size ${tuning.backgroundSizeX}% ${tuning.backgroundSizeY}%, background-position ${tuning.backgroundPositionX}% ${tuning.backgroundPositionY}%`;
    });
    return [
      'Create a single transparent WebP sprite sheet for LearnWeaver Lumi.',
      `Canvas layout: ${LUMI_SPRITE_ROWS} rows x ${LUMI_SPRITE_COLUMNS} columns, cell size ${LUMI_SPRITE_CELL_WIDTH}x${LUMI_SPRITE_CELL_HEIGHT}px, total sheet ${LUMI_SPRITE_SHEET_WIDTH}x${LUMI_SPRITE_SHEET_HEIGHT}px, transparent background, no outer frame, no background scene.`,
      'Character direction: soft sci-fi learning companion, calm Korean edtech product tone, not cartoonishly childish, consistent body proportions and lighting across all cells.',
      'Each cell must contain only Lumi, centered and fully visible, with no cropped limbs.',
      'Keep the existing LearnWeaver Lumi visual identity and color family.',
      'Required states and exact cell mapping:',
      ...rows,
      'Export as a single sprite sheet that can replace /public/images/lumi.webp directly.',
    ].join('\n');
  }, []);

  const selectedStatePrompt = useMemo(() => {
    const coord = spriteCoordinate;
    const tuning = getDefaultSpriteTuning(selectedState);
    return [
      `Generate the LearnWeaver Lumi "${selectedState}" pose for a ${LUMI_SPRITE_ROWS}x${LUMI_SPRITE_COLUMNS} transparent sprite sheet (${LUMI_SPRITE_SHEET_WIDTH}x${LUMI_SPRITE_SHEET_HEIGHT}px, cell ${LUMI_SPRITE_CELL_WIDTH}x${LUMI_SPRITE_CELL_HEIGHT}px).`,
      `Target cell: row ${coord.row}, col ${coord.col}.`,
      `Current preview calibration: background-size ${tuning.backgroundSizeX}% ${tuning.backgroundSizeY}%, background-position ${tuning.backgroundPositionX}% ${tuning.backgroundPositionY}%.`,
      'Match the existing Lumi style exactly: soft explorer companion, polished expression, clean silhouette, transparent background, no props outside the intended pose.',
      'Keep the face angle, body scale, and lighting consistent with the rest of the sprite sheet.',
    ].join('\n');
  }, [selectedState, spriteCoordinate]);

  return (
    <>
      <section style={assetPanelStyle}>
        <div style={assetColumnStyle}>
          <h2 style={sectionTitleStyle}>Asset Control</h2>
          <div style={assetMetaCardStyle}>
            <div style={assetPreviewWrapStyle}>
              <img src={`/images/lumi.webp?v=${assetVersion}`} alt="현재 Lumi 스프라이트 시트" style={assetPreviewStyle} />
            </div>
            <div style={assetMetaTextStyle}>
              <div><strong>현재 자산:</strong> {assetMeta?.path ?? '미확인'}</div>
              <div><strong>업데이트:</strong> {assetMeta?.updated_at ?? '미확인'}</div>
              <div><strong>파일 크기:</strong> {assetMeta ? `${Math.round(assetMeta.size_bytes / 1024)} KB` : '미확인'}</div>
              <div style={gridSpecStyle}>
                <strong>그리드 규격:</strong>
                {` ${LUMI_SPRITE_SHEET_WIDTH} × ${LUMI_SPRITE_SHEET_HEIGHT}px`}
                {` / 셀 ${LUMI_SPRITE_CELL_WIDTH} × ${LUMI_SPRITE_CELL_HEIGHT}px`}
                {` / ${LUMI_SPRITE_COLUMNS}열 × ${LUMI_SPRITE_ROWS}행 = ${LUMI_SPRITE_COLUMNS * LUMI_SPRITE_ROWS}상태`}
              </div>
              <div><strong>운영 규칙:</strong> {`${LUMI_SPRITE_CELL_WIDTH}×${LUMI_SPRITE_CELL_HEIGHT}px 균일 셀 기준 WebP — 교체 시 동일 규격 유지`}</div>
            </div>
          </div>
          <div style={uploadBoxStyle}>
            <label style={labelStyle}>Lumi 스프라이트 업로드</label>
            <input type="file" accept=".webp,image/webp" onChange={(event) => setSelectedUpload(event.target.files?.[0] ?? null)} style={fileInputStyle} />
            <div style={uploadActionRowStyle}>
              <span style={uploadHintStyle}>
                {selectedUpload ? `${selectedUpload.name} 선택됨` : `${LUMI_SPRITE_SHEET_WIDTH}×${LUMI_SPRITE_SHEET_HEIGHT}px · ${LUMI_SPRITE_COLUMNS}×${LUMI_SPRITE_ROWS} 투명 배경 WebP를 선택하세요.`}
              </span>
              <button type="button" onClick={handleUpload} disabled={isUploading} style={primaryButtonStyle}>
                {isUploading ? '업로드 중...' : '스프라이트 교체'}
              </button>
            </div>
            {assetStatus ? <div style={assetStatusStyle}>{assetStatus}</div> : null}
          </div>
        </div>

        <div style={assetColumnStyle}>
          <h2 style={sectionTitleStyle}>Prompt Library</h2>
          <div style={promptCardStyle}>
            <div style={promptHeaderStyle}>
              <strong style={promptTitleStyle}>전체 스프라이트 시트 생성 프롬프트</strong>
              <button type="button" onClick={() => handleCopyPrompt('sheet', spriteSheetPrompt)} style={copyButtonStyle}>
                {copiedPromptId === 'sheet' ? '복사됨' : '복사'}
              </button>
            </div>
            <textarea value={spriteSheetPrompt} readOnly rows={14} style={promptTextareaStyle} />
          </div>
          <div style={promptCardStyle}>
            <div style={promptHeaderStyle}>
              <strong style={promptTitleStyle}>{selectedState} 상태 개별 프롬프트</strong>
              <button type="button" onClick={() => handleCopyPrompt('state', selectedStatePrompt)} style={copyButtonStyle}>
                {copiedPromptId === 'state' ? '복사됨' : '복사'}
              </button>
            </div>
            <textarea value={selectedStatePrompt} readOnly rows={8} style={promptTextareaStyle} />
            <div style={promptFootnoteStyle}>
              현재 보정값:{' '}
              {selectedStateTuning ? `background-size ${selectedStateTuning.backgroundSizeX}% ${selectedStateTuning.backgroundSizeY}%, background-position ${selectedStateTuning.backgroundPositionX}% ${selectedStateTuning.backgroundPositionY}%` : '기본 grid 계산값'}
            </div>
          </div>
        </div>
      </section>

      <div style={layoutStyle}>
        <section style={controlsPanelStyle}>
          <h2 style={sectionTitleStyle}>Controls</h2>
          <label style={fieldStyle}>
            <span style={labelStyle}>Lumi 상태</span>
            <select value={selectedState} onChange={(event) => setSelectedState(event.target.value as LumiState)} style={selectStyle}>
              {stateOptions.map((option) => <option key={option} value={option}>{option}</option>)}
            </select>
          </label>
          <label style={fieldStyle}>
            <span style={labelStyle}>Context</span>
            <select value={selectedContext} onChange={(event) => setSelectedContext(event.target.value as typeof selectedContext)} style={selectStyle}>
              {contextOptions.map((option) => <option key={option} value={option}>{option}</option>)}
            </select>
          </label>
          <label style={fieldStyle}>
            <span style={labelStyle}>Dock Slot</span>
            <select value={selectedDockSlot} onChange={(event) => setSelectedDockSlot(event.target.value as typeof selectedDockSlot)} style={selectStyle}>
              {dockSlotOptions.map((option) => <option key={option} value={option}>{option}</option>)}
            </select>
          </label>
          <label style={fieldStyle}>
            <span style={labelStyle}>Message Type</span>
            <select value={messageType} onChange={(event) => setMessageType(event.target.value as typeof messageType)} style={selectStyle}>
              {messageTypeOptions.map((option) => <option key={option} value={option}>{option}</option>)}
            </select>
          </label>
          <label style={fieldStyle}>
            <span style={labelStyle}>메시지</span>
            <textarea value={message} onChange={(event) => setMessage(event.target.value)} rows={4} style={textareaStyle} />
          </label>
          <div style={rangeGroupStyle}>
            <RangeField label={`프리뷰 크기 ${size}px`} min={56} max={160} value={size} onChange={setSize} />
            <RangeField label={`background-size X ${backgroundSizeX}%`} min={350} max={650} value={backgroundSizeX} onChange={setBackgroundSizeX} />
            <RangeField label={`background-size Y ${backgroundSizeY}%`} min={220} max={420} value={backgroundSizeY} onChange={setBackgroundSizeY} />
            <RangeField label={`background-position X ${backgroundPositionX}%`} min={-10} max={120} value={backgroundPositionX} onChange={setBackgroundPositionX} />
            <RangeField label={`background-position Y ${backgroundPositionY}%`} min={-10} max={120} value={backgroundPositionY} onChange={setBackgroundPositionY} />
          </div>
          <div style={toggleGridStyle}>
            <ToggleChip label="보이기" checked={manualVisible} onToggle={() => setManualVisible((v) => !v)} />
            <ToggleChip label="확장" checked={manualExpanded} onToggle={() => setManualExpanded((v) => !v)} />
            <ToggleChip label="모바일 모드" checked={manualMobile} onToggle={() => setManualMobile((v) => !v)} />
          </div>
          <div style={triggerGroupStyle}>
            <span style={labelStyle}>Trigger Presets</span>
            <div style={triggerButtonsStyle}>
              {triggerOptions.map((trigger) => (
                <button key={trigger} type="button" onClick={() => handleTrigger(trigger)} style={triggerButtonStyle}>{trigger}</button>
              ))}
            </div>
          </div>
        </section>

        <section style={previewPanelStyle}>
          <h2 style={sectionTitleStyle}>Preview</h2>
          <div style={previewGridStyle}>
            <div style={previewCardStyle}>
              <span style={miniLabelStyle}>공통 시스템 기본 계산</span>
              <LumiAvatar state={selectedState} size={size} reducedMotion={false} />
              <code style={codeStyle}>row {spriteCoordinate.row}, col {spriteCoordinate.col}</code>
            </div>
            <div style={previewCardStyle}>
              <span style={miniLabelStyle}>수동 스프라이트 튜닝</span>
              <div style={manualSpriteStyle} />
              <code style={codeStyle}>
                background-size: {backgroundSizeX}% {backgroundSizeY}%{'\n'}
                background-position: {backgroundPositionX}% {backgroundPositionY}%
              </code>
            </div>
          </div>
          <div style={dockStageStyle}>
            <div style={dockHeaderStyle}>
              <span style={miniLabelStyle}>Docked Bubble</span>
              <span style={codePillStyle}>{selectedDockSlot}</span>
            </div>
            <LumiDock slot={selectedDockSlot} visible={manualVisible}>
              <LumiAvatar state={selectedState} size={56} />
              <LumiBubble message={message} messageType={messageType} visible={manualVisible} />
            </LumiDock>
          </div>
          <div style={panelPreviewRowStyle}>
            <div style={panelColumnStyle}>
              <span style={miniLabelStyle}>Panel</span>
              <LumiPanel
                open={manualExpanded}
                state={selectedState}
                message={message}
                quickActions={lumi.viewState.quickActions}
                reducedMotion={lumi.prefersReducedMotion}
              />
            </div>
            <div style={sheetColumnStyle}>
              <span style={miniLabelStyle}>원본 스프라이트 시트</span>
              <div style={spriteSheetFrameStyle}><div style={spriteSheetStyle} /></div>
              <p style={sheetHelpStyle}>{`${LUMI_SPRITE_SHEET_WIDTH}×${LUMI_SPRITE_SHEET_HEIGHT}px · ${LUMI_SPRITE_COLUMNS}열×${LUMI_SPRITE_ROWS}행 · 셀 ${LUMI_SPRITE_CELL_WIDTH}×${LUMI_SPRITE_CELL_HEIGHT}px`}</p>
              <p style={sheetHelpStyle}>균일 셀 기준이므로 기본 grid 계산값(background-size {LUMI_SPRITE_COLUMNS * 100}% {LUMI_SPRITE_ROWS * 100}%)이 정확합니다. 캐릭터 위치 보정이 필요할 때만 슬라이더를 조정하세요.</p>
            </div>
          </div>
        </section>
      </div>
    </>
  );
}
