'use client';

import { useEffect, useState } from 'react';
import { PLANET_TEXTURE_MAP_FALLBACK_PATH } from '@/components/dashboard/planetSpriteSystem';

interface ActivePlanetTextureMapResponse {
  planet_texture_map?: {
    public_url?: string;
    rotation_duration_seconds?: number;
    rotation_direction?: 'left' | 'right';
    version?: number;
  } | null;
}

export interface ActivePlanetTextureMap {
  atlasURL: string;
  rotationDurationSeconds: number;
  rotationDirection: 'left' | 'right';
}

export const fallbackTextureMap: ActivePlanetTextureMap = {
  atlasURL: PLANET_TEXTURE_MAP_FALLBACK_PATH,
  rotationDurationSeconds: 36,
  rotationDirection: 'left',
};

export function useActivePlanetTextureMap(): ActivePlanetTextureMap {
  const [textureMap, setTextureMap] = useState<ActivePlanetTextureMap>(fallbackTextureMap);

  useEffect(() => {
    let cancelled = false;

    fetch('/api/v1/public/planet-texture-maps/active', { cache: 'no-store' })
      .then((response) => response.ok ? response.json() : null)
      .then((payload: ActivePlanetTextureMapResponse | null) => {
        if (cancelled || !payload?.planet_texture_map?.public_url) return;
        const asset = payload.planet_texture_map;
        const versionSuffix = asset.version ? `?v=${asset.version}` : '';
        setTextureMap({
          atlasURL: `${asset.public_url}${versionSuffix}`,
          rotationDurationSeconds: asset.rotation_duration_seconds ?? 36,
          rotationDirection: asset.rotation_direction ?? 'left',
        });
      })
      .catch(() => undefined);

    return () => {
      cancelled = true;
    };
  }, []);

  return textureMap;
}
