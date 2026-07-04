import { API_BASE } from './constants'
import type {
  DebugLabel,
  DebugScenario,
  DebugScenarioDetail,
  GeneratedLessonItem,
  GeneratedLessonsSnapshot,
  GoalProfile,
  LabelSummary,
  LLMJobObservations,
  LessonSearchPrerunReportsResponse,
  RankerArtifactStatus,
  SavedExternalCandidate,
  RecommendationComparison,
  RecommendationSpecMetrics,
  RolloutState,
  ScenarioInput,
} from './types'

export type AuthHeaders = Record<string, string>

const readJSON = async <T,>(res: Response): Promise<T> => res.json().catch(() => ({} as T))

const ensureOK = async <T,>(res: Response, fallbackMessage: string): Promise<T> => {
  const data = await readJSON<T & { error?: string }>(res)
  if (!res.ok) {
    throw new Error(data.error ?? fallbackMessage)
  }
  return data
}

type LLMJobResponse = {
  job_id?: string
  feature?: string
  status?: 'queued' | 'running' | 'succeeded' | 'failed' | 'canceled' | 'expired'
  poll_url?: string
  result_ref?: {
    snapshot?: GeneratedLessonsSnapshot
    scenario?: Record<string, unknown>
    status?: string
    goal_profile_version?: number
    message_count?: number
  }
  error?: string
  error_code?: string
}

const sleep = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

const fetchLLMJob = async (headers: AuthHeaders, jobId: string) => {
  const res = await fetch(`${API_BASE}/api/v1/llm-jobs/${jobId}`, { headers, credentials: 'include' })
  return ensureOK<LLMJobResponse>(res, '작업 상태 조회 실패')
}

const waitForLLMJobSnapshot = async (headers: AuthHeaders, jobId: string) => {
  const startedAt = Date.now()
  while (Date.now() - startedAt < 90000) {
    const job = await fetchLLMJob(headers, jobId)
    if (job.status === 'succeeded') {
      if (!job.result_ref?.snapshot) {
        throw new Error('리슨 생성 결과가 비어 있습니다.')
      }
      return job.result_ref.snapshot
    }
    if (job.status === 'failed' || job.status === 'canceled' || job.status === 'expired') {
      throw new Error(job.error_code ?? job.error ?? '리슨 생성 작업 실패')
    }
    await sleep(1000)
  }
  throw new Error('리슨 생성 작업이 아직 완료되지 않았습니다. 잠시 후 다시 시도해주세요.')
}

const waitForGoalJob = async (headers: AuthHeaders, scenarioId: string, jobId: string) => {
  const startedAt = Date.now()
  while (Date.now() - startedAt < 60000) {
    const job = await fetchLLMJob(headers, jobId)
    if (job.status === 'succeeded') {
      const scenario = await fetchScenarioDetail(headers, scenarioId)
      if (!scenario.goal_profile_snapshot) {
        throw new Error('목표 채팅 결과가 비어 있습니다.')
      }
      return JSON.parse(scenario.goal_profile_snapshot) as GoalProfile
    }
    if (job.status === 'failed' || job.status === 'canceled' || job.status === 'expired') {
      throw new Error(job.error_code ?? job.error ?? '목표 채팅 작업 실패')
    }
    await sleep(1000)
  }
  throw new Error('목표 채팅 작업이 아직 완료되지 않았습니다. 잠시 후 다시 시도해주세요.')
}

