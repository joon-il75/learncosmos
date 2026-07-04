'use client';

import type { LandingPageCopy } from '@/lib/i18n/pages/landing';

const YOUTUBE_VIDEO_ID = '76bq5ixUzmE';

export default function LearningGuideVideoGate({
  copy,
  locale,
  disabled,
}: {
  copy: LandingPageCopy['introVideo'];
  locale?: string | null;
  disabled?: boolean;
}) {
  void locale;
  void disabled;

  const origin = typeof window === 'undefined' ? '' : `&origin=${window.location.origin}`;
  const playerSrc = `https://www.youtube.com/embed/${YOUTUBE_VIDEO_ID}?rel=0&playsinline=1${origin}`;

  return (
    <div
      className="landingLearningGuideVideoCard"
      data-video-gate="inline"
      data-video-ready="true"
      aria-label={copy.videoTitle}
    >
      <div className="landingLearningGuideVideoFrame">
        <iframe
          className="landingLearningGuideVideoEmbed"
          src={playerSrc}
          title={copy.videoTitle}
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
          allowFullScreen
        />
      </div>
      <p className="landingLearningGuideFlow">{copy.simulationFlowLabel}</p>
    </div>
  );
}
