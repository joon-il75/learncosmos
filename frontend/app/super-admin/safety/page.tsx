'use client';

import type { CSSProperties } from 'react';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { useRouter } from 'next/navigation';
import SuperAdminPanelHeader from '@/components/super-admin/SuperAdminPanelHeader';
import SuperAdminPanelNav from '@/components/super-admin/SuperAdminPanelNav';
import { superAdminShellStyle } from '@/components/super-admin/layout';

const PAGE_SIZE = 50;

type ModerationAction = 'allow' | 'soft_warn' | 'block';
type RuleAction = 'soft_warn' | 'block';
type RuleType = 'keyword' | 'regex';
type SafetyTab = 'logs' | 'rules';
type AIOutputSeverity = 'none' | 'low' | 'caution' | 'warning' | 'critical';
type AIOutputReviewStatus = 'unreviewed' | 'false_positive' | 'needs_prompt_guard' | 'needs_rule_tuning' | 'needs_masking' | 'needs_regeneration' | 'resolved';

type ModerationLog = {
  id: string;
  user_id?: string;
  user_email?: string;
  target_type: string;
  target_id?: string;
  input_text_hash: string;
  risk_type: string;
  action: ModerationAction;
  matched_rule_id?: string;
  route: string;
  locale?: string;
  metadata?: Record<string, unknown>;
  created_at: string;
  matched_pattern?: string;
  rule_type?: string;
};

