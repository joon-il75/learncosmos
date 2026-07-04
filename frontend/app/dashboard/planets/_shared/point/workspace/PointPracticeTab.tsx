'use client';

import { useEffect, useState } from 'react';
import LumiModalShell from '@/components/common/LumiModalShell';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../../pointPageUtils';
import type { UsePointLearningResult } from '../usePointLearning';
import {
  formGridStyle, sectionMiniHeaderStyle, blockMetaStyle,
  blockCardStyle, blockActionRowStyle,
  placeholderCardStyle, emptyCardStyle, formActionRowStyle,
  inputStyle, textareaStyle,
  pointShellTitleStyle,
} from '../../pointPageStyles';
import { makeWorkspaceStyles } from './workspaceStyles';
import { scrollToLearningRecordSection } from './workspaceTypes';

interface Props {
  copy: PointLearningCopy['workspace']['learningWork']['recordPanels']['practice'];
  commonCopy: PointLearningCopy['workspace']['learningWork']['recordPanels']['common'];
  learning: UsePointLearningResult;
  themeTokens: PointPageThemeTokens;
  canEditLearningRecords: boolean;
  onEditingChange: (editing: boolean) => void;
}

const RECORDS_PER_PAGE = 5;
const PRACTICE_REFLECTION_SEPARATOR = '\n\n--- 추가 기록 ---\n\n';

