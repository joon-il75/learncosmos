'use client'
import type { Dispatch, SetStateAction } from 'react'
import type { AdSlot, AffiliateSettings } from '../settingsTypes'
import { SLOT_TYPES, SLOT_LABELS, CATEGORIES } from '../settingsConstants'
import { S, OPTION_STYLE } from '../settingsStyles'
import { SlotTag, Toggle } from '../settingsComponents'

interface Props {
  adSlots: AdSlot[]
  setAdSlots: Dispatch<SetStateAction<AdSlot[]>>
  affiliate: AffiliateSettings
  setAffiliate: Dispatch<SetStateAction<AffiliateSettings>>
  saving: boolean
  saveAdSlots: () => Promise<void>
  saveAffiliate: () => Promise<void>
}

export function AdsSettingsTab({
  adSlots, setAdSlots, affiliate, setAffiliate,
  saving, saveAdSlots, saveAffiliate,
}: Props) {
  return (
    <div style={{ paddingTop:'20px' }}>
      <div style={S.warnBox}>
        ⚡ 모든 광고의 위치·시점·내용은 슈퍼관리자가 직접 제어합니다.<br/>
        ❌ 구글 애드센스 자동 광고 미도입 · 팝업·인터스티셜 · 자동재생 영상·음성 금지
      </div>

      <div style={S.card}>
        <div style={S.cardTitle}>광고 슬롯 관리</div>
        <div style={S.cardSub}>ad_slots 테이블 · 슬롯별 제목·링크·카테고리·활성화 직접 제어</div>
        <div style={{ display:'grid', gridTemplateColumns:'110px 1fr 150px 110px 60px 50px', gap:'10px', paddingBottom:'9px', borderBottom:'1px solid rgba(120,140,200,0.15)', fontSize:'11px', fontWeight:700, color:'rgba(200,210,235,0.35)', letterSpacing:'.4px' }}>
          <span>슬롯 유형</span><span>제목</span><span>링크 URL</span><span>카테고리</span><span>활성화</span><span>삭제</span>
        </div>
        {adSlots.map((s, i) => (
          <div key={s.id || i} style={{ display:'grid', gridTemplateColumns:'110px 1fr 150px 110px 60px 50px', gap:'10px', alignItems:'center', padding:'10px 0', borderBottom:'1px solid rgba(120,140,200,0.07)' }}>
            <div>
              <SlotTag type={s.slot_type} />
              <select style={{ ...S.sel, marginTop:'4px', fontSize:'10px', padding:'3px 6px' }}
                value={s.slot_type}
                onChange={e => setAdSlots(prev => prev.map((x,j) => j===i ? {...x,slot_type:e.target.value} : x))}>
                {SLOT_TYPES.map(t => <option style={OPTION_STYLE} key={t} value={t}>{SLOT_LABELS[t]}</option>)}
              </select>
            </div>
            <input style={{ ...S.inp, fontSize:'12px', padding:'7px 10px' }} value={s.title}
              onChange={e => setAdSlots(prev => prev.map((x,j) => j===i ? {...x,title:e.target.value} : x))} />
            <input style={{ ...S.inp, fontSize:'11px', padding:'7px 10px' }} value={s.link_url} placeholder="https://..."
              onChange={e => setAdSlots(prev => prev.map((x,j) => j===i ? {...x,link_url:e.target.value} : x))} />
            <select style={{ ...S.sel, padding:'7px 8px', fontSize:'11px' }}
              value={s.category_slug || ''}
              onChange={e => setAdSlots(prev => prev.map((x,j) => j===i ? {...x,category_slug:e.target.value} : x))}>
              <option style={OPTION_STYLE} value="">전체</option>
              {CATEGORIES.map(c => <option style={OPTION_STYLE} key={c} value={c}>{c}</option>)}
            </select>
            <Toggle value={s.is_active} onChange={v => setAdSlots(prev => prev.map((x,j) => j===i ? {...x,is_active:v} : x))} />
            <button style={S.btnDel} onClick={() => setAdSlots(prev => prev.filter((_,j) => j!==i))}>삭제</button>
          </div>
        ))}
        <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginTop:'14px' }}>
          <button style={S.btnAdd} onClick={() => setAdSlots(prev => [...prev, { slot_type:'category_sponsor', title:'', link_url:'', is_active:false, sort_order:prev.length+1 }])}>+ 슬롯 추가</button>
          <button style={S.btnSave} onClick={saveAdSlots} disabled={saving}>슬롯 저장</button>
        </div>
      </div>

      <div style={S.card}>
        <div style={S.cardTitle}>제휴 프로그램 설정</div>
        <div style={S.cardSub}>MVP 5종 제휴 — 계약 없이 즉시 시작 가능</div>
        <div style={S.grid2}>
          {[
            { key:'coupang',  label:'🛒 쿠팡 파트너스 추적 ID',    ph:'af_id=learnweaver01' },
            { key:'naver',    label:'🛍️ 네이버 쇼핑 파트너스 ID',  ph:'파트너스 ID' },
            { key:'class101', label:'🎓 클래스101 파트너스 코드',   ph:'ref=...' },
            { key:'kyobo',    label:'📚 교보문고 제휴 코드',        ph:'제휴 코드' },
            { key:'adpick',   label:'🔗 애드픽 트래킹 ID',         ph:'애드픽 ID' },
          ].map(({ key, label, ph }) => (
            <div key={key}>
              <label style={S.label}>{label}</label>
              <input style={S.inp} placeholder={ph}
                value={affiliate[key as keyof AffiliateSettings]}
                onChange={e => setAffiliate(prev => ({ ...prev, [key]: e.target.value }))} />
            </div>
          ))}
        </div>
        <div style={S.saveRow}><button style={S.btnSave} onClick={saveAffiliate} disabled={saving}>제휴 설정 저장</button></div>
      </div>

      <div style={S.card}>
        <div style={S.cardTitle}>📵 광고 절대 금지 위치 (시스템 고정)</div>
        <div style={S.cardSub}>아래 위치에는 어떠한 광고도 표시되지 않습니다</div>
        <div style={{ display:'flex', flexDirection:'column', gap:'8px' }}>
          {[
            '레슨 학습 화면 — 집중 방해 금지',
            'Feynman Test · 퀴즈 · 메타인지 평가 화면',
            '결제 · 인증 화면',
            '에러 · 로딩 화면',
            '구독 플랜 가입자 (모든 광고 제거)',
            '관리자 · 슈퍼관리자 계정',
          ].map(item => (
            <div key={item} style={{ display:'flex', alignItems:'center', gap:'10px', background:'rgba(226,75,74,0.06)', border:'1px solid rgba(226,75,74,0.15)', borderRadius:'8px', padding:'9px 14px', fontSize:'12px', color:'rgba(200,210,235,0.6)' }}>
              <span style={{ color:'#F09595', fontSize:'14px' }}>❌</span>
              {item}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
