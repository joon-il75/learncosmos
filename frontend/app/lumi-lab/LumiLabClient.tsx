'use client';

import { useState } from 'react';
import LumiBottomSheet from '@/components/lumi/LumiBottomSheet';
import SuperAdminPanelNav from '@/components/super-admin/SuperAdminPanelNav';
import SuperAdminPanelHeader from '@/components/super-admin/SuperAdminPanelHeader';
import { LumiProvider } from '@/providers/LumiProvider';
import { type LumiLabTab } from './lumiLabTypes';
import { useLumiLabAssets } from './useLumiLabAssets';
import { useLumiLabRuntimeConfig } from './useLumiLabRuntimeConfig';
import { useLumiLabUIPreview } from './useLumiLabUIPreview';
import LumiSpritesTab from './tabs/LumiSpritesTab';
import LumiScenariosTab from './tabs/LumiScenariosTab';
import LumiRulesTab from './tabs/LumiRulesTab';
import LumiMessagesTab from './tabs/LumiMessagesTab';
import LumiUIPresetsTab from './tabs/LumiUIPresetsTab';
import {
  pageStyle, backdropStyle, shellStyle, topChromeStyle, introCardStyle, introTextStyle,
  introTitleStyle, introCopyStyle, tabSectionStyle, tabGridStyle, tabButtonStyle,
  tabButtonActiveStyle, tabLabelStyle, tabDescriptionStyle, tabMetaStyle, tabMetaCopyStyle,
  runtimeConfigPanelStyle, runtimeConfigMetaStyle, ruleActionRowStyle, runtimeConfigStatusStyle,
  runtimeConfigDirtyStyle, codePillStyle, sectionTitleStyle, authGateStyle, authCardStyle,
  primaryButtonStyle, secondaryActionButtonStyle,
} from './lumiLabStyles';

const lumiLabTabs: Array<{ id: LumiLabTab; label: string; description: string }> = [
  { id: 'sprites', label: 'Sprites', description: '스프라이트 시트 업로드와 좌표 미세 조정' },
  { id: 'scenarios', label: 'Scenarios', description: '페이지별 runtime context와 적용 순서 점검' },
  { id: 'rules', label: 'Rules', description: 'trigger, adapter, runtime engine 연결 상태 확인' },
  { id: 'messages', label: 'Messages', description: '메시지 타입과 현재 preview 문구 검토' },
  { id: 'ui', label: 'UI Presets', description: 'dock slot, panel, mobile sheet 프리셋 정리' },
];

export default function LumiLabClientPage() {
  return (
    <LumiProvider>
      <LumiLabClientInner />
    </LumiProvider>
  );
}

