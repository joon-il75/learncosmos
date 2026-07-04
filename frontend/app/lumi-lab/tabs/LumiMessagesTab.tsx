'use client';

import type { LumiMessageType, LumiTrigger } from '@/lib/lumi/lumiTypes';
import { messageTypeOptions, triggerOptions } from '../lumiLabTypes';
import type { LumiLabRuntimeConfigResult } from '../useLumiLabRuntimeConfig';
import type { LumiLabUIPreviewResult } from '../useLumiLabUIPreview';
import { DetailRow } from '../LumiLabComponents';
import {
  rulesLayoutStyle, ruleListCardStyle, ruleEditorCardStyle, rulePreviewCardStyle,
  ruleCardHeaderStyle, ruleActionRowStyle, ruleListStyle, ruleListItemStyle, ruleListItemActiveStyle,
  ruleListTitleStyle, ruleListMetaStyle, ruleEditorGridStyle, rulePreviewBubbleStyle,
  rulePreviewTypeStyle, rulePreviewMetaStyle, ruleNoteCardStyle, ruleNoteCopyStyle,
  diffListStyle, diffLineStyle, draftToolbarStyle, infoCardSubtitleStyle, messagePreviewStyle,
  promptFootnoteStyle, sectionTitleStyle, fieldStyle, labelStyle, selectStyle, textareaStyle,
  codePillStyle, miniLabelStyle, secondaryActionButtonStyle, primaryButtonStyle,
} from '../lumiLabStyles';

interface Props {
  rc: LumiLabRuntimeConfigResult;
  uip: LumiLabUIPreviewResult;
  onApplyToPreview: () => void;
}

