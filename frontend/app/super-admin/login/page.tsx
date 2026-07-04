'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import BrandLogo from '@/components/common/BrandLogo';

const API_BASE = '';

export default function SuperAdminLoginPage() {
  const router = useRouter();
  const [id, setId] = useState('');
  const [pw, setPw] = useState('');
  const [totp, setTotp] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const res = await fetch(`${API_BASE}/api/v1/super-admin/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id, password: pw, totp }),
      });
      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        setError(data.error || data.message || '로그인에 실패했습니다.');
      } else {
		await res.json();
        router.replace('/super-admin');
      }
    } catch {
      setError('서버 연결에 실패했습니다.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <main style={{
      minHeight: '100vh',
      background: '#081826',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '24px',
      position: 'relative',
      overflow: 'hidden',
      fontFamily: "'Noto Sans KR', -apple-system, sans-serif",
    }}>

      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Noto+Sans+KR:wght@400;500;600;700&display=swap');

        @keyframes bf1 { 0%,100%{transform:translate(0,0)} 50%{transform:translate(-24px,18px)} }
        @keyframes bf2 { 0%,100%{transform:translate(0,0)} 50%{transform:translate(20px,-22px)} }
        @keyframes bf3 { 0%,100%{transform:translate(0,0)} 50%{transform:translate(14px,20px)} }
        @keyframes cardIn { from{opacity:0;transform:translateY(24px)} to{opacity:1;transform:translateY(0)} }
        @keyframes pulse { 0%,100%{opacity:1} 50%{opacity:0.35} }
        @keyframes spin { from{transform:rotate(0deg)} to{transform:rotate(360deg)} }

        .sa-blob { position:absolute; border-radius:50%; filter:blur(80px); pointer-events:none; }
        .sa-b1 { width:420px; height:420px; background:radial-gradient(circle,rgba(55,138,221,0.13) 0%,transparent 70%); top:-120px; right:-80px; animation:bf1 14s ease-in-out infinite; }
        .sa-b2 { width:320px; height:320px; background:radial-gradient(circle,rgba(127,119,221,0.11) 0%,transparent 70%); bottom:-60px; left:-60px; animation:bf2 18s ease-in-out infinite; }
        .sa-b3 { width:220px; height:220px; background:radial-gradient(circle,rgba(239,159,39,0.08) 0%,transparent 70%); top:38%; left:8%; animation:bf3 22s ease-in-out infinite; }

        .sa-card { animation: cardIn 0.65s cubic-bezier(.22,1,.36,1) forwards; }

        .sa-input {
          width:100%; background:rgba(255,255,255,0.04);
          border:1px solid rgba(120,140,200,0.2); border-radius:12px;
          padding:13px 14px 13px 44px; font-size:14px;
          color:#E8EAF2; outline:none; transition:all 0.2s;
          font-family:inherit; box-sizing:border-box;
        }
        .sa-input::placeholder { color:rgba(200,210,235,0.3); font-size:13px; }
        .sa-input:focus {
          border-color:rgba(55,138,221,0.5);
          background:rgba(55,138,221,0.06);
          box-shadow:0 0 0 3px rgba(55,138,221,0.08);
        }
        .sa-totp {
          text-align:center; letter-spacing:10px; font-size:20px; font-weight:600;
          padding:13px 14px; background:rgba(239,159,39,0.06);
          border:1px solid rgba(239,159,39,0.2); border-radius:12px;
          color:#E8EAF2; outline:none; transition:all 0.2s;
          width:100%; font-family:inherit; box-sizing:border-box;
        }
        .sa-totp::placeholder { letter-spacing:6px; font-size:16px; color:rgba(200,210,235,0.25); }
        .sa-totp:focus {
          border-color:rgba(239,159,39,0.5);
          background:rgba(239,159,39,0.09);
          box-shadow:0 0 0 3px rgba(239,159,39,0.08);
        }
        .sa-btn {
          width:100%; padding:14px;
          background:linear-gradient(135deg,#378ADD,#7F77DD);
          color:#fff; border:none; border-radius:12px;
          font-size:15px; font-weight:700; cursor:pointer;
          font-family:inherit; position:relative; overflow:hidden;
          transition:transform 0.2s, box-shadow 0.2s, opacity 0.2s;
          display:flex; align-items:center; justify-content:center; gap:8px;
          margin-top:8px;
        }
        .sa-btn:hover:not(:disabled) {
          transform:translateY(-2px);
          box-shadow:0 12px 32px rgba(55,138,221,0.35),0 4px 12px rgba(127,119,221,0.2);
        }
        .sa-btn:active:not(:disabled) { transform:translateY(0); }
        .sa-btn:disabled { opacity:0.6; cursor:not-allowed; }
        .sa-pulse { animation:pulse 2s infinite; }

        .sa-logo-text {
          font-size:20px; font-weight:700; letter-spacing:-0.5px;
          background:linear-gradient(135deg,#378ADD 0%,#7F77DD 50%,#EF9F27 100%);
          -webkit-background-clip:text; -webkit-text-fill-color:transparent;
          background-clip:text;
        }
      `}</style>

      {/* blob 배경 */}
      <div className="sa-blob sa-b1" />
      <div className="sa-blob sa-b2" />
      <div className="sa-blob sa-b3" />

      {/* 카드 */}
      <div className="sa-card" style={{
        width: '100%', maxWidth: '420px',
        padding: '48px 40px 40px',
        background: 'rgba(17,30,53,0.88)',
        backdropFilter: 'blur(24px)',
        WebkitBackdropFilter: 'blur(24px)',
        borderRadius: '24px',
        border: '1px solid rgba(120,140,200,0.22)',
        boxShadow: '0 24px 64px rgba(0,0,0,0.5), 0 0 0 1px rgba(55,138,221,0.06)',
        position: 'relative', zIndex: 10,
      }}>

        {/* 로고 */}
        <div style={{ display:'flex', alignItems:'center', justifyContent:'center', marginBottom:'10px' }}>
          <BrandLogo href="/" iconSize={38} textSize="34px" textClassName="sa-logo-text" />
        </div>

        {/* SUPER ADMIN 뱃지 */}
        <div style={{ display:'flex', justifyContent:'center', marginBottom:'28px' }}>
          <div style={{
            display:'inline-flex', alignItems:'center', gap:'5px',
            background:'rgba(239,159,39,0.12)', border:'1px solid rgba(239,159,39,0.25)',
            color:'#FAC775', fontSize:'11px', fontWeight:600,
            padding:'4px 14px', borderRadius:'100px', letterSpacing:'0.5px',
          }}>
            🔐 SUPER ADMIN PANEL
          </div>
        </div>

        {/* 타이틀 */}
        <div style={{ textAlign:'center', marginBottom:'32px' }}>
          <h1 style={{ fontSize:'20px', fontWeight:700, color:'#E8EAF2', marginBottom:'4px', letterSpacing:'-0.3px' }}>
            슈퍼관리자 로그인
          </h1>
          <p style={{ fontSize:'13px', color:'rgba(200,210,235,0.55)' }}>
            관리자 인증 정보를 입력해주세요
          </p>
        </div>

        <form onSubmit={handleLogin} style={{ display:'flex', flexDirection:'column', gap:'16px' }}>

          {/* 관리자 ID */}
          <div>
            <label style={{ display:'block', fontSize:'12px', fontWeight:600, color:'rgba(200,210,235,0.65)', marginBottom:'8px', letterSpacing:'0.3px' }}>
              관리자 ID
            </label>
            <div style={{ position:'relative' }}>
              <span style={{ position:'absolute', left:'14px', top:'50%', transform:'translateY(-50%)', opacity:0.4, pointerEvents:'none' }}>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#E8EAF2" strokeWidth="2" strokeLinecap="round">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                  <circle cx="12" cy="7" r="4"/>
                </svg>
              </span>
              <input
                type="text"
                className="sa-input"
                placeholder="관리자 ID 입력"
                value={id}
                onChange={e => setId(e.target.value)}
                required
                autoComplete="username"
              />
            </div>
          </div>

          {/* 비밀번호 */}
          <div>
            <label style={{ display:'block', fontSize:'12px', fontWeight:600, color:'rgba(200,210,235,0.65)', marginBottom:'8px', letterSpacing:'0.3px' }}>
              비밀번호
            </label>
            <div style={{ position:'relative' }}>
              <span style={{ position:'absolute', left:'14px', top:'50%', transform:'translateY(-50%)', opacity:0.4, pointerEvents:'none' }}>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#E8EAF2" strokeWidth="2" strokeLinecap="round">
                  <rect x="3" y="11" width="18" height="11" rx="2"/>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                </svg>
              </span>
              <input
                type="password"
                className="sa-input"
                placeholder="비밀번호 입력"
                value={pw}
                onChange={e => setPw(e.target.value)}
                required
                autoComplete="current-password"
              />
            </div>
          </div>

          {/* 구분선 */}
          <div style={{ height:'1px', background:'rgba(120,140,200,0.15)' }} />

          {/* TOTP */}
          <div>
            <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', marginBottom:'8px' }}>
              <label style={{ fontSize:'12px', fontWeight:600, color:'rgba(200,210,235,0.65)', letterSpacing:'0.3px' }}>
                TOTP 인증번호
              </label>
              <span style={{ fontSize:'11px', color:'rgba(250,199,117,0.8)' }}>Google Authenticator</span>
            </div>
            <input
              type="text"
              className="sa-totp"
              placeholder="000000"
              maxLength={6}
              value={totp}
              onChange={e => setTotp(e.target.value.replace(/\D/g, ''))}
              required
              autoComplete="one-time-code"
              inputMode="numeric"
            />
          </div>

          {/* 에러 메시지 */}
          {error && (
            <div style={{
              background:'rgba(226,75,74,0.1)', border:'1px solid rgba(226,75,74,0.25)',
              borderRadius:'10px', padding:'10px 14px',
              fontSize:'13px', color:'#F09595',
              display:'flex', alignItems:'center', gap:'8px',
            }}>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                <circle cx="12" cy="12" r="10"/>
                <line x1="12" y1="8" x2="12" y2="12"/>
                <line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
              {error}
            </div>
          )}

          {/* 로그인 버튼 */}
          <button type="submit" className="sa-btn" disabled={loading}>
            {loading ? (
              <>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="2.5" strokeLinecap="round" style={{ animation:'spin 1s linear infinite' }}>
                  <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
                </svg>
                로그인 중...
              </>
            ) : (
              <>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="2.5" strokeLinecap="round">
                  <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/>
                  <polyline points="10 17 15 12 10 7"/>
                  <line x1="15" y1="12" x2="3" y2="12"/>
                </svg>
                로그인
              </>
            )}
          </button>

        </form>


      </div>
    </main>
  );
}
