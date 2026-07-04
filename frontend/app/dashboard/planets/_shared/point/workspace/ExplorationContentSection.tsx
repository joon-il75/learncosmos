'use client';

import { useState } from 'react';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../../pointPageUtils';
import {
  sectionMiniHeaderStyle, pointShellTitleStyle,
  emptyCardStyle, noticeCardStyle, aiSummarySourceStyle, aiSummaryTextStyle,
  formGridStyle, inputRowStyle, inputStyle, textareaStyle, formActionRowStyle,
  blockMetaStyle, disclosureStyle, disclosureSummaryStyle,
  readingPanelStyle, readingBodyStyle,
} from '../../pointPageStyles';
import type { UsePointLearningResult } from '../usePointLearning';
import { makeWorkspaceStyles } from './workspaceStyles';

interface Props {
  copy: PointLearningCopy['workspace']['explorationContent'];
  learning: UsePointLearningResult;
  themeTokens: PointPageThemeTokens;
  isExplorationPoint: boolean;
  canEditLearningRecords: boolean;
  onSourceContentOpen: () => void;
}

export default function ExplorationContentSection({
  copy,
  learning,
  themeTokens,
  isExplorationPoint,
  canEditLearningRecords,
  onSourceContentOpen,
}: Props) {
  const [replacementLinkTitle, setReplacementLinkTitle] = useState('');
  const [replacementLinkURL, setReplacementLinkURL] = useState('');
  const [contentIssueNotice, setContentIssueNotice] = useState('');

  const {
    pointDetail,
    pointRuntimeMessage,
    pointEntryMessage,
    isSavingAttachment,
    isReportingMaterial,
    materialReportType, setMaterialReportType,
    materialReportMessage, setMaterialReportMessage,
    handleReportMaterial,
    replacementCandidates, replacementMessage, isLoadingReplacementCandidates,
    handleLoadReplacementCandidates, handleUseReplacementCandidate, handleUseReplacementLink, handleCancelMaterialReport,
    isGeneratingAISummary, handleGenerateAISummary,
  } = learning;

  const formatMaterialReportTime = (value: string) => {
    const parsed = Date.parse(value);
    if (!Number.isFinite(parsed)) return value;
    return new Intl.DateTimeFormat(copy.dateLocale, { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }).format(parsed);
  };
  const materialReports = [...(pointDetail?.material_reports ?? [])].sort((left, right) => {
    const r = Date.parse(right.created_at);
    const l = Date.parse(left.created_at);
    return (Number.isFinite(r) ? r : 0) - (Number.isFinite(l) ? l : 0);
  });
  const eligibleSourceReplacementReports = materialReports.filter((report) => (
    report.target_type === 'source'
    && (report.status === 'open' || report.status === 'reviewing')
    && !report.replaced_at
  ));
  const activeSourceReplacementReport = eligibleSourceReplacementReports[0] ?? null;
  const activeSourceReplacementReportID = activeSourceReplacementReport?.id ?? '';
  const hasSourceReplacementGate = Boolean(activeSourceReplacementReportID);
  const latestReplacementReport = materialReports.find((r) => r.replaced_at && r.replacement_url);

  const getYoutubeVideoId = (value: string | null | undefined) => {
    if (!value) return null
    try {
      const parsed = new URL(value)
      const host = parsed.hostname.toLowerCase()
      if (host === 'youtu.be' || host.endsWith('.youtu.be')) {
        const raw = parsed.pathname.replace(/^\//, '')
        return raw.split('/')[0] || null
      }
      if (host === 'youtube.com' || host.endsWith('.youtube.com') || host.endsWith('.youtube-nocookie.com')) {
        const videoParam = parsed.searchParams.get('v')
        if (videoParam) return videoParam
        const pathMatch = parsed.pathname.match(/\/(?:embed|shorts)\/([\w-]+)/)
        if (pathMatch?.[1]) return pathMatch[1]
      }
    } catch {
      return null
    }
    return null
  }
  const getYoutubeWatchLink = (value: string | null | undefined) => {
    const videoId = getYoutubeVideoId(value)
    return videoId ? `https://www.youtube.com/watch?v=${videoId}` : value || '#'
  }
  const getYoutubeEmbedLink = (value: string | null | undefined) => {
    const videoId = getYoutubeVideoId(value)
    return videoId ? `https://www.youtube.com/embed/${videoId}` : value || '#'
  }

  const {
    materialToolbarButtonStyle,
    materialToolbarPrimaryButtonStyle,
    listBoardStyle,
    getListHeaderRowStyle,
    listColumnHeaderStyle,
  } = makeWorkspaceStyles(themeTokens);

  const sourceContentButtonStyle = {
    ...materialToolbarButtonStyle,
    minHeight: '42px',
    justifyContent: 'center',
    background: 'linear-gradient(180deg, #93C5FD, #60A5FA)',
    borderColor: 'rgba(30, 64, 175, 0.35)',
    color: '#0B1220',
    boxShadow: '0 8px 20px rgba(30, 64, 175, 0.16), 0 0 0 2px rgba(191, 219, 254, 0.22)',
  } as const;
  const sourceSecondaryButtonStyle = {
    ...materialToolbarButtonStyle,
    minHeight: '42px',
    justifyContent: 'center',
    background: 'linear-gradient(180deg, #E2E8F0, #CBD5E1)',
    borderColor: 'rgba(71, 85, 105, 0.35)',
    color: '#0F172A',
    boxShadow: '0 8px 20px rgba(15, 23, 42, 0.14), 0 0 0 2px rgba(148, 163, 184, 0.18)',
  } as const;
  const contentIssueDisclosureStyle = {
    ...disclosureStyle,
    padding: '10px 12px',
    background: themeTokens.pageBackground === '#F8FAFC' ? '#FFF7ED' : 'rgba(146, 64, 14, 0.18)',
    borderColor: themeTokens.pageBackground === '#F8FAFC' ? 'rgba(217, 119, 6, 0.42)' : 'rgba(251, 191, 36, 0.34)',
  } as const;
  const contentIssueSummaryStyle = {
    ...disclosureSummaryStyle,
    color: themeTokens.pageBackground === '#F8FAFC' ? '#92400E' : '#FCD34D',
    fontSize: '13px',
    fontWeight: 850,
    lineHeight: 1.35,
  } as const;
  const contentIssueMetaStyle = {
    ...blockMetaStyle,
    color: themeTokens.pageBackground === '#F8FAFC' ? '#9A3412' : '#FDBA74',
  } as const;

  const handleCreateDirectReplacementLink = async () => {
    if (!activeSourceReplacementReportID) return;
    setContentIssueNotice('');
    const saved = await handleUseReplacementLink(replacementLinkTitle, replacementLinkURL, activeSourceReplacementReportID);
    if (saved) {
      setReplacementLinkTitle('');
      setReplacementLinkURL('');
    }
  };
  const handleCancelSourceReplacementReport = async () => {
    if (!activeSourceReplacementReportID) return;
    const cancelled = await handleCancelMaterialReport(activeSourceReplacementReportID);
    if (cancelled) {
      setReplacementLinkTitle('');
      setReplacementLinkURL('');
      setContentIssueNotice(copy.cancelReportSuccessNotice);
    }
  };

  const aiPointDescriptionPanel = isExplorationPoint ? (
    <div style={{ ...disclosureStyle, display: 'grid', gap: '12px', background: themeTokens.placeholderBackground, borderColor: themeTokens.placeholderBorder }}>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: '12px', flexWrap: 'wrap' }}>
        <strong style={{ ...disclosureSummaryStyle, cursor: 'default', color: themeTokens.title }}>{copy.aiSummaryTitle}</strong>
        {canEditLearningRecords ? (
          <button
            type="button"
            onClick={handleGenerateAISummary}
            disabled={isGeneratingAISummary}
            style={{ ...materialToolbarPrimaryButtonStyle, minHeight: '38px', cursor: isGeneratingAISummary ? 'progress' : 'pointer' }}
          >
            {isGeneratingAISummary ? copy.aiSummaryGenerating : pointDetail?.ai_summary_entry?.summary ? copy.aiSummaryRegenerate : copy.aiSummaryGenerate}
          </button>
        ) : null}
      </div>
      {pointDetail?.ai_summary_entry?.summary ? (
        <div style={{ display: 'grid', gap: '8px' }}>
          <div style={{ ...aiSummarySourceStyle, color: themeTokens.metaLabel }}>{pointDetail.ai_summary_entry.source_title || pointDetail.point.title}</div>
          <div style={{ ...aiSummaryTextStyle, color: themeTokens.description }}>{pointDetail.ai_summary_entry.summary}</div>
        </div>
      ) : (
        <div style={{ color: themeTokens.mutedText, fontSize: '15px', lineHeight: 1.65 }}>
          {copy.aiSummaryEmpty}
        </div>
      )}
    </div>
  ) : null;

  return (
    <div
      style={{
        ...readingPanelStyle,
        display: isExplorationPoint ? 'grid' : 'none',
        width: '100%',
        maxWidth: '860px',
        margin: '0 auto',
        background: themeTokens.surfaceBackground,
        borderColor: themeTokens.surfaceBorder,
      }}
    >
      <div style={{ display: 'grid', gap: '8px' }}>
        {isExplorationPoint ? (
          <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.providedContentTitle}</strong>
        ) : null}
        {pointDetail?.point.description ? (
          <p style={{ ...readingBodyStyle, color: themeTokens.description }}>{pointDetail.point.description}</p>
        ) : null}
        {aiPointDescriptionPanel}
      </div>

      {isExplorationPoint ? (
        <div style={{ display: 'grid', gap: '16px' }}>
          {pointDetail?.point.external_url ? (
            <div style={{ display: 'grid', gap: '10px' }}>
              {getYoutubeVideoId(pointDetail.point.external_url) ? (
                <>
                  <p style={{ margin: 0, color: themeTokens.title, fontSize: '13px', fontWeight: 750 }}>{copy.youtubeSourceLabel}</p>
                  <div style={{ position: 'relative', width: '100%', paddingTop: '56.25%', borderRadius: '12px', overflow: 'hidden', border: `1px solid ${themeTokens.surfaceBorder}`, background: '#000' }}>
                    <iframe
                      src={getYoutubeEmbedLink(pointDetail.point.external_url)}
                      title={copy.openYoutubePlayer}
                      loading="lazy"
                      allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                      allowFullScreen
                      style={{
                        position: 'absolute',
                        inset: 0,
                        width: '100%',
                        height: '100%',
                        border: 0,
                        background: '#000',
                      }}
                    />
                  </div>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                    <a href={getYoutubeWatchLink(pointDetail.point.external_url)} target="_blank" rel="noopener noreferrer" onClick={onSourceContentOpen} style={sourceContentButtonStyle}>
                      {copy.openYoutubePlayer}
                    </a>
                    <a href={pointDetail.point.external_url} target="_blank" rel="noopener noreferrer" onClick={onSourceContentOpen} style={sourceSecondaryButtonStyle}>
                      {copy.openOriginalContent}
                    </a>
                  </div>
                  <p style={{ margin: 0, color: themeTokens.mutedText, fontSize: '12px', lineHeight: 1.5, fontWeight: 650 }}>{copy.sourceLayoutHint}</p>
                </>
              ) : (
                <>
                  <a href={pointDetail.point.external_url} target="_blank" rel="noopener noreferrer" onClick={onSourceContentOpen} style={sourceContentButtonStyle}>
                    {copy.openCurrentContent}
                  </a>
                  <p style={{ margin: 0, color: themeTokens.mutedText, fontSize: '12px', lineHeight: 1.5, fontWeight: 650 }}>
                    {copy.sourceLayoutHint}
                  </p>
                </>
              )}
            </div>
          ) : (
            <div style={{ ...emptyCardStyle, background: themeTokens.emptyBackground, borderColor: themeTokens.emptyBorder, color: themeTokens.mutedText }}>{copy.emptyExternalLink}</div>
          )}
          {latestReplacementReport?.replacement_url ? (
            <div style={{ ...noticeCardStyle, background: themeTokens.noticeBackground, borderColor: themeTokens.noticeBorder, color: themeTokens.description }}>
              {copy.replacementCompletedNotice}
            </div>
          ) : null}
          {canEditLearningRecords ? (
            <details style={contentIssueDisclosureStyle}>
              <summary style={contentIssueSummaryStyle}>{copy.issueSummary}</summary>
              <div style={{ ...formGridStyle, marginTop: '14px' }}>
                <div style={sectionMiniHeaderStyle}>
                  <span style={contentIssueMetaStyle}>{copy.issueHelp}</span>
                </div>
                {contentIssueNotice ? (
                  <div style={{ ...noticeCardStyle, background: themeTokens.noticeBackground, borderColor: themeTokens.noticeBorder, color: themeTokens.description }}>{contentIssueNotice}</div>
                ) : null}
                {!hasSourceReplacementGate ? (
                  <div style={{ display: 'grid', gap: '10px' }}>
                    <div style={sectionMiniHeaderStyle}>
                      <strong style={{ ...pointShellTitleStyle, color: themeTokens.title, fontSize: '14px' }}>{copy.reportTitle}</strong>
                      <span style={contentIssueMetaStyle}>{copy.reportHelp}</span>
                    </div>
                    <div style={inputRowStyle}>
                      <select style={{ ...inputStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={materialReportType} onChange={(e) => setMaterialReportType(e.target.value as typeof materialReportType)}>
                        <option value="broken_link">{copy.reportTypeLabels.broken_link}</option>
                        <option value="wrong_content">{copy.reportTypeLabels.wrong_content}</option>
                        <option value="unsafe_content">{copy.reportTypeLabels.unsafe_content}</option>
                        <option value="copyright">{copy.reportTypeLabels.copyright}</option>
                        <option value="low_quality">{copy.reportTypeLabels.low_quality}</option>
                        <option value="other">{copy.reportTypeLabels.other}</option>
                      </select>
                    </div>
                    <textarea style={{ ...textareaStyle, minHeight: '76px', background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={materialReportMessage} onChange={(e) => setMaterialReportMessage(e.target.value)} placeholder={copy.reportPlaceholder} />
                    <div style={formActionRowStyle}>
                      <button type="button" onClick={() => { setContentIssueNotice(''); void handleReportMaterial(); }} disabled={isReportingMaterial || materialReportMessage.trim().length < 5} style={{ ...materialToolbarButtonStyle, opacity: isReportingMaterial || materialReportMessage.trim().length < 5 ? 0.62 : 1, cursor: isReportingMaterial ? 'progress' : materialReportMessage.trim().length < 5 ? 'default' : 'pointer' }}>
                        {isReportingMaterial ? copy.reportSubmitting : copy.reportAction}
                      </button>
                    </div>
                  </div>
                ) : null}
                {!hasSourceReplacementGate ? (
                  <div style={{ ...noticeCardStyle, background: themeTokens.noticeBackground, borderColor: themeTokens.noticeBorder, color: themeTokens.description }}>
                    {copy.reportGateHelp}
                  </div>
                ) : null}
                {hasSourceReplacementGate ? (
                  <>
                    <div style={{ ...noticeCardStyle, background: themeTokens.noticeBackground, borderColor: themeTokens.noticeBorder, color: themeTokens.description }}>
                      {copy.replacementGateNotice}
                    </div>
                    <div style={{ ...formActionRowStyle, justifyContent: 'space-between' }}>
                      <button type="button" onClick={() => void handleCancelSourceReplacementReport()} disabled={isReportingMaterial || isSavingAttachment} style={{ ...materialToolbarButtonStyle, opacity: isReportingMaterial || isSavingAttachment ? 0.62 : 1, cursor: isReportingMaterial || isSavingAttachment ? 'progress' : 'pointer' }}>
                        {copy.cancelReport}
                      </button>
                      <button type="button" onClick={() => { setContentIssueNotice(''); void handleLoadReplacementCandidates(); }} disabled={isLoadingReplacementCandidates || isSavingAttachment} style={{ ...materialToolbarPrimaryButtonStyle, cursor: isLoadingReplacementCandidates ? 'progress' : 'pointer', opacity: isLoadingReplacementCandidates || isSavingAttachment ? 0.68 : 1 }}>
                        {isLoadingReplacementCandidates ? copy.replacementLoading : copy.loadReplacement}
                      </button>
                    </div>
                    {replacementMessage ? (
                      <div style={{ ...noticeCardStyle, background: themeTokens.noticeBackground, borderColor: themeTokens.noticeBorder, color: themeTokens.description }}>{replacementMessage}</div>
                    ) : null}
                    {replacementCandidates.length ? (
                      <div role="table" aria-label={copy.replacementCandidatesAria} style={listBoardStyle}>
                        <div role="row" style={getListHeaderRowStyle('minmax(0, 1fr) 118px')}>
                          <span role="columnheader" style={listColumnHeaderStyle}>{copy.candidateColumn}</span>
                          <span role="columnheader" style={listColumnHeaderStyle}>{copy.replaceColumn}</span>
                        </div>
                        {replacementCandidates.map((candidate, index) => (
                          <div key={`${candidate.url ?? candidate.title}-${index}`} role="row" style={{ display: 'grid', gridTemplateColumns: 'minmax(0, 1fr) 118px', gap: '10px', alignItems: 'center', padding: '12px 16px', borderBottom: `1px solid ${themeTokens.surfaceBorder}` }}>
                            <div role="cell" style={{ minWidth: 0, display: 'grid', gap: '4px' }}>
                              {candidate.url ? (
                                <a href={candidate.url} target="_blank" rel="noopener noreferrer" style={{ color: themeTokens.title, fontSize: '14px', lineHeight: 1.45, fontWeight: 800, textDecoration: 'underline', textUnderlineOffset: '3px' }}>
                                  {candidate.title || copy.untitledCandidate}
                                </a>
                              ) : (
                                <strong style={{ color: themeTokens.title, fontSize: '14px', lineHeight: 1.45 }}>{candidate.title || copy.untitledCandidate}</strong>
                              )}
                              {candidate.description ? (
                                <span style={{ color: themeTokens.description, fontSize: '13px', lineHeight: 1.45, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{candidate.description}</span>
                              ) : null}
                              {candidate.url ? (
                                <span style={{ color: themeTokens.metaLabel, fontSize: '12px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{candidate.url}</span>
                              ) : null}
                            </div>
                            <button type="button" onClick={() => void handleUseReplacementCandidate(candidate, activeSourceReplacementReportID)} disabled={isSavingAttachment || !candidate.url} style={{ ...materialToolbarButtonStyle, justifyContent: 'center', opacity: isSavingAttachment || !candidate.url ? 0.62 : 1, cursor: isSavingAttachment ? 'progress' : candidate.url ? 'pointer' : 'default' }}>
                              {copy.replace}
                            </button>
                          </div>
                        ))}
                      </div>
                    ) : null}
                    <div style={{ display: 'grid', gap: '10px', paddingTop: '4px' }}>
                      <div style={sectionMiniHeaderStyle}>
                        <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.directReplacementTitle}</strong>
                        <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.directReplacementHelp}</span>
                      </div>
                      <div style={inputRowStyle}>
                        <input style={{ ...inputStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={replacementLinkTitle} onChange={(e) => setReplacementLinkTitle(e.target.value)} placeholder={copy.directReplacementTitlePlaceholder} />
                        <input style={{ ...inputStyle, background: themeTokens.inputBackground, borderColor: themeTokens.inputBorder, color: themeTokens.inputText }} value={replacementLinkURL} onChange={(e) => setReplacementLinkURL(e.target.value)} placeholder="https://..." />
                      </div>
                      <div style={formActionRowStyle}>
                        <button type="button" onClick={() => void handleCreateDirectReplacementLink()} disabled={isSavingAttachment || !replacementLinkTitle.trim() || !replacementLinkURL.trim()} style={{ ...materialToolbarButtonStyle, opacity: isSavingAttachment || !replacementLinkTitle.trim() || !replacementLinkURL.trim() ? 0.62 : 1, cursor: isSavingAttachment ? 'progress' : replacementLinkTitle.trim() && replacementLinkURL.trim() ? 'pointer' : 'default' }}>
                          {copy.directReplacementAction}
                        </button>
                      </div>
                    </div>
                  </>
                ) : null}
                <div style={{ display: 'grid', gap: '10px' }}>
                  <div style={sectionMiniHeaderStyle}>
                    <strong style={{ ...pointShellTitleStyle, color: themeTokens.title }}>{copy.myReportsTitle}</strong>
                    <span style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{copy.reportCount(materialReports.length)}</span>
                  </div>
                  {materialReports.length ? (
                    <div role="table" aria-label={copy.reportHistoryAria} style={listBoardStyle}>
                      <div role="row" style={getListHeaderRowStyle('minmax(0, 1fr) 92px 132px')}>
                        <span role="columnheader" style={listColumnHeaderStyle}>{copy.reportColumn}</span>
                        <span role="columnheader" style={listColumnHeaderStyle}>{copy.statusColumn}</span>
                        <span role="columnheader" style={listColumnHeaderStyle}>{copy.receivedAtColumn}</span>
                      </div>
                      {materialReports.map((report) => (
                        <div key={report.id} role="row" style={{ display: 'grid', gridTemplateColumns: 'minmax(0, 1fr) 92px 132px', gap: '10px', alignItems: 'center', minHeight: '44px', padding: '9px 16px', borderBottom: `1px solid ${themeTokens.surfaceBorder}` }}>
                          <div role="cell" style={{ minWidth: 0, display: 'grid', gap: '2px' }}>
                            <strong style={{ color: themeTokens.title, fontSize: '13px', lineHeight: 1.35, fontWeight: 800, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                              {copy.reportTypeLabels[report.report_type] ?? report.report_type}
                            </strong>
                            <span style={{ color: themeTokens.description, fontSize: '12px', lineHeight: 1.35, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                              {report.replaced_at && report.replacement_url ? copy.replacedPrefix(report.replacement_title || report.replacement_url) : report.message}
                            </span>
                          </div>
                          <span role="cell" style={{ ...blockMetaStyle, color: themeTokens.metaLabel, fontWeight: 800 }}>{copy.reportStatusLabels[report.status] ?? report.status}</span>
                          <span role="cell" style={{ ...blockMetaStyle, color: themeTokens.metaLabel }}>{formatMaterialReportTime(report.created_at)}</span>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div style={{ ...emptyCardStyle, background: themeTokens.emptyBackground, borderColor: themeTokens.emptyBorder, color: themeTokens.mutedText }}>{copy.emptyReports}</div>
                  )}
                </div>
              </div>
            </details>
          ) : null}
        </div>
      ) : null}

      {pointRuntimeMessage ? <div style={{ ...noticeCardStyle, background: themeTokens.noticeBackground, borderColor: themeTokens.noticeBorder, color: themeTokens.description }}>{pointRuntimeMessage}</div> : null}
      {pointEntryMessage ? <div style={{ ...noticeCardStyle, background: themeTokens.noticeBackground, borderColor: themeTokens.noticeBorder, color: themeTokens.description }}>{pointEntryMessage}</div> : null}
    </div>
  );
}
