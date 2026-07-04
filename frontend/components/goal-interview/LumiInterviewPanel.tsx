'use client'

import { useEffect, useRef, useState, type CSSProperties, type KeyboardEvent, type ReactNode } from 'react'
import { getDefaultSpriteTuning } from '@/lib/lumi/lumiSpriteMap'
import type { LumiState } from '@/lib/lumi/lumiTypes'
import type { DashboardGoalCopy } from '@/lib/i18n/pages/dashboardGoal'
import { getDashboardGoalCopy } from '@/lib/i18n/pages/dashboardGoal'
import type { InterviewMessage, InterviewState, RebuildDecision } from './useGoalInterview'
import { GoalProposalCard } from './GoalProposalCard'

export type RebuildOption = {
  value: RebuildDecision
  label: string
  desc: string
}

export type LumiInterviewProfile = {
  interview_state: InterviewState | string
  confirmed_goal?: string | null
  messages: InterviewMessage[]
  version?: number
}

interface LumiInterviewPanelProps {
  profile: LumiInterviewProfile | null
  isSending: boolean
  error: string | null
  onSendMessage: (message: string) => void
  onConfirmGoal: (goal: string) => void
  onRetryGoal?: () => void
  onRebuildDecision?: (decision: RebuildDecision) => void
  onCancelRebuildDecision?: () => void
  rebuildOptions?: RebuildOption[]
  rebuildTitle?: string
  rebuildDescription?: string
  onConfirmedAction?: () => void
  confirmedActionLabel?: string
  isConfirmedActionBusy?: boolean
  busyMessage?: string | null
  pointStatusSlot?: ReactNode
  bottomActionSlot?: ReactNode
  copy?: DashboardGoalCopy['interview']
  proposalCopy?: DashboardGoalCopy['proposal']
}

function interviewStateToLumiState(state: InterviewState, isSending: boolean): LumiState {
  if (isSending) return 'thinking'
  switch (state) {
    case 'listening': return 'idle'
    case 'clarifying': return 'curious'
    case 'proposing_goal': return 'happy'
    case 'confirmed': return 'celebrate'
    case 'revising_goal': return 'encourage'
    case 'awaiting_rebuild_decision': return 'curious'
    default: return 'idle'
  }
}

