'use client'

import { useEffect, useRef, useState } from 'react'
import type { CSSProperties } from 'react'

import type { LandingPageCopy } from '@/lib/i18n/pages/landing';

const guideImages = [
  '/images/lumi_scroll_step1__alpha_v2.webp',
  '/images/lumi_scroll_step2_goal_chat_ship_behind_alpha_v2.png',
  '/images/lumi_scroll_step3_content_search_no_ship_alpha_v2.png',
  '/images/lumi_scroll_step4_learning_display_alpha_v1.png',
] as const;

const stepTwoSuccessImage = '/images/lumi_step2_success_pose_flipped_alpha_v1.png';

const guideImageFacingCorrections = [
  1,
  1,
  1,
  1,
] as const;

type GuideStep = LandingPageCopy['lumiScrollGuide']['steps'][number];
type StepOnePhase = 'intro' | 'typing' | 'click' | 'next';
type StepTwoPhase = 'question' | 'answer' | 'summary' | 'card' | 'click' | 'success';

const stepOneMessages = {
  intro: '주제를 입력하면 시작점이 열려요',
  next: '목표 채팅으로 이어집니다.',
} as const;

function getStepOnePhase(progress: number): StepOnePhase {
  if (progress < 0.12) return 'intro'
  if (progress < 0.70) return 'typing'
  if (progress < 0.84) return 'click'
  return 'next'
}

function getStepOneTypedText(progress: number) {
  const typingProgress = Math.min(1, Math.max(0, (progress - 0.12) / 0.58))
  if (typingProgress < 0.14) return ''
  if (typingProgress < 0.28) return '기'
  if (typingProgress < 0.46) return '기타'
  if (typingProgress < 0.62) return '기타 '
  if (typingProgress < 0.80) return '기타 입'
  return '기타 입문'
}

function getStepTwoPhase(progress: number): StepTwoPhase {
  if (progress < 0.23) return 'question'
  if (progress < 0.41) return 'answer'
  if (progress < 0.55) return 'summary'
  if (progress < 0.65) return 'card'
  if (progress < 0.80) return 'click'
  return 'success'
}

function StepOneCtaDemo({
  step,
  phase,
  typedText,
  actionActive,
  pinned = false,
}: {
  step: GuideStep
  phase: StepOnePhase
  typedText: string
  actionActive: boolean
  pinned?: boolean
}) {
  return (
    <div
      className={`lumiScrollGuideDemo lumiStepOneCtaDemo${pinned ? ' lumiStepOneCtaDemoPinned' : ''}`}
      data-step-one-phase={phase}
      data-step-one-action={actionActive ? 'active' : 'idle'}
      aria-label="학습탐험 시작 입력 예시"
    >
      <div className="lumiScrollGuideDemoTop lumiStepOneCtaTop">
        <span>{step.badge}</span>
        <strong>{step.title}</strong>
      </div>
      <div className="lumiStepOneInputMock" aria-hidden="true">
        <span className="lumiStepOneTypedText">{typedText}</span>
        <span className="lumiStepOneCaret" />
        <span className="lumiStepOneSubmitMock">{step.demoAction}</span>
        <svg className="lumiStepOnePointer" viewBox="0 0 54 58" aria-hidden="true">
          <circle className="lumiStepOnePointerRing" cx="39" cy="40" r="10" />
          <path className="lumiStepOnePointerShape" d="M11 5 L39 31 L27 34 L34 50 L27 53 L20 37 L11 47 Z" />
        </svg>
      </div>
      <div className="lumiStepOneHintStack" aria-hidden="true">
        <span className="lumiStepOneHint lumiStepOneHintTopic">{stepOneMessages.intro}</span>
        <span className="lumiStepOneHint lumiStepOneHintNext">{stepOneMessages.next}</span>
      </div>
    </div>
  );
}


