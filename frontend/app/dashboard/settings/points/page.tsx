'use client';

import type React from 'react';
import { useRouter } from 'next/navigation';

import DashboardSettingsShell, {
  loadingCardStyle,
  loadingPageStyle,
  pillStyle,
  primaryActionButtonStyle,
  profileMetaCardStyle,
  profileMetaGridStyle,
  profileMetaLabelStyle,
  profileMetaValueStyle,
  secondaryActionButtonStyle,
  sectionCardStyle,
  sectionHeaderStyle,
  sectionSubtitleStyle,
  sectionTitleStyle,
  settingsInputStyle,
  summaryRowStyle,
  summaryTextStyle,
} from '@/components/dashboard/settings/DashboardSettingsShell';
import { formatDateTime, type UserPointTransaction, useDashboardSettingsData } from '@/components/dashboard/settings/useDashboardSettingsData';
import { getDashboardSettingsPointsCopy } from '@/lib/i18n/pages/dashboardSettingsPoints';

export default function DashboardSettingsPointsPage() {
  const router = useRouter();
  const {
    user,
    pointUsageTransactions,
    pointUsageTotal,
    pointUsagePage,
    pointUsageTotalPages,
    pointUsageSummary,
    pointUsageStartDate,
    setPointUsageStartDate,
    pointUsageEndDate,
    setPointUsageEndDate,
    isLoadingPointUsage,
    pointUsageMessage,
    isLoading,
    handleLogout,
    handlePointUsageSearch,
    handlePointUsagePageChange,
  } = useDashboardSettingsData(router, '/dashboard/settings/points');
  const copy = getDashboardSettingsPointsCopy(user?.ui_locale);
  const uiLocale = user?.ui_locale === 'en' ? 'en' : 'ko';
  const maxUsageEndDate = getMaxUsageEndDate(pointUsageStartDate);
  const isUsageRangeInvalid = !pointUsageStartDate || !pointUsageEndDate || pointUsageEndDate < pointUsageStartDate || pointUsageEndDate > maxUsageEndDate;
  const usagePageStart = pointUsageTotal === 0 ? 0 : (pointUsagePage - 1) * 10 + 1;
  const usagePageEnd = Math.min(pointUsagePage * 10, pointUsageTotal);

  const handleUsageStartDateChange = (value: string) => {
    setPointUsageStartDate(value);
    if (!value) return;
    const nextMaxEndDate = getMaxUsageEndDate(value);
    if (!pointUsageEndDate || pointUsageEndDate < value) {
      setPointUsageEndDate(value);
    } else if (pointUsageEndDate > nextMaxEndDate) {
      setPointUsageEndDate(nextMaxEndDate);
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
      activeTab="points"
      onLogout={handleLogout}
    >
      <section style={sectionCardStyle}>
        <div style={sectionHeaderStyle}>
          <div>
            <h2 style={sectionTitleStyle}>{copy.points.title}</h2>
            <p style={sectionSubtitleStyle}>{copy.points.subtitle}</p>
          </div>
          <span style={pillStyle(user?.premium_access ? '#EF9F27' : '#48BB78')}>
            {user?.premium_access ? 'LearnCosmos Pro' : copy.points.standard}
          </span>
        </div>

        <div style={summaryRowStyle}>
          <span style={pointBalanceStyle}>{formatPoint(user?.total_points ?? 0, uiLocale)}</span>
          <span style={summaryTextStyle}>{copy.points.account(user?.email)}</span>
        </div>

        <div style={pointBreakdownHeaderStyle}>{copy.points.breakdown}</div>
        <div style={profileMetaGridStyle}>
          <div style={profileMetaCardStyle}>
            <span style={profileMetaLabelStyle}>{copy.points.free}</span>
            <span style={profileMetaValueStyle}>{formatPoint(user?.free_points ?? 0, uiLocale)}</span>
          </div>
          <div style={profileMetaCardStyle}>
            <span style={profileMetaLabelStyle}>{copy.points.paid}</span>
            <span style={profileMetaValueStyle}>{formatPoint(user?.paid_points ?? 0, uiLocale)}</span>
          </div>
        </div>
      </section>

      <section style={sectionCardStyle}>
        <div style={sectionHeaderStyle}>
          <div>
            <h2 style={sectionTitleStyle}>{copy.usage.title}</h2>
            <p style={sectionSubtitleStyle}>{copy.usage.subtitle}</p>
          </div>
        </div>

        <div style={usageFilterPanelStyle}>
          <label style={usageDateLabelStyle}>
            {copy.usage.startDate}
            <input type="date" value={pointUsageStartDate} onChange={(event) => handleUsageStartDateChange(event.target.value)} style={usageDateInputStyle} />
          </label>
          <label style={usageDateLabelStyle}>
            {copy.usage.endDate}
            <input type="date" value={pointUsageEndDate} min={pointUsageStartDate || undefined} max={maxUsageEndDate || undefined} onChange={(event) => setPointUsageEndDate(event.target.value)} style={usageDateInputStyle} />
          </label>
          <button type="button" onClick={handlePointUsageSearch} style={primaryActionButtonStyle(isLoadingPointUsage || isUsageRangeInvalid)} disabled={isLoadingPointUsage || isUsageRangeInvalid}>
            {isLoadingPointUsage ? copy.usage.searching : copy.usage.search}
          </button>
          <div style={usageRangeHintStyle}>
            {isUsageRangeInvalid ? copy.usage.invalidRange : copy.usage.rangeResult(pointUsageStartDate, pointUsageEndDate, pointUsageTotal)}
          </div>
        </div>

        {pointUsageMessage ? <div style={usageErrorStyle}>{pointUsageMessage}</div> : null}

        <div style={usageSummaryGridStyle}>
          <SummaryItem label={copy.usage.granted} value={formatSignedPoint(pointUsageSummary.granted_points, uiLocale, 'positive')} tone="positive" />
          <SummaryItem label={copy.usage.purchased} value={formatSignedPoint(pointUsageSummary.purchased_points, uiLocale, 'positive')} tone="positive" />
          <SummaryItem label={copy.usage.used} value={formatSignedPoint(pointUsageSummary.used_points, uiLocale, 'negative')} tone="negative" />
          <SummaryItem label={copy.usage.refunded} value={formatSignedPoint(pointUsageSummary.refunded_points, uiLocale, 'positive')} tone="positive" />
          <SummaryItem label={copy.usage.netChange} value={formatSignedPoint(pointUsageSummary.net_change, uiLocale)} tone={pointUsageSummary.net_change < 0 ? 'negative' : 'positive'} />
        </div>

        {pointUsageTransactions.length > 0 ? (
          <div style={usageListStyle}>
            {pointUsageTransactions.map((transaction) => (
              <div key={transaction.id} style={usageItemStyle}>
                <div style={usageItemMainStyle}>
                  <span style={usageFeatureStyle}>{getPointTransactionTitle(transaction, copy.usage.transactionTypes)}</span>
                  <span style={usageMetaStyle}>{getPointTransactionMeta(transaction, uiLocale)}</span>
                </div>
                <div style={usageAmountStyle(transaction.amount)}>{formatSignedPoint(transaction.amount, uiLocale)}</div>
              </div>
            ))}
          </div>
        ) : (
          <div style={usageEmptyStyle}>{isLoadingPointUsage ? copy.usage.loading : copy.usage.empty}</div>
        )}

        <div style={usagePaginationStyle}>
          <span style={usagePaginationTextStyle}>{pointUsageTotal > 0 ? copy.usage.count(usagePageStart, usagePageEnd, pointUsageTotal) : copy.usage.zero}</span>
          <div style={usagePaginationButtonRowStyle}>
            <button type="button" style={secondaryActionButtonStyle(isLoadingPointUsage || pointUsagePage <= 1)} disabled={isLoadingPointUsage || pointUsagePage <= 1} onClick={() => handlePointUsagePageChange(pointUsagePage - 1)}>
              {copy.usage.previous}
            </button>
            <span style={usagePaginationTextStyle}>{Math.max(pointUsagePage, 1)} / {Math.max(pointUsageTotalPages, 1)}</span>
            <button type="button" style={secondaryActionButtonStyle(isLoadingPointUsage || pointUsageTotalPages === 0 || pointUsagePage >= pointUsageTotalPages)} disabled={isLoadingPointUsage || pointUsageTotalPages === 0 || pointUsagePage >= pointUsageTotalPages} onClick={() => handlePointUsagePageChange(pointUsagePage + 1)}>
              {copy.usage.next}
            </button>
          </div>
        </div>
      </section>
    </DashboardSettingsShell>
  );
}

function SummaryItem({ label, value, tone }: { label: string; value: string; tone: 'positive' | 'negative' }) {
  return (
    <div style={usageSummaryItemStyle(tone)}>
      <span style={usageSummaryLabelStyle(tone)}>{label}</span>
      <strong style={usageSummaryValueStyle(tone)}>{value}</strong>
    </div>
  );
}

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
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return year + '-' + month + '-' + day;
}

