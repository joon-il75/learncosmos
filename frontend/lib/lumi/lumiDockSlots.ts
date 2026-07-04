import type { LumiDockSlot } from './lumiTypes';

export const lumiDockSlots: readonly LumiDockSlot[] = [
  'dashboard-cta',
  'dashboard-map-top',
  'dashboard-map-bottom',
  'dashboard-panel',
  'diary-header',
  'diary-sidebar',
  'player-inline',
  'player-bottom',
  'planet-scene-anchor',
  'completion-center',
  'mobile-bottom-sheet-trigger',
  'mobile-inline',
] as const;

export function resolveMobileDockSlot(slot: LumiDockSlot): LumiDockSlot {
  if (slot === 'dashboard-cta') return 'mobile-inline';
  if (slot === 'dashboard-map-top' || slot === 'dashboard-map-bottom') return 'mobile-bottom-sheet-trigger';
  if (slot === 'diary-sidebar') return 'mobile-bottom-sheet-trigger';
  if (slot === 'player-bottom') return 'mobile-bottom-sheet-trigger';
  if (slot === 'completion-center') return 'mobile-inline';
  return slot;
}
