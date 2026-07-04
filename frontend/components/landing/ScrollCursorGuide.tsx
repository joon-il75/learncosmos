'use client';

import { useEffect, useRef } from 'react';

export default function ScrollCursorGuide({ disabled = false }: { disabled?: boolean }) {
  const cursorRef = useRef<HTMLImageElement>(null);

  useEffect(() => {
    const cursor = cursorRef.current;
    if (!cursor || disabled) {
      if (cursor) cursor.dataset.active = 'false';
      document.body.classList.remove('landingScrollCursorActive');
      return;
    }

    let frame = 0;
    let pointerX = window.innerWidth / 2;
    let pointerY = window.innerHeight / 2;
    let isPointerInside = false;

    const isFinePointer = () => window.matchMedia('(pointer: fine)').matches && window.innerWidth >= 1024;

    const getScrollRange = () => {
      const hero = document.querySelector<HTMLElement>('.heroShell');
      const guide = document.querySelector<HTMLElement>('#lumi-scroll-guide');
      if (!hero || !guide) return null;

      const heroTop = hero.getBoundingClientRect().top + window.scrollY;
      const guideBottom = guide.getBoundingClientRect().bottom + window.scrollY;
      return { start: heroTop, end: guideBottom };
    };

    const isInsideHeroCTA = () => {
      const cta = document.querySelector<HTMLElement>('.heroCTAStack');
      if (!cta) return false;

      const rect = cta.getBoundingClientRect();
      const padding = 10;
      return pointerX >= rect.left - padding
        && pointerX <= rect.right + padding
        && pointerY >= rect.top - padding
        && pointerY <= rect.bottom + padding;
    };

    const setActive = (active: boolean) => {
      cursor.dataset.active = active ? 'true' : 'false';
      document.body.classList.toggle('landingScrollCursorActive', active);
    };

    const update = () => {
      frame = 0;
      const range = getScrollRange();
      const scrollCenter = window.scrollY + window.innerHeight * 0.45;
      const active = Boolean(range && isFinePointer() && isPointerInside && !isInsideHeroCTA() && scrollCenter >= range.start && scrollCenter <= range.end);

      cursor.style.setProperty('--landing-scroll-cursor-x', pointerX.toFixed(1) + 'px');
      cursor.style.setProperty('--landing-scroll-cursor-y', pointerY.toFixed(1) + 'px');
      setActive(active);
    };

    const schedule = () => {
      if (frame) return;
      frame = window.requestAnimationFrame(update);
    };

    const handleMouseMove = (event: MouseEvent) => {
      pointerX = event.clientX;
      pointerY = event.clientY;
      isPointerInside = true;
      schedule();
    };

    const handleMouseLeave = () => {
      isPointerInside = false;
      schedule();
    };

    window.addEventListener('mousemove', handleMouseMove, { passive: true });
    window.addEventListener('mouseleave', handleMouseLeave);
    window.addEventListener('scroll', schedule, { passive: true });
    window.addEventListener('resize', schedule);
    schedule();

    return () => {
      if (frame) window.cancelAnimationFrame(frame);
      window.removeEventListener('mousemove', handleMouseMove);
      window.removeEventListener('mouseleave', handleMouseLeave);
      window.removeEventListener('scroll', schedule);
      window.removeEventListener('resize', schedule);
      document.body.classList.remove('landingScrollCursorActive');
    };
  }, [disabled]);

  return (
    <img
      src="/images/lumi_scroll_cursor.gif"
      alt=""
      aria-hidden="true"
      className="landingScrollCursor"
      data-active="false"
      ref={cursorRef}
    />
  );
}
