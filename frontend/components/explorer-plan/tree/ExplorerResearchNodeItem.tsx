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
  treeItemBlockStyle,
  treeRowStyle,
  subResearchRowStyle,
  researchButtonStyle,
  researchBadgeStyle,
  researchTitleStyle,
  researchMetaStyle,
  selectedButtonStyle,
  treeButtonDisabledStyle,
  treeInputStyle,
} from './explorerTreeStyles'

interface ExplorerResearchNodeItemProps {
  node: ExplorerNode
  selected: boolean
  onSelectNode: (nodeId: string) => void
  onUpdateResearchNode: (nodeId: string, title: string) => Promise<void>
  onDeleteResearchNode: (nodeId: string) => Promise<void>
  isMutating: boolean
  showResearchType?: boolean
}

const researchIconStyle: CSSProperties = {
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

export function ExplorerResearchNodeItem({
  node,
  selected,
  onSelectNode,
  onUpdateResearchNode,
  onDeleteResearchNode,
  isMutating,
  showResearchType = false,
}: ExplorerResearchNodeItemProps) {
  const [isEditing, setIsEditing] = useState(false)
  const [draftTitle, setDraftTitle] = useState(node.title)

  useEffect(() => {
    if (!isEditing) setDraftTitle(node.title)
  }, [isEditing, node.title])

  const trimmedTitle = draftTitle.trim()
  const editDisabled = isMutating || trimmedTitle.length === 0

  const handleSave = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (editDisabled) return
    await onUpdateResearchNode(node.id, trimmedTitle)
    setIsEditing(false)
  }

  const handleCancel = () => {
    setDraftTitle(node.title)
    setIsEditing(false)
  }

  const handleDelete = async () => {
    if (isMutating) return
    const ok = window.confirm(`"${node.title}" 연구지점을 비활성화할까요?`)
    if (!ok) return
    await onDeleteResearchNode(node.id)
  }

  return (
    <div style={showResearchType ? subResearchRowStyle : treeItemBlockStyle}>
      <div style={treeRowStyle}>
        <span style={researchIconStyle}>🔬</span>

        <div style={connectorWrapStyle}>
          {showResearchType ? <div style={connectorStyle} /> : null}

          {isEditing ? (
            <form style={regionEditFormStyle} onSubmit={handleSave}>
              <input
                value={draftTitle}
                onChange={(event) => setDraftTitle(event.target.value)}
                maxLength={200}
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
                <span style={researchBadgeStyle}>연구지점</span>
                <span style={researchTitleStyle}>{node.title}</span>

                {showResearchType && node.research_type ? (
                  <span style={researchMetaStyle}>유형 · {node.research_type}</span>
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
