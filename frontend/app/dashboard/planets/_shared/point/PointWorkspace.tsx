'use client';

import { useEffect, useState } from 'react';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../pointPageUtils';
import { NodeSection } from '../pointPageUtils';
import type { UsePointLearningResult } from './usePointLearning';
import type { LearningWorkGuideDetail } from './workspace/workspaceTypes';
import ResearchMaterialSection from './workspace/ResearchMaterialSection';
import ExplorationContentSection from './workspace/ExplorationContentSection';
import LearningWorkSection from './workspace/LearningWorkSection';

interface Props {
  copy: PointLearningCopy['workspace'];
  learning: UsePointLearningResult;
  themeTokens: PointPageThemeTokens;
  isExplorationPoint: boolean;
  isSharedRoute: boolean;
  isLearningContentOpen: boolean;
  isLearningWorkOpen: boolean;
  editorToolbarStickyTop: string;
  onResearchMaterialConfirmed: () => void;
  onLearningWorkGuideChange: (guide: LearningWorkGuideDetail) => void;
  onSourceContentOpen: () => void;
  researchEditModeSignal?: number;
  initialResearchMaterialEntry?: 'video' | 'content' | 'attachments';
  hideResearchMaterialEntryNav?: boolean;
  researchMaterialSaveSignal?: number;
  onResearchMaterialEntrySaved?: () => void;
  onResearchMaterialEntrySaveDisabledChange?: (disabled: boolean) => void;
  showLearningContentTitle?: boolean;
}

export default function PointWorkspace({
  copy,
  learning,
  themeTokens,
  isExplorationPoint,
  isSharedRoute,
  isLearningContentOpen,
  isLearningWorkOpen,
  editorToolbarStickyTop,
  onResearchMaterialConfirmed,
  onLearningWorkGuideChange,
  onSourceContentOpen,
  researchEditModeSignal = 0,
  initialResearchMaterialEntry,
  hideResearchMaterialEntryNav = false,
  researchMaterialSaveSignal = 0,
  onResearchMaterialEntrySaved,
  onResearchMaterialEntrySaveDisabledChange,
  showLearningContentTitle = true,
}: Props) {
  const [isSectionMiniNavSide, setIsSectionMiniNavSide] = useState(false);

  const { routeKind, pointDetail, planetID, pointID, isResearchMaterialConfirmed } = learning;
  const isPointCompleted = pointDetail?.point.status === 'completed';
  const canEditLearningRecords = routeKind === 'learning' && !isPointCompleted;
  const canEditResearchMaterial = routeKind === 'learning' && !isPointCompleted;
  const isResearchWorkLocked = !isExplorationPoint && !isResearchMaterialConfirmed;
  const workspaceSubtitle = isExplorationPoint ? copy.explorationSubtitle : '';

  useEffect(() => {
    const query = window.matchMedia('(min-width: 1120px)');
    const updateLayout = () => setIsSectionMiniNavSide(query.matches);
    updateLayout();
    query.addEventListener('change', updateLayout);
    return () => query.removeEventListener('change', updateLayout);
  }, []);

  if (!pointDetail) return null;

  return (
    <>
      {isLearningContentOpen ? (
        <NodeSection id="point-material-section" title={copy.materialTitle} subtitle={workspaceSubtitle} tokens={themeTokens} scrollMarginTop="430px" showHeader={showLearningContentTitle}>
          {!isExplorationPoint ? (
            <ResearchMaterialSection
              copy={copy.researchContentActions}
              editorCopy={copy.editor}
              materialCopy={copy.researchMaterial}
              learning={learning}
              themeTokens={themeTokens}
              canEditResearchMaterial={canEditResearchMaterial}
              isSectionMiniNavSide={isSectionMiniNavSide}
              editorToolbarStickyTop={editorToolbarStickyTop}
              planetID={planetID ?? ''}
              pointID={pointID ?? ''}
              onResearchMaterialConfirmed={onResearchMaterialConfirmed}
              editModeSignal={researchEditModeSignal}
              initialResearchEntry={initialResearchMaterialEntry}
              hideEntryNav={hideResearchMaterialEntryNav}
              saveSignal={researchMaterialSaveSignal}
              onEntrySaved={onResearchMaterialEntrySaved}
              onEntrySaveDisabledChange={onResearchMaterialEntrySaveDisabledChange}
            />
          ) : null}
          <ExplorationContentSection
            copy={copy.explorationContent}
            learning={learning}
            themeTokens={themeTokens}
            isExplorationPoint={isExplorationPoint}
            canEditLearningRecords={canEditLearningRecords}
            onSourceContentOpen={onSourceContentOpen}
          />
        </NodeSection>
      ) : null}

      {isLearningWorkOpen ? (
        <LearningWorkSection
          copy={copy.learningWork}
          editorCopy={copy.editor}
          learning={learning}
          themeTokens={themeTokens}
          canEditLearningRecords={canEditLearningRecords}
          isSharedRoute={isSharedRoute}
          isResearchWorkLocked={isResearchWorkLocked}
          isSectionMiniNavSide={isSectionMiniNavSide}
          editorToolbarStickyTop={editorToolbarStickyTop}
          onLearningWorkGuideChange={onLearningWorkGuideChange}
        />
      ) : null}

      <style jsx global>{`
        .lw-tiptap-editor {
          width: 100%;
          max-width: 860px;
          margin-left: auto;
          margin-right: auto;
          min-height: 220px;
          outline: none;
          font-size: 16px;
          line-height: 1.72;
        }

        .lw-tiptap-editor p {
          margin: 0 0 0.85em;
        }

        .lw-tiptap-editor h1,
        .lw-tiptap-editor h2,
        .lw-tiptap-editor h3 {
          margin: 1em 0 0.45em;
          line-height: 1.25;
          font-weight: 850;
        }

        .lw-tiptap-editor h1 {
          font-size: 1.65em;
        }

        .lw-tiptap-editor h2 {
          font-size: 1.38em;
        }

        .lw-tiptap-editor h3 {
          font-size: 1.16em;
        }

        .lw-tiptap-editor ul,
        .lw-tiptap-editor ol {
          margin: 0.7em 0 0.9em 1.35em;
          padding-left: 1.2em;
        }

        .lw-tiptap-editor ul {
          list-style-type: disc;
        }

        .lw-tiptap-editor ol {
          list-style-type: decimal;
        }

        .lw-tiptap-editor blockquote {
          margin: 0.9em 0;
          padding-left: 1em;
          border-left: 3px solid currentColor;
          opacity: 0.78;
        }

        .lw-tiptap-editor img {
          max-width: 100%;
          height: auto;
          border-radius: 8px;
        }

        .lw-tiptap-editor .tableWrapper {
          margin: 0.8em 0 1em;
          overflow-x: auto;
        }

        .lw-tiptap-editor table {
          width: 100%;
          border-collapse: collapse;
          table-layout: fixed;
          font-size: 0.94em;
        }

        .lw-tiptap-editor th,
        .lw-tiptap-editor td {
          min-width: 88px;
          border: 1px solid currentColor;
          padding: 8px 10px;
          vertical-align: top;
        }

        .lw-tiptap-editor th {
          font-weight: 850;
          background: color-mix(in srgb, currentColor 8%, transparent);
        }

        .lw-tiptap-editor th p,
        .lw-tiptap-editor td p {
          margin: 0;
        }
      `}</style>
    </>
  );
}
