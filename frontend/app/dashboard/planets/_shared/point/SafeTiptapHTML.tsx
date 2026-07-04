'use client';

import type { CSSProperties } from 'react';

interface Props {
  html: string;
  className?: string;
  style?: CSSProperties;
}

const emptyHTML = '<p></p>';
const allowedElements = new Set([
  'a',
  'blockquote',
  'br',
  'em',
  'h1',
  'h2',
  'h3',
  'li',
  'mark',
  'ol',
  'p',
  's',
  'strong',
  'sub',
  'sup',
  'table',
  'tbody',
  'td',
  'th',
  'thead',
  'tr',
  'u',
  'ul',
  'img',
]);
const dropElements = new Set([
  'script',
  'style',
  'iframe',
  'object',
  'embed',
  'svg',
  'math',
  'meta',
  'link',
  'base',
  'form',
  'input',
  'button',
  'textarea',
  'select',
]);
const blockedElementPattern = /<\s*(script|style|iframe|object|embed|svg|math|meta|link|base|form|input|button|textarea|select)\b[^>]*>[\s\S]*?<\s*\/\s*\1\s*>/gi;
const blockedVoidElementPattern = /<\s*(script|style|iframe|object|embed|svg|math|meta|link|base|form|input|button|textarea|select)\b[^>]*\/?\s*>/gi;
const eventAttributePattern = /\s+on[a-z]+\s*=\s*("[^"]*"|'[^']*'|[^\s>]+)/gi;
const dangerousURLAttributePattern = /\s+(href|src)\s*=\s*("|')?\s*(javascript:|data:|vbscript:)[^"'\s>]*(\2)?/gi;

export default function SafeTiptapHTML({ html, className, style }: Props) {
  const sanitized = sanitizeTiptapHTMLForRender(html);
  return <div className={className} style={style} dangerouslySetInnerHTML={{ __html: sanitized }} />;
}

export function sanitizeTiptapHTMLForRender(value: string): string {
  const source = value.trim().length ? value : emptyHTML;
  if (typeof DOMParser === 'undefined') {
    return sanitizeTiptapHTMLFallback(source);
  }

  const parser = new DOMParser();
  const doc = parser.parseFromString(`<!doctype html><html><body>${source}</body></html>`, 'text/html');
  sanitizeChildren(doc.body);
  const sanitized = doc.body.innerHTML.trim();
  return sanitized.length ? sanitized : emptyHTML;
}

function sanitizeChildren(parent: ParentNode) {
  for (const child of Array.from(parent.childNodes)) {
    if (child.nodeType !== Node.ELEMENT_NODE) {
      continue;
    }

    const element = child as HTMLElement;
    const tagName = element.tagName.toLowerCase();
    if (dropElements.has(tagName)) {
      element.remove();
      continue;
    }
    if (!allowedElements.has(tagName)) {
      unwrapElement(element);
      continue;
    }

    sanitizeElementAttributes(element, tagName);
    sanitizeChildren(element);
  }
}

function unwrapElement(element: HTMLElement) {
  const parent = element.parentNode;
  if (!parent) return;
  while (element.firstChild) {
    parent.insertBefore(element.firstChild, element);
  }
  element.remove();
}

function sanitizeElementAttributes(element: HTMLElement, tagName: string) {
  const nextAttributes: Array<[string, string]> = [];
  for (const attr of Array.from(element.attributes)) {
    const key = attr.name.toLowerCase().trim();
    const attrValue = attr.value.trim();
    if (!key || key.startsWith('on')) continue;

    if (tagName === 'a' && key === 'href') {
      if (!isSafeHTTPSURL(attrValue)) {
        unwrapElement(element);
        return;
      }
      nextAttributes.push(['href', attrValue]);
      continue;
    }

    if (tagName === 'img') {
      if (key === 'src') {
        if (!isSafeHTTPSURL(attrValue)) {
          element.remove();
          return;
        }
        nextAttributes.push(['src', attrValue]);
      } else if (key === 'alt' || key === 'title') {
        nextAttributes.push([key, attrValue]);
      } else if (key === 'loading' && attrValue.toLowerCase() === 'lazy') {
        nextAttributes.push(['loading', 'lazy']);
      }
      continue;
    }

    if (tagName === 'table' && key === 'class' && attrValue.split(/\s+/).includes('lw-tiptap-table')) {
      nextAttributes.push(['class', 'lw-tiptap-table']);
      continue;
    }

    if ((tagName === 'td' || tagName === 'th') && (key === 'colspan' || key === 'rowspan')) {
      if (isSafePositiveInt(attrValue, 1, 24)) {
        nextAttributes.push([key, attrValue]);
      }
      continue;
    }

    if (key === 'style') {
      const style = sanitizeTiptapStyle(attrValue);
      if (style) nextAttributes.push(['style', style]);
    }
  }

  for (const attr of Array.from(element.attributes)) {
    element.removeAttribute(attr.name);
  }
  if (tagName === 'a') {
    nextAttributes.push(['rel', 'noopener noreferrer']);
    nextAttributes.push(['target', '_blank']);
  }
  for (const [key, attrValue] of nextAttributes) {
    element.setAttribute(key, attrValue);
  }
  if (tagName === 'img' && !element.hasAttribute('src')) {
    element.remove();
  }
}

function sanitizeTiptapStyle(value: string): string {
  for (const part of value.split(';')) {
    const [rawKey, ...rawValueParts] = part.split(':');
    if (!rawKey || rawValueParts.length === 0) continue;
    if (rawKey.trim().toLowerCase() !== 'text-align') continue;
    const align = rawValueParts.join(':').trim().toLowerCase();
    if (['left', 'right', 'center', 'justify'].includes(align)) {
      return `text-align: ${align}`;
    }
  }
  return '';
}

function isSafeHTTPSURL(value: string): boolean {
  try {
    const parsed = new URL(value);
    return parsed.protocol === 'https:' && parsed.hostname.length > 0 && !parsed.username && !parsed.password;
  } catch {
    return false;
  }
}

function isSafePositiveInt(value: string, min: number, max: number): boolean {
  const parsed = Number.parseInt(value, 10);
  return Number.isInteger(parsed) && String(parsed) === value && parsed >= min && parsed <= max;
}

function sanitizeTiptapHTMLFallback(value: string): string {
  let sanitized = value.trim().length ? value : emptyHTML;
  sanitized = sanitized.replace(blockedElementPattern, '');
  sanitized = sanitized.replace(blockedVoidElementPattern, '');
  sanitized = sanitized.replace(eventAttributePattern, '');
  sanitized = sanitized.replace(dangerousURLAttributePattern, '');
  sanitized = sanitized.trim();
  return sanitized.length ? sanitized : emptyHTML;
}
