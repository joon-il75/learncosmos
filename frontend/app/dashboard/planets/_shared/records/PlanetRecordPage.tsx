'use client'

import { PlanetRecordFeedPage } from './PlanetRecordFeedPage'

type PlanetRouteKind = 'learning' | 'shared'

export function PlanetRecordPage({
  planetId,
  routeKind,
}: {
  planetId: string
  routeKind: PlanetRouteKind
}) {
  return <PlanetRecordFeedPage planetId={planetId} routeKind={routeKind} />
}
