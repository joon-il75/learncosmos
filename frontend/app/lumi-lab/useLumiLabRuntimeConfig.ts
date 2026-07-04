'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { createLumiViewStateFromRuntime } from '@/lib/lumi/lumiRuntimeView';
import {
  type LumiLabTab,
  type LumiRuntimeConfigResponse,
  type LumiRuleDraft,
  type LumiMessageDraft,
  normalizeRuntimeConfigResponse,
  serializeRuntimeConfigDrafts,
  createBlankRuleDraft,
  createBlankMessageDraft,
  createInitialLumiRuleDrafts,
  getDefaultLumiRuleDraftId,
  createInitialLumiMessageDrafts,
  getDefaultLumiMessageDraftId,
} from './lumiLabTypes';

export type LumiLabRuntimeConfigResult = {
  ruleDrafts: LumiRuleDraft[];
  selectedRuleId: string;
  setSelectedRuleId: (id: string) => void;
  messageDrafts: LumiMessageDraft[];
  selectedMessageDraftId: string;
  setSelectedMessageDraftId: (id: string) => void;
  runtimeConfigStatus: string | null;
  setRuntimeConfigStatus: (status: string | null) => void;
  runtimeConfigVersion: number | null;
  isSavingRuntimeConfig: boolean;
  ruleSearchQuery: string;
  setRuleSearchQuery: (q: string) => void;
  ruleSortMode: 'label' | 'runtime' | 'state';
  setRuleSortMode: (m: 'label' | 'runtime' | 'state') => void;
  messageSearchQuery: string;
  setMessageSearchQuery: (q: string) => void;
  messageSortMode: 'label' | 'trigger' | 'type';
  setMessageSortMode: (m: 'label' | 'trigger' | 'type') => void;
  selectedRuleDraft: LumiRuleDraft | null;
  selectedMessageDraft: LumiMessageDraft | null;
  filteredRuleDrafts: LumiRuleDraft[];
  filteredMessageDrafts: LumiMessageDraft[];
  hasUnsavedRuntimeChanges: boolean;
  selectedRuleDiffLines: string[];
  selectedMessageDiffLines: string[];
  selectedRuleEngineBaseline: ReturnType<typeof createLumiViewStateFromRuntime> | null;
  loadRuntimeConfig: () => Promise<void>;
  saveRuntimeConfig: () => Promise<void>;
  updateSelectedRuleDraft: (next: Partial<LumiRuleDraft>) => void;
  resetSelectedRuleDraft: () => void;
  updateSelectedMessageDraft: (next: Partial<LumiMessageDraft>) => void;
  resetSelectedMessageDraft: () => void;
  addRuleDraft: () => void;
  deleteSelectedRuleDraft: () => void;
  addMessageDraft: () => void;
  deleteSelectedMessageDraft: () => void;
};

