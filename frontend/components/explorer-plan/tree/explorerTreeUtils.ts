import type { ExplorerNode, RegionAggregate, SubRegionAggregate } from '../explorerPlanTypes'

export type TreeResearchNodeItem = {
  kind: 'research-node'
  node: ExplorerNode
}

export type TreeExplorationNodeItem = {
  kind: 'exploration-node'
  node: ExplorerNode
}

export type TreeSubRegionItem = {
  kind: 'subregion'
  subAgg: SubRegionAggregate
}

export type TreeChildItem = TreeSubRegionItem | TreeResearchNodeItem | TreeExplorationNodeItem

export function formatRegionTitle(name: string) {
  return '"' + name + '" 지역'
}

export function formatRegionMeta(explorationCount: number, researchCount: number) {
  const explorationLabel = explorationCount > 0 ? `탐험 ${explorationCount}` : '탐험 없음'
  const researchLabel = researchCount > 0 ? `연구 ${researchCount}` : '연구 없음'
  return `${explorationLabel} · ${researchLabel}`
}

export function formatSubRegionMeta(explorationCount: number, researchCount: number) {
  const explorationLabel = explorationCount > 0 ? `탐험 ${explorationCount}` : '탐험 없음'
  const researchLabel = researchCount > 0 ? `연구 ${researchCount}` : '연구 없음'
  return `${explorationLabel} · ${researchLabel}`
}

export function getRegionCounts(regionAgg: RegionAggregate) {
  const allNodes: ExplorerNode[] = [
    ...regionAgg.nodes,
    ...regionAgg.subregions.flatMap(s => s.nodes),
  ]
  return {
    explorationCount: allNodes.filter(n => n.node_type === 'exploration').length,
    researchCount: allNodes.filter(n => n.node_type === 'research').length,
  }
}

export function getSubRegionCounts(subAgg: SubRegionAggregate) {
  return {
    explorationCount: subAgg.nodes.filter(n => n.node_type === 'exploration').length,
    researchCount: subAgg.nodes.filter(n => n.node_type === 'research').length,
  }
}
