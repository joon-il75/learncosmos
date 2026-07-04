'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import SuperAdminPanelNav from '@/components/super-admin/SuperAdminPanelNav'
import SuperAdminPanelHeader from '@/components/super-admin/SuperAdminPanelHeader'
import { superAdminShellStyle } from '@/components/super-admin/layout'
import { useSettingsData } from './useSettingsData'
import { S } from './settingsStyles'
import { LLMSettingsTab } from './tabs/LLMSettingsTab'
import { PointPolicyTab } from './tabs/PointPolicyTab'
import { ProductsTab } from './tabs/ProductsTab'
import { PolicyDocsTab } from './tabs/PolicyDocsTab'
import { AdsSettingsTab } from './tabs/AdsSettingsTab'
import { UIEngineTab } from './tabs/UIEngineTab'

export default function SettingsPage() {
  const router = useRouter()
  const [tab, setTab] = useState<'llm'|'point'|'product'|'ads'|'ui'|'policy'>('llm')
  const data = useSettingsData()
  const { isAuthorized, toast } = data

  if (isAuthorized === false) {
    return (
      <main style={{
        minHeight: '100vh', background: '#081826', color: '#E8EAF2',
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        padding: '24px', fontFamily: "'Noto Sans KR', sans-serif",
      }}>
        <div style={{
          width: '100%', maxWidth: '560px',
          background: 'rgba(17,30,53,0.88)', border: '1px solid rgba(120,140,200,0.2)',
          borderRadius: '20px', padding: '28px', boxShadow: '0 24px 64px rgba(0,0,0,0.35)',
        }}>
          <h1 style={{ fontSize: '24px', fontWeight: 700, marginBottom: '10px' }}>슈퍼관리자 설정 페이지</h1>
          <p style={{ fontSize: '14px', lineHeight: 1.7, color: 'rgba(200,210,235,0.72)', marginBottom: '14px' }}>
            이 페이지는 웹 초기화 페이지가 아니라 로그인된 슈퍼관리자 전용 설정 화면입니다.
          </p>
          <div style={S.infoBox}>
            현재 구현에서 슈퍼관리자 초기값은 웹에서 생성하지 않고 서버 환경변수
            `SUPER_ADMIN_ID`, `SUPER_ADMIN_PW_HASH`, `SUPER_ADMIN_TOTP_SECRET`
            로 부트스트랩합니다.
          </div>
          <div style={{ display: 'flex', gap: '10px', marginTop: '18px' }}>
            <button type="button" onClick={() => router.push('/super-admin/login')} style={S.btnSave}>로그인으로 이동</button>
            <button type="button" onClick={() => router.push('/')} style={S.btnNav}>홈으로 이동</button>
          </div>
        </div>
      </main>
    )
  }

  return (
    <main style={S.page}>
      <div style={superAdminShellStyle}>
        {toast && (
          <div style={{ position:'fixed', bottom:'24px', right:'24px', background:'rgba(82,183,136,0.9)', color:'white', padding:'12px 20px', borderRadius:'10px', fontSize:'13px', fontWeight:600, zIndex:999 }}>
            ✓ {toast}
          </div>
        )}

        <SuperAdminPanelHeader subtitle="시스템 설정" description="LLM, API 키, 포인트 정책, 상품, 광고, 제휴, 정책 문서를 조정합니다." />
        <SuperAdminPanelNav activeSection="settings" />

        <div style={S.tabs}>
          {(['llm','point','product','ads','ui','policy'] as const).map(t => (
            <div key={t} style={S.tab(tab===t)} onClick={() => setTab(t)}>
              {{ llm:'🤖 LLM 설정', point:'💎 포인트 정책', product:'💳 구매 상품', ads:'📢 광고 설정', ui:'🪐 UI 엔진', policy:'📄 약관 문서' }[t]}
            </div>
          ))}
        </div>

        {tab === 'llm' && (
          <LLMSettingsTab
            llmSettings={data.llmSettings}
            apiKeys={data.apiKeys} setApiKeys={data.setApiKeys}
            apiKey2s={data.apiKey2s} setApiKey2s={data.setApiKey2s}
            endpoints={data.endpoints} setEndpoints={data.setEndpoints}
            apiKeyStatuses={data.apiKeyStatuses}
            apiKey2Statuses={data.apiKey2Statuses}
            apiKeyUpdatedAts={data.apiKeyUpdatedAts}
            saving={data.saving}
            saveLLM={data.saveLLM} saveAPIKey={data.saveAPIKey}
            keyStatusText={data.keyStatusText} updateLLM={data.updateLLM}
          />
        )}
        {tab === 'point' && (
          <PointPolicyTab
            policy={data.policy} setPolicy={data.setPolicy}
            saving={data.saving} savePolicy={data.savePolicy}
          />
        )}
        {tab === 'product' && (
          <ProductsTab
            products={data.products} setProducts={data.setProducts}
            saving={data.saving} saveProducts={data.saveProducts}
          />
        )}
        {tab === 'policy' && (
          <PolicyDocsTab
            policyDocs={data.policyDocs} policyDocsLoaded={data.policyDocsLoaded}
            saving={data.saving} settingsLoadStatus={data.settingsLoadStatus}
            savePolicyDrafts={data.savePolicyDrafts} confirmPolicyDraft={data.confirmPolicyDraft}
            updatePolicyDraftField={data.updatePolicyDraftField}
          />
        )}
        {tab === 'ads' && (
          <AdsSettingsTab
            adSlots={data.adSlots} setAdSlots={data.setAdSlots}
            affiliate={data.affiliate} setAffiliate={data.setAffiliate}
            saving={data.saving} saveAdSlots={data.saveAdSlots} saveAffiliate={data.saveAffiliate}
          />
        )}
        {tab === 'ui' && (
          <UIEngineTab
            uiEngineSettings={data.uiEngineSettings} setUiEngineSettings={data.setUiEngineSettings}
            saving={data.saving} settingsLoadStatus={data.settingsLoadStatus}
            saveUIEngineSettings={data.saveUIEngineSettings}
          />
        )}
      </div>
    </main>
  )
}
