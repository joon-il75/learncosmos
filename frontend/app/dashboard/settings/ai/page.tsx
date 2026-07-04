'use client';

import type React from 'react';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

import DashboardSettingsShell, {
  fieldGridStyle,
  fieldLabelStyle,
  helperTextStyle,
  inlineInfoStyle,
  loadingCardStyle,
  loadingPageStyle,
  pillStyle,
  primaryActionButtonStyle,
  secondaryActionButtonStyle,
  sectionCardStyle,
  sectionHeaderStyle,
  sectionSubtitleStyle,
  sectionTitleStyle,
  settingsInputStyle,
  summaryRowStyle,
  summaryTextStyle,
} from '@/components/dashboard/settings/DashboardSettingsShell';
import { AI_PROVIDERS, formatDateTime, useDashboardSettingsData } from '@/components/dashboard/settings/useDashboardSettingsData';
import { getDashboardSettingsAICopy, type DashboardSettingsAICopy } from '@/lib/i18n/pages/dashboardSettingsAI';

export default function DashboardSettingsAIPage() {
  const router = useRouter();
  const [isDeletePanelOpen, setIsDeletePanelOpen] = useState(false);
  const [deleteConfirmText, setDeleteConfirmText] = useState('');
  const {
    user,
    aiSettings,
    aiUsageEvents,
    aiUsageTotal,
    aiUsagePage,
    aiUsageTotalPages,
    aiUsageSummary,
    aiUsageStartDate,
    setAIUsageStartDate,
    aiUsageEndDate,
    setAIUsageEndDate,
    isLoadingAIUsage,
    aiUsageMessage,
    aiProvider,
    setAIProvider,
    aiAPIKey,
    setAIAPIKey,
    aiEndpointURL,
    setAIEndpointURL,
    isValidatingAI,
    isSavingAI,
    isTogglingAIEnabled,
    isDeletingAIKey,
    isCurrentAIKeyValidated,
    aiKeyMessage,
    isLoading,
    handleLogout,
    handleValidateAI,
    handleSaveAI,
    handleToggleAIEnabled,
    handleDeleteAIKey,
    handleAIUsageSearch,
    handleAIUsagePageChange,
  } = useDashboardSettingsData(router, '/dashboard/settings/ai');
  const copy = getDashboardSettingsAICopy(user?.ui_locale);
  const uiLocale = user?.ui_locale === 'en' ? 'en' : 'ko';

  const requiresEndpointURL = aiProvider === 'llama';
  const hasValidKeyStatus = aiSettings?.last_validation_status === 'valid';
  const hasSavedValidKeyStatus = Boolean(aiSettings?.has_api_key && hasValidKeyStatus);
  const savedProviderLabel = getAIProviderLabel(aiSettings?.provider || aiProvider);
  const isDeleteConfirmed = deleteConfirmText.trim() === copy.byok.deleteConfirmValue;
  const maxUsageEndDate = getMaxUsageEndDate(aiUsageStartDate);
  const isUsageRangeInvalid =
    !aiUsageStartDate ||
    !aiUsageEndDate ||
    aiUsageEndDate < aiUsageStartDate ||
    aiUsageEndDate > maxUsageEndDate;
  const usagePageStart = aiUsageTotal === 0 ? 0 : (aiUsagePage - 1) * 10 + 1;
  const usagePageEnd = Math.min(aiUsagePage * 10, aiUsageTotal);

  useEffect(() => {
    if (!aiSettings?.has_api_key || !isDeletePanelOpen) setDeleteConfirmText('');
  }, [aiSettings?.has_api_key, isDeletePanelOpen]);

  const handleUsageStartDateChange = (value: string) => {
    setAIUsageStartDate(value);
    if (!value) return;
    const nextMaxEndDate = getMaxUsageEndDate(value);
    if (!aiUsageEndDate || aiUsageEndDate < value) {
      setAIUsageEndDate(value);
    } else if (aiUsageEndDate > nextMaxEndDate) {
      setAIUsageEndDate(nextMaxEndDate);
    }
  };

  if (isLoading) {
    return (
      <div style={loadingPageStyle}>
        <div style={loadingCardStyle}>{copy.loading}</div>
      </div>
    );
  }

  return (
    <DashboardSettingsShell
      eyebrow={copy.shell.eyebrow}
      title={copy.shell.title}
      copy={copy.shell.copy}
      tabLabels={copy.shell.tabs}
      uiLocale={uiLocale}
      activeTab="ai"
      onLogout={handleLogout}
    >
      <section style={sectionCardStyle}>
        <div style={sectionHeaderStyle}>
          <div>
            <h2 style={sectionTitleStyle}>{copy.byok.title}</h2>
            <p style={sectionSubtitleStyle}>{copy.byok.subtitle}</p>
          </div>
          <span style={pillStyle(aiSettings?.has_api_key ? (aiSettings.is_enabled ? '#48BB78' : '#9CA3AF') : '#8FA7D6')}>
            {aiSettings?.has_api_key ? (aiSettings.is_enabled ? copy.byok.statusUsing : copy.byok.statusDisabled) : copy.byok.statusManaged}
          </span>
        </div>

        <div style={summaryRowStyle}>
          <span style={summaryTextStyle}>
            {aiSettings?.has_api_key
              ? copy.byok.summarySaved(savedProviderLabel, aiSettings.is_enabled)
              : copy.byok.summaryManaged}
          </span>
          {aiSettings?.has_api_key && aiSettings?.last_validated_at ? <span style={summaryTextStyle}>{copy.byok.recentlyValidated(formatDateTime(aiSettings.last_validated_at, uiLocale))}</span> : null}
        </div>

        <div style={aiStatusPanelStyle}>
          <div>
            <div style={byokControlTitleStyle}>{copy.byok.usageStateTitle}</div>
            <div style={byokControlTextStyle}>
              {aiSettings?.has_api_key
                ? aiSettings.is_enabled
                  ? copy.byok.usageEnabled
                  : copy.byok.usageDisabled
                : copy.byok.usageEmpty}
            </div>
          </div>
          <label style={toggleWrapStyle(aiSettings?.has_api_key && aiSettings.is_enabled)}>
            <input
              type="checkbox"
              checked={Boolean(aiSettings?.has_api_key && aiSettings.is_enabled)}
              disabled={!aiSettings?.has_api_key || isTogglingAIEnabled}
              onChange={(event) => handleToggleAIEnabled(event.target.checked)}
              style={toggleInputStyle}
            />
            <span style={toggleKnobStyle(aiSettings?.has_api_key && aiSettings.is_enabled)} />
          </label>
        </div>

        <div style={formPanelStyle}>
          <div style={panelHeadingStyle}>{copy.byok.registrationTitle}</div>
          {aiSettings?.has_api_key ? (
            <div style={savedKeyInfoStyle}>
              <div style={savedKeyInfoTitleStyle}>{copy.byok.savedKeyTitle}</div>
              <div style={savedKeyInfoGridStyle}>
                <span>{copy.byok.savedProvider(savedProviderLabel)}</span>
                <span>{hasSavedValidKeyStatus ? copy.byok.savedValid : copy.byok.validationStatus(aiSettings.last_validation_status || '')}</span>
                {aiSettings.last_validated_at ? <span>{copy.byok.recentlyValidated(formatDateTime(aiSettings.last_validated_at, uiLocale))}</span> : null}
              </div>
            </div>
          ) : null}
          <div style={fieldGridStyle}>
            <label style={fieldLabelStyle}>
              {copy.byok.providerLabel}
              <select
                value={aiProvider}
                onChange={(event) => {
                  const nextProvider = event.target.value;
                  setAIProvider(nextProvider);
                  if (nextProvider !== 'llama') setAIEndpointURL('');
                }}
                style={settingsInputStyle}
              >
                {AI_PROVIDERS.map((provider) => (
                  <option key={provider.value} value={provider.value}>
                    {provider.label}
                  </option>
                ))}
              </select>
            </label>

            {requiresEndpointURL ? (
              <label style={fieldLabelStyle}>
                {copy.byok.endpointLabel}
                <input
                  type="text"
                  value={aiEndpointURL}
                  onChange={(event) => setAIEndpointURL(event.target.value)}
                  placeholder={copy.byok.endpointPlaceholder}
                  style={settingsInputStyle}
                />
              </label>
            ) : null}

            <label style={fieldLabelStyle}>
              {copy.byok.apiKeyLabel}
              <input
                type="password"
                value={aiAPIKey}
                onChange={(event) => setAIAPIKey(event.target.value)}
                placeholder="sk-..."
                style={settingsInputStyle}
              />
            </label>
          </div>
          <div style={keyActionPanelStyle}>
            <div style={keyActionStatusStyle(isCurrentAIKeyValidated)}>
              <span style={keyActionStatusMarkStyle(isCurrentAIKeyValidated)}>{isCurrentAIKeyValidated ? '✓' : '!'}</span>
              <span>
                {isValidatingAI
                  ? copy.byok.validatingStatus
                  : aiKeyMessage
                    ? aiKeyMessage
                    : isCurrentAIKeyValidated
                        ? copy.byok.validStatus
                        : copy.byok.needsValidationStatus}
              </span>
            </div>
            <div style={keyActionRowStyle}>
              <button type="button" onClick={handleValidateAI} style={secondaryActionButtonStyle(isValidatingAI)} disabled={isValidatingAI}>
                {isValidatingAI ? copy.byok.validating : copy.byok.validate}
              </button>
              <button type="button" onClick={handleSaveAI} style={primaryActionButtonStyle(isSavingAI || !isCurrentAIKeyValidated)} disabled={isSavingAI || !isCurrentAIKeyValidated}>
                {isSavingAI ? copy.byok.saving : copy.byok.save}
              </button>
            </div>
          </div>
          <div style={formHelpStyle}>{copy.byok.inputHelp}</div>
          {aiSettings?.has_api_key ? (
            <details style={deletePanelStyle} onToggle={(event) => setIsDeletePanelOpen(event.currentTarget.open)}>
              <summary style={deleteSummaryStyle}>
                <span>{copy.byok.deleteSummary}</span>
                <span style={deleteSummaryHintStyle} aria-hidden="true">{isDeletePanelOpen ? '▴' : '▾'}</span>
              </summary>
              <div style={deletePanelBodyStyle}>
                <div>
                  <div style={deleteTitleStyle}>{copy.byok.deleteTitle}</div>
                  <div style={deleteTextStyle}>{copy.byok.deleteText}</div>
                </div>
                <label style={deleteConfirmLabelStyle}>
                  {copy.byok.deleteConfirmLabel}
                  <input
                    type="text"
                    value={deleteConfirmText}
                    onChange={(event) => setDeleteConfirmText(event.target.value)}
                    placeholder={copy.byok.deleteConfirmValue}
                    aria-label={copy.byok.deleteConfirmAria}
                    style={deleteConfirmInputStyle}
                  />
                </label>
                <button type="button" onClick={handleDeleteAIKey} style={deleteButtonStyle(isDeletingAIKey || !isDeleteConfirmed)} disabled={isDeletingAIKey || !isDeleteConfirmed}>
                  {isDeletingAIKey ? copy.byok.deleting : copy.byok.deleteAction}
                </button>
              </div>
            </details>
          ) : null}
        </div>

        <div style={usageNoticeStyle}>
          <div style={usageNoticeTitleStyle}>{copy.usage.noticeTitle}</div>
          <p style={usageNoticeTextStyle}>
            {copy.usage.noticeBody1}
          </p>
          <p style={usageNoticeTextStyle}>
            {copy.usage.noticeBody2}
          </p>
        </div>

        <div style={usageHistoryStyle}>
          <div style={usageHistoryHeaderStyle}>
            <div>
              <div style={usageNoticeTitleStyle}>{copy.usage.historyTitle}</div>
              <p style={usageHistorySubStyle}>{copy.usage.historySubtitle}</p>
            </div>
          </div>
          <div style={usageFilterPanelStyle}>
            <label style={usageDateLabelStyle}>
              {copy.usage.startDate}
              <input
                type="date"
                value={aiUsageStartDate}
                onChange={(event) => handleUsageStartDateChange(event.target.value)}
                style={usageDateInputStyle}
              />
            </label>
            <label style={usageDateLabelStyle}>
              {copy.usage.endDate}
              <input
                type="date"
                value={aiUsageEndDate}
                min={aiUsageStartDate || undefined}
                max={maxUsageEndDate || undefined}
                onChange={(event) => setAIUsageEndDate(event.target.value)}
                style={usageDateInputStyle}
              />
            </label>
            <button
              type="button"
              onClick={handleAIUsageSearch}
              style={primaryActionButtonStyle(isLoadingAIUsage || isUsageRangeInvalid)}
              disabled={isLoadingAIUsage || isUsageRangeInvalid}
            >
              {isLoadingAIUsage ? copy.usage.searching : copy.usage.search}
            </button>
            <div style={usageRangeHintStyle}>
              {isUsageRangeInvalid
                ? copy.usage.invalidRange
                : copy.usage.rangeResult(aiUsageStartDate, aiUsageEndDate, aiUsageTotal)}
            </div>
          </div>
          {aiUsageMessage ? <div style={usageErrorStyle}>{aiUsageMessage}</div> : null}
          <div style={usageSummaryGridStyle}>
            <div style={usageSummaryItemStyle}>
              <span style={usageSummaryLabelStyle}>{copy.usage.totalTokens}</span>
              <strong style={usageSummaryValueStyle}>{formatSummaryTokenCount(aiUsageSummary.total_tokens, uiLocale)}</strong>
            </div>
            <div style={usageSummaryItemStyle}>
              <span style={usageSummaryLabelStyle}>{copy.usage.dailyTokens}</span>
              <strong style={usageSummaryValueStyle}>{formatSummaryTokenCount(aiUsageSummary.average_daily_tokens, uiLocale)}</strong>
            </div>
            <div style={usageSummaryItemStyle}>
              <span style={usageSummaryLabelStyle}>{copy.usage.totalCost}</span>
              <strong style={usageSummaryValueStyle}>{formatSummaryCost(aiUsageSummary.total_estimated_cost_usd, copy)}</strong>
            </div>
            <div style={usageSummaryItemStyle}>
              <span style={usageSummaryLabelStyle}>{copy.usage.dailyCost}</span>
              <strong style={usageSummaryValueStyle}>{formatSummaryCost(aiUsageSummary.average_daily_estimated_cost_usd, copy)}</strong>
            </div>
          </div>
          {aiUsageEvents.length > 0 ? (
            <div style={usageListStyle}>
              {aiUsageEvents.map((event) => (
                <div key={event.id} style={usageItemStyle}>
                  <div style={usageItemMainStyle}>
                    <span style={usageFeatureStyle}>{formatAIUsageFeature(event.feature, copy)}</span>
                    <span style={usageMetaStyle}>{formatDateTime(event.created_at, uiLocale)} · {getAIProviderLabel(event.provider)}{event.model ? ` · ${event.model}` : ''}</span>
                  </div>
                  <div style={usageStatsStyle}>
                    <span>{copy.usage.input} {formatTokenCount(event.input_tokens, uiLocale, copy)}</span>
                    <span>{copy.usage.output} {formatTokenCount(event.output_tokens, uiLocale, copy)}</span>
                    <span>{formatEstimatedCost(event.estimated_cost_usd, copy)}</span>
                    <span style={event.success ? usageSuccessStyle : usageFailStyle}>{event.success ? copy.usage.success : copy.usage.fail}</span>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div style={usageEmptyStyle}>
              {isLoadingAIUsage ? copy.usage.loading : copy.usage.empty}
            </div>
          )}
          <div style={usagePaginationStyle}>
            <span style={usagePaginationTextStyle}>
              {aiUsageTotal > 0
                ? copy.usage.count(usagePageStart, usagePageEnd, aiUsageTotal)
                : copy.usage.zero}
            </span>
            <div style={usagePaginationButtonRowStyle}>
              <button
                type="button"
                style={secondaryActionButtonStyle(isLoadingAIUsage || aiUsagePage <= 1)}
                disabled={isLoadingAIUsage || aiUsagePage <= 1}
                onClick={() => handleAIUsagePageChange(aiUsagePage - 1)}
              >
                {copy.usage.previous}
              </button>
              <span style={usagePaginationTextStyle}>
                {Math.max(aiUsagePage, 1)} / {Math.max(aiUsageTotalPages, 1)}
              </span>
              <button
                type="button"
                style={secondaryActionButtonStyle(isLoadingAIUsage || aiUsageTotalPages === 0 || aiUsagePage >= aiUsageTotalPages)}
                disabled={isLoadingAIUsage || aiUsageTotalPages === 0 || aiUsagePage >= aiUsageTotalPages}
                onClick={() => handleAIUsagePageChange(aiUsagePage + 1)}
              >
                {copy.usage.next}
              </button>
            </div>
          </div>
        </div>

        {aiSettings?.has_api_key && aiSettings?.last_validation_error ? <div style={inlineInfoStyle('#FFB4A2')}>{aiSettings.last_validation_error}</div> : null}
      </section>
    </DashboardSettingsShell>
  );
}

const BYOK_COST_KRW_PER_USD = 1400;

const byokControlStyle: React.CSSProperties = {
  padding: '20px',
  borderRadius: 8,
  border: '1px solid rgba(15, 23, 42, 0.14)',
  background: '#FFFFFF',
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: 20,
};

const byokControlTitleStyle: React.CSSProperties = {
  color: '#111827',
  fontSize: 18,
  fontWeight: 800,
  marginBottom: 8,
};

const byokControlTextStyle: React.CSSProperties = {
  color: '#1F2937',
  fontSize: 16,
  lineHeight: 1.75,
  maxWidth: 680,
};

const toggleInputStyle: React.CSSProperties = {
  position: 'absolute',
  opacity: 0,
  pointerEvents: 'none',
};

const toggleWrapStyle = (enabled?: boolean): React.CSSProperties => ({
  width: 54,
  minWidth: 54,
  height: 30,
  borderRadius: 999,
  border: `1px solid ${enabled ? '#047857' : '#9CA3AF'}`,
  background: enabled ? '#D1FAE5' : '#E5E7EB',
  position: 'relative',
  cursor: 'pointer',
  marginTop: 2,
});

const toggleKnobStyle = (enabled?: boolean): React.CSSProperties => ({
  position: 'absolute',
  top: 4,
  left: enabled ? 28 : 4,
  width: 20,
  height: 20,
  borderRadius: '50%',
  background: enabled ? '#047857' : '#6B7280',
  transition: 'left 160ms ease, background 160ms ease',
});

const usageNoticeStyle: React.CSSProperties = {
  padding: '20px',
  borderRadius: 8,
  border: '1px solid rgba(13, 49, 78, 0.2)',
  background: '#F8FAFC',
};

const usageNoticeTitleStyle: React.CSSProperties = {
  color: '#111827',
  fontSize: 18,
  fontWeight: 800,
  marginBottom: 10,
};

const usageNoticeTextStyle: React.CSSProperties = {
  color: '#1F2937',
  fontSize: 16,
  lineHeight: 1.75,
  margin: '8px 0',
};

const aiStatusPanelStyle: React.CSSProperties = {
  ...byokControlStyle,
  marginTop: 4,
};

const formPanelStyle: React.CSSProperties = {
  display: 'grid',
  gap: 14,
  padding: '20px',
  borderRadius: 8,
  border: '1px solid rgba(146, 111, 62, 0.18)',
  background: '#FFFBF3',
};

const panelHeadingStyle: React.CSSProperties = {
  color: '#111827',
  fontSize: 18,
  lineHeight: 1.35,
  fontWeight: 800,
};

const savedKeyInfoStyle: React.CSSProperties = {
  display: 'grid',
  gap: 8,
  padding: '14px 16px',
  borderRadius: 8,
  border: '1px solid rgba(4, 120, 87, 0.22)',
  background: '#F0FDF4',
};

const savedKeyInfoTitleStyle: React.CSSProperties = {
  color: '#064E3B',
  fontSize: 15,
  fontWeight: 900,
};

const savedKeyInfoGridStyle: React.CSSProperties = {
  display: 'flex',
  gap: 10,
  flexWrap: 'wrap',
  color: '#065F46',
  fontSize: 14,
  lineHeight: 1.6,
  fontWeight: 800,
};

const formHelpStyle: React.CSSProperties = {
  ...helperTextStyle,
  paddingTop: 2,
  color: '#374151',
  fontSize: 15,
};

const keyActionPanelStyle: React.CSSProperties = {
  display: 'grid',
  gap: 12,
  padding: '14px 16px',
  borderRadius: 8,
  border: '1px solid rgba(4, 120, 87, 0.22)',
  background: '#FFFFFF',
};

const keyActionStatusStyle = (valid: boolean): React.CSSProperties => ({
  display: 'flex',
  alignItems: 'center',
  gap: 10,
  minHeight: 44,
  padding: '10px 12px',
  borderRadius: 8,
  border: `2px solid ${valid ? '#047857' : 'rgba(55, 65, 81, 0.18)'}`,
  background: valid ? '#D1FAE5' : '#F9FAFB',
  color: valid ? '#064E3B' : '#374151',
  fontSize: 16,
  lineHeight: 1.5,
  fontWeight: 900,
});

const keyActionStatusMarkStyle = (valid: boolean): React.CSSProperties => ({
  width: 26,
  minWidth: 26,
  height: 26,
  borderRadius: '50%',
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  background: valid ? '#047857' : '#6B7280',
  color: '#FFFFFF',
  fontSize: 16,
  fontWeight: 900,
});

const keyActionRowStyle: React.CSSProperties = {
  display: 'flex',
  gap: 10,
  flexWrap: 'wrap',
  alignItems: 'center',
};

const usageHistoryStyle: React.CSSProperties = {
  ...usageNoticeStyle,
  background: '#FFFFFF',
};

const usageHistoryHeaderStyle: React.CSSProperties = {
  display: 'flex',
  justifyContent: 'space-between',
  gap: 16,
  alignItems: 'flex-start',
};

const usageHistorySubStyle: React.CSSProperties = {
  margin: '6px 0 0',
  color: '#4B5563',
  fontSize: 15,
  lineHeight: 1.65,
};

const usageFilterPanelStyle: React.CSSProperties = {
  display: 'flex',
  gap: 10,
  flexWrap: 'wrap',
  alignItems: 'flex-end',
  marginTop: 16,
  padding: '14px 16px',
  borderRadius: 8,
  border: '1px solid rgba(17, 24, 39, 0.10)',
  background: '#F8FAFC',
};

const usageDateLabelStyle: React.CSSProperties = {
  display: 'grid',
  gap: 6,
  color: '#374151',
  fontSize: 13,
  fontWeight: 800,
};

const usageDateInputStyle: React.CSSProperties = {
  ...settingsInputStyle,
  minWidth: 160,
  height: 42,
  padding: '0 12px',
};

const usageRangeHintStyle: React.CSSProperties = {
  flexBasis: '100%',
  color: '#4B5563',
  fontSize: 13,
  lineHeight: 1.5,
};

const usageErrorStyle: React.CSSProperties = {
  marginTop: 12,
  padding: '10px 12px',
  borderRadius: 8,
  border: '1px solid rgba(185, 28, 28, 0.18)',
  background: '#FEF2F2',
  color: '#991B1B',
  fontSize: 14,
  fontWeight: 800,
};

const usageSummaryGridStyle: React.CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))',
  gap: 10,
  marginTop: 14,
};

