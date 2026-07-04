'use client'

import { useState, type CSSProperties } from 'react'
import ModalPortal from '@/components/common/ModalPortal'
import type { ExplorerContentCandidate } from './useExplorerEditForm'
import { getDashboardCourseDraftCopy, type DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'
import {
  treeButtonDisabledStyle,
  treePrimaryButtonStyle,
  addExplorationOverlayStyle,
  addExplorationModalBoxStyle,
  explorationChoiceBtnStyle,
  addExplorationButtonStyle,
  childObjectTabRowStyle,
  childObjectTabButtonStyle,
  childObjectTabButtonActiveStyle,
  addItemModalHeaderStyle,
  addItemModalTitleStyle,
  addItemModalLabelStyle,
  addItemModalInputStyle,
  addItemModalRowStyle,
  addItemModalPrimaryBtnStyle,
  addItemModalCancelBtnStyle,
  addItemModalCheckBtnStyle,
  addItemModalErrorStyle,
  addItemModalInfoStyle,
  parentContextInfoStyle,
  recommendationListStyle,
  recommendationItemStyle,
  recommendationItemTitleStyle,
  recommendationItemMetaStyle,
} from './explorerTreeStyles'
import { RecommendationStatusMessage } from './ExplorerRecommendModal'

const defaultModalCopy = getDashboardCourseDraftCopy('ko').tree.editPanel.modals

// ── helpers ──

const addChildObjectPanelStyle: CSSProperties = {
  display: 'grid',
  gap: 14,
  minHeight: 0,
  flex: 1,
}

type DescriptionCalloutTone = 'exploration' | 'research' | 'subregion' | 'recommendation' | 'url'

const descriptionCalloutTheme: Record<DescriptionCalloutTone, {
  accent: string
  border: string
  background: string
  labelBackground: string
  labelColor: string
  color: string
  shadow: string
}> = {
  exploration: {
    accent: '#2563eb',
    border: 'rgba(37, 99, 235, 0.34)',
    background: 'linear-gradient(135deg, rgba(239, 246, 255, 0.96), rgba(219, 234, 254, 0.78))',
    labelBackground: 'rgba(37, 99, 235, 0.12)',
    labelColor: '#1d4ed8',
    color: '#1e3a8a',
    shadow: '0 10px 22px rgba(37, 99, 235, 0.14)',
  },
  research: {
    accent: '#7c3aed',
    border: 'rgba(124, 58, 237, 0.34)',
    background: 'linear-gradient(135deg, rgba(245, 243, 255, 0.98), rgba(237, 233, 254, 0.80))',
    labelBackground: 'rgba(124, 58, 237, 0.12)',
    labelColor: '#6d28d9',
    color: '#4c1d95',
    shadow: '0 10px 22px rgba(124, 58, 237, 0.14)',
  },
  subregion: {
    accent: '#d97706',
    border: 'rgba(217, 119, 6, 0.34)',
    background: 'linear-gradient(135deg, rgba(255, 251, 235, 0.98), rgba(254, 243, 199, 0.82))',
    labelBackground: 'rgba(217, 119, 6, 0.13)',
    labelColor: '#b45309',
    color: '#78350f',
    shadow: '0 10px 22px rgba(217, 119, 6, 0.13)',
  },
  recommendation: {
    accent: '#059669',
    border: 'rgba(5, 150, 105, 0.34)',
    background: 'linear-gradient(135deg, rgba(236, 253, 245, 0.98), rgba(209, 250, 229, 0.78))',
    labelBackground: 'rgba(5, 150, 105, 0.12)',
    labelColor: '#047857',
    color: '#064e3b',
    shadow: '0 10px 22px rgba(5, 150, 105, 0.13)',
  },
  url: {
    accent: '#0891b2',
    border: 'rgba(8, 145, 178, 0.34)',
    background: 'linear-gradient(135deg, rgba(236, 254, 255, 0.98), rgba(207, 250, 254, 0.78))',
    labelBackground: 'rgba(8, 145, 178, 0.12)',
    labelColor: '#0e7490',
    color: '#164e63',
    shadow: '0 10px 22px rgba(8, 145, 178, 0.13)',
  },
}

const addObjectDescriptionStyle = (tone: DescriptionCalloutTone, nested = false): CSSProperties => {
  const theme = descriptionCalloutTheme[tone]
  return {
    ...addItemModalInfoStyle,
    position: 'relative',
    display: 'grid',
    gap: 5,
    marginTop: nested ? -4 : 0,
    padding: '10px 12px 10px 16px',
    borderRadius: 10,
    border: `1px solid ${theme.border}`,
    borderLeft: `5px solid ${theme.accent}`,
    background: theme.background,
    color: theme.color,
    fontSize: 12,
    lineHeight: 1.5,
    fontWeight: 750,
    boxShadow: theme.shadow,
  }
}

const descriptionCalloutLabelStyle = (tone: DescriptionCalloutTone): CSSProperties => {
  const theme = descriptionCalloutTheme[tone]
  return {
    justifySelf: 'start',
    padding: '2px 7px',
    borderRadius: 999,
    background: theme.labelBackground,
    color: theme.labelColor,
    fontSize: 10,
    lineHeight: 1.2,
    fontWeight: 900,
    letterSpacing: '0.03em',
  }
}

function DescriptionCallout({
  tone,
  label,
  nested = false,
  children,
}: {
  tone: DescriptionCalloutTone
  label: string
  nested?: boolean
  children: string
}) {
  return (
    <div style={addObjectDescriptionStyle(tone, nested)}>
      <span style={descriptionCalloutLabelStyle(tone)}>{label}</span>
      <span>{children}</span>
    </div>
  )
}

const nestedTabRowStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '1fr 1fr',
  gap: 0,
  padding: '3px 3px 0',
  borderBottom: '1px solid rgba(196, 142, 54, 0.24)',
  background: 'rgba(196, 142, 54, 0.08)',
  borderRadius: '11px 11px 0 0',
  flexShrink: 0,
}