export default function PointPracticeTab({ copy, commonCopy, learning, themeTokens, canEditLearningRecords, onEditingChange }: Props) {
  const {
    practiceLogDrafts, newPracticeLog,
    isSavingPracticeLog,
    handleCreatePracticeLog, handleUpdatePracticeLog, handleDeletePracticeLog,
    handleChangePracticeLogDraft, handleChangeNewPracticeLog,
  } = learning;

  const {
    materialToolbarButtonStyle, materialToolbarPrimaryButtonStyle, materialToolbarDangerButtonStyle,
    listBoardStyle, getListHeaderRowStyle, getListDataRowStyle, listColumnHeaderStyle,
    getListIndexCellStyle, getListTextCellStyle, readonlyRecordFieldStyle, readonlyRecordValueStyle,
  } = makeWorkspaceStyles(themeTokens);

  const [selectedPracticeLogID, setSelectedPracticeLogID] = useState<string | null>(null);
  const [practiceMode, setPracticeMode] = useState<'view' | 'edit'>('view');
  const [practicePage, setPracticePage] = useState(1);
  const [isAddingPracticeLog, setIsAddingPracticeLog] = useState(false);
  const [visiblePracticeReflectionLogIDs, setVisiblePracticeReflectionLogIDs] = useState<Set<string>>(() => new Set());
  const [isNewPracticeReflectionVisible, setIsNewPracticeReflectionVisible] = useState(false);
  const [isConfirmingPracticeListReturn, setIsConfirmingPracticeListReturn] = useState(false);
  const [isConfirmingPracticeAddCancel, setIsConfirmingPracticeAddCancel] = useState(false);
  const [pendingDeletePracticeLog, setPendingDeletePracticeLog] = useState<typeof practiceLogDrafts[number] | null>(null);
  const [hoveredListRowID, setHoveredListRowID] = useState<string | null>(null);

  const isEditing = Boolean((selectedPracticeLogID && practiceMode === 'edit') || isAddingPracticeLog);
  useEffect(() => { onEditingChange(isEditing); }, [isEditing, onEditingChange]);

  useEffect(() => {
    const maxPage = Math.max(1, Math.ceil(practiceLogDrafts.length / RECORDS_PER_PAGE));
    setPracticePage((current) => Math.min(current, maxPage));
    if (selectedPracticeLogID && !practiceLogDrafts.some((log) => log.id === selectedPracticeLogID)) {
      setSelectedPracticeLogID(null);
      setPracticeMode('view');
    }
  }, [practiceLogDrafts, selectedPracticeLogID]);

  const getRecentTime = (r: { updatedAt?: string; createdAt?: string }) => {
    const parsed = Date.parse(r.updatedAt || r.createdAt || '');
    return Number.isFinite(parsed) ? parsed : 0;
  };
  const sortedPracticeLogs = [...practiceLogDrafts].sort((a, b) => getRecentTime(b) - getRecentTime(a));
  const practicePageCount = Math.max(1, Math.ceil(sortedPracticeLogs.length / RECORDS_PER_PAGE));
  const pagedPracticeLogs = sortedPracticeLogs.slice((practicePage - 1) * RECORDS_PER_PAGE, practicePage * RECORDS_PER_PAGE);

  const selectedPracticeLog = selectedPracticeLogID ? practiceLogDrafts.find((log) => log.id === selectedPracticeLogID) ?? null : null;
  const isPracticeEditing = Boolean(selectedPracticeLog && practiceMode === 'edit' && canEditLearningRecords);

  const getPracticeReflectionValue = (log: typeof practiceLogDrafts[number]) => {
    if (log.blockedPart.includes(PRACTICE_REFLECTION_SEPARATOR)) {
      return log.blockedPart.split(PRACTICE_REFLECTION_SEPARATOR).filter((v) => v.trim()).join('\n\n');
    }
    return [log.blockedPart, log.nextPlan].filter((v) => v.trim()).join('\n\n');
  };
  const getNewPracticeReflectionValue = () => getPracticeReflectionValue(newPracticeLog);

  const handleChangePracticeReflectionRecord = (logID: string, value: string) => {
    handleChangePracticeLogDraft(logID, 'blockedPart', value);
    handleChangePracticeLogDraft(logID, 'nextPlan', '');
  };

  const handleChangeNewPracticeReflectionRecord = (value: string) => {
    handleChangeNewPracticeLog('blockedPart', value);
    handleChangeNewPracticeLog('nextPlan', '');
  };

  const handleUpdatePracticeLogAndClose = async (log: typeof practiceLogDrafts[number]) => {
    const reflectionValue = getPracticeReflectionValue(log).trim();
    const updated = await handleUpdatePracticeLog({ ...log, blockedPart: reflectionValue, nextPlan: '', nextPractice: '' });
    if (updated) {
      setSelectedPracticeLogID(null);
      setPracticeMode('view');
      setVisiblePracticeReflectionLogIDs((s) => { const next = new Set(s); next.delete(log.id); return next; });
      scrollToLearningRecordSection();
    }
  };

  const handleDeletePracticeLogAndClose = async (log: typeof practiceLogDrafts[number]) => {
    const deleted = await handleDeletePracticeLog(log);
    if (deleted) { setSelectedPracticeLogID(null); setPracticeMode('view'); scrollToLearningRecordSection(); }
  };

  const handleCreatePracticeLogAndClose = async () => {
    const created = await handleCreatePracticeLog();
    if (created) { setIsAddingPracticeLog(false); setIsNewPracticeReflectionVisible(false); scrollToLearningRecordSection(); }
  };

  const handleConfirmPracticeListReturn = () => {
    if (!selectedPracticeLog) return;
    setIsConfirmingPracticeListReturn(false);
    setVisiblePracticeReflectionLogIDs((s) => { const next = new Set(s); next.delete(selectedPracticeLog.id); return next; });
    setSelectedPracticeLogID(null);
    setPracticeMode('view');
    scrollToLearningRecordSection();
  };

  const handleConfirmPracticeAddCancel = () => {
    setIsConfirmingPracticeAddCancel(false);
    setIsNewPracticeReflectionVisible(false);
    setIsAddingPracticeLog(false);
    scrollToLearningRecordSection();
  };

  const handleConfirmDeletePracticeLog = async () => {
    if (!pendingDeletePracticeLog || isSavingPracticeLog) return;
    const log = pendingDeletePracticeLog;
    setPendingDeletePracticeLog(null);
    await handleDeletePracticeLogAndClose(log);
  };

  return (
    <>
      {pendingDeletePracticeLog ? (
        <LumiModalShell
          title={copy.detail.deleteModalTitle} eyebrow="Lumi Confirm" lumiState="curious" tone="alert" width={440}
          onClose={() => { if (!isSavingPracticeLog) setPendingDeletePracticeLog(null); }}
          message={copy.detail.deleteModalMessage(pendingDeletePracticeLog.title || copy.title)}
          actions={(
            <>
              <button type="button" onClick={() => setPendingDeletePracticeLog(null)} disabled={isSavingPracticeLog} style={{ ...materialToolbarButtonStyle, opacity: isSavingPracticeLog ? 0.62 : 1 }}>{copy.detail.cancel}</button>
              <button type="button" onClick={() => void handleConfirmDeletePracticeLog()} disabled={isSavingPracticeLog} style={{ ...materialToolbarDangerButtonStyle, opacity: isSavingPracticeLog ? 0.72 : 1 }}>{isSavingPracticeLog ? copy.detail.deleting : copy.detail.delete}</button>
            </>
          )}
        />
      ) : null}
      {isConfirmingPracticeListReturn ? (
        <LumiModalShell
          title={copy.detail.listReturnTitle} eyebrow="Lumi Confirm" lumiState="curious" tone="alert" width={440}
          onClose={() => setIsConfirmingPracticeListReturn(false)}
          message={copy.detail.listReturnMessage}
          actions={(
            <>
              <button type="button" onClick={() => setIsConfirmingPracticeListReturn(false)} style={materialToolbarButtonStyle}>{copy.detail.back}</button>
              <button type="button" onClick={handleConfirmPracticeListReturn} style={materialToolbarButtonStyle}>{copy.detail.list}</button>
            </>
          )}
        />
      ) : null}
      {isConfirmingPracticeAddCancel ? (
        <LumiModalShell
          title={copy.detail.addCancelTitle} eyebrow="Lumi Confirm" lumiState="curious" tone="alert" width={440}
          onClose={() => setIsConfirmingPracticeAddCancel(false)}
          message={copy.detail.addCancelMessage}
          actions={(
            <>
              <button type="button" onClick={() => setIsConfirmingPracticeAddCancel(false)} style={materialToolbarButtonStyle}>{copy.detail.back}</button>
              <button type="button" onClick={handleConfirmPracticeAddCancel} style={materialToolbarButtonStyle}>{copy.detail.cancel}</button>
            </>
          )}
        />
      ) : null}
      <div style={formGridStyle}>
        <div style={sectionMiniHeaderStyle}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px', flexWrap: 'wrap' }}>
            <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.title}</strong>
            <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.subtitle}</span>
          </div>
          {canEditLearningRecords && !selectedPracticeLog && !isAddingPracticeLog ? (
            <button
              type="button"
              onClick={() => { setIsAddingPracticeLog(true); setSelectedPracticeLogID(null); setIsNewPracticeReflectionVisible(false); scrollToLearningRecordSection(); }}
              style={materialToolbarPrimaryButtonStyle}
            >
              {copy.add}
            </button>
          ) : null}
        </div>
        {!selectedPracticeLog && !isAddingPracticeLog && practiceLogDrafts.length ? (
          <div style={{ display: 'grid', gap: '8px' }}>
            <div style={sectionMiniHeaderStyle}>
              <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{commonCopy.countRange(practiceLogDrafts.length, (practicePage - 1) * RECORDS_PER_PAGE + 1, Math.min(practicePage * RECORDS_PER_PAGE, practiceLogDrafts.length))}</span>
            </div>
            <div role="table" aria-label={copy.listAria} style={listBoardStyle}>
              <div role="row" style={getListHeaderRowStyle('58px minmax(0, 1fr) 120px')}>
                <span role="columnheader" style={listColumnHeaderStyle}>{commonCopy.number}</span>
                <span role="columnheader" style={listColumnHeaderStyle}>{commonCopy.title}</span>
                <span role="columnheader" style={listColumnHeaderStyle}>{commonCopy.status}</span>
              </div>
              {pagedPracticeLogs.map((log, index) => {
                const hasDetail = Boolean(log.durationMinutes.trim() || log.achievementNote.trim() || log.blockedPart.trim() || log.nextPlan.trim());
                const rowID = `practice-${log.id}`;
                const hovered = hoveredListRowID === rowID;
                const practiceIndex = sortedPracticeLogs.length - ((practicePage - 1) * RECORDS_PER_PAGE + index);
                const statusColor = hasDetail ? (themeTokens.pageBackground === '#F8FAFC' ? '#047857' : '#A7F3D0') : (themeTokens.pageBackground === '#F8FAFC' ? '#B45309' : '#FDE68A');
                return (
                  <button
                    key={log.id} type="button" role="row"
                    onClick={() => { setSelectedPracticeLogID(log.id); setPracticeMode('view'); setIsAddingPracticeLog(false); scrollToLearningRecordSection(); setVisiblePracticeReflectionLogIDs((s) => { const next = new Set(s); next.delete(log.id); return next; }); }}
                    onMouseEnter={() => setHoveredListRowID(rowID)}
                    onMouseLeave={() => setHoveredListRowID(null)}
                    style={getListDataRowStyle('58px minmax(0, 1fr) 120px', hovered)}
                  >
                    <span role="cell" style={getListIndexCellStyle(hovered)}>{practiceIndex}</span>
                    <span role="cell" style={getListTextCellStyle(hovered)}>{log.title.trim() || copy.fallbackTitle(practiceIndex)}</span>
                    <span role="cell" style={{ ...blockMetaStyle, color: hovered ? themeTokens.title : statusColor, fontWeight: 900 }}>{hasDetail ? copy.hasRecord : copy.waiting}</span>
                  </button>
                );
              })}
            </div>
            {practicePageCount > 1 ? (
              <div style={{ ...formActionRowStyle, justifyContent: 'center' }}>
                <button type="button" onClick={() => setPracticePage((p) => Math.max(1, p - 1))} disabled={practicePage <= 1} style={{ ...materialToolbarButtonStyle, opacity: practicePage <= 1 ? 0.55 : 1 }}>{commonCopy.previous}</button>
                <span style={{ color: themeTokens.mutedText, fontSize: '14px', fontWeight: 800 }}>{practicePage} / {practicePageCount}</span>
                <button type="button" onClick={() => setPracticePage((p) => Math.min(practicePageCount, p + 1))} disabled={practicePage >= practicePageCount} style={{ ...materialToolbarButtonStyle, opacity: practicePage >= practicePageCount ? 0.55 : 1 }}>{commonCopy.next}</button>
              </div>
            ) : null}
          </div>
        ) : !selectedPracticeLog && !isAddingPracticeLog ? (
          <div style={{ ...emptyCardStyle, background: themeTokens.emptyBackground, borderColor: themeTokens.emptyBorder, color: themeTokens.mutedText }}>{copy.empty}</div>
        ) : null}
        {selectedPracticeLog ? (
          <div style={{ ...blockCardStyle, background: themeTokens.surfaceBackground, borderColor: themeTokens.surfaceBorder }}>
            <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{isPracticeEditing ? copy.detail.editTitle : copy.detail.viewTitle}</strong>
            {isPracticeEditing ? (
              <>
                <label style={{ display: 'grid', gap: '7px' }}>
                  <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.title}</span>
                  <input style={{ ...inputStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={selectedPracticeLog.title} onChange={(e) => handleChangePracticeLogDraft(selectedPracticeLog.id, 'title', e.target.value)} placeholder={copy.detail.titlePlaceholder} readOnly={!canEditLearningRecords} />
                </label>
                <div style={{ display: 'grid', gap: '10px' }}>
                  <label style={{ display: 'grid', gap: '7px', maxWidth: '220px' }}>
                    <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.duration}</span>
                    <input style={{ ...inputStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={selectedPracticeLog.durationMinutes} onChange={(e) => handleChangePracticeLogDraft(selectedPracticeLog.id, 'durationMinutes', e.target.value)} placeholder={copy.detail.durationPlaceholder} inputMode="numeric" pattern="[0-9]*" readOnly={!canEditLearningRecords} />
                  </label>
                  <label style={{ display: 'grid', gap: '7px' }}>
                    <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.achievement}</span>
                    <textarea style={{ ...textareaStyle, minHeight: '86px', background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={selectedPracticeLog.achievementNote} onChange={(e) => handleChangePracticeLogDraft(selectedPracticeLog.id, 'achievementNote', e.target.value)} placeholder={copy.detail.achievementPlaceholder} readOnly={!canEditLearningRecords} />
                  </label>
                </div>
                {(() => {
                  const reflectionValue = getPracticeReflectionValue(selectedPracticeLog);
                  const isReflectionVisible = Boolean(reflectionValue.trim() || visiblePracticeReflectionLogIDs.has(selectedPracticeLog.id));
                  return (
                    <div style={{ display: 'grid', gap: '10px' }}>
                      {isReflectionVisible ? (
                        <label style={{ display: 'grid', gap: '7px' }}>
                          <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.reflection}</span>
                          <textarea style={{ ...textareaStyle, minHeight: '86px', background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={reflectionValue} onChange={(e) => handleChangePracticeReflectionRecord(selectedPracticeLog.id, e.target.value)} placeholder={copy.detail.reflectionPlaceholder} readOnly={!canEditLearningRecords} />
                        </label>
                      ) : null}
                      {canEditLearningRecords && !isReflectionVisible ? (
                        <div style={{ ...formActionRowStyle, justifyContent: 'flex-start' }}>
                          <button type="button" onClick={() => { setVisiblePracticeReflectionLogIDs((s) => new Set(s).add(selectedPracticeLog.id)); scrollToLearningRecordSection(); }} disabled={isSavingPracticeLog} style={{ ...materialToolbarButtonStyle, opacity: isSavingPracticeLog ? 0.62 : 1 }}>{copy.detail.addReflection}</button>
                        </div>
                      ) : null}
                    </div>
                  );
                })()}
              </>
            ) : (() => {
              const reflectionValue = getPracticeReflectionValue(selectedPracticeLog).trim();
              return (
                <div style={{ display: 'grid', gap: '0' }}>
                  <div style={readonlyRecordFieldStyle}>
                    <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.title}</span>
                    <p style={readonlyRecordValueStyle}>{selectedPracticeLog.title.trim() || copy.fallbackTitle(1)}</p>
                  </div>
                  {selectedPracticeLog.durationMinutes.trim() ? (<div style={readonlyRecordFieldStyle}><span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.duration}</span><p style={readonlyRecordValueStyle}>{selectedPracticeLog.durationMinutes}</p></div>) : null}
                  {selectedPracticeLog.achievementNote.trim() ? (<div style={readonlyRecordFieldStyle}><span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.achievement}</span><p style={readonlyRecordValueStyle}>{selectedPracticeLog.achievementNote}</p></div>) : null}
                  {reflectionValue ? (<div style={{ ...readonlyRecordFieldStyle, borderBottom: 'none' }}><span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.reflection}</span><p style={readonlyRecordValueStyle}>{reflectionValue}</p></div>) : null}
                  {!selectedPracticeLog.durationMinutes.trim() && !selectedPracticeLog.achievementNote.trim() && !reflectionValue ? (
                    <div style={{ ...emptyCardStyle, background: themeTokens.emptyBackground, borderColor: themeTokens.emptyBorder, color: themeTokens.mutedText }}>{copy.detail.noDetail}</div>
                  ) : null}
                </div>
              );
            })()}
            {canEditLearningRecords ? (
              <div style={{ ...blockActionRowStyle, alignItems: 'center' }}>
                <div style={{ display: 'flex', justifyContent: 'flex-start', gap: '8px', flexWrap: 'wrap' }}>
                  <button type="button" onClick={() => { if (isPracticeEditing) setIsConfirmingPracticeListReturn(true); else handleConfirmPracticeListReturn(); }} disabled={isSavingPracticeLog} style={{ ...materialToolbarButtonStyle, opacity: isSavingPracticeLog ? 0.62 : 1 }}>{copy.detail.list}</button>
                  <button type="button" onClick={() => setPendingDeletePracticeLog(selectedPracticeLog)} disabled={isSavingPracticeLog} style={{ ...materialToolbarDangerButtonStyle, opacity: isSavingPracticeLog ? 0.62 : 1 }}>{copy.detail.delete}</button>
                </div>
                <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
                  {isPracticeEditing ? (
                    <button type="button" onClick={() => void handleUpdatePracticeLogAndClose(selectedPracticeLog)} disabled={isSavingPracticeLog || !selectedPracticeLog.title.trim()} style={{ ...materialToolbarPrimaryButtonStyle, cursor: isSavingPracticeLog ? 'progress' : 'pointer', opacity: isSavingPracticeLog || !selectedPracticeLog.title.trim() ? 0.62 : 1 }}>{isSavingPracticeLog ? copy.detail.saving : copy.detail.save}</button>
                  ) : (
                    <button type="button" onClick={() => { setPracticeMode('edit'); scrollToLearningRecordSection(); }} style={materialToolbarButtonStyle}>{copy.detail.edit}</button>
                  )}
                </div>
              </div>
            ) : null}
          </div>
        ) : null}
        {canEditLearningRecords && isAddingPracticeLog ? (
          <div style={{ ...placeholderCardStyle, background: themeTokens.placeholderBackground, borderColor: themeTokens.placeholderBorder }}>
            <div style={formGridStyle}>
              <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.detail.addTitle}</strong>
              <label style={{ display: 'grid', gap: '7px' }}>
                <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.title}</span>
                <input style={{ ...inputStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={newPracticeLog.title} onChange={(e) => handleChangeNewPracticeLog('title', e.target.value)} placeholder={copy.detail.titlePlaceholder} />
              </label>
              <div style={{ display: 'grid', gap: '10px' }}>
                <label style={{ display: 'grid', gap: '7px', maxWidth: '220px' }}>
                  <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.duration}</span>
                  <input style={{ ...inputStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={newPracticeLog.durationMinutes} onChange={(e) => handleChangeNewPracticeLog('durationMinutes', e.target.value)} placeholder={copy.detail.durationPlaceholder} inputMode="numeric" pattern="[0-9]*" />
                </label>
                <label style={{ display: 'grid', gap: '7px' }}>
                  <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.achievement}</span>
                  <textarea style={{ ...textareaStyle, minHeight: '86px', background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={newPracticeLog.achievementNote} onChange={(e) => handleChangeNewPracticeLog('achievementNote', e.target.value)} placeholder={copy.detail.achievementPlaceholder} />
                </label>
              </div>
              {(() => {
                const reflectionValue = getNewPracticeReflectionValue();
                const isReflectionVisible = Boolean(reflectionValue.trim() || isNewPracticeReflectionVisible);
                return (
                  <div style={{ display: 'grid', gap: '10px' }}>
                    {isReflectionVisible ? (
                      <label style={{ display: 'grid', gap: '7px' }}>
                        <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.detail.reflection}</span>
                        <textarea style={{ ...textareaStyle, minHeight: '86px', background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={reflectionValue} onChange={(e) => handleChangeNewPracticeReflectionRecord(e.target.value)} placeholder={copy.detail.reflectionPlaceholder} />
                      </label>
                    ) : null}
                    {!isReflectionVisible ? (
                      <div style={{ ...formActionRowStyle, justifyContent: 'flex-start' }}>
                        <button type="button" onClick={() => { setIsNewPracticeReflectionVisible(true); scrollToLearningRecordSection(); }} disabled={isSavingPracticeLog} style={{ ...materialToolbarButtonStyle, opacity: isSavingPracticeLog ? 0.62 : 1 }}>{copy.detail.addReflection}</button>
                      </div>
                    ) : null}
                  </div>
                );
              })()}
              <div style={formActionRowStyle}>
                <button type="button" onClick={() => setIsConfirmingPracticeAddCancel(true)} disabled={isSavingPracticeLog} style={{ ...materialToolbarButtonStyle, opacity: isSavingPracticeLog ? 0.62 : 1 }}>{copy.detail.cancel}</button>
                <button type="button" onClick={() => void handleCreatePracticeLogAndClose()} disabled={isSavingPracticeLog || !newPracticeLog.title.trim()} style={{ ...materialToolbarPrimaryButtonStyle, cursor: isSavingPracticeLog ? 'progress' : 'pointer', opacity: isSavingPracticeLog || !newPracticeLog.title.trim() ? 0.62 : 1 }}>{isSavingPracticeLog ? copy.detail.saving : copy.detail.save}</button>
              </div>
            </div>
          </div>
        ) : null}
      </div>
    </>
  );
}