const usageSummaryItemStyle: React.CSSProperties = {
  display: 'grid',
  gap: 6,
  minHeight: 86,
  padding: '14px 16px',
  borderRadius: 8,
  border: '1px solid rgba(4, 120, 87, 0.16)',
  background: '#F0FDF4',
};

const usageSummaryLabelStyle: React.CSSProperties = {
  color: '#065F46',
  fontSize: 13,
  lineHeight: 1.4,
  fontWeight: 800,
};

const usageSummaryValueStyle: React.CSSProperties = {
  color: '#064E3B',
  fontSize: 20,
  lineHeight: 1.2,
  fontWeight: 900,
};

const usageListStyle: React.CSSProperties = {
  display: 'grid',
  gap: 10,
  marginTop: 14,
};

const usageItemStyle: React.CSSProperties = {
  display: 'grid',
  gap: 10,
  padding: '14px 16px',
  borderRadius: 8,
  border: '1px solid rgba(17, 24, 39, 0.12)',
  background: '#F9FAFB',
};

const usageItemMainStyle: React.CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 12,
  flexWrap: 'wrap',
};

const usageFeatureStyle: React.CSSProperties = {
  color: '#111827',
  fontSize: 15,
  fontWeight: 800,
};

const usageMetaStyle: React.CSSProperties = {
  color: '#4B5563',
  fontSize: 13,
  lineHeight: 1.5,
};

