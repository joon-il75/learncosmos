'use client';

import LumiAvatar from '@/components/lumi/LumiAvatar';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageTheme, getPointPageThemeTokens } from '../pointPageUtils';
import { secondaryButtonStyle } from '../pointPageStyles';

type ThemeTokens = ReturnType<typeof getPointPageThemeTokens>;

export function PointLumiGuideBox({
  copy,
  theme,
  themeTokens,
  pointGoalContextText,
  currentFlowGuideMessage,
  lumiGuideTitle,
  lumiGuideMessage,
  isRestNoticeOpen,
  onToggleLumiGuide,
  onConfirmRestNotice,
}: {
  copy: PointLearningCopy['lumiGuide'];
  theme: PointPageTheme;
  themeTokens: ThemeTokens;
  pointGoalContextText: string;
  currentFlowGuideMessage: string;
  lumiGuideTitle: string;
  lumiGuideMessage: string;
  isRestNoticeOpen: boolean;
  onToggleLumiGuide: () => void;
  onConfirmRestNotice: () => void;
}) {
  return (
    <>
      <div
        aria-hidden="true"
        style={{
          height: '0px',
          width: '100%',
          background: theme === 'light'
            ? 'linear-gradient(90deg, rgba(148, 163, 184, 0.12), rgba(20, 184, 166, 0.34), rgba(148, 163, 184, 0.12))'
            : 'linear-gradient(90deg, rgba(51, 65, 85, 0.28), rgba(94, 234, 212, 0.32), rgba(51, 65, 85, 0.28))',
        }}
      />
      <div
        role="status"
        aria-live="polite"
        style={{
          display: 'grid',
          gap: '8px',
          padding: '10px 12px',
          borderRadius: '8px',
          border: `1px solid ${theme === 'light' ? 'rgba(20, 184, 166, 0.24)' : 'rgba(94, 234, 212, 0.20)'}`,
          background: theme === 'light'
            ? 'rgba(240, 253, 250, 0.94)'
            : 'rgba(15, 23, 42, 0.92)',
          boxShadow: theme === 'light'
            ? '0 16px 38px rgba(15, 23, 42, 0.18), inset 0 1px 0 rgba(255, 255, 255, 0.78)'
            : '0 18px 42px rgba(0, 0, 0, 0.38), inset 0 1px 0 rgba(255, 255, 255, 0.08)',
          backdropFilter: 'blur(14px)',
          WebkitBackdropFilter: 'blur(14px)',
          maxHeight: isRestNoticeOpen ? 'calc(100vh - 170px)' : '42vh',
          overflowY: 'auto',
        }}
      >
        <PointGoalLine
          theme={theme}
          themeTokens={themeTokens}
          text={pointGoalContextText}
          copy={copy}
          onToggleLumiGuide={onToggleLumiGuide}
        />
        <div style={{ display: 'grid', gridTemplateColumns: 'auto minmax(0, 1fr)', gap: '10px', alignItems: 'center' }}>
          <LumiAvatar state="exploring" size={34} reducedMotion={false} />
          <div style={{ display: 'grid', gap: '6px', minWidth: 0 }}>
            <CurrentTaskBox
              theme={theme}
              themeTokens={themeTokens}
              label={copy.currentTaskLabel}
              message={currentFlowGuideMessage}
            />
            <div
              aria-hidden="true"
              style={{
                height: '1px',
                width: '100%',
                margin: '2px 0',
                background: theme === 'light'
                  ? 'linear-gradient(90deg, rgba(20, 184, 166, 0.20), rgba(148, 163, 184, 0.18))'
                  : 'linear-gradient(90deg, rgba(94, 234, 212, 0.22), rgba(51, 65, 85, 0.34))',
              }}
            />
            <strong style={{ fontSize: '12px', lineHeight: 1.35, color: themeTokens.title }}>
              {copy.featureGuideLabel} · {lumiGuideTitle}
            </strong>
            <span style={{ fontSize: '13px', lineHeight: 1.42, color: themeTokens.mutedText }}>
              {lumiGuideMessage}
            </span>
          </div>
        </div>
        {isRestNoticeOpen ? (
          <RestNotice
            theme={theme}
            copy={copy}
            onConfirm={onConfirmRestNotice}
          />
        ) : null}
      </div>
    </>
  );
}

