'use client';

import { useRef, useState } from 'react';
import type { CSSProperties } from 'react';

interface VideoUploadPanelProps {
  accept: string;
  disabled: boolean;
  label: string;
  helpText: string;
  buttonStyle: CSSProperties;
  inputStyle: CSSProperties;
  mutedColor: string;
  onSelect: (file: File, title: string) => void;
}

export default function VideoUploadPanel({ accept, disabled, label, helpText, buttonStyle, inputStyle, mutedColor, onSelect }: VideoUploadPanelProps) {
  const inputRef = useRef<HTMLInputElement | null>(null);
  const [videoTitle, setVideoTitle] = useState('');

  return (
    <div style={{ width: '100%', maxWidth: '860px', margin: '0 auto', display: 'grid', gap: '8px' }}>
      <input
        ref={inputRef}
        type="file"
        accept={accept}
        style={{ display: 'none' }}
        onChange={(event) => {
          const file = event.target.files?.[0] ?? null;
          event.currentTarget.value = '';
          if (file) onSelect(file, videoTitle);
        }}
      />
      <div style={{ display: 'flex', flexWrap: 'wrap', alignItems: 'center', gap: '8px' }}>
        <input
          type="text"
          value={videoTitle}
          onChange={(event) => setVideoTitle(event.target.value)}
          disabled={disabled}
          placeholder="영상 제목"
          style={{ ...inputStyle, flex: '1 1 280px', minWidth: 0 }}
        />
        <button
          type="button"
          onClick={() => inputRef.current?.click()}
          disabled={disabled}
          style={{
            ...buttonStyle,
            flex: '0 0 auto',
            opacity: disabled ? 0.62 : 1,
            cursor: disabled ? 'default' : 'pointer',
          }}
        >
          {label}
        </button>
      </div>
      <span style={{ color: mutedColor, fontSize: '13px', lineHeight: 1.5 }}>{helpText}</span>
    </div>
  );
}
