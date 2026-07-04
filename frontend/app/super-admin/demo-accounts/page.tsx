'use client';

import { useEffect, useRef, useState } from 'react';
import type * as React from 'react';
import { useRouter } from 'next/navigation';
import SuperAdminPanelHeader from '@/components/super-admin/SuperAdminPanelHeader';
import SuperAdminPanelNav from '@/components/super-admin/SuperAdminPanelNav';
import { superAdminShellStyle } from '@/components/super-admin/layout';

interface DemoAccount {
  id: string;
  user_id: string;
  label: string;
  email: string;
  nickname: string;
  assigned_to: string;
  assignment_note: string;
  login_code?: string;
  status: 'active' | 'disabled' | 'expired' | 'unavailable';
  expires_at?: string;
  activated_at?: string;
  last_used_at?: string;
  disabled_at?: string;
  free_points: number;
  paid_points: number;
  created_at: string;
  updated_at: string;
}

interface CreatedDemoAccount {
  id: string;
  user_id: string;
  label: string;
  email: string;
  nickname: string;
  login_code: string;
}

interface PurgePreview {
  demo_account_id: string;
  user_id: string;
  label: string;
  course_drafts: number;
  courses: number;
  goal_profiles: number;
  goal_revision_logs: number;
  course_points: number;
  course_draft_points: number;
  attachments: number;
  object_storage_files: number;
  ai_usage_events: number;
  recommendation_events: number;
  recommendation_learner_events: number;
  contents: number;
  safe_deletable_contents: number;
  object_keys?: string[];
}

interface PurgeResult extends PurgePreview {
  deleted_object_keys: string[];
  object_delete_failed_keys: string[];
  object_delete_failure_text: string[];
  account_deleted: boolean;
}

const statusLabels: Record<string, string> = {
  active: '활성',
  disabled: '비활성',
  expired: '만료',
  unavailable: '사용자 비활성',
};

