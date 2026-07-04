import { useState } from 'react'
import { S } from '../styles'
import { LABEL_OPTIONS } from '../constants'
import type { DebugLabel, DebugScenario, GeneratedLessonItem, GeneratedLessonsSnapshot, RecommendationCompareLesson, RecommendationComparison } from '../types'
import { candidateLabelKey } from '../utils'

const CANDIDATES_PER_PAGE = 5

type Props = {
  selectedScenario: DebugScenario
  generatedLessons: GeneratedLessonsSnapshot | null
  recommendationComparison: RecommendationComparison | null
  lessonRecommendationQueries: Record<string, string>
  labelsByCandidateKey: Record<string, DebugLabel>
  savingLabelKey: string | null
  savingExternalCandidateKey: string | null
  generatingLessons: boolean
  runningComparison: boolean
  runningComparisonLessonId: string | null
  onGenerateLessons: () => void
  onUpdateLessonRecommendationQuery: (lessonId: string, value: string) => void
  onCompareLesson: (lesson: GeneratedLessonItem) => void
  onSaveCandidateLabel: (
    lesson: RecommendationCompareLesson['lesson'],
    candidate: RecommendationCompareLesson['baseline']['candidates'][number],
    candidateIndex: number,
    label: string,
    runIDOverride?: string,
  ) => void
  onSaveExternalCandidate: (
    lesson: RecommendationCompareLesson['lesson'],
    candidate: RecommendationCompareLesson['baseline']['candidates'][number],
    candidateIndex: number,
  ) => void
}

