'use client';

import { useEffect, useRef } from 'react';
import type { PointPageTheme, getPointPageThemeTokens } from '../pointPageUtils';
import { PointPathBreadcrumb, type CurrentPointPathItem, type PointPathItem } from './PointPathBreadcrumb';

type ThemeTokens = ReturnType<typeof getPointPageThemeTokens>;

export function PointLearningToolbar({
  theme,
  themeTokens,
  pointPathItems,
  currentPointPathItem,
  fullPointPathText,
  onCurrentPathClick,
  onHeightChange,
}: {
  theme: PointPageTheme;
  themeTokens: ThemeTokens;
  pointPathItems: PointPathItem[];
  currentPointPathItem: CurrentPointPathItem;
  fullPointPathText: string;
  onCurrentPathClick: () => void;
  onHeightChange?: (height: number) => void;
}) {
  const toolbarRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const toolbarElement = toolbarRef.current;
    if (!toolbarElement || !onHeightChange) return;

    const notifyHeight = () => {
      onHeightChange(Math.ceil(toolbarElement.getBoundingClientRect().height));
    };
    notifyHeight();

    if (typeof ResizeObserver === 'undefined') {
      window.addEventListener('resize', notifyHeight);
      return () => window.removeEventListener('resize', notifyHeight);
    }

    const observer = new ResizeObserver(notifyHeight);
    observer.observe(toolbarElement);
    return () => observer.disconnect();
  }, [onHeightChange]);

  return (
    <div
      ref={toolbarRef}
      className="lw-point-learning-toolbar"
      style={{
        position: 'fixed',
        top: '100px',
        left: 0,
        right: 0,
        zIndex: 49,
        borderBottom: `1px solid ${themeTokens.navBorder}`,
        background: theme === 'light'
          ? 'rgba(248, 250, 252, 0.96)'
          : 'rgba(11, 15, 20, 0.92)',
        backdropFilter: 'blur(12px)',
        WebkitBackdropFilter: 'blur(12px)',
      }}
    >
      <div
        style={{
          width: '100%',
          maxWidth: '1360px',
          margin: '0 auto',
          minWidth: 0,
          display: 'flex',
          alignItems: 'center',
          padding: '10px 24px 12px',
        }}
      >
        <PointPathBreadcrumb
          items={pointPathItems}
          currentItem={currentPointPathItem}
          fullText={fullPointPathText}
          theme={theme}
          themeTokens={themeTokens}
          onCurrentClick={onCurrentPathClick}
        />
      </div>
    </div>
  );
}