const nestedTabPanelStyle: CSSProperties = {
  display: 'grid',
  gap: 14,
  minHeight: 0,
  flex: 1,
}

// ── 탐험요소 추가 모달 ──

export interface AddChildObjectModalProps {
  canAddSubRegion: boolean
  isRegionSelected: boolean
  parentRegionName?: string | null
  parentSubRegionName?: string | null
  subRegionNameInput: string
  nodeTitleInput: string
  nodeUrlInput: string
  urlCheckStatus: 'idle' | 'checking' | 'ok' | 'error'
  urlCheckMsg: string
  recommendationQuery: string
  recommendationCandidates: ExplorerContentCandidate[]
  isLoadingRecommendations: boolean
  recommendationMessage: string | null
  recommendationPointError: string | null
  onSubRegionNameChange: (value: string) => void
  onNodeTitleChange: (value: string) => void
  onNodeUrlChange: (value: string) => void
  onRecommendationQueryChange: (value: string) => void
  onAddSubRegion: () => void
  onCheckUrl: () => void
  onSearchRecommendations: () => void
  onApplyRecommendation: (candidate: ExplorerContentCandidate) => void
  onAddExploration: () => void
  onAddResearch: () => void
  onClose: () => void
  isMutating: boolean
  copy?: DashboardCourseDraftCopy['tree']['editPanel']['modals']
}