function PointGoalLine({
  theme,
  themeTokens,
  text,
  copy,
  onToggleLumiGuide,
}: {
  theme: PointPageTheme;
  themeTokens: ThemeTokens;
  text: string;
  copy: PointLearningCopy['lumiGuide'];
  onToggleLumiGuide: () => void;
}) {
  return (
    <div
      title={text}
      style={{
        minWidth: 0,
        display: 'grid',
        gridTemplateColumns: 'auto minmax(0, 1fr) auto',
        alignItems: 'center',
        gap: '8px',
        padding: '3px 8px',
        borderRadius: '4px',
        background: theme === 'light' ? 'rgba(255, 255, 255, 0.56)' : 'rgba(30, 41, 59, 0.58)',
        border: `1px solid ${theme === 'light' ? 'rgba(20, 184, 166, 0.18)' : 'rgba(94, 234, 212, 0.14)'}`,
      }}
    >
      <span style={{ flex: '0 0 auto', fontSize: '13px', fontWeight: 900, color: themeTokens.title }}>
        🎯 {copy.goalLabel}
      </span>
      <span
        style={{
          minWidth: 0,
          overflow: 'hidden',
          textOverflow: 'ellipsis',
          whiteSpace: 'nowrap',
          fontSize: '13px',
          lineHeight: 1.35,
          fontWeight: 650,
          color: themeTokens.mutedText,
        }}
      >
        {text}
      </span>
      <button
        type="button"
        aria-label={copy.close}
        title={copy.close}
        onClick={onToggleLumiGuide}
        style={{
          minHeight: '26px',
          display: 'inline-flex',
          alignItems: 'center',
          justifyContent: 'center',
          padding: '0 9px',
          borderRadius: '999px',
          border: `1px solid ${theme === 'light' ? 'rgba(20, 184, 166, 0.26)' : 'rgba(94, 234, 212, 0.22)'}`,
          background: theme === 'light' ? 'rgba(255, 255, 255, 0.72)' : 'rgba(15, 23, 42, 0.58)',
          color: themeTokens.title,
          fontSize: '12px',
          lineHeight: 1,
          fontWeight: 850,
          whiteSpace: 'nowrap',
          fontFamily: 'inherit',
        }}
      >
        {copy.closeButton}
      </button>
    </div>
  );
}

function CurrentTaskBox({
  theme,
  themeTokens,
  label,
  message,
}: {
  theme: PointPageTheme;
  themeTokens: ThemeTokens;
  label: string;
  message: string;
}) {
  return (
    <div
      style={{
        display: 'grid',
        gap: '3px',
        padding: '7px 9px',
        borderRadius: '5px',
        border: `1px solid ${theme === 'light' ? 'rgba(217, 119, 6, 0.34)' : 'rgba(251, 191, 36, 0.28)'}`,
        background: theme === 'light'
          ? 'linear-gradient(180deg, rgba(255, 251, 235, 0.98), rgba(254, 243, 199, 0.84))'
          : 'linear-gradient(180deg, rgba(69, 26, 3, 0.48), rgba(30, 41, 59, 0.72))',
        boxShadow: theme === 'light'
          ? 'inset 1px 1px 0 rgba(255,255,255,0.88), 0 8px 18px rgba(217, 119, 6, 0.10)'
          : 'inset 1px 1px 0 rgba(255,255,255,0.08), 0 10px 20px rgba(0, 0, 0, 0.16)',
      }}
    >
      <strong style={{ fontSize: '13px', lineHeight: 1.25, color: theme === 'light' ? '#92400E' : '#FDE68A' }}>
        {label}
      </strong>
      <span style={{ fontSize: '14px', lineHeight: 1.45, fontWeight: 850, color: themeTokens.title }}>
        {message}
      </span>
    </div>
  );
}

function RestNotice({
  theme,
  copy,
  onConfirm,
}: {
  theme: PointPageTheme;
  copy: PointLearningCopy['lumiGuide'];
  onConfirm: () => void;
}) {
  return (
    <div
      role="status"
      aria-live="polite"
      style={{
        display: 'grid',
        gridTemplateColumns: 'minmax(0, 1fr) auto',
        gap: '10px',
        alignItems: 'center',
        padding: '10px 12px',
        borderRadius: '7px',
        border: `1px solid ${theme === 'light' ? 'rgba(217, 119, 6, 0.46)' : 'rgba(251, 191, 36, 0.38)'}`,
        background: theme === 'light'
          ? 'linear-gradient(180deg, rgba(255, 247, 214, 0.98), rgba(254, 232, 138, 0.92))'
          : 'linear-gradient(180deg, rgba(92, 58, 12, 0.72), rgba(69, 26, 3, 0.52))',
        boxShadow: theme === 'light'
          ? '0 10px 24px rgba(217, 119, 6, 0.14), inset 1px 1px 0 rgba(255,255,255,0.82)'
          : '0 12px 24px rgba(0, 0, 0, 0.22), inset 1px 1px 0 rgba(255,255,255,0.08)',
      }}
    >
      <div style={{ display: 'grid', gap: '3px', minWidth: 0 }}>
        <strong style={{ fontSize: '13px', lineHeight: 1.25, color: theme === 'light' ? '#7C2D12' : '#FDE68A' }}>
          {copy.restTitle}
        </strong>
        <span style={{ fontSize: '13px', lineHeight: 1.45, fontWeight: 750, color: theme === 'light' ? '#5F3812' : '#FEF3C7' }}>
          {copy.restMessage}
        </span>
      </div>
      <button
        type="button"
        onClick={onConfirm}
        style={{
          ...secondaryButtonStyle,
          minHeight: '34px',
          padding: '0 14px',
          borderRadius: '999px',
          background: theme === 'light' ? '#F59E0B' : 'rgba(251, 191, 36, 0.88)',
          borderColor: theme === 'light' ? '#B45309' : 'rgba(253, 230, 138, 0.68)',
          color: theme === 'light' ? '#431407' : '#422006',
          fontSize: '13px',
          fontWeight: 900,
          whiteSpace: 'nowrap',
        }}
      >
        {copy.restConfirm}
      </button>
    </div>
  );
}