type ModerationRule = {
  id: string;
  rule_type: RuleType;
  pattern: string;
  risk_type: string;
  action: RuleAction;
  locale?: string | null;
  description: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

type AIOutputActionCounts = {
  allow: number;
  soft_warn: number;
  block: number;
};

type AIOutputDimensionCount = {
  key: string;
  count: number;
};

type AIOutputRuleCount = {
  rule_id?: string;
  pattern?: string;
  count: number;
};

type AIOutputTargetSummary = {
  target_type: string;
  total: number;
  risky_total: number;
  risky_rate: number;
  severity: AIOutputSeverity;
  action_counts: AIOutputActionCounts;
};

type AIOutputSummaryWindow = {
  label: string;
  hours: number;
  since: string;
  total: number;
  risky_total: number;
  risky_rate: number;
  severity: AIOutputSeverity;
  message: string;
  action_counts: AIOutputActionCounts;
  targets: AIOutputTargetSummary[];
  risk_types: AIOutputDimensionCount[];
  source_features: AIOutputDimensionCount[];
  matched_rules: AIOutputRuleCount[];
};

type AIOutputSummary = {
  generated_at: string;
  windows: AIOutputSummaryWindow[];
};

type AIOutputReviewDecision = {
  candidate_key: string;
  status: AIOutputReviewStatus;
  note: string;
  reviewed_by: string;
  reviewed_at?: string;
  metadata?: Record<string, unknown>;
  created_at?: string;
  updated_at?: string;
};

type AIOutputReviewCandidate = {
  key: string;
  target_type: string;
  risk_type: string;
  source_feature: string;
  matched_rule_id?: string;
  matched_pattern?: string;
  action_counts: AIOutputActionCounts;
  total: number;
  latest_seen_at: string;
  sample_target_id?: string;
  review: AIOutputReviewDecision;
};

type AIOutputReviewCandidates = {
  generated_at: string;
  window_hours: number;
  since: string;
  total_candidates: number;
  candidates: AIOutputReviewCandidate[];
};

type SelectOption = { value: string; label: string };

type RuleForm = {
  rule_type: RuleType;
  pattern: string;
  risk_type: string;
  action: RuleAction;
  locale: string;
  description: string;
};

const emptyRuleForm: RuleForm = {
  rule_type: 'keyword',
  pattern: '',
  risk_type: 'unsafe_instruction',
  action: 'block',
  locale: '',
  description: '',
};

const actionOptions = [
  { value: '', label: 'Action 전체' },
  { value: 'block', label: 'Block' },
  { value: 'soft_warn', label: 'Soft warn' },
  { value: 'allow', label: 'Allow' },
] as const;

const ruleActionOptions = [
  { value: '', label: 'Action 전체' },
  { value: 'block', label: 'Block' },
  { value: 'soft_warn', label: 'Soft warn' },
] as const;

const reviewRubricItems = [
  { status: 'false_positive', description: '정상 맥락 감지' },
  { status: 'needs_prompt_guard', description: '기능 prompt 보강' },
  { status: 'needs_rule_tuning', description: 'rule 범위 조정' },
  { status: 'needs_masking', description: '일부 표현 가림' },
  { status: 'needs_regeneration', description: '출력 재생성 필요' },
  { status: 'resolved', description: '조치 완료' },
] as const;

const aiOutputReviewStatusOptions = [
  { value: 'unreviewed', label: '미검토' },
  { value: 'false_positive', label: 'False positive' },
  { value: 'needs_prompt_guard', label: 'Prompt guard 필요' },
  { value: 'needs_rule_tuning', label: 'Rule tuning 필요' },
  { value: 'needs_masking', label: 'Masking 필요' },
  { value: 'needs_regeneration', label: 'Regeneration 필요' },
  { value: 'resolved', label: 'Resolved' },
] as const;

const ruleFormActionOptions = [
  { value: 'block', label: 'block' },
  { value: 'soft_warn', label: 'soft_warn' },
] as const;

const ruleTypeOptions = [
  { value: 'keyword', label: 'keyword' },
  { value: 'regex', label: 'regex' },
] as const;

const ruleTypeFilterOptions = [
  { value: '', label: 'type 전체' },
  ...ruleTypeOptions,
] as const;

const ruleActiveOptions = [
  { value: '', label: '활성 전체' },
  { value: 'true', label: 'active' },
  { value: 'false', label: 'inactive' },
] as const;

const targetOptions = [
  { value: '', label: 'Target 전체' },
  { value: 'goal_start', label: 'goal_start' },
  { value: 'goal_interview', label: 'goal_interview' },
  { value: 'course_draft_source_query', label: 'course_draft_source_query' },
  { value: 'point_observation_note', label: 'point_observation_note' },
  { value: 'point_journal', label: 'point_journal' },
  { value: 'point_goal', label: 'point_goal' },
  { value: 'point_question', label: 'point_question' },
  { value: 'point_practice_log', label: 'point_practice_log' },
  { value: 'point_artifact', label: 'point_artifact' },
  { value: 'point_attachment', label: 'point_attachment' },
  { value: 'point_self_evaluation', label: 'point_self_evaluation' },
  { value: 'ai_goal_interview_output', label: 'ai_goal_interview_output' },
  { value: 'ai_course_draft_output', label: 'ai_course_draft_output' },
  { value: 'ai_point_summary_output', label: 'ai_point_summary_output' },
  { value: 'ai_point_question_output', label: 'ai_point_question_output' },
  { value: 'ai_point_feedback_output', label: 'ai_point_feedback_output' },
  { value: 'ai_point_self_evaluation_output', label: 'ai_point_self_evaluation_output' },
] as const;

export default function SuperAdminSafetyPage() {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState<SafetyTab>('logs');
  const [logs, setLogs] = useState<ModerationLog[]>([]);
  const [rules, setRules] = useState<ModerationRule[]>([]);
  const [aiOutputSummary, setAIOutputSummary] = useState<AIOutputSummary | null>(null);
  const [aiOutputCandidates, setAIOutputCandidates] = useState<AIOutputReviewCandidates | null>(null);
  const [action, setAction] = useState('block');
  const [targetType, setTargetType] = useState('');
  const [riskType, setRiskType] = useState('');
  const [offset, setOffset] = useState(0);
  const [ruleActiveFilter, setRuleActiveFilter] = useState('');
  const [ruleActionFilter, setRuleActionFilter] = useState('');
  const [ruleTypeFilter, setRuleTypeFilter] = useState('');
  const [ruleQuery, setRuleQuery] = useState('');
  const [ruleOffset, setRuleOffset] = useState(0);
  const [ruleForm, setRuleForm] = useState<RuleForm>(emptyRuleForm);
  const [editingRuleID, setEditingRuleID] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [aiOutputSummaryLoading, setAIOutputSummaryLoading] = useState(true);
  const [aiOutputCandidatesLoading, setAIOutputCandidatesLoading] = useState(true);
  const [rulesLoading, setRulesLoading] = useState(false);
  const [savingRule, setSavingRule] = useState(false);
  const [aiOutputReviewSavingKey, setAIOutputReviewSavingKey] = useState<string | null>(null);
  const [error, setError] = useState('');
  const [aiOutputSummaryError, setAIOutputSummaryError] = useState('');
  const [aiOutputCandidatesError, setAIOutputCandidatesError] = useState('');
  const [rulesError, setRulesError] = useState('');
  const [rulesMessage, setRulesMessage] = useState('');
  const [aiOutputCandidatesMessage, setAIOutputCandidatesMessage] = useState('');
  const [openSelectID, setOpenSelectID] = useState<string | null>(null);

  const fetchWithAuth = useCallback(async (url: string, init?: RequestInit) => {
    const res = await fetch(url, {
      ...init,
      credentials: 'include',
      headers: { ...(init?.headers ?? {}) },
      cache: 'no-store',
    });
    if (res.status === 401 || res.status === 403) {
      router.replace('/super-admin/login');
      return null;
    }
    return res;
  }, [router]);

  const fetchLogs = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const params = new URLSearchParams({ limit: String(PAGE_SIZE), offset: String(offset) });
      if (action) params.set('action', action);
      if (targetType) params.set('target_type', targetType);
      if (riskType.trim()) params.set('risk_type', riskType.trim());
      const res = await fetchWithAuth('/api/v1/super-admin/safety/moderation-logs?' + params.toString());
      if (!res) return;
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(payload.error ?? 'Safety 로그를 불러오지 못했습니다.');
      setLogs((payload.logs ?? []) as ModerationLog[]);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Safety 로그를 불러오지 못했습니다.');
    } finally {
      setLoading(false);
    }
  }, [action, fetchWithAuth, offset, riskType, targetType]);

  const fetchAIOutputSummary = useCallback(async () => {
    setAIOutputSummaryLoading(true);
    setAIOutputSummaryError('');
    try {
      const res = await fetchWithAuth('/api/v1/super-admin/safety/ai-output-summary');
      if (!res) return;
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(payload.error ?? 'AI output 관찰 집계를 불러오지 못했습니다.');
      setAIOutputSummary(payload as AIOutputSummary);
    } catch (err) {
      setAIOutputSummaryError(err instanceof Error ? err.message : 'AI output 관찰 집계를 불러오지 못했습니다.');
    } finally {
      setAIOutputSummaryLoading(false);
    }
  }, [fetchWithAuth]);

  const fetchAIOutputCandidates = useCallback(async () => {
    setAIOutputCandidatesLoading(true);
    setAIOutputCandidatesError('');
    try {
      const res = await fetchWithAuth('/api/v1/super-admin/safety/ai-output-review-candidates');
      if (!res) return;
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(payload.error ?? 'AI output 검토 후보를 불러오지 못했습니다.');
      setAIOutputCandidates(payload as AIOutputReviewCandidates);
    } catch (err) {
      setAIOutputCandidatesError(err instanceof Error ? err.message : 'AI output 검토 후보를 불러오지 못했습니다.');
    } finally {
      setAIOutputCandidatesLoading(false);
    }
  }, [fetchWithAuth]);
  const saveAIOutputReviewDecision = useCallback(async (candidateKey: string, status: AIOutputReviewStatus, note: string) => {
    setAIOutputReviewSavingKey(candidateKey);
    setAIOutputCandidatesError('');
    setAIOutputCandidatesMessage('');
    try {
      const res = await fetchWithAuth('/api/v1/super-admin/safety/ai-output-review-candidates/review', {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ candidate_key: candidateKey, status, note }),
      });
      if (!res) return;
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(payload.error ?? 'AI output 후보 판정을 저장하지 못했습니다.');
      setAIOutputCandidatesMessage('AI output 후보 판정을 저장했습니다.');
      await fetchAIOutputCandidates();
    } catch (err) {
      setAIOutputCandidatesError(err instanceof Error ? err.message : 'AI output 후보 판정을 저장하지 못했습니다.');
    } finally {
      setAIOutputReviewSavingKey(null);
    }
  }, [fetchAIOutputCandidates, fetchWithAuth]);

  const fetchRules = useCallback(async () => {
    setRulesLoading(true);
    setRulesError('');
    try {
      const params = new URLSearchParams({ limit: String(PAGE_SIZE), offset: String(ruleOffset) });
      if (ruleActiveFilter) params.set('is_active', ruleActiveFilter);
      if (ruleActionFilter) params.set('action', ruleActionFilter);
      if (ruleTypeFilter) params.set('rule_type', ruleTypeFilter);
      if (ruleQuery.trim()) params.set('q', ruleQuery.trim());
      const res = await fetchWithAuth('/api/v1/super-admin/safety/moderation-rules?' + params.toString());
      if (!res) return;
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(payload.error ?? 'Safety rule을 불러오지 못했습니다.');
      setRules((payload.rules ?? []) as ModerationRule[]);
    } catch (err) {
      setRulesError(err instanceof Error ? err.message : 'Safety rule을 불러오지 못했습니다.');
    } finally {
      setRulesLoading(false);
    }
  }, [fetchWithAuth, ruleActionFilter, ruleActiveFilter, ruleOffset, ruleQuery, ruleTypeFilter]);

  useEffect(() => { fetchLogs(); }, [fetchLogs]);
  useEffect(() => { if (activeTab === 'logs') fetchAIOutputSummary(); }, [activeTab, fetchAIOutputSummary]);
  useEffect(() => { if (activeTab === 'logs') fetchAIOutputCandidates(); }, [activeTab, fetchAIOutputCandidates]);
  useEffect(() => { if (activeTab === 'rules') fetchRules(); }, [activeTab, fetchRules]);
  useEffect(() => { setOffset(0); }, [action, targetType, riskType]);
  useEffect(() => { setRuleOffset(0); }, [ruleActiveFilter, ruleActionFilter, ruleTypeFilter, ruleQuery]);

  const visibleDescription = useMemo(() => activeTab === 'logs'
    ? '학습자 입력 moderation 결과를 원문 없이 hash, action, risk, route, target 기준으로 확인합니다.'
    : 'DB 기반 moderation rule을 생성, 수정, 비활성화하고 변경 이력을 audit event로 남깁니다.', [activeTab]);

  const normalizedRulePreview = useMemo(() => ruleForm.pattern.trim().toLowerCase().replace(/\s+/g, ' '), [ruleForm.pattern]);
  const duplicateRule = useMemo(() => rules.find((rule) => (
    rule.id !== editingRuleID
    && rule.rule_type === ruleForm.rule_type
    && rule.pattern.trim() === ruleForm.pattern.trim()
    && (rule.locale ?? '') === ruleForm.locale.trim()
  )) ?? null, [editingRuleID, ruleForm.locale, ruleForm.pattern, ruleForm.rule_type, rules]);

  const updateRuleForm = (field: keyof RuleForm, value: string) => {
    setRuleForm((current) => ({ ...current, [field]: value }));
  };

  const renderSelect = (id: string, value: string, onChange: (value: string) => void, options: readonly SelectOption[], minWidth?: CSSProperties['minWidth']) => (
    <SelectControl
      id={id}
      value={value}
      options={options}
      minWidth={minWidth}
      open={openSelectID === id}
      onToggle={() => setOpenSelectID((current) => (current === id ? null : id))}
      onClose={() => setOpenSelectID(null)}
      onChange={onChange}
    />
  );

  const startEditRule = (rule: ModerationRule) => {
    setEditingRuleID(rule.id);
    setRuleForm({
      rule_type: rule.rule_type,
      pattern: rule.pattern,
      risk_type: rule.risk_type,
      action: rule.action,
      locale: rule.locale ?? '',
      description: rule.description ?? '',
    });
    setRulesMessage('');
    setRulesError('');
  };

  const resetRuleForm = () => {
    setEditingRuleID(null);
    setRuleForm(emptyRuleForm);
  };

  const saveRule = async () => {
    setSavingRule(true);
    setRulesError('');
    setRulesMessage('');
    try {
      const body = JSON.stringify(ruleForm);
      const url = editingRuleID ? '/api/v1/super-admin/safety/moderation-rules/' + editingRuleID : '/api/v1/super-admin/safety/moderation-rules';
      const res = await fetchWithAuth(url, { method: editingRuleID ? 'PATCH' : 'POST', headers: { 'Content-Type': 'application/json' }, body });
      if (!res) return;
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(payload.error ?? 'Safety rule을 저장하지 못했습니다.');
      setRulesMessage(editingRuleID ? 'Rule을 수정했습니다.' : 'Rule을 생성했습니다.');
      resetRuleForm();
      await fetchRules();
    } catch (err) {
      setRulesError(err instanceof Error ? err.message : 'Safety rule을 저장하지 못했습니다.');
    } finally {
      setSavingRule(false);
    }
  };

  const setRuleActive = async (rule: ModerationRule, active: boolean) => {
    setSavingRule(true);
    setRulesError('');
    setRulesMessage('');
    try {
      const actionPath = active ? 'reactivate' : 'deactivate';
      const res = await fetchWithAuth('/api/v1/super-admin/safety/moderation-rules/' + rule.id + '/' + actionPath, { method: 'POST' });
      if (!res) return;
      const payload = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(payload.error ?? 'Safety rule 상태를 변경하지 못했습니다.');
      setRulesMessage(active ? 'Rule을 재활성화했습니다.' : 'Rule을 비활성화했습니다.');
      await fetchRules();
    } catch (err) {
      setRulesError(err instanceof Error ? err.message : 'Safety rule 상태를 변경하지 못했습니다.');
    } finally {
      setSavingRule(false);
    }
  };

  return (
    <main style={pageStyle}>
      <div style={superAdminShellStyle}>
        <div style={headerStackStyle}>
          <SuperAdminPanelHeader subtitle={activeTab === 'logs' ? 'Safety 로그' : 'Safety rule 관리'} description={visibleDescription} />
          <SuperAdminPanelNav activeSection="safety" />
        </div>

        <section style={tabBarStyle}>
          <button type="button" onClick={() => setActiveTab('logs')} style={tabButtonStyle(activeTab === 'logs')}>로그</button>
          <button type="button" onClick={() => setActiveTab('rules')} style={tabButtonStyle(activeTab === 'rules')}>룰 관리</button>
        </section>

        {activeTab === 'logs' ? (
          <>
            {error ? <div style={errorStyle}>{error}</div> : null}
            {aiOutputSummaryError ? <div style={errorStyle}>{aiOutputSummaryError}</div> : null}
            {aiOutputCandidatesError ? <div style={errorStyle}>{aiOutputCandidatesError}</div> : null}
            {aiOutputCandidatesMessage ? <div style={successStyle}>{aiOutputCandidatesMessage}</div> : null}
            <AIOutputSummaryPanel summary={aiOutputSummary} loading={aiOutputSummaryLoading} onRefresh={fetchAIOutputSummary} />
            <AIOutputCandidatesPanel candidates={aiOutputCandidates} loading={aiOutputCandidatesLoading} savingKey={aiOutputReviewSavingKey} onRefresh={fetchAIOutputCandidates} onSaveReview={saveAIOutputReviewDecision} />
            <section style={filterBarStyle}>
              {renderSelect('log-action', action, setAction, actionOptions, 124)}
              {renderSelect('log-target', targetType, setTargetType, targetOptions, 220)}
              <input value={riskType} onChange={(event) => setRiskType(event.target.value)} placeholder="risk_type 필터" style={inputStyle} />
              <button type="button" onClick={fetchLogs} style={buttonStyle}>새로고침</button>
            </section>
            <section style={panelStyle}>
              <div style={logTableHeaderStyle}><span>시간</span><span>Action</span><span>Risk</span><span>Target / Route</span><span>User / Hash</span></div>
              {loading ? <div style={emptyStyle}>불러오는 중입니다.</div> : null}
              {!loading && logs.length === 0 ? <div style={emptyStyle}>표시할 로그가 없습니다.</div> : null}
              {!loading && logs.map((log) => (
                <article key={log.id} style={logRowStyle}>
                  <div style={mutedStyle}>{formatDate(log.created_at)}</div>
                  <div><span style={badgeStyle(log.action)}>{log.action}</span></div>
                  <div style={monoStyle}>{log.risk_type || '-'}</div>
                  <div style={stackStyle}><strong>{log.target_type}</strong><span style={mutedStyle}>{log.route || '-'}</span>{log.target_id ? <span style={mutedStyle}>{log.target_id}</span> : null}</div>
                  <div style={stackStyle}><span>{log.user_email || log.user_id || '-'}</span><span style={monoStyle}>{log.input_text_hash}</span>{log.matched_pattern ? <span style={mutedStyle}>{log.rule_type}: {log.matched_pattern}</span> : null}</div>
                </article>
              ))}
            </section>
            <div style={pagerStyle}><button type="button" disabled={offset === 0} onClick={() => setOffset(Math.max(0, offset - PAGE_SIZE))} style={buttonStyle}>이전</button><span style={mutedStyle}>offset {offset}</span><button type="button" disabled={logs.length < PAGE_SIZE} onClick={() => setOffset(offset + PAGE_SIZE)} style={buttonStyle}>다음</button></div>
          </>
        ) : (
          <>
            {rulesError ? <div style={errorStyle}>{rulesError}</div> : null}
            {rulesMessage ? <div style={successStyle}>{rulesMessage}</div> : null}
            <section style={ruleFormStyle}>
              <strong>{editingRuleID ? 'Rule 수정' : 'Rule 생성'}</strong>
              {renderSelect('rule-form-type', ruleForm.rule_type, (value) => updateRuleForm('rule_type', value), ruleTypeOptions, 118)}
              <input value={ruleForm.pattern} onChange={(event) => updateRuleForm('pattern', event.target.value)} placeholder="pattern" style={inputStyle} />
              <input value={ruleForm.risk_type} onChange={(event) => updateRuleForm('risk_type', event.target.value)} placeholder="risk_type" style={inputStyle} />
              {renderSelect('rule-form-action', ruleForm.action, (value) => updateRuleForm('action', value), ruleFormActionOptions, 122)}
              <input value={ruleForm.locale} onChange={(event) => updateRuleForm('locale', event.target.value)} placeholder="locale 선택" style={smallInputStyle} />
              <input value={ruleForm.description} onChange={(event) => updateRuleForm('description', event.target.value)} placeholder="description" style={inputStyle} />
              <button type="button" disabled={savingRule || Boolean(duplicateRule) || !ruleForm.pattern.trim() || !ruleForm.risk_type.trim()} onClick={saveRule} style={buttonStyle}>{savingRule ? '저장 중' : editingRuleID ? '수정 저장' : '생성'}</button>
              {editingRuleID ? <button type="button" onClick={resetRuleForm} style={buttonStyle}>취소</button> : null}
              <span style={mutedStyle}>normalized: {normalizedRulePreview || '-'}</span>
              {duplicateRule ? <span style={warningInlineStyle}>중복 rule: {duplicateRule.id}</span> : null}
            </section>
            <section style={filterBarStyle}>
              {renderSelect('rule-active', ruleActiveFilter, setRuleActiveFilter, ruleActiveOptions, 118)}
              {renderSelect('rule-action', ruleActionFilter, setRuleActionFilter, ruleActionOptions, 124)}
              {renderSelect('rule-type', ruleTypeFilter, setRuleTypeFilter, ruleTypeFilterOptions, 118)}
              <input value={ruleQuery} onChange={(event) => setRuleQuery(event.target.value)} placeholder="pattern/description 검색" style={inputStyle} />
              <button type="button" onClick={fetchRules} style={buttonStyle}>새로고침</button>
            </section>
            <section style={panelStyle}>
              <div style={ruleTableHeaderStyle}><span>상태</span><span>Rule</span><span>Risk / Action</span><span>설명</span><span>관리</span></div>
              {rulesLoading ? <div style={emptyStyle}>불러오는 중입니다.</div> : null}
              {!rulesLoading && rules.length === 0 ? <div style={emptyStyle}>표시할 rule이 없습니다.</div> : null}
              {!rulesLoading && rules.map((rule) => (
                <article key={rule.id} style={ruleRowStyle}>
                  <div style={stackStyle}><span style={badgeStyle(rule.is_active ? 'allow' : 'soft_warn')}>{rule.is_active ? 'active' : 'inactive'}</span><span style={mutedStyle}>{formatDate(rule.updated_at)}</span></div>
                  <div style={stackStyle}><strong style={monoStyle}>{rule.pattern}</strong><span style={mutedStyle}>{rule.rule_type} / {rule.locale || 'all locales'}</span></div>
                  <div style={stackStyle}><span style={monoStyle}>{rule.risk_type}</span><span style={badgeStyle(rule.action)}>{rule.action}</span></div>
                  <div style={mutedStyle}>{rule.description || '-'}</div>
                  <div style={rowActionStyle}><button type="button" onClick={() => startEditRule(rule)} style={buttonStyle}>수정</button>{rule.is_active ? <button type="button" onClick={() => setRuleActive(rule, false)} style={dangerButtonStyle}>비활성</button> : <button type="button" onClick={() => setRuleActive(rule, true)} style={buttonStyle}>재활성</button>}</div>
                </article>
              ))}
            </section>
            <div style={pagerStyle}><button type="button" disabled={ruleOffset === 0} onClick={() => setRuleOffset(Math.max(0, ruleOffset - PAGE_SIZE))} style={buttonStyle}>이전</button><span style={mutedStyle}>offset {ruleOffset}</span><button type="button" disabled={rules.length < PAGE_SIZE} onClick={() => setRuleOffset(ruleOffset + PAGE_SIZE)} style={buttonStyle}>다음</button></div>
          </>
        )}
      </div>
    </main>
  );
}

