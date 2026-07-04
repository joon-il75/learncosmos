'use client';

import { useEffect, useRef, type CSSProperties } from 'react';
import {
  getProgressGridPosition,
  normalizePlanetProgressPercent,
} from '@/lib/world-ui-engine/planetProgressUtils';

type RotationDirection = 'left' | 'right';

interface RotatingPlanetCanvasProps {
  atlasURL: string | null;
  progressPercent: number;
  size: number;
  durationSeconds?: number;
  direction?: RotationDirection;
  playing?: boolean;
  className?: string;
  style?: CSSProperties;
}

const ATLAS_COLUMNS = 4;
const ATLAS_ROWS = 3;
const DEFAULT_ROTATION_SECONDS = 36;
const imageCache = new Map<string, HTMLImageElement>();

export default function RotatingPlanetCanvas({
  atlasURL,
  progressPercent,
  size,
  durationSeconds = DEFAULT_ROTATION_SECONDS,
  direction = 'left',
  playing = true,
  className,
  style,
}: RotatingPlanetCanvasProps) {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas || !atlasURL) return;

    const context = canvas.getContext('2d', { willReadFrequently: true });
    if (!context) return;

    let cancelled = false;
    let animationFrame = 0;
    let lastTime = performance.now();
    let rotationOffset = 0;

    const image = loadAtlasImage(atlasURL);
    const startDrawing = () => {
      if (cancelled) return;
      const sourceCanvas = document.createElement('canvas');
      sourceCanvas.width = image.naturalWidth || image.width;
      sourceCanvas.height = image.naturalHeight || image.height;
      const sourceContext = sourceCanvas.getContext('2d', { willReadFrequently: true });
      if (!sourceContext || sourceCanvas.width <= 0 || sourceCanvas.height <= 0) return;

      sourceContext.drawImage(image, 0, 0, sourceCanvas.width, sourceCanvas.height);
      const cellWidth = Math.floor(sourceCanvas.width / ATLAS_COLUMNS);
      const cellHeight = Math.floor(sourceCanvas.height / ATLAS_ROWS);
      const frame = getProgressGridPosition(normalizePlanetProgressPercent(progressPercent));
      const sourceX = (frame.col - 1) * cellWidth;
      const sourceY = (frame.row - 1) * cellHeight;
      const cellData = sourceContext.getImageData(sourceX, sourceY, cellWidth, cellHeight);

      const render = (now: number) => {
        if (cancelled) return;
        const elapsed = Math.max(0, now - lastTime);
        lastTime = now;
        if (playing && durationSeconds > 0) {
          const delta = (elapsed / 1000) / durationSeconds;
          rotationOffset = wrapUnit(rotationOffset + (direction === 'left' ? delta : -delta));
        }
        drawPlanet(context, cellData, canvas.width, canvas.height, rotationOffset);
        animationFrame = window.requestAnimationFrame(render);
      };

      render(performance.now());
    };

    image.addEventListener('load', startDrawing);
    if (image.complete && image.naturalWidth > 0) {
      startDrawing();
    }

    return () => {
      cancelled = true;
      image.removeEventListener('load', startDrawing);
      if (animationFrame) window.cancelAnimationFrame(animationFrame);
    };
  }, [atlasURL, direction, durationSeconds, playing, progressPercent]);

  return (
    <canvas
      ref={canvasRef}
      className={className}
      width={Math.max(96, Math.round(size))}
      height={Math.max(96, Math.round(size))}
      style={{
        width: size,
        height: size,
        display: 'block',
        ...style,
      }}
    />
  );
}

function loadAtlasImage(src: string): HTMLImageElement {
  const cached = imageCache.get(src);
  if (cached) return cached;
  const image = new Image();
  image.crossOrigin = 'anonymous';
  image.src = src;
  imageCache.set(src, image);
  return image;
}

function wrapUnit(value: number): number {
  return ((value % 1) + 1) % 1;
}

function drawPlanet(
  context: CanvasRenderingContext2D,
  cellData: ImageData,
  width: number,
  height: number,
  rotationOffset: number,
) {
  const output = context.createImageData(width, height);
  const radius = Math.min(width, height) * 0.46;
  const centerX = width / 2;
  const centerY = height / 2;
  const radiusSquared = radius * radius;

  for (let y = 0; y < height; y += 1) {
    const dy = (y - centerY) / radius;
    for (let x = 0; x < width; x += 1) {
      const dx = (x - centerX) / radius;
      const distanceSquared = dx * dx + dy * dy;
      const targetIndex = (y * width + x) * 4;

      if (distanceSquared > 1) {
        output.data[targetIndex + 3] = 0;
        continue;
      }

      const dz = Math.sqrt(Math.max(0, 1 - distanceSquared));
      const longitude = Math.atan2(dx, dz);
      const latitude = Math.asin(clamp(dy, -1, 1));
      const u = wrapUnit((longitude / (Math.PI * 2)) + 0.5 + rotationOffset);
      const v = clamp(0.5 - (latitude / Math.PI), 0, 1);
      const sourceX = Math.min(cellData.width - 1, Math.floor(u * cellData.width));
      const sourceY = Math.min(cellData.height - 1, Math.floor(v * cellData.height));
      const sourceIndex = (sourceY * cellData.width + sourceX) * 4;

      const rim = Math.pow(1 - dz, 2.2);
      const light = clamp(0.62 + dz * 0.62 - dx * 0.16 - dy * 0.08 - rim * 0.24, 0.28, 1.22);
      const atmosphere = Math.max(0, 1 - Math.sqrt(distanceSquared));

      output.data[targetIndex] = clampColor(cellData.data[sourceIndex] * light + atmosphere * 10);
      output.data[targetIndex + 1] = clampColor(cellData.data[sourceIndex + 1] * light + atmosphere * 14);
      output.data[targetIndex + 2] = clampColor(cellData.data[sourceIndex + 2] * light + atmosphere * 20);
      output.data[targetIndex + 3] = clampColor(cellData.data[sourceIndex + 3] * (1 - rim * 0.18));
    }
  }

  context.clearRect(0, 0, width, height);
  context.putImageData(output, 0, 0);
}

function clamp(value: number, min: number, max: number): number {
  return Math.max(min, Math.min(max, value));
}

function clampColor(value: number): number {
  return Math.max(0, Math.min(255, Math.round(value)));
}
