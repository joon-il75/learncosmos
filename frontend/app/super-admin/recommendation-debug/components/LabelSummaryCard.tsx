import { LABEL_OPTIONS } from '../constants'
import { S } from '../styles'
import type { LabelSummary } from '../types'

type Props = {
  labelSummary: LabelSummary | null
  exportingLabelFormat: 'jsonl' | 'csv' | null
  onExportLabelDataset: (format: 'jsonl' | 'csv') => void
  onRefreshLabelSummary: () => void
}

export default function LabelSummaryCard({
  labelSummary,
  exportingLabelFormat,
  onExportLabelDataset,
  onRefreshLabelSummary,
}: Props) {
  return (
    <div style={S.card}>
      <div style={{ display:'flex', justifyContent:'space-between', gap:'12px', alignItems:'flex-start', flexWrap:'wrap', marginBottom:'14px' }}>
        <div>
          <div style={{ fontSize:'15px', fontWeight:800 }}>5. Ranker 학습 라벨</div>
          <div style={{ fontSize:'12px', color:'rgba(200,210,235,0.58)', marginTop:'5px', lineHeight:1.7 }}>
            추천 후보 라벨을 기준으로 LightGBM 학습 데이터 품질을 추적합니다.
          </div>
        </div>
        <div style={{ display:'flex', gap:'8px', flexWrap:'wrap', justifyContent:'flex-end' }}>
          <button style={S.secondary} onClick={() => onExportLabelDataset('jsonl')} disabled={exportingLabelFormat !== null || !labelSummary?.total_labels}>
            {exportingLabelFormat === 'jsonl' ? 'JSONL 생성 중...' : 'JSONL export'}
          </button>
          <button style={S.secondary} onClick={() => onExportLabelDataset('csv')} disabled={exportingLabelFormat !== null || !labelSummary?.total_labels}>
            {exportingLabelFormat === 'csv' ? 'CSV 생성 중...' : 'CSV export'}
          </button>
          <button style={S.secondary} onClick={onRefreshLabelSummary}>
            라벨 요약 새로고침
          </button>
        </div>
      </div>

      {labelSummary ? (
        <>
          <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(150px, 1fr))', gap:'12px' }}>
            <div>
              <div style={S.label}>quality status</div>
              <div style={{ fontSize:'16px', fontWeight:800, color: labelSummary.quality_status === 'passed' ? '#8EE59A' : labelSummary.quality_status === 'failed' ? '#FFB4A2' : '#FAC775' }}>
                {labelSummary.quality_status}
              </div>
            </div>
            <div>
              <div style={S.label}>labels</div>
              <div style={{ fontSize:'16px', fontWeight:800 }}>{labelSummary.total_labels} / {labelSummary.minimum_required}</div>
            </div>
            <div>
              <div style={S.label}>good fit</div>
              <div style={{ fontSize:'16px', fontWeight:800 }}>{(labelSummary.good_fit_rate * 100).toFixed(0)}%</div>
            </div>
            <div>
              <div style={S.label}>irrelevant</div>
              <div style={{ fontSize:'16px', fontWeight:800 }}>{(labelSummary.irrelevant_rate * 100).toFixed(0)}%</div>
            </div>
            <div>
              <div style={S.label}>duplicate</div>
              <div style={{ fontSize:'16px', fontWeight:800 }}>{(labelSummary.duplicate_rate * 100).toFixed(0)}%</div>
            </div>
            <div>
              <div style={S.label}>problem</div>
              <div style={{ fontSize:'16px', fontWeight:800 }}>{(labelSummary.problem_rate * 100).toFixed(0)}%</div>
            </div>
          </div>
          <div style={{ ...S.info, marginTop:'14px' }}>
            {labelSummary.quality_reason}
            {labelSummary.latest_run_id && (
              <div style={{ marginTop:'4px', wordBreak:'break-all' }}>latest run: {labelSummary.latest_run_id}</div>
            )}
          </div>
          <div style={{ marginTop:'14px', display:'flex', gap:'8px', flexWrap:'wrap' }}>
            {LABEL_OPTIONS.map(option => (
              <div key={option.value} style={{ border:'1px solid rgba(120,140,200,0.14)', background:'rgba(255,255,255,0.04)', borderRadius:'999px', padding:'7px 10px', fontSize:'11px', color:'rgba(200,210,235,0.78)' }}>
                {option.label} {labelSummary.label_counts[option.value] ?? 0}
              </div>
            ))}
          </div>
        </>
      ) : (
        <div style={S.warn}>라벨 요약이 아직 없습니다. 추천 비교 실행 후 후보에 라벨을 저장하면 품질 신호가 표시됩니다.</div>
      )}
    </div>
  )
}