export function useLumiLabRuntimeConfig(params: {
  onNavigateToTab: (tab: LumiLabTab) => void;
}): LumiLabRuntimeConfigResult {
  const router = useRouter();
  const [ruleDrafts, setRuleDrafts] = useState<LumiRuleDraft[]>(() => createInitialLumiRuleDrafts());
  const [selectedRuleId, setSelectedRuleId] = useState(() => getDefaultLumiRuleDraftId());
  const [messageDrafts, setMessageDrafts] = useState<LumiMessageDraft[]>(() => createInitialLumiMessageDrafts());
  const [selectedMessageDraftId, setSelectedMessageDraftId] = useState(() => getDefaultLumiMessageDraftId());
  const [runtimeConfigStatus, setRuntimeConfigStatus] = useState<string | null>(null);
  const [runtimeConfigVersion, setRuntimeConfigVersion] = useState<number | null>(null);
  const [isSavingRuntimeConfig, setIsSavingRuntimeConfig] = useState(false);
  const [ruleSearchQuery, setRuleSearchQuery] = useState('');
  const [ruleSortMode, setRuleSortMode] = useState<'label' | 'runtime' | 'state'>('runtime');
  const [messageSearchQuery, setMessageSearchQuery] = useState('');
  const [messageSortMode, setMessageSortMode] = useState<'label' | 'trigger' | 'type'>('trigger');
  const lastSavedRuntimeConfigRef = useRef(
    serializeRuntimeConfigDrafts(createInitialLumiRuleDrafts(), createInitialLumiMessageDrafts()),
  );
  const onNavigateToTabRef = useRef(params.onNavigateToTab);
  onNavigateToTabRef.current = params.onNavigateToTab;

  const selectedRuleDraft = ruleDrafts.find((rule) => rule.id === selectedRuleId) ?? ruleDrafts[0] ?? null;
  const selectedMessageDraft = messageDrafts.find((draft) => draft.id === selectedMessageDraftId) ?? messageDrafts[0] ?? null;

  const hasUnsavedRuntimeChanges = useMemo(() => (
    serializeRuntimeConfigDrafts(ruleDrafts, messageDrafts) !== lastSavedRuntimeConfigRef.current
  ), [messageDrafts, ruleDrafts]);

  const savedRuntimeConfig = useMemo(() => {
    try {
      return JSON.parse(lastSavedRuntimeConfigRef.current) as { rules: LumiRuleDraft[]; messages: LumiMessageDraft[] };
    } catch {
      return { rules: [], messages: [] };
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [hasUnsavedRuntimeChanges, runtimeConfigVersion]);

  const filteredRuleDrafts = useMemo(() => {
    const query = ruleSearchQuery.trim().toLowerCase();
    const filtered = ruleDrafts.filter((rule) => {
      if (!query) return true;
      return [rule.label, rule.page, rule.action, rule.scene, rule.trigger, rule.state, rule.messageType, rule.message]
        .join(' ').toLowerCase().includes(query);
    });
    const sorted = [...filtered];
    sorted.sort((l, r) => {
      if (ruleSortMode === 'label') return l.label.localeCompare(r.label);
      if (ruleSortMode === 'state') return `${l.state} ${l.messageType}`.localeCompare(`${r.state} ${r.messageType}`);
      return `${l.page} ${l.action} ${l.scene}`.localeCompare(`${r.page} ${r.action} ${r.scene}`);
    });
    return sorted;
  }, [ruleDrafts, ruleSearchQuery, ruleSortMode]);

  const filteredMessageDrafts = useMemo(() => {
    const query = messageSearchQuery.trim().toLowerCase();
    const filtered = messageDrafts.filter((draft) => {
      if (!query) return true;
      return [draft.label, draft.trigger, draft.messageType, draft.message, draft.note]
        .join(' ').toLowerCase().includes(query);
    });
    const sorted = [...filtered];
    sorted.sort((l, r) => {
      if (messageSortMode === 'label') return l.label.localeCompare(r.label);
      if (messageSortMode === 'type') return `${l.messageType} ${l.label}`.localeCompare(`${r.messageType} ${r.label}`);
      return `${l.trigger} ${l.label}`.localeCompare(`${r.trigger} ${r.label}`);
    });
    return sorted;
  }, [messageDrafts, messageSearchQuery, messageSortMode]);

  const savedSelectedRuleDraft = useMemo(() => (
    savedRuntimeConfig.rules.find((rule) => rule.id === selectedRuleId) ?? null
  ), [savedRuntimeConfig.rules, selectedRuleId]);

  const savedSelectedMessageDraft = useMemo(() => (
    savedRuntimeConfig.messages.find((draft) => draft.id === selectedMessageDraftId) ?? null
  ), [savedRuntimeConfig.messages, selectedMessageDraftId]);

  const selectedRuleDiffLines = useMemo(() => {
    if (!selectedRuleDraft) return [];
    const saved = savedSelectedRuleDraft;
    const lines: string[] = [];
    if (!saved) { lines.push('새 rule draft: 아직 저장되지 않았습니다.'); return lines; }
    if (saved.label !== selectedRuleDraft.label) lines.push(`label: ${saved.label} -> ${selectedRuleDraft.label}`);
    if (saved.trigger !== selectedRuleDraft.trigger) lines.push(`trigger: ${saved.trigger} -> ${selectedRuleDraft.trigger}`);
    if (`${saved.page}/${saved.action}/${saved.scene}` !== `${selectedRuleDraft.page}/${selectedRuleDraft.action}/${selectedRuleDraft.scene}`) lines.push(`runtime: ${saved.page}/${saved.action}/${saved.scene} -> ${selectedRuleDraft.page}/${selectedRuleDraft.action}/${selectedRuleDraft.scene}`);
    if (`${saved.context}/${saved.dockSlot}` !== `${selectedRuleDraft.context}/${selectedRuleDraft.dockSlot}`) lines.push(`legacy: ${saved.context}/${saved.dockSlot} -> ${selectedRuleDraft.context}/${selectedRuleDraft.dockSlot}`);
    if (`${saved.state}/${saved.messageType}` !== `${selectedRuleDraft.state}/${selectedRuleDraft.messageType}`) lines.push(`output: ${saved.state}/${saved.messageType} -> ${selectedRuleDraft.state}/${selectedRuleDraft.messageType}`);
    if (saved.message !== selectedRuleDraft.message) lines.push('message: 저장본과 현재 편집본이 다릅니다.');
    if (saved.note !== selectedRuleDraft.note) lines.push('note: 저장본과 현재 편집본이 다릅니다.');
    return lines;
  }, [savedSelectedRuleDraft, selectedRuleDraft]);

  const selectedMessageDiffLines = useMemo(() => {
    if (!selectedMessageDraft) return [];
    const saved = savedSelectedMessageDraft;
    const lines: string[] = [];
    if (!saved) { lines.push('새 message draft: 아직 저장되지 않았습니다.'); return lines; }
    if (saved.label !== selectedMessageDraft.label) lines.push(`label: ${saved.label} -> ${selectedMessageDraft.label}`);
    if (saved.trigger !== selectedMessageDraft.trigger) lines.push(`trigger: ${saved.trigger} -> ${selectedMessageDraft.trigger}`);
    if (saved.messageType !== selectedMessageDraft.messageType) lines.push(`type: ${saved.messageType} -> ${selectedMessageDraft.messageType}`);
    if (saved.preview !== selectedMessageDraft.preview) lines.push('preview: 저장본과 현재 편집본이 다릅니다.');
    if (saved.message !== selectedMessageDraft.message) lines.push('message: 저장본과 현재 편집본이 다릅니다.');
    if (saved.note !== selectedMessageDraft.note) lines.push('note: 저장본과 현재 편집본이 다릅니다.');
    return lines;
  }, [savedSelectedMessageDraft, selectedMessageDraft]);

  const selectedRuleEngineBaseline = useMemo(() => {
    if (!selectedRuleDraft) return null;
    return createLumiViewStateFromRuntime({
      runtimeContext: {
        mode: selectedRuleDraft.mode,
        page: selectedRuleDraft.page,
        action: selectedRuleDraft.action,
        scene: selectedRuleDraft.scene,
      },
      trigger: selectedRuleDraft.trigger,
      payload: { context: selectedRuleDraft.context, dockSlot: selectedRuleDraft.dockSlot },
    });
  }, [selectedRuleDraft]);

  useEffect(() => {
    if (!hasUnsavedRuntimeChanges) return;
    const handleBeforeUnload = (event: BeforeUnloadEvent) => {
      event.preventDefault();
      event.returnValue = '';
    };
    window.addEventListener('beforeunload', handleBeforeUnload);
    return () => window.removeEventListener('beforeunload', handleBeforeUnload);
  }, [hasUnsavedRuntimeChanges]);

  const loadRuntimeConfig = useCallback(async () => {
    const res = await fetch('/api/v1/super-admin/settings/lumi-runtime', {
      credentials: 'include',
      cache: 'no-store',
    });
    if (res.status === 401 || res.status === 403) {
      router.replace('/super-admin/login');
      return;
    }
    if (!res.ok) throw new Error('Lumi runtime 설정을 불러오지 못했습니다.');
    const payload = normalizeRuntimeConfigResponse(await res.json() as LumiRuntimeConfigResponse);
    setRuleDrafts(payload.rules);
    setMessageDrafts(payload.messages);
    setSelectedRuleId(payload.rules[0]?.id ?? getDefaultLumiRuleDraftId());
    setSelectedMessageDraftId(payload.messages[0]?.id ?? getDefaultLumiMessageDraftId());
    lastSavedRuntimeConfigRef.current = serializeRuntimeConfigDrafts(payload.rules, payload.messages);
    setRuntimeConfigVersion(payload.version ?? null);
    setRuntimeConfigStatus(payload.exists ? '저장된 Lumi runtime 설정을 불러왔습니다.' : '저장된 Lumi runtime 설정이 없어 기본값을 사용합니다.');
  }, [router]);

  const saveRuntimeConfig = async () => {
    setIsSavingRuntimeConfig(true);
    setRuntimeConfigStatus(null);
    try {
      const res = await fetch('/api/v1/super-admin/settings/lumi-runtime', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ rules: ruleDrafts, messages: messageDrafts }),
      });
      if (res.status === 401 || res.status === 403) { router.replace('/super-admin/login'); return; }
      if (!res.ok) throw new Error('Lumi runtime 설정을 저장하지 못했습니다.');
      const payload = await res.json() as { config?: LumiRuntimeConfigResponse };
      const normalized = normalizeRuntimeConfigResponse(
        payload.config ?? { exists: true, rules: ruleDrafts, messages: messageDrafts },
      );
      setRuleDrafts(normalized.rules);
      setMessageDrafts(normalized.messages);
      setSelectedRuleId((current) => normalized.rules.find((rule) => rule.id === current)?.id ?? normalized.rules[0]?.id ?? '');
      setSelectedMessageDraftId((current) => normalized.messages.find((draft) => draft.id === current)?.id ?? normalized.messages[0]?.id ?? '');
      lastSavedRuntimeConfigRef.current = serializeRuntimeConfigDrafts(normalized.rules, normalized.messages);
      setRuntimeConfigVersion(normalized.version ?? Date.now());
      setRuntimeConfigStatus('Lumi runtime 설정을 저장했습니다.');
    } catch (error) {
      setRuntimeConfigStatus(error instanceof Error ? error.message : 'Lumi runtime 설정을 저장하지 못했습니다.');
    } finally {
      setIsSavingRuntimeConfig(false);
    }
  };

  const updateSelectedRuleDraft = (next: Partial<LumiRuleDraft>) => {
    if (!selectedRuleDraft) return;
    setRuleDrafts((current) => current.map((rule) => rule.id === selectedRuleDraft.id ? { ...rule, ...next } : rule));
  };

  const resetSelectedRuleDraft = () => {
    const baseline = createInitialLumiRuleDrafts().find((rule) => rule.id === selectedRuleId);
    if (!baseline) return;
    setRuleDrafts((current) => current.map((rule) => rule.id === selectedRuleId ? baseline : rule));
  };

  const updateSelectedMessageDraft = (next: Partial<LumiMessageDraft>) => {
    if (!selectedMessageDraft) return;
    setMessageDrafts((current) => current.map((draft) => draft.id === selectedMessageDraft.id ? { ...draft, ...next } : draft));
  };

  const resetSelectedMessageDraft = () => {
    const baseline = createInitialLumiMessageDrafts().find((draft) => draft.id === selectedMessageDraftId);
    if (!baseline) return;
    setMessageDrafts((current) => current.map((draft) => draft.id === selectedMessageDraftId ? baseline : draft));
  };

  const addRuleDraft = () => {
    const nextDraft = createBlankRuleDraft();
    setRuleDrafts((current) => [...current, nextDraft]);
    setSelectedRuleId(nextDraft.id);
    onNavigateToTabRef.current('rules');
    setRuntimeConfigStatus('새 rule draft를 추가했습니다. 저장 전까지는 editor 상태로만 유지됩니다.');
  };

  const deleteSelectedRuleDraft = () => {
    if (!selectedRuleDraft || ruleDrafts.length <= 1) return;
    if (!window.confirm(`'${selectedRuleDraft.label}' rule draft를 삭제할까요? 저장 전까지는 되돌릴 수 있습니다.`)) return;
    const nextRules = ruleDrafts.filter((rule) => rule.id !== selectedRuleDraft.id);
    setRuleDrafts(nextRules);
    setSelectedRuleId(nextRules[0]?.id ?? '');
    setRuntimeConfigStatus('선택한 rule draft를 삭제했습니다.');
  };

  const addMessageDraft = () => {
    const nextDraft = createBlankMessageDraft();
    setMessageDrafts((current) => [...current, nextDraft]);
    setSelectedMessageDraftId(nextDraft.id);
    onNavigateToTabRef.current('messages');
    setRuntimeConfigStatus('새 message draft를 추가했습니다. 저장 전까지는 editor 상태로만 유지됩니다.');
  };

  const deleteSelectedMessageDraft = () => {
    if (!selectedMessageDraft || messageDrafts.length <= 1) return;
    if (!window.confirm(`'${selectedMessageDraft.label}' message draft를 삭제할까요? 저장 전까지는 되돌릴 수 있습니다.`)) return;
    const nextMessages = messageDrafts.filter((draft) => draft.id !== selectedMessageDraft.id);
    setMessageDrafts(nextMessages);
    setSelectedMessageDraftId(nextMessages[0]?.id ?? '');
    setRuntimeConfigStatus('선택한 message draft를 삭제했습니다.');
  };

  return {
    ruleDrafts, selectedRuleId, setSelectedRuleId,
    messageDrafts, selectedMessageDraftId, setSelectedMessageDraftId,
    runtimeConfigStatus, setRuntimeConfigStatus, runtimeConfigVersion, isSavingRuntimeConfig,
    ruleSearchQuery, setRuleSearchQuery, ruleSortMode, setRuleSortMode,
    messageSearchQuery, setMessageSearchQuery, messageSortMode, setMessageSortMode,
    selectedRuleDraft, selectedMessageDraft,
    filteredRuleDrafts, filteredMessageDrafts,
    hasUnsavedRuntimeChanges, selectedRuleDiffLines, selectedMessageDiffLines, selectedRuleEngineBaseline,
    loadRuntimeConfig, saveRuntimeConfig,
    updateSelectedRuleDraft, resetSelectedRuleDraft,
    updateSelectedMessageDraft, resetSelectedMessageDraft,
    addRuleDraft, deleteSelectedRuleDraft,
    addMessageDraft, deleteSelectedMessageDraft,
  };
}
