'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import {
  buildDashboardCourseRecordFromPlanet,
  planetStatusSortWeight,
} from '@/lib/world-ui-engine/stateNormalizer';
import type {
  UserInfo,
  UserAISettings,
  CourseItem,
  PlanetListItem,
  TodayTask,
  TodayTaskResponse,
} from './_types';

interface DashboardLoaderResult {
  user: UserInfo | null;
  aiSettings: UserAISettings | null;
  courses: CourseItem[];
  todayTask: TodayTask | null;
  isLoading: boolean;
}

export function useDashboardLoader(): DashboardLoaderResult {
  const router = useRouter();
  const [user, setUser] = useState<UserInfo | null>(null);
  const [aiSettings, setAISettings] = useState<UserAISettings | null>(null);
  const [courses, setCourses] = useState<CourseItem[]>([]);
  const [todayTask, setTodayTask] = useState<TodayTask | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      const refreshRes = await fetch('/api/v1/auth/refresh', {
        method: 'POST',
        credentials: 'include',
      });
      if (!refreshRes.ok) {
        router.push('/login?redirect_after=/dashboard');
        return;
      }

      const meRes = await fetch('/api/v1/auth/me', {
        credentials: 'include',
        cache: 'no-store',
      });
      if (!meRes.ok) {
        router.push('/login?redirect_after=/dashboard');
        return;
      }

      const meData = (await meRes.json()) as UserInfo;
      if (meData.language_setup_required) {
        const setupPath = meData.ui_locale === 'en' ? '/en/language-setup' : '/language-setup';
        router.replace(`${setupPath}?redirect_after=/dashboard`);
        return;
      }
      if (meData.required_consent_pending) {
        const agreementsPath = meData.ui_locale === 'en' ? '/en/agreements' : '/agreements';
        router.replace(`${agreementsPath}?redirect_after=/dashboard`);
        return;
      }
      setUser(meData);

      const aiRes = await fetch('/api/v1/users/me/ai-settings', {
        credentials: 'include',
        cache: 'no-store',
      });
      if (!aiRes.ok) throw new Error('BYOK 설정을 불러오지 못했습니다.');
      const aiData = (await aiRes.json()) as UserAISettings;
      setAISettings(aiData);

      const [learningRes, sharedRes, todayTaskRes] = await Promise.all([
        fetch('/api/v1/planets/learning', {
          credentials: 'include',
          cache: 'no-store',
        }),
        fetch('/api/v1/planets/shared', {
          credentials: 'include',
          cache: 'no-store',
        }),
        fetch('/api/v1/planets/today-task', {
          credentials: 'include',
          cache: 'no-store',
        }),
      ]);

      if (!learningRes.ok || !sharedRes.ok) {
        throw new Error('행성 API 목록을 불러오지 못했습니다.');
      }

      const learningPayload = (await learningRes.json()) as { planets?: PlanetListItem[] };
      const sharedPayload = (await sharedRes.json()) as { planets?: PlanetListItem[] };
      if (todayTaskRes.ok) {
        const todayTaskPayload = (await todayTaskRes.json()) as TodayTaskResponse;
        setTodayTask(todayTaskPayload.task ?? null);
      } else {
        setTodayTask(null);
      }

      const activePlanets = [...(learningPayload.planets ?? []), ...(sharedPayload.planets ?? [])];

      const planetSummaries = activePlanets.map((planet) =>
        buildDashboardCourseRecordFromPlanet(
          planet,
          0,
          planet.lesson_count,
          planet.completed_lesson_count,
          planet.progress ?? null,
        ),
      );

      setCourses(
        planetSummaries.sort((a, b) => {
          const aLastAccessed = a.lastAccessedAt ? new Date(a.lastAccessedAt).getTime() : 0;
          const bLastAccessed = b.lastAccessedAt ? new Date(b.lastAccessedAt).getTime() : 0;
          if (aLastAccessed !== bLastAccessed) {
            return bLastAccessed - aLastAccessed;
          }
          const diff = planetStatusSortWeight(a.status) - planetStatusSortWeight(b.status);
          if (diff !== 0) return diff;
          return new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime();
        }),
      );
      setIsLoading(false);
    };

    load().catch(() => router.push('/login?redirect_after=/dashboard'));
  }, [router]);

  return { user, aiSettings, courses, todayTask, isLoading };
}
