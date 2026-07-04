'use client'

import { PlanetResultFeedPage } from './PlanetResultFeedPage'

type PlanetRouteKind = 'learning' | 'shared'

export function PlanetResultPage({
  planetId,
  routeKind,
}: {
  planetId: string
  routeKind: PlanetRouteKind
}) {
  return <PlanetResultFeedPage planetId={planetId} routeKind={routeKind} />
}
