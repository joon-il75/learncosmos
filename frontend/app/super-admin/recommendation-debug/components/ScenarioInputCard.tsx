import { S } from '../styles'
import type { ScenarioInput } from '../types'

type Props = {
  scenarioInput: ScenarioInput
  onChangeScenarioInput: (patch: Partial<ScenarioInput>) => void
}

export default function ScenarioInputCard({ scenarioInput, onChangeScenarioInput }: Props) {
  return (
    <div style={S.card}>
      <div style={{ fontSize:'15px', fontWeight:700, marginBottom:'14px' }}>1. 코스 생성 입력</div>
      <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(260px, 1fr))', gap:'14px' }}>
        <div>
          <label style={S.label}>코스 제목</label>
          <input
            style={S.input}
            value={scenarioInput.course_title}
            onChange={e => onChangeScenarioInput({ course_title: e.target.value })}
            placeholder="예: 4주 안에 어반스케치 시작하기"
          />
        </div>
        <div>
          <label style={S.label}>운영 메모</label>
          <input
            style={S.input}
            value={scenarioInput.notes}
            onChange={e => onChangeScenarioInput({ notes: e.target.value })}
            placeholder="예: 초보자 목표 생성 품질 확인"
          />
        </div>
      </div>
      <div style={{ ...S.info, marginTop:'14px' }}>
        코스 제목을 입력한 뒤 아래 목표 채팅을 시작하면 점검 시나리오가 자동으로 저장됩니다.
      </div>
    </div>
  )
}
