import { describe, it, expect, vi } from 'vitest';
import {
  buildDashboardCreateDraftRequestBody,
  computeDashboardLaunchDelayMs,
  executeDashboardCreateDraft,
  normalizeDashboardCTAQuery,
  resolveDashboardCreateDraftErrorMessage,
  validateDashboardCTAQuery,
} from '@/lib/world-ui-engine/ctaEngine';

describe('normalizeDashboardCTAQuery', () => {
  it('앞뒤 공백을 제거한다', () => {
    expect(normalizeDashboardCTAQuery('  기타 독학  ')).toBe('기타 독학');
  });
});

describe('validateDashboardCTAQuery', () => {
  it('빈 입력은 에러를 반환한다', () => {
    expect(validateDashboardCTAQuery('   ')).toEqual({
      ok: false,
      message: '탐험할 행성 주제를 먼저 입력해 주세요.',
    });
  });

  it('유효한 입력은 null을 반환한다', () => {
    expect(validateDashboardCTAQuery('기타 독학')).toBeNull();
  });
});

describe('buildDashboardCreateDraftRequestBody', () => {
  it('trim된 source_query를 JSON body로 만든다', () => {
    expect(buildDashboardCreateDraftRequestBody('  기타 독학  ')).toBe(
      JSON.stringify({ source_query: '기타 독학' }),
    );
  });
});

describe('resolveDashboardCreateDraftErrorMessage', () => {
  it('insufficient_points를 사용자 메시지로 변환한다', () => {
    expect(resolveDashboardCreateDraftErrorMessage({ error: 'insufficient_points' })).toBe(
      '포인트가 부족해 새 행성탐험계획을 만들 수 없습니다.',
    );
  });

  it('busy 생성 실패를 혼잡 메시지로 변환한다', () => {
    expect(resolveDashboardCreateDraftErrorMessage({ error: 'course_generation_busy' })).toBe(
      '지금은 탐험계획 생성 요청이 많아 잠시 대기 중이에요. 잠시 후 다시 시도해 주세요.',
    );
  });

  it('retryable 생성 실패를 재시도 메시지로 변환한다', () => {
    expect(resolveDashboardCreateDraftErrorMessage({ error: 'draft_generation_retry' })).toBe(
      '탐험계획 생성 응답이 불안정해 이번 생성은 마무리되지 않았어요. 잠시 후 다시 시도해 주세요.',
    );
  });

  it('알 수 없는 에러는 원문 또는 기본 메시지를 사용한다', () => {
    expect(resolveDashboardCreateDraftErrorMessage({ error: 'server_error' })).toBe('server_error');
    expect(resolveDashboardCreateDraftErrorMessage({})).toBe('행성탐험계획 생성에 실패했습니다.');
  });
});

describe('computeDashboardLaunchDelayMs', () => {
  it('최소 launch 시간보다 덜 지났으면 남은 시간을 반환한다', () => {
    expect(
      computeDashboardLaunchDelayMs({
        startedAt: 1000,
        now: 1500,
        minLaunchDurationMs: 920,
      }),
    ).toBe(420);
  });

  it('이미 최소 시간을 넘겼으면 0을 반환한다', () => {
    expect(
      computeDashboardLaunchDelayMs({
        startedAt: 1000,
        now: 2500,
        minLaunchDurationMs: 920,
      }),
    ).toBe(0);
  });
});

