'use client'

import { useState, type CSSProperties } from 'react'
import ModalPortal from '@/components/common/ModalPortal'
import type { ExplorerContentCandidate } from './useExplorerEditForm'
import { getDashboardCourseDraftCopy, type DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'
import {
  treeButtonDisabledStyle,
  treePrimaryButtonStyle,
  treeInputStyle,
  addItemModalCancelBtnStyle,
  addItemModalErrorStyle,
  recommendationOverlayStyle,
  recommendationModalStyle,
  recommendationHeaderStyle,
  recommendationTitleStyle,
  recommendationSearchRowStyle,
  recommendationListStyle,
  recommendationItemStyle,
  recommendationItemTitleStyle,
  recommendationItemMetaStyle,
  recommendationMessageStyle,
} from './explorerTreeStyles'

const defaultModalCopy = getDashboardCourseDraftCopy('ko').tree.editPanel.modals

// ── 콘텐츠 추천 모달 ──────────────────────────────────────────────────────────

export interface ContentRecommendationModalProps {
  query: string
  candidates: ExplorerContentCandidate[]
  isLoading: boolean
  message: string | null
  pointError: string | null
  onQueryChange: (value: string) => void
  onSearch: () => void
  onRecommendAgain: () => void
  onClose: () => void
  onApply: (candidate: ExplorerContentCandidate) => void
  copy?: DashboardCourseDraftCopy['tree']['editPanel']['modals']['recommendation']
  commonCopy?: DashboardCourseDraftCopy['tree']['editPanel']['modals']
}

export function ContentRecommendationModal({
  query, candidates, isLoading, message, pointError,
  onQueryChange, onSearch, onRecommendAgain, onClose, onApply,
  copy = defaultModalCopy.recommendation,
  commonCopy = defaultModalCopy,
}: ContentRecommendationModalProps) {
  const [page, setPage] = useState(0)
  const PAGE_SIZE = 3
  const totalPages = Math.ceil(candidates.length / PAGE_SIZE)
  const paged = candidates.slice(page * PAGE_SIZE, (page + 1) * PAGE_SIZE)
  const handleSearch = () => { setPage(0); onSearch() }
  const handleRecommendAgain = () => { setPage(0); onRecommendAgain() }
  return (
    <ModalPortal overlayStyle={recommendationOverlayStyle} onMouseDown={onClose}>
      <div
        role="dialog"
        aria-modal="true"
        aria-label={copy.aria}
        style={{ ...recommendationModalStyle, display: 'flex', flexDirection: 'column' }}
        onMouseDown={(event) => event.stopPropagation()}
      >
        <div style={{ ...recommendationHeaderStyle, flexShrink: 0 }}>
          <h3 style={recommendationTitleStyle}>{copy.title}</h3>
          <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose}>{commonCopy.close}</button>
        </div>

        {pointError && (
          <div style={{ ...addItemModalErrorStyle, flexShrink: 0 }}>{pointError}</div>
        )}

        <div style={{ ...recommendationSearchRowStyle, flexShrink: 0 }}>
          <input
            style={treeInputStyle}
            value={query}
            onChange={(event) => onQueryChange(event.target.value)}
            onKeyDown={(event) => { if (event.key === 'Enter') handleSearch() }}
            placeholder={copy.queryPlaceholder}
            autoFocus
          />
          <button
            type="button"
            style={{ ...treePrimaryButtonStyle, ...(isLoading ? treeButtonDisabledStyle : undefined) }}
            onClick={handleSearch}
            disabled={isLoading}
          >
            {isLoading ? copy.searching : copy.search}
          </button>
        </div>

        <div style={{ display: 'flex', justifyContent: 'flex-end', flexShrink: 0 }}>
          <button
            type="button"
            style={{
              ...treePrimaryButtonStyle,
              fontSize: '12px',
              height: 28,
              ...(isLoading ? treeButtonDisabledStyle : undefined),
            }}
            onClick={handleRecommendAgain}
            disabled={isLoading}
          >
            {copy.recommendAgain}
          </button>
        </div>

        <div style={{ ...recommendationListStyle, flex: 1, overflowY: 'auto', minHeight: 0 }}>
          {candidates.length === 0 && (
            <RecommendationStatusMessage
              isLoading={isLoading}
              message={
                pointError
                  ? copy.pointErrorEmpty
                  : message ?? copy.empty
              }
              copy={copy}
            />
          )}
          {paged.map((candidate) => (
            <div key={candidate.id} style={recommendationItemStyle}>
              <div style={{ minWidth: 0 }}>
                <p style={recommendationItemTitleStyle} title={candidate.title}>{candidate.title}</p>
                <div style={recommendationItemMetaStyle}>
                  {candidate.content_type}
                  {candidate.content_id ? ` · ${commonCopy.internalContent}` : ''}
                  {candidate.author ? ` · ${candidate.author}` : ''}
                  {candidate.url ? ` · ${candidate.url}` : ''}
                </div>
              </div>
              <div style={{ display: 'flex', gap: 4, flexShrink: 0 }}>
                {(candidate.url ?? candidate.canonical_url) && (
                  <button
                    type="button"
                    style={{ ...treePrimaryButtonStyle, padding: '0 8px' }}
                    title={commonCopy.openLink}
                    onClick={() => window.open(candidate.url ?? candidate.canonical_url ?? '', '_blank', 'noopener,noreferrer')}
                  >
                    🔗
                  </button>
                )}
                <button
                  type="button"
                  style={treePrimaryButtonStyle}
                  onClick={() => onApply(candidate)}
                  disabled={!candidate.content_id && !candidate.url && !candidate.canonical_url}
                >
                  {commonCopy.apply}
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
      </div>
    </ModalPortal>
  )
}

// ── 추천 상태 메시지 헬퍼 (ContentRecommendationModal 내부 + ExplorerAddModals에서 재사용) ──

export function RecommendationStatusMessage({
  isLoading,
  message,
  copy,
}: {
  isLoading: boolean
  message: string
  copy?: DashboardCourseDraftCopy['tree']['editPanel']['modals']['recommendation']
}) {
  const resolvedCopy = copy ?? defaultModalCopy.recommendation
  if (isLoading) {
    return (
      <div style={{ ...recommendationMessageStyle, ...recommendationLoadingMessageStyle }}>
        <style>{`
          @keyframes explorerRecommendationSpin {
            to { transform: rotate(360deg); }
          }
        `}</style>
        <span>{resolvedCopy.loading}</span>
        <span style={recommendationSpinnerStyle} aria-hidden="true" />
      </div>
    )
  }

  return <div style={recommendationMessageStyle}>{message}</div>
}

const recommendationLoadingMessageStyle: CSSProperties = {
  display: 'inline-flex',
  alignItems: 'center',
  gap: 7,
}

const recommendationSpinnerStyle: CSSProperties = {
  width: 12,
  height: 12,
  borderRadius: 999,
  border: '2px solid rgba(122, 90, 36, 0.18)',
  borderTopColor: '#7A5A24',
  animation: 'explorerRecommendationSpin 0.8s linear infinite',
}
