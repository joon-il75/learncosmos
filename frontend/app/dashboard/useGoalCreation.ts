'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import type { GoalFlowStage } from '@/components/goal-interview/GoalFlowLoadingOverlay';
import {
  normalizeDashboardCTAQuery,
  validateDashboardCTAQuery,
} from '@/lib/world-ui-engine/ctaEngine';

interface UseGoalCreationOptions {
  queryInputFocusFn?: () => void;
  errorMessages?: {
    required: string;
    createFailed: string;
  };
}

export interface UseGoalCreationReturn {
  isGoalSubmitting: boolean;
  goalFlowStage: GoalFlowStage | null;
  goalCreationError: string | null;
  handleGoalSubmit: (query: string) => Promise<void>;
  clearGoalError: () => void;
}

export function useGoalCreation({ queryInputFocusFn, errorMessages }: UseGoalCreationOptions): UseGoalCreationReturn {
  const router = useRouter();
  const [isGoalSubmitting, setIsGoalSubmitting] = useState(false);
  const [goalFlowStage, setGoalFlowStage] = useState<GoalFlowStage | null>(null);
  const [goalCreationError, setGoalCreationError] = useState<string | null>(null);

  const handleGoalSubmit = async (query: string) => {
    if (isGoalSubmitting) return;

    const err = validateDashboardCTAQuery(query);
    if (err) {
      setGoalCreationError(errorMessages?.required ?? err.message);
      queryInputFocusFn?.();
      return;
    }

    setGoalCreationError(null);
    setIsGoalSubmitting(true);
    setGoalFlowStage('bootstrapping_lumi');
    try {
      const intent = encodeURIComponent(normalizeDashboardCTAQuery(query));
      setGoalFlowStage('starting_interview');
      router.push(`/dashboard/goal?intent=${intent}`);
    } catch (e) {
      setGoalCreationError(e instanceof Error ? e.message : errorMessages?.createFailed ?? '행성탐험계획 생성에 실패했습니다.');
      setGoalFlowStage(null);
    } finally {
      setIsGoalSubmitting(false);
    }
  };

  return {
    isGoalSubmitting,
    goalFlowStage,
    goalCreationError,
    handleGoalSubmit,
    clearGoalError: () => {
      setGoalCreationError(null);
      setGoalFlowStage(null);
    },
  };
}
