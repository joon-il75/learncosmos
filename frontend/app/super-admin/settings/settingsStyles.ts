import type { CSSProperties } from 'react'

export const S = {
  page:         { minHeight:'100vh', background:'#081826', color:'#E8EAF2', padding:'24px', fontFamily:"'Noto Sans KR', sans-serif", overflowX:'hidden' } as CSSProperties,
  topbar:       { display:'flex', justifyContent:'space-between', alignItems:'flex-start', marginBottom:'6px', gap:'12px', flexWrap:'wrap' as const } as CSSProperties,
  title:        { fontSize:'22px', fontWeight:700 } as CSSProperties,
  sub:          { fontSize:'13px', color:'rgba(200,210,235,0.5)', marginTop:'3px' } as CSSProperties,
  btnNav:       { background:'rgba(255,255,255,0.05)', border:'1px solid rgba(120,140,200,0.2)', color:'rgba(200,210,235,0.7)', borderRadius:'8px', padding:'7px 14px', fontSize:'12px', cursor:'pointer', marginLeft:'8px' } as CSSProperties,
  btnLogout:    { background:'transparent', border:'1px solid rgba(120,140,200,0.15)', color:'rgba(200,210,235,0.4)', borderRadius:'8px', padding:'7px 14px', fontSize:'12px', cursor:'pointer', marginLeft:'8px' } as CSSProperties,
  tabs:         { display:'flex', gap:'6px', borderBottom:'1px solid rgba(120,140,200,0.15)', marginTop:'20px', flexWrap:'wrap' as const, paddingBottom:'6px' } as CSSProperties,
  tab:          (on:boolean): CSSProperties => ({ padding:'10px 18px', fontSize:'13px', fontWeight:500, cursor:'pointer', borderRadius:'8px 8px 0 0', color: on ? '#E8EAF2' : 'rgba(200,210,235,0.45)', background: on ? 'rgba(17,30,53,0.95)' : 'transparent', border: on ? '1px solid rgba(120,140,200,0.2)' : '1px solid transparent', borderBottom: on ? '1px solid rgba(17,30,53,0.95)' : '1px solid transparent', marginBottom:'-1px', transition:'all .15s', whiteSpace:'nowrap' }),
  card:         { background:'rgba(17,30,53,0.7)', border:'1px solid rgba(120,140,200,0.15)', borderRadius:'14px', padding:'22px', marginBottom:'14px' } as CSSProperties,
  cardTitle:    { fontSize:'14px', fontWeight:700, marginBottom:'3px' } as CSSProperties,
  cardSub:      { fontSize:'12px', color:'rgba(200,210,235,0.4)', marginBottom:'16px' } as CSSProperties,
  label:        { display:'block', fontSize:'12px', fontWeight:600, color:'rgba(200,210,235,0.6)', marginBottom:'6px', letterSpacing:'0.3px' } as CSSProperties,
  inp:          { width:'100%', background:'rgba(255,255,255,0.04)', border:'1px solid rgba(120,140,200,0.2)', borderRadius:'9px', padding:'9px 13px', fontSize:'13px', color:'#E8EAF2', fontFamily:'inherit', outline:'none', boxSizing:'border-box' as const } as CSSProperties,
  sel:          { width:'100%', background:'rgba(255,255,255,0.04)', border:'1px solid rgba(120,140,200,0.2)', borderRadius:'9px', padding:'9px 13px', fontSize:'12px', color:'#E8EAF2', fontFamily:'inherit', outline:'none' } as CSSProperties,
  grid2:        { display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(260px, 1fr))', gap:'12px' } as CSSProperties,
  infoBox:      { background:'rgba(55,138,221,0.07)', border:'1px solid rgba(55,138,221,0.18)', borderRadius:'10px', padding:'11px 15px', fontSize:'12px', color:'rgba(200,210,235,0.6)', marginBottom:'14px', lineHeight:1.75 } as CSSProperties,
  warnBox:      { background:'rgba(239,159,39,0.07)', border:'1px solid rgba(239,159,39,0.2)', borderRadius:'10px', padding:'11px 15px', fontSize:'12px', color:'rgba(250,199,117,0.75)', marginBottom:'14px', lineHeight:1.75 } as CSSProperties,
  btnSave:      { background:'linear-gradient(135deg,#378ADD,#7F77DD)', color:'#fff', border:'none', borderRadius:'9px', padding:'9px 22px', fontSize:'13px', fontWeight:700, cursor:'pointer', fontFamily:'inherit' } as CSSProperties,
  btnSecondary: { background:'rgba(255,255,255,0.05)', border:'1px solid rgba(120,140,200,0.2)', color:'rgba(200,210,235,0.78)', borderRadius:'9px', padding:'9px 18px', fontSize:'13px', fontWeight:600, cursor:'pointer', fontFamily:'inherit' } as CSSProperties,
  btnAdd:       { background:'rgba(82,183,136,0.12)', color:'#52B788', border:'1px solid rgba(82,183,136,0.28)', borderRadius:'8px', padding:'7px 16px', fontSize:'12px', fontWeight:600, cursor:'pointer', fontFamily:'inherit' } as CSSProperties,
  btnDel:       { background:'rgba(226,75,74,0.1)', color:'#F09595', border:'1px solid rgba(226,75,74,0.22)', borderRadius:'6px', padding:'4px 10px', fontSize:'11px', fontWeight:600, cursor:'pointer', fontFamily:'inherit' } as CSSProperties,
  saveRow:      { display:'flex', justifyContent:'flex-end', marginTop:'14px' } as CSSProperties,
  toggle:       (on:boolean): CSSProperties => ({ width:'34px', height:'19px', background: on ? 'rgba(82,183,136,0.8)' : 'rgba(120,140,200,0.2)', borderRadius:'10px', cursor:'pointer', position:'relative', flexShrink:0, display:'inline-block' }),
  toggleKnob:   (on:boolean): CSSProperties => ({ position:'absolute', width:'15px', height:'15px', background:'white', borderRadius:'50%', top:'2px', left: on ? 'auto' : '2px', right: on ? '2px' : 'auto', transition:'all .2s' }),
}

export const OPTION_STYLE: CSSProperties = {
  background: '#10233A',
  color: '#E8EAF2',
}
