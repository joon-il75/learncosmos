'use client'

import type { CSSProperties } from 'react'
import type { ExplorerEditFormState } from './useExplorerEditForm'
import type { DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'
import {
  treeButtonDisabledStyle,
  saveButtonStyle,
  cancelButtonStyle,
  deleteButtonStyle,
  inactiveToggleButtonStyle,
} from './explorerTreeStyles'

const barStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(2, minmax(0, 1fr))',
  alignItems: 'center',
  gap: 6,
  padding: '10px 8px',
  background: 'linear-gradient(180deg, rgba(86, 50, 17, 0.82) 0%, rgba(50, 28, 10, 0.88) 100%)',
  borderRadius: 10,
  border: '1px solid rgba(196, 142, 54, 0.34)',
  boxShadow: 'inset 0 1px 0 rgba(255, 220, 150, 0.18), 0 6px 14px rgba(28, 14, 4, 0.34)',
  backdropFilter: 'blur(4px)',
}

const msgStyle: CSSProperties = {
  gridColumn: '1 / -1',
  minWidth: 0,
  fontSize: '11px',
  fontWeight: 600,
  color: 'rgba(255, 239, 196, 0.92)',
  textShadow: '0 1px 1px rgba(0, 0, 0, 0.32)',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
  lineHeight: 1.35,
}

const inactiveCourseBarStyle: CSSProperties = {
  ...barStyle,
  gridTemplateColumns: '1fr',
  background: 'linear-gradient(180deg, rgba(39, 58, 82, 0.92) 0%, rgba(18, 30, 48, 0.94) 100%)',
  border: '1px solid rgba(148, 197, 255, 0.46)',
  boxShadow: 'inset 0 1px 0 rgba(226, 239, 255, 0.16), 0 10px 24px rgba(15, 23, 42, 0.30)',
}

const activatePlanetButtonStyle: CSSProperties = {
  ...inactiveToggleButtonStyle,
  height: 38,
  border: '1px solid rgba(116, 255, 194, 0.76)',
  background: 'linear-gradient(180deg, rgba(30, 178, 117, 0.96), rgba(13, 116, 86, 0.98))',
  color: '#F0FFF9',
  fontSize: '13px',
  fontWeight: 900,
  boxShadow: 'inset 0 1px 0 rgba(255,255,255,0.22), 0 0 0 1px rgba(125, 255, 202, 0.14), 0 12px 24px rgba(16, 185, 129, 0.24)',
}

interface Props {
  editForm: ExplorerEditFormState
  onAfterSave?: () => void
  journalLimitedEdit?: boolean
  inactiveCourseMode?: boolean
  copy: DashboardCourseDraftCopy['tree']['actionBar']
}

export function ExplorerTreeActionBar({ editForm, onAfterSave, journalLimitedEdit = false, inactiveCourseMode = false, copy }: Props) {
  const {
    canDelete,
    canDeleteCourse,
    canActivateCourse,
    hasChildren,
    isSelectedInactive,
    isSelectedPendingDeleted,
    isActionBarBlocked,
    showInactiveItems,
    handleSave,
    handleCancel,
    handleDelete,
    handleActivateCourse,
    handleInactiveToggle,
    isMutating,
    mutationMessage,
  } = editForm
  const isDisabled = isMutating || isActionBarBlocked

  if (inactiveCourseMode) {
    return (
      <div style={inactiveCourseBarStyle}>
        <button
          type="button"
          style={{
            ...activatePlanetButtonStyle,
            ...(isDisabled || !canActivateCourse ? treeButtonDisabledStyle : undefined),
          }}
          onClick={() => void handleActivateCourse()}
          disabled={isDisabled || !canActivateCourse}
          title={copy.inactiveActivateTitle}
        >
          {copy.activatePlanet}
        </button>
        {mutationMessage && (
          <span style={msgStyle} title={mutationMessage}>
            {mutationMessage}
          </span>
        )}
      </div>
    )
  }

  return (
    <div style={barStyle}>
      <button
        type="button"
        style={{ ...saveButtonStyle, ...(isDisabled ? treeButtonDisabledStyle : undefined) }}
        onClick={() => {
          void handleSave().then(() => {
            onAfterSave?.()
          })
        }}
        disabled={isDisabled}
      >
        {copy.save}
      </button>
      {!journalLimitedEdit && (
        <button
          type="button"
          style={{ ...cancelButtonStyle, ...(isDisabled ? treeButtonDisabledStyle : undefined) }}
          onClick={handleCancel}
          disabled={isDisabled}
        >
          {copy.cancel}
        </button>
      )}
      {canDelete && !journalLimitedEdit && (
        <button
          type="button"
          style={{
            ...deleteButtonStyle,
            ...(isDisabled ? treeButtonDisabledStyle : undefined),
          }}
          onClick={() => void handleDelete()}
          disabled={isDisabled}
          title={
            isSelectedPendingDeleted
              ? copy.undoDeleteTitle
              : canDeleteCourse
              ? copy.deleteCourseTitle
              : hasChildren
              ? copy.deleteWithChildrenTitle
              : undefined
          }
        >
          {isSelectedPendingDeleted ? copy.undoDelete : copy.delete}
        </button>
      )}
      {canActivateCourse && !journalLimitedEdit && (
        <button
          type="button"
          style={{
            ...inactiveToggleButtonStyle,
            ...(isDisabled ? treeButtonDisabledStyle : undefined),
          }}
          onClick={() => void handleActivateCourse()}
          disabled={isDisabled}
          title={copy.activateTitle}
        >
          {copy.activatePlanet}
        </button>
      )}
      {!journalLimitedEdit && (
        <button
          type="button"
          style={{
            ...inactiveToggleButtonStyle,
            ...(isDisabled || isSelectedPendingDeleted ? treeButtonDisabledStyle : undefined),
          }}
          onClick={() => void handleInactiveToggle()}
          disabled={isDisabled || isSelectedPendingDeleted}
          title={isSelectedInactive ? copy.reactivateTitle : copy.inactiveToggleTitle}
        >
          {isSelectedInactive ? copy.activate : showInactiveItems ? copy.hideInactive : copy.showInactive}
        </button>
      )}
      {mutationMessage && (
        <span style={msgStyle} title={mutationMessage}>
          {mutationMessage}
        </span>
      )}
    </div>
  )
}
