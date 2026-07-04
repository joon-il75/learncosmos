'use client'

import { useEffect, useMemo, useState } from 'react'
import { useRouter } from 'next/navigation'
import type { PlanetRecordCard, PlanetRecordFeedResponse, RecordFeedFilterKey } from './recordTypes'

type PlanetRouteKind = 'learning' | 'shared'

export function PlanetRecordFeedPage({
  planetId,
  routeKind,
}: {
  planetId: string
  routeKind: PlanetRouteKind
}) {
  const router = useRouter()
  const [feed, setFeed] = useState<PlanetRecordFeedResponse | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [activeFilter, setActiveFilter] = useState<RecordFeedFilterKey>('all')
  const [expandedCardID, setExpandedCardID] = useState<string | null>(null)
  const [copiedCardID, setCopiedCardID] = useState<string | null>(null)

  useEffect(() => {
    let ignore = false

    async function loadRecords() {
      setIsLoading(true)
      setError(null)
      try {
        const response = await fetch(`/api/v1/planets/${routeKind}/${planetId}/records/feed`, {
          credentials: 'include',
          cache: 'no-store',
        })
        const payload = await response.json().catch(() => null)
        if (!response.ok) throw new Error(payload?.error ?? '탐험기록을 불러오지 못했습니다.')
        if (!ignore) {
          setFeed(payload as PlanetRecordFeedResponse)
          setActiveFilter('all')
          setExpandedCardID(null)
        }
      } catch (loadError) {
        if (!ignore) setError(loadError instanceof Error ? loadError.message : '탐험기록을 불러오지 못했습니다.')
      } finally {
        if (!ignore) setIsLoading(false)
      }
    }

    void loadRecords()

    return () => {
      ignore = true
    }
  }, [planetId, routeKind])

  const filteredCards = useMemo(() => {
    const cards = feed?.cards ?? []
    if (activeFilter === 'all') return cards
    if (activeFilter === 'share_candidate') return cards.filter((card) => card.visibility === 'share_candidate')
    return cards.filter((card) => card.category === activeFilter)
  }, [activeFilter, feed?.cards])

  const handleCopy = async (card: PlanetRecordCard) => {
    try {
      await navigator.clipboard.writeText(card.share_text)
      setCopiedCardID(card.id)
      window.setTimeout(() => setCopiedCardID((current) => (current === card.id ? null : current)), 1600)
    } catch {
      setCopiedCardID(null)
    }
  }

  const handleOpenLearning = (card: PlanetRecordCard) => {
    if (!card.point_id || !card.actions.can_open_learning_page) return
    router.push(`/dashboard/planets/learning/${planetId}/points/${card.point_id}`)
  }

  if (isLoading) return <RecordStatus>탐험기록을 불러오는 중...</RecordStatus>
  if (error) return <RecordStatus>{error}</RecordStatus>
  if (!feed) return <RecordStatus>표시할 탐험기록이 없습니다.</RecordStatus>

  return (
    <div className="recordFeedPage">
      <style jsx>{`
        .recordFeedPage {
          box-sizing: border-box;
          display: grid;
          grid-template-rows: auto auto minmax(0, 1fr);
          gap: 14px;
          width: 100%;
          max-width: 100%;
          height: 100%;
          min-height: 0;
          min-width: 0;
          padding: 22px;
          color: #28190a;
          overflow: hidden;
        }
        .recordHero {
          box-sizing: border-box;
          display: grid;
          gap: 12px;
          width: 100%;
          max-width: 100%;
          min-width: 0;
          padding: 18px;
          border-radius: 8px;
          background: linear-gradient(180deg, rgba(255, 250, 239, 0.94), rgba(246, 232, 205, 0.88));
          border: 1px solid rgba(104, 72, 35, 0.16);
          box-shadow: 0 12px 28px rgba(55, 35, 12, 0.08);
        }
        .eyebrow {
          font-size: 11px;
          font-weight: 900;
          color: #805817;
          letter-spacing: 0;
          text-transform: uppercase;
        }
        .title {
          margin: 0;
          font-size: 24px;
          line-height: 1.25;
          color: #211407;
        }
        .summaryGrid {
          display: grid;
          grid-template-columns: repeat(4, minmax(0, 1fr));
          gap: 8px;
          min-width: 0;
        }
        .summaryItem {
          display: grid;
          gap: 3px;
          min-width: 0;
          padding: 10px;
          border-radius: 8px;
          background: rgba(126, 86, 31, 0.10);
        }
        .summaryLabel {
          font-size: 11px;
          font-weight: 800;
          color: rgba(92, 61, 18, 0.70);
        }
        .summaryValue {
          font-size: 15px;
          font-weight: 900;
          color: #3a260c;
        }
        .filterRow {
          display: flex;
          gap: 8px;
          max-width: 100%;
          min-width: 0;
          overflow-x: auto;
          padding: 2px 0 6px;
          scrollbar-width: thin;
        }
        .filterButton {
          flex: 0 0 auto;
          min-height: 36px;
          padding: 0 12px;
          border: 1px solid rgba(104, 72, 35, 0.16);
          border-radius: 999px;
          background: rgba(255, 250, 239, 0.76);
          color: rgba(47, 31, 13, 0.74);
          font-size: 13px;
          font-weight: 850;
          cursor: pointer;
          white-space: nowrap;
        }
        .filterButtonActive {
          border-color: rgba(42, 98, 182, 0.40);
          background: rgba(219, 234, 254, 0.86);
          color: #1f4e8a;
        }
        .feedScroll {
          min-height: 0;
          min-width: 0;
          overflow-y: auto;
          display: grid;
          gap: 12px;
          padding-right: 6px;
        }
        .recordCard {
          box-sizing: border-box;
          display: grid;
          gap: 12px;
          width: 100%;
          max-width: 100%;
          min-width: 0;
          padding: 16px;
          border-radius: 8px;
          border: 1px solid rgba(104, 72, 35, 0.16);
          background: rgba(255, 251, 242, 0.92);
          box-shadow: 0 8px 18px rgba(55, 35, 12, 0.06);
        }
        .cardHeader {
          display: grid;
          grid-template-columns: minmax(0, 1fr) auto;
          gap: 12px;
          align-items: start;
        }
        .cardTitle {
          margin: 0;
          font-size: 17px;
          line-height: 1.35;
          color: #26190b;
          overflow-wrap: anywhere;
        }
        .cardMeta {
          display: flex;
          flex-wrap: wrap;
          gap: 6px;
          color: rgba(83, 55, 18, 0.68);
          font-size: 12px;
          font-weight: 800;
          min-width: 0;
        }
        .badge {
          padding: 5px 9px;
          border-radius: 999px;
          background: rgba(139, 92, 36, 0.11);
          color: #704816;
          font-size: 11px;
          font-weight: 900;
          white-space: nowrap;
        }
        .summary {
          margin: 0;
          color: rgba(38, 25, 11, 0.78);
          font-size: 14px;
          line-height: 1.6;
          overflow-wrap: anywhere;
        }
        .detail {
          display: grid;
          gap: 10px;
          padding: 12px;
          border-radius: 8px;
          background: rgba(255, 255, 255, 0.62);
          border: 1px solid rgba(122, 90, 36, 0.10);
        }
        .detailText {
          margin: 0;
          color: rgba(38, 25, 11, 0.82);
          font-size: 13px;
          line-height: 1.55;
          white-space: pre-wrap;
          overflow-wrap: anywhere;
        }
        .progressTrack {
          height: 9px;
          border-radius: 999px;
          background: rgba(126, 86, 31, 0.14);
          overflow: hidden;
        }
        .progressFill {
          height: 100%;
          border-radius: inherit;
          background: linear-gradient(90deg, #2f7dd1, #f2b84b);
        }
        .actions {
          display: flex;
          flex-wrap: wrap;
          gap: 8px;
        }
        .actionButton {
          min-height: 36px;
          padding: 0 12px;
          border: 1px solid rgba(104, 72, 35, 0.18);
          border-radius: 8px;
          background: rgba(255, 250, 239, 0.88);
          color: #4a3210;
          font-size: 13px;
          font-weight: 850;
          cursor: pointer;
        }
        .primaryAction {
          border-color: rgba(47, 125, 209, 0.34);
          background: rgba(219, 234, 254, 0.88);
          color: #1f4e8a;
        }
        .empty {
          box-sizing: border-box;
          display: grid;
          place-items: center;
          width: 100%;
          max-width: 100%;
          min-width: 0;
          min-height: 220px;
          padding: 24px;
          border-radius: 8px;
          background: rgba(255, 248, 235, 0.86);
          border: 1px solid rgba(104, 72, 35, 0.14);
          color: rgba(43, 29, 12, 0.72);
          font-size: 14px;
          text-align: center;
        }
        @media (max-width: 768px) {
          .recordFeedPage {
            padding: 10px 8px 12px;
            gap: 10px;
            overflow-x: hidden;
          }
          .recordHero {
            gap: 10px;
            padding: 12px;
          }
          .title {
            font-size: 18px;
          }
          .summaryGrid {
            grid-template-columns: repeat(2, minmax(0, 1fr));
          }
          .feedScroll {
            padding-right: 0;
          }
          .cardHeader {
            grid-template-columns: minmax(0, 1fr);
          }
          .badge {
            justify-self: start;
          }
          .actions {
            display: grid;
            grid-template-columns: repeat(2, minmax(0, 1fr));
          }
          .actionButton {
            width: 100%;
            padding: 0 8px;
          }
        }
        @media (max-width: 520px) {
          .recordFeedPage {
            grid-template-rows: auto auto minmax(120px, 1fr);
          }
          .summaryGrid {
            grid-template-columns: 1fr;
          }
          .summaryItem {
            padding: 8px;
          }
          .summaryLabel {
            font-size: 10px;
          }
          .summaryValue {
            font-size: 14px;
          }
          .filterButton {
            min-height: 32px;
            padding: 0 10px;
            font-size: 12px;
          }
          .recordCard {
            padding: 12px;
          }
          .empty {
            min-height: 150px;
            padding: 16px;
          }
        }
        @media (max-width: 420px) {
          .actions {
            grid-template-columns: 1fr;
          }
        }
      `}</style>

      <header className="recordHero">
        <span className="eyebrow">{feed.route_kind === 'shared' ? '완료 탐험기록 보기' : '학습 탐험기록 보기'}</span>
        <h2 className="title">{feed.course.title} 탐험기록</h2>
        <div className="summaryGrid">
          <SummaryItem label="완료 학습콘텐츠" value={`${feed.course.completed_points}/${feed.course.total_points}`} />
          <SummaryItem label="완료 리슨" value={`${feed.course.completed_lessons}/${feed.course.total_lessons}`} />
          <SummaryItem label="기록 카드" value={`${feed.course.record_card_count}개`} />
          <SummaryItem label="공유 추천" value={`${feed.course.share_candidate_count}개`} />
        </div>
      </header>

      <div className="filterRow" aria-label="탐험기록 필터">
        {feed.filters.map((filter) => (
          <button
            key={filter.key}
            type="button"
            className={`filterButton ${activeFilter === filter.key ? 'filterButtonActive' : ''}`}
            onClick={() => setActiveFilter(filter.key)}
          >
            {filter.label} {filter.count}
          </button>
        ))}
      </div>

      <div className="feedScroll">
        {filteredCards.length > 0 ? (
          filteredCards.map((card) => {
            const expanded = expandedCardID === card.id
            return (
              <article key={card.id} className="recordCard">
                <div className="cardHeader">
                  <div style={{ display: 'grid', gap: 6, minWidth: 0 }}>
                    <h3 className="cardTitle" title={card.title}>{card.title}</h3>
                    <div className="cardMeta">
                      {card.lesson_title ? <span title={card.lesson_title}>{card.lesson_title}</span> : null}
                      {card.point_title ? <span title={card.point_title}>{card.point_title}</span> : null}
                      <span>{formatRecordDate(card.occurred_at)}</span>
                    </div>
                  </div>
                  <span className="badge">{card.badge}</span>
                </div>
                <p className="summary">{card.summary}</p>

                {expanded ? <RecordCardDetail card={card} /> : null}

                <div className="actions">
                  <button
                    type="button"
                    className="actionButton"
                    onClick={() => setExpandedCardID(expanded ? null : card.id)}
                  >
                    {expanded ? '접기' : '자세히'}
                  </button>
                  {card.actions.can_copy_share_text ? (
                    <button type="button" className="actionButton" onClick={() => void handleCopy(card)}>
                      {copiedCardID === card.id ? '복사됨' : '문구 복사'}
                    </button>
                  ) : null}
                  {card.actions.can_open_learning_page && card.point_id ? (
                    <button type="button" className="actionButton primaryAction" onClick={() => handleOpenLearning(card)}>
                      학습페이지 보기
                    </button>
                  ) : null}
                </div>
              </article>
            )
          })
        ) : (
          <div className="empty">이 필터에 표시할 기록 카드가 없습니다.</div>
        )}
      </div>
    </div>
  )
}

