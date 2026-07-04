'use client'

import type { CSSProperties } from 'react'
import type { DashboardGoalCopy } from '@/lib/i18n/pages/dashboardGoal'
import { getDashboardGoalCopy } from '@/lib/i18n/pages/dashboardGoal'

interface GoalProposalCardProps {
  proposedGoal: string
  isConfirming: boolean
  onConfirm: (goal: string) => void
  onRetry: () => void
  copy?: DashboardGoalCopy['proposal']
}

export function GoalProposalCard({
  proposedGoal,
  isConfirming,
  onConfirm,
  onRetry,
  copy = getDashboardGoalCopy('ko').proposal,
}: GoalProposalCardProps) {
  return (
    <div style={cardStyle}>
      <div style={badgeStyle}>{copy.badge}</div>
      <p style={goalTextStyle}>{proposedGoal}</p>
      <div style={buttonRowStyle}>
        <button
          className="goalInterviewButton goalInterviewButtonPrimary"
          onClick={() => onConfirm(proposedGoal)}
          disabled={isConfirming}
          style={confirmButtonStyle}
        >
          {isConfirming ? copy.confirming : copy.confirm}
        </button>
        <button
          className="goalInterviewButton goalInterviewButtonSecondary"
          onClick={onRetry}
          disabled={isConfirming}
          style={retryButtonStyle}
        >
          {copy.retry}
        </button>
      </div>
    </div>
  )
}

const cardStyle: CSSProperties = {
  background: 'linear-gradient(135deg, #FFF9F0 0%, #FFF3E0 100%)',
  border: '2px solid #F59E0B',
  borderRadius: 16,
  padding: '20px 24px',
  margin: '8px 0 4px',
  boxShadow: '0 4px 16px rgba(245,158,11,0.15)',
}

const badgeStyle: CSSProperties = {
  fontSize: 13,
  fontWeight: 700,
  color: '#D97706',
  marginBottom: 10,
  letterSpacing: '0.02em',
}

const goalTextStyle: CSSProperties = {
  fontSize: 16,
  fontWeight: 600,
  color: '#1C1400',
  lineHeight: 1.6,
  margin: '0 0 16px',
}

const buttonRowStyle: CSSProperties = {
  display: 'flex',
  gap: 10,
  flexWrap: 'wrap',
}

const confirmButtonStyle: CSSProperties = {
  flex: 1,
  minWidth: 140,
  padding: '10px 16px',
  background: '#F59E0B',
  color: '#fff',
  border: 'none',
  borderRadius: 8,
  fontSize: 14,
  fontWeight: 700,
  cursor: 'pointer',
}

const retryButtonStyle: CSSProperties = {
  padding: '10px 16px',
  background: 'transparent',
  color: '#92400E',
  border: '1.5px solid #D97706',
  borderRadius: 8,
  fontSize: 14,
  fontWeight: 600,
  cursor: 'pointer',
}
