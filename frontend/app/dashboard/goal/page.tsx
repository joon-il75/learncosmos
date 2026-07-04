'use client'

import { Suspense, useCallback, useEffect, useRef, useState, type CSSProperties } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { GoalFlowLoadingOverlay, type GoalFlowStage, type PlanetTextureMapItem } from '@/components/goal-interview/GoalFlowLoadingOverlay'
import { LumiInterviewPanel } from '@/components/goal-interview/LumiInterviewPanel'
import GoalPointStatusBar from '@/components/goal-interview/GoalPointStatusBar'
import { type GoalProfile, useGoalInterview } from '@/components/goal-interview/useGoalInterview'
import { executeDashboardCreateDraft, pollDashboardCreateDraftJob, type DashboardCourseGenerationWaitEstimate, type DashboardCreateDraftFailure } from '@/lib/world-ui-engine/ctaEngine'
import { normalizeLocale, type Locale } from '@/lib/i18n/locales'
import { getDashboardGoalCopy } from '@/lib/i18n/pages/dashboardGoal'
import AppHeaderShell from '@/components/common/AppHeaderShell'
import LearnerHeaderActions, { getLearnerHeaderActionsCopy } from '@/components/common/LearnerHeaderActions'
import LumiModalShell, {
  lumiModalPrimaryButtonStyle,
  lumiModalSecondaryButtonStyle,
} from '@/components/common/LumiModalShell'

function resolveCreateDraftErrorMessage(
  reason: DashboardCreateDraftFailure['reason'],
  fallback: string,
  rawMessage: string,
): string {
  if (reason === 'busy' || reason === 'retry') return fallback
  return rawMessage || fallback
}

const COURSE_GENERATION_JOB_STORAGE_KEY = 'learnweaver:course-generation-job:v1'
const COURSE_GENERATION_JOB_MAX_AGE_MS = 30 * 60 * 1000

type StoredCourseGenerationJob = {
  jobId: string
  pollURL?: string
  query: string
  learningGoal?: string
  intent: string
  startedAt: number
}

function readStoredCourseGenerationJob(intent: string): StoredCourseGenerationJob | null {
  if (typeof window === 'undefined') return null
  try {
    const raw = window.sessionStorage.getItem(COURSE_GENERATION_JOB_STORAGE_KEY)
    if (!raw) return null
    const parsed = JSON.parse(raw) as Partial<StoredCourseGenerationJob>
    if (!parsed.jobId || !parsed.query || parsed.intent !== intent || !parsed.startedAt) return null
    if (Date.now() - parsed.startedAt > COURSE_GENERATION_JOB_MAX_AGE_MS) {
      window.sessionStorage.removeItem(COURSE_GENERATION_JOB_STORAGE_KEY)
      return null
    }
    return parsed as StoredCourseGenerationJob
  } catch {
    return null
  }
}

function writeStoredCourseGenerationJob(job: StoredCourseGenerationJob) {
  if (typeof window === 'undefined') return
  try {
    window.sessionStorage.setItem(COURSE_GENERATION_JOB_STORAGE_KEY, JSON.stringify(job))
  } catch {
    // Storage failure should not block generation.
  }
}

function clearStoredCourseGenerationJob() {
  if (typeof window === 'undefined') return
  try {
    window.sessionStorage.removeItem(COURSE_GENERATION_JOB_STORAGE_KEY)
  } catch {
    // Ignore storage cleanup failure.
  }
}

export default function GoalInterviewPage() {
  return (
    <Suspense fallback={
      <div style={loadingStyle}>
        <GoalFlowLoadingOverlay stage="bootstrapping_lumi" />
      </div>
    }
    >
      <GoalInterviewPageContent />
    </Suspense>
  )
}

function GoalInterviewPageContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const intent = searchParams.get('intent') ?? ''
  const decodedIntent = intent ? decodeURIComponent(intent) : ''
  const [uiLocale, setUiLocale] = useState<Locale>('ko')
  const copy = getDashboardGoalCopy(uiLocale)
  const learnerHeaderCopy = getLearnerHeaderActionsCopy(uiLocale)

  const {
    profile,
    isLoading,
    isSending,
    error: interviewError,
    startInterview,
    sendMessage,
    confirmGoal,
    reviseGoal,
    applyRebuildDecision,
    resetActiveGoal,
  } = useGoalInterview(undefined, { autoLoad: false, errorCopy: copy.hookErrors })

  const startCalledRef = useRef(false)
  const [startAttempted, setStartAttempted] = useState(false)
  const [isCreatingDraft, setIsCreatingDraft] = useState(false)
  const [goalFlowStage, setGoalFlowStage] = useState<GoalFlowStage | null>(decodedIntent ? 'bootstrapping_lumi' : null)
  const [creationError, setCreationError] = useState<string | null>(null)
  const [creationErrorReason, setCreationErrorReason] = useState<DashboardCreateDraftFailure['reason']>(undefined)
  const [pendingPlanetTitle, setPendingPlanetTitle] = useState<string | null>(null)
  const [availablePlanetTextureMaps, setAvailablePlanetTextureMaps] = useState<PlanetTextureMapItem[]>([])
  const [selectedPlanetTextureMapId, setSelectedPlanetTextureMapId] = useState<string | null>(null)
  const [isDraftReadyToContinue, setIsDraftReadyToContinue] = useState(false)
  const [pendingDestination, setPendingDestination] = useState<string | null>(null)
  const [pendingDraftId, setPendingDraftId] = useState<string | null>(null)
  const [courseGenerationWait, setCourseGenerationWait] = useState<DashboardCourseGenerationWaitEstimate | null>(null)
  const [showExitConfirm, setShowExitConfirm] = useState(false)
  const selectedTextureMapRef = useRef<string | null>(null)

  const loadPlanetTextureMaps = useCallback(() => {
    void fetch('/api/v1/public/planet-texture-maps', { credentials: 'include' })
      .then((r) => r.json())
      .then((data: { planet_texture_maps?: PlanetTextureMapItem[] }) => {
        setAvailablePlanetTextureMaps(data.planet_texture_maps ?? [])
      })
      .catch(() => { /* 실패해도 진행 */ })
  }, [])

  useEffect(() => {
    let cancelled = false
    const loadLocale = async () => {
      try {
        await fetch('/api/v1/auth/refresh', { method: 'POST', credentials: 'include' })
        const res = await fetch('/api/v1/auth/me', { credentials: 'include', cache: 'no-store' })
        if (!res.ok) return
        const data = (await res.json()) as { ui_locale?: string }
        if (!cancelled) setUiLocale(normalizeLocale(data.ui_locale))
      } catch {
        // Keep Korean fallback.
      }
    }
    void loadLocale()
    return () => {
      cancelled = true
    }
  }, [])

  useEffect(() => {
    if (!decodedIntent || startCalledRef.current) return
    const storedJob = readStoredCourseGenerationJob(decodedIntent)
    if (!storedJob) return

    startCalledRef.current = true
    setStartAttempted(true)
    setIsCreatingDraft(true)
    setCreationError(null)
    setCreationErrorReason(undefined)
    setPendingPlanetTitle(storedJob.learningGoal || storedJob.query || decodedIntent)
    setSelectedPlanetTextureMapId(null)
    selectedTextureMapRef.current = null
    setIsDraftReadyToContinue(false)
    setPendingDestination(null)
    setPendingDraftId(null)
    setCourseGenerationWait(null)
    setGoalFlowStage('creating_planet')
    loadPlanetTextureMaps()

    let cancelled = false
    void (async () => {
      try {
        const result = await pollDashboardCreateDraftJob({
          query: storedJob.query,
          learningGoal: storedJob.learningGoal,
          jobId: storedJob.jobId,
          pollURL: storedJob.pollURL,
          onJobProgress: (job) => setCourseGenerationWait(job.estimatedWait ?? null),
        }, storedJob.startedAt)
        if (cancelled) return
        if (!result.ok) {
          clearStoredCourseGenerationJob()
          setCreationErrorReason(result.reason)
          throw new Error(resolveCreateDraftErrorMessage(result.reason, copy.page.errors.createFailed, result.message))
        }

        setPendingPlanetTitle(result.draftTitle ?? storedJob.learningGoal ?? storedJob.query ?? decodedIntent)
        const destination = result.courseId
          ? `/dashboard/planets/learning/${result.courseId}`
          : `/dashboard/course-drafts/${result.draftId}`
        setPendingDestination(destination)
        setPendingDraftId(result.draftId)
        setCourseGenerationWait(null)
        setIsDraftReadyToContinue(true)
      } catch (err) {
        if (cancelled) return
        setCreationError(err instanceof Error ? err.message : copy.page.errors.createFailed)
        setPendingPlanetTitle(null)
        setPendingDestination(null)
        setPendingDraftId(null)
        setCourseGenerationWait(null)
        setIsDraftReadyToContinue(false)
        setGoalFlowStage(null)
      } finally {
        if (!cancelled) setIsCreatingDraft(false)
      }
    })()

    return () => {
      cancelled = true
    }
  }, [copy.page.errors.createFailed, decodedIntent, loadPlanetTextureMaps])

  useEffect(() => {
    if (startCalledRef.current) return
    if (!decodedIntent) return
    startCalledRef.current = true
    setStartAttempted(true)
    setCreationError(null)
    setGoalFlowStage('starting_interview')
    void (async () => {
      try {
        await resetActiveGoal()
        await startInterview(decodedIntent)
        setGoalFlowStage(null)
      } catch (err) {
        setCreationError(err instanceof Error ? err.message : copy.page.errors.resetFailed)
        setGoalFlowStage(null)
      }
    })()
  }, [copy.page.errors.resetFailed, decodedIntent, resetActiveGoal, startInterview])

  useEffect(() => {
    if (!decodedIntent) {
      router.replace('/dashboard')
    }
  }, [decodedIntent, router])

  const createDraftFromGoal = useCallback(async (goalProfile: GoalProfile | null) => {
    if (!goalProfile?.confirmed_goal) {
      setCreationError(copy.page.errors.missingConfirmedGoal)
      return
    }

    setIsCreatingDraft(true)
    setCreationError(null)
    setCreationErrorReason(undefined)
    setPendingPlanetTitle(null)
    setSelectedPlanetTextureMapId(null)
    selectedTextureMapRef.current = null
    setIsDraftReadyToContinue(false)
    setPendingDestination(null)
    setPendingDraftId(null)
    setCourseGenerationWait(null)
    setGoalFlowStage('creating_planet')

    clearStoredCourseGenerationJob()
    // 신규 자전용 텍스처맵만 로드한다. 기존 planet type 목록은 로딩 선택 UI에서 쓰지 않는다.
    loadPlanetTextureMaps()

    try {
      const query = goalProfile.user_intent || decodedIntent
      const learningGoal = goalProfile.confirmed_goal
      const startedAt = Date.now()
      const result = await executeDashboardCreateDraft({
        query,
        learningGoal,
        startedAt,
        onJobAccepted: (job) => {
          setCourseGenerationWait(job.estimatedWait ?? null)
          writeStoredCourseGenerationJob({
            jobId: job.jobId,
            pollURL: job.pollURL,
            query,
            learningGoal,
            intent: decodedIntent,
            startedAt,
          })
        },
        onJobProgress: (job) => setCourseGenerationWait(job.estimatedWait ?? null),
      })
      if (!result.ok) {
        setCreationErrorReason(result.reason)
        throw new Error(resolveCreateDraftErrorMessage(result.reason, copy.page.errors.createFailed, result.message))
      }

      let finalTitle = result.draftTitle ?? goalProfile.user_intent ?? decodedIntent ?? null

      setPendingPlanetTitle(finalTitle)
      const destination = result.courseId
        ? `/dashboard/planets/learning/${result.courseId}`
        : `/dashboard/course-drafts/${result.draftId}`
      setPendingDestination(destination)
      setPendingDraftId(result.draftId)
      setCourseGenerationWait(null)
      setIsDraftReadyToContinue(true)
    } catch (err) {
      clearStoredCourseGenerationJob()
      setCreationError(err instanceof Error ? err.message : copy.page.errors.createFailed)
      setPendingPlanetTitle(null)
      setPendingDestination(null)
      setPendingDraftId(null)
      setCourseGenerationWait(null)
      setIsDraftReadyToContinue(false)
      setGoalFlowStage(null)
    } finally {
      setIsCreatingDraft(false)
    }
  }, [copy.page.errors.createFailed, copy.page.errors.missingConfirmedGoal, decodedIntent, loadPlanetTextureMaps])

  const handleSelectPlanetTextureMap = useCallback((textureMap: PlanetTextureMapItem) => {
    setSelectedPlanetTextureMapId(textureMap.id)
    selectedTextureMapRef.current = textureMap.id
    setPendingPlanetTitle(textureMap.name)
  }, [])

  const handleConfirmPlanetTextureMap = useCallback(() => {
    if (!selectedTextureMapRef.current || !pendingDestination || !isDraftReadyToContinue) return

    void (async () => {
      const draftId = pendingDraftId ?? pendingDestination.match(/\/dashboard\/course-drafts\/([^/?#]+)/)?.[1]
      if (draftId) {
        await fetch(`/api/v1/course-drafts/${draftId}`, {
          method: 'PATCH',
          credentials: 'include',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ planet_texture_map_id: selectedTextureMapRef.current }),
        }).catch(() => undefined)
      }
      clearStoredCourseGenerationJob()
      setGoalFlowStage('redirecting')
      router.push(pendingDestination)
    })()
  }, [isDraftReadyToContinue, pendingDestination, pendingDraftId, router])

  const handleConfirm = useCallback(async (goal: string) => {
    const confirmed = await confirmGoal(goal)
    if (!confirmed) return
    await createDraftFromGoal(confirmed)
  }, [confirmGoal, createDraftFromGoal])

  const handleRetryGoal = useCallback(() => {
    if (!decodedIntent || isSending || isCreatingDraft) return
    clearStoredCourseGenerationJob()
    setCreationError(null)
    setCreationErrorReason(undefined)
    setPendingPlanetTitle(null)
    setSelectedPlanetTextureMapId(null)
    selectedTextureMapRef.current = null
    setIsDraftReadyToContinue(false)
    setPendingDestination(null)
    setPendingDraftId(null)
    setCourseGenerationWait(null)
    setGoalFlowStage('starting_interview')
    void (async () => {
      try {
        await resetActiveGoal()
        await startInterview(decodedIntent)
      } catch (err) {
        setCreationError(err instanceof Error ? err.message : copy.page.errors.resetFailed)
      } finally {
        setGoalFlowStage(null)
      }
    })()
  }, [copy.page.errors.resetFailed, decodedIntent, isCreatingDraft, isSending, resetActiveGoal, startInterview])

  const handleExitToDashboard = useCallback(() => {
    setShowExitConfirm(false)
    router.replace('/dashboard')
  }, [router])

  const combinedError = creationError ?? interviewError
  const isWaitingForLumi = isLoading || (!!decodedIntent && (profile === null && (!startAttempted || isSending)) && !combinedError)

  if (isWaitingForLumi) {
    return (
      <div style={loadingStyle}>
        <GoalFlowLoadingOverlay
          stage={goalFlowStage ?? 'bootstrapping_lumi'}
          copy={copy.loading}
          courseGenerationWait={courseGenerationWait}
        />
      </div>
    )
  }

  if (profile === null && !!decodedIntent && combinedError) {
    return (
      <div style={loadingStyle}>
        <div style={errorBoxStyle}>
          <div style={{ fontSize: 15, fontWeight: 700, color: '#92400E', marginBottom: 8 }}>{copy.page.loadErrorTitle}</div>
          <div style={{ fontSize: 13, color: '#78350F', marginBottom: 16 }}>{combinedError}</div>
          <button onClick={() => router.replace('/dashboard')} style={backDashButtonStyle}>
            {copy.page.backToDashboard}
          </button>
        </div>
      </div>
    )
  }

  return (
    <div style={pageStyle}>
      {goalFlowStage ? (
        <GoalFlowLoadingOverlay
          stage={goalFlowStage}
          copy={copy.loading}
          planetTitle={pendingPlanetTitle}
          planetTextureMaps={availablePlanetTextureMaps}
          selectedPlanetTextureMapId={selectedPlanetTextureMapId}
          isReadyToContinue={isDraftReadyToContinue}
          onSelectPlanetTextureMap={handleSelectPlanetTextureMap}
          onConfirmPlanetTextureMap={handleConfirmPlanetTextureMap}
          courseGenerationWait={courseGenerationWait}
        />
      ) : null}
      <AppHeaderShell
        logoHref="/"
        logoIconSize={26}
        logoTextSize="17px"
        maxWidth="1080px"
        headerStyle={goalHeaderStyle}
        innerStyle={goalHeaderInnerStyle}
        leftMeta={
          <span style={goalHeaderMetaStyle}>
            <span style={goalHeaderEyebrowStyle}>{copy.page.headerEyebrow}</span>
            <span style={goalHeaderTitleStyle}>{copy.page.headerTitle}</span>
          </span>
        }
        rightSlot={<LearnerHeaderActions palette={goalHeaderActionsPalette} copy={learnerHeaderCopy} />}
      />

      <div style={contentStyle}>
        {creationErrorReason === 'busy' ? (
          <div style={busyNoticeStyle}>
            <div style={busyNoticeTitleStyle}>{copy.page.notices.busyTitle}</div>
            <div style={busyNoticeBodyStyle}>
              {copy.page.notices.busyBody}
            </div>
          </div>
        ) : creationErrorReason === 'retry' ? (
          <div style={busyNoticeStyle}>
            <div style={busyNoticeTitleStyle}>{copy.page.notices.retryTitle}</div>
            <div style={busyNoticeBodyStyle}>
              {copy.page.notices.retryBody}
            </div>
          </div>
        ) : null}
        <LumiInterviewPanel
          profile={profile}
          isSending={isSending}
          error={combinedError}
          onSendMessage={sendMessage}
          onConfirmGoal={handleConfirm}
          onRetryGoal={handleRetryGoal}
          onRebuildDecision={applyRebuildDecision}
          onConfirmedAction={profile?.interview_state === 'confirmed' ? () => void createDraftFromGoal(profile) : undefined}
          confirmedActionLabel={copy.page.confirmedAction}
          isConfirmedActionBusy={isCreatingDraft}
          pointStatusSlot={<GoalPointStatusBar copy={copy.points} />}
          copy={copy.interview}
          proposalCopy={copy.proposal}
          bottomActionSlot={
            <button
              type="button"
              onClick={() => setShowExitConfirm(true)}
              style={goalExitButtonStyle}
            >
              {copy.page.backToDashboard}
            </button>
          }
        />
      </div>
      {showExitConfirm ? (
        <LumiModalShell
          ariaLabel={copy.page.exitAria}
          eyebrow={copy.page.exitEyebrow}
          lumiState="curious"
          title={copy.page.exitTitle}
          message={copy.page.exitMessage}
          onClose={() => setShowExitConfirm(false)}
          actions={
            <>
              <button
                type="button"
                onClick={() => setShowExitConfirm(false)}
                style={lumiModalSecondaryButtonStyle}
              >
                {copy.page.continue}
              </button>
              <button
                type="button"
                onClick={handleExitToDashboard}
                style={lumiModalPrimaryButtonStyle}
              >
                {copy.page.dashboardShort}
              </button>
            </>
          }
        />
      ) : null}
    </div>
  )
}

const pageStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  height: '100dvh',
  background: '#FFF8E7',
}

const goalHeaderStyle: CSSProperties = {
  background: 'rgba(255, 253, 247, 0.94)',
  borderBottom: '1px solid rgba(245,158,11,0.18)',
  backdropFilter: 'blur(12px)',
  WebkitBackdropFilter: 'blur(12px)',
}

const goalHeaderInnerStyle: CSSProperties = {
  padding: '12px 18px',
  gap: '14px',
}

const goalHeaderMetaStyle: CSSProperties = {
  display: 'grid',
  gap: '2px',
}

const goalHeaderEyebrowStyle: CSSProperties = {
  fontSize: 11,
  fontWeight: 700,
  letterSpacing: '0.08em',
  textTransform: 'uppercase',
  color: 'rgba(146, 64, 14, 0.62)',
}

const goalHeaderTitleStyle: CSSProperties = {
  fontSize: 15,
  fontWeight: 800,
  color: '#3D2000',
  letterSpacing: '0.02em',
}

const goalHeaderActionsPalette = {
  text: '#3D2000',
  galaxyText: '#FFFFFF',
  galaxyBackground: '#4F36A6',
  galaxyBorder: 'rgba(79, 54, 166, 0.42)',
  galaxyShadow: '0 9px 22px rgba(79, 54, 166, 0.22)',
  userMenu: {
    triggerText: '#2B1A00',
    pointText: '#78350F',
    pointBackground: 'rgba(255, 247, 214, 0.98)',
    pointBorder: 'rgba(245,158,11,0.62)',
  },
}

