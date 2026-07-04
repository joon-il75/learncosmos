export interface RecordJournalEntry {
  observation: string
  reflection: string
  next_step: string
  updated_at?: string
}

export interface RecordEntry {
  study_minutes: number
  practice_count: number
  confidence_level: number
  application_note: string
  updated_at?: string
}

export interface RecordQuestion {
  id: string
  question: string
  question_type: string
  answer?: string | null
  ai_feedback?: string | null
  created_by?: string
  status?: string
}

export interface RecordSelfEvaluation {
  understanding: number
  application_note: string
  proficiency: number
  goal_alignment_note: string
  updated_at?: string
}

export interface RecordBlock {
  id: string
  block_type: string
  content: unknown
  order_index: number
}

export interface PlanetRecordPoint {
  point_id: string
  point_title: string
  point_type: 'exploration' | 'research'
  blocks?: RecordBlock[]
  journal_entry?: RecordJournalEntry | null
  application_note: string
  self_evaluation_application_note: string
  goal_alignment_note: string
  questions?: RecordQuestion[]
  record_entry?: RecordEntry | null
  self_evaluation?: RecordSelfEvaluation | null
}

export interface PlanetRecordLesson {
  lesson_id: string
  lesson_title: string
  points: PlanetRecordPoint[]
}

export interface PlanetRecordAggregate {
  course: {
    id: string
    title: string
    total_points: number
    recorded_points: number
  }
  lessons: PlanetRecordLesson[]
}

export type RecordFeedFilterKey = 'all' | 'completed' | 'work' | 'artifact' | 'share_candidate'

export interface PlanetRecordFeedResponse {
  route_kind: 'learning' | 'shared'
  course: {
    id: string
    title: string
    status: 'learning' | 'completed'
    total_points: number
    completed_points: number
    total_lessons: number
    completed_lessons: number
    record_card_count: number
    share_candidate_count: number
  }
  filters: Array<{
    key: RecordFeedFilterKey
    label: string
    count: number
  }>
  cards: PlanetRecordCard[]
}

export interface PlanetRecordCard {
  id: string
  card_type:
    | 'point_completed'
    | 'lesson_completed'
    | 'course_completed'
    | 'artifact_submitted'
    | 'question_answered'
    | 'self_evaluation_saved'
    | 'practice_log_added'
    | 'journal_saved'
  title: string
  summary: string
  badge: string
  category: 'completed' | 'work' | 'artifact' | 'share_candidate'
  visibility: 'private' | 'share_candidate' | 'shared'
  course_id: string
  lesson_id?: string
  point_id?: string
  lesson_title?: string
  point_title?: string
  point_type?: 'exploration' | 'research'
  occurred_at: string
  share_candidate_score: number
  share_text: string
  detail: {
    primary_text?: string
    journal_excerpt?: string
    question_excerpt?: string
    artifact_title?: string
    self_evaluation_summary?: string
    practice_summary?: string
    lesson_progress?: {
      completed_points: number
      total_points: number
      percent: number
    }
  }
  actions: {
    can_open_learning_page: boolean
    can_copy_share_text: boolean
    can_create_share_card: boolean
  }
}
