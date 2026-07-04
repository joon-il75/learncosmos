'use client';

import type { CSSProperties } from 'react';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { useRouter } from 'next/navigation';
import SuperAdminPanelHeader from '@/components/super-admin/SuperAdminPanelHeader';
import SuperAdminPanelNav from '@/components/super-admin/SuperAdminPanelNav';
import { superAdminShellStyle } from '@/components/super-admin/layout';

const API_BASE = '';
const PAGE_SIZE = 5;

type ReportStatus = 'open' | 'reviewing' | 'resolved' | 'dismissed' | 'cancelled';

interface MaterialReport {
  id: string;
  course_point_id: string;
  user_id: string;
  reporter_display: string;
  reporter_email: string;
  attachment_id?: string;
  target_type: string;
  report_type: string;
  message: string;
  status: ReportStatus;
  created_at: string;
  updated_at: string;
  admin_note: string;
  reviewed_by?: string;
  reviewer_display: string;
  reviewed_at?: string;
  course_id: string;
  course_title: string;
  level_title: string;
  lesson_title: string;
  point_title: string;
  point_type: string;
  point_external_url: string;
  attachment_title: string;
  attachment_url: string;
  attachment_file_path: string;
  attachment_file_size?: number;
  attachment_mime_type: string;
  attachment_source_context: string;
  attachment_type: string;
  recommendation_blocked: boolean;
}

const statusOptions = [
  { value: 'open', label: '열림' },
  { value: 'reviewing', label: '검토중' },
  { value: 'resolved', label: '해결' },
  { value: 'dismissed', label: '반려' },
  { value: 'cancelled', label: '사용자 취소' },
  { value: 'all', label: '전체' },
] as const;

const targetTypeOptions = [
  { value: 'all', label: '대상 전체' },
  { value: 'source', label: '원문 자료' },
  { value: 'ai_summary', label: 'AI 요약' },
  { value: 'attachment', label: '첨부/보조자료' },
  { value: 'other', label: '기타' },
] as const;

const reportTypeOptions = [
  { value: 'all', label: '유형 전체' },
  { value: 'broken_link', label: '링크 깨짐' },
  { value: 'wrong_content', label: '내용 불일치' },
  { value: 'unsafe_content', label: '부적절 자료' },
  { value: 'copyright', label: '저작권' },
  { value: 'low_quality', label: '품질 낮음' },
  { value: 'other', label: '기타' },
] as const;

const statusLabel: Record<ReportStatus, string> = {
  open: '열림',
  reviewing: '검토중',
  resolved: '해결',
  dismissed: '반려',
  cancelled: '사용자 취소',
};

const statusColor: Record<ReportStatus, { color: string; background: string; border: string }> = {
  open: { color: '#FBBF24', background: 'rgba(251, 191, 36, 0.12)', border: 'rgba(251, 191, 36, 0.32)' },
  reviewing: { color: '#93C5FD', background: 'rgba(96, 165, 250, 0.13)', border: 'rgba(96, 165, 250, 0.32)' },
  resolved: { color: '#86EFAC', background: 'rgba(34, 197, 94, 0.13)', border: 'rgba(34, 197, 94, 0.32)' },
  dismissed: { color: '#FDA4AF', background: 'rgba(244, 63, 94, 0.13)', border: 'rgba(244, 63, 94, 0.32)' },
  cancelled: { color: '#CBD5E1', background: 'rgba(148, 163, 184, 0.13)', border: 'rgba(148, 163, 184, 0.32)' },
};