function SummaryItem({ label, value }: { label: string; value: string }) {
  return (
    <div className="summaryItem">
      <span className="summaryLabel">{label}</span>
      <strong className="summaryValue">{value}</strong>
    </div>
  )
}

function RecordCardDetail({ card }: { card: PlanetRecordCard }) {
  const detailItems = [
    card.detail.primary_text,
    card.detail.journal_excerpt,
    card.detail.question_excerpt,
    card.detail.artifact_title ? `결과물: ${card.detail.artifact_title}` : '',
    card.detail.self_evaluation_summary,
    card.detail.practice_summary,
  ].filter((value): value is string => Boolean(value?.trim()))

  return (
    <div className="detail">
      {card.detail.lesson_progress ? (
        <div style={{ display: 'grid', gap: 8 }}>
          <p className="detailText">
            리슨 진행률 {card.detail.lesson_progress.completed_points}/{card.detail.lesson_progress.total_points}개 완료
          </p>
          <div className="progressTrack" aria-hidden="true">
            <div className="progressFill" style={{ width: `${card.detail.lesson_progress.percent}%` }} />
          </div>
        </div>
      ) : null}
      {detailItems.length > 0 ? (
        detailItems.map((item, index) => <p key={`${card.id}-detail-${index}`} className="detailText">{item}</p>)
      ) : (
        <p className="detailText">이 기록의 세부 내용은 학습페이지에서 확인할 수 있습니다.</p>
      )}
    </div>
  )
}

function RecordStatus({ children }: { children: string }) {
  return (
    <div style={{
      display: 'grid',
      placeItems: 'center',
      minHeight: 260,
      padding: 24,
      borderRadius: 8,
      background: 'rgba(255, 248, 235, 0.86)',
      border: '1px solid rgba(104, 72, 35, 0.14)',
      color: 'rgba(43, 29, 12, 0.72)',
      fontSize: 14,
      textAlign: 'center',
    }}>
      {children}
    </div>
  )
}

function formatRecordDate(value: string): string {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleDateString('ko-KR', {
    month: '2-digit',
    day: '2-digit',
  })
}
