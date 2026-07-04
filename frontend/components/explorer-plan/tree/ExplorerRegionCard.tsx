'use client'

import { useEffect, useState, type FormEvent } from 'react'

import type { RegionAggregate } from '../explorerPlanTypes'
import { ExplorerExplorationNodeItem } from './ExplorerExplorationNodeItem'
import { ExplorerResearchNodeItem } from './ExplorerResearchNodeItem'
import { ExplorerSubRegionItem } from './ExplorerSubRegionItem'
import {
  childrenWrapStyle,
  collapseButtonStyle,
  createSubRegionFormStyle,
  emptyChildrenStyle,
  regionBadgeStyle,
  regionButtonStyle,
  regionCardStyle,
  regionActionButtonStyle,
  regionActionColumnStyle,
  regionDangerButtonStyle,
  regionEditActionRowStyle,
  regionEditFormStyle,
  regionHeaderRowStyle,
  regionMetaStyle,
  regionTitleStyle,
  regionTopRowStyle,
  selectedButtonStyle,
  subRegionLimitStyle,
  treeButtonDisabledStyle,
  treeInputStyle,
  treePrimaryButtonStyle,
} from './explorerTreeStyles'
import type { TreeChildItem, TreeExplorationNodeItem, TreeResearchNodeItem, TreeSubRegionItem } from './explorerTreeUtils'
import { formatRegionMeta, formatRegionTitle, getRegionCounts } from './explorerTreeUtils'

interface ExplorerRegionCardProps {
  regionAgg: RegionAggregate
  regionIndex: number
  selectedRegionId: string | null
  selectedSubRegionId: string | null
  selectedNodeId: string | null
  collapsed: boolean
  onToggle: (regionId: string) => void
  collapsedSubRegionIds: Set<string>
  onToggleSubRegion: (subRegionId: string) => void
  onSelectRegion: (regionId: string) => void
  onSelectSubRegion: (regionId: string, subRegionId: string) => void
  onSelectNode: (nodeId: string) => void
  onUpdateRegion: (regionId: string, name: string) => Promise<void>
  onDeleteRegion: (regionId: string) => Promise<void>
  onCreateSubRegion: (regionId: string, name: string) => Promise<void>
  onUpdateSubRegion: (regionId: string, subRegionId: string, name: string) => Promise<void>
  onDeleteSubRegion: (subRegionId: string) => Promise<void>
  onCreateResearchNode: (parentKind: 'region' | 'subregion', parentId: string, title: string) => Promise<void>
  onUpdateResearchNode: (nodeId: string, title: string) => Promise<void>
  onDeleteResearchNode: (nodeId: string) => Promise<void>
  onCreateExplorationNode: (parentKind: 'region' | 'subregion', parentId: string, title: string, sourceUrl: string) => Promise<void>
  onUpdateExplorationNode: (nodeId: string, title: string, sourceUrl: string) => Promise<void>
  onDeleteExplorationNode: (nodeId: string) => Promise<void>
  isMutating: boolean
}