function StepTwoGoalChatDemo({
  step,
  phase,
  actionActive,
  planetFrameIndex,
}: {
  step: GuideStep
  phase: StepTwoPhase
  actionActive: boolean
  planetFrameIndex: number
}) {
  const planetFrameSrc = `/images/lumi_step2_planet_success_frame_${String(planetFrameIndex).padStart(2, '0')}.png`
  return (
    <div
      className="lumiScrollGuideDemo lumiScrollGuideDemoChat lumiStepTwoGoalChatDemo"
      data-step-two-phase={phase}
      data-step-two-action={actionActive ? 'active' : 'idle'}
      aria-label="목표 채팅 예시"
    >
      <div className="lumiScrollGuideDemoTop">
        <span>{step.badge}</span>
        <strong>{step.demoTitle}</strong>
      </div>
      <div className="lumiStepTwoChatStack">
        <div className="lumiScrollGuideChatLine lumiScrollGuideChatLineLumi lumiStepTwoChatQuestion">가장 해내고 싶은 목표는 무엇인가요?</div>
        <div className="lumiScrollGuideChatLine lumiScrollGuideChatLineUser lumiStepTwoChatAnswer">3개월 안에 좋아하는 곡을 연주하고 싶어</div>
        <div className="lumiScrollGuideChatLine lumiScrollGuideChatLineLumi lumiStepTwoChatSummary">좋아요. 이 목표로 탐험계획을 만들어볼까요?</div>
      </div>
      <div className="lumiStepTwoGoalOffer">
        <div className="lumiStepTwoGoalOfferLabel">목표 제안</div>
        <strong>3개월 안에 좋아하는 곡 연주하기</strong>
        <div className="lumiStepTwoGoalActions">
          <span className="lumiStepTwoGoalStart">이 목표로 시작하기</span>
          <span className="lumiStepTwoGoalReset">다시 정하기</span>
        </div>
        <svg className="lumiStepTwoPointer" viewBox="0 0 54 58" aria-hidden="true">
          <circle className="lumiStepTwoPointerRing" cx="39" cy="40" r="10" />
          <path className="lumiStepTwoPointerShape" d="M11 5 L39 31 L27 34 L34 50 L27 53 L20 37 L11 47 Z" />
        </svg>
      </div>
      <div className="lumiStepTwoSuccessScene" aria-hidden={phase !== 'success'}>
        <div className="lumiStepTwoPlanetStage">
          <span className="lumiStepTwoPlanetGlow" />
          <span className="lumiStepTwoPlanetSpark" />
          <img src={planetFrameSrc} alt="" className="lumiStepTwoPlanetFrame" />
        </div>
        <div className="lumiStepTwoSuccessMessage">3개월 안에 좋아하는 곡 연주하기 행성 생성 성공</div>
      </div>
    </div>
  );
}

function StepDemo({
  step,
  index,
  stepOnePhase = 'intro',
  stepOneTypedText = '',
  stepOneActionActive = false,
  stepTwoPhase = 'question',
  stepTwoActionActive = false,
  stepTwoPlanetFrameIndex = 0,
  pinned = false,
}: {
  step: GuideStep
  index: number
  stepOnePhase?: StepOnePhase
  stepOneTypedText?: string
  stepOneActionActive?: boolean
  stepTwoPhase?: StepTwoPhase
  stepTwoActionActive?: boolean
  stepTwoPlanetFrameIndex?: number
  pinned?: boolean
}) {
  if (index === 0) {
    return <StepOneCtaDemo step={step} phase={stepOnePhase} typedText={stepOneTypedText} actionActive={stepOneActionActive} pinned={pinned} />;
  }

  if (index === 1) {
    return <StepTwoGoalChatDemo step={step} phase={stepTwoPhase} actionActive={stepTwoActionActive} planetFrameIndex={stepTwoPlanetFrameIndex} />;
  }

  if (index === 2) {
    return (
      <div className="lumiScrollGuideDemo lumiScrollGuideDemoPath">
        <div className="lumiScrollGuideDemoTop">
          <span>{step.badge}</span>
          <strong>{step.demoTitle}</strong>
        </div>
        <div className="lumiScrollGuidePathList">
          {step.demoItems.map((item, itemIndex) => (
            <div className="lumiScrollGuidePathItem" key={item}>
              <span>{itemIndex + 1}</span>
              <p>{item}</p>
            </div>
          ))}
        </div>
        <button className="lumiScrollGuideDemoAction" type="button">{step.demoAction}</button>
      </div>
    );
  }

  return (
    <div className="lumiScrollGuideDemo lumiScrollGuideDemoProgress">
      <div className="lumiScrollGuideDemoTop">
        <span>{step.badge}</span>
        <strong>{step.demoTitle}</strong>
      </div>
      <div className="lumiScrollGuideProgressRing" aria-hidden="true">68%</div>
      <div className="lumiScrollGuideProgressList">
        {step.demoItems.map((item) => (
          <div key={item}>{item}</div>
        ))}
      </div>
      <button className="lumiScrollGuideDemoAction" type="button">{step.demoAction}</button>
    </div>
  );
}

