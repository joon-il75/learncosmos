'use client'

import type { CSSProperties } from 'react'
import LumiModalShell from '@/components/common/LumiModalShell'
import type { PlanetDiaryCopy } from '@/lib/i18n/pages/planetDiary'
import { getStatusSummary } from './diaryStageUtils'
import { getDiaryNodeStatusLabel, type DiaryNode } from './planetDiaryTypes'

export type PlannedBookmark = 'community' | 'civilization'

type PointContext = {
  node: DiaryNode
  regionName: string
  subRegionName: string | null
}

export function DiaryStageModals({
  showExitConfirm,
  showPlanningConfirm,
  planningEditHref,
  pendingPointContext,
  blockedPointContext,
  plannedBookmark,
  onCloseExitConfirm,
  onConfirmExit,
  onClosePlanningConfirm,
  onConfirmPlanning,
  onClosePendingPoint,
  onConfirmPendingPoint,
  onClosePlannedBookmark,
  onCloseBlockedPoint,
  copy,
}: {
  showExitConfirm: boolean
  showPlanningConfirm: boolean
  planningEditHref: string | null
  pendingPointContext: PointContext | null
  blockedPointContext: PointContext | null
  plannedBookmark: PlannedBookmark | null
  onCloseExitConfirm: () => void
  onConfirmExit: () => void
  onClosePlanningConfirm: () => void
  onConfirmPlanning: () => void
  onClosePendingPoint: () => void
  onConfirmPendingPoint: () => void
  onClosePlannedBookmark: () => void
  onCloseBlockedPoint: () => void
  copy: PlanetDiaryCopy
}) {
  const plannedBookmarkInfo: Record<PlannedBookmark, { title: string; message: string }> = {
    community: {
      title: copy.modals.plannedBookmark.communityTitle,
      message: copy.modals.plannedBookmark.communityMessage,
    },
    civilization: {
      title: copy.modals.plannedBookmark.civilizationTitle,
      message: copy.modals.plannedBookmark.civilizationMessage,
    },
  }

  return (
    <>
      {showExitConfirm ? (
        <LumiModalShell
          title={copy.modals.exitTitle}
          eyebrow="Explorer Exit"
          lumiState="planet-hold"
          message={copy.modals.exitMessage}
          onClose={onCloseExitConfirm}
          actions={
            <>
              <button type="button" onClick={onCloseExitConfirm} style={modalSecondaryButtonStyle}>
                {copy.modals.exitCancel}
              </button>
              <button type="button" onClick={onConfirmExit} style={modalPrimaryButtonStyle}>
                {copy.modals.exitConfirm}
              </button>
            </>
          }
        />
      ) : null}
      {showPlanningConfirm && planningEditHref ? (
        <LumiModalShell
          title={copy.modals.planningTitle}
          eyebrow="Journal Bookmark"
          lumiState="curious"
          message={copy.modals.planningMessage}
          onClose={onClosePlanningConfirm}
          actions={
            <>
              <button type="button" onClick={onClosePlanningConfirm} style={modalSecondaryButtonStyle}>
                {copy.modals.cancel}
              </button>
              <button type="button" onClick={onConfirmPlanning} style={modalPrimaryButtonStyle}>
                {copy.modals.confirm}
              </button>
            </>
          }
        />
      ) : null}
      {pendingPointContext ? (
        <LumiModalShell
          title={pendingPointContext.node.title}
          eyebrow={pendingPointContext.node.node_type === 'research' ? 'Research Point' : 'Exploration Point'}
          lumiState="curious"
          message={
            <>
              <div>{copy.modals.pendingPoint.region(pendingPointContext.regionName)}</div>
              {pendingPointContext.subRegionName ? <div>{copy.modals.pendingPoint.subRegion(pendingPointContext.subRegionName)}</div> : null}
              <div>
                {copy.modals.pendingPoint.type(
                  pendingPointContext.node.node_type === 'research'
                    ? copy.modals.pendingPoint.researchType
                    : copy.modals.pendingPoint.explorationType,
                )}
              </div>
              <div>{copy.modals.pendingPoint.status(getDiaryNodeStatusLabel(pendingPointContext.node.learning_status, copy.statusLabels))}</div>
              <div style={{ marginTop: 8 }}>{getStatusSummary(pendingPointContext.node.learning_status, copy.statusSummaries)}</div>
              <div style={{ marginTop: 8 }}>{copy.modals.pendingPoint.prompt}</div>
            </>
          }
          onClose={onClosePendingPoint}
          actions={
            <>
              <button type="button" onClick={onClosePendingPoint} style={modalSecondaryButtonStyle}>
                {copy.modals.cancel}
              </button>
              <button type="button" onClick={onConfirmPendingPoint} style={modalPrimaryButtonStyle}>
                {copy.modals.pendingPoint.action}
              </button>
            </>
          }
        />
      ) : null}
      {plannedBookmark ? (
        <LumiModalShell
          title={plannedBookmarkInfo[plannedBookmark].title}
          eyebrow="Coming Soon"
          lumiState="curious"
          message={plannedBookmarkInfo[plannedBookmark].message}
          onClose={onClosePlannedBookmark}
          actions={
            <button type="button" onClick={onClosePlannedBookmark} style={modalPrimaryButtonStyle}>
              {copy.modals.confirm}
            </button>
          }
        />
      ) : null}
      {blockedPointContext ? (
        <LumiModalShell
          title={copy.modals.blockedPoint.title}
          eyebrow={blockedPointContext.node.node_type === 'research' ? 'Research Point' : 'Exploration Point'}
          lumiState="curious"
          message={
            <>
              <div>{blockedPointContext.node.title}</div>
              <div>{copy.modals.blockedPoint.unsaved}</div>
              <div style={{ marginTop: 8 }}>{copy.modals.blockedPoint.message}</div>
            </>
          }
          onClose={onCloseBlockedPoint}
          actions={
            <button type="button" onClick={onCloseBlockedPoint} style={modalPrimaryButtonStyle}>
              {copy.modals.confirm}
            </button>
          }
        />
      ) : null}
    </>
  )
}

const modalSecondaryButtonStyle: CSSProperties = {
  minHeight: 40,
  padding: '0 14px',
  borderRadius: 999,
  border: '1px solid rgba(122, 90, 36, 0.22)',
  background: 'rgba(255,255,255,0.68)',
  color: '#5F4319',
  fontSize: 13,
  fontWeight: 800,
  cursor: 'pointer',
}

const modalPrimaryButtonStyle: CSSProperties = {
  minHeight: 40,
  padding: '0 16px',
  borderRadius: 999,
  border: '1px solid rgba(173, 120, 36, 0.34)',
  background: 'linear-gradient(180deg, rgba(244, 214, 153, 0.96), rgba(205, 154, 67, 0.96))',
  color: '#4A2C00',
  fontSize: 13,
  fontWeight: 900,
  cursor: 'pointer',
}