export default function DemoAccountsPage() {
  const router = useRouter();
  const expiresAtInputRef = useRef<HTMLInputElement | null>(null);
  const editExpiresAtInputRef = useRef<HTMLInputElement | null>(null);
  const [items, setItems] = useState<DemoAccount[]>([]);
  const [createdItems, setCreatedItems] = useState<CreatedDemoAccount[]>([]);
  const [query, setQuery] = useState('');
  const [status, setStatus] = useState('all');
  const [count, setCount] = useState(1);
  const [startNumber, setStartNumber] = useState(1);
  const [expiresAt, setExpiresAt] = useState(defaultExpiresAtDate());
  const [assignedTo, setAssignedTo] = useState('');
  const [assignmentNote, setAssignmentNote] = useState('');
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [purgePreview, setPurgePreview] = useState<PurgePreview | null>(null);
  const [purgeConfirm, setPurgeConfirm] = useState('');
  const [purgeLoading, setPurgeLoading] = useState(false);
  const [editingExpiresId, setEditingExpiresId] = useState('');
  const [editingExpiresValue, setEditingExpiresValue] = useState('');
  const [updatingExpiresId, setUpdatingExpiresId] = useState('');
  const loadItems = async () => {
    setLoading(true);
    setError('');
    try {
      const params = new URLSearchParams({ q: query, status, page: '1', limit: '100' });
      const res = await fetch(`/api/v1/super-admin/demo-accounts?${params.toString()}`, {
        credentials: 'include',
        cache: 'no-store',
      });
      if (res.status === 401 || res.status === 403) {
        router.replace('/super-admin/login');
        return;
      }
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(typeof payload.error === 'string' ? payload.error : '데모 계정 목록을 불러오지 못했습니다.');
      setItems(Array.isArray(payload.items) ? payload.items : []);
    } catch (err) {
      setError(err instanceof Error ? err.message : '데모 계정 목록을 불러오지 못했습니다.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadItems();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [status]);

  const handleCreate = async () => {
    setSubmitting(true);
    setError('');
    setMessage('');
    setCreatedItems([]);
    try {
      const res = await fetch('/api/v1/super-admin/demo-accounts', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          count,
          start_number: startNumber,
          email_prefix: 'demo',
          email_domain: 'learnweavr.local',
          label_prefix: 'demo',
          ui_locale: 'ko',
          learning_language: 'ko',
          initial_points: 300,
          expires_at: expiresAt,
          assigned_to: assignedTo,
          assignment_note: assignmentNote,
        }),
      });
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(typeof payload.error === 'string' ? payload.error : '데모 계정 생성에 실패했습니다.');
      const created = Array.isArray(payload.items) ? payload.items : [];
      setCreatedItems(created);
      setMessage(`${created.length}개 데모 계정을 생성했습니다. 로그인 코드는 목록에서도 확인할 수 있습니다.`);
      await loadItems();
    } catch (err) {
      setError(err instanceof Error ? err.message : '데모 계정 생성에 실패했습니다.');
    } finally {
      setSubmitting(false);
    }
  };

  const handleRotate = async (item: DemoAccount) => {
    setError('');
    setMessage('');
    setCreatedItems([]);
    try {
      const res = await fetch(`/api/v1/super-admin/demo-accounts/${item.id}/rotate-code`, {
        method: 'POST',
        credentials: 'include',
      });
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(typeof payload.error === 'string' ? payload.error : '로그인 코드 재발급에 실패했습니다.');
      setCreatedItems([{
        id: item.id,
        user_id: item.user_id,
        label: item.label,
        email: item.email,
        nickname: item.nickname,
        login_code: payload.login_code,
      }]);
      setMessage(`${item.label} 로그인 코드를 재발급했습니다. 새 코드는 목록에서도 확인할 수 있습니다.`);
      await loadItems();
    } catch (err) {
      setError(err instanceof Error ? err.message : '로그인 코드 재발급에 실패했습니다.');
    }
  };

  const handleDisabled = async (item: DemoAccount, disabled: boolean) => {
    setError('');
    setMessage('');
    try {
      const res = await fetch(`/api/v1/super-admin/demo-accounts/${item.id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          disabled,
          disabled_reason: disabled ? 'super admin disabled' : '',
        }),
      });
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(typeof payload.error === 'string' ? payload.error : '데모 계정 상태 변경에 실패했습니다.');
      setMessage(`${item.label} 상태를 변경했습니다.`);
      await loadItems();
    } catch (err) {
      setError(err instanceof Error ? err.message : '데모 계정 상태 변경에 실패했습니다.');
    }
  };

  const startEditExpires = (item: DemoAccount) => {
    setEditingExpiresId(item.id);
    setEditingExpiresValue(toDateInputValue(item.expires_at) || defaultExpiresAtDate());
    setError('');
    setMessage('');
  };

  const cancelEditExpires = () => {
    setEditingExpiresId('');
    setEditingExpiresValue('');
  };

  const handleUpdateExpires = async (item: DemoAccount) => {
    if (!editingExpiresValue) return;
    setUpdatingExpiresId(item.id);
    setError('');
    setMessage('');
    try {
      const res = await fetch(`/api/v1/super-admin/demo-accounts/${item.id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          expires_at: editingExpiresValue,
        }),
      });
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(typeof payload.error === 'string' ? payload.error : '만료일 변경에 실패했습니다.');
      setMessage(`${item.label} 만료일을 변경했습니다.`);
      cancelEditExpires();
      await loadItems();
    } catch (err) {
      setError(err instanceof Error ? err.message : '만료일 변경에 실패했습니다.');
    } finally {
      setUpdatingExpiresId('');
    }
  };

  const openDatePicker = (inputRef: React.RefObject<HTMLInputElement | null>) => {
    const input = inputRef.current as (HTMLInputElement & { showPicker?: () => void }) | null;
    if (!input) return;
    input.focus();
    if (typeof input.showPicker === 'function') {
      try {
        input.showPicker();
      } catch {
        // Some browsers allow the native date picker only on direct pointer activation.
      }
    }
  };

  const handlePreviewPurge = async (item: DemoAccount) => {
    setError('');
    setMessage('');
    setPurgeConfirm('');
    setPurgeLoading(true);
    try {
      const res = await fetch(`/api/v1/super-admin/demo-accounts/${item.id}/purge-preview`, {
        credentials: 'include',
        cache: 'no-store',
      });
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(typeof payload.error === 'string' ? payload.error : '삭제 미리보기에 실패했습니다.');
      setPurgePreview(payload as PurgePreview);
    } catch (err) {
      setError(err instanceof Error ? err.message : '삭제 미리보기에 실패했습니다.');
    } finally {
      setPurgeLoading(false);
    }
  };

  const handlePurge = async () => {
    if (!purgePreview) return;
    setPurgeLoading(true);
    setError('');
    setMessage('');
    try {
      const res = await fetch(`/api/v1/super-admin/demo-accounts/${purgePreview.demo_account_id}/purge`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          confirm: purgeConfirm,
          include_account: false,
        }),
      });
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(typeof payload.error === 'string' ? payload.error : '데모 데이터 삭제에 실패했습니다.');
      const result = payload as PurgeResult;
      const failedCount = result.object_delete_failed_keys?.length ?? 0;
      setMessage(failedCount > 0
        ? `${result.label} 학습 데이터를 삭제했습니다. Object Storage 삭제 실패 ${failedCount}건은 결과를 확인해 주세요.`
        : `${result.label} 학습 데이터를 삭제했습니다.`);
      setPurgePreview(null);
      setPurgeConfirm('');
      await loadItems();
    } catch (err) {
      setError(err instanceof Error ? err.message : '데모 데이터 삭제에 실패했습니다.');
    } finally {
      setPurgeLoading(false);
    }
  };

  return (
    <main style={pageStyle}>
      <div style={superAdminShellStyle}>
        <SuperAdminPanelHeader
          subtitle="데모 계정"
          description="외부 검수자와 알파 테스터에게 전달할 데모 학습자 계정을 생성하고 로그인 코드를 관리합니다."
        />
        <SuperAdminPanelNav activeSection="demo-accounts" />

        <section style={panelStyle}>
          <div>
            <h2 style={sectionTitleStyle}>데모 계정 생성</h2>
            <p style={sectionCopyStyle}>기본값은 한국어 학습자, 300pt, 30일 만료입니다. 생성된 로그인 코드는 목록에서 확인할 수 있습니다.</p>
          </div>
          <div style={formGridStyle}>
            <label style={fieldStyle}>
              <span style={labelStyle}>생성 개수</span>
              <input type="number" min={1} max={100} value={count} onChange={(event) => setCount(clampNumber(event.target.value, 1, 100))} style={inputStyle} />
            </label>
            <label style={fieldStyle}>
              <span style={labelStyle}>시작 번호</span>
              <input type="number" min={1} value={startNumber} onChange={(event) => setStartNumber(clampNumber(event.target.value, 1, 9999))} style={inputStyle} />
            </label>
            <label style={fieldStyle}>
              <span style={labelStyle}>만료일</span>
              <span style={dateInputShellStyle}>
                <input
                  ref={expiresAtInputRef}
                  type="date"
                  value={expiresAt}
                  onClick={() => openDatePicker(expiresAtInputRef)}
                  onFocus={() => openDatePicker(expiresAtInputRef)}
                  onChange={(event) => setExpiresAt(event.target.value)}
                  style={dateInputStyle}
                />
                <button type="button" onClick={() => openDatePicker(expiresAtInputRef)} style={datePickerButtonStyle}>
                  달력 열기
                </button>
              </span>
            </label>
            <label style={fieldStyle}>
              <span style={labelStyle}>배정 대상</span>
              <input value={assignedTo} onChange={(event) => setAssignedTo(event.target.value)} placeholder="예: YouTube API reviewer" style={inputStyle} />
            </label>
            <label style={{ ...fieldStyle, gridColumn: '1 / -1' }}>
              <span style={labelStyle}>배정 메모</span>
              <textarea value={assignmentNote} onChange={(event) => setAssignmentNote(event.target.value)} rows={3} placeholder="전달 경로, 목적, 회수 일정" style={textareaStyle} />
            </label>
          </div>
          <button type="button" onClick={handleCreate} disabled={submitting} style={primaryButtonStyle}>
            {submitting ? '생성 중...' : '데모 계정 생성'}
          </button>
          {message ? <p style={successStyle}>{message}</p> : null}
          {error ? <p style={errorStyle}>{error}</p> : null}
        </section>

        {createdItems.length > 0 ? (
          <section style={panelStyle}>
            <h2 style={sectionTitleStyle}>새 로그인 코드</h2>
            <p style={warningStyle}>생성 직후 바로 전달할 수 있도록 새 로그인 코드를 모아 보여줍니다. 같은 코드는 아래 목록에서도 확인할 수 있습니다.</p>
            <div style={codeListStyle}>
              {createdItems.map((item) => (
                <div key={`${item.id}-${item.login_code}`} style={codeCardStyle}>
                  <div>
                    <strong style={{ color: '#F7FAFF' }}>{item.label}</strong>
                    <p style={mutedStyle}>{item.email}</p>
                  </div>
                  <code style={codeStyle}>{item.login_code}</code>
                  <button type="button" onClick={() => navigator.clipboard?.writeText(item.login_code)} style={secondaryButtonStyle}>복사</button>
                </div>
              ))}
            </div>
          </section>
        ) : null}

        <section style={panelStyle}>
          <div style={listHeaderStyle}>
            <div>
              <h2 style={sectionTitleStyle}>데모 계정 목록</h2>
              <p style={sectionCopyStyle}>최근 생성 계정 100개까지 표시합니다.</p>
            </div>
            <div style={filterRowStyle}>
              <input value={query} onChange={(event) => setQuery(event.target.value)} onKeyDown={(event) => { if (event.key === 'Enter') loadItems(); }} placeholder="label, email, 배정 대상 검색" style={searchInputStyle} />
              <select value={status} onChange={(event) => setStatus(event.target.value)} style={selectStyle}>
                <option value="all">전체</option>
                <option value="active">활성</option>
                <option value="disabled">비활성</option>
                <option value="expired">만료</option>
              </select>
              <button type="button" onClick={loadItems} style={secondaryButtonStyle}>조회</button>
            </div>
          </div>

          {loading ? (
            <p style={emptyStyle}>목록을 불러오는 중입니다.</p>
          ) : items.length === 0 ? (
            <p style={emptyStyle}>데모 계정이 없습니다.</p>
          ) : (
            <div style={tableWrapStyle}>
              <table style={tableStyle}>
                <thead>
                  <tr>
                    <th style={thStyle}>계정</th>
                    <th style={thStyle}>로그인 코드</th>
                    <th style={thStyle}>상태</th>
                    <th style={thStyle}>배정</th>
                    <th style={thStyle}>포인트</th>
                    <th style={thStyle}>활성화/사용</th>
                    <th style={thStyle}>작업</th>
                  </tr>
                </thead>
                <tbody>
                  {items.map((item) => (
                    <tr key={item.id}>
                      <td style={tdStyle}>
                        <strong style={{ color: '#F7FAFF' }}>{item.label}</strong>
                        <p style={mutedStyle}>{item.email}</p>
                        <p style={mutedStyle}>{item.nickname}</p>
                      </td>
                      <td style={tdStyle}>
                        {item.login_code ? (
                          <div style={inlineCodeStyle}>
                            <code style={codeStyle}>{item.login_code}</code>
                            <button type="button" onClick={() => navigator.clipboard?.writeText(item.login_code || '')} style={tinyButtonStyle}>복사</button>
                          </div>
                        ) : (
                          <div style={missingCodeStyle}>
                            <span>재발급 필요</span>
                            <button type="button" onClick={() => handleRotate(item)} style={tinyButtonStyle}>코드 재발급</button>
                          </div>
                        )}
                      </td>
                      <td style={tdStyle}>
                        <span style={statusBadgeStyle(item.status)}>{statusLabels[item.status] ?? item.status}</span>
                      </td>
                      <td style={tdStyle}>
                        <p style={cellMainStyle}>{item.assigned_to || '-'}</p>
                        <p style={mutedStyle}>{item.assignment_note || ''}</p>
                      </td>
                      <td style={tdStyle}>{item.free_points + item.paid_points}pt</td>
                      <td style={tdStyle}>
                        <p style={cellMainStyle}>활성화 {formatDateTime(item.activated_at)}</p>
                        <p style={mutedStyle}>최근 사용 {formatDateTime(item.last_used_at)}</p>
                        {editingExpiresId === item.id ? (
                          <div style={expiresEditStyle}>
                            <label style={compactFieldStyle}>
                              <span style={labelStyle}>만료일</span>
                              <span style={dateInputShellStyle}>
                                <input
                                  ref={editExpiresAtInputRef}
                                  type="date"
                                  value={editingExpiresValue}
                                  onClick={() => openDatePicker(editExpiresAtInputRef)}
                                  onFocus={() => openDatePicker(editExpiresAtInputRef)}
                                  onChange={(event) => setEditingExpiresValue(event.target.value)}
                                  style={dateInputStyle}
                                />
                                <button type="button" onClick={() => openDatePicker(editExpiresAtInputRef)} style={datePickerButtonStyle}>
                                  달력 열기
                                </button>
                              </span>
                            </label>
                            <div style={actionRowStyle}>
                              <button type="button" onClick={() => handleUpdateExpires(item)} disabled={updatingExpiresId === item.id} style={smallButtonStyle}>
                                {updatingExpiresId === item.id ? '저장 중...' : '저장'}
                              </button>
                              <button type="button" onClick={cancelEditExpires} style={smallButtonStyle}>취소</button>
                            </div>
                          </div>
                        ) : (
                          <div style={expiresDisplayStyle}>
                            <p style={cellMainStyle}>만료 {formatDate(item.expires_at)}</p>
                            <button type="button" onClick={() => startEditExpires(item)} style={tinyButtonStyle}>만료일 연장</button>
                          </div>
                        )}
                      </td>
                      <td style={tdStyle}>
                        <div style={actionRowStyle}>
                          <button type="button" onClick={() => handleRotate(item)} style={smallButtonStyle}>코드 재발급</button>
                          {item.status === 'disabled' ? (
                            <button type="button" onClick={() => handleDisabled(item, false)} style={smallButtonStyle}>활성화</button>
                          ) : (
                            <button type="button" onClick={() => handleDisabled(item, true)} style={dangerButtonStyle}>비활성화</button>
                          )}
                          <button type="button" onClick={() => handlePreviewPurge(item)} style={dangerButtonStyle}>삭제 미리보기</button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>

        {purgePreview ? (
          <div style={modalBackdropStyle} role="dialog" aria-modal="true" aria-labelledby="demo-purge-title">
            <section style={modalPanelStyle}>
              <h2 id="demo-purge-title" style={sectionTitleStyle}>데모 데이터 삭제</h2>
              <p style={warningStyle}>
                계정 row는 유지하고 학습 데이터만 삭제합니다. 일반 사용자 데이터는 삭제하지 않도록 demo 계정 user_id 기준으로 제한합니다.
              </p>
              <div style={previewGridStyle}>
                <PreviewStat label="코스 초안" value={purgePreview.course_drafts} />
                <PreviewStat label="학습 코스" value={purgePreview.courses} />
                <PreviewStat label="목표 프로필" value={purgePreview.goal_profiles} />
                <PreviewStat label="지점" value={purgePreview.course_points + purgePreview.course_draft_points} />
                <PreviewStat label="첨부" value={purgePreview.attachments} />
                <PreviewStat label="Object 파일" value={purgePreview.object_storage_files} />
                <PreviewStat label="AI 사용 기록" value={purgePreview.ai_usage_events} />
                <PreviewStat label="추천 이벤트" value={purgePreview.recommendation_events + purgePreview.recommendation_learner_events} />
                <PreviewStat label="콘텐츠 전체" value={purgePreview.contents} />
                <PreviewStat label="삭제 가능 콘텐츠" value={purgePreview.safe_deletable_contents} />
              </div>
              <label style={fieldStyle}>
                <span style={labelStyle}>확인 문자열</span>
                <input
                  value={purgeConfirm}
                  onChange={(event) => setPurgeConfirm(event.target.value)}
                  placeholder="PURGE_DEMO_ACCOUNT_DATA"
                  style={inputStyle}
                />
              </label>
              <div style={modalActionRowStyle}>
                <button type="button" onClick={() => setPurgePreview(null)} style={secondaryButtonStyle}>취소</button>
                <button
                  type="button"
                  onClick={handlePurge}
                  disabled={purgeLoading || purgeConfirm !== 'PURGE_DEMO_ACCOUNT_DATA'}
                  style={{ ...dangerButtonStyle, minHeight: '38px', opacity: purgeConfirm === 'PURGE_DEMO_ACCOUNT_DATA' ? 1 : 0.55 }}
                >
                  {purgeLoading ? '삭제 중...' : '학습 데이터 삭제'}
                </button>
              </div>
            </section>
          </div>
        ) : null}
      </div>
    </main>
  );
}

function PreviewStat({ label, value }: { label: string; value: number }) {
  return (
    <div style={previewStatStyle}>
      <span style={mutedStyle}>{label}</span>
      <strong style={{ color: '#F7FAFF', fontSize: '18px' }}>{value}</strong>
    </div>
  );
}

function defaultExpiresAtDate() {
  const date = new Date();
  date.setDate(date.getDate() + 30);
  return date.toISOString().slice(0, 10);
}

function clampNumber(value: string, min: number, max: number) {
  const parsed = Number(value);
  if (!Number.isFinite(parsed)) return min;
  return Math.min(max, Math.max(min, Math.floor(parsed)));
}

function formatDate(value?: string) {
  if (!value) return '-';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '-';
  return date.toLocaleDateString('ko-KR');
}

function formatDateTime(value?: string) {
  if (!value) return '-';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '-';
  return date.toLocaleString('ko-KR', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
}

function toDateInputValue(value?: string) {
  if (!value) return '';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '';
  return date.toISOString().slice(0, 10);
}

const pageStyle: React.CSSProperties = {
  minHeight: '100vh',
  padding: '32px 24px 64px',
  background: 'linear-gradient(135deg, #0F172A 0%, #1E293B 52%, #26364F 100%)',
  color: '#E8EAF2',
};

const panelStyle: React.CSSProperties = {
  display: 'grid',
  gap: '18px',
  marginBottom: '18px',
  padding: '22px',
  border: '1px solid rgba(148, 163, 184, 0.22)',
  borderRadius: '8px',
  background: 'rgba(15, 23, 42, 0.72)',
};

const sectionTitleStyle: React.CSSProperties = { margin: 0, fontSize: '18px', fontWeight: 800, color: '#F7FAFF' };
const sectionCopyStyle: React.CSSProperties = { margin: '6px 0 0', fontSize: '13px', lineHeight: 1.7, color: 'rgba(226,232,240,0.72)' };
const formGridStyle: React.CSSProperties = { display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(190px, 1fr))', gap: '14px' };
const fieldStyle: React.CSSProperties = { display: 'grid', gap: '7px' };
const labelStyle: React.CSSProperties = { fontSize: '12px', fontWeight: 700, color: 'rgba(226,232,240,0.74)' };
const inputStyle: React.CSSProperties = { minHeight: '42px', borderRadius: '8px', border: '1px solid rgba(148,163,184,0.32)', background: 'rgba(15,23,42,0.88)', color: '#F7FAFF', padding: '0 12px' };
const dateInputShellStyle: React.CSSProperties = { display: 'flex', alignItems: 'center', minHeight: '42px', borderRadius: '8px', border: '1px solid rgba(148,163,184,0.32)', background: 'rgba(15,23,42,0.88)', overflow: 'hidden' };
const dateInputStyle: React.CSSProperties = { ...inputStyle, flex: '1 1 auto', minWidth: 0, border: 'none', borderRadius: 0, background: 'transparent' };
const datePickerButtonStyle: React.CSSProperties = { alignSelf: 'stretch', flex: '0 0 auto', border: 'none', borderLeft: '1px solid rgba(148,163,184,0.22)', background: 'rgba(110,168,255,0.16)', color: '#E8EAF2', padding: '0 12px', fontSize: '12px', fontWeight: 800, cursor: 'pointer' };
const searchInputStyle: React.CSSProperties = { ...inputStyle, width: '260px' };
const selectStyle: React.CSSProperties = { ...inputStyle, width: '120px' };
const textareaStyle: React.CSSProperties = { ...inputStyle, minHeight: '88px', padding: '12px', resize: 'vertical' };
const primaryButtonStyle: React.CSSProperties = { minHeight: '42px', borderRadius: '8px', border: '0', background: '#6EA8FF', color: '#061525', fontWeight: 800, cursor: 'pointer' };
const secondaryButtonStyle: React.CSSProperties = { minHeight: '36px', borderRadius: '8px', border: '1px solid rgba(148,163,184,0.32)', background: 'rgba(255,255,255,0.06)', color: '#E8EAF2', fontWeight: 700, cursor: 'pointer', padding: '0 12px' };
const smallButtonStyle: React.CSSProperties = { ...secondaryButtonStyle, minHeight: '32px', fontSize: '12px' };
const tinyButtonStyle: React.CSSProperties = { ...secondaryButtonStyle, minHeight: '28px', fontSize: '11px', padding: '0 9px' };
const dangerButtonStyle: React.CSSProperties = { ...smallButtonStyle, color: '#FCA5A5', border: '1px solid rgba(248,113,113,0.38)' };
const successStyle: React.CSSProperties = { margin: 0, color: '#86EFAC', fontSize: '13px' };
const errorStyle: React.CSSProperties = { margin: 0, color: '#FCA5A5', fontSize: '13px' };
const warningStyle: React.CSSProperties = { margin: 0, color: '#FDE68A', fontSize: '13px', lineHeight: 1.7 };
const mutedStyle: React.CSSProperties = { margin: '4px 0 0', fontSize: '12px', color: 'rgba(226,232,240,0.58)' };
const cellMainStyle: React.CSSProperties = { margin: 0, fontSize: '13px', color: 'rgba(248,250,252,0.86)' };
const listHeaderStyle: React.CSSProperties = { display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', gap: '16px', flexWrap: 'wrap' };
const filterRowStyle: React.CSSProperties = { display: 'flex', alignItems: 'center', gap: '8px', flexWrap: 'wrap' };
const emptyStyle: React.CSSProperties = { margin: 0, padding: '22px', borderRadius: '8px', background: 'rgba(255,255,255,0.04)', color: 'rgba(226,232,240,0.62)' };
const tableWrapStyle: React.CSSProperties = { overflowX: 'auto' };
const tableStyle: React.CSSProperties = { width: '100%', borderCollapse: 'collapse', minWidth: '1080px' };
const thStyle: React.CSSProperties = { padding: '10px', textAlign: 'left', fontSize: '12px', color: 'rgba(226,232,240,0.62)', borderBottom: '1px solid rgba(148,163,184,0.2)' };
const tdStyle: React.CSSProperties = { padding: '12px 10px', verticalAlign: 'top', fontSize: '13px', borderBottom: '1px solid rgba(148,163,184,0.12)' };
const actionRowStyle: React.CSSProperties = { display: 'flex', gap: '8px', flexWrap: 'wrap' };
const codeListStyle: React.CSSProperties = { display: 'grid', gap: '10px' };
const codeCardStyle: React.CSSProperties = { display: 'grid', gridTemplateColumns: '1fr auto auto', alignItems: 'center', gap: '12px', padding: '12px', borderRadius: '8px', background: 'rgba(255,255,255,0.05)' };
const codeStyle: React.CSSProperties = { padding: '8px 10px', borderRadius: '6px', background: 'rgba(6,21,37,0.86)', color: '#FDE68A', fontWeight: 800 };
const inlineCodeStyle: React.CSSProperties = { display: 'flex', alignItems: 'center', gap: '8px', flexWrap: 'wrap' };
const missingCodeStyle: React.CSSProperties = { display: 'grid', justifyItems: 'start', gap: '7px', color: 'rgba(226,232,240,0.62)', fontSize: '12px' };
const modalBackdropStyle: React.CSSProperties = { position: 'fixed', inset: 0, zIndex: 120, display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '24px', background: 'rgba(2, 6, 23, 0.72)', backdropFilter: 'blur(8px)' };
const modalPanelStyle: React.CSSProperties = { width: 'min(720px, 100%)', maxHeight: 'min(820px, 92vh)', overflow: 'auto', display: 'grid', gap: '16px', padding: '22px', borderRadius: '8px', border: '1px solid rgba(248,113,113,0.28)', background: '#111827', boxShadow: '0 30px 80px rgba(0,0,0,0.46)' };
const previewGridStyle: React.CSSProperties = { display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(130px, 1fr))', gap: '10px' };
const previewStatStyle: React.CSSProperties = { display: 'grid', gap: '4px', padding: '12px', borderRadius: '8px', background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(148,163,184,0.12)' };
const modalActionRowStyle: React.CSSProperties = { display: 'flex', justifyContent: 'flex-end', gap: '10px', flexWrap: 'wrap' };
const compactFieldStyle: React.CSSProperties = { display: 'grid', gap: '5px' };
const expiresDisplayStyle: React.CSSProperties = { display: 'grid', justifyItems: 'start', gap: '6px', marginTop: '8px' };
const expiresEditStyle: React.CSSProperties = { display: 'grid', gap: '8px', marginTop: '8px' };

function statusBadgeStyle(status: string): React.CSSProperties {
  const color = status === 'active' ? '#86EFAC' : status === 'expired' ? '#FDE68A' : '#FCA5A5';
  return { display: 'inline-flex', padding: '4px 9px', borderRadius: '999px', background: 'rgba(255,255,255,0.06)', color, fontSize: '12px', fontWeight: 800 };
}