export function ExplorerRegionCard({
  regionAgg,
  regionIndex,
  selectedRegionId,
  selectedSubRegionId,
  selectedNodeId,
  collapsed,
  onToggle,
  collapsedSubRegionIds,
  onToggleSubRegion,
  onSelectRegion,
  onSelectSubRegion,
  onSelectNode,
  onUpdateRegion,
  onDeleteRegion,
  onCreateSubRegion,
  onUpdateSubRegion,
  onDeleteSubRegion,
  onCreateResearchNode,
  onUpdateResearchNode,
  onDeleteResearchNode,
  onCreateExplorationNode,
  onUpdateExplorationNode,
  onDeleteExplorationNode,
  isMutating,
}: ExplorerRegionCardProps) {
  const regionId = regionAgg.region.id
  const regionSelected = selectedRegionId === regionId
  const { explorationCount, researchCount } = getRegionCounts(regionAgg)
  const [isEditing, setIsEditing] = useState(false)
  const [draftName, setDraftName] = useState(regionAgg.region.name)
  const [newSubRegionName, setNewSubRegionName] = useState('')
  const [newResearchNodeTitle, setNewResearchNodeTitle] = useState('')
  const [newExplorationNodeTitle, setNewExplorationNodeTitle] = useState('')
  const [newExplorationNodeUrl, setNewExplorationNodeUrl] = useState('')

  useEffect(() => {
    if (!isEditing) setDraftName(regionAgg.region.name)
  }, [isEditing, regionAgg.region.name])

  const directResearchNodes = regionAgg.nodes.filter(n => n.node_type === 'research')
  const directExplorationNodes = regionAgg.nodes.filter(n => n.node_type === 'exploration')

  const regionChildren: TreeChildItem[] = [
    ...regionAgg.subregions.map(
      (subAgg): TreeSubRegionItem => ({ kind: 'subregion', subAgg })
    ),
    ...directExplorationNodes.map(
      (node): TreeExplorationNodeItem => ({ kind: 'exploration-node', node })
    ),
    ...directResearchNodes.map(
      (node): TreeResearchNodeItem => ({ kind: 'research-node', node })
    ),
  ]

  const trimmedDraftName = draftName.trim()
  const editDisabled = isMutating || trimmedDraftName.length === 0
  const subRegionLimitReached = regionAgg.subregions.length >= 3
  const trimmedSubRegionName = newSubRegionName.trim()
  const createSubRegionDisabled = isMutating || subRegionLimitReached || trimmedSubRegionName.length === 0
  const trimmedResearchNodeTitle = newResearchNodeTitle.trim()
  const createResearchNodeDisabled = isMutating || trimmedResearchNodeTitle.length === 0
  const trimmedExplorationNodeTitle = newExplorationNodeTitle.trim()
  const trimmedExplorationNodeUrl = newExplorationNodeUrl.trim()
  const createExplorationNodeDisabled = isMutating || trimmedExplorationNodeTitle.length === 0 || trimmedExplorationNodeUrl.length === 0

  const handleSave = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (editDisabled) return
    await onUpdateRegion(regionId, trimmedDraftName)
    setIsEditing(false)
  }

  const handleCancel = () => {
    setDraftName(regionAgg.region.name)
    setIsEditing(false)
  }

  const handleDelete = async () => {
    if (isMutating) return
    const ok = window.confirm(`"${regionAgg.region.name}" 지역을 비활성화할까요?`)
    if (!ok) return
    await onDeleteRegion(regionId)
  }

  const handleCreateSubRegion = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (createSubRegionDisabled) return
    await onCreateSubRegion(regionId, trimmedSubRegionName)
    setNewSubRegionName('')
  }

  const handleCreateResearchNode = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (createResearchNodeDisabled) return
    await onCreateResearchNode('region', regionId, trimmedResearchNodeTitle)
    setNewResearchNodeTitle('')
  }

  const handleCreateExplorationNode = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (createExplorationNodeDisabled) return
    await onCreateExplorationNode('region', regionId, trimmedExplorationNodeTitle, trimmedExplorationNodeUrl)
    setNewExplorationNodeTitle('')
    setNewExplorationNodeUrl('')
  }

  return (
    <section style={regionCardStyle}>
      <div style={regionHeaderRowStyle}>
        <button
          type="button"
          onClick={() => onToggle(regionId)}
          aria-label={collapsed ? '지역 펼치기' : '지역 접기'}
          style={collapseButtonStyle}
        >
          {collapsed ? '▶' : '▼'}
        </button>

        {isEditing ? (
          <form style={regionEditFormStyle} onSubmit={handleSave}>
            <input
              value={draftName}
              onChange={(event) => setDraftName(event.target.value)}
              maxLength={100}
              style={treeInputStyle}
            />
            <div style={regionEditActionRowStyle}>
              <button
                type="submit"
                disabled={editDisabled}
                style={{
                  ...regionActionButtonStyle,
                  ...(editDisabled ? treeButtonDisabledStyle : null),
                }}
              >
                저장
              </button>
              <button
                type="button"
                onClick={handleCancel}
                disabled={isMutating}
                style={{
                  ...regionActionButtonStyle,
                  ...(isMutating ? treeButtonDisabledStyle : null),
                }}
              >
                취소
              </button>
            </div>
          </form>
        ) : (
          <>
            <button
              type="button"
              onClick={() => onSelectRegion(regionId)}
              style={{
                ...regionButtonStyle,
                ...(regionSelected ? selectedButtonStyle : null),
              }}
            >
              <span style={regionTopRowStyle}>
                <span style={regionBadgeStyle}>지역 {regionIndex + 1}</span>
              </span>

              <span style={regionTitleStyle}>
                {formatRegionTitle(regionAgg.region.name)}
              </span>

              <span style={regionMetaStyle}>
                {formatRegionMeta(explorationCount, researchCount)}
              </span>
            </button>

            <div style={regionActionColumnStyle}>
              <button
                type="button"
                onClick={() => setIsEditing(true)}
                disabled={isMutating}
                style={{
                  ...regionActionButtonStyle,
                  ...(isMutating ? treeButtonDisabledStyle : null),
                }}
              >
                수정
              </button>
              <button
                type="button"
                onClick={handleDelete}
                disabled={isMutating}
                style={{
                  ...regionActionButtonStyle,
                  ...regionDangerButtonStyle,
                  ...(isMutating ? treeButtonDisabledStyle : null),
                }}
              >
                삭제
              </button>
            </div>
          </>
        )}
      </div>

      {!collapsed && (
        <div style={childrenWrapStyle}>
          <form style={createSubRegionFormStyle} onSubmit={handleCreateSubRegion}>
            <input
              value={newSubRegionName}
              onChange={(event) => setNewSubRegionName(event.target.value)}
              maxLength={100}
              placeholder={subRegionLimitReached ? '서브지역은 최대 3개입니다' : '새 서브지역 이름'}
              disabled={isMutating || subRegionLimitReached}
              style={treeInputStyle}
            />
            <button
              type="submit"
              disabled={createSubRegionDisabled}
              style={{
                ...treePrimaryButtonStyle,
                ...(createSubRegionDisabled ? treeButtonDisabledStyle : null),
              }}
            >
              추가
            </button>
          </form>
          {subRegionLimitReached && (
            <div style={subRegionLimitStyle}>지역마다 서브지역은 최대 3개까지 만들 수 있습니다.</div>
          )}

          <form style={createSubRegionFormStyle} onSubmit={handleCreateExplorationNode}>
            <input
              value={newExplorationNodeTitle}
              onChange={(event) => setNewExplorationNodeTitle(event.target.value)}
              maxLength={200}
              placeholder="새 탐험지점 제목"
              disabled={isMutating}
              style={treeInputStyle}
            />
            <input
              value={newExplorationNodeUrl}
              onChange={(event) => setNewExplorationNodeUrl(event.target.value)}
              placeholder="https://..."
              disabled={isMutating}
              style={treeInputStyle}
            />
            <button
              type="submit"
              disabled={createExplorationNodeDisabled}
              style={{
                ...treePrimaryButtonStyle,
                ...(createExplorationNodeDisabled ? treeButtonDisabledStyle : null),
              }}
            >
              탐험 추가
            </button>
          </form>

          <form style={createSubRegionFormStyle} onSubmit={handleCreateResearchNode}>
            <input
              value={newResearchNodeTitle}
              onChange={(event) => setNewResearchNodeTitle(event.target.value)}
              maxLength={200}
              placeholder="새 연구지점 제목"
              disabled={isMutating}
              style={treeInputStyle}
            />
            <button
              type="submit"
              disabled={createResearchNodeDisabled}
              style={{
                ...treePrimaryButtonStyle,
                ...(createResearchNodeDisabled ? treeButtonDisabledStyle : null),
              }}
            >
              연구 추가
            </button>
          </form>

          {regionChildren.length === 0 ? (
            <div style={emptyChildrenStyle}>서브지역이나 연구지점이 아직 없습니다.</div>
          ) : (
            regionChildren.map((child) => {
              if (child.kind === 'subregion') {
                return (
                  <ExplorerSubRegionItem
                    key={child.subAgg.subregion.id}
                    regionId={regionId}
                    subAgg={child.subAgg}
                    selectedSubRegionId={selectedSubRegionId}
                    selectedNodeId={selectedNodeId}
                    collapsedSubRegionIds={collapsedSubRegionIds}
                    onToggleSubRegion={onToggleSubRegion}
                    onSelectSubRegion={onSelectSubRegion}
                    onSelectNode={onSelectNode}
                    onUpdateSubRegion={onUpdateSubRegion}
                    onDeleteSubRegion={onDeleteSubRegion}
                    onCreateResearchNode={onCreateResearchNode}
                    onUpdateResearchNode={onUpdateResearchNode}
                    onDeleteResearchNode={onDeleteResearchNode}
                    onCreateExplorationNode={onCreateExplorationNode}
                    onUpdateExplorationNode={onUpdateExplorationNode}
                    onDeleteExplorationNode={onDeleteExplorationNode}
                    isMutating={isMutating}
                  />
                )
              }

              if (child.kind === 'exploration-node') {
                return (
                  <ExplorerExplorationNodeItem
                    key={child.node.id}
                    node={child.node}
                    selected={selectedNodeId === child.node.id}
                    onSelectNode={onSelectNode}
                    onUpdateExplorationNode={onUpdateExplorationNode}
                    onDeleteExplorationNode={onDeleteExplorationNode}
                    isMutating={isMutating}
                    showSource
                  />
                )
              }

              return (
                <ExplorerResearchNodeItem
                  key={child.node.id}
                  node={child.node}
                  selected={selectedNodeId === child.node.id}
                  onSelectNode={onSelectNode}
                  onUpdateResearchNode={onUpdateResearchNode}
                  onDeleteResearchNode={onDeleteResearchNode}
                  isMutating={isMutating}
                  showResearchType
                />
              )
            })
          )}
        </div>
      )}
    </section>
  )
}
