import { S } from '../styles'
import type { LessonSearchPrerunReport, LessonSearchPrerunReportsResponse, LessonSearchPrerunTokenCount } from '../types'

type Props = {
  reports: LessonSearchPrerunReportsResponse | null
  loading: boolean
  onRefresh: () => void
}

const statusLabel: Record<string, string> = {
  succeeded: '정상',
  failed: '실패',
  running: '실행 중',
  skipped_missing_credentials: 'credential 없음',
  skipped_lock_not_acquired: 'lock skip',
}

const statusColor = (status?: string) => {
  switch (status) {
    case 'succeeded': return '#8EE59A'
    case 'failed': return '#FFB4A2'
    case 'running': return '#8CC7FF'
    default: return '#FAC775'
  }
}

const formatTime = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleString('ko-KR')
}

const formatNumber = (value?: number | null) => {
  if (typeof value !== 'number' || Number.isNaN(value)) return '-'
  return String(value)
}

const tokenEntries = (values?: LessonSearchPrerunTokenCount[]) => [...(values ?? [])]
  .filter(item => item.token && item.count > 0)
  .sort((a, b) => b.count - a.count || a.token.localeCompare(b.token))
  .slice(0, 6)

export default function LessonSearchPrerunCard({ reports, loading, onRefresh }: Props) {
  const latest = reports?.reports?.[0] ?? null
  const history = reports?.reports?.slice(1, 5) ?? []
  const unregisteredHints = tokenEntries(latest?.noise_taxonomy_summary?.unregistered_noise_hints)
  const registeredAvoids = tokenEntries(latest?.noise_taxonomy_summary?.registered_avoid_hits)

  return (
    <div style={S.card}>
      <div style={{ display:'flex', justifyContent:'space-between', gap:'12px', alignItems:'flex-start', flexWrap:'wrap' }}>
        <div>
          <div style={{ fontSize:'15px', fontWeight:700, marginBottom:'6px' }}>리슨 검색 프리런</div>
          <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.58)', lineHeight:1.7 }}>
            YouTube/Naver provider fixture 반복 점검 결과와 quota, noise 후보를 확인합니다.
          </div>
        </div>
        <button style={S.secondary} onClick={onRefresh} disabled={loading}>
          {loading ? '조회 중...' : '새로고침'}
        </button>
      </div>

      {latest ? (
        <div style={{ display:'grid', gap:'14px', marginTop:'16px' }}>
          <div style={{ ...S.info, display:'flex', justifyContent:'space-between', gap:'10px', flexWrap:'wrap' }}>
            <span>최근 실행 {formatTime(latest.started_at)}</span>
            <span>완료 {formatTime(latest.finished_at)}</span>
            <span>provider {latest.provider_filter}</span>
            <span>top {latest.top_n} / risk {latest.risk_top_n}</span>
          </div>

          <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(150px, 1fr))', gap:'12px' }}>
            <Metric label="status" value={statusLabel[latest.status] ?? latest.status} color={statusColor(latest.status)} />
            <Metric label="provider errors" value={formatNumber(latest.summary.provider_errors)} color={latest.summary.provider_errors ? '#FFB4A2' : '#8EE59A'} />
            <Metric label="noise suspects" value={formatNumber(latest.summary.noise_suspects)} color={latest.summary.noise_suspects ? '#FAC775' : '#8EE59A'} />
            <Metric label="requests" value={formatNumber(latest.summary.estimated_provider_requests)} color="#8CC7FF" />
            <Metric label="candidate max" value={formatNumber(latest.summary.estimated_candidates_max)} color="#E8EAF2" />
          </div>

          {(latest.error_code || latest.error_message) && (
            <div style={S.warn}>{latest.error_code || 'error'} · {latest.error_message || '상세 메시지 없음'}</div>
          )}

          <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(240px, 1fr))', gap:'12px' }}>
            <TokenPanel title="unregistered noise hints" entries={unregisteredHints} emptyText="미등록 noise hint 없음" tone={unregisteredHints.length > 0 ? 'warn' : 'ok'} />
            <TokenPanel title="registered avoid hits" entries={registeredAvoids} emptyText="등록 avoid hit 없음" tone="info" />
          </div>

          <RiskFixturePanel latest={latest} />

          {history.length > 0 && (
            <div style={{ background:'rgba(255,255,255,0.03)', border:'1px solid rgba(120,140,200,0.12)', borderRadius:'12px', padding:'12px' }}>
              <div style={{ fontSize:'13px', fontWeight:800, marginBottom:'10px' }}>최근 이력</div>
              <div style={{ display:'grid', gap:'8px' }}>
                {history.map(report => (
                  <div key={report.id} style={{ display:'grid', gridTemplateColumns:'minmax(0, 1fr) auto auto', gap:'10px', alignItems:'center', borderBottom:'1px solid rgba(120,140,200,0.08)', paddingBottom:'8px' }}>
                    <span style={{ fontSize:'12px', color:'rgba(200,210,235,0.72)', overflowWrap:'anywhere' }}>{formatTime(report.started_at)}</span>
                    <strong style={{ fontSize:'12px', color:statusColor(report.status) }}>{statusLabel[report.status] ?? report.status}</strong>
                    <span style={{ fontSize:'12px', color:'rgba(200,210,235,0.62)' }}>noise {formatNumber(report.summary.noise_suspects)}</span>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      ) : (
        <div style={{ ...S.warn, marginTop:'16px' }}>저장된 리슨 검색 프리런 report가 아직 없습니다.</div>
      )}
    </div>
  )
}

function Metric({ label, value, color }: { label: string; value: number | string; color: string }) {
  return (
    <div style={{ background:'rgba(255,255,255,0.03)', border:'1px solid rgba(120,140,200,0.12)', borderRadius:'12px', padding:'12px' }}>
      <div style={S.label}>{label}</div>
      <div style={{ fontSize:'22px', fontWeight:900, color, overflowWrap:'anywhere' }}>{value}</div>
    </div>
  )
}

function TokenPanel({ title, entries, emptyText, tone }: { title: string; entries: LessonSearchPrerunTokenCount[]; emptyText: string; tone: 'ok' | 'warn' | 'info' }) {
  const color = tone === 'warn' ? '#FAC775' : tone === 'ok' ? '#8EE59A' : '#8CC7FF'
  return (
    <div style={{ background:'rgba(255,255,255,0.03)', border:'1px solid rgba(120,140,200,0.12)', borderRadius:'12px', padding:'12px' }}>
      <div style={{ fontSize:'13px', fontWeight:800, marginBottom:'10px' }}>{title}</div>
      {entries.length > 0 ? (
        <div style={{ display:'grid', gap:'8px' }}>
          {entries.map(item => (
            <div key={item.token} style={{ display:'flex', justifyContent:'space-between', gap:'10px', alignItems:'center', borderBottom:'1px solid rgba(120,140,200,0.08)', paddingBottom:'8px' }}>
              <span style={{ fontSize:'12px', color:'rgba(200,210,235,0.76)', overflowWrap:'anywhere' }}>{item.token}</span>
              <strong style={{ fontSize:'15px', color }}>{item.count}</strong>
            </div>
          ))}
        </div>
      ) : (
        <div style={{ fontSize:'12px', color }}>{emptyText}</div>
      )}
    </div>
  )
}

function RiskFixturePanel({ latest }: { latest: LessonSearchPrerunReport }) {
  const processed = latest.summary.risk_fixture_ids_processed ?? []
  const configured = latest.summary.risk_fixture_ids_configured ?? latest.summary.quota?.risk_fixture_ids ?? []
  return (
    <div style={{ background:'rgba(255,255,255,0.03)', border:'1px solid rgba(120,140,200,0.12)', borderRadius:'12px', padding:'12px' }}>
      <div style={{ fontSize:'13px', fontWeight:800, marginBottom:'10px' }}>risk fixture</div>
      <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(220px, 1fr))', gap:'10px' }}>
        <FixtureList title="processed" values={processed} emptyText="처리된 risk fixture 없음" />
        <FixtureList title="configured" values={configured} emptyText="설정된 risk fixture 없음" />
      </div>
    </div>
  )
}

function FixtureList({ title, values, emptyText }: { title: string; values: string[]; emptyText: string }) {
  return (
    <div>
      <div style={S.label}>{title}</div>
      {values.length > 0 ? (
        <div style={{ display:'flex', flexWrap:'wrap', gap:'8px' }}>
          {values.slice(0, 8).map(value => (
            <span key={value} style={{ border:'1px solid rgba(120,140,200,0.18)', borderRadius:'999px', padding:'6px 9px', fontSize:'11px', color:'rgba(232,234,242,0.82)', background:'rgba(0,0,0,0.12)', overflowWrap:'anywhere' }}>{value}</span>
          ))}
        </div>
      ) : (
        <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.52)' }}>{emptyText}</div>
      )}
    </div>
  )
}