export const fetchScenarioConsole = async (headers: AuthHeaders) => {
  const [rolloutRes, artifactRes, scenariosRes, observationsRes, specMetricsRes, prerunReportsRes] = await Promise.all([
    fetch(`${API_BASE}/api/v1/super-admin/recommendation-debug/rollout-state`, { headers, credentials: 'include' }),
    fetch(`${API_BASE}/api/v1/super-admin/recommendation-debug/ranker-artifact`, { headers, credentials: 'include' }),
    fetch(`${API_BASE}/api/v1/super-admin/recommendation-debug/scenarios`, { headers, credentials: 'include' }),
    fetch(`${API_BASE}/api/v1/llm-jobs/observations?queued_timeout=10m&running_timeout=30m`, { headers, credentials: 'include' }),
    fetch(API_BASE + "/api/v1/super-admin/recommendation-debug/spec-metrics?days=7", { headers, credentials: 'include' }),
    fetch(API_BASE + "/api/v1/super-admin/recommendation-debug/lesson-search-prerun-reports?limit=5", { headers, credentials: "include" }),
  ])
  const rolloutData = await ensureOK<{ state?: RolloutState | null }>(rolloutRes, '전환 상태 조회 실패')
  const artifactData = await ensureOK<{ artifact?: RankerArtifactStatus | null }>(artifactRes, 'Ranker artifact 조회 실패')
  const scenariosData = await ensureOK<{ scenarios?: DebugScenario[] }>(scenariosRes, '시나리오 목록 조회 실패')
  const observationsData = await ensureOK<LLMJobObservations>(observationsRes, 'worker 관측 조회 실패')
  const specMetricsData = await ensureOK<{ metrics?: RecommendationSpecMetrics }>(specMetricsRes, '명세 관측 조회 실패')
  const prerunReportsData = await ensureOK<LessonSearchPrerunReportsResponse>(prerunReportsRes, '리슨 검색 프리런 report 조회 실패')
  return {
    rolloutState: rolloutData.state ?? null,
    rankerArtifact: artifactData.artifact ?? null,
    scenarios: scenariosData.scenarios ?? [],
    workerObservations: observationsData,
    specMetrics: specMetricsData.metrics ?? null,
    lessonSearchPrerunReports: prerunReportsData,
  }
}

export const fetchWorkerObservations = async (headers: AuthHeaders) => {
  const res = await fetch(`${API_BASE}/api/v1/llm-jobs/observations?queued_timeout=10m&running_timeout=30m`, { headers, credentials: 'include' })
  return ensureOK<LLMJobObservations>(res, 'worker 관측 조회 실패')
}

export const fetchRecommendationSpecMetrics = async (headers: AuthHeaders, days = 7) => {
  const res = await fetch(`${API_BASE}/api/v1/super-admin/recommendation-debug/spec-metrics?days=${days}`, { headers, credentials: 'include' })
  const data = await ensureOK<{ metrics?: RecommendationSpecMetrics }>(res, '명세 관측 조회 실패')
  return data.metrics ?? null
}

