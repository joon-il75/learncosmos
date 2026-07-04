import type { LumiRuntimeContext } from './lumiEngineTypes';

export type LumiState =
  | 'idle'
  | 'thinking'
  | 'happy'
  | 'celebrate'
  | 'victory'
  | 'encourage'
  | 'surprise'
  | 'curious'
  | 'raise-hand'
  | 'sleepy'
  | 'exploring'
  | 'focus'
  | 'planet-point'
  | 'planet-hold'
  | 'note-read';

export type LumiPageContext =
  | 'dashboard'
  | 'planet-map'
  | 'planet-scene'
  | 'diary'
  | 'player'
  | 'completion'
  | 'mobile';

export type LumiDockSlot =
  | 'dashboard-cta'
  | 'dashboard-map-top'
  | 'dashboard-map-bottom'
  | 'dashboard-floating'
  | 'dashboard-panel'
  | 'diary-header'
  | 'diary-sidebar'
  | 'player-inline'
  | 'player-bottom'
  | 'planet-scene-anchor'
  | 'completion-center'
  | 'mobile-bottom-sheet-trigger'
  | 'mobile-inline';

export type LumiMessageType =
  | 'summary'
  | 'hint'
  | 'question'
  | 'encourage'
  | 'success'
  | 'guide'
  | 'discovery';

export type LumiTrigger =
  | 'page_enter'
  | 'cta_focus'
  | 'cta_submit'
  | 'generation_loading'
  | 'generation_success'
  | 'planet_hover'
  | 'planet_select'
  | 'idle'
  | 'search_start'
  | 'search_success'
  | 'evaluation_start'
  | 'evaluation_success'
  | 'lesson_enter'
  | 'note_saved'
  | 'lesson_complete'
  | 'level_complete'
  | 'course_complete'
  | 'base_built'
  | 'manual_open';

export type LumiQuickAction = {
  id: string;
  label: string;
  action: () => void;
};

export type LumiQuickActionTemplate = {
  id: string;
  label: string;
};

export interface LumiViewState {
  state: LumiState;
  context: LumiPageContext;
  dockSlot: LumiDockSlot;
  visible: boolean;
  expanded: boolean;
  message: string;
  messageType: LumiMessageType;
  quickActions?: LumiQuickAction[];
  runtimeContext?: LumiRuntimeContext;
}

export interface LumiTriggerPayload {
  context?: LumiPageContext;
  dockSlot?: LumiDockSlot;
  courseTitle?: string;
  lessonTitle?: string;
  status?: 'draft' | 'ready' | 'learning' | 'completed';
  message?: string;
  preview?: string;
  actionHandlers?: Record<string, () => void>;
}

export interface LumiTriggerDefinition {
  priority: number;
  state: LumiState;
  context?: LumiPageContext;
  dockSlot?: LumiDockSlot;
  messageType: LumiMessageType;
  getMessage: (payload?: LumiTriggerPayload) => string;
}
