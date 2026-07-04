import { fromTriggerToLumiRuntimeContext } from './lumiContextAdapters';
import { resolveLumiRuntime } from './lumiRuntimeEngine';
import type { LumiQuickAction, LumiTrigger, LumiTriggerPayload, LumiViewState } from './lumiTypes';

const lumiTriggerPriorityMap: Record<LumiTrigger, number> = {
  page_enter: 10,
  cta_focus: 70,
  cta_submit: 65,
  generation_loading: 60,
  generation_success: 80,
  planet_hover: 30,
  planet_select: 40,
  idle: 20,
  search_start: 55,
  search_success: 61,
  evaluation_start: 50,
  evaluation_success: 62,
  lesson_enter: 45,
  note_saved: 63,
  lesson_complete: 75,
  level_complete: 76,
  course_complete: 90,
  base_built: 100,
  manual_open: 15,
};

function toQuickActions(payload: LumiTriggerPayload | undefined, templates: Array<{ id: string; label: string }>): LumiQuickAction[] {
  const handlers = payload?.actionHandlers ?? {};
  return templates.map((template) => ({
    ...template,
    action: handlers[template.id] ?? (() => {}),
  }));
}

export function createLumiViewState(
  trigger: LumiTrigger,
  payload?: LumiTriggerPayload,
  options?: { isMobile?: boolean; expanded?: boolean },
): LumiViewState {
  const runtimeContext = fromTriggerToLumiRuntimeContext(trigger, payload);
  const resolved = resolveLumiRuntime({
    kind: 'trigger',
    runtimeContext,
    trigger,
    payload,
    isMobile: options?.isMobile ?? false,
    expanded: options?.expanded ?? false,
  });

  return {
    state: resolved.state,
    context: resolved.context,
    dockSlot: resolved.dockSlot,
    visible: resolved.visible,
    expanded: resolved.expanded,
    message: payload?.message ?? resolved.message,
    messageType: resolved.messageType,
    quickActions: toQuickActions(payload, resolved.quickActionTemplates),
    runtimeContext: resolved.runtimeContext,
  };
}

export function getLumiTriggerPriority(trigger: LumiTrigger): number {
  return lumiTriggerPriorityMap[trigger];
}
