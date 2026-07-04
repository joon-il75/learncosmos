'use client'

import type { CSSProperties } from 'react'
import { useState } from 'react'

import LumiModalShell from '@/components/common/LumiModalShell'
import ExplorerTreePanel from '@/components/explorer-plan/ExplorerTreePanel'
import type { ExplorerPlanState } from '@/components/explorer-plan/useExplorerPlan'
import { useExplorerEditForm } from '@/components/explorer-plan/tree/useExplorerEditForm'
import { ExplorerTreeActionBar } from '@/components/explorer-plan/tree/ExplorerTreeActionBar'
import type { PlanningLumiGuideFocus } from './PlanningLumiGuide'
import type { DashboardCourseDraftCopy } from '@/lib/i18n/pages/dashboardCourseDraft'

// ── 위치 상수 — 픽셀 분석 기준 ───────────────────────────────────────────────
// Explorer_Diary_Design.webp (1395×757) — 펼침: 우측 작은 디바이스 화면
// x≈1058~1307  y≈103~540
const OPEN_LEFT   = '75.8%'
const OPEN_TOP    = '13.6%'
const OPEN_WIDTH  = '17.8%'
const OPEN_HEIGHT = '57.8%'
// 디바이스 바텀: 13.6 + 57.8 = 71.4%

// Explorer_Diary_Design_wide.webp (1395×757) — 접힘: 큰 디바이스 화면
// x=313~1212  y=150~645
const CLOSED_LEFT   = '20.4%'
const CLOSED_TOP    = '16.7%'
const CLOSED_WIDTH  = '70.2%'
const CLOSED_HEIGHT = '59.4%'
// 디바이스 바텀: 16.7 + 59.4 = 76.1%

// ── 스타일 ────────────────────────────────────────────────────────────────────
const openShellStyle: CSSProperties = {
  position: 'absolute',
  left: OPEN_LEFT,
  top: OPEN_TOP,
  width: OPEN_WIDTH,
  height: OPEN_HEIGHT,
  zIndex: 3,
  minHeight: 0,
  overflow: 'hidden',
  padding: '6px 4px',
  boxSizing: 'border-box',
}

const closedShellStyle: CSSProperties = {
  position: 'absolute',
  left: CLOSED_LEFT,
  top: CLOSED_TOP,
  width: CLOSED_WIDTH,
  height: CLOSED_HEIGHT,
  zIndex: 3,
  minHeight: 0,
  overflow: 'hidden',
  padding: '14px 16px',
  boxSizing: 'border-box',
}

// 버튼 바 — device shell 바로 아래에 배치
const openActionBarStyle: CSSProperties = {
  position: 'absolute',
  left: '75.4%',
  top: '78.4%',
  width: OPEN_WIDTH,
  zIndex: 3,
  boxSizing: 'border-box',
  padding: '0 2px',
}

const closedActionBarStyle: CSSProperties = {
  position: 'absolute',
  left: CLOSED_LEFT,
  top: '77.5%', // 76.1% + 1.4% 여백
  width: CLOSED_WIDTH,
  zIndex: 3,
  boxSizing: 'border-box',
  padding: '0 4px',
}

// ── 오른쪽 책갈피 인덱스 ─────────────────────────────────────────────────────
// Explorer Diary 페이지 인덱스: 계획 → 일지 → 기록 → 결과물 → 커뮤니티 → 문명
type BookmarkState = 'active' | 'inactive' | 'locked'

type PlannedBookmark = 'community' | 'civilization'

function bookmarkRailStyle(mapExpanded: boolean): CSSProperties {
  return {
    position: 'absolute',
    left: mapExpanded ? '92.9%' : '91.2%',
    top: '12.5%',
    display: 'flex',
    flexDirection: 'column',
    gap: 5,
    zIndex: 6,
  }
}

