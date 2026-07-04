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
  recommendationListStyle,
  recommendationItemStyle,
  recommendationItemTitleStyle,
  recommendationItemMetaStyle,
  recommendationMessageStyle,
} from './explorerTreeStyles'
import { RecommendationStatusMessage } from './ExplorerRecommendModal'

const defaultModalCopy = getDashboardCourseDraftCopy('ko').tree.editPanel.modals

// ── 탐험지점 추가 모달 ────────────────────────────────────────────────────────

export interface AddExplorationNodeModalProps {
  query: string
  candidates: ExplorerContentCandidate[]
  isLoading: boolean
  message: string | null
  pointError: string | null
  titleInput: string
  urlInput: string
  urlCheckStatus: 'idle' | 'checking' | 'ok' | 'error'
  urlCheckMsg: string
  onQueryChange: (value: string) => void
  onSearch: () => void
  onApply: (candidate: ExplorerContentCandidate) => void
  onTitleChange: (value: string) => void
  onUrlChange: (value: string) => void
  onCheckUrl: () => void
  onAdd: () => void
  onClose: () => void
  isMutating: boolean
  copy?: DashboardCourseDraftCopy['tree']['editPanel']['modals']
}

export function AddExplorationNodeModal({
  query,
  candidates,
  isLoading,
  message,
  pointError,
  titleInput,
  urlInput,
  urlCheckStatus,
  urlCheckMsg,
  onQueryChange,
  onSearch,
  onApply,
  onTitleChange,
  onUrlChange,
  onCheckUrl,
  onAdd,
  onClose,
  isMutating,
  copy = defaultModalCopy,
}: AddExplorationNodeModalProps) {
  const [activeTab, setActiveTab] = useState<'recommendation' | 'url'>('recommendation')
  const [page, setPage] = useState(0)
  const PAGE_SIZE = 3
  const totalPages = Math.ceil(candidates.length / PAGE_SIZE)
  const paged = candidates.slice(page * PAGE_SIZE, (page + 1) * PAGE_SIZE)
  const handleSearch = () => { setPage(0); onSearch() }
  const urlPassed = urlCheckStatus === 'ok'
  const canAddUrl = titleInput.trim() !== '' && urlPassed
  const checkBtnStyle: CSSProperties = {
    ...addItemModalCheckBtnStyle,
    ...(urlCheckStatus === 'ok'
      ? { color: '#136857', borderColor: 'rgba(20,140,120,0.5)', background: 'rgba(20,140,120,0.10)' }
      : urlCheckStatus === 'error'
      ? { color: '#9B2C2C', borderColor: 'rgba(155,44,44,0.4)', background: 'rgba(155,44,44,0.06)' }
      : undefined),
  }

  const recommendationTabStyle: CSSProperties = {
    ...explorationChoiceBtnStyle,
    ...(activeTab === 'recommendation' ? addExplorationButtonStyle : undefined),
    fontWeight: activeTab === 'recommendation' ? 800 : 700,
  }
  const urlTabStyle: CSSProperties = {
    ...explorationChoiceBtnStyle,
    ...(activeTab === 'url' ? addExplorationButtonStyle : undefined),
    fontWeight: activeTab === 'url' ? 800 : 700,
  }

  return (
    <ModalPortal overlayStyle={addExplorationOverlayStyle} onMouseDown={onClose}>
      <div
        role="dialog"
        aria-modal="true"
        aria-label={copy.addExploration.aria}
        style={addExplorationModalBoxStyle}
        onMouseDown={(event) => event.stopPropagation()}
      >
        <div style={{ ...addItemModalHeaderStyle, flexShrink: 0 }}>
          <h3 style={addItemModalTitleStyle}>{copy.addExploration.title}</h3>
          <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
            {copy.close}
          </button>
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 6, flexShrink: 0 }}>
          <button
            type="button"
            style={{ ...recommendationTabStyle, ...(isMutating ? treeButtonDisabledStyle : undefined) }}
            onClick={() => setActiveTab('recommendation')}
            disabled={isMutating}
          >
            {copy.addExploration.recommendationTab}
          </button>
          <button
            type="button"
            style={{ ...urlTabStyle, ...(isMutating ? treeButtonDisabledStyle : undefined) }}
            onClick={() => setActiveTab('url')}
            disabled={isMutating}
          >
            {copy.addExploration.urlTab}
          </button>
        </div>

        {activeTab === 'recommendation' ? (
          <>
            {pointError && <div style={{ ...addItemModalErrorStyle, flexShrink: 0 }}>{pointError}</div>}
            <div style={{ flexShrink: 0 }}>
              <label style={addItemModalLabelStyle}>{copy.recommendation.queryLabel}</label>
              <div style={addItemModalRowStyle}>
                <input
                  style={{ ...addItemModalInputStyle, flex: 1 }}
                  value={query}
                  onChange={(event) => onQueryChange(event.target.value)}
                  onKeyDown={(event) => { if (event.key === 'Enter' && !isLoading) handleSearch() }}
                  placeholder={copy.recommendation.queryPlaceholder}
                  disabled={isMutating || isLoading}
                  autoFocus
                />
                <button
                  type="button"
                  style={{
                    ...addItemModalPrimaryBtnStyle,
                    minWidth: 108,
                    ...((isMutating || isLoading) ? treeButtonDisabledStyle : undefined),
                  }}
                  onClick={handleSearch}
                  disabled={isMutating || isLoading}
                >
                  {isLoading ? copy.recommendation.searching : copy.recommendation.search}
                </button>
              </div>
            </div>
            <div style={{ ...recommendationListStyle, flex: 1, overflowY: 'auto', minHeight: 0 }}>
              {candidates.length === 0 && (
                <RecommendationStatusMessage
                  isLoading={isLoading}
                  message={
                    pointError
                      ? copy.recommendation.pointErrorEmpty
                      : message ?? copy.recommendation.empty
                  }
                  copy={copy.recommendation}
                />
              )}
              {paged.map((candidate) => (
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
                      onClick={() => onApply(candidate)}
                      disabled={isMutating || (!candidate.content_id && !candidate.url && !candidate.canonical_url)}
                    >
                      {copy.apply}
                    </button>
                  </div>
                </div>
              ))}
            </div>
            {totalPages > 1 && (
              <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: 8, padding: '4px 0', flexShrink: 0 }}>
                <button
                  type="button"
                  style={{ ...treePrimaryButtonStyle, ...(page === 0 ? treeButtonDisabledStyle : undefined) }}
                  onClick={() => setPage((p) => p - 1)}
                  disabled={page === 0}
                >
                  ◀
                </button>
                <span style={{ fontSize: '11px', color: '#4A3520' }}>{page + 1} / {totalPages}</span>
                <button
                  type="button"
                  style={{ ...treePrimaryButtonStyle, ...(page >= totalPages - 1 ? treeButtonDisabledStyle : undefined) }}
                  onClick={() => setPage((p) => p + 1)}
                  disabled={page >= totalPages - 1}
                >
                  ▶
                </button>
              </div>
            )}
            {message && candidates.length > 0 && !pointError && (
              <div style={{ ...recommendationMessageStyle, flexShrink: 0 }}>{message}</div>
            )}
          </>
        ) : (
          <>
            <div>
              <label style={addItemModalLabelStyle}>{copy.addExploration.titleLabel}</label>
              <input
                style={addItemModalInputStyle}
                value={titleInput}
                onChange={(event) => onTitleChange(event.target.value)}
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
                  value={urlInput}
                  onChange={(event) => onUrlChange(event.target.value)}
                  placeholder="https://..."
                  disabled={isMutating}
                />
                <button
                  type="button"
                  style={{ ...checkBtnStyle, ...((!urlInput.trim() || urlCheckStatus === 'checking' || isMutating) ? treeButtonDisabledStyle : undefined) }}
                  onClick={onCheckUrl}
                  disabled={!urlInput.trim() || urlCheckStatus === 'checking' || isMutating}
                >
                  {urlCheckStatus === 'checking' ? copy.checking : urlCheckStatus === 'ok' ? copy.valid : urlCheckStatus === 'error' ? copy.error : copy.check}
                </button>
              </div>
              {urlCheckMsg && (
                <div style={{ ...addItemModalInfoStyle, marginTop: 4, color: urlCheckStatus === 'ok' ? '#136857' : '#9B2C2C' }}>
                  {urlCheckMsg}
                </div>
              )}
              {!urlPassed && urlCheckStatus === 'idle' && urlInput.trim() && (
                <div style={{ ...addItemModalInfoStyle, marginTop: 4 }}>
                  {copy.urlMustPass}
                </div>
              )}
            </div>
            <div style={addItemModalRowStyle}>
              <button
                type="button"
                style={{ ...addItemModalPrimaryBtnStyle, ...(!canAddUrl || isMutating ? treeButtonDisabledStyle : undefined) }}
                onClick={onAdd}
                disabled={!canAddUrl || isMutating}
              >
                {copy.addExploration.addButton}
              </button>
              <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
                {copy.cancel}
              </button>
            </div>
          </>
        )}
      </div>
    </ModalPortal>
  )
}
