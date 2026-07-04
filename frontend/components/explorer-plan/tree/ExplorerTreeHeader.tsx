'use client'

import { useState, type FormEvent } from 'react'

import {
  countBadgeStyle,
  createRegionFormStyle,
  eyebrowStyle,
  headerStyle,
  headerTopRowStyle,
  mutationMessageStyle,
  titleStyle,
  treeButtonDisabledStyle,
  treeInputStyle,
  treePrimaryButtonStyle,
} from './explorerTreeStyles'

interface ExplorerTreeHeaderProps {
  regionCount: number
  isDirty: boolean
  isMutating: boolean
  mutationMessage: string | null
  onCreateRegion: (name: string) => Promise<void>
  onSavePlanChanges: () => Promise<void>
}

export function ExplorerTreeHeader({
  regionCount,
  isDirty,
  isMutating,
  mutationMessage,
  onCreateRegion,
  onSavePlanChanges,
}: ExplorerTreeHeaderProps) {
  const [regionName, setRegionName] = useState('')
  const trimmedName = regionName.trim()
  const disabled = isMutating || trimmedName.length === 0

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (disabled) return
    await onCreateRegion(trimmedName)
    setRegionName('')
  }

  return (
    <div style={headerStyle}>
      <div style={headerTopRowStyle}>
        <div>
          <div style={eyebrowStyle}>Explorer Tree</div>
          <h3 style={titleStyle}>탐험 경로</h3>
        </div>
        <div style={countBadgeStyle}>{isDirty ? '미저장 변경 있음' : `${regionCount}개 지역`}</div>
      </div>

      <form style={createRegionFormStyle} onSubmit={handleSubmit}>
        <input
          value={regionName}
          onChange={(event) => setRegionName(event.target.value)}
          placeholder="새 지역 이름"
          maxLength={100}
          style={treeInputStyle}
        />
        <button
          type="submit"
          disabled={disabled}
          style={{
            ...treePrimaryButtonStyle,
            ...(disabled ? treeButtonDisabledStyle : null),
          }}
        >
          추가
        </button>
      </form>

      <button
        type="button"
        onClick={onSavePlanChanges}
        disabled={isMutating || !isDirty}
        style={{
          ...treePrimaryButtonStyle,
          width: '100%',
          ...(isMutating || !isDirty ? treeButtonDisabledStyle : null),
        }}
      >
        {isMutating ? '탐험계획 저장 중...' : '탐험계획 저장'}
      </button>

      <div style={mutationMessageStyle}>{mutationMessage ?? ''}</div>
    </div>
  )
}
