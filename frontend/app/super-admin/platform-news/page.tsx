'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import SuperAdminPanelHeader from '@/components/super-admin/SuperAdminPanelHeader';
import SuperAdminPanelNav from '@/components/super-admin/SuperAdminPanelNav';
import { superAdminShellStyle } from '@/components/super-admin/layout';

interface PlatformNotice {
  id: string;
  slug: string;
  locale: 'ko' | 'en';
  title: string;
  summary: string;
  body: string;
  status: 'draft' | 'published' | 'archived';
  pinned: boolean;
  published_at?: string;
  updated_at: string;
}

interface NoticeFormState {
  id: string;
  slug: string;
  locale: 'ko' | 'en';
  title: string;
  summary: string;
  body: string;
  status: 'draft' | 'published' | 'archived';
  pinned: boolean;
}

const emptyForm: NoticeFormState = {
  id: '',
  slug: '',
  locale: 'ko',
  title: '',
  summary: '',
  body: '',
  status: 'draft',
  pinned: false,
};

export default function SuperAdminPlatformNewsPage() {
  const router = useRouter();
  const [notices, setNotices] = useState<PlatformNotice[]>([]);
  const [form, setForm] = useState<NoticeFormState>(emptyForm);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  const isEditing = useMemo(() => form.id !== '', [form.id]);

  const fetchNotices = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const res = await fetch('/api/v1/super-admin/platform-notices?limit=50', { credentials: 'include', cache: 'no-store' });
      if (res.status === 401 || res.status === 403) { router.replace('/super-admin/login'); return; }
      if (!res.ok) throw new Error('조회 실패');
      const data = await res.json();
      setNotices(Array.isArray(data.notices) ? data.notices : []);
    } catch (err) {
      setError(err instanceof Error ? err.message : '공지사항을 불러오지 못했습니다.');
    } finally {
      setLoading(false);
    }
  }, [router]);

  useEffect(() => { fetchNotices(); }, [fetchNotices]);

  const resetForm = () => {
    setForm(emptyForm);
    setMessage('');
    setError('');
  };

  const editNotice = (notice: PlatformNotice) => {
    setForm({
      id: notice.id,
      slug: notice.slug,
      locale: notice.locale,
      title: notice.title,
      summary: notice.summary,
      body: notice.body,
      status: notice.status,
      pinned: notice.pinned,
    });
    setMessage('수정할 공지를 불러왔습니다.');
    setError('');
  };

  const saveNotice = async () => {
    if (!form.title.trim()) {
      setError('제목을 입력해 주세요.');
      return;
    }
    setSaving(true);
    setError('');
    setMessage('');
    try {
      const endpoint = isEditing ? '/api/v1/super-admin/platform-notices/' + form.id : '/api/v1/super-admin/platform-notices';
      const res = await fetch(endpoint, {
        method: isEditing ? 'PATCH' : 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          slug: form.slug,
          locale: form.locale,
          title: form.title,
          summary: form.summary,
          body: form.body,
          status: form.status,
          pinned: form.pinned,
        }),
      });
      if (res.status === 401 || res.status === 403) { router.replace('/super-admin/login'); return; }
      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        throw new Error(data.error || '저장 실패');
      }
      await fetchNotices();
      setForm(emptyForm);
      setMessage(isEditing ? '공지사항을 수정했습니다.' : '공지사항을 생성했습니다.');
    } catch (err) {
      setError(err instanceof Error ? err.message : '저장하지 못했습니다.');
    } finally {
      setSaving(false);
    }
  };

  const deleteNotice = async (notice: PlatformNotice) => {
    if (!window.confirm(`'${notice.title}' 공지를 삭제할까요?`)) return;
    setSaving(true);
    setError('');
    setMessage('');
    try {
      const res = await fetch('/api/v1/super-admin/platform-notices/' + notice.id, { method: 'DELETE', credentials: 'include' });
      if (res.status === 401 || res.status === 403) { router.replace('/super-admin/login'); return; }
      if (!res.ok) throw new Error('삭제 실패');
      await fetchNotices();
      if (form.id === notice.id) setForm(emptyForm);
      setMessage('공지사항을 삭제했습니다.');
    } catch (err) {
      setError(err instanceof Error ? err.message : '삭제하지 못했습니다.');
    } finally {
      setSaving(false);
    }
  };

  return (
    <main style={pageStyle}>
      <div style={superAdminShellStyle}>
        <div style={headerStackStyle}>
          <SuperAdminPanelHeader
            subtitle="새소식 관리"
            description="플랫폼 새소식(공지사항)의 작성, 수정, 공개 상태 관리를 위한 슈퍼관리자 전용 화면입니다. 공개 /platform/news는 읽기 전용으로 유지됩니다."
          />
          <SuperAdminPanelNav activeSection="platform-news" />
        </div>

        <section style={summaryGridStyle} aria-label="새소식 관리 범위">
          <article style={summaryCardStyle}>
            <span style={summaryKickerStyle}>Public</span>
            <strong style={summaryTitleStyle}>공개 페이지</strong>
            <p style={summaryCopyStyle}>/platform/news는 published 공지만 읽기 전용으로 보여줍니다.</p>
            <Link href="/platform/news" style={summaryLinkStyle}>공개 페이지 보기</Link>
          </article>
          <article style={summaryCardStyle}>
            <span style={summaryKickerStyle}>Super Admin</span>
            <strong style={summaryTitleStyle}>작성 권한</strong>
            <p style={summaryCopyStyle}>작성, 수정, 삭제, 고정 공지, 공개 상태 변경은 이 화면에서만 처리합니다.</p>
            <span style={statusPillStyle}>서버 권한 검사 적용</span>
          </article>
          <article style={summaryCardStyle}>
            <span style={summaryKickerStyle}>Status</span>
            <strong style={summaryTitleStyle}>현재 등록 {notices.length}건</strong>
            <p style={summaryCopyStyle}>draft는 관리자만 볼 수 있고 published만 공개 페이지에 표시됩니다.</p>
            <span style={statusPillStyle}>API 연결 완료</span>
          </article>
        </section>

        {message ? <div style={messageStyle}>{message}</div> : null}
        {error ? <div style={errorStyle}>{error}</div> : null}

        <section style={workspaceGridStyle}>
          <article style={editorCardStyle} aria-label="공지 작성 폼">
            <div style={sectionHeadingStyle}>
              <span style={sectionEyebrowStyle}>{isEditing ? 'Edit Notice' : 'Create Notice'}</span>
              <h2 style={sectionTitleStyle}>{isEditing ? '공지 수정' : '공지 작성'}</h2>
              <p style={sectionCopyStyle}>published로 저장하면 공개 /platform/news에 표시됩니다.</p>
            </div>

            <div style={formGridStyle}>
              <label style={fieldStyle}>
                <span style={labelStyle}>제목</span>
                <input value={form.title} onChange={(e) => setForm({ ...form, title: e.target.value })} placeholder="예: 알파 테스트 2차 모집 안내" style={inputStyle} />
              </label>
              <label style={fieldStyle}>
                <span style={labelStyle}>슬러그</span>
                <input value={form.slug} onChange={(e) => setForm({ ...form, slug: e.target.value })} placeholder="비워두면 notice-타임스탬프 생성" style={inputStyle} />
              </label>
              <label style={fieldStyle}>
                <span style={labelStyle}>요약</span>
                <input value={form.summary} onChange={(e) => setForm({ ...form, summary: e.target.value })} placeholder="목록과 홈 preview에 표시될 짧은 설명" style={inputStyle} />
              </label>
              <label style={fieldStyle}>
                <span style={labelStyle}>본문</span>
                <textarea value={form.body} onChange={(e) => setForm({ ...form, body: e.target.value })} placeholder="공지 상세 본문 입력 영역" rows={8} style={textareaStyle} />
              </label>
              <div style={optionRowStyle} aria-label="공지 옵션">
                <label style={selectLabelStyle}>언어
                  <select value={form.locale} onChange={(e) => setForm({ ...form, locale: e.target.value as 'ko' | 'en' })} style={selectStyle}>
                    <option value="ko">한국어</option>
                    <option value="en">English</option>
                  </select>
                </label>
                <label style={selectLabelStyle}>상태
                  <select value={form.status} onChange={(e) => setForm({ ...form, status: e.target.value as NoticeFormState['status'] })} style={selectStyle}>
                    <option value="draft">draft</option>
                    <option value="published">published</option>
                    <option value="archived">archived</option>
                  </select>
                </label>
                <label style={checkboxStyle}>
                  <input type="checkbox" checked={form.pinned} onChange={(e) => setForm({ ...form, pinned: e.target.checked })} />
                  고정 공지
                </label>
              </div>
              <div style={actionRowStyle}>
                <button type="button" onClick={saveNotice} disabled={saving} style={primaryButtonStyle}>{saving ? '저장 중...' : isEditing ? '수정 저장' : '공지 생성'}</button>
                <button type="button" onClick={resetForm} disabled={saving} style={secondaryButtonStyle}>새 작성</button>
              </div>
            </div>
          </article>

          <aside style={listCardStyle} aria-label="공지 목록">
            <div style={sectionHeadingStyle}>
              <span style={sectionEyebrowStyle}>Notice List</span>
              <h2 style={sectionTitleStyle}>등록된 공지</h2>
              <p style={sectionCopyStyle}>상태와 고정 여부를 확인하고 수정할 수 있습니다.</p>
            </div>

            <div style={tableShellStyle}>
              <div style={tableHeaderStyle}>
                <span>제목</span>
                <span>상태</span>
                <span>관리</span>
              </div>
              {loading ? (
                <div style={emptyStateStyle}><strong style={emptyTitleStyle}>불러오는 중...</strong></div>
              ) : notices.length === 0 ? (
                <div style={emptyStateStyle}>
                  <strong style={emptyTitleStyle}>등록된 공지사항이 없습니다</strong>
                  <p style={emptyCopyStyle}>첫 공지를 작성하면 이곳에 표시됩니다.</p>
                </div>
              ) : notices.map((notice) => (
                <div key={notice.id} style={tableRowStyle}>
                  <span style={{ display: 'grid', gap: 4 }}>
                    <strong>{notice.pinned ? '★ ' : ''}{notice.title}</strong>
                    <small style={{ color: '#8FA0C4' }}>{notice.locale} · {notice.slug} · {new Date(notice.updated_at).toLocaleString('ko-KR')}</small>
                  </span>
                  <span style={statusCellStyle(notice.status)}>{notice.status}</span>
                  <span style={{ display: 'flex', gap: 6, flexWrap: 'wrap' }}>
                    <button type="button" onClick={() => editNotice(notice)} style={smallButtonStyle}>수정</button>
                    <button type="button" onClick={() => deleteNotice(notice)} style={dangerButtonStyle}>삭제</button>
                  </span>
                </div>
              ))}
            </div>
          </aside>
        </section>
      </div>
    </main>
  );
}

