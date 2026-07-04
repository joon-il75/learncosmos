import { resolveLumiRuntime } from './lumiRuntimeEngine';
import { lumiQuickActionTemplates } from './lumiQuickActions';
import type { LumiRuntimeContext } from './lumiEngineTypes';
import type { LumiQuickAction, LumiTrigger, LumiTriggerPayload, LumiViewState } from './lumiTypes';

interface RuntimeViewOptions {
  runtimeContext: LumiRuntimeContext;
  trigger: LumiTrigger;
  payload?: LumiTriggerPayload;
  isMobile?: boolean;
  expanded?: boolean;
}

function toQuickActions(payload: LumiTriggerPayload | undefined, templates: Array<{ id: string; label: string }>): LumiQuickAction[] {
  const handlers = payload?.actionHandlers ?? {};
  return templates.map((template) => ({
    ...template,
    action: handlers[template.id] ?? (() => {}),
  }));
}

export function createLumiViewStateFromRuntime({
  runtimeContext,
  trigger,
  payload,
  isMobile = false,
  expanded = false,
}: RuntimeViewOptions): LumiViewState {
  const resolved = resolveLumiRuntime({
    kind: 'trigger',
    runtimeContext,
    trigger,
    payload,
    isMobile,
    expanded,
  });

  const contextTemplates = lumiQuickActionTemplates[resolved.context] ?? [];

  return {
    state: resolved.state,
    context: resolved.context,
    dockSlot: resolved.dockSlot,
    visible: resolved.visible,
    expanded: resolved.expanded,
    message: payload?.message ?? resolved.message,
    messageType: resolved.messageType,
    quickActions: toQuickActions(payload, contextTemplates),
    runtimeContext: resolved.runtimeContext,
  };
}