function getPointTransactionTitle(transaction: UserPointTransaction, labels: Record<string, string>) {
  const description = typeof transaction.description === 'string' ? transaction.description.trim() : '';
  return description || labels[transaction.type] || transaction.type;
}

function getPointTransactionMeta(transaction: UserPointTransaction, locale: 'ko' | 'en') {
  const parts = [formatDateTime(transaction.created_at, locale)];
  const feature = typeof transaction.feature === 'string' ? transaction.feature.trim() : '';
  const referenceType = typeof transaction.reference_type === 'string' ? transaction.reference_type.trim() : '';
  if (feature) parts.push(formatPointFeature(feature, locale));
  if (referenceType) parts.push(formatReferenceType(referenceType, locale));
  return parts.join(' · ');
}

function formatPointFeature(feature: string, locale: 'ko' | 'en') {
  const labels: Record<string, { ko: string; en: string }> = {
    welcome_points: { ko: '가입 지급', en: 'Welcome grant' },
    course_draft_create: { ko: '탐험계획 생성', en: 'Plan generation' },
    course_draft_rebuild_all: { ko: '탐험계획 재구성', en: 'Plan rebuild' },
    course_draft_rebuild: { ko: '탐험계획 재구성', en: 'Plan rebuild' },
    lesson_candidate_search: { ko: '리슨 추천 검색', en: 'Lesson search' },
    explorer_recommendation_search: { ko: '탐험지점 추천 검색', en: 'Explorer search' },
    admin_point_adjustment: { ko: '관리자 조정', en: 'Admin adjustment' },
  };
  return labels[feature]?.[locale] ?? feature;
}

