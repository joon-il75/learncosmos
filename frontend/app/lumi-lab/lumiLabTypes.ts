import { applyLumiRuntimeConfig, type LumiPersistedRuntimeConfig } from '@/lib/lumi/lumiRuntimeConfigStore';
import { createInitialLumiRuleDrafts, getDefaultLumiRuleDraftId, type LumiRuleDraft } from '@/lib/lumi/lumiRuleDrafts';
import { createInitialLumiMessageDrafts, getDefaultLumiMessageDraftId, type LumiMessageDraft } from '@/lib/lumi/lumiMessageTemplates';
import type { LumiAction, LumiMode, LumiPage, LumiScene } from '@/lib/lumi/lumiEngineTypes';
import type { LumiPageContext, LumiDockSlot, LumiMessageType, LumiState, LumiTrigger } from '@/lib/lumi/lumiTypes';

export type { LumiRuleDraft, LumiMessageDraft };

export const stateOptions: LumiState[] = [
  'idle',
  'thinking',
  'happy',
  'celebrate',
  'victory',
  'encourage',
  'surprise',
  'curious',
  'raise-hand',
  'sleepy',
  'exploring',
  'focus',
  'planet-point',
  'planet-hold',
  'note-read',
];

export const contextOptions: LumiPageContext[] = [
  'dashboard',
  'planet-map',
  'planet-scene',
  'diary',
  'player',
  'completion',
  'mobile',
];

export const dockSlotOptions: LumiDockSlot[] = [
  'dashboard-cta',
  'dashboard-map-top',
  'dashboard-map-bottom',
  'dashboard-panel',
  'diary-header',
  'diary-sidebar',
  'player-inline',
  'player-bottom',
  'planet-scene-anchor',
  'completion-center',
  'mobile-bottom-sheet-trigger',
  'mobile-inline',
];

export const messageTypeOptions: LumiMessageType[] = [
  'summary',
  'hint',
  'question',
  'encourage',
  'success',
  'guide',
  'discovery',
];

export const triggerOptions: LumiTrigger[] = [
  'page_enter',
  'cta_focus',
  'generation_loading',
  'generation_success',
  'planet_hover',
  'planet_select',
  'idle',
  'search_start',
  'search_success',
  'evaluation_start',
  'evaluation_success',
  'lesson_enter',
  'note_saved',
  'lesson_complete',
  'course_complete',
  'base_built',
  'manual_open',
];

export type LumiLabTab = 'sprites' | 'scenarios' | 'rules' | 'messages' | 'ui';

export const runtimeModeOptions: LumiMode[] = ['general', 'ai'];
export const runtimePageOptions: LumiPage[] = ['hero', 'dashboard', 'planet_map', 'course_draft_editor', 'course_player'];
export const runtimeActionOptions: LumiAction[] = [
  'page_enter',
  'open_lumi',
  'goal_submit',
  'create_course',
  'search_content',
  'ai_generate',
  'ai_response',
  'lesson_complete',
  'course_complete',
  'error',
];
export const runtimeSceneOptions: LumiScene[] = ['empty', 'overview', 'editing', 'searching', 'waiting', 'success', 'error'];

export interface LumiAssetMeta {
  path: string;
  public_url: string;
  exists: boolean;
  size_bytes: number;
  updated_at?: string;
  version?: number;
}

export interface LumiRuntimeConfigResponse extends LumiPersistedRuntimeConfig {}

export function normalizeRuntimeConfigResponse(payload?: Partial<LumiRuntimeConfigResponse> | null): LumiRuntimeConfigResponse {
  const normalized = applyLumiRuntimeConfig(payload ?? { exists: false });
  return {
    ...normalized,
    updated_at: payload?.updated_at,
    version: payload?.version,
  };
}

export function serializeRuntimeConfigDrafts(rules: LumiRuleDraft[], messages: LumiMessageDraft[]): string {
  return JSON.stringify({ rules, messages });
}

export function createRuleDraftId(): string {
  return `rule-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;
}

export function createMessageDraftId(): string {
  return `message-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;
}

export function createBlankRuleDraft(): LumiRuleDraft {
  return {
    id: createRuleDraftId(),
    label: 'New Rule',
    trigger: 'page_enter',
    mode: 'general',
    page: 'dashboard',
    action: 'page_enter',
    scene: 'overview',
    context: 'dashboard',
    dockSlot: 'dashboard-map-top',
    state: 'idle',
    messageType: 'summary',
    message: '새 Lumi rule을 편집해 주세요.',
    note: '',
  };
}

export function createBlankMessageDraft(): LumiMessageDraft {
  return {
    id: createMessageDraftId(),
    trigger: 'page_enter',
    label: 'New Message',
    note: '',
    messageType: 'summary',
    message: '새 Lumi message를 편집해 주세요.',
    preview: '',
  };
}

export { createInitialLumiRuleDrafts, getDefaultLumiRuleDraftId, createInitialLumiMessageDrafts, getDefaultLumiMessageDraftId };
