'use client';

import { type CSSProperties, useEffect, useMemo, useRef, useState } from 'react';
import { useLumi } from '@/hooks/useLumi';
import {
  getDefaultSpriteTuning,
  lumiSpriteMap,
  LUMI_SPRITE_COLUMNS,
  LUMI_SPRITE_ROWS,
} from '@/lib/lumi/lumiSpriteMap';
import { createLumiViewState } from '@/lib/lumi/lumiTriggers';
import type { LumiPageContext, LumiDockSlot, LumiMessageType, LumiState, LumiTrigger } from '@/lib/lumi/lumiTypes';
import { stateOptions, contextOptions, dockSlotOptions } from './lumiLabTypes';

export type LumiLabUIPreviewResult = {
  lumi: ReturnType<typeof useLumi>;
  selectedState: LumiState;
  setSelectedState: (s: LumiState) => void;
  selectedContext: LumiPageContext;
  setSelectedContext: (c: LumiPageContext) => void;
  selectedDockSlot: LumiDockSlot;
  setSelectedDockSlot: (s: LumiDockSlot) => void;
  messageType: LumiMessageType;
  setMessageType: (t: LumiMessageType) => void;
  message: string;
  setMessage: (m: string) => void;
  size: number;
  setSize: (s: number) => void;
  sheetOpen: boolean;
  setSheetOpen: (o: boolean) => void;
  manualVisible: boolean;
  setManualVisible: (v: boolean | ((prev: boolean) => boolean)) => void;
  manualExpanded: boolean;
  setManualExpanded: (e: boolean | ((prev: boolean) => boolean)) => void;
  manualMobile: boolean;
  setManualMobile: (m: boolean | ((prev: boolean) => boolean)) => void;
  backgroundSizeX: number;
  setBackgroundSizeX: (x: number) => void;
  backgroundSizeY: number;
  setBackgroundSizeY: (y: number) => void;
  backgroundPositionX: number;
  setBackgroundPositionX: (x: number) => void;
  backgroundPositionY: number;
  setBackgroundPositionY: (y: number) => void;
  manualSpriteStyle: CSSProperties;
  spriteCoordinate: (typeof lumiSpriteMap)[LumiState];
  selectedStateTuning: ReturnType<typeof getDefaultSpriteTuning>;
  handleTrigger: (trigger: LumiTrigger) => void;
  stateOptions: typeof stateOptions;
  contextOptions: typeof contextOptions;
  dockSlotOptions: typeof dockSlotOptions;
};

export function useLumiLabUIPreview(): LumiLabUIPreviewResult {
  const lumi = useLumi();
  const [selectedState, setSelectedState] = useState<LumiState>('idle');
  const [selectedContext, setSelectedContext] = useState<LumiPageContext>('dashboard');
  const [selectedDockSlot, setSelectedDockSlot] = useState<LumiDockSlot>('dashboard-map-top');
  const [messageType, setMessageType] = useState<LumiMessageType>('summary');
  const [message, setMessage] = useState('새로운 행성을 만들거나, 이어서 탐험할 수 있어요.');
  const [size, setSize] = useState(88);
  const [sheetOpen, setSheetOpen] = useState(false);
  const [manualVisible, setManualVisible] = useState(true);
  const [manualExpanded, setManualExpanded] = useState(true);
  const [manualMobile, setManualMobile] = useState(false);
  const [backgroundSizeX, setBackgroundSizeX] = useState(LUMI_SPRITE_COLUMNS * 100);
  const [backgroundSizeY, setBackgroundSizeY] = useState(LUMI_SPRITE_ROWS * 100);
  const [backgroundPositionX, setBackgroundPositionX] = useState(0);
  const [backgroundPositionY, setBackgroundPositionY] = useState(0);

  const lumiRef = useRef(lumi);
  lumiRef.current = lumi;

  useEffect(() => {
    lumiRef.current.show({
      state: selectedState,
      context: selectedContext,
      dockSlot: selectedDockSlot,
      visible: manualVisible,
      expanded: manualExpanded,
      message,
      messageType,
    });
  // lumiRef는 ref이므로 deps에서 제외 — lumi를 직접 넣으면 show()→setViewState→context변경→effect재실행 무한루프 발생
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedState, selectedContext, selectedDockSlot, manualVisible, manualExpanded, message, messageType]);

  useEffect(() => {
    const tuning = getDefaultSpriteTuning(selectedState);
    setBackgroundSizeX(tuning.backgroundSizeX);
    setBackgroundSizeY(tuning.backgroundSizeY);
    setBackgroundPositionX(tuning.backgroundPositionX);
    setBackgroundPositionY(tuning.backgroundPositionY);
  }, [selectedState]);

  const spriteCoordinate = lumiSpriteMap[selectedState];
  const selectedStateTuning = getDefaultSpriteTuning(selectedState);

  const manualSpriteStyle = useMemo<CSSProperties>(() => ({
    width: size,
    height: size,
    borderRadius: '50%',
    backgroundImage: "url('/images/lumi.webp')",
    backgroundRepeat: 'no-repeat',
    backgroundSize: `${backgroundSizeX}% ${backgroundSizeY}%`,
    backgroundPosition: `${backgroundPositionX}% ${backgroundPositionY}%`,
    boxShadow: '0 16px 32px rgba(0, 0, 0, 0.25)',
    backgroundColor: 'rgba(5, 12, 22, 0.65)',
  }), [backgroundPositionX, backgroundPositionY, backgroundSizeX, backgroundSizeY, size]);

  const handleTrigger = (trigger: LumiTrigger) => {
    const triggered = createLumiViewState(trigger, {
      context: selectedContext,
      dockSlot: selectedDockSlot,
      courseTitle: '수채화 첫 행성',
      lessonTitle: '색 혼합 기초',
      status: 'learning',
      actionHandlers: {
        'new-course': () => setMessage('새 행성 만들기 액션 테스트'),
        'resume-course': () => setMessage('진행 중인 탐험 보기 액션 테스트'),
        'status-summary': () => setMessage('지금 상태 요약 액션 테스트'),
      },
    }, { isMobile: manualMobile, expanded: manualExpanded });

    setSelectedState(triggered.state);
    setSelectedContext(triggered.context);
    setSelectedDockSlot(triggered.dockSlot);
    setMessageType(triggered.messageType);
    setMessage(triggered.message);
    setManualVisible(triggered.visible);
    setManualExpanded(triggered.expanded);
    lumi.trigger(trigger, {
      context: selectedContext,
      dockSlot: selectedDockSlot,
      courseTitle: '수채화 첫 행성',
      lessonTitle: '색 혼합 기초',
      status: 'learning',
      actionHandlers: {
        'new-course': () => setMessage('새 행성 만들기 액션 테스트'),
        'resume-course': () => setMessage('진행 중인 탐험 보기 액션 테스트'),
        'status-summary': () => setMessage('지금 상태 요약 액션 테스트'),
      },
    });
  };

  return {
    lumi, selectedState, setSelectedState, selectedContext, setSelectedContext,
    selectedDockSlot, setSelectedDockSlot, messageType, setMessageType,
    message, setMessage, size, setSize, sheetOpen, setSheetOpen,
    manualVisible, setManualVisible, manualExpanded, setManualExpanded,
    manualMobile, setManualMobile,
    backgroundSizeX, setBackgroundSizeX, backgroundSizeY, setBackgroundSizeY,
    backgroundPositionX, setBackgroundPositionX, backgroundPositionY, setBackgroundPositionY,
    manualSpriteStyle, spriteCoordinate, selectedStateTuning, handleTrigger,
    stateOptions, contextOptions, dockSlotOptions,
  };
}
