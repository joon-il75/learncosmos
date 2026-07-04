'use client'

import { useState, useEffect } from 'react'
import type { LandingPageCopy } from '@/lib/i18n/pages/landing'

const API_BASE = ''

export default function AIUsage({ copy }: { copy: LandingPageCopy['aiUsage'] }) {
  const [welcomePoints, setWelcomePoints] = useState(30)
  const [courseGenCost, setCourseGenCost] = useState(5)
  const [lessonRecCost, setLessonRecCost] = useState(1)

  useEffect(() => {
    fetch(`${API_BASE}/api/v1/public/point-policy`)
      .then(r => r.ok ? r.json() : null)
      .then(d => {
        if (!d) return
        if (d.welcome_points)  setWelcomePoints(d.welcome_points)
        if (d.course_gen_cost) setCourseGenCost(d.course_gen_cost)
        if (d.lesson_rec_cost) setLessonRecCost(d.lesson_rec_cost)
      })
      .catch(() => {})
  }, [])

  return (
    <section
      id="ai-usage"
      className="bg-[linear-gradient(180deg,#0D314E_0%,#222C4D_100%)] px-6 py-24"
    >
      <div className="max-w-6xl mx-auto">

        {/* 섹션 헤더 */}
        <div className="text-center mb-16">
          <span
            className="inline-block text-sm font-medium px-4 py-1.5 rounded-full mb-4"
            style={{
              backgroundColor: 'rgba(106, 210, 193, 0.12)',
              color: '#6AD2C1',
              border: '1px solid rgba(106, 210, 193, 0.26)',
            }}
          >
            {copy.eyebrow}
          </span>
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">
            {copy.titlePrefix} <span style={{ color: '#6AD2C1' }}>{copy.titleHighlight}</span>
          </h2>
          <p className="text-lg leading-relaxed max-w-2xl mx-auto" style={{ color: '#9DB4CC' }}>
            {copy.descriptionLines[0]}
            <br className="hidden md:block" />
            {' '}
            {copy.descriptionLines[1]}
          </p>
        </div>

        {/* 2열 카드 */}
        <div className="grid md:grid-cols-2 gap-6 mb-8">

          {/* BYOK 카드 (강조) */}
          <div
            className="relative rounded-2xl p-8 overflow-hidden"
            style={{
              background: 'linear-gradient(135deg, rgba(106,210,193,0.13) 0%, rgba(140,93,152,0.12) 100%)',
              border: '1px solid rgba(106, 210, 193, 0.30)',
            }}
          >
            {/* 추천 뱃지 */}
            <div className="absolute top-4 right-4">
              <span
                className="text-xs font-semibold px-3 py-1 rounded-full"
                style={{ backgroundColor: '#1C7D79', color: '#ffffff' }}
              >
                {copy.byokBadge}
              </span>
            </div>

            {/* 아이콘 */}
            <div
              className="w-12 h-12 rounded-xl flex items-center justify-center mb-5"
              style={{ backgroundColor: 'rgba(106, 210, 193, 0.16)' }}
            >
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <path
                  d="M15 7C15 4.79 13.21 3 11 3C8.79 3 7 4.79 7 7C7 9.21 8.79 11 11 11C13.21 11 15 9.21 15 7Z"
                  stroke="#6AD2C1"
                  strokeWidth="1.8"
                  strokeLinecap="round"
                />
                <path
                  d="M11 11L11 21M8 18H14M8 14H14"
                  stroke="#6AD2C1"
                  strokeWidth="1.8"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            </div>

            <h3 className="text-xl font-bold text-white mb-2">
              {copy.byokTitle}
            </h3>
            <p className="text-sm mb-6" style={{ color: '#9DB4CC', lineHeight: '1.7' }}>
              {copy.byokDescriptionLines.map((line) => (
                <span key={line} className="block">{line}</span>
              ))}
            </p>

            <ul className="space-y-2.5">
              {copy.byokBullets.map((item) => (
                <li key={item} className="flex items-center gap-2.5 text-sm" style={{ color: '#C8D8E8' }}>
                  <svg width="16" height="16" viewBox="0 0 16 16" fill="none" className="shrink-0">
                    <circle cx="8" cy="8" r="7" fill="rgba(106,210,193,0.16)" stroke="rgba(106,210,193,0.42)" strokeWidth="0.8"/>
                    <path d="M5 8L7 10L11 6" stroke="#6AD2C1" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
                  </svg>
                  {item}
                </li>
              ))}
            </ul>
          </div>

          {/* 크레딧 카드 */}
          <div
            className="rounded-2xl p-8"
            style={{
              background: 'rgba(255, 255, 255, 0.04)',
              border: '1px solid rgba(245, 165, 36, 0.20)',
            }}
          >
            {/* 아이콘 */}
            <div
              className="w-12 h-12 rounded-xl flex items-center justify-center mb-5"
              style={{ backgroundColor: 'rgba(245, 165, 36, 0.14)' }}
            >
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="9" stroke="#F5A524" strokeWidth="1.8"/>
                <path
                  d="M12 7V17M9.5 9.5C9.5 8.4 10.6 7.5 12 7.5C13.4 7.5 14.5 8.4 14.5 9.5C14.5 11.5 9.5 11 9.5 13C9.5 14.2 10.6 15 12 15C13.4 15 14.5 14.2 14.5 13"
                  stroke="#F5A524"
                  strokeWidth="1.8"
                  strokeLinecap="round"
                />
              </svg>
            </div>

            <h3 className="text-xl font-bold text-white mb-2">
              {copy.pointsTitle}
            </h3>
            <p className="text-sm mb-6" style={{ color: '#9DB4CC', lineHeight: '1.7' }}>
              {copy.pointsDescription(welcomePoints).map((line, index) => (
                <span key={line} className="block">
                  {index === 1 ? <strong style={{ color: '#F5A524' }}>{line}</strong> : line}
                </span>
              ))}
            </p>

            <div
              className="rounded-xl p-4 mb-5 space-y-3"
              style={{ backgroundColor: 'rgba(255,255,255,0.05)' }}
            >
              <p className="text-xs font-medium mb-3" style={{ color: '#7A94AA' }}>
                {copy.pointExampleLabel}
              </p>
              {copy.pointExamples.map(({ key, label, icon }) => {
                const point = `${key === 'course' ? courseGenCost : lessonRecCost}pt`
                return (
                <div key={label} className="flex items-center justify-between">
                  <span className="text-sm flex items-center gap-2" style={{ color: '#C8D8E8' }}>
                    <span>{icon}</span>
                    {label}
                  </span>
                  <span
                    className="text-xs font-bold px-2.5 py-0.5 rounded-full"
                    style={{ backgroundColor: 'rgba(245,165,36,0.14)', color: '#F5A524' }}
                  >
                    {point}
                  </span>
                </div>
              )})}
            </div>

            <ul className="space-y-2.5">
              {copy.pointBullets.map((item) => (
                <li key={item} className="flex items-center gap-2.5 text-sm" style={{ color: '#C8D8E8' }}>
                  <svg width="16" height="16" viewBox="0 0 16 16" fill="none" className="shrink-0">
                    <circle cx="8" cy="8" r="7" fill="rgba(245,165,36,0.14)" stroke="rgba(245,165,36,0.38)" strokeWidth="0.8"/>
                    <path d="M5 8L7 10L11 6" stroke="#F5A524" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
                  </svg>
                  {item}
                </li>
              ))}
            </ul>
          </div>
        </div>

        <div
          className="mb-12 rounded-2xl px-5 py-4 text-center text-sm leading-relaxed"
          style={{
            background: 'rgba(255,255,255,0.035)',
            border: '1px solid rgba(157, 180, 204, 0.14)',
            color: '#9DB4CC',
          }}
        >
          {copy.sensitiveNoticeLines[0]}
          <br className="hidden md:block" />
          {' '}
          {copy.sensitiveNoticeLines[1]}
        </div>

        {/* 지원 LLM 배지 행 */}
        <div className="text-center">
          <p className="text-xs mb-5" style={{ color: '#5A7A94' }}>
            {copy.futureProvidersLabel}
          </p>
          <div className="flex flex-wrap justify-center items-center gap-3">
            {[
              { name: 'OpenAI',        color: '#10A37F', bg: 'rgba(16,163,127,0.1)',  border: 'rgba(16,163,127,0.25)' },
              { name: 'Claude',        color: '#D97757', bg: 'rgba(217,119,87,0.1)',  border: 'rgba(217,119,87,0.25)' },
              { name: 'Gemini',        color: '#4285F4', bg: 'rgba(66,133,244,0.1)',  border: 'rgba(66,133,244,0.25)' },
              { name: 'Grok',          color: '#E8EAF2', bg: 'rgba(232,234,242,0.08)', border: 'rgba(232,234,242,0.2)' },
              { name: 'HyperCLOVA X', color: '#03C75A', bg: 'rgba(3,199,90,0.1)',    border: 'rgba(3,199,90,0.25)'  },
              { name: 'Solar',         color: '#FF6B35', bg: 'rgba(255,107,53,0.1)',  border: 'rgba(255,107,53,0.25)' },
              { name: 'Llama',         color: '#9B59B6', bg: 'rgba(155,89,182,0.1)',  border: 'rgba(155,89,182,0.25)' },
              { name: 'EXAONE',        color: '#66BDF2', bg: 'rgba(102,189,242,0.1)',  border: 'rgba(102,189,242,0.25)' },
            ].map(({ name, color, bg, border }) => (
              <span
                key={name}
                className="text-xs font-medium px-3 py-1.5 rounded-full"
                style={{ backgroundColor: bg, color, border: `1px solid ${border}` }}
              >
                {name}
              </span>
            ))}
          </div>
        </div>

      </div>
    </section>
  )
}
