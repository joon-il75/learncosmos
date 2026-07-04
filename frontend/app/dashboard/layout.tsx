'use client';

import type { ReactNode } from 'react';
import { useEffect, useState } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import { getDashboardGateCopy } from '@/lib/i18n/pages/dashboardGate';

type GateState = 'checking' | 'allowed';

interface AlphaStatus {
  granted: boolean;
  requires_code: boolean;
  role: string;
}

interface DashboardGateUser {
  ui_locale?: 'ko' | 'en';
  language_setup_required?: boolean;
  required_consent_pending?: boolean;
}

export default function DashboardLayout({ children }: { children: ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const [state, setState] = useState<GateState>('checking');
  const [loadingLocale, setLoadingLocale] = useState<'ko' | 'en'>('ko');

  useEffect(() => {
    let cancelled = false;

    const checkAccess = async () => {
      const redirectAfter = pathname || '/dashboard';
      const refreshRes = await fetch('/api/v1/auth/refresh', {
        method: 'POST',
        credentials: 'include',
      });
      if (!refreshRes.ok) {
        router.replace(`/login?redirect_after=${encodeURIComponent(redirectAfter)}`);
        return;
      }

      const meRes = await fetch('/api/v1/auth/me', {
        credentials: 'include',
        cache: 'no-store',
      });
      if (!meRes.ok) {
        router.replace(`/login?redirect_after=${encodeURIComponent(redirectAfter)}`);
        return;
      }

      const meData = (await meRes.json()) as DashboardGateUser;
      if (!cancelled) setLoadingLocale(meData.ui_locale === 'en' ? 'en' : 'ko');
      if (meData.language_setup_required) {
        const setupPath = meData.ui_locale === 'en' ? '/en/language-setup' : '/language-setup';
        router.replace(`${setupPath}?redirect_after=${encodeURIComponent(redirectAfter)}`);
        return;
      }
      if (meData.required_consent_pending) {
        const agreementsPath = meData.ui_locale === 'en' ? '/en/agreements' : '/agreements';
        router.replace(`${agreementsPath}?redirect_after=${encodeURIComponent(redirectAfter)}`);
        return;
      }

      if (pathname === '/dashboard/settings' || pathname === '/dashboard/settings/profile') {
        if (!cancelled) setState('allowed');
        return;
      }

      const statusRes = await fetch('/api/v1/alpha-access/status', {
        credentials: 'include',
        cache: 'no-store',
      });
      if (!statusRes.ok) {
        router.replace(`/login?redirect_after=${encodeURIComponent(redirectAfter)}`);
        return;
      }

      const statusData = (await statusRes.json()) as AlphaStatus;
      if (!statusData.granted && statusData.requires_code) {
        const alphaInvitePath = meData.ui_locale === 'en' ? '/en/alpha-invite' : '/alpha-invite';
        router.replace(`${alphaInvitePath}?redirect_after=${encodeURIComponent(redirectAfter)}`);
        return;
      }

      if (!cancelled) setState('allowed');
    };

    setState('checking');
    checkAccess().catch(() => {
      if (!cancelled) router.replace(`/login?redirect_after=${encodeURIComponent(pathname || '/dashboard')}`);
    });

    return () => {
      cancelled = true;
    };
  }, [pathname, router]);

  if (state === 'checking') {
    const copy = getDashboardGateCopy(loadingLocale);
    return (
      <main style={loadingPageStyle}>
        <div style={loadingCardStyle}>
          <strong>{copy.checkingTitle}</strong>
          <span>{copy.checkingBody}</span>
        </div>
      </main>
    );
  }

  return <>{children}</>;
}

const loadingPageStyle = {
  minHeight: '100vh',
  display: 'grid',
  placeItems: 'center',
  padding: '24px',
  background: '#0B1629',
  color: '#E8EAF2',
};

const loadingCardStyle = {
  display: 'grid',
  gap: '8px',
  width: 'min(420px, 100%)',
  padding: '24px',
  borderRadius: '18px',
  border: '1px solid rgba(160,186,224,0.18)',
  background: 'rgba(11,22,41,0.9)',
  boxShadow: '0 24px 80px rgba(3,8,20,0.36)',
  textAlign: 'center' as const,
  fontSize: '14px',
};
