import type { LumiQuickAction } from '@/lib/lumi/lumiTypes';

export interface DashboardUIEngineRuntimeSettings {
  compact_breakpoint: number;
  phone_breakpoint: number;
  short_viewport_height: number;
  cta_min_launch_duration_ms: number;
  selection_cta_launch_delay_ms: number;
  lumi_quick_action_limit: number;
  lumi_desktop_panel_enabled: boolean;
  lumi_mobile_sheet_enabled: boolean;
}

export const DEFAULT_DASHBOARD_UI_ENGINE_SETTINGS: DashboardUIEngineRuntimeSettings = {
  compact_breakpoint: 1180,
  phone_breakpoint: 680,
  short_viewport_height: 460,
  cta_min_launch_duration_ms: 920,
  selection_cta_launch_delay_ms: 720,
  lumi_quick_action_limit: 2,
  lumi_desktop_panel_enabled: true,
  lumi_mobile_sheet_enabled: true,
};

export interface DashboardUIEngineRuntimeSettingsResponse {
  settings?: Partial<DashboardUIEngineRuntimeSettings>;
  updated_at?: string;
}

export function normalizeDashboardUIEngineRuntimeSettings(
  settings?: Partial<DashboardUIEngineRuntimeSettings> | null,
): DashboardUIEngineRuntimeSettings {
  return {
    ...DEFAULT_DASHBOARD_UI_ENGINE_SETTINGS,
    ...(settings ?? {}),
  };
}

export function limitDashboardLumiQuickActions(
  actions: LumiQuickAction[],
  limit: number,
): LumiQuickAction[] {
  if (limit <= 0) return [];
  return actions.slice(0, limit);
}