export function AddChildObjectModal({
  canAddSubRegion,
  isRegionSelected,
  parentRegionName,
  parentSubRegionName,
  subRegionNameInput,
  nodeTitleInput,
  nodeUrlInput,
  urlCheckStatus,
  urlCheckMsg,
  recommendationQuery,
  recommendationCandidates,
  isLoadingRecommendations,
  recommendationMessage,
  recommendationPointError,
  onSubRegionNameChange,
  onNodeTitleChange,
  onNodeUrlChange,
  onRecommendationQueryChange,
  onAddSubRegion,
  onCheckUrl,
  onSearchRecommendations,
  onApplyRecommendation,
  onAddExploration,
  onAddResearch,
  onClose,
  isMutating,
  copy = defaultModalCopy,
}: AddChildObjectModalProps) {
  const [activeTab, setActiveTab] = useState<'subregion' | 'exploration' | 'research'>(
    'exploration'
  )
  const [explorationInputTab, setExplorationInputTab] = useState<'recommendation' | 'url'>('recommendation')
  const [recommendationPage, setRecommendationPage] = useState(0)
  const recommendationPageSize = 3
  const recommendationTotalPages = Math.ceil(recommendationCandidates.length / recommendationPageSize)
  const pagedCandidates = recommendationCandidates.slice(
    recommendationPage * recommendationPageSize,
    (recommendationPage + 1) * recommendationPageSize
  )
  const urlPassed = urlCheckStatus === 'ok'
  const canAddExploration = nodeTitleInput.trim() !== '' && urlPassed
  const canAddResearch = nodeTitleInput.trim() !== ''
  const canAddSubRegionNow = isRegionSelected && canAddSubRegion && subRegionNameInput.trim() !== ''
  const parentPath = [parentRegionName?.trim(), parentSubRegionName?.trim()].filter(Boolean).join(' > ')
  const checkBtnStyle: CSSProperties = {
    ...addItemModalCheckBtnStyle,
    ...(urlCheckStatus === 'ok'
      ? { color: '#136857', borderColor: 'rgba(20,140,120,0.5)', background: 'rgba(20,140,120,0.10)' }
      : urlCheckStatus === 'error'
      ? { color: '#9B2C2C', borderColor: 'rgba(155,44,44,0.4)', background: 'rgba(155,44,44,0.06)' }
      : undefined),
  }

  const renderTabButton = (
    key: 'subregion' | 'exploration' | 'research',
    label: string,
    disabled: boolean,
  ) => (
    <button
      type="button"
      style={{
        ...childObjectTabButtonStyle,
        ...(activeTab === key ? childObjectTabButtonActiveStyle : undefined),
        ...(disabled || isMutating ? treeButtonDisabledStyle : undefined),
      }}
      onClick={() => setActiveTab(key)}
      disabled={disabled || isMutating}
      role="tab"
      aria-selected={activeTab === key}
      aria-controls={`explorer-add-${key}-panel`}
    >
      {label}
    </button>
  )

  return (
    <ModalPortal overlayStyle={addExplorationOverlayStyle} onMouseDown={onClose}>
      <div
        role="dialog"
        aria-modal="true"
        aria-label={copy.addChild.aria}
        style={addExplorationModalBoxStyle}
        onMouseDown={(event) => event.stopPropagation()}
      >
        <div style={{ ...addItemModalHeaderStyle, flexShrink: 0 }}>
          <h3 style={addItemModalTitleStyle}>{copy.addChild.title}</h3>
          <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
            {copy.close}
          </button>
        </div>
        {parentPath ? (
          <div style={{ ...parentContextInfoStyle, flexShrink: 0 }}>
            {copy.parentPath(parentPath)}
          </div>
        ) : null}

        <div style={childObjectTabRowStyle} role="tablist" aria-label={copy.addChild.tabAria}>
          {renderTabButton('exploration', copy.addChild.explorationTab, false)}
          {renderTabButton('research', copy.addChild.researchTab, false)}
          {renderTabButton('subregion', copy.addChild.subregionTab, !isRegionSelected || !canAddSubRegion)}
        </div>

        {activeTab === 'subregion' ? (
          <div id="explorer-add-subregion-panel" role="tabpanel" style={addChildObjectPanelStyle}>
            <DescriptionCallout tone="subregion" label={copy.addChild.subregionTab}>
              {copy.addChild.subregionDescription}
            </DescriptionCallout>
            {!canAddSubRegion && (
              <div style={addItemModalInfoStyle}>
                {copy.addChild.subregionLimit}
              </div>
            )}
            <div>
              <label style={addItemModalLabelStyle}>{copy.addChild.subregionNameLabel}</label>
              <input
                style={addItemModalInputStyle}
                value={subRegionNameInput}
                onChange={(event) => onSubRegionNameChange(event.target.value)}
                placeholder={copy.addChild.subregionNamePlaceholder}
                disabled={isMutating || !isRegionSelected || !canAddSubRegion}
                autoFocus
                onKeyDown={(event) => {
                  if (event.key === 'Enter' && canAddSubRegionNow) onAddSubRegion()
                }}
              />
            </div>
            <div style={addItemModalRowStyle}>
              <button
                type="button"
                style={{
                  ...addItemModalPrimaryBtnStyle,
                  ...(!canAddSubRegionNow || isMutating ? treeButtonDisabledStyle : undefined),
                }}
                onClick={onAddSubRegion}
                disabled={!canAddSubRegionNow || isMutating}
              >
                {copy.addChild.addSubregion}
              </button>
              <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
                {copy.cancel}
              </button>
            </div>
          </div>
        ) : activeTab === 'exploration' ? (
          <div id="explorer-add-exploration-panel" role="tabpanel" style={addChildObjectPanelStyle}>
            <DescriptionCallout tone="exploration" label={copy.addChild.explorationTab}>
              {copy.addChild.explorationDescription}
            </DescriptionCallout>
            <div style={nestedTabRowStyle} role="tablist" aria-label={copy.addChild.explorationMethodAria}>
              <button
                type="button"
                style={{
                  ...explorationChoiceBtnStyle,
                  ...(explorationInputTab === 'recommendation' ? addExplorationButtonStyle : undefined),
                  ...(isMutating ? treeButtonDisabledStyle : undefined),
                }}
                onClick={() => {
                  setExplorationInputTab('recommendation')
                  setRecommendationPage(0)
                }}
                disabled={isMutating}
                role="tab"
                aria-selected={explorationInputTab === 'recommendation'}
                aria-controls="explorer-add-recommendation-panel"
              >
                {copy.addExploration.recommendationTab}
              </button>
              <button
                type="button"
                style={{
                  ...explorationChoiceBtnStyle,
                  ...(explorationInputTab === 'url' ? addExplorationButtonStyle : undefined),
                  ...(isMutating ? treeButtonDisabledStyle : undefined),
                }}
                onClick={() => setExplorationInputTab('url')}
                disabled={isMutating}
                role="tab"
                aria-selected={explorationInputTab === 'url'}
                aria-controls="explorer-add-url-panel"
              >
                {copy.addChild.urlDirectTab}
              </button>
            </div>
            {explorationInputTab === 'recommendation' ? (
              <div id="explorer-add-recommendation-panel" role="tabpanel" style={nestedTabPanelStyle}>
                <DescriptionCallout tone="recommendation" label={copy.addExploration.recommendationTab} nested>
                  {copy.addChild.recommendationDescription}
                </DescriptionCallout>
                {recommendationPointError && <div style={{ ...addItemModalErrorStyle, flexShrink: 0 }}>{recommendationPointError}</div>}
                <div>
                  <label style={addItemModalLabelStyle}>{copy.recommendation.queryLabel}</label>
                  <div style={addItemModalRowStyle}>
                    <input
                      style={{ ...addItemModalInputStyle, flex: 1 }}
                      value={recommendationQuery}
                      onChange={(event) => onRecommendationQueryChange(event.target.value)}
                      onKeyDown={(event) => {
                        if (event.key === 'Enter' && !isLoadingRecommendations) {
                          setRecommendationPage(0)
                          onSearchRecommendations()
                        }
                      }}
                      placeholder={copy.recommendation.queryPlaceholder}
                      disabled={isMutating || isLoadingRecommendations}
                      autoFocus
                    />
                    <button
                      type="button"
                      style={{
                        ...addItemModalPrimaryBtnStyle,
                        minWidth: 108,
                        ...((isMutating || isLoadingRecommendations) ? treeButtonDisabledStyle : undefined),
                      }}
                      onClick={() => {
                        setRecommendationPage(0)
                        onSearchRecommendations()
                      }}
                      disabled={isMutating || isLoadingRecommendations}
                    >
                      {isLoadingRecommendations ? copy.recommendation.searching : copy.recommendation.search}
                    </button>
                  </div>
                </div>
                <div style={{ ...recommendationListStyle, flex: 1, overflowY: 'auto', minHeight: 0 }}>
                  {recommendationCandidates.length === 0 && (
                    <RecommendationStatusMessage
                      isLoading={isLoadingRecommendations}
                      message={
                        recommendationPointError
                          ? copy.recommendation.pointErrorEmpty
                          : recommendationMessage ?? copy.recommendation.empty
                      }
                      copy={copy.recommendation}
                    />
                  )}
                  {pagedCandidates.map((candidate) => (
                    <div key={candidate.id} style={recommendationItemStyle}>
                      <div style={{ minWidth: 0 }}>
                        <p style={recommendationItemTitleStyle} title={candidate.title}>{candidate.title}</p>
                        <div style={recommendationItemMetaStyle}>
                          {candidate.content_type}
                          {candidate.content_id ? ` · ${copy.internalContent}` : ''}
                          {candidate.author ? ` · ${candidate.author}` : ''}
                          {candidate.url ? ` · ${candidate.url}` : ''}
                        </div>
                      </div>
                      <div style={{ display: 'flex', gap: 4, flexShrink: 0 }}>
                        {(candidate.url ?? candidate.canonical_url) && (
                          <button
                            type="button"
                            style={{ ...treePrimaryButtonStyle, padding: '0 8px' }}
                            title={copy.openLink}
                            onClick={() => window.open(candidate.url ?? candidate.canonical_url ?? '', '_blank', 'noopener,noreferrer')}
                          >
                            🔗
                          </button>
                        )}
                        <button
                          type="button"
                          style={treePrimaryButtonStyle}
                          onClick={() => onApplyRecommendation(candidate)}
                          disabled={isMutating || (!candidate.content_id && !candidate.url && !candidate.canonical_url)}
                        >
                          {copy.apply}
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
                {recommendationTotalPages > 1 && (
                  <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: 8, padding: '4px 0', flexShrink: 0 }}>
                    <button
                      type="button"
                      style={{ ...treePrimaryButtonStyle, ...(recommendationPage === 0 ? treeButtonDisabledStyle : undefined) }}
                      onClick={() => setRecommendationPage((page) => page - 1)}
                      disabled={recommendationPage === 0}
                    >
                      ◀
                    </button>
                    <span style={{ fontSize: '11px', color: '#4A3520' }}>{recommendationPage + 1} / {recommendationTotalPages}</span>
                    <button
                      type="button"
                      style={{ ...treePrimaryButtonStyle, ...(recommendationPage >= recommendationTotalPages - 1 ? treeButtonDisabledStyle : undefined) }}
                      onClick={() => setRecommendationPage((page) => page + 1)}
                      disabled={recommendationPage >= recommendationTotalPages - 1}
                    >
                      ▶
                    </button>
                  </div>
                )}
              </div>
            ) : (
              <div id="explorer-add-url-panel" role="tabpanel" style={nestedTabPanelStyle}>
                <DescriptionCallout tone="url" label={copy.addChild.urlDirectTab} nested>
                  {copy.addChild.urlDescription}
                </DescriptionCallout>
                <div>
                  <label style={addItemModalLabelStyle}>{copy.addExploration.titleLabel}</label>
                  <input
                    style={addItemModalInputStyle}
                    value={nodeTitleInput}
                    onChange={(event) => onNodeTitleChange(event.target.value)}
                    placeholder={copy.addExploration.titlePlaceholder}
                    disabled={isMutating}
                    autoFocus
                  />
                </div>
                <div>
                  <label style={addItemModalLabelStyle}>{copy.addExploration.urlLabel}</label>
                  <div style={addItemModalRowStyle}>
                    <input
                      style={{ ...addItemModalInputStyle, flex: 1 }}
                      value={nodeUrlInput}
                      onChange={(event) => onNodeUrlChange(event.target.value)}
                      placeholder="https://..."
                      disabled={isMutating}
                    />
                    <button
                      type="button"
                      style={{ ...checkBtnStyle, ...((!nodeUrlInput.trim() || urlCheckStatus === 'checking' || isMutating) ? treeButtonDisabledStyle : undefined) }}
                      onClick={onCheckUrl}
                      disabled={!nodeUrlInput.trim() || urlCheckStatus === 'checking' || isMutating}
                    >
                      {urlCheckStatus === 'checking' ? copy.checking : urlCheckStatus === 'ok' ? copy.valid : urlCheckStatus === 'error' ? copy.error : copy.check}
                    </button>
                  </div>
                  {urlCheckMsg && (
                    <div style={{ ...addItemModalInfoStyle, marginTop: 4, color: urlCheckStatus === 'ok' ? '#136857' : '#9B2C2C' }}>
                      {urlCheckMsg}
                    </div>
                  )}
                </div>
                <div style={addItemModalRowStyle}>
                  <button
                    type="button"
                    style={{ ...addItemModalPrimaryBtnStyle, ...(!canAddExploration || isMutating ? treeButtonDisabledStyle : undefined) }}
                    onClick={onAddExploration}
                    disabled={!canAddExploration || isMutating}
                  >
                    {copy.addExploration.addButton}
                  </button>
                  <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
                    {copy.cancel}
                  </button>
                </div>
              </div>
            )}
          </div>
        ) : (
          <div id="explorer-add-research-panel" role="tabpanel" style={addChildObjectPanelStyle}>
            <DescriptionCallout tone="research" label={copy.addChild.researchTab}>
              {copy.addChild.researchDescription}
            </DescriptionCallout>
            <div>
              <label style={addItemModalLabelStyle}>{copy.addChild.researchTitleLabel}</label>
              <input
                style={addItemModalInputStyle}
                value={nodeTitleInput}
                onChange={(event) => onNodeTitleChange(event.target.value)}
                placeholder={copy.addChild.researchTitlePlaceholder}
                disabled={isMutating}
                autoFocus
              />
            </div>
            <div style={addItemModalRowStyle}>
              <button
                type="button"
                style={{ ...addItemModalPrimaryBtnStyle, ...(!canAddResearch || isMutating ? treeButtonDisabledStyle : undefined) }}
                onClick={onAddResearch}
                disabled={!canAddResearch || isMutating}
              >
                {copy.addChild.addResearch}
              </button>
              <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
                {copy.cancel}
              </button>
            </div>
          </div>
        )}
      </div>
    </ModalPortal>
  )
}
