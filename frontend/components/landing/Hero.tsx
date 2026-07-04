'use client'

import { useEffect, useRef, useState } from 'react'
import { useRouter } from 'next/navigation'

import { GoalFlowLoadingOverlay } from '@/components/goal-interview/GoalFlowLoadingOverlay'
import { useGoalCreation } from '@/app/dashboard/useGoalCreation'
import type { Locale } from '@/lib/i18n/locales'
import type { LandingPageCopy } from '@/lib/i18n/pages/landing'

const floatingPlanets = [
  { src: '/images/hero_floating_planet_desert_v1.webp', className: 'heroFloatingPlanet heroFloatingPlanetDesert', alt: '사막 크레이터 행성' },
  { src: '/images/hero_floating_planet_lava_v1.webp', className: 'heroFloatingPlanet heroFloatingPlanetLava', alt: '화산 용암 행성' },
  { src: '/images/hero_floating_planet_ringed_v1.webp', className: 'heroFloatingPlanet heroFloatingPlanetRinged', alt: '파스텔 고리 행성' },
  { src: '/images/hero_floating_planet_ice_v1.webp', className: 'heroFloatingPlanet heroFloatingPlanetIce', alt: '얼음 행성' },
  { src: '/images/hero_floating_planet_forest_v1.webp', className: 'heroFloatingPlanet heroFloatingPlanetForest', alt: '숲과 바다 행성' },
] as const

