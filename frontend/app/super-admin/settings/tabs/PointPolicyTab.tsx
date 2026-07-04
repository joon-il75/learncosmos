'use client'
import type { Dispatch, SetStateAction } from 'react'
import type { PointPolicy } from '../settingsTypes'
import { S } from '../settingsStyles'

interface Props {
  policy: PointPolicy | null
  setPolicy: Dispatch<SetStateAction<PointPolicy | null>>
  saving: boolean
  savePolicy: () => Promise<void>
}

export function PointPolicyTab({ policy, setPolicy, saving, savePolicy }: Props) {
  return (
    <div style={{ paddingTop:'20px' }}>
      <div style={S.warnBox}>⚠️ 변경 즉시 적용됩니다. 웰컴 포인트는 기존 유저에게 소급 적용되지 않습니다. Pro 월 포인트와 기능별 차감도 여기서 운영 기준값을 직접 조정합니다.</div>
      {!policy ? (
        <div style={S.card}>
          <div style={S.cardTitle}>AI 포인트 정책값</div>
          <div style={S.cardSub}>DB에서 현재 운영값을 불러오는 중입니다.</div>
        </div>
      ) : (
        <div style={S.card}>
          <div style={S.cardTitle}>AI 포인트 정책값</div>
          <div style={S.cardSub}>point_settings 테이블에서 관리 · 일반 포인트 정책과 LearnCosmos Pro 월간 정책을 함께 관리</div>
          <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:'12px' }}>
            {[
              { key:'welcome_points', label:'🎁 웰컴 포인트', unit:'pt', desc:'신규 가입 시 1회 지급 · 소급 없음' },
              { key:'course_gen_cost', label:'🤖 AI 코스 생성 단가', unit:'pt', desc:'코스 1개 생성 시 차감 · OpenAI BYOK만 무료' },
              { key:'lesson_rec_cost', label:'🔍 레슨 AI 추천 단가', unit:'pt', desc:'AI 호출 비용 + 서버 사용 비용 기준 · BYOK 사용량 기록 없음' },
              { key:'admin_max_grant', label:'🛡️ 관리자 최대 지급량', unit:'pt', desc:'관리자가 학습자에게 지급 가능한 1회 한도' },
              { key:'pro_monthly_points', label:'💎 Pro 월 제공 포인트', unit:'pt', desc:'LearnCosmos Pro 월간/연간 가입자에게 매월 제공하는 기준 포인트' },
              { key:'tutor_cost', label:'💬 Pro 튜터 단가', unit:'pt', desc:'실시간 튜터링 1회 기준 차감 포인트' },
              { key:'vision_cost', label:'🖼️ Pro 비전 코칭 단가', unit:'pt', desc:'비전 코칭 1회 기준 차감 포인트 · 긴 변 1280px 제한' },
              { key:'quiz_cost', label:'📝 Pro 퀴즈 단가', unit:'pt', desc:'4지선다 + 이해도/숙련도 점검 질문 생성 1회 기준 차감 포인트' },
            ].map(({ key, label, unit, desc }) => (
              <div key={key} style={{ background:'rgba(255,255,255,0.02)', border:'1px solid rgba(120,140,200,0.1)', borderRadius:'11px', padding:'15px 17px' }}>
                <div style={{ fontSize:'11px', fontWeight:700, color:'rgba(200,210,235,0.45)', marginBottom:'9px', letterSpacing:'.5px', textTransform:'uppercase' as const }}>{label}</div>
                <div style={{ display:'flex', alignItems:'center', gap:'8px' }}>
                  <input type="number" min={0} style={{ width:'74px', background:'rgba(255,255,255,0.05)', border:'1px solid rgba(120,140,200,0.2)', borderRadius:'8px', padding:'7px 10px', fontSize:'18px', fontWeight:700, color:'#E8EAF2', fontFamily:'inherit', textAlign:'center' as const, outline:'none' }}
                    value={policy[key as keyof PointPolicy]}
                    onChange={e => setPolicy({ ...policy, [key]: parseInt(e.target.value) || 0 })} />
                  <span style={{ fontSize:'13px', color:'rgba(200,210,235,0.4)' }}>{unit}</span>
                </div>
                <div style={{ fontSize:'11px', color:'rgba(200,210,235,0.3)', marginTop:'5px' }}>{desc}</div>
              </div>
            ))}
          </div>
          <div style={S.saveRow}><button style={S.btnSave} onClick={savePolicy} disabled={saving}>포인트 정책 저장</button></div>
        </div>
      )}
    </div>
  )
}