const usageStatsStyle: React.CSSProperties = {
  display: 'flex',
  gap: 8,
  flexWrap: 'wrap',
  color: '#374151',
  fontSize: 13,
  fontWeight: 700,
};

const usageSuccessStyle: React.CSSProperties = {
  color: '#047857',
};

const usageFailStyle: React.CSSProperties = {
  color: '#B91C1C',
};

const usageEmptyStyle: React.CSSProperties = {
  marginTop: 14,
  padding: '14px 16px',
  borderRadius: 8,
  border: '1px dashed rgba(17, 24, 39, 0.18)',
  color: '#4B5563',
  fontSize: 15,
  lineHeight: 1.65,
  background: '#F9FAFB',
};

const usagePaginationStyle: React.CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 12,
  flexWrap: 'wrap',
  marginTop: 14,
};

const usagePaginationButtonRowStyle: React.CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: 8,
  flexWrap: 'wrap',
};

const usagePaginationTextStyle: React.CSSProperties = {
  color: '#374151',
  fontSize: 14,
  fontWeight: 800,
};

function getMaxUsageEndDate(startDate: string) {
  const parsed = parseDateInput(startDate);
  if (!parsed) return '';
  parsed.setMonth(parsed.getMonth() + 3);
  parsed.setDate(parsed.getDate() - 1);
  return formatDateInput(parsed);
}

