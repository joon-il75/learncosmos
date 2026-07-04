'use client'

import { useEffect, useRef, useState } from 'react'
import { useRouter } from 'next/navigation'
import { QRCodeSVG } from 'qrcode.react'
import { authFetch } from '@/lib/auth/store'

const API_BASE = ''

export default function AdminTotpSetupPage() {
  const router = useRouter()
  const [qrURL, setQrURL] = useState('')
  const [secret, setSecret] = useState('')
  const [code, setCode] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [initializing, setInitializing] = useState(true)
  const codeRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    const setup = async () => {
      try {
        // refresh는 HttpOnly 쿠키만 갱신한다.
        const refreshRes = await fetch(`${API_BASE}/api/v1/auth/refresh`, {
          method: 'POST',
          credentials: 'include',
        })
        if (!refreshRes.ok) {
          router.replace('/login')
          return
        }

        // authFetch는 쿠키 기반 인증을 사용하고, 401 시 refresh 후 재시도한다.
        const res = await authFetch(`${API_BASE}/api/v1/admin/totp/setup`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
        })

        if (!res.ok) {
          setError('TOTP 설정을 시작할 수 없습니다. 다시 로그인해 주세요.')
          return
        }

        const data = await res.json()
        setQrURL(data.qr_url)
        setSecret(data.secret)
        codeRef.current?.focus()
      } catch {
        setError('서버 오류가 발생했습니다.')
      } finally {
        setInitializing(false)
      }
    }

    setup()
  }, [router])

  const handleVerify = async () => {
    if (code.length !== 6) return
    setError('')
    setLoading(true)

    try {
      const res = await authFetch(`${API_BASE}/api/v1/admin/totp/verify`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code }),
      })

      if (!res.ok) {
        setError('코드가 올바르지 않습니다. 다시 시도해주세요.')
        setCode('')
        codeRef.current?.focus()
        return
      }

      router.replace('/admin/dashboard')
    } catch {
      setError('서버 오류가 발생했습니다.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-[#0B1629] px-6">
      <div className="w-full max-w-sm">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-bold text-[#E8EAF2]">2단계 인증 설정</h1>
          <p className="mt-2 text-sm text-[rgba(200,210,235,0.5)]">LearnCosmos Admin Panel</p>
        </div>

        {initializing ? (
          <div className="text-center text-sm text-[rgba(200,210,235,0.5)]">초기화 중...</div>
        ) : error && !qrURL ? (
          <p className="text-center text-sm text-red-400">{error}</p>
        ) : (
          <div className="flex flex-col gap-6">
            {/* QR 코드 */}
            <div className="flex flex-col items-center gap-3">
              <p className="text-sm text-[rgba(200,210,235,0.7)]">
                Google Authenticator에서 QR을 스캔하세요
              </p>
              {qrURL && (
                <div className="rounded-xl bg-white p-3">
                  <QRCodeSVG value={qrURL} size={200} />
                </div>
              )}
            </div>

            {/* 평문 시크릿 (수동 입력용) */}
            {secret && (
              <div className="rounded-xl border border-[rgba(239,159,39,0.3)] bg-[#1a2540] p-4">
                <p className="mb-1.5 text-xs font-medium text-[#EF9F27]">
                  수동 입력 코드 (이 코드는 지금만 확인할 수 있습니다)
                </p>
                <p className="break-all font-mono text-sm tracking-wider text-[#E8EAF2]">
                  {secret}
                </p>
              </div>
            )}

            {/* 6자리 코드 입력 */}
            <div>
              <label className="mb-1.5 block text-xs font-medium text-[rgba(200,210,235,0.6)]">
                인증 코드 6자리 입력
              </label>
              <input
                ref={codeRef}
                type="text"
                inputMode="numeric"
                value={code}
                onChange={(e) => setCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                onKeyDown={(e) => e.key === 'Enter' && handleVerify()}
                maxLength={6}
                placeholder="000000"
                autoFocus
                className="w-full rounded-xl border border-[rgba(120,140,200,0.2)] bg-[#111E35] px-4 py-3 text-center text-xl font-semibold tracking-[0.5em] text-[#E8EAF2] outline-none focus:border-[#378ADD]"
              />
            </div>

            {error && <p className="text-center text-sm text-red-400">{error}</p>}

            <button
              type="button"
              onClick={handleVerify}
              disabled={loading || code.length !== 6}
              className="mt-2 rounded-xl bg-[#EF9F27] px-6 py-3 text-sm font-semibold text-[#0B1629] transition-opacity hover:opacity-90 disabled:opacity-50"
            >
              {loading ? '확인 중...' : '설정 완료'}
            </button>
          </div>
        )}
      </div>
    </main>
  )
}
