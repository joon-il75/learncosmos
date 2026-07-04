'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import LumiModalShell from '@/components/common/LumiModalShell';
import VideoPlayer from '@/components/media/VideoPlayer';
import VideoUploadPanel from '@/components/media/VideoUploadPanel';
import { hasResearchBlockBodyContent, type PointPageThemeTokens } from '../../pointPageUtils';
import {
  blockMetaStyle, blockActionRowStyle, formGridStyle, inputStyle,
  pointShellTitleStyle, emptyCardStyle, primaryButtonStyle,
} from '../../pointPageStyles';
import type { UsePointLearningResult } from '../usePointLearning';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import TiptapResearchEditor from '../TiptapResearchEditor';
import SafeTiptapHTML from '../SafeTiptapHTML';
import { makeWorkspaceStyles } from './workspaceStyles';
import { scrollToWorkspaceTarget } from './workspaceTypes';
import { ResearchSectionHeader, WorkspaceDivider, WorkspaceMiniNav } from './WorkspaceSectionHeader';

interface Props {
  copy: PointLearningCopy['workspace']['researchContentActions'];
  editorCopy: PointLearningCopy['workspace']['editor'];
  materialCopy: PointLearningCopy['workspace']['researchMaterial'];
  learning: UsePointLearningResult;
  themeTokens: PointPageThemeTokens;
  canEditResearchMaterial: boolean;
  isSectionMiniNavSide: boolean;
  editorToolbarStickyTop: string;
  planetID: string;
  pointID: string;
  onResearchMaterialConfirmed: () => void;
  editModeSignal?: number;
  initialResearchEntry?: 'video' | 'content' | 'attachments';
  hideEntryNav?: boolean;
  saveSignal?: number;
  onEntrySaved?: () => void;
  onEntrySaveDisabledChange?: (disabled: boolean) => void;
}

