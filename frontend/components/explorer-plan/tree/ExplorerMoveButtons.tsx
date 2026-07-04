'use client'

import {
  treeButtonDisabledStyle,
  addNodeButtonStyle,
} from './explorerTreeStyles'

// ── 위치 조정 버튼 ────────────────────────────────────────────────────────────

export interface MoveButtonsProps {
  canMoveUp: boolean
  canMoveDown: boolean
  isMutating: boolean
  onMoveUp: () => void
  onMoveDown: () => void
  moveUpLabel?: string
  moveDownLabel?: string
}

export function MoveButtons({
  canMoveUp,
  canMoveDown,
  isMutating,
  onMoveUp,
  onMoveDown,
  moveUpLabel = '한 칸 위로 이동',
  moveDownLabel = '한 칸 아래로 이동',
}: MoveButtonsProps) {
  return (
    <div style={{ display: 'flex', gap: 5, flexWrap: 'wrap' }}>
      <button
        type="button"
        style={{ ...addNodeButtonStyle, ...(!canMoveUp || isMutating ? treeButtonDisabledStyle : undefined) }}
        onClick={onMoveUp}
        disabled={!canMoveUp || isMutating}
        aria-label={moveUpLabel}
      >
        ▲
      </button>
      <button
        type="button"
        style={{ ...addNodeButtonStyle, ...(!canMoveDown || isMutating ? treeButtonDisabledStyle : undefined) }}
        onClick={onMoveDown}
        disabled={!canMoveDown || isMutating}
        aria-label={moveDownLabel}
      >
        ▼
      </button>
    </div>
  )
}
