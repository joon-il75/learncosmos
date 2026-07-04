'use client'

import type { CSSProperties } from 'react'
import LumiAvatar from '@/components/lumi/LumiAvatar'
import type { DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'

const containerStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'auto 1fr',
  alignItems: 'center',
  gap: 14,
  padding: '16px 18px',
  background: 'linear-gradient(135deg, rgba(47, 31, 13, 0.42), rgba(12, 29, 48, 0.58))',
  border: '1px solid rgba(228, 199, 122, 0.18)',
  borderRadius: 22,
  boxShadow: 'inset 0 1px 0 rgba(255, 238, 201, 0.08)',
  color: '#F7F2E2',
  fontSize: 14,
  lineHeight: 1.7,
}

const avatarWrapStyle: CSSProperties = {
  display: 'grid',
  placeItems: 'center',
  width: 70,
  height: 70,
  borderRadius: 20,
  background: 'rgba(255, 244, 214, 0.10)',
  border: '1px solid rgba(255, 236, 194, 0.16)',
}

const dirtyContainerStyle: CSSProperties = {
  ...containerStyle,
  position: 'relative',
  overflow: 'hidden',
  border: '1px solid rgba(255, 203, 96, 0.58)',
  background: 'linear-gradient(135deg, rgba(91, 49, 10, 0.62), rgba(24, 39, 55, 0.66))',
  boxShadow: 'inset 0 1px 0 rgba(255, 238, 201, 0.16), 0 0 0 1px rgba(255, 208, 106, 0.08), 0 16px 38px rgba(255, 170, 70, 0.18)',
  animation: 'planningLumiSavePulse 1.7s ease-in-out infinite',
}

const inactiveContainerStyle: CSSProperties = {
  ...containerStyle,
  background: 'linear-gradient(135deg, rgba(29, 37, 50, 0.78), rgba(51, 58, 68, 0.70))',
  border: '1px solid rgba(180, 193, 210, 0.38)',
  boxShadow: 'inset 0 1px 0 rgba(236, 242, 248, 0.10), 0 14px 34px rgba(7, 12, 20, 0.18)',
}

const inactiveAvatarWrapStyle: CSSProperties = {
  ...avatarWrapStyle,
  background: 'rgba(226, 232, 240, 0.12)',
  border: '1px solid rgba(203, 213, 225, 0.24)',
  filter: 'grayscale(0.35)',
}

const textWrapStyle: CSSProperties = {
  display: 'grid',
  gap: 5,
  minWidth: 0,
}

const eyebrowStyle: CSSProperties = {
  fontSize: 11,
  fontWeight: 900,
  letterSpacing: '0.12em',
  textTransform: 'uppercase',
  color: '#E8D39A',
}

const dirtyMessageStyle: CSSProperties = {
  margin: 0,
  color: '#FFF6D6',
  fontWeight: 900,
  fontSize: 14,
  lineHeight: 1.7,
}

const dirtyMessageBadgeStyle: CSSProperties = {
  display: 'inline-flex',
  alignItems: 'center',
  width: 'fit-content',
  maxWidth: '100%',
  padding: '4px 10px',
  borderRadius: 999,
  background: 'rgba(255, 202, 84, 0.16)',
  border: '1px solid rgba(255, 215, 124, 0.34)',
  boxShadow: '0 0 18px rgba(255, 191, 73, 0.16)',
  textShadow: '0 1px 8px rgba(116, 62, 0, 0.40)',
  animation: 'planningLumiMessageGlow 1.7s ease-in-out infinite',
}

const messageStyle: CSSProperties = {
  ...dirtyMessageStyle,
  fontWeight: 600,
  color: 'rgba(247, 242, 226, 0.88)',
}

export type PlanningLumiGuideFocus = 'course' | 'region' | 'subregion' | 'exploration-node' | 'research-node' | null

interface PlanningLumiGuideProps {
  hasUnsavedChanges: boolean
  focus?: PlanningLumiGuideFocus
  isInactiveCourse?: boolean
  copy: DashboardCourseDraftCopy['planning']['lumiGuide']
}

export function PlanningLumiGuide({ hasUnsavedChanges, focus = null, isInactiveCourse = false, copy }: PlanningLumiGuideProps) {
  const isCourseFocused = focus === 'course'
  const isRegionFocused = focus === 'region'
  const isSubRegionFocused = focus === 'subregion'
  const isExplorationNodeFocused = focus === 'exploration-node'
  const isResearchNodeFocused = focus === 'research-node'
  const message = isInactiveCourse
    ? copy.inactive
    : hasUnsavedChanges
    ? copy.unsaved
    : isCourseFocused
      ? copy.course
      : isRegionFocused
        ? copy.region
        : isSubRegionFocused
          ? copy.subregion
        : isExplorationNodeFocused
          ? copy.explorationNode
        : isResearchNodeFocused
          ? copy.researchNode
          : copy.default

  return (
    <div style={isInactiveCourse ? inactiveContainerStyle : hasUnsavedChanges ? dirtyContainerStyle : containerStyle}>
      <div style={isInactiveCourse ? inactiveAvatarWrapStyle : avatarWrapStyle}>
        <LumiAvatar state={isInactiveCourse ? 'surprise' : hasUnsavedChanges ? 'thinking' : 'exploring'} size={58} reducedMotion={false} />
      </div>
      <div style={textWrapStyle}>
        <span style={eyebrowStyle}>{copy.eyebrow}</span>
        <p style={hasUnsavedChanges && !isInactiveCourse ? dirtyMessageStyle : messageStyle}>
          {hasUnsavedChanges && !isInactiveCourse ? <span style={dirtyMessageBadgeStyle}>{message}</span> : message}
        </p>
      </div>
      {hasUnsavedChanges && !isInactiveCourse ? (
        <style jsx>{`
          @keyframes planningLumiSavePulse {
            0%, 100% {
              border-color: rgba(255, 203, 96, 0.52);
              box-shadow:
                inset 0 1px 0 rgba(255, 238, 201, 0.16),
                0 0 0 1px rgba(255, 208, 106, 0.08),
                0 16px 38px rgba(255, 170, 70, 0.18);
            }
            50% {
              border-color: rgba(255, 231, 154, 0.92);
              box-shadow:
                inset 0 1px 0 rgba(255, 244, 212, 0.24),
                0 0 0 3px rgba(255, 210, 99, 0.14),
                0 18px 48px rgba(255, 181, 71, 0.30);
            }
          }

          @keyframes planningLumiMessageGlow {
            0%, 100% {
              transform: translateY(0);
              background: rgba(255, 202, 84, 0.16);
            }
            50% {
              transform: translateY(-1px);
              background: rgba(255, 216, 112, 0.25);
            }
          }
        `}</style>
      ) : null}
    </div>
  )
}
