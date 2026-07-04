'use client'

import { emptyStateStyle } from './explorerTreeStyles'

export function EmptyExplorerTreeState({ message }: { message: string }) {
  return <div style={emptyStateStyle}>{message}</div>
}
