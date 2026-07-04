'use client'

import { useState, type CSSProperties } from 'react'
import { AddChildObjectModal } from '@/components/explorer-plan/tree/ExplorerEditModals'
import {
  addNodeButtonStyle,
  addSubRegionButtonStyle,
  cancelButtonStyle,
  saveButtonStyle,
  treeButtonDisabledStyle,
} from '@/components/explorer-plan/tree/explorerTreeStyles'
import { useExplorerPlan } from '@/components/explorer-plan/useExplorerPlan'
import { useExplorerEditForm } from '@/components/explorer-plan/tree/useExplorerEditForm'
import type { PlanetDiaryCopy } from '@/lib/i18n/pages/planetDiary'
import type { DiaryRegionAggregate, DiarySubRegionAggregate } from './planetDiaryTypes'

type PlanetRouteKind = 'learning' | 'shared'

const journalSecondaryActionButtonStyle: CSSProperties = {
  ...addNodeButtonStyle,
  border: '1px solid rgba(39, 94, 77, 0.32)',
  background: 'linear-gradient(180deg, rgba(232, 247, 241, 0.94), rgba(204, 232, 220, 0.96))',
  color: '#1E5A48',
}

const compactDiaryActionPanelStyle: CSSProperties = {
  alignSelf: 'end',
  borderRadius: 8,
  border: '1px solid rgba(39, 94, 77, 0.24)',
  background: 'rgba(224, 247, 239, 0.78)',
  boxShadow: 'inset 0 1px 0 rgba(255,255,255,0.30)',
  padding: '5px 6px',
  minHeight: 0,
  maxHeight: 96,
  overflow: 'hidden',
}

const compactDiaryActionFormStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'auto minmax(0, 1fr)',
  alignItems: 'center',
  gap: 6,
  minHeight: 0,
}

const compactDiaryActionLabelStyle: CSSProperties = {
  fontSize: 10,
  fontWeight: 800,
  color: '#2C6A56',
  whiteSpace: 'nowrap',
}

const compactDiaryActionButtonStyle: CSSProperties = {
  height: 24,
  minHeight: 24,
  minWidth: 92,
  flexShrink: 0,
  padding: '0 9px',
  borderRadius: 7,
  fontSize: 11,
  lineHeight: 1,
}

const compactDiaryMutationMsgStyle: CSSProperties = {
  margin: '3px 0 0',
  fontSize: 10,
  lineHeight: 1.2,
  fontWeight: 700,
  color: 'rgba(30, 90, 72, 0.82)',
  whiteSpace: 'nowrap',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
}

export const diaryActionBarStyles = {
  saveButtonStyle,
  cancelButtonStyle,
  treeButtonDisabledStyle,
}

export function DiaryGuidePanel({
  explorerPlan,
  editForm,
  selectedRegion,
  selectedSubRegion,
  routeKind,
  copy,
  isOnboardingAddTarget = false,
  onOnboardingAddClick,
  onOnboardingElementAdded,
}: {
  explorerPlan: ReturnType<typeof useExplorerPlan>
  editForm: ReturnType<typeof useExplorerEditForm>
  selectedRegion: DiaryRegionAggregate | null
  selectedSubRegion: DiarySubRegionAggregate | null
  routeKind: PlanetRouteKind
  copy: PlanetDiaryCopy['guidePanel']
  isOnboardingAddTarget?: boolean
  onOnboardingAddClick?: () => void
  onOnboardingElementAdded?: () => void
}) {
  const [actionMessage, setActionMessage] = useState<string | null>(null)

  const canAdd = Boolean(selectedRegion || selectedSubRegion)
  const canEdit = routeKind === 'learning' && canAdd
  const isActionDisabled = explorerPlan.isMutating || editForm.isMutating

  if (!canEdit) return null

  const handleAddSubRegion = async () => {
    await editForm.handleAddSubRegion()
    onOnboardingElementAdded?.()
    setActionMessage(copy.unsavedPointNotice)
  }

  const handleAddExploration = async () => {
    await editForm.handleAddNode('exploration')
    onOnboardingElementAdded?.()
    setActionMessage(copy.unsavedPointNotice)
  }

  const handleAddResearch = async () => {
    await editForm.handleAddNode('research')
    onOnboardingElementAdded?.()
    setActionMessage(copy.unsavedPointNotice)
  }

  const handleApplyRecommendation = async (candidate: Parameters<typeof editForm.handleApplyRecommendation>[0]) => {
    await editForm.handleApplyRecommendation(candidate)
    onOnboardingElementAdded?.()
    setActionMessage(copy.recommendationAdded)
  }

  return (
    <div style={{
      ...compactDiaryActionPanelStyle,
      ...(isOnboardingAddTarget ? onboardingGuidePanelTargetStyle : undefined),
    }}>
      <div style={compactDiaryActionFormStyle}>
        <span style={compactDiaryActionLabelStyle}>{copy.action}</span>
        <button
          type="button"
          style={{
            ...journalSecondaryActionButtonStyle,
            ...addSubRegionButtonStyle,
            ...compactDiaryActionButtonStyle,
            ...(isActionDisabled ? treeButtonDisabledStyle : undefined),
          }}
          disabled={isActionDisabled}
          onClick={() => {
            onOnboardingAddClick?.()
            editForm.setAddModalType('child-object')
          }}
        >
          {copy.addExplorationElement}
        </button>
      </div>
      {actionMessage || explorerPlan.mutationMessage ? (
        <p style={compactDiaryMutationMsgStyle}>{actionMessage ?? explorerPlan.mutationMessage}</p>
      ) : null}

      {editForm.addModalType === 'child-object' ? (
        <AddChildObjectModal
          canAddSubRegion={editForm.canAddSubRegion}
          isRegionSelected={editForm.isRegionSelected}
          parentRegionName={selectedRegion?.region.name ?? null}
          parentSubRegionName={selectedSubRegion?.subregion.name ?? null}
          subRegionNameInput={editForm.newSubRegionNameInput}
          nodeTitleInput={editForm.newNodeTitleInput}
          nodeUrlInput={editForm.newNodeUrlInput}
          urlCheckStatus={editForm.newUrlCheckStatus}
          urlCheckMsg={editForm.newUrlCheckMsg}
          recommendationQuery={editForm.recommendationQuery}
          recommendationCandidates={editForm.recommendationCandidates}
          isLoadingRecommendations={editForm.isLoadingRecommendations}
          recommendationMessage={editForm.recommendationMessage}
          recommendationPointError={editForm.recommendationPointError}
          onSubRegionNameChange={editForm.setNewSubRegionNameInput}
          onNodeTitleChange={editForm.setNewNodeTitleInput}
          onNodeUrlChange={editForm.setNewNodeUrlInput}
          onRecommendationQueryChange={editForm.setRecommendationQuery}
          onAddSubRegion={() => void handleAddSubRegion()}
          onCheckUrl={() => void editForm.handleCheckUrl('new')}
          onSearchRecommendations={() => void editForm.handleSearchRecommendations()}
          onApplyRecommendation={(candidate) => void handleApplyRecommendation(candidate)}
          onAddExploration={() => void handleAddExploration()}
          onAddResearch={() => void handleAddResearch()}
          onClose={() => editForm.setAddModalType(null)}
          isMutating={editForm.isMutating}
        />
      ) : null}
    </div>
  )
}

const onboardingGuidePanelTargetStyle: CSSProperties = {
  position: 'relative',
  zIndex: 12,
  boxShadow: '0 0 0 2px rgba(16, 185, 129, 0.82), 0 0 22px rgba(16, 185, 129, 0.42)',
}
