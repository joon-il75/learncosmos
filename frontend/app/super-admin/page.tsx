'use client';

import Link from 'next/link';
import SuperAdminPanelNav from '@/components/super-admin/SuperAdminPanelNav';
import SuperAdminPanelHeader from '@/components/super-admin/SuperAdminPanelHeader';
import { SUPER_ADMIN_PAGE_WIDTH } from '@/components/super-admin/layout';

const cards = [
  {
    title: '사용자 관리',
    copy: '학습자 계정, 역할, Pro 권한, TOTP 리셋, 포인트 지급을 관리합니다.',
    href: '/super-admin/users',
  },
  {
    title: '초대 코드',
    copy: '알파 테스트 접근 코드를 발급하고 발송 대상 메모와 사용 이력을 확인합니다.',
    href: '/super-admin/alpha-invite-codes',
  },
  {
    title: '시스템 설정',
    copy: 'LLM, API 키, 포인트 정책, 상품, 광고, 제휴, 정책 문서를 조정합니다.',
    href: '/super-admin/settings',
  },
  {
    title: '새소식 관리',
    copy: '플랫폼 새소식과 공지사항의 작성, 수정, 공개 상태 관리를 준비합니다.',
    href: '/super-admin/platform-news',
  },
  {
    title: '추천 점검',
    copy: '추천 품질을 직접 입력 또는 draft 기반으로 점검합니다.',
    href: '/super-admin/recommendation-debug',
  },
  {
    title: '오류 신고',
    copy: '학습탐험 지점 자료 오류신고를 확인하고 검토 상태와 처리 메모를 관리합니다.',
    href: '/super-admin/material-reports',
  },
  {
    title: 'Safety 로그',
    copy: '학습자 입력 moderation 결과를 원문 없이 hash, action, risk, route 기준으로 확인합니다.',
    href: '/super-admin/safety',
  },
  {
    title: '루미설정',
    copy: '루미 상태, 메시지, 스프라이트와 자산을 운영 기준으로 조정합니다.',
    href: '/super-admin/lumi-lab',
  },
  {
    title: '행성 텍스처맵',
    copy: '4x3 자전용 행성 텍스처맵을 등록 전 검토하고 변화 시점별 자전 프리뷰를 확인합니다.',
    href: '/super-admin/planet-texture-maps',
  },
  {
    title: '운영 로그',
    copy: '포인트 지급 이력을 중심으로 운영 변경 내역을 확인합니다.',
    href: '/super-admin/audit',
  },
] as const;

export default function SuperAdminHubPage() {
  return (
    <main style={pageStyle}>
      <div style={shellStyle}>
        <SuperAdminPanelHeader
          subtitle="운영 허브"
          description="현재 운영 업무를 `사용자 관리 / 시스템 설정 / 추천 점검 / 경험 설정 / 운영 로그`로 분리합니다. 이 화면은 전체 진입점과 빠른 이동 허브 역할을 합니다."
        />
        <SuperAdminPanelNav activeSection="hub" />

        <section style={gridStyle}>
          {cards.map((card) => (
            <Link key={card.href} href={card.href} style={cardStyle}>
              <strong style={cardTitleStyle}>{card.title}</strong>
              <p style={cardCopyStyle}>{card.copy}</p>
              <span style={cardActionStyle}>바로 이동</span>
            </Link>
          ))}
        </section>
      </div>
    </main>
  );
}

const pageStyle = {
  minHeight: '100vh',
  background: 'radial-gradient(circle at top, rgba(45, 78, 132, 0.2), transparent 42%), linear-gradient(180deg, #07111f, #0b1321 42%, #060b12)',
  color: '#eff6ff',
  padding: '28px 20px 48px',
};

const shellStyle = {
  width: SUPER_ADMIN_PAGE_WIDTH,
  margin: '0 auto',
  display: 'grid',
  gap: '18px',
};

const gridStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
  gap: '14px',
};

const cardStyle = {
  display: 'grid',
  gap: '10px',
  padding: '20px',
  borderRadius: '20px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(255,255,255,0.05)',
  textDecoration: 'none',
  color: 'inherit',
};

const cardTitleStyle = {
  fontSize: '18px',
  color: '#F7FAFF',
};

const cardCopyStyle = {
  margin: 0,
  fontSize: '13px',
  lineHeight: 1.7,
  color: '#B8C5E2',
};

const cardActionStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: '32px',
  width: 'fit-content',
  padding: '0 12px',
  borderRadius: '999px',
  background: 'rgba(95, 131, 255, 0.18)',
  border: '1px solid rgba(125, 160, 255, 0.34)',
  color: '#E8EEFF',
  fontSize: '12px',
  fontWeight: 700,
};
