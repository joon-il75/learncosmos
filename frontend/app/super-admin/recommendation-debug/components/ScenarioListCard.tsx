import { S } from '../styles'
import type { DebugScenario } from '../types'

type Props = {
  scenarios: DebugScenario[]
  selectedScenario: DebugScenario | null
  loadingScenarioDetailId: string | null
  onSelectScenario: (scenario: DebugScenario) => void
}

export default function ScenarioListCard({ scenarios, selectedScenario, loadingScenarioDetailId, onSelectScenario }: Props) {
  return (
    <div style={S.card}>
      <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', gap:'12px', flexWrap:'wrap', marginBottom:'14px' }}>
        <div style={{ fontSize:'15px', fontWeight:700 }}>테스트 시나리오 선택</div>
        <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.55)' }}>{scenarios.length}개</div>
      </div>
      <div style={{ display:'grid', gap:'10px' }}>
        {scenarios.length === 0 ? (
          <div style={S.warn}>아직 저장된 시나리오가 없습니다.</div>
        ) : scenarios.map(item => (
          <button
            key={item.id}
            type="button"
            onClick={() => onSelectScenario(item)}
            disabled={loadingScenarioDetailId === item.id}
            style={{
              width:'100%',
              textAlign:'left',
              background:selectedScenario?.id === item.id ? 'rgba(55,138,221,0.14)' : 'rgba(255,255,255,0.03)',
              border:selectedScenario?.id === item.id ? '1px solid rgba(55,138,221,0.45)' : '1px solid rgba(120,140,200,0.12)',
              borderRadius:'12px',
              padding:'13px',
              color:'#E8EAF2',
              cursor:'pointer',
            }}
          >
            <div style={{ display:'flex', justifyContent:'space-between', gap:'12px', flexWrap:'wrap' }}>
              <div style={{ fontSize:'14px', fontWeight:800 }}>{item.course_title}</div>
              <div style={{ fontSize:'11px', color:'#8CC7FF' }}>{loadingScenarioDetailId === item.id ? '불러오는 중...' : item.status}</div>
            </div>
            <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.62)', marginTop:'6px', lineHeight:1.6 }}>
              {item.initial_user_intent || '목표 요청은 목표 확정 채팅에서 입력 대기'}
            </div>
            <div style={{ fontSize:'11px', color:'rgba(200,210,235,0.45)', marginTop:'8px' }}>
              runs {item.run_count ?? 0} · labels {item.label_count ?? 0} · updated {new Date(item.updated_at).toLocaleString('ko-KR')}
            </div>
          </button>
        ))}
      </div>
    </div>
  )
}