function AIOutputSummaryPanel({ summary, loading, onRefresh }: { summary: AIOutputSummary | null; loading: boolean; onRefresh: () => void }) {
  const windows = summary?.windows ?? [];
  return (
    <section style={aiSummaryPanelStyle}>
      <div style={aiSummaryHeaderStyle}>
        <div style={stackStyle}>
          <strong>AI output 관찰 상태</strong>
          <span style={mutedStyle}>최근 24시간/7일 기준 log-only moderation 반복 신호를 집계합니다.</span>
        </div>
        <button type="button" onClick={onRefresh} style={buttonStyle}>집계 새로고침</button>
      </div>
      {loading ? <div style={emptyStyle}>AI output 집계를 불러오는 중입니다.</div> : null}
      {!loading && windows.length === 0 ? <div style={emptyStyle}>관찰 데이터가 없습니다.</div> : null}
      {!loading && windows.length > 0 ? (
        <div style={aiSummaryGridStyle}>
          {windows.map((window) => (
            <article key={window.label} style={aiSummaryWindowStyle}>
              <div style={aiSummaryWindowHeaderStyle}>
                <span style={aiWindowTitleStyle}>{window.label}</span>
                <span style={severityBadgeStyle(window.severity)}>{severityLabel(window.severity)}</span>
              </div>
              <div style={aiMetricGridStyle}>
                <Metric label="전체" value={window.total} />
                <Metric label="위험" value={window.risky_total} />
                <Metric label="위험률" value={formatPercent(window.risky_rate)} />
              </div>
              <div style={aiActionLineStyle}>
                <span>allow {window.action_counts.allow}</span>
                <span>soft {window.action_counts.soft_warn}</span>
                <span>block {window.action_counts.block}</span>
              </div>
              <p style={aiSummaryMessageStyle}>{window.message}</p>
              <div style={aiSummaryListsStyle}>
                <SummaryList title="Target" items={window.targets.slice(0, 3).map((item) => ({ key: item.target_type, count: item.risky_total || item.total }))} empty="target 없음" />
                <SummaryList title="Risk" items={window.risk_types.slice(0, 3)} empty="위험 반복 없음" />
                <SummaryList title="Source" items={window.source_features.slice(0, 3)} empty="source 반복 없음" />
              </div>
            </article>
          ))}
        </div>
      ) : null}
      {summary?.generated_at ? <div style={aiSummaryFooterStyle}>생성 시각 {formatDate(summary.generated_at)}</div> : null}
    </section>
  );
}

