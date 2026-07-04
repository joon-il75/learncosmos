'use client';

import Link from 'next/link';
import type { CSSProperties } from 'react';

type ActiveSection = 'hub' | 'users' | 'alpha-invite' | 'demo-accounts' | 'settings' | 'platform-news' | 'recommendation' | 'material-reports' | 'safety' | 'experience' | 'audit';
type ExperienceSection = 'lumi' | 'planet-texture-maps' | null;

interface SuperAdminPanelNavProps {
  activeSection: ActiveSection;
  activeExperience?: ExperienceSection;
}

const primaryItems = [
  { key: 'hub', label: '운영 허브', href: '/super-admin' },
  { key: 'users', label: '사용자 관리', href: '/super-admin/users' },
  { key: 'alpha-invite', label: '초대 코드', href: '/super-admin/alpha-invite-codes' },
  { key: 'demo-accounts', label: '데모 계정', href: '/super-admin/demo-accounts' },
  { key: 'settings', label: '시스템 설정', href: '/super-admin/settings' },
  { key: 'platform-news', label: '새소식 관리', href: '/super-admin/platform-news' },
  { key: 'recommendation', label: '추천 점검', href: '/super-admin/recommendation-debug' },
  { key: 'material-reports', label: '오류 신고', href: '/super-admin/material-reports' },
  { key: 'safety', label: 'Safety 로그', href: '/super-admin/safety' },
  { key: 'experience', label: '경험 설정', href: '/super-admin/lumi-lab' },
  { key: 'audit', label: '운영 로그', href: '/super-admin/audit' },
] as const;

const experienceItems = [
  { key: 'lumi', label: '루미설정', href: '/super-admin/lumi-lab' },
  { key: 'planet-texture-maps', label: '행성 텍스처맵', href: '/super-admin/planet-texture-maps' },
] as const;

export default function SuperAdminPanelNav({
  activeSection,
  activeExperience = null,
}: SuperAdminPanelNavProps) {
  const handleLogout = () => {
    sessionStorage.removeItem('super_admin_token');
  };

  return (
    <div style={shellStyle}>
      <div style={primaryRowStyle}>
        {primaryItems.map((item) => (
          <Link
            key={item.key}
            href={item.href}
            style={navButtonStyle(item.key === activeSection)}
          >
            {item.label}
          </Link>
        ))}
        <Link href="/super-admin/login" onClick={handleLogout} style={utilityButtonStyle(false)}>
          로그아웃
        </Link>
      </div>

      {activeSection === 'experience' ? (
        <div style={secondaryRowStyle}>
          <span style={secondaryLabelStyle}>경험 설정</span>
          {experienceItems.map((item) => (
            <Link
              key={item.key}
              href={item.href}
              style={subNavButtonStyle(item.key === activeExperience)}
            >
              {item.label}
            </Link>
          ))}
        </div>
      ) : null}
    </div>
  );
}

const shellStyle: CSSProperties = {
  position: 'relative',
  zIndex: 80,
  isolation: 'isolate',
  pointerEvents: 'auto',
  display: 'grid',
  gap: '10px',
  marginTop: '12px',
  marginBottom: '16px',
};

const primaryRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: '8px',
  flexWrap: 'wrap',
};

const secondaryRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: '8px',
  flexWrap: 'wrap',
};

const navButtonStyle = (active: boolean): CSSProperties => ({
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: '36px',
  padding: '0 14px',
  borderRadius: '10px',
  textDecoration: 'none',
  fontSize: '13px',
  fontWeight: active ? 700 : 600,
  color: active ? '#F7FAFF' : 'rgba(200,210,235,0.82)',
  background: active ? 'rgba(95, 131, 255, 0.18)' : 'rgba(255,255,255,0.05)',
  border: active ? '1px solid rgba(125, 160, 255, 0.34)' : '1px solid rgba(120,140,200,0.2)',
  fontFamily: 'inherit',
  cursor: 'pointer',
  position: 'relative',
  zIndex: 1,
  pointerEvents: 'auto',
});

const utilityButtonStyle = (_planetMap: boolean): CSSProperties => ({
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: '36px',
  padding: '0 14px',
  borderRadius: '10px',
  textDecoration: 'none',
  fontSize: '13px',
  fontWeight: 600,
  color: 'rgba(200,210,235,0.82)',
  background: 'rgba(255,255,255,0.08)',
  border: '1px solid rgba(120,140,200,0.2)',
  fontFamily: 'inherit',
  cursor: 'pointer',
  position: 'relative',
  zIndex: 1,
  pointerEvents: 'auto',
});

const secondaryLabelStyle: CSSProperties = {
  fontSize: '12px',
  fontWeight: 700,
  color: 'rgba(160,178,214,0.78)',
  letterSpacing: '0.06em',
  textTransform: 'uppercase',
  marginRight: '2px',
};

const subNavButtonStyle = (active: boolean): CSSProperties => ({
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: '32px',
  padding: '0 12px',
  borderRadius: '999px',
  textDecoration: 'none',
  fontSize: '12px',
  fontWeight: active ? 700 : 600,
  color: active ? '#DFFFE8' : 'rgba(200,210,235,0.8)',
  background: active ? 'rgba(50, 200, 100, 0.16)' : 'rgba(255,255,255,0.04)',
  border: active ? '1px solid rgba(50, 200, 100, 0.34)' : '1px solid rgba(120,140,200,0.16)',
  fontFamily: 'inherit',
  cursor: 'pointer',
  position: 'relative',
  zIndex: 1,
  pointerEvents: 'auto',
});