export default function SuperAdminMaterialReportsPage() {
  const router = useRouter();
  const [reports, setReports] = useState<MaterialReport[]>([]);
  const [selectedID, setSelectedID] = useState('');
  const [statusFilter, setStatusFilter] = useState('open');
  const [targetFilter, setTargetFilter] = useState('all');
  const [reportTypeFilter, setReportTypeFilter] = useState('all');
  const [query, setQuery] = useState('');
  const [appliedQuery, setAppliedQuery] = useState('');
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [adminNote, setAdminNote] = useState('');

  const selectedReport = useMemo(
    () => reports.find((report) => report.id === selectedID) ?? reports[0] ?? null,
    [reports, selectedID],
  );
  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE));

  const fetchReports = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const params = new URLSearchParams({
        page: String(page),
        limit: String(PAGE_SIZE),
        status: statusFilter,
        target_type: targetFilter,
        report_type: reportTypeFilter,
      });
      if (appliedQuery.trim()) {
        params.set('q', appliedQuery.trim());
      }
      const res = await fetch(`${API_BASE}/api/v1/super-admin/material-reports?${params.toString()}`, {
        credentials: 'include',
      });
      if (res.status === 401 || res.status === 403) {
        router.replace('/super-admin/login');
        return;
      }
      if (!res.ok) {
        throw new Error('오류신고 내역을 불러오지 못했습니다.');
      }
      const data = await res.json();
      const nextReports = (data.reports ?? []) as MaterialReport[];
      setReports(nextReports);
      setTotal(data.total ?? 0);
      setSelectedID((current) => {
        if (current && nextReports.some((report) => report.id === current)) return current;
        return nextReports[0]?.id ?? '';
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : '오류가 발생했습니다.');
    } finally {
      setLoading(false);
    }
  }, [appliedQuery, page, reportTypeFilter, router, statusFilter, targetFilter]);

  useEffect(() => {
    fetchReports();
  }, [fetchReports]);

  useEffect(() => {
    setPage(1);
  }, [statusFilter, targetFilter, reportTypeFilter, appliedQuery]);

  useEffect(() => {
    setAdminNote(selectedReport?.admin_note ?? '');
  }, [selectedReport?.id, selectedReport?.admin_note]);

  const updateReport = async (status: ReportStatus) => {
    if (!selectedReport) return;
    setSaving(true);
    setError('');
    try {
      const res = await fetch(`${API_BASE}/api/v1/super-admin/material-reports/${selectedReport.id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ status, admin_note: adminNote }),
      });
      if (res.status === 401 || res.status === 403) {
        router.replace('/super-admin/login');
        return;
      }
      if (!res.ok) {
        throw new Error('처리 상태를 저장하지 못했습니다.');
      }
      await fetchReports();
    } catch (err) {
      setError(err instanceof Error ? err.message : '오류가 발생했습니다.');
    } finally {
      setSaving(false);
    }
  };

  const formatDate = (value?: string) => {
    if (!value) return '-';
    return new Date(value).toLocaleString('ko-KR', {
      timeZone: 'Asia/Seoul',
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <main style={pageStyle}>
      <div style={superAdminShellStyle}>
        <div style={headerStackStyle}>
          <SuperAdminPanelHeader
            subtitle="오류 신고"
            description="학습탐험 지점 자료 오류신고를 게시판형 목록으로 확인하고 검토 상태와 관리자 메모를 남깁니다."
          />
          <SuperAdminPanelNav activeSection="material-reports" />
        </div>

        {error ? <div style={errorStyle}>{error}</div> : null}

        <section style={filterBarStyle}>
          <select value={statusFilter} onChange={(event) => setStatusFilter(event.target.value)} style={selectStyle}>
            {statusOptions.map((option) => (
              <option key={option.value} value={option.value}>{option.label}</option>
            ))}
          </select>
          <select value={targetFilter} onChange={(event) => setTargetFilter(event.target.value)} style={selectStyle}>
            {targetTypeOptions.map((option) => (
              <option key={option.value} value={option.value}>{option.label}</option>
            ))}
          </select>
          <select value={reportTypeFilter} onChange={(event) => setReportTypeFilter(event.target.value)} style={selectStyle}>
            {reportTypeOptions.map((option) => (
              <option key={option.value} value={option.value}>{option.label}</option>
            ))}
          </select>
          <form
            style={searchFormStyle}
            onSubmit={(event) => {
              event.preventDefault();
              setAppliedQuery(query);
            }}
          >
            <input
              value={query}
              onChange={(event) => setQuery(event.target.value)}
              placeholder="코스, 지점, 신고자, 메시지 검색"
              style={searchInputStyle}
            />
            <button type="submit" style={primaryButtonStyle}>검색</button>
          </form>
        </section>

        <section style={contentGridStyle}>
          <div style={tablePanelStyle}>
            <div style={tableHeaderStyle}>
              <strong>신고 내역</strong>
              <span style={summaryStyle}>총 {total.toLocaleString('ko-KR')}건 · 페이지당 5개</span>
            </div>
            <div style={tableScrollStyle}>
              <table style={tableStyle}>
                <thead>
                  <tr>
                    {['상태', '신고유형', '대상', '메시지', '신고자', '접수일'].map((heading) => (
                      <th key={heading} style={thStyle}>{heading}</th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {loading ? (
                    <tr><td colSpan={6} style={emptyCellStyle}>불러오는 중...</td></tr>
                  ) : reports.length === 0 ? (
                    <tr><td colSpan={6} style={emptyCellStyle}>조건에 맞는 신고 내역이 없습니다.</td></tr>
                  ) : reports.map((report) => (
                    <tr
                      key={report.id}
                      className="material-report-row"
                      data-selected={selectedReport?.id === report.id ? 'true' : 'false'}
                      onClick={() => setSelectedID(report.id)}
                    >
                      <td style={tdStyle}>
                        <span style={badgeStyle(report.status)}>{statusLabel[report.status]}</span>
                      </td>
                      <td style={tdStyle}>{getReportTypeLabel(report.report_type)}</td>
                      <td style={tdStyle}>
                        <div style={targetStackStyle}>
                          <strong style={targetTitleStyle}>{report.point_title}</strong>
                          <span style={mutedTextStyle}>{getTargetTypeLabel(report.target_type)} · {getPointTypeLabel(report.point_type)}</span>
                        </div>
                      </td>
                      <td style={messageCellStyle}>{report.message}</td>
                      <td style={tdStyle}>{report.reporter_display}</td>
                      <td style={tdStyle}>{formatDate(report.created_at)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div style={paginationStyle}>
              <button
                type="button"
                onClick={() => setPage((current) => Math.max(1, current - 1))}
                disabled={page === 1}
                style={pageButtonStyle(page === 1)}
              >
                이전
              </button>
              <span style={pageTextStyle}>{page} / {totalPages}</span>
              <button
                type="button"
                onClick={() => setPage((current) => Math.min(totalPages, current + 1))}
                disabled={page === totalPages}
                style={pageButtonStyle(page === totalPages)}
              >
                다음
              </button>
            </div>
          </div>

          <aside style={detailPanelStyle}>
            {selectedReport ? (
              <>
                <div style={detailHeaderStyle}>
                  <div style={badgeRowStyle}>
                    <span style={badgeStyle(selectedReport.status)}>{statusLabel[selectedReport.status]}</span>
                    {selectedReport.recommendation_blocked ? (
                      <span style={recommendationBlockedBadgeStyle}>추천 제외됨</span>
                    ) : null}
                  </div>
                  <strong style={detailTitleStyle}>{selectedReport.point_title}</strong>
                  <span style={mutedTextStyle}>{selectedReport.course_title}</span>
                </div>

                <dl style={metaGridStyle}>
                  <Meta label="지역" value={selectedReport.level_title || '-'} />
                  <Meta label="서브지역" value={selectedReport.lesson_title || '-'} />
                  <Meta label="신고자" value={`${selectedReport.reporter_display}${selectedReport.reporter_email ? ` · ${selectedReport.reporter_email}` : ''}`} />
                  <Meta label="접수일" value={formatDate(selectedReport.created_at)} />
                  <Meta label="처리자" value={selectedReport.reviewer_display || '-'} />
                  <Meta label="처리일" value={formatDate(selectedReport.reviewed_at)} />
                </dl>

                <div style={detailBlockStyle}>
                  <span style={detailLabelStyle}>신고 메시지</span>
                  <p style={messageBoxStyle}>{selectedReport.message}</p>
                </div>

                <div style={detailBlockStyle}>
                  <span style={detailLabelStyle}>대상 자료</span>
                  <div style={resourceBoxStyle}>
                    <strong>{selectedReport.attachment_title || selectedReport.point_title}</strong>
                    <span style={mutedTextStyle}>{getTargetTypeLabel(selectedReport.target_type)} · {getReportTypeLabel(selectedReport.report_type)}</span>
                    {selectedReport.report_type === 'broken_link' ? (
                      <span style={selectedReport.recommendation_blocked ? blockedNoticeStyle : mutedTextStyle}>
                        {selectedReport.recommendation_blocked
                          ? '해결 처리된 링크 깨짐 신고라 추천 후보에서 제외됩니다.'
                          : '해결 처리하면 이 대상은 추천 후보에서 제외됩니다.'}
                      </span>
                    ) : null}
                    {getTargetURL(selectedReport) ? (
                      <a href={getTargetURL(selectedReport)} target="_blank" rel="noopener noreferrer" style={linkStyle}>
                        새 창에서 대상 열기
                      </a>
                    ) : (
                      <span style={mutedTextStyle}>바로 열 수 있는 URL이 없습니다.</span>
                    )}
                  </div>
                </div>

                <label style={noteLabelStyle}>
                  관리자 메모
                  <textarea
                    value={adminNote}
                    onChange={(event) => setAdminNote(event.target.value)}
                    rows={6}
                    maxLength={2000}
                    style={textareaStyle}
                    placeholder="처리 내용, 확인 결과, 후속 조치 등을 기록합니다."
                  />
                </label>

                <div style={actionGridStyle}>
                  <button type="button" disabled={saving} onClick={() => updateReport('reviewing')} style={actionButtonStyle('reviewing')}>검토중</button>
                  <button type="button" disabled={saving} onClick={() => updateReport('resolved')} style={actionButtonStyle('resolved')}>해결</button>
                  <button type="button" disabled={saving} onClick={() => updateReport('dismissed')} style={actionButtonStyle('dismissed')}>반려</button>
                  <button type="button" disabled={saving} onClick={() => updateReport('open')} style={actionButtonStyle('open')}>열림으로 되돌리기</button>
                </div>
              </>
            ) : (
              <p style={emptyDetailStyle}>선택된 신고 내역이 없습니다.</p>
            )}
          </aside>
        </section>
      </div>

      <style jsx>{`
        .material-report-row {
          border-bottom: 1px solid rgba(78, 104, 148, 0.42);
          cursor: pointer;
          transition: background 140ms ease, box-shadow 140ms ease;
        }
        .material-report-row:hover {
          background: rgba(96, 165, 250, 0.12);
          box-shadow: inset 3px 0 0 rgba(96, 165, 250, 0.75);
        }
        .material-report-row[data-selected='true'] {
          background: rgba(34, 197, 94, 0.1);
          box-shadow: inset 3px 0 0 rgba(34, 197, 94, 0.72);
        }
      `}</style>
    </main>
  );
}

function Meta({ label, value }: { label: string; value: string }) {
  return (
    <div style={metaItemStyle}>
      <dt style={metaLabelStyle}>{label}</dt>
      <dd style={metaValueStyle}>{value}</dd>
    </div>
  );
}

function getTargetURL(report: MaterialReport) {
  return report.attachment_url || report.point_external_url || '';
}

function getPointTypeLabel(value: string) {
  if (value === 'exploration') return '탐험지점';
  if (value === 'research') return '연구지점';
  return value || '-';
}

function getTargetTypeLabel(value: string) {
  return targetTypeOptions.find((option) => option.value === value)?.label ?? value;
}

function getReportTypeLabel(value: string) {
  return reportTypeOptions.find((option) => option.value === value)?.label ?? value;
}

const badgeStyle = (status: ReportStatus): CSSProperties => ({
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '24px',
  padding: '0 9px',
  borderRadius: '999px',
  fontSize: '12px',
  fontWeight: 800,
  color: statusColor[status].color,
  background: statusColor[status].background,
  border: `1px solid ${statusColor[status].border}`,
});

const pageStyle: CSSProperties = {
  minHeight: '100vh',
  background: 'radial-gradient(circle at top, rgba(45, 78, 132, 0.2), transparent 42%), linear-gradient(180deg, #07111f, #0b1321 42%, #060b12)',
  color: '#E2E8F0',
  padding: '28px 20px 48px',
};

const headerStackStyle: CSSProperties = {
  display: 'grid',
  gap: '16px',
  marginBottom: '18px',
};

const errorStyle: CSSProperties = {
  marginBottom: '14px',
  padding: '12px 16px',
  borderRadius: '8px',
  background: 'rgba(127, 29, 29, 0.54)',
  border: '1px solid rgba(248, 113, 113, 0.42)',
  color: '#FCA5A5',
  fontSize: '14px',
};

const filterBarStyle: CSSProperties = {
  display: 'flex',
  gap: '10px',
  flexWrap: 'wrap',
  alignItems: 'center',
  marginBottom: '14px',
  padding: '12px',
  border: '1px solid rgba(78, 104, 148, 0.42)',
  borderRadius: '8px',
  background: 'rgba(15, 35, 55, 0.72)',
};

const selectStyle: CSSProperties = {
  height: '38px',
  borderRadius: '7px',
  border: '1px solid rgba(100, 116, 139, 0.6)',
  background: '#0B1629',
  color: '#E2E8F0',
  padding: '0 10px',
  fontSize: '13px',
};

const searchFormStyle: CSSProperties = {
  display: 'flex',
  gap: '8px',
  flex: '1 1 260px',
};

const searchInputStyle: CSSProperties = {
  minWidth: 0,
  flex: 1,
  height: '38px',
  borderRadius: '7px',
  border: '1px solid rgba(100, 116, 139, 0.6)',
  background: '#07111F',
  color: '#E2E8F0',
  padding: '0 12px',
  fontSize: '13px',
};

const primaryButtonStyle: CSSProperties = {
  minHeight: '38px',
  borderRadius: '7px',
  border: '1px solid rgba(96, 165, 250, 0.45)',
  background: 'rgba(37, 99, 235, 0.3)',
  color: '#DBEAFE',
  padding: '0 14px',
  fontWeight: 800,
  cursor: 'pointer',
};

const contentGridStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))',
  gap: '14px',
  alignItems: 'start',
};

const tablePanelStyle: CSSProperties = {
  border: '1px solid rgba(78, 104, 148, 0.48)',
  borderRadius: '8px',
  background: 'rgba(15, 35, 55, 0.72)',
  overflow: 'hidden',
};

const tableHeaderStyle: CSSProperties = {
  display: 'flex',
  justifyContent: 'space-between',
  gap: '12px',
  alignItems: 'center',
  padding: '14px 16px',
  borderBottom: '1px solid rgba(78, 104, 148, 0.48)',
};

const summaryStyle: CSSProperties = {
  fontSize: '12px',
  color: '#94A3B8',
};

const tableScrollStyle: CSSProperties = {
  overflowX: 'auto',
};

const tableStyle: CSSProperties = {
  width: '100%',
  minWidth: '760px',
  borderCollapse: 'collapse',
  fontSize: '13px',
};

const thStyle: CSSProperties = {
  padding: '11px 12px',
  textAlign: 'left',
  color: '#94A3B8',
  background: '#0B1629',
  fontWeight: 800,
  borderBottom: '1px solid rgba(78, 104, 148, 0.48)',
};

const tdStyle: CSSProperties = {
  padding: '12px',
  verticalAlign: 'middle',
};

const messageCellStyle: CSSProperties = {
  ...tdStyle,
  maxWidth: '260px',
  whiteSpace: 'nowrap',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  color: '#CBD5E1',
};

const emptyCellStyle: CSSProperties = {
  padding: '44px 16px',
  textAlign: 'center',
  color: '#94A3B8',
};

const targetStackStyle: CSSProperties = {
  display: 'grid',
  gap: '4px',
};

const targetTitleStyle: CSSProperties = {
  color: '#F8FAFC',
  fontSize: '13px',
};

const mutedTextStyle: CSSProperties = {
  color: '#94A3B8',
  fontSize: '12px',
};

const paginationStyle: CSSProperties = {
  display: 'flex',
  justifyContent: 'center',
  gap: '8px',
  padding: '14px',
  borderTop: '1px solid rgba(78, 104, 148, 0.48)',
};

const pageButtonStyle = (disabled: boolean): CSSProperties => ({
  minHeight: '34px',
  padding: '0 13px',
  borderRadius: '7px',
  border: '1px solid rgba(78, 104, 148, 0.7)',
  background: disabled ? '#0B1629' : 'rgba(30, 58, 95, 0.86)',
  color: disabled ? '#475569' : '#E2E8F0',
  cursor: disabled ? 'not-allowed' : 'pointer',
  fontSize: '13px',
});

const pageTextStyle: CSSProperties = {
  minHeight: '34px',
  display: 'inline-flex',
  alignItems: 'center',
  color: '#94A3B8',
  fontSize: '13px',
};

const detailPanelStyle: CSSProperties = {
  display: 'grid',
  gap: '16px',
  padding: '16px',
  border: '1px solid rgba(78, 104, 148, 0.48)',
  borderRadius: '8px',
  background: 'rgba(9, 23, 39, 0.86)',
};

const detailHeaderStyle: CSSProperties = {
  display: 'grid',
  gap: '8px',
};

const badgeRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: '8px',
  flexWrap: 'wrap',
};

const recommendationBlockedBadgeStyle: CSSProperties = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '24px',
  padding: '0 9px',
  borderRadius: '999px',
  fontSize: '12px',
  fontWeight: 800,
  color: '#FDE68A',
  background: 'rgba(217, 119, 6, 0.16)',
  border: '1px solid rgba(217, 119, 6, 0.38)',
};

const detailTitleStyle: CSSProperties = {
  color: '#F8FAFC',
  fontSize: '18px',
  lineHeight: 1.35,
};

const metaGridStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '1fr 1fr',
  gap: '10px',
  margin: 0,
};

const metaItemStyle: CSSProperties = {
  display: 'grid',
  gap: '4px',
  padding: '10px',
  border: '1px solid rgba(78, 104, 148, 0.36)',
  borderRadius: '7px',
  background: 'rgba(15, 35, 55, 0.5)',
};

const metaLabelStyle: CSSProperties = {
  margin: 0,
  color: '#94A3B8',
  fontSize: '12px',
  fontWeight: 700,
};

const metaValueStyle: CSSProperties = {
  margin: 0,
  color: '#E2E8F0',
  fontSize: '13px',
  lineHeight: 1.45,
  wordBreak: 'break-word',
};

const detailBlockStyle: CSSProperties = {
  display: 'grid',
  gap: '8px',
};

const detailLabelStyle: CSSProperties = {
  color: '#CBD5E1',
  fontSize: '13px',
  fontWeight: 800,
};

const messageBoxStyle: CSSProperties = {
  margin: 0,
  padding: '12px',
  minHeight: '84px',
  borderRadius: '7px',
  border: '1px solid rgba(78, 104, 148, 0.36)',
  background: 'rgba(15, 35, 55, 0.5)',
  color: '#E2E8F0',
  lineHeight: 1.7,
  whiteSpace: 'pre-wrap',
};

const resourceBoxStyle: CSSProperties = {
  display: 'grid',
  gap: '6px',
  padding: '12px',
  borderRadius: '7px',
  border: '1px solid rgba(78, 104, 148, 0.36)',
  background: 'rgba(15, 35, 55, 0.5)',
};

const blockedNoticeStyle: CSSProperties = {
  color: '#FDE68A',
  fontSize: '12px',
  lineHeight: 1.5,
};

const linkStyle: CSSProperties = {
  color: '#93C5FD',
  fontSize: '13px',
  fontWeight: 800,
};

const noteLabelStyle: CSSProperties = {
  display: 'grid',
  gap: '8px',
  color: '#CBD5E1',
  fontSize: '13px',
  fontWeight: 800,
};

const textareaStyle: CSSProperties = {
  width: '100%',
  resize: 'vertical',
  borderRadius: '7px',
  border: '1px solid rgba(100, 116, 139, 0.6)',
  background: '#07111F',
  color: '#E2E8F0',
  padding: '10px 12px',
  fontSize: '13px',
  lineHeight: 1.6,
  fontFamily: 'inherit',
};

const actionGridStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '1fr 1fr',
  gap: '8px',
};

const actionButtonStyle = (status: ReportStatus): CSSProperties => ({
  minHeight: '38px',
  borderRadius: '7px',
  border: `1px solid ${statusColor[status].border}`,
  background: statusColor[status].background,
  color: statusColor[status].color,
  fontWeight: 900,
  cursor: 'pointer',
});

const emptyDetailStyle: CSSProperties = {
  margin: 0,
  padding: '36px 12px',
  textAlign: 'center',
  color: '#94A3B8',
};
