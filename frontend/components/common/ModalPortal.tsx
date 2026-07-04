'use client';

import { useEffect, useRef, useState, type CSSProperties, type MouseEvent, type ReactNode } from 'react';
import { createPortal } from 'react-dom';

interface ModalPortalProps {
  children: ReactNode;
  overlayStyle?: CSSProperties;
  onMouseDown?: (event: MouseEvent<HTMLDivElement>) => void;
  className?: string;
}

const modalRootAttribute = 'data-learnweaver-modal-root';
let activeModalCount = 0;
let scrollYBeforeLock = 0;
const inertState = new Map<HTMLElement, { inert: boolean; ariaHidden: string | null }>();
let pageStyleBeforeLock: {
  bodyPosition: string;
  bodyTop: string;
  bodyLeft: string;
  bodyRight: string;
  bodyWidth: string;
  bodyOverflow: string;
  bodyPaddingRight: string;
  htmlOverflow: string;
} | null = null;

export const modalOverlayBaseStyle: CSSProperties = {
  position: 'fixed',
  inset: 0,
  width: '100vw',
  height: '100dvh',
  zIndex: 2147483000,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  boxSizing: 'border-box',
  padding: '24px 16px',
  overflow: 'hidden',
  overscrollBehavior: 'contain',
  pointerEvents: 'auto',
};

export default function ModalPortal({ children, overlayStyle, onMouseDown, className }: ModalPortalProps) {
  const [isMounted, setIsMounted] = useState(false);
  const portalNodeRef = useRef<HTMLDivElement | null>(null);

  if (typeof document !== 'undefined' && portalNodeRef.current === null) {
    const node = document.createElement('div');
    node.setAttribute(modalRootAttribute, 'true');
    portalNodeRef.current = node;
  }

  useEffect(() => {
    const node = portalNodeRef.current;
    if (!node) return undefined;

    document.body.appendChild(node);
    activeModalCount += 1;
    if (activeModalCount === 1) {
      lockPage();
    }
    applyBackgroundInert();
    setIsMounted(true);

    return () => {
      setIsMounted(false);
      activeModalCount = Math.max(0, activeModalCount - 1);
      if (node.parentNode) {
        node.parentNode.removeChild(node);
      }
      if (activeModalCount === 0) {
        unlockPage();
        restoreBackgroundInert();
      } else {
        applyBackgroundInert();
      }
    };
  }, []);

  if (!isMounted || !portalNodeRef.current) return null;

  return createPortal(
    <div
      role="presentation"
      className={className}
      style={{ ...modalOverlayBaseStyle, ...overlayStyle }}
      onMouseDown={onMouseDown}
    >
      {children}
    </div>,
    portalNodeRef.current,
  );
}

function lockPage() {
  const { body, documentElement } = document;
  scrollYBeforeLock = window.scrollY;
  const scrollbarWidth = window.innerWidth - documentElement.clientWidth;
  pageStyleBeforeLock = {
    bodyPosition: body.style.position,
    bodyTop: body.style.top,
    bodyLeft: body.style.left,
    bodyRight: body.style.right,
    bodyWidth: body.style.width,
    bodyOverflow: body.style.overflow,
    bodyPaddingRight: body.style.paddingRight,
    htmlOverflow: documentElement.style.overflow,
  };

  body.dataset.modalScrollLock = 'true';
  body.style.position = 'fixed';
  body.style.top = `-${scrollYBeforeLock}px`;
  body.style.left = '0';
  body.style.right = '0';
  body.style.width = '100%';
  body.style.overflow = 'hidden';
  documentElement.style.overflow = 'hidden';
  if (scrollbarWidth > 0) {
    body.style.paddingRight = `${scrollbarWidth}px`;
  }
}

function unlockPage() {
  const { body, documentElement } = document;
  body.style.position = pageStyleBeforeLock?.bodyPosition ?? '';
  body.style.top = pageStyleBeforeLock?.bodyTop ?? '';
  body.style.left = pageStyleBeforeLock?.bodyLeft ?? '';
  body.style.right = pageStyleBeforeLock?.bodyRight ?? '';
  body.style.width = pageStyleBeforeLock?.bodyWidth ?? '';
  body.style.overflow = pageStyleBeforeLock?.bodyOverflow ?? '';
  body.style.paddingRight = pageStyleBeforeLock?.bodyPaddingRight ?? '';
  documentElement.style.overflow = pageStyleBeforeLock?.htmlOverflow ?? '';
  pageStyleBeforeLock = null;
  delete body.dataset.modalScrollLock;
  window.scrollTo(0, scrollYBeforeLock);
}

function applyBackgroundInert() {
  const bodyChildren = Array.from(document.body.children);
  bodyChildren.forEach((child) => {
    if (!(child instanceof HTMLElement) || child.hasAttribute(modalRootAttribute)) return;
    if (!inertState.has(child)) {
      inertState.set(child, {
        inert: child.inert,
        ariaHidden: child.getAttribute('aria-hidden'),
      });
    }
    child.inert = true;
    child.setAttribute('aria-hidden', 'true');
  });
}

function restoreBackgroundInert() {
  inertState.forEach((state, element) => {
    element.inert = state.inert;
    if (state.ariaHidden === null) {
      element.removeAttribute('aria-hidden');
    } else {
      element.setAttribute('aria-hidden', state.ariaHidden);
    }
  });
  inertState.clear();
}
