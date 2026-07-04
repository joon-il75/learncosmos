import { useEffect, useState, type CSSProperties } from 'react'
import ModalPortal from '@/components/common/ModalPortal'
import { LumiInterviewPanel } from '@/components/goal-interview/LumiInterviewPanel'
import { S } from '../styles'
import type { DebugScenario, GoalProfile, ScenarioInput } from '../types'

type Props = {
  selectedScenario: DebugScenario | null
  scenarioInput: ScenarioInput
  goalProfile: GoalProfile | null
  goalMessage: string
  setGoalMessage: React.Dispatch<React.SetStateAction<string>>
  runningGoalAction: boolean
  creatingScenario: boolean
  generatingLessons: boolean
  onGoalAction: (action: 'start' | 'message' | 'confirm', goalOverride?: string) => void
}

export default function GoalChatCard({
  selectedScenario,
  scenarioInput,
  goalProfile,
  goalMessage,
  setGoalMessage,
  runningGoalAction,
  creatingScenario,
  generatingLessons,
  onGoalAction,
}: Props) {
  const [goalModalOpen, setGoalModalOpen] = useState(false)
  const isGoalBusy = runningGoalAction || creatingScenario || generatingLessons

  useEffect(() => {
    if (goalProfile) {
      setGoalModalOpen(true)
    }
  }, [goalProfile?.messages.length])

  const startGoalChat = () => {
    setGoalModalOpen(true)
    onGoalAction('start')
  }

  return (
    <div style={S.card}>
      <div style={{ display:'flex', justifyContent:'space-between', gap:'12px', alignItems:'flex-start', flexWrap:'wrap', marginBottom:'14px' }}>
        <div>
          <div style={{ fontSize:'15px', fontWeight:800 }}>2. 목표 확정 채팅</div>
          <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.58)', marginTop:'5px', lineHeight:1.7 }}>
            실제 학습자 목표를 바꾸지 않고, 점검 시나리오의 goal profile snapshot에만 저장합니다.
          </div>
        </div>
        {selectedScenario && (
          <div style={{ fontSize:'11px', color:'#8CC7FF', fontWeight:800 }}>{selectedScenario.status}</div>
        )}
      </div>

      {selectedScenario ? (
        <div style={{ ...S.info, marginBottom:'14px' }}>
          <div style={{ fontWeight:700, marginBottom:'6px' }}>{selectedScenario.course_title}</div>
          <div>{selectedScenario.initial_user_intent || '목표 요청은 아래 입력 후 목표 채팅 시작 시 저장됩니다.'}</div>
        </div>
      ) : (
        <div style={{ ...S.info, marginBottom:'14px' }}>
          <div style={{ fontWeight:700, marginBottom:'6px' }}>{scenarioInput.course_title || '새 코스 제목 입력 대기'}</div>
          <div>목표 요청을 입력하고 채팅을 시작하면 새 점검 시나리오가 생성됩니다.</div>
        </div>
      )}

      {goalProfile ? (
        <div style={{ display:'grid', gap:'14px' }}>
          <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(180px, 1fr))', gap:'12px' }}>
            <div>
              <div style={S.label}>state</div>
              <div style={{ fontSize:'14px', fontWeight:800, color:'#8CC7FF' }}>{goalProfile.interview_state}</div>
            </div>
            <div>
              <div style={S.label}>goal</div>
              <div style={{ fontSize:'13px', color:'rgba(232,234,242,0.9)' }}>
                {goalProfile.interview_state === 'confirmed' && goalProfile.confirmed_goal ? goalProfile.confirmed_goal : '확정 전'}
              </div>
            </div>
            <div>
              <div style={S.label}>usage context</div>
              <div style={{ fontSize:'13px', color:'rgba(232,234,242,0.9)' }}>{goalProfile.usage_context || '아직 없음'}</div>
            </div>
            <div>
              <div style={S.label}>difficulty</div>
              <div style={{ fontSize:'13px', color:'rgba(232,234,242,0.9)' }}>{goalProfile.difficulty_level || '아직 없음'}</div>
            </div>
          </div>
          <div style={{ display:'flex', justifyContent:'flex-end' }}>
            <button style={S.primary} onClick={() => setGoalModalOpen(true)} disabled={isGoalBusy}>
              목표 채팅 열기
            </button>
          </div>
        </div>
      ) : (
        <div style={{ display:'grid', gap:'12px' }}>
          <div style={{ display:'flex', justifyContent:'flex-end' }}>
            <button style={S.primary} onClick={startGoalChat} disabled={runningGoalAction || creatingScenario}>
              {runningGoalAction || creatingScenario ? '목표 채팅 시작 중...' : '목표 채팅 시작'}
            </button>
          </div>
          <div style={S.warn}>코스 제목을 첫 메시지로 사용해 루미 목표 채팅을 시작합니다.</div>
        </div>
      )}

      <div style={{ ...S.info, marginTop:'14px' }}>
        목표 확정이 완료되면 리슨 생성까지 자동으로 진행하고, 생성된 리슨을 `generated_lessons_snapshot`에 저장합니다.
      </div>

      {goalModalOpen && (
        <ModalPortal overlayStyle={goalModalOverlayStyle} onMouseDown={() => !isGoalBusy && setGoalModalOpen(false)}>
          <div
            role="dialog"
            aria-modal="true"
            aria-label="추천 점검 목표 채팅"
            style={goalModalCardStyle}
            onMouseDown={event => event.stopPropagation()}
          >
            <div style={goalModalHeaderStyle}>
              <div style={{ minWidth:0 }}>
                <div style={goalModalEyebrowStyle}>Goal Interview</div>
                <h3 style={goalModalTitleStyle}>목표 확정 채팅</h3>
                <p style={goalModalDescriptionStyle}>
                  학습자 목표 대화와 같은 Lumi 패널로 점검 시나리오 목표를 확정합니다.
                </p>
              </div>
              <button
                type="button"
                style={{ ...goalModalCloseButtonStyle, ...(isGoalBusy ? goalModalCloseButtonDisabledStyle : null) }}
                onClick={() => setGoalModalOpen(false)}
                disabled={isGoalBusy}
              >
                x
              </button>
            </div>
            <div style={goalModalPanelStyle}>
              {goalProfile ? (
                <LumiInterviewPanel
                  profile={goalProfile}
                  isSending={isGoalBusy}
                  error={null}
                  onSendMessage={message => onGoalAction('message', message)}
                  onConfirmGoal={goal => onGoalAction('confirm', goal)}
                  busyMessage={generatingLessons ? '리슨을 생성하고 있습니다.' : null}
                />
              ) : (
                <div style={goalModalEmptyStyle}>
                  {isGoalBusy ? '목표 채팅을 준비하고 있습니다.' : '목표 채팅 시작을 눌러 대화를 시작해주세요.'}
                </div>
              )}
            </div>
          </div>
        </ModalPortal>
      )}
    </div>
  )
}