export default function Hero({ copy, locale }: { copy: LandingPageCopy['hero']; locale: Locale }) {
  const router = useRouter()
  const [query, setQuery] = useState('')
  const [isLoggedIn, setIsLoggedIn] = useState(false)
  const [inputHighlight, setInputHighlight] = useState(false)
  const [isInputFocused, setIsInputFocused] = useState(false)
  const heroShellRef = useRef<HTMLElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)
  const inputPulseTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const {
    isGoalSubmitting,
    goalFlowStage,
    goalCreationError: creationError,
    handleGoalSubmit,
    clearGoalError: clearCreationError,
  } = useGoalCreation({
    queryInputFocusFn: () => inputRef.current?.focus(),
  })

  const pulseHeroInput = (delayMs = 0, durationMs = 1200) => {
    if (inputPulseTimeoutRef.current) clearTimeout(inputPulseTimeoutRef.current)
    inputPulseTimeoutRef.current = setTimeout(() => {
      setInputHighlight(true)
      inputPulseTimeoutRef.current = setTimeout(() => {
        setInputHighlight(false)
        inputPulseTimeoutRef.current = null
      }, durationMs)
    }, delayMs)
  }

  const getGoalIntentPath = () => {
    const trimmedQuery = query.trim()
    if (!trimmedQuery) return '/dashboard/goal'
    return `/dashboard/goal?intent=${encodeURIComponent(trimmedQuery)}`
  }

  const getLoginPath = () => (locale === 'en' ? '/en/login' : '/login')

  const routeToLoginForGoal = () => {
    router.push(`${getLoginPath()}?redirect_after=${encodeURIComponent(getGoalIntentPath())}`)
  }

  const updateHeroParallax = (clientX: number, clientY: number) => {
    const shell = heroShellRef.current
    if (!shell) return

    const rect = shell.getBoundingClientRect()
    const x = ((clientX - rect.left) / rect.width - 0.5) * 2
    const y = ((clientY - rect.top) / rect.height - 0.5) * 2
    const depths = [3, -5, 7, -4, 5]

    depths.forEach((depth, index) => {
      shell.style.setProperty(`--hero-planet-x-${index + 1}`, `${(x * depth).toFixed(2)}px`)
      shell.style.setProperty(`--hero-planet-y-${index + 1}`, `${(y * depth * 0.55).toFixed(2)}px`)
    })
    shell.style.setProperty('--hero-beam-x', `${(x * -3.8).toFixed(2)}px`)
    shell.style.setProperty('--hero-beam-y', `${(y * 1.1).toFixed(2)}px`)
  }

  const resetHeroParallax = () => {
    const shell = heroShellRef.current
    if (!shell) return

    for (let index = 1; index <= floatingPlanets.length; index += 1) {
      shell.style.setProperty(`--hero-planet-x-${index}`, '0px')
      shell.style.setProperty(`--hero-planet-y-${index}`, '0px')
    }
    shell.style.setProperty('--hero-beam-x', '0px')
    shell.style.setProperty('--hero-beam-y', '0px')
  }

  const updateHeroLumiHandoff = () => {
    const shell = heroShellRef.current
    if (!shell) return

    const rect = shell.getBoundingClientRect()
    const viewportHeight = window.innerHeight || 1
    const scrollInsideHero = Math.max(0, -rect.top)
    const speechStart = viewportHeight * 0.08
    const speechDistance = viewportHeight * 0.34
    const speechProgress = Math.min(1, Math.max(0, (scrollInsideHero - speechStart) / speechDistance))
    const speechPresence = Math.max(0, 1 - speechProgress)

    if (window.innerWidth < 1024) {
      shell.style.setProperty('--hero-lumi-handoff-progress', '0')
      shell.style.setProperty('--hero-lumi-presence', '1')
      shell.style.setProperty('--hero-lumi-speech-presence', speechPresence.toFixed(3))
      shell.dataset.lumiHandoff = speechProgress > 0 ? 'active' : 'idle'
      return
    }

    const start = viewportHeight * 1.04
    const distance = viewportHeight * 0.52
    const progress = Math.min(1, Math.max(0, (start - rect.bottom) / distance))
    const shipPresence = Math.max(0, 1 - progress * 1.22)
    shell.style.setProperty('--hero-lumi-handoff-progress', progress.toFixed(3))
    shell.style.setProperty('--hero-lumi-presence', shipPresence.toFixed(3))
    shell.style.setProperty('--hero-lumi-speech-presence', speechPresence.toFixed(3))
    shell.dataset.lumiHandoff = progress > 0 || speechProgress > 0 ? 'active' : 'idle'
  }

  useEffect(() => {
    setIsLoggedIn(document.cookie.includes('is_logged_in=1'))
    pulseHeroInput(650, 900)

    return () => {
      if (inputPulseTimeoutRef.current) {
        clearTimeout(inputPulseTimeoutRef.current)
        inputPulseTimeoutRef.current = null
      }
    }
  }, [])

  useEffect(() => {
    const handleMouseMove = (event: MouseEvent) => updateHeroParallax(event.clientX, event.clientY)

    window.addEventListener('mousemove', handleMouseMove)
    window.addEventListener('mouseleave', resetHeroParallax)

    return () => {
      window.removeEventListener('mousemove', handleMouseMove)
      window.removeEventListener('mouseleave', resetHeroParallax)
    }
  }, [])

  useEffect(() => {
    let frame = 0

    const scheduleHandoff = () => {
      if (frame) return
      frame = window.requestAnimationFrame(() => {
        frame = 0
        updateHeroLumiHandoff()
      })
    }

    updateHeroLumiHandoff()
    window.addEventListener('scroll', scheduleHandoff, { passive: true })
    window.addEventListener('resize', scheduleHandoff)

    return () => {
      if (frame) window.cancelAnimationFrame(frame)
      window.removeEventListener('scroll', scheduleHandoff)
      window.removeEventListener('resize', scheduleHandoff)
    }
  }, [])

  useEffect(() => {
    const handleCategorySelect = (event: Event) => {
      const detail = (event as CustomEvent<{ query?: string }>).detail
      const nextQuery = detail?.query?.trim() ?? ''
      setQuery(nextQuery)
      clearCreationError()
      pulseHeroInput(0, 900)
      inputRef.current?.scrollIntoView({ behavior: 'smooth', block: 'center' })
      window.setTimeout(() => inputRef.current?.focus(), 250)
    }

    window.addEventListener('learnweaver:category-select', handleCategorySelect)
    return () => window.removeEventListener('learnweaver:category-select', handleCategorySelect)
  }, [clearCreationError])

  const heroHeadline = locale === 'en'
    ? 'Enter what you want to learn today.'
    : '배우고 싶은 것을 입력해 보세요.'
  const heroWorldLine = locale === 'en'
    ? 'Enter a goal, and your own learning planet opens.'
    : '목표를 입력하면, 나만의 학습 행성이 열립니다.'
  const lumiSpeechLines = locale === 'en'
    ? [
        'Great, shall we build a planet together?',
        'Scroll down and I will show you how.',
      ]
    : [
        '좋아요, 같이 행성을 만들어 볼까요?',
        '아래에서 사용법을 알려드릴께요.',
      ]


  const handleLaunch = async () => {
    if (isGoalSubmitting) return

    if (!query.trim()) {
      await handleGoalSubmit(query)
      return
    }

    if (!isLoggedIn) {
      routeToLoginForGoal()
      return
    }

    try {
      const refreshRes = await fetch('/api/v1/auth/refresh', {
        method: 'POST',
        credentials: 'include',
      })

      if (!refreshRes.ok) {
        document.cookie = 'is_logged_in=; Max-Age=0; path=/;'
        routeToLoginForGoal()
        return
      }

      await handleGoalSubmit(query)
    } catch {
      pulseHeroInput(0, 900)
    }
  }

  return (
    <section
      ref={heroShellRef}
      className="heroShell relative flex min-h-[100svh] items-start justify-center overflow-hidden pb-10 pt-20 md:items-center md:pb-0"
    >
      {isGoalSubmitting && goalFlowStage ? (
        <GoalFlowLoadingOverlay stage={goalFlowStage} />
      ) : null}

      <div aria-hidden="true" className="heroGeneratedSpaceBg" />
      <div aria-hidden="true" className="heroSkyWash" />
      <div aria-hidden="true" className="heroStarSprinkles" />
      <div aria-hidden="true" className="heroFloatingPlanetLayer">
        {floatingPlanets.map((planet) => (
          <img key={planet.src} src={planet.src} alt="" className={planet.className} />
        ))}
      </div>
      <div aria-hidden="true" className="heroHorizonGlow" />
      <svg
        aria-hidden="true"
        className="heroGuideBridgeSvg"
        viewBox="0 0 1120 820"
        preserveAspectRatio="none"
      >
        <defs>
          <linearGradient id="heroGuideBeamStroke" x1="0" x2="1" y1="0" y2="1">
            <stop offset="0%" stopColor="#FFFFFF" stopOpacity="0.18" />
            <stop offset="28%" stopColor="#AEEFFF" stopOpacity="0.78" />
            <stop offset="58%" stopColor="#6ED4F3" stopOpacity="0.9" />
            <stop offset="100%" stopColor="#FFFFFF" stopOpacity="0.3" />
          </linearGradient>
        </defs>
        <circle className="heroGuideBeamSparkGlow" cx="548" cy="78" r="7" />
        <circle className="heroGuideBeamSpark" cx="548" cy="78" r="2" />
        <path className="heroGuideBeamGlow heroGuideBeamFar" d="M 572 78 C 500 104 384 124 300 166" />
        <path className="heroGuideBeamGlow heroGuideBeamMid" d="M 300 166 C 152 250 86 382 208 486 C 270 540 376 555 460 564" />
        <path className="heroGuideBeamGlow heroGuideBeamNear" d="M 460 564 C 580 574 670 606 682 688 C 692 760 552 790 514 820" />
        <path className="heroGuideBeamBase heroGuideBeamFar" d="M 572 78 C 500 104 384 124 300 166" />
        <path className="heroGuideBeamBase heroGuideBeamMid" d="M 300 166 C 152 250 86 382 208 486 C 270 540 376 555 460 564" />
        <path className="heroGuideBeamBase heroGuideBeamNear" d="M 460 564 C 580 574 670 606 682 688 C 692 760 552 790 514 820" />
        <path className="heroGuideBeamCore heroGuideBeamFar" d="M 572 78 C 500 104 384 124 300 166" />
        <path className="heroGuideBeamCore heroGuideBeamMid" d="M 300 166 C 152 250 86 382 208 486 C 270 540 376 555 460 564" />
        <path className="heroGuideBeamCore heroGuideBeamNear" d="M 460 564 C 580 574 670 606 682 688 C 692 760 552 790 514 820" />
      </svg>

      <div className="heroStage relative z-10 mx-auto grid w-full max-w-7xl items-start gap-7 px-4 pb-8 pt-7 sm:px-5 md:px-8 md:py-12 lg:grid-cols-[0.92fr_1.08fr] lg:items-center lg:gap-6 lg:py-10">
        <div className="heroCopyColumn">
          <h1 className="sr-only">LearnCosmos</h1>
          <p className="heroLeadText heroLeadTextSingle">{heroHeadline}</p>
          <div className="heroBrandGlowLine" />

          <div className="heroCTAStack w-full max-w-[580px]">
            <div
              className={`heroMysticInput mb-[11px] flex items-center rounded-[16px] pl-[19px] pr-2 py-2 ${inputHighlight ? 'heroMysticInputPulse' : ''} ${isInputFocused ? 'heroMysticInputActive' : ''}`}
              style={{
                border: isInputFocused
                  ? '1px solid rgba(246,183,60,0.9)'
                  : inputHighlight
                    ? '1px solid rgba(58,169,199,0.72)'
                    : '1px solid rgba(61,126,162,0.28)',
                boxShadow: isInputFocused
                  ? '0 0 0 4px rgba(246,183,60,0.18), 0 18px 42px rgba(42,94,139,0.18)'
                  : inputHighlight
                    ? '0 0 0 4px rgba(106,210,193,0.16), 0 16px 34px rgba(42,94,139,0.14)'
                    : '0 12px 28px rgba(42,94,139,0.12)',
                transition: 'border 0.4s ease, box-shadow 0.4s ease, transform 0.4s ease, filter 0.4s ease',
              }}
            >
              <input
                ref={inputRef}
                id="hero-query-input"
                type="text"
                value={query}
                placeholder={copy.placeholder}
                className="relative z-[1] w-full flex-1 border-none bg-transparent py-[12px] text-[15.5px] outline-none"
                style={{
                  color: '#000000',
                  caretColor: '#111111',
                }}
                onChange={(event) => {
                  setQuery(event.target.value)
                  clearCreationError()
                }}
                onFocus={() => {
                  clearCreationError()
                  setIsInputFocused(true)
                }}
                onBlur={() => setIsInputFocused(false)}
                onKeyDown={(event) => {
                  if (event.key === 'Enter') {
                    event.preventDefault()
                    void handleLaunch()
                  }
                }}
              />
              <button
                type="button"
                className="relative z-[1] flex-shrink-0 rounded-[11px] px-[19px] py-[11px] text-[14px] font-bold text-white whitespace-nowrap"
                style={{
                  background: 'linear-gradient(135deg, #FFB84D 0%, #FF7A3D 52%, #4DB7E8 100%)',
                  boxShadow: '0 10px 20px rgba(255,122,61,0.24), inset 0 1px 0 rgba(255,255,255,0.38)',
                }}
                onClick={() => void handleLaunch()}
              >
                {copy.submitLabel}
              </button>
            </div>


            <p className="heroWorldLine">{heroWorldLine}</p>

            {creationError ? (
              <div className="heroCreationError">
                {creationError}
              </div>
            ) : null}

            <div className="flex flex-wrap gap-[7px]">
              {copy.chips.map((chip) => (
                <button
                  key={chip}
                  type="button"
                  className="rounded-full px-[14px] py-[6px] text-[13px] transition hover:-translate-y-0.5"
                  style={{
                    background: 'rgba(255,255,255,0.78)',
                    border: '1px solid rgba(77,183,232,0.24)',
                    color: '#3B6478',
                    boxShadow: '0 7px 16px rgba(42,94,139,0.08)',
                  }}
                  onClick={() => {
                    setQuery(chip)
                    clearCreationError()
                    pulseHeroInput(0, 720)
                    inputRef.current?.focus()
                  }}
                >
                  {chip}
                </button>
              ))}
            </div>

            <div className="heroCompanionGuide" aria-label={lumiSpeechLines.join(' ')}>
              <div className="heroCompanionSpeech">
                {lumiSpeechLines.map((line) => (
                  <span key={line}>{line}</span>
                ))}
              </div>
              <img
                src="/images/mainlumi.webp"
                alt="우주선 위에 앉아 손을 흔드는 루미"
                className="heroCompanionLumi"
              />
            </div>
          </div>
        </div>


      </div>
        <div className="heroVisualStage" aria-label="테라포밍 중인 행성 위를 루미가 우주선으로 탐험하는 장면">
          <div aria-hidden="true" className="heroTerraformAura" />
          <svg
            aria-hidden="true"
            className="heroLearningRouteSvg"
            viewBox="0 0 660 760"
            preserveAspectRatio="none"
          >
            <defs>
              <linearGradient id="heroRouteStroke" x1="170" y1="135" x2="460" y2="720" gradientUnits="userSpaceOnUse">
                <stop offset="0" stopColor="#4DB7E8" stopOpacity="0.72" />
                <stop offset="0.45" stopColor="#9B7CFF" stopOpacity="0.42" />
                <stop offset="1" stopColor="#FFB84D" stopOpacity="0.1" />
              </linearGradient>
              <filter id="heroRouteGlow" x="-20%" y="-20%" width="140%" height="140%">
                <feGaussianBlur stdDeviation="3.2" result="blur" />
                <feMerge>
                  <feMergeNode in="blur" />
                  <feMergeNode in="SourceGraphic" />
                </feMerge>
              </filter>
            </defs>
            <ellipse
              cx="330"
              cy="330"
              rx="244"
              ry="92"
              transform="rotate(-14 330 330)"
              className="heroLearningRouteOrbit"
            />
          </svg>
          <div aria-hidden="true" className="heroOrbitLine heroOrbitLineOne" />
          <div aria-hidden="true" className="heroOrbitLine heroOrbitLineTwo" />
          <img
            src="/images/hero_terraforming_planet.webp"
            alt="테라포밍 중인 행성"
            className="heroTerraformPlanet"
          />
          <div className="heroLumiSpeech" aria-label={lumiSpeechLines.join(' ')}>
            {lumiSpeechLines.map((line) => (
              <span key={line}>{line}</span>
            ))}
          </div>
          <img
            src="/images/mainlumi.webp"
            alt="우주선 위에 앉아 손을 흔드는 루미"
            className="heroLumiShip"
          />
        </div>

    </section>
  )
}
