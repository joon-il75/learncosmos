'use client'

import LumiModalShell from '@/components/common/LumiModalShell'
import ModalPortal from '@/components/common/ModalPortal'
import { getDashboardCourseDraftCopy, type DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'
import {
  treeButtonDisabledStyle,
  addItemOverlayStyle,
  addItemModalBoxStyle,
  addExplorationModalBoxStyle,
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
} from './explorerTreeStyles'

const defaultModalCopy = getDashboardCourseDraftCopy('ko').tree.editPanel.modals

// ── 지점 수정 모달 ────────────────────────────────────────────────────────────

export interface NodeEditModalProps {
  mode: 'exploration' | 'research'
  titleInput: string
  urlInput: string
  urlCheckStatus: 'idle' | 'checking' | 'ok' | 'error'
  urlCheckMsg: string
  parentRegionName?: string | null
  parentSubRegionName?: string | null
  onTitleChange: (value: string) => void
  onUrlChange: (value: string) => void
  onCheckUrl: () => void
  onOpenRecommendation: () => void
  onSubmit: () => void
  onClose: () => void
  isMutating: boolean
  copy?: DashboardCourseDraftCopy['tree']['editPanel']['modals']
}

export function NodeEditModal({
  mode,
  titleInput,
  urlInput,
  urlCheckStatus,
  urlCheckMsg,
  parentRegionName,
  parentSubRegionName,
  onTitleChange,
  onUrlChange,
  onCheckUrl,
  onOpenRecommendation,
  onSubmit,
  onClose,
  isMutating,
  copy = defaultModalCopy,
}: NodeEditModalProps) {
  const checkBtnColor =
    urlCheckStatus === 'ok'
      ? { color: '#136857', borderColor: 'rgba(20,140,120,0.4)', background: 'rgba(20,140,120,0.1)' }
      : urlCheckStatus === 'error'
      ? { color: '#9B2C2C', borderColor: 'rgba(155,44,44,0.3)', background: 'rgba(155,44,44,0.06)' }
      : undefined
  const canSubmit = Boolean(titleInput.trim()) && !isMutating && (mode === 'research' || Boolean(urlInput.trim()))
  const parentPath = [parentRegionName, parentSubRegionName].filter(Boolean).join(' > ')

  return (
    <ModalPortal overlayStyle={addItemOverlayStyle} onMouseDown={onClose}>
      <div
        role="dialog"
        aria-modal="true"
        aria-label={mode === 'exploration' ? copy.nodeEdit.explorationAria : copy.nodeEdit.researchAria}
        style={addExplorationModalBoxStyle}
        onMouseDown={(event) => event.stopPropagation()}
      >
        <div style={addItemModalHeaderStyle}>
          <h3 style={addItemModalTitleStyle}>{mode === 'exploration' ? copy.nodeEdit.explorationTitle : copy.nodeEdit.researchTitle}</h3>
          <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
            {copy.close}
          </button>
        </div>

        {parentPath && (
          <div style={{ ...parentContextInfoStyle, marginTop: -4 }}>
            {copy.parentPath(parentPath)}
          </div>
        )}

        <div>
          <label style={addItemModalLabelStyle}>{copy.nodeEdit.titleLabel}</label>
          <input
            style={addItemModalInputStyle}
            value={titleInput}
            onChange={(event) => onTitleChange(event.target.value)}
            placeholder={mode === 'exploration' ? copy.nodeEdit.explorationPlaceholder : copy.nodeEdit.researchPlaceholder}
            disabled={isMutating}
            autoFocus
          />
        </div>

        {mode === 'exploration' && (
          <>
            <div>
              <label style={addItemModalLabelStyle}>{copy.nodeEdit.linkLabel}</label>
              <div style={addItemModalRowStyle}>
                <input
                  style={{ ...addItemModalInputStyle, flex: 1 }}
                  value={urlInput}
                  onChange={(event) => onUrlChange(event.target.value)}
                  placeholder="https://..."
                  disabled={isMutating}
                />
                {urlInput.trim() && (
                  <button
                    type="button"
                    style={{ ...addItemModalCheckBtnStyle, padding: '0 10px' }}
                    onClick={() => window.open(urlInput.trim(), '_blank', 'noopener,noreferrer')}
                    title={copy.openLink}
                  >
                    🔗
                  </button>
                )}
                <button
                  type="button"
                  style={{ ...addItemModalCheckBtnStyle, ...checkBtnColor }}
                  onClick={onCheckUrl}
                  disabled={!urlInput.trim() || urlCheckStatus === 'checking' || isMutating}
                >
                  {urlCheckStatus === 'checking' ? '...' : urlCheckStatus === 'ok' ? '✓' : urlCheckStatus === 'error' ? '✗' : copy.check}
                </button>
                <button
                  type="button"
                  style={addItemModalCheckBtnStyle}
                  onClick={onOpenRecommendation}
                  disabled={isMutating}
                >
                  {copy.nodeEdit.recommend}
                </button>
              </div>
              {urlCheckMsg && (
                <div style={{ fontSize: '11px', marginTop: 6, color: urlCheckStatus === 'ok' ? '#136857' : '#9B2C2C' }}>
                  {urlCheckMsg}
                </div>
              )}
            </div>
          </>
        )}

        <div style={addItemModalRowStyle}>
          <button
            type="button"
            style={{
              ...addItemModalPrimaryBtnStyle,
              ...(!canSubmit ? treeButtonDisabledStyle : undefined),
            }}
            onClick={onSubmit}
            disabled={!canSubmit}
          >
            {copy.edit}
          </button>
          <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
            {copy.cancel}
          </button>
        </div>
      </div>
    </ModalPortal>
  )
}

// ── 삭제 확인 모달 ────────────────────────────────────────────────────────────

export interface DeleteConfirmModalProps {
  mode: 'confirm' | 'alert'
  action?: 'delete' | 'activate'
  title: string
  message: string
  requireTitleInput: boolean
  confirmInput: string
  onConfirmInputChange: (value: string) => void
  onConfirm: () => void
  onClose: () => void
  isMutating: boolean
  copy?: DashboardCourseDraftCopy['tree']['editPanel']['modals']
}

export function DeleteConfirmModal({
  mode,
  action = 'delete',
  title,
  message,
  requireTitleInput,
  confirmInput,
  onConfirmInputChange,
  onConfirm,
  onClose,
  isMutating,
  copy = defaultModalCopy,
}: DeleteConfirmModalProps) {
  const canConfirm = mode === 'confirm' && !isMutating && (!requireTitleInput || confirmInput.trim() === title.trim())
  const isActivate = action === 'activate'

  if (mode === 'alert') {
    return (
      <LumiModalShell
        ariaLabel={copy.deleteConfirm.alertAria}
        eyebrow={copy.deleteConfirm.warningEyebrow}
        lumiState="surprise"
        tone="alert"
        title={copy.deleteConfirm.alertTitle}
        message={
          <>
            {message}
            {'\n'}{copy.deleteConfirm.alertSuffix}
          </>
        }
        onClose={onClose}
        actions={
          <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
            {copy.apply}
          </button>
        }
      />
    )
  }

  return (
    <ModalPortal overlayStyle={addItemOverlayStyle} onMouseDown={onClose}>
      <div
        role="dialog"
        aria-modal="true"
        aria-label={isActivate ? copy.deleteConfirm.activateAria : copy.deleteConfirm.deleteAria}
        style={addItemModalBoxStyle}
        onMouseDown={(event) => event.stopPropagation()}
      >
        <div style={addItemModalHeaderStyle}>
          <h3 style={addItemModalTitleStyle}>{isActivate ? copy.deleteConfirm.activateTitle : copy.deleteConfirm.deleteTitle}</h3>
          <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
            {copy.close}
          </button>
        </div>
        <div style={addItemModalErrorStyle}>{message}</div>
        <div>
          <label style={addItemModalLabelStyle}>{copy.deleteConfirm.targetLabel}</label>
          <div style={addItemModalInfoStyle}>{title}</div>
        </div>
        {requireTitleInput && (
          <div>
            <label style={addItemModalLabelStyle}>{copy.deleteConfirm.courseConfirmLabel}</label>
            <input
              style={addItemModalInputStyle}
              value={confirmInput}
              onChange={(event) => onConfirmInputChange(event.target.value)}
              placeholder={title}
              disabled={isMutating}
              autoFocus
              onKeyDown={(event) => { if (event.key === 'Enter' && canConfirm) onConfirm() }}
            />
            <div style={{ ...addItemModalInfoStyle, marginTop: 4 }}>
              {copy.deleteConfirm.exactTitleHelp(isActivate ? 'activate' : 'delete')}
            </div>
          </div>
        )}
        <div style={addItemModalRowStyle}>
          <button
            type="button"
            style={{
              ...addItemModalPrimaryBtnStyle,
              ...(isActivate
                ? {
                    background: 'linear-gradient(180deg, rgba(25, 122, 88, 0.96), rgba(14, 84, 65, 0.98))',
                    borderColor: 'rgba(164, 240, 205, 0.72)',
                  }
                : {
                    background: 'linear-gradient(180deg, rgba(185, 58, 42, 0.94), rgba(128, 32, 24, 0.98))',
                    borderColor: 'rgba(255, 170, 150, 0.72)',
                  }),
              ...(!canConfirm ? treeButtonDisabledStyle : undefined),
            }}
            onClick={onConfirm}
            disabled={!canConfirm}
          >
            {isActivate ? copy.deleteConfirm.activate : copy.deleteConfirm.delete}
          </button>
          <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
            {copy.cancel}
          </button>
        </div>
      </div>
    </ModalPortal>
  )
}
