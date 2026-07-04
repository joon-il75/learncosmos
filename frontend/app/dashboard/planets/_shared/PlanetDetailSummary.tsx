'use client';

import type { ReactNode } from 'react';
import LumiAvatar from '@/components/lumi/LumiAvatar';
import type { PlanetDiaryCopy } from '@/lib/i18n/pages/planetDiary';
import type {
  DiaryContextDraft,
  PlanetAggregate,
  PlanetRouteKind,
  PlanetStatus,
} from './planetDetailUtils';

export function DiaryContextBand({
  planet,
  draftContext,
  pointProgress,
  lessonProgress,
  actionSlot,
  copy,
}: {
  routeKind: PlanetRouteKind;
  planet: PlanetAggregate;
  draftContext: DiaryContextDraft | null;
  pointProgress: { total: number; completed: number };
  lessonProgress: { total: number; completed: number; learning: number };
  actionSlot?: ReactNode;
  copy: PlanetDiaryCopy['summary'];
}) {
  const completedLessonCount = lessonProgress.completed;
  const totalLessonCount = lessonProgress.total;
  const lessonCompletionPercent =
    totalLessonCount > 0 ? Math.round((completedLessonCount / totalLessonCount) * 100) : 0;
  const goalText =
    planet.goal_context?.confirmed_goal?.trim() ||
    planet.goal_context?.learning_goal?.trim() ||
    draftContext?.draft.source_query?.trim() ||
    draftContext?.draft.title?.trim() ||
    planet.planet.title.trim() ||
    copy.fallbackGoal;

  return (
    <div style={journalContextBandStyle}>
      <div style={journalContextHeaderStyle}>
        <div style={{ display: 'grid', gap: '6px', minWidth: 0 }}>
          <div style={journalContextEyebrowStyle}>{copy.goalEyebrow}</div>
          <div style={journalContextTitleStyle}>{copy.goalTitle}</div>
        </div>
      </div>

      <p style={journalContextBodyStyle}>{goalText}</p>

      <div style={journalProgressWrapStyle}>
        <div style={journalProgressHeaderStyle}>
          <span>{copy.completionRate(lessonCompletionPercent)}</span>
          <span>{copy.completedLessons(completedLessonCount, totalLessonCount)}</span>
          <span>{copy.completedPoints(pointProgress.completed, pointProgress.total)}</span>
        </div>
        <div style={journalProgressTrackStyle}>
          <div
            style={{
              ...journalProgressFillStyle,
              width: `${lessonCompletionPercent}%`,
            }}
          />
        </div>
        {actionSlot ? <div style={journalProgressActionStyle}>{actionSlot}</div> : null}
      </div>
    </div>
  );
}

export function LearningDiaryProgressAction({
  planetStatus,
  canCompletePlanet,
  isStartingPlanet,
  isCompletingPlanet,
  onStart,
  onOpenCompleteConfirm,
  copy,
}: {
  planetStatus: PlanetStatus;
  canCompletePlanet: boolean;
  isStartingPlanet: boolean;
  isCompletingPlanet: boolean;
  onStart: () => void;
  onOpenCompleteConfirm: () => void;
  copy: PlanetDiaryCopy['summary']['progressAction'];
}) {
  return (
    <div style={progressActionWrapStyle}>
      {planetStatus === 'ready' ? (
        <>
          <span style={startReadyBadgeStyle}>{copy.readyBadge}</span>
          <button
            type="button"
            onClick={onStart}
            disabled={isStartingPlanet}
            style={{
              ...startButtonStyle,
              ...(isStartingPlanet ? startButtonDisabledStyle : null),
            }}
          >
            {isStartingPlanet ? copy.starting : copy.start}
          </button>
        </>
      ) : (
        <>
          <span style={startReadyBadgeStyle}>{copy.learningBadge}</span>
          <button
            type="button"
            onClick={onOpenCompleteConfirm}
            disabled={!canCompletePlanet || isCompletingPlanet}
            style={canCompletePlanet ? completeButtonStyle : completionBlockedButtonStyle}
          >
            {canCompletePlanet ? copy.complete : copy.completeBlocked}
          </button>
        </>
      )}
    </div>
  );
}

