import { isSafetyInputBlocked, resolveSafetyInputMessage } from '../safetyErrors';

export interface DashboardCourseGenerationWaitEstimate {
  active_jobs?: number;
  queue_position?: number;
  max_active_jobs?: number;
  min_seconds?: number;
  max_seconds?: number;
  label?: string;
}

export interface DashboardDraftCreateResponse {
  job_id?: string;
  status?: DashboardCourseGenerationJobStatus;
  poll_url?: string;
  draft_id?: string;
  course_id?: string;
  draft?: {
    draft?: {
      id: string;
      title?: string;
      planet_type_asset?: string | null;
    };
  };
  error?: string;
  error_code?: string;
  message?: string;
  estimated_wait?: DashboardCourseGenerationWaitEstimate | null;
}

export type DashboardCourseGenerationJobStatus =
  | 'queued'
  | 'running'
  | 'succeeded'
  | 'failed'
  | 'canceled'
  | 'expired';

export interface DashboardCourseGenerationJobResponse {
  job_id?: string;
  status?: DashboardCourseGenerationJobStatus;
  error?: string;
  error_code?: string;
  message?: string;
  estimated_wait?: DashboardCourseGenerationWaitEstimate | null;
  draft_id?: string;
  course_id?: string;
  draft?: DashboardDraftCreateResponse['draft'];
}

export interface DashboardCreateDraftSuccess {
  ok: true;
  draftId: string;
  courseId?: string;
  remainingDelayMs: number;
  draftTitle?: string;
  planetTypeAsset?: string | null;
}

export interface DashboardCreateDraftFailure {
  ok: false;
  message: string;
  reason?: 'busy' | 'retry' | 'insufficient_points' | 'unknown';
}

export type DashboardCreateDraftResult =
  | DashboardCreateDraftSuccess
  | DashboardCreateDraftFailure;

export interface DashboardCTARequestOptions {
  query: string;
  learningGoal?: string;
  startedAt?: number;
  minLaunchDurationMs?: number;
  pollIntervalMs?: number;
  maxPollDurationMs?: number;
  fetchImpl?: typeof fetch;
  onJobAccepted?: (job: {
    jobId: string;
    pollURL?: string;
    status?: DashboardCourseGenerationJobStatus;
    estimatedWait?: DashboardCourseGenerationWaitEstimate | null;
  }) => void;
  onJobProgress?: (job: {
    jobId: string;
    status?: DashboardCourseGenerationJobStatus;
    estimatedWait?: DashboardCourseGenerationWaitEstimate | null;
  }) => void;
}

const DEFAULT_MIN_LAUNCH_DURATION_MS = 920;
const DEFAULT_JOB_POLL_INTERVAL_MS = 1200;
const DEFAULT_JOB_MAX_POLL_DURATION_MS = 120000;

export function normalizeDashboardCTAQuery(query: string): string {
  return query.trim();
}

export function validateDashboardCTAQuery(query: string): DashboardCreateDraftFailure | null {
  if (!normalizeDashboardCTAQuery(query)) {
    return {
      ok: false,
      message: '탐험할 행성 주제를 먼저 입력해 주세요.',
    };
  }

  return null;
}

export function buildDashboardCreateDraftRequestBody(query: string): string {
  return JSON.stringify({
    source_query: normalizeDashboardCTAQuery(query),
  });
}

export function resolveDashboardCreateDraftErrorMessage(
  payload: DashboardDraftCreateResponse | DashboardCourseGenerationJobResponse,
): string {
  const code = payload.error_code || payload.error;
  if (code === 'insufficient_points') {
    return '포인트가 부족해 새 행성탐험계획을 만들 수 없습니다.';
  }

  if ((code === 'draft_generation_busy' || code === 'course_generation_busy')) {
    return '지금은 탐험계획 생성 요청이 많아 잠시 대기 중이에요. 잠시 후 다시 시도해 주세요.';
  }

  if (code === 'draft_generation_retry') {
    return '탐험계획 생성 응답이 불안정해 이번 생성은 마무리되지 않았어요. 잠시 후 다시 시도해 주세요.';
  }

  if (isSafetyInputBlocked(payload)) {
    return resolveSafetyInputMessage(payload);
  }

  return payload.message || payload.error || payload.error_code || '행성탐험계획 생성에 실패했습니다.';
}

export function resolveDashboardCreateDraftErrorReason(
  payload: DashboardDraftCreateResponse | DashboardCourseGenerationJobResponse,
): DashboardCreateDraftFailure['reason'] {
  const code = payload.error_code || payload.error;
  if (code === 'insufficient_points') return 'insufficient_points';
  if (code === 'draft_generation_busy' || code === 'course_generation_busy') return 'busy';
  if (code === 'draft_generation_retry') return 'retry';
  if (isSafetyInputBlocked(payload)) return 'unknown';
  return 'unknown';
}

