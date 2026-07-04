'use client';

import { useEffect, useRef, useState } from 'react';
import LumiModalShell from '@/components/common/LumiModalShell';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../../pointPageUtils';
import type { UsePointLearningResult } from '../usePointLearning';
import {
  formGridStyle, sectionMiniHeaderStyle, blockMetaStyle,
  blockCardStyle, blockHeaderStyle, blockActionRowStyle,
  placeholderCardStyle, emptyCardStyle, formActionRowStyle,
  inputStyle,
  pointShellTitleStyle,
} from '../../pointPageStyles';
import { makeWorkspaceStyles } from './workspaceStyles';
import { ResearchSectionHeader } from './WorkspaceSectionHeader';
import { scrollToLearningRecordSection } from './workspaceTypes';

const RESEARCH_ATTACHMENT_ACCEPT = '.jpg,.jpeg,.png,.webp,.gif,.pdf,.doc,.docx,.hwp,.hwpx,.txt,.xls,.xlsx,.csv,.ppt,.pptx,.rtf,.odt,.md,.zip';
const RECORDS_PER_PAGE = 5;

interface Props {
  copy: PointLearningCopy['workspace']['learningWork']['recordPanels']['attachments'];
  commonCopy: PointLearningCopy['workspace']['learningWork']['recordPanels']['common'];
  learning: UsePointLearningResult;
  themeTokens: PointPageThemeTokens;
  canEditLearningRecords: boolean;
  onEditingChange: (editing: boolean) => void;
}

