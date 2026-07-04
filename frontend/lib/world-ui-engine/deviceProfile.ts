import type {
  DashboardDeviceProfile,
  DashboardDeviceState,
} from './types';
import {
  DEFAULT_DASHBOARD_UI_ENGINE_SETTINGS,
  type DashboardUIEngineRuntimeSettings,
} from './uiEngineConfig';

export interface DashboardViewportInput {
  width: number;
  height: number;
}

export function buildDashboardDeviceProfile(
  input: DashboardViewportInput,
  settings: DashboardUIEngineRuntimeSettings = DEFAULT_DASHBOARD_UI_ENGINE_SETTINGS,
): DashboardDeviceProfile {
  const isCompactLayout = input.width < settings.compact_breakpoint;
  const isPhoneLayout = input.width < settings.phone_breakpoint;
  const isShortViewport = input.height < settings.short_viewport_height;
  const isSE = input.width < 390 || (isPhoneLayout && isShortViewport);
  const isWideDesktop = input.width >= 1600;
  const isTablet = input.width >= settings.phone_breakpoint && input.width < settings.compact_breakpoint;

  return {
    deviceKind: isWideDesktop
      ? 'wide-desktop'
      : isSE
        ? 'se'
        : isPhoneLayout
          ? 'mobile'
          : isTablet
            ? 'tablet'
            : 'desktop',
    viewportKind: isShortViewport ? 'short' : 'default',
    useCompactMapLabels: isCompactLayout || isShortViewport,
    mapBreakpoint: isCompactLayout || isSE ? 'mobile' : 'desktop',
  };
}

export function buildDashboardDeviceState(
  input: DashboardViewportInput,
  settings: DashboardUIEngineRuntimeSettings = DEFAULT_DASHBOARD_UI_ENGINE_SETTINGS,
): DashboardDeviceState {
  const profile = buildDashboardDeviceProfile(input, settings);

  return {
    viewportWidth: input.width,
    viewportHeight: input.height,
    isCompactLayout: profile.mapBreakpoint === 'mobile',
    isPhoneLayout: profile.deviceKind === 'se' || profile.deviceKind === 'mobile',
    isShortViewport: profile.viewportKind === 'short',
  };
}
