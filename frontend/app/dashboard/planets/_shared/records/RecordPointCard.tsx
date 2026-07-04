'use client'

import type { CSSProperties } from 'react'
import { RecordBlocksCollapsible } from './RecordBlocksCollapsible'
import type { PlanetRecordPoint } from './recordTypes'

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

const sectionStyle: CSSProperties = {
  display: 'grid',
  gap: 8,
  padding: 12,
  borderRadius: 8,
  background: 'rgba(255, 255, 255, 0.58)',
  border: '1px solid rgba(122, 90, 36, 0.10)',
}

const sectionTitleStyle: CSSProperties = {
  margin: 0,
  fontSize: 13,
  fontWeight: 900,
  color: '#5e3d13',
}

const textGridStyle: CSSProperties = {
  display: 'grid',
  gap: 6,
}

const fieldStyle: CSSProperties = {
  display: 'grid',
  gap: 3,
}

const labelStyle: CSSProperties = {
  fontSize: 11,
  fontWeight: 800,
  color: 'rgba(83, 55, 18, 0.66)',
}

const valueStyle: CSSProperties = {
  margin: 0,
  whiteSpace: 'pre-wrap',
  wordBreak: 'break-word',
  fontSize: 13,
  lineHeight: 1.55,
  color: 'rgba(38, 25, 11, 0.84)',
}

const metricGridStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(120px, 1fr))',
  gap: 8,
}

const metricStyle: CSSProperties = {
  display: 'grid',
  gap: 3,
  padding: 10,
  borderRadius: 8,
  background: 'rgba(246, 232, 204, 0.42)',
}

const emptyStyle: CSSProperties = {
  margin: 0,
  fontSize: 13,
  lineHeight: 1.5,
  color: 'rgba(76, 54, 25, 0.62)',
}

function hasText(value?: string | null): boolean {
  return Boolean(value?.trim())
}

function Field({ label, value }: { label: string; value?: string | null }) {
  return (
    <div style={fieldStyle}>
      <span style={labelStyle}>{label}</span>
      <p style={valueStyle}>{hasText(value) ? value : '아직 기록이 없습니다.'}</p>
    </div>
  )
}

export function RecordPointCard({
  point,
  lessonTitle,
  pointIndex,
}: {
  point: PlanetRecordPoint
  lessonTitle: string
  pointIndex: number
}) {
  const journal = point.journal_entry
  const record = point.record_entry
  const selfEvaluation = point.self_evaluation
  const questions = point.questions ?? []
  const blocks = point.blocks ?? []

  return (
    <article style={cardStyle}>
      <div style={headerStyle}>
        <div style={{ display: 'grid', gap: 4, minWidth: 0 }}>
          <h4 style={titleStyle}>{point.point_title}</h4>
          <span style={labelStyle}>{lessonTitle} · 지점 {pointIndex}</span>
        </div>
        <span style={badgeStyle}>{point.point_type === 'research' ? '연구지점' : '탐험지점'}</span>
      </div>

      <section style={sectionStyle}>
        <h5 style={sectionTitleStyle}>내용정리</h5>
        {point.point_type === 'research' ? <RecordBlocksCollapsible blocks={blocks} /> : null}
        <div style={textGridStyle}>
          <Field label="관찰한 내용" value={journal?.observation} />
          <Field label="느낀 점 / 배운 점" value={journal?.reflection} />
          <Field label="다음에 해볼 것" value={journal?.next_step} />
          <Field label="적용 메모" value={point.application_note || point.self_evaluation_application_note} />
          <Field label="목표 연결 메모" value={point.goal_alignment_note} />
        </div>
      </section>

      <section style={sectionStyle}>
        <h5 style={sectionTitleStyle}>질문관리</h5>
        {questions.length > 0 ? (
          <div style={textGridStyle}>
            {questions.map((question) => (
              <div key={question.id} style={fieldStyle}>
                <span style={labelStyle}>
                  {question.created_by === 'ai' ? 'AI 질문' : '학습자 질문'} · {question.question_type}
                </span>
                <p style={valueStyle}>{question.question}</p>
                <Field label="답변" value={question.answer} />
                {hasText(question.ai_feedback) ? <Field label="AI 참고 피드백" value={question.ai_feedback} /> : null}
              </div>
            ))}
          </div>
        ) : (
          <p style={emptyStyle}>아직 질문 기록이 없습니다.</p>
        )}
      </section>

      <section style={sectionStyle}>
        <h5 style={sectionTitleStyle}>연습/활동기록</h5>
        {record || selfEvaluation ? (
          <>
            <div style={metricGridStyle}>
              <div style={metricStyle}>
                <span style={labelStyle}>학습시간</span>
                <strong>{record ? `${record.study_minutes}분` : '미기록'}</strong>
              </div>
              <div style={metricStyle}>
                <span style={labelStyle}>반복횟수</span>
                <strong>{record ? `${record.practice_count}회` : '미기록'}</strong>
              </div>
              <div style={metricStyle}>
                <span style={labelStyle}>자신감</span>
                <strong>{record ? `${record.confidence_level}/5` : '미기록'}</strong>
              </div>
              <div style={metricStyle}>
                <span style={labelStyle}>이해도</span>
                <strong>{selfEvaluation ? `${selfEvaluation.understanding}/5` : '미기록'}</strong>
              </div>
              <div style={metricStyle}>
                <span style={labelStyle}>적용 자신감</span>
                <strong>{selfEvaluation ? `${selfEvaluation.proficiency}/5` : '미기록'}</strong>
              </div>
            </div>
            <Field label="실제 적용 메모" value={record?.application_note} />
          </>
        ) : (
          <p style={emptyStyle}>아직 연습/활동기록이 없습니다.</p>
        )}
      </section>
    </article>
  )
}