function AIOutputCandidatesPanel({
  candidates,
  loading,
  savingKey,
  onRefresh,
  onSaveReview,
}: {
  candidates: AIOutputReviewCandidates | null;
  loading: boolean;
  savingKey: string | null;
  onRefresh: () => void;
  onSaveReview: (candidateKey: string, status: AIOutputReviewStatus, note: string) => Promise<void>;
}) {
  const items = candidates?.candidates ?? [];
  return (
    <section style={aiCandidatesPanelStyle}>
      <div style={aiSummaryHeaderStyle}>
        <div style={stackStyle}>
          <strong>AI output 검토 후보</strong>
          <span style={mutedStyle}>최근 30일 soft_warn/block 로그를 원문 없이 target, risk, source, rule 기준으로 묶습니다.</span>
        </div>
        <button type="button" onClick={onRefresh} style={buttonStyle}>후보 새로고침</button>
      </div>
      <div style={reviewRubricStyle}>
        {reviewRubricItems.map((item) => (
          <span key={item.status} style={reviewRubricItemStyle}>
            <strong>{reviewStatusLabel(item.status)}</strong>
            <span>{item.description}</span>
          </span>
        ))}
      </div>
      {loading ? <div style={emptyStyle}>AI output 후보를 불러오는 중입니다.</div> : null}
      {!loading && items.length === 0 ? (
        <div style={aiCandidateEmptyStyle}>
          <strong>현재 검토할 AI output 후보가 없습니다.</strong>
          <span>soft_warn 또는 block AI output 로그가 쌓이면 이 영역에 후보가 표시됩니다.</span>
        </div>
      ) : null}
      {!loading && items.length > 0 ? (
        <div style={aiCandidateListStyle}>
          {items.map((item) => (
            <article key={item.key} style={aiCandidateRowStyle}>
              <div style={stackStyle}>
                <strong style={monoStyle}>{item.target_type}</strong>
                <span style={mutedStyle}>{item.source_feature} / {item.risk_type}</span>
                <span style={reviewStatusBadgeStyle(item.review.status)}>{reviewStatusLabel(item.review.status)}</span>
              </div>
              <div style={aiCandidateCountsStyle}>
                <Metric label="전체" value={item.total} />
                <Metric label="soft" value={item.action_counts.soft_warn} />
                <Metric label="block" value={item.action_counts.block} />
              </div>
              <div style={stackStyle}>
                <span style={mutedStyle}>최근 {formatDate(item.latest_seen_at)}</span>
                {item.matched_pattern || item.matched_rule_id ? <span style={monoStyle}>{item.matched_pattern || item.matched_rule_id}</span> : <span style={mutedStyle}>matched rule 없음</span>}
                {item.sample_target_id ? <span style={mutedStyle}>{item.sample_target_id}</span> : null}
                {item.review.reviewed_at ? <span style={mutedStyle}>검토 {formatDate(item.review.reviewed_at)} / {item.review.reviewed_by || '-'}</span> : null}
              </div>
              <AIOutputCandidateReviewForm item={item} saving={savingKey === item.key} onSave={onSaveReview} />
            </article>
          ))}
        </div>
      ) : null}
      {candidates?.generated_at ? <div style={aiSummaryFooterStyle}>생성 시각 {formatDate(candidates.generated_at)} / 기준 {candidates.window_hours}h</div> : null}
    </section>
  );
}

