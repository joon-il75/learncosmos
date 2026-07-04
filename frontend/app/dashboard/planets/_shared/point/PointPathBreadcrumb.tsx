'use client';

import { useState } from 'react';
import Link from 'next/link';
import type { PointPageTheme, getPointPageThemeTokens } from '../pointPageUtils';

export type PointPathItem = {
  key: string;
  label: string;
  href: string;
  title: string;
};

export type CurrentPointPathItem = {
  key: string;
  label: string;
  title: string;
};

export function PointPathBreadcrumb({
  items,
  currentItem,
  fullText,
  theme,
  themeTokens,
}: {
  items: PointPathItem[];
  currentItem: CurrentPointPathItem;
  fullText: string;
  theme: PointPageTheme;
  themeTokens: ReturnType<typeof getPointPageThemeTokens>;
  onCurrentClick: () => void;
}) {
  const [isOpen, setIsOpen] = useState(false);
  const visibleItems = items.filter((item) => item.key !== 'system');

  return (
    <div
      style={{
        flex: '1 1 auto',
        minWidth: 0,
        display: 'grid',
        gap: isOpen ? '8px' : '0px',
        color: themeTokens.mutedText,
        lineHeight: 1.4,
      }}
      title={fullText}
    >
      <button
        type="button"
        title={isOpen ? currentItem.title : fullText}
        aria-expanded={isOpen}
        onClick={() => setIsOpen((open) => !open)}
        style={{
          width: 'fit-content',
          maxWidth: '100%',
          display: 'inline-flex',
          alignItems: 'flex-start',
          gap: '8px',
          minHeight: '38px',
          padding: '8px 10px',
          borderRadius: '6px',
          border: `1px solid ${theme === 'light' ? 'rgba(107, 151, 216, 0.58)' : 'rgba(148, 163, 184, 0.58)'}`,
          background: theme === 'light'
            ? 'rgba(238, 246, 255, 0.88)'
            : 'rgba(30, 41, 59, 0.78)',
          color: themeTokens.buttonText,
          font: 'inherit',
          fontSize: '16px',
          lineHeight: 1.2,
          fontWeight: 900,
          textAlign: 'left',
          cursor: 'pointer',
        }}
      >
        <span style={{ minWidth: 0, whiteSpace: 'normal', overflowWrap: 'anywhere', wordBreak: 'break-word', lineHeight: 1.35 }}>{currentItem.label}</span>
        <span aria-hidden="true" style={{ flex: '0 0 auto', paddingTop: '2px', color: themeTokens.mutedText, fontSize: '13px', fontWeight: 900 }}>{isOpen ? '▲' : '▼'}</span>
      </button>

      {isOpen ? (
        <div
          role="tree"
          aria-label="현재 지점 경로"
          style={{
            display: 'grid',
            gap: '2px',
            padding: '4px 0 2px',
            color: themeTokens.mutedText,
            fontSize: '13px',
            fontWeight: 780,
          }}
        >
          {visibleItems.map((item, index) => (
            <div key={item.key} role="treeitem" style={{ display: 'flex', minWidth: 0, paddingLeft: `${index * 18}px` }}>
              {index > 0 ? <span aria-hidden="true" style={{ flex: '0 0 auto', marginRight: '6px', color: themeTokens.mutedText }}>└</span> : null}
              <Link
                href={item.href}
                title={item.title}
                style={{
                  minWidth: 0,
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  whiteSpace: 'nowrap',
                  color: themeTokens.mutedText,
                  textDecoration: 'none',
                }}
              >
                {item.label}
              </Link>
            </div>
          ))}
          <div role="treeitem" aria-current="page" style={{ display: 'flex', alignItems: 'flex-start', minWidth: 0, paddingLeft: `${visibleItems.length * 18}px`, color: themeTokens.buttonText, fontWeight: 900 }}>
            {visibleItems.length > 0 ? <span aria-hidden="true" style={{ flex: '0 0 auto', marginRight: '6px', color: themeTokens.mutedText }}>└</span> : null}
            <span style={{ minWidth: 0, whiteSpace: 'normal', overflowWrap: 'anywhere', wordBreak: 'break-word', lineHeight: 1.45 }}>{currentItem.label}</span>
          </div>
        </div>
      ) : null}
    </div>
  );
}
