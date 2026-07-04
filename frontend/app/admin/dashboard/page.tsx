'use client'

import { useEffect, useRef, useState, useCallback } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { authFetch, logout } from '@/lib/auth/store'
import AppHeaderShell, { appHeaderActionButtonStyle, appHeaderActionLinkStyle } from '@/components/common/AppHeaderShell'

const API_BASE = ''

interface UserInfo {
  id: string
  email: string
  nickname: string
  role: string
}

interface AdminUser {
  id: string
  email: string
  nickname: string
  role: string
  points: number
  created_at: string
}

type GrantStep = 'edit' | 'confirm'

type Stage = 'loading' | 'not-setup' | 'totp-gate' | 'reset-requested' | 'dashboard'

export default function AdminDashboardPage() {
  const router = useRouter()
  const [stage, setStage] = useState<Stage>('loading')
  const [user, setUser] = useState<UserInfo | null>(null)
  const [totpCode, setTotpCode] = useState('')
  const [totpError, setTotpError] = useState('')
  const [validating, setValidating] = useState(false)
  const [requesting, setRequesting] = useState(false)
  const codeRef = useRef<HTMLInputElement>(null)
  const hasInit = useRef(false)

  // ── 사용자 관리 상태 ──
  const [adminUsers, setAdminUsers] = useState<AdminUser[]>([])
  const [adminTotal, setAdminTotal] = useState(0)
  const [adminPage, setAdminPage] = useState(1)
  const [searchInput, setSearchInput] = useState('')
  const [searchQuery, setSearchQuery] = useState('')
  const [usersLoading, setUsersLoading] = useState(false)

  // ── 포인트 정책 ──
  const [adminMaxGrant, setAdminMaxGrant] = useState(50)

  // ── 포인트 지급 모달 상태 ──
  const [grantModal, setGrantModal] = useState<{ userId: string; userName: string } | null>(null)
  const [grantAmount, setGrantAmount] = useState(10)
  const [grantMemo, setGrantMemo] = useState('')
  const [grantLoading, setGrantLoading] = useState(false)
  const [grantError, setGrantError] = useState('')
  const [grantStep, setGrantStep] = useState<GrantStep>('edit')

  const LIMIT = 20

  const fetchUsers = useCallback(async (q: string, p: number) => {
    setUsersLoading(true)
    try {
      const params = new URLSearchParams({ q, page: String(p), limit: String(LIMIT) })
      const res = await authFetch(`${API_BASE}/api/v1/admin/users?${params}`)
      if (!res.ok) return
      const data = await res.json()
      setAdminUsers(data.users ?? [])
      setAdminTotal(data.total ?? 0)
    } finally {
      setUsersLoading(false)
    }
  }, [])

  useEffect(() => {
    if (hasInit.current) return
    hasInit.current = true

    const init = async () => {
      try {
        const refreshRes = await fetch(`${API_BASE}/api/v1/auth/refresh`, {
          method: 'POST',
          credentials: 'include',
        })
        if (!refreshRes.ok) { router.replace('/login'); return }

        const meRes = await fetch(`${API_BASE}/api/v1/auth/me`, {
          credentials: 'include',
          cache: 'no-store',
        })
        if (!meRes.ok) { router.replace('/login'); return }

        const userData: UserInfo = await meRes.json()
        if (userData.role !== 'admin') { router.replace('/dashboard'); return }
        setUser(userData)

        const statusRes = await fetch(`${API_BASE}/api/v1/admin/totp/status`, {
          credentials: 'include',
          cache: 'no-store',
        })
        if (!statusRes.ok) { router.replace('/login'); return }
        const status = await statusRes.json()

        if (!status.totp_enabled) {
          setStage(status.reset_requested ? 'reset-requested' : 'not-setup')
          return
        }

        const totpSessionRes = await fetch(`/api/v1/admin/point-policy`, {
          credentials: 'include',
          cache: 'no-store',
        })
        if (totpSessionRes.ok) {
          setStage('dashboard')
          return
        }

        setStage('totp-gate')
        setTimeout(() => codeRef.current?.focus(), 100)
      } catch {
        router.replace('/login')
      }
    }
    init()
  }, [router])

  useEffect(() => {
    if (stage === 'dashboard') {
      fetchUsers(searchQuery, adminPage)
      authFetch(`${API_BASE}/api/v1/admin/point-policy`)
        .then(r => r.ok ? r.json() : null)
        .then(d => { if (d?.admin_max_grant) setAdminMaxGrant(d.admin_max_grant) })
        .catch(() => {})
    }
  }, [stage, searchQuery, adminPage, fetchUsers])

  const handleTotpValidate = async () => {
    if (totpCode.length !== 6 || validating) return
    setValidating(true)
    setTotpError('')
    try {
      const res = await authFetch(`${API_BASE}/api/v1/admin/totp/validate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code: totpCode }),
      })
      if (res.ok) {
        setStage('dashboard')
        return
      }
      const data = await res.json().catch(() => ({}))
      if (data.error === 'totp_not_initialized') { setStage('not-setup'); return }
      setTotpError('코드가 올바르지 않습니다. 다시 시도해주세요.')
      setTotpCode('')
      setTimeout(() => codeRef.current?.focus(), 50)
    } catch {
      setTotpError('서버 오류가 발생했습니다.')
    } finally {
      setValidating(false)
    }
  }

  const handleResetRequest = async () => {
    if (requesting) return
    setRequesting(true)
    try {
      await authFetch(`${API_BASE}/api/v1/admin/totp/reset-request`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      })
      setStage('reset-requested')
    } catch {
      alert('오류가 발생했습니다.')
    } finally {
      setRequesting(false)
    }
  }

  const handleSearch = () => {
    setAdminPage(1)
    setSearchQuery(searchInput)
  }

  const openGrantModal = (u: AdminUser) => {
    setGrantModal({ userId: u.id, userName: u.nickname || u.email || u.id.slice(0, 8) })
    setGrantAmount(10)
    setGrantMemo('')
    setGrantError('')
    setGrantStep('edit')
  }

  const handleGrant = async () => {
    if (!grantModal || grantLoading) return
    if (grantAmount === 0 || grantAmount < -adminMaxGrant || grantAmount > adminMaxGrant) {
      setGrantError(`-${adminMaxGrant}~${adminMaxGrant}pt 범위로 입력해주세요. 0은 불가합니다.`)
      return
    }
    setGrantLoading(true)
    setGrantError('')
    try {
      const res = await authFetch(`${API_BASE}/api/v1/admin/users/${grantModal.userId}/grant-points`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ amount: grantAmount, memo: grantMemo }),
      })
      if (!res.ok) {
        const d = await res.json().catch(() => ({}))
        setGrantError(d.error || '지급에 실패했습니다.')
        return
      }
      setGrantModal(null)
      setGrantStep('edit')
      fetchUsers(searchQuery, adminPage)
    } catch {
      setGrantError('서버 오류가 발생했습니다.')
    } finally {
      setGrantLoading(false)
    }
  }

  // ── 로딩 ──
  if (stage === 'loading') {
    return (
      <main className="flex min-h-screen items-center justify-center bg-[#0B1629]">
        <p className="text-sm text-[rgba(200,210,235,0.5)]">불러오는 중...</p>
      </main>
    )
  }

  // ── TOTP 미설정 ──
  if (stage === 'not-setup') {
    return (
      <main className="flex min-h-screen items-center justify-center bg-[#0B1629] px-6">
        <div className="w-full max-w-sm text-center">
          <div className="mb-6 text-4xl">🔐</div>
          <h1 className="mb-2 text-xl font-bold text-[#E8EAF2]">2단계 인증 설정 필요</h1>
          <p className="mb-8 text-sm text-[rgba(200,210,235,0.5)]">
            관리자 패널 접근을 위해 Google Authenticator 등록이 필요합니다.
          </p>
          <a
            href="/admin/totp-setup"
            className="block rounded-xl bg-[#EF9F27] px-6 py-3 text-sm font-semibold text-[#0B1629] transition-opacity hover:opacity-90"
          >
            기기 등록하러 가기
          </a>
          <button
            onClick={() => logout()}
            className="mt-4 text-xs text-[rgba(200,210,235,0.3)] hover:text-[rgba(200,210,235,0.6)]"
          >
            로그아웃
          </button>
        </div>
      </main>
    )
  }

  // ── TOTP 게이트 ──
  if (stage === 'totp-gate') {
    return (
      <main className="flex min-h-screen items-center justify-center bg-[#0B1629] px-6">
        <div className="w-full max-w-sm">
          <div className="mb-8 text-center">
            <h1 className="text-2xl font-bold text-[#E8EAF2]">2단계 인증</h1>
            <p className="mt-2 text-sm text-[rgba(200,210,235,0.5)]">{user?.email}</p>
          </div>
          <div className="flex flex-col gap-4">
            <div>
              <label className="mb-1.5 block text-xs font-medium text-[rgba(200,210,235,0.6)]">
                인증 앱의 6자리 코드를 입력하세요
              </label>
              <input
                ref={codeRef}
                type="text"
                inputMode="numeric"
                value={totpCode}
                onChange={(e) => { setTotpCode(e.target.value.replace(/\D/g, '').slice(0, 6)); setTotpError('') }}
                onKeyDown={(e) => e.key === 'Enter' && handleTotpValidate()}
                maxLength={6}
                placeholder="000000"
                className="w-full rounded-xl border border-[rgba(120,140,200,0.2)] bg-[#111E35] px-4 py-3 text-center text-xl font-semibold tracking-[0.5em] text-[#E8EAF2] outline-none focus:border-[#378ADD]"
              />
            </div>
            {totpError && <p className="text-center text-sm text-red-400">{totpError}</p>}
            <button
              onClick={handleTotpValidate}
              disabled={totpCode.length !== 6 || validating}
              className="rounded-xl bg-[#378ADD] px-6 py-3 text-sm font-semibold text-white transition-opacity hover:opacity-90 disabled:opacity-40"
            >
              {validating ? '확인 중...' : '확인'}
            </button>
            <div className="mt-2 border-t border-[rgba(120,140,200,0.1)] pt-4 text-center">
              <p className="mb-2 text-xs text-[rgba(200,210,235,0.35)]">기기를 분실하셨나요?</p>
              <button
                onClick={handleResetRequest}
                disabled={requesting}
                className="text-xs text-[rgba(239,159,39,0.7)] hover:text-[#EF9F27] disabled:opacity-40"
              >
                {requesting ? '신청 중...' : '슈퍼관리자에게 초기화 신청'}
              </button>
            </div>
            <button onClick={() => logout()} className="text-center text-xs text-[rgba(200,210,235,0.25)] hover:text-[rgba(200,210,235,0.5)]">
              로그아웃
            </button>
          </div>
        </div>
      </main>
    )
  }

  // ── 초기화 신청 완료 ──
  if (stage === 'reset-requested') {
    return (
      <main className="flex min-h-screen items-center justify-center bg-[#0B1629] px-6">
        <div className="w-full max-w-sm text-center">
          <div className="mb-6 text-4xl">⏳</div>
          <h1 className="mb-2 text-xl font-bold text-[#E8EAF2]">초기화 신청 완료</h1>
          <p className="mb-2 text-sm text-[rgba(200,210,235,0.5)]">
            슈퍼관리자가 초기화 처리 후 알림을 드립니다.
          </p>
          <p className="mb-8 text-sm text-[rgba(200,210,235,0.35)]">
            처리 완료 후 이 페이지에서 기기를 재등록하실 수 있습니다.
          </p>
          <button
            onClick={() => window.location.reload()}
            className="block w-full rounded-xl border border-[rgba(120,140,200,0.2)] px-6 py-3 text-sm text-[rgba(200,210,235,0.6)] transition-colors hover:bg-[#111E35]"
          >
            새로고침
          </button>
          <button onClick={() => logout()} className="mt-4 text-xs text-[rgba(200,210,235,0.25)] hover:text-[rgba(200,210,235,0.5)]">
            로그아웃
          </button>
        </div>
      </main>
    )
  }

  // ── 대시보드 ──
  const totalPages = Math.ceil(adminTotal / LIMIT)
  const grantActionLabel = grantAmount < 0 ? '차감' : '지급'
  const absGrantAmount = Math.abs(grantAmount)

  return (
    <main className="min-h-screen bg-[#0B1629]">
      {/* 포인트 지급 모달 */}
      {grantModal && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4"
          onClick={(e) => { if (e.target === e.currentTarget) setGrantModal(null) }}
        >
          <div className="w-full max-w-sm rounded-2xl border border-[rgba(120,140,200,0.2)] bg-[#111E35] p-6">
            <h2 className="mb-1 text-base font-bold text-[#E8EAF2]">포인트 지급</h2>
            <p className="mb-5 text-xs text-[rgba(200,210,235,0.5)]">{grantModal.userName}</p>

            {grantStep === 'edit' ? (
              <>
                <label className="mb-1.5 block text-xs font-medium text-[rgba(200,210,235,0.6)]">
                  지급/차감량 (최대 {adminMaxGrant}pt)
                </label>
                <div className="mb-4 flex items-center gap-2">
                  <input
                    type="number"
                    min={-adminMaxGrant}
                    max={adminMaxGrant}
                    value={grantAmount}
                    onChange={(e) => {
                      const value = Number(e.target.value)
                      if (Number.isNaN(value)) { setGrantAmount(0); return }
                      setGrantAmount(Math.min(adminMaxGrant, Math.max(-adminMaxGrant, value)))
                    }}
                    className="w-24 rounded-xl border border-[rgba(120,140,200,0.2)] bg-[#0B1629] px-3 py-2 text-center text-sm font-semibold text-[#E8EAF2] outline-none focus:border-[#378ADD]"
                  />
                  <span className="text-sm text-[rgba(200,210,235,0.5)]">pt</span>
                  <div className="ml-auto flex gap-1">
                    {[-10, 10, 20, adminMaxGrant].filter((v, i, arr) => arr.indexOf(v) === i && Math.abs(v) <= adminMaxGrant).map(v => (
                      <button
                        key={v}
                        onClick={() => setGrantAmount(v)}
                        className={`rounded-lg px-2.5 py-1 text-xs font-medium transition-colors ${
                          grantAmount === v
                            ? 'bg-[#378ADD] text-white'
                            : 'border border-[rgba(120,140,200,0.2)] text-[rgba(200,210,235,0.6)] hover:bg-[#1a2540]'
                        }`}
                      >
                        {v > 0 ? `+${v}` : `${v}`}
                      </button>
                    ))}
                  </div>
                </div>

                <label className="mb-1.5 block text-xs font-medium text-[rgba(200,210,235,0.6)]">
                  메모 (선택)
                </label>
                <input
                  type="text"
                  value={grantMemo}
                  onChange={(e) => setGrantMemo(e.target.value)}
                  placeholder="지급/차감 사유"
                  maxLength={100}
                  className="mb-4 w-full rounded-xl border border-[rgba(120,140,200,0.2)] bg-[#0B1629] px-3 py-2 text-sm text-[#E8EAF2] outline-none focus:border-[#378ADD] placeholder:text-[rgba(200,210,235,0.25)]"
                />
              </>
            ) : (
              <div className="mb-4 rounded-xl border border-[rgba(120,140,200,0.18)] bg-[#0B1629] px-4 py-4">
                <p className="mb-2 text-sm font-semibold text-[#E8EAF2]">
                  {grantModal.userName}에게 {absGrantAmount}pt를 {grantActionLabel}합니다.
                </p>
                <p className="text-xs text-[rgba(200,210,235,0.5)]">
                  메모: {grantMemo || '없음'}
                </p>
              </div>
            )}

            {grantError && <p className="mb-3 text-xs text-red-400">{grantError}</p>}

            <div className="flex gap-2">
              <button
                onClick={() => {
                  if (grantStep === 'confirm') {
                    setGrantStep('edit')
                    return
                  }
                  setGrantModal(null)
                }}
                className="flex-1 rounded-xl border border-[rgba(120,140,200,0.2)] py-2.5 text-sm text-[rgba(200,210,235,0.6)] transition-colors hover:bg-[#1a2540]"
              >
                {grantStep === 'confirm' ? '수정' : '취소'}
              </button>
              <button
                onClick={() => {
                  if (grantStep === 'edit') {
                    if (grantAmount === 0 || grantAmount < -adminMaxGrant || grantAmount > adminMaxGrant) {
                      setGrantError(`-${adminMaxGrant}~${adminMaxGrant}pt 범위로 입력해주세요. 0은 불가합니다.`)
                      return
                    }
                    setGrantError('')
                    setGrantStep('confirm')
                    return
                  }
                  handleGrant()
                }}
                disabled={grantLoading}
                className="flex-1 rounded-xl bg-[#EF9F27] py-2.5 text-sm font-semibold text-[#0B1629] transition-opacity hover:opacity-90 disabled:opacity-40"
              >
                {grantLoading ? '처리 중...' : grantStep === 'edit' ? '확인' : `${absGrantAmount}pt ${grantActionLabel} 확정`}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* 헤더 */}
      <AppHeaderShell
        logoHref="/"
        logoIconSize={28}
        logoTextSize="16px"
        position="sticky"
        maxWidth="80rem"
        headerStyle={{
          zIndex: 40,
          borderBottom: '1px solid rgba(120,140,200,0.1)',
          background: 'rgba(11,22,41,0.95)',
        }}
        innerStyle={{ maxWidth: '80rem', padding: '16px 24px' }}
        logoWrapperStyle={{ opacity: 0.8 }}
        leftMeta={
          <>
            <span className="text-[rgba(120,140,200,0.3)] text-sm">/</span>
            <div>
              <h1 className="text-sm font-bold text-[#E8EAF2]">관리자 패널</h1>
              <p className="text-xs text-[rgba(200,210,235,0.45)]">{user?.nickname || user?.email}</p>
            </div>
          </>
        }
        rightSlot={
          <>
            <Link
              href="/dashboard"
              style={appHeaderActionLinkStyle(true)}
            >
              Planet Map
            </Link>
            <button
              onClick={() => logout()}
              style={appHeaderActionButtonStyle()}
            >
              로그아웃
            </button>
          </>
        }
      />

      <div className="mx-auto max-w-5xl px-6 py-8">
        {/* 검색 + 통계 */}
        <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h2 className="text-base font-semibold text-[#E8EAF2]">학습자 목록</h2>
            <p className="mt-0.5 text-xs text-[rgba(200,210,235,0.4)]">총 {adminTotal}명</p>
          </div>
          <div className="flex gap-2">
            <input
              type="text"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
              placeholder="닉네임 또는 이메일 검색"
              className="w-56 rounded-xl border border-[rgba(120,140,200,0.2)] bg-[#111E35] px-3 py-2 text-sm text-[#E8EAF2] outline-none focus:border-[#378ADD] placeholder:text-[rgba(200,210,235,0.25)]"
            />
            <button
              onClick={handleSearch}
              className="rounded-xl bg-[#378ADD] px-4 py-2 text-sm font-medium text-white transition-opacity hover:opacity-90"
            >
              검색
            </button>
          </div>
        </div>

        {/* 테이블 */}
        <div className="overflow-hidden rounded-2xl border border-[rgba(120,140,200,0.15)] bg-[#111E35]">
          {/* 테이블 헤더 */}
          <div className="grid grid-cols-[2fr_2fr_1fr_1fr_1fr] border-b border-[rgba(120,140,200,0.1)] px-5 py-3">
            {['닉네임', '이메일', '역할', '포인트', '가입일'].map(h => (
              <span key={h} className="text-xs font-semibold text-[rgba(200,210,235,0.4)] uppercase tracking-wide">{h}</span>
            ))}
          </div>

          {/* 테이블 바디 */}
          {usersLoading ? (
            <div className="py-16 text-center text-sm text-[rgba(200,210,235,0.4)]">불러오는 중...</div>
          ) : adminUsers.length === 0 ? (
            <div className="py-16 text-center text-sm text-[rgba(200,210,235,0.4)]">사용자가 없습니다.</div>
          ) : (
            adminUsers.map((u) => (
              <div
                key={u.id}
                className="grid grid-cols-[2fr_2fr_1fr_1fr_1fr] items-center border-b border-[rgba(120,140,200,0.07)] px-5 py-3.5 transition-colors last:border-b-0 hover:bg-[#1a2540]"
              >
                <span className="truncate text-sm font-medium text-[#E8EAF2]">
                  {u.nickname || <span className="text-[rgba(200,210,235,0.35)]">—</span>}
                </span>
                <span className="truncate text-sm text-[rgba(200,210,235,0.6)]">
                  {u.email || <span className="text-[rgba(200,210,235,0.3)]">—</span>}
                </span>
                <span>
                  <span className="inline-block rounded-full bg-[rgba(255,216,77,0.15)] px-2.5 py-0.5 text-xs font-medium text-[#FFD84D]">
                    {u.role}
                  </span>
                </span>
                <span className="text-sm font-semibold text-[#EF9F27]">{u.points}pt</span>
                <div className="flex items-center gap-3">
                  <span className="text-xs text-[rgba(200,210,235,0.4)]">{u.created_at}</span>
                  <button
                    onClick={() => openGrantModal(u)}
                    className="rounded-lg border border-[rgba(239,159,39,0.3)] px-2.5 py-1 text-xs font-medium text-[#EF9F27] transition-colors hover:bg-[rgba(239,159,39,0.1)]"
                  >
                    지급
                  </button>
                </div>
              </div>
            ))
          )}
        </div>

        {/* 페이지네이션 */}
        {totalPages > 1 && (
          <div className="mt-6 flex items-center justify-center gap-2">
            <button
              onClick={() => setAdminPage(p => Math.max(1, p - 1))}
              disabled={adminPage === 1}
              className="rounded-lg border border-[rgba(120,140,200,0.2)] px-3 py-1.5 text-xs text-[rgba(200,210,235,0.6)] transition-colors hover:bg-[#111E35] disabled:opacity-30"
            >
              이전
            </button>
            <span className="text-xs text-[rgba(200,210,235,0.5)]">{adminPage} / {totalPages}</span>
            <button
              onClick={() => setAdminPage(p => Math.min(totalPages, p + 1))}
              disabled={adminPage === totalPages}
              className="rounded-lg border border-[rgba(120,140,200,0.2)] px-3 py-1.5 text-xs text-[rgba(200,210,235,0.6)] transition-colors hover:bg-[#111E35] disabled:opacity-30"
            >
              다음
            </button>
          </div>
        )}
      </div>
    </main>
  )
}
