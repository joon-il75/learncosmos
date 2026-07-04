'use client';

import { useEffect, useRef, useState } from 'react';

interface VideoPlayerProps {
  title: string;
  meta?: string;
  sourceKey?: string;
  posterKey?: string;
  subtitleKey?: string;
  loadSource: () => Promise<string | null>;
  loadPosterSource?: () => Promise<string | null>;
  loadSubtitleSource?: () => Promise<string | null>;
  borderColor: string;
  background: string;
  textColor: string;
  mutedColor: string;
}

export default function VideoPlayer({ title, meta, sourceKey, posterKey, subtitleKey, loadSource, loadPosterSource, loadSubtitleSource, borderColor, background, textColor, mutedColor }: VideoPlayerProps) {
  const videoRef = useRef<HTMLVideoElement | null>(null);
  const [source, setSource] = useState<string | null>(null);
  const [posterSource, setPosterSource] = useState<string | null>(null);
  const [subtitleSource, setSubtitleSource] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    setSource(null);
    setError(null);
    loadSource()
      .then((nextSource) => {
        if (cancelled) return;
        if (!nextSource) {
          setError('동영상 URL을 불러오지 못했습니다.');
          return;
        }
        setSource(nextSource);
      })
      .catch((e) => {
        if (cancelled) return;
        setError(e instanceof Error ? e.message : '동영상 URL을 불러오지 못했습니다.');
      });
    return () => {
      cancelled = true;
    };
  }, [sourceKey]);

  useEffect(() => {
    if (!source || posterSource) return;
    const video = videoRef.current;
    if (!video) return;
    video.load();
  }, [source, posterSource]);

  useEffect(() => {
    let cancelled = false;
    setPosterSource(null);
    if (!loadPosterSource) return () => {
      cancelled = true;
    };
    loadPosterSource()
      .then((nextSource) => {
        if (!cancelled) setPosterSource(nextSource);
      })
      .catch(() => {
        if (!cancelled) setPosterSource(null);
      });
    return () => {
      cancelled = true;
    };
  }, [posterKey]);

  useEffect(() => {
    let cancelled = false;
    setSubtitleSource(null);
    if (!loadSubtitleSource) return () => {
      cancelled = true;
    };
    loadSubtitleSource()
      .then((nextSource) => {
        if (!cancelled) setSubtitleSource(nextSource);
      })
      .catch(() => {
        if (!cancelled) setSubtitleSource(null);
      });
    return () => {
      cancelled = true;
    };
  }, [subtitleKey]);

  return (
    <div style={{ width: '100%', maxWidth: '860px', margin: '0 auto', display: 'grid', gap: '8px' }}>
      <div style={{ display: 'grid', gap: '4px' }}>
        <strong style={{ color: textColor, fontSize: '15px', lineHeight: 1.35 }}>{title}</strong>
        {meta ? <span style={{ color: mutedColor, fontSize: '13px', lineHeight: 1.45 }}>{meta}</span> : null}
      </div>
      {source ? (
        <video
          ref={videoRef}
          controls
          preload={posterSource ? 'metadata' : 'auto'}
          src={source}
          poster={posterSource ?? undefined}
          onLoadedMetadata={(event) => {
            if (posterSource) return;
            const video = event.currentTarget;
            if (!Number.isFinite(video.duration) || video.duration <= 0) return;
            try {
              video.currentTime = Math.min(0.001, video.duration);
            } catch {
              // Some browsers/codecs reject programmatic seek before enough data is buffered.
            }
          }}
          onLoadedData={(event) => {
            if (!posterSource) event.currentTarget.pause();
          }}
          style={{
            width: '100%',
            maxHeight: '520px',
            borderRadius: '8px',
            border: `1px solid ${borderColor}`,
            background,
          }}
        >
          {subtitleSource ? <track kind="subtitles" src={subtitleSource} srcLang="ko" label="한국어" default /> : null}
        </video>
      ) : (
        <div
          style={{
            minHeight: '120px',
            display: 'grid',
            placeItems: 'center',
            borderRadius: '8px',
            border: `1px solid ${borderColor}`,
            background,
            color: error ? '#B91C1C' : mutedColor,
            fontSize: '14px',
            fontWeight: 700,
            textAlign: 'center',
            padding: '18px',
          }}
        >
          {error ?? '동영상을 불러오는 중입니다.'}
        </div>
      )}
    </div>
  );
}
