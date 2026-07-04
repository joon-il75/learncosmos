'use client';

import { runtimeModeOptions, runtimePageOptions, runtimeActionOptions, runtimeSceneOptions } from '../lumiLabTypes';
import type { LumiLabUIPreviewResult } from '../useLumiLabUIPreview';
import { InfoCard, DetailRow } from '../LumiLabComponents';
import { infoGridStyle, bulletListStyle } from '../lumiLabStyles';

interface Props {
  uip: LumiLabUIPreviewResult;
}

export default function LumiScenariosTab({ uip }: Props) {
  const { lumi, selectedContext, selectedDockSlot, messageType } = uip;
  const activeRuntimeContext = lumi.viewState.runtimeContext;

  return (
    <section style={infoGridStyle}>
      <InfoCard title="Runtime Context Model" subtitle="LumiRuntimeContext 기준 입력 모델">
        <DetailRow label="mode" value={runtimeModeOptions.join(', ')} />
        <DetailRow label="page" value={runtimePageOptions.join(', ')} />
        <DetailRow label="action" value={runtimeActionOptions.join(', ')} />
        <DetailRow label="scene" value={runtimeSceneOptions.join(', ')} />
      </InfoCard>

      <InfoCard title="Current Preview Scenario" subtitle="지금 Lab 프리뷰가 만드는 상태">
        <DetailRow label="legacy context" value={selectedContext} />
        <DetailRow label="dock slot" value={selectedDockSlot} />
        <DetailRow label="message type" value={messageType} />
        <DetailRow label="runtime mode" value={activeRuntimeContext?.mode ?? '미연결'} />
        <DetailRow label="runtime page" value={activeRuntimeContext?.page ?? '미연결'} />
        <DetailRow label="runtime action" value={activeRuntimeContext?.action ?? '미연결'} />
        <DetailRow label="runtime scene" value={activeRuntimeContext?.scene ?? '미연결'} />
      </InfoCard>

      <InfoCard title="Safe Rollout Order" subtitle="문서 기준 페이지 적용 순서">
        <ul style={bulletListStyle}>
          <li>Dashboard를 기준 엔진으로 먼저 안정화하고, Hero는 action 체계만 얹어서 따라갑니다.</li>
          <li>Draft Editor와 Player는 page/action/scene이 더 세분화된 뒤에 붙이는 것이 안전합니다.</li>
          <li>이 탭은 아직 저장 기능 없이 운영 기준만 보여 주므로 기존 동작을 깨지 않습니다.</li>
        </ul>
      </InfoCard>
    </section>
  );
}