export function DiaryLumiGuideCard({
  routeKind,
  planetStatus,
  pointProgress,
  canCompletePlanet,
  copy,
}: {
  routeKind: PlanetRouteKind;
  planetStatus: PlanetStatus;
  pointProgress: { total: number; completed: number };
  canCompletePlanet: boolean;
  copy: PlanetDiaryCopy['summary']['guide'];
}) {
  const remaining = Math.max(0, pointProgress.total - pointProgress.completed);
  const message =
    routeKind === 'shared'
      ? copy.shared
      : planetStatus === 'ready'
        ? copy.ready
        : !canCompletePlanet
          ? copy.learning
        : remaining === 0
          ? copy.allPointsDone
          : copy.inProgress;

  return (
    <div style={lumiGuideCardStyle}>
      <div style={lumiGuideAvatarStyle}>
        <LumiAvatar state={remaining === 0 ? 'celebrate' : 'curious'} size={58} />
      </div>
      <div style={lumiGuideTextStyle}>
        <span style={lumiGuideEyebrowStyle}>Lumi Guide</span>
        <p style={lumiGuideMessageStyle}>{message}</p>
      </div>
    </div>
  );
}

export function SharedDiaryNoticeCard({ copy }: { copy: PlanetDiaryCopy['summary']['sharedNotice'] }) {
  return (
    <div style={sharedResultCardStyle}>
      <div style={{ display: 'grid', gap: '6px' }}>
        <div style={sharedResultEyebrowStyle}>Shared Read-only Diary</div>
        <strong style={sharedResultTitleStyle}>{copy.title}</strong>
        <p style={sharedResultMessageStyle}>{copy.message}</p>
      </div>
    </div>
  );
}

const sharedResultCardStyle = {
  display: 'grid',
  gap: '10px',
  padding: '16px 18px',
  borderRadius: '20px',
  border: '1px solid rgba(131, 198, 255, 0.18)',
  background: 'rgba(109, 176, 238, 0.08)',
} as const;

const sharedResultEyebrowStyle = {
  fontSize: '12px',
  fontWeight: 800,
  letterSpacing: '0.08em',
  textTransform: 'uppercase',
  color: '#B9DFFF',
} as const;

const sharedResultTitleStyle = {
  fontSize: '16px',
  color: '#F4F8FF',
} as const;

const sharedResultMessageStyle = {
  margin: 0,
  fontSize: '13px',
  lineHeight: 1.7,
  color: '#D9E8FF',
} as const;

const journalContextBandStyle = {
  display: 'grid',
  gap: '12px',
  padding: '18px 18px 16px',
  borderRadius: '22px',
  border: '1px solid rgba(228, 199, 122, 0.16)',
  background: 'linear-gradient(180deg, rgba(45, 33, 14, 0.30), rgba(14, 23, 38, 0.46))',
  boxShadow: 'inset 0 1px 0 rgba(255, 238, 201, 0.08)',
} as const;

const journalContextHeaderStyle = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '12px',
  flexWrap: 'wrap',
} as const;

const journalContextEyebrowStyle = {
  fontSize: '11px',
  fontWeight: 800,
  letterSpacing: '0.16em',
  textTransform: 'uppercase' as const,
  color: '#E8D39A',
} as const;

const journalContextTitleStyle = {
  color: '#FFF4D2',
  fontSize: '15px',
  fontWeight: 800,
  lineHeight: 1.45,
} as const;

const journalContextBodyStyle = {
  margin: 0,
  color: '#F7F2E2',
  fontSize: '18px',
  lineHeight: 1.75,
} as const;

