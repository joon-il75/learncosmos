import { S } from '../styles'
import type { LLMJobObservationCount, LLMJobObservations, LLMJobObservationStatus } from '../types'

type Props = {
  observations: LLMJobObservations | null
  loading: boolean
  onRefresh: () => void
}

const statusOrder: LLMJobObservationStatus[] = ['queued', 'running', 'failed', 'succeeded', 'canceled', 'expired']

const statusColor = (status: LLMJobObservationStatus, staleCount: number) => {
  if (staleCount > 0) return '#FFB4A2'
  if (status === 'failed' || status === 'expired') return '#FAC775'
  if (status === 'queued' || status === 'running') return '#8CC7FF'
  if (status === 'succeeded') return '#8EE59A'
  return 'rgba(200,210,235,0.72)'
}

const formatTime = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleString('ko-KR')
}

const formatDuration = (ms?: number) => {
  if (!ms || ms <= 0) return '비활성'
  const minutes = Math.round(ms / 60000)
  if (minutes < 60) return `${minutes}분`
  const hours = Math.floor(minutes / 60)
  const rest = minutes % 60
  return rest ? `${hours}시간 ${rest}분` : `${hours}시간`
}

const summarize = (counts: LLMJobObservationCount[]) => counts.reduce(
  (acc, item) => {
    acc.total += item.total_count
    acc.stale += item.stale_count
    if (item.status === 'queued') acc.queued += item.total_count
    if (item.status === 'running') acc.running += item.total_count
    if (item.status === 'failed') acc.failed += item.total_count
    return acc
  },
  { total:0, stale:0, queued:0, running:0, failed:0 },
)

const groupedCounts = (counts: LLMJobObservationCount[]) => {
  const grouped = new Map<string, LLMJobObservationCount[]>()
  for (const item of counts) {
    grouped.set(item.feature, [...(grouped.get(item.feature) ?? []), item])
  }
  return [...grouped.entries()].map(([feature, items]) => ({
    feature,
    items: items.sort((a, b) => statusOrder.indexOf(a.status) - statusOrder.indexOf(b.status)),
  }))
}

export default function WorkerObservationCard({ observations, loading, onRefresh }: Props) {
  const summary = summarize(observations?.counts ?? [])

  return (
    <div style={S.card}>
      <div style={{ display:'flex', justifyContent:'space-between', gap:'12px', alignItems:'flex-start', flexWrap:'wrap' }}>
        <div>
          <div style={{ fontSize:'15px', fontWeight:700, marginBottom:'6px' }}>LLM Worker 관측</div>
          <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.58)', lineHeight:1.7 }}>
            공통 worker job의 상태, stale 후보, feature별 적체를 확인합니다.
          </div>
        </div>
        <button style={S.secondary} onClick={onRefresh} disabled={loading}>
          {loading ? '조회 중...' : 'worker 새로고침'}
        </button>
      </div>

      {observations ? (
        <div style={{ display:'grid', gap:'14px', marginTop:'16px' }}>
          <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(150px, 1fr))', gap:'12px' }}>
            <Metric label="total" value={summary.total} color="#E8EAF2" />
            <Metric label="queued" value={summary.queued} color={summary.queued > 0 ? '#8CC7FF' : 'rgba(200,210,235,0.72)'} />
            <Metric label="running" value={summary.running} color={summary.running > 0 ? '#8CC7FF' : 'rgba(200,210,235,0.72)'} />
            <Metric label="failed" value={summary.failed} color={summary.failed > 0 ? '#FAC775' : 'rgba(200,210,235,0.72)'} />
            <Metric label="stale" value={summary.stale} color={summary.stale > 0 ? '#FFB4A2' : '#8EE59A'} />
          </div>

          <div style={{ ...S.info, display:'flex', justifyContent:'space-between', gap:'10px', flexWrap:'wrap' }}>
            <span>queued stale {formatDuration(observations.policy.queued_timeout_ms)}</span>
            <span>running stale {formatDuration(observations.policy.running_timeout_ms)}</span>
            <span>updated {formatTime(observations.generated_at)}</span>
          </div>

          {observations.counts.length > 0 ? (
            <div style={{ display:'grid', gap:'10px' }}>
              {groupedCounts(observations.counts).map(group => (
                <div key={group.feature} style={{ background:'rgba(255,255,255,0.03)', border:'1px solid rgba(120,140,200,0.12)', borderRadius:'12px', padding:'12px' }}>
                  <div style={{ fontSize:'13px', fontWeight:800, marginBottom:'10px', wordBreak:'break-all' }}>{group.feature}</div>
                  <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(150px, 1fr))', gap:'8px' }}>
                    {group.items.map(item => (
                      <div key={`${item.feature}-${item.status}`} style={{ border:'1px solid rgba(120,140,200,0.12)', borderRadius:'10px', padding:'10px', background:'rgba(0,0,0,0.12)' }}>
                        <div style={{ ...S.label, marginBottom:'4px' }}>{item.status}</div>
                        <div style={{ fontSize:'18px', fontWeight:900, color:statusColor(item.status, item.stale_count) }}>
                          {item.total_count}
                          {item.stale_count > 0 && <span style={{ fontSize:'12px', marginLeft:'6px' }}>stale {item.stale_count}</span>}
                        </div>
                        <div style={{ fontSize:'11px', color:'rgba(200,210,235,0.52)', marginTop:'6px', lineHeight:1.5 }}>
                          <div>created {formatTime(item.oldest_created_at)}</div>
                          <div>activity {formatTime(item.oldest_activity_at)}</div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div style={S.warn}>집계된 worker job이 없습니다.</div>
          )}
        </div>
      ) : (
        <div style={{ ...S.warn, marginTop:'16px' }}>worker 관측 데이터를 아직 불러오지 못했습니다.</div>
      )}
    </div>
  )
}

function Metric({ label, value, color }: { label: string; value: number; color: string }) {
  return (
    <div style={{ background:'rgba(255,255,255,0.03)', border:'1px solid rgba(120,140,200,0.12)', borderRadius:'12px', padding:'12px' }}>
      <div style={S.label}>{label}</div>
      <div style={{ fontSize:'22px', fontWeight:900, color }}>{value}</div>
    </div>
  )
}