export function LumiInterviewPanel({
  profile,
  isSending,
  error,
  onSendMessage,
  onConfirmGoal,
  onRetryGoal,
  onRebuildDecision,
  onCancelRebuildDecision,
  rebuildOptions,
  rebuildTitle,
  rebuildDescription,
  onConfirmedAction,
  confirmedActionLabel,
  isConfirmedActionBusy = false,
  busyMessage,
  pointStatusSlot,
  bottomActionSlot,
  copy = getDashboardGoalCopy('ko').interview,
  proposalCopy = getDashboardGoalCopy('ko').proposal,
}: LumiInterviewPanelProps) {
  const [input, setInput] = useState('')
  const bottomRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLTextAreaElement>(null)

  const interviewState = (profile?.interview_state ?? 'listening') as InterviewState
  const lumiState = interviewStateToLumiState(interviewState, isSending)
  const tuning = getDefaultSpriteTuning(lumiState)
  const isConfirmed = profile?.interview_state === 'confirmed'
  const isAwaitingRebuild = profile?.interview_state === 'awaiting_rebuild_decision'
  const proposedGoal = profile?.interview_state === 'proposing_goal' ? profile.confirmed_goal : null
  const effectiveRebuildOptions = rebuildOptions ?? [
    { value: 'rebuild_all' as RebuildDecision, label: copy.rebuildOptions.rebuildAll.label, desc: copy.rebuildOptions.rebuildAll.desc },
    { value: 'keep_structure' as RebuildDecision, label: copy.rebuildOptions.keepStructure.label, desc: copy.rebuildOptions.keepStructure.desc },
  ]

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [profile?.messages])

  const handleSend = () => {
    const trimmed = input.trim()
    if (!trimmed || isSending || isConfirmed || isAwaitingRebuild) return
    setInput('')
    onSendMessage(trimmed)
  }

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const handleRetry = () => {
    if (onRetryGoal) {
      onRetryGoal()
      return
    }
    if ((profile?.version ?? 1) > 1 && onCancelRebuildDecision) {
      onCancelRebuildDecision()
      return
    }
    onSendMessage(copy.retryMessage)
  }

  return (
    <div style={panelStyle}>
      {/* Lumi 아바타 */}
      <div style={avatarAreaStyle}>
        <div style={avatarContainerStyle}>
          <div style={lumiAvatarStyle(tuning)} />
          {isSending && <div style={thinkingRingStyle} />}
        </div>
        <span style={lumiNameStyle}>Lumi</span>
        <span style={lumiSubtitleStyle}>{copy.subtitle}</span>
      </div>

      {/* 대화 영역 */}
      <div style={chatAreaStyle}>
        {(profile?.messages ?? []).map((msg, i) => {
          const isLumi = msg.role === 'lumi'
          return (
            <div key={i} style={messageRowStyle(isLumi)}>
              {isLumi && (
                <div style={smallAvatarStyle(tuning)} />
              )}
              <div style={bubbleStyle(isLumi)}>
                {msg.content}
              </div>
            </div>
          )
        })}

        {/* 목표 제안 카드 */}
        {proposedGoal && (
          <div style={{ padding: '0 8px' }}>
            <GoalProposalCard
              proposedGoal={proposedGoal}
              isConfirming={isSending}
              onConfirm={onConfirmGoal}
              onRetry={handleRetry}
              copy={proposalCopy}
            />
          </div>
        )}

        {/* 구조 재구성 선택 카드 */}
        {isAwaitingRebuild && onRebuildDecision && (
          <div style={rebuildCardStyle}>
            <div style={rebuildCardTitleStyle}>
              {rebuildTitle ?? copy.rebuildTitle}
            </div>
            <div style={rebuildCardDescriptionStyle}>
              {rebuildDescription ?? copy.rebuildDescription}
            </div>
            {effectiveRebuildOptions.map((opt) => (
              <button
                className="goalInterviewButton goalInterviewButtonSecondary"
                key={opt.value}
                onClick={() => onRebuildDecision(opt.value)}
                disabled={isSending}
                style={rebuildOptionStyle(isSending)}
              >
                <span style={rebuildOptionLabelStyle}>{opt.label}</span>
                <span style={rebuildOptionDescStyle}>{opt.desc}</span>
              </button>
            ))}
            {onCancelRebuildDecision && (
              <button
                className="goalInterviewButton goalInterviewButtonSecondary"
                type="button"
                onClick={onCancelRebuildDecision}
                disabled={isSending}
                style={rebuildCancelButtonStyle(isSending)}
              >
                {copy.cancel}
              </button>
            )}
          </div>
        )}

        {/* 확정 완료 배너 */}
        {isConfirmed && profile?.confirmed_goal && (
          <div style={confirmedBannerStyle}>
            <div style={{ fontSize: 22, marginBottom: 6 }}>🎉</div>
            <div style={{ fontWeight: 700, fontSize: 15, color: '#1C1400', marginBottom: 4 }}>{copy.confirmedTitle}</div>
            <div style={{ fontSize: 14, color: '#78350F' }}>{profile.confirmed_goal}</div>
            {onConfirmedAction && confirmedActionLabel && (
              <button
                className="goalInterviewButton goalInterviewButtonPrimary"
                onClick={onConfirmedAction}
                disabled={isConfirmedActionBusy}
                style={confirmedActionButtonStyle(isConfirmedActionBusy)}
              >
                {isConfirmedActionBusy ? copy.confirmedActionBusy : confirmedActionLabel}
              </button>
            )}
          </div>
        )}

        {isSending && (
          <div style={messageRowStyle(true)}>
            <div style={smallAvatarStyle(tuning)} />
            <div style={bubbleStyle(true)}>
              {busyMessage ? (
                <span>{busyMessage}</span>
              ) : (
                <span style={typingDotsStyle}>{copy.typing}</span>
              )}
            </div>
          </div>
        )}

        {error && (
          <div style={errorStyle}>{error}</div>
        )}

        <div ref={bottomRef} />
      </div>

      {/* 입력 영역 */}
      {!isConfirmed && !isAwaitingRebuild && (
        <>
          {pointStatusSlot}
          <div style={inputAreaStyle}>
            <textarea
              ref={inputRef}
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder={copy.inputPlaceholder}
              disabled={isSending}
              rows={2}
              style={textareaStyle(isSending)}
            />
            <button
              className="goalInterviewButton goalInterviewButtonPrimary"
              onClick={handleSend}
              disabled={!input.trim() || isSending}
              style={sendButtonStyle(!input.trim() || isSending)}
            >
              {copy.send}
            </button>
          </div>
          {bottomActionSlot ? (
            <div style={bottomActionWrapStyle}>
              {bottomActionSlot}
            </div>
          ) : null}
        </>
      )}
      <style jsx global>{`
        .goalInterviewButton {
          transition:
            transform 140ms ease,
            box-shadow 140ms ease,
            filter 140ms ease,
            background-color 140ms ease,
            border-color 140ms ease;
          user-select: none;
        }

        .goalInterviewButton:not(:disabled):hover {
          transform: translateY(-1px);
          filter: brightness(1.04);
          box-shadow: 0 8px 18px rgba(120, 53, 15, 0.16);
        }

        .goalInterviewButton:not(:disabled):active {
          transform: translateY(1px) scale(0.99);
          filter: brightness(0.96);
          box-shadow: 0 3px 8px rgba(120, 53, 15, 0.12);
        }

        .goalInterviewButtonPrimary:not(:disabled):hover {
          background-color: #d97706 !important;
        }

        .goalInterviewButtonSecondary:not(:disabled):hover {
          border-color: rgba(146, 64, 14, 0.86) !important;
          background-color: #fff8e8 !important;
        }
      `}</style>
    </div>
  )
}

const panelStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  height: '100%',
  background: '#FFFDF7',
  borderRadius: 20,
  overflow: 'hidden',
  boxShadow: '0 8px 32px rgba(120,80,20,0.10)',
}

const avatarAreaStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  padding: '24px 16px 12px',
  background: 'linear-gradient(180deg, #FFF8E7 0%, #FFFDF7 100%)',
  borderBottom: '1px solid rgba(245,158,11,0.15)',
}

const avatarContainerStyle: CSSProperties = {
  position: 'relative',
  width: 100,
  height: 100,
  marginBottom: 6,
}

function lumiAvatarStyle(tuning: ReturnType<typeof getDefaultSpriteTuning>): CSSProperties {
  return {
    width: '100%',
    height: '100%',
    backgroundImage: `url(/images/lumi.webp)`,
    backgroundSize: `${tuning.backgroundSizeX}% ${tuning.backgroundSizeY}%`,
    backgroundPosition: `${tuning.backgroundPositionX}% ${tuning.backgroundPositionY}%`,
    backgroundRepeat: 'no-repeat',
    transition: 'background-position 0.3s ease',
    imageRendering: 'pixelated',
  }
}

const thinkingRingStyle: CSSProperties = {
  position: 'absolute',
  inset: -6,
  borderRadius: '50%',
  border: '3px solid #F59E0B',
  opacity: 0.6,
  animation: 'pulse 1.2s ease-in-out infinite',
}

const lumiNameStyle: CSSProperties = {
  fontSize: 16,
  fontWeight: 800,
  color: '#78350F',
  letterSpacing: '0.02em',
}

const lumiSubtitleStyle: CSSProperties = {
  fontSize: 12,
  color: '#A16207',
  marginTop: 2,
}

const chatAreaStyle: CSSProperties = {
  flex: 1,
  overflowY: 'auto',
  padding: '16px 16px 8px',
  display: 'flex',
  flexDirection: 'column',
  gap: 10,
}

function messageRowStyle(isLumi: boolean): CSSProperties {
  return {
    display: 'flex',
    flexDirection: isLumi ? 'row' : 'row-reverse',
    alignItems: 'flex-end',
    gap: 8,
  }
}

function smallAvatarStyle(tuning: ReturnType<typeof getDefaultSpriteTuning>): CSSProperties {
  return {
    flexShrink: 0,
    width: 32,
    height: 32,
    backgroundImage: `url(/images/lumi.webp)`,
    backgroundSize: `${tuning.backgroundSizeX}% ${tuning.backgroundSizeY}%`,
    backgroundPosition: `${tuning.backgroundPositionX}% ${tuning.backgroundPositionY}%`,
    backgroundRepeat: 'no-repeat',
    imageRendering: 'pixelated',
    borderRadius: '50%',
    border: '1.5px solid rgba(245,158,11,0.3)',
  }
}

function bubbleStyle(isLumi: boolean): CSSProperties {
  return {
    maxWidth: '72%',
    padding: '10px 14px',
    borderRadius: isLumi ? '4px 16px 16px 16px' : '16px 4px 16px 16px',
    background: isLumi ? '#FFF3CD' : '#1A0E03',
    color: isLumi ? '#3D2000' : '#F5ECD7',
    fontSize: 14,
    lineHeight: 1.6,
    boxShadow: isLumi
      ? '0 2px 8px rgba(245,158,11,0.12)'
      : '0 2px 8px rgba(0,0,0,0.20)',
    whiteSpace: 'pre-wrap',
  }
}

