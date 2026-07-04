import { displayRolloutMode } from '../constants'
import { S } from '../styles'
import type { RankerArtifactStatus, RolloutState } from '../types'

type Props = {
  rolloutState: RolloutState | null
  rankerInput: { ranker_model_version: string; feature_schema_version: string }
  setRankerInput: React.Dispatch<React.SetStateAction<{ ranker_model_version: string; feature_schema_version: string }>>
  savingRanker: boolean
  rankerArtifact: RankerArtifactStatus | null
  loadingScenarios: boolean
  onRefresh: () => void
  onSaveRanker: () => void
  onLoadRankerArtifact: () => void
}

export default function RankerConditionCard({
  rolloutState,
  rankerInput,
  setRankerInput,
  savingRanker,
  rankerArtifact,
  loadingScenarios,
  onRefresh,
  onSaveRanker,
  onLoadRankerArtifact,
}: Props) {
  return (
    <div style={S.card}>
      <div style={{ display:'flex', justifyContent:'space-between', gap:'12px', alignItems:'flex-start', flexWrap:'wrap' }}>
        <div>
          <div style={{ fontSize:'15px', fontWeight:700, marginBottom:'6px' }}>LightGBM Ranker 모델 사용 조건</div>
          <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.58)', lineHeight:1.7 }}>
            현재 추천 엔진 설정과 Ranker artifact 파일, 버전, feature schema 일치 여부를 확인합니다.
          </div>
        </div>
        <button style={S.secondary} onClick={onRefresh} disabled={loadingScenarios}>
          {loadingScenarios ? '새로고침 중...' : '새로고침'}
        </button>
      </div>
      {rolloutState ? (
        <div style={{ display:'grid', gap:'14px', marginTop:'16px' }}>
          <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(190px, 1fr))', gap:'12px' }}>
            <div>
              <div style={S.label}>mode</div>
              <div style={{ fontSize:'16px', fontWeight:800, color:'#8CC7FF' }}>{displayRolloutMode(rolloutState.mode)}</div>
            </div>
            <div>
              <div style={S.label}>embedding</div>
              <div style={{ fontSize:'13px', color:'rgba(232,234,242,0.9)' }}>{rolloutState.embedding_provider} / {rolloutState.embedding_model}</div>
              <div style={{ fontSize:'11px', color:'rgba(200,210,235,0.5)', marginTop:'3px' }}>dim {rolloutState.embedding_dimension}</div>
            </div>
            <div>
              <div style={S.label}>ranker</div>
              <div style={{ fontSize:'13px', color:'rgba(232,234,242,0.9)' }}>{rolloutState.ranker_model_version || '미적용'}</div>
              <div style={{ fontSize:'11px', color:'rgba(200,210,235,0.5)', marginTop:'3px' }}>{rolloutState.feature_schema_version}</div>
            </div>
            <div>
              <div style={S.label}>quality gate</div>
              <div style={{ fontSize:'16px', fontWeight:800, color: rolloutState.quality_gate_status === 'passed' ? '#8EE59A' : '#FAC775' }}>{rolloutState.quality_gate_status}</div>
            </div>
            <div>
              <div style={S.label}>updated</div>
              <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.7)' }}>{new Date(rolloutState.updated_at).toLocaleString('ko-KR')}</div>
            </div>
          </div>
          <div style={{ background:'rgba(255,255,255,0.03)', border:'1px solid rgba(120,140,200,0.12)', borderRadius:'12px', padding:'12px' }}>
            <div style={{ display:'flex', justifyContent:'space-between', gap:'10px', alignItems:'center', flexWrap:'wrap', marginBottom:'10px' }}>
              <div style={{ fontSize:'13px', fontWeight:800 }}>LightGBM Ranker artifact</div>
              <button style={{ ...S.secondary, padding:'8px 11px', fontSize:'11px' }} onClick={onLoadRankerArtifact}>
                artifact 확인
              </button>
            </div>
            <div style={{ display:'grid', gridTemplateColumns:'minmax(0, 1.2fr) minmax(0, 1fr) auto', gap:'10px', alignItems:'end' }}>
              <label>
                <span style={S.label}>model version</span>
                <input
                  style={S.input}
                  value={rankerInput.ranker_model_version}
                  onChange={e => setRankerInput(prev => ({ ...prev, ranker_model_version:e.target.value }))}
                  placeholder="예: lgbm-ranker-2026-05-v1"
                />
              </label>
              <label>
                <span style={S.label}>feature schema</span>
                <input
                  style={S.input}
                  value={rankerInput.feature_schema_version}
                  onChange={e => setRankerInput(prev => ({ ...prev, feature_schema_version:e.target.value }))}
                  placeholder="ranker-feature-v1"
                />
              </label>
              <button style={S.secondary} onClick={onSaveRanker} disabled={savingRanker}>
                {savingRanker ? '저장 중...' : '저장'}
              </button>
            </div>
            {rankerArtifact && (
              <div style={{ ...(rankerArtifact.ready ? S.info : S.warn), marginTop:'12px' }}>
                <div style={{ display:'flex', justifyContent:'space-between', gap:'10px', flexWrap:'wrap', marginBottom:'6px' }}>
                  <strong>{rankerArtifact.ready ? 'artifact ready' : 'artifact not ready'}</strong>
                  <span>{rankerArtifact.provider || 'lightgbm'} · {rankerArtifact.config_source || 'none'}</span>
                </div>
                <div>{rankerArtifact.reason}</div>
                <div style={{ marginTop:'6px', wordBreak:'break-all' }}>
                  model file: {rankerArtifact.model_path || '미설정'}
                </div>
                <div style={{ wordBreak:'break-all' }}>
                  manifest file: {rankerArtifact.manifest_path || '미설정'} · {rankerArtifact.manifest_exists ? 'found' : 'optional/missing'}
                </div>
                <div style={{ marginTop:'6px' }}>
                  version {rankerArtifact.ranker_model_version || '미설정'} · schema {rankerArtifact.feature_schema_version || 'ranker-feature-v1'} · size {rankerArtifact.model_size_bytes || 0} bytes
                </div>
                {rankerArtifact.checked_paths?.length > 0 && (
                  <details style={{ marginTop:'8px' }}>
                    <summary style={{ cursor:'pointer', fontWeight:800 }}>checked files</summary>
                    <div style={{ marginTop:'6px', display:'grid', gap:'3px' }}>
                      {rankerArtifact.checked_paths.map(path => (
                        <div key={path} style={{ wordBreak:'break-all', color:'rgba(200,210,235,0.62)' }}>{path}</div>
                      ))}
                    </div>
                  </details>
                )}
              </div>
            )}
          </div>
        </div>
      ) : (
        <div style={{ ...S.warn, marginTop:'16px' }}>활성 전환 상태가 없습니다.</div>
      )}
    </div>
  )
}
