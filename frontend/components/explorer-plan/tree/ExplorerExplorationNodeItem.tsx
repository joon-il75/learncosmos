'use client'

import { useEffect, useState, type FormEvent } from 'react'
import type { CSSProperties } from 'react'
import type { ExplorerNode } from '../explorerPlanTypes'
import {
  regionActionButtonStyle,
  regionActionColumnStyle,
  regionDangerButtonStyle,
  regionEditActionRowStyle,
  regionEditFormStyle,
  researchBadgeStyle,
  researchButtonStyle,
  researchMetaStyle,
  researchTitleStyle,
  selectedButtonStyle,
  subResearchRowStyle,
  treeButtonDisabledStyle,
  treeInputStyle,
  treeItemBlockStyle,
  treeRowStyle,
} from './explorerTreeStyles'

interface ExplorerExplorationNodeItemProps {
  node: ExplorerNode
  selected: boolean
  onSelectNode: (nodeId: string) => void
  onUpdateExplorationNode: (nodeId: string, title: string, sourceUrl: string) => Promise<void>
  onDeleteExplorationNode: (nodeId: string) => Promise<void>
  isMutating: boolean
  showSource?: boolean
}

const explorationIconStyle: CSSProperties = {
  flexShrink: 0,
  width: '18px',
  paddingTop: '5px',
  textAlign: 'center',
  fontSize: '12px',
}

const connectorWrapStyle: CSSProperties = {
  position: 'relative',
  flex: 1,
  minWidth: 0,
}

const connectorStyle: CSSProperties = {
  position: 'absolute',
  left: '-10px',
  top: '12px',
  width: '10px',
  height: '1px',
  background: 'rgba(120, 98, 64, 0.25)',
}

const explorationBadgeStyle: CSSProperties = {
  ...researchBadgeStyle,
  background: 'rgba(180, 120, 20, 0.12)',
  color: '#8A5A10',
}

export function ExplorerExplorationNodeItem({
  node,
  selected,
  onSelectNode,
  onUpdateExplorationNode,
  onDeleteExplorationNode,
  isMutating,
  showSource = false,
}: ExplorerExplorationNodeItemProps) {
  const [isEditing, setIsEditing] = useState(false)
  const [draftTitle, setDraftTitle] = useState(node.title)
  const [draftUrl, setDraftUrl] = useState(node.source_url ?? '')

  useEffect(() => {
    if (!isEditing) {
      setDraftTitle(node.title)
      setDraftUrl(node.source_url ?? '')
    }
  }, [isEditing, node.source_url, node.title])

  const trimmedTitle = draftTitle.trim()
  const trimmedUrl = draftUrl.trim()
  const editDisabled = isMutating || trimmedTitle.length === 0 || trimmedUrl.length === 0

  const handleSave = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (editDisabled) return
    await onUpdateExplorationNode(node.id, trimmedTitle, trimmedUrl)
    setIsEditing(false)
  }

  const handleCancel = () => {
    setDraftTitle(node.title)
    setDraftUrl(node.source_url ?? '')
    setIsEditing(false)
  }

  const handleDelete = async () => {
    if (isMutating) return
    const ok = window.confirm(`"${node.title}" 탐험지점을 비활성화할까요?`)
    if (!ok) return
    await onDeleteExplorationNode(node.id)
  }

  return (
    <div style={showSource ? subResearchRowStyle : treeItemBlockStyle}>
      <div style={treeRowStyle}>
        <span style={explorationIconStyle}>◉</span>

        <div style={connectorWrapStyle}>
          {showSource ? <div style={connectorStyle} /> : null}

          {isEditing ? (
            <form style={regionEditFormStyle} onSubmit={handleSave}>
              <input
                value={draftTitle}
                onChange={(event) => setDraftTitle(event.target.value)}
                maxLength={200}
                placeholder="탐험지점 제목"
                style={treeInputStyle}
              />
              <input
                value={draftUrl}
                onChange={(event) => setDraftUrl(event.target.value)}
                placeholder="https://..."
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
                onClick={() => onSelectNode(node.id)}
                style={{
                  ...researchButtonStyle,
                  ...(selected ? selectedButtonStyle : null),
                }}
              >
                <span style={explorationBadgeStyle}>탐험지점</span>
                <span style={researchTitleStyle}>{node.title}</span>
                {showSource && node.source_url ? (
                  <span style={researchMetaStyle}>{node.source_type ?? 'web'} · {node.source_url}</span>
                ) : null}
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
      </div>
    </div>
  )
}
