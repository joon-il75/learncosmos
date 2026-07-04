'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { buildDefaultRecommendationQuery } from '@/lib/recommendation/recommendationQueryBuilder'
import {
  compareRecommendations,
  createScenario,
  fetchLabelDataset,
  fetchLabelSummary,
  fetchLessonSearchPrerunReports,
  fetchRankerArtifact,
  fetchRecommendationSpecMetrics,
  fetchScenarioConsole,
  fetchWorkerObservations,
  fetchScenarioDetail,
  generateLessons,
  patchRolloutRanker,
  runGoalAction,
  saveExternalCandidate,
  saveLabel,
  type AuthHeaders,
} from './api'
import type {
  Candidate,
  DebugLabel,
  DebugScenario,
  GeneratedLessonItem,
  GeneratedLessonsSnapshot,
  GoalProfile,
  LabelSummary,
  LLMJobObservations,
  LessonSearchPrerunReportsResponse,
  RankerArtifactStatus,
  RecommendationCompareLesson,
  RecommendationComparison,
  RecommendationSpecMetrics,
  RolloutState,
  ScenarioInput,
} from './types'
import {
  buildStoredRecommendationComparison,
  candidateFeatureSnapshotForLabel,
  candidateLabelKey,
  parseScenarioJSON,
  recommendationQueriesByLesson,
} from './utils'

