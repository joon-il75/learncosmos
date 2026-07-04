
import { S } from '../styles'
import type { RecommendationSpecMetrics } from '../types'

type Props = {
  metrics: RecommendationSpecMetrics | null
  loading: boolean
  activeDays: number
  onRefresh: (days?: number) => void
}

const sourceLabel: Record<string, string> = {
  pattern_template: 'pattern template',
  fallback: 'fallback',
  empty: 'empty',
}

const resolutionLabel: Record<string, string> = {
  direct_lesson: 'direct lesson',
  explorer_region: 'explorer region',
  explorer_subregion: 'explorer subregion',
  empty: 'empty',
}

const formatTime = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleString('ko-KR')
}

const formatPercent = (value?: number | null) => {
  if (typeof value !== 'number' || Number.isNaN(value)) return '-'
  return `${Math.round(value * 1000) / 10}%`
}

const formatNumber = (value?: number | null) => {
  if (typeof value !== 'number' || Number.isNaN(value)) return '-'
  return Number.isInteger(value) ? String(value) : value.toFixed(1)
}

const countEntries = (counts?: Record<string, number>) => Object.entries(counts ?? {})
  .sort((a, b) => b[1] - a[1] || a[0].localeCompare(b[0]))

export default function SpecMetricsCard({ metrics, loading, activeDays, onRefresh }: Props) {
  return (
    <div style={S.card}>
      <div style={{ display:'flex', justifyContent:'space-between', gap:'12px', alignItems:'flex-start', flexWrap:'wrap' }}>
        <div>
          <div style={{ fontSize:'15px', fontWeight:700, marginBottom:'6px' }}>리슨 명세 관측</div>
          <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.58)', lineHeight:1.7 }}>
            추천 검색 명세 source, lesson resolution, effective query, learner event를 확인합니다.
          </div>
        </div>
        <div style={{ display:'flex', gap:'8px', flexWrap:'wrap' }}>
          {[7, 30].map(days => (
            <button
              key={days}
              style={days === activeDays ? S.primary : S.secondary}
              onClick={() => onRefresh(days)}
              disabled={loading}
            >
              {days}일
            </button>
          ))}
          <button style={S.secondary} onClick={() => onRefresh(activeDays)} disabled={loading}>
            {loading ? '조회 중...' : '새로고침'}
          </button>
        </div>
      </div>

      {metrics ? (
        <div style={{ display:'grid', gap:'14px', marginTop:'16px' }}>
          <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(150px, 1fr))', gap:'12px' }}>
            <Metric label="searches" value={metrics.recommendation_searches} color="#E8EAF2" />
            <Metric label="empty query" value={formatPercent(metrics.empty_effective_query_rate)} color={metrics.empty_effective_query_rate ? '#FAC775' : '#8EE59A'} />
            <Metric label="short query" value={formatPercent(metrics.short_effective_query_rate)} color={metrics.short_effective_query_rate ? '#FAC775' : '#8EE59A'} />
            <Metric label="exposed" value={metrics.rollout_events.recommendation_exposed} color="#8CC7FF" />
            <Metric label="avg candidates" value={formatNumber(metrics.rollout_events.average_candidate_count)} color="#E8EAF2" />
          </div>

          <div style={{ ...S.info, display:'flex', justifyContent:'space-between', gap:'10px', flexWrap:'wrap' }}>
            <span>{formatTime(metrics.window_started_at)} - {formatTime(metrics.window_ended_at)}</span>
            <span>교체 {metrics.rollout_events.material_replaced}</span>
            <span>broken {metrics.rollout_events.broken_link_reported}</span>
            <span>wrong {metrics.rollout_events.wrong_content_reported}</span>
          </div>

          <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(240px, 1fr))', gap:'12px' }}>
            <CountPanel title="search spec source" entries={countEntries(metrics.search_spec_source_counts)} labels={sourceLabel} />
            <CountPanel title="resolution source" entries={countEntries(metrics.resolution_source_counts)} labels={resolutionLabel} />
          </div>

          {metrics.top_effective_queries.length > 0 ? (
            <div style={{ background:'rgba(255,255,255,0.03)', border:'1px solid rgba(120,140,200,0.12)', borderRadius:'12px', padding:'12px' }}>
              <div style={{ fontSize:'13px', fontWeight:800, marginBottom:'10px' }}>top effective query</div>
              <div style={{ display:'grid', gap:'8px' }}>
                {metrics.top_effective_queries.map(item => (
                  <div key={item.query} style={{ display:'grid', gridTemplateColumns:'minmax(0, 1fr) auto', gap:'12px', alignItems:'center', border:'1px solid rgba(120,140,200,0.1)', borderRadius:'10px', padding:'10px', background:'rgba(0,0,0,0.12)' }}>
                    <div style={{ fontSize:'12px', color:'rgba(232,234,242,0.84)', lineHeight:1.5, wordBreak:'keep-all', overflowWrap:'anywhere' }}>{item.query}</div>
                    <div style={{ fontSize:'16px', fontWeight:900, color:'#8CC7FF' }}>{item.count}</div>
                  </div>
                ))}
              </div>
            </div>
          ) : (
            <div style={S.warn}>집계된 effective query가 없습니다.</div>
          )}
        </div>
      ) : (
        <div style={{ ...S.warn, marginTop:'16px' }}>리슨 명세 관측 데이터를 아직 불러오지 못했습니다.</div>
      )}
    </div>
  )
}

function Metric({ label, value, color }: { label: string; value: number | string; color: string }) {
  return (
    <div style={{ background:'rgba(255,255,255,0.03)', border:'1px solid rgba(120,140,200,0.12)', borderRadius:'12px', padding:'12px' }}>
      <div style={S.label}>{label}</div>
      <div style={{ fontSize:'22px', fontWeight:900, color }}>{value}</div>
    </div>
  )
}

function CountPanel({ title, entries, labels }: { title: string; entries: [string, number][]; labels: Record<string, string> }) {
  return (
    <div style={{ background:'rgba(255,255,255,0.03)', border:'1px solid rgba(120,140,200,0.12)', borderRadius:'12px', padding:'12px' }}>
      <div style={{ fontSize:'13px', fontWeight:800, marginBottom:'10px' }}>{title}</div>
      {entries.length > 0 ? (
        <div style={{ display:'grid', gap:'8px' }}>
          {entries.map(([key, count]) => (
            <div key={key} style={{ display:'flex', justifyContent:'space-between', gap:'10px', alignItems:'center', borderBottom:'1px solid rgba(120,140,200,0.08)', paddingBottom:'8px' }}>
              <span style={{ fontSize:'12px', color:'rgba(200,210,235,0.72)' }}>{labels[key] ?? key}</span>
              <strong style={{ fontSize:'15px', color:'#E8EAF2' }}>{count}</strong>
            </div>
          ))}
        </div>
      ) : (
        <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.52)' }}>집계 없음</div>
      )}
    </div>
  )
}
