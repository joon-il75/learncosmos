export type SafetyAPIErrorPayload = {
  error?: string;
  error_code?: string | null;
  message?: string | null;
  safety?: {
    action?: string;
    risk_type?: string;
    error_code?: string | null;
    message?: string | null;
  } | null;
};

export function isSafetyInputBlocked(payload: SafetyAPIErrorPayload | null | undefined): boolean {
  const code = payload?.safety?.error_code || payload?.error_code || payload?.error;
  return code === 'safety_input_blocked' || code === 'safety_input_soft_warn';
}

export function resolveSafetyInputMessage(
  payload: SafetyAPIErrorPayload | null | undefined,
  fallback = '이 내용은 학습 목적과 맞지 않을 수 있어요. 학습 목표나 과제 형태로 바꿔서 다시 작성해 주세요.',
): string {
  if (!isSafetyInputBlocked(payload)) return fallback;
  return payload?.safety?.message || payload?.message || fallback;
}
