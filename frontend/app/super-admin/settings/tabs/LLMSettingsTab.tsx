'use client'
import type { Dispatch, SetStateAction } from 'react'
import type { LLMSetting } from '../settingsTypes'
import { PROVIDERS, PROVIDER_LABELS, MODELS, DEFAULT_FEATURE, PRO_FEATURES } from '../settingsConstants'
import { S, OPTION_STYLE } from '../settingsStyles'
import { FeatureTag } from '../settingsComponents'

interface Props {
  llmSettings: LLMSetting[]
  apiKeys: Record<string, string>
  setApiKeys: Dispatch<SetStateAction<Record<string, string>>>
  apiKey2s: Record<string, string>
  setApiKey2s: Dispatch<SetStateAction<Record<string, string>>>
  endpoints: Record<string, string>
  setEndpoints: Dispatch<SetStateAction<Record<string, string>>>
  apiKeyStatuses: Record<string, string>
  apiKey2Statuses: Record<string, string>
  apiKeyUpdatedAts: Record<string, string>
  saving: boolean
  saveLLM: () => Promise<void>
  saveAPIKey: (provider: string) => Promise<void>
  keyStatusText: (provider: string) => string
  updateLLM: (feature: string, field: string, value: string) => void
}

export function LLMSettingsTab({
  llmSettings, apiKeys, setApiKeys, apiKey2s, setApiKey2s,
  endpoints, setEndpoints, apiKeyStatuses, apiKey2Statuses, apiKeyUpdatedAts,
  saving, saveLLM, saveAPIKey, keyStatusText, updateLLM,
}: Props) {
  return (
    <div style={{ paddingTop:'20px' }}>
      <div style={S.infoBox}>
        🔑 API 키는 저장 후 다시 표시되지 않습니다.<br/>
        🔐 저장 시 서버에서 암호화되어 DB에 보관되고 런타임에 즉시 적용됩니다.<br/>
        ⚡ 일반 사용자: 기본값 1개를 모든 기능에 공통 적용 · LearnCosmos Pro: 기능별 전용 설정 사용
      </div>

      <div style={S.card}>
        <div style={S.cardTitle}>기본 LLM 설정</div>
        <div style={S.cardSub}>일반 사용자에게는 아래 기본값 1개가 커리큘럼·튜터·비전·퀴즈 전체에 공통 적용됩니다.</div>
        <div style={{ display:'grid', gridTemplateColumns:'110px 1fr 1fr 80px', gap:'10px', paddingBottom:'9px', borderBottom:'1px solid rgba(120,140,200,0.12)', fontSize:'11px', fontWeight:700, color:'rgba(200,210,235,0.35)', letterSpacing:'.5px' }}>
          <span>기능</span><span>제공자</span><span>모델</span><span>상태</span>
        </div>
        {[DEFAULT_FEATURE].map(feature => {
          const s = llmSettings.find(x => x.feature === feature)
          const isSet = !!s?.provider
          return (
            <div key={feature} style={{ display:'grid', gridTemplateColumns:'110px 1fr 1fr 80px', gap:'10px', alignItems:'center', padding:'10px 0', borderBottom:'1px solid rgba(120,140,200,0.07)' }}>
              <FeatureTag feature={feature} />
              <select style={S.sel} value={s?.provider || 'openai'} onChange={e => updateLLM(feature, 'provider', e.target.value)}>
                {PROVIDERS.map(p => <option style={OPTION_STYLE} key={p} value={p}>{PROVIDER_LABELS[p]}</option>)}
              </select>
              <select style={S.sel} value={s?.model || ''} onChange={e => updateLLM(feature, 'model', e.target.value)}>
                {(MODELS[s?.provider || 'openai'] || []).map(m => <option style={OPTION_STYLE} key={m} value={m}>{m}</option>)}
              </select>
              <span style={{ fontSize:'11px', color: isSet ? '#52B788' : 'rgba(200,210,235,0.3)' }}>
                {isSet ? '● 적용 중' : '○ 미설정'}
              </span>
            </div>
          )
        })}
        <div style={{ ...S.cardSub, marginTop:'16px', marginBottom:'10px' }}>LearnCosmos Pro 가입자에게만 아래 기능별 전용 LLM 설정이 적용됩니다.</div>
        <div style={{ display:'grid', gridTemplateColumns:'110px 1fr 1fr 80px', gap:'10px', paddingBottom:'9px', borderBottom:'1px solid rgba(120,140,200,0.12)', fontSize:'11px', fontWeight:700, color:'rgba(200,210,235,0.35)', letterSpacing:'.5px' }}>
          <span>Pro 기능</span><span>제공자</span><span>모델</span><span>상태</span>
        </div>
        {PRO_FEATURES.map(feature => {
          const s = llmSettings.find(x => x.feature === feature)
          const isSet = !!s?.provider
          return (
            <div key={feature} style={{ display:'grid', gridTemplateColumns:'110px 1fr 1fr 80px', gap:'10px', alignItems:'center', padding:'10px 0', borderBottom:'1px solid rgba(120,140,200,0.07)' }}>
              <FeatureTag feature={feature} />
              <select style={S.sel} value={s?.provider || 'openai'} onChange={e => updateLLM(feature, 'provider', e.target.value)}>
                {PROVIDERS.map(p => <option style={OPTION_STYLE} key={p} value={p}>{PROVIDER_LABELS[p]}</option>)}
              </select>
              <select style={S.sel} value={s?.model || ''} onChange={e => updateLLM(feature, 'model', e.target.value)}>
                {(MODELS[s?.provider || 'openai'] || []).map(m => <option style={OPTION_STYLE} key={m} value={m}>{m}</option>)}
              </select>
              <span style={{ fontSize:'11px', color: isSet ? '#52B788' : 'rgba(200,210,235,0.3)' }}>
                {isSet ? '● 적용 중' : '○ 미설정'}
              </span>
            </div>
          )
        })}
        <div style={S.saveRow}><button style={S.btnSave} onClick={saveLLM} disabled={saving}>기능별 설정 저장</button></div>
      </div>

      <div style={S.card}>
        <div style={S.cardTitle}>제공자별 API 키</div>
        <div style={S.cardSub}>저장 후 키는 다시 표시되지 않습니다. AES 암호화 저장되며 변경 시에만 입력하세요.</div>
        <div style={S.grid2}>
          {['openai','anthropic','google','grok','solar'].map(p => (
            <div key={p} style={{ marginBottom:'12px' }}>
              <label style={S.label}>
                {PROVIDER_LABELS[p]} API Key
                {apiKeyStatuses[p] === 'SET' && <span style={{ color:'#52B788', fontSize:'10px', marginLeft:'8px' }}>● 설정됨</span>}
              </label>
              <div style={{ fontSize:'11px', color:'rgba(200,210,235,0.42)', marginBottom:'6px' }}>
                {keyStatusText(p)}{apiKeyUpdatedAts[p] ? ` · 최근 저장 ${apiKeyUpdatedAts[p]}` : ''}
              </div>
              <div style={{ display:'flex', gap:'8px' }}>
                <input style={{ ...S.inp, letterSpacing: apiKeys[p] ? '0' : '2px' }} type="password"
                  placeholder={apiKeyStatuses[p] === 'SET' ? '변경 시에만 입력' : '키 입력...'}
                  value={apiKeys[p] || ''}
                  onChange={e => setApiKeys(prev => ({ ...prev, [p]: e.target.value }))} />
                <button style={{ ...S.btnSave, padding:'9px 14px', fontSize:'12px', whiteSpace:'nowrap' }}
                  onClick={() => saveAPIKey(p)}>저장</button>
              </div>
            </div>
          ))}
        </div>

        <div style={{ marginBottom:'12px' }}>
          <label style={S.label}>
            HyperCLOVA X API Key
            {apiKeyStatuses['hyperclova'] === 'SET' && <span style={{ color:'#52B788', fontSize:'10px', marginLeft:'8px' }}>● 설정됨</span>}
          </label>
          <div style={{ fontSize:'11px', color:'rgba(200,210,235,0.42)', marginBottom:'6px' }}>
            {keyStatusText('hyperclova')}{apiKeyUpdatedAts['hyperclova'] ? ` · 최근 저장 ${apiKeyUpdatedAts['hyperclova']}` : ''}
          </div>
          <div style={S.grid2}>
            <input style={S.inp} type="password" placeholder="NCP API Key"
              value={apiKeys['hyperclova'] || ''} onChange={e => setApiKeys(prev => ({ ...prev, hyperclova: e.target.value }))} />
            <input style={S.inp} type="password" placeholder="NCP APIGW Key"
              value={apiKey2s['hyperclova'] || ''} onChange={e => setApiKey2s(prev => ({ ...prev, hyperclova: e.target.value }))} />
          </div>
          <div style={{ ...S.saveRow, marginTop:'8px' }}>
            <button style={{ ...S.btnSave, padding:'9px 14px', fontSize:'12px' }} onClick={() => saveAPIKey('hyperclova')}>저장</button>
          </div>
        </div>

        <div style={{ background:'rgba(255,255,255,0.02)', border:'1px solid rgba(120,140,200,0.1)', borderRadius:'10px', padding:'13px 15px', marginTop:'4px' }}>
          <div style={{ fontSize:'12px', fontWeight:700, color:'rgba(200,210,235,0.5)', marginBottom:'11px' }}>🖥️ Ollama 자체 호스팅 (Llama · EXAONE)</div>
          <div style={S.grid2}>
            {['llama','exaone'].map(p => (
              <div key={p}>
                <label style={S.label}>{PROVIDER_LABELS[p]} 엔드포인트</label>
                <div style={{ fontSize:'11px', color:'rgba(200,210,235,0.42)', marginBottom:'6px' }}>
                  {keyStatusText(p)}{apiKeyUpdatedAts[p] ? ` · 최근 저장 ${apiKeyUpdatedAts[p]}` : ''}
                </div>
                <input style={S.inp} type="text" placeholder="http://localhost:11434"
                  value={endpoints[p] || ''}
                  onChange={e => setEndpoints(prev => ({ ...prev, [p]: e.target.value }))} />
              </div>
            ))}
          </div>
          <div style={{ ...S.saveRow, marginTop:'10px' }}>
            <button style={{ ...S.btnSave, padding:'9px 14px', fontSize:'12px' }}
              onClick={() => { saveAPIKey('llama'); saveAPIKey('exaone') }}>엔드포인트 저장</button>
          </div>
        </div>
      </div>
    </div>
  )
}
