'use client'
import { S } from './settingsStyles'
import { FEATURE_LABELS, SLOT_LABELS } from './settingsConstants'

export function FeatureTag({ feature }: { feature: string }) {
  const normalized = feature.replace(/^pro_/, '')
  const colors: Record<string,{bg:string,color:string}> = {
    default:    { bg:'rgba(120,140,200,0.12)', color:'rgba(200,210,235,0.55)' },
    curriculum: { bg:'rgba(239,159,39,0.15)',  color:'#FAC775' },
    tutor:      { bg:'rgba(127,119,221,0.15)', color:'#9B96E8' },
    vision:     { bg:'rgba(255,140,50,0.15)',  color:'#FF8C32' },
    quiz:       { bg:'rgba(55,138,221,0.15)',  color:'#378ADD' },
  }
  const c = colors[normalized] || colors.default
  return <span style={{ display:'inline-block', padding:'2px 9px', borderRadius:'100px', fontSize:'10px', fontWeight:700, background:c.bg, color:c.color }}>{FEATURE_LABELS[feature] || feature}</span>
}

export function SlotTag({ type }: { type: string }) {
  const colors: Record<string,{bg:string,color:string}> = {
    category_sponsor: { bg:'rgba(127,119,221,0.15)', color:'#9B96E8' },
    step_complete:    { bg:'rgba(82,183,136,0.15)',  color:'#52B788' },
    dashboard_banner: { bg:'rgba(55,138,221,0.15)',  color:'#378ADD' },
    creator_sponsor:  { bg:'rgba(255,140,50,0.15)',  color:'#FF8C32' },
  }
  const c = colors[type] || { bg:'rgba(120,140,200,0.12)', color:'rgba(200,210,235,0.55)' }
  return <span style={{ display:'inline-block', padding:'2px 8px', borderRadius:'100px', fontSize:'10px', fontWeight:700, background:c.bg, color:c.color, whiteSpace:'nowrap' }}>{SLOT_LABELS[type] || type}</span>
}

export function Toggle({ value, onChange }: { value: boolean; onChange: (v: boolean) => void }) {
  return (
    <div style={S.toggle(value)} onClick={() => onChange(!value)}>
      <div style={S.toggleKnob(value)} />
    </div>
  )
}

export function DocStateTag({ label, active }: { label: string; active?: boolean }) {
  return (
    <span style={{
      display: 'inline-block',
      padding: '3px 10px',
      borderRadius: '999px',
      fontSize: '11px',
      fontWeight: 700,
      background: active ? 'rgba(82,183,136,0.16)' : 'rgba(239,159,39,0.16)',
      color: active ? '#52B788' : '#FAC775',
      whiteSpace: 'nowrap',
    }}>
      {label}
    </span>
  )
}
