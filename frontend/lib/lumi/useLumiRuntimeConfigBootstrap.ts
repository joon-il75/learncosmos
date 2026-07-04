'use client';

import { useEffect, useSyncExternalStore } from 'react';

import { getLumiRuntimeConfigVersion, loadPublicLumiRuntimeConfig, subscribeLumiRuntimeConfig } from './lumiRuntimeConfigStore';

export function useLumiRuntimeConfigBootstrap(): number {
  const version = useSyncExternalStore(subscribeLumiRuntimeConfig, getLumiRuntimeConfigVersion, getLumiRuntimeConfigVersion);

  useEffect(() => {
    void loadPublicLumiRuntimeConfig();
  }, []);

  return version;
}
