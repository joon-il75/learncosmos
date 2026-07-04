'use client';

import { Suspense, useEffect, useMemo, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import AppHeaderShell from '@/components/common/AppHeaderShell';
import type { Locale } from '@/lib/i18n/locales';

interface AlphaStatus {
  granted: boolean;
  requires_code: boolean;
  role: string;
}

type ApiErrorPayload = {
  error?: string;
  error_code?: string;
}

type AlphaInviteCopy = {
  tag: string;
  eyebrow: string;
  title: string;
  description: string;
  loading: string;
  statusError: string;
  emptyError: string;
  submitError: string;
  invalidRequestError: string;
  notFoundError: string;
  revokedError: string;
  expiredError: string;
  usedUpError: string;
  label: string;
  placeholder: string;
  submitting: string;
  submit: string;
  profileLink: string;
  fallback: string;
}

const alphaInviteCopy: Record<Locale, AlphaInviteCopy> = {
  ko: {
    tag: 'Alpha Test',
    eyebrow: 'LearnCosmos Alpha',
    title: '알파 테스트 초대 코드가 필요합니다',
    description: '가입과 약관 동의를 마친 뒤, 운영자가 전달한 초대 코드를 입력하면 학습탐험 기능에 접근할 수 있습니다.',
    loading: '접근 상태를 확인하는 중입니다.',
    statusError: '접근 상태를 확인하지 못했습니다. 잠시 후 다시 시도해 주세요.',
    emptyError: '초대 코드를 입력해 주세요.',
    submitError: '초대 코드 확인에 실패했습니다.',
    invalidRequestError: '잘못된 요청입니다.',
    notFoundError: '초대 코드를 찾을 수 없습니다.',
    revokedError: '회수된 초대 코드입니다.',
    expiredError: '입력 기간이 지난 초대 코드입니다.',
    usedUpError: '이미 사용이 완료된 초대 코드입니다.',
    label: '초대 코드',
    placeholder: 'LW-XXXX-XXXX-XXXX',
    submitting: '확인 중...',
    submit: '코드 확인',
    profileLink: '회원정보 / 계정 탈퇴',
    fallback: '초대 코드 화면을 준비하는 중입니다.',
  },
  en: {
    tag: 'Alpha Test',
    eyebrow: 'LearnCosmos Alpha',
    title: 'Alpha invite code required',
    description: 'After sign-up and required policy consent, enter the invite code from the operator to access learning features.',
    loading: 'Checking access status.',
    statusError: 'Could not check access status. Please try again later.',
    emptyError: 'Enter your invite code.',
    submitError: 'Could not verify the invite code.',
    invalidRequestError: 'Invalid request.',
    notFoundError: 'Could not find this invite code.',
    revokedError: 'This invite code has been revoked.',
    expiredError: 'This invite code has expired.',
    usedUpError: 'This invite code has already been used.',
    label: 'Invite code',
    placeholder: 'LW-XXXX-XXXX-XXXX',
    submitting: 'Checking...',
    submit: 'Verify code',
    profileLink: 'Profile / Delete Account',
    fallback: 'Preparing the invite code page.',
  },
};

function getAlphaInviteErrorMessage(payload: ApiErrorPayload, copy: AlphaInviteCopy) {
  if (payload.error_code === 'invalid_request') return copy.invalidRequestError;
  if (payload.error_code === 'alpha_code_not_found') return copy.notFoundError;
  if (payload.error_code === 'alpha_code_revoked') return copy.revokedError;
  if (payload.error_code === 'alpha_code_expired') return copy.expiredError;
  if (payload.error_code === 'alpha_code_used_up') return copy.usedUpError;
  if (payload.error_code === 'alpha_redeem_failed') return copy.submitError;
  return payload.error || copy.submitError;
}

function AlphaInviteContent({ locale }: { locale: Locale }) {
  const copy = alphaInviteCopy[locale];
  const router = useRouter();
  const searchParams = useSearchParams();
  const currentPath = locale === 'en' ? '/en/alpha-invite' : '/alpha-invite';
  const loginPath = locale === 'en' ? '/en/login' : '/login';
  const redirectAfter = normalizeRedirectAfter(searchParams.get('redirect_after'));
  const [code, setCode] = useState('');
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const normalizedCode = useMemo(() => code.trim().toUpperCase(), [code]);

  useEffect(() => {
    let cancelled = false;

    const load = async () => {
      const refreshRes = await fetch('/api/v1/auth/refresh', {
        method: 'POST',
        credentials: 'include',
      });
      if (!refreshRes.ok) {
        router.replace(`${loginPath}?redirect_after=${encodeURIComponent(currentPath)}`);
        return;
      }

      const meRes = await fetch('/api/v1/auth/me', {
        credentials: 'include',
        cache: 'no-store',
      });
      if (!meRes.ok) {
        router.replace(`${loginPath}?redirect_after=${encodeURIComponent(currentPath)}`);
        return;
      }
      const meData = await meRes.json();
      if (meData.required_consent_pending) {
        const agreementsPath = locale === 'en' ? '/en/agreements' : '/agreements';
        router.replace(`${agreementsPath}?redirect_after=${encodeURIComponent(currentPath + '?redirect_after=' + redirectAfter)}`);
        return;
      }

      const statusRes = await fetch('/api/v1/alpha-access/status', {
        credentials: 'include',
        cache: 'no-store',
      });
      if (!statusRes.ok) {
        router.replace(`${loginPath}?redirect_after=${encodeURIComponent(currentPath)}`);
        return;
      }
      const status = (await statusRes.json()) as AlphaStatus;
      if (status.granted) {
        router.replace(status.role === 'super_admin' ? '/super-admin' : redirectAfter);
        return;
      }
      if (!cancelled) setLoading(false);
    };

    load().catch(() => {
      if (!cancelled) {
        setError(copy.statusError);
        setLoading(false);
      }
    });

    return () => {
      cancelled = true;
    };
  }, [copy.statusError, currentPath, loginPath, redirectAfter, router]);

  const handleSubmit = async () => {
    if (!normalizedCode) {
      setError(copy.emptyError);
      return;
    }
    setSubmitting(true);
    setError('');
    try {
      const res = await fetch('/api/v1/alpha-access/redeem', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ code: normalizedCode }),
      });
      const payload = (await res.json().catch(() => ({}))) as ApiErrorPayload;
      if (!res.ok) {
        throw new Error(getAlphaInviteErrorMessage(payload, copy));
      }
      router.replace(redirectAfter);
    } catch (err) {
      setError(err instanceof Error ? err.message : copy.submitError);
      setSubmitting(false);
    }
  };

  return (
    <main style={pageStyle}>
      <AppHeaderShell
        position="sticky"
        maxWidth="980px"
        headerStyle={{ background: 'rgba(11,22,41,0.86)' }}
        rightSlot={<span style={tagStyle}>{copy.tag}</span>}
      />

      <section style={cardStyle}>
        <p style={eyebrowStyle}>{copy.eyebrow}</p>
        <h1 style={titleStyle}>{copy.title}</h1>
        <p style={copyStyle}>{copy.description}</p>

        {loading ? (
          <div style={noticeStyle}>{copy.loading}</div>
        ) : (
          <div style={formStyle}>
            <label style={labelStyle} htmlFor="alpha-code">
              {copy.label}
            </label>
            <input
              id="alpha-code"
              value={code}
              onChange={(event) => setCode(event.target.value)}
              placeholder={copy.placeholder}
              autoCapitalize="characters"
              style={inputStyle}
              onKeyDown={(event) => {
                if (event.key === 'Enter') handleSubmit();
              }}
            />
            {error ? <p style={errorStyle}>{error}</p> : null}
            <button
              type="button"
              onClick={handleSubmit}
              disabled={submitting}
              style={{
                ...submitButtonStyle,
                opacity: submitting ? 0.62 : 1,
                cursor: submitting ? 'not-allowed' : 'pointer',
              }}
            >
              {submitting ? copy.submitting : copy.submit}
            </button>
            <a href="/dashboard/settings/profile" style={profileLinkStyle}>
              {copy.profileLink}
            </a>
          </div>
        )}
      </section>
    </main>
  );
}

