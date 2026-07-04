'use client';

import { useEffect, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import SuperAdminPanelHeader from '@/components/super-admin/SuperAdminPanelHeader';
import SuperAdminPanelNav from '@/components/super-admin/SuperAdminPanelNav';
import { SUPER_ADMIN_PAGE_WIDTH } from '@/components/super-admin/layout';

interface InviteCodeUseRecord {
  id: string;
  user_id: string;
  email: string;
  nickname: string;
  display_id?: string;
  provider: string;
  provider_id: string;
  used_at: string;
}

interface InviteCodeItem {
  id: string;
  code: string;
  status: string;
  state: 'available' | 'used_up' | 'expired' | 'revoked';
  max_uses: number;
  used_count: number;
  expires_at: string;
  sent_to_note: string;
  admin_note: string;
  created_at: string;
  updated_at: string;
  uses: InviteCodeUseRecord[];
}

const stateLabels: Record<InviteCodeItem['state'], string> = {
  available: '사용 가능',
  used_up: '사용 완료',
  expired: '기간 만료',
  revoked: '회수',
};

export default function AlphaInviteCodesPage() {
  const router = useRouter();
  const expiresAtInputRef = useRef<HTMLInputElement | null>(null);
  const [items, setItems] = useState<InviteCodeItem[]>([]);
  const [expiresAt, setExpiresAt] = useState(() => defaultExpiresAtDate());
  const [maxUses, setMaxUses] = useState(1);
  const [sentToNote, setSentToNote] = useState('');
  const [adminNote, setAdminNote] = useState('');
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  const loadItems = async () => {
    setLoading(true);
    setError('');
    try {
      const res = await fetch('/api/v1/super-admin/alpha-invite-codes', {
        credentials: 'include',
        cache: 'no-store',
      });
      if (res.status === 401 || res.status === 403) {
        router.replace('/super-admin/login');
        return;
      }
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(typeof payload.error === 'string' ? payload.error : '초대 코드 목록을 불러오지 못했습니다.');
      setItems(Array.isArray(payload.items) ? payload.items : []);
    } catch (err) {
      setError(err instanceof Error ? err.message : '초대 코드 목록을 불러오지 못했습니다.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadItems();
  }, []);

  const handleCreate = async () => {
    setSubmitting(true);
    setError('');
    setMessage('');
    try {
      const res = await fetch('/api/v1/super-admin/alpha-invite-codes', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          expires_at: expiresAt,
          max_uses: maxUses,
          sent_to_note: sentToNote,
          admin_note: adminNote,
        }),
      });
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(typeof payload.error === 'string' ? payload.error : '초대 코드 발급에 실패했습니다.');
      setMessage(`초대 코드 ${payload.code}가 발급되었습니다.`);
      setSentToNote('');
      setAdminNote('');
      setMaxUses(1);
      await loadItems();
    } catch (err) {
      setError(err instanceof Error ? err.message : '초대 코드 발급에 실패했습니다.');
    } finally {
      setSubmitting(false);
    }
  };

  const handleRevoke = async (item: InviteCodeItem) => {
    if (item.state === 'revoked') return;
    setError('');
    setMessage('');
    try {
      const res = await fetch(`/api/v1/super-admin/alpha-invite-codes/${item.id}/revoke`, {
        method: 'PATCH',
        credentials: 'include',
      });
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(typeof payload.error === 'string' ? payload.error : '초대 코드 회수에 실패했습니다.');
      setMessage(`${item.code} 코드를 회수했습니다.`);
      await loadItems();
    } catch (err) {
      setError(err instanceof Error ? err.message : '초대 코드 회수에 실패했습니다.');
    }
  };

  const openExpiresAtPicker = () => {
    const input = expiresAtInputRef.current as (HTMLInputElement & { showPicker?: () => void }) | null;
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

  return (
    <main style={pageStyle}>
      <div style={shellStyle}>
        <SuperAdminPanelHeader
          subtitle="초대 코드"
          description="알파 테스트 접근 코드를 발급하고, 코드를 보낸 대상 메모와 실제 사용 계정을 확인합니다."
        />
        <SuperAdminPanelNav activeSection="alpha-invite" />

        <section style={panelStyle}>
          <div>
            <h2 style={sectionTitleStyle}>코드 발급</h2>
            <p style={sectionCopyStyle}>코드는 특정 이메일에 묶지 않습니다. 전달받은 사용자가 먼저 입력하면 해당 계정에 접근권이 부여됩니다.</p>
          </div>
          <div style={formGridStyle}>
            <label style={fieldStyle}>
              <span style={labelStyle}>입력 기간 종료일</span>
              <span style={dateInputShellStyle}>
                <input
                  ref={expiresAtInputRef}
                  type="date"
                  value={expiresAt}
                  onClick={openExpiresAtPicker}
                  onFocus={openExpiresAtPicker}
                  onChange={(event) => setExpiresAt(event.target.value)}
                  style={dateInputStyle}
                />
                <button type="button" onClick={openExpiresAtPicker} style={datePickerButtonStyle}>
                  달력 열기
                </button>
              </span>
            </label>
            <label style={fieldStyle}>
              <span style={labelStyle}>사용 가능 횟수</span>
              <input
                type="number"
                min={1}
                max={100}
                value={maxUses}
                onChange={(event) => setMaxUses(Math.max(1, Number(event.target.value) || 1))}
                style={inputStyle}
              />
            </label>
            <label style={{ ...fieldStyle, gridColumn: '1 / -1' }}>
              <span style={labelStyle}>발송 대상 메모</span>
              <input value={sentToNote} onChange={(event) => setSentToNote(event.target.value)} placeholder="예: 유튜브 API 검수 계정, 1차 알파 테스터 김OO" style={inputStyle} />
            </label>
            <label style={{ ...fieldStyle, gridColumn: '1 / -1' }}>
              <span style={labelStyle}>운영 메모</span>
              <textarea value={adminNote} onChange={(event) => setAdminNote(event.target.value)} rows={3} placeholder="발송 경로, 목적, 참고 사항" style={textareaStyle} />
            </label>
          </div>
          <button type="button" onClick={handleCreate} disabled={submitting} style={primaryButtonStyle}>
            {submitting ? '발급 중...' : '초대 코드 발급'}
          </button>
          {message ? <p style={successStyle}>{message}</p> : null}
          {error ? <p style={errorStyle}>{error}</p> : null}
        </section>

        <section style={panelStyle}>
          <div style={listHeaderStyle}>
            <div>
              <h2 style={sectionTitleStyle}>발급 코드 목록</h2>
              <p style={sectionCopyStyle}>최근 200개까지 표시합니다.</p>
            </div>
            <button type="button" onClick={loadItems} style={secondaryButtonStyle}>새로고침</button>
          </div>

          {loading ? (
            <p style={emptyStyle}>목록을 불러오는 중입니다.</p>
          ) : items.length === 0 ? (
            <p style={emptyStyle}>발급된 초대 코드가 없습니다.</p>
          ) : (
            <div style={tableWrapStyle}>
              <table style={tableStyle}>
                <thead>
                  <tr>
                    <th style={thStyle}>코드</th>
                    <th style={thStyle}>상태</th>
                    <th style={thStyle}>사용</th>
                    <th style={thStyle}>만료</th>
                    <th style={thStyle}>발송 대상 메모</th>
                    <th style={thStyle}>사용 이력</th>
                    <th style={thStyle}>관리</th>
                  </tr>
                </thead>
                <tbody>
                  {items.map((item) => (
                    <tr key={item.id}>
                      <td style={tdStyle}><code style={codeStyle}>{item.code}</code></td>
                      <td style={tdStyle}><span style={stateBadgeStyle(item.state)}>{stateLabels[item.state]}</span></td>
                      <td style={tdStyle}>{item.used_count} / {item.max_uses}</td>
                      <td style={tdStyle}>{formatKoreanDate(item.expires_at)}</td>
                      <td style={tdStyle}>{item.sent_to_note || '-'}</td>
                      <td style={tdStyle}>
                        {item.uses.length === 0 ? (
                          <span style={mutedStyle}>미사용</span>
                        ) : (
                          <div style={usesStyle}>
                            {item.uses.map((use) => (
                              <div key={use.id} style={useRecordStyle}>
                                <strong>{use.nickname || use.email || use.display_id || use.user_id}</strong>
                                <span>{use.provider || 'unknown'} · {formatKoreanDateTime(use.used_at)}</span>
                              </div>
                            ))}
                          </div>
                        )}
                      </td>
                      <td style={tdStyle}>
                        <button
                          type="button"
                          onClick={() => handleRevoke(item)}
                          disabled={item.state === 'revoked'}
                          style={dangerButtonStyle(item.state === 'revoked')}
                        >
                          회수
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      </div>
    </main>
  );
}

function defaultExpiresAtDate() {
  const date = new Date();
  date.setDate(date.getDate() + 14);
  return date.toISOString().slice(0, 10);
}

function formatKoreanDate(value: string) {
  if (!value) return '-';
  return new Date(value).toLocaleDateString('ko-KR');
}

function formatKoreanDateTime(value: string) {
  if (!value) return '-';
  return new Date(value).toLocaleString('ko-KR', { dateStyle: 'short', timeStyle: 'short' });
}

const pageStyle = {
  minHeight: '100vh',
  background: 'radial-gradient(circle at top, rgba(45, 78, 132, 0.2), transparent 42%), linear-gradient(180deg, #07111f, #0b1321 42%, #060b12)',
  color: '#eff6ff',
  padding: '28px 20px 48px',
};

const shellStyle = {
  width: SUPER_ADMIN_PAGE_WIDTH,
  margin: '0 auto',
  display: 'grid',
  gap: '18px',
};

const panelStyle = {
  display: 'grid',
  gap: '16px',
  padding: '20px',
  borderRadius: '18px',
  border: '1px solid rgba(194, 210, 245, 0.14)',
  background: 'rgba(255,255,255,0.05)',
};

const sectionTitleStyle = {
  margin: 0,
  fontSize: '20px',
  color: '#F7FAFF',
};

const sectionCopyStyle = {
  margin: '6px 0 0',
  fontSize: '13px',
  lineHeight: 1.7,
  color: '#B8C5E2',
};

const formGridStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))',
  gap: '12px',
};

const fieldStyle = {
  display: 'grid',
  gap: '7px',
};

const labelStyle = {
  fontSize: '12px',
  fontWeight: 800,
  color: 'rgba(220,228,245,0.78)',
};

const inputStyle = {
  minHeight: '40px',
  borderRadius: '10px',
  border: '1px solid rgba(160,186,224,0.2)',
  background: 'rgba(7,18,34,0.72)',
  color: '#F7FAFF',
  padding: '0 12px',
  fontSize: '13px',
  fontFamily: 'inherit',
};

const dateInputShellStyle = {
  display: 'flex',
  alignItems: 'center',
  minHeight: '40px',
  borderRadius: '10px',
  border: '1px solid rgba(160,186,224,0.2)',
  background: 'rgba(7,18,34,0.72)',
  overflow: 'hidden',
};

const dateInputStyle = {
  ...inputStyle,
  flex: '1 1 auto',
  minWidth: 0,
  border: 'none',
  borderRadius: 0,
  background: 'transparent',
};

const datePickerButtonStyle = {
  alignSelf: 'stretch',
  flex: '0 0 auto',
  border: 'none',
  borderLeft: '1px solid rgba(160,186,224,0.18)',
  background: 'rgba(95, 131, 255, 0.16)',
  color: '#E8EEFF',
  padding: '0 12px',
  fontSize: '12px',
  fontWeight: 800,
  fontFamily: 'inherit',
  cursor: 'pointer',
};

const textareaStyle = {
  ...inputStyle,
  minHeight: '82px',
  padding: '10px 12px',
  resize: 'vertical' as const,
};

const primaryButtonStyle = {
  justifySelf: 'start',
  minHeight: '40px',
  borderRadius: '10px',
  border: '1px solid rgba(125, 160, 255, 0.34)',
  background: 'rgba(95, 131, 255, 0.22)',
  color: '#F7FAFF',
  padding: '0 14px',
  fontSize: '13px',
  fontWeight: 800,
  fontFamily: 'inherit',
  cursor: 'pointer',
};

const secondaryButtonStyle = {
  ...primaryButtonStyle,
  background: 'rgba(255,255,255,0.07)',
  color: 'rgba(220,228,245,0.88)',
};

const listHeaderStyle = {
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'flex-start',
  gap: '12px',
  flexWrap: 'wrap' as const,
};

const tableWrapStyle = {
  overflowX: 'auto' as const,
  borderRadius: '14px',
  border: '1px solid rgba(160,186,224,0.12)',
};

const tableStyle = {
  width: '100%',
  borderCollapse: 'collapse' as const,
  minWidth: '980px',
};

const thStyle = {
  textAlign: 'left' as const,
  padding: '12px',
  background: 'rgba(255,255,255,0.05)',
  color: 'rgba(220,228,245,0.68)',
  fontSize: '12px',
};

const tdStyle = {
  padding: '12px',
  borderTop: '1px solid rgba(160,186,224,0.1)',
  verticalAlign: 'top' as const,
  fontSize: '13px',
  color: 'rgba(232,238,255,0.88)',
};

const codeStyle = {
  color: '#F7FAFF',
  fontWeight: 800,
  letterSpacing: '0.04em',
};

const stateBadgeStyle = (state: InviteCodeItem['state']) => ({
  display: 'inline-flex',
  alignItems: 'center',
  minHeight: '24px',
  padding: '0 9px',
  borderRadius: '999px',
  fontSize: '12px',
  fontWeight: 800,
  color: state === 'available' ? '#DFFFE8' : state === 'used_up' ? '#E8EEFF' : '#FFD8D8',
  background: state === 'available' ? 'rgba(50, 200, 100, 0.16)' : state === 'used_up' ? 'rgba(95, 131, 255, 0.18)' : 'rgba(255, 110, 110, 0.14)',
  border: state === 'available' ? '1px solid rgba(50, 200, 100, 0.3)' : state === 'used_up' ? '1px solid rgba(125, 160, 255, 0.28)' : '1px solid rgba(255, 110, 110, 0.26)',
});

const usesStyle = {
  display: 'grid',
  gap: '8px',
};

const useRecordStyle = {
  display: 'grid',
  gap: '3px',
  color: 'rgba(220,228,245,0.82)',
};

const mutedStyle = {
  color: 'rgba(190,204,232,0.56)',
};

const emptyStyle = {
  margin: 0,
  padding: '18px',
  borderRadius: '14px',
  background: 'rgba(255,255,255,0.04)',
  color: 'rgba(220,228,245,0.66)',
  fontSize: '13px',
};

const successStyle = {
  margin: 0,
  color: '#B9FFC9',
  fontSize: '13px',
};

const errorStyle = {
  margin: 0,
  color: '#FFB4B4',
  fontSize: '13px',
};

const dangerButtonStyle = (disabled: boolean) => ({
  minHeight: '32px',
  borderRadius: '9px',
  border: '1px solid rgba(255,110,110,0.24)',
  background: 'rgba(255,110,110,0.12)',
  color: '#FFD8D8',
  padding: '0 10px',
  fontSize: '12px',
  fontWeight: 800,
  fontFamily: 'inherit',
  opacity: disabled ? 0.44 : 1,
  cursor: disabled ? 'not-allowed' : 'pointer',
});
