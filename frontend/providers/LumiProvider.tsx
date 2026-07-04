'use client';

import { createContext, useEffect, useRef, useState, type ReactNode } from 'react';

import { createLumiViewState, getLumiTriggerPriority } from '@/lib/lumi/lumiTriggers';
import type { LumiPageContext, LumiTrigger, LumiTriggerPayload, LumiViewState } from '@/lib/lumi/lumiTypes';

type LumiPageContextValue = {
  viewState: LumiViewState;
  isMobile: boolean;
  prefersReducedMotion: boolean;
  show: (next: Partial<LumiViewState>) => void;
  hide: () => void;
  update: (next: Partial<LumiViewState>) => void;
  resetForContext: (context: LumiPageContext) => void;
  trigger: (event: LumiTrigger, payload?: LumiTriggerPayload) => void;
};

const initialViewState: LumiViewState = {
  state: 'idle',
  context: 'dashboard',
  dockSlot: 'dashboard-map-top',
  visible: false,
  expanded: false,
  message: '새로운 행성을 만들거나, 이어서 탐험할 수 있어요.',
  messageType: 'summary',
  quickActions: [],
  runtimeContext: {
    mode: 'general',
    page: 'dashboard',
    action: 'page_enter',
    scene: 'overview',
  },
};

export const LumiContextObject = createContext<LumiPageContextValue | null>(null);

interface LumiProviderProps {
  children: ReactNode;
}

export function LumiProvider({ children }: LumiProviderProps) {
  const [viewState, setViewState] = useState<LumiViewState>(initialViewState);
  const [isMobile, setIsMobile] = useState(false);
  const [prefersReducedMotion, setPrefersReducedMotion] = useState(false);
  const activePriorityRef = useRef(0);
  const hideTimerRef = useRef<number | null>(null);

  useEffect(() => {
    const mediaQuery = window.matchMedia('(max-width: 820px)');
    const motionQuery = window.matchMedia('(prefers-reduced-motion: reduce)');

    const handleViewportChange = () => setIsMobile(mediaQuery.matches);
    const handleMotionChange = () => setPrefersReducedMotion(motionQuery.matches);

    handleViewportChange();
    handleMotionChange();

    mediaQuery.addEventListener('change', handleViewportChange);
    motionQuery.addEventListener('change', handleMotionChange);

    return () => {
      mediaQuery.removeEventListener('change', handleViewportChange);
      motionQuery.removeEventListener('change', handleMotionChange);
    };
  }, []);

  useEffect(() => () => {
    if (hideTimerRef.current !== null) {
      window.clearTimeout(hideTimerRef.current);
    }
  }, []);

  const show = (next: Partial<LumiViewState>) => {
    setViewState((current) => ({
      ...current,
      ...next,
      visible: next.visible ?? true,
    }));
  };

  const hide = () => {
    setViewState((current) => ({ ...current, visible: false, expanded: false }));
    activePriorityRef.current = 0;
  };

  const update = (next: Partial<LumiViewState>) => {
    setViewState((current) => ({ ...current, ...next }));
  };

  const resetForContext = (context: LumiPageContext) => {
    activePriorityRef.current = 0;
    setViewState(
      createLumiViewState('page_enter', { context }, { isMobile, expanded: false }),
    );
  };

  const trigger = (event: LumiTrigger, payload?: LumiTriggerPayload) => {
    const nextPriority = getLumiTriggerPriority(event);
    if (nextPriority < activePriorityRef.current) {
      return;
    }

    activePriorityRef.current = nextPriority;
    const nextState = createLumiViewState(event, payload, {
      isMobile,
      expanded: payload?.context === 'mobile' ? true : false,
    });

    setViewState(nextState);

    if (hideTimerRef.current !== null) {
      window.clearTimeout(hideTimerRef.current);
      hideTimerRef.current = null;
    }

    if (event !== 'manual_open' && event !== 'course_complete' && event !== 'base_built') {
      hideTimerRef.current = window.setTimeout(() => {
        activePriorityRef.current = 0;
        setViewState((current) => ({ ...current, visible: false }));
      }, prefersReducedMotion ? 1800 : 3200);
    }
  };

  return (
    <LumiContextObject.Provider
      value={{
        viewState,
        isMobile,
        prefersReducedMotion,
        show,
        hide,
        update,
        resetForContext,
        trigger,
      }}
    >
      {children}
    </LumiContextObject.Provider>
  );
}
