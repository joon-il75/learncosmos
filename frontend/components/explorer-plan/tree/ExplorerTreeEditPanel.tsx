'use client'

import { useState, type CSSProperties, type FormEvent } from 'react'
import ModalPortal from '@/components/common/ModalPortal'
import type { ExplorerEditFormState } from './useExplorerEditForm'
import type { DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'
import {
  treePrimaryButtonStyle,
  treeButtonDisabledStyle,
  editPanelStyle,
  editFormAreaStyle,
  editSectionLabelStyle,
  editMutationMsgStyle,
  pendingCreateNoticeStyle,
  pendingDeleteNoticeStyle,
  addNodeButtonStyle,
  addSubRegionButtonStyle,
  addItemOverlayStyle,
  addItemModalBoxStyle,
  addItemModalHeaderStyle,
  addItemModalTitleStyle,
  addItemModalLabelStyle,
  addItemModalInputStyle,
  addItemModalRowStyle,
  addItemModalPrimaryBtnStyle,
  addItemModalCancelBtnStyle,
  addItemModalInfoStyle,
  parentContextInfoStyle,
} from './explorerTreeStyles'
import {
  MoveButtons,
  NodeEditModal,
  DeleteConfirmModal,
  AddChildObjectModal,
  AddExplorationNodeModal,
  AddUrlNodeModal,
  ContentRecommendationModal,
} from './ExplorerEditModals'

interface Props {
  editForm: ExplorerEditFormState
  journalLimitedEdit?: boolean
  copy: DashboardCourseDraftCopy['tree']['editPanel']
}

const secondaryActionButtonStyle: CSSProperties = {
  ...addNodeButtonStyle,
  border: '1px solid rgba(39, 94, 77, 0.32)',
  background: 'linear-gradient(180deg, rgba(232, 247, 241, 0.94), rgba(204, 232, 220, 0.96))',
  color: '#1E5A48',
}

const journalCompactPanelStyle: CSSProperties = {
  ...editPanelStyle,
  marginTop: 'auto',
  marginBottom: 2,
  marginLeft: 6,
  marginRight: 6,
  paddingRight: 6,
  maxHeight: 46,
  borderRadius: 6,
}

const journalCompactFormAreaStyle: CSSProperties = {
  ...editFormAreaStyle,
  padding: '4px 5px',
  gap: 0,
  overflow: 'hidden',
}

const journalCompactButtonStyle: CSSProperties = {
  ...secondaryActionButtonStyle,
  ...addSubRegionButtonStyle,
  height: 26,
  minHeight: 26,
  padding: '0 9px',
  fontSize: 11,
  borderRadius: 7,
}

export function ExplorerTreeEditPanel({ editForm, journalLimitedEdit = false, copy }: Props) {
  const [courseRenameOpen, setCourseRenameOpen] = useState(false)
  const [regionRenameOpen, setRegionRenameOpen] = useState(false)
  const [regionAddOpen, setRegionAddOpen] = useState(false)
  const {
    isCourseSelected,
    isRegionSelected,
    isSubRegionSelected,
    isNodeSelected,
    selectedNode,
    selectedRegionAgg,
    selectedSubAgg,
    selectedNodeRegionAgg,
    selectedNodeSubAgg,
    isSelectedPendingCreated,
    isSelectedPendingDeleted,
    canAddSubRegion,
    canMoveUp, canMoveDown,
    courseTitleInput, setCourseTitleInput,
    newRegionNameInput, setNewRegionNameInput,
    regionNameInput, setRegionNameInput,
    subRegionNameInput, setSubRegionNameInput,
    newSubRegionNameInput, setNewSubRegionNameInput,
    newNodeTitleInput, setNewNodeTitleInput,
    newNodeUrlInput, setNewNodeUrlInput,
    nodeTitleInput, setNodeTitleInput,
    nodeUrlInput, setNodeUrlInput,
    urlCheckStatus, setUrlCheckStatus,
    urlCheckMsg,
    newUrlCheckStatus,
    newUrlCheckMsg,
    isRecommendationOpen,
    recommendationQuery, setRecommendationQuery,
    recommendationCandidates,
    isLoadingRecommendations,
    recommendationMessage,
    recommendationPointError,
    deleteConfirm,
    setDeleteConfirmInput,
    addModalType, setAddModalType,
    handleAddRegion,
    handleAddSubRegion,
    handleAddNode,
    handleCheckUrl,
    handleOpenRecommendationModal,
    handleCloseRecommendationModal,
    handleSearchRecommendations,
    handleApplyRecommendation,
    editNodeModalType,
    handleOpenNodeEditModal,
    handleCloseNodeEditModal,
    handleSubmitNodeEdit,
    handleRenameCourse,
    handleRenameRegion,
    handleConfirmDelete,
    handleCloseDeleteConfirm,
    handleMoveUp,
    handleMoveDown,
    isMutating,
    mutationMessage,
  } = editForm

  if (!isCourseSelected && !isRegionSelected && !isSubRegionSelected && !isNodeSelected) {
    return null
  }

  const isEditDisabled = isMutating || isSelectedPendingDeleted
  const panelStyle = journalLimitedEdit ? journalCompactPanelStyle : editPanelStyle
  const formAreaStyle = journalLimitedEdit ? journalCompactFormAreaStyle : editFormAreaStyle
  const modalCopy = copy.modals
  const renameModalTitle = isSubRegionSelected ? modalCopy.regionRename.subregionTitle : modalCopy.regionRename.regionTitle
  const renameModalLabel = isSubRegionSelected ? modalCopy.regionRename.subregionLabel : modalCopy.regionRename.regionLabel
  const renameParentRegionName = selectedRegionAgg?.region.name ?? null
  const renameParentSubRegionName = isSubRegionSelected ? selectedSubAgg?.subregion.name ?? null : null

  return (
    <div style={panelStyle}>
      <div style={formAreaStyle}>
        {isSelectedPendingDeleted && (
          <div style={pendingDeleteNoticeStyle}>
            {copy.pendingDeleteNotice}
          </div>
        )}
        {isSelectedPendingCreated && !isSelectedPendingDeleted && (
          <div style={pendingCreateNoticeStyle}>
            {copy.pendingCreateNotice}
          </div>
        )}

        {/* 코스 선택 */}
        {isCourseSelected && !journalLimitedEdit && (
          <div>
            <div style={editSectionLabelStyle}>{copy.action}</div>
            <div style={{ display: 'flex', gap: 5, flexWrap: 'wrap', alignItems: 'center' }}>
              <button
                type="button"
                style={{
                  ...secondaryActionButtonStyle,
                  ...addSubRegionButtonStyle,
                  ...(isEditDisabled ? treeButtonDisabledStyle : undefined),
                }}
                onClick={() => setCourseRenameOpen(true)}
                disabled={isEditDisabled}
              >
                {copy.renamePlanet}
              </button>
              <button
                type="button"
                style={{
                  ...secondaryActionButtonStyle,
                  ...addSubRegionButtonStyle,
                  ...(isEditDisabled ? treeButtonDisabledStyle : undefined),
                }}
                onClick={() => setRegionAddOpen(true)}
                disabled={isEditDisabled}
              >
                {copy.addRegion}
              </button>
            </div>
          </div>
        )}

        {/* 지역 선택 */}
        {isRegionSelected && (
          <>
            <div>
              {!journalLimitedEdit && <div style={editSectionLabelStyle}>{copy.action}</div>}
              <div style={{ display: 'grid', gap: 6 }}>
                {!journalLimitedEdit && (
                  <div style={{ display: 'flex', gap: 5, flexWrap: 'wrap', alignItems: 'center' }}>
                    <MoveButtons
                      canMoveUp={canMoveUp}
                      canMoveDown={canMoveDown}
                      isMutating={isEditDisabled}
                      onMoveUp={handleMoveUp}
                      onMoveDown={handleMoveDown}
                      moveUpLabel={copy.moveUp}
                      moveDownLabel={copy.moveDown}
                    />
                  </div>
                )}
                <div style={{ display: 'flex', gap: 5, flexWrap: 'wrap', alignItems: 'center' }}>
                  {!journalLimitedEdit && (
                    <button
                      type="button"
                      style={{
                        ...secondaryActionButtonStyle,
                        ...(isEditDisabled ? treeButtonDisabledStyle : undefined),
                      }}
                      onClick={() => setRegionRenameOpen(true)}
                      disabled={isEditDisabled}
                    >
                      {copy.renameRegion}
                    </button>
                  )}
                  <button
                    type="button"
                    style={{
                      ...(journalLimitedEdit ? journalCompactButtonStyle : { ...secondaryActionButtonStyle, ...addSubRegionButtonStyle }),
                      ...(isEditDisabled ? treeButtonDisabledStyle : undefined),
                    }}
                    onClick={() => setAddModalType('child-object')}
                    disabled={isEditDisabled}
                  >
                    {copy.addChild}
                  </button>
                </div>
              </div>
            </div>
          </>
        )}

        {/* 서브지역 선택 */}
        {isSubRegionSelected && (
          <>
            <div>
              {!journalLimitedEdit && <div style={editSectionLabelStyle}>{copy.action}</div>}
              <div style={{ display: 'flex', gap: 5, flexWrap: 'wrap', alignItems: 'center' }}>
                {!journalLimitedEdit && (
                  <button
                    type="button"
                    style={{
                      ...secondaryActionButtonStyle,
                      ...(isEditDisabled ? treeButtonDisabledStyle : undefined),
                    }}
                    onClick={() => setRegionRenameOpen(true)}
                    disabled={isEditDisabled}
                  >
                    {copy.renameSubRegion}
                  </button>
                )}
                <button
                  type="button"
                  style={{
                    ...(journalLimitedEdit ? journalCompactButtonStyle : { ...secondaryActionButtonStyle, ...addSubRegionButtonStyle }),
                    ...(isEditDisabled ? treeButtonDisabledStyle : undefined),
                  }}
                  onClick={() => setAddModalType('child-object')}
                  disabled={isEditDisabled}
                >
                  {copy.addChild}
                </button>
              </div>
            </div>
          </>
        )}

        {/* 지점 선택 */}
        {isNodeSelected && selectedNode && !journalLimitedEdit && (
          <>
            <div>
              <div style={editSectionLabelStyle}>{copy.action}</div>
              <div style={{ display: 'flex', gap: 5, flexWrap: 'wrap', alignItems: 'center' }}>
                <MoveButtons
                  canMoveUp={canMoveUp}
                  canMoveDown={canMoveDown}
                  isMutating={isEditDisabled}
                  onMoveUp={handleMoveUp}
                  onMoveDown={handleMoveDown}
                  moveUpLabel={copy.moveUp}
                  moveDownLabel={copy.moveDown}
                />
                <button
                  type="button"
                  style={{
                    ...secondaryActionButtonStyle,
                    ...(isEditDisabled ? treeButtonDisabledStyle : undefined),
                  }}
                  onClick={handleOpenNodeEditModal}
                  disabled={isEditDisabled}
                >
                  {copy.editNode}
                </button>
              </div>
            </div>
          </>
        )}
      </div>

      {mutationMessage && <div style={editMutationMsgStyle}>{mutationMessage}</div>}

      {courseRenameOpen && (
        <CourseRenameModal
          value={courseTitleInput}
          onChange={setCourseTitleInput}
          onSubmit={async () => {
            await handleRenameCourse()
            setCourseRenameOpen(false)
          }}
          onClose={() => setCourseRenameOpen(false)}
          isMutating={isMutating}
          copy={modalCopy}
        />
      )}

      {regionAddOpen && (
        <RegionAddModal
          value={newRegionNameInput}
          onChange={setNewRegionNameInput}
          onSubmit={async () => {
            await handleAddRegion()
            setRegionAddOpen(false)
          }}
          onClose={() => setRegionAddOpen(false)}
          isMutating={isMutating}
          copy={modalCopy}
        />
      )}

      {regionRenameOpen && (
        <RegionRenameModal
          title={renameModalTitle}
          label={renameModalLabel}
          parentRegionName={renameParentRegionName}
          parentSubRegionName={renameParentSubRegionName}
          value={isSubRegionSelected ? subRegionNameInput : regionNameInput}
          onChange={isSubRegionSelected ? setSubRegionNameInput : setRegionNameInput}
          onSubmit={async () => {
            await handleRenameRegion()
            setRegionRenameOpen(false)
          }}
          onClose={() => setRegionRenameOpen(false)}
          isMutating={isMutating}
          copy={modalCopy}
        />
      )}

      {addModalType === 'child-object' && (
        <AddChildObjectModal
          canAddSubRegion={canAddSubRegion}
          isRegionSelected={isRegionSelected}
          parentRegionName={selectedRegionAgg?.region.name ?? null}
          parentSubRegionName={selectedSubAgg?.subregion.name ?? null}
          subRegionNameInput={newSubRegionNameInput}
          nodeTitleInput={newNodeTitleInput}
          nodeUrlInput={newNodeUrlInput}
          urlCheckStatus={newUrlCheckStatus}
          urlCheckMsg={newUrlCheckMsg}
          recommendationQuery={recommendationQuery}
          recommendationCandidates={recommendationCandidates}
          isLoadingRecommendations={isLoadingRecommendations}
          recommendationMessage={recommendationMessage}
          recommendationPointError={recommendationPointError}
          onSubRegionNameChange={setNewSubRegionNameInput}
          onNodeTitleChange={setNewNodeTitleInput}
          onNodeUrlChange={setNewNodeUrlInput}
          onRecommendationQueryChange={setRecommendationQuery}
          onAddSubRegion={() => void handleAddSubRegion()}
          onCheckUrl={() => void handleCheckUrl('new')}
          onSearchRecommendations={() => void handleSearchRecommendations()}
          onApplyRecommendation={(candidate) => void handleApplyRecommendation(candidate)}
          onAddExploration={() => void handleAddNode('exploration')}
          onAddResearch={() => void handleAddNode('research')}
          onClose={() => setAddModalType(null)}
          isMutating={isMutating}
          copy={modalCopy}
        />
      )}

      {addModalType === 'exploration-choice' && (
        <AddExplorationNodeModal
          query={recommendationQuery}
          candidates={recommendationCandidates}
          isLoading={isLoadingRecommendations}
          message={recommendationMessage}
          pointError={recommendationPointError}
          titleInput={newNodeTitleInput}
          urlInput={newNodeUrlInput}
          urlCheckStatus={newUrlCheckStatus}
          urlCheckMsg={newUrlCheckMsg}
          onQueryChange={setRecommendationQuery}
          onSearch={() => void handleSearchRecommendations()}
          onApply={(candidate) => void handleApplyRecommendation(candidate)}
          onTitleChange={setNewNodeTitleInput}
          onUrlChange={setNewNodeUrlInput}
          onCheckUrl={() => void handleCheckUrl('new')}
          onAdd={() => void handleAddNode()}
          onClose={() => setAddModalType(null)}
          isMutating={isMutating}
          copy={modalCopy}
        />
      )}

      {addModalType === 'exploration-url' && (
        <AddUrlNodeModal
          mode="exploration"
          titleInput={newNodeTitleInput}
          urlInput={newNodeUrlInput}
          urlCheckStatus={newUrlCheckStatus}
          urlCheckMsg={newUrlCheckMsg}
          onTitleChange={setNewNodeTitleInput}
          onUrlChange={setNewNodeUrlInput}
          onCheckUrl={() => void handleCheckUrl('new')}
          onAdd={() => void handleAddNode()}
          onClose={() => setAddModalType(null)}
          isMutating={isMutating}
          copy={modalCopy}
        />
      )}

      {addModalType === 'research' && (
        <AddUrlNodeModal
          mode="research"
          titleInput={newNodeTitleInput}
          urlInput=""
          urlCheckStatus="idle"
          urlCheckMsg=""
          onTitleChange={setNewNodeTitleInput}
          onUrlChange={setNewNodeUrlInput}
          onCheckUrl={() => void handleCheckUrl('new')}
          onAdd={() => void handleAddNode()}
          onClose={() => setAddModalType(null)}
          isMutating={isMutating}
          copy={modalCopy}
        />
      )}

      {isRecommendationOpen && (
        <ContentRecommendationModal
          query={recommendationQuery}
          candidates={recommendationCandidates}
          isLoading={isLoadingRecommendations}
          message={recommendationMessage}
          pointError={recommendationPointError}
          onQueryChange={setRecommendationQuery}
          onSearch={() => void handleSearchRecommendations()}
          onRecommendAgain={() => void handleSearchRecommendations()}
          onClose={handleCloseRecommendationModal}
          onApply={(c) => void handleApplyRecommendation(c)}
          copy={modalCopy.recommendation}
          commonCopy={modalCopy}
        />
      )}

      {editNodeModalType && selectedNode && (
        <NodeEditModal
          mode={editNodeModalType}
          titleInput={nodeTitleInput}
          urlInput={nodeUrlInput}
          urlCheckStatus={urlCheckStatus}
          urlCheckMsg={urlCheckMsg}
          parentRegionName={selectedNodeRegionAgg?.region.name ?? null}
          parentSubRegionName={selectedNodeSubAgg?.subregion.name ?? null}
          onTitleChange={setNodeTitleInput}
          onUrlChange={(value) => { setNodeUrlInput(value); setUrlCheckStatus('idle') }}
          onCheckUrl={() => void handleCheckUrl('selected')}
          onOpenRecommendation={() => void handleOpenRecommendationModal('selected')}
          onSubmit={() => void handleSubmitNodeEdit()}
          onClose={handleCloseNodeEditModal}
          isMutating={isMutating}
          copy={modalCopy}
        />
      )}

      {deleteConfirm && (
        <DeleteConfirmModal
          mode={deleteConfirm.mode}
          action={deleteConfirm.action}
          title={deleteConfirm.title}
          message={deleteConfirm.message}
          requireTitleInput={deleteConfirm.requireTitleInput}
          confirmInput={deleteConfirm.confirmInput}
          onConfirmInputChange={setDeleteConfirmInput}
          onConfirm={() => void handleConfirmDelete()}
          onClose={handleCloseDeleteConfirm}
          isMutating={isMutating}
          copy={modalCopy}
        />
      )}
    </div>
  )
}

interface SimpleTextModalProps {
  value: string
  onChange: (value: string) => void
  onSubmit: () => Promise<void>
  onClose: () => void
  isMutating: boolean
  copy: DashboardCourseDraftCopy['tree']['editPanel']['modals']
}

interface RegionRenameModalProps extends SimpleTextModalProps {
  title: string
  label: string
  parentRegionName?: string | null
  parentSubRegionName?: string | null
}

function CourseRenameModal({ value, onChange, onSubmit, onClose, isMutating, copy }: SimpleTextModalProps) {
  const trimmed = value.trim()
  const disabled = isMutating || trimmed.length === 0
  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (disabled) return
    void onSubmit()
  }

  return (
    <ModalPortal overlayStyle={addItemOverlayStyle} onMouseDown={onClose}>
      <form
        role="dialog"
        aria-modal="true"
        aria-label={copy.courseRename.aria}
        style={addItemModalBoxStyle}
        onMouseDown={(event) => event.stopPropagation()}
        onSubmit={handleSubmit}
      >
        <div style={addItemModalHeaderStyle}>
          <h3 style={addItemModalTitleStyle}>{copy.courseRename.title}</h3>
          <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
            {copy.close}
          </button>
        </div>
        <div>
          <label style={addItemModalLabelStyle}>{copy.courseRename.label}</label>
          <input
            style={addItemModalInputStyle}
            value={value}
            onChange={(event) => onChange(event.target.value)}
            placeholder={copy.courseRename.placeholder}
            maxLength={160}
            disabled={isMutating}
            autoFocus
          />
          <div style={{ ...addItemModalInfoStyle, marginTop: 6 }}>
            {copy.courseRename.help}
          </div>
        </div>
        <div style={addItemModalRowStyle}>
          <button
            type="submit"
            style={{
              ...addItemModalPrimaryBtnStyle,
              ...(disabled ? treeButtonDisabledStyle : undefined),
            }}
            disabled={disabled}
          >
            {copy.change}
          </button>
          <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
            {copy.cancel}
          </button>
        </div>
      </form>
    </ModalPortal>
  )
}