export default function LumiScrollGuide({ copy }: { copy: LandingPageCopy['lumiScrollGuide'] }) {
  const [activeStep, setActiveStep] = useState(0)
  const [stepOnePhase, setStepOnePhase] = useState<StepOnePhase>('intro')
  const [stepOneTypedText, setStepOneTypedText] = useState('')
  const [stepOneActionActive, setStepOneActionActive] = useState(false)
  const [stepTwoPhase, setStepTwoPhase] = useState<StepTwoPhase>('question')
  const [stepTwoActionActive, setStepTwoActionActive] = useState(false)
  const [stepTwoSuccessComplete, setStepTwoSuccessComplete] = useState(false)
  const [stepTwoPlanetFrameIndex, setStepTwoPlanetFrameIndex] = useState(0)
  const sectionRef = useRef<HTMLElement | null>(null)
  const stickySlotRef = useRef<HTMLElement | null>(null)
  const stickyGuideRef = useRef<HTMLDivElement | null>(null)
  const desktopLayoutRef = useRef<HTMLDivElement | null>(null)
  const pinnedDemoRef = useRef<HTMLDivElement | null>(null)
  const stepRefs = useRef<Array<HTMLElement | null>>([])
  const activeGuideStep = copy.steps[activeStep] ?? copy.steps[0]
  const activePinnedStep = copy.steps[activeStep] ?? copy.steps[0]
  const activeGuideImage = activeStep === 1 && stepTwoSuccessComplete
    ? stepTwoSuccessImage
    : guideImages[activeStep] ?? guideImages[0]
  const activeGuideImageFacingCorrection = guideImageFacingCorrections[activeStep] ?? 1
  const isStepOneActive = activeStep === 0
  const activeGuideDescription = isStepOneActive
    ? stepOneMessages.intro
    : activeGuideStep.description
  const showStepOneNextNote = isStepOneActive && stepOnePhase === 'next'

  useEffect(() => {
    const steps = stepRefs.current.filter((step): step is HTMLElement => Boolean(step))
    if (!steps.length || typeof IntersectionObserver === 'undefined') return

    const visibility = new Map<number, number>()
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          const index = Number((entry.target as HTMLElement).dataset.stepIndex ?? 0)
          visibility.set(index, entry.isIntersecting ? entry.intersectionRatio : 0)
        })

        const [nextActive, nextRatio] = Array.from(visibility.entries()).reduce(
          (best, current) => (current[1] > best[1] ? current : best),
          [0, 0] as [number, number],
        )

        if (nextRatio > 0) {
          setActiveStep((current) => (current === nextActive ? current : nextActive))
        }
      },
      {
        root: null,
        rootMargin: '-28% 0px -38% 0px',
        threshold: [0, 0.18, 0.32, 0.48, 0.64, 0.82],
      },
    )

    steps.forEach((step, index) => {
      visibility.set(index, index === 0 ? 1 : 0)
      observer.observe(step)
    })

    return () => observer.disconnect()
  }, [])

  useEffect(() => {
    const section = sectionRef.current
    if (!section) return

    let frame = 0

    const updateScrollState = () => {
      frame = 0
      const viewportHeight = window.innerHeight || 1
      const pinTop = 104
      const layoutRectForProgress = desktopLayoutRef.current?.getBoundingClientRect()
      const stepOneCanAdvance = window.innerWidth < 1024 || !layoutRectForProgress || layoutRectForProgress.top <= pinTop
      section.style.setProperty('--lumi-guide-image-facing-x', activeGuideImageFacingCorrection.toString())

      let nextStepOnePhase: StepOnePhase = 'intro'
      let nextStepOneTypedText = ''
      let nextStepOneActionActive = false
      let stepOneProgress = 0
      let nextStepTwoPhase: StepTwoPhase = 'question'
      let nextStepTwoActionActive = false
      let stepTwoProgress = 0
      let stepTwoSuccessProgress = 0
      let nextStepTwoSuccessComplete = false
      let nextStepTwoPlanetFrameIndex = 0

      stepRefs.current.forEach((step, index) => {
        if (!step) return

        const rect = step.getBoundingClientRect()
        const start = viewportHeight * 0.82
        const distance = viewportHeight * 0.64
        const progress = Math.min(1, Math.max(0, (start - rect.top) / distance))
        const dashOffset = (100 - progress * 100).toFixed(3)
        section.style.setProperty(`--lumi-route-progress-${index + 1}`, progress.toFixed(3))
        section.style.setProperty(`--lumi-route-offset-${index + 1}`, dashOffset)

        if (index === 0) {
          if (stepOneCanAdvance) {
            const pinProgressDistance = Math.max(420, viewportHeight * 0.5)
            if (window.innerWidth >= 1024 && layoutRectForProgress) {
              stepOneProgress = Math.min(1, Math.max(0, (pinTop - layoutRectForProgress.top) / pinProgressDistance))
            } else {
              const phaseStart = viewportHeight * 0.78
              const phaseDistance = Math.max(1, viewportHeight * 1.08)
              stepOneProgress = Math.min(1, Math.max(0, (phaseStart - rect.top) / phaseDistance))
            }
          } else {
            stepOneProgress = 0
          }
          nextStepOnePhase = getStepOnePhase(stepOneProgress)
          nextStepOneTypedText = getStepOneTypedText(stepOneProgress)
          nextStepOneActionActive = nextStepOnePhase === 'click' || nextStepOnePhase === 'next'
        }

        if (index === 1) {
          const phaseStart = viewportHeight * 0.64
          const phaseDistance = Math.max(1, viewportHeight * 1.62)
          stepTwoProgress = Math.min(1, Math.max(0, (phaseStart - rect.top) / phaseDistance))
          stepTwoSuccessProgress = Math.min(1, Math.max(0, (stepTwoProgress - 0.8) / 0.2))
          nextStepTwoPlanetFrameIndex = Math.min(11, Math.max(0, Math.round(stepTwoSuccessProgress * 11)))
          nextStepTwoSuccessComplete = stepTwoSuccessProgress >= 0.96
          nextStepTwoPhase = getStepTwoPhase(stepTwoProgress)
          nextStepTwoActionActive = nextStepTwoPhase === 'click'
        }
      })

      section.style.setProperty('--lumi-step1-progress', stepOneProgress.toFixed(3))
      section.style.setProperty('--lumi-step2-progress', stepTwoProgress.toFixed(3))
      section.style.setProperty('--lumi-step2-success-progress', stepTwoSuccessProgress.toFixed(3))
      section.style.setProperty('--lumi-step1-demo-presence', activeStep === 0 ? '1' : '0')
      section.dataset.stepOnePhase = activeStep === 0 ? nextStepOnePhase : 'intro'
      section.dataset.stepOneAction = activeStep === 0 && nextStepOneActionActive ? 'active' : 'idle'
      section.dataset.stepTwoPhase = activeStep === 1 ? nextStepTwoPhase : 'question'
      section.dataset.stepTwoSuccess = activeStep === 1 && nextStepTwoSuccessComplete ? 'complete' : 'forming'
      section.dataset.stepTwoAction = activeStep === 1 && nextStepTwoActionActive ? 'active' : 'idle'
      setStepOnePhase((current) => (current === nextStepOnePhase ? current : nextStepOnePhase))
      setStepOneTypedText((current) => (current === nextStepOneTypedText ? current : nextStepOneTypedText))
      setStepOneActionActive((current) => (current === nextStepOneActionActive ? current : nextStepOneActionActive))
      setStepTwoPhase((current) => (current === nextStepTwoPhase ? current : nextStepTwoPhase))
      setStepTwoActionActive((current) => (current === nextStepTwoActionActive ? current : nextStepTwoActionActive))
      setStepTwoSuccessComplete((current) => (current === nextStepTwoSuccessComplete ? current : nextStepTwoSuccessComplete))
      setStepTwoPlanetFrameIndex((current) => (current === nextStepTwoPlanetFrameIndex ? current : nextStepTwoPlanetFrameIndex))

      const slot = stickySlotRef.current
      const guide = stickyGuideRef.current
      const pinnedDemo = pinnedDemoRef.current
      if (!slot || !guide || window.innerWidth < 1024) {
        section.style.setProperty('--lumi-guide-handoff-progress', '1')
        section.style.setProperty('--lumi-guide-handoff-presence', '1')
        section.style.setProperty('--lumi-guide-handoff-x', '0px')
        section.style.setProperty('--lumi-guide-handoff-y', '0px')
        section.style.setProperty('--lumi-guide-state-scale', '1')
        section.style.setProperty('--lumi-guide-facing-x', '-1')
        section.style.setProperty('--lumi-guide-image-arrival-x', '0px')
        section.style.setProperty('--lumi-pinned-demo-left', '0px')
        section.style.setProperty('--lumi-pinned-demo-width', '0px')
        section.style.setProperty('--lumi-pinned-demo-height', '0px')
        section.style.setProperty('--lumi-pinned-demo-top', '188px')
        section.dataset.guidePin = 'flow'
        return
      }

      const sectionRect = section.getBoundingClientRect()
      const slotRect = slot.getBoundingClientRect()
      const guideRect = guide.getBoundingClientRect()
      const layoutRect = desktopLayoutRef.current?.getBoundingClientRect()
      const firstStepRect = stepRefs.current[0]?.getBoundingClientRect()
      const pinnedDemoRect = pinnedDemo?.getBoundingClientRect()
      const handoffStart = viewportHeight * 0.70
      const handoffDistance = Math.max(1, handoffStart - pinTop)
      const guideHandoffProgress = Math.min(1, Math.max(0, (handoffStart - sectionRect.top) / handoffDistance))
      const slotCenter = slotRect.left + slotRect.width / 2
      const heroShip = document.querySelector<HTMLElement>('.heroLumiShip')
      const heroShipRect = heroShip?.getBoundingClientRect()
      const guideStartCenter = heroShipRect ? heroShipRect.left + heroShipRect.width / 2 : window.innerWidth / 2
      const guideHandoffXRemain = Math.pow(1 - guideHandoffProgress, 0.55)
      const guideHandoffX = (guideStartCenter - slotCenter) * guideHandoffXRemain
      const guideHandoffY = -Math.max(112, viewportHeight * 0.18) * (1 - guideHandoffProgress)
      const guideStateScale = 0.92 + guideHandoffProgress * 0.08
      const guideFacingX = guideHandoffProgress >= 0.985 ? 1 : -1
      const firstGuideArrivalShiftX = activeStep === 0 && guideHandoffProgress >= 0.985
        ? -Math.min(132, Math.max(96, slotRect.width * 0.21))
        : 0
      section.style.setProperty('--lumi-guide-handoff-progress', guideHandoffProgress.toFixed(3))
      section.style.setProperty('--lumi-guide-handoff-presence', guideHandoffProgress.toFixed(3))
      section.style.setProperty('--lumi-guide-handoff-x', `${guideHandoffX.toFixed(2)}px`)
      section.style.setProperty('--lumi-guide-handoff-y', `${guideHandoffY.toFixed(2)}px`)
      section.style.setProperty('--lumi-guide-state-scale', guideStateScale.toFixed(3))
      section.style.setProperty('--lumi-guide-facing-x', guideFacingX.toFixed(3))
      section.style.setProperty('--lumi-guide-image-arrival-x', `${firstGuideArrivalShiftX.toFixed(2)}px`)
      section.style.setProperty('--lumi-sticky-left', `${slotRect.left}px`)
      section.style.setProperty('--lumi-sticky-width', `${slotRect.width}px`)
      section.style.setProperty('--lumi-sticky-height', `${guideRect.height}px`)

      if (firstStepRect) {
        section.style.setProperty('--lumi-pinned-demo-left', `${firstStepRect.left.toFixed(2)}px`)
        section.style.setProperty('--lumi-pinned-demo-width', `${firstStepRect.width.toFixed(2)}px`)
      }
      if (pinnedDemoRect) {
        const demoTop = Math.max(pinTop + 24, pinTop + guideRect.height / 2 - pinnedDemoRect.height / 2)
        section.style.setProperty('--lumi-pinned-demo-height', `${pinnedDemoRect.height.toFixed(2)}px`)
        section.style.setProperty('--lumi-pinned-demo-top', `${demoTop.toFixed(2)}px`)
      }

      const pinReady = layoutRect ? layoutRect.top <= pinTop : sectionRect.top <= pinTop
      if (pinReady && sectionRect.bottom > pinTop + guideRect.height) {
        section.dataset.guidePin = 'fixed'
      } else if (sectionRect.bottom <= pinTop + guideRect.height) {
        section.dataset.guidePin = 'after'
      } else {
        section.dataset.guidePin = 'flow'
      }
    }

    const scheduleScrollState = () => {
      if (frame) return
      frame = window.requestAnimationFrame(updateScrollState)
    }

    updateScrollState()
    window.addEventListener('scroll', scheduleScrollState, { passive: true })
    window.addEventListener('resize', scheduleScrollState)

    return () => {
      if (frame) window.cancelAnimationFrame(frame)
      window.removeEventListener('scroll', scheduleScrollState)
      window.removeEventListener('resize', scheduleScrollState)
    }
  }, [activeGuideImageFacingCorrection, activeStep])

  return (
    <section id="lumi-scroll-guide" className="lumiScrollGuideSection" data-active-step={activeStep} ref={sectionRef}>
      <svg
        aria-hidden="true"
        className="lumiScrollGuideRouteSvg"
        viewBox="0 0 1120 2200"
        preserveAspectRatio="none"
      >
        <defs>
          <linearGradient id="lumiScrollGuideBeamStroke" x1="0" x2="1" y1="0" y2="1">
            <stop offset="0%" stopColor="#FFFFFF" stopOpacity="0.34" />
            <stop offset="24%" stopColor="#AEEFFF" stopOpacity="0.74" />
            <stop offset="58%" stopColor="#6ED4F3" stopOpacity="0.86" />
            <stop offset="100%" stopColor="#FFFFFF" stopOpacity="0.26" />
          </linearGradient>
        </defs>
        <g className="lumiScrollGuideSegmentGroup lumiScrollGuideSegmentGroup1">
          <path className="lumiScrollGuideBeamGlow lumiScrollGuideSegmentGlow" d="M 606 0 C 742 210 902 260 760 440" pathLength="100" />
          <path className="lumiScrollGuideBeamBase lumiScrollGuideSegmentBase" d="M 606 0 C 742 210 902 260 760 440" pathLength="100" />
          <path className="lumiScrollGuideBeamCore lumiScrollGuideSegmentCore" d="M 606 0 C 742 210 902 260 760 440" pathLength="100" />
        </g>
        <g className="lumiScrollGuideSegmentGroup lumiScrollGuideSegmentGroup2">
          <path className="lumiScrollGuideBeamGlow lumiScrollGuideSegmentGlow" d="M 760 440 C 616 622 286 640 410 880" pathLength="100" />
          <path className="lumiScrollGuideBeamBase lumiScrollGuideSegmentBase" d="M 760 440 C 616 622 286 640 410 880" pathLength="100" />
          <path className="lumiScrollGuideBeamCore lumiScrollGuideSegmentCore" d="M 760 440 C 616 622 286 640 410 880" pathLength="100" />
        </g>
        <g className="lumiScrollGuideSegmentGroup lumiScrollGuideSegmentGroup3">
          <path className="lumiScrollGuideBeamGlow lumiScrollGuideSegmentGlow" d="M 410 880 C 548 1146 876 1106 704 1392" pathLength="100" />
          <path className="lumiScrollGuideBeamBase lumiScrollGuideSegmentBase" d="M 410 880 C 548 1146 876 1106 704 1392" pathLength="100" />
          <path className="lumiScrollGuideBeamCore lumiScrollGuideSegmentCore" d="M 410 880 C 548 1146 876 1106 704 1392" pathLength="100" />
        </g>
        <g className="lumiScrollGuideSegmentGroup lumiScrollGuideSegmentGroup4">
          <path className="lumiScrollGuideBeamGlow lumiScrollGuideSegmentGlow" d="M 704 1392 C 552 1646 310 1704 486 1940 C 560 2040 606 2118 610 2200" pathLength="100" />
          <path className="lumiScrollGuideBeamBase lumiScrollGuideSegmentBase" d="M 704 1392 C 552 1646 310 1704 486 1940 C 560 2040 606 2118 610 2200" pathLength="100" />
          <path className="lumiScrollGuideBeamCore lumiScrollGuideSegmentCore" d="M 704 1392 C 552 1646 310 1704 486 1940 C 560 2040 606 2118 610 2200" pathLength="100" />
        </g>
      </svg>
      <div className="lumiScrollGuideInner">
        <div className="lumiScrollGuideIntro animate-on-scroll fade-up">
          <span>{copy.eyebrow}</span>
          <h2>{copy.title}</h2>
        </div>

        <div className="lumiScrollGuideDesktopLayout" ref={desktopLayoutRef}>
          <aside className="lumiStickyGuideSlot" aria-live="polite" ref={stickySlotRef}>
            <div className="lumiStickyGuide" ref={stickyGuideRef}>
              <div className="lumiStickyGuideImageWrap">
                <img src={activeGuideImage} alt="" className="lumiStickyGuideImage" data-guide-image-step={activeStep} data-guide-image-phase={activeStep === 1 && stepTwoSuccessComplete ? 'success' : activeStep === 1 ? 'forming' : undefined} key={activeGuideImage} />
              </div>
              <div className="lumiStickyGuideBubble">
                <span>{activeGuideStep.badge}</span>
                <h3>{activeGuideStep.title}</h3>
                <p key={activeGuideDescription}>{activeGuideDescription}</p>
                <div className="lumiStickyGuideNextNote" data-visible={showStepOneNextNote ? 'true' : 'false'}>
                  {stepOneMessages.next}
                </div>
              </div>
            </div>
          </aside>

          {activePinnedStep ? (
            <div
              className="lumiPinnedDemo"
              aria-hidden="true"
              ref={pinnedDemoRef}
            >
              <StepDemo
                key={`pinned-${activeStep}`}
                step={activePinnedStep}
                index={activeStep}
                stepOnePhase={stepOnePhase}
                stepOneTypedText={stepOneTypedText}
                stepOneActionActive={stepOneActionActive}
                stepTwoPhase={stepTwoPhase}
                stepTwoActionActive={stepTwoActionActive}
                stepTwoPlanetFrameIndex={stepTwoPlanetFrameIndex}
                pinned
              />
            </div>
          ) : null}

          <div className="lumiScrollGuideSteps">
            {copy.steps.map((step, index) => (
              <article
                className="lumiScrollGuideStep animate-on-scroll fade-up"
                data-step-active={activeStep === index}
                data-step-index={index}
                style={{ '--lumi-step-image-facing-x': guideImageFacingCorrections[index] ?? 1 } as CSSProperties}
                key={step.title}
                ref={(node) => {
                  stepRefs.current[index] = node
                }}
              >
                <div className="lumiScrollGuideScene">
                  <div className="lumiScrollGuideImageWrap">
                    <img src={guideImages[index] ?? guideImages[0]} alt="" className="lumiScrollGuideImage" />
                  </div>
                  <div className="lumiScrollGuideBubble">
                    <span>{step.badge}</span>
                    <h3>{step.title}</h3>
                    <p>{step.description}</p>
                  </div>
                </div>
                <StepDemo
                  step={step}
                  index={index}
                  stepOnePhase={stepOnePhase}
                  stepOneTypedText={stepOneTypedText}
                  stepOneActionActive={stepOneActionActive}
                  stepTwoPhase={stepTwoPhase}
                  stepTwoActionActive={stepTwoActionActive}
                  stepTwoPlanetFrameIndex={stepTwoPlanetFrameIndex}
                />
              </article>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