function parseDateInput(value: string) {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(value)) return null;
  const [year, month, day] = value.split('-').map(Number);
  return new Date(year, month - 1, day);
}

function formatDateInput(date: Date) {
  const year = date.getFullYear();
  const month = `${date.getMonth() + 1}`.padStart(2, '0');
  const day = `${date.getDate()}`.padStart(2, '0');
  return `${year}-${month}-${day}`;
}

function formatTokenCount(value: number | undefined, locale: 'ko' | 'en', copy: DashboardSettingsAICopy) {
  if (typeof value !== 'number') return copy.usage.unknownTokens;
  return `${value.toLocaleString(locale === 'en' ? 'en-US' : 'ko-KR')} tokens`;
}

function formatSummaryTokenCount(value: number | undefined, locale: 'ko' | 'en') {
  if (typeof value !== 'number') return '0 tokens';
  return `${Math.round(value).toLocaleString(locale === 'en' ? 'en-US' : 'ko-KR')} tokens`;
}

function formatEstimatedCost(value: number | undefined, copy: DashboardSettingsAICopy) {
  if (typeof value !== 'number') return copy.usage.unknownCost;
  return copy.usage.estimatedCost(formatKRW(value, copy));
}

function formatSummaryCost(value: number | undefined, copy: DashboardSettingsAICopy) {
  if (typeof value !== 'number') return copy.usage.won('0');
  return formatKRW(value, copy);
}