export default function LessonGenerationCard({
  selectedScenario,
  generatedLessons,
  recommendationComparison,
  lessonRecommendationQueries,
  labelsByCandidateKey,
  savingLabelKey,
  savingExternalCandidateKey,
  generatingLessons,
  runningComparison,
  runningComparisonLessonId,
  onGenerateLessons,
  onUpdateLessonRecommendationQuery,
  onCompareLesson,
  onSaveCandidateLabel,
  onSaveExternalCandidate,
}: Props) {
  const [candidatePages, setCandidatePages] = useState<Record<string, number>>({})

  const changeCandidatePage = (lessonId: string, page: number) => {
    setCandidatePages(prev => ({ ...prev, [lessonId]: page }))
  }

  return (
    <div style={S.card}>
      <div style={{ display:'flex', justifyContent:'space-between', gap:'12px', alignItems:'flex-start', flexWrap:'wrap', marginBottom:'14px' }}>
        <div>
          <div style={{ fontSize:'15px', fontWeight:800 }}>3. 리슨 생성</div>
          <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.58)', marginTop:'5px', lineHeight:1.7 }}>
            확정된 목표 프로필을 사용해 실제 코스 생성과 같은 리슨 구조를 만들되, 사용자 draft는 생성하지 않습니다.
          </div>
        </div>
        <button
          style={S.primary}
          onClick={onGenerateLessons}
          disabled={generatingLessons || selectedScenario.status === 'draft'}
        >
          {generatingLessons ? '생성 중...' : '리슨 생성'}
        </button>
      </div>

      {selectedScenario.status === 'draft' && (
        <div style={S.warn}>목표 확정 후 리슨 생성을 실행할 수 있습니다.</div>
      )}

      {generatedLessons ? (
        <div style={{ display:'grid', gap:'14px' }}>
          <div style={S.info}>
            <div style={{ fontWeight:800, marginBottom:'6px' }}>{generatedLessons.draft_title}</div>
            <div style={{ marginBottom:'8px' }}>{generatedLessons.draft_description}</div>
            <div>provider: {generatedLessons.provider} · model: {generatedLessons.model} · latency: {generatedLessons.latency_ms}ms · lessons: {generatedLessons.lesson_count}</div>
          </div>
          <div style={{ display:'grid', gap:'10px' }}>
            {generatedLessons.lessons.map((lesson, index) => {
              const lessonComparison = recommendationComparison?.lessons.find(item => item.lesson.lesson_id === lesson.lesson_id)
              return (
                <div key={lesson.lesson_id} style={{ background:'rgba(255,255,255,0.03)', border:'1px solid rgba(120,140,200,0.12)', borderRadius:'12px', padding:'12px' }}>
                  <div style={{ display:'flex', justifyContent:'space-between', gap:'10px', alignItems:'flex-start', flexWrap:'wrap', marginBottom:'6px' }}>
                    <div>
                      <div style={{ fontSize:'11px', color:'rgba(200,210,235,0.5)', marginBottom:'4px' }}>LESSON {index + 1} · {lesson.source_type}</div>
                      <div style={{ fontSize:'15px', fontWeight:800 }}>{lesson.title}</div>
                    </div>
                  </div>
                  <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.72)', lineHeight:1.65 }}>{lesson.objective}</div>
                  <div style={{ marginTop:'10px', display:'grid', gridTemplateColumns:'minmax(0, 1fr)', gap:'8px', alignItems:'start' }}>
                    <label>
                      <span style={S.label}>추천 기준</span>
                      <input
                        value={lessonRecommendationQueries[lesson.lesson_id] ?? ''}
                        onChange={event => onUpdateLessonRecommendationQuery(lesson.lesson_id, event.target.value)}
                        placeholder="예: 시장조사 강의 주제 선정, 유튜브 강의 촬영 편집"
                        style={{ ...S.input, padding:'9px 11px', fontSize:'12px' }}
                      />
                    </label>
                    <button
                      style={{ ...S.secondary, padding:'8px 11px', fontSize:'11px', justifySelf:'start' }}
                      onClick={() => onCompareLesson(lesson)}
                      disabled={runningComparison || !['lessons_generated', 'recommendation_tested', 'labelled'].includes(selectedScenario.status)}
                    >
                      {runningComparisonLessonId === lesson.lesson_id ? '찾는 중...' : '추천 콘텐츠 찾기'}
                    </button>
                  </div>
                  {lessonComparison && (
                    <div style={{ marginTop:'12px', borderTop:'1px solid rgba(120,140,200,0.12)', paddingTop:'12px' }}>
                      <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(190px, 1fr))', gap:'10px', marginBottom:'10px' }}>
                        <div style={S.info}>
                          <div style={{ fontWeight:800, marginBottom:'4px' }}>추천 후보</div>
                          <div>후보 {lessonComparison.baseline.candidates.length}개</div>
                          <div>lexical {lessonComparison.baseline.stages.search.lexical_hit_count} · vector {lessonComparison.baseline.stages.search.vector_hit_count}</div>
                        </div>
                        <div style={lessonComparison.shadow.enabled ? S.info : S.warn}>
                          <div style={{ fontWeight:800, marginBottom:'4px' }}>Embedding route</div>
                          <div>{lessonComparison.shadow.enabled ? 'EmbeddingGemma 응답 확인' : 'EmbeddingGemma 확인 대기'}</div>
                          <div>{lessonComparison.shadow.reason}</div>
                          {lessonComparison.shadow.ranker && (
                            <div style={{ marginTop:'4px' }}>
                              ranker: {lessonComparison.shadow.ranker.provider || 'noop'} · {lessonComparison.shadow.ranker.candidate_count ?? 0}개 · {lessonComparison.shadow.ranker.reason}
                            </div>
                          )}
                        </div>
                      </div>

                      <div style={{ display:'grid', gap:'8px' }}>
                        {lessonComparison.baseline.candidates.length === 0 ? (
                          <div style={S.warn}>추천 후보가 없습니다.</div>
                        ) : (() => {
                          const totalCandidates = lessonComparison.baseline.candidates.length
                          const totalPages = Math.max(1, Math.ceil(totalCandidates / CANDIDATES_PER_PAGE))
                          const currentPage = Math.min(Math.max(candidatePages[lessonComparison.lesson.lesson_id] ?? 1, 1), totalPages)
                          const startIndex = (currentPage - 1) * CANDIDATES_PER_PAGE
                          const visibleCandidates = lessonComparison.baseline.candidates.slice(startIndex, startIndex + CANDIDATES_PER_PAGE)
                          return (
                            <>
                              {visibleCandidates.map((candidate, visibleIndex) => {
                                const candidateIndex = startIndex + visibleIndex
                                const labelKey = candidateLabelKey(lessonComparison.lesson.lesson_id, candidate, candidateIndex)
                                const savedLabel = labelsByCandidateKey[labelKey]
                                const isExternalCandidate = candidate.source_type === 'external' || candidate.resource_type === 'external'
                                return (
                                  <div key={candidate.content_id || candidate.external_url || `${candidate.title}-${candidateIndex}`} style={{ background:'rgba(255,255,255,0.025)', border:'1px solid rgba(120,140,200,0.1)', borderRadius:'10px', padding:'10px' }}>
                                    <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.5)' }}>#{candidateIndex + 1} · score {candidate.rank_score?.toFixed?.(4) ?? candidate.rank_score}</div>
                                    <div style={{ fontSize:'13px', fontWeight:800, marginTop:'3px' }}>{candidate.title}</div>
                                    {candidate.selection_reason && (
                                      <div style={{ fontSize:'11px', color:'rgba(200,210,235,0.62)', marginTop:'3px', lineHeight:1.55 }}>{candidate.selection_reason}</div>
                                    )}
                                    {candidate.external_url && (
                                      <a href={candidate.external_url} target="_blank" rel="noreferrer" style={{ display:'inline-block', color:'#9BC3FF', fontSize:'11px', marginTop:'5px' }}>
                                        원문 열기
                                      </a>
                                    )}
                                    {isExternalCandidate && (
                                      <div style={{ marginTop:'8px' }}>
                                        <button
                                          type="button"
                                          onClick={() => onSaveExternalCandidate(lessonComparison.lesson, candidate, candidateIndex)}
                                          disabled={candidate.already_saved || candidate.eligible_for_save === false || savingExternalCandidateKey === labelKey}
                                          style={{ ...S.secondary, padding:'7px 10px', fontSize:'11px' }}
                                        >
                                          {candidate.already_saved ? '콘텐츠 저장됨' : savingExternalCandidateKey === labelKey ? '저장 중...' : '콘텐츠로 저장'}
                                        </button>
                                      </div>
                                    )}
                                    <div style={{ marginTop:'8px', display:'flex', gap:'6px', flexWrap:'wrap' }}>
                                      {LABEL_OPTIONS.map(option => {
                                        const selected = savedLabel?.label === option.value
                                        return (
                                          <button
                                            key={option.value}
                                            type="button"
                                            onClick={() => onSaveCandidateLabel(lessonComparison.lesson, candidate, candidateIndex, option.value, lessonComparison.run_id)}
                                            disabled={!(lessonComparison.run_id || recommendationComparison?.run_id) || savingLabelKey === labelKey}
                                            style={{
                                              border:selected ? '1px solid rgba(142,229,154,0.65)' : '1px solid rgba(120,140,200,0.16)',
                                              background:selected ? 'rgba(142,229,154,0.14)' : 'rgba(255,255,255,0.04)',
                                              color:selected ? '#8EE59A' : 'rgba(200,210,235,0.76)',
                                              borderRadius:'999px',
                                              padding:'6px 9px',
                                              fontSize:'11px',
                                              fontWeight:700,
                                              cursor:'pointer',
                                            }}
                                          >
                                            {savingLabelKey === labelKey && !selected ? '저장 중...' : option.label}
                                          </button>
                                        )
                                      })}
                                    </div>
                                    {savedLabel && (
                                      <div style={{ fontSize:'11px', color:'#8EE59A', marginTop:'6px' }}>
                                        저장됨: {LABEL_OPTIONS.find(option => option.value === savedLabel.label)?.label ?? savedLabel.label}
                                      </div>
                                    )}
                                  </div>
                                )
                              })}
                              {totalPages > 1 && (
                                <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', gap:'10px', flexWrap:'wrap', marginTop:'4px' }}>
                                  <div style={{ fontSize:'11px', color:'rgba(200,210,235,0.58)' }}>
                                    {currentPage} / {totalPages} 페이지 · {CANDIDATES_PER_PAGE}개씩 표시
                                  </div>
                                  <div style={{ display:'flex', gap:'6px' }}>
                                    <button
                                      type="button"
                                      style={{ ...S.secondary, padding:'7px 10px', fontSize:'11px' }}
                                      onClick={() => changeCandidatePage(lessonComparison.lesson.lesson_id, currentPage - 1)}
                                      disabled={currentPage <= 1}
                                    >
                                      이전
                                    </button>
                                    <button
                                      type="button"
                                      style={{ ...S.secondary, padding:'7px 10px', fontSize:'11px' }}
                                      onClick={() => changeCandidatePage(lessonComparison.lesson.lesson_id, currentPage + 1)}
                                      disabled={currentPage >= totalPages}
                                    >
                                      다음
                                    </button>
                                  </div>
                                </div>
                              )}
                            </>
                          )
                        })()}
                      </div>
                    </div>
                  )}
                </div>
              )
            })}
          </div>
          {generatedLessons.completion_criteria.length > 0 && (
            <div style={S.info}>
              <div style={{ fontWeight:800, marginBottom:'6px' }}>완료 기준</div>
              {generatedLessons.completion_criteria.map((item, index) => (
                <div key={`${item}-${index}`}>- {item}</div>
              ))}
            </div>
          )}
        </div>
      ) : selectedScenario.status !== 'draft' ? (
        <div style={S.info}>아직 이 화면에서 생성한 리슨 snapshot이 없습니다. 테스트를 실행하면 결과가 바로 표시됩니다.</div>
      ) : null}

      <div style={{ ...S.info, marginTop:'14px' }}>
        다음 단계에서는 생성된 리슨 snapshot을 기준으로 primary 추천 결과와 ranker preview를 나란히 비교합니다.
      </div>
    </div>
  )
}