function RegionRenameModal({
  title,
  label,
  parentRegionName,
  parentSubRegionName,
  value,
  onChange,
  onSubmit,
  onClose,
  isMutating,
  copy,
}: RegionRenameModalProps) {
  const trimmed = value.trim()
  const disabled = isMutating || trimmed.length === 0
  const parentPath = [parentRegionName, parentSubRegionName].filter(Boolean).join(' > ')
  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (disabled) return
    void onSubmit()
  }

  return (
    <ModalPortal overlayStyle={addItemOverlayStyle} onMouseDown={onClose}>
      <form
        role="dialog"
        aria-modal="true"
        aria-label={copy.regionRename.aria}
        style={addItemModalBoxStyle}
        onMouseDown={(event) => event.stopPropagation()}
        onSubmit={handleSubmit}
      >
        <div style={addItemModalHeaderStyle}>
          <h3 style={addItemModalTitleStyle}>{title}</h3>
          <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
            {copy.close}
          </button>
        </div>
        {parentPath && (
          <div style={{ ...parentContextInfoStyle, marginTop: -4, marginBottom: 8 }}>
            {copy.parentPath(parentPath)}
          </div>
        )}
        <div>
          <label style={addItemModalLabelStyle}>{label}</label>
          <input
            style={addItemModalInputStyle}
            value={value}
            onChange={(event) => onChange(event.target.value)}
            placeholder={copy.regionRename.placeholder(label)}
            maxLength={100}
            disabled={isMutating}
            autoFocus
          />
          <div style={{ ...addItemModalInfoStyle, marginTop: 6 }}>
            {copy.regionRename.help}
          </div>
        </div>
        <div style={addItemModalRowStyle}>
          <button
            type="submit"
            style={{
              ...addItemModalPrimaryBtnStyle,
              ...(disabled ? treeButtonDisabledStyle : undefined),
            }}
            disabled={disabled}
          >
            {copy.edit}
          </button>
          <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
            {copy.cancel}
          </button>
        </div>
      </form>
    </ModalPortal>
  )
}