export const patchRolloutRanker = async (
  headers: AuthHeaders,
  input: { ranker_model_version: string; feature_schema_version: string },
) => {
  const res = await fetch(`${API_BASE}/api/v1/super-admin/recommendation-debug/rollout-ranker`, {
    method: 'PATCH',
    headers: {
      ...headers,
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify(input),
  })
  const data = await ensureOK<{ state?: RolloutState | null }>(res, 'Ranker 설정 저장 실패')
  return data.state ?? null
}

export const fetchRankerArtifact = async (headers: AuthHeaders) => {
  const res = await fetch(`${API_BASE}/api/v1/super-admin/recommendation-debug/ranker-artifact`, { headers, credentials: 'include' })
  const data = await ensureOK<{ artifact?: RankerArtifactStatus | null }>(res, 'Ranker artifact 조회 실패')
  return data.artifact ?? null
}

export const fetchLabelSummary = async (headers: AuthHeaders, scenarioId: string) => {
  const res = await fetch(`${API_BASE}/api/v1/super-admin/recommendation-debug/scenarios/${scenarioId}/labels/summary`, { headers, credentials: 'include' })
  const data = await ensureOK<{ summary?: LabelSummary }>(res, '라벨 요약 조회 실패')
  return data.summary ?? null
}

export const fetchScenarioDetail = async (headers: AuthHeaders, scenarioId: string) => {
  const res = await fetch(`${API_BASE}/api/v1/super-admin/recommendation-debug/scenarios/${scenarioId}`, { headers, credentials: 'include' })
  const data = await ensureOK<{ scenario?: DebugScenarioDetail }>(res, '시나리오 상세 조회 실패')
  if (!data.scenario) {
    throw new Error('시나리오 상세가 비어 있습니다.')
  }
  return data.scenario
}

export const createScenario = async (headers: AuthHeaders, input: ScenarioInput) => {
  const res = await fetch(`${API_BASE}/api/v1/super-admin/recommendation-debug/scenarios`, {
    method: 'POST',
    headers: {
      ...headers,
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({
      course_title: input.course_title,
      initial_user_intent: '',
      notes: input.notes,
    }),
  })
  const data = await ensureOK<{ scenario?: DebugScenario }>(res, '시나리오 생성 실패')
  return data.scenario ?? null
}

export const runGoalAction = async (
  headers: AuthHeaders,
  scenarioId: string,
  action: 'start' | 'message' | 'confirm',
  body: Record<string, string>,
) => {
  const res = await fetch(`${API_BASE}/api/v1/super-admin/recommendation-debug/scenarios/${scenarioId}/goal/${action}`, {
    method: 'POST',
    headers: {
      ...headers,
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify(body),
  })
  const data = await ensureOK<{ goal?: GoalProfile } & LLMJobResponse>(res, '목표 채팅 처리 실패')
  if (data.goal) return data.goal
  if (data.status && data.job_id) return waitForGoalJob(headers, scenarioId, data.job_id)
  return null
}

export const generateLessons = async (headers: AuthHeaders, scenarioId: string) => {
  const res = await fetch(`${API_BASE}/api/v1/super-admin/recommendation-debug/scenarios/${scenarioId}/lessons/generate`, {
    method: 'POST',
    headers: {
      ...headers,
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({}),
  })
  const data = await ensureOK<{ snapshot?: GeneratedLessonsSnapshot } & LLMJobResponse>(res, '리슨 생성 실패')
  if (data.snapshot) return data.snapshot
  if (data.status && data.job_id) return waitForLLMJobSnapshot(headers, data.job_id)
  return null
}

export const compareRecommendations = async (
  headers: AuthHeaders,
  scenarioId: string,
  lesson?: GeneratedLessonItem,
  recommendationQuery = '',
) => {
  const res = await fetch(`${API_BASE}/api/v1/super-admin/recommendation-debug/scenarios/${scenarioId}/recommendations/compare`, {
    method: 'POST',
    headers: {
      ...headers,
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({
      max_per_lesson: 12,
      lesson_id: lesson?.lesson_id ?? '',
      recommendation_query: recommendationQuery.trim(),
    }),
  })
  const data = await ensureOK<{ comparison?: RecommendationComparison }>(res, '추천 비교 실패')
  return data.comparison ?? null
}

export const saveLabel = async (
  headers: AuthHeaders,
  scenarioId: string,
  payload: {
    run_id: string
    candidate_key: string
    content_id: string
    url: string
    label: string
    note: string
    baseline_rank: number
    feature_snapshot: Record<string, unknown>
  },
) => {
  const res = await fetch(`${API_BASE}/api/v1/super-admin/recommendation-debug/scenarios/${scenarioId}/labels`, {
    method: 'POST',
    headers: {
      ...headers,
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify(payload),
  })
  const data = await ensureOK<{ label?: DebugLabel }>(res, '라벨 저장 실패')
  return data.label ?? null
}

export const saveExternalCandidate = async (
  headers: AuthHeaders,
  scenarioId: string,
  payload: {
    lesson_id: string
    title: string
    description?: string | null
    url: string
    source?: string
    external_content_id?: string
    thumbnail_url?: string | null
    author?: string
    language?: string
  },
) => {
  const res = await fetch(`${API_BASE}/api/v1/super-admin/recommendation-debug/scenarios/${scenarioId}/external-candidates/save`, {
    method: 'POST',
    headers: {
      ...headers,
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify(payload),
  })
  return ensureOK<SavedExternalCandidate>(res, '외부 후보 저장 실패')
}

export const fetchLabelDataset = async (headers: AuthHeaders, scenarioId: string, format: 'jsonl' | 'csv') => {
  const res = await fetch(`${API_BASE}/api/v1/super-admin/recommendation-debug/scenarios/${scenarioId}/labels/export?format=${format}`, { headers, credentials: 'include' })
  if (!res.ok) {
    const data = await readJSON<{ error?: string }>(res)
    throw new Error(data.error ?? '라벨 데이터셋 export 실패')
  }
  return res.blob()
}

export const fetchLessonSearchPrerunReports = async (headers: AuthHeaders, limit = 5) => {
  const safeLimit = Math.min(Math.max(Math.trunc(limit), 1), 50)
  const res = await fetch(`${API_BASE}/api/v1/super-admin/recommendation-debug/lesson-search-prerun-reports?limit=${safeLimit}`, { headers, credentials: 'include' })
  return ensureOK<LessonSearchPrerunReportsResponse>(res, '리슨 검색 프리런 report 조회 실패')
}
