'use client'

import type { CSSProperties } from 'react'
import { ResultPointCard } from './ResultPointCard'
import type { PlanetResultLesson } from './resultTypes'

const sectionStyle: CSSProperties = {
  display: 'grid',
  gap: 12,
}

const headerStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 12,
  paddingBottom: 4,
  borderBottom: '1px solid rgba(96, 67, 32, 0.13)',
}

const titleStyle: CSSProperties = {
  margin: 0,
  fontSize: 18,
  color: '#211407',
}

const countStyle: CSSProperties = {
  padding: '4px 9px',
  borderRadius: 999,
  background: 'rgba(93, 67, 28, 0.10)',
  color: '#5c421b',
  fontSize: 12,
  fontWeight: 800,
}

export function ResultLessonSection({ lesson }: { lesson: PlanetResultLesson }) {
  const artifactCount = lesson.points.reduce((sum, point) => sum + point.artifacts.length, 0)

  return (
    <section style={sectionStyle}>
      <div style={headerStyle}>
        <h3 style={titleStyle}>{lesson.lesson_title}</h3>
        <span style={countStyle}>결과물 {artifactCount}개</span>
      </div>
      {lesson.points.map((point, index) => (
        <ResultPointCard key={point.point_id} point={point} lessonTitle={lesson.lesson_title} pointIndex={index + 1} />
      ))}
    </section>
  )
}
