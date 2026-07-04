'use client';

import { type ReactNode } from 'react';
import {
  infoCardStyle,
  infoCardHeaderStyle,
  sectionTitleStyle,
  infoCardSubtitleStyle,
  detailRowStyle,
  detailLabelStyle,
  detailValueStyle,
  fieldStyle,
  labelStyle,
  rangeFieldRowStyle,
  rangeInputStyle,
  numberInputStyle,
  toggleChipStyle,
  toggleChipActiveStyle,
} from './lumiLabStyles';

export function InfoCard({
  title,
  subtitle,
  children,
}: {
  title: string;
  subtitle: string;
  children: ReactNode;
}) {
  return (
    <section style={infoCardStyle}>
      <div style={infoCardHeaderStyle}>
        <h2 style={sectionTitleStyle}>{title}</h2>
        <p style={infoCardSubtitleStyle}>{subtitle}</p>
      </div>
      {children}
    </section>
  );
}

export function DetailRow({ label, value }: { label: string; value: string }) {
  return (
    <div style={detailRowStyle}>
      <span style={detailLabelStyle}>{label}</span>
      <code style={detailValueStyle}>{value}</code>
    </div>
  );
}

export function RangeField({
  label,
  min,
  max,
  value,
  onChange,
  step = 1,
}: {
  label: string;
  min: number;
  max: number;
  value: number;
  onChange: (next: number) => void;
  step?: number;
}) {
  return (
    <label style={fieldStyle}>
      <span style={labelStyle}>{label}</span>
      <div style={rangeFieldRowStyle}>
        <input type="range" min={min} max={max} step={step} value={value} onChange={(event) => onChange(Number(event.target.value))} style={rangeInputStyle} />
        <input type="number" min={min} max={max} step={step} value={value} onChange={(event) => onChange(Number(event.target.value) || 0)} style={numberInputStyle} />
      </div>
    </label>
  );
}

export function ToggleChip({ label, checked, onToggle }: { label: string; checked: boolean; onToggle: () => void }) {
  return (
    <button type="button" onClick={onToggle} style={{ ...toggleChipStyle, ...(checked ? toggleChipActiveStyle : null) }}>
      {label}
    </button>
  );
}