function formatKRW(usdValue: number, copy: DashboardSettingsAICopy) {
  const krwValue = usdValue * BYOK_COST_KRW_PER_USD;
  if (krwValue > 0 && krwValue < 1) return copy.usage.lessThanOneWon;
  return copy.usage.won(Math.round(krwValue).toLocaleString(copy.shell.tabs.profile === 'Profile' ? 'en-US' : 'ko-KR'));
}

function formatAIUsageFeature(feature: string, copy: DashboardSettingsAICopy) {
  return copy.usage.features[feature] ?? feature;
}

function getAIProviderLabel(provider?: string) {
  return AI_PROVIDERS.find((item) => item.value === provider)?.label ?? provider ?? 'OpenAI';
}

const deletePanelStyle: React.CSSProperties = {
  borderRadius: 8,
  border: '1px solid rgba(185, 28, 28, 0.2)',
  background: '#FFF7F7',
  overflow: 'hidden',
};

const deleteSummaryStyle: React.CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 12,
  padding: '14px 16px',
  color: '#7F1D1D',
  fontSize: 15,
  fontWeight: 900,
  cursor: 'pointer',
  listStyle: 'none',
};

const deleteSummaryHintStyle: React.CSSProperties = {
  color: '#991B1B',
  fontSize: 18,
  fontWeight: 800,
  whiteSpace: 'nowrap',
  lineHeight: 1,
};