const pageStyle = { minHeight: '100vh', background: 'radial-gradient(circle at top, rgba(45, 78, 132, 0.2), transparent 42%), linear-gradient(180deg, #07111f, #0b1321 42%, #060b12)', color: '#eff6ff', padding: '32px 20px 56px' };
const headerStackStyle = { display: 'grid', gap: '16px', marginBottom: '18px' };
const summaryGridStyle = { display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: '14px', marginBottom: '18px' };
const summaryCardStyle = { display: 'grid', gap: '10px', padding: '18px', borderRadius: '20px', border: '1px solid rgba(194, 210, 245, 0.14)', background: 'rgba(255,255,255,0.055)' };
const summaryKickerStyle = { color: '#8AB4FF', fontSize: '11px', fontWeight: 800, letterSpacing: '0.12em', textTransform: 'uppercase' as const };
const summaryTitleStyle = { color: '#F7FAFF', fontSize: '18px' };
const summaryCopyStyle = { margin: 0, color: '#B8C5E2', fontSize: '13px', lineHeight: 1.68 };
const summaryLinkStyle = { width: 'fit-content', minHeight: '32px', display: 'inline-flex', alignItems: 'center', padding: '0 12px', borderRadius: '999px', textDecoration: 'none', color: '#E8EEFF', background: 'rgba(95, 131, 255, 0.18)', border: '1px solid rgba(125, 160, 255, 0.34)', fontSize: '12px', fontWeight: 800 };
const statusPillStyle = { width: 'fit-content', minHeight: '30px', display: 'inline-flex', alignItems: 'center', padding: '0 11px', borderRadius: '999px', color: '#DFFFE8', background: 'rgba(50, 200, 100, 0.14)', border: '1px solid rgba(50, 200, 100, 0.28)', fontSize: '12px', fontWeight: 800 };
const messageStyle = { marginBottom: 14, padding: '12px 14px', borderRadius: 12, background: 'rgba(50, 200, 100, 0.14)', border: '1px solid rgba(50, 200, 100, 0.3)', color: '#DFFFE8', fontSize: 13, fontWeight: 800 };
const errorStyle = { ...messageStyle, background: 'rgba(185, 28, 28, 0.24)', border: '1px solid rgba(248, 113, 113, 0.34)', color: '#FCA5A5' };
const workspaceGridStyle = { display: 'grid', gridTemplateColumns: 'minmax(0, 1.12fr) minmax(320px, 0.88fr)', gap: '16px' };
const editorCardStyle = { display: 'grid', gap: '18px', padding: '22px', borderRadius: '24px', border: '1px solid rgba(194, 210, 245, 0.14)', background: 'rgba(7, 18, 34, 0.72)', boxShadow: '0 24px 64px rgba(0,0,0,0.18)' };
const listCardStyle = { ...editorCardStyle, alignSelf: 'start' };
const sectionHeadingStyle = { display: 'grid', gap: '8px' };
const sectionEyebrowStyle = { color: '#8AB4FF', fontSize: '12px', fontWeight: 800, letterSpacing: '0.12em', textTransform: 'uppercase' as const };
const sectionTitleStyle = { margin: 0, color: '#F7FAFF', fontSize: '24px', letterSpacing: '-0.03em' };
const sectionCopyStyle = { margin: 0, color: '#B8C5E2', fontSize: '13px', lineHeight: 1.7 };
const formGridStyle = { display: 'grid', gap: '14px' };
const fieldStyle = { display: 'grid', gap: '7px' };
const labelStyle = { color: '#D7E0F6', fontSize: '13px', fontWeight: 800 };
const inputStyle = { width: '100%', minHeight: '42px', borderRadius: '12px', border: '1px solid rgba(120,140,200,0.22)', background: 'rgba(255,255,255,0.06)', color: '#E8EEFF', padding: '0 13px', fontSize: '13px' };
const textareaStyle = { ...inputStyle, minHeight: '156px', padding: '12px 13px', resize: 'vertical' as const };
const optionRowStyle = { display: 'flex', flexWrap: 'wrap' as const, gap: '8px', alignItems: 'end' };
const selectLabelStyle = { display: 'grid', gap: 6, color: '#D7E0F6', fontSize: 12, fontWeight: 800 };
const selectStyle = { minHeight: 34, borderRadius: 10, border: '1px solid rgba(120,140,200,0.22)', background: '#101B2D', color: '#E8EEFF', padding: '0 10px' };
const checkboxStyle = { minHeight: 34, display: 'inline-flex', alignItems: 'center', gap: 8, padding: '0 12px', borderRadius: 999, color: '#D7E0F6', background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(120,140,200,0.18)', fontSize: 12, fontWeight: 800 };
const actionRowStyle = { display: 'flex', flexWrap: 'wrap' as const, gap: '10px' };
const primaryButtonStyle = { minHeight: '40px', border: '0', borderRadius: '12px', padding: '0 15px', color: '#fff', background: 'rgba(95, 131, 255, 0.82)', fontWeight: 900, cursor: 'pointer' };
const secondaryButtonStyle = { ...primaryButtonStyle, background: 'rgba(255,255,255,0.08)', color: '#E8EEFF', border: '1px solid rgba(120,140,200,0.2)' };
const tableShellStyle = { display: 'grid', overflow: 'hidden', borderRadius: '18px', border: '1px solid rgba(120,140,200,0.18)' };
const tableHeaderStyle = { display: 'grid', gridTemplateColumns: '1fr 86px 112px', gap: '10px', padding: '12px 14px', background: 'rgba(255,255,255,0.06)', color: '#8FA0C4', fontSize: '12px', fontWeight: 800 };
const tableRowStyle = { ...tableHeaderStyle, background: 'transparent', color: '#D7E0F6', alignItems: 'center', borderTop: '1px solid rgba(120,140,200,0.12)' };
const emptyStateStyle = { display: 'grid', placeItems: 'center', gap: '8px', minHeight: '220px', padding: '24px', textAlign: 'center' as const };
const emptyTitleStyle = { color: '#F7FAFF', fontSize: '18px' };
const emptyCopyStyle = { margin: 0, maxWidth: '360px', color: '#B8C5E2', fontSize: '13px', lineHeight: 1.7 };
const smallButtonStyle = { minHeight: 30, borderRadius: 9, border: '1px solid rgba(125,160,255,0.28)', background: 'rgba(95,131,255,0.14)', color: '#E8EEFF', fontSize: 12, fontWeight: 800, cursor: 'pointer' };
const dangerButtonStyle = { ...smallButtonStyle, border: '1px solid rgba(248,113,113,0.3)', background: 'rgba(185,28,28,0.2)', color: '#FCA5A5' };
const statusCellStyle = (status: PlatformNotice['status']) => ({ width: 'fit-content', borderRadius: 999, padding: '5px 9px', background: status === 'published' ? 'rgba(50,200,100,0.14)' : status === 'archived' ? 'rgba(148,163,184,0.14)' : 'rgba(239,159,39,0.15)', color: status === 'published' ? '#DFFFE8' : status === 'archived' ? '#CBD5E1' : '#FDD68A', fontSize: 12, fontWeight: 900 });
