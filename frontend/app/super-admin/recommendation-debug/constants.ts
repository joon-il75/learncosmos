export const API_BASE = ''

export const LABEL_OPTIONS = [
  { value:'good_fit', label:'적합' },
  { value:'goal_fit_stage_weak', label:'목표 적합/단계 약함' },
  { value:'stage_fit_goal_weak', label:'단계 적합/목표 약함' },
  { value:'irrelevant', label:'무관' },
  { value:'duplicate', label:'중복' },
  { value:'broken_link', label:'링크 문제' },
  { value:'low_quality', label:'품질 낮음' },
  { value:'unsafe', label:'부적절' },
] as const

export const displayRolloutMode = (mode?: string | null) => {
  switch (mode) {
    case 'shadow_only':
      return '점검 준비'
    case 'admin_preview':
      return 'admin preview'
    case 'learner_10_percent':
      return 'learner 10%'
    case 'learner_50_percent':
      return 'learner 50%'
    case 'learner_100_percent':
      return 'learner 100%'
    case 'rolled_back':
      return 'rolled back'
    default:
      return mode || '없음'
  }
}

export const displayEmbeddingRoute = (route?: string | null) => {
  switch (route) {
    case 'shadow_vector':
      return 'primary vector'
    case 'embedding_vector':
      return 'embedding vector'
    default:
      return route || 'primary vector'
  }
}

export const normalizeEmbeddingReason = (reason?: string | null) => {
  if (!reason) return ''
  return reason
    .replaceAll('shadow vector', 'primary vector')
    .replaceAll('shadow_embedding', 'primary_embedding')
    .replaceAll('shadow', 'primary')
    .replaceAll('baseline route remains active', 'primary route keeps current ordering')
}
