'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import LumiModalShell from '@/components/common/LumiModalShell';
import VideoPlayer from '@/components/media/VideoPlayer';
import VideoUploadPanel from '@/components/media/VideoUploadPanel';
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
import TiptapResearchEditor from '../TiptapResearchEditor';
import SafeTiptapHTML from '../SafeTiptapHTML';
import { makeWorkspaceStyles } from './workspaceStyles';
import { ResearchSectionHeader } from './WorkspaceSectionHeader';
import { normalizeArtifactKind, artifactTypeOptions, scrollToLearningRecordSection } from './workspaceTypes';
import type { ArtifactKind } from './workspaceTypes';

const RECORDS_PER_PAGE = 5;
const RESEARCH_VIDEO_ACCEPT = '.mp4,.webm,.mov,video/mp4,video/webm,video/quicktime';
const RESEARCH_SUBTITLE_ACCEPT = '.vtt,.srt,text/vtt,text/plain';
const RESEARCH_THUMBNAIL_ACCEPT = '.jpg,.jpeg,.png,.webp,.gif,image/jpeg,image/png,image/webp,image/gif';
const RESEARCH_ATTACHMENT_ACCEPT = '.jpg,.jpeg,.png,.webp,.gif,.pdf,.doc,.docx,.hwp,.hwpx,.txt,.xls,.xlsx,.csv,.ppt,.pptx,.rtf,.odt,.md,.zip';

interface Props {
  copy: PointLearningCopy['workspace']['learningWork']['recordPanels']['artifacts'];
  commonCopy: PointLearningCopy['workspace']['learningWork']['recordPanels']['common'];
  editorCopy: PointLearningCopy['workspace']['editor'];
  learning: UsePointLearningResult;
  themeTokens: PointPageThemeTokens;
  canEditLearningRecords: boolean;
  isSharedRoute: boolean;
  editorToolbarStickyTop: string;
  onEditingChange: (editing: boolean) => void;
}