function AIOutputCandidateReviewForm({ item, saving, onSave }: { item: AIOutputReviewCandidate; saving: boolean; onSave: (candidateKey: string, status: AIOutputReviewStatus, note: string) => Promise<void> }) {
  const [status, setStatus] = useState<AIOutputReviewStatus>(item.review.status);
  const [note, setNote] = useState(item.review.note ?? '');

  useEffect(() => {
    setStatus(item.review.status);
    setNote(item.review.note ?? '');
  }, [item.key, item.review.note, item.review.status]);

  const changed = status !== item.review.status || note.trim() !== (item.review.note ?? '').trim();

  return (
    <div style={aiCandidateReviewStyle}>
      <label style={reviewLabelStyle}>
        <span>판정</span>
        <select value={status} onChange={(event) => setStatus(event.target.value as AIOutputReviewStatus)} style={reviewSelectStyle}>
          {aiOutputReviewStatusOptions.map((option) => <option key={option.value} value={option.value}>{option.label}</option>)}
        </select>
      </label>
      <label style={reviewLabelStyle}>
        <span>메모</span>
        <textarea value={note} onChange={(event) => setNote(event.target.value)} placeholder="판정 근거 또는 후속 조치 메모" rows={3} style={reviewTextareaStyle} />
      </label>
      <button type="button" disabled={saving || !changed} onClick={() => onSave(item.key, status, note)} style={buttonStyle}>{saving ? '저장 중' : '판정 저장'}</button>
    </div>
  );
}