function RegionAddModal({ value, onChange, onSubmit, onClose, isMutating, copy }: SimpleTextModalProps) {
  const trimmed = value.trim()
  const disabled = isMutating || trimmed.length === 0
  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (disabled) return
    void onSubmit()
  }

  return (
    <ModalPortal overlayStyle={addItemOverlayStyle} onMouseDown={onClose}>
      <form
        role="dialog"
        aria-modal="true"
        aria-label={copy.regionAdd.aria}
        style={addItemModalBoxStyle}
        onMouseDown={(event) => event.stopPropagation()}
        onSubmit={handleSubmit}
      >
        <div style={addItemModalHeaderStyle}>
          <h3 style={addItemModalTitleStyle}>{copy.regionAdd.title}</h3>
          <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
            {copy.close}
          </button>
        </div>
        <div>
          <label style={addItemModalLabelStyle}>{copy.regionAdd.label}</label>
          <div style={{ ...addItemModalInfoStyle, marginBottom: 8, color: '#5D3D17', fontWeight: 650 }}>
            {copy.regionAdd.description}
          </div>
          <input
            style={addItemModalInputStyle}
            value={value}
            onChange={(event) => onChange(event.target.value)}
            placeholder={copy.regionAdd.placeholder}
            maxLength={100}
            disabled={isMutating}
            autoFocus
          />
          <div style={{ ...addItemModalInfoStyle, marginTop: 6 }}>
            {copy.regionAdd.help}
          </div>
        </div>
        <div style={addItemModalRowStyle}>
          <button
            type="submit"
            style={{
              ...addItemModalPrimaryBtnStyle,
              ...(disabled ? treeButtonDisabledStyle : undefined),
            }}
            disabled={disabled}
          >
            {copy.regionAdd.addButton}
          </button>
          <button type="button" style={addItemModalCancelBtnStyle} onClick={onClose} disabled={isMutating}>
            {copy.cancel}
          </button>
        </div>
      </form>
    </ModalPortal>
  )
}
