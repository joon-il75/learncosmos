'use client';

import type { PageSystemChunk } from '@/components/dashboard/GalaxyMap';
import type { DashboardMainCopy } from '@/lib/i18n/pages/dashboardMain';
import {
  mobileInstructionStyle,
  mobileInstructionTitleStyle,
  mobileInstructionTextStyle,
  mobileInstructionSubtitleStyle,
} from '@/components/dashboard/GalaxyMapStyles';

interface Props {
  mode: 'galaxy' | 'star-system';
  selectedSystemChunk: PageSystemChunk | null;
  copy: DashboardMainCopy['galaxy']['mobile'];
}

export function GalaxyMobileInstruction({ mode, selectedSystemChunk, copy }: Props) {
  const key = mode === 'star-system' ? 'starSystem' : 'galaxy';
  const title = copy.title[key];
  const primaryText = copy.primary[key];
  const secondaryText = copy.secondary[key];
  const subtitle = selectedSystemChunk ? copy.current(selectedSystemChunk.summaryTitle) : undefined;

  return (
    <section style={mobileInstructionStyle}>
      <div style={mobileInstructionTitleStyle}>{title}</div>
      <div style={mobileInstructionTextStyle}>{primaryText}</div>
      <div style={mobileInstructionTextStyle}>{secondaryText}</div>
      {subtitle ? <div style={mobileInstructionSubtitleStyle}>{subtitle}</div> : null}
    </section>
  );
}
