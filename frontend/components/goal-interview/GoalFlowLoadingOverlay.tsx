'use client'

import { useEffect, useRef, useState, type CSSProperties } from 'react'
import LumiAvatar from '@/components/lumi/LumiAvatar'
import AtmosphericPlanetPreview from '@/components/dashboard/AtmosphericPlanetPreview'
import { useActivePlanetTextureMap } from '@/components/dashboard/useActivePlanetTextureMap'
import type { DashboardGoalCopy } from '@/lib/i18n/pages/dashboardGoal'
import type { DashboardCourseGenerationWaitEstimate } from '@/lib/world-ui-engine/ctaEngine'
import { getDashboardGoalCopy } from '@/lib/i18n/pages/dashboardGoal'

export type GoalFlowStage =
  | 'bootstrapping_lumi'
  | 'starting_interview'
  | 'creating_planet'
  | 'redirecting'

export interface PlanetTypeItem {
  id: string
  name: string
  asset_path: string
}

export interface PlanetTextureMapItem {
  id: string
  name: string
  public_url: string
  rotation_duration_seconds: number
  rotation_direction: 'left' | 'right'
}

interface GoalFlowLoadingOverlayProps {
  stage: GoalFlowStage
  variant?: 'page' | 'panel'
  copy?: DashboardGoalCopy['loading']
  planetTitle?: string | null
  planetTextureMaps?: PlanetTextureMapItem[]
  selectedPlanetTextureMapId?: string | null
  isReadyToContinue?: boolean
  onSelectPlanetTextureMap?: (textureMap: PlanetTextureMapItem) => void
  onConfirmPlanetTextureMap?: () => void
  courseGenerationWait?: DashboardCourseGenerationWaitEstimate | null
}


function useCreatingPlanetElapsedMS(stage: GoalFlowStage, isReadyToContinue: boolean): number {
  const [elapsedMS, setElapsedMS] = useState(0)

  useEffect(() => {
    if (stage !== 'creating_planet' || isReadyToContinue) {
      setElapsedMS(0)
      return
    }

    const startedAt = Date.now()
    const timer = window.setInterval(() => {
      setElapsedMS(Date.now() - startedAt)
    }, 500)

    return () => window.clearInterval(timer)
  }, [stage, isReadyToContinue])

  return elapsedMS
}

function resolveCreatingPlanetDescription(copy: DashboardGoalCopy['loading'], elapsedMS: number): string {
  if (elapsedMS >= 15000) return copy.creatingPlanetWait.busy
  if (elapsedMS >= 12000) return copy.creatingPlanetWait.busyHint
  if (elapsedMS >= 8000) return copy.creatingPlanetWait.long
  return copy.creatingPlanetWait.normal
}

function resolveWaitEstimateLabel(
  copy: DashboardGoalCopy['loading'],
  estimate: DashboardCourseGenerationWaitEstimate,
): string {
  if (typeof estimate.min_seconds === 'number' && typeof estimate.max_seconds === 'number') {
    return copy.waitEstimate.duration(estimate.min_seconds, estimate.max_seconds)
  }
  return estimate.label || '잠시 후'
}

function resolveCourseGenerationWaitLine(
  copy: DashboardGoalCopy['loading'],
  estimate: DashboardCourseGenerationWaitEstimate | null | undefined,
): string | null {
  if (!estimate) return null
  const label = resolveWaitEstimateLabel(copy, estimate)
  if (typeof estimate.queue_position === 'number') {
    if (estimate.queue_position <= 0) return copy.waitEstimate.running(label)
    return copy.waitEstimate.queued(estimate.queue_position, label)
  }
  if (typeof estimate.active_jobs === 'number') {
    return copy.waitEstimate.activeOnly(estimate.active_jobs, label)
  }
  return null
}

