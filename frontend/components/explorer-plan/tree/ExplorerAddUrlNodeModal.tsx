'use client'

import { type CSSProperties } from 'react'
import ModalPortal from '@/components/common/ModalPortal'
import { getDashboardCourseDraftCopy, type DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'
import {
  treeButtonDisabledStyle,
  addItemOverlayStyle,
  addItemModalBoxStyle,
  addItemModalHeaderStyle,
  addItemModalTitleStyle,
  addItemModalLabelStyle,
  addItemModalInputStyle,
  addItemModalRowStyle,
  addItemModalPrimaryBtnStyle,
  addItemModalCancelBtnStyle,
  addItemModalCheckBtnStyle,
  addItemModalInfoStyle,
} from './explorerTreeStyles'

const defaultModalCopy = getDashboardCourseDraftCopy('ko').tree.editPanel.modals

// ── URL / 연구지점 추가 모달 ──

export interface AddUrlNodeModalProps {
  mode: 'exploration' | 'research'
  titleInput: string
  urlInput: string
  urlCheckStatus: 'idle' | 'checking' | 'ok' | 'error'
  urlCheckMsg: string
  onTitleChange: (v: string) => void
  onUrlChange: (v: string) => void
  onCheckUrl: () => void
  onAdd: () => void
  onClose: () => void
  isMutating: boolean
  copy?: DashboardCourseDraftCopy['tree']['editPanel']['modals']
}

export function AddUrlNodeModal({
  mode, titleInput, urlInput, urlCheckStatus, urlCheckMsg,
  onTitleChange, onUrlChange, onCheckUrl, onAdd, onClose, isMutating,
  copy = defaultModalCopy,
}: AddUrlNodeModalProps) {
  const urlPassed = urlCheckStatus === 'ok'
  const canAdd = titleInput.trim() !== '' && (mode === 'research' || urlPassed)

  const checkBtnStyle: CSSProperties = {
    ...addItemModalCheckBtnStyle,
    ...(urlCheckStatus === 'ok'
      ? { color: '#136857', borderColor: 'rgba(20,140,120,0.5)', background: 'rgba(20,140,120,0.10)' }
      : urlCheckStatus === 'error'
      ? { color: '#9B2C2C', borderColor: 'rgba(155,44,44,0.4)', background: 'rgba(155,44,44,0.06)' }
      : undefined),
  }

  const title = mode === 'exploration' ? copy.addUrl.explorationTitle : copy.addUrl.researchTitle

  return (
    <ModalPortal overlayStyle={addItemOverlayStyle} onMouseDown={onClose}>
      <div
        role="dialog"
        aria-modal="true"
        aria-label={title}
        style={addItemModalBoxStyle}
        onMouseDown={(e) => e.stopPropagation()}
      >
        <div style={addItemModalHeaderStyle}>
          <h3 style={addItemModalTitleStyle}>{title}</h3>
          <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose}>{copy.close}</button>
        </div>

        <div>
          <label style={addItemModalLabelStyle}>
            {mode === 'exploration' ? copy.addExploration.titleLabel : copy.addChild.researchTitleLabel}
          </label>
          <input
            style={addItemModalInputStyle}
            value={titleInput}
            onChange={(e) => onTitleChange(e.target.value)}
            placeholder={mode === 'exploration' ? copy.addExploration.titlePlaceholder : copy.addChild.researchTitlePlaceholder}
            disabled={isMutating}
            autoFocus
          />
        </div>

        {mode === 'exploration' && (
          <div>
            <label style={addItemModalLabelStyle}>{copy.addExploration.urlLabel}</label>
            <div style={addItemModalRowStyle}>
              <input
                style={{ ...addItemModalInputStyle, flex: 1 }}
                value={urlInput}
                onChange={(e) => onUrlChange(e.target.value)}
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
        )}

        <div style={addItemModalRowStyle}>
          <button
            type="button"
            style={{ ...addItemModalPrimaryBtnStyle, ...(!canAdd || isMutating ? treeButtonDisabledStyle : undefined) }}
            onClick={onAdd}
            disabled={!canAdd || isMutating}
          >
            {mode === 'exploration' ? copy.addExploration.addButton : copy.addChild.addResearch}
          </button>
          <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
            {copy.cancel}
          </button>
        </div>
      </div>
    </ModalPortal>
  )
}