function Metric({ label, value }: { label: string; value: number | string }) {
  return <div style={metricStyle}><span style={mutedStyle}>{label}</span><strong>{value}</strong></div>;
}

function SummaryList({ title, items, empty }: { title: string; items: { key: string; count: number }[]; empty: string }) {
  return (
    <div style={stackStyle}>
      <span style={aiSummaryListTitleStyle}>{title}</span>
      {items.length === 0 ? <span style={mutedStyle}>{empty}</span> : items.map((item) => (
        <span key={item.key} style={aiSummaryItemStyle}><span>{item.key}</span><strong>{item.count}</strong></span>
      ))}
    </div>
  );
}

function reviewStatusLabel(status: AIOutputReviewStatus) {
  if (status === 'false_positive') return 'False positive';
  if (status === 'needs_prompt_guard') return 'Prompt guard 필요';
  if (status === 'needs_rule_tuning') return 'Rule tuning 필요';
  if (status === 'needs_masking') return 'Masking 필요';
  if (status === 'needs_regeneration') return 'Regeneration 필요';
  if (status === 'resolved') return 'Resolved';
  return '미검토';
}

function severityLabel(severity: AIOutputSeverity) {
  if (severity === 'critical') return '심각';
  if (severity === 'warning') return '경고';
  if (severity === 'caution') return '주의';
  if (severity === 'low') return '낮음';
  return '데이터 없음';
}

function formatPercent(value: number) {
  return value.toFixed(value % 1 === 0 ? 0 : 1) + '%';
}

function SelectControl({
  id,
  value,
  options,
  minWidth,
  open,
  onToggle,
  onClose,
  onChange,
}: {
  id: string;
  value: string;
  options: readonly SelectOption[];
  minWidth?: CSSProperties['minWidth'];
  open: boolean;
  onToggle: () => void;
  onClose: () => void;
  onChange: (value: string) => void;
}) {
  const selectedOption = options.find((option) => option.value === value) ?? options[0];

  return (
    <div
      style={{ ...selectWrapStyle, minWidth }}
      onBlur={(event) => {
        const nextTarget = event.relatedTarget as Node | null;
        if (!nextTarget || !event.currentTarget.contains(nextTarget)) onClose();
      }}
    >
      <button
        id={id}
        type="button"
        aria-haspopup="listbox"
        aria-expanded={open}
        aria-controls={id + '-listbox'}
        onClick={onToggle}
        style={selectButtonStyle}
      >
        <span style={selectLabelStyle}>{selectedOption?.label ?? value}</span>
        <span aria-hidden="true" style={selectArrowStyle}>v</span>
      </button>
      {open ? (
        <div id={id + '-listbox'} role="listbox" aria-labelledby={id} style={selectMenuStyle}>
          {options.map((option) => {
            const selected = option.value === value;
            return (
              <button
                key={option.value}
                type="button"
                role="option"
                aria-selected={selected}
                onClick={() => {
                  onChange(option.value);
                  onClose();
                }}
                style={selectOptionStyle(selected)}
              >
                {option.label}
              </button>
            );
          })}
        </div>
      ) : null}
    </div>
  );
}

