import type React from 'react'
import { SUPER_ADMIN_PAGE_WIDTH } from '@/components/super-admin/layout'

export const S = {
  page: { minHeight:'100vh', background:'#081826', color:'#E8EAF2', padding:'32px 24px 56px', fontFamily:"'Noto Sans KR', sans-serif" } as React.CSSProperties,
  shell: { width:SUPER_ADMIN_PAGE_WIDTH, margin:'0 auto', display:'grid', gap:'22px' } as React.CSSProperties,
  topbar: { display:'flex', justifyContent:'space-between', alignItems:'flex-start', gap:'18px', flexWrap:'wrap' as const, marginBottom:'2px' } as React.CSSProperties,
  title: { fontSize:'24px', fontWeight:700 } as React.CSSProperties,
  sub: { fontSize:'13px', color:'rgba(200,210,235,0.55)', marginTop:'6px', lineHeight:1.7, maxWidth:'760px' } as React.CSSProperties,
  btnNav: { background:'rgba(255,255,255,0.05)', border:'1px solid rgba(120,140,200,0.2)', color:'rgba(200,210,235,0.8)', borderRadius:'10px', padding:'9px 14px', fontSize:'13px', cursor:'pointer' } as React.CSSProperties,
  card: { background:'rgba(17,30,53,0.72)', border:'1px solid rgba(120,140,200,0.14)', borderRadius:'18px', padding:'22px', marginBottom:0 } as React.CSSProperties,
  label: { display:'block', fontSize:'12px', fontWeight:700, color:'rgba(200,210,235,0.6)', marginBottom:'7px' } as React.CSSProperties,
  input: { width:'100%', background:'rgba(255,255,255,0.05)', border:'1px solid rgba(120,140,200,0.18)', borderRadius:'10px', padding:'11px 12px', color:'#E8EAF2', fontSize:'13px', outline:'none', boxSizing:'border-box' as const, resize:'vertical' as const } as React.CSSProperties,
  primary: { background:'linear-gradient(135deg,#378ADD,#7F77DD)', border:'none', color:'#fff', borderRadius:'10px', padding:'11px 16px', fontWeight:700, fontSize:'13px', cursor:'pointer' } as React.CSSProperties,
  secondary: { background:'rgba(255,255,255,0.05)', border:'1px solid rgba(120,140,200,0.18)', color:'rgba(200,210,235,0.8)', borderRadius:'10px', padding:'11px 16px', fontWeight:700, fontSize:'13px', cursor:'pointer' } as React.CSSProperties,
  info: { background:'rgba(55,138,221,0.08)', border:'1px solid rgba(55,138,221,0.2)', borderRadius:'14px', padding:'14px 16px', fontSize:'12px', color:'rgba(200,210,235,0.75)', lineHeight:1.8 } as React.CSSProperties,
  warn: { background:'rgba(239,159,39,0.08)', border:'1px solid rgba(239,159,39,0.2)', borderRadius:'14px', padding:'14px 16px', fontSize:'12px', color:'#FAC775', lineHeight:1.8 } as React.CSSProperties,
}
