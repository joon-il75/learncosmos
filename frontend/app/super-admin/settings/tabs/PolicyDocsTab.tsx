'use client'
import type { PolicyDocument } from '../settingsTypes'
import { S } from '../settingsStyles'
import { DocStateTag, Toggle } from '../settingsComponents'

interface Props {
  policyDocs: PolicyDocument[]
  policyDocsLoaded: boolean
  saving: boolean
  settingsLoadStatus: Record<string, string>
  savePolicyDrafts: () => Promise<void>
  confirmPolicyDraft: (policyType: 'terms' | 'privacy', locale: 'ko' | 'en') => Promise<void>
  updatePolicyDraftField: (
    policyType: 'terms' | 'privacy',
    locale: 'ko' | 'en',
    field: keyof Pick<PolicyDocument, 'title' | 'effective_at' | 'required' | 'content'>,
    value: string | boolean,
  ) => void
}

export function PolicyDocsTab({
  policyDocs, policyDocsLoaded, saving, settingsLoadStatus,
  savePolicyDrafts, confirmPolicyDraft, updatePolicyDraftField,
}: Props) {
  return (
    <div style={{ paddingTop:'20px' }}>
      <div style={{ ...S.infoBox, marginBottom:'14px' }}>
        설정 API 상태:
        <br />`/settings/llm`: {settingsLoadStatus['/api/v1/super-admin/settings/llm'] || '대기'}
        <br />`/settings/api-keys`: {settingsLoadStatus['/api/v1/super-admin/settings/api-keys'] || '대기'}
        <br />`/settings/point-policy`: {settingsLoadStatus['/api/v1/super-admin/settings/point-policy'] || '대기'}
        <br />`/settings/products`: {settingsLoadStatus['/api/v1/super-admin/settings/products'] || '대기'}
        <br />`/settings/ad-slots`: {settingsLoadStatus['/api/v1/super-admin/settings/ad-slots'] || '대기'}
        <br />`/settings/affiliate`: {settingsLoadStatus['/api/v1/super-admin/settings/affiliate'] || '대기'}
        <br />`/settings/policies`: {settingsLoadStatus['/api/v1/super-admin/settings/policies'] || '대기'}
      </div>

      <div style={S.infoBox}>
        📄 활성 문서만 사용자에게 노출됩니다.<br/>
        ✍️ 먼저 초안을 저장한 뒤, 검토 후 `확정`해야 새 버전이 공개됩니다.<br/>
        ✅ 사용자는 현재 활성 필수 버전에 동의해야 서비스를 사용할 수 있습니다.
      </div>

      <div style={{ ...S.warnBox, marginBottom:'14px' }}>
        정책 문서 로드 상태: {policyDocsLoaded ? '로드 완료' : '로드 중'} · 현재 메모리 문서 수: {policyDocs.length}
      </div>

      {(['terms', 'privacy'] as const).flatMap((policyType) => (['ko', 'en'] as const).map((locale) => {
        const activeDoc = policyDocs.find(item => item.type === policyType && item.locale === locale && item.active)
        const draftDoc = policyDocs.find(item => item.type === policyType && item.locale === locale && item.draft)
        const doc = draftDoc || activeDoc
        const localeLabel = locale === 'ko' ? '한국어' : 'English'
        return (
          <div key={`${policyType}-${locale}`} style={S.card}>
            <div style={S.cardTitle}>
              {policyType === 'terms' ? '서비스 이용약관' : '개인정보처리방침'} · {localeLabel}
            </div>
            <div style={{ ...S.cardSub, marginBottom:'12px' }}>
              아래 편집 영역은 {draftDoc ? '초안 버전' : '운영 버전'} 기준으로 표시됩니다. locale={locale}
            </div>

            <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(240px, 1fr))', gap:'10px', marginBottom:'16px' }}>
              <div style={{ background:'rgba(255,255,255,0.03)', border:'1px solid rgba(120,140,200,0.12)', borderRadius:'10px', padding:'12px 14px' }}>
                <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', gap:'10px', marginBottom:'8px' }}>
                  <strong style={{ fontSize:'12px', color:'#E8EAF2' }}>현재 운영 버전</strong>
                  <DocStateTag label={activeDoc ? `운영 v${activeDoc.version}` : '운영본 없음'} active />
                </div>
                <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.58)', lineHeight:1.7 }}>
                  {activeDoc?.title || '운영 중인 문서가 없습니다'}
                  <br />
                  {activeDoc?.published_at ? `확정 시각: ${new Date(activeDoc.published_at).toLocaleString('ko-KR')}` : '확정 시각 정보 없음'}
                </div>
              </div>
              <div style={{ background:'rgba(255,255,255,0.03)', border:'1px solid rgba(239,159,39,0.18)', borderRadius:'10px', padding:'12px 14px' }}>
                <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', gap:'10px', marginBottom:'8px' }}>
                  <strong style={{ fontSize:'12px', color:'#E8EAF2' }}>확정 전 초안</strong>
                  <DocStateTag label={draftDoc ? `초안 v${draftDoc.version}` : '초안 없음'} />
                </div>
                <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.58)', lineHeight:1.7 }}>
                  {draftDoc?.title || '아직 저장된 초안이 없습니다'}
                  <br />
                  {draftDoc?.updated_at ? `최근 수정: ${new Date(draftDoc.updated_at).toLocaleString('ko-KR')}` : '초안 저장 후 여기 표시됩니다'}
                </div>
              </div>
            </div>

            <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(320px, 1fr))', gap:'12px', marginBottom:'16px' }}>
              <div style={{ background:'rgba(255,255,255,0.02)', border:'1px solid rgba(120,140,200,0.12)', borderRadius:'12px', padding:'14px' }}>
                <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', gap:'10px', marginBottom:'10px' }}>
                  <strong style={{ fontSize:'12px', color:'#E8EAF2' }}>현재 운영 버전 본문</strong>
                  <DocStateTag label={activeDoc ? `운영 v${activeDoc.version}` : '운영본 없음'} active />
                </div>
                <textarea readOnly style={{ ...S.inp, minHeight:'220px', resize:'vertical' as const, lineHeight:1.7, opacity: activeDoc ? 1 : 0.6 }}
                  value={activeDoc?.content || '현재 운영 중인 문서가 없습니다.'} />
              </div>
              <div style={{ background:'rgba(255,255,255,0.02)', border:'1px solid rgba(239,159,39,0.18)', borderRadius:'12px', padding:'14px' }}>
                <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', gap:'10px', marginBottom:'10px' }}>
                  <strong style={{ fontSize:'12px', color:'#E8EAF2' }}>초안 편집 대상 본문</strong>
                  <DocStateTag label={draftDoc ? `초안 v${draftDoc.version}` : activeDoc ? `운영 v${activeDoc.version}` : '문서 없음'} />
                </div>
                <textarea readOnly style={{ ...S.inp, minHeight:'220px', resize:'vertical' as const, lineHeight:1.7, opacity: doc ? 1 : 0.6 }}
                  value={doc?.content || '편집할 문서가 없습니다.'} />
              </div>
            </div>

            <div style={S.grid2}>
              <div>
                <label style={S.label}>문서 제목</label>
                <input style={S.inp} value={doc?.title || ''}
                  onChange={e => updatePolicyDraftField(policyType, locale, 'title', e.target.value)} />
              </div>
              <div>
                <label style={S.label}>시행 시각 (RFC3339)</label>
                <input style={S.inp} value={doc?.effective_at || ''}
                  onChange={e => updatePolicyDraftField(policyType, locale, 'effective_at', e.target.value)}
                  placeholder="2026-03-31T00:00:00Z" />
              </div>
            </div>

            <div style={{ marginTop:'12px' }}>
              <label style={S.label}>필수 동의 여부</label>
              <div style={{ display:'flex', alignItems:'center', gap:'10px' }}>
                <Toggle value={doc?.required ?? true} onChange={(value) => updatePolicyDraftField(policyType, locale, 'required', value)} />
                <span style={{ fontSize:'12px', color:'rgba(200,210,235,0.58)' }}>
                  {doc?.required ? '필수 동의 문서' : '선택 문서'}
                </span>
              </div>
            </div>

            <div style={{ marginTop:'12px' }}>
              <label style={S.label}>문서 본문</label>
              <textarea style={{ ...S.inp, minHeight:'320px', resize:'vertical' as const, lineHeight:1.7 }}
                value={doc?.content || ''}
                onChange={e => updatePolicyDraftField(policyType, locale, 'content', e.target.value)} />
            </div>

            {!draftDoc && activeDoc ? (
              <div style={{ ...S.infoBox, marginTop:'14px', marginBottom:'0' }}>
                초안이 없어서 현재 운영 버전을 편집할 수 없습니다. 아래 `초안 저장`을 누르면 다음 버전 초안을 만들고 그 초안을 수정하게 됩니다.
              </div>
            ) : null}

            <div style={{ display:'flex', justifyContent:'flex-end', gap:'10px', marginTop:'14px' }}>
              <button style={S.btnSecondary} onClick={savePolicyDrafts} disabled={saving}>초안 저장</button>
              <button style={S.btnSave} onClick={() => confirmPolicyDraft(policyType, locale)} disabled={saving || !draftDoc}>확정</button>
            </div>
          </div>
        )
      }))}
    </div>
  )
}