export default function PointArtifactsTab({ copy, commonCopy, editorCopy, learning, themeTokens, canEditLearningRecords, isSharedRoute, editorToolbarStickyTop, onEditingChange }: Props) {
  const {
    planetID, pointID,
    artifactType, setArtifactType,
    artifactTitle, setArtifactTitle,
    artifactDescription, setArtifactDescription,
    artifactDrafts, attachmentDrafts,
    isSavingArtifact, isSavingAttachment,
    handleSaveArtifact, handleChangeArtifactDraft, handleUpdateArtifact, handleDeleteArtifact,
    handleUploadArtifactAttachment, handleOpenAttachment, handleDeleteAttachment,
  } = learning;

  const {
    materialToolbarButtonStyle, materialToolbarPrimaryButtonStyle, materialToolbarDangerButtonStyle,
    listBoardStyle, getListHeaderRowStyle, getListDataRowStyle, listColumnHeaderStyle,
    getListIndexCellStyle, getListTextCellStyle,
  } = makeWorkspaceStyles(themeTokens);

  const artifactAttachmentInputRef = useRef<HTMLInputElement | null>(null);
  const artifactSubtitleInputRef = useRef<HTMLInputElement | null>(null);
  const artifactThumbnailInputRef = useRef<HTMLInputElement | null>(null);

  const [selectedArtifactID, setSelectedArtifactID] = useState<string | null>(null);
  const [artifactMode, setArtifactMode] = useState<'view' | 'edit'>('view');
  const [artifactPage, setArtifactPage] = useState(1);
  const [isAddingArtifact, setIsAddingArtifact] = useState(false);
  const [isConfirmingArtifactListReturn, setIsConfirmingArtifactListReturn] = useState(false);
  const [isConfirmingArtifactAddCancel, setIsConfirmingArtifactAddCancel] = useState(false);
  const [pendingDeleteArtifact, setPendingDeleteArtifact] = useState<typeof artifactDrafts[number] | null>(null);
  const [pendingDeleteArtifactAttachment, setPendingDeleteArtifactAttachment] = useState<typeof attachmentDrafts[number] | null>(null);
  const [hoveredListRowID, setHoveredListRowID] = useState<string | null>(null);

  const isEditing = Boolean(selectedArtifactID || isAddingArtifact);
  useEffect(() => { onEditingChange(isEditing); }, [isEditing, onEditingChange]);

  useEffect(() => {
    const maxPage = Math.max(1, Math.ceil(artifactDrafts.length / RECORDS_PER_PAGE));
    setArtifactPage((p) => Math.min(p, maxPage));
    if (selectedArtifactID && !artifactDrafts.some((a) => a.id === selectedArtifactID)) setSelectedArtifactID(null);
  }, [artifactDrafts, selectedArtifactID]);

  const getRecentTime = (r: { updatedAt?: string; createdAt?: string }) => {
    const parsed = Date.parse(r.updatedAt || r.createdAt || '');
    return Number.isFinite(parsed) ? parsed : 0;
  };
  const sortedArtifacts = [...artifactDrafts].sort((a, b) => {
    const byTime = getRecentTime(b) - getRecentTime(a);
    return byTime || b.orderIndex - a.orderIndex;
  });
  const artifactPageCount = Math.max(1, Math.ceil(sortedArtifacts.length / RECORDS_PER_PAGE));
  const pagedArtifacts = sortedArtifacts.slice((artifactPage - 1) * RECORDS_PER_PAGE, artifactPage * RECORDS_PER_PAGE);

  const selectedArtifact = selectedArtifactID ? artifactDrafts.find((a) => a.id === selectedArtifactID) ?? null : null;
  const isArtifactEditing = Boolean(selectedArtifact && artifactMode === 'edit' && canEditLearningRecords);

  const artifactAttachments = selectedArtifact
    ? attachmentDrafts.filter((a) => a.sourceContext === 'artifact' && a.artifactID === selectedArtifact.id)
    : [];
  const artifactVideoAttachment = artifactAttachments.find((a) => a.attachmentType === 'video' || a.mimeType.toLowerCase().startsWith('video/')) ?? null;
  const artifactSubtitleAttachment = artifactAttachments.find((a) => a.attachmentType === 'subtitle') ?? null;
  const artifactThumbnailAttachment = artifactAttachments.find((a) => a.attachmentType === 'thumbnail') ?? null;
  const artifactFileAttachments = artifactAttachments.filter((a) => (
    a.id !== artifactVideoAttachment?.id && a.id !== artifactSubtitleAttachment?.id && a.id !== artifactThumbnailAttachment?.id
  ));

  const getResearchAttachmentLabel = (attachment: typeof attachmentDrafts[number]) =>
    attachment.title || attachment.filePath || attachment.url || copy.detail.attachmentFallback;

  const formatAttachmentSize = (value: string) => {
    const bytes = Number(value);
    if (!Number.isFinite(bytes) || bytes <= 0) return '';
    if (bytes >= 1024 * 1024) return `${Math.round((bytes / 1024 / 1024) * 10) / 10}MB`;
    if (bytes >= 1024) return `${Math.round((bytes / 1024) * 10) / 10}KB`;
    return `${bytes}B`;
  };

  const loadAttachmentSource = useCallback(async (attachmentID: string) => {
    if (!planetID || !pointID) return null;
    const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/attachments/${attachmentID}/open`, { credentials: 'include' });
    const payload = (await res.json().catch(() => ({}))) as { url?: string; error?: string };
    if (!res.ok || !payload.url) throw new Error(payload.error ?? 'Failed to create video URL.');
    return payload.url;
  }, [planetID, pointID]);

  const handleSelectArtifactVideoFile = (file: File, title: string) => {
    if (!selectedArtifact) return;
    void handleUploadArtifactAttachment(selectedArtifact.id, file, 'video', title).then((uploaded) => {
      if (uploaded) scrollToLearningRecordSection();
    });
  };
  const handleSelectArtifactSubtitleFile = (file: File) => {
    if (!selectedArtifact) return;
    void handleUploadArtifactAttachment(selectedArtifact.id, file, 'subtitle').then((uploaded) => {
      if (uploaded) scrollToLearningRecordSection();
    });
  };
  const handleSelectArtifactThumbnailFile = (file: File) => {
    if (!selectedArtifact) return;
    void handleUploadArtifactAttachment(selectedArtifact.id, file, 'thumbnail').then((uploaded) => {
      if (uploaded) scrollToLearningRecordSection();
    });
  };
  const handleSelectArtifactAttachmentFile = (file: File | null) => {
    if (!selectedArtifact || !file) return;
    void handleUploadArtifactAttachment(selectedArtifact.id, file).then((uploaded) => {
      if (uploaded) scrollToLearningRecordSection();
    });
  };

  const handleConfirmDeleteArtifactAttachment = async () => {
    if (!pendingDeleteArtifactAttachment || isSavingAttachment) return;
    const attachment = pendingDeleteArtifactAttachment;
    setPendingDeleteArtifactAttachment(null);
    await handleDeleteAttachment(attachment);
  };

  const handleConfirmArtifactListReturn = () => {
    setIsConfirmingArtifactListReturn(false);
    setSelectedArtifactID(null);
    setArtifactMode('view');
    scrollToLearningRecordSection();
  };

  const handleConfirmArtifactAddCancel = () => {
    setIsConfirmingArtifactAddCancel(false);
    setIsAddingArtifact(false);
    scrollToLearningRecordSection();
  };

  const handleConfirmDeleteArtifact = async () => {
    if (!pendingDeleteArtifact || isSavingArtifact) return;
    const artifact = pendingDeleteArtifact;
    setPendingDeleteArtifact(null);
    const deleted = await handleDeleteArtifact(artifact);
    if (deleted) { setSelectedArtifactID(null); scrollToLearningRecordSection(); }
  };

  const handleUpdateArtifactAndClose = async (artifact: typeof artifactDrafts[number]) => {
    const updated = await handleUpdateArtifact(artifact);
    if (updated) { setArtifactMode('view'); scrollToLearningRecordSection(); }
  };

  const handleCreateArtifactAndClose = async () => {
    const created = await handleSaveArtifact();
    if (created) { setIsAddingArtifact(false); setSelectedArtifactID(created.id); setArtifactMode('edit'); scrollToLearningRecordSection(); }
  };

  const handleCreateVideoArtifactAndUpload = async (file: File, title: string) => {
    const nextTitle = title.trim() || file.name;
    const created = await handleSaveArtifact({ artifactType: '영상자료', title: nextTitle });
    if (!created) return;
    setSelectedArtifactID(created.id);
    setArtifactMode('edit');
    await handleUploadArtifactAttachment(created.id, file, 'video', nextTitle);
    setIsAddingArtifact(false);
    scrollToLearningRecordSection();
  };

  const handleCreateAttachmentArtifactAndUpload = async (file: File | null) => {
    if (!file) return;
    const created = await handleSaveArtifact({ artifactType: '첨부자료', title: file.name });
    if (!created) return;
    setSelectedArtifactID(created.id);
    setArtifactMode('edit');
    await handleUploadArtifactAttachment(created.id, file);
    setIsAddingArtifact(false);
    scrollToLearningRecordSection();
  };

  const renderArtifactVideoSlot = () => (
    <div style={{ display: 'grid', gap: '10px' }}>
      <ResearchSectionHeader title={copy.detail.videoSection} themeTokens={themeTokens} />
      {artifactVideoAttachment ? (
        <div style={{ display: 'grid', gap: '10px' }}>
          <VideoPlayer
            title={copy.detail.videoPlayerTitle}
            meta={[getResearchAttachmentLabel(artifactVideoAttachment), formatAttachmentSize(artifactVideoAttachment.fileSize)].filter(Boolean).join(' · ')}
            sourceKey={artifactVideoAttachment.id}
            posterKey={artifactThumbnailAttachment?.id}
            subtitleKey={artifactSubtitleAttachment?.id}
            loadSource={() => loadAttachmentSource(artifactVideoAttachment.id)}
            loadPosterSource={artifactThumbnailAttachment ? () => loadAttachmentSource(artifactThumbnailAttachment.id) : undefined}
            loadSubtitleSource={artifactSubtitleAttachment ? () => Promise.resolve(`/api/v1/planets/learning/${planetID}/points/${pointID}/attachments/${artifactSubtitleAttachment.id}/inline`) : undefined}
            borderColor={themeTokens.surfaceBorder}
            background={themeTokens.surfaceBackground}
            textColor={themeTokens.title}
            mutedColor={themeTokens.mutedText}
          />
          {canEditLearningRecords && artifactMode === 'edit' ? (
            <div style={{ ...blockActionRowStyle, width: '100%', maxWidth: '860px', margin: '0 auto', justifyContent: 'space-between', alignItems: 'center' }}>
              <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px', alignItems: 'center' }}>
                <input ref={artifactSubtitleInputRef} type="file" accept={RESEARCH_SUBTITLE_ACCEPT} style={{ display: 'none' }} onChange={(e) => { const file = e.target.files?.[0] ?? null; e.currentTarget.value = ''; if (file) handleSelectArtifactSubtitleFile(file); }} />
                <input ref={artifactThumbnailInputRef} type="file" accept={RESEARCH_THUMBNAIL_ACCEPT} style={{ display: 'none' }} onChange={(e) => { const file = e.target.files?.[0] ?? null; e.currentTarget.value = ''; if (file) handleSelectArtifactThumbnailFile(file); }} />
                {!artifactSubtitleAttachment ? <button type="button" onClick={() => artifactSubtitleInputRef.current?.click()} disabled={isSavingAttachment} style={{ ...materialToolbarButtonStyle, minHeight: '30px', padding: '0 9px', fontSize: '12px', opacity: isSavingAttachment ? 0.62 : 1 }}>{copy.detail.uploadSubtitle}</button> : null}
                {!artifactThumbnailAttachment ? <button type="button" onClick={() => artifactThumbnailInputRef.current?.click()} disabled={isSavingAttachment} style={{ ...materialToolbarButtonStyle, minHeight: '30px', padding: '0 9px', fontSize: '12px', opacity: isSavingAttachment ? 0.62 : 1 }}>{copy.detail.uploadThumbnail}</button> : null}
                {artifactSubtitleAttachment ? <button type="button" onClick={() => setPendingDeleteArtifactAttachment(artifactSubtitleAttachment)} disabled={isSavingAttachment} style={{ ...materialToolbarDangerButtonStyle, minHeight: '30px', padding: '0 9px', fontSize: '12px', opacity: isSavingAttachment ? 0.62 : 1 }}>{copy.detail.deleteSubtitle}</button> : null}
                {artifactThumbnailAttachment ? <button type="button" onClick={() => setPendingDeleteArtifactAttachment(artifactThumbnailAttachment)} disabled={isSavingAttachment} style={{ ...materialToolbarDangerButtonStyle, minHeight: '30px', padding: '0 9px', fontSize: '12px', opacity: isSavingAttachment ? 0.62 : 1 }}>{copy.detail.deleteThumbnail}</button> : null}
              </div>
              <button type="button" onClick={() => setPendingDeleteArtifactAttachment(artifactVideoAttachment)} disabled={isSavingAttachment} style={{ ...materialToolbarDangerButtonStyle, minHeight: '30px', padding: '0 9px', fontSize: '12px', opacity: isSavingAttachment ? 0.62 : 1 }}>{copy.detail.deleteVideo}</button>
            </div>
          ) : null}
        </div>
      ) : canEditLearningRecords && artifactMode === 'edit' ? (
        <VideoUploadPanel
          accept={RESEARCH_VIDEO_ACCEPT}
          disabled={isSavingAttachment}
          label={isSavingAttachment ? copy.detail.uploadingVideo : copy.detail.uploadVideo}
          helpText={copy.detail.videoHelp}
          buttonStyle={materialToolbarButtonStyle}
          inputStyle={{ ...inputStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText, fontSize: '15px', fontWeight: 800 }}
          mutedColor={themeTokens.mutedText}
          onSelect={handleSelectArtifactVideoFile}
        />
      ) : (
        <div style={{ color: themeTokens.mutedText, fontSize: '15px', lineHeight: 1.65, padding: '10px 0', textAlign: 'center' }}>{copy.detail.noVideo}</div>
      )}
    </div>
  );

  const renderArtifactAttachmentList = () => (
    <div style={{ display: 'grid', gap: '12px' }}>
      <ResearchSectionHeader title={copy.detail.attachmentSection} themeTokens={themeTokens} />
      {canEditLearningRecords && artifactMode === 'edit' ? (
        <div style={{ width: '100%', maxWidth: '860px', margin: '0 auto', display: 'flex', justifyContent: 'flex-end', alignItems: 'center' }}>
          <input ref={artifactAttachmentInputRef} type="file" accept={RESEARCH_ATTACHMENT_ACCEPT} style={{ display: 'none' }} onChange={(e) => { const file = e.target.files?.[0] ?? null; e.currentTarget.value = ''; handleSelectArtifactAttachmentFile(file); }} />
          <button type="button" onClick={() => artifactAttachmentInputRef.current?.click()} disabled={isSavingAttachment || artifactFileAttachments.length >= 10} style={{ ...materialToolbarButtonStyle, opacity: isSavingAttachment || artifactFileAttachments.length >= 10 ? 0.62 : 1 }}>{isSavingAttachment ? copy.detail.uploadingFile : copy.detail.uploadFile}</button>
        </div>
      ) : null}
      {artifactFileAttachments.length ? (
        <div style={{ width: '100%', maxWidth: '860px', margin: '0 auto', display: 'grid', gap: '8px' }}>
          {artifactFileAttachments.map((attachment) => (
            <div key={attachment.id} style={{ ...blockCardStyle, background: themeTokens.surfaceBackground, borderColor: themeTokens.surfaceBorder }}>
              <div style={blockHeaderStyle}>
                <button type="button" onClick={() => void handleOpenAttachment(attachment)} style={{ border: 'none', background: 'transparent', padding: 0, color: themeTokens.title, fontSize: '15px', fontWeight: 900, textAlign: 'left', cursor: 'pointer' }}>{getResearchAttachmentLabel(attachment)}</button>
                {canEditLearningRecords && artifactMode === 'edit' ? <button type="button" onClick={() => setPendingDeleteArtifactAttachment(attachment)} disabled={isSavingAttachment} style={{ ...materialToolbarDangerButtonStyle, minHeight: '30px', padding: '0 9px', fontSize: '12px', opacity: isSavingAttachment ? 0.62 : 1 }}>{copy.detail.delete}</button> : null}
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div style={{ color: themeTokens.mutedText, fontSize: '15px', lineHeight: 1.65, padding: '10px 0', textAlign: 'center' }}>{copy.detail.noFiles}</div>
      )}
    </div>
  );

  const renderArtifactTitleField = (artifact: typeof artifactDrafts[number]) => (
    <label style={{ display: 'grid', gap: '7px', width: '100%', maxWidth: '860px', margin: '0 auto' }}>
      <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.title}</span>
      <input style={{ ...inputStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText, fontSize: '15px', fontWeight: 800 }} value={artifact.title} onChange={(e) => handleChangeArtifactDraft(artifact.id, 'title', e.target.value)} placeholder={copy.detail.titlePlaceholder} readOnly={!isArtifactEditing} />
    </label>
  );

  const renderArtifactEditableTitleField = (artifact: typeof artifactDrafts[number]) =>
    isArtifactEditing ? (
      <label style={{ display: 'grid', gap: '7px', width: '100%', maxWidth: '860px', margin: '0 auto' }}>
        <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.editTitleLabel}</span>
        <input style={{ ...inputStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText, fontSize: '15px', fontWeight: 800 }} value={artifact.title} onChange={(e) => handleChangeArtifactDraft(artifact.id, 'title', e.target.value)} placeholder={copy.detail.titlePlaceholder} />
      </label>
    ) : null;

  const renderSelectedArtifactBody = (artifact: typeof artifactDrafts[number]) => {
    const kind = normalizeArtifactKind(artifact.artifactType);
    if (kind === '영상자료') {
      return (<div style={{ display: 'grid', gap: '14px' }}>{renderArtifactEditableTitleField(artifact)}{renderArtifactVideoSlot()}</div>);
    }
    if (kind === '첨부자료') {
      return (<div style={{ display: 'grid', gap: '14px' }}>{renderArtifactEditableTitleField(artifact)}{renderArtifactAttachmentList()}</div>);
    }
    return (
      <div style={{ display: 'grid', gap: '14px' }}>
        <ResearchSectionHeader title={copy.detail.documentSection} themeTokens={themeTokens} />
        {renderArtifactTitleField(artifact)}
        {isArtifactEditing ? (
          <TiptapResearchEditor value={artifact.description} editable placeholder={copy.detail.documentPlaceholder} copy={editorCopy} themeTokens={themeTokens} toolbarStickyTop={editorToolbarStickyTop} onChange={(nextValue) => handleChangeArtifactDraft(artifact.id, 'description', nextValue)} />
        ) : (
          <SafeTiptapHTML html={artifact.description || '<p></p>'} className="lw-tiptap-editor" style={{ width: '100%', maxWidth: '860px', margin: '0 auto', minHeight: 'auto', padding: '0', border: 'none', background: 'transparent', color: themeTokens.description }} />
        )}
      </div>
    );
  };

  const renderNewArtifactBody = () => {
    const kind = normalizeArtifactKind(artifactType);
    if (kind === '영상자료') {
      return (
        <div style={{ display: 'grid', gap: '14px' }}>
          <ResearchSectionHeader title={copy.detail.videoSection} themeTokens={themeTokens} />
          <VideoUploadPanel accept={RESEARCH_VIDEO_ACCEPT} disabled={isSavingArtifact || isSavingAttachment} label={isSavingArtifact || isSavingAttachment ? copy.detail.uploadingVideo : copy.detail.uploadVideo} helpText={copy.detail.videoHelp} buttonStyle={materialToolbarButtonStyle} inputStyle={{ ...inputStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText, fontSize: '15px', fontWeight: 800 }} mutedColor={themeTokens.mutedText} onSelect={handleCreateVideoArtifactAndUpload} />
        </div>
      );
    }
    if (kind === '첨부자료') {
      return (
        <div style={{ display: 'grid', gap: '14px' }}>
          <ResearchSectionHeader title={copy.detail.attachmentSection} themeTokens={themeTokens} />
          <div style={{ width: '100%', maxWidth: '860px', margin: '0 auto', display: 'flex', justifyContent: 'flex-end', alignItems: 'center' }}>
            <input ref={artifactAttachmentInputRef} type="file" accept={RESEARCH_ATTACHMENT_ACCEPT} style={{ display: 'none' }} onChange={(e) => { const file = e.target.files?.[0] ?? null; e.currentTarget.value = ''; void handleCreateAttachmentArtifactAndUpload(file); }} />
            <button type="button" onClick={() => artifactAttachmentInputRef.current?.click()} disabled={isSavingArtifact || isSavingAttachment} style={{ ...materialToolbarButtonStyle, opacity: isSavingArtifact || isSavingAttachment ? 0.62 : 1 }}>{isSavingArtifact || isSavingAttachment ? copy.detail.uploadingFile : copy.detail.uploadFile}</button>
          </div>
        </div>
      );
    }
    return (
      <div style={{ display: 'grid', gap: '14px' }}>
        <ResearchSectionHeader title={copy.detail.documentSection} themeTokens={themeTokens} />
        <label style={{ display: 'grid', gap: '7px', width: '100%', maxWidth: '860px', margin: '0 auto' }}>
          <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.title}</span>
          <input style={{ ...inputStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText, fontSize: '15px', fontWeight: 800 }} value={artifactTitle} onChange={(e) => setArtifactTitle(e.target.value)} placeholder={copy.detail.documentTitlePlaceholder} />
        </label>
        <TiptapResearchEditor value={artifactDescription} editable placeholder={copy.detail.documentPlaceholder} copy={editorCopy} themeTokens={themeTokens} toolbarStickyTop={editorToolbarStickyTop} onChange={setArtifactDescription} />
      </div>
    );
  };

  return (
    <>
      {pendingDeleteArtifact ? (
        <LumiModalShell
          title={copy.detail.deleteModalTitle} eyebrow="Lumi Confirm" lumiState="curious" tone="alert" width={440}
          onClose={() => { if (!isSavingArtifact) setPendingDeleteArtifact(null); }}
          message={copy.detail.deleteModalMessage(pendingDeleteArtifact.title || copy.title)}
          actions={(
            <>
              <button type="button" onClick={() => setPendingDeleteArtifact(null)} disabled={isSavingArtifact} style={{ ...materialToolbarButtonStyle, opacity: isSavingArtifact ? 0.62 : 1 }}>{copy.detail.cancel}</button>
              <button type="button" onClick={() => void handleConfirmDeleteArtifact()} disabled={isSavingArtifact} style={{ ...materialToolbarDangerButtonStyle, opacity: isSavingArtifact ? 0.72 : 1 }}>{isSavingArtifact ? copy.detail.deleting : copy.detail.delete}</button>
            </>
          )}
        />
      ) : null}
      {pendingDeleteArtifactAttachment ? (
        <LumiModalShell
          title={copy.detail.deleteAttachmentTitle} eyebrow="Lumi Confirm" lumiState="curious" tone="alert" width={440}
          onClose={() => { if (!isSavingAttachment) setPendingDeleteArtifactAttachment(null); }}
          message={copy.detail.deleteAttachmentMessage(getResearchAttachmentLabel(pendingDeleteArtifactAttachment))}
          actions={(
            <>
              <button type="button" onClick={() => setPendingDeleteArtifactAttachment(null)} disabled={isSavingAttachment} style={{ ...materialToolbarButtonStyle, opacity: isSavingAttachment ? 0.62 : 1 }}>{copy.detail.cancel}</button>
              <button type="button" onClick={() => void handleConfirmDeleteArtifactAttachment()} disabled={isSavingAttachment} style={{ ...materialToolbarDangerButtonStyle, opacity: isSavingAttachment ? 0.72 : 1 }}>{isSavingAttachment ? copy.detail.deleting : copy.detail.delete}</button>
            </>
          )}
        />
      ) : null}
      {isConfirmingArtifactListReturn ? (
        <LumiModalShell
          title={copy.detail.listReturnTitle} eyebrow="Lumi Confirm" lumiState="curious" tone="alert" width={440}
          onClose={() => setIsConfirmingArtifactListReturn(false)}
          message={copy.detail.listReturnMessage}
          actions={(
            <>
              <button type="button" onClick={() => setIsConfirmingArtifactListReturn(false)} style={materialToolbarButtonStyle}>{copy.detail.back}</button>
              <button type="button" onClick={handleConfirmArtifactListReturn} style={materialToolbarButtonStyle}>{copy.detail.list}</button>
            </>
          )}
        />
      ) : null}
      {isConfirmingArtifactAddCancel ? (
        <LumiModalShell
          title={copy.detail.addCancelTitle} eyebrow="Lumi Confirm" lumiState="curious" tone="alert" width={440}
          onClose={() => setIsConfirmingArtifactAddCancel(false)}
          message={copy.detail.addCancelMessage}
          actions={(
            <>
              <button type="button" onClick={() => setIsConfirmingArtifactAddCancel(false)} style={materialToolbarButtonStyle}>{copy.detail.back}</button>
              <button type="button" onClick={handleConfirmArtifactAddCancel} style={materialToolbarButtonStyle}>{copy.detail.cancel}</button>
            </>
          )}
        />
      ) : null}
      <div id="point-artifact-submit-section" style={{ ...formGridStyle, scrollMarginTop: '430px' }}>
        <div style={sectionMiniHeaderStyle}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px', flexWrap: 'wrap' }}>
            <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.title}</strong>
            <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.subtitle}</span>
          </div>
          {canEditLearningRecords && !selectedArtifact && !isAddingArtifact ? (
            <button
              type="button"
              onClick={() => { setArtifactType('문서 편집자료'); setIsAddingArtifact(true); setSelectedArtifactID(null); setArtifactMode('edit'); scrollToLearningRecordSection(); }}
              style={materialToolbarPrimaryButtonStyle}
            >
              {copy.add}
            </button>
          ) : null}
        </div>
        {!selectedArtifact && !isAddingArtifact && artifactDrafts.length ? (
          <div style={formGridStyle}>
            <div style={sectionMiniHeaderStyle}>
              <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{commonCopy.countRange(artifactDrafts.length, (artifactPage - 1) * RECORDS_PER_PAGE + 1, Math.min(artifactPage * RECORDS_PER_PAGE, artifactDrafts.length))}</span>
            </div>
            <div role="table" aria-label={copy.listAria} style={listBoardStyle}>
              <div role="row" style={getListHeaderRowStyle('58px minmax(0, 1fr) 132px')}>
                <span role="columnheader" style={listColumnHeaderStyle}>{commonCopy.number}</span>
                <span role="columnheader" style={listColumnHeaderStyle}>{commonCopy.title}</span>
                <span role="columnheader" style={listColumnHeaderStyle}>{commonCopy.type}</span>
              </div>
              {pagedArtifacts.map((artifact, index) => {
                const rowID = `artifact-${artifact.id}`;
                const hovered = hoveredListRowID === rowID;
                const artifactIndex = sortedArtifacts.length - ((artifactPage - 1) * RECORDS_PER_PAGE + index);
                return (
                  <button key={artifact.id} type="button" role="row" onClick={() => { setSelectedArtifactID(artifact.id); setArtifactMode('view'); setIsAddingArtifact(false); scrollToLearningRecordSection(); }} onMouseEnter={() => setHoveredListRowID(rowID)} onMouseLeave={() => setHoveredListRowID(null)} style={getListDataRowStyle('58px minmax(0, 1fr) 132px', hovered)}>
                    <span role="cell" style={getListIndexCellStyle(hovered)}>{artifactIndex}</span>
                    <span role="cell" style={getListTextCellStyle(hovered)}>{artifact.title.trim() || copy.fallbackTitle(artifactIndex)}</span>
                    <span role="cell" style={{ ...blockMetaStyle, color: hovered ? themeTokens.title : themeTokens.metaLabel, fontWeight: 900 }}>{copy.kindLabels[normalizeArtifactKind(artifact.artifactType)]}</span>
                  </button>
                );
              })}
            </div>
            {artifactPageCount > 1 ? (
              <div style={{ ...formActionRowStyle, justifyContent: 'center' }}>
                <button type="button" onClick={() => setArtifactPage((p) => Math.max(1, p - 1))} disabled={artifactPage <= 1} style={{ ...materialToolbarButtonStyle, opacity: artifactPage <= 1 ? 0.55 : 1 }}>{commonCopy.previous}</button>
                <span style={{ color: themeTokens.mutedText, fontSize: '14px', fontWeight: 800 }}>{artifactPage} / {artifactPageCount}</span>
                <button type="button" onClick={() => setArtifactPage((p) => Math.min(artifactPageCount, p + 1))} disabled={artifactPage >= artifactPageCount} style={{ ...materialToolbarButtonStyle, opacity: artifactPage >= artifactPageCount ? 0.55 : 1 }}>{commonCopy.next}</button>
              </div>
            ) : null}
          </div>
        ) : !selectedArtifact && !isAddingArtifact ? (
          <div style={{ ...emptyCardStyle, background: themeTokens.emptyBackground, borderColor: themeTokens.emptyBorder, color: themeTokens.mutedText }}>{copy.empty}</div>
        ) : null}
        {selectedArtifact ? (
          <div id={`point-artifact-${selectedArtifact.id}`} style={{ ...blockCardStyle, background: themeTokens.surfaceBackground, borderColor: themeTokens.surfaceBorder, scrollMarginTop: '430px' }}>
            <div style={blockHeaderStyle}>
              <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{artifactMode === 'edit' ? copy.detail.editTitle : copy.detail.viewTitle}</strong>
            </div>
            {isArtifactEditing ? (
              <label style={{ display: 'grid', gap: '7px', width: '100%', maxWidth: '860px', margin: '0 auto' }}>
                <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.kind}</span>
                <select style={{ ...inputStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={normalizeArtifactKind(selectedArtifact.artifactType)} onChange={(e) => { handleChangeArtifactDraft(selectedArtifact.id, 'artifactType', e.target.value); scrollToLearningRecordSection(); }}>
                  {artifactTypeOptions.map((option) => <option key={option} value={option}>{copy.kindLabels[option]}</option>)}
                </select>
              </label>
            ) : null}
            {renderSelectedArtifactBody(selectedArtifact)}
            <div style={{ ...blockActionRowStyle, alignItems: 'center' }}>
              <div style={{ display: 'flex', justifyContent: 'flex-start', gap: '8px', flexWrap: 'wrap' }}>
                <button type="button" onClick={() => { if (isArtifactEditing) setIsConfirmingArtifactListReturn(true); else handleConfirmArtifactListReturn(); }} disabled={isSavingArtifact} style={{ ...materialToolbarButtonStyle, opacity: isSavingArtifact ? 0.62 : 1 }}>{copy.detail.list}</button>
                {canEditLearningRecords ? <button type="button" onClick={() => setPendingDeleteArtifact(selectedArtifact)} disabled={isSavingArtifact} style={{ ...materialToolbarDangerButtonStyle, opacity: isSavingArtifact ? 0.62 : 1 }}>{copy.detail.delete}</button> : null}
              </div>
              {canEditLearningRecords && artifactMode === 'view' ? <button type="button" onClick={() => { setArtifactMode('edit'); scrollToLearningRecordSection(); }} style={materialToolbarButtonStyle}>{copy.detail.edit}</button> : null}
              {isArtifactEditing ? <button type="button" onClick={() => void handleUpdateArtifactAndClose(selectedArtifact)} disabled={isSavingArtifact} style={{ ...materialToolbarPrimaryButtonStyle, cursor: isSavingArtifact ? 'progress' : 'pointer', opacity: isSavingArtifact ? 0.62 : 1 }}>{isSavingArtifact ? copy.detail.saving : copy.detail.save}</button> : null}
            </div>
          </div>
        ) : null}
        {canEditLearningRecords && isAddingArtifact ? (
          <div style={{ ...placeholderCardStyle, background: themeTokens.placeholderBackground, borderColor: themeTokens.placeholderBorder }}>
            <div style={formGridStyle}>
              <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.detail.addTitle}</strong>
              <label style={{ display: 'grid', gap: '7px', width: '100%', maxWidth: '860px', margin: '0 auto' }}>
                <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.kind}</span>
                <select style={{ ...inputStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={normalizeArtifactKind(artifactType)} onChange={(e) => { setArtifactType(e.target.value as ArtifactKind); scrollToLearningRecordSection(); }}>
                  {artifactTypeOptions.map((option) => <option key={option} value={option}>{copy.kindLabels[option]}</option>)}
                </select>
              </label>
              {renderNewArtifactBody()}
              <div style={formActionRowStyle}>
                <button type="button" onClick={() => setIsConfirmingArtifactAddCancel(true)} disabled={isSavingArtifact} style={{ ...materialToolbarButtonStyle, opacity: isSavingArtifact ? 0.62 : 1 }}>{copy.detail.cancel}</button>
                {normalizeArtifactKind(artifactType) === '문서 편집자료' ? (
                  <button type="button" onClick={() => void handleCreateArtifactAndClose()} disabled={isSavingArtifact || !artifactTitle.trim()} style={{ ...materialToolbarPrimaryButtonStyle, cursor: isSavingArtifact ? 'progress' : 'pointer', opacity: isSavingArtifact || !artifactTitle.trim() ? 0.62 : 1 }}>{isSavingArtifact ? copy.detail.saving : copy.detail.save}</button>
                ) : null}
              </div>
            </div>
          </div>
        ) : null}
      </div>
    </>
  );
}
