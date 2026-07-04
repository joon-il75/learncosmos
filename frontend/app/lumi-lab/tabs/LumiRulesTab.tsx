'use client';

import LumiAvatar from '@/components/lumi/LumiAvatar';
import type { LumiAction, LumiMode, LumiPage, LumiScene } from '@/lib/lumi/lumiEngineTypes';
import type { LumiPageContext, LumiDockSlot, LumiMessageType, LumiState, LumiTrigger } from '@/lib/lumi/lumiTypes';
import {
  stateOptions, contextOptions, dockSlotOptions, messageTypeOptions,
  triggerOptions, runtimeModeOptions, runtimePageOptions, runtimeActionOptions, runtimeSceneOptions,
} from '../lumiLabTypes';
import type { LumiLabRuntimeConfigResult } from '../useLumiLabRuntimeConfig';
import { DetailRow } from '../LumiLabComponents';
import {
  rulesLayoutStyle, ruleListCardStyle, ruleEditorCardStyle, rulePreviewCardStyle,
  ruleCardHeaderStyle, ruleActionRowStyle, ruleListStyle, ruleListItemStyle, ruleListItemActiveStyle,
  ruleListTitleStyle, ruleListMetaStyle, ruleEditorGridStyle, rulePreviewHeroStyle,
  rulePreviewAvatarStyle, rulePreviewBubbleStyle, rulePreviewTypeStyle, rulePreviewMetaStyle,
  ruleNoteCardStyle, ruleNoteCopyStyle, diffListStyle, diffLineStyle, draftToolbarStyle,
  infoCardSubtitleStyle, messagePreviewStyle, sectionTitleStyle, fieldStyle, labelStyle,
  selectStyle, textareaStyle, codePillStyle, miniLabelStyle, secondaryActionButtonStyle, primaryButtonStyle,
} from '../lumiLabStyles';

interface Props {
  rc: LumiLabRuntimeConfigResult;
  onApplyToPreview: () => void;
}

