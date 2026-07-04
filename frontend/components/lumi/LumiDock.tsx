'use client';

import type { CSSProperties, ReactNode } from 'react';

import type { LumiDockSlot } from '@/lib/lumi/lumiTypes';

interface LumiDockProps {
  slot: LumiDockSlot;
  children: ReactNode;
  visible?: boolean;
}

export default function LumiDock({ slot, children, visible = true }: LumiDockProps) {
  return (
    <div
      data-lumi-slot={slot}
      style={{
        ...dockStyle,
        opacity: visible ? 1 : 0,
        pointerEvents: visible ? 'auto' : 'none',
      }}
    >
      {children}
    </div>
  );
}

const dockStyle: CSSProperties = {
  display: 'inline-flex',
  alignItems: 'flex-start',
  gap: 10,
  transition: 'opacity 180ms ease',
};
