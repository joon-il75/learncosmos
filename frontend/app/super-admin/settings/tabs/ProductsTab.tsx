'use client'
import type { Dispatch, SetStateAction } from 'react'
import type { Product } from '../settingsTypes'
import { S, OPTION_STYLE } from '../settingsStyles'
import { Toggle } from '../settingsComponents'

interface Props {
  products: Product[]
  setProducts: Dispatch<SetStateAction<Product[]>>
  saving: boolean
  saveProducts: () => Promise<void>
}

export function ProductsTab({ products, setProducts, saving, saveProducts }: Props) {
  return (
    <div style={{ paddingTop:'20px' }}>
      <div style={S.infoBox}>💳 활성화된 상품만 학습자 구매 화면에 표시됩니다. 포인트 단건 상품과 Pro 월간/연간 구독 상품을 함께 관리합니다. 현재 권장값은 Pro 월 300pt, 프리미엄 팩 11,900원입니다.</div>
      <div style={S.card}>
        <div style={S.cardTitle}>포인트 구매 상품 목록</div>
        <div style={S.cardSub}>단건 포인트 상품과 Pro 구독 상품(월간/연간) 관리 · Pro 권장 설명: 커리큘럼 5pt · 튜터 1pt · 비전 5pt(긴 변 1280px 제한) · 퀴즈 2pt</div>
        <div style={{ display:'grid', gridTemplateColumns:'1.2fr 90px 90px 90px 80px 70px 70px 50px', gap:'10px', paddingBottom:'9px', borderBottom:'1px solid rgba(120,140,200,0.15)', fontSize:'11px', fontWeight:700, color:'rgba(200,210,235,0.35)', letterSpacing:'.5px' }}>
          <span>상품명</span><span>유형</span><span>주기</span><span>포인트</span><span>가격(원)</span><span>광고 제거</span><span>활성화</span><span>삭제</span>
        </div>
        {products.map((p, i) => (
          <div key={p.id || i} style={{ padding:'12px 0', borderBottom:'1px solid rgba(120,140,200,0.07)' }}>
            <div style={{ display:'grid', gridTemplateColumns:'1.2fr 90px 90px 90px 80px 70px 70px 50px', gap:'10px', alignItems:'center' }}>
              <input style={{ ...S.inp, fontSize:'12px', padding:'7px 10px' }} value={p.name}
                onChange={e => setProducts(prev => prev.map((x,j) => j===i ? {...x,name:e.target.value} : x))} />
              <select style={{ ...S.sel, padding:'7px 8px', fontSize:'11px' }}
                value={p.product_type}
                onChange={e => setProducts(prev => prev.map((x,j) => j===i ? {
                  ...x,
                  product_type: e.target.value as Product['product_type'],
                  billing_period: e.target.value === 'subscription' ? (x.billing_period || 'monthly') : undefined,
                } : x))}>
                <option style={OPTION_STYLE} value="one_time">단건</option>
                <option style={OPTION_STYLE} value="subscription">구독</option>
              </select>
              <select style={{ ...S.sel, padding:'7px 8px', fontSize:'11px' }}
                value={p.billing_period || ''}
                disabled={p.product_type !== 'subscription'}
                onChange={e => setProducts(prev => prev.map((x,j) => j===i ? {...x,billing_period:e.target.value as Product['billing_period']} : x))}>
                <option style={OPTION_STYLE} value="">-</option>
                <option style={OPTION_STYLE} value="monthly">월간</option>
                <option style={OPTION_STYLE} value="yearly">연간</option>
              </select>
              <div style={{ display:'flex', alignItems:'center', gap:'4px' }}>
                <input type="number" style={{ width:'70px', background:'rgba(255,255,255,0.04)', border:'1px solid rgba(120,140,200,0.2)', borderRadius:'8px', padding:'7px 8px', fontSize:'13px', fontWeight:700, color:'#FAC775', textAlign:'center' as const, outline:'none', fontFamily:'inherit' }}
                  value={p.product_type === 'subscription' ? p.monthly_points : p.points}
                  onChange={e => setProducts(prev => prev.map((x,j) => j===i ? {
                    ...x,
                    ...(x.product_type === 'subscription'
                      ? { monthly_points: parseInt(e.target.value) || 0 }
                      : { points: parseInt(e.target.value) || 0 }),
                  } : x))} />
                <span style={{ fontSize:'11px', color:'rgba(200,210,235,0.35)' }}>pt</span>
              </div>
              <input type="number" style={{ width:'80px', background:'rgba(255,255,255,0.04)', border:'1px solid rgba(120,140,200,0.2)', borderRadius:'8px', padding:'7px 8px', fontSize:'13px', textAlign:'center' as const, outline:'none', color:'#E8EAF2', fontFamily:'inherit' }}
                value={p.price_krw}
                onChange={e => setProducts(prev => prev.map((x,j) => j===i ? {...x,price_krw:parseInt(e.target.value)||0} : x))} />
              <Toggle value={p.ad_free} onChange={v => setProducts(prev => prev.map((x,j) => j===i ? {...x,ad_free:v} : x))} />
              <Toggle value={p.is_active} onChange={v => setProducts(prev => prev.map((x,j) => j===i ? {...x,is_active:v} : x))} />
              <button style={S.btnDel} onClick={() => setProducts(prev => prev.filter((_,j) => j!==i))}>삭제</button>
            </div>
            {p.product_type === 'subscription' && (
              <div style={{ display:'grid', gridTemplateColumns:'1fr 90px', gap:'10px', marginTop:'10px' }}>
                <input style={{ ...S.inp, fontSize:'12px', padding:'7px 10px' }}
                  placeholder="혜택 설명"
                  value={p.description || ''}
                  onChange={e => setProducts(prev => prev.map((x,j) => j===i ? {...x,description:e.target.value} : x))} />
                <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', background:'rgba(255,255,255,0.02)', border:'1px solid rgba(120,140,200,0.1)', borderRadius:'9px', padding:'8px 10px' }}>
                  <span style={{ fontSize:'11px', color:'rgba(200,210,235,0.55)' }}>프리미엄</span>
                  <Toggle value={p.premium_access} onChange={v => setProducts(prev => prev.map((x,j) => j===i ? {...x,premium_access:v} : x))} />
                </div>
              </div>
            )}
          </div>
        ))}
        <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginTop:'14px' }}>
          <div style={{ display:'flex', gap:'8px' }}>
            <button style={S.btnAdd} onClick={() => setProducts(prev => [...prev, {
              name:'새 포인트 상품', points:0, price_krw:0, is_active:false, sort_order:prev.length+1,
              product_type:'one_time', monthly_points:0, ad_free:false, premium_access:false,
            }])}>+ 단건 상품</button>
            <button style={S.btnAdd} onClick={() => setProducts(prev => [...prev, {
              name:'LearnCosmos Pro', points:0, price_krw:0, is_active:false, sort_order:prev.length+1,
              product_type:'subscription', billing_period:'monthly', monthly_points:300, ad_free:true, premium_access:true,
              description:'월 300pt · 커리큘럼 5pt · 튜터 1pt · 비전 코칭 5pt(긴 변 1280px 제한) · 퀴즈 2pt · 커리큘럼: Claude 3.5 Sonnet · 튜터: Grok · 비전 코칭: GPT-4o · 퀴즈: GPT-4o-mini · 광고 제거',
            }])}>+ Pro 구독</button>
          </div>
          <button style={S.btnSave} onClick={saveProducts} disabled={saving}>상품 목록 저장</button>
        </div>
      </div>
    </div>
  )
}