function bookmarkTabStyle(state: BookmarkState): CSSProperties {
  const active = state === 'active'
  const locked = state === 'locked'

  return {
    appearance: 'none',
    width: 40,
    minHeight: 72,
    padding: '12px 0',
    background: active
      ? 'linear-gradient(135deg, #c9920e 0%, #a06c08 60%, #7a5005 100%)'
      : 'linear-gradient(135deg, rgba(92, 60, 19, 0.90) 0%, rgba(47, 30, 10, 0.94) 100%)',
    color: active ? '#fffaf0' : locked ? 'rgba(248, 223, 168, 0.86)' : 'rgba(255, 227, 170, 0.92)',
    fontSize: 12,
    fontWeight: 900,
    letterSpacing: '0.12em',
    lineHeight: 1.05,
    writingMode: 'vertical-rl',
    textOrientation: 'mixed',
    textAlign: 'center',
    borderRadius: '0 9px 9px 0',
    cursor: active ? 'default' : locked ? 'not-allowed' : 'pointer',
    border: `1px solid ${active ? '#e8b020' : locked ? 'rgba(176, 126, 45, 0.50)' : 'rgba(168, 111, 37, 0.42)'}`,
    borderLeft: 'none',
    opacity: 1,
    boxShadow: active
      ? 'inset 0 1px 0 rgba(255,240,180,0.35), inset 0 -2px 4px rgba(0,0,0,0.35), 3px 2px 10px rgba(0,0,0,0.60), 0 0 12px rgba(200,150,10,0.40)'
      : 'inset 0 1px 0 rgba(255,220,140,0.15), inset 0 -2px 4px rgba(0,0,0,0.40), 3px 2px 10px rgba(0,0,0,0.60)',
    textShadow: active ? '0 1px 3px rgba(0,0,0,0.50)' : '0 1px 2px rgba(0,0,0,0.60)',
    userSelect: 'none',
    transition: 'all 180ms ease',
    marginBottom: 0,
    outline: 'none',
    position: 'relative',
    overflow: 'hidden',
  }
}

const inactiveSlashStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  display: 'grid',
  placeItems: 'center',
  color: 'rgba(226, 232, 240, 0.92)',
  fontSize: 42,
  fontWeight: 900,
  lineHeight: 1,
  writingMode: 'horizontal-tb',
  textOrientation: 'mixed',
  textShadow: '0 2px 6px rgba(2, 6, 23, 0.70)',
  pointerEvents: 'none',
}

// ── 컴포넌트 ──────────────────────────────────────────────────────────────────
interface PlanningRightPanelProps {
  explorerPlan: ExplorerPlanState
  mapExpanded: boolean
  isNarrowViewport: boolean
  readOnly?: boolean
  journalLimitedEdit?: boolean
  activeTab?: 'planning' | 'journal'
  isJournalAvailable?: boolean
  onTabSelect?: (tab: 'planning' | 'journal') => void
  onNodeActivate?: (nodeId: string) => void
  onSelectionFocusChange?: (focus: PlanningLumiGuideFocus) => void
  onAfterSave?: () => void
  copy: DashboardCourseDraftCopy['planning']['rightPanel']
  treeCopy: DashboardCourseDraftCopy['tree']
}