export default function LumiMessagesTab({ rc, uip, onApplyToPreview }: Props) {
  const {
    filteredMessageDrafts, messageDrafts, selectedMessageDraftId, setSelectedMessageDraftId,
    messageSearchQuery, setMessageSearchQuery, messageSortMode, setMessageSortMode,
    selectedMessageDraft, updateSelectedMessageDraft, resetSelectedMessageDraft,
    addMessageDraft, deleteSelectedMessageDraft, selectedMessageDiffLines,
  } = rc;
  const { messageType, lumi } = uip;

  return (
    <section style={rulesLayoutStyle}>
      <section style={ruleListCardStyle}>
        <div style={ruleCardHeaderStyle}>
          <h2 style={sectionTitleStyle}>Message Drafts</h2>
          <span style={codePillStyle}>{filteredMessageDrafts.length} / {messageDrafts.length}</span>
        </div>
        <p style={infoCardSubtitleStyle}>분리된 message storage를 기준으로 message draft를 편집합니다. 변경 사항은 상단 저장 버튼으로 실제 runtime source에 반영할 수 있고, 저장 전에는 Lab editor 상태로만 유지됩니다.</p>
        <div style={draftToolbarStyle}>
          <input value={messageSearchQuery} onChange={(e) => setMessageSearchQuery(e.target.value)} placeholder="message 검색" style={{ padding: '6px 10px', borderRadius: '6px', border: '1px solid rgba(255,255,255,0.12)', background: 'rgba(255,255,255,0.06)', color: 'inherit', fontSize: '13px', width: '100%' }} />
          <select value={messageSortMode} onChange={(e) => setMessageSortMode(e.target.value as typeof messageSortMode)} style={selectStyle}>
            <option value="trigger">trigger 순</option>
            <option value="label">label 순</option>
            <option value="type">type 순</option>
          </select>
        </div>
        <div style={ruleListStyle}>
          {filteredMessageDrafts.map((draft) => (
            <button key={draft.id} type="button" onClick={() => setSelectedMessageDraftId(draft.id)} style={{ ...ruleListItemStyle, ...(selectedMessageDraftId === draft.id ? ruleListItemActiveStyle : null) }}>
              <strong style={ruleListTitleStyle}>{draft.label}</strong>
              <span style={ruleListMetaStyle}>{draft.trigger}</span>
              <span style={ruleListMetaStyle}>{draft.messageType}</span>
            </button>
          ))}
        </div>
      </section>

      <section style={ruleEditorCardStyle}>
        <div style={ruleCardHeaderStyle}>
          <h2 style={sectionTitleStyle}>Message Editor</h2>
          <div style={ruleActionRowStyle}>
            <button type="button" onClick={addMessageDraft} style={secondaryActionButtonStyle}>Message 추가</button>
            <button type="button" onClick={deleteSelectedMessageDraft} disabled={messageDrafts.length <= 1} style={secondaryActionButtonStyle}>선택 삭제</button>
            <button type="button" onClick={resetSelectedMessageDraft} style={secondaryActionButtonStyle}>기본값 복원</button>
            <button type="button" onClick={onApplyToPreview} style={primaryButtonStyle}>Preview 메시지 적용</button>
          </div>
        </div>
        {selectedMessageDraft ? (
          <div style={ruleEditorGridStyle}>
            <label style={fieldStyle}><span style={labelStyle}>Draft Label</span><input value={selectedMessageDraft.label} onChange={(e) => updateSelectedMessageDraft({ label: e.target.value })} style={{ padding: '6px 10px', borderRadius: '6px', border: '1px solid rgba(255,255,255,0.12)', background: 'rgba(255,255,255,0.06)', color: 'inherit', fontSize: '13px', width: '100%' }} /></label>
            <label style={fieldStyle}><span style={labelStyle}>Trigger</span>
              <select value={selectedMessageDraft.trigger} onChange={(e) => updateSelectedMessageDraft({ trigger: e.target.value as LumiTrigger })} style={selectStyle}>
                {triggerOptions.filter((o) => o !== 'planet_hover' && o !== 'planet_select' && o !== 'lesson_enter').map((o) => <option key={o} value={o}>{o}</option>)}
              </select>
            </label>
            <label style={fieldStyle}><span style={labelStyle}>Message Type</span>
              <select value={selectedMessageDraft.messageType} onChange={(e) => updateSelectedMessageDraft({ messageType: e.target.value as LumiMessageType })} style={selectStyle}>
                {messageTypeOptions.map((o) => <option key={o} value={o}>{o}</option>)}
              </select>
            </label>
            <label style={fieldStyle}><span style={labelStyle}>Preview Hint</span><input value={selectedMessageDraft.preview} onChange={(e) => updateSelectedMessageDraft({ preview: e.target.value })} style={{ padding: '6px 10px', borderRadius: '6px', border: '1px solid rgba(255,255,255,0.12)', background: 'rgba(255,255,255,0.06)', color: 'inherit', fontSize: '13px', width: '100%' }} placeholder="preview text" /></label>
            <label style={{ ...fieldStyle, gridColumn: '1 / -1' }}><span style={labelStyle}>Message</span><textarea value={selectedMessageDraft.message} onChange={(e) => updateSelectedMessageDraft({ message: e.target.value })} rows={5} style={textareaStyle} /></label>
            <label style={{ ...fieldStyle, gridColumn: '1 / -1' }}><span style={labelStyle}>Note</span><textarea value={selectedMessageDraft.note} onChange={(e) => updateSelectedMessageDraft({ note: e.target.value })} rows={3} style={textareaStyle} /></label>
          </div>
        ) : null}
      </section>

      <section style={rulePreviewCardStyle}>
        <div style={ruleCardHeaderStyle}>
          <h2 style={sectionTitleStyle}>Message Preview</h2>
          <span style={codePillStyle}>{selectedMessageDraft?.messageType ?? 'none'}</span>
        </div>
        {selectedMessageDraft ? (
          <>
            <div style={rulePreviewBubbleStyle}>
              <div style={rulePreviewTypeStyle}>{selectedMessageDraft.trigger}</div>
              <div style={messagePreviewStyle}>{selectedMessageDraft.message}</div>
              {selectedMessageDraft.preview ? <div style={promptFootnoteStyle}>preview: {selectedMessageDraft.preview}</div> : null}
            </div>
            <div style={rulePreviewMetaStyle}>
              <DetailRow label="draft type" value={selectedMessageDraft.messageType} />
              <DetailRow label="current type" value={messageType} />
              <DetailRow label="runtime type" value={lumi.viewState.messageType} />
            </div>
            <div style={ruleNoteCardStyle}>
              <strong style={miniLabelStyle}>저장 전 Diff</strong>
              <div style={diffListStyle}>
                {(selectedMessageDiffLines.length > 0 ? selectedMessageDiffLines : ['저장본과 현재 편집본이 같습니다.']).map((line) => (
                  <code key={line} style={diffLineStyle}>{line}</code>
                ))}
              </div>
            </div>
            <div style={ruleNoteCardStyle}>
              <strong style={miniLabelStyle}>운영 메모</strong>
              <p style={ruleNoteCopyStyle}>{selectedMessageDraft.note || '메모 없음'}</p>
            </div>
          </>
        ) : null}
      </section>
    </section>
  );
}