export default function PointAttachmentsTab({ copy, commonCopy, learning, themeTokens, canEditLearningRecords, onEditingChange }: Props) {
  const {
    planetID, pointID,
    attachmentDrafts, newAttachment, selectedAttachmentFile, setSelectedAttachmentFile,
    isSavingAttachment,
    handleUploadAttachment, handleReplaceAttachmentFile, handleOpenAttachment,
    handleUpdateAttachment, handleDeleteAttachment,
    handleChangeAttachmentDraft, handleChangeNewAttachment,
  } = learning;

  const {
    materialToolbarButtonStyle, materialToolbarPrimaryButtonStyle, materialToolbarDangerButtonStyle,
    listBoardStyle, getListHeaderRowStyle, getListDataRowStyle, listColumnHeaderStyle,
    getListIndexCellStyle, getListTextCellStyle,
  } = makeWorkspaceStyles(themeTokens);

  const workAttachmentInputRef = useRef<HTMLInputElement | null>(null);
  const replacementWorkAttachmentInputRef = useRef<HTMLInputElement | null>(null);

  const [selectedAttachmentID, setSelectedAttachmentID] = useState<string | null>(null);
  const [attachmentMode, setAttachmentMode] = useState<'view' | 'edit'>('view');
  const [attachmentPage, setAttachmentPage] = useState(1);
  const [isAddingAttachment, setIsAddingAttachment] = useState(false);
  const [isConfirmingAttachmentListReturn, setIsConfirmingAttachmentListReturn] = useState(false);
  const [isConfirmingAttachmentAddCancel, setIsConfirmingAttachmentAddCancel] = useState(false);
  const [pendingDeleteWorkAttachment, setPendingDeleteWorkAttachment] = useState<typeof attachmentDrafts[number] | null>(null);
  const [workAttachmentUploadFile, setWorkAttachmentUploadFile] = useState<File | null>(null);
  const [replacementWorkAttachmentFile, setReplacementWorkAttachmentFile] = useState<File | null>(null);
  const [hoveredListRowID, setHoveredListRowID] = useState<string | null>(null);

  const isEditing = Boolean(selectedAttachmentID || isAddingAttachment);
  useEffect(() => { onEditingChange(isEditing); }, [isEditing, onEditingChange]);

  const researchMaterialAttachments = attachmentDrafts.filter((a) => a.sourceContext === 'research_material');
  const workAttachments = attachmentDrafts.filter((a) => a.sourceContext === 'work_attachment');

  const getRecentTime = (r: { updatedAt?: string; createdAt?: string }) => {
    const parsed = Date.parse(r.updatedAt || r.createdAt || '');
    return Number.isFinite(parsed) ? parsed : 0;
  };
  const sortedWorkAttachments = [...workAttachments].sort((a, b) => getRecentTime(b) - getRecentTime(a));
  const attachmentPageCount = Math.max(1, Math.ceil(sortedWorkAttachments.length / RECORDS_PER_PAGE));
  const pagedAttachments = sortedWorkAttachments.slice((attachmentPage - 1) * RECORDS_PER_PAGE, attachmentPage * RECORDS_PER_PAGE);

  const selectedAttachment = selectedAttachmentID ? attachmentDrafts.find((a) => a.id === selectedAttachmentID) ?? null : null;
  const isAttachmentEditing = Boolean(selectedAttachment && attachmentMode === 'edit' && canEditLearningRecords);

  useEffect(() => {
    const maxPage = Math.max(1, Math.ceil(sortedWorkAttachments.length / RECORDS_PER_PAGE));
    setAttachmentPage((p) => Math.min(p, maxPage));
    if (selectedAttachmentID && !attachmentDrafts.some((a) => a.id === selectedAttachmentID)) setSelectedAttachmentID(null);
  }, [attachmentDrafts, selectedAttachmentID, sortedWorkAttachments.length]);

  const getResearchAttachmentLabel = (attachment: typeof attachmentDrafts[number]) =>
    attachment.title || attachment.filePath || attachment.url || copy.detail.attachmentFallback;

  const handleUploadAttachmentAndClose = async () => {
    const uploadFile = workAttachmentUploadFile ?? selectedAttachmentFile;
    const uploaded = await handleUploadAttachment(uploadFile, newAttachment.title.trim() || uploadFile?.name);
    if (uploaded) { setIsAddingAttachment(false); setWorkAttachmentUploadFile(null); scrollToLearningRecordSection(); }
  };

  const handleSaveAttachmentEdit = async (attachment: typeof attachmentDrafts[number]) => {
    const saved = replacementWorkAttachmentFile
      ? await handleReplaceAttachmentFile(attachment, replacementWorkAttachmentFile, attachment.title)
      : await handleUpdateAttachment(attachment);
    if (saved) { setReplacementWorkAttachmentFile(null); setAttachmentMode('view'); scrollToLearningRecordSection(); }
  };

  const handleDeleteAttachmentAndClose = async (attachment: typeof attachmentDrafts[number]) => {
    const deleted = await handleDeleteAttachment(attachment);
    if (deleted) { setSelectedAttachmentID(null); setAttachmentMode('view'); setReplacementWorkAttachmentFile(null); scrollToLearningRecordSection(); }
  };

  const handleConfirmAttachmentListReturn = () => {
    setIsConfirmingAttachmentListReturn(false);
    setSelectedAttachmentID(null);
    setAttachmentMode('view');
    setReplacementWorkAttachmentFile(null);
    scrollToLearningRecordSection();
  };

  const handleConfirmAttachmentAddCancel = () => {
    setIsConfirmingAttachmentAddCancel(false);
    setIsAddingAttachment(false);
    setWorkAttachmentUploadFile(null);
    setSelectedAttachmentFile(null);
    scrollToLearningRecordSection();
  };

  const handleConfirmDeleteWorkAttachment = async () => {
    if (!pendingDeleteWorkAttachment || isSavingAttachment) return;
    const attachment = pendingDeleteWorkAttachment;
    setPendingDeleteWorkAttachment(null);
    await handleDeleteAttachmentAndClose(attachment);
  };

  const selectedWorkAttachmentFile = workAttachmentUploadFile ?? selectedAttachmentFile;

  return (
    <>
      {pendingDeleteWorkAttachment ? (
        <LumiModalShell
          title={copy.detail.deleteModalTitle} eyebrow="Lumi Confirm" lumiState="curious" tone="alert" width={440}
          onClose={() => { if (!isSavingAttachment) setPendingDeleteWorkAttachment(null); }}
          message={copy.detail.deleteModalMessage(getResearchAttachmentLabel(pendingDeleteWorkAttachment))}
          actions={(
            <>
              <button type="button" onClick={() => setPendingDeleteWorkAttachment(null)} disabled={isSavingAttachment} style={{ ...materialToolbarButtonStyle, opacity: isSavingAttachment ? 0.62 : 1 }}>{copy.detail.cancel}</button>
              <button type="button" onClick={() => void handleConfirmDeleteWorkAttachment()} disabled={isSavingAttachment} style={{ ...materialToolbarDangerButtonStyle, opacity: isSavingAttachment ? 0.72 : 1 }}>{isSavingAttachment ? copy.detail.deleting : copy.detail.delete}</button>
            </>
          )}
        />
      ) : null}
      {isConfirmingAttachmentListReturn ? (
        <LumiModalShell
          title={copy.detail.listReturnTitle} eyebrow="Lumi Confirm" lumiState="curious" tone="alert" width={440}
          onClose={() => setIsConfirmingAttachmentListReturn(false)}
          message={copy.detail.listReturnMessage}
          actions={(
            <>
              <button type="button" onClick={() => setIsConfirmingAttachmentListReturn(false)} style={materialToolbarButtonStyle}>{copy.detail.back}</button>
              <button type="button" onClick={handleConfirmAttachmentListReturn} style={materialToolbarButtonStyle}>{copy.detail.list}</button>
            </>
          )}
        />
      ) : null}
      {isConfirmingAttachmentAddCancel ? (
        <LumiModalShell
          title={copy.detail.addCancelTitle} eyebrow="Lumi Confirm" lumiState="curious" tone="alert" width={440}
          onClose={() => setIsConfirmingAttachmentAddCancel(false)}
          message={copy.detail.addCancelMessage}
          actions={(
            <>
              <button type="button" onClick={() => setIsConfirmingAttachmentAddCancel(false)} style={materialToolbarButtonStyle}>{copy.detail.back}</button>
              <button type="button" onClick={handleConfirmAttachmentAddCancel} style={materialToolbarButtonStyle}>{copy.detail.cancel}</button>
            </>
          )}
        />
      ) : null}
      <div style={formGridStyle}>
        <div style={sectionMiniHeaderStyle}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px', flexWrap: 'wrap' }}>
            <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.title}</strong>
            <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.subtitle(researchMaterialAttachments.length, workAttachments.length)}</span>
          </div>
          {canEditLearningRecords && !selectedAttachment && !isAddingAttachment ? (
            <button
              type="button"
              onClick={() => { setIsAddingAttachment(true); setSelectedAttachmentID(null); setWorkAttachmentUploadFile(null); setSelectedAttachmentFile(null); scrollToLearningRecordSection(); }}
              style={materialToolbarPrimaryButtonStyle}
            >
              {copy.add}
            </button>
          ) : null}
        </div>
        {!selectedAttachment && !isAddingAttachment && workAttachments.length ? (
          <div style={formGridStyle}>
            <div style={sectionMiniHeaderStyle}>
              <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{commonCopy.countRange(workAttachments.length, (attachmentPage - 1) * RECORDS_PER_PAGE + 1, Math.min(attachmentPage * RECORDS_PER_PAGE, workAttachments.length))}</span>
            </div>
            <div role="table" aria-label={copy.listAria} style={listBoardStyle}>
              <div role="row" style={getListHeaderRowStyle('58px minmax(0, 1fr) 118px')}>
                <span role="columnheader" style={listColumnHeaderStyle}>{commonCopy.number}</span>
                <span role="columnheader" style={listColumnHeaderStyle}>{commonCopy.title}</span>
                <span role="columnheader" style={listColumnHeaderStyle}>{commonCopy.type}</span>
              </div>
              {pagedAttachments.map((attachment, index) => {
                const rowID = `attachment-${attachment.id}`;
                const hovered = hoveredListRowID === rowID;
                const attachmentIndex = sortedWorkAttachments.length - ((attachmentPage - 1) * RECORDS_PER_PAGE + index);
                return (
                  <button key={attachment.id} type="button" role="row" onClick={() => { setSelectedAttachmentID(attachment.id); setAttachmentMode('view'); setReplacementWorkAttachmentFile(null); setIsAddingAttachment(false); scrollToLearningRecordSection(); }} onMouseEnter={() => setHoveredListRowID(rowID)} onMouseLeave={() => setHoveredListRowID(null)} style={getListDataRowStyle('58px minmax(0, 1fr) 118px', hovered)}>
                    <span role="cell" style={getListIndexCellStyle(hovered)}>{attachmentIndex}</span>
                    <span role="cell" style={getListTextCellStyle(hovered)}>{attachment.title.trim() || copy.fallbackTitle(attachmentIndex)}</span>
                    <span role="cell" style={{ ...blockMetaStyle, color: hovered ? themeTokens.title : themeTokens.metaLabel, fontWeight: 900 }}>{attachment.sourceContext === 'research_material' ? copy.learningContent : copy.reference}</span>
                  </button>
                );
              })}
            </div>
            {attachmentPageCount > 1 ? (
              <div style={{ ...formActionRowStyle, justifyContent: 'center' }}>
                <button type="button" onClick={() => setAttachmentPage((p) => Math.max(1, p - 1))} disabled={attachmentPage <= 1} style={{ ...materialToolbarButtonStyle, opacity: attachmentPage <= 1 ? 0.55 : 1 }}>{commonCopy.previous}</button>
                <span style={{ color: themeTokens.mutedText, fontSize: '14px', fontWeight: 800 }}>{attachmentPage} / {attachmentPageCount}</span>
                <button type="button" onClick={() => setAttachmentPage((p) => Math.min(attachmentPageCount, p + 1))} disabled={attachmentPage >= attachmentPageCount} style={{ ...materialToolbarButtonStyle, opacity: attachmentPage >= attachmentPageCount ? 0.55 : 1 }}>{commonCopy.next}</button>
              </div>
            ) : null}
          </div>
        ) : !selectedAttachment && !isAddingAttachment ? (
          <div style={{ ...emptyCardStyle, background: themeTokens.emptyBackground, borderColor: themeTokens.emptyBorder, color: themeTokens.mutedText }}>{copy.empty}</div>
        ) : null}
        {selectedAttachment ? (() => {
          const isResearchMaterialAttachment = selectedAttachment.sourceContext === 'research_material';
          return (
            <div style={{ ...blockCardStyle, background: themeTokens.surfaceBackground, borderColor: themeTokens.surfaceBorder }}>
              <div style={blockHeaderStyle}>
                <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{isAttachmentEditing ? copy.detail.editTitle : copy.detail.viewTitle}</strong>
                <span style={{ ...blockMetaStyle, padding: '5px 9px', borderRadius: '999px', border: `1px solid ${themeTokens.surfaceBorder}`, background: isResearchMaterialAttachment ? themeTokens.noticeBackground : themeTokens.placeholderBackground, color: themeTokens.metaValue }}>
                  {isResearchMaterialAttachment ? copy.learningContent : copy.reference}
                </span>
              </div>
              <div style={{ display: 'grid', gap: '12px' }}>
                <ResearchSectionHeader title={copy.detail.sectionTitle} themeTokens={themeTokens} />
                {isAttachmentEditing && !isResearchMaterialAttachment ? (
                  <label style={{ display: 'grid', gap: '7px', width: '100%', maxWidth: '860px', margin: '0 auto' }}>
                    <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.editTitleLabel}</span>
                    <input style={{ ...inputStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={selectedAttachment.title} onChange={(e) => handleChangeAttachmentDraft(selectedAttachment.id, 'title', e.target.value)} placeholder={copy.detail.titlePlaceholder} />
                  </label>
                ) : null}
                {isAttachmentEditing && !isResearchMaterialAttachment ? (
                  <div style={{ width: '100%', maxWidth: '860px', margin: '0 auto', display: 'flex', justifyContent: 'flex-end', alignItems: 'center' }}>
                    <input ref={replacementWorkAttachmentInputRef} style={{ display: 'none' }} type="file" accept={RESEARCH_ATTACHMENT_ACCEPT} onChange={(e) => setReplacementWorkAttachmentFile(e.target.files?.[0] ?? null)} />
                    <button type="button" onClick={() => replacementWorkAttachmentInputRef.current?.click()} disabled={isSavingAttachment} style={{ ...materialToolbarButtonStyle, opacity: isSavingAttachment ? 0.62 : 1 }}>{isSavingAttachment ? copy.detail.uploadingFile : copy.detail.uploadFile}</button>
                  </div>
                ) : null}
                <div style={{ width: '100%', maxWidth: '860px', margin: '0 auto', display: 'grid', gap: '8px' }}>
                  <div style={{ ...blockCardStyle, background: themeTokens.surfaceBackground, borderColor: themeTokens.surfaceBorder }}>
                    <div style={blockHeaderStyle}>
                      <button type="button" onClick={() => void handleOpenAttachment(selectedAttachment)} style={{ border: 'none', background: 'transparent', padding: 0, color: themeTokens.title, fontSize: '15px', fontWeight: 900, textAlign: 'left', cursor: 'pointer' }}>
                        {replacementWorkAttachmentFile ? replacementWorkAttachmentFile.name : getResearchAttachmentLabel(selectedAttachment)}
                      </button>
                      {isAttachmentEditing && !isResearchMaterialAttachment && replacementWorkAttachmentFile ? (
                        <span style={{ ...blockMetaStyle, color: themeTokens.metaValue, fontWeight: 850 }}>{copy.detail.replacementPending}</span>
                      ) : null}
                    </div>
                  </div>
                </div>
              </div>
              {isResearchMaterialAttachment ? (
                <div style={{ color: themeTokens.mutedText, fontSize: '14px', lineHeight: 1.6 }}>{copy.detail.researchMaterialNotice}</div>
              ) : null}
              <div style={{ ...formActionRowStyle, justifyContent: 'space-between', alignItems: 'center' }}>
                <div style={{ display: 'flex', justifyContent: 'flex-start', gap: '8px' }}>
                  <button type="button" onClick={() => { if (isAttachmentEditing) setIsConfirmingAttachmentListReturn(true); else { setSelectedAttachmentID(null); setReplacementWorkAttachmentFile(null); scrollToLearningRecordSection(); } }} disabled={isSavingAttachment} style={{ ...materialToolbarButtonStyle, opacity: isSavingAttachment ? 0.62 : 1 }}>{copy.detail.list}</button>
                  {canEditLearningRecords && !isResearchMaterialAttachment ? <button type="button" onClick={() => setPendingDeleteWorkAttachment(selectedAttachment)} disabled={isSavingAttachment} style={{ ...materialToolbarDangerButtonStyle, opacity: isSavingAttachment ? 0.62 : 1 }}>{copy.detail.delete}</button> : null}
                </div>
                <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '8px' }}>
                  {!isAttachmentEditing && canEditLearningRecords && !isResearchMaterialAttachment ? <button type="button" onClick={() => { setAttachmentMode('edit'); setReplacementWorkAttachmentFile(null); scrollToLearningRecordSection(); }} style={materialToolbarButtonStyle}>{copy.detail.edit}</button> : null}
                  {isAttachmentEditing && canEditLearningRecords && !isResearchMaterialAttachment ? <button type="button" onClick={() => void handleSaveAttachmentEdit(selectedAttachment)} disabled={isSavingAttachment || !selectedAttachment.title.trim() || (!selectedAttachment.filePath.trim() && !selectedAttachment.url.trim() && !replacementWorkAttachmentFile)} style={{ ...materialToolbarPrimaryButtonStyle, cursor: isSavingAttachment ? 'progress' : 'pointer', opacity: isSavingAttachment || !selectedAttachment.title.trim() || (!selectedAttachment.filePath.trim() && !selectedAttachment.url.trim() && !replacementWorkAttachmentFile) ? 0.62 : 1 }}>{isSavingAttachment ? copy.detail.saving : copy.detail.save}</button> : null}
                </div>
              </div>
            </div>
          );
        })() : null}
        {canEditLearningRecords && isAddingAttachment ? (
          <div style={{ ...placeholderCardStyle, background: themeTokens.placeholderBackground, borderColor: themeTokens.placeholderBorder }}>
            <div style={formGridStyle}>
              <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.detail.addTitle}</strong>
              <label style={{ display: 'grid', gap: '7px' }}>
                <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.title}</span>
                <input style={{ ...inputStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={newAttachment.title} onChange={(e) => handleChangeNewAttachment('title', e.target.value)} placeholder={copy.detail.titlePlaceholder} />
              </label>
              <label style={{ display: 'grid', gap: '7px' }}>
                <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.file}</span>
                <input
                  ref={workAttachmentInputRef}
                  style={{ display: 'none' }}
                  type="file"
                  accept={RESEARCH_ATTACHMENT_ACCEPT}
                  onChange={(e) => {
                    const file = e.target.files?.[0] ?? null;
                    setWorkAttachmentUploadFile(file);
                    setSelectedAttachmentFile(file);
                    handleChangeNewAttachment('attachmentType', 'file');
                    if (file && !newAttachment.title.trim()) handleChangeNewAttachment('title', file.name);
                  }}
                />
                {selectedWorkAttachmentFile ? (
                  <span style={{ ...blockMetaStyle, color: themeTokens.metaValue, fontWeight: 850 }}>{copy.detail.selectedFile(selectedWorkAttachmentFile.name)}</span>
                ) : null}
              </label>
              <div style={formActionRowStyle}>
                <button type="button" onClick={() => setIsConfirmingAttachmentAddCancel(true)} disabled={isSavingAttachment} style={{ ...materialToolbarButtonStyle, opacity: isSavingAttachment ? 0.62 : 1 }}>{copy.detail.cancel}</button>
                <button
                  type="button"
                  onClick={() => { if (isSavingAttachment) return; if (!selectedWorkAttachmentFile) { workAttachmentInputRef.current?.click(); return; } void handleUploadAttachmentAndClose(); }}
                  disabled={isSavingAttachment}
                  style={{
                    ...materialToolbarPrimaryButtonStyle,
                    background: selectedWorkAttachmentFile ? (themeTokens.pageBackground === '#F8FAFC' ? '#2563EB' : '#3B82F6') : materialToolbarPrimaryButtonStyle.background,
                    borderColor: selectedWorkAttachmentFile ? (themeTokens.pageBackground === '#F8FAFC' ? '#1D4ED8' : '#93C5FD') : materialToolbarPrimaryButtonStyle.borderColor,
                    color: selectedWorkAttachmentFile ? '#FFFFFF' : materialToolbarPrimaryButtonStyle.color,
                    cursor: isSavingAttachment ? 'progress' : 'pointer',
                    opacity: isSavingAttachment ? 0.62 : 1,
                    boxShadow: selectedWorkAttachmentFile ? '0 10px 22px rgba(37, 99, 235, 0.30)' : '0 4px 10px rgba(15, 23, 42, 0.10)',
                  }}
                >
                  {isSavingAttachment ? copy.detail.uploadingFile : selectedWorkAttachmentFile ? copy.detail.uploadSelected : copy.detail.uploadFile}
                </button>
              </div>
            </div>
          </div>
        ) : null}
      </div>
    </>
  );
}
