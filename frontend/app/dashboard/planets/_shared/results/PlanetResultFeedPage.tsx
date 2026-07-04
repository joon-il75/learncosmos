'use client'

import { useEffect, useMemo, useState } from 'react'
import { useRouter } from 'next/navigation'
import type { PlanetResultAggregate, PlanetResultArtifact, PlanetResultLesson, PlanetResultPoint } from './resultTypes'

type PlanetRouteKind = 'learning' | 'shared'
type ResultFilterKey = 'all' | 'link' | 'image' | 'document' | 'text' | 'other'

interface ResultCard {
  id: string
  lessonId: string
  lessonTitle: string
  pointId: string
  pointTitle: string
  pointType: PlanetResultPoint['point_type']
  artifact: PlanetResultArtifact
  category: ResultFilterKey
}

const filterLabels: Record<ResultFilterKey, string> = {
  all: '전체',
  link: '링크',
  image: '이미지',
  document: '문서/파일',
  text: '텍스트',
  other: '기타',
}

export function PlanetResultFeedPage({
  planetId,
  routeKind,
}: {
  planetId: string
  routeKind: PlanetRouteKind
}) {
  const router = useRouter()
  const [results, setResults] = useState<PlanetResultAggregate | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [activeFilter, setActiveFilter] = useState<ResultFilterKey>('all')
  const [expandedCardID, setExpandedCardID] = useState<string | null>(null)

  useEffect(() => {
    let ignore = false

    async function loadResults() {
      setIsLoading(true)
      setError(null)
      try {
        const response = await fetch(`/api/v1/planets/${routeKind}/${planetId}/results`, {
          credentials: 'include',
          cache: 'no-store',
        })
        const payload = await response.json().catch(() => null)
        if (!response.ok) throw new Error(payload?.error ?? '탐험결과물을 불러오지 못했습니다.')
        if (!ignore) {
          setResults(payload as PlanetResultAggregate)
          setActiveFilter('all')
          setExpandedCardID(null)
        }
      } catch (loadError) {
        if (!ignore) setError(loadError instanceof Error ? loadError.message : '탐험결과물을 불러오지 못했습니다.')
      } finally {
        if (!ignore) setIsLoading(false)
      }
    }

    void loadResults()

    return () => {
      ignore = true
    }
  }, [planetId, routeKind])

  const cards = useMemo(() => flattenResultCards(results), [results])
  const filteredCards = useMemo(
    () => (activeFilter === 'all' ? cards : cards.filter((card) => card.category === activeFilter)),
    [activeFilter, cards],
  )
  const filterCounts = useMemo(() => buildFilterCounts(cards), [cards])
  const lessonCount = useMemo(() => new Set(cards.map((card) => card.lessonId)).size, [cards])
  const pointCount = useMemo(() => new Set(cards.map((card) => card.pointId)).size, [cards])

  if (isLoading) return <ResultStatus>탐험결과물을 불러오는 중...</ResultStatus>
  if (error) return <ResultStatus>{error}</ResultStatus>
  if (!results) return <ResultStatus>표시할 탐험결과물이 없습니다.</ResultStatus>

  return (
    <div className="resultFeedPage">
      <style jsx>{`
        .resultFeedPage {
          box-sizing: border-box;
          display: grid;
          grid-template-rows: auto auto minmax(0, 1fr);
          gap: 14px;
          width: 100%;
          max-width: 100%;
          height: 100%;
          min-width: 0;
          min-height: 0;
          padding: 22px;
          color: #28190a;
          overflow: hidden;
        }
        .resultHero,
        .resultCard,
        .empty {
          box-sizing: border-box;
          width: 100%;
          max-width: 100%;
          min-width: 0;
        }
        .resultHero {
          display: grid;
          gap: 12px;
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
          color: #211407;
          font-size: 24px;
          line-height: 1.25;
          overflow-wrap: anywhere;
        }
        .notice {
          margin: 0;
          color: rgba(64, 43, 16, 0.72);
          font-size: 13px;
          line-height: 1.5;
          overflow-wrap: anywhere;
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
          color: rgba(92, 61, 18, 0.70);
          font-size: 11px;
          font-weight: 800;
        }
        .summaryValue {
          color: #3a260c;
          font-size: 15px;
          font-weight: 900;
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
          white-space: nowrap;
          cursor: pointer;
        }
        .filterButtonActive {
          border-color: rgba(39, 94, 77, 0.36);
          background: rgba(220, 242, 234, 0.88);
          color: #1e5a48;
        }
        .feedScroll {
          display: grid;
          gap: 12px;
          min-width: 0;
          min-height: 0;
          overflow-y: auto;
          padding-right: 6px;
        }
        .resultCard {
          display: grid;
          gap: 12px;
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
          color: #26190b;
          font-size: 17px;
          line-height: 1.35;
          overflow-wrap: anywhere;
        }
        .cardMeta {
          display: flex;
          flex-wrap: wrap;
          gap: 6px;
          min-width: 0;
          color: rgba(83, 55, 18, 0.68);
          font-size: 12px;
          font-weight: 800;
        }
        .badge {
          padding: 5px 9px;
          border-radius: 999px;
          background: rgba(39, 94, 77, 0.11);
          color: #1e5a48;
          font-size: 11px;
          font-weight: 900;
          white-space: nowrap;
        }
        .description {
          margin: 0;
          color: rgba(38, 25, 11, 0.78);
          font-size: 14px;
          line-height: 1.6;
          white-space: pre-wrap;
          overflow-wrap: anywhere;
        }
        .detail {
          display: grid;
          gap: 8px;
          padding: 12px;
          border-radius: 8px;
          background: rgba(255, 255, 255, 0.62);
          border: 1px solid rgba(122, 90, 36, 0.10);
        }
        .detailRow {
          display: grid;
          gap: 3px;
          min-width: 0;
        }
        .detailLabel {
          color: rgba(78, 54, 19, 0.62);
          font-size: 11px;
          font-weight: 900;
        }
        .detailValue {
          color: rgba(38, 25, 11, 0.82);
          font-size: 13px;
          line-height: 1.55;
          overflow-wrap: anywhere;
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
          text-decoration: none;
          display: inline-flex;
          align-items: center;
          justify-content: center;
        }
        .primaryAction {
          border-color: rgba(39, 94, 77, 0.36);
          background: rgba(220, 242, 234, 0.90);
          color: #1e5a48;
        }
        .disabledNote {
          align-self: center;
          color: rgba(83, 55, 18, 0.62);
          font-size: 12px;
          font-weight: 800;
        }
        .empty {
          display: grid;
          place-items: center;
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
          .resultFeedPage {
            padding: 10px 8px 12px;
            gap: 10px;
            overflow-x: hidden;
          }
          .resultHero {
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
          .resultCard {
            padding: 12px;
          }
          .actions {
            grid-template-columns: 1fr;
          }
          .empty {
            min-height: 150px;
            padding: 16px;
          }
        }
      `}</style>

      <header className="resultHero">
        <span className="eyebrow">{routeKind === 'shared' ? '완료 탐험결과물 보기' : '학습 탐험결과물 보기'}</span>
        <h2 className="title">{results.course.title} 탐험결과물</h2>
        <p className="notice">학습페이지의 탐험결과물 탭에서 제출한 콘텐츠를 확인할 수 있습니다. 커뮤니티 피드백은 정식 오픈 때 연결됩니다.</p>
        <div className="summaryGrid">
          <SummaryItem label="제출 결과물" value={`${cards.length}개`} />
          <SummaryItem label="결과물 지점" value={`${pointCount}/${results.course.total_points}`} />
          <SummaryItem label="연결 리슨" value={`${lessonCount}개`} />
          <SummaryItem label="커뮤니티 피드백" value="준비 중" />
        </div>
      </header>

      <div className="filterRow" aria-label="탐험결과물 필터">
        {(Object.keys(filterLabels) as ResultFilterKey[]).map((filterKey) => (
          <button
            key={filterKey}
            type="button"
            className={`filterButton ${activeFilter === filterKey ? 'filterButtonActive' : ''}`}
            onClick={() => setActiveFilter(filterKey)}
          >
            {filterLabels[filterKey]} {filterCounts[filterKey] ?? 0}
          </button>
        ))}
      </div>

      <div className="feedScroll">
        {filteredCards.length > 0 ? (
          filteredCards.map((card) => {
            const expanded = expandedCardID === card.id
            const title = card.artifact.title.trim() || '제목 없는 결과물'
            const description = card.artifact.description.trim() || '설명이 아직 없습니다.'
            const url = card.artifact.url.trim()
            return (
              <article key={card.id} className="resultCard">
                <div className="cardHeader">
                  <div style={{ display: 'grid', gap: 6, minWidth: 0 }}>
                    <h3 className="cardTitle" title={title}>{title}</h3>
                    <div className="cardMeta">
                      <span title={card.lessonTitle}>{card.lessonTitle}</span>
                      <span title={card.pointTitle}>{card.pointTitle}</span>
                      <span>{formatResultDate(card.artifact.updated_at || card.artifact.created_at)}</span>
                    </div>
                  </div>
                  <span className="badge">{formatArtifactType(card.artifact.artifact_type)}</span>
                </div>
                <p className="description">{description}</p>

                {expanded ? <ResultCardDetail card={card} /> : null}

                <div className="actions">
                  <button
                    type="button"
                    className="actionButton"
                    onClick={() => setExpandedCardID(expanded ? null : card.id)}
                  >
                    {expanded ? '접기' : '자세히'}
                  </button>
                  {url ? (
                    <a href={url} target="_blank" rel="noreferrer" className="actionButton primaryAction">
                      결과물 열기
                    </a>
                  ) : null}
                  <button
                    type="button"
                    className="actionButton"
                    onClick={() => router.push(`/dashboard/planets/${routeKind}/${planetId}/points/${card.pointId}`)}
                  >
                    제출 위치 보기
                  </button>
                  <span className="disabledNote">커뮤니티 피드백 준비 중</span>
                </div>
              </article>
            )
          })
        ) : (
          <div className="empty">
            {cards.length > 0 ? '이 필터에 표시할 결과물이 없습니다.' : '아직 제출한 탐험결과물이 없습니다.'}
          </div>
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

function ResultCardDetail({ card }: { card: ResultCard }) {
  return (
    <div className="detail">
      <DetailRow label="제출 위치" value={`${card.lessonTitle} · ${card.pointTitle}`} />
      <DetailRow label="지점 유형" value={card.pointType === 'research' ? '연구지점' : '탐험지점'} />
      <DetailRow label="결과물 유형" value={formatArtifactType(card.artifact.artifact_type)} />
      {card.artifact.url.trim() ? <DetailRow label="링크" value={card.artifact.url.trim()} /> : null}
      <DetailRow label="제출일" value={formatResultDate(card.artifact.created_at)} />
      <DetailRow label="최근 수정" value={formatResultDate(card.artifact.updated_at)} />
    </div>
  )
}

function DetailRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="detailRow">
      <span className="detailLabel">{label}</span>
      <span className="detailValue">{value || '기록 없음'}</span>
    </div>
  )
}

function ResultStatus({ children }: { children: string }) {
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

function flattenResultCards(results: PlanetResultAggregate | null): ResultCard[] {
  if (!results) return []
  const cards: ResultCard[] = []
  results.lessons.forEach((lesson: PlanetResultLesson) => {
    lesson.points.forEach((point) => {
      point.artifacts.forEach((artifact) => {
        cards.push({
          id: artifact.id,
          lessonId: lesson.lesson_id,
          lessonTitle: lesson.lesson_title,
          pointId: point.point_id,
          pointTitle: point.point_title,
          pointType: point.point_type,
          artifact,
          category: classifyArtifact(artifact),
        })
      })
    })
  })
  return cards.sort((left, right) => {
    const leftTime = new Date(left.artifact.updated_at || left.artifact.created_at || '').getTime()
    const rightTime = new Date(right.artifact.updated_at || right.artifact.created_at || '').getTime()
    return (Number.isNaN(rightTime) ? 0 : rightTime) - (Number.isNaN(leftTime) ? 0 : leftTime)
  })
}

function buildFilterCounts(cards: ResultCard[]): Record<ResultFilterKey, number> {
  const counts: Record<ResultFilterKey, number> = {
    all: cards.length,
    link: 0,
    image: 0,
    document: 0,
    text: 0,
    other: 0,
  }
  cards.forEach((card) => {
    counts[card.category] += 1
  })
  return counts
}

function classifyArtifact(artifact: PlanetResultArtifact): ResultFilterKey {
  const type = artifact.artifact_type.trim().toLowerCase()
  const url = artifact.url.trim().toLowerCase()
  if (type.includes('이미지') || /\.(png|jpe?g|webp|gif|svg)(\?|#|$)/.test(url)) return 'image'
  if (type.includes('파일') || type.includes('문서') || /\.(pdf|docx?|hwp|hwpx|pptx?|xlsx?|csv|zip)(\?|#|$)/.test(url)) return 'document'
  if (type.includes('텍스트') || type === 'note' || type.includes('노트')) return 'text'
  if (type.includes('링크') || Boolean(url)) return 'link'
  return 'other'
}

function formatArtifactType(value: string): string {
  const type = value.trim()
  if (!type) return '결과물'
  if (type === 'note') return '텍스트'
  return type
}

function formatResultDate(value?: string): string {
  if (!value) return '기록 없음'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleDateString('ko-KR', {
    month: '2-digit',
    day: '2-digit',
  })
}
