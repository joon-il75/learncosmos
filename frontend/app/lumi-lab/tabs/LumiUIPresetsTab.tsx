'use client';

import { dockSlotOptions } from '../lumiLabTypes';
import type { LumiLabUIPreviewResult } from '../useLumiLabUIPreview';
import { InfoCard, DetailRow } from '../LumiLabComponents';
import { infoGridStyle, tagGroupStyle, tagStyle, bulletListStyle } from '../lumiLabStyles';

interface Props {
  uip: LumiLabUIPreviewResult;
}

export default function LumiUIPresetsTab({ uip }: Props) {
  const { manualVisible, manualExpanded, manualMobile, size, selectedDockSlot } = uip;

  return (
    <section style={infoGridStyle}>
      <InfoCard title="Dock Slots" subtitle="현재 배치 가능한 프리셋 위치">
        <div style={tagGroupStyle}>
          {dockSlotOptions.map((slot) => <span key={slot} style={tagStyle}>{slot}</span>)}
        </div>
      </InfoCard>

      <InfoCard title="Component Roles" subtitle="UI preset이 다루는 현재 요소">
        <ul style={bulletListStyle}>
          <li>`LumiAvatar`: sprite sheet 기반 상태 렌더링</li>
          <li>`LumiBubble` / `LumiDock`: 인라인 안내와 위치 배치</li>
          <li>`LumiPanel` / `LumiBottomSheet`: desktop 확장 패널과 mobile sheet</li>
        </ul>
      </InfoCard>

      <InfoCard title="Current UI State" subtitle="Lab에서 즉시 확인 가능한 값">
        <DetailRow label="visible" value={manualVisible ? 'true' : 'false'} />
        <DetailRow label="expanded" value={manualExpanded ? 'true' : 'false'} />
        <DetailRow label="mobile mode" value={manualMobile ? 'true' : 'false'} />
        <DetailRow label="avatar size" value={`${size}px`} />
        <DetailRow label="active dock" value={selectedDockSlot} />
      </InfoCard>
    </section>
  );
}
