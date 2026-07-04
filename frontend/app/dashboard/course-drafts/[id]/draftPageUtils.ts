// ── 최소 타입 ─────────────────────────────────────────────────────────────────

type GoalProfileForPills = {
  usage_context?: string | null
  motivation?: string | null
  version: number
  rebuild_decision?: string | null
} | null

type LessonPointStat = {
  points?: { point: { status: string; point_type: string } }[]
  sub_lessons?: LessonPointStat[]
}

// ── 유틸리티 함수 ──────────────────────────────────────────────────────────────

export function formatPlanetExplorationPlanTitle(title?: string | null) {
  const trimmed = title?.trim()
  return trimmed ? `행성 ${trimmed} 탐험 계획(학습계획)` : '행성 탐험 계획(학습계획)'
}

export function formatPlanetExplorationPlanBaseTitle(title?: string | null) {
  const trimmed = title?.trim()
  return trimmed ? `행성 ${trimmed} 탐험 계획` : '행성 탐험 계획'
}

export function getGoalMetaPills(profile: GoalProfileForPills) {
  if (!profile) return []
  const pills: string[] = []
  if (profile.usage_context) pills.push(`상황 ${profile.usage_context}`)
  if (profile.motivation) pills.push(`이유 ${profile.motivation}`)
  pills.push(`Goal v${profile.version}.0`)
  if (profile.rebuild_decision === 'keep_structure') pills.push('선택: 학습자 편집')
  if (profile.rebuild_decision === 'rebuild_remaining' || profile.rebuild_decision === 'rebuild_all') {
    pills.push('선택: 모두 재구성')
  }
  return pills
}

export function countLessonPointStats(tree: LessonPointStat) {
  const stats = {
    totalPoints: 0,
    learningPoints: 0,
    explorationPoints: 0,
    completedExplorationPoints: 0,
    researchPoints: 0,
    completedResearchPoints: 0,
  }

  const visit = (node: LessonPointStat) => {
    ;(node.points ?? []).forEach(({ point }) => {
      stats.totalPoints += 1
      if (point.status === 'learning') stats.learningPoints += 1
      if (point.point_type === 'exploration') {
        stats.explorationPoints += 1
        if (point.status === 'completed') stats.completedExplorationPoints += 1
      }
      if (point.point_type === 'research') {
        stats.researchPoints += 1
        if (point.status === 'completed') stats.completedResearchPoints += 1
      }
    })
    ;(node.sub_lessons ?? []).forEach(visit)
  }

  visit(tree)
  return stats
}

export function countDraftLessonProgress(lessons: LessonPointStat[]) {
  return lessons.reduce(
    (acc, lessonTree) => {
      const stats = countLessonPointStats(lessonTree)
      acc.total += 1
      acc.totalPoints += stats.totalPoints
      acc.explorationPoints += stats.explorationPoints
      acc.completedExplorationPoints += stats.completedExplorationPoints
      acc.researchPoints += stats.researchPoints
      acc.completedResearchPoints += stats.completedResearchPoints
      const completedPoints = stats.completedExplorationPoints + stats.completedResearchPoints
      if (stats.totalPoints > 0 && completedPoints === stats.totalPoints) {
        acc.completed += 1
      } else if (stats.learningPoints > 0 || completedPoints > 0) {
        acc.learning += 1
      }
      return acc
    },
    {
      total: 0,
      completed: 0,
      learning: 0,
      totalPoints: 0,
      explorationPoints: 0,
      completedExplorationPoints: 0,
      researchPoints: 0,
      completedResearchPoints: 0,
    },
  )
}

export function getJournalStatusLabel(status?: string) {
  if (status === 'learning') return '탐험중(learning)'
  if (status === 'confirmed') return 'confirmed'
  if (status === 'archived') return '탐험완료(archived)'
  return '탐험계획중(draft)'
}
