'use client';

import type { ReactNode } from 'react';
import type { PointPageThemeTokens } from '../../pointPageUtils';
import { pointShellTitleStyle } from '../../pointPageStyles';
import { pointTabButtonStyle } from '../../pointPageStyles';

interface ResearchSectionHeaderProps {
  title: string;
  action?: ReactNode;
  themeTokens: PointPageThemeTokens;
}

export function ResearchSectionHeader({ title, action, themeTokens }: ResearchSectionHeaderProps) {
  return (
    <div style={{ width: '100%', maxWidth: '860px', margin: '0 auto', display: 'flex', alignItems: 'center', gap: '10px' }}>
      <span
        aria-hidden="true"
        style={{
          width: '5px',
          height: '22px',
          borderRadius: '999px',
          background: themeTokens.pageBackground === '#F8FAFC' ? '#2563EB' : '#7BA7E8',
          boxShadow: themeTokens.pageBackground === '#F8FAFC' ? '0 8px 18px rgba(37, 99, 235, 0.16)' : '0 8px 18px rgba(123, 167, 232, 0.20)',
        }}
      />
      <strong style={{ ...pointShellTitleStyle, color: themeTokens.title, fontSize: '18px', letterSpacing: 0 }}>{title}</strong>
      {action ? <div style={{ marginLeft: 'auto', display: 'flex', alignItems: 'center' }}>{action}</div> : null}
    </div>
  );
}

interface WorkspaceDividerProps {
  themeTokens: PointPageThemeTokens;
}

export function WorkspaceDivider({ themeTokens }: WorkspaceDividerProps) {
  return (
    <div
      aria-hidden="true"
      style={{
        height: '1px',
        width: '100%',
        margin: '10px 0',
        background: `linear-gradient(90deg, transparent, ${themeTokens.sectionBorder}, transparent)`,
      }}
    />
  );
}

interface WorkspaceMiniNavItem {
  key: string;
  label: string;
  active?: boolean;
  checked?: boolean;
  disabled?: boolean;
  onClick: () => void;
}

interface WorkspaceMiniNavProps {
  ariaLabel: string;
  items: WorkspaceMiniNavItem[];
  sticky?: boolean;
  layout?: 'top' | 'side';
  variant?: 'default' | 'register';
  themeTokens: PointPageThemeTokens;
}

export function WorkspaceMiniNav({ ariaLabel, items, sticky = false, layout = 'top', variant = 'default', themeTokens }: WorkspaceMiniNavProps) {
  const isRegisterVariant = variant === 'register';
  const registerAccentStyles: Record<string, { background: string; borderColor: string; color: string; boxShadow: string }> = {
    video: { background: '#378ADD', borderColor: 'rgba(55, 138, 221, 0.76)', color: '#FFFDF7', boxShadow: '0 10px 20px rgba(55, 138, 221, 0.20)' },
    content: { background: '#EF9F27', borderColor: 'rgba(239, 159, 39, 0.76)', color: '#241505', boxShadow: '0 10px 20px rgba(239, 159, 39, 0.20)' },
    attachments: { background: '#1C7D79', borderColor: 'rgba(28, 125, 121, 0.76)', color: '#F8FFFC', boxShadow: '0 10px 20px rgba(28, 125, 121, 0.20)' },
  };
  return (
    <nav
      aria-label={ariaLabel}
      style={{
        width: '100%',
        maxWidth: layout === 'side' ? '112px' : '860px',
        margin: layout === 'side' ? '0' : '0 auto',
        display: isRegisterVariant ? 'grid' : 'flex',
        gridTemplateColumns: isRegisterVariant ? 'repeat(3, minmax(0, 1fr))' : undefined,
        flexDirection: layout === 'side' ? 'column' : 'row',
        flexWrap: layout === 'side' ? 'nowrap' : 'wrap',
        justifyContent: isRegisterVariant ? undefined : undefined,
        gap: '8px',
        padding: layout === 'side' ? '10px' : '10px 12px',
        boxSizing: 'border-box',
        borderRadius: '8px',
        border: isRegisterVariant ? 'none' : `1px solid ${themeTokens.surfaceBorder}`,
        background: isRegisterVariant ? 'transparent' : themeTokens.pageBackground === '#F8FAFC' ? 'rgba(241, 245, 249, 0.92)' : 'rgba(15, 23, 42, 0.58)',
        boxShadow: isRegisterVariant ? 'none' : themeTokens.pageBackground === '#F8FAFC' ? '0 8px 18px rgba(15, 23, 42, 0.06)' : '0 10px 24px rgba(0, 0, 0, 0.16)',
        alignSelf: layout === 'side' ? 'start' : 'stretch',
        ...(sticky ? {
          position: 'sticky',
          top: '430px',
          zIndex: 8,
          backdropFilter: 'blur(10px)',
        } : {}),
      }}
    >
      {items.map((item) => {
        const registerAccentStyle = isRegisterVariant ? registerAccentStyles[item.key] : null;
        return (
        <button
          key={item.key}
          type="button"
          onClick={item.onClick}
          disabled={item.disabled}
          style={{
            ...pointTabButtonStyle,
            minHeight: isRegisterVariant ? '48px' : '34px',
            padding: isRegisterVariant ? '8px 6px' : '6px 10px',
            fontSize: isRegisterVariant ? '13px' : '12px',
            borderRadius: isRegisterVariant ? '12px' : undefined,
            gap: '6px',
            width: layout === 'side' || isRegisterVariant ? '100%' : undefined,
            justifyContent: layout === 'side' ? 'flex-start' : 'center',
            background: registerAccentStyle ? registerAccentStyle.background : item.active ? themeTokens.primaryButtonBackground : themeTokens.secondaryButtonBackground,
            borderColor: registerAccentStyle ? registerAccentStyle.borderColor : item.active ? themeTokens.primaryButtonBorder : themeTokens.secondaryButtonBorder,
            color: registerAccentStyle ? registerAccentStyle.color : themeTokens.buttonText,
            cursor: item.disabled ? 'default' : 'pointer',
            opacity: item.disabled ? 0.56 : 1,
            boxShadow: registerAccentStyle ? registerAccentStyle.boxShadow : item.active ? '0 8px 18px rgba(37, 99, 235, 0.16)' : '0 4px 10px rgba(15, 23, 42, 0.06)',
          }}
        >
          {typeof item.checked === 'boolean' ? (
            <span aria-hidden="true" style={{ fontSize: '13px', lineHeight: 1 }}>
              {item.checked ? '☑' : '☐'}
            </span>
          ) : null}
          <strong style={{ whiteSpace: isRegisterVariant ? 'nowrap' : undefined }}>{item.label}</strong>
        </button>
        );
      })}
    </nav>
  );
}