const deletePanelBodyStyle: React.CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 16,
  padding: '0 16px 16px',
  flexWrap: 'wrap',
};

const deleteTitleStyle: React.CSSProperties = {
  color: '#7F1D1D',
  fontSize: 15,
  fontWeight: 800,
  marginBottom: 4,
};

const deleteTextStyle: React.CSSProperties = {
  color: '#374151',
  fontSize: 14,
  lineHeight: 1.65,
  maxWidth: 620,
};

const deleteConfirmLabelStyle: React.CSSProperties = {
  display: 'grid',
  gap: 6,
  color: '#7F1D1D',
  fontSize: 13,
  fontWeight: 800,
  minWidth: 180,
};

const deleteConfirmInputStyle: React.CSSProperties = {
  ...settingsInputStyle,
  minHeight: 40,
  borderColor: 'rgba(185, 28, 28, 0.32)',
  background: '#FFFFFF',
  color: '#111827',
};

const deleteButtonStyle = (disabled: boolean): React.CSSProperties => ({
  minHeight: 40,
  padding: '0 16px',
  borderRadius: 999,
  border: '1px solid rgba(185, 28, 28, 0.34)',
  background: disabled ? 'rgba(185, 28, 28, 0.12)' : '#B91C1C',
  color: disabled ? '#7F1D1D' : '#FFFFFF',
  fontWeight: 800,
  fontSize: 14,
  cursor: disabled ? 'default' : 'pointer',
  fontFamily: 'inherit',
  whiteSpace: 'nowrap',
});