function useSimulatedProgress(stage: GoalFlowStage, isReadyToContinue: boolean): number {
  const [progress, setProgress] = useState(0)
  const rafRef = useRef<number | null>(null)
  const startRef = useRef<number | null>(null)

  useEffect(() => {
    if (stage === 'redirecting' || (stage === 'creating_planet' && isReadyToContinue)) {
      setProgress(100)
      if (rafRef.current !== null) cancelAnimationFrame(rafRef.current)
      return
    }
    if (stage !== 'creating_planet') {
      setProgress(0)
      return
    }
    // 0 → 90% over ~22s with ease-out (asymptote)
    const TARGET = 90
    const DURATION = 22000
    startRef.current = null

    const tick = (now: number) => {
      if (startRef.current === null) startRef.current = now
      const elapsed = now - startRef.current
      // ease-out: 1 - e^(-k*t)
      const k = 3 / DURATION
      const pct = TARGET * (1 - Math.exp(-k * elapsed))
      setProgress(Math.min(pct, TARGET))
      if (pct < TARGET) rafRef.current = requestAnimationFrame(tick)
    }

    rafRef.current = requestAnimationFrame(tick)
    return () => {
      if (rafRef.current !== null) cancelAnimationFrame(rafRef.current)
    }
  }, [stage, isReadyToContinue])

  return progress
}

