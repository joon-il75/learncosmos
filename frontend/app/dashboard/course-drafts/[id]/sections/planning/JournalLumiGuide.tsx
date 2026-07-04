'use client'

import type { CSSProperties } from 'react'
import LumiAvatar from '@/components/lumi/LumiAvatar'

const containerStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: 8,
  padding: '6px 12px 6px 8px',
  background: 'rgba(6, 10, 18, 0.72)',
  borderRadius: 10,
  color: '#fff',
  fontSize: 13,
  lineHeight: 1.5,
  maxWidth: 420,
}

const messageStyle: CSSProperties = {
  color: 'rgba(220, 230, 255, 0.84)',
  fontSize: 12,
}

interface JournalLumiGuideProps {
  canOpenPoint: boolean
  copy: {
    canOpen: string
    explore: string
  }
}

export function JournalLumiGuide({ canOpenPoint, copy }: JournalLumiGuideProps) {
  return (
    <div style={containerStyle}>
      <LumiAvatar state={canOpenPoint ? 'curious' : 'exploring'} size={28} reducedMotion={false} />
      <span style={messageStyle}>
        {canOpenPoint ? copy.canOpen : copy.explore}
      </span>
    </div>
  )
}
