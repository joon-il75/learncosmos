'use client'

import type { CSSProperties } from 'react'
import type { PlanetResultPoint } from './resultTypes'

const cardStyle: CSSProperties = {
  display: 'grid',
  gap: 12,
  padding: 16,
  borderRadius: 8,
  border: '1px solid rgba(104, 72, 35, 0.16)',
  background: 'rgba(255, 251, 242, 0.86)',
  boxShadow: '0 8px 18px rgba(55, 35, 12, 0.06)',
}

const headerStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: 12,
}

const titleStyle: CSSProperties = {
  margin: 0,
  fontSize: 16,
  lineHeight: 1.35,
  color: '#26190b',
}

const badgeStyle: CSSProperties = {
  flex: '0 0 auto',
  padding: '5px 9px',
  borderRadius: 999,
  background: 'rgba(139, 92, 36, 0.11)',
  color: '#704816',
  fontSize: 11,
  fontWeight: 900,
}

const listStyle: CSSProperties = {
  display: 'grid',
  gap: 8,
}

const artifactStyle: CSSProperties = {
  display: 'grid',
  gap: 6,
  padding: 12,
  borderRadius: 8,
  background: 'rgba(255, 255, 255, 0.62)',
  border: '1px solid rgba(122, 90, 36, 0.10)',
}

const artifactTitleStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 10,
  color: '#2b1c0c',
  fontSize: 14,
  fontWeight: 900,
}

const typeStyle: CSSProperties = {
  padding: '4px 8px',
  borderRadius: 999,
  background: 'rgba(176, 117, 39, 0.12)',
  color: '#6f4612',
  fontSize: 11,
  fontWeight: 900,
}

const linkStyle: CSSProperties = {
  color: '#8b4f12',
  fontSize: 13,
  fontWeight: 800,
  wordBreak: 'break-all',
}

const descriptionStyle: CSSProperties = {
  margin: 0,
  whiteSpace: 'pre-wrap',
  wordBreak: 'break-word',
  fontSize: 13,
  lineHeight: 1.55,
  color: 'rgba(38, 25, 11, 0.78)',
}

const emptyStyle: CSSProperties = {
  margin: 0,
  padding: 12,
  borderRadius: 8,
  background: 'rgba(255,255,255,0.50)',
  color: 'rgba(76, 54, 25, 0.62)',
  fontSize: 13,
}

export function ResultPointCard({
  point,
  lessonTitle,
  pointIndex,
}: {
  point: PlanetResultPoint
  lessonTitle: string
  pointIndex: number
}) {
  return (
    <article style={cardStyle}>
      <div style={headerStyle}>
        <div style={{ display: 'grid', gap: 4, minWidth: 0 }}>
          <h4 style={titleStyle}>{point.point_title}</h4>
          <span style={descriptionStyle}>{lessonTitle} · 지점 {pointIndex}</span>
        </div>
        <span style={badgeStyle}>{point.point_type === 'research' ? '연구지점' : '탐험지점'}</span>
      </div>

      {point.artifacts.length > 0 ? (
        <div style={listStyle}>
          {point.artifacts.map((artifact) => (
            <div key={artifact.id} style={artifactStyle}>
              <div style={artifactTitleStyle}>
                <span>{artifact.title.trim() || '제목 없는 결과물'}</span>
                <span style={typeStyle}>{artifact.artifact_type.trim() || '결과물'}</span>
              </div>
              {artifact.url.trim() ? (
                <a href={artifact.url} target="_blank" rel="noreferrer" style={linkStyle}>
                  결과물 링크 열기
                </a>
              ) : null}
              <p style={descriptionStyle}>{artifact.description.trim() || '설명이 아직 없습니다.'}</p>
            </div>
          ))}
        </div>
      ) : (
        <p style={emptyStyle}>아직 결과물이 없습니다.</p>
      )}
    </article>
  )
}
