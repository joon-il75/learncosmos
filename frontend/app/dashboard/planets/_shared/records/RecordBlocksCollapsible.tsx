'use client'

import { useState, type CSSProperties } from 'react'
import type { RecordBlock } from './recordTypes'

const wrapStyle: CSSProperties = {
  display: 'grid',
  gap: 8,
}

const toggleStyle: CSSProperties = {
  width: 'fit-content',
  minHeight: 32,
  padding: '0 12px',
  borderRadius: 999,
  border: '1px solid rgba(138, 92, 38, 0.24)',
  background: 'rgba(255, 247, 231, 0.82)',
  color: '#5b3a13',
  fontSize: 12,
  fontWeight: 800,
  cursor: 'pointer',
}

const blockStyle: CSSProperties = {
  display: 'grid',
  gap: 6,
  padding: 12,
  borderRadius: 8,
  border: '1px solid rgba(120, 84, 38, 0.14)',
  background: 'rgba(255, 255, 255, 0.62)',
}

const blockMetaStyle: CSSProperties = {
  fontSize: 11,
  fontWeight: 800,
  color: '#7a5a24',
  textTransform: 'uppercase',
}

const blockTextStyle: CSSProperties = {
  margin: 0,
  whiteSpace: 'pre-wrap',
  wordBreak: 'break-word',
  fontSize: 13,
  lineHeight: 1.55,
  color: 'rgba(42, 31, 18, 0.86)',
}

function stringifyBlockContent(content: unknown): string {
  if (content == null) return '내용 없음'
  if (typeof content === 'string') return htmlToPlainText(content)
  if (typeof content === 'number' || typeof content === 'boolean') return String(content)
  if (typeof content === 'object') {
    const record = content as Record<string, unknown>
    const directText = record.text ?? record.body ?? record.title ?? record.url ?? record.src
    if (typeof directText === 'string' && directText.trim()) return htmlToPlainText(directText)
    return JSON.stringify(content, null, 2)
  }
  return '내용 없음'
}

function decodeHtmlEntity(entity: string): string {
  if (entity.startsWith('&#x') || entity.startsWith('&#X')) {
    const codePoint = Number.parseInt(entity.slice(3, -1), 16)
    return Number.isFinite(codePoint) ? String.fromCodePoint(codePoint) : entity
  }
  if (entity.startsWith('&#')) {
    const codePoint = Number.parseInt(entity.slice(2, -1), 10)
    return Number.isFinite(codePoint) ? String.fromCodePoint(codePoint) : entity
  }
  const namedEntities: Record<string, string> = {
    '&amp;': '&',
    '&lt;': '<',
    '&gt;': '>',
    '&quot;': '"',
    '&#39;': "'",
    '&nbsp;': ' ',
  }
  return namedEntities[entity] ?? entity
}

function htmlToPlainText(value: string): string {
  const stripped = value
    .replace(/<\s*(script|style)\b[^>]*>[\s\S]*?<\s*\/\s*\1\s*>/gi, '')
    .replace(/<\s*br\s*\/?>/gi, '\n')
    .replace(/<\/(p|div|li|h[1-6]|blockquote)>/gi, '\n')
    .replace(/<[^>]*>/g, '')
    .replace(/&(amp|lt|gt|quot|nbsp|#39|#[0-9]+|#x[0-9a-f]+);/gi, decodeHtmlEntity)
    .replace(/\n{3,}/g, '\n\n')
    .trim()

  return stripped || '내용 없음'
}

export function RecordBlocksCollapsible({ blocks }: { blocks: RecordBlock[] }) {
  const [open, setOpen] = useState(false)

  if (blocks.length === 0) return null

  return (
    <div style={wrapStyle}>
      <button type="button" style={toggleStyle} onClick={() => setOpen((prev) => !prev)}>
        연구 블록 원문 {open ? '접기' : '펼치기'} · {blocks.length}개
      </button>
      {open ? (
        <div style={wrapStyle}>
          {blocks.map((block) => (
            <div key={block.id} style={blockStyle}>
              <span style={blockMetaStyle}>{block.block_type}</span>
              <p style={blockTextStyle}>{stringifyBlockContent(block.content)}</p>
            </div>
          ))}
        </div>
      ) : null}
    </div>
  )
}
