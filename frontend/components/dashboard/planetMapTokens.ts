export type PlanetStatus = 'draft' | 'ready' | 'learning' | 'completed';
export type PlanetStatusCanonical = 'draft' | 'ready' | 'exploring' | 'completed';

export const dashboardTokens = {
  colors: {
    planetDraft:      '#94A3B8',
    planetReady:      '#3B82F6',  // 파란색 — 학습준비중
    planetExploring:  '#22C55E',  // 녹색  — 학습중 (learning)
    planetCompleted:  '#F59E0B',
    panelGlassBg: 'rgba(8, 15, 32, 0.58)',
    panelDarkBg:  'rgba(10, 18, 36, 0.82)',
    panelBorder:  'rgba(148, 163, 184, 0.18)',
    panelGlow:    'rgba(96, 165, 250, 0.18)',
    textPrimary:   '#E5EEFF',
    textSecondary: '#A8B4C7',
    textMuted:     '#7E8AA0',
  },
  radius: {
    lg: '24px',
    md: '18px',
    sm: '12px',
    chip: '999px',
  },
  shadow: {
    panel:         '0 20px 60px rgba(0,0,0,0.24)',
    hover:         '0 24px 80px rgba(55, 125, 255, 0.18)',
    glowPlanet:    '0 0 40px rgba(34, 197, 94, 0.28)',   // 녹색 glow
    glowCompleted: '0 0 48px rgba(245, 158, 11, 0.22)',
  },
  duration: {
    fast:      '160ms',
    base:      '220ms',
    slow:      '420ms',
    floatLoop: '8s',
  },
} as const;

export const statusLabel: Record<PlanetStatus, string> = {
  draft:     '이전 버전',
  ready:     '준비중',
  learning:  '학습중',
  completed: '완료',
};

export const statusDescription: Record<PlanetStatus, string> = {
  draft:     '이전 버전 계획 데이터',
  ready:     '학습 준비 완료',
  learning:  '행성 학습 진행중',
  completed: '행성 학습 완료',
};

export const statusColor: Record<PlanetStatus, string> = {
  draft:     dashboardTokens.colors.planetDraft,
  ready:     dashboardTokens.colors.planetReady,
  learning:  dashboardTokens.colors.planetExploring,
  completed: dashboardTokens.colors.planetCompleted,
};

export function normalizePlanetStatus(raw: PlanetStatus): PlanetStatusCanonical {
  if (raw === 'learning') return 'exploring';
  return raw;
}

export function canonicalStatusColor(status: PlanetStatusCanonical): string {
  if (status === 'draft')     return dashboardTokens.colors.planetDraft;
  if (status === 'ready')     return dashboardTokens.colors.planetReady;
  if (status === 'exploring') return dashboardTokens.colors.planetExploring;
  return dashboardTokens.colors.planetCompleted;
}

export function canonicalStatusSurface(status: PlanetStatusCanonical): string {
  if (status === 'draft') {
    return 'radial-gradient(circle at 35% 32%, rgba(255,255,255,0.28), rgba(148,163,184,0.82) 22%, rgba(84,97,118,0.96) 60%, rgba(34,43,57,0.98) 100%)';
  }
  if (status === 'ready') {
    // 파란색 계열
    return 'radial-gradient(circle at 35% 32%, rgba(239,246,255,0.4), rgba(96,165,250,0.9) 22%, rgba(37,99,235,0.98) 60%, rgba(15,40,100,1) 100%)';
  }
  if (status === 'exploring') {
    // 녹색 계열 — 학습중
    return 'radial-gradient(circle at 35% 32%, rgba(240,255,246,0.4), rgba(74,222,128,0.9) 22%, rgba(22,163,74,0.98) 60%, rgba(10,57,30,1) 100%)';
  }
  return 'radial-gradient(circle at 35% 32%, rgba(255,249,230,0.45), rgba(255,195,94,0.96) 18%, rgba(223,141,23,0.98) 56%, rgba(99,58,10,1) 100%)';
}

export function canonicalStatusGlow(status: PlanetStatusCanonical): string {
  if (status === 'draft')     return '0 0 26px rgba(148, 163, 184, 0.18)';
  if (status === 'ready')     return '0 0 34px rgba(59, 130, 246, 0.26)';
  if (status === 'exploring') return dashboardTokens.shadow.glowPlanet;
  return dashboardTokens.shadow.glowCompleted;
}

export function canonicalStatusRing(status: PlanetStatusCanonical): string {
  if (status === 'draft')     return 'rgba(218, 226, 237, 0.18)';
  if (status === 'ready')     return 'rgba(96, 165, 250, 0.46)';
  if (status === 'exploring') return 'rgba(74, 222, 128, 0.56)';
  return 'rgba(255, 220, 134, 0.62)';
}
