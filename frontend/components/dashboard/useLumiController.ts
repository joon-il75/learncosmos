'use client';

import { useMemo } from 'react';

import { fromDashboardStateToLumiRuntimeContext } from '@/lib/lumi/lumiContextAdapters';
import type {
  DashboardHoveredSystemSummary,
  DashboardHoverSource,
  LumiRuntimeContext,
} from '@/lib/lumi/lumiEngineTypes';
import {
  resolveLumiRuntime,
  type DashboardLumiActionId,
  type DashboardLumiCourse,
  type DashboardStatusFilter,
} from '@/lib/lumi/lumiRuntimeEngine';
import { lumiSpriteMap } from '@/lib/lumi/lumiSpriteMap';
import { useLumiRuntimeConfigBootstrap } from '@/lib/lumi/useLumiRuntimeConfigBootstrap';
import type {
  LumiDockSlot,
  LumiQuickAction,
  LumiState,
} from '@/lib/lumi/lumiTypes';

export type { LumiDockSlot } from '@/lib/lumi/lumiTypes';

export type PlanetStatus = 'draft' | 'ready' | 'learning' | 'completed';

export interface LumiCourseItem extends DashboardLumiCourse {}

export interface LumiViewState {
  dockSlot: LumiDockSlot;
  mode: LumiState;
  visible: boolean;
  expanded: boolean;
  message: string;
  preview: string;
  sprite: string;
  quickActions: LumiQuickAction[];
  relatedCourseId: string | null;
  runtimeContext: LumiRuntimeContext;
}

interface LumiActionHandlers {
  newCourse: () => void;
  viewRelatedCourse: (courseId: string) => void;
  changeFilter: (status: DashboardStatusFilter) => void;
  clearFilter: () => void;
  openSettings: () => void;
}

interface UseLumiControllerArgs {
  courses: LumiCourseItem[];
  selectedCourseId: string | null;
  hoveredCourseId: string | null;
  activeStatusFilter: DashboardStatusFilter;
  userName: string;
  hasByok: boolean;
  isPhoneLayout: boolean;
  isExpanded: boolean;
  actionHandlers: LumiActionHandlers;
  hoverSource: DashboardHoverSource | null;
  hoveredSystemSummary: DashboardHoveredSystemSummary | null;
}

function spriteForMood(mood: LumiState): string {
  const sprite = lumiSpriteMap[mood];
  return `[${sprite.row},${sprite.col}]`;
}

function viewRelatedCourseLabel(
  relatedCourseId: string | null,
  courses: LumiCourseItem[],
): string {
  const course = relatedCourseId ? courses.find((c) => c.id === relatedCourseId) : null;
  if (course?.status === 'draft') return '이전 버전 보기';
  if (course?.status === 'ready') return '학습 준비 보기';
  if (course?.status === 'learning') return '학습 이어가기';
  return '행성으로 이동';
}

function buildQuickActions(
  actionIds: DashboardLumiActionId[],
  relatedCourseId: string | null,
  courses: LumiCourseItem[],
  handlers: LumiActionHandlers,
): LumiQuickAction[] {
  return actionIds.slice(0, 2).reduce<LumiQuickAction[]>((actions, actionId) => {
    if (actionId === 'view-related-course' && relatedCourseId) {
      actions.push({
        id: actionId,
        label: viewRelatedCourseLabel(relatedCourseId, courses),
        action: () => handlers.viewRelatedCourse(relatedCourseId),
      });
      return actions;
    }
    if (actionId === 'new-course') {
      actions.push({ id: actionId, label: '새 행성 만들기', action: handlers.newCourse });
      return actions;
    }
    if (actionId === 'show-draft') {
      return actions;
    }
    if (actionId === 'show-ready') {
      return actions;
    }
    if (actionId === 'show-learning') {
      actions.push({ id: actionId, label: '학습중만 보기', action: () => handlers.changeFilter('learning') });
      return actions;
    }
    if (actionId === 'show-completed') {
      actions.push({ id: actionId, label: '완료만 보기', action: () => handlers.changeFilter('completed') });
      return actions;
    }
    if (actionId === 'clear-filter') {
      actions.push({ id: actionId, label: '전체 보기', action: handlers.clearFilter });
      return actions;
    }
    if (actionId === 'open-settings') {
      actions.push({ id: actionId, label: '설정 열기', action: handlers.openSettings });
      return actions;
    }
    // PHASE 5 placeholders — 기능 미구현, UI 노출만
    if (actionId === 'share-course') {
      actions.push({ id: actionId, label: '행성 공유하기', action: () => {} });
      return actions;
    }
    if (actionId === 'view-expedition-base') {
      actions.push({ id: actionId, label: '학습기지 보기', action: () => {} });
      return actions;
    }
    return actions;
  }, []);
}


export function useLumiController({
  courses,
  selectedCourseId,
  hoveredCourseId,
  activeStatusFilter,
  userName,
  hasByok,
  isPhoneLayout,
  isExpanded,
  actionHandlers,
  hoverSource,
  hoveredSystemSummary,
}: UseLumiControllerArgs): LumiViewState {
  const lumiConfigVersion = useLumiRuntimeConfigBootstrap();

  return useMemo(() => {
      const runtimeContext = fromDashboardStateToLumiRuntimeContext({
        courses,
        selectedCourseId,
        hoveredCourseId,
        activeStatusFilter,
        hasByok,
        hoverSource,
        hoveredSystemSummary,
      });

    const resolved = resolveLumiRuntime({
      kind: 'dashboard',
      runtimeContext,
      courses,
      selectedCourseId,
      hoveredCourseId,
      activeStatusFilter,
      userName,
      hasByok,
      isMobile: isPhoneLayout,
      expanded: isExpanded,
    });

    return {
      dockSlot: resolved.dockSlot,
      mode: resolved.state,
      visible: resolved.visible,
      expanded: resolved.expanded,
      message: resolved.message,
      preview: resolved.preview,
      sprite: spriteForMood(resolved.state),
      quickActions: buildQuickActions(resolved.dashboardActionIds, resolved.relatedCourseId, courses, actionHandlers),
      relatedCourseId: resolved.relatedCourseId,
      runtimeContext: resolved.runtimeContext,
    };
  }, [activeStatusFilter, actionHandlers, courses, hasByok, hoveredCourseId, isExpanded, isPhoneLayout, selectedCourseId, userName, lumiConfigVersion, hoverSource, hoveredSystemSummary]);
}
