export interface PlanetResultArtifact {
  id: string
  course_point_id: string
  artifact_type: string
  title: string
  url: string
  description: string
  order_index: number
  created_at?: string
  updated_at?: string
}

export interface PlanetResultPoint {
  point_id: string
  point_title: string
  point_type: 'exploration' | 'research'
  artifacts: PlanetResultArtifact[]
}

export interface PlanetResultLesson {
  lesson_id: string
  lesson_title: string
  points: PlanetResultPoint[]
}

export interface PlanetResultAggregate {
  course: {
    id: string
    title: string
    total_points: number
    points_with_artifacts: number
  }
  lessons: PlanetResultLesson[]
}
