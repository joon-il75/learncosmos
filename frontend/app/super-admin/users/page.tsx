'use client'

import React, { useEffect, useState, useCallback } from 'react'
import { useRouter } from 'next/navigation'
import SuperAdminPanelNav from '@/components/super-admin/SuperAdminPanelNav'
import SuperAdminPanelHeader from '@/components/super-admin/SuperAdminPanelHeader'
import { superAdminShellStyle } from '@/components/super-admin/layout'

const API_BASE = ''

interface User {
  id: string
  email: string
  nickname: string
  role: string
  premium_access: boolean
  provider: 'google' | 'kakao' | 'naver' | 'unknown'
  display_id: string | null
  created_at: string
  totp_reset_requested: boolean
}

const PROVIDER_ICON: Record<string, string> = {
  google: '🔵',
  kakao:  '🟡',
  naver:  '🟢',
}

function getIdentifier(user: User): string {
  if (user.email && user.email !== '') return user.email
  if (user.display_id) return user.display_id
  return `Unknown #${user.id.slice(0, 6)}`
}

function getIdentifierColor(user: User): string {
  return (user.email && user.email !== '') ? '#E8EAF2' : 'rgba(200,210,235,0.6)'
}

function getRoleBadgeStyle(role: string): React.CSSProperties {
  switch (role) {
    case 'admin':
      return { background: 'rgba(82,183,136,0.18)', color: '#52B788', padding: '3px 10px', borderRadius: '100px', fontSize: '11px', fontWeight: 600, display: 'inline-block', whiteSpace: 'nowrap' }
    case 'learner':
      return { background: 'rgba(255,220,80,0.15)', color: '#FFD84D', padding: '3px 10px', borderRadius: '100px', fontSize: '11px', fontWeight: 600, display: 'inline-block', whiteSpace: 'nowrap' }
    case 'creator':
      return { background: 'rgba(255,140,50,0.18)', color: '#FF8C32', padding: '3px 10px', borderRadius: '100px', fontSize: '11px', fontWeight: 600, display: 'inline-block', whiteSpace: 'nowrap' }
    default:
      return { background: 'rgba(120,140,200,0.15)', color: 'rgba(200,210,235,0.7)', padding: '3px 10px', borderRadius: '100px', fontSize: '11px', fontWeight: 600, display: 'inline-block' }
  }
}

interface UsersResponse {
  users: User[]
  total: number
  page: number
  limit: number
}