function LumiLabClientInner() {
  const [activeTab, setActiveTab] = useState<LumiLabTab>('sprites');

  const rc = useLumiLabRuntimeConfig({ onNavigateToTab: setActiveTab });
  const uip = useLumiLabUIPreview();
  const assets = useLumiLabAssets({
    onRuntimeConfigLoad: rc.loadRuntimeConfig,
    onRuntimeConfigError: rc.setRuntimeConfigStatus,
  });

  const applyRuleToPreview = () => {
    const rule = rc.selectedRuleDraft;
    if (!rule) return;
    uip.setSelectedState(rule.state);
    uip.setSelectedContext(rule.context);
    uip.setSelectedDockSlot(rule.dockSlot);
    uip.setMessageType(rule.messageType);
    uip.setMessage(rule.message);
    uip.setManualVisible(true);
    uip.setManualExpanded(true);
    uip.setManualMobile(rule.context === 'mobile');
  };

  const applyMessageToPreview = () => {
    const msg = rc.selectedMessageDraft;
    if (!msg) return;
    uip.setMessageType(msg.messageType);
    uip.setMessage(msg.message);
    uip.setManualVisible(true);
  };

  const activeTabMeta = lumiLabTabs.find((tab) => tab.id === activeTab) ?? lumiLabTabs[0];

  if (assets.isAuthorized !== true) {
    return (
      <main style={authGateStyle}>
        <div style={authCardStyle}>슈퍼관리자 권한을 확인하는 중입니다.</div>
      </main>
    );
  }

  return (
    <main style={pageStyle}>
      <div style={backdropStyle} />
      <div style={shellStyle}>
        <div style={topChromeStyle}>
          <SuperAdminPanelHeader subtitle="루미설정" description="Lumi 자산, 상태, 메시지, 스프라이트 업로드를 슈퍼관리자 전용 설정 흐름에서 관리합니다." />
          <SuperAdminPanelNav activeSection="experience" activeExperience="lumi" />
        </div>

        <section style={introCardStyle}>
          <div style={introTextStyle}>
            <strong style={introTitleStyle}>Lumi 상태 미세 조정 페이지</strong>
            <p style={introCopyStyle}>
              슈퍼관리자 전용 페이지입니다. 상태, 메시지, 도킹 슬롯, 모바일 모드, 스프라이트 위치를 바꿔 보면서 공통 Lumi 시스템용 기준값을 확인하고, 생성 프롬프트와 실제 스프라이트 업로드까지 관리할 수 있습니다.
            </p>
          </div>
          <button type="button" onClick={() => { uip.setSheetOpen(true); uip.setManualExpanded(true); }} style={primaryButtonStyle}>
            Bottom Sheet 보기
          </button>
        </section>

        <section style={tabSectionStyle}>
          <div style={tabGridStyle}>
            {lumiLabTabs.map((tab) => (
              <button key={tab.id} type="button" onClick={() => setActiveTab(tab.id)} style={{ ...tabButtonStyle, ...(activeTab === tab.id ? tabButtonActiveStyle : null) }}>
                <strong style={tabLabelStyle}>{tab.label}</strong>
                <span style={tabDescriptionStyle}>{tab.description}</span>
              </button>
            ))}
          </div>
          <div style={tabMetaStyle}>
            <span style={codePillStyle}>{activeTabMeta.label}</span>
            <span style={tabMetaCopyStyle}>{activeTabMeta.description}</span>
          </div>
        </section>

        {(activeTab === 'rules' || activeTab === 'messages') ? (
          <section style={runtimeConfigPanelStyle}>
            <div style={runtimeConfigMetaStyle}>
              <strong style={sectionTitleStyle}>Runtime Config</strong>
              <span style={tabMetaCopyStyle}>Rules와 Messages를 저장하면 실제 Lumi runtime source에 반영됩니다.</span>
              <span style={runtimeConfigStatusStyle}>{rc.runtimeConfigStatus ?? '저장 전 변경은 editor 로컬 상태에만 남습니다.'}</span>
              <span style={runtimeConfigDirtyStyle}>{rc.hasUnsavedRuntimeChanges ? '저장되지 않은 변경 있음' : '저장 상태 동기화됨'}</span>
            </div>
            <div style={ruleActionRowStyle}>
              <span style={codePillStyle}>{rc.runtimeConfigVersion ? `v${rc.runtimeConfigVersion}` : 'unsaved'}</span>
              <button
                type="button"
                onClick={() => {
                  if (rc.hasUnsavedRuntimeChanges && !window.confirm('저장하지 않은 변경을 버리고 마지막 저장 상태를 다시 불러올까요?')) return;
                  void rc.loadRuntimeConfig();
                }}
                style={secondaryActionButtonStyle}
              >
                다시 불러오기
              </button>
              <button type="button" onClick={rc.saveRuntimeConfig} disabled={rc.isSavingRuntimeConfig || !rc.hasUnsavedRuntimeChanges} style={primaryButtonStyle}>
                {rc.isSavingRuntimeConfig ? '저장 중...' : rc.hasUnsavedRuntimeChanges ? 'Rules + Messages 저장' : '저장 완료됨'}
              </button>
            </div>
          </section>
        ) : null}

        {activeTab === 'sprites' ? <LumiSpritesTab assets={assets} uip={uip} /> : null}
        {activeTab === 'scenarios' ? <LumiScenariosTab uip={uip} /> : null}
        {activeTab === 'rules' ? <LumiRulesTab rc={rc} onApplyToPreview={applyRuleToPreview} /> : null}
        {activeTab === 'messages' ? <LumiMessagesTab rc={rc} uip={uip} onApplyToPreview={applyMessageToPreview} /> : null}
        {activeTab === 'ui' ? <LumiUIPresetsTab uip={uip} /> : null}
      </div>

      <LumiBottomSheet
        open={uip.sheetOpen}
        state={uip.selectedState}
        message={uip.message}
        quickActions={uip.lumi.viewState.quickActions}
        onClose={() => uip.setSheetOpen(false)}
        reducedMotion={uip.lumi.prefersReducedMotion}
      />
    </main>
  );
}