export function useRecommendationDebugConsole() {
  const router = useRouter()
  const [error, setError] = useState('')
  const [rolloutState, setRolloutState] = useState<RolloutState | null>(null)
  const [rankerInput, setRankerInput] = useState({ ranker_model_version:'', feature_schema_version:'ranker-feature-v1' })
  const [savingRanker, setSavingRanker] = useState(false)
  const [rankerArtifact, setRankerArtifact] = useState<RankerArtifactStatus | null>(null)
  const [workerObservations, setWorkerObservations] = useState<LLMJobObservations | null>(null)
  const [loadingWorkerObservations, setLoadingWorkerObservations] = useState(false)
  const [specMetrics, setSpecMetrics] = useState<RecommendationSpecMetrics | null>(null)
  const [specMetricsDays, setSpecMetricsDays] = useState(7)
  const [loadingSpecMetrics, setLoadingSpecMetrics] = useState(false)
  const [lessonSearchPrerunReports, setLessonSearchPrerunReports] = useState<LessonSearchPrerunReportsResponse | null>(null)
  const [loadingLessonSearchPrerunReports, setLoadingLessonSearchPrerunReports] = useState(false)
  const [scenarios, setScenarios] = useState<DebugScenario[]>([])
  const [selectedScenario, setSelectedScenario] = useState<DebugScenario | null>(null)
  const [scenarioInput, setScenarioInput] = useState<ScenarioInput>({
    course_title: '',
    notes: '',
  })
  const [goalInitialIntent, setGoalInitialIntent] = useState('')
  const [goalProfile, setGoalProfile] = useState<GoalProfile | null>(null)
  const [goalMessage, setGoalMessage] = useState('')
  const [confirmGoalText, setConfirmGoalText] = useState('')
  const [generatedLessons, setGeneratedLessons] = useState<GeneratedLessonsSnapshot | null>(null)
  const [recommendationComparison, setRecommendationComparison] = useState<RecommendationComparison | null>(null)
  const [lessonRecommendationQueries, setLessonRecommendationQueries] = useState<Record<string, string>>({})
  const [labelsByCandidateKey, setLabelsByCandidateKey] = useState<Record<string, DebugLabel>>({})
  const [labelSummary, setLabelSummary] = useState<LabelSummary | null>(null)
  const [savingLabelKey, setSavingLabelKey] = useState<string | null>(null)
  const [savingExternalCandidateKey, setSavingExternalCandidateKey] = useState<string | null>(null)
  const [exportingLabelFormat, setExportingLabelFormat] = useState<'jsonl' | 'csv' | null>(null)
  const [loadingScenarios, setLoadingScenarios] = useState(false)
  const [loadingScenarioDetailId, setLoadingScenarioDetailId] = useState<string | null>(null)
  const [creatingScenario, setCreatingScenario] = useState(false)
  const [runningGoalAction, setRunningGoalAction] = useState(false)
  const [generatingLessons, setGeneratingLessons] = useState(false)
  const [runningComparison, setRunningComparison] = useState(false)
  const [runningComparisonLessonId, setRunningComparisonLessonId] = useState<string | null>(null)

  const authHeaders = (): AuthHeaders => ({})

  useEffect(() => {
    if (!rolloutState) return
    setRankerInput({
      ranker_model_version: rolloutState.ranker_model_version || '',
      feature_schema_version: rolloutState.feature_schema_version || 'ranker-feature-v1',
    })
  }, [rolloutState?.id, rolloutState?.ranker_model_version, rolloutState?.feature_schema_version])

  useEffect(() => {
    if (!loadingScenarios) void loadScenarioConsole()
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    if (selectedScenario?.id) {
      setGoalInitialIntent(selectedScenario.initial_user_intent || '')
      void loadLabelSummary(selectedScenario.id)
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedScenario?.id])

  const loadScenarioConsole = async () => {
    const headers = authHeaders()
    setLoadingScenarios(true)
    setError('')
    try {
      const data = await fetchScenarioConsole(headers)
      setRolloutState(data.rolloutState)
      setRankerArtifact(data.rankerArtifact)
      setWorkerObservations(data.workerObservations)
      setSpecMetrics(data.specMetrics)
      setLessonSearchPrerunReports(data.lessonSearchPrerunReports)
      setSpecMetricsDays(data.specMetrics?.window_days ?? 7)
      setScenarios(data.scenarios)
      setSelectedScenario(prev => {
        if (prev && data.scenarios.some(item => item.id === prev.id)) return prev
        return null
      })
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : '추천 실험 콘솔 조회 중 오류가 발생했습니다.')
    } finally {
      setLoadingScenarios(false)
    }
  }

  const saveRolloutRanker = async () => {
    const headers = authHeaders()
    setSavingRanker(true)
    setError('')
    try {
      const state = await patchRolloutRanker(headers, {
        ranker_model_version: rankerInput.ranker_model_version.trim(),
        feature_schema_version: rankerInput.feature_schema_version.trim() || 'ranker-feature-v1',
      })
      setRolloutState(state)
      await loadRankerArtifact()
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Ranker 설정 저장 중 오류가 발생했습니다.')
    } finally {
      setSavingRanker(false)
    }
  }

  const loadRankerArtifact = async () => {
    const headers = authHeaders()
    try {
      setRankerArtifact(await fetchRankerArtifact(headers))
    } catch (e: unknown) {
      setRankerArtifact(null)
      setError(e instanceof Error ? e.message : 'Ranker artifact 조회 중 오류가 발생했습니다.')
    }
  }

  const loadWorkerObservations = async () => {
    const headers = authHeaders()
    setLoadingWorkerObservations(true)
    setError('')
    try {
      setWorkerObservations(await fetchWorkerObservations(headers))
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'worker 관측 조회 중 오류가 발생했습니다.')
    } finally {
      setLoadingWorkerObservations(false)
    }
  }


  const loadSpecMetrics = async (days = specMetricsDays) => {
    const headers = authHeaders()
    setLoadingSpecMetrics(true)
    setError('')
    try {
      const metrics = await fetchRecommendationSpecMetrics(headers, days)
      setSpecMetrics(metrics)
      setSpecMetricsDays(metrics?.window_days ?? days)
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : '명세 관측 조회 중 오류가 발생했습니다.')
    } finally {
      setLoadingSpecMetrics(false)
    }
  }

  const loadLessonSearchPrerunReports = async () => {
    const headers = authHeaders()
    setLoadingLessonSearchPrerunReports(true)
    setError('')
    try {
      setLessonSearchPrerunReports(await fetchLessonSearchPrerunReports(headers, 5))
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : '리슨 검색 프리런 report 조회 중 오류가 발생했습니다.')
    } finally {
      setLoadingLessonSearchPrerunReports(false)
    }
  }

  const loadLabelSummary = async (scenarioId: string) => {
    const headers = authHeaders()
    try {
      setLabelSummary(await fetchLabelSummary(headers, scenarioId))
    } catch (e: unknown) {
      setLabelSummary(null)
      setError(e instanceof Error ? e.message : '라벨 요약 조회 중 오류가 발생했습니다.')
    }
  }

  const resetSelectedScenarioProgress = () => {
    setSelectedScenario(null)
    setGoalInitialIntent('')
    setGoalProfile(null)
    setGoalMessage('')
    setConfirmGoalText('')
    setGeneratedLessons(null)
    setRecommendationComparison(null)
    setLessonRecommendationQueries({})
    setLabelsByCandidateKey({})
    setLabelSummary(null)
    setError('')
  }

  const updateScenarioInput = (patch: Partial<ScenarioInput>) => {
    if (selectedScenario && typeof patch.course_title === 'string' && patch.course_title !== scenarioInput.course_title) {
      resetSelectedScenarioProgress()
    }
    setScenarioInput(prev => ({ ...prev, ...patch }))
  }

  const buildLessonRecommendationQueries = (
    snapshot: GeneratedLessonsSnapshot | null,
    scenario: DebugScenario | null,
  ): Record<string, string> => {
    if (!snapshot) return {}
    return Object.fromEntries(snapshot.lessons.map(lesson => [
      lesson.lesson_id,
      buildDefaultRecommendationQuery({
        target: 'new',
        courseTitle: snapshot.draft_title || scenario?.course_title,
        courseGoal: snapshot.confirmed_goal || scenario?.initial_user_intent,
        regionTitle: lesson.title,
        regionDescription: lesson.objective,
        nodeTitleInput: lesson.title,
      }),
    ]))
  }

  const loadScenarioDetail = async (scenario: DebugScenario) => {
    const headers = authHeaders()
    setLoadingScenarioDetailId(scenario.id)
    setError('')
    try {
      const detail = await fetchScenarioDetail(headers, scenario.id)
      setSelectedScenario(detail)
      setScenarios(prev => prev.map(item => item.id === detail.id ? { ...item, status: detail.status, initial_user_intent: detail.initial_user_intent, updated_at: detail.updated_at } : item))
      setGoalInitialIntent(detail.initial_user_intent || '')
      setGoalMessage('')
      setConfirmGoalText('')
      setGoalProfile(parseScenarioJSON<GoalProfile>(detail.goal_profile_snapshot))
      const lessonsSnapshot = parseScenarioJSON<GeneratedLessonsSnapshot>(detail.generated_lessons_snapshot)
      setGeneratedLessons(lessonsSnapshot)
      const storedComparison = buildStoredRecommendationComparison(detail.runs ?? [])
      setRecommendationComparison(storedComparison)
      const storedQueries = recommendationQueriesByLesson(storedComparison)
      setLessonRecommendationQueries({
        ...buildLessonRecommendationQueries(lessonsSnapshot, detail),
        ...storedQueries,
      })
      const labels = detail.labels ?? []
      setLabelsByCandidateKey(Object.fromEntries(labels.map(label => [label.candidate_key, label])))
      void loadLabelSummary(detail.id)
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : '시나리오 상세 조회 중 오류가 발생했습니다.')
    } finally {
      setLoadingScenarioDetailId(null)
    }
  }

  const createScenarioFromInput = async (headers: AuthHeaders) => {
    if (!scenarioInput.course_title.trim()) {
      setError('코스 제목을 입력해주세요.')
      return null
    }
    setCreatingScenario(true)
    setError('')
    try {
      const created = await createScenario(headers, scenarioInput)
      if (created) {
        setScenarios(prev => [created, ...prev])
        setSelectedScenario(created)
        setGoalInitialIntent('')
        setScenarioInput({ course_title: '', notes: '' })
        return created
      }
      return null
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : '시나리오 생성 중 오류가 발생했습니다.')
      return null
    } finally {
      setCreatingScenario(false)
    }
  }

  const runScenarioGoalAction = async (action: 'start' | 'message' | 'confirm', goalOverride = '') => {
    const headers = authHeaders()
    if (action !== 'start' && !selectedScenario) return
    if (action === 'start' && !selectedScenario && !scenarioInput.course_title.trim()) {
      setError('코스 제목을 입력해주세요.')
      return
    }
    const startMessage = goalInitialIntent.trim() || scenarioInput.course_title.trim()
    const messageText = (goalOverride || goalMessage).trim()
    let body: Record<string, string>
    if (action === 'message') {
      body = { message: messageText }
    } else if (action === 'confirm') {
      body = { confirmed_goal: (goalOverride || confirmGoalText).trim() }
    } else {
      body = { message: startMessage }
    }
    if (action === 'message' && !messageText) {
      setError('루미에게 보낼 메시지를 입력해주세요.')
      return
    }
    if (action === 'confirm' && !String(body.confirmed_goal ?? '').trim()) {
      setError('확정할 목표 문장을 입력해주세요.')
      return
    }

    setRunningGoalAction(true)
    setError('')
    try {
      let scenario = selectedScenario
      if (action === 'start' && !scenario) {
        scenario = await createScenarioFromInput(headers)
        if (!scenario) return
      }
      if (!scenario) return

      const nextGoal = await runGoalAction(headers, scenario.id, action, body)
      setGoalProfile(nextGoal)
      if (action === 'start') {
        setConfirmGoalText('')
        const nextIntent = startMessage
        setScenarios(prev => prev.map(item => item.id === scenario.id ? { ...item, initial_user_intent: nextIntent } : item))
        setSelectedScenario(prev => prev ? { ...prev, initial_user_intent: nextIntent } : prev)
      }
      if (action === 'message') {
        setGoalMessage('')
      }
      if (action === 'confirm') {
        const confirmedScenario = { ...scenario, status: 'goal_confirmed' }
        setScenarios(prev => prev.map(item => item.id === scenario.id ? confirmedScenario : item))
        setSelectedScenario(prev => prev ? { ...prev, status: 'goal_confirmed' } : prev)
        await generateScenarioLessons(confirmedScenario)
      }
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : '목표 채팅 처리 중 오류가 발생했습니다.')
    } finally {
      setRunningGoalAction(false)
    }
  }

  const generateScenarioLessons = async (targetScenario?: DebugScenario) => {
    const headers = authHeaders()
    const scenario = targetScenario ?? selectedScenario
    if (!headers || !scenario) return
    setGeneratingLessons(true)
    setError('')
    try {
      const snapshot = await generateLessons(headers, scenario.id)
      setGeneratedLessons(snapshot)
      setRecommendationComparison(null)
      setLessonRecommendationQueries(buildLessonRecommendationQueries(snapshot, scenario))
      setLabelsByCandidateKey({})
      setScenarios(prev => prev.map(item => item.id === scenario.id ? { ...item, status: 'lessons_generated' } : item))
      setSelectedScenario(prev => prev ? { ...prev, status: 'lessons_generated' } : prev)
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : '리슨 생성 중 오류가 발생했습니다.')
    } finally {
      setGeneratingLessons(false)
    }
  }

  const updateLessonRecommendationQuery = (lessonId: string, value: string) => {
    setLessonRecommendationQueries(prev => ({ ...prev, [lessonId]: value }))
  }

  const compareScenarioRecommendations = async (lesson?: GeneratedLessonItem) => {
    const headers = authHeaders()
    if (!headers || !selectedScenario) return
    setRunningComparison(true)
    setRunningComparisonLessonId(lesson?.lesson_id ?? null)
    setError('')
    try {
      const recommendationQuery = lesson ? lessonRecommendationQueries[lesson.lesson_id] ?? '' : ''
      const comparison = await compareRecommendations(headers, selectedScenario.id, lesson, recommendationQuery)
      const lessonsWithRun = (comparison?.lessons ?? []).map(item => ({ ...item, run_id: comparison?.run_id }))
      const nextComparison = comparison ? { ...comparison, lessons: lessonsWithRun } : null
      setLessonRecommendationQueries(prev => ({ ...prev, ...recommendationQueriesByLesson(nextComparison) }))
      setRecommendationComparison(prev => {
        if (!prev || !nextComparison || prev.scenario_id !== nextComparison.scenario_id || !lesson) {
          return nextComparison
        }
        const nextLessons = [...prev.lessons]
        for (const nextLesson of lessonsWithRun) {
          const index = nextLessons.findIndex(item => item.lesson.lesson_id === nextLesson.lesson.lesson_id)
          if (index >= 0) {
            nextLessons[index] = nextLesson
          } else {
            nextLessons.push(nextLesson)
          }
        }
        const totalCandidates = nextLessons.reduce((sum, item) => sum + item.baseline.candidates.length, 0)
        const totalLexicalHits = nextLessons.reduce((sum, item) => sum + item.baseline.stages.search.lexical_hit_count, 0)
        const totalVectorHits = nextLessons.reduce((sum, item) => sum + item.baseline.stages.search.vector_hit_count, 0)
        return {
          ...nextComparison,
          lessons: nextLessons,
          lesson_count: nextLessons.length,
          total_candidates: totalCandidates,
          total_lexical_hits: totalLexicalHits,
          total_vector_hits: totalVectorHits,
        }
      })
      setLabelsByCandidateKey({})
      setScenarios(prev => prev.map(item => item.id === selectedScenario.id ? { ...item, status: 'recommendation_tested' } : item))
      setSelectedScenario(prev => prev ? { ...prev, status: 'recommendation_tested' } : prev)
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : '추천 비교 중 오류가 발생했습니다.')
    } finally {
      setRunningComparison(false)
      setRunningComparisonLessonId(null)
    }
  }

  const saveCandidateLabel = async (
    lesson: RecommendationCompareLesson['lesson'],
    candidate: Candidate,
    candidateIndex: number,
    label: string,
    runIDOverride?: string,
  ) => {
    const headers = authHeaders()
    const runID = runIDOverride || recommendationComparison?.run_id
    if (!headers || !selectedScenario || !runID) return
    const candidateKey = candidateLabelKey(lesson.lesson_id, candidate, candidateIndex)
    const hadLabel = Boolean(labelsByCandidateKey[candidateKey])
    setSavingLabelKey(candidateKey)
    setError('')
    try {
      const saved = await saveLabel(headers, selectedScenario.id, {
        run_id: runID,
        candidate_key: candidateKey,
        content_id: candidate.content_id || '',
        url: candidate.external_url || '',
        label,
        note: '',
        baseline_rank: candidateIndex + 1,
        feature_snapshot: candidateFeatureSnapshotForLabel(lesson, candidate, candidateIndex),
      })
      if (saved) {
        setLabelsByCandidateKey(prev => ({ ...prev, [saved.candidate_key]: saved }))
        setScenarios(prev => prev.map(item => item.id === selectedScenario.id ? { ...item, status: 'labelled', label_count: (item.label_count ?? 0) + (hadLabel ? 0 : 1) } : item))
        setSelectedScenario(prev => prev ? { ...prev, status: 'labelled' } : prev)
        void loadLabelSummary(selectedScenario.id)
      }
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : '라벨 저장 중 오류가 발생했습니다.')
    } finally {
      setSavingLabelKey(null)
    }
  }

  const saveExternalCandidateContent = async (
    lesson: RecommendationCompareLesson['lesson'],
    candidate: Candidate,
    candidateIndex: number,
  ) => {
    const headers = authHeaders()
    if (!headers || !selectedScenario) return
    const candidateKey = candidateLabelKey(lesson.lesson_id, candidate, candidateIndex)
    const url = String(candidate.external_url || '').trim()
    if (!url) {
      setError('외부 후보 URL이 없어 저장할 수 없습니다.')
      return
    }
    setSavingExternalCandidateKey(candidateKey)
    setError('')
    try {
      const saved = await saveExternalCandidate(headers, selectedScenario.id, {
        lesson_id: lesson.lesson_id,
        title: candidate.title,
        description: candidate.description || '',
        url,
        source: candidate.external_source || candidate.content_type || '',
        external_content_id: candidate.external_content_id || '',
        thumbnail_url: candidate.thumbnail_url || '',
        author: candidate.author || '',
        language: candidate.language || 'ko',
      })
      const savedContentID = saved.content?.id || candidate.content_id || ''
      setRecommendationComparison(prev => {
        if (!prev) return prev
        return {
          ...prev,
          lessons: prev.lessons.map(item => {
            if (item.lesson.lesson_id !== lesson.lesson_id) return item
            return {
              ...item,
              baseline: {
                ...item.baseline,
                candidates: item.baseline.candidates.map((current, idx) => {
                  if (idx !== candidateIndex) return current
                  return {
                    ...current,
                    content_id: savedContentID,
                    already_saved: true,
                    eligible_for_save: false,
                    selection_reason: current.selection_reason || '운영자가 외부 후보를 콘텐츠 저장소에 저장함',
                  }
                }),
              },
            }
          }),
        }
      })
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : '외부 후보 저장 중 오류가 발생했습니다.')
    } finally {
      setSavingExternalCandidateKey(null)
    }
  }

  const exportLabelDataset = async (format: 'jsonl' | 'csv') => {
    const headers = authHeaders()
    if (!headers || !selectedScenario) return
    setExportingLabelFormat(format)
    setError('')
    try {
      const blob = await fetchLabelDataset(headers, selectedScenario.id, format)
      const objectUrl = URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = objectUrl
      link.download = `recommendation-labels-${selectedScenario.id}.${format}`
      document.body.appendChild(link)
      link.click()
      link.remove()
      URL.revokeObjectURL(objectUrl)
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : '라벨 데이터셋 export 중 오류가 발생했습니다.')
    } finally {
      setExportingLabelFormat(null)
    }
  }

  return {
    error,
    rolloutState,
    rankerInput,
    setRankerInput,
    savingRanker,
    rankerArtifact,
    workerObservations,
    loadingWorkerObservations,
    specMetrics,
    specMetricsDays,
    loadingSpecMetrics,
    lessonSearchPrerunReports,
    loadingLessonSearchPrerunReports,
    scenarios,
    selectedScenario,
    scenarioInput,
    updateScenarioInput,
    goalInitialIntent,
    goalProfile,
    goalMessage,
    setGoalMessage,
    confirmGoalText,
    generatedLessons,
    recommendationComparison,
    lessonRecommendationQueries,
    labelsByCandidateKey,
    labelSummary,
    savingLabelKey,
    savingExternalCandidateKey,
    exportingLabelFormat,
    loadingScenarios,
    loadingScenarioDetailId,
    creatingScenario,
    runningGoalAction,
    generatingLessons,
    runningComparison,
    runningComparisonLessonId,
    loadScenarioConsole,
    saveRolloutRanker,
    loadRankerArtifact,
    loadWorkerObservations,
    loadSpecMetrics,
    loadLessonSearchPrerunReports,
    loadLabelSummary,
    loadScenarioDetail,
    runScenarioGoalAction,
    generateScenarioLessons,
    updateLessonRecommendationQuery,
    compareScenarioRecommendations,
    saveCandidateLabel,
    saveExternalCandidateContent,
    exportLabelDataset,
  }
}
