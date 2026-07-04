'use client';

import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageTheme, getPointPageThemeTokens } from '../pointPageUtils';
import { getPointToolbarButtonStyle } from './pointToolbarStyles';

type ThemeTokens = ReturnType<typeof getPointPageThemeTokens>;

export function ToolbarSectionButtons({
  copy,
  theme,
  themeTokens,
  isLearningContentOpen,
  isLearningWorkOpen,
  isCompletionOpen,
  canOpenLearningWork,
  canOpenSelfEvaluation,
  isExplorationPoint,
  isContentViewChecked,
  isLearningWorkChecked,
  isSelfEvaluationChecked,
  isLearningContentToolbarActive,
  isLearningWorkToolbarActive,
  isSelfEvaluationToolbarActive,
  onToggleLearningContent,
  onToggleLearningWork,
  onToggleSelfEvaluation,
}: {
  copy: PointLearningCopy['toolbar'];
  theme: PointPageTheme;
  themeTokens: ThemeTokens;
  isLearningContentOpen: boolean;
  isLearningWorkOpen: boolean;
  isCompletionOpen: boolean;
  canOpenLearningWork: boolean;
  canOpenSelfEvaluation: boolean;
  isExplorationPoint: boolean;
  isContentViewChecked: boolean;
  isLearningWorkChecked: boolean;
  isSelfEvaluationChecked: boolean;
  isLearningContentToolbarActive: boolean;
  isLearningWorkToolbarActive: boolean;
  isSelfEvaluationToolbarActive: boolean;
  onToggleLearningContent: () => void;
  onToggleLearningWork: () => void;
  onToggleSelfEvaluation: () => void;
}) {
  const contentLabel = isExplorationPoint ? copy.content : copy.researchContent;

  return (
    <div
      className="lw-point-learning-toolbar-row"
      role="tablist"
      aria-label={copy.ariaLabel}
      style={{
        display: 'flex',
        justifyContent: 'flex-start',
        gap: 0,
        alignItems: 'center',
        overflowX: 'auto',
        overflowY: 'hidden',
        flex: '1 1 0',
        maxWidth: '100%',
        flexWrap: 'nowrap',
        minWidth: 0,
        padding: '2px 2px 2px 8px',
        scrollbarWidth: 'thin',
        WebkitOverflowScrolling: 'touch',
        border: `1px solid ${theme === 'light' ? '#AAB6C8' : 'rgba(100, 116, 139, 0.72)'}`,
        borderRadius: '8px',
        background: theme === 'light' ? 'rgba(248, 250, 252, 0.88)' : 'rgba(15, 23, 42, 0.72)',
      }}
    >
      <button
        type="button"
        role="tab"
        aria-label={contentLabel}
        aria-selected={isLearningContentToolbarActive}
        aria-pressed={isLearningContentOpen}
        title={contentLabel}
        onClick={onToggleLearningContent}
        style={getPointToolbarButtonStyle({
          theme,
          themeTokens,
          tone: 'content',
          isActive: isLearningContentToolbarActive,
        })}
      >
        <span>{isContentViewChecked ? '☑' : '☐'} {contentLabel}</span>
      </button>
      <button
        type="button"
        role="tab"
        aria-label={copy.work}
        aria-selected={isLearningWorkToolbarActive}
        aria-pressed={isLearningWorkOpen}
        disabled={!canOpenLearningWork}
        title={copy.work}
        onClick={onToggleLearningWork}
        style={getPointToolbarButtonStyle({
          theme,
          themeTokens,
          tone: 'work',
          isActive: isLearningWorkToolbarActive,
          disabled: !canOpenLearningWork,
        })}
      >
        <span>{isLearningWorkChecked ? '☑' : '☐'} {copy.work}</span>
      </button>
      <button
        type="button"
        role="tab"
        aria-label={canOpenSelfEvaluation ? copy.evaluation : copy.evaluationDisabled}
        aria-selected={isSelfEvaluationToolbarActive}
        aria-pressed={isCompletionOpen}
        disabled={!canOpenSelfEvaluation}
        title={canOpenSelfEvaluation ? copy.evaluation : copy.evaluationDisabledTitle}
        onClick={onToggleSelfEvaluation}
        style={getPointToolbarButtonStyle({
          theme,
          themeTokens,
          tone: 'evaluation',
          isActive: isSelfEvaluationToolbarActive,
          disabled: !canOpenSelfEvaluation,
        })}
      >
        <span>{isSelfEvaluationChecked ? '☑' : '☐'} {copy.evaluation}</span>
      </button>
    </div>
  );
}
