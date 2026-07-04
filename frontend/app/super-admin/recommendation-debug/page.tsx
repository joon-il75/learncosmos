'use client'

import SuperAdminPanelHeader from '@/components/super-admin/SuperAdminPanelHeader'
import SuperAdminPanelNav from '@/components/super-admin/SuperAdminPanelNav'
import GoalChatCard from './components/GoalChatCard'
import LabelSummaryCard from './components/LabelSummaryCard'
import LessonSearchPrerunCard from './components/LessonSearchPrerunCard'
import LessonGenerationCard from './components/LessonGenerationCard'
import RankerConditionCard from './components/RankerConditionCard'
import ScenarioInputCard from './components/ScenarioInputCard'
import ScenarioListCard from './components/ScenarioListCard'
import SpecMetricsCard from './components/SpecMetricsCard'
import WorkerObservationCard from './components/WorkerObservationCard'
import { S } from './styles'
import { useRecommendationDebugConsole } from './useRecommendationDebugConsole'

export default function RecommendationDebugPage() {
  const consoleState = useRecommendationDebugConsole()

  return (
    <main style={S.page}>
      <div style={S.shell}>
        <SuperAdminPanelHeader
          subtitle="추천 점검"
          description="코스 목표와 추천 결과를 운영자 기준으로 점검하고, 랭커 전환 상태를 단계별로 추적합니다."
        />
        <SuperAdminPanelNav activeSection="recommendation" />

        <div style={S.info}>
          추천 점검은 LightGBM Ranker 사용 조건을 확인한 뒤, 코스 생성 입력부터 목표 확정 채팅, 리슨 생성, 추천 콘텐츠까지 학습자 생성 흐름과 같은 순서로 테스트합니다.
        </div>

        {consoleState.error && <div style={{ ...S.warn, color:'#FFB4A2' }}>{consoleState.error}</div>}

        <RankerConditionCard
          rolloutState={consoleState.rolloutState}
          rankerInput={consoleState.rankerInput}
          setRankerInput={consoleState.setRankerInput}
          savingRanker={consoleState.savingRanker}
          rankerArtifact={consoleState.rankerArtifact}
          loadingScenarios={consoleState.loadingScenarios}
          onRefresh={consoleState.loadScenarioConsole}
          onSaveRanker={consoleState.saveRolloutRanker}
          onLoadRankerArtifact={consoleState.loadRankerArtifact}
        />

        <WorkerObservationCard
          observations={consoleState.workerObservations}
          loading={consoleState.loadingWorkerObservations}
          onRefresh={consoleState.loadWorkerObservations}
        />

        <SpecMetricsCard
          metrics={consoleState.specMetrics}
          loading={consoleState.loadingSpecMetrics}
          activeDays={consoleState.specMetricsDays}
          onRefresh={consoleState.loadSpecMetrics}
        />

        <LessonSearchPrerunCard
          reports={consoleState.lessonSearchPrerunReports}
          loading={consoleState.loadingLessonSearchPrerunReports}
          onRefresh={consoleState.loadLessonSearchPrerunReports}
        />

        <ScenarioInputCard
          scenarioInput={consoleState.scenarioInput}
          onChangeScenarioInput={consoleState.updateScenarioInput}
        />

        <ScenarioListCard
          scenarios={consoleState.scenarios}
          selectedScenario={consoleState.selectedScenario}
          loadingScenarioDetailId={consoleState.loadingScenarioDetailId}
          onSelectScenario={consoleState.loadScenarioDetail}
        />

        <GoalChatCard
          selectedScenario={consoleState.selectedScenario}
          scenarioInput={consoleState.scenarioInput}
          goalProfile={consoleState.goalProfile}
          goalMessage={consoleState.goalMessage}
          setGoalMessage={consoleState.setGoalMessage}
          runningGoalAction={consoleState.runningGoalAction}
          creatingScenario={consoleState.creatingScenario}
          generatingLessons={consoleState.generatingLessons}
          onGoalAction={consoleState.runScenarioGoalAction}
        />

        {consoleState.selectedScenario && (
          <LessonGenerationCard
            selectedScenario={consoleState.selectedScenario}
            generatedLessons={consoleState.generatedLessons}
            recommendationComparison={consoleState.recommendationComparison}
            lessonRecommendationQueries={consoleState.lessonRecommendationQueries}
            labelsByCandidateKey={consoleState.labelsByCandidateKey}
            savingLabelKey={consoleState.savingLabelKey}
            savingExternalCandidateKey={consoleState.savingExternalCandidateKey}
            generatingLessons={consoleState.generatingLessons}
            runningComparison={consoleState.runningComparison}
            runningComparisonLessonId={consoleState.runningComparisonLessonId}
            onGenerateLessons={() => consoleState.generateScenarioLessons()}
            onUpdateLessonRecommendationQuery={consoleState.updateLessonRecommendationQuery}
            onCompareLesson={consoleState.compareScenarioRecommendations}
            onSaveCandidateLabel={consoleState.saveCandidateLabel}
            onSaveExternalCandidate={consoleState.saveExternalCandidateContent}
          />
        )}

        {consoleState.selectedScenario && (
          <LabelSummaryCard
            labelSummary={consoleState.labelSummary}
            exportingLabelFormat={consoleState.exportingLabelFormat}
            onExportLabelDataset={consoleState.exportLabelDataset}
            onRefreshLabelSummary={() => consoleState.loadLabelSummary(consoleState.selectedScenario!.id)}
          />
        )}
      </div>
    </main>
  )
}