export default function AlphaInvitePageClient({ locale }: { locale: Locale }) {
  const copy = alphaInviteCopy[locale];

  return (
    <Suspense
      fallback={
        <main style={pageStyle}>
          <div style={cardStyle}>{copy.fallback}</div>
        </main>
      }
    >
      <AlphaInviteContent locale={locale} />
    </Suspense>
  );
}

function normalizeRedirectAfter(value: string | null) {
  if (!value || !value.startsWith('/') || value.startsWith('//')) return '/dashboard';
  return value;
}

const pageStyle = {
  minHeight: '100vh',
  background: 'radial-gradient(circle at 20% 0%, rgba(72, 121, 255, 0.24), transparent 34%), linear-gradient(180deg, #07111f, #0B1629)',
  color: '#E8EAF2',
  padding: '32px 20px 56px',
};

const cardStyle = {
  width: 'min(560px, 100%)',
  margin: '64px auto 0',
  display: 'grid',
  gap: '18px',
  padding: '28px',
  borderRadius: '22px',
  border: '1px solid rgba(160,186,224,0.18)',
  background: 'rgba(9, 20, 38, 0.84)',
  boxShadow: '0 30px 100px rgba(3, 8, 20, 0.42)',
};

const eyebrowStyle = {
  margin: 0,
  fontSize: '12px',
  fontWeight: 800,
  letterSpacing: '0.08em',
  color: '#8FB4FF',
  textTransform: 'uppercase' as const,
};