export function PlanningRightPanel({
  explorerPlan,
  mapExpanded,
  isNarrowViewport,
  readOnly = false,
  journalLimitedEdit = false,
  activeTab = 'planning',
  isJournalAvailable = false,
  onTabSelect,
  onNodeActivate,
  onSelectionFocusChange,
  onAfterSave,
  copy,
  treeCopy,
}: PlanningRightPanelProps) {
  const [plannedBookmark, setPlannedBookmark] = useState<PlannedBookmark | null>(null)
  const editForm = useExplorerEditForm(explorerPlan, treeCopy.editPanel.modals)
  const isInactiveCourse = Boolean(explorerPlan.courseAggregate?.is_inactive)
  const bookmarkTabs = copy.bookmarks.map((tab) => {
    if (tab.key === 'planning') {
      const state: BookmarkState = activeTab === 'planning' ? 'active' : 'inactive'
      return { ...tab, state }
    }
    if (isInactiveCourse) {
      return { ...tab, state: 'locked' as BookmarkState }
    }
    if (tab.key === 'journal') {
      const state: BookmarkState = activeTab === 'journal' ? 'active' : isJournalAvailable ? 'inactive' : 'locked'
      return { ...tab, state }
    }
    return { ...tab, state: tab.key === 'civilization' ? 'locked' as BookmarkState : 'inactive' as BookmarkState }
  })

  return (
    <>
      {/* 디바이스 화면 shell */}
      <div style={mapExpanded ? openShellStyle : closedShellStyle}>
        <ExplorerTreePanel
          explorerPlan={explorerPlan}
          editForm={editForm}
          readOnly={readOnly}
          journalLimitedEdit={journalLimitedEdit}
          compact={mapExpanded}
          height="100%"
          onNodeActivate={onNodeActivate}
          onSelectionFocusChange={onSelectionFocusChange}
          hideEditPanel={isInactiveCourse}
          suppressSelectionFocus={isInactiveCourse}
          copy={treeCopy}
        />
      </div>

      {/* 저장/취소/삭제/비활성 버튼 바 — device shell 아래 별도 레이어 */}
      {!readOnly && (
        <div style={mapExpanded ? openActionBarStyle : closedActionBarStyle}>
          <ExplorerTreeActionBar
            editForm={editForm}
            onAfterSave={onAfterSave}
            journalLimitedEdit={journalLimitedEdit}
            inactiveCourseMode={isInactiveCourse}
            copy={treeCopy.actionBar}
          />
        </div>
      )}

      {/* 오른쪽 책갈피 인덱스 */}
      <div style={bookmarkRailStyle(mapExpanded)}>
        {bookmarkTabs.map((tab) => {
          const plannedKey =
            tab.key === 'community'
              ? 'community'
              : tab.key === 'civilization'
                ? 'civilization'
                : null
          const isDisabled = tab.state === 'locked' && plannedKey == null
          return (
            <button
              type="button"
              key={tab.label}
              style={{
                ...bookmarkTabStyle(tab.state),
                ...(plannedKey ? { cursor: 'pointer' } : null),
              }}
              role="tab"
              aria-label={tab.label}
              aria-selected={tab.state === 'active'}
              aria-disabled={tab.state === 'locked'}
              disabled={isDisabled}
              title={isInactiveCourse && tab.key !== 'planning' ? copy.inactiveTitle : tab.state === 'locked' ? copy.lockedTitle(tab.label) : tab.label}
              onClick={() => {
                if (plannedKey) {
                  setPlannedBookmark(plannedKey)
                  return
                }
                if (isInactiveCourse && tab.key !== 'planning') return
                if (!onTabSelect) return
                if (tab.key === 'planning') {
                  onTabSelect('planning')
                  return
                }
                if (tab.key === 'journal' && isJournalAvailable) {
                  onTabSelect('journal')
                }
              }}
            >
              {tab.shortLabel}
              {isInactiveCourse && tab.key !== 'planning' ? (
                <span aria-hidden="true" style={inactiveSlashStyle}>/</span>
              ) : null}
            </button>
          )
        })}
      </div>
      {plannedBookmark ? (
        <LumiModalShell
          title={copy.planned[plannedBookmark].title}
          eyebrow={copy.planned.eyebrow}
          lumiState="curious"
          message={copy.planned[plannedBookmark].message}
          onClose={() => setPlannedBookmark(null)}
          actions={
            <button type="button" style={plannedModalButtonStyle} onClick={() => setPlannedBookmark(null)}>
              {copy.planned.confirm}
            </button>
          }
        />
      ) : null}
    </>
  )
}

const plannedModalButtonStyle: CSSProperties = {
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