export function GoalFlowLoadingOverlay({
  stage,
  variant = 'page',
  copy = getDashboardGoalCopy('ko').loading,
  planetTitle = null,
  planetTextureMaps,
  selectedPlanetTextureMapId,
  isReadyToContinue = false,
  onSelectPlanetTextureMap,
  onConfirmPlanetTextureMap,
  courseGenerationWait,
}: GoalFlowLoadingOverlayProps) {
  const stageCopy = copy.stages[stage]
  const creatingElapsedMS = useCreatingPlanetElapsedMS(stage, isReadyToContinue)
  const stageDescription = stage === 'creating_planet'
    ? resolveCreatingPlanetDescription(copy, creatingElapsedMS)
    : stageCopy.description
  const waitLine = stage === 'creating_planet' && !isReadyToContinue
    ? resolveCourseGenerationWaitLine(copy, courseGenerationWait)
    : null
  const activeTextureMap = useActivePlanetTextureMap()
  const selectedTextureMap = planetTextureMaps?.find((item) => item.id === selectedPlanetTextureMapId) ?? null
  const previewTextureMap = selectedTextureMap
    ? {
        atlasURL: selectedTextureMap.public_url,
        rotationDurationSeconds: selectedTextureMap.rotation_duration_seconds,
        rotationDirection: selectedTextureMap.rotation_direction,
      }
    : activeTextureMap
  const showPicker = stage === 'creating_planet' && !!onSelectPlanetTextureMap
  const showProgressBar = stage === 'creating_planet' || stage === 'redirecting'
  const progress = useSimulatedProgress(stage, isReadyToContinue)
  const showAtlasPreview = stage === 'creating_planet' || stage === 'redirecting'

  return (
    <div style={overlayStyle(variant)}>
      <style>{`
        @keyframes goalFlowPulse {
          0%, 100% { transform: scale(0.96); opacity: 0.42; }
          50% { transform: scale(1.04); opacity: 0.9; }
        }
        @keyframes goalFlowDot {
          0%, 100% { transform: translateY(0); opacity: 0.36; }
          50% { transform: translateY(-5px); opacity: 1; }
        }
      `}</style>
      <div style={cardStyle(variant, showPicker)}>
        <div style={avatarWrapStyle}>
          <LumiAvatar state="thinking" size={variant === 'page' ? 84 : 68} />
          <div style={pulseRingStyle} />
        </div>
        <div style={titleStyle(variant)}>{stageCopy.title}</div>
        <div style={descriptionStyle(variant)}>{stageDescription}</div>
        {waitLine ? <div style={waitLineStyle(variant)}>{waitLine}</div> : null}

        {showAtlasPreview ? (
          <div style={atlasPreviewWrapStyle}>
              <AtmosphericPlanetPreview
              atlasURL={previewTextureMap.atlasURL}
              size={variant === 'page' ? 118 : 96}
              progressPercent={Math.min(100, progress)}
              durationSeconds={previewTextureMap.rotationDurationSeconds}
              direction={previewTextureMap.rotationDirection}
              playing
              preset="goal-loading"
            />
            <div style={atlasPreviewCaptionStyle(variant)}>
              {isReadyToContinue && !selectedTextureMap
                ? copy.atlasCaption.readyNoSelection
                : selectedTextureMap
                  ? copy.atlasCaption.selected(selectedTextureMap.name)
                  : stage === 'redirecting'
                    ? copy.atlasCaption.redirecting
                    : copy.atlasCaption.creating}
            </div>
          </div>
        ) : null}

        {showPicker ? (
          <div style={pickerSectionStyle}>
            <div style={pickerLabelStyle}>
              {selectedTextureMap
                ? copy.picker.selected
                : isReadyToContinue
                  ? copy.picker.ready
                  : copy.picker.idle}
            </div>
            {!planetTextureMaps || planetTextureMaps.length === 0 ? (
              <div style={pickerEmptyStyle}>{copy.picker.empty}</div>
            ) : (
              <div style={pickerGridStyle}>
                {planetTextureMaps.map((textureMap) => {
                  const isSelected = textureMap.id === selectedPlanetTextureMapId
                  return (
                    <button
                      key={textureMap.id}
                      type="button"
                      onClick={() => onSelectPlanetTextureMap(textureMap)}
                      style={pickerItemStyle(isSelected)}
                      title={textureMap.name}
                      aria-label={copy.picker.ariaLabel(textureMap.name, isSelected)}
                    >
                      <AtmosphericPlanetPreview
                        atlasURL={textureMap.public_url}
                        size={58}
                        progressPercent={0}
                        durationSeconds={textureMap.rotation_duration_seconds}
                        direction={textureMap.rotation_direction}
                        playing={isSelected}
                        preset="galaxy-mini"
                      />
                      <div style={pickerItemNameStyle(isSelected)}>{textureMap.name}</div>
                    </button>
                  )
                })}
              </div>
            )}
            <button
              type="button"
              onClick={onConfirmPlanetTextureMap}
              disabled={!selectedTextureMap || !isReadyToContinue}
              style={confirmTextureButtonStyle(Boolean(selectedTextureMap && isReadyToContinue))}
            >
              {!selectedTextureMap
                ? copy.picker.chooseTexture
                : isReadyToContinue
                  ? copy.picker.startWithTexture
                  : copy.picker.preparing}
            </button>
          </div>
        ) : null}

        {showProgressBar ? (
          <div style={progressWrapStyle}>
            <div style={progressTrackStyle}>
              <div
                style={{
                  ...progressFillStyle,
                  width: `${progress.toFixed(1)}%`,
                  transition: progress >= 100 ? 'width 0.4s ease-out' : undefined,
                }}
              />
            </div>
            <div style={progressLabelStyle}>
              {progress >= 100 ? copy.complete : `${Math.floor(progress)}%`}
            </div>
          </div>
        ) : (
          <div style={dotsStyle} aria-hidden="true">
            <span style={dotStyle(0)} />
            <span style={dotStyle(0.18)} />
            <span style={dotStyle(0.36)} />
          </div>
        )}
      </div>
    </div>
  )
}

function overlayStyle(variant: 'page' | 'panel'): CSSProperties {
  return {
    position: variant === 'page' ? 'fixed' : 'absolute',
    inset: 0,
    zIndex: variant === 'page' ? 1200 : 30,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    padding: variant === 'page' ? '24px' : '18px',
    background: variant === 'page'
      ? 'rgba(255, 248, 231, 0.9)'
      : 'rgba(255, 253, 247, 0.86)',
    backdropFilter: 'blur(8px)',
    overflowY: 'auto',
  }
}

