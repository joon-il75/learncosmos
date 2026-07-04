'use client';

import { useEffect, useRef } from 'react';

import type { LandingPageCopy } from '@/lib/i18n/pages/landing';
import LearningGuideVideoGate from '@/components/landing/LearningGuideVideoGate';

const journeySections = [
  {
    id: 'intro-video',
    label: '학습 안내',
    marker: '01',
  },
  {
    id: 'public-documents',
    label: '운영기준',
    marker: '02',
  },
  {
    id: 'ai-points-byok',
    label: 'AI포인트,BYOK 안내',
    marker: '03',
  },
  {
    id: 'platform-news',
    label: '새소식(공지사항)',
    marker: '04',
  },
];

type LandingJourneySectionsCopy = {
  introVideo: LandingPageCopy['introVideo'];
  policies: LandingPageCopy['policies'];
  aiUsage: LandingPageCopy['aiUsage'];
};

const clampProgress = (value: number) => Math.min(1, Math.max(0, value));
const stageProgress = (progress: number, from: number, to: number) => clampProgress((progress - from) / (to - from));

export default function LandingJourneySections({ copy, locale, videoGateDisabled = false }: { copy: LandingJourneySectionsCopy; locale?: string | null; videoGateDisabled?: boolean }) {
  const journeyRef = useRef<HTMLDivElement | null>(null);
  const learningGuideRef = useRef<HTMLElement | null>(null);
  const publicDocumentsRef = useRef<HTMLElement | null>(null);
  const aiPointsByokRef = useRef<HTMLElement | null>(null);
  const platformNewsRef = useRef<HTMLElement | null>(null);
  const introCopy = copy.introVideo;
  const policiesCopy = copy.policies;
  const aiUsageCopy = copy.aiUsage;
  const platformNewsCopy = locale === 'en'
    ? {
        eyebrow: 'News Board',
        title: 'News and notices will be collected here',
        speech: 'I will bring important platform updates here.',
        description: 'This board area is ready for notices, release updates, and alpha test announcements.',
        emptyTitle: 'No notices yet',
        emptyDescription: 'The list is intentionally empty for now. New posts can be connected here later.',
        ctaLabel: 'View all news',
      }
    : {
        eyebrow: '새소식 게시판',
        title: '새소식과 공지사항을 이곳에 모아둘게요',
        speech: '중요한 소식은 제가 이곳에서 알려드릴게요.',
        description: '서비스 공지, 업데이트 안내, 알파 테스트 소식을 게시판 형식으로 연결할 예정입니다.',
        emptyTitle: '아직 등록된 공지사항이 없습니다',
        emptyDescription: '게시판 리스트 영역은 비워두었습니다. 이후 공지 데이터가 준비되면 이 영역에 연결합니다.',
        ctaLabel: '전체 새소식 보기',
      };

  useEffect(() => {
    const element = learningGuideRef.current;
    if (!element) return;

    const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)');
    let frame = 0;

    const setProgress = () => {
      frame = 0;
      if (reduceMotion.matches) {
        element.style.setProperty('--learning-guide-progress', '1');
        element.style.setProperty('--learning-guide-vector-progress', '1');
        element.style.setProperty('--learning-guide-lumi-progress', '1');
        element.style.setProperty('--learning-guide-bubble-progress', '1');
        element.style.setProperty('--learning-guide-video-progress', '1');
        element.style.setProperty('--learning-guide-details-progress', '1');
        element.style.setProperty('--learning-guide-lumi-x', '0px');
        element.style.setProperty('--learning-guide-lumi-y', '0px');
        element.style.setProperty('--learning-guide-bubble-y', '0px');
        element.style.setProperty('--learning-guide-video-x', '0px');
        element.style.setProperty('--learning-guide-video-y', '0px');
        element.style.setProperty('--learning-guide-details-y', '0px');
        element.style.setProperty('--learning-guide-lumi-scale', '1');
        element.style.setProperty('--learning-guide-video-scale', '1');
        element.dataset.learningGuideActive = 'true';
        return;
      }

      const rect = element.getBoundingClientRect();
      const viewportHeight = window.innerHeight || 1;
      const isStackedJourney = window.innerWidth < 1024;
      const start = viewportHeight * (isStackedJourney ? 1.08 : 0.92);
      const distance = isStackedJourney
        ? Math.min(Math.max(viewportHeight * 0.9, rect.height * 0.34), viewportHeight * 1.18)
        : Math.min(Math.max(viewportHeight * 1.28, rect.height * 0.46), viewportHeight * 1.62);
      const progress = clampProgress((start - rect.top) / distance);
      const vectorProgress = isStackedJourney ? stageProgress(progress, 0, 0.18) : stageProgress(progress, 0, 0.28);
      const lumiProgress = isStackedJourney ? stageProgress(progress, 0.06, 0.34) : stageProgress(progress, 0.18, 0.52);
      const bubbleProgress = isStackedJourney ? stageProgress(progress, 0.16, 0.46) : stageProgress(progress, 0.38, 0.72);
      const videoProgress = isStackedJourney ? stageProgress(progress, 0.3, 0.66) : stageProgress(progress, 0.56, 0.9);
      const detailsProgress = isStackedJourney ? stageProgress(progress, 0.48, 0.82) : stageProgress(progress, 0.78, 1);

      element.style.setProperty('--learning-guide-progress', progress.toFixed(4));
      element.style.setProperty('--learning-guide-vector-progress', vectorProgress.toFixed(4));
      element.style.setProperty('--learning-guide-lumi-progress', lumiProgress.toFixed(4));
      element.style.setProperty('--learning-guide-bubble-progress', bubbleProgress.toFixed(4));
      element.style.setProperty('--learning-guide-video-progress', videoProgress.toFixed(4));
      element.style.setProperty('--learning-guide-details-progress', detailsProgress.toFixed(4));
      element.style.setProperty('--learning-guide-lumi-x', `${(-96 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty('--learning-guide-lumi-y', `${(44 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty('--learning-guide-bubble-y', `${(-48 * (1 - bubbleProgress)).toFixed(2)}px`);
      element.style.setProperty('--learning-guide-video-x', `${(128 * (1 - videoProgress)).toFixed(2)}px`);
      element.style.setProperty('--learning-guide-video-y', `${(28 * (1 - videoProgress)).toFixed(2)}px`);
      element.style.setProperty('--learning-guide-details-y', `${(18 * (1 - detailsProgress)).toFixed(2)}px`);
      element.style.setProperty('--learning-guide-lumi-scale', (0.94 + lumiProgress * 0.06).toFixed(4));
      element.style.setProperty('--learning-guide-video-scale', (0.965 + videoProgress * 0.035).toFixed(4));
      element.dataset.learningGuideActive = progress > 0.02 ? 'true' : 'false';
    };

    const requestUpdate = () => {
      if (frame) return;
      frame = window.requestAnimationFrame(setProgress);
    };

    setProgress();
    window.addEventListener('scroll', requestUpdate, { passive: true });
    window.addEventListener('resize', requestUpdate);
    reduceMotion.addEventListener('change', requestUpdate);

    return () => {
      if (frame) window.cancelAnimationFrame(frame);
      window.removeEventListener('scroll', requestUpdate);
      window.removeEventListener('resize', requestUpdate);
      reduceMotion.removeEventListener('change', requestUpdate);
    };
  }, []);

  useEffect(() => {
    const element = publicDocumentsRef.current;
    if (!element) return;

    const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)');
    let frame = 0;

    const setProgress = () => {
      frame = 0;
      if (reduceMotion.matches) {
        element.style.setProperty('--public-docs-progress', '1');
        element.style.setProperty('--public-docs-marker-progress', '1');
        element.style.setProperty('--public-docs-lumi-progress', '1');
        element.style.setProperty('--public-docs-bubble-progress', '1');
        element.style.setProperty('--public-docs-copy-progress', '1');
        element.style.setProperty('--public-docs-cards-progress', '1');
        element.style.setProperty('--public-docs-lumi-x', '0px');
        element.style.setProperty('--public-docs-lumi-y', '0px');
        element.style.setProperty('--public-docs-bubble-y', '0px');
        element.style.setProperty('--public-docs-copy-y', '0px');
        element.style.setProperty('--public-docs-card-y', '0px');
        element.style.setProperty('--public-docs-lumi-scale', '1');
        element.dataset.publicDocumentsActive = 'true';
        return;
      }

      const rect = element.getBoundingClientRect();
      const viewportHeight = window.innerHeight || 1;
      const isStackedJourney = window.innerWidth < 1024;
      const start = viewportHeight * (isStackedJourney ? 1.06 : 0.9);
      const distance = isStackedJourney
        ? Math.min(Math.max(viewportHeight * 0.86, rect.height * 0.32), viewportHeight * 1.14)
        : Math.min(Math.max(viewportHeight * 1.18, rect.height * 0.42), viewportHeight * 1.56);
      const progress = clampProgress((start - rect.top) / distance);
      const markerProgress = isStackedJourney ? stageProgress(progress, 0, 0.12) : stageProgress(progress, 0, 0.18);
      const lumiProgress = isStackedJourney ? stageProgress(progress, 0.04, 0.32) : stageProgress(progress, 0.08, 0.46);
      const bubbleProgress = isStackedJourney ? stageProgress(progress, 0.14, 0.44) : stageProgress(progress, 0.26, 0.62);
      const copyProgress = isStackedJourney ? stageProgress(progress, 0.28, 0.58) : stageProgress(progress, 0.42, 0.76);
      const cardsProgress = isStackedJourney ? stageProgress(progress, 0.38, 0.72) : stageProgress(progress, 0.52, 1);

      element.style.setProperty('--public-docs-progress', progress.toFixed(4));
      element.style.setProperty('--public-docs-marker-progress', markerProgress.toFixed(4));
      element.style.setProperty('--public-docs-lumi-progress', lumiProgress.toFixed(4));
      element.style.setProperty('--public-docs-bubble-progress', bubbleProgress.toFixed(4));
      element.style.setProperty('--public-docs-copy-progress', copyProgress.toFixed(4));
      element.style.setProperty('--public-docs-cards-progress', cardsProgress.toFixed(4));
      element.style.setProperty('--public-docs-lumi-x', `${(-128 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty('--public-docs-lumi-y', `${(68 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty('--public-docs-bubble-y', `${(-64 * (1 - bubbleProgress)).toFixed(2)}px`);
      element.style.setProperty('--public-docs-copy-y', `${(36 * (1 - copyProgress)).toFixed(2)}px`);
      element.style.setProperty('--public-docs-card-y', `${(54 * (1 - cardsProgress)).toFixed(2)}px`);
      element.style.setProperty('--public-docs-lumi-scale', (0.88 + lumiProgress * 0.12).toFixed(4));
      element.dataset.publicDocumentsActive = progress > 0.02 ? 'true' : 'false';
    };

    const requestUpdate = () => {
      if (frame) return;
      frame = window.requestAnimationFrame(setProgress);
    };

    setProgress();
    window.addEventListener('scroll', requestUpdate, { passive: true });
    window.addEventListener('resize', requestUpdate);
    reduceMotion.addEventListener('change', requestUpdate);

    return () => {
      if (frame) window.cancelAnimationFrame(frame);
      window.removeEventListener('scroll', requestUpdate);
      window.removeEventListener('resize', requestUpdate);
      reduceMotion.removeEventListener('change', requestUpdate);
    };
  }, []);

  useEffect(() => {
    const element = aiPointsByokRef.current;
    if (!element) return;

    const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)');
    let frame = 0;

    const setProgress = () => {
      frame = 0;
      if (reduceMotion.matches) {
        element.style.setProperty('--ai-points-progress', '1');
        element.style.setProperty('--ai-points-lumi-progress', '1');
        element.style.setProperty('--ai-points-bubble-progress', '1');
        element.style.setProperty('--ai-points-copy-progress', '1');
        element.style.setProperty('--ai-points-cards-progress', '1');
        element.style.setProperty('--ai-points-lumi-x', '0px');
        element.style.setProperty('--ai-points-lumi-y', '0px');
        element.style.setProperty('--ai-points-bubble-y', '0px');
        element.style.setProperty('--ai-points-copy-y', '0px');
        element.style.setProperty('--ai-points-card-y', '0px');
        element.style.setProperty('--ai-points-lumi-scale', '1');
        element.dataset.aiPointsActive = 'true';
        return;
      }

      const rect = element.getBoundingClientRect();
      const viewportHeight = window.innerHeight || 1;
      const isStackedJourney = window.innerWidth < 1024;
      const start = viewportHeight * (isStackedJourney ? 1.06 : 0.92);
      const distance = isStackedJourney
        ? Math.min(Math.max(viewportHeight * 0.88, rect.height * 0.33), viewportHeight * 1.16)
        : Math.min(Math.max(viewportHeight * 1.2, rect.height * 0.44), viewportHeight * 1.62);
      const progress = clampProgress((start - rect.top) / distance);
      const lumiProgress = isStackedJourney ? stageProgress(progress, 0.04, 0.32) : stageProgress(progress, 0.06, 0.42);
      const bubbleProgress = isStackedJourney ? stageProgress(progress, 0.14, 0.44) : stageProgress(progress, 0.24, 0.58);
      const copyProgress = isStackedJourney ? stageProgress(progress, 0.28, 0.58) : stageProgress(progress, 0.4, 0.74);
      const cardsProgress = isStackedJourney ? stageProgress(progress, 0.4, 0.74) : stageProgress(progress, 0.56, 1);

      element.style.setProperty('--ai-points-progress', progress.toFixed(4));
      element.style.setProperty('--ai-points-lumi-progress', lumiProgress.toFixed(4));
      element.style.setProperty('--ai-points-bubble-progress', bubbleProgress.toFixed(4));
      element.style.setProperty('--ai-points-copy-progress', copyProgress.toFixed(4));
      element.style.setProperty('--ai-points-cards-progress', cardsProgress.toFixed(4));
      element.style.setProperty('--ai-points-lumi-x', `${(-118 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty('--ai-points-lumi-y', `${(72 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty('--ai-points-bubble-y', `${(-62 * (1 - bubbleProgress)).toFixed(2)}px`);
      element.style.setProperty('--ai-points-copy-y', `${(40 * (1 - copyProgress)).toFixed(2)}px`);
      element.style.setProperty('--ai-points-card-y', `${(58 * (1 - cardsProgress)).toFixed(2)}px`);
      element.style.setProperty('--ai-points-lumi-scale', (0.88 + lumiProgress * 0.12).toFixed(4));
      element.dataset.aiPointsActive = progress > 0.02 ? 'true' : 'false';
    };

    const requestUpdate = () => {
      if (frame) return;
      frame = window.requestAnimationFrame(setProgress);
    };

    setProgress();
    window.addEventListener('scroll', requestUpdate, { passive: true });
    window.addEventListener('resize', requestUpdate);
    reduceMotion.addEventListener('change', requestUpdate);

    return () => {
      if (frame) window.cancelAnimationFrame(frame);
      window.removeEventListener('scroll', requestUpdate);
      window.removeEventListener('resize', requestUpdate);
      reduceMotion.removeEventListener('change', requestUpdate);
    };
  }, []);

  useEffect(() => {
    const element = platformNewsRef.current;
    if (!element) return;

    const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)');
    let frame = 0;

    const setProgress = () => {
      frame = 0;
      if (reduceMotion.matches) {
        element.style.setProperty('--platform-news-progress', '1');
        element.style.setProperty('--platform-news-lumi-progress', '1');
        element.style.setProperty('--platform-news-bubble-progress', '1');
        element.style.setProperty('--platform-news-copy-progress', '1');
        element.style.setProperty('--platform-news-board-progress', '1');
        element.style.setProperty('--platform-news-lumi-x', '0px');
        element.style.setProperty('--platform-news-lumi-y', '0px');
        element.style.setProperty('--platform-news-bubble-y', '0px');
        element.style.setProperty('--platform-news-copy-y', '0px');
        element.style.setProperty('--platform-news-board-y', '0px');
        element.style.setProperty('--platform-news-lumi-scale', '1');
        element.dataset.platformNewsActive = 'true';
        return;
      }

      const rect = element.getBoundingClientRect();
      const viewportHeight = window.innerHeight || 1;
      const isStackedJourney = window.innerWidth < 1024;
      const start = viewportHeight * (isStackedJourney ? 1.06 : 0.92);
      const distance = isStackedJourney
        ? Math.min(Math.max(viewportHeight * 0.86, rect.height * 0.32), viewportHeight * 1.14)
        : Math.min(Math.max(viewportHeight * 1.18, rect.height * 0.42), viewportHeight * 1.58);
      const progress = clampProgress((start - rect.top) / distance);
      const lumiProgress = isStackedJourney ? stageProgress(progress, 0.04, 0.32) : stageProgress(progress, 0.06, 0.42);
      const bubbleProgress = isStackedJourney ? stageProgress(progress, 0.12, 0.42) : stageProgress(progress, 0.22, 0.58);
      const copyProgress = isStackedJourney ? stageProgress(progress, 0.26, 0.56) : stageProgress(progress, 0.38, 0.72);
      const boardProgress = isStackedJourney ? stageProgress(progress, 0.38, 0.72) : stageProgress(progress, 0.52, 1);

      element.style.setProperty('--platform-news-progress', progress.toFixed(4));
      element.style.setProperty('--platform-news-lumi-progress', lumiProgress.toFixed(4));
      element.style.setProperty('--platform-news-bubble-progress', bubbleProgress.toFixed(4));
      element.style.setProperty('--platform-news-copy-progress', copyProgress.toFixed(4));
      element.style.setProperty('--platform-news-board-progress', boardProgress.toFixed(4));
      element.style.setProperty('--platform-news-lumi-x', `${(-116 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty('--platform-news-lumi-y', `${(70 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty('--platform-news-bubble-y', `${(-60 * (1 - bubbleProgress)).toFixed(2)}px`);
      element.style.setProperty('--platform-news-copy-y', `${(36 * (1 - copyProgress)).toFixed(2)}px`);
      element.style.setProperty('--platform-news-board-y', `${(58 * (1 - boardProgress)).toFixed(2)}px`);
      element.style.setProperty('--platform-news-lumi-scale', (0.88 + lumiProgress * 0.12).toFixed(4));
      element.dataset.platformNewsActive = progress > 0.02 ? 'true' : 'false';
    };

    const requestUpdate = () => {
      if (frame) return;
      frame = window.requestAnimationFrame(setProgress);
    };

    setProgress();
    window.addEventListener('scroll', requestUpdate, { passive: true });
    window.addEventListener('resize', requestUpdate);
    reduceMotion.addEventListener('change', requestUpdate);

    return () => {
      if (frame) window.cancelAnimationFrame(frame);
      window.removeEventListener('scroll', requestUpdate);
      window.removeEventListener('resize', requestUpdate);
      reduceMotion.removeEventListener('change', requestUpdate);
    };
  }, []);



  useEffect(() => {
    const runway = journeyRef.current;
    if (!runway) return;

    const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)');
    const panels = [learningGuideRef, publicDocumentsRef, aiPointsByokRef, platformNewsRef];
    let frame = 0;

    const setPanelVars = (element: HTMLElement, prefix: string, localProgress: number, activeOpacity: number) => {
      const p = clampProgress(localProgress);
      element.style.setProperty('--journey-panel-opacity', activeOpacity.toFixed(4));
      element.dataset.journeyPanelActive = activeOpacity > 0.5 ? 'true' : 'false';

      if (prefix === 'learning-guide') {
        const vectorProgress = stageProgress(p, 0, 0.24);
        const lumiProgress = stageProgress(p, 0.08, 0.4);
        const bubbleProgress = stageProgress(p, 0.24, 0.58);
        const videoProgress = stageProgress(p, 0.42, 0.78);
        const detailsProgress = stageProgress(p, 0.68, 1);
        element.style.setProperty('--learning-guide-progress', p.toFixed(4));
        element.style.setProperty('--learning-guide-vector-progress', vectorProgress.toFixed(4));
        element.style.setProperty('--learning-guide-lumi-progress', lumiProgress.toFixed(4));
        element.style.setProperty('--learning-guide-bubble-progress', bubbleProgress.toFixed(4));
        element.style.setProperty('--learning-guide-video-progress', videoProgress.toFixed(4));
        element.style.setProperty('--learning-guide-details-progress', detailsProgress.toFixed(4));
        element.style.setProperty('--learning-guide-lumi-x', `${(-84 * (1 - lumiProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-lumi-y', `${(38 * (1 - lumiProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-bubble-y', `${(-42 * (1 - bubbleProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-video-x', `${(118 * (1 - videoProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-video-y', `${(26 * (1 - videoProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-details-y', `${(16 * (1 - detailsProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-lumi-scale', (0.94 + lumiProgress * 0.06).toFixed(4));
        element.style.setProperty('--learning-guide-video-scale', (0.965 + videoProgress * 0.035).toFixed(4));
        return;
      }

      const names = prefix === 'public-docs'
        ? ['--public-docs-progress', '--public-docs-lumi-progress', '--public-docs-bubble-progress', '--public-docs-copy-progress', '--public-docs-cards-progress', '--public-docs-lumi-x', '--public-docs-lumi-y', '--public-docs-bubble-y', '--public-docs-copy-y', '--public-docs-card-y', '--public-docs-lumi-scale']
        : prefix === 'ai-points'
          ? ['--ai-points-progress', '--ai-points-lumi-progress', '--ai-points-bubble-progress', '--ai-points-copy-progress', '--ai-points-cards-progress', '--ai-points-lumi-x', '--ai-points-lumi-y', '--ai-points-bubble-y', '--ai-points-copy-y', '--ai-points-card-y', '--ai-points-lumi-scale']
          : ['--platform-news-progress', '--platform-news-lumi-progress', '--platform-news-bubble-progress', '--platform-news-copy-progress', '--platform-news-board-progress', '--platform-news-lumi-x', '--platform-news-lumi-y', '--platform-news-bubble-y', '--platform-news-copy-y', '--platform-news-board-y', '--platform-news-lumi-scale'];
      const lumiProgress = stageProgress(p, 0.05, 0.38);
      const bubbleProgress = stageProgress(p, 0.2, 0.54);
      const copyProgress = stageProgress(p, 0.36, 0.7);
      const objectProgress = stageProgress(p, 0.52, 1);
      element.style.setProperty(names[0], p.toFixed(4));
      element.style.setProperty(names[1], lumiProgress.toFixed(4));
      element.style.setProperty(names[2], bubbleProgress.toFixed(4));
      element.style.setProperty(names[3], copyProgress.toFixed(4));
      element.style.setProperty(names[4], objectProgress.toFixed(4));
      element.style.setProperty(names[5], `${(-112 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty(names[6], `${(68 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty(names[7], `${(-58 * (1 - bubbleProgress)).toFixed(2)}px`);
      element.style.setProperty(names[8], `${(36 * (1 - copyProgress)).toFixed(2)}px`);
      element.style.setProperty(names[9], `${(54 * (1 - objectProgress)).toFixed(2)}px`);
      element.style.setProperty(names[10], (0.9 + lumiProgress * 0.1).toFixed(4));
    };

    const setProgress = () => {
      frame = 0;
      const panelElements = panels.map((ref) => ref.current).filter(Boolean) as HTMLElement[];
      if (!panelElements.length) return;

      if (reduceMotion.matches || window.innerWidth < 1024) {
        panelElements.forEach((panel) => {
          panel.style.setProperty('--journey-panel-opacity', '1');
          panel.dataset.journeyPanelActive = 'true';
        });
        return;
      }

      const rect = runway.getBoundingClientRect();
      const viewportHeight = window.innerHeight || 1;
      const travel = Math.max(1, rect.height - viewportHeight);
      const progress = clampProgress((-rect.top) / travel);
      const windows = [
        [0, 0.28],
        [0.28, 0.56],
        [0.56, 0.84],
        [0.84, 1],
      ] as const;
      const assemblyWindows = [
        [0, 0.12],
        [0.28, 0.4],
        [0.56, 0.68],
        [0.84, 0.96],
      ] as const;
      const prefixes = ['learning-guide', 'public-docs', 'ai-points', 'platform-news'] as const;

      panelElements.forEach((panel, index) => {
        const [from, to] = windows[index];
        const [assemblyFrom, assemblyTo] = assemblyWindows[index];
        const local = clampProgress((progress - assemblyFrom) / Math.max(0.001, assemblyTo - assemblyFrom));
        const enter = stageProgress(progress, from, from + 0.055);
        const exit = index === windows.length - 1 ? 1 : 1 - stageProgress(progress, to - 0.055, to);
        const opacity = clampProgress(Math.min(enter, exit));
        setPanelVars(panel, prefixes[index], local, opacity);
      });
    };

    const requestUpdate = () => {
      if (frame) return;
      frame = window.requestAnimationFrame(setProgress);
    };

    setProgress();
    window.addEventListener('scroll', requestUpdate, { passive: true });
    window.addEventListener('resize', requestUpdate);
    reduceMotion.addEventListener('change', requestUpdate);

    return () => {
      if (frame) window.cancelAnimationFrame(frame);
      window.removeEventListener('scroll', requestUpdate);
      window.removeEventListener('resize', requestUpdate);
      reduceMotion.removeEventListener('change', requestUpdate);
    };
  }, []);

  return (
    <div ref={journeyRef} className="landingJourneySections landingStickyJourneyRunway" aria-label="LearnCosmos learning journey sections">
      <div className="landingStickyJourneyStage">
      {journeySections.map((section, index) => {
        if (index === 0) {


  useEffect(() => {
    const runway = journeyRef.current;
    if (!runway) return;

    const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)');
    const panels = [learningGuideRef, publicDocumentsRef, aiPointsByokRef, platformNewsRef];
    let frame = 0;

    const setPanelVars = (element: HTMLElement, prefix: string, localProgress: number, activeOpacity: number) => {
      const p = clampProgress(localProgress);
      element.style.setProperty('--journey-panel-opacity', activeOpacity.toFixed(4));
      element.dataset.journeyPanelActive = activeOpacity > 0.5 ? 'true' : 'false';

      if (prefix === 'learning-guide') {
        const vectorProgress = stageProgress(p, 0, 0.24);
        const lumiProgress = stageProgress(p, 0.08, 0.4);
        const bubbleProgress = stageProgress(p, 0.24, 0.58);
        const videoProgress = stageProgress(p, 0.42, 0.78);
        const detailsProgress = stageProgress(p, 0.68, 1);
        element.style.setProperty('--learning-guide-progress', p.toFixed(4));
        element.style.setProperty('--learning-guide-vector-progress', vectorProgress.toFixed(4));
        element.style.setProperty('--learning-guide-lumi-progress', lumiProgress.toFixed(4));
        element.style.setProperty('--learning-guide-bubble-progress', bubbleProgress.toFixed(4));
        element.style.setProperty('--learning-guide-video-progress', videoProgress.toFixed(4));
        element.style.setProperty('--learning-guide-details-progress', detailsProgress.toFixed(4));
        element.style.setProperty('--learning-guide-lumi-x', `${(-84 * (1 - lumiProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-lumi-y', `${(38 * (1 - lumiProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-bubble-y', `${(-42 * (1 - bubbleProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-video-x', `${(118 * (1 - videoProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-video-y', `${(26 * (1 - videoProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-details-y', `${(16 * (1 - detailsProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-lumi-scale', (0.94 + lumiProgress * 0.06).toFixed(4));
        element.style.setProperty('--learning-guide-video-scale', (0.965 + videoProgress * 0.035).toFixed(4));
        return;
      }

      const names = prefix === 'public-docs'
        ? ['--public-docs-progress', '--public-docs-lumi-progress', '--public-docs-bubble-progress', '--public-docs-copy-progress', '--public-docs-cards-progress', '--public-docs-lumi-x', '--public-docs-lumi-y', '--public-docs-bubble-y', '--public-docs-copy-y', '--public-docs-card-y', '--public-docs-lumi-scale']
        : prefix === 'ai-points'
          ? ['--ai-points-progress', '--ai-points-lumi-progress', '--ai-points-bubble-progress', '--ai-points-copy-progress', '--ai-points-cards-progress', '--ai-points-lumi-x', '--ai-points-lumi-y', '--ai-points-bubble-y', '--ai-points-copy-y', '--ai-points-card-y', '--ai-points-lumi-scale']
          : ['--platform-news-progress', '--platform-news-lumi-progress', '--platform-news-bubble-progress', '--platform-news-copy-progress', '--platform-news-board-progress', '--platform-news-lumi-x', '--platform-news-lumi-y', '--platform-news-bubble-y', '--platform-news-copy-y', '--platform-news-board-y', '--platform-news-lumi-scale'];
      const lumiProgress = stageProgress(p, 0.05, 0.38);
      const bubbleProgress = stageProgress(p, 0.2, 0.54);
      const copyProgress = stageProgress(p, 0.36, 0.7);
      const objectProgress = stageProgress(p, 0.52, 1);
      element.style.setProperty(names[0], p.toFixed(4));
      element.style.setProperty(names[1], lumiProgress.toFixed(4));
      element.style.setProperty(names[2], bubbleProgress.toFixed(4));
      element.style.setProperty(names[3], copyProgress.toFixed(4));
      element.style.setProperty(names[4], objectProgress.toFixed(4));
      element.style.setProperty(names[5], `${(-112 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty(names[6], `${(68 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty(names[7], `${(-58 * (1 - bubbleProgress)).toFixed(2)}px`);
      element.style.setProperty(names[8], `${(36 * (1 - copyProgress)).toFixed(2)}px`);
      element.style.setProperty(names[9], `${(54 * (1 - objectProgress)).toFixed(2)}px`);
      element.style.setProperty(names[10], (0.9 + lumiProgress * 0.1).toFixed(4));
    };

    const setProgress = () => {
      frame = 0;
      const panelElements = panels.map((ref) => ref.current).filter(Boolean) as HTMLElement[];
      if (!panelElements.length) return;

      if (reduceMotion.matches || window.innerWidth < 1024) {
        panelElements.forEach((panel) => {
          panel.style.setProperty('--journey-panel-opacity', '1');
          panel.dataset.journeyPanelActive = 'true';
        });
        return;
      }

      const rect = runway.getBoundingClientRect();
      const viewportHeight = window.innerHeight || 1;
      const travel = Math.max(1, rect.height - viewportHeight);
      const progress = clampProgress((-rect.top) / travel);
      const windows = [
        [0, 0.28],
        [0.28, 0.56],
        [0.56, 0.84],
        [0.84, 1],
      ] as const;
      const assemblyWindows = [
        [0, 0.12],
        [0.28, 0.4],
        [0.56, 0.68],
        [0.84, 0.96],
      ] as const;
      const prefixes = ['learning-guide', 'public-docs', 'ai-points', 'platform-news'] as const;

      panelElements.forEach((panel, index) => {
        const [from, to] = windows[index];
        const [assemblyFrom, assemblyTo] = assemblyWindows[index];
        const local = clampProgress((progress - assemblyFrom) / Math.max(0.001, assemblyTo - assemblyFrom));
        const enter = stageProgress(progress, from, from + 0.055);
        const exit = index === windows.length - 1 ? 1 : 1 - stageProgress(progress, to - 0.055, to);
        const opacity = clampProgress(Math.min(enter, exit));
        setPanelVars(panel, prefixes[index], local, opacity);
      });
    };

    const requestUpdate = () => {
      if (frame) return;
      frame = window.requestAnimationFrame(setProgress);
    };

    setProgress();
    window.addEventListener('scroll', requestUpdate, { passive: true });
    window.addEventListener('resize', requestUpdate);
    reduceMotion.addEventListener('change', requestUpdate);

    return () => {
      if (frame) window.cancelAnimationFrame(frame);
      window.removeEventListener('scroll', requestUpdate);
      window.removeEventListener('resize', requestUpdate);
      reduceMotion.removeEventListener('change', requestUpdate);
    };
  }, []);

  return (
            <section
              key={section.id}
              id={section.id}
              ref={learningGuideRef}
              data-ready-progress="0.18"
              data-mobile-ready-progress="1"
              className="landingJourneySection landingStickyJourneyPanel landingLearningGuideSection"
              aria-label={section.label}
            >
              <span id="intro-video-ready" className="landingJourneyReadyAnchor landingJourneyReadyAnchorLearningGuide" aria-hidden="true" />
              <div className="landingLearningGuideShell">
                <div className="landingLearningGuideLumiPane">
                  <div className="landingJourneySectionMarker landingLearningGuideMarker" aria-hidden="true">
                    <span>{section.marker}</span>
                    <strong>{introCopy.title}</strong>
                  </div>
                  <div className="landingLearningGuideLumiScene">
                    <img
                      src="/images/lumi_journey_learning_guide_pose_alpha_v1.png"
                      alt=""
                      className="landingLearningGuideLumiImage"
                      aria-hidden="true"
                    />
                    <div className="landingLearningGuideBubble" aria-label={introCopy.title}>
                      {introCopy.lumiSpeechLines.map((line) => (
                        <p key={line}>{line}</p>
                      ))}
                    </div>
                  </div>
                </div>

                <div className="landingLearningGuideVideoPane">
                  <LearningGuideVideoGate copy={introCopy} locale={locale} disabled={videoGateDisabled} />
                  <div className="landingLearningGuideAlphaCopy">
                    {introCopy.descriptionLines.map((line) => (
                      <p key={line}>{line}</p>
                    ))}
                  </div>
                  <a className="landingLearningGuideAlphaCta" href="/alpha">
                    {introCopy.alphaCta}
                  </a>
                </div>
              </div>
            </section>
          );
        }

        if (index === 1) {


  useEffect(() => {
    const runway = journeyRef.current;
    if (!runway) return;

    const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)');
    const panels = [learningGuideRef, publicDocumentsRef, aiPointsByokRef, platformNewsRef];
    let frame = 0;

    const setPanelVars = (element: HTMLElement, prefix: string, localProgress: number, activeOpacity: number) => {
      const p = clampProgress(localProgress);
      element.style.setProperty('--journey-panel-opacity', activeOpacity.toFixed(4));
      element.dataset.journeyPanelActive = activeOpacity > 0.5 ? 'true' : 'false';

      if (prefix === 'learning-guide') {
        const vectorProgress = stageProgress(p, 0, 0.24);
        const lumiProgress = stageProgress(p, 0.08, 0.4);
        const bubbleProgress = stageProgress(p, 0.24, 0.58);
        const videoProgress = stageProgress(p, 0.42, 0.78);
        const detailsProgress = stageProgress(p, 0.68, 1);
        element.style.setProperty('--learning-guide-progress', p.toFixed(4));
        element.style.setProperty('--learning-guide-vector-progress', vectorProgress.toFixed(4));
        element.style.setProperty('--learning-guide-lumi-progress', lumiProgress.toFixed(4));
        element.style.setProperty('--learning-guide-bubble-progress', bubbleProgress.toFixed(4));
        element.style.setProperty('--learning-guide-video-progress', videoProgress.toFixed(4));
        element.style.setProperty('--learning-guide-details-progress', detailsProgress.toFixed(4));
        element.style.setProperty('--learning-guide-lumi-x', `${(-84 * (1 - lumiProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-lumi-y', `${(38 * (1 - lumiProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-bubble-y', `${(-42 * (1 - bubbleProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-video-x', `${(118 * (1 - videoProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-video-y', `${(26 * (1 - videoProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-details-y', `${(16 * (1 - detailsProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-lumi-scale', (0.94 + lumiProgress * 0.06).toFixed(4));
        element.style.setProperty('--learning-guide-video-scale', (0.965 + videoProgress * 0.035).toFixed(4));
        return;
      }

      const names = prefix === 'public-docs'
        ? ['--public-docs-progress', '--public-docs-lumi-progress', '--public-docs-bubble-progress', '--public-docs-copy-progress', '--public-docs-cards-progress', '--public-docs-lumi-x', '--public-docs-lumi-y', '--public-docs-bubble-y', '--public-docs-copy-y', '--public-docs-card-y', '--public-docs-lumi-scale']
        : prefix === 'ai-points'
          ? ['--ai-points-progress', '--ai-points-lumi-progress', '--ai-points-bubble-progress', '--ai-points-copy-progress', '--ai-points-cards-progress', '--ai-points-lumi-x', '--ai-points-lumi-y', '--ai-points-bubble-y', '--ai-points-copy-y', '--ai-points-card-y', '--ai-points-lumi-scale']
          : ['--platform-news-progress', '--platform-news-lumi-progress', '--platform-news-bubble-progress', '--platform-news-copy-progress', '--platform-news-board-progress', '--platform-news-lumi-x', '--platform-news-lumi-y', '--platform-news-bubble-y', '--platform-news-copy-y', '--platform-news-board-y', '--platform-news-lumi-scale'];
      const lumiProgress = stageProgress(p, 0.05, 0.38);
      const bubbleProgress = stageProgress(p, 0.2, 0.54);
      const copyProgress = stageProgress(p, 0.36, 0.7);
      const objectProgress = stageProgress(p, 0.52, 1);
      element.style.setProperty(names[0], p.toFixed(4));
      element.style.setProperty(names[1], lumiProgress.toFixed(4));
      element.style.setProperty(names[2], bubbleProgress.toFixed(4));
      element.style.setProperty(names[3], copyProgress.toFixed(4));
      element.style.setProperty(names[4], objectProgress.toFixed(4));
      element.style.setProperty(names[5], `${(-112 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty(names[6], `${(68 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty(names[7], `${(-58 * (1 - bubbleProgress)).toFixed(2)}px`);
      element.style.setProperty(names[8], `${(36 * (1 - copyProgress)).toFixed(2)}px`);
      element.style.setProperty(names[9], `${(54 * (1 - objectProgress)).toFixed(2)}px`);
      element.style.setProperty(names[10], (0.9 + lumiProgress * 0.1).toFixed(4));
    };

    const setProgress = () => {
      frame = 0;
      const panelElements = panels.map((ref) => ref.current).filter(Boolean) as HTMLElement[];
      if (!panelElements.length) return;

      if (reduceMotion.matches || window.innerWidth < 1024) {
        panelElements.forEach((panel) => {
          panel.style.setProperty('--journey-panel-opacity', '1');
          panel.dataset.journeyPanelActive = 'true';
        });
        return;
      }

      const rect = runway.getBoundingClientRect();
      const viewportHeight = window.innerHeight || 1;
      const travel = Math.max(1, rect.height - viewportHeight);
      const progress = clampProgress((-rect.top) / travel);
      const windows = [
        [0, 0.28],
        [0.28, 0.56],
        [0.56, 0.84],
        [0.84, 1],
      ] as const;
      const assemblyWindows = [
        [0, 0.12],
        [0.28, 0.4],
        [0.56, 0.68],
        [0.84, 0.96],
      ] as const;
      const prefixes = ['learning-guide', 'public-docs', 'ai-points', 'platform-news'] as const;

      panelElements.forEach((panel, index) => {
        const [from, to] = windows[index];
        const [assemblyFrom, assemblyTo] = assemblyWindows[index];
        const local = clampProgress((progress - assemblyFrom) / Math.max(0.001, assemblyTo - assemblyFrom));
        const enter = stageProgress(progress, from, from + 0.055);
        const exit = index === windows.length - 1 ? 1 : 1 - stageProgress(progress, to - 0.055, to);
        const opacity = clampProgress(Math.min(enter, exit));
        setPanelVars(panel, prefixes[index], local, opacity);
      });
    };

    const requestUpdate = () => {
      if (frame) return;
      frame = window.requestAnimationFrame(setProgress);
    };

    setProgress();
    window.addEventListener('scroll', requestUpdate, { passive: true });
    window.addEventListener('resize', requestUpdate);
    reduceMotion.addEventListener('change', requestUpdate);

    return () => {
      if (frame) window.cancelAnimationFrame(frame);
      window.removeEventListener('scroll', requestUpdate);
      window.removeEventListener('resize', requestUpdate);
      reduceMotion.removeEventListener('change', requestUpdate);
    };
  }, []);

  return (
            <section
              key={section.id}
              id={section.id}
              ref={publicDocumentsRef}
              data-ready-progress="0.46"
              data-mobile-ready-progress="1"
              className="landingJourneySection landingStickyJourneyPanel landingPublicDocumentsSection"
              aria-label={section.label}
            >
              <span id="public-documents-ready" className="landingJourneyReadyAnchor landingJourneyReadyAnchorPublicDocuments" aria-hidden="true" />
              <svg
                className="landingPublicDocumentsVector"
                viewBox="0 0 980 620"
                preserveAspectRatio="none"
                aria-hidden="true"
              >
                <defs>
                  <linearGradient id="landingPublicDocumentsPath" x1="120" y1="20" x2="860" y2="590" gradientUnits="userSpaceOnUse">
                    <stop offset="0%" stopColor="#9FE8FF" stopOpacity="0.08" />
                    <stop offset="45%" stopColor="#6AD2C1" stopOpacity="0.42" />
                    <stop offset="100%" stopColor="#88DAFF" stopOpacity="0.18" />
                  </linearGradient>
                  <filter id="landingPublicDocumentsGlow" x="-20%" y="-20%" width="140%" height="140%">
                    <feGaussianBlur stdDeviation="9" result="blur" />
                    <feMerge>
                      <feMergeNode in="blur" />
                      <feMergeNode in="SourceGraphic" />
                    </feMerge>
                  </filter>
                </defs>
                <path
                  className="landingPublicDocumentsVectorGlow"
                  d="M 94 536 C 218 408 180 266 340 230 C 492 196 556 336 692 286 C 808 244 810 118 920 72"
                />
                <path
                  className="landingPublicDocumentsVectorCore"
                  d="M 94 536 C 218 408 180 266 340 230 C 492 196 556 336 692 286 C 808 244 810 118 920 72"
                />
                <path
                  className="landingPublicDocumentsVectorFine"
                  d="M 176 612 C 292 482 258 360 392 318 C 520 278 626 394 770 322"
                />
              </svg>

              <div className="landingPublicDocumentsShell">
                <div className="landingJourneySectionMarker landingPublicDocumentsMarker" aria-hidden="true">
                  <span>{section.marker}</span>
                  <strong>{section.label}</strong>
                </div>

                <div className="landingPublicDocumentsLumiPane">
                  <img
                    src="/images/lumi_journey_operation_standards_pose_alpha_v1.png"
                    alt=""
                    className="landingPublicDocumentsLumiImage"
                    aria-hidden="true"
                  />
                  <div className="landingPublicDocumentsSpeech">
                    <h2>{policiesCopy.title}</h2>
                  </div>
                </div>

                <div className="landingPublicDocumentsCopy">
                  <p>{policiesCopy.eyebrow}</p>
                  <div>
                    {policiesCopy.descriptionLines.map((line) => (
                      <span key={line}>{line}</span>
                    ))}
                  </div>
                </div>

                <div className="landingPublicDocumentsGrid">
                  {policiesCopy.docs.map((doc) => (
                    <a key={doc.href} className="landingPublicDocumentCard" href={doc.href}>
                      <span>{doc.badge}</span>
                      <strong>{doc.title}</strong>
                      <p>{doc.description}</p>
                      <em>{policiesCopy.linkLabel}</em>
                    </a>
                  ))}
                </div>
              </div>
            </section>
          );
        }

        if (index === 2) {


  useEffect(() => {
    const runway = journeyRef.current;
    if (!runway) return;

    const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)');
    const panels = [learningGuideRef, publicDocumentsRef, aiPointsByokRef, platformNewsRef];
    let frame = 0;

    const setPanelVars = (element: HTMLElement, prefix: string, localProgress: number, activeOpacity: number) => {
      const p = clampProgress(localProgress);
      element.style.setProperty('--journey-panel-opacity', activeOpacity.toFixed(4));
      element.dataset.journeyPanelActive = activeOpacity > 0.5 ? 'true' : 'false';

      if (prefix === 'learning-guide') {
        const vectorProgress = stageProgress(p, 0, 0.24);
        const lumiProgress = stageProgress(p, 0.08, 0.4);
        const bubbleProgress = stageProgress(p, 0.24, 0.58);
        const videoProgress = stageProgress(p, 0.42, 0.78);
        const detailsProgress = stageProgress(p, 0.68, 1);
        element.style.setProperty('--learning-guide-progress', p.toFixed(4));
        element.style.setProperty('--learning-guide-vector-progress', vectorProgress.toFixed(4));
        element.style.setProperty('--learning-guide-lumi-progress', lumiProgress.toFixed(4));
        element.style.setProperty('--learning-guide-bubble-progress', bubbleProgress.toFixed(4));
        element.style.setProperty('--learning-guide-video-progress', videoProgress.toFixed(4));
        element.style.setProperty('--learning-guide-details-progress', detailsProgress.toFixed(4));
        element.style.setProperty('--learning-guide-lumi-x', `${(-84 * (1 - lumiProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-lumi-y', `${(38 * (1 - lumiProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-bubble-y', `${(-42 * (1 - bubbleProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-video-x', `${(118 * (1 - videoProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-video-y', `${(26 * (1 - videoProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-details-y', `${(16 * (1 - detailsProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-lumi-scale', (0.94 + lumiProgress * 0.06).toFixed(4));
        element.style.setProperty('--learning-guide-video-scale', (0.965 + videoProgress * 0.035).toFixed(4));
        return;
      }

      const names = prefix === 'public-docs'
        ? ['--public-docs-progress', '--public-docs-lumi-progress', '--public-docs-bubble-progress', '--public-docs-copy-progress', '--public-docs-cards-progress', '--public-docs-lumi-x', '--public-docs-lumi-y', '--public-docs-bubble-y', '--public-docs-copy-y', '--public-docs-card-y', '--public-docs-lumi-scale']
        : prefix === 'ai-points'
          ? ['--ai-points-progress', '--ai-points-lumi-progress', '--ai-points-bubble-progress', '--ai-points-copy-progress', '--ai-points-cards-progress', '--ai-points-lumi-x', '--ai-points-lumi-y', '--ai-points-bubble-y', '--ai-points-copy-y', '--ai-points-card-y', '--ai-points-lumi-scale']
          : ['--platform-news-progress', '--platform-news-lumi-progress', '--platform-news-bubble-progress', '--platform-news-copy-progress', '--platform-news-board-progress', '--platform-news-lumi-x', '--platform-news-lumi-y', '--platform-news-bubble-y', '--platform-news-copy-y', '--platform-news-board-y', '--platform-news-lumi-scale'];
      const lumiProgress = stageProgress(p, 0.05, 0.38);
      const bubbleProgress = stageProgress(p, 0.2, 0.54);
      const copyProgress = stageProgress(p, 0.36, 0.7);
      const objectProgress = stageProgress(p, 0.52, 1);
      element.style.setProperty(names[0], p.toFixed(4));
      element.style.setProperty(names[1], lumiProgress.toFixed(4));
      element.style.setProperty(names[2], bubbleProgress.toFixed(4));
      element.style.setProperty(names[3], copyProgress.toFixed(4));
      element.style.setProperty(names[4], objectProgress.toFixed(4));
      element.style.setProperty(names[5], `${(-112 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty(names[6], `${(68 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty(names[7], `${(-58 * (1 - bubbleProgress)).toFixed(2)}px`);
      element.style.setProperty(names[8], `${(36 * (1 - copyProgress)).toFixed(2)}px`);
      element.style.setProperty(names[9], `${(54 * (1 - objectProgress)).toFixed(2)}px`);
      element.style.setProperty(names[10], (0.9 + lumiProgress * 0.1).toFixed(4));
    };

    const setProgress = () => {
      frame = 0;
      const panelElements = panels.map((ref) => ref.current).filter(Boolean) as HTMLElement[];
      if (!panelElements.length) return;

      if (reduceMotion.matches || window.innerWidth < 1024) {
        panelElements.forEach((panel) => {
          panel.style.setProperty('--journey-panel-opacity', '1');
          panel.dataset.journeyPanelActive = 'true';
        });
        return;
      }

      const rect = runway.getBoundingClientRect();
      const viewportHeight = window.innerHeight || 1;
      const travel = Math.max(1, rect.height - viewportHeight);
      const progress = clampProgress((-rect.top) / travel);
      const windows = [
        [0, 0.28],
        [0.28, 0.56],
        [0.56, 0.84],
        [0.84, 1],
      ] as const;
      const assemblyWindows = [
        [0, 0.12],
        [0.28, 0.4],
        [0.56, 0.68],
        [0.84, 0.96],
      ] as const;
      const prefixes = ['learning-guide', 'public-docs', 'ai-points', 'platform-news'] as const;

      panelElements.forEach((panel, index) => {
        const [from, to] = windows[index];
        const [assemblyFrom, assemblyTo] = assemblyWindows[index];
        const local = clampProgress((progress - assemblyFrom) / Math.max(0.001, assemblyTo - assemblyFrom));
        const enter = stageProgress(progress, from, from + 0.055);
        const exit = index === windows.length - 1 ? 1 : 1 - stageProgress(progress, to - 0.055, to);
        const opacity = clampProgress(Math.min(enter, exit));
        setPanelVars(panel, prefixes[index], local, opacity);
      });
    };

    const requestUpdate = () => {
      if (frame) return;
      frame = window.requestAnimationFrame(setProgress);
    };

    setProgress();
    window.addEventListener('scroll', requestUpdate, { passive: true });
    window.addEventListener('resize', requestUpdate);
    reduceMotion.addEventListener('change', requestUpdate);

    return () => {
      if (frame) window.cancelAnimationFrame(frame);
      window.removeEventListener('scroll', requestUpdate);
      window.removeEventListener('resize', requestUpdate);
      reduceMotion.removeEventListener('change', requestUpdate);
    };
  }, []);

  return (
            <section
              key={section.id}
              id={section.id}
              ref={aiPointsByokRef}
              data-ready-progress="0.74"
              data-mobile-ready-progress="1"
              className="landingJourneySection landingStickyJourneyPanel landingAiPointsByokSection"
              aria-label={section.label}
            >
              <span id="ai-points-byok-ready" className="landingJourneyReadyAnchor landingJourneyReadyAnchorAiPoints" aria-hidden="true" />
              <div className="landingAiPointsByokShell">
                <div className="landingAiPointsByokLumiPane">
                  <img
                    src="/images/lumi_journey_ai_points_byok_pose_alpha_v1.png"
                    alt=""
                    className="landingAiPointsByokLumiImage"
                    aria-hidden="true"
                  />
                  <div className="landingAiPointsByokSpeech">
                    <p>{aiUsageCopy.pointsTitle}</p>
                    <p>{aiUsageCopy.byokTitle}</p>
                  </div>
                </div>

                <div className="landingAiPointsByokCopy">
                  <p>{aiUsageCopy.eyebrow}</p>
                  <h2>
                    {aiUsageCopy.titlePrefix}
                    <span>{aiUsageCopy.titleHighlight}</span>
                  </h2>
                  <div>
                    {aiUsageCopy.descriptionLines.map((line) => (
                      <span key={line}>{line}</span>
                    ))}
                  </div>
                </div>

                <div className="landingAiPointsByokCards">
                  <article className="landingAiPointsByokCard">
                    <span>POINT</span>
                    <strong>{aiUsageCopy.pointsTitle}</strong>
                    {aiUsageCopy.pointsDescription(30).map((line) => (
                      <p key={line}>{line}</p>
                    ))}
                  </article>
                  <article className="landingAiPointsByokCard">
                    <span>{aiUsageCopy.byokBadge}</span>
                    <strong>{aiUsageCopy.byokTitle}</strong>
                    {aiUsageCopy.byokDescriptionLines.map((line) => (
                      <p key={line}>{line}</p>
                    ))}
                  </article>
                </div>
              </div>
            </section>
          );
        }

        if (index === 3) {


  useEffect(() => {
    const runway = journeyRef.current;
    if (!runway) return;

    const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)');
    const panels = [learningGuideRef, publicDocumentsRef, aiPointsByokRef, platformNewsRef];
    let frame = 0;

    const setPanelVars = (element: HTMLElement, prefix: string, localProgress: number, activeOpacity: number) => {
      const p = clampProgress(localProgress);
      element.style.setProperty('--journey-panel-opacity', activeOpacity.toFixed(4));
      element.dataset.journeyPanelActive = activeOpacity > 0.5 ? 'true' : 'false';

      if (prefix === 'learning-guide') {
        const vectorProgress = stageProgress(p, 0, 0.24);
        const lumiProgress = stageProgress(p, 0.08, 0.4);
        const bubbleProgress = stageProgress(p, 0.24, 0.58);
        const videoProgress = stageProgress(p, 0.42, 0.78);
        const detailsProgress = stageProgress(p, 0.68, 1);
        element.style.setProperty('--learning-guide-progress', p.toFixed(4));
        element.style.setProperty('--learning-guide-vector-progress', vectorProgress.toFixed(4));
        element.style.setProperty('--learning-guide-lumi-progress', lumiProgress.toFixed(4));
        element.style.setProperty('--learning-guide-bubble-progress', bubbleProgress.toFixed(4));
        element.style.setProperty('--learning-guide-video-progress', videoProgress.toFixed(4));
        element.style.setProperty('--learning-guide-details-progress', detailsProgress.toFixed(4));
        element.style.setProperty('--learning-guide-lumi-x', `${(-84 * (1 - lumiProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-lumi-y', `${(38 * (1 - lumiProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-bubble-y', `${(-42 * (1 - bubbleProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-video-x', `${(118 * (1 - videoProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-video-y', `${(26 * (1 - videoProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-details-y', `${(16 * (1 - detailsProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-lumi-scale', (0.94 + lumiProgress * 0.06).toFixed(4));
        element.style.setProperty('--learning-guide-video-scale', (0.965 + videoProgress * 0.035).toFixed(4));
        return;
      }

      const names = prefix === 'public-docs'
        ? ['--public-docs-progress', '--public-docs-lumi-progress', '--public-docs-bubble-progress', '--public-docs-copy-progress', '--public-docs-cards-progress', '--public-docs-lumi-x', '--public-docs-lumi-y', '--public-docs-bubble-y', '--public-docs-copy-y', '--public-docs-card-y', '--public-docs-lumi-scale']
        : prefix === 'ai-points'
          ? ['--ai-points-progress', '--ai-points-lumi-progress', '--ai-points-bubble-progress', '--ai-points-copy-progress', '--ai-points-cards-progress', '--ai-points-lumi-x', '--ai-points-lumi-y', '--ai-points-bubble-y', '--ai-points-copy-y', '--ai-points-card-y', '--ai-points-lumi-scale']
          : ['--platform-news-progress', '--platform-news-lumi-progress', '--platform-news-bubble-progress', '--platform-news-copy-progress', '--platform-news-board-progress', '--platform-news-lumi-x', '--platform-news-lumi-y', '--platform-news-bubble-y', '--platform-news-copy-y', '--platform-news-board-y', '--platform-news-lumi-scale'];
      const lumiProgress = stageProgress(p, 0.05, 0.38);
      const bubbleProgress = stageProgress(p, 0.2, 0.54);
      const copyProgress = stageProgress(p, 0.36, 0.7);
      const objectProgress = stageProgress(p, 0.52, 1);
      element.style.setProperty(names[0], p.toFixed(4));
      element.style.setProperty(names[1], lumiProgress.toFixed(4));
      element.style.setProperty(names[2], bubbleProgress.toFixed(4));
      element.style.setProperty(names[3], copyProgress.toFixed(4));
      element.style.setProperty(names[4], objectProgress.toFixed(4));
      element.style.setProperty(names[5], `${(-112 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty(names[6], `${(68 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty(names[7], `${(-58 * (1 - bubbleProgress)).toFixed(2)}px`);
      element.style.setProperty(names[8], `${(36 * (1 - copyProgress)).toFixed(2)}px`);
      element.style.setProperty(names[9], `${(54 * (1 - objectProgress)).toFixed(2)}px`);
      element.style.setProperty(names[10], (0.9 + lumiProgress * 0.1).toFixed(4));
    };

    const setProgress = () => {
      frame = 0;
      const panelElements = panels.map((ref) => ref.current).filter(Boolean) as HTMLElement[];
      if (!panelElements.length) return;

      if (reduceMotion.matches || window.innerWidth < 1024) {
        panelElements.forEach((panel) => {
          panel.style.setProperty('--journey-panel-opacity', '1');
          panel.dataset.journeyPanelActive = 'true';
        });
        return;
      }

      const rect = runway.getBoundingClientRect();
      const viewportHeight = window.innerHeight || 1;
      const travel = Math.max(1, rect.height - viewportHeight);
      const progress = clampProgress((-rect.top) / travel);
      const windows = [
        [0, 0.28],
        [0.28, 0.56],
        [0.56, 0.84],
        [0.84, 1],
      ] as const;
      const assemblyWindows = [
        [0, 0.12],
        [0.28, 0.4],
        [0.56, 0.68],
        [0.84, 0.96],
      ] as const;
      const prefixes = ['learning-guide', 'public-docs', 'ai-points', 'platform-news'] as const;

      panelElements.forEach((panel, index) => {
        const [from, to] = windows[index];
        const [assemblyFrom, assemblyTo] = assemblyWindows[index];
        const local = clampProgress((progress - assemblyFrom) / Math.max(0.001, assemblyTo - assemblyFrom));
        const enter = stageProgress(progress, from, from + 0.055);
        const exit = index === windows.length - 1 ? 1 : 1 - stageProgress(progress, to - 0.055, to);
        const opacity = clampProgress(Math.min(enter, exit));
        setPanelVars(panel, prefixes[index], local, opacity);
      });
    };

    const requestUpdate = () => {
      if (frame) return;
      frame = window.requestAnimationFrame(setProgress);
    };

    setProgress();
    window.addEventListener('scroll', requestUpdate, { passive: true });
    window.addEventListener('resize', requestUpdate);
    reduceMotion.addEventListener('change', requestUpdate);

    return () => {
      if (frame) window.cancelAnimationFrame(frame);
      window.removeEventListener('scroll', requestUpdate);
      window.removeEventListener('resize', requestUpdate);
      reduceMotion.removeEventListener('change', requestUpdate);
    };
  }, []);

  return (
            <section
              key={section.id}
              id={section.id}
              ref={platformNewsRef}
              data-ready-progress="0.97"
              data-mobile-ready-progress="1"
              className="landingJourneySection landingStickyJourneyPanel landingPlatformNewsSection"
              aria-label={section.label}
            >
              <span id="platform-news-ready" className="landingJourneyReadyAnchor landingJourneyReadyAnchorPlatformNews" aria-hidden="true" />
              <div className="landingPlatformNewsShell">
                <div className="landingPlatformNewsLumiPane">
                  <img
                    src="/images/lumi_journey_platform_news_pose_alpha_v1.png"
                    alt=""
                    className="landingPlatformNewsLumiImage"
                    aria-hidden="true"
                  />
                  <div className="landingPlatformNewsSpeech">
                    <p>{platformNewsCopy.speech}</p>
                  </div>
                </div>

                <div className="landingPlatformNewsCopy">
                  <p>{platformNewsCopy.eyebrow}</p>
                  <h2>{platformNewsCopy.title}</h2>
                  <span>{platformNewsCopy.description}</span>
                </div>

                <div className="landingPlatformNewsBoard" aria-label={section.label}>
                  <div className="landingPlatformNewsBoardHeader">
                    <span>Notice</span>
                    <strong>{section.label}</strong>
                  </div>
                  <div className="landingPlatformNewsEmptyList">
                    <strong>{platformNewsCopy.emptyTitle}</strong>
                    <p>{platformNewsCopy.emptyDescription}</p>
                    <a className="landingPlatformNewsCta" href="/platform/news">{platformNewsCopy.ctaLabel}</a>
                  </div>
                </div>
              </div>
            </section>
          );
        }



  useEffect(() => {
    const runway = journeyRef.current;
    if (!runway) return;

    const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)');
    const panels = [learningGuideRef, publicDocumentsRef, aiPointsByokRef, platformNewsRef];
    let frame = 0;

    const setPanelVars = (element: HTMLElement, prefix: string, localProgress: number, activeOpacity: number) => {
      const p = clampProgress(localProgress);
      element.style.setProperty('--journey-panel-opacity', activeOpacity.toFixed(4));
      element.dataset.journeyPanelActive = activeOpacity > 0.5 ? 'true' : 'false';

      if (prefix === 'learning-guide') {
        const vectorProgress = stageProgress(p, 0, 0.24);
        const lumiProgress = stageProgress(p, 0.08, 0.4);
        const bubbleProgress = stageProgress(p, 0.24, 0.58);
        const videoProgress = stageProgress(p, 0.42, 0.78);
        const detailsProgress = stageProgress(p, 0.68, 1);
        element.style.setProperty('--learning-guide-progress', p.toFixed(4));
        element.style.setProperty('--learning-guide-vector-progress', vectorProgress.toFixed(4));
        element.style.setProperty('--learning-guide-lumi-progress', lumiProgress.toFixed(4));
        element.style.setProperty('--learning-guide-bubble-progress', bubbleProgress.toFixed(4));
        element.style.setProperty('--learning-guide-video-progress', videoProgress.toFixed(4));
        element.style.setProperty('--learning-guide-details-progress', detailsProgress.toFixed(4));
        element.style.setProperty('--learning-guide-lumi-x', `${(-84 * (1 - lumiProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-lumi-y', `${(38 * (1 - lumiProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-bubble-y', `${(-42 * (1 - bubbleProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-video-x', `${(118 * (1 - videoProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-video-y', `${(26 * (1 - videoProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-details-y', `${(16 * (1 - detailsProgress)).toFixed(2)}px`);
        element.style.setProperty('--learning-guide-lumi-scale', (0.94 + lumiProgress * 0.06).toFixed(4));
        element.style.setProperty('--learning-guide-video-scale', (0.965 + videoProgress * 0.035).toFixed(4));
        return;
      }

      const names = prefix === 'public-docs'
        ? ['--public-docs-progress', '--public-docs-lumi-progress', '--public-docs-bubble-progress', '--public-docs-copy-progress', '--public-docs-cards-progress', '--public-docs-lumi-x', '--public-docs-lumi-y', '--public-docs-bubble-y', '--public-docs-copy-y', '--public-docs-card-y', '--public-docs-lumi-scale']
        : prefix === 'ai-points'
          ? ['--ai-points-progress', '--ai-points-lumi-progress', '--ai-points-bubble-progress', '--ai-points-copy-progress', '--ai-points-cards-progress', '--ai-points-lumi-x', '--ai-points-lumi-y', '--ai-points-bubble-y', '--ai-points-copy-y', '--ai-points-card-y', '--ai-points-lumi-scale']
          : ['--platform-news-progress', '--platform-news-lumi-progress', '--platform-news-bubble-progress', '--platform-news-copy-progress', '--platform-news-board-progress', '--platform-news-lumi-x', '--platform-news-lumi-y', '--platform-news-bubble-y', '--platform-news-copy-y', '--platform-news-board-y', '--platform-news-lumi-scale'];
      const lumiProgress = stageProgress(p, 0.05, 0.38);
      const bubbleProgress = stageProgress(p, 0.2, 0.54);
      const copyProgress = stageProgress(p, 0.36, 0.7);
      const objectProgress = stageProgress(p, 0.52, 1);
      element.style.setProperty(names[0], p.toFixed(4));
      element.style.setProperty(names[1], lumiProgress.toFixed(4));
      element.style.setProperty(names[2], bubbleProgress.toFixed(4));
      element.style.setProperty(names[3], copyProgress.toFixed(4));
      element.style.setProperty(names[4], objectProgress.toFixed(4));
      element.style.setProperty(names[5], `${(-112 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty(names[6], `${(68 * (1 - lumiProgress)).toFixed(2)}px`);
      element.style.setProperty(names[7], `${(-58 * (1 - bubbleProgress)).toFixed(2)}px`);
      element.style.setProperty(names[8], `${(36 * (1 - copyProgress)).toFixed(2)}px`);
      element.style.setProperty(names[9], `${(54 * (1 - objectProgress)).toFixed(2)}px`);
      element.style.setProperty(names[10], (0.9 + lumiProgress * 0.1).toFixed(4));
    };

    const setProgress = () => {
      frame = 0;
      const panelElements = panels.map((ref) => ref.current).filter(Boolean) as HTMLElement[];
      if (!panelElements.length) return;

      if (reduceMotion.matches || window.innerWidth < 1024) {
        panelElements.forEach((panel) => {
          panel.style.setProperty('--journey-panel-opacity', '1');
          panel.dataset.journeyPanelActive = 'true';
        });
        return;
      }

      const rect = runway.getBoundingClientRect();
      const viewportHeight = window.innerHeight || 1;
      const travel = Math.max(1, rect.height - viewportHeight);
      const progress = clampProgress((-rect.top) / travel);
      const windows = [
        [0, 0.28],
        [0.28, 0.56],
        [0.56, 0.84],
        [0.84, 1],
      ] as const;
      const assemblyWindows = [
        [0, 0.12],
        [0.28, 0.4],
        [0.56, 0.68],
        [0.84, 0.96],
      ] as const;
      const prefixes = ['learning-guide', 'public-docs', 'ai-points', 'platform-news'] as const;

      panelElements.forEach((panel, index) => {
        const [from, to] = windows[index];
        const [assemblyFrom, assemblyTo] = assemblyWindows[index];
        const local = clampProgress((progress - assemblyFrom) / Math.max(0.001, assemblyTo - assemblyFrom));
        const enter = stageProgress(progress, from, from + 0.055);
        const exit = index === windows.length - 1 ? 1 : 1 - stageProgress(progress, to - 0.055, to);
        const opacity = clampProgress(Math.min(enter, exit));
        setPanelVars(panel, prefixes[index], local, opacity);
      });
    };

    const requestUpdate = () => {
      if (frame) return;
      frame = window.requestAnimationFrame(setProgress);
    };

    setProgress();
    window.addEventListener('scroll', requestUpdate, { passive: true });
    window.addEventListener('resize', requestUpdate);
    reduceMotion.addEventListener('change', requestUpdate);

    return () => {
      if (frame) window.cancelAnimationFrame(frame);
      window.removeEventListener('scroll', requestUpdate);
      window.removeEventListener('resize', requestUpdate);
      reduceMotion.removeEventListener('change', requestUpdate);
    };
  }, []);

  return (
          <section
            key={section.id}
            id={section.id}
            className="landingJourneySection"
            aria-label={section.label}
          >
            <span id={`${section.id}-ready`} className="landingJourneyReadyAnchor landingJourneyReadyAnchorDefault" aria-hidden="true" />
            <div className="landingJourneySectionMarker" aria-hidden="true">
              <span>{section.marker}</span>
              <strong>{section.label}</strong>
            </div>
          </section>
        );
      })}
      </div>
    </div>
  );
}
