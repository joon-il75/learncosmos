'use client';

import { useState } from 'react';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../pointPageUtils';
import { pointLearningStages, type PointLearningStage } from './skillFlowTypes';

interface Props {
  copy: PointLearningCopy['skillFlow'];
  activeStage: PointLearningStage;
  themeTokens: PointPageThemeTokens;
  checkedStages: Record<PointLearningStage, boolean>;
  unlockedStages: Record<PointLearningStage, boolean>;
  onStageSelect: (stage: PointLearningStage) => void;
}

export function PointFlowProgress({ copy, activeStage, themeTokens, checkedStages, unlockedStages, onStageSelect }: Props) {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const activeIndex = pointLearningStages.indexOf(activeStage);
  const activeStageCopy = copy.stages[activeStage];
  const stageAccentStyles: Record<PointLearningStage, { border: string; background: string; color: string; boxShadow: string }> = {
    point: {
      border: '1px solid rgba(55, 138, 221, 0.76)',
      background: '#378ADD',
      color: '#FFFDF7',
      boxShadow: '0 10px 20px rgba(55, 138, 221, 0.20)',
    },
    goal: {
      border: '1px solid rgba(239, 159, 39, 0.76)',
      background: '#EF9F27',
      color: '#241505',
      boxShadow: '0 10px 20px rgba(239, 159, 39, 0.20)',
    },
    skill_discovery: {
      border: '1px solid rgba(127, 119, 221, 0.76)',
      background: '#7F77DD',
      color: '#FFFDF7',
      boxShadow: '0 10px 20px rgba(127, 119, 221, 0.20)',
    },
    goal_iteration: {
      border: '1px solid rgba(28, 125, 121, 0.76)',
      background: '#1C7D79',
      color: '#F8FFFC',
      boxShadow: '0 10px 20px rgba(28, 125, 121, 0.20)',
    },
    final_result: {
      border: '1px solid rgba(204, 82, 22, 0.76)',
      background: '#CC5216',
      color: '#FFF7ED',
      boxShadow: '0 10px 20px rgba(204, 82, 22, 0.20)',
    },
    complete: {
      border: '1px solid rgba(82, 183, 136, 0.76)',
      background: '#52B788',
      color: '#061B13',
      boxShadow: '0 10px 20px rgba(82, 183, 136, 0.20)',
    },
  };
  const activeStageAccentStyle = stageAccentStyles[activeStage];
  const renderStageButton = (stage: PointLearningStage, index: number) => {
    const active = stage === activeStage;
    const complete = checkedStages[stage];
    const unlocked = unlockedStages[stage];
    const stageAccentStyle = stageAccentStyles[stage];
    const numberColor = stage === 'goal' || stage === 'complete' ? 'rgba(36, 21, 5, 0.68)' : 'rgba(255, 253, 247, 0.82)';
    return (
      <button
        key={stage}
        type="button"
        aria-current={active ? 'step' : undefined}
        disabled={!unlocked}
        onClick={() => {
          if (!unlocked) return;
          onStageSelect(stage);
          setIsMobileMenuOpen(false);
        }}
        style={{
          scrollMarginTop: '340px',
          minHeight: '52px',
          display: 'grid',
          alignContent: 'center',
          gap: '3px',
          padding: '8px 10px',
          borderRadius: '8px',
          border: stageAccentStyle.border,
          background: stageAccentStyle.background,
          color: stageAccentStyle.color,
          boxShadow: active ? stageAccentStyle.boxShadow : complete ? '0 5px 12px rgba(18, 24, 38, 0.08)' : undefined,
          opacity: !unlocked ? 0.48 : active ? 1 : complete ? 0.86 : 0.94,
          cursor: unlocked ? 'pointer' : 'not-allowed',
          filter: unlocked ? undefined : 'grayscale(0.38)',
          fontFamily: 'inherit',
          textAlign: 'left',
        }}
      >
        <span style={{ display: 'inline-flex', alignItems: 'center', gap: '6px', fontSize: '11px', fontWeight: 800, color: numberColor }}>
          <input
            type="checkbox"
            checked={complete}
            readOnly
            aria-label={copy.stages[stage].label + (complete ? ' checked' : ' unchecked')}
            tabIndex={-1}
            style={{ width: '13px', height: '13px', accentColor: stageAccentStyle.color === '#241505' || stageAccentStyle.color === '#061B13' ? '#1C7D79' : '#FFFDF7', pointerEvents: 'none' }}
          />
          {index + 1}
        </span>
        <strong style={{ fontSize: '13px', lineHeight: 1.25 }}>{copy.stages[stage].label}</strong>
      </button>
    );
  };

  return (
    <nav aria-label={copy.progressLabel} style={{ display: 'grid', gap: '10px', scrollMarginTop: '340px' }}>
      <button
        type="button"
        className="lw-point-flow-mobile-trigger"
        aria-expanded={isMobileMenuOpen}
        onClick={() => setIsMobileMenuOpen((open) => !open)}
        style={{
          display: 'none',
          width: '100%',
          minHeight: '52px',
          alignItems: 'center',
          justifyContent: 'space-between',
          gap: '10px',
          padding: '10px 12px',
          borderRadius: '8px',
          border: activeStageAccentStyle.border,
          background: activeStageAccentStyle.background,
          color: activeStageAccentStyle.color,
          boxShadow: activeStageAccentStyle.boxShadow,
          font: 'inherit',
          textAlign: 'left',
        }}
      >
        <span style={{ display: 'grid', gap: '2px', minWidth: 0 }}>
          <span style={{ display: 'inline-flex', alignItems: 'center', gap: '6px', color: activeStage === 'goal' || activeStage === 'complete' ? 'rgba(36, 21, 5, 0.68)' : 'rgba(255, 253, 247, 0.82)', fontSize: '11px', fontWeight: 850 }}>
            <input type="checkbox" checked={checkedStages[activeStage]} readOnly tabIndex={-1} style={{ width: '13px', height: '13px', pointerEvents: 'none' }} />
            {activeIndex + 1} / {pointLearningStages.length}
          </span>
          <strong style={{ fontSize: '15px', lineHeight: 1.25 }}>{activeStageCopy.label}</strong>
        </span>
        <span aria-hidden="true" style={{ flex: '0 0 auto', color: activeStage === 'goal' || activeStage === 'complete' ? 'rgba(36, 21, 5, 0.72)' : 'rgba(255, 253, 247, 0.86)', fontSize: '14px', fontWeight: 900 }}>{isMobileMenuOpen ? '▲' : '▼'}</span>
      </button>
      <div className={`lw-point-flow-progress ${isMobileMenuOpen ? 'lw-point-flow-progress--open' : ''}`} style={{ display: 'grid', gridTemplateColumns: 'repeat(6, minmax(112px, 1fr))', gap: '8px', overflowX: 'auto', paddingBottom: '2px' }}>
        {pointLearningStages.map(renderStageButton)}
      </div>
      <style jsx>{`
        @media (max-width: 760px) {
          .lw-point-flow-mobile-trigger {
            display: flex !important;
          }
          .lw-point-flow-progress {
            display: none !important;
            grid-template-columns: repeat(2, minmax(0, 1fr)) !important;
            overflow-x: visible !important;
            padding-bottom: 0 !important;
          }
          .lw-point-flow-progress--open {
            display: grid !important;
          }
          .lw-point-flow-progress > button {
            scroll-margin-top: 340px;
          }
        }
      `}</style>
    </nav>
  );
}
