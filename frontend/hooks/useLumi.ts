'use client';

import { useContext } from 'react';

import { LumiContextObject } from '@/providers/LumiProvider';

export function useLumi() {
  const context = useContext(LumiContextObject);
  if (!context) {
    throw new Error('useLumi must be used within a LumiProvider');
  }
  return context;
}