function formatDate(value: string) {
  return new Date(value).toLocaleString('ko-KR', { timeZone: 'Asia/Seoul', year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' });
}

const pageStyle: CSSProperties = { minHeight: '100vh', background: 'linear-gradient(180deg, #07111f, #0b1321 48%, #060b12)', color: '#eff6ff', padding: '28px 20px 48px' };
const headerStackStyle: CSSProperties = { display: 'grid', gap: '10px' };
const tabBarStyle: CSSProperties = { display: 'flex', gap: '8px', margin: '4px 0 14px' };
const tabButtonStyle = (active: boolean): CSSProperties => ({ minHeight: '38px', borderRadius: '8px', border: active ? '1px solid rgba(147,197,253,0.72)' : '1px solid rgba(125,160,255,0.28)', background: active ? 'rgba(59,130,246,0.34)' : 'rgba(255,255,255,0.06)', color: '#eff6ff', padding: '0 16px', fontWeight: 800, cursor: 'pointer' });
const filterBarStyle: CSSProperties = { display: 'flex', flexWrap: 'wrap', gap: '10px', marginBottom: '14px' };
const aiSummaryPanelStyle: CSSProperties = { border: '1px solid rgba(125,160,255,0.22)', borderRadius: '8px', padding: '14px', marginBottom: '14px', background: 'rgba(5,12,22,0.72)' };
const aiSummaryHeaderStyle: CSSProperties = { display: 'flex', justifyContent: 'space-between', gap: '12px', alignItems: 'flex-start', marginBottom: '12px', flexWrap: 'wrap' };
const aiSummaryGridStyle: CSSProperties = { display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: '12px' };
const aiSummaryWindowStyle: CSSProperties = { border: '1px solid rgba(125,160,255,0.18)', borderRadius: '8px', padding: '12px', background: 'rgba(255,255,255,0.04)', display: 'grid', gap: '10px' };
const aiSummaryWindowHeaderStyle: CSSProperties = { display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: '10px' };
const aiWindowTitleStyle: CSSProperties = { fontSize: '16px', fontWeight: 900 };
const aiMetricGridStyle: CSSProperties = { display: 'grid', gridTemplateColumns: 'repeat(3, minmax(0, 1fr))', gap: '8px' };
const metricStyle: CSSProperties = { minHeight: '58px', borderRadius: '8px', border: '1px solid rgba(125,160,255,0.14)', padding: '8px', background: 'rgba(0,0,0,0.14)', display: 'grid', gap: '4px' };
const aiActionLineStyle: CSSProperties = { display: 'flex', flexWrap: 'wrap', gap: '8px', color: 'rgba(219,234,254,0.82)', fontSize: '12px', fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' };
const aiSummaryMessageStyle: CSSProperties = { margin: 0, color: 'rgba(219,234,254,0.86)', fontSize: '13px', lineHeight: 1.5 };
const aiSummaryListsStyle: CSSProperties = { display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(130px, 1fr))', gap: '10px' };
const aiSummaryListTitleStyle: CSSProperties = { color: 'rgba(200,210,235,0.72)', fontSize: '11px', fontWeight: 900, textTransform: 'uppercase' };
const aiSummaryItemStyle: CSSProperties = { display: 'flex', justifyContent: 'space-between', gap: '8px', color: '#dbeafe', fontSize: '12px', overflowWrap: 'anywhere' };
const aiSummaryFooterStyle: CSSProperties = { marginTop: '10px', color: 'rgba(200,210,235,0.58)', fontSize: '11px' };
const aiCandidatesPanelStyle: CSSProperties = { ...aiSummaryPanelStyle, background: 'rgba(8,18,28,0.72)' };
const aiCandidateEmptyStyle: CSSProperties = { display: 'grid', gap: '6px', padding: '18px 14px', border: '1px dashed rgba(125,160,255,0.26)', borderRadius: '8px', color: 'rgba(219,234,254,0.82)', background: 'rgba(255,255,255,0.03)' };
const aiCandidateListStyle: CSSProperties = { display: 'grid', gap: '10px' };
const aiCandidateRowStyle: CSSProperties = { display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '12px', alignItems: 'start', border: '1px solid rgba(125,160,255,0.16)', borderRadius: '8px', padding: '12px', background: 'rgba(255,255,255,0.04)' };
const aiCandidateCountsStyle: CSSProperties = { display: 'grid', gridTemplateColumns: 'repeat(3, minmax(0, 1fr))', gap: '8px' };
const reviewRubricStyle: CSSProperties = { display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))', gap: '8px', marginBottom: '12px' };
const reviewRubricItemStyle: CSSProperties = { minHeight: '46px', borderRadius: '8px', border: '1px solid rgba(125,160,255,0.14)', background: 'rgba(255,255,255,0.035)', padding: '8px 10px', display: 'grid', gap: '3px', color: 'rgba(219,234,254,0.78)', fontSize: '11px' };
const aiCandidateReviewStyle: CSSProperties = { display: 'grid', gap: '8px', minWidth: 0 };
const reviewLabelStyle: CSSProperties = { display: 'grid', gap: '5px', color: 'rgba(219,234,254,0.82)', fontSize: '12px', fontWeight: 800 };
const reviewSelectStyle: CSSProperties = { minHeight: '36px', borderRadius: '8px', border: '1px solid rgba(125,160,255,0.28)', background: '#0b1321', color: '#eff6ff', padding: '0 10px' };
const reviewTextareaStyle: CSSProperties = { minHeight: '76px', borderRadius: '8px', border: '1px solid rgba(125,160,255,0.28)', background: 'rgba(255,255,255,0.06)', color: '#eff6ff', padding: '9px 10px', resize: 'vertical', lineHeight: 1.45 };
const reviewStatusBadgeStyle = (status: AIOutputReviewStatus): CSSProperties => ({
  display: 'inline-flex', alignItems: 'center', width: 'fit-content', minHeight: '24px', borderRadius: '999px', padding: '0 9px', fontSize: '11px', fontWeight: 900,
  color: status === 'false_positive' || status === 'resolved' ? '#bbf7d0' : status === 'unreviewed' ? 'rgba(219,234,254,0.76)' : '#fde68a',
  background: status === 'false_positive' || status === 'resolved' ? 'rgba(22,163,74,0.16)' : status === 'unreviewed' ? 'rgba(148,163,184,0.12)' : 'rgba(217,119,6,0.16)',
  border: status === 'false_positive' || status === 'resolved' ? '1px solid rgba(74,222,128,0.3)' : status === 'unreviewed' ? '1px solid rgba(148,163,184,0.22)' : '1px solid rgba(251,191,36,0.3)',
});
const selectStyle: CSSProperties = { minHeight: '38px', borderRadius: '8px', border: '1px solid rgba(125,160,255,0.28)', background: 'rgba(255,255,255,0.06)', color: '#eff6ff', padding: '0 10px' };
const selectWrapStyle: CSSProperties = { position: 'relative', display: 'inline-flex' };
const selectButtonStyle: CSSProperties = { ...selectStyle, width: '100%', display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: '10px', textAlign: 'left', cursor: 'pointer' };
const selectLabelStyle: CSSProperties = { minWidth: 0, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' };
const selectArrowStyle: CSSProperties = { color: 'rgba(200,210,235,0.76)', fontSize: '11px', lineHeight: 1 };
const selectMenuStyle: CSSProperties = { position: 'absolute', zIndex: 80, top: 'calc(100% + 6px)', left: 0, minWidth: '100%', maxHeight: '260px', overflowY: 'auto', borderRadius: '8px', border: '1px solid rgba(125,160,255,0.36)', background: '#0b1321', color: '#eff6ff', boxShadow: '0 18px 48px rgba(0,0,0,0.45)', padding: '6px' };
const selectOptionStyle = (selected: boolean): CSSProperties => ({ width: '100%', minHeight: '32px', border: 0, borderRadius: '6px', background: selected ? 'rgba(95,131,255,0.28)' : '#0b1321', color: '#eff6ff', padding: '0 10px', textAlign: 'left', fontWeight: selected ? 800 : 600, cursor: 'pointer', whiteSpace: 'nowrap' });
const inputStyle: CSSProperties = { ...selectStyle, minWidth: '180px' };
const smallInputStyle: CSSProperties = { ...selectStyle, width: '118px' };
const buttonStyle: CSSProperties = { minHeight: '38px', borderRadius: '8px', border: '1px solid rgba(125,160,255,0.32)', background: 'rgba(95,131,255,0.18)', color: '#eff6ff', padding: '0 14px', fontWeight: 700, cursor: 'pointer' };
const dangerButtonStyle: CSSProperties = { ...buttonStyle, borderColor: 'rgba(248,113,113,0.36)', background: 'rgba(127,29,29,0.26)' };
const panelStyle: CSSProperties = { border: '1px solid rgba(125,160,255,0.22)', borderRadius: '8px', overflow: 'hidden', background: 'rgba(5,12,22,0.72)' };
const logTableHeaderStyle: CSSProperties = { display: 'grid', gridTemplateColumns: '150px 110px 120px minmax(220px,1.1fr) minmax(240px,1fr)', gap: '12px', padding: '12px 14px', borderBottom: '1px solid rgba(125,160,255,0.18)', color: 'rgba(200,210,235,0.76)', fontSize: '12px', fontWeight: 700 };
const logRowStyle: CSSProperties = { display: 'grid', gridTemplateColumns: '150px 110px 120px minmax(220px,1.1fr) minmax(240px,1fr)', gap: '12px', padding: '13px 14px', borderBottom: '1px solid rgba(125,160,255,0.12)', alignItems: 'start' };
const ruleTableHeaderStyle: CSSProperties = { ...logTableHeaderStyle, gridTemplateColumns: '120px minmax(220px,1.1fr) 150px minmax(190px,0.9fr) 180px' };
const ruleRowStyle: CSSProperties = { ...logRowStyle, gridTemplateColumns: '120px minmax(220px,1.1fr) 150px minmax(190px,0.9fr) 180px' };
const ruleFormStyle: CSSProperties = { display: 'flex', flexWrap: 'wrap', gap: '10px', alignItems: 'center', marginBottom: '14px', padding: '12px', border: '1px solid rgba(125,160,255,0.22)', borderRadius: '8px', background: 'rgba(5,12,22,0.72)' };
const rowActionStyle: CSSProperties = { display: 'flex', gap: '8px', flexWrap: 'wrap' };
const stackStyle: CSSProperties = { display: 'grid', gap: '4px', minWidth: 0 };
const mutedStyle: CSSProperties = { color: 'rgba(200,210,235,0.68)', fontSize: '12px', overflowWrap: 'anywhere' };
const monoStyle: CSSProperties = { color: '#dbeafe', fontSize: '12px', fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace', overflowWrap: 'anywhere' };
const emptyStyle: CSSProperties = { padding: '26px 14px', color: 'rgba(200,210,235,0.72)' };
const errorStyle: CSSProperties = { marginBottom: '14px', border: '1px solid rgba(248,113,113,0.38)', borderRadius: '8px', padding: '12px', color: '#fecaca', background: 'rgba(127,29,29,0.2)' };
const successStyle: CSSProperties = { marginBottom: '14px', border: '1px solid rgba(74,222,128,0.34)', borderRadius: '8px', padding: '12px', color: '#bbf7d0', background: 'rgba(22,101,52,0.2)' };
const warningInlineStyle: CSSProperties = { color: '#fde68a', fontSize: '12px', fontWeight: 800, overflowWrap: 'anywhere' };
const pagerStyle: CSSProperties = { display: 'flex', justifyContent: 'flex-end', alignItems: 'center', gap: '10px', marginTop: '12px' };
const severityBadgeStyle = (severity: AIOutputSeverity): CSSProperties => ({
  display: 'inline-flex', alignItems: 'center', width: 'fit-content', minHeight: '26px', borderRadius: '999px', padding: '0 10px', fontSize: '12px', fontWeight: 900,
  color: severity === 'critical' ? '#fecaca' : severity === 'warning' ? '#fed7aa' : severity === 'caution' ? '#fde68a' : severity === 'low' ? '#bbf7d0' : 'rgba(219,234,254,0.72)',
  background: severity === 'critical' ? 'rgba(127,29,29,0.28)' : severity === 'warning' ? 'rgba(154,52,18,0.24)' : severity === 'caution' ? 'rgba(217,119,6,0.16)' : severity === 'low' ? 'rgba(22,163,74,0.16)' : 'rgba(148,163,184,0.12)',
  border: severity === 'critical' ? '1px solid rgba(248,113,113,0.34)' : severity === 'warning' ? '1px solid rgba(251,146,60,0.34)' : severity === 'caution' ? '1px solid rgba(251,191,36,0.3)' : severity === 'low' ? '1px solid rgba(74,222,128,0.3)' : '1px solid rgba(148,163,184,0.22)',
});
const badgeStyle = (action: ModerationAction | RuleAction | 'active' | 'inactive'): CSSProperties => ({
  display: 'inline-flex', alignItems: 'center', width: 'fit-content', minHeight: '26px', borderRadius: '999px', padding: '0 10px', fontSize: '12px', fontWeight: 800,
  color: action === 'block' || action === 'inactive' ? '#fecaca' : action === 'soft_warn' ? '#fde68a' : '#bbf7d0',
  background: action === 'block' || action === 'inactive' ? 'rgba(220,38,38,0.16)' : action === 'soft_warn' ? 'rgba(217,119,6,0.16)' : 'rgba(22,163,74,0.16)',
  border: action === 'block' || action === 'inactive' ? '1px solid rgba(248,113,113,0.3)' : action === 'soft_warn' ? '1px solid rgba(251,191,36,0.3)' : '1px solid rgba(74,222,128,0.3)',
});