const startButtonStyle = {
  border: 'none',
  borderRadius: '999px',
  padding: '12px 18px',
  background: 'linear-gradient(135deg, #39d0ff, #4f8dff)',
  color: '#04111f',
  fontWeight: 800,
  fontSize: '14px',
  cursor: 'pointer',
} as const;

const startButtonDisabledStyle = {
  opacity: 0.7,
  cursor: 'progress',
} as const;

const startReadyBadgeStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  padding: '10px 14px',
  borderRadius: '999px',
  background: 'rgba(73, 214, 163, 0.16)',
  border: '1px solid rgba(73, 214, 163, 0.28)',
  color: '#93f7ce',
  fontWeight: 700,
  fontSize: '13px',
} as const;

const lumiGuideCardStyle = {
  display: 'grid',
  gridTemplateColumns: 'auto 1fr',
  alignItems: 'center',
  gap: '14px',
  padding: '16px 18px',
  borderRadius: '22px',
  border: '1px solid rgba(228, 199, 122, 0.18)',
  background: 'linear-gradient(135deg, rgba(47, 31, 13, 0.42), rgba(12, 29, 48, 0.58))',
  boxShadow: 'inset 0 1px 0 rgba(255, 238, 201, 0.08)',
} as const;

const lumiGuideAvatarStyle = {
  display: 'grid',
  placeItems: 'center',
  width: '70px',
  height: '70px',
  borderRadius: '20px',
  background: 'rgba(255, 244, 214, 0.10)',
  border: '1px solid rgba(255, 236, 194, 0.16)',
} as const;

const lumiGuideTextStyle = {
  display: 'grid',
  gap: '5px',
  minWidth: 0,
} as const;

const lumiGuideEyebrowStyle = {
  fontSize: '11px',
  fontWeight: 900,
  letterSpacing: '0.12em',
  textTransform: 'uppercase' as const,
  color: '#E8D39A',
} as const;

const lumiGuideMessageStyle = {
  margin: 0,
  color: '#F7F2E2',
  fontSize: '14px',
  lineHeight: 1.7,
} as const;

const journalProgressWrapStyle = {
  display: 'grid',
  gap: '8px',
} as const;

const journalProgressHeaderStyle = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: '12px',
  flexWrap: 'wrap',
  color: '#F6E9C5',
  fontSize: '12px',
  fontWeight: 700,
} as const;

const journalProgressTrackStyle = {
  position: 'relative',
  width: '100%',
  height: '12px',
  borderRadius: '999px',
  overflow: 'hidden',
  background: 'rgba(255,255,255,0.12)',
  border: '1px solid rgba(255,255,255,0.14)',
} as const;

const journalProgressFillStyle = {
  position: 'absolute',
  inset: 0,
  borderRadius: '999px',
  background: 'linear-gradient(90deg, rgba(72, 196, 130, 0.96), rgba(231, 182, 64, 0.96))',
  boxShadow: '0 0 18px rgba(134, 211, 136, 0.22)',
} as const;

const journalProgressActionStyle = {
  display: 'grid',
  gap: '10px',
  paddingTop: '8px',
  borderTop: '1px solid rgba(255, 236, 194, 0.12)',
} as const;

const progressActionWrapStyle = {
  display: 'flex',
  alignItems: 'center',
  gap: '10px',
  flexWrap: 'wrap',
} as const;

const completeButtonStyle = {
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '38px',
  padding: '0 16px',
  borderRadius: '999px',
  background: 'rgba(250, 180, 50, 0.14)',
  border: '1px solid rgba(250, 180, 50, 0.3)',
  color: '#FAE0A0',
  fontSize: '14px',
  fontWeight: 700,
  cursor: 'pointer',
} as const;

const completionBlockedButtonStyle = {
  ...completeButtonStyle,
  opacity: 0.45,
  cursor: 'not-allowed',
} as const;
