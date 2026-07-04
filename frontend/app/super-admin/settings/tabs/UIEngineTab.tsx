'use client'
import type { Dispatch, SetStateAction } from 'react'
import type { UIEngineSettings } from '../settingsTypes'
import { DEFAULT_UI_ENGINE_SETTINGS } from '../settingsConstants'
import { S } from '../settingsStyles'
import { Toggle } from '../settingsComponents'

interface Props {
  uiEngineSettings: UIEngineSettings | null
  setUiEngineSettings: Dispatch<SetStateAction<UIEngineSettings | null>>
  saving: boolean
  settingsLoadStatus: Record<string, string>
  saveUIEngineSettings: () => Promise<void>
}

export function UIEngineTab({
  uiEngineSettings, setUiEngineSettings, saving, settingsLoadStatus, saveUIEngineSettings,
}: Props) {
  return (
    <div style={{ paddingTop:'20px' }}>
      <div style={S.infoBox}>
        1차 공개 범위만 노출합니다.
        <br />
        Dashboard UI 엔진의 운영값만 조정하고, 상태 정규화·selection policy·route sync 같은 핵심 규칙은 코드 고정으로 유지합니다.
      </div>
      <div style={S.warnBox}>
        변경 직후 영향 범위:
        <br />`compact_breakpoint`, `phone_breakpoint`, `short_viewport_height` → 레이아웃 반응형 규칙
        <br />`cta_min_launch_duration_ms`, `selection_cta_launch_delay_ms` → CTA 연출 시간
        <br />`lumi_quick_action_limit`, `lumi_*_enabled` → Dashboard Lumi 표시 정책
      </div>
      <div style={S.card}>
        <div style={S.cardTitle}>Dashboard UI 엔진 운영값</div>
        <div style={S.cardSub}>1차로 안전한 숫자 임계값과 토글만 운영 설정으로 분리합니다.</div>
        <div style={S.grid2}>
          {([
            { key:'compact_breakpoint', label:'Compact Breakpoint' },
            { key:'phone_breakpoint', label:'Phone Breakpoint' },
            { key:'short_viewport_height', label:'Short Viewport Height' },
            { key:'cta_min_launch_duration_ms', label:'CTA Min Launch Duration (ms)' },
            { key:'selection_cta_launch_delay_ms', label:'Selection CTA Launch Delay (ms)' },
            { key:'lumi_quick_action_limit', label:'Lumi Quick Action Limit' },
          ] as { key: Exclude<keyof UIEngineSettings, 'lumi_desktop_panel_enabled' | 'lumi_mobile_sheet_enabled'>; label: string }[]).map(({ key, label }) => (
            <div key={key}>
              <label style={S.label}>{label}</label>
              <input type="number" style={S.inp}
                value={uiEngineSettings?.[key] ?? DEFAULT_UI_ENGINE_SETTINGS[key] as number}
                onChange={e => setUiEngineSettings(prev => ({
                  ...(prev ?? DEFAULT_UI_ENGINE_SETTINGS),
                  [key]: parseInt(e.target.value) || 0,
                }))} />
            </div>
          ))}
        </div>

        <div style={{ ...S.grid2, marginTop:'14px' }}>
          <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', background:'rgba(255,255,255,0.03)', border:'1px solid rgba(120,140,200,0.12)', borderRadius:'10px', padding:'12px 14px' }}>
            <div>
              <div style={{ fontSize:'13px', fontWeight:700, color:'#E8EAF2' }}>Lumi Desktop Panel</div>
              <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.5)', marginTop:'4px' }}>데스크톱 Dashboard 우측 Lumi 패널 표시 여부</div>
            </div>
            <Toggle
              value={uiEngineSettings?.lumi_desktop_panel_enabled ?? DEFAULT_UI_ENGINE_SETTINGS.lumi_desktop_panel_enabled}
              onChange={value => setUiEngineSettings(prev => ({ ...(prev ?? DEFAULT_UI_ENGINE_SETTINGS), lumi_desktop_panel_enabled: value }))}
            />
          </div>
          <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', background:'rgba(255,255,255,0.03)', border:'1px solid rgba(120,140,200,0.12)', borderRadius:'10px', padding:'12px 14px' }}>
            <div>
              <div style={{ fontSize:'13px', fontWeight:700, color:'#E8EAF2' }}>Lumi Mobile Sheet</div>
              <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.5)', marginTop:'4px' }}>모바일 Dashboard Lumi 패널 표시 여부</div>
            </div>
            <Toggle
              value={uiEngineSettings?.lumi_mobile_sheet_enabled ?? DEFAULT_UI_ENGINE_SETTINGS.lumi_mobile_sheet_enabled}
              onChange={value => setUiEngineSettings(prev => ({ ...(prev ?? DEFAULT_UI_ENGINE_SETTINGS), lumi_mobile_sheet_enabled: value }))}
            />
          </div>
        </div>

        <div style={{ ...S.infoBox, marginTop:'14px', marginBottom:'0' }}>
          설정 API 상태:
          <br />`/settings/ui-engine`: {settingsLoadStatus['/api/v1/super-admin/settings/ui-engine'] || '대기'}
        </div>

        <div style={S.saveRow}>
          <button style={S.btnSave} onClick={saveUIEngineSettings} disabled={saving || !uiEngineSettings}>
            UI 엔진 설정 저장
          </button>
        </div>
      </div>
    </div>
  )
}