export default function LumiRulesTab({ rc, onApplyToPreview }: Props) {
  const {
    filteredRuleDrafts, ruleDrafts, selectedRuleId, setSelectedRuleId,
    ruleSearchQuery, setRuleSearchQuery, ruleSortMode, setRuleSortMode,
    selectedRuleDraft, updateSelectedRuleDraft, resetSelectedRuleDraft,
    addRuleDraft, deleteSelectedRuleDraft,
    selectedRuleDiffLines, selectedRuleEngineBaseline,
  } = rc;

  return (
    <section style={rulesLayoutStyle}>
      <section style={ruleListCardStyle}>
        <div style={ruleCardHeaderStyle}>
          <h2 style={sectionTitleStyle}>Rule Drafts</h2>
          <span style={codePillStyle}>{filteredRuleDrafts.length} / {ruleDrafts.length}</span>
        </div>
        <p style={infoCardSubtitleStyle}>현재 runtime 엔진에서 자주 쓰는 규칙을 편집합니다. 변경 사항은 상단 저장 버튼으로 실제 runtime source에 반영할 수 있고, 저장 전에는 Lab editor 상태로만 유지됩니다.</p>
        <div style={draftToolbarStyle}>
          <input value={ruleSearchQuery} onChange={(e) => setRuleSearchQuery(e.target.value)} placeholder="rule 검색" style={{ padding: '6px 10px', borderRadius: '6px', border: '1px solid rgba(255,255,255,0.12)', background: 'rgba(255,255,255,0.06)', color: 'inherit', fontSize: '13px', width: '100%' }} />
          <select value={ruleSortMode} onChange={(e) => setRuleSortMode(e.target.value as typeof ruleSortMode)} style={selectStyle}>
            <option value="runtime">runtime 순</option>
            <option value="label">label 순</option>
            <option value="state">state 순</option>
          </select>
        </div>
        <div style={ruleListStyle}>
          {filteredRuleDrafts.map((rule) => (
            <button key={rule.id} type="button" onClick={() => setSelectedRuleId(rule.id)} style={{ ...ruleListItemStyle, ...(selectedRuleId === rule.id ? ruleListItemActiveStyle : null) }}>
              <strong style={ruleListTitleStyle}>{rule.label}</strong>
              <span style={ruleListMetaStyle}>{rule.page} / {rule.action} / {rule.scene}</span>
              <span style={ruleListMetaStyle}>{rule.state} / {rule.messageType}</span>
            </button>
          ))}
        </div>
      </section>

      <section style={ruleEditorCardStyle}>
        <div style={ruleCardHeaderStyle}>
          <h2 style={sectionTitleStyle}>Rule Editor</h2>
          <div style={ruleActionRowStyle}>
            <button type="button" onClick={addRuleDraft} style={secondaryActionButtonStyle}>Rule 추가</button>
            <button type="button" onClick={deleteSelectedRuleDraft} disabled={ruleDrafts.length <= 1} style={secondaryActionButtonStyle}>선택 삭제</button>
            <button type="button" onClick={resetSelectedRuleDraft} style={secondaryActionButtonStyle}>기본값 복원</button>
            <button type="button" onClick={onApplyToPreview} style={primaryButtonStyle}>Sprites preview 적용</button>
          </div>
        </div>
        {selectedRuleDraft ? (
          <div style={ruleEditorGridStyle}>
            <label style={fieldStyle}><span style={labelStyle}>Rule Label</span><input value={selectedRuleDraft.label} onChange={(e) => updateSelectedRuleDraft({ label: e.target.value })} style={{ padding: '6px 10px', borderRadius: '6px', border: '1px solid rgba(255,255,255,0.12)', background: 'rgba(255,255,255,0.06)', color: 'inherit', fontSize: '13px', width: '100%' }} /></label>
            <label style={fieldStyle}><span style={labelStyle}>Trigger</span>
              <select value={selectedRuleDraft.trigger} onChange={(e) => updateSelectedRuleDraft({ trigger: e.target.value as LumiTrigger })} style={selectStyle}>
                {triggerOptions.map((o) => <option key={o} value={o}>{o}</option>)}
              </select>
            </label>
            <label style={fieldStyle}><span style={labelStyle}>Mode</span>
              <select value={selectedRuleDraft.mode} onChange={(e) => updateSelectedRuleDraft({ mode: e.target.value as LumiMode })} style={selectStyle}>
                {runtimeModeOptions.map((o) => <option key={o} value={o}>{o}</option>)}
              </select>
            </label>
            <label style={fieldStyle}><span style={labelStyle}>Page</span>
              <select value={selectedRuleDraft.page} onChange={(e) => updateSelectedRuleDraft({ page: e.target.value as LumiPage })} style={selectStyle}>
                {runtimePageOptions.map((o) => <option key={o} value={o}>{o}</option>)}
              </select>
            </label>
            <label style={fieldStyle}><span style={labelStyle}>Action</span>
              <select value={selectedRuleDraft.action} onChange={(e) => updateSelectedRuleDraft({ action: e.target.value as LumiAction })} style={selectStyle}>
                {runtimeActionOptions.map((o) => <option key={o} value={o}>{o}</option>)}
              </select>
            </label>
            <label style={fieldStyle}><span style={labelStyle}>Scene</span>
              <select value={selectedRuleDraft.scene} onChange={(e) => updateSelectedRuleDraft({ scene: e.target.value as LumiScene })} style={selectStyle}>
                {runtimeSceneOptions.map((o) => <option key={o} value={o}>{o}</option>)}
              </select>
            </label>
            <label style={fieldStyle}><span style={labelStyle}>Legacy Context</span>
              <select value={selectedRuleDraft.context} onChange={(e) => updateSelectedRuleDraft({ context: e.target.value as LumiPageContext })} style={selectStyle}>
                {contextOptions.map((o) => <option key={o} value={o}>{o}</option>)}
              </select>
            </label>
            <label style={fieldStyle}><span style={labelStyle}>Dock Slot</span>
              <select value={selectedRuleDraft.dockSlot} onChange={(e) => updateSelectedRuleDraft({ dockSlot: e.target.value as LumiDockSlot })} style={selectStyle}>
                {dockSlotOptions.map((o) => <option key={o} value={o}>{o}</option>)}
              </select>
            </label>
            <label style={fieldStyle}><span style={labelStyle}>State</span>
              <select value={selectedRuleDraft.state} onChange={(e) => updateSelectedRuleDraft({ state: e.target.value as LumiState })} style={selectStyle}>
                {stateOptions.map((o) => <option key={o} value={o}>{o}</option>)}
              </select>
            </label>
            <label style={fieldStyle}><span style={labelStyle}>Message Type</span>
              <select value={selectedRuleDraft.messageType} onChange={(e) => updateSelectedRuleDraft({ messageType: e.target.value as LumiMessageType })} style={selectStyle}>
                {messageTypeOptions.map((o) => <option key={o} value={o}>{o}</option>)}
              </select>
            </label>
            <label style={{ ...fieldStyle, gridColumn: '1 / -1' }}><span style={labelStyle}>Message</span><textarea value={selectedRuleDraft.message} onChange={(e) => updateSelectedRuleDraft({ message: e.target.value })} rows={4} style={textareaStyle} /></label>
            <label style={{ ...fieldStyle, gridColumn: '1 / -1' }}><span style={labelStyle}>Rule Note</span><textarea value={selectedRuleDraft.note} onChange={(e) => updateSelectedRuleDraft({ note: e.target.value })} rows={3} style={textareaStyle} /></label>
          </div>
        ) : null}
      </section>

      <section style={rulePreviewCardStyle}>
        <div style={ruleCardHeaderStyle}>
          <h2 style={sectionTitleStyle}>Rule Preview</h2>
          <span style={codePillStyle}>{selectedRuleDraft?.trigger ?? 'none'}</span>
        </div>
        {selectedRuleDraft ? (
          <>
            <div style={rulePreviewHeroStyle}>
              <div style={rulePreviewAvatarStyle}><LumiAvatar state={selectedRuleDraft.state} size={68} reducedMotion={false} /></div>
              <div style={rulePreviewBubbleStyle}>
                <div style={rulePreviewTypeStyle}>{selectedRuleDraft.messageType}</div>
                <div style={messagePreviewStyle}>{selectedRuleDraft.message}</div>
              </div>
            </div>
            <div style={rulePreviewMetaStyle}>
              <DetailRow label="runtime" value={`${selectedRuleDraft.page} / ${selectedRuleDraft.action} / ${selectedRuleDraft.scene}`} />
              <DetailRow label="legacy" value={`${selectedRuleDraft.context} / ${selectedRuleDraft.dockSlot}`} />
              <DetailRow label="baseline" value={selectedRuleEngineBaseline ? `${selectedRuleEngineBaseline.state} / ${selectedRuleEngineBaseline.messageType}` : '미연결'} />
            </div>
            <div style={ruleNoteCardStyle}>
              <strong style={miniLabelStyle}>저장 전 Diff</strong>
              <div style={diffListStyle}>
                {(selectedRuleDiffLines.length > 0 ? selectedRuleDiffLines : ['저장본과 현재 편집본이 같습니다.']).map((line) => (
                  <code key={line} style={diffLineStyle}>{line}</code>
                ))}
              </div>
            </div>
            <div style={ruleNoteCardStyle}>
              <strong style={miniLabelStyle}>운영 메모</strong>
              <p style={ruleNoteCopyStyle}>{selectedRuleDraft.note}</p>
            </div>
          </>
        ) : null}
      </section>
    </section>
  );
}