function formatReferenceType(referenceType: string, locale: 'ko' | 'en') {
  const labels: Record<string, { ko: string; en: string }> = {
    course_draft: { ko: '탐험계획', en: 'Plan' },
    admin_action: { ko: '운영 처리', en: 'Admin action' },
  };
  return labels[referenceType]?.[locale] ?? referenceType;
}

function formatPoint(value: number, locale: 'ko' | 'en') {
  return value.toLocaleString(locale === 'en' ? 'en-US' : 'ko-KR') + 'pt';
}

function formatSignedPoint(value: number, locale: 'ko' | 'en', forcedSign?: 'positive' | 'negative') {
  const normalized = forcedSign === 'negative' ? -Math.abs(value) : forcedSign === 'positive' ? Math.abs(value) : value;
  const sign = normalized > 0 ? '+' : '';
  return sign + formatPoint(normalized, locale);
}

const pointBalanceStyle = {
  fontSize: '34px',
  lineHeight: 1.15,
  fontWeight: 900,
  color: '#160E08',
} satisfies React.CSSProperties;

const pointBreakdownHeaderStyle = {
  marginTop: '4px',
  paddingTop: '16px',
  borderTop: '1px solid rgba(146, 111, 62, 0.14)',
  color: '#4C4033',
  fontSize: '14px',
  fontWeight: 800,
} satisfies React.CSSProperties;

const usageFilterPanelStyle: React.CSSProperties = {
  display: 'flex',
  gap: 10,
  flexWrap: 'wrap',
  alignItems: 'flex-end',
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
  gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))',
  gap: 10,
};

const usageSummaryItemStyle = (tone: 'positive' | 'negative'): React.CSSProperties => ({
  display: 'grid',
  gap: 6,
  minHeight: 82,
  padding: '14px 16px',
  borderRadius: 8,
  border: '1px solid ' + (tone === 'negative' ? 'rgba(185, 28, 28, 0.18)' : 'rgba(4, 120, 87, 0.16)'),
  background: tone === 'negative' ? '#FEF2F2' : '#F0FDF4',
});

const usageSummaryLabelStyle = (tone: 'positive' | 'negative'): React.CSSProperties => ({
  color: tone === 'negative' ? '#991B1B' : '#065F46',
  fontSize: 13,
  lineHeight: 1.4,
  fontWeight: 800,
});

const usageSummaryValueStyle = (tone: 'positive' | 'negative'): React.CSSProperties => ({
  color: tone === 'negative' ? '#7F1D1D' : '#064E3B',
  fontSize: 20,
  lineHeight: 1.2,
  fontWeight: 900,
});

const usageListStyle: React.CSSProperties = {
  display: 'grid',
  gap: 10,
};

const usageItemStyle: React.CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 12,
  flexWrap: 'wrap',
  padding: '14px 16px',
  borderRadius: 8,
  border: '1px solid rgba(17, 24, 39, 0.12)',
  background: '#F9FAFB',
};

const usageItemMainStyle: React.CSSProperties = {
  display: 'grid',
  gap: 5,
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

const usageAmountStyle = (amount: number): React.CSSProperties => ({
  color: amount < 0 ? '#B91C1C' : '#047857',
  fontSize: 16,
  fontWeight: 900,
  whiteSpace: 'nowrap',
});

const usageEmptyStyle: React.CSSProperties = {
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
