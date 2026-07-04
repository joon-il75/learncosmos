'use client';

import type { CSSProperties } from 'react';

import ModalPortal from '@/components/common/ModalPortal';
import LumiPanel from './LumiPanel';
import type { LumiQuickAction, LumiState } from '@/lib/lumi/lumiTypes';

interface LumiBottomSheetProps {
  open: boolean;
  state: LumiState;
  message: string;
  quickActions?: LumiQuickAction[];
  onClose: () => void;
  reducedMotion?: boolean;
}

export default function LumiBottomSheet({
  open,
  state,
  message,
  quickActions,
  onClose,
  reducedMotion = false,
}: LumiBottomSheetProps) {
  if (!open) return null;

  return (
    <ModalPortal overlayStyle={overlayStyle}>
      <div style={sheetStyle}>
        <button type="button" onClick={onClose} style={sheetCloseStyle} aria-label="Lumi 모바일 패널 닫기">
          닫기
        </button>
        <div role="dialog" aria-modal="true" aria-label="Lumi 모바일 패널">
          <LumiPanel
            open
            state={state}
            message={message}
            quickActions={quickActions}
            onClose={onClose}
            reducedMotion={reducedMotion}
          />
        </div>
      </div>
    </ModalPortal>
  );
}

const overlayStyle: CSSProperties = {
  padding: '24px 16px',
  background: 'rgba(4, 8, 15, 0.46)',
};

const sheetStyle: CSSProperties = {
  width: 'min(420px, 100%)',
  maxHeight: 'calc(100dvh - 48px)',
  overflowY: 'auto',
};

const sheetCloseStyle: CSSProperties = {
  position: 'absolute',
  width: 1,
  height: 1,
  overflow: 'hidden',
  clipPath: 'inset(50%)',
};
