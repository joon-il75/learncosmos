'use client';

import { useEffect, useRef } from 'react';

import type { PlanetRouteKind } from '../pointPageTypes';
import { learningSessionHeartbeatMs } from './pointLearningUtils';

type UsePointLearningSessionArgs = {
  routeKind: PlanetRouteKind;
  planetID: string;
  pointID: string;
  pointDetailID?: string | null;
};

export function usePointLearningSession({
  routeKind,
  planetID,
  pointID,
  pointDetailID,
}: UsePointLearningSessionArgs) {
  const learningSessionIDRef = useRef<string | null>(null);
  const learningSessionEndedRef = useRef(false);
  const learningSessionActiveStartedAtRef = useRef<number | null>(null);
  const learningSessionPendingActiveMsRef = useRef(0);

  useEffect(() => {
    if (routeKind !== 'learning' || !planetID || !pointID || !pointDetailID) return;
    if (typeof window === 'undefined' || typeof document === 'undefined') return;

    let cancelled = false;
    let heartbeatTimer: number | null = null;
    const handlePageHide = () => endSession('pagehide');
    const handleBeforeUnload = () => endSession('beforeunload');

    learningSessionIDRef.current = null;
    learningSessionEndedRef.current = false;
    learningSessionActiveStartedAtRef.current = document.visibilityState === 'visible' ? Date.now() : null;
    learningSessionPendingActiveMsRef.current = 0;

    const getVisibility = () => (document.visibilityState === 'hidden' ? 'hidden' : 'visible');
    const endpointBase = `/api/v1/planets/learning/${planetID}/points/${pointID}/sessions`;
    const collectActiveMs = () => {
      const activeStartedAt = learningSessionActiveStartedAtRef.current;
      if (activeStartedAt == null) return;
      const now = Date.now();
      if (now > activeStartedAt) learningSessionPendingActiveMsRef.current += now - activeStartedAt;
      learningSessionActiveStartedAtRef.current = document.visibilityState === 'visible' ? now : null;
    };
    const takeActiveSeconds = (mode: 'heartbeat' | 'end') => {
      collectActiveMs();
      const pendingMs = learningSessionPendingActiveMsRef.current;
      const seconds = mode === 'end' ? Math.ceil(pendingMs / 1000) : Math.floor(pendingMs / 1000);
      if (seconds > 0) learningSessionPendingActiveMsRef.current = Math.max(0, pendingMs - seconds * 1000);
      return Math.max(0, seconds);
    };
    const buildPayload = (mode: 'heartbeat' | 'end', reason: string) => ({
      active_seconds_delta: takeActiveSeconds(mode),
      visibility: getVisibility(),
      reason,
    });
    const sendHeartbeat = (reason: string) => {
      const sessionID = learningSessionIDRef.current;
      if (!sessionID || learningSessionEndedRef.current) return;
      void fetch(`${endpointBase}/${sessionID}/heartbeat`, {
        method: 'PATCH',
        credentials: 'include',
        keepalive: true,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(buildPayload('heartbeat', reason)),
      }).catch(() => undefined);
    };
    function endSession(reason: string) {
      const sessionID = learningSessionIDRef.current;
      if (!sessionID || learningSessionEndedRef.current) return;
      learningSessionEndedRef.current = true;
      const body = JSON.stringify(buildPayload('end', reason));
      const url = `${endpointBase}/${sessionID}/end`;
      if (navigator.sendBeacon) {
        const blob = new Blob([body], { type: 'application/json' });
        if (navigator.sendBeacon(url, blob)) return;
      }
      void fetch(url, {
        method: 'POST',
        credentials: 'include',
        keepalive: true,
        headers: { 'Content-Type': 'application/json' },
        body,
      }).catch(() => undefined);
    }
    const handleVisibilityChange = () => {
      if (document.visibilityState === 'hidden') {
        sendHeartbeat('hidden');
        learningSessionActiveStartedAtRef.current = null;
      } else {
        learningSessionActiveStartedAtRef.current = Date.now();
      }
    };
    const startSession = async () => {
      try {
        const res = await fetch(endpointBase, {
          method: 'POST',
          credentials: 'include',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ visibility: getVisibility(), user_agent: navigator.userAgent }),
        });
        const payload = (await res.json().catch(() => ({}))) as { session?: { id?: string } };
        if (cancelled || !res.ok || !payload.session?.id) return;
        learningSessionIDRef.current = payload.session.id;
        heartbeatTimer = window.setInterval(() => sendHeartbeat('interval'), learningSessionHeartbeatMs);
      } catch {
        // 체류 시간 기록 실패는 학습 화면 사용을 막지 않는다.
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);
    window.addEventListener('pagehide', handlePageHide);
    window.addEventListener('beforeunload', handleBeforeUnload);
    void startSession();

    return () => {
      cancelled = true;
      if (heartbeatTimer) window.clearInterval(heartbeatTimer);
      document.removeEventListener('visibilitychange', handleVisibilityChange);
      window.removeEventListener('pagehide', handlePageHide);
      window.removeEventListener('beforeunload', handleBeforeUnload);
      endSession('unmount');
    };
  }, [planetID, pointDetailID, pointID, routeKind]);
}
