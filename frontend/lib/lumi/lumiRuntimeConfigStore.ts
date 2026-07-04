import { createInitialLumiMessageDrafts, setLumiMessageDraftOverrides, type LumiMessageDraft } from './lumiMessageTemplates';
import { createInitialLumiRuleDrafts, setLumiRuleDraftOverrides, type LumiRuleDraft } from './lumiRuleDrafts';

export interface LumiPersistedRuntimeConfig {
  exists: boolean;
  rules: LumiRuleDraft[];
  messages: LumiMessageDraft[];
  updated_at?: string;
  version?: number;
}

type LumiPersistedRuntimeConfigInput = Partial<LumiPersistedRuntimeConfig>;

type LumiRuntimeConfigState = {
  version: number;
  hydrated: boolean;
  hydrationPromise: Promise<void> | null;
};

const listeners = new Set<() => void>();
const state: LumiRuntimeConfigState = {
  version: 0,
  hydrated: false,
  hydrationPromise: null,
};

function cloneRules(rules: LumiRuleDraft[]): LumiRuleDraft[] {
  return rules.map((rule) => ({ ...rule }));
}

function cloneMessages(messages: LumiMessageDraft[]): LumiMessageDraft[] {
  return messages.map((message) => ({ ...message }));
}

function normalizePersistedConfig(input?: LumiPersistedRuntimeConfigInput | null): LumiPersistedRuntimeConfig {
  if (!input || input.exists === false) {
    return {
      exists: false,
      rules: createInitialLumiRuleDrafts(),
      messages: createInitialLumiMessageDrafts(),
    };
  }

  return {
    exists: true,
    rules: Array.isArray(input.rules) ? cloneRules(input.rules) : createInitialLumiRuleDrafts(),
    messages: Array.isArray(input.messages) ? cloneMessages(input.messages) : createInitialLumiMessageDrafts(),
    updated_at: input.updated_at,
    version: input.version,
  };
}

function notifyRuntimeConfigListeners() {
  listeners.forEach((listener) => listener());
}

export function subscribeLumiRuntimeConfig(listener: () => void): () => void {
  listeners.add(listener);
  return () => listeners.delete(listener);
}

export function getLumiRuntimeConfigVersion(): number {
  return state.version;
}

export function applyLumiRuntimeConfig(input?: LumiPersistedRuntimeConfigInput | null): LumiPersistedRuntimeConfig {
  const config = normalizePersistedConfig(input);

  setLumiRuleDraftOverrides(config.rules);
  setLumiMessageDraftOverrides(config.messages);

  state.version = config.version ?? Date.now();
  state.hydrated = true;
  notifyRuntimeConfigListeners();

  return config;
}

export async function loadPublicLumiRuntimeConfig(force = false): Promise<LumiPersistedRuntimeConfig | null> {
  if (state.hydrated && !force) {
    return null;
  }

  if (state.hydrationPromise && !force) {
    await state.hydrationPromise;
    return null;
  }

  state.hydrationPromise = (async () => {
    try {
      const res = await fetch('/api/v1/public/lumi/runtime-config', { cache: 'no-store' });
      if (!res.ok) {
        applyLumiRuntimeConfig({ exists: false });
        return;
      }

      const payload = await res.json() as LumiPersistedRuntimeConfig;
      applyLumiRuntimeConfig(payload);
    } catch {
      applyLumiRuntimeConfig({ exists: false });
    } finally {
      state.hydrationPromise = null;
    }
  })();

  await state.hydrationPromise;
  return null;
}