const goalModalOverlayStyle: CSSProperties = {
  padding:'24px',
  background:'rgba(4, 8, 18, 0.68)',
  backdropFilter:'blur(8px)',
}

const goalModalCardStyle: CSSProperties = {
  width:'min(880px, 100%)',
  maxHeight:'min(860px, calc(100vh - 48px))',
  display:'grid',
  gridTemplateRows:'auto minmax(0, 1fr)',
  gap:'14px',
  padding:'18px',
  borderRadius:'24px',
  border:'1px solid rgba(120,140,200,0.24)',
  background:'linear-gradient(180deg, rgba(18,25,42,0.98), rgba(10,15,26,0.98))',
  boxShadow:'0 28px 80px rgba(0,0,0,0.42)',
}

const goalModalHeaderStyle: CSSProperties = {
  display:'flex',
  alignItems:'flex-start',
  justifyContent:'space-between',
  gap:'16px',
}

const goalModalEyebrowStyle: CSSProperties = {
  fontSize:'11px',
  fontWeight:800,
  textTransform:'uppercase',
  color:'rgba(155,195,255,0.8)',
}

const goalModalTitleStyle: CSSProperties = {
  margin:'4px 0 0',
  fontSize:'22px',
  lineHeight:1.35,
  color:'#F4F7FF',
}

const goalModalDescriptionStyle: CSSProperties = {
  margin:'8px 0 0',
  fontSize:'13px',
  lineHeight:1.65,
  color:'rgba(218,226,246,0.72)',
}

const goalModalCloseButtonStyle: CSSProperties = {
  minWidth:'42px',
  minHeight:'42px',
  borderRadius:'999px',
  border:'1px solid rgba(120,140,200,0.24)',
  background:'rgba(255,255,255,0.06)',
  color:'#E7EEFF',
  fontSize:'18px',
  fontWeight:700,
  cursor:'pointer',
}

const goalModalCloseButtonDisabledStyle: CSSProperties = {
  opacity:0.5,
  cursor:'not-allowed',
}

const goalModalPanelStyle: CSSProperties = {
  minHeight:0,
  overflow:'hidden',
  borderRadius:'20px',
  background:'#FFFDF7',
}

const goalModalEmptyStyle: CSSProperties = {
  display:'grid',
  placeItems:'center',
  minHeight:'360px',
  padding:'24px',
  color:'#5B3912',
  fontSize:'14px',
  fontWeight:800,
  background:'#FFFDF7',
}