const contentStyle: CSSProperties = {
  flex: 1,
  overflow: 'hidden',
  padding: '84px 20px 16px',
  maxWidth: 640,
  width: '100%',
  margin: '0 auto',
  display: 'flex',
  flexDirection: 'column',
  gap: 12,
}

const loadingStyle: CSSProperties = {
  position: 'relative',
  height: '100dvh',
  background: '#FFF8E7',
}

const errorBoxStyle: CSSProperties = {
  background: '#FFFDF7',
  border: '1px solid rgba(245,158,11,0.3)',
  borderRadius: 16,
  padding: '24px 28px',
  textAlign: 'center',
  maxWidth: 320,
}

const backDashButtonStyle: CSSProperties = {
  background: '#F59E0B',
  color: '#fff',
  border: 'none',
  borderRadius: 8,
  padding: '10px 20px',
  fontSize: 14,
  fontWeight: 700,
  cursor: 'pointer',
}

const busyNoticeStyle: CSSProperties = {
  background: 'linear-gradient(180deg, rgba(255,248,231,0.98) 0%, rgba(255,243,214,0.98) 100%)',
  border: '1px solid rgba(245,158,11,0.28)',
  borderRadius: 16,
  padding: '14px 16px',
  boxShadow: '0 8px 20px rgba(120,80,20,0.08)',
}

const busyNoticeTitleStyle: CSSProperties = {
  fontSize: 14,
  fontWeight: 800,
  color: '#78350F',
  marginBottom: 6,
}

const busyNoticeBodyStyle: CSSProperties = {
  fontSize: 13,
  lineHeight: 1.6,
  color: '#92400E',
}

const goalExitButtonStyle: CSSProperties = {
  minHeight: 36,
  padding: '0 14px',
  borderRadius: 999,
  border: '1px solid rgba(146, 64, 14, 0.34)',
  background: '#FFF8E7',
  color: '#78350F',
  fontSize: 13,
  fontWeight: 900,
  cursor: 'pointer',
  boxShadow: '0 6px 16px rgba(120, 53, 15, 0.08)',
}
