'use client'

import { useEffect, useState, type FormEvent } from 'react'

import type { SubRegionAggregate } from '../explorerPlanTypes'
import { ExplorerExplorationNodeItem } from './ExplorerExplorationNodeItem'
import { ExplorerResearchNodeItem } from './ExplorerResearchNodeItem'
import {
  regionActionButtonStyle,
  regionActionColumnStyle,
  regionDangerButtonStyle,
  regionEditActionRowStyle,
  regionEditFormStyle,
  createSubRegionFormStyle,
  smallCollapseButtonStyle,
  subRegionButtonStyle,
  subRegionTopRowStyle,
  subRegionBadgeStyle,
  subRegionTitleStyle,
  subRegionMetaStyle,
  subRegionInnerWrapStyle,
  treeItemBlockStyle,
  treeRowStyle,
  selectedButtonStyle,
  treeButtonDisabledStyle,
  treeInputStyle,
  treePrimaryButtonStyle,
} from './explorerTreeStyles'
import { formatSubRegionMeta, getSubRegionCounts } from './explorerTreeUtils'

interface ExplorerSubRegionItemProps {
  regionId: string
  subAgg: SubRegionAggregate
  selectedSubRegionId: string | null
  selectedNodeId: string | null
  collapsedSubRegionIds: Set<string>
  onToggleSubRegion: (subRegionId: string) => void
  onSelectSubRegion: (regionId: string, subRegionId: string) => void
  onSelectNode: (nodeId: string) => void
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

export function ExplorerSubRegionItem({
  regionId,
  subAgg,
  selectedSubRegionId,
  selectedNodeId,
  collapsedSubRegionIds,
  onToggleSubRegion,
  onSelectSubRegion,
  onSelectNode,
  onUpdateSubRegion,
  onDeleteSubRegion,
  onCreateResearchNode,
  onUpdateResearchNode,
  onDeleteResearchNode,
  onCreateExplorationNode,
  onUpdateExplorationNode,
  onDeleteExplorationNode,
  isMutating,
}: ExplorerSubRegionItemProps) {
  const subId = subAgg.subregion.id
  const subCollapsed = collapsedSubRegionIds.has(subId)
  const subSelected = selectedSubRegionId === subId
  const explorationNodes = subAgg.nodes.filter(n => n.node_type === 'exploration')
  const researchNodes = subAgg.nodes.filter(n => n.node_type === 'research')
  const { explorationCount, researchCount } = getSubRegionCounts(subAgg)
  const [isEditing, setIsEditing] = useState(false)
  const [draftName, setDraftName] = useState(subAgg.subregion.name)
  const [newResearchNodeTitle, setNewResearchNodeTitle] = useState('')
  const [newExplorationNodeTitle, setNewExplorationNodeTitle] = useState('')
  const [newExplorationNodeUrl, setNewExplorationNodeUrl] = useState('')

  useEffect(() => {
    if (!isEditing) setDraftName(subAgg.subregion.name)
  }, [isEditing, subAgg.subregion.name])

  const trimmedDraftName = draftName.trim()
  const editDisabled = isMutating || trimmedDraftName.length === 0
  const trimmedResearchNodeTitle = newResearchNodeTitle.trim()
  const createResearchNodeDisabled = isMutating || trimmedResearchNodeTitle.length === 0
  const trimmedExplorationNodeTitle = newExplorationNodeTitle.trim()
  const trimmedExplorationNodeUrl = newExplorationNodeUrl.trim()
  const createExplorationNodeDisabled = isMutating || trimmedExplorationNodeTitle.length === 0 || trimmedExplorationNodeUrl.length === 0

  const handleSave = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (editDisabled) return
    await onUpdateSubRegion(regionId, subId, trimmedDraftName)
    setIsEditing(false)
  }

  const handleCancel = () => {
    setDraftName(subAgg.subregion.name)
    setIsEditing(false)
  }

  const handleDelete = async () => {
    if (isMutating) return
    const ok = window.confirm(`"${subAgg.subregion.name}" 서브지역을 비활성화할까요?`)
    if (!ok) return
    await onDeleteSubRegion(subId)
  }

  const handleCreateResearchNode = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (createResearchNodeDisabled) return
    await onCreateResearchNode('subregion', subId, trimmedResearchNodeTitle)
    setNewResearchNodeTitle('')
  }

  const handleCreateExplorationNode = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (createExplorationNodeDisabled) return
    await onCreateExplorationNode('subregion', subId, trimmedExplorationNodeTitle, trimmedExplorationNodeUrl)
    setNewExplorationNodeTitle('')
    setNewExplorationNodeUrl('')
  }

  return (
    <div style={treeItemBlockStyle}>
      <div style={treeRowStyle}>
        <button
          type="button"
          onClick={() => onToggleSubRegion(subId)}
          style={smallCollapseButtonStyle}
        >
          {subCollapsed ? '▶' : '▼'}
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
              onClick={() => onSelectSubRegion(regionId, subId)}
              style={{
                ...subRegionButtonStyle,
                ...(subSelected ? selectedButtonStyle : null),
              }}
            >
              <span style={subRegionTopRowStyle}>
                <span>🗂️</span>
                <span style={subRegionBadgeStyle}>서브지역</span>
              </span>

              <span style={subRegionTitleStyle}>{subAgg.subregion.name}</span>
              <span style={subRegionMetaStyle}>{formatSubRegionMeta(explorationCount, researchCount)}</span>
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

      {!subCollapsed && (
        <div style={subRegionInnerWrapStyle}>
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

          {explorationNodes.map((node) => (
            <ExplorerExplorationNodeItem
              key={node.id}
              node={node}
              selected={selectedNodeId === node.id}
              onSelectNode={onSelectNode}
              onUpdateExplorationNode={onUpdateExplorationNode}
              onDeleteExplorationNode={onDeleteExplorationNode}
              isMutating={isMutating}
              showSource
            />
          ))}

          {researchNodes.map((node) => (
            <ExplorerResearchNodeItem
              key={node.id}
              node={node}
              selected={selectedNodeId === node.id}
              onSelectNode={onSelectNode}
              onUpdateResearchNode={onUpdateResearchNode}
              onDeleteResearchNode={onDeleteResearchNode}
              isMutating={isMutating}
              showResearchType
            />
          ))}
        </div>
      )}
    </div>
  )
}