export function computeDashboardLaunchDelayMs(options: {
  startedAt: number;
  now?: number;
  minLaunchDurationMs?: number;
}): number {
  const now = options.now ?? Date.now();
  const minimum = options.minLaunchDurationMs ?? DEFAULT_MIN_LAUNCH_DURATION_MS;
  const elapsed = now - options.startedAt;
  return Math.max(0, minimum - elapsed);
}

export async function executeGoalInterviewFlow(
  options: DashboardCTARequestOptions,
): Promise<DashboardCreateDraftResult> {
  return { ok: true, draftId: '', remainingDelayMs: 0 };
}


function sleep(ms: number): Promise<void> {
  if (ms <= 0) return Promise.resolve();
  return new Promise((resolve) => globalThis.setTimeout(resolve, ms));
}

function buildCreateDraftSuccessFromPayload(
  payload: DashboardDraftCreateResponse | DashboardCourseGenerationJobResponse,
  options: DashboardCTARequestOptions,
  startedAt: number,
): DashboardCreateDraftResult {
  const draftId = payload.draft?.draft?.id ?? payload.draft_id;
  if (!draftId) {
    return {
      ok: false,
      message: '생성된 행성탐험계획 ID를 확인하지 못했습니다.',
      reason: 'unknown',
    };
  }

  return {
    ok: true,
    draftId,
    courseId: payload.course_id,
    draftTitle: payload.draft?.draft?.title,
    planetTypeAsset: payload.draft?.draft?.planet_type_asset ?? null,
    remainingDelayMs: computeDashboardLaunchDelayMs({
      startedAt,
      minLaunchDurationMs: options.minLaunchDurationMs,
    }),
  };
}

export async function pollDashboardCreateDraftJob(
  options: DashboardCTARequestOptions & { jobId: string; pollURL?: string },
  startedAt: number,
): Promise<DashboardCreateDraftResult> {
  const fetchImpl = options.fetchImpl ?? fetch;
  const pollURL = options.pollURL || `/api/v1/course-generation-jobs/${options.jobId}`;
  const intervalMS = options.pollIntervalMs ?? DEFAULT_JOB_POLL_INTERVAL_MS;
  const maxDurationMS = options.maxPollDurationMs ?? DEFAULT_JOB_MAX_POLL_DURATION_MS;
  const deadline = Date.now() + maxDurationMS;

  while (Date.now() <= deadline) {
    const response = await fetchImpl(pollURL, {
      method: 'GET',
      credentials: 'include',
      cache: 'no-store',
    });
    const payload = (await response
      .json()
      .catch(() => ({}))) as DashboardCourseGenerationJobResponse;

    if (!response.ok) {
      return {
        ok: false,
        message: resolveDashboardCreateDraftErrorMessage(payload),
        reason: resolveDashboardCreateDraftErrorReason(payload),
      };
    }

    options.onJobProgress?.({
      jobId: payload.job_id || options.jobId,
      status: payload.status,
      estimatedWait: payload.estimated_wait ?? null,
    });

    if (payload.status === 'succeeded') {
      return buildCreateDraftSuccessFromPayload(payload, options, startedAt);
    }

    if (payload.status === 'failed' || payload.status === 'canceled' || payload.status === 'expired') {
      return {
        ok: false,
        message: resolveDashboardCreateDraftErrorMessage(payload),
        reason: resolveDashboardCreateDraftErrorReason(payload),
      };
    }

    await sleep(intervalMS);
  }

  return {
    ok: false,
    message: '행성탐험계획 생성 상태 확인 시간이 초과되었습니다. 잠시 후 다시 시도해 주세요.',
    reason: 'unknown',
  };
}

export async function executeDashboardCreateDraft(
  options: DashboardCTARequestOptions,
): Promise<DashboardCreateDraftResult> {
  const fetchImpl = options.fetchImpl ?? fetch;
  const startedAt = options.startedAt ?? Date.now();

  const response = await fetchImpl('/api/v1/course-drafts', {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      source_query: normalizeDashboardCTAQuery(options.query),
      learning_goal: options.learningGoal?.trim() || undefined,
    }),
  });

  const payload = (await response
    .json()
    .catch(() => ({}))) as DashboardDraftCreateResponse;

  if (!response.ok) {
    return {
      ok: false,
      message: resolveDashboardCreateDraftErrorMessage(payload),
      reason: resolveDashboardCreateDraftErrorReason(payload),
    };
  }

  if (response.status === 202 || payload.job_id) {
    if (!payload.job_id) {
      return {
        ok: false,
        message: '행성탐험계획 생성 작업 ID를 확인하지 못했습니다.',
        reason: 'unknown',
      };
    }
    options.onJobAccepted?.({
      jobId: payload.job_id,
      pollURL: payload.poll_url,
      status: payload.status,
      estimatedWait: payload.estimated_wait ?? null,
    });
    options.onJobProgress?.({
      jobId: payload.job_id,
      status: payload.status,
      estimatedWait: payload.estimated_wait ?? null,
    });
    return pollDashboardCreateDraftJob({
      ...options,
      jobId: payload.job_id,
      pollURL: payload.poll_url,
    }, startedAt);
  }

  return buildCreateDraftSuccessFromPayload(payload, options, startedAt);
}
