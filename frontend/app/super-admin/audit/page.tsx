'use client'

import { useEffect, useState, useCallback } from 'react'
import { useRouter } from 'next/navigation'
import SuperAdminPanelNav from '@/components/super-admin/SuperAdminPanelNav'
import SuperAdminPanelHeader from '@/components/super-admin/SuperAdminPanelHeader'
import { superAdminShellStyle } from '@/components/super-admin/layout'

const API_BASE = ''

interface AuditLog {
  id: string
  admin_id: string
  admin_display: string
  target_user_id: string
  target_identifier: string
  amount: number
  memo: string
  created_at: string
}

export default function AuditPage() {
  const router = useRouter()
  const [logs, setLogs] = useState<AuditLog[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const limit = 20

  const fetchLogs = useCallback(async (p: number) => {
    setLoading(true)
    try {
      const res = await fetch(
        `${API_BASE}/api/v1/super-admin/audit/point-grants?page=${p}&limit=${limit}`,
        {
          credentials: 'include',
        }
      )
      if (res.status === 401 || res.status === 403) { router.replace('/super-admin/login'); return }
      if (!res.ok) throw new Error('조회 실패')
      const data = await res.json()
      setLogs(data.logs ?? [])
      setTotalPages(Math.ceil((data.total ?? 0) / limit) || 1)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : '오류가 발생했습니다.')
    } finally {
      setLoading(false)
    }
  }, [router])

  useEffect(() => {
    fetchLogs(page)
  }, [page, fetchLogs])

  const formatDate = (iso: string) =>
    new Date(iso).toLocaleString('ko-KR', { timeZone: 'Asia/Seoul' })

  return (
    <div style={{ minHeight:'100vh', background:'#081826', color:'#E2E8F0', padding:'32px' }}>
      <div style={superAdminShellStyle}>
      {/* 헤더 */}
      <div style={{ display:'grid', gap:'16px', marginBottom:'32px' }}>
        <SuperAdminPanelHeader subtitle="운영 로그" description="관리자의 수동 포인트 지급 내역을 중심으로 운영 변경 이력을 확인합니다." />
        <SuperAdminPanelNav activeSection="audit" />
      </div>

      {/* 에러 */}
      {error && (
        <div style={{ background:'#7F1D1D', border:'1px solid #991B1B', borderRadius:'8px', padding:'12px 16px', marginBottom:'20px', color:'#FCA5A5', fontSize:'14px' }}>
          {error}
        </div>
      )}

      {/* 테이블 */}
      <div style={{ background:'#0F2337', border:'1px solid #1E3A5F', borderRadius:'12px', overflow:'hidden' }}>
        <table style={{ width:'100%', borderCollapse:'collapse', fontSize:'13px' }}>
          <thead>
            <tr style={{ background:'#0B1629', borderBottom:'1px solid #1E3A5F' }}>
              {['지급 일시', '지급 관리자', '수신 사용자', '지급량', '메모'].map(h => (
                <th key={h} style={{ padding:'12px 16px', textAlign:'left', color:'#64748B', fontWeight:600 }}>{h}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr><td colSpan={5} style={{ padding:'40px', textAlign:'center', color:'#64748B' }}>불러오는 중...</td></tr>
            ) : logs.length === 0 ? (
              <tr><td colSpan={5} style={{ padding:'40px', textAlign:'center', color:'#64748B' }}>지급 이력이 없습니다.</td></tr>
            ) : logs.map((log, idx) => (
              <tr key={log.id} style={{ borderBottom:'1px solid #1E3A5F', background: idx % 2 === 0 ? 'transparent' : 'rgba(30,58,95,0.2)' }}>
                <td style={{ padding:'12px 16px', color:'#94A3B8' }}>{formatDate(log.created_at)}</td>
                <td style={{ padding:'12px 16px' }}>{log.admin_display}</td>
                <td style={{ padding:'12px 16px', color:'#CBD5E1' }}>{log.target_identifier}</td>
                <td style={{ padding:'12px 16px' }}>
                  <span style={{ background:'rgba(239,159,39,0.15)', color:'#EF9F27', border:'1px solid rgba(239,159,39,0.3)', borderRadius:'4px', padding:'2px 8px', fontWeight:700 }}>
                    +{log.amount}pt
                  </span>
                </td>
                <td style={{ padding:'12px 16px', color:'#94A3B8' }}>{log.memo || '—'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* 페이지네이션 */}
      {totalPages > 1 && (
        <div style={{ display:'flex', justifyContent:'center', gap:'8px', marginTop:'24px' }}>
          <button
            onClick={() => setPage(p => Math.max(1, p - 1))}
            disabled={page === 1}
            style={{ padding:'6px 14px', borderRadius:'6px', background: page === 1 ? '#0F2337' : '#1E3A5F', border:'1px solid #1E3A5F', color: page === 1 ? '#475569' : '#E2E8F0', cursor: page === 1 ? 'not-allowed' : 'pointer', fontSize:'13px' }}
          >이전</button>
          <span style={{ padding:'6px 14px', color:'#94A3B8', fontSize:'13px' }}>{page} / {totalPages}</span>
          <button
            onClick={() => setPage(p => Math.min(totalPages, p + 1))}
            disabled={page === totalPages}
            style={{ padding:'6px 14px', borderRadius:'6px', background: page === totalPages ? '#0F2337' : '#1E3A5F', border:'1px solid #1E3A5F', color: page === totalPages ? '#475569' : '#E2E8F0', cursor: page === totalPages ? 'not-allowed' : 'pointer', fontSize:'13px' }}
          >다음</button>
        </div>
      )}
      </div>
    </div>
  )
}