describe('executeDashboardCreateDraft', () => {
  it('성공 시 draft id와 남은 delay를 반환한다', async () => {
    const fetchImpl = vi.fn().mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue({
        draft: {
          draft: {
            id: 'draft-1',
          },
        },
      }),
    });

    const result = await executeDashboardCreateDraft({
      query: '  기타 독학  ',
      startedAt: Date.now(),
      minLaunchDurationMs: 0,
      fetchImpl: fetchImpl as typeof fetch,
    });

    expect(fetchImpl).toHaveBeenCalledWith('/api/v1/course-drafts', {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ source_query: '기타 독학' }),
    });
    expect(result).toMatchObject({
      ok: true,
      draftId: 'draft-1',
      remainingDelayMs: 0,
    });
  });

  it('응답이 실패하면 사용자 메시지를 반환한다', async () => {
    const fetchImpl = vi.fn().mockResolvedValue({
      ok: false,
      json: vi.fn().mockResolvedValue({
        error: 'insufficient_points',
      }),
    });

    await expect(
      executeDashboardCreateDraft({
        query: '기타 독학',
        fetchImpl: fetchImpl as typeof fetch,
      }),
    ).resolves.toEqual({
      ok: false,
      message: '포인트가 부족해 새 행성탐험계획을 만들 수 없습니다.',
      reason: 'insufficient_points',
    });
  });

  it('혼잡 생성 실패면 busy reason을 반환한다', async () => {
    const fetchImpl = vi.fn().mockResolvedValue({
      ok: false,
      json: vi.fn().mockResolvedValue({
        error: 'course_generation_busy',
      }),
    });

    await expect(
      executeDashboardCreateDraft({
        query: '기타 독학',
        fetchImpl: fetchImpl as typeof fetch,
      }),
    ).resolves.toEqual({
      ok: false,
      message: '지금은 탐험계획 생성 요청이 많아 잠시 대기 중이에요. 잠시 후 다시 시도해 주세요.',
      reason: 'busy',
    });
  });

  it('재시도 가능한 생성 실패면 retry reason을 반환한다', async () => {
    const fetchImpl = vi.fn().mockResolvedValue({
      ok: false,
      json: vi.fn().mockResolvedValue({
        error: 'draft_generation_retry',
      }),
    });

    await expect(
      executeDashboardCreateDraft({
        query: '기타 독학',
        fetchImpl: fetchImpl as typeof fetch,
      }),
    ).resolves.toEqual({
      ok: false,
      message: '탐험계획 생성 응답이 불안정해 이번 생성은 마무리되지 않았어요. 잠시 후 다시 시도해 주세요.',
      reason: 'retry',
    });
  });

  it('202 job 응답이면 상태 API를 polling한 뒤 성공 결과를 반환한다', async () => {
    const fetchImpl = vi
      .fn()
      .mockResolvedValueOnce({
        ok: true,
        status: 202,
        json: vi.fn().mockResolvedValue({
          job_id: 'job-1',
          poll_url: '/api/v1/course-generation-jobs/job-1',
        }),
      })
      .mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: vi.fn().mockResolvedValue({
          job_id: 'job-1',
          status: 'running',
        }),
      })
      .mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: vi.fn().mockResolvedValue({
          job_id: 'job-1',
          status: 'succeeded',
          draft_id: 'draft-async-1',
          course_id: 'course-async-1',
        }),
      });

    const result = await executeDashboardCreateDraft({
      query: '기타 독학',
      startedAt: Date.now(),
      minLaunchDurationMs: 0,
      pollIntervalMs: 0,
      fetchImpl: fetchImpl as typeof fetch,
    });

    expect(fetchImpl).toHaveBeenNthCalledWith(2, '/api/v1/course-generation-jobs/job-1', {
      method: 'GET',
      credentials: 'include',
      cache: 'no-store',
    });
    expect(result).toMatchObject({
      ok: true,
      draftId: 'draft-async-1',
      courseId: 'course-async-1',
      remainingDelayMs: 0,
    });
  });

  it('202 job 응답을 받으면 복구용 callback에 job 정보를 전달한다', async () => {
    const onJobAccepted = vi.fn();
    const fetchImpl = vi
      .fn()
      .mockResolvedValueOnce({
        ok: true,
        status: 202,
        json: vi.fn().mockResolvedValue({
          job_id: 'job-1',
          status: 'queued',
          poll_url: '/api/v1/course-generation-jobs/job-1',
        }),
      })
      .mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: vi.fn().mockResolvedValue({
          job_id: 'job-1',
          status: 'succeeded',
          draft_id: 'draft-async-1',
        }),
      });

    await executeDashboardCreateDraft({
      query: '기타 독학',
      minLaunchDurationMs: 0,
      pollIntervalMs: 0,
      fetchImpl: fetchImpl as typeof fetch,
      onJobAccepted,
    });

    expect(onJobAccepted).toHaveBeenCalledWith({
      jobId: 'job-1',
      pollURL: '/api/v1/course-generation-jobs/job-1',
      status: 'queued',
    });
  });

  it('job 실패 상태는 error_code 기준으로 실패 reason을 반환한다', async () => {
    const fetchImpl = vi
      .fn()
      .mockResolvedValueOnce({
        ok: true,
        status: 202,
        json: vi.fn().mockResolvedValue({ job_id: 'job-1' }),
      })
      .mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: vi.fn().mockResolvedValue({
          job_id: 'job-1',
          status: 'failed',
          error_code: 'course_generation_busy',
        }),
      });

    await expect(
      executeDashboardCreateDraft({
        query: '기타 독학',
        pollIntervalMs: 0,
        fetchImpl: fetchImpl as typeof fetch,
      }),
    ).resolves.toEqual({
      ok: false,
      message: '지금은 탐험계획 생성 요청이 많아 잠시 대기 중이에요. 잠시 후 다시 시도해 주세요.',
      reason: 'busy',
    });
  });

  it('202 응답에 job id가 없으면 실패로 처리한다', async () => {
    const fetchImpl = vi.fn().mockResolvedValue({
      ok: true,
      status: 202,
      json: vi.fn().mockResolvedValue({ status: 'queued' }),
    });

    await expect(
      executeDashboardCreateDraft({
        query: '기타 독학',
        fetchImpl: fetchImpl as typeof fetch,
      }),
    ).resolves.toEqual({
      ok: false,
      message: '행성탐험계획 생성 작업 ID를 확인하지 못했습니다.',
      reason: 'unknown',
    });
  });

  it('poll API 실패는 payload error_code 기준으로 실패 처리한다', async () => {
    const fetchImpl = vi
      .fn()
      .mockResolvedValueOnce({
        ok: true,
        status: 202,
        json: vi.fn().mockResolvedValue({ job_id: 'job-1' }),
      })
      .mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: vi.fn().mockResolvedValue({ error_code: 'llm_job_not_found' }),
      });

    await expect(
      executeDashboardCreateDraft({
        query: '기타 독학',
        pollIntervalMs: 0,
        fetchImpl: fetchImpl as typeof fetch,
      }),
    ).resolves.toEqual({
      ok: false,
      message: 'llm_job_not_found',
      reason: 'unknown',
    });
  });

  it('expired job 상태는 실패로 처리한다', async () => {
    const fetchImpl = vi
      .fn()
      .mockResolvedValueOnce({
        ok: true,
        status: 202,
        json: vi.fn().mockResolvedValue({ job_id: 'job-1' }),
      })
      .mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: vi.fn().mockResolvedValue({
          job_id: 'job-1',
          status: 'expired',
          message: '코스 초안 생성 요청이 만료되었습니다.',
        }),
      });

    await expect(
      executeDashboardCreateDraft({
        query: '기타 독학',
        pollIntervalMs: 0,
        fetchImpl: fetchImpl as typeof fetch,
      }),
    ).resolves.toEqual({
      ok: false,
      message: '코스 초안 생성 요청이 만료되었습니다.',
      reason: 'unknown',
    });
  });

  it('polling 제한 시간을 넘기면 timeout 실패로 처리한다', async () => {
    const fetchImpl = vi.fn().mockResolvedValueOnce({
      ok: true,
      status: 202,
      json: vi.fn().mockResolvedValue({ job_id: 'job-1' }),
    });

    await expect(
      executeDashboardCreateDraft({
        query: '기타 독학',
        maxPollDurationMs: -1,
        fetchImpl: fetchImpl as typeof fetch,
      }),
    ).resolves.toEqual({
      ok: false,
      message: '행성탐험계획 생성 상태 확인 시간이 초과되었습니다. 잠시 후 다시 시도해 주세요.',
      reason: 'unknown',
    });
  });

  it('성공 응답에 draft id가 없으면 실패로 처리한다', async () => {
    const fetchImpl = vi.fn().mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue({
        draft: {},
      }),
    });

    await expect(
      executeDashboardCreateDraft({
        query: '기타 독학',
        fetchImpl: fetchImpl as typeof fetch,
      }),
    ).resolves.toEqual({
      ok: false,
      message: '생성된 행성탐험계획 ID를 확인하지 못했습니다.',
      reason: 'unknown',
    });
  });

  it('json 파싱이 실패해도 기본 실패 메시지로 처리한다', async () => {
    const fetchImpl = vi.fn().mockResolvedValue({
      ok: false,
      json: vi.fn().mockRejectedValue(new Error('invalid json')),
    });

    await expect(
      executeDashboardCreateDraft({
        query: '기타 독학',
        fetchImpl: fetchImpl as typeof fetch,
      }),
    ).resolves.toEqual({
      ok: false,
      message: '행성탐험계획 생성에 실패했습니다.',
      reason: 'unknown',
    });
  });
});