export default function ResearchMaterialSection({
  copy,
  editorCopy,
  materialCopy,
  learning,
  themeTokens,
  canEditResearchMaterial,
  isSectionMiniNavSide,
  editorToolbarStickyTop,
  planetID,
  pointID,
  onResearchMaterialConfirmed,
  editModeSignal = 0,
  initialResearchEntry,
  hideEntryNav = false,
  saveSignal = 0,
  onEntrySaved,
  onEntrySaveDisabledChange,
}: Props) {
  const router = useRouter();
  const [editingResearchBlockIDs, setEditingResearchBlockIDs] = useState<Set<string>>(() => new Set());
  const [pendingDeleteResearchBlock, setPendingDeleteResearchBlock] = useState<typeof learning.researchBlocks[number] | null>(null);
  const [pendingDeleteResearchAttachment, setPendingDeleteResearchAttachment] = useState<typeof learning.attachmentDrafts[number] | null>(null);
  const [researchVideoTitleDraft, setResearchVideoTitleDraft] = useState('');
  const [isResearchEditMode, setIsResearchEditMode] = useState(false);
  const [activeResearchEntry, setActiveResearchEntry] = useState<'video' | 'content' | 'attachments' | null>(() => initialResearchEntry ?? null);
  const [pendingResearchMaterialStateAction, setPendingResearchMaterialStateAction] = useState<'confirm' | 'unconfirm' | null>(null);
  const researchAttachmentInputRef = useRef<HTMLInputElement | null>(null);
  const researchSubtitleInputRef = useRef<HTMLInputElement | null>(null);
  const researchThumbnailInputRef = useRef<HTMLInputElement | null>(null);
  const didPrepareInitialEntryRef = useRef(false);
  const lastHandledSaveSignalRef = useRef(0);

  const {
    routeKind,
    pointDetail,
    attachmentDrafts,
    isSavingAttachment,
    selectedResearchMaterialFile: _selectedResearchMaterialFile,
    setSelectedResearchMaterialFile,
    canAddResearchMaterialAttachment,
    researchMaterialAttachmentLimitMessage,
    researchMaterialAttachmentValidationMessage,
    setResearchMaterialAttachmentValidationMessage,
    researchMaterialUploadProgress,
    handleOpenAttachment,
    handleUploadResearchMaterialAttachment,
    handleUploadResearchMaterialInlineImage,
    handleUpdateAttachment,
    handleDeleteAttachment,
    researchBlocks,
    canAddResearchBlock,
    isSavingResearchBlock,
    isResearchMaterialConfirmed,
    canConfirmResearchMaterial,
    isConfirmingResearchMaterial,
    handleConfirmResearchMaterial,
    handleUnconfirmResearchMaterial,
    handleAddResearchBlock,
    handleChangeResearchBlock,
    handleSaveResearchBlock,
    handleDeleteResearchBlock,
  } = learning;

  const canUseResearchEditMode = canEditResearchMaterial && isResearchEditMode && !isResearchMaterialConfirmed;
  const researchMaterialAttachments = attachmentDrafts.filter((a) => a.sourceContext === 'research_material');
  const researchVideoAttachment = researchMaterialAttachments.find((a) => (
    a.attachmentType === 'video' || a.mimeType.toLowerCase().startsWith('video/')
  )) ?? null;
  const researchSubtitleAttachment = researchMaterialAttachments.find((a) => a.attachmentType === 'subtitle') ?? null;
  const researchThumbnailAttachment = researchMaterialAttachments.find((a) => a.attachmentType === 'thumbnail') ?? null;
  const researchMaterialFileAttachments = researchMaterialAttachments.filter((a) => (
    a.id !== researchVideoAttachment?.id
    && a.id !== researchSubtitleAttachment?.id
    && a.id !== researchThumbnailAttachment?.id
  ));
  const primaryResearchBlock = researchBlocks.find((b) => b.blockType === 'text') ?? null;
  const hasResearchContent = Boolean(primaryResearchBlock && hasResearchBlockBodyContent(primaryResearchBlock));
  const hasSavedResearchContent = Boolean(primaryResearchBlock && !primaryResearchBlock.isNew && hasResearchBlockBodyContent(primaryResearchBlock));
  const hasSavedResearchMaterial = hasSavedResearchContent || researchMaterialAttachments.length > 0;
  const isPrimaryResearchBlockEditing = Boolean(primaryResearchBlock && (
    primaryResearchBlock.isNew || primaryResearchBlock.isDirty || editingResearchBlockIDs.has(primaryResearchBlock.id)
  ));
  const isPointCompleted = pointDetail?.point.status === 'completed';
  const researchMaterialDeleteLocked = isResearchMaterialConfirmed;
  const researchMaterialDeleteLockedMessage = materialCopy.deleteLockedMessage;

  const researchMaterialStatus = isResearchMaterialConfirmed ? 'confirmed' : 'draft';
  const researchMaterialStatusLabel = researchMaterialStatus === 'confirmed' ? materialCopy.statusConfirmed : materialCopy.statusDraft;
  const researchMaterialStatusStyle = {
    borderColor: researchMaterialStatus === 'confirmed'
      ? (themeTokens.pageBackground === '#F8FAFC' ? '#16A34A' : 'rgba(74, 222, 128, 0.48)')
      : themeTokens.surfaceBorder,
    background: researchMaterialStatus === 'confirmed'
      ? (themeTokens.pageBackground === '#F8FAFC' ? '#EAF8EF' : 'rgba(22, 163, 74, 0.18)')
      : themeTokens.placeholderBackground,
    color: researchMaterialStatus === 'confirmed'
      ? (themeTokens.pageBackground === '#F8FAFC' ? '#166534' : '#BBF7D0')
      : themeTokens.description,
  } as const;
  const researchMaterialConfirmDescription = isResearchMaterialConfirmed
    ? materialCopy.confirmedHelp
    : hasSavedResearchMaterial
        ? materialCopy.confirmReadyHelp
        : materialCopy.confirmBlockedHelp;
  const canActivateResearchMaterialConfirm = canUseResearchEditMode && canConfirmResearchMaterial && !isResearchMaterialConfirmed && !isPointCompleted;
  const canActivateResearchMaterialUnconfirm = canEditResearchMaterial && isResearchMaterialConfirmed && !isPointCompleted;
  const canActivateResearchMaterialAction = canActivateResearchMaterialConfirm || canActivateResearchMaterialUnconfirm;
  const confirmResearchMaterialButtonStyle = canActivateResearchMaterialAction
    ? {
        ...primaryButtonStyle,
        border: 'none',
        background: isResearchMaterialConfirmed
          ? (themeTokens.pageBackground === '#F8FAFC' ? '#EA580C' : '#C2410C')
          : (themeTokens.pageBackground === '#F8FAFC' ? '#2563EB' : '#1D4ED8'),
        borderColor: isResearchMaterialConfirmed
          ? (themeTokens.pageBackground === '#F8FAFC' ? '#C2410C' : 'rgba(251, 146, 60, 0.62)')
          : (themeTokens.pageBackground === '#F8FAFC' ? '#1D4ED8' : 'rgba(147, 197, 253, 0.58)'),
        color: '#FFFFFF',
      }
    : {
        ...primaryButtonStyle,
        border: 'none',
        background: themeTokens.disabledButtonBackground,
        borderColor: themeTokens.disabledButtonBorder,
        color: themeTokens.disabledButtonText,
      };
  const researchAttachmentAccept = '.jpg,.jpeg,.png,.webp,.gif,.pdf,.doc,.docx,.hwp,.hwpx,.txt,.xls,.xlsx,.csv,.ppt,.pptx,.rtf,.odt,.md,.zip';
  const researchVideoAccept = '.mp4,.webm,.mov,video/mp4,video/webm,video/quicktime';
  const researchSubtitleAccept = '.vtt,.srt,text/vtt,text/plain';
  const researchThumbnailAccept = '.jpg,.jpeg,.png,.webp,.gif,image/jpeg,image/png,image/webp,image/gif';
  const isResearchAttachmentImage = (attachment: typeof attachmentDrafts[number]) => (
    attachment.attachmentType === 'image' || attachment.mimeType.toLowerCase().startsWith('image/')
  );
  const getResearchAttachmentInlineSrc = (attachmentID: string) => (
    `/api/v1/planets/learning/${planetID}/points/${pointID}/attachments/${attachmentID}/inline`
  );
  const getResearchAttachmentLabel = (attachment: typeof attachmentDrafts[number]) => (
    attachment.title || attachment.filePath || attachment.url || materialCopy.attachmentFallback
  );
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
    if (!res.ok || !payload.url) throw new Error(payload.error ?? materialCopy.videoUrlFailed);
    return payload.url;
  }, [planetID, pointID]);
  const getResearchBlockSaveLabel = (block: typeof researchBlocks[number]) => {
    if (isSavingResearchBlock) return copy.saving;
    if (block.isNew) return copy.saveNew;
    if (block.isDirty) return copy.saveChanges;
    return copy.saved;
  };
  const isResearchBlockTitleMissing = (block: typeof researchBlocks[number]) => (
    block.blockType === 'text' && !block.title.trim()
  );
  const isResearchBlockBodyMissing = (block: typeof researchBlocks[number]) => (
    block.blockType === 'text' && !hasResearchBlockBodyContent(block)
  );
  const isResearchBlockSaveDisabled = (block: typeof researchBlocks[number]) => (
    isResearchMaterialConfirmed || isSavingResearchBlock || isResearchBlockTitleMissing(block) || isResearchBlockBodyMissing(block) || (!block.isNew && !block.isDirty)
  );
  const getResearchBlockDeleteLabel = (block: typeof researchBlocks[number]) => (
    block.title.trim() || materialCopy.contentFallback
  );
  const isEntrySaveDisabled = activeResearchEntry === 'content' && (
    isResearchMaterialConfirmed || !primaryResearchBlock || isResearchBlockTitleMissing(primaryResearchBlock) || isResearchBlockBodyMissing(primaryResearchBlock)
  );
  const {
    materialToolbarButtonStyle,
    materialToolbarPrimaryButtonStyle,
    materialToolbarDangerButtonStyle,
  } = makeWorkspaceStyles(themeTokens);
  const sectionMiniNavBodyLayoutStyle = isSectionMiniNavSide
    ? {
        width: '100%',
        maxWidth: '1064px',
        margin: '0 auto',
        display: 'grid',
        gridTemplateColumns: '124px minmax(0, 920px)',
        gap: '20px',
        alignItems: 'start',
      } as const
    : { display: 'grid', gap: '14px' } as const;

  useEffect(() => {
    onEntrySaveDisabledChange?.(isEntrySaveDisabled);
  }, [isEntrySaveDisabled, onEntrySaveDisabledChange]);

  useEffect(() => {
    setIsResearchEditMode(false);
    setActiveResearchEntry(initialResearchEntry ?? null);
  }, [initialResearchEntry, pointID]);

  useEffect(() => {
    if (editModeSignal > 0 && canEditResearchMaterial && !isResearchMaterialConfirmed && !isPointCompleted) {
      setIsResearchEditMode(true);
      setActiveResearchEntry(initialResearchEntry ?? null);
    }
  }, [canEditResearchMaterial, editModeSignal, initialResearchEntry, isPointCompleted, isResearchMaterialConfirmed]);

  useEffect(() => {
    if (isResearchMaterialConfirmed) {
      setIsResearchEditMode(false);
    }
  }, [isResearchMaterialConfirmed]);

  useEffect(() => {
    setEditingResearchBlockIDs((current) => {
      const liveIDs = new Set(researchBlocks.map((b) => b.id));
      const next = new Set([...current].filter((id) => liveIDs.has(id)));
      researchBlocks.forEach((b) => {
        if (b.isNew) next.add(b.id);
        if (!b.isNew && !b.isDirty) next.delete(b.id);
      });
      return next;
    });
  }, [researchBlocks]);

  useEffect(() => {
    setResearchVideoTitleDraft(researchVideoAttachment?.title ?? '');
  }, [researchVideoAttachment?.id, researchVideoAttachment?.title]);

  useEffect(() => {
    if (didPrepareInitialEntryRef.current || initialResearchEntry !== 'content' || !canUseResearchEditMode) return;
    if (primaryResearchBlock) {
      didPrepareInitialEntryRef.current = true;
      setEditingResearchBlockIDs((current) => new Set(current).add(primaryResearchBlock.id));
      scrollToWorkspaceTarget('point-research-block-' + primaryResearchBlock.id);
      return;
    }
    if (canAddResearchBlock) {
      const nextBlockID = handleAddResearchBlock('text');
      if (nextBlockID) {
        didPrepareInitialEntryRef.current = true;
        setEditingResearchBlockIDs((current) => new Set(current).add(nextBlockID));
        scrollToWorkspaceTarget('point-research-block-' + nextBlockID);
      }
    }
  }, [canAddResearchBlock, canUseResearchEditMode, handleAddResearchBlock, initialResearchEntry, primaryResearchBlock]);

  const handlePickResearchAttachmentFile = () => {
    if (!canAddResearchMaterialAttachment || isSavingAttachment) return;
    researchAttachmentInputRef.current?.click();
  };
  const handleSelectResearchAttachmentFile = (file: File | null) => {
    setSelectedResearchMaterialFile(file);
    if (!file) return;
    void handleUploadResearchMaterialAttachment(file).then((uploaded) => {
      if (uploaded) scrollToWorkspaceTarget('point-research-attachment-section');
    });
  };
  const handleSelectResearchVideoFile = (file: File, title: string) => {
    setSelectedResearchMaterialFile(file);
    void handleUploadResearchMaterialAttachment(file, 'video', title).then((uploaded) => {
      if (uploaded) scrollToWorkspaceTarget('point-research-video-section');
    });
  };
  const handleSaveResearchVideoTitle = async () => {
    if (!researchVideoAttachment || isSavingAttachment) return;
    const nextTitle = researchVideoTitleDraft.trim();
    if (!nextTitle || nextTitle === researchVideoAttachment.title.trim()) return;
    await handleUpdateAttachment({ ...researchVideoAttachment, title: nextTitle, isDirty: true });
  };
  const handleSelectResearchSubtitleFile = (file: File) => {
    setSelectedResearchMaterialFile(file);
    void handleUploadResearchMaterialAttachment(file, 'subtitle').then((uploaded) => {
      if (uploaded) scrollToWorkspaceTarget('point-research-video-section');
    });
  };
  const handleSelectResearchThumbnailFile = (file: File) => {
    setSelectedResearchMaterialFile(file);
    void handleUploadResearchMaterialAttachment(file, 'thumbnail').then((uploaded) => {
      if (uploaded) scrollToWorkspaceTarget('point-research-video-section');
    });
  };
  const handleConfirmDeleteResearchAttachment = async () => {
    if (!pendingDeleteResearchAttachment || isSavingAttachment) return;
    if (pendingDeleteResearchAttachment.sourceContext === 'research_material' && researchMaterialDeleteLocked) {
      setPendingDeleteResearchAttachment(null);
      return;
    }
    const attachment = pendingDeleteResearchAttachment;
    setPendingDeleteResearchAttachment(null);
    await handleDeleteAttachment(attachment);
  };
  const handleConfirmDeleteResearchBlock = async () => {
    if (!pendingDeleteResearchBlock || isSavingResearchBlock) return;
    if (researchMaterialDeleteLocked) {
      setPendingDeleteResearchBlock(null);
      return;
    }
    const block = pendingDeleteResearchBlock;
    setPendingDeleteResearchBlock(null);
    await handleDeleteResearchBlock(block);
  };
  const handleOpenResearchMaterialStateConfirm = () => {
    if (!canActivateResearchMaterialAction || isConfirmingResearchMaterial) return;
    setPendingResearchMaterialStateAction(isResearchMaterialConfirmed ? 'unconfirm' : 'confirm');
  };
  const handleConfirmResearchMaterialStateChange = async () => {
    if (!pendingResearchMaterialStateAction || isConfirmingResearchMaterial) return;
    const action = pendingResearchMaterialStateAction;
    if (action === 'unconfirm') {
      await handleUnconfirmResearchMaterial();
      setPendingResearchMaterialStateAction(null);
      return;
    }
    const confirmed = await handleConfirmResearchMaterial();
    if (confirmed) {
      onResearchMaterialConfirmed();
      setPendingResearchMaterialStateAction(null);
    }
  };
  const handleSelectResearchEntry = (entry: 'video' | 'content' | 'attachments') => {
    if (planetID && pointID && initialResearchEntry !== entry) {
      router.push('/dashboard/planets/learning/' + planetID + '/points/' + pointID + '/research-material/' + entry);
      return;
    }
    setActiveResearchEntry(entry);
    if (entry === 'video') {
      scrollToWorkspaceTarget('point-research-video-section');
      return;
    }
    if (entry === 'attachments') {
      scrollToWorkspaceTarget('point-research-attachment-section');
      return;
    }
    if (primaryResearchBlock) {
      setEditingResearchBlockIDs((current) => new Set(current).add(primaryResearchBlock.id));
      scrollToWorkspaceTarget('point-research-block-' + primaryResearchBlock.id);
      return;
    }
    if (canUseResearchEditMode && canAddResearchBlock) {
      const nextBlockID = handleAddResearchBlock('text');
      if (nextBlockID) {
        setEditingResearchBlockIDs((current) => new Set(current).add(nextBlockID));
        scrollToWorkspaceTarget('point-research-block-' + nextBlockID);
        return;
      }
    }
    scrollToWorkspaceTarget('point-research-content-section');
  };

  useEffect(() => {
    if (!saveSignal || saveSignal === lastHandledSaveSignalRef.current) return;
    lastHandledSaveSignalRef.current = saveSignal;
    const saveAndReturn = async () => {
      if (activeResearchEntry === 'video' && researchVideoAttachment && researchVideoTitleDraft.trim() && researchVideoTitleDraft.trim() !== researchVideoAttachment.title.trim()) {
        await handleSaveResearchVideoTitle();
      }
      if (activeResearchEntry === 'content') {
        if (!primaryResearchBlock || isResearchBlockTitleMissing(primaryResearchBlock) || isResearchBlockBodyMissing(primaryResearchBlock)) return;
        if (!isResearchBlockSaveDisabled(primaryResearchBlock)) {
          await handleSaveResearchBlockAndClose(primaryResearchBlock, 0);
        }
      }
      onEntrySaved?.();
    };
    void saveAndReturn();
  }, [activeResearchEntry, onEntrySaved, primaryResearchBlock, researchVideoAttachment, researchVideoTitleDraft, saveSignal]);

  const handleAddContentAndScroll = () => {
    if (primaryResearchBlock) {
      setEditingResearchBlockIDs((current) => new Set(current).add(primaryResearchBlock.id));
      scrollToWorkspaceTarget(`point-research-block-${primaryResearchBlock.id}`);
      return;
    }
    const nextBlockID = handleAddResearchBlock('text');
    if (nextBlockID) {
      setEditingResearchBlockIDs((current) => new Set(current).add(nextBlockID));
      scrollToWorkspaceTarget(`point-research-block-${nextBlockID}`);
    }
  };
  const handleSaveResearchBlockAndClose = async (block: typeof researchBlocks[number], orderIndex: number) => {
    await handleSaveResearchBlock({ ...block, orderIndex });
  };

  const renderResearchAttachmentList = () => (
    <div id="point-research-attachment-section" style={{ display: 'grid', gap: '12px' }}>
      <ResearchSectionHeader title={materialCopy.attachmentsTitle} themeTokens={themeTokens} />
      {canUseResearchEditMode && !researchMaterialFileAttachments.length ? (
        <div style={{ width: '100%', maxWidth: '860px', margin: '0 auto', display: 'flex', justifyContent: 'flex-end', alignItems: 'center' }}>
          <input
            ref={researchAttachmentInputRef}
            type="file"
            accept={researchAttachmentAccept}
            style={{ display: 'none' }}
            onChange={(event) => {
              const file = event.target.files?.[0] ?? null;
              handleSelectResearchAttachmentFile(file);
              event.currentTarget.value = '';
            }}
          />
          <button
            type="button"
            onClick={handlePickResearchAttachmentFile}
            disabled={!canAddResearchMaterialAttachment || isSavingAttachment}
            title={researchMaterialAttachments.length >= 10 ? materialCopy.attachmentLimitTitle : materialCopy.attachmentAddTitle}
            style={{
              ...materialToolbarButtonStyle,
              opacity: canAddResearchMaterialAttachment && !isSavingAttachment ? 1 : 0.62,
              cursor: canAddResearchMaterialAttachment && !isSavingAttachment ? 'pointer' : 'default',
            }}
          >
            <span>{isSavingAttachment ? materialCopy.uploadingAttachment : materialCopy.addAttachment}</span>
          </button>
        </div>
      ) : null}
      {researchMaterialAttachmentLimitMessage ? (
        <div style={{ color: themeTokens.metaValue, fontSize: '13px', fontWeight: 750, lineHeight: 1.5, textAlign: 'center' }}>
          {researchMaterialAttachmentLimitMessage}
        </div>
      ) : null}
      {researchMaterialFileAttachments.length ? (
        <div style={{ display: 'grid', gap: '10px', width: '100%', maxWidth: '860px', margin: '0 auto', padding: '12px', borderRadius: '8px', border: '1px solid ' + themeTokens.surfaceBorder, background: themeTokens.surfaceBackground }}>
          <strong style={{ color: themeTokens.title, fontSize: '14px', fontWeight: 900, lineHeight: 1.35 }}>{materialCopy.uploadedFileListTitle}</strong>
          {researchMaterialFileAttachments.map((attachment) => (
            <div key={attachment.id} style={{ display: 'grid', gap: '8px', width: '100%' }}>
              <div style={{ display: 'flex', flexWrap: 'wrap', alignItems: 'center', justifyContent: 'space-between', gap: '8px' }}>
                <span style={{ minWidth: 0, flex: '1 1 180px', color: themeTokens.title, fontSize: '14px', fontWeight: 850, lineHeight: 1.35, overflowWrap: 'anywhere' }}>
                  {getResearchAttachmentLabel(attachment)}
                </span>
                {formatAttachmentSize(attachment.fileSize) ? (
                  <span style={{ color: themeTokens.mutedText, fontSize: '13px', fontWeight: 800, lineHeight: 1.4 }}>{formatAttachmentSize(attachment.fileSize)}</span>
                ) : null}
                <button type="button" onClick={() => void handleOpenAttachment(attachment)} style={{ ...materialToolbarButtonStyle, minHeight: '34px', padding: '0 12px' }}>
                  🔗 {materialCopy.openFileLink}
                </button>
                {canUseResearchEditMode ? (
                  <button
                    type="button"
                    onClick={() => setPendingDeleteResearchAttachment(attachment)}
                    disabled={isSavingAttachment || researchMaterialDeleteLocked}
                    title={researchMaterialDeleteLocked ? researchMaterialDeleteLockedMessage : materialCopy.attachmentDeleteTitle}
                    style={{
                      ...materialToolbarDangerButtonStyle,
                      minHeight: '30px',
                      padding: '0 9px',
                      fontSize: '12px',
                      opacity: isSavingAttachment || researchMaterialDeleteLocked ? 0.62 : 1,
                      cursor: isSavingAttachment || researchMaterialDeleteLocked ? 'default' : 'pointer',
                    }}
                  >
                    {materialCopy.delete}
                  </button>
                ) : null}
              </div>
              {isResearchAttachmentImage(attachment) ? (
                <img
                  src={getResearchAttachmentInlineSrc(attachment.id)}
                  alt={attachment.title || materialCopy.attachmentImageAlt}
                  style={{
                    maxWidth: 'min(100%, 720px)',
                    maxHeight: '420px',
                    objectFit: 'contain',
                    borderRadius: '6px',
                    border: `1px solid ${themeTokens.surfaceBorder}`,
                    background: themeTokens.surfaceBackground,
                  }}
                />
              ) : null}
            </div>
          ))}
        </div>
      ) : (
        <div style={{ color: themeTokens.mutedText, fontSize: '15px', lineHeight: 1.65, padding: '10px 0', textAlign: 'center' }}>
          {materialCopy.emptyAttachments}
        </div>
      )}
      <div style={{ color: themeTokens.mutedText, fontSize: '13px', lineHeight: 1.55, textAlign: 'left' }}>
        {materialCopy.attachmentHelp}
      </div>
    </div>
  );

  const renderResearchVideoSlot = () => (
    <div id="point-research-video-section" style={{ display: 'grid', gap: '10px' }}>
      <ResearchSectionHeader title={materialCopy.videoSectionTitle} themeTokens={themeTokens} />
      {researchVideoAttachment ? (
        <div style={{ display: 'grid', gap: '10px' }}>
          <div style={{ display: 'grid', gap: '8px', width: '100%', maxWidth: '860px', margin: '0 auto', padding: '12px', borderRadius: '8px', border: '1px solid ' + themeTokens.surfaceBorder, background: themeTokens.surfaceBackground }}>
            <strong style={{ color: themeTokens.title, fontSize: '14px', fontWeight: 900, lineHeight: 1.35 }}>{materialCopy.uploadedVideoListTitle}</strong>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: '8px', flexWrap: 'wrap' }}>
              <span style={{ minWidth: 0, flex: '1 1 180px', color: themeTokens.title, fontSize: '14px', fontWeight: 850, lineHeight: 1.35, overflowWrap: 'anywhere' }}>
                {getResearchAttachmentLabel(researchVideoAttachment)}
              </span>
              {formatAttachmentSize(researchVideoAttachment.fileSize) ? (
                <span style={{ color: themeTokens.mutedText, fontSize: '13px', fontWeight: 800, lineHeight: 1.4 }}>{formatAttachmentSize(researchVideoAttachment.fileSize)}</span>
              ) : null}
              <button type="button" onClick={() => void handleOpenAttachment(researchVideoAttachment)} style={{ ...materialToolbarButtonStyle, minHeight: '34px', padding: '0 12px' }}>
                🔗 {materialCopy.openVideoLink}
              </button>
            </div>
          </div>
          {canUseResearchEditMode ? (
            <div style={{ display: 'grid', gap: '7px', width: '100%', maxWidth: '860px', margin: '0 auto' }}>
              <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{materialCopy.videoTitleLabel}</span>
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px', flexWrap: 'wrap' }}>
                <input
                  type="text"
                  value={researchVideoTitleDraft}
                  onChange={(event) => setResearchVideoTitleDraft(event.target.value)}
                  placeholder={materialCopy.videoTitlePlaceholder}
                  style={{
                    ...inputStyle,
                    flex: '1 1 280px',
                    minWidth: 0,
                    background: themeTokens.inputBackground,
                    borderColor: themeTokens.inputBorder,
                    color: themeTokens.inputText,
                    fontSize: '15px',
                    fontWeight: 800,
                  }}
                />
                <button
                  type="button"
                  onClick={() => void handleSaveResearchVideoTitle()}
                  disabled={isSavingAttachment || !researchVideoTitleDraft.trim() || researchVideoTitleDraft.trim() === researchVideoAttachment.title.trim()}
                  style={{
                    ...materialToolbarPrimaryButtonStyle,
                    flex: '0 0 auto',
                    cursor: isSavingAttachment ? 'progress' : researchVideoTitleDraft.trim() && researchVideoTitleDraft.trim() !== researchVideoAttachment.title.trim() ? 'pointer' : 'default',
                    opacity: isSavingAttachment || !researchVideoTitleDraft.trim() || researchVideoTitleDraft.trim() === researchVideoAttachment.title.trim() ? 0.72 : 1,
                  }}
                >
                  {isSavingAttachment ? materialCopy.saving : materialCopy.saveTitle}
                </button>
              </div>
            </div>
          ) : null}
          <VideoPlayer
            title={materialCopy.videoPlayerTitle}
            meta={[getResearchAttachmentLabel(researchVideoAttachment), formatAttachmentSize(researchVideoAttachment.fileSize)].filter(Boolean).join(' · ')}
            sourceKey={researchVideoAttachment.id}
            posterKey={researchThumbnailAttachment?.id}
            subtitleKey={researchSubtitleAttachment?.id}
            loadSource={() => loadAttachmentSource(researchVideoAttachment.id)}
            loadPosterSource={researchThumbnailAttachment ? () => loadAttachmentSource(researchThumbnailAttachment.id) : undefined}
            loadSubtitleSource={researchSubtitleAttachment ? () => Promise.resolve(getResearchAttachmentInlineSrc(researchSubtitleAttachment.id)) : undefined}
            borderColor={themeTokens.surfaceBorder}
            background={themeTokens.surfaceBackground}
            textColor={themeTokens.title}
            mutedColor={themeTokens.mutedText}
          />
          {canUseResearchEditMode ? (
            <div style={{ ...blockActionRowStyle, width: '100%', maxWidth: '860px', margin: '0 auto', justifyContent: 'space-between', alignItems: 'center' }}>
              <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px', alignItems: 'center' }}>
                <input ref={researchSubtitleInputRef} type="file" accept={researchSubtitleAccept} style={{ display: 'none' }} onChange={(event) => { const file = event.target.files?.[0] ?? null; event.currentTarget.value = ''; if (file) handleSelectResearchSubtitleFile(file); }} />
                <input ref={researchThumbnailInputRef} type="file" accept={researchThumbnailAccept} style={{ display: 'none' }} onChange={(event) => { const file = event.target.files?.[0] ?? null; event.currentTarget.value = ''; if (file) handleSelectResearchThumbnailFile(file); }} />
                {!researchSubtitleAttachment ? (
                  <button type="button" onClick={() => researchSubtitleInputRef.current?.click()} disabled={isSavingAttachment} style={{ ...materialToolbarButtonStyle, minHeight: '30px', padding: '0 9px', fontSize: '12px', opacity: isSavingAttachment ? 0.62 : 1 }}>{materialCopy.uploadSubtitle}</button>
                ) : null}
                {!researchThumbnailAttachment ? (
                  <button type="button" onClick={() => researchThumbnailInputRef.current?.click()} disabled={isSavingAttachment} style={{ ...materialToolbarButtonStyle, minHeight: '30px', padding: '0 9px', fontSize: '12px', opacity: isSavingAttachment ? 0.62 : 1 }}>{materialCopy.uploadThumbnail}</button>
                ) : null}
                {researchSubtitleAttachment ? (
                  <button type="button" onClick={() => setPendingDeleteResearchAttachment(researchSubtitleAttachment)} disabled={isSavingAttachment || researchMaterialDeleteLocked} title={researchMaterialDeleteLocked ? researchMaterialDeleteLockedMessage : materialCopy.deleteSubtitle} style={{ ...materialToolbarDangerButtonStyle, minHeight: '30px', padding: '0 9px', fontSize: '12px', opacity: isSavingAttachment || researchMaterialDeleteLocked ? 0.62 : 1 }}>{materialCopy.deleteSubtitle}</button>
                ) : null}
                {researchThumbnailAttachment ? (
                  <button type="button" onClick={() => setPendingDeleteResearchAttachment(researchThumbnailAttachment)} disabled={isSavingAttachment || researchMaterialDeleteLocked} title={researchMaterialDeleteLocked ? researchMaterialDeleteLockedMessage : materialCopy.deleteThumbnail} style={{ ...materialToolbarDangerButtonStyle, minHeight: '30px', padding: '0 9px', fontSize: '12px', opacity: isSavingAttachment || researchMaterialDeleteLocked ? 0.62 : 1 }}>{materialCopy.deleteThumbnail}</button>
                ) : null}
              </div>
              <button type="button" onClick={() => setPendingDeleteResearchAttachment(researchVideoAttachment)} disabled={isSavingAttachment || researchMaterialDeleteLocked} title={researchMaterialDeleteLocked ? researchMaterialDeleteLockedMessage : materialCopy.deleteVideo} style={{ ...materialToolbarDangerButtonStyle, minHeight: '30px', padding: '0 9px', fontSize: '12px', opacity: isSavingAttachment || researchMaterialDeleteLocked ? 0.62 : 1, cursor: isSavingAttachment || researchMaterialDeleteLocked ? 'default' : 'pointer' }}>{materialCopy.deleteVideo}</button>
            </div>
          ) : null}
        </div>
      ) : canUseResearchEditMode ? (
        <VideoUploadPanel
          accept={researchVideoAccept}
          disabled={isSavingAttachment || Boolean(researchVideoAttachment)}
          label={isSavingAttachment ? materialCopy.uploadingVideo : materialCopy.uploadVideo}
          helpText={materialCopy.videoHelp}
          buttonStyle={materialToolbarButtonStyle}
          inputStyle={{
            ...inputStyle,
            background: themeTokens.inputBackground,
            borderColor: themeTokens.inputBorder,
            color: themeTokens.inputText,
            fontSize: '15px',
            fontWeight: 800,
          }}
          mutedColor={themeTokens.mutedText}
          onSelect={handleSelectResearchVideoFile}
        />
      ) : (
        <div style={{ color: themeTokens.mutedText, fontSize: '15px', lineHeight: 1.65, padding: '10px 0', textAlign: 'center' }}>
          {materialCopy.noVideo}
        </div>
      )}
    </div>
  );

  const renderResearchObservationDraftNotice = () => (
    <div style={{ display: 'grid', gap: '12px', padding: '18px', borderRadius: '8px', border: `1px solid ${themeTokens.noticeBorder}`, background: themeTokens.noticeBackground }}>
      <div style={{ display: 'grid', gap: '6px' }}>
        <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{materialCopy.observeDraftTitle}</strong>
        <p style={{ margin: 0, color: themeTokens.mutedText, fontSize: '15px', lineHeight: 1.65 }}>{materialCopy.observeDraftDescription}</p>
      </div>
      <button
        type="button"
        onClick={() => setIsResearchEditMode(true)}
        disabled={!canEditResearchMaterial || isPointCompleted}
        style={{
          ...materialToolbarButtonStyle,
          justifySelf: 'start',
          opacity: canEditResearchMaterial && !isPointCompleted ? 1 : 0.62,
          cursor: canEditResearchMaterial && !isPointCompleted ? 'pointer' : 'default',
        }}
      >
        {materialCopy.openEditMode}
      </button>
      <div style={{ color: themeTokens.mutedText, fontSize: '13px', lineHeight: 1.55 }}>{researchMaterialConfirmDescription}</div>
    </div>
  );

  const renderResearchSingleContentEditor = () => {
    const block = primaryResearchBlock;
    const isBlockEditing = Boolean(block && isPrimaryResearchBlockEditing);
    return (
      <div id="point-research-content-editor" style={{ ...formGridStyle, gap: '14px' }}>
        {activeResearchEntry === 'video' ? (
          <>
            {renderResearchVideoSlot()}
            <WorkspaceDivider themeTokens={themeTokens} />
          </>
        ) : null}
        {activeResearchEntry === 'content' ? <div id="point-research-content-section" style={{ display: 'grid', gap: '12px', scrollMarginTop: '430px' }}>
          {block ? (
            <div id={`point-research-block-${block.id}`} style={{ scrollMarginTop: '430px', display: 'grid', gap: '14px' }}>
              {isBlockEditing ? (
                <>
                  <div style={{ display: 'grid', gap: '7px', width: '100%', maxWidth: '860px', margin: '0 auto' }}>
                    <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{materialCopy.contentTitleLabel}</span>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px', flexWrap: 'wrap' }}>
                      <input
                        style={{ ...inputStyle, flex: '1 1 280px', minWidth: 0, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText, fontSize: '15px', fontWeight: 800 }}
                        value={block.title}
                        onChange={(e) => handleChangeResearchBlock(block.id, 'title', e.target.value)}
                        placeholder={materialCopy.contentTitlePlaceholder}
                        readOnly={!canUseResearchEditMode}
                      />
                    </div>
                  </div>
                  <TiptapResearchEditor
                    value={block.text}
                    editable={canUseResearchEditMode}
                    placeholder={materialCopy.contentPlaceholder}
                    themeTokens={themeTokens}
                    copy={editorCopy}
                    toolbarStickyTop={editorToolbarStickyTop}
                    onChange={(nextValue) => handleChangeResearchBlock(block.id, 'text', nextValue)}
                    onUploadImage={canUseResearchEditMode ? handleUploadResearchMaterialInlineImage : undefined}
                  />
                </>
              ) : (
                <div style={{ display: 'grid', gap: '12px' }}>
                  <div style={{ width: '100%', maxWidth: '860px', margin: '0 auto', display: 'flex', alignItems: 'center', justifyContent: 'flex-start', gap: '10px', flexWrap: 'wrap' }}>
                    {block.title.trim() ? <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{materialCopy.contentDisplayTitle(block.title)}</strong> : <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{materialCopy.contentFallback}</strong>}
                  </div>
                  <SafeTiptapHTML html={block.text} className="lw-tiptap-editor" style={{ width: '100%', maxWidth: '860px', margin: '0 auto', minHeight: 'auto', padding: '0', border: 'none', background: 'transparent', color: themeTokens.description }} />
                </div>
              )}
              {canUseResearchEditMode ? (
                <div style={{ ...blockActionRowStyle, width: '100%', maxWidth: '860px', margin: '0 auto', alignItems: 'center', justifyContent: 'space-between' }}>
                  <button
                    type="button"
                    onClick={() => setPendingDeleteResearchBlock(block)}
                    disabled={isSavingResearchBlock || researchMaterialDeleteLocked}
                    title={researchMaterialDeleteLocked ? researchMaterialDeleteLockedMessage : materialCopy.deleteContentTitle}
                    style={{
                      ...materialToolbarDangerButtonStyle,
                      opacity: isSavingResearchBlock || researchMaterialDeleteLocked ? 0.72 : 1,
                      cursor: isSavingResearchBlock || researchMaterialDeleteLocked ? 'default' : 'pointer',
                    }}
                  >
                    {materialCopy.delete}
                  </button>
                  <div style={{ marginLeft: 'auto', display: 'flex', justifyContent: 'flex-end', alignItems: 'center', gap: '8px' }}>
                    {!isBlockEditing ? (
                      <button
                        type="button"
                        onClick={() => setEditingResearchBlockIDs((current) => new Set(current).add(block.id))}
                        style={materialToolbarButtonStyle}
                      >
                        {copy.edit}
                      </button>
                    ) : null}
                  </div>
                </div>
              ) : null}
            </div>
          ) : (
            <div style={{ display: 'grid', gap: '12px', textAlign: 'center' }}>
              <div style={{ color: themeTokens.mutedText, fontSize: '15px', lineHeight: 1.65, padding: '18px 0', textAlign: 'center' }}>{materialCopy.emptyContent}</div>
            </div>
          )}
        </div> : null}
        {activeResearchEntry === 'attachments' ? (
          <>
            {renderResearchAttachmentList()}
            <WorkspaceDivider themeTokens={themeTokens} />
          </>
        ) : null}
      </div>
    );
  };

  return (
    <>
      {pendingDeleteResearchBlock ? (
        <LumiModalShell
          title={materialCopy.deleteContentTitle}
          eyebrow="Lumi Confirm"
          lumiState="curious"
          tone="alert"
          width={440}
          onClose={() => { if (!isSavingResearchBlock) setPendingDeleteResearchBlock(null); }}
          message={materialCopy.deleteContentMessage(getResearchBlockDeleteLabel(pendingDeleteResearchBlock))}
          actions={(
            <>
              <button type="button" onClick={() => setPendingDeleteResearchBlock(null)} disabled={isSavingResearchBlock} style={{ ...materialToolbarButtonStyle, opacity: isSavingResearchBlock ? 0.62 : 1 }}>{materialCopy.cancel}</button>
              <button type="button" onClick={() => void handleConfirmDeleteResearchBlock()} disabled={isSavingResearchBlock} style={{ ...materialToolbarDangerButtonStyle, opacity: isSavingResearchBlock ? 0.72 : 1, cursor: isSavingResearchBlock ? 'default' : 'pointer' }}>
                {isSavingResearchBlock ? materialCopy.deleting : materialCopy.delete}
              </button>
            </>
          )}
        />
      ) : null}
      {pendingDeleteResearchAttachment ? (
        <LumiModalShell
          title={materialCopy.deleteAttachmentModalTitle}
          eyebrow="Lumi Confirm"
          lumiState="curious"
          tone="alert"
          width={440}
          onClose={() => { if (!isSavingAttachment) setPendingDeleteResearchAttachment(null); }}
          message={materialCopy.deleteAttachmentMessage(getResearchAttachmentLabel(pendingDeleteResearchAttachment))}
          actions={(
            <>
              <button type="button" onClick={() => setPendingDeleteResearchAttachment(null)} disabled={isSavingAttachment} style={{ ...materialToolbarButtonStyle, opacity: isSavingAttachment ? 0.62 : 1 }}>{materialCopy.cancel}</button>
              <button type="button" onClick={() => void handleConfirmDeleteResearchAttachment()} disabled={isSavingAttachment} style={{ ...materialToolbarDangerButtonStyle, opacity: isSavingAttachment ? 0.72 : 1, cursor: isSavingAttachment ? 'default' : 'pointer' }}>
                {isSavingAttachment ? materialCopy.deleting : materialCopy.delete}
              </button>
            </>
          )}
        />
      ) : null}
      {pendingResearchMaterialStateAction ? (
        <LumiModalShell
          title={pendingResearchMaterialStateAction === 'confirm' ? materialCopy.confirmModalTitle : materialCopy.unconfirmModalTitle}
          eyebrow="Lumi Confirm"
          lumiState="curious"
          tone="alert"
          width={440}
          onClose={() => { if (!isConfirmingResearchMaterial) setPendingResearchMaterialStateAction(null); }}
          message={pendingResearchMaterialStateAction === 'confirm' ? materialCopy.confirmModalMessage : materialCopy.unconfirmModalMessage}
          actions={(
            <>
              <button type="button" onClick={() => setPendingResearchMaterialStateAction(null)} disabled={isConfirmingResearchMaterial} style={{ ...materialToolbarButtonStyle, opacity: isConfirmingResearchMaterial ? 0.62 : 1 }}>{materialCopy.cancel}</button>
              <button type="button" onClick={() => void handleConfirmResearchMaterialStateChange()} disabled={isConfirmingResearchMaterial} style={{ ...confirmResearchMaterialButtonStyle, opacity: isConfirmingResearchMaterial ? 0.72 : 1, cursor: isConfirmingResearchMaterial ? 'default' : 'pointer' }}>
                {isConfirmingResearchMaterial
                  ? (pendingResearchMaterialStateAction === 'confirm' ? materialCopy.confirmSaving : materialCopy.unconfirmSaving)
                  : (pendingResearchMaterialStateAction === 'confirm' ? materialCopy.confirmAction : materialCopy.unconfirmAction)}
              </button>
            </>
          )}
        />
      ) : null}
      {researchMaterialAttachmentValidationMessage ? (
        <LumiModalShell
          title={materialCopy.validationTitle}
          eyebrow="Lumi Notice"
          lumiState="curious"
          tone="alert"
          width={440}
          onClose={() => setResearchMaterialAttachmentValidationMessage(null)}
          message={researchMaterialAttachmentValidationMessage}
          actions={<button type="button" onClick={() => setResearchMaterialAttachmentValidationMessage(null)} style={materialToolbarPrimaryButtonStyle}>{materialCopy.validationConfirm}</button>}
        />
      ) : null}
      {researchMaterialUploadProgress !== null ? (
        <LumiModalShell
          title={materialCopy.uploadProgressTitle}
          eyebrow="Lumi Upload"
          lumiState="encourage"
          width={440}
          onClose={() => undefined}
          message={(
            <div style={{ display: 'grid', gap: '10px' }}>
              <span>{materialCopy.uploadProgressMessage}</span>
              <div aria-label={materialCopy.uploadProgressAria(researchMaterialUploadProgress)} role="progressbar" aria-valuemin={0} aria-valuemax={100} aria-valuenow={researchMaterialUploadProgress} style={{ height: '12px', borderRadius: '999px', overflow: 'hidden', background: 'rgba(120, 70, 18, 0.16)', border: '1px solid rgba(120, 70, 18, 0.18)' }}>
                <div style={{ width: `${researchMaterialUploadProgress}%`, height: '100%', borderRadius: '999px', background: 'linear-gradient(90deg, #7BA7E8, #14B8A6)', transition: 'width 120ms ease' }} />
              </div>
              <strong style={{ fontSize: '16px', lineHeight: 1.2 }}>{researchMaterialUploadProgress}%</strong>
            </div>
          )}
          actions={null}
        />
      ) : null}
      <div style={{ display: 'grid', gap: '10px' }}>
        <div
          data-testid="research-material-status"
          style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: '10px', borderRadius: '8px', border: `1px solid ${researchMaterialStatusStyle.borderColor}`, background: researchMaterialStatusStyle.background, color: researchMaterialStatusStyle.color, fontSize: '13px', fontWeight: 800, lineHeight: 1.35, padding: '9px 12px' }}
        >
          <span>{materialCopy.statusLabel}</span>
          <span style={{ display: 'inline-flex', alignItems: 'center', gap: '8px', flexWrap: 'wrap', justifyContent: 'flex-end' }}>
            <strong style={{ fontSize: '14px', whiteSpace: 'nowrap' }}>{researchMaterialStatusLabel}</strong>
          </span>
        </div>
        {!isResearchMaterialConfirmed && !isResearchEditMode ? (
          renderResearchObservationDraftNotice()
        ) : (
          <div style={{ display: 'grid', gap: '14px' }}>
            {!hideEntryNav ? (
              <WorkspaceMiniNav
              ariaLabel={materialCopy.navAria}
              items={[
                { key: 'video', label: materialCopy.navVideo, active: activeResearchEntry === 'video', checked: Boolean(researchVideoAttachment), onClick: () => handleSelectResearchEntry('video') },
                { key: 'content', label: materialCopy.navContent, active: activeResearchEntry === 'content', checked: hasResearchContent, onClick: () => handleSelectResearchEntry('content') },
                { key: 'attachments', label: materialCopy.navAttachments, active: activeResearchEntry === 'attachments', checked: researchMaterialFileAttachments.length > 0, onClick: () => handleSelectResearchEntry('attachments') },
              ]}
              layout="top"
              variant="register"
              themeTokens={themeTokens}
              />
            ) : null}
            {!hideEntryNav && routeKind === 'learning' ? (
              <div style={{ display: 'grid', gap: '10px', justifyItems: 'center', width: '100%', maxWidth: '860px', margin: '0 auto' }}>
                <button
                  type="button"
                  onClick={handleOpenResearchMaterialStateConfirm}
                  disabled={isConfirmingResearchMaterial || isPointCompleted || (!isResearchMaterialConfirmed && !canConfirmResearchMaterial)}
                  style={{ ...confirmResearchMaterialButtonStyle, minHeight: '42px', minWidth: 'min(260px, 100%)', cursor: isConfirmingResearchMaterial ? 'progress' : canActivateResearchMaterialAction ? 'pointer' : 'default', opacity: canActivateResearchMaterialAction ? 1 : 0.72 }}
                >
                  {isConfirmingResearchMaterial ? (isResearchMaterialConfirmed ? materialCopy.unconfirmSaving : materialCopy.confirmSaving) : isResearchMaterialConfirmed ? materialCopy.unconfirmAction : materialCopy.confirmAction}
                </button>
                <div style={{ justifySelf: 'stretch', color: themeTokens.mutedText, fontSize: '14px', lineHeight: 1.55, textAlign: 'center' }}>
                  {researchMaterialConfirmDescription}
                </div>
              </div>
            ) : null}
            <div style={{ display: 'grid', gap: '10px', minWidth: 0 }}>
              {activeResearchEntry ? renderResearchSingleContentEditor() : null}
            </div>
          </div>
        )}
      </div>
    </>
  );
}