function cardStyle(variant: 'page' | 'panel', expanded: boolean): CSSProperties {
  return {
    width: '100%',
    maxWidth: expanded ? (variant === 'page' ? 520 : 400) : (variant === 'page' ? 420 : 360),
    padding: variant === 'page' ? '28px 26px' : '22px 20px',
    borderRadius: 24,
    background: 'linear-gradient(180deg, #FFFDF7 0%, #FFF4D8 100%)',
    border: '1px solid rgba(245, 158, 11, 0.25)',
    boxShadow: '0 18px 48px rgba(120, 80, 20, 0.18)',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    textAlign: 'center',
    marginTop: expanded ? 'auto' : undefined,
    marginBottom: expanded ? 'auto' : undefined,
  }
}

const avatarWrapStyle: CSSProperties = {
  position: 'relative',
  width: 92,
  height: 92,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  marginBottom: 14,
}

const pulseRingStyle: CSSProperties = {
  position: 'absolute',
  inset: 2,
  borderRadius: '50%',
  border: '2px solid rgba(245, 158, 11, 0.55)',
  animation: 'goalFlowPulse 1.4s ease-in-out infinite',
}

const atlasPreviewWrapStyle: CSSProperties = {
  width: '100%',
  marginTop: 16,
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  gap: 4,
}

function atlasPreviewCaptionStyle(variant: 'page' | 'panel'): CSSProperties {
  return {
    fontSize: variant === 'page' ? 12 : 11,
    color: '#A16207',
    fontWeight: 800,
    letterSpacing: '0.01em',
  }
}

function titleStyle(variant: 'page' | 'panel'): CSSProperties {
  return {
    fontSize: variant === 'page' ? 20 : 18,
    fontWeight: 800,
    color: '#3D2000',
    marginBottom: 8,
    letterSpacing: '0.01em',
  }
}

function descriptionStyle(variant: 'page' | 'panel'): CSSProperties {
  return {
    fontSize: variant === 'page' ? 14 : 13,
    lineHeight: 1.6,
    color: '#92400E',
    maxWidth: 300,
  }
}

function waitLineStyle(variant: 'page' | 'panel'): CSSProperties {
  return {
    marginTop: 10,
    padding: variant === 'page' ? '8px 12px' : '7px 10px',
    borderRadius: 999,
    background: 'rgba(255, 247, 214, 0.92)',
    border: '1px solid rgba(245, 158, 11, 0.28)',
    color: '#78350F',
    fontSize: variant === 'page' ? 12 : 11,
    fontWeight: 900,
    lineHeight: 1.4,
    maxWidth: 320,
  }
}

const pickerSectionStyle: CSSProperties = {
  width: '100%',
  marginTop: 20,
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  gap: 12,
}

const pickerLabelStyle: CSSProperties = {
  fontSize: 12,
  fontWeight: 700,
  color: '#92400E',
  letterSpacing: '0.02em',
  textAlign: 'center',
}

const pickerEmptyStyle: CSSProperties = {
  fontSize: 12,
  color: '#B45309',
  opacity: 0.7,
}

const pickerGridStyle: CSSProperties = {
  display: 'flex',
  flexWrap: 'wrap',
  gap: 10,
  justifyContent: 'center',
  width: '100%',
}

const confirmTextureButtonStyle = (enabled: boolean): CSSProperties => ({
  width: '100%',
  minHeight: 44,
  borderRadius: 14,
  border: enabled ? '1px solid rgba(37, 99, 235, 0.34)' : '1px solid rgba(120, 100, 70, 0.18)',
  background: enabled
    ? 'linear-gradient(135deg, #2563eb, #14b8a6)'
    : 'rgba(120, 100, 70, 0.12)',
  color: enabled ? '#ffffff' : 'rgba(92, 65, 28, 0.62)',
  font: 'inherit',
  fontSize: 14,
  fontWeight: 900,
  cursor: enabled ? 'pointer' : 'not-allowed',
  boxShadow: enabled ? '0 10px 24px rgba(37, 99, 235, 0.22)' : 'none',
})

