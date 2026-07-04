import { displayEmbeddingRoute, LABEL_OPTIONS, normalizeEmbeddingReason } from '../constants'
import { S } from '../styles'
import type { DebugLabel, DebugScenario, RecommendationCompareLesson, RecommendationComparison } from '../types'
import { candidateLabelKey } from '../utils'

type Props = {
  selectedScenario: DebugScenario
  recommendationComparison: RecommendationComparison | null
  labelsByCandidateKey: Record<string, DebugLabel>
  savingLabelKey: string | null
  onSaveCandidateLabel: (
    lesson: RecommendationCompareLesson['lesson'],
    candidate: RecommendationCompareLesson['baseline']['candidates'][number],
    candidateIndex: number,
    label: string,
    runIDOverride?: string,
  ) => void
}

export default function RecommendationComparisonCard({
  selectedScenario,
  recommendationComparison,
  labelsByCandidateKey,
  savingLabelKey,
  onSaveCandidateLabel,
}: Props) {
  return (
    <div style={S.card}>
      <div style={{ display:'flex', justifyContent:'space-between', gap:'12px', alignItems:'flex-start', flexWrap:'wrap', marginBottom:'14px' }}>
        <div>
          <div style={{ fontSize:'15px', fontWeight:800 }}>4. 리슨별 추천 콘텐츠</div>
          <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.58)', marginTop:'5px', lineHeight:1.7 }}>
            3단계의 각 리슨에서 `추천 콘텐츠 찾기`를 실행하면 해당 리슨 기준 추천 결과가 여기에 쌓입니다.
          </div>
        </div>
      </div>

      {!['lessons_generated', 'recommendation_tested', 'labelled'].includes(selectedScenario.status) && (
        <div style={S.warn}>리슨 생성 후 각 리슨의 `추천 콘텐츠 찾기`를 실행할 수 있습니다.</div>
      )}

      {recommendationComparison ? (
        <div style={{ display:'grid', gap:'14px' }}>
          <div style={S.info}>
            <div style={{ fontWeight:800, marginBottom:'6px' }}>{recommendationComparison.course_title}</div>
            <div style={{ marginBottom:'8px' }}>{recommendationComparison.confirmed_goal}</div>
            <div>
              lessons: {recommendationComparison.lesson_count} · candidates: {recommendationComparison.total_candidates} · lexical: {recommendationComparison.total_lexical_hits} · vector: {recommendationComparison.total_vector_hits} · latency: {recommendationComparison.latency_ms}ms
            </div>
            <div style={{ marginTop:'6px' }}>
              embedding check: attempted {recommendationComparison.shadow_attempted ?? 0} · succeeded {recommendationComparison.shadow_succeeded ?? 0}
            </div>
          </div>
          <div style={{ display:'grid', gap:'12px' }}>
            {recommendationComparison.lessons.map((item, index) => (
              <div key={item.lesson.lesson_id || `${item.lesson.title}-${index}`} style={{ background:'rgba(255,255,255,0.03)', border:'1px solid rgba(120,140,200,0.12)', borderRadius:'12px', padding:'14px' }}>
                <div style={{ fontSize:'11px', color:'rgba(200,210,235,0.5)', marginBottom:'4px' }}>LESSON {index + 1}</div>
                <div style={{ fontSize:'15px', fontWeight:800, marginBottom:'6px' }}>{item.lesson.title}</div>
                <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.68)', lineHeight:1.65, marginBottom:'10px' }}>{item.lesson.objective}</div>
                <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(190px, 1fr))', gap:'10px', marginBottom:'10px' }}>
                  <div style={S.info}>
                    <div style={{ fontWeight:800, marginBottom:'4px' }}>Primary 추천</div>
                    <div>후보 {item.baseline.candidates.length}개</div>
                    <div>lexical {item.baseline.stages.search.lexical_hit_count} · vector {item.baseline.stages.search.vector_hit_count}</div>
                  </div>
                  <div style={item.shadow.enabled ? S.info : S.warn}>
                    <div style={{ fontWeight:800, marginBottom:'4px' }}>Embedding route</div>
                    <div>{item.shadow.enabled ? 'EmbeddingGemma 응답 확인' : 'EmbeddingGemma 확인 대기'}</div>
                    <div>{normalizeEmbeddingReason(item.shadow.embedding?.reason || item.shadow.reason)}</div>
                    {item.shadow.embedding && (
                      <div style={{ marginTop:'4px' }}>
                        dim {item.shadow.embedding.dimension} · latency {item.shadow.embedding.latency_ms}ms · search {item.shadow.embedding.search_used ? 'yes' : 'no'}
                      </div>
                    )}
                    {item.shadow.search && (
                      <div style={{ marginTop:'4px' }}>
                        route {displayEmbeddingRoute(item.shadow.search.route)} · candidates {item.shadow.search.candidate_count}
                        <br />
                        {normalizeEmbeddingReason(item.shadow.search.reason)}
                        {item.shadow.search.fallback_reason && (
                          <>
                            <br />
                            fallback: {normalizeEmbeddingReason(item.shadow.search.fallback_reason)}
                          </>
                        )}
                      </div>
                    )}
                    {typeof item.metrics.shadow_overlap === 'number' && (
                      <div style={{ marginTop:'4px' }}>primary/vector overlap {(item.metrics.shadow_overlap * 100).toFixed(0)}%</div>
                    )}
                    {item.shadow.ranker && (
                      <div style={{ marginTop:'4px' }}>
                        ranker: {item.shadow.ranker.provider || 'noop'} · {item.shadow.ranker.candidate_count ?? 0}개 · {item.shadow.ranker.reason}
                      </div>
                    )}
                  </div>
                </div>
                {item.shadow.candidates && item.shadow.candidates.length > 0 && (
                  <div style={{ margin:'8px 0 12px', borderTop:'1px solid rgba(120,140,200,0.12)', paddingTop:'10px' }}>
                    <div style={{ fontSize:'12px', fontWeight:800, marginBottom:'7px', color:'rgba(232,234,242,0.9)' }}>Ranker preview 비교</div>
                    <div style={{ display:'grid', gap:'6px' }}>
                      {item.shadow.candidates.map((candidate, shadowIndex) => (
                        <div
                          key={candidate.content_id || candidate.external_url || `${candidate.title}-${shadowIndex}`}
                          style={{
                            display:'grid',
                            gridTemplateColumns:'72px minmax(0, 1fr) 120px',
                            gap:'10px',
                            alignItems:'center',
                            background:'rgba(255,255,255,0.025)',
                            border:'1px solid rgba(120,140,200,0.1)',
                            borderRadius:'10px',
                            padding:'8px 10px',
                          }}
                        >
                          <div style={{ fontSize:'11px', color:'rgba(200,210,235,0.58)', lineHeight:1.5 }}>
                            primary #{candidate.ranker_route_rank ?? shadowIndex + 1}
                            <br />
                            rerank #{candidate.ranker_rerank_rank ?? candidate.ranker_rank ?? shadowIndex + 1}
                          </div>
                          <div style={{ minWidth:0 }}>
                            <div style={{ fontSize:'12px', fontWeight:800, whiteSpace:'nowrap', overflow:'hidden', textOverflow:'ellipsis' }}>{candidate.title}</div>
                            <div style={{ fontSize:'11px', color:'rgba(200,210,235,0.52)', marginTop:'2px', whiteSpace:'nowrap', overflow:'hidden', textOverflow:'ellipsis' }}>
                              {candidate.ranker_reason || candidate.selection_reason || 'ranker score preview'}
                            </div>
                          </div>
                          <div style={{ fontSize:'11px', color:'rgba(200,210,235,0.66)', textAlign:'right' }}>
                            score {typeof candidate.ranker_score === 'number' ? candidate.ranker_score.toFixed(4) : 'n/a'}
                            <br />
                            delta {typeof candidate.ranker_rank_delta === 'number' ? candidate.ranker_rank_delta : 0}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
                <div style={{ display:'grid', gap:'8px' }}>
                  {item.baseline.candidates.length === 0 ? (
                    <div style={S.warn}>Primary 추천 후보가 없습니다.</div>
                  ) : item.baseline.candidates.map((candidate, candidateIndex) => {
                    const labelKey = candidateLabelKey(item.lesson.lesson_id, candidate, candidateIndex)
                    const savedLabel = labelsByCandidateKey[labelKey]
                    return (
                      <div key={candidate.content_id || candidate.external_url || `${candidate.title}-${candidateIndex}`} style={{ borderTop:'1px solid rgba(120,140,200,0.12)', paddingTop:'8px' }}>
                        <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.5)' }}>#{candidateIndex + 1} · score {candidate.rank_score?.toFixed?.(4) ?? candidate.rank_score}</div>
                        <div style={{ fontSize:'13px', fontWeight:800, marginTop:'3px' }}>{candidate.title}</div>
                        {candidate.selection_reason && (
                          <div style={{ fontSize:'11px', color:'rgba(200,210,235,0.62)', marginTop:'3px', lineHeight:1.55 }}>{candidate.selection_reason}</div>
                        )}
                        <div style={{ marginTop:'8px', display:'flex', gap:'6px', flexWrap:'wrap' }}>
                          {LABEL_OPTIONS.map(option => {
                            const selected = savedLabel?.label === option.value
                            return (
                              <button
                                key={option.value}
                                type="button"
                                onClick={() => onSaveCandidateLabel(item.lesson, candidate, candidateIndex, option.value, item.run_id)}
                                disabled={!(item.run_id || recommendationComparison.run_id) || savingLabelKey === labelKey}
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
                </div>
              </div>
            ))}
          </div>
        </div>
      ) : ['lessons_generated', 'recommendation_tested', 'labelled'].includes(selectedScenario.status) ? (
        <div style={S.info}>아직 실행한 리슨별 추천 결과가 없습니다. 3단계의 리슨 카드에서 `추천 콘텐츠 찾기`를 눌러주세요.</div>
      ) : null}

      <div style={{ ...S.info, marginTop:'14px' }}>
        다음 단계에서는 이 비교표에 EmbeddingGemma primary provider와 LightGBM Ranker preview 순위를 채웁니다.
      </div>
    </div>
  )
}