const typingDotsStyle: CSSProperties = {
  fontSize: 20,
  letterSpacing: 4,
  color: '#A16207',
}

const confirmedBannerStyle: CSSProperties = {
  background: 'linear-gradient(135deg, #FFF9F0, #FFF3E0)',
  border: '2px solid #F59E0B',
  borderRadius: 14,
  padding: '16px 20px',
  textAlign: 'center',
  margin: '8px 0',
}

function confirmedActionButtonStyle(disabled: boolean): CSSProperties {
  return {
    marginTop: 14,
    padding: '10px 16px',
    background: disabled ? '#D1D5DB' : '#1A0E03',
    color: '#FFF8E7',
    border: 'none',
    borderRadius: 10,
    fontSize: 14,
    fontWeight: 700,
    cursor: disabled ? 'default' : 'pointer',
  }
}

const errorStyle: CSSProperties = {
  fontSize: 13,
  color: '#B91C1C',
  background: '#FEF2F2',
  borderRadius: 8,
  padding: '8px 12px',
  margin: '4px 0',
}

const inputAreaStyle: CSSProperties = {
  display: 'flex',
  gap: 8,
  padding: '12px 16px',
  borderTop: '1px solid rgba(245,158,11,0.15)',
  background: '#FFFDF7',
}

const bottomActionWrapStyle: CSSProperties = {
  display: 'flex',
  justifyContent: 'center',
  padding: '0 16px 14px',
  background: '#FFFDF7',
}

function textareaStyle(disabled: boolean): CSSProperties {
  return {
    flex: 1,
    resize: 'none',
    border: '1.5px solid rgba(245,158,11,0.4)',
    borderRadius: 10,
    padding: '8px 12px',
    fontSize: 14,
    color: '#1C1400',
    background: disabled ? '#FFFBF0' : '#fff',
    outline: 'none',
    fontFamily: 'inherit',
    lineHeight: 1.5,
  }
}

function sendButtonStyle(disabled: boolean): CSSProperties {
  return {
    padding: '0 16px',
    background: disabled ? '#D1D5DB' : '#F59E0B',
    color: '#fff',
    border: 'none',
    borderRadius: 10,
    fontSize: 14,
    fontWeight: 700,
    cursor: disabled ? 'default' : 'pointer',
    flexShrink: 0,
    alignSelf: 'stretch',
  }
}

const rebuildCardStyle: CSSProperties = {
  background: '#FFF7E1',
  border: '1.5px solid rgba(180, 83, 9, 0.42)',
  borderRadius: 14,
  padding: '16px 18px',
  margin: '8px 0',
  display: 'flex',
  flexDirection: 'column',
  gap: 8,
}

const rebuildCardTitleStyle: CSSProperties = {
  fontWeight: 900,
  fontSize: 14,
  color: '#1F1300',
  marginBottom: 4,
}

const rebuildCardDescriptionStyle: CSSProperties = {
  fontSize: 13,
  lineHeight: 1.55,
  color: '#3B2200',
  marginBottom: 12,
}

const rebuildOptionLabelStyle: CSSProperties = {
  fontWeight: 900,
  fontSize: 13,
  color: '#1F1300',
}

const rebuildOptionDescStyle: CSSProperties = {
  fontSize: 12,
  lineHeight: 1.45,
  color: '#4A2A00',
  marginTop: 3,
}

function rebuildOptionStyle(disabled: boolean): CSSProperties {
  return {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'flex-start',
    padding: '11px 14px',
    background: disabled ? '#E5E7EB' : '#FFFDF7',
    border: '1.5px solid rgba(146, 64, 14, 0.58)',
    borderRadius: 10,
    cursor: disabled ? 'default' : 'pointer',
    textAlign: 'left',
    gap: 3,
    transition: 'background 0.15s',
    boxShadow: disabled ? 'none' : '0 1px 0 rgba(120, 53, 15, 0.12)',
  }
}

function rebuildCancelButtonStyle(disabled: boolean): CSSProperties {
  return {
    marginTop: 2,
    padding: '10px 14px',
    background: disabled ? '#E5E7EB' : 'transparent',
    border: '1.5px solid rgba(87, 55, 20, 0.36)',
    borderRadius: 10,
    color: disabled ? '#6B7280' : '#2B1A00',
    fontSize: 13,
    fontWeight: 900,
    cursor: disabled ? 'default' : 'pointer',
    textAlign: 'center',
  }
}
