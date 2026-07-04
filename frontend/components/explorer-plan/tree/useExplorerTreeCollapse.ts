'use client'

import { useState } from 'react'

export function useExplorerTreeCollapse() {
  const [collapsedRegionIds, setCollapsedRegionIds] = useState<Set<string>>(new Set())
  const [collapsedSubRegionIds, setCollapsedSubRegionIds] = useState<Set<string>>(new Set())

  const toggleRegion = (regionId: string) => {
    setCollapsedRegionIds((current) => {
      const next = new Set(current)
      if (next.has(regionId)) next.delete(regionId)
      else next.add(regionId)
      return next
    })
  }

  const toggleSubRegion = (subRegionId: string) => {
    setCollapsedSubRegionIds((current) => {
      const next = new Set(current)
      if (next.has(subRegionId)) next.delete(subRegionId)
      else next.add(subRegionId)
      return next
    })
  }

  return {
    collapsedRegionIds,
    collapsedSubRegionIds,
    toggleRegion,
    toggleSubRegion,
  }
}