function pickerItemStyle(selected: boolean): CSSProperties {
  return {
    background: selected ? 'rgba(245, 158, 11, 0.12)' : 'transparent',
    border: selected ? '2px solid #F59E0B' : '2px solid rgba(245, 158, 11, 0.2)',
    borderRadius: 12,
    padding: '8px 6px 6px',
    cursor: 'pointer',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: 5,
    transition: 'border-color 0.15s, background 0.15s',
    width: 80,
  }
}

function pickerPlanetSphereStyle(selected: boolean): CSSProperties {
  return {
    width: 64,
    height: 64,
    borderRadius: '50%',
    boxShadow: selected
      ? '0 0 0 2px #F59E0B, 0 8px 18px rgba(120, 80, 20, 0.3), inset -8px -10px 14px rgba(40, 20, 0, 0.28)'
      : '0 8px 18px rgba(120, 80, 20, 0.2), inset -8px -10px 14px rgba(40, 20, 0, 0.22)',
    flexShrink: 0,
  }
}

function pickerItemNameStyle(selected: boolean): CSSProperties {
  return {
    fontSize: 10,
    fontWeight: selected ? 800 : 600,
    color: selected ? '#92400E' : '#B45309',
    lineHeight: 1.3,
    maxWidth: 70,
    wordBreak: 'keep-all',
  }
}

const progressWrapStyle: CSSProperties = {
  width: '100%',
  marginTop: 18,
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  gap: 6,
}

const progressTrackStyle: CSSProperties = {
  width: '100%',
  maxWidth: 280,
  height: 8,
  borderRadius: 999,
  background: 'rgba(245, 158, 11, 0.15)',
  overflow: 'hidden',
}

const progressFillStyle: CSSProperties = {
  height: '100%',
  borderRadius: 999,
  background: 'linear-gradient(90deg, #F59E0B 0%, #FBBF24 100%)',
  boxShadow: '0 0 6px rgba(245, 158, 11, 0.5)',
}

const progressLabelStyle: CSSProperties = {
  fontSize: 12,
  fontWeight: 700,
  color: '#A16207',
  letterSpacing: '0.03em',
  minWidth: 36,
  textAlign: 'center',
}

const dotsStyle: CSSProperties = {
  display: 'flex',
  gap: 8,
  marginTop: 16,
}

const planetPreviewWrapStyle: CSSProperties = {
  position: 'relative',
  width: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  marginTop: 18,
}

const planetPreviewHaloStyle: CSSProperties = {
  position: 'absolute',
  top: 4,
  width: 132,
  height: 132,
  borderRadius: '50%',
  background: 'radial-gradient(circle, rgba(245, 158, 11, 0.24) 0%, rgba(245, 158, 11, 0.08) 48%, rgba(245, 158, 11, 0) 78%)',
  filter: 'blur(3px)',
}

const planetPreviewSphereStyle: CSSProperties = {
  position: 'relative',
  width: 116,
  height: 116,
  borderRadius: '50%',
  boxShadow: '0 16px 32px rgba(120, 80, 20, 0.24), inset -14px -16px 22px rgba(40, 20, 0, 0.28)',
  border: '1px solid rgba(255, 255, 255, 0.5)',
}

function planetPreviewCaptionStyle(variant: 'page' | 'panel'): CSSProperties {
  return {
    marginTop: 12,
    fontSize: variant === 'page' ? 12 : 11,
    color: '#A16207',
    fontWeight: 700,
    letterSpacing: '0.01em',
  }
}

function dotStyle(delay: number): CSSProperties {
  return {
    width: 9,
    height: 9,
    borderRadius: '50%',
    background: '#F59E0B',
    opacity: 0.38,
    animation: `goalFlowDot 1.1s ease-in-out ${delay}s infinite`,
  }
}