export default function SuperAdminUsersPage() {
  const router = useRouter()
  const [users, setUsers] = useState<User[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [query, setQuery] = useState('')
  const [searchInput, setSearchInput] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const limit = 20

  const [pointModal, setPointModal] = useState<{
    open: boolean; userId: string; userName: string; amount: number; memo: string;
  }>({ open: false, userId: '', userName: '', amount: 10, memo: '' })
  const [pointLoading, setPointLoading] = useState(false)
  const [pointError, setPointError] = useState('')
  const [welcomePoints, setWelcomePoints] = useState(30)
  const [courseGenCost, setCourseGenCost] = useState(5)
  const [lessonRecCost, setLessonRecCost] = useState(1)

  const fetchUsers = useCallback(async (q: string, p: number) => {
    setLoading(true)
    setError('')
    try {
      const params = new URLSearchParams({ q, page: String(p), limit: String(limit) })
      const res = await fetch(`${API_BASE}/api/v1/super-admin/users?${params}`, {
        credentials: 'include',
      })

      if (res.status === 401 || res.status === 403) {
        router.replace('/super-admin/login')
        return
      }

      if (!res.ok) {
        const data = await res.json().catch(() => null)
        throw new Error(data?.error ?? `사용자 목록을 불러오지 못했습니다. (${res.status})`)
      }

      const data: UsersResponse = await res.json()
      setUsers(data.users ?? [])
      setTotal(data.total ?? 0)
    } catch (err: unknown) {
      setUsers([])
      setTotal(0)
      setError(err instanceof Error ? err.message : '사용자 목록을 불러오지 못했습니다.')
    } finally {
      setLoading(false)
    }
  }, [router])

  useEffect(() => {
    fetchUsers(query, page)
    fetch(`${API_BASE}/api/v1/super-admin/settings/point-policy`, {
      credentials: 'include',
    })
      .then(r => r.ok ? r.json() : null)
      .then(d => {
        const policy = d?.policy
        if (!policy) return
        if (policy.welcome_points)  setWelcomePoints(policy.welcome_points)
        if (policy.course_gen_cost) setCourseGenCost(policy.course_gen_cost)
        if (policy.lesson_rec_cost) setLessonRecCost(policy.lesson_rec_cost)
      })
      .catch(() => {})
  }, [query, page, fetchUsers])

  const handleSearch = () => {
    setPage(1)
    setQuery(searchInput)
  }

  const handleRoleChange = async (userId: string, newRole: string) => {
    await fetch(`${API_BASE}/api/v1/super-admin/users/${userId}/role`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ role: newRole }),
    })

    fetchUsers(query, page)
  }

  const handlePremiumAccessChange = async (userId: string, premiumAccess: boolean) => {
    await fetch(`${API_BASE}/api/v1/super-admin/users/${userId}/premium-access`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ premium_access: premiumAccess }),
    })

    fetchUsers(query, page)
  }

  const handleLogout = async () => {
    await fetch('/api/v1/auth/logout', { method: 'POST', credentials: 'include' }).catch(() => undefined)
    window.location.href = '/super-admin/login'
  }

  const handleTotpReset = async (userId: string) => {
    if (!confirm('이 관리자의 TOTP를 초기화하시겠습니까?')) return

    await fetch(`${API_BASE}/api/v1/super-admin/users/${userId}/totp-reset`, {
      method: 'POST',
      credentials: 'include',
    })
    fetchUsers(query, page)
  }

  const handleGrantPoint = async () => {
    if (pointModal.amount < 1 || pointModal.amount > 50) {
      setPointError('지급량은 1~50pt 사이여야 합니다.')
      return
    }
    setPointLoading(true)
    setPointError('')
    try {
      const res = await fetch(`${API_BASE}/api/v1/super-admin/users/${pointModal.userId}/grant-points`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ amount: pointModal.amount, memo: pointModal.memo }),
      })
      if (!res.ok) {
        const data = await res.json()
        throw new Error(data.error ?? '포인트 지급 실패')
      }
      setPointModal({ open: false, userId: '', userName: '', amount: 10, memo: '' })
      fetchUsers(query, page)
    } catch (err: unknown) {
      setPointError(err instanceof Error ? err.message : '오류가 발생했습니다.')
    } finally {
      setPointLoading(false)
    }
  }

  const totalPages = Math.ceil(total / limit)

  return (
    <main className="min-h-screen bg-[#0B1629] px-6 py-8">
      {/* 포인트 지급 모달 */}
      {pointModal.open && (
        <div
          style={{ position:'fixed', inset:0, background:'rgba(0,0,0,0.6)', display:'flex', alignItems:'center', justifyContent:'center', zIndex:1000 }}
          onClick={() => setPointModal(prev => ({ ...prev, open: false }))}
        >
          <div
            style={{ background:'#0F2337', border:'1px solid #1E3A5F', borderRadius:'12px', padding:'28px', width:'360px', color:'#E2E8F0' }}
            onClick={e => e.stopPropagation()}
          >
            <h3 style={{ margin:'0 0 4px', fontSize:'16px', fontWeight:700 }}>💎 포인트 지급</h3>
            <p style={{ margin:'0 0 20px', fontSize:'13px', color:'#94A3B8' }}>{pointModal.userName}</p>

            <label style={{ display:'block', fontSize:'13px', marginBottom:'6px', color:'#94A3B8' }}>지급량 (1~50pt)</label>
            <input
              type="number" min={1} max={50}
              value={pointModal.amount}
              onChange={e => setPointModal(prev => ({ ...prev, amount: Number(e.target.value) }))}
              style={{ width:'100%', padding:'8px 12px', background:'#0B1629', border:'1px solid #1E3A5F', borderRadius:'6px', color:'#E2E8F0', fontSize:'14px', boxSizing:'border-box', marginBottom:'14px' }}
            />

            <label style={{ display:'block', fontSize:'13px', marginBottom:'6px', color:'#94A3B8' }}>메모 (선택)</label>
            <input
              type="text" placeholder="지급 사유를 입력하세요"
              value={pointModal.memo}
              onChange={e => setPointModal(prev => ({ ...prev, memo: e.target.value }))}
              style={{ width:'100%', padding:'8px 12px', background:'#0B1629', border:'1px solid #1E3A5F', borderRadius:'6px', color:'#E2E8F0', fontSize:'14px', boxSizing:'border-box', marginBottom:'20px' }}
            />

            {pointError && <p style={{ color:'#FC8181', fontSize:'13px', marginBottom:'12px' }}>{pointError}</p>}

            <div style={{ display:'flex', gap:'10px', justifyContent:'flex-end' }}>
              <button
                onClick={() => setPointModal(prev => ({ ...prev, open: false }))}
                style={{ padding:'8px 16px', borderRadius:'6px', background:'transparent', border:'1px solid #1E3A5F', color:'#94A3B8', cursor:'pointer', fontSize:'13px' }}
              >취소</button>
              <button
                onClick={handleGrantPoint}
                disabled={pointLoading}
                style={{ padding:'8px 20px', borderRadius:'6px', background: pointLoading ? '#7C4A00' : '#EF9F27', border:'none', color:'#fff', cursor: pointLoading ? 'not-allowed' : 'pointer', fontWeight:700, fontSize:'13px' }}
              >{pointLoading ? '지급 중...' : `${pointModal.amount}pt 지급`}</button>
            </div>
          </div>
        </div>
      )}
      <div style={superAdminShellStyle}>
        {/* Header */}
        <div className="mb-8 grid gap-4">
          <SuperAdminPanelHeader subtitle="사용자 관리" description="학습자 계정, 역할, Pro 권한, TOTP 리셋, 포인트 지급을 운영합니다." />
          <SuperAdminPanelNav activeSection="users" />
        </div>

        {/* 역할 권한 안내 카드 */}
        <div style={{ display:'grid', gridTemplateColumns:'repeat(3, 1fr)', gap:'12px', marginBottom:'24px' }}>

          {/* super_admin — 골드 */}
          <div style={{ background:'rgba(239,159,39,0.08)', border:'1px solid rgba(239,159,39,0.25)', borderRadius:'14px', padding:'18px 20px' }}>
            <div style={{ marginBottom:'12px' }}>
              <span style={{ background:'rgba(239,159,39,0.2)', color:'#FAC775', fontSize:'11px', fontWeight:700, padding:'3px 12px', borderRadius:'100px', letterSpacing:'0.3px', display:'inline-block' }}>
                👑 super_admin
              </span>
            </div>
            <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.55)', lineHeight:'1.85' }}>
              <div>• 전체 사용자 조회 및 권한 변경</div>
              <div>• 시스템 LLM 설정 및 API 키 관리</div>
              <div>• AI 포인트 정책값 설정</div>
              <div>• 광고 슬롯 관리</div>
              <div>• 관리자 포인트 지급 이력 감사</div>
              <div>• DB 미저장 · 환경변수로만 존재</div>
            </div>
          </div>

          {/* admin — 초록 */}
          <div style={{ background:'rgba(82,183,136,0.08)', border:'1px solid rgba(82,183,136,0.25)', borderRadius:'14px', padding:'18px 20px' }}>
            <div style={{ marginBottom:'12px' }}>
              <span style={{ background:'rgba(82,183,136,0.2)', color:'#52B788', fontSize:'11px', fontWeight:700, padding:'3px 12px', borderRadius:'100px', letterSpacing:'0.3px', display:'inline-block' }}>
                🛡️ admin
              </span>
            </div>
            <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.55)', lineHeight:'1.85' }}>
              <div>• 학습자 검색 및 AI 포인트 수동 지급</div>
              <div>• 1회 최대 지급량: 슈퍼관리자 설정값</div>
              <div>• 지급 시 사유 메모 필수 입력</div>
              <div>• 소셜 로그인 + TOTP 2단계 인증 필수</div>
              <div>• 첫 로그인 후 Google Authenticator 설정</div>
              <div>• /admin/* 경로 접근 가능</div>
            </div>
          </div>

          {/* learner + creator — 노랑 + 주황 */}
          <div style={{ background:'rgba(255,255,204,0.06)', border:'1px solid rgba(255,220,80,0.22)', borderRadius:'14px', padding:'18px 20px' }}>
            <div style={{ display:'flex', alignItems:'center', gap:'6px', flexWrap:'wrap', marginBottom:'12px' }}>
              <span style={{ background:'rgba(255,220,80,0.18)', color:'#FFD84D', fontSize:'11px', fontWeight:700, padding:'3px 12px', borderRadius:'100px', letterSpacing:'0.3px', display:'inline-block' }}>
                📚 learner
              </span>
              <span style={{ background:'rgba(255,140,50,0.18)', color:'#FF8C32', fontSize:'11px', fontWeight:700, padding:'3px 12px', borderRadius:'100px', letterSpacing:'0.3px', display:'inline-block' }}>
                🎨 creator
              </span>
            </div>
            <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.55)', lineHeight:'1.85' }}>
              <div>• 소셜 로그인 (구글·카카오·네이버)</div>
              <div>• 가입 시 웰컴 AI 포인트 {welcomePoints}pt 지급</div>
              <div>• AI 코스 생성 ({courseGenCost}pt), 레슨 추천 ({lessonRecCost}pt)</div>
              <div>• Pro Access 사용자는 `pro_curriculum` 런타임 사용</div>
              <div>• BYOK 등록 시 포인트 소비 없음</div>
              <div style={{ color:'#FF8C32' }}>• 🎨 크리에이터: 학습자 역할 포함</div>
              <div style={{ color:'#FF8C32' }}>• 콘텐츠 등록·공개·포크 허용 (Phase 2)</div>
            </div>
          </div>

        </div>

        {/* Search */}
        <div className="mb-6 flex gap-3">
          <input
            type="text"
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
            placeholder="이메일 또는 닉네임 검색"
            className="flex-1 rounded-xl border border-[rgba(120,140,200,0.2)] bg-[#111E35] px-4 py-2.5 text-sm text-[#E8EAF2] outline-none focus:border-[#378ADD]"
          />
          <button
            onClick={handleSearch}
            className="rounded-xl bg-[#378ADD] px-5 py-2.5 text-sm font-medium text-white transition-opacity hover:opacity-90"
          >
            검색
          </button>
        </div>

        {error && (
          <div className="mb-4 rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-200">
            {error}
          </div>
        )}

        {/* Table */}
        <div className="overflow-hidden rounded-2xl border border-[rgba(120,140,200,0.2)]">
          <table className="w-full">
            <thead>
              <tr className="border-b border-[rgba(120,140,200,0.2)] bg-[#111E35]">
                <th className="px-4 py-3 text-left text-xs font-medium text-[rgba(200,210,235,0.5)]">계정</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-[rgba(200,210,235,0.5)]">닉네임</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-[rgba(200,210,235,0.5)]">역할</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-[rgba(200,210,235,0.5)]">가입일</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-[rgba(200,210,235,0.5)]">관리</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={5} className="py-12 text-center text-sm text-[rgba(200,210,235,0.5)]">불러오는 중...</td>
                </tr>
              ) : users.length === 0 ? (
                <tr>
                  <td colSpan={5} className="py-12 text-center text-sm text-[rgba(200,210,235,0.5)]">사용자가 없습니다.</td>
                </tr>
              ) : (
                users.map((user) => (
                  <tr key={user.id} className="border-b border-[rgba(120,140,200,0.1)] hover:bg-[#111E35]"
                    style={{ background: user.totp_reset_requested ? 'rgba(239,68,68,0.05)' : '#0B1629' }}
                  >
                    <td className="px-4 py-3">
                      <div style={{ display:'flex', flexDirection:'column', gap:'3px' }}>
                        <div style={{ display:'flex', alignItems:'center', gap:'6px', flexWrap:'wrap' }}>
                          <span style={{ fontSize:'15px' }}>{PROVIDER_ICON[user.provider] ?? '⚪'}</span>
                          <span style={{ fontSize:'13px', color: getIdentifierColor(user) }}>{getIdentifier(user)}</span>
                          {user.totp_reset_requested && (
                            <span style={{ fontSize:'10px', background:'rgba(239,68,68,0.2)', color:'#FC8181', padding:'2px 7px', borderRadius:'100px', fontWeight:600 }}>
                              🔴 TOTP 초기화 신청
                            </span>
                          )}
                        </div>
                        <span style={{ fontSize:'10px', color:'rgba(200,210,235,0.28)', fontFamily:'monospace', letterSpacing:'0.5px' }}>
                          ID: {user.id.slice(0, 8)}
                        </span>
                      </div>
                    </td>
                    <td className="px-4 py-3 text-sm text-[#E8EAF2]">{user.nickname || '-'}</td>
                    <td className="px-4 py-3">
                      <div style={{ display:'flex', alignItems:'center', gap:'6px', flexWrap:'wrap' }}>
                        <span style={getRoleBadgeStyle(user.role)}>{user.role}</span>
                        {user.premium_access && (
                          <span style={{ background:'rgba(239,159,39,0.16)', color:'#FAC775', padding:'3px 10px', borderRadius:'100px', fontSize:'11px', fontWeight:600, display:'inline-block', whiteSpace:'nowrap' }}>
                            pro access
                          </span>
                        )}
                      </div>
                    </td>
                    <td className="px-4 py-3 text-sm text-[rgba(200,210,235,0.5)]">{user.created_at}</td>
                    <td className="px-4 py-3">
                      <div style={{ display: 'flex', gap: '6px', alignItems: 'center', flexWrap: 'wrap' }}>
                        {user.role === 'learner' ? (
                          <button
                            onClick={() => handleRoleChange(user.id, 'admin')}
                            className="rounded-lg bg-[#378ADD]/20 px-3 py-1.5 text-xs font-medium text-[#378ADD] transition-colors hover:bg-[#378ADD]/30"
                          >
                            관리자 권한 부여
                          </button>
                        ) : user.role === 'admin' ? (
                          <button
                            onClick={() => handleRoleChange(user.id, 'learner')}
                            className="rounded-lg bg-red-500/20 px-3 py-1.5 text-xs font-medium text-red-400 transition-colors hover:bg-red-500/30"
                          >
                            권한 회수
                          </button>
                        ) : null}
                        {user.role === 'admin' && user.totp_reset_requested && (
                          <button
                            onClick={() => handleTotpReset(user.id)}
                            style={{ background:'rgba(239,68,68,0.15)', color:'#FC8181', border:'1px solid rgba(239,68,68,0.3)', borderRadius:'6px', padding:'4px 10px', fontSize:'11px', fontWeight:600, cursor:'pointer' }}
                          >
                            🔐 TOTP 초기화
                          </button>
                        )}
                        <button
                          onClick={() => setPointModal({
                            open: true,
                            userId: user.id,
                            userName: user.display_id ?? user.email ?? user.id.slice(0, 8),
                            amount: 10,
                            memo: '',
                          })}
                          style={{ background:'rgba(239,159,39,0.15)', color:'#EF9F27', border:'1px solid rgba(239,159,39,0.3)', borderRadius:'6px', padding:'4px 10px', fontSize:'11px', fontWeight:600, cursor:'pointer' }}
                        >
                          💎 포인트 지급
                        </button>
                        <button
                          onClick={() => handlePremiumAccessChange(user.id, !user.premium_access)}
                          style={{
                            background: user.premium_access ? 'rgba(239,159,39,0.15)' : 'rgba(55,138,221,0.12)',
                            color: user.premium_access ? '#EF9F27' : '#8CC7FF',
                            border: user.premium_access ? '1px solid rgba(239,159,39,0.3)' : '1px solid rgba(55,138,221,0.25)',
                            borderRadius:'6px',
                            padding:'4px 10px',
                            fontSize:'11px',
                            fontWeight:600,
                            cursor:'pointer',
                          }}
                        >
                          {user.premium_access ? 'Pro Access 해제' : 'Pro Access 부여'}
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="mt-4 flex items-center justify-between">
            <p className="text-sm text-[rgba(200,210,235,0.5)]">
              총 {total}명 · {page}/{totalPages} 페이지
            </p>
            <div className="flex gap-2">
              <button
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page === 1}
                className="rounded-xl border border-[rgba(120,140,200,0.2)] px-4 py-2 text-sm text-[rgba(200,210,235,0.6)] transition-colors hover:bg-[#111E35] disabled:opacity-30"
              >
                이전
              </button>
              <button
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={page === totalPages}
                className="rounded-xl border border-[rgba(120,140,200,0.2)] px-4 py-2 text-sm text-[rgba(200,210,235,0.6)] transition-colors hover:bg-[#111E35] disabled:opacity-30"
              >
                다음
              </button>
            </div>
          </div>
        )}
      </div>
    </main>
  )
}
