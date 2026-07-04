'use client';

import type { CSSProperties } from 'react';
import BrandLogo from '@/components/common/BrandLogo';

interface SuperAdminPanelHeaderProps {
  subtitle: string;
  description?: string;
}

export default function SuperAdminPanelHeader({
  subtitle,
  description,
}: SuperAdminPanelHeaderProps) {
  return (
    <section style={cardStyle}>
      <div style={logoRowStyle}>
        <BrandLogo href="/" iconSize={30} textSize="20px" />
        <span style={dividerStyle}>/</span>
        <span style={panelLabelStyle}>슈퍼관리자 패널</span>
      </div>
      <div style={contentStyle}>
        <h1 style={titleStyle}>{subtitle}</h1>
        {description ? <p style={descriptionStyle}>{description}</p> : null}
      </div>
    </section>
  );
}

const cardStyle: CSSProperties = {
  display: 'grid',
  gap: '12px',
  padding: '22px 24px',
  borderRadius: '24px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(7, 18, 34, 0.62)',
  boxShadow: '0 24px 64px rgba(0,0,0,0.22)',
};

const logoRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: '12px',
  flexWrap: 'wrap',
};

const dividerStyle: CSSProperties = {
  fontSize: '14px',
  color: 'rgba(160, 178, 214, 0.4)',
};

const panelLabelStyle: CSSProperties = {
  fontSize: '18px',
  fontWeight: 700,
  color: '#F7FAFF',
};

const contentStyle: CSSProperties = {
  display: 'grid',
  gap: '8px',
};

const titleStyle: CSSProperties = {
  margin: 0,
  fontSize: '28px',
  color: '#F7FAFF',
};

const descriptionStyle: CSSProperties = {
  margin: 0,
  fontSize: '14px',
  lineHeight: 1.7,
  color: '#C8D1E8',
};