const titleStyle = {
  margin: 0,
  fontSize: 'clamp(26px, 5vw, 38px)',
  lineHeight: 1.16,
  letterSpacing: 0,
};

const copyStyle = {
  margin: 0,
  fontSize: '14px',
  lineHeight: 1.8,
  color: 'rgba(220,228,245,0.74)',
};

const formStyle = {
  display: 'grid',
  gap: '12px',
};

const labelStyle = {
  fontSize: '13px',
  fontWeight: 800,
  color: '#DCE6FF',
};

const inputStyle = {
  width: '100%',
  minHeight: '48px',
  borderRadius: '12px',
  border: '1px solid rgba(160,186,224,0.24)',
  background: 'rgba(255,255,255,0.06)',
  color: '#F7FAFF',
  padding: '0 14px',
  fontSize: '16px',
  fontWeight: 800,
  letterSpacing: '0.04em',
  outline: 'none',
};

const submitButtonStyle = {
  minHeight: '46px',
  borderRadius: '12px',
  border: '1px solid rgba(124, 166, 255, 0.42)',
  background: 'linear-gradient(180deg, rgba(92, 135, 255, 0.94), rgba(58, 99, 220, 0.94))',
  color: '#fff',
  fontSize: '14px',
  fontWeight: 800,
  fontFamily: 'inherit',
};

const profileLinkStyle = {
  minHeight: '42px',
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  borderRadius: '12px',
  border: '1px solid rgba(160,186,224,0.24)',
  background: 'rgba(255,255,255,0.06)',
  color: '#DCE6FF',
  fontSize: '13px',
  fontWeight: 800,
  textDecoration: 'none',
};

const errorStyle = {
  margin: 0,
  fontSize: '13px',
  lineHeight: 1.6,
  color: '#FFB4B4',
};

const noticeStyle = {
  borderRadius: '14px',
  border: '1px solid rgba(160,186,224,0.16)',
  background: 'rgba(255,255,255,0.06)',
  padding: '14px 16px',
  fontSize: '14px',
  color: 'rgba(220,228,245,0.78)',
};

const tagStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: '34px',
  padding: '0 12px',
  borderRadius: '999px',
  border: '1px solid rgba(143,180,255,0.28)',
  background: 'rgba(143,180,255,0.1)',
  color: '#CFE0FF',
  fontSize: '12px',
  fontWeight: 800,
};
